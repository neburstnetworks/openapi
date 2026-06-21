package neburst

// GetBalance returns the authenticated user's account balance.
func (c *Client) GetBalance() (*Balance, error) {
	var result Balance
	err := c.doRequest("GET", "/open/v1/billing/balance", nil, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// ListInvoices returns invoices for the authenticated user, with pagination.
func (c *Client) ListInvoices(opts ...ListOption) (*PaginatedResult[Invoice], error) {
	q := applyListOpts(opts)
	var result PaginatedResult[Invoice]
	err := c.doRequest("GET", "/open/v1/billing/invoices", q, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// GetInvoice returns a single invoice by its UUID.
func (c *Client) GetInvoice(id string) (*Invoice, error) {
	var result Invoice
	err := c.doRequest("GET", "/open/v1/billing/invoices/"+id, nil, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
