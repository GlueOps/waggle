package api

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/glueops/waggle/internal/models/tenant"
	"github.com/glueops/waggle/internal/service"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
)

// ---- views ----

type poolView struct {
	ID           uuid.UUID       `json:"id"`
	DatacenterID uuid.UUID       `json:"datacenter_id"`
	SlotID       uuid.UUID       `json:"slot_id"`
	Name         string          `json:"name"`
	DesiredCount int             `json:"desired_count"`
	Metadata     json.RawMessage `json:"metadata,omitempty"`
	CreatedAt    time.Time       `json:"created_at" format:"date-time"`
	UpdatedAt    time.Time       `json:"updated_at" format:"date-time"`
}

func toPoolView(p *tenant.Pool) poolView {
	v := poolView{
		ID:           p.ID,
		DatacenterID: p.DatacenterID,
		SlotID:       p.SlotID,
		Name:         p.Name,
		DesiredCount: p.DesiredCount,
		CreatedAt:    p.CreatedAt,
		UpdatedAt:    p.UpdatedAt,
	}
	if len(p.Metadata) > 0 {
		v.Metadata = json.RawMessage(p.Metadata)
	}
	return v
}

type placementView struct {
	ID             uuid.UUID `json:"id"`
	HypervisorID   uuid.UUID `json:"hypervisor_id"`
	HypervisorName string    `json:"hypervisor_name"`
	VMID           *int      `json:"vmid,omitempty"`
	CreatedAt      time.Time `json:"created_at" format:"date-time"`
}

func toPlacementView(p service.PlacementView) placementView {
	return placementView{
		ID:             p.ID,
		HypervisorID:   p.HypervisorID,
		HypervisorName: p.HypervisorName,
		VMID:           p.VMID,
		CreatedAt:      p.CreatedAt,
	}
}

func toPlacementViews(ps []service.PlacementView) []placementView {
	out := make([]placementView, 0, len(ps))
	for _, p := range ps {
		out = append(out, toPlacementView(p))
	}
	return out
}

// ---- I/O ----

type createPoolInput struct {
	Body struct {
		DatacenterID uuid.UUID       `json:"datacenter_id" required:"true" format:"uuid"`
		SlotID       uuid.UUID       `json:"slot_id" required:"true" format:"uuid"`
		Name         string          `json:"name" required:"true" maxLength:"255"`
		DesiredCount int             `json:"desired_count" minimum:"0"`
		Metadata     json.RawMessage `json:"metadata,omitempty"`
	}
}

type resizePoolInput struct {
	ID   uuid.UUID `path:"id"`
	Body struct {
		DesiredCount int `json:"desired_count" minimum:"0"`
	}
}

type poolIDInput struct {
	ID uuid.UUID `path:"id"`
}

// poolViewOutput is the pool entity on its own, without placements. Used by
// create/get/resize so the read and write representations match: generated
// SDKs/Terraform then model the pool as a flat resource that round-trips.
// Placements live at GET /pools/{id}/placements (placementListOutput).
type poolViewOutput struct {
	Body poolView
}

type poolListOutput struct {
	Body struct {
		Items []poolView `json:"items"`
	}
}

type placementListOutput struct {
	Body struct {
		Items []placementView `json:"items"`
	}
}

type fleetPlacementView struct {
	ID             uuid.UUID `json:"id"`
	PoolID         uuid.UUID `json:"pool_id"`
	PoolName       string    `json:"pool_name"`
	HypervisorID   uuid.UUID `json:"hypervisor_id"`
	HypervisorName string    `json:"hypervisor_name"`
	SlotName       string    `json:"slot_name"`
	VCPU           int       `json:"vcpu"`
	RAMGB          int       `json:"ram_gb"`
	DiskGB         int       `json:"disk_gb"`
	VMID           *int      `json:"vmid,omitempty"`
	CreatedAt      time.Time `json:"created_at" format:"date-time"`
}

type fleetPlacementListOutput struct {
	Body struct {
		Items []fleetPlacementView `json:"items"`
	}
}

type backfillVMIDInput struct {
	ID   uuid.UUID `path:"id"`
	Body struct {
		VMID int `json:"vmid" required:"true" minimum:"1"`
	}
}

type placementOutput struct {
	Body placementView
}

func (s *Server) registerPool(fleet *service.FleetService, tokens *service.TokenService, keys *service.APIKeyService) {
	if fleet == nil || tokens == nil {
		return
	}
	// Dual auth: user JWT or organization API key (Terraform).
	requireAuth := RequireBearer(s.API, tokens, keys)
	// Separate tags for the two collections this file serves. Tag-driven
	// generators map one tag -> one resource, so pools and placements must not
	// share a tag (a conflated tag produces an unusable empty resource).
	tagged := func(tag string) func(huma.Operation) huma.Operation {
		return func(op huma.Operation) huma.Operation {
			op.Tags = []string{tag}
			op.Security = []map[string][]string{{"bearer": {}}}
			op.Middlewares = huma.Middlewares{requireAuth}
			return op
		}
	}
	// The pool tag matches its URL collection segment ("/pools"), so generators
	// bind it to a full CRUD resource.
	pools := tagged("pools")
	// Placements are NOT a standalone managed resource: they have no create or
	// read-by-id and are produced by pool create/resize. The tag is left as the
	// singular "placement" so it deliberately does NOT match "/placements" —
	// otherwise the generator emits a broken partial-CRUD resource (an Update
	// against an empty model: `var result client.` with no type).
	placements := tagged("placement")

	huma.Register(s.API, pools(huma.Operation{
		OperationID:   "create-pool",
		Method:        http.MethodPost,
		Path:          "/pools",
		Summary:       "Create a node pool and place its VMs across hypervisors (anti-affinity spread, all-or-nothing). Placements are available at GET /pools/{id}/placements.",
		DefaultStatus: http.StatusCreated,
	}), func(ctx context.Context, in *createPoolInput) (*poolViewOutput, error) {
		orgID, err := orgFromContext(ctx)
		if err != nil {
			return nil, err
		}
		res, err := fleet.CreatePool(ctx, orgID, service.PoolInput{
			DatacenterID: in.Body.DatacenterID,
			SlotID:       in.Body.SlotID,
			Name:         in.Body.Name,
			DesiredCount: in.Body.DesiredCount,
			Metadata:     in.Body.Metadata,
		})
		if err != nil {
			return nil, mapPoolError(err)
		}
		return &poolViewOutput{Body: toPoolView(&res.Pool)}, nil
	})

	huma.Register(s.API, pools(huma.Operation{
		OperationID: "list-pools",
		Method:      http.MethodGet,
		Path:        "/pools",
		Summary:     "List node pools in the caller's tenant.",
	}), func(ctx context.Context, _ *struct{}) (*poolListOutput, error) {
		orgID, err := orgFromContext(ctx)
		if err != nil {
			return nil, err
		}
		pools, err := fleet.ListPools(ctx, orgID)
		if err != nil {
			return nil, mapPoolError(err)
		}
		out := &poolListOutput{}
		out.Body.Items = make([]poolView, 0, len(pools))
		for i := range pools {
			out.Body.Items = append(out.Body.Items, toPoolView(&pools[i]))
		}
		return out, nil
	})

	huma.Register(s.API, pools(huma.Operation{
		OperationID: "get-pool",
		Method:      http.MethodGet,
		Path:        "/pools/{id}",
		Summary:     "Fetch a pool. Its placements are available at GET /pools/{id}/placements.",
	}), func(ctx context.Context, in *poolIDInput) (*poolViewOutput, error) {
		orgID, err := orgFromContext(ctx)
		if err != nil {
			return nil, err
		}
		res, err := fleet.GetPool(ctx, orgID, in.ID)
		if err != nil {
			return nil, mapPoolError(err)
		}
		return &poolViewOutput{Body: toPoolView(&res.Pool)}, nil
	})

	huma.Register(s.API, pools(huma.Operation{
		OperationID: "resize-pool",
		Method:      http.MethodPatch,
		Path:        "/pools/{id}",
		Summary:     "Resize a pool's desired count. Grow places new VMs (all-or-nothing); shrink removes newest placements (LIFO). Placements are available at GET /pools/{id}/placements.",
	}), func(ctx context.Context, in *resizePoolInput) (*poolViewOutput, error) {
		orgID, err := orgFromContext(ctx)
		if err != nil {
			return nil, err
		}
		res, err := fleet.ResizePool(ctx, orgID, in.ID, in.Body.DesiredCount)
		if err != nil {
			return nil, mapPoolError(err)
		}
		return &poolViewOutput{Body: toPoolView(&res.Pool)}, nil
	})

	huma.Register(s.API, pools(huma.Operation{
		OperationID:   "delete-pool",
		Method:        http.MethodDelete,
		Path:          "/pools/{id}",
		Summary:       "Delete a pool and release all its placements.",
		DefaultStatus: http.StatusNoContent,
	}), func(ctx context.Context, in *poolIDInput) (*struct{}, error) {
		orgID, err := orgFromContext(ctx)
		if err != nil {
			return nil, err
		}
		if err := fleet.DeletePool(ctx, orgID, in.ID); err != nil {
			return nil, mapPoolError(err)
		}
		return nil, nil
	})

	huma.Register(s.API, pools(huma.Operation{
		OperationID: "list-pool-placements",
		Method:      http.MethodGet,
		Path:        "/pools/{id}/placements",
		Summary:     "List a pool's placements (hypervisor + optional vmid).",
	}), func(ctx context.Context, in *poolIDInput) (*placementListOutput, error) {
		orgID, err := orgFromContext(ctx)
		if err != nil {
			return nil, err
		}
		placements, err := fleet.ListPlacements(ctx, orgID, in.ID)
		if err != nil {
			return nil, mapPoolError(err)
		}
		out := &placementListOutput{}
		out.Body.Items = toPlacementViews(placements)
		return out, nil
	})

	huma.Register(s.API, placements(huma.Operation{
		OperationID: "list-placements",
		Method:      http.MethodGet,
		Path:        "/placements",
		Summary:     "List all placements in the tenant with pool, slot, and hypervisor context (fleet overview).",
	}), func(ctx context.Context, _ *struct{}) (*fleetPlacementListOutput, error) {
		orgID, err := orgFromContext(ctx)
		if err != nil {
			return nil, err
		}
		placements, err := fleet.ListAllPlacements(ctx, orgID)
		if err != nil {
			return nil, mapPoolError(err)
		}
		out := &fleetPlacementListOutput{}
		out.Body.Items = make([]fleetPlacementView, 0, len(placements))
		for _, p := range placements {
			out.Body.Items = append(out.Body.Items, fleetPlacementView{
				ID:             p.ID,
				PoolID:         p.PoolID,
				PoolName:       p.PoolName,
				HypervisorID:   p.HypervisorID,
				HypervisorName: p.HypervisorName,
				SlotName:       p.SlotName,
				VCPU:           p.VCPU,
				RAMGB:          p.RAMGB,
				DiskGB:         p.DiskGB,
				VMID:           p.VMID,
				CreatedAt:      p.CreatedAt,
			})
		}
		return out, nil
	})

	huma.Register(s.API, placements(huma.Operation{
		OperationID: "backfill-placement-vmid",
		Method:      http.MethodPatch,
		Path:        "/placements/{id}",
		Summary:     "Attach the externally-assigned Proxmox vmid to a placement.",
	}), func(ctx context.Context, in *backfillVMIDInput) (*placementOutput, error) {
		orgID, err := orgFromContext(ctx)
		if err != nil {
			return nil, err
		}
		pl, err := fleet.BackfillVMID(ctx, orgID, in.ID, in.Body.VMID)
		if err != nil {
			return nil, mapPoolError(err)
		}
		return &placementOutput{Body: placementView{
			ID:           pl.ID,
			HypervisorID: pl.HypervisorID,
			VMID:         pl.VMID,
			CreatedAt:    pl.CreatedAt,
		}}, nil
	})
}

func mapPoolError(err error) error {
	var ce *service.CapacityError
	if errors.As(err, &ce) {
		return huma.Error422UnprocessableEntity(ce.Error())
	}
	return mapFleetError(err)
}
