use std::collections::BTreeMap;
use crate::client::NeburstClient;
use crate::types::{Balance, Invoice, PaginatedResult, Result};

impl NeburstClient {
    pub async fn get_balance(&self) -> Result<Balance> {
        self.request("GET", "/open/v1/billing/balance", None, None::<&()>).await
    }

    pub async fn list_invoices(&self, page: i32, page_size: i32) -> Result<PaginatedResult<Invoice>> {
        let mut q = BTreeMap::new();
        q.insert("page".into(), page.to_string());
        q.insert("page_size".into(), page_size.to_string());
        self.request("GET", "/open/v1/billing/invoices", Some(&q), None::<&()>).await
    }

    pub async fn get_invoice(&self, id: &str) -> Result<Invoice> {
        self.request("GET", &format!("/open/v1/billing/invoices/{}", id), None, None::<&()>).await
    }
}
