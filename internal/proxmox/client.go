// Package proxmox is a minimal read-only client for the Proxmox VE API. Waggle
// uses it only to discover hypervisor capacity (nodes); it never creates or
// mutates anything on the cluster, in keeping with Waggle's role as a placement
// oracle rather than a provisioner.
package proxmox

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const defaultTimeout = 15 * time.Second

// Config configures a Client. Token is the full Proxmox API token credential in
// the form "USER@REALM!TOKENID=SECRET"; it is sent verbatim in the
// "Authorization: PVEAPIToken=..." header.
type Config struct {
	BaseURL string
	Token   string
	// InsecureSkipVerify disables TLS verification. Proxmox clusters often use
	// self-signed certs; callers opt in explicitly rather than defaulting to
	// insecure.
	InsecureSkipVerify bool
	HTTPClient         *http.Client
}

type Client struct {
	baseURL string
	token   string
	hc      *http.Client
}

func New(cfg Config) (*Client, error) {
	base, err := normalizeBaseURL(cfg.BaseURL)
	if err != nil {
		return nil, err
	}
	if strings.TrimSpace(cfg.Token) == "" {
		return nil, errors.New("proxmox: token is required")
	}
	hc := cfg.HTTPClient
	if hc == nil {
		hc = &http.Client{
			Timeout: defaultTimeout,
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: cfg.InsecureSkipVerify},
			},
		}
	}
	return &Client{baseURL: base, token: cfg.Token, hc: hc}, nil
}

// normalizeBaseURL reduces a stored datacenter URL to a clean API origin
// (scheme://host[:port]), tolerating a trailing slash, path, or query string
// (e.g. a legacy "https://host:8006/?insecure=1") so we never build malformed
// request URLs like ".../?insecure=1/api2/json/nodes".
func normalizeBaseURL(raw string) (string, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return "", errors.New("proxmox: base url is required")
	}
	u, err := url.Parse(raw)
	if err != nil {
		return "", fmt.Errorf("proxmox: invalid base url %q: %w", raw, err)
	}
	if u.Scheme == "" || u.Host == "" {
		return "", fmt.Errorf("proxmox: base url %q must be absolute (scheme://host)", raw)
	}
	return u.Scheme + "://" + u.Host, nil
}

// getJSON performs an authenticated GET and decodes the JSON body into out.
func (c *Client) getJSON(ctx context.Context, path string, out any) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+path, nil)
	if err != nil {
		return fmt.Errorf("proxmox: build request: %w", err)
	}
	req.Header.Set("Authorization", "PVEAPIToken="+c.token)
	req.Header.Set("Accept", "application/json")

	resp, err := c.hc.Do(req)
	if err != nil {
		return fmt.Errorf("proxmox: GET %s: %w", path, err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(io.LimitReader(resp.Body, 4<<20))
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("proxmox: GET %s: unexpected status %d: %s", path, resp.StatusCode, strings.TrimSpace(string(body)))
	}
	if err := json.Unmarshal(body, out); err != nil {
		return fmt.Errorf("proxmox: decode %s: %w", path, err)
	}
	return nil
}

// Node is a discovered Proxmox cluster node with its total capacity.
type Node struct {
	Name    string
	Status  string
	MaxCPU  int   // logical CPU cores
	MaxMem  int64 // bytes
	MaxDisk int64 // bytes (images-capable storage / root filesystem)
}

// ListNodes returns the cluster's nodes via GET /api2/json/nodes.
func (c *Client) ListNodes(ctx context.Context) ([]Node, error) {
	var payload struct {
		Data []struct {
			Node    string `json:"node"`
			Status  string `json:"status"`
			MaxCPU  int    `json:"maxcpu"`
			MaxMem  int64  `json:"maxmem"`
			MaxDisk int64  `json:"maxdisk"`
		} `json:"data"`
	}
	if err := c.getJSON(ctx, "/api2/json/nodes", &payload); err != nil {
		return nil, err
	}
	nodes := make([]Node, 0, len(payload.Data))
	for _, n := range payload.Data {
		nodes = append(nodes, Node{
			Name:    n.Node,
			Status:  n.Status,
			MaxCPU:  n.MaxCPU,
			MaxMem:  n.MaxMem,
			MaxDisk: n.MaxDisk,
		})
	}
	return nodes, nil
}

// NodeUsage is the capacity already committed to existing guests on a node:
// the sum of every VM's and container's configured (allocated) resources,
// which is the meaningful "used" figure for placement — distinct from live
// runtime utilization.
type NodeUsage struct {
	VCPU      int   // sum of guest vCPUs
	MemBytes  int64 // sum of guest configured RAM
	DiskBytes int64 // sum of guest configured disk
	Guests    int
	// Managed is the set of excluded (Waggle-managed) vmids actually found
	// resident on this node. It is the only evidence Waggle ever gets that a
	// guest is not on the hypervisor it was assigned: backfill reports the vmid
	// but never the node, and Waggle never touches the cluster itself. Without
	// it a guest built in the wrong place stays charged to the node its
	// placement names and is invisible on the node hosting it, because the
	// exclude set skipped it here.
	Managed []int
}

// guestRow is the shared shape of /qemu and /lxc list entries we sum over.
type guestRow struct {
	VMID    int   `json:"vmid"`
	CPUs    int   `json:"cpus"`
	MaxMem  int64 `json:"maxmem"`
	MaxDisk int64 `json:"maxdisk"`
}

// NodeUsage sums allocated resources across a node's QEMU VMs and LXC
// containers.
//
// Guests whose vmid is in exclude are skipped for the capacity sums: those are
// VMs Waggle already accounts for via its placement ledger, so counting them
// here too would double-subtract their capacity at scheduling time. A nil
// exclude counts every guest. vmid is cluster-unique in Proxmox, so a flat set
// suffices — no per-node keying is needed. Skipped vmids are still reported in
// NodeUsage.Managed so the caller can tell WHERE each managed guest actually
// runs, which is the only signal that catches a VM migrated out from under its
// placement row.
//
// A guest endpoint that errors (e.g. permissions, timeout) contributes nothing
// to the sums; if EVERY kind fails the node's usage is unknowable and an error
// is returned rather than a spuriously low zero. Reporting zero as though it
// were measured is what lets discovery overwrite a correct ram_gb_used with a
// value that makes a full node look empty.
func (c *Client) NodeUsage(ctx context.Context, node string, exclude map[int]struct{}) (NodeUsage, error) {
	var usage NodeUsage
	kinds := []string{"qemu", "lxc"}
	var errs []error
	for _, kind := range kinds {
		var payload struct {
			Data []guestRow `json:"data"`
		}
		if err := c.getJSON(ctx, "/api2/json/nodes/"+node+"/"+kind, &payload); err != nil {
			// Tolerate per-kind failures; report only if every kind is unusable.
			errs = append(errs, fmt.Errorf("%s: %w", kind, err))
			continue
		}
		for _, g := range payload.Data {
			if _, managed := exclude[g.VMID]; managed {
				usage.Managed = append(usage.Managed, g.VMID)
				continue
			}
			usage.VCPU += g.CPUs
			usage.MemBytes += g.MaxMem
			usage.DiskBytes += g.MaxDisk
			usage.Guests++
		}
	}
	if len(errs) == len(kinds) {
		return NodeUsage{}, fmt.Errorf("proxmox: node %s: no guest listing succeeded: %w", node, errors.Join(errs...))
	}
	return usage, nil
}
