package api

import (
	"context"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/glueops/waggle/internal/database"
	"github.com/glueops/waggle/internal/models/tenant"
	"github.com/glueops/waggle/internal/service"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
)

// ---- views ----

type datacenterView struct {
	ID                 uuid.UUID `json:"id"`
	Name               string    `json:"name"`
	URL                string    `json:"url"`
	HasToken           bool      `json:"has_token" doc:"Whether a Proxmox API token is configured (the token itself is never returned)."`
	InsecureSkipVerify bool      `json:"insecure_skip_verify" doc:"Whether TLS verification is disabled for this cluster (self-signed certs)."`
	CreatedAt          time.Time `json:"created_at" format:"date-time"`
	UpdatedAt          time.Time `json:"updated_at" format:"date-time"`
}

func toDatacenterView(d *tenant.Datacenter) datacenterView {
	return datacenterView{
		ID:                 d.ID,
		Name:               d.Name,
		URL:                d.Url,
		HasToken:           d.EncryptedTokenKey != "",
		InsecureSkipVerify: d.InsecureSkipVerify,
		CreatedAt:          d.CreatedAt,
		UpdatedAt:          d.UpdatedAt,
	}
}

type slotView struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	VCPU      int       `json:"vcpu"`
	RAMGB     int       `json:"ram_gb"`
	DiskGB    int       `json:"disk_gb"`
	CreatedAt time.Time `json:"created_at" format:"date-time"`
	UpdatedAt time.Time `json:"updated_at" format:"date-time"`
}

func toSlotView(s *tenant.Slot) slotView {
	return slotView{ID: s.ID, Name: s.Name, VCPU: s.VCPU, RAMGB: s.RAMGB, DiskGB: s.DiskGB, CreatedAt: s.CreatedAt, UpdatedAt: s.UpdatedAt}
}

type hypervisorView struct {
	ID             uuid.UUID  `json:"id"`
	DatacenterID   uuid.UUID  `json:"datacenter_id"`
	Name           string     `json:"name"`
	CPUTotal       int        `json:"cpu_total"`
	CPUReserved    int        `json:"cpu_reserved"`
	CPUUsed        int        `json:"cpu_used" doc:"vCPU allocated to existing guests (from discovery)."`
	CPUBookable    int        `json:"cpu_bookable"`
	RAMGBTotal     int        `json:"ram_gb_total"`
	RAMGBReserved  int        `json:"ram_gb_reserved"`
	RAMGBUsed      int        `json:"ram_gb_used" doc:"RAM (GB) allocated to existing guests (from discovery)."`
	RAMGBBookable  int        `json:"ram_gb_bookable"`
	DiskGBTotal    int        `json:"disk_gb_total"`
	DiskGBReserved int        `json:"disk_gb_reserved"`
	DiskGBUsed     int        `json:"disk_gb_used" doc:"Disk (GB) allocated to existing guests (from discovery)."`
	DiskGBBookable int        `json:"disk_gb_bookable"`
	Schedulable    bool       `json:"schedulable" doc:"When false, placement excludes this hypervisor."`
	LastSyncedAt   *time.Time `json:"last_synced_at,omitempty" format:"date-time"`
	CreatedAt      time.Time  `json:"created_at" format:"date-time"`
	UpdatedAt      time.Time  `json:"updated_at" format:"date-time"`
}

func toHypervisorView(h *tenant.Hypervisor) hypervisorView {
	return hypervisorView{
		ID:             h.ID,
		DatacenterID:   h.DatacenterID,
		Name:           h.Name,
		CPUTotal:       h.CPUTotal,
		CPUReserved:    h.CPUReserved,
		CPUUsed:        h.CPUUsed,
		CPUBookable:    h.CPUTotal - h.CPUReserved - h.CPUUsed,
		RAMGBTotal:     h.RAMGBTotal,
		RAMGBReserved:  h.RAMGBReserved,
		RAMGBUsed:      h.RAMGBUsed,
		RAMGBBookable:  h.RAMGBTotal - h.RAMGBReserved - h.RAMGBUsed,
		DiskGBTotal:    h.DiskGBTotal,
		DiskGBReserved: h.DiskGBReserved,
		DiskGBUsed:     h.DiskGBUsed,
		DiskGBBookable: h.DiskGBTotal - h.DiskGBReserved - h.DiskGBUsed,
		Schedulable:    h.Schedulable,
		LastSyncedAt:   h.LastSyncedAt,
		CreatedAt:      h.CreatedAt,
		UpdatedAt:      h.UpdatedAt,
	}
}

// ---- datacenter I/O ----

type datacenterBody struct {
	Name string `json:"name" required:"true" maxLength:"255"`
	URL  string `json:"url" required:"true" maxLength:"255"`
	// Token is the Proxmox API token ("USER@REALM!TOKENID=SECRET") used for
	// hypervisor discovery. Write-only: encrypted with the tenant DEK and never
	// returned. Omit to leave an existing token unchanged.
	Token string `json:"token,omitempty" maxLength:"512"`
	// InsecureSkipVerify disables TLS verification for self-signed Proxmox
	// clusters. Omit to default false (create) / leave unchanged (update).
	InsecureSkipVerify *bool `json:"insecure_skip_verify,omitempty"`
}

// tokenPtr returns nil when no token was supplied so the service leaves any
// existing token untouched, or a pointer to the value to (re)set it.
func (b datacenterBody) tokenPtr() *string {
	if b.Token == "" {
		return nil
	}
	t := b.Token
	return &t
}

type createDatacenterInput struct {
	Body datacenterBody
}
type updateDatacenterInput struct {
	ID   uuid.UUID `path:"id"`
	Body datacenterBody
}
type datacenterIDInput struct {
	ID uuid.UUID `path:"id"`
}
type datacenterOutput struct {
	Body datacenterView
}
type datacenterListOutput struct {
	Body struct {
		Items []datacenterView `json:"items"`
	}
}

// ---- slot I/O ----

type slotBody struct {
	Name   string `json:"name" required:"true" maxLength:"255"`
	VCPU   int    `json:"vcpu" required:"true" minimum:"1"`
	RAMGB  int    `json:"ram_gb" required:"true" minimum:"1"`
	DiskGB int    `json:"disk_gb" required:"true" minimum:"1"`
}

type createSlotInput struct {
	Body slotBody
}
type updateSlotInput struct {
	ID   uuid.UUID `path:"id"`
	Body slotBody
}
type slotIDInput struct {
	ID uuid.UUID `path:"id"`
}
type listSlotsInput struct {
	Name string `query:"name" doc:"Filter to the slot with this exact (unique) name."`
}
type slotOutput struct {
	Body slotView
}
type slotListOutput struct {
	Body struct {
		Items []slotView `json:"items"`
	}
}

// ---- hypervisor I/O ----

type hypervisorBody struct {
	DatacenterID   uuid.UUID `json:"datacenter_id" required:"true" format:"uuid"`
	Name           string    `json:"name" required:"true" maxLength:"255"`
	CPUTotal       int       `json:"cpu_total" minimum:"0"`
	CPUReserved    int       `json:"cpu_reserved" minimum:"0"`
	RAMGBTotal     int       `json:"ram_gb_total" minimum:"0"`
	RAMGBReserved  int       `json:"ram_gb_reserved" minimum:"0"`
	DiskGBTotal    int       `json:"disk_gb_total" minimum:"0"`
	DiskGBReserved int       `json:"disk_gb_reserved" minimum:"0"`
	// Schedulable: omit to default true (create) / leave unchanged (update).
	Schedulable *bool `json:"schedulable,omitempty"`
}

type createHypervisorInput struct {
	Body hypervisorBody
}
type updateHypervisorInput struct {
	ID   uuid.UUID `path:"id"`
	Body hypervisorBody
}
type hypervisorIDInput struct {
	ID uuid.UUID `path:"id"`
}
type listHypervisorsInput struct {
	// huma rejects pointer params, so this is a plain string and parsed in
	// the handler; empty means "no datacenter filter".
	DatacenterID string `query:"datacenter_id"`
}
type hypervisorOutput struct {
	Body hypervisorView
}
type hypervisorListOutput struct {
	Body struct {
		Items []hypervisorView `json:"items"`
	}
}

// ---- discovery I/O ----

type discoverInput struct {
	ID   uuid.UUID `path:"id"`
	Body struct {
		// Async queues a background discovery job and returns immediately
		// instead of calling Proxmox inline.
		Async bool `json:"async,omitempty"`
	}
}

type discoverOutput struct {
	Body struct {
		// Queued is true when the request was handled asynchronously; Items is
		// then empty and the result lands later (see hypervisor last_synced_at).
		Queued bool             `json:"queued"`
		Items  []hypervisorView `json:"items"`
	}
}

// discoveryEnqueuer is the subset of the jobs enqueuer the tenant API needs for
// async discovery. Declared as an interface so api stays decoupled from jobs.
type discoveryEnqueuer interface {
	EnqueueDiscoverHypervisors(ctx context.Context, orgID, datacenterID uuid.UUID) error
}

func (s *Server) registerTenant(fleet *service.FleetService, tokens *service.TokenService, keys *service.APIKeyService, enq discoveryEnqueuer) {
	if fleet == nil || tokens == nil {
		return
	}
	// Dual auth: user JWT or organization API key (Terraform). Both resolve to
	// the tenant's org context.
	requireAuth := RequireBearer(s.API, tokens, keys)
	// One OpenAPI tag per entity (not a single "fleet" tag): tag-driven SDK and
	// Terraform-provider generators map one tag -> one resource, so a tag that
	// spans datacenters, slots, and hypervisors yields an unusable empty
	// resource. Tag each operation by the collection it manages instead.
	tagged := func(tag string) func(huma.Operation) huma.Operation {
		return func(op huma.Operation) huma.Operation {
			op.Tags = []string{tag}
			op.Security = []map[string][]string{{"bearer": {}}}
			op.Middlewares = huma.Middlewares{requireAuth}
			return op
		}
	}
	// Tag names must match the URL collection segment (e.g. "/datacenters"):
	// the terraform-provider/SDK generators bind a tag's resource to the path
	// whose last segment equals the tag, so a singular tag never binds.
	datacenters := tagged("datacenters")
	slots := tagged("slots")
	hypervisors := tagged("hypervisors")

	huma.Register(s.API, datacenters(huma.Operation{
		OperationID:   "create-datacenter",
		Method:        http.MethodPost,
		Path:          "/datacenters",
		Summary:       "Create a datacenter in the caller's tenant.",
		DefaultStatus: http.StatusCreated,
	}), func(ctx context.Context, in *createDatacenterInput) (*datacenterOutput, error) {
		orgID, err := orgFromContext(ctx)
		if err != nil {
			return nil, err
		}
		dc, err := fleet.CreateDatacenter(ctx, orgID, service.DatacenterInput{Name: in.Body.Name, URL: in.Body.URL, Token: in.Body.tokenPtr(), InsecureSkipVerify: in.Body.InsecureSkipVerify})
		if err != nil {
			return nil, mapFleetError(err)
		}
		return &datacenterOutput{Body: toDatacenterView(dc)}, nil
	})

	huma.Register(s.API, datacenters(huma.Operation{
		OperationID: "list-datacenters",
		Method:      http.MethodGet,
		Path:        "/datacenters",
		Summary:     "List datacenters in the caller's tenant.",
	}), func(ctx context.Context, _ *struct{}) (*datacenterListOutput, error) {
		orgID, err := orgFromContext(ctx)
		if err != nil {
			return nil, err
		}
		dcs, err := fleet.ListDatacenters(ctx, orgID)
		if err != nil {
			return nil, mapFleetError(err)
		}
		out := &datacenterListOutput{}
		out.Body.Items = make([]datacenterView, 0, len(dcs))
		for i := range dcs {
			out.Body.Items = append(out.Body.Items, toDatacenterView(&dcs[i]))
		}
		return out, nil
	})

	huma.Register(s.API, datacenters(huma.Operation{
		OperationID: "get-datacenter",
		Method:      http.MethodGet,
		Path:        "/datacenters/{id}",
		Summary:     "Fetch a datacenter by ID.",
	}), func(ctx context.Context, in *datacenterIDInput) (*datacenterOutput, error) {
		orgID, err := orgFromContext(ctx)
		if err != nil {
			return nil, err
		}
		dc, err := fleet.GetDatacenter(ctx, orgID, in.ID)
		if err != nil {
			return nil, mapFleetError(err)
		}
		return &datacenterOutput{Body: toDatacenterView(dc)}, nil
	})

	huma.Register(s.API, datacenters(huma.Operation{
		OperationID: "update-datacenter",
		Method:      http.MethodPut,
		Path:        "/datacenters/{id}",
		Summary:     "Update a datacenter.",
	}), func(ctx context.Context, in *updateDatacenterInput) (*datacenterOutput, error) {
		orgID, err := orgFromContext(ctx)
		if err != nil {
			return nil, err
		}
		dc, err := fleet.UpdateDatacenter(ctx, orgID, in.ID, service.DatacenterInput{Name: in.Body.Name, URL: in.Body.URL, Token: in.Body.tokenPtr(), InsecureSkipVerify: in.Body.InsecureSkipVerify})
		if err != nil {
			return nil, mapFleetError(err)
		}
		return &datacenterOutput{Body: toDatacenterView(dc)}, nil
	})

	huma.Register(s.API, datacenters(huma.Operation{
		OperationID:   "delete-datacenter",
		Method:        http.MethodDelete,
		Path:          "/datacenters/{id}",
		Summary:       "Delete a datacenter.",
		DefaultStatus: http.StatusNoContent,
	}), func(ctx context.Context, in *datacenterIDInput) (*struct{}, error) {
		orgID, err := orgFromContext(ctx)
		if err != nil {
			return nil, err
		}
		if err := fleet.DeleteDatacenter(ctx, orgID, in.ID); err != nil {
			return nil, mapFleetError(err)
		}
		return nil, nil
	})

	huma.Register(s.API, slots(huma.Operation{
		OperationID:   "create-slot",
		Method:        http.MethodPost,
		Path:          "/slots",
		Summary:       "Create a slot (t-shirt-size VM template) in the caller's tenant.",
		DefaultStatus: http.StatusCreated,
	}), func(ctx context.Context, in *createSlotInput) (*slotOutput, error) {
		orgID, err := orgFromContext(ctx)
		if err != nil {
			return nil, err
		}
		slot, err := fleet.CreateSlot(ctx, orgID, service.SlotInput{Name: in.Body.Name, VCPU: in.Body.VCPU, RAMGB: in.Body.RAMGB, DiskGB: in.Body.DiskGB})
		if err != nil {
			return nil, mapFleetError(err)
		}
		return &slotOutput{Body: toSlotView(slot)}, nil
	})

	huma.Register(s.API, slots(huma.Operation{
		OperationID: "list-slots",
		Method:      http.MethodGet,
		Path:        "/slots",
		Summary:     "List slots in the caller's tenant.",
	}), func(ctx context.Context, in *listSlotsInput) (*slotListOutput, error) {
		orgID, err := orgFromContext(ctx)
		if err != nil {
			return nil, err
		}
		slots, err := fleet.ListSlots(ctx, orgID, in.Name)
		if err != nil {
			return nil, mapFleetError(err)
		}
		out := &slotListOutput{}
		out.Body.Items = make([]slotView, 0, len(slots))
		for i := range slots {
			out.Body.Items = append(out.Body.Items, toSlotView(&slots[i]))
		}
		return out, nil
	})

	huma.Register(s.API, slots(huma.Operation{
		OperationID: "get-slot",
		Method:      http.MethodGet,
		Path:        "/slots/{id}",
		Summary:     "Fetch a slot by ID.",
	}), func(ctx context.Context, in *slotIDInput) (*slotOutput, error) {
		orgID, err := orgFromContext(ctx)
		if err != nil {
			return nil, err
		}
		slot, err := fleet.GetSlot(ctx, orgID, in.ID)
		if err != nil {
			return nil, mapFleetError(err)
		}
		return &slotOutput{Body: toSlotView(slot)}, nil
	})

	huma.Register(s.API, slots(huma.Operation{
		OperationID: "update-slot",
		Method:      http.MethodPut,
		Path:        "/slots/{id}",
		Summary:     "Update a slot.",
	}), func(ctx context.Context, in *updateSlotInput) (*slotOutput, error) {
		orgID, err := orgFromContext(ctx)
		if err != nil {
			return nil, err
		}
		slot, err := fleet.UpdateSlot(ctx, orgID, in.ID, service.SlotInput{Name: in.Body.Name, VCPU: in.Body.VCPU, RAMGB: in.Body.RAMGB, DiskGB: in.Body.DiskGB})
		if err != nil {
			return nil, mapFleetError(err)
		}
		return &slotOutput{Body: toSlotView(slot)}, nil
	})

	huma.Register(s.API, slots(huma.Operation{
		OperationID:   "delete-slot",
		Method:        http.MethodDelete,
		Path:          "/slots/{id}",
		Summary:       "Delete a slot.",
		DefaultStatus: http.StatusNoContent,
	}), func(ctx context.Context, in *slotIDInput) (*struct{}, error) {
		orgID, err := orgFromContext(ctx)
		if err != nil {
			return nil, err
		}
		if err := fleet.DeleteSlot(ctx, orgID, in.ID); err != nil {
			return nil, mapFleetError(err)
		}
		return nil, nil
	})

	huma.Register(s.API, hypervisors(huma.Operation{
		OperationID:   "create-hypervisor",
		Method:        http.MethodPost,
		Path:          "/hypervisors",
		Summary:       "Create a hypervisor in the caller's tenant.",
		DefaultStatus: http.StatusCreated,
	}), func(ctx context.Context, in *createHypervisorInput) (*hypervisorOutput, error) {
		orgID, err := orgFromContext(ctx)
		if err != nil {
			return nil, err
		}
		hv, err := fleet.CreateHypervisor(ctx, orgID, service.HypervisorInput{
			DatacenterID:   in.Body.DatacenterID,
			Name:           in.Body.Name,
			CPUTotal:       in.Body.CPUTotal,
			CPUReserved:    in.Body.CPUReserved,
			RAMGBTotal:     in.Body.RAMGBTotal,
			RAMGBReserved:  in.Body.RAMGBReserved,
			DiskGBTotal:    in.Body.DiskGBTotal,
			DiskGBReserved: in.Body.DiskGBReserved,
			Schedulable:    in.Body.Schedulable,
		})
		if err != nil {
			return nil, mapFleetError(err)
		}
		return &hypervisorOutput{Body: toHypervisorView(hv)}, nil
	})

	huma.Register(s.API, hypervisors(huma.Operation{
		OperationID: "list-hypervisors",
		Method:      http.MethodGet,
		Path:        "/hypervisors",
		Summary:     "List hypervisors in the caller's tenant.",
	}), func(ctx context.Context, in *listHypervisorsInput) (*hypervisorListOutput, error) {
		orgID, err := orgFromContext(ctx)
		if err != nil {
			return nil, err
		}
		var dcFilter *uuid.UUID
		if in.DatacenterID != "" {
			id, perr := uuid.Parse(in.DatacenterID)
			if perr != nil {
				return nil, huma.Error422UnprocessableEntity("invalid datacenter_id")
			}
			dcFilter = &id
		}
		hvs, err := fleet.ListHypervisors(ctx, orgID, dcFilter)
		if err != nil {
			return nil, mapFleetError(err)
		}
		out := &hypervisorListOutput{}
		out.Body.Items = make([]hypervisorView, 0, len(hvs))
		for i := range hvs {
			out.Body.Items = append(out.Body.Items, toHypervisorView(&hvs[i]))
		}
		return out, nil
	})

	huma.Register(s.API, hypervisors(huma.Operation{
		OperationID: "get-hypervisor",
		Method:      http.MethodGet,
		Path:        "/hypervisors/{id}",
		Summary:     "Fetch a hypervisor by ID.",
	}), func(ctx context.Context, in *hypervisorIDInput) (*hypervisorOutput, error) {
		orgID, err := orgFromContext(ctx)
		if err != nil {
			return nil, err
		}
		hv, err := fleet.GetHypervisor(ctx, orgID, in.ID)
		if err != nil {
			return nil, mapFleetError(err)
		}
		return &hypervisorOutput{Body: toHypervisorView(hv)}, nil
	})

	huma.Register(s.API, hypervisors(huma.Operation{
		OperationID: "update-hypervisor",
		Method:      http.MethodPut,
		Path:        "/hypervisors/{id}",
		Summary:     "Update a hypervisor.",
	}), func(ctx context.Context, in *updateHypervisorInput) (*hypervisorOutput, error) {
		orgID, err := orgFromContext(ctx)
		if err != nil {
			return nil, err
		}
		hv, err := fleet.UpdateHypervisor(ctx, orgID, in.ID, service.HypervisorInput{
			DatacenterID:   in.Body.DatacenterID,
			Name:           in.Body.Name,
			CPUTotal:       in.Body.CPUTotal,
			CPUReserved:    in.Body.CPUReserved,
			RAMGBTotal:     in.Body.RAMGBTotal,
			RAMGBReserved:  in.Body.RAMGBReserved,
			DiskGBTotal:    in.Body.DiskGBTotal,
			DiskGBReserved: in.Body.DiskGBReserved,
			Schedulable:    in.Body.Schedulable,
		})
		if err != nil {
			return nil, mapFleetError(err)
		}
		return &hypervisorOutput{Body: toHypervisorView(hv)}, nil
	})

	huma.Register(s.API, hypervisors(huma.Operation{
		OperationID:   "delete-hypervisor",
		Method:        http.MethodDelete,
		Path:          "/hypervisors/{id}",
		Summary:       "Delete a hypervisor.",
		DefaultStatus: http.StatusNoContent,
	}), func(ctx context.Context, in *hypervisorIDInput) (*struct{}, error) {
		orgID, err := orgFromContext(ctx)
		if err != nil {
			return nil, err
		}
		if err := fleet.DeleteHypervisor(ctx, orgID, in.ID); err != nil {
			return nil, mapFleetError(err)
		}
		return nil, nil
	})

	huma.Register(s.API, datacenters(huma.Operation{
		OperationID: "discover-hypervisors",
		Method:      http.MethodPost,
		Path:        "/datacenters/{id}/discover",
		Summary:     "Discover hypervisors from the datacenter's Proxmox cluster and upsert them (preserving reserved capacity and schedulable). Set async to run in the background.",
	}), func(ctx context.Context, in *discoverInput) (*discoverOutput, error) {
		orgID, err := orgFromContext(ctx)
		if err != nil {
			return nil, err
		}
		out := &discoverOutput{}
		if in.Body.Async {
			if enq == nil {
				return nil, huma.Error503ServiceUnavailable("async discovery unavailable: job queue not configured")
			}
			if err := enq.EnqueueDiscoverHypervisors(ctx, orgID, in.ID); err != nil {
				return nil, huma.Error500InternalServerError("failed to enqueue discovery")
			}
			out.Body.Queued = true
			out.Body.Items = []hypervisorView{}
			return out, nil
		}
		hvs, err := fleet.DiscoverHypervisors(ctx, orgID, in.ID)
		if err != nil {
			return nil, mapFleetError(err)
		}
		out.Body.Items = make([]hypervisorView, 0, len(hvs))
		for i := range hvs {
			out.Body.Items = append(out.Body.Items, toHypervisorView(&hvs[i]))
		}
		return out, nil
	})
}

func orgFromContext(ctx context.Context) (uuid.UUID, error) {
	claims, ok := ClaimsFromContext(ctx)
	if !ok {
		return uuid.Nil, huma.Error401Unauthorized("missing claims")
	}
	if claims.OrganizationID == uuid.Nil {
		return uuid.Nil, huma.Error400BadRequest("token has no organization context")
	}
	return claims.OrganizationID, nil
}

func mapFleetError(err error) error {
	switch {
	case errors.Is(err, service.ErrNotFound):
		return huma.Error404NotFound("not found")
	case errors.Is(err, service.ErrInvalidInput):
		return huma.Error422UnprocessableEntity("invalid input")
	case errors.Is(err, database.ErrTenantNotProvisioned), errors.Is(err, database.ErrTenantNotActive):
		return huma.Error503ServiceUnavailable("tenant database not ready; provisioning may still be in progress")
	case errors.Is(err, database.ErrOrgNotFound):
		return huma.Error401Unauthorized("organization no longer exists; please re-authenticate")
	case errors.Is(err, service.ErrDiscovery):
		return huma.Error502BadGateway("could not reach the Proxmox cluster for discovery")
	}
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" {
		return huma.Error409Conflict("a resource with that unique value already exists")
	}
	// Unmapped: surface a generic 500 to the client but log the real cause so
	// these aren't opaque in the API logs.
	log.Printf("fleet: unhandled error -> 500: %v", err)
	return huma.Error500InternalServerError("internal error")
}
