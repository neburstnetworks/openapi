package com.neburst.api;

import com.google.gson.reflect.TypeToken;
import com.neburst.api.model.Balance;
import com.neburst.api.model.Invoice;

import java.io.IOException;
import java.util.List;

/**
 * API client for billing operations.
 */
public class BillingApi {

    private final NeburstClient client;

    public BillingApi(NeburstClient client) {
        this.client = client;
    }

    /**
     * Gets the current account balance.
     *
     * @return the account balance
     */
    public Balance getBalance() throws NeburstApiException, IOException, InterruptedException {
        return client.doRequest("GET", "/open/v1/billing/balance", null, null,
                Balance.class);
    }

    /**
     * Lists all invoices.
     *
     * @return list of invoices
     */
    public List<Invoice> listInvoices() throws NeburstApiException, IOException, InterruptedException {
        return client.doRequest("GET", "/open/v1/billing/invoices", null, null,
                new TypeToken<List<Invoice>>() {}.getType());
    }

    /**
     * Gets a single invoice by ID.
     *
     * @param id the invoice UUID
     * @return the invoice details
     */
    public Invoice getInvoice(String id) throws NeburstApiException, IOException, InterruptedException {
        return client.doRequest("GET", "/open/v1/billing/invoices/" + id, null, null,
                Invoice.class);
    }
}
