package neburst

import "fmt"

// Instance represents a compute instance (cloud or bare-metal).
type Instance struct {
	UUID        string        `json:"uuid"`
	Name        string        `json:"name"`
	Type        string        `json:"type"`
	Status      string        `json:"status"`
	Region      string        `json:"region,omitempty"`
	Hostname    string        `json:"hostname,omitempty"`
	PayCycle    string        `json:"pay_cycle,omitempty"`
	AutoRenew   bool          `json:"auto_renew"`
	NextPayAt   *string       `json:"next_pay_at,omitempty"`
	CreatedAt   string        `json:"created_at"`
	PrimaryIPv4 string        `json:"primary_ipv4,omitempty"`
	IPv4List    []string      `json:"ipv4_list,omitempty"`
	IPv6List    []string      `json:"ipv6_list,omitempty"`
	Specs       *InstanceSpecs `json:"specs,omitempty"`
	OSName      string        `json:"os_name,omitempty"`
}

// InstanceSpecs contains hardware specifications.
type InstanceSpecs struct {
	CPUModel         string   `json:"cpu_model,omitempty"`
	CPUCores         int      `json:"cpu_cores,omitempty"`
	MemoryGB         int      `json:"memory_gb,omitempty"`
	Disks            []Disk   `json:"disks,omitempty"`
	NetworkSpeedGbps float64  `json:"network_speed_gbps,omitempty"`
}

// Disk represents a disk in the instance configuration.
type Disk struct {
	Type     string `json:"type"`
	SizeGB   int    `json:"size_gb"`
	Quantity int    `json:"quantity"`
}

// PowerStatus represents the current power state of an instance.
type PowerStatus struct {
	Status       string `json:"status"`
	IsInstalling bool   `json:"is_installing"`
}

// Traffic contains all traffic packages for an instance.
type Traffic struct {
	Packages []TrafficPackage `json:"packages"`
}

// TrafficPackage represents a single traffic quota package.
type TrafficPackage struct {
	Name       string  `json:"name"`
	CapacityGB int     `json:"capacity_gb"`
	UsedGB     float64 `json:"used_gb"`
	ResetCycle string  `json:"reset_cycle,omitempty"`
}

// OSProfile represents an available OS profile for rebuild or rescue.
type OSProfile struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Category string `json:"category"`
	IsRescue bool   `json:"is_rescue"`
	Features struct {
		AllowSSHKeys     bool `json:"allow_ssh_keys"`
		AllowSetHostname bool `json:"allow_set_hostname"`
	} `json:"features"`
}

// Balance represents the user's account balance.
type Balance struct {
	Available float64 `json:"available"`
	Locked    float64 `json:"locked"`
	Currency  string  `json:"currency"`
}

// Invoice represents a billing invoice.
type Invoice struct {
	UUID      string  `json:"uuid"`
	Amount    float64 `json:"amount"`
	Status    string  `json:"status"`
	Category  string  `json:"category"`
	CreatedAt string  `json:"created_at"`
	DueAt     string  `json:"due_at,omitempty"`
}

// Metrics contains instance performance metrics.
type Metrics struct {
	CPU struct {
		Percentage float64 `json:"percentage"`
	} `json:"cpu"`
	Memory    ResourceUsage  `json:"memory"`
	Disk      ResourceUsage  `json:"disk"`
	Bandwidth BandwidthUsage `json:"bandwidth"`
	Network   NetworkUsage   `json:"network"`
}

// ResourceUsage represents usage of a limited resource (memory, disk).
type ResourceUsage struct {
	Limit      float64 `json:"limit"`
	Usage      float64 `json:"usage"`
	Free       float64 `json:"free"`
	Percentage float64 `json:"percentage"`
	Unit       string  `json:"unit"`
}

// BandwidthUsage represents bandwidth quota usage.
type BandwidthUsage struct {
	Limit       float64 `json:"limit"`
	Allowance   float64 `json:"allowance"`
	Usage       float64 `json:"usage"`
	Inbound     float64 `json:"inbound"`
	Outbound    float64 `json:"outbound"`
	Free        float64 `json:"free"`
	Percentage  float64 `json:"percentage"`
	UsageUnit   string  `json:"usage_unit"`
	LimitUnit   string  `json:"limit_unit"`
	StartedTime string  `json:"started_time"`
	EndTime     string  `json:"end_time"`
}

// NetworkUsage represents current network throughput.
type NetworkUsage struct {
	Inbound  float64 `json:"inbound"`
	Outbound float64 `json:"outbound"`
	Unit     string  `json:"unit"`
}

// PaginatedResult holds a page of results with metadata.
type PaginatedResult[T any] struct {
	Items    []T `json:"items"`
	Total    int `json:"total"`
	Page     int `json:"page"`
	PageSize int `json:"page_size"`
}

// apiResponse is the generic API response envelope.
type apiResponse[T any] struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data T      `json:"data"`
}

// APIError is returned when the API responds with a non-zero code.
type APIError struct {
	Code    int
	Message string
}

func (e *APIError) Error() string {
	return "neburst: api error " + e.Message
}

// Power action constants.
const (
	PowerOn    = "power-on"
	PowerOff   = "power-off"
	PowerCycle = "power-cycle"
	PowerReset = "power-reset"
)

// RebuildOption configures optional parameters for RebuildInstance.
type RebuildOption func(m map[string]any)

// WithHostname sets the hostname for a rebuild operation.
func WithHostname(hostname string) RebuildOption {
	return func(m map[string]any) {
		m["hostname"] = hostname
	}
}

// WithPublicKeys sets the SSH public keys for a rebuild operation.
func WithPublicKeys(keys []string) RebuildOption {
	return func(m map[string]any) {
		m["public_keys"] = keys
	}
}

// ListOption configures optional parameters for list operations.
type ListOption func(q map[string]string)

// WithPage sets the page number for a list operation.
func WithPage(page int) ListOption {
	return func(q map[string]string) {
		q["page"] = fmt.Sprintf("%d", page)
	}
}

// WithPageSize sets the page size for a list operation.
func WithPageSize(size int) ListOption {
	return func(q map[string]string) {
		q["page_size"] = fmt.Sprintf("%d", size)
	}
}
