use serde::{Deserialize, Serialize};
use thiserror::Error;

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct ApiResponse<T> {
    pub code: i32,
    pub msg: String,
    pub data: Option<T>,
}

#[derive(Error, Debug)]
pub enum NeburstError {
    #[error("API error (code={code}): {message}")]
    Api { code: i32, message: String },
    #[error("HTTP error: {0}")]
    Http(#[from] reqwest::Error),
    #[error("JSON error: {0}")]
    Json(#[from] serde_json::Error),
    #[error("{0}")]
    Other(String),
}

pub type Result<T> = std::result::Result<T, NeburstError>;

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct Instance {
    pub uuid: String,
    pub name: String,
    #[serde(rename = "type")]
    pub instance_type: String,
    pub status: String,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub region: Option<String>,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub hostname: Option<String>,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub pay_cycle: Option<String>,
    pub auto_renew: bool,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub next_pay_at: Option<String>,
    pub created_at: String,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub primary_ipv4: Option<String>,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub ipv4_list: Option<Vec<String>>,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub ipv6_list: Option<Vec<String>>,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub os_name: Option<String>,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub specs: Option<InstanceSpecs>,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct InstanceSpecs {
    #[serde(skip_serializing_if = "Option::is_none")]
    pub cpu_model: Option<String>,
    #[serde(default)]
    pub cpu_cores: i32,
    #[serde(default)]
    pub memory_gb: i32,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub disks: Option<Vec<Disk>>,
    #[serde(default)]
    pub network_speed_gbps: f64,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct Disk {
    #[serde(rename = "type")]
    pub disk_type: String,
    pub size_gb: i32,
    pub quantity: i32,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct PowerStatus {
    pub status: String,
    pub is_installing: bool,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct Traffic {
    pub packages: Vec<TrafficPackage>,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct TrafficPackage {
    pub name: String,
    pub capacity_gb: i64,
    pub used_gb: f64,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub reset_cycle: Option<String>,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct Metrics {
    pub cpu: CpuMetrics,
    pub memory: ResourceMetrics,
    pub disk: ResourceMetrics,
    pub bandwidth: BandwidthMetrics,
    pub network: NetworkMetrics,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct CpuMetrics {
    pub percentage: f64,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct ResourceMetrics {
    pub limit: f64,
    pub usage: f64,
    pub free: f64,
    pub percentage: f64,
    #[serde(default)]
    pub unit: String,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct BandwidthMetrics {
    pub limit: f64,
    pub allowance: f64,
    pub usage: f64,
    pub inbound: f64,
    pub outbound: f64,
    pub free: f64,
    pub percentage: f64,
    #[serde(default)]
    pub usage_unit: String,
    #[serde(default)]
    pub limit_unit: String,
    #[serde(default)]
    pub started_time: String,
    #[serde(default)]
    pub end_time: String,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct NetworkMetrics {
    pub inbound: f64,
    pub outbound: f64,
    #[serde(default)]
    pub unit: String,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct OSProfile {
    pub id: i32,
    pub name: String,
    pub category: String,
    pub is_rescue: bool,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct PaginatedResult<T> {
    pub items: Vec<T>,
    pub total: i32,
    pub page: i32,
    pub page_size: i32,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct Balance {
    pub available: f64,
    pub locked: f64,
    pub currency: String,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct Invoice {
    pub uuid: String,
    pub amount: f64,
    pub status: String,
    pub category: String,
    pub created_at: String,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub due_at: Option<String>,
}
