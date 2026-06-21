package neburst

import "net/url"

// ── Cloud Instance (/compute/instance/*) ──

func (c *Client) ListInstances(opts ...ListOption) (*PaginatedResult[Instance], error) {
	q := applyListOpts(opts)
	var result PaginatedResult[Instance]
	err := c.doRequest("GET", "/open/v1/compute/instance/list", q, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) GetInstance(id string) (*Instance, error) {
	var result Instance
	err := c.doRequest("GET", "/open/v1/compute/instance/"+id, nil, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) GetInstanceStatus(id string) (*PowerStatus, error) {
	var result PowerStatus
	err := c.doRequest("GET", "/open/v1/compute/instance/"+id+"/status", nil, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) GetInstanceTraffic(id string) (*Traffic, error) {
	var result Traffic
	err := c.doRequest("GET", "/open/v1/compute/instance/"+id+"/traffic", nil, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) CloudPowerAction(id, action string) error {
	body := map[string]string{"action": action}
	return c.doRequest("POST", "/open/v1/compute/instance/"+id+"/power", nil, body, nil)
}

func (c *Client) GetCloudMetrics(id string) (*Metrics, error) {
	var result Metrics
	err := c.doRequest("GET", "/open/v1/compute/instance/"+id+"/metrics", nil, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// ── Bare Metal (/compute/bare-metal/*) ──

func (c *Client) ListBareMetalInstances(opts ...ListOption) (*PaginatedResult[Instance], error) {
	q := applyListOpts(opts)
	var result PaginatedResult[Instance]
	err := c.doRequest("GET", "/open/v1/compute/bare-metal/list", q, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) GetBareMetalInstance(id string) (*Instance, error) {
	var result Instance
	err := c.doRequest("GET", "/open/v1/compute/bare-metal/"+id, nil, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) GetBareMetalStatus(id string) (*PowerStatus, error) {
	var result PowerStatus
	err := c.doRequest("GET", "/open/v1/compute/bare-metal/"+id+"/status", nil, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) GetBareMetalTraffic(id string) (*Traffic, error) {
	var result Traffic
	err := c.doRequest("GET", "/open/v1/compute/bare-metal/"+id+"/traffic", nil, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) BareMetalPowerAction(id, action string) error {
	body := map[string]string{"action": action}
	return c.doRequest("POST", "/open/v1/compute/bare-metal/"+id+"/power", nil, body, nil)
}

func (c *Client) GetBareMetalMetrics(id string) (*Metrics, error) {
	var result Metrics
	err := c.doRequest("GET", "/open/v1/compute/bare-metal/"+id+"/metrics", nil, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) GetReinstallProfiles(id string) ([]OSProfile, error) {
	var result []OSProfile
	err := c.doRequest("GET", "/open/v1/compute/bare-metal/"+id+"/profiles", nil, nil, &result)
	return result, err
}

func (c *Client) GetRescueProfiles(id string) ([]OSProfile, error) {
	var result []OSProfile
	err := c.doRequest("GET", "/open/v1/compute/bare-metal/"+id+"/rescue-profiles", nil, nil, &result)
	return result, err
}

func (c *Client) RebuildInstance(id string, profileID int, opts ...RebuildOption) error {
	body := map[string]any{"profile_id": profileID}
	for _, opt := range opts {
		opt(body)
	}
	return c.doRequest("POST", "/open/v1/compute/bare-metal/"+id+"/rebuild", nil, body, nil)
}

func (c *Client) RescueInstance(id string, profileID int) error {
	body := map[string]any{"profile_id": profileID}
	return c.doRequest("POST", "/open/v1/compute/bare-metal/"+id+"/rescue", nil, body, nil)
}

func applyListOpts(opts []ListOption) url.Values {
	if len(opts) == 0 {
		return nil
	}
	m := make(map[string]string)
	for _, opt := range opts {
		opt(m)
	}
	q := url.Values{}
	for k, v := range m {
		q.Set(k, v)
	}
	return q
}
