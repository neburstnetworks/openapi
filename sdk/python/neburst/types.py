"""Data types for the Neburst OpenAPI SDK."""

from dataclasses import dataclass, field
from typing import Any, Dict, List, Optional


@dataclass
class APIError(Exception):
    """Error returned by the Neburst API."""

    code: int
    msg: str

    def __str__(self) -> str:
        return "APIError(code={}, msg={!r})".format(self.code, self.msg)


@dataclass
class DiskInfo:
    """Disk configuration."""

    type: str = ""
    size_gb: int = 0
    quantity: int = 0

    @classmethod
    def from_dict(cls, data: Dict[str, Any]) -> "DiskInfo":
        return cls(
            type=data.get("type", ""),
            size_gb=data.get("size_gb", 0),
            quantity=data.get("quantity", 0),
        )


@dataclass
class InstanceSpecs:
    """Hardware specifications of an instance."""

    cpu_model: Optional[str] = None
    cpu_cores: int = 0
    memory_gb: int = 0
    disks: List[DiskInfo] = field(default_factory=list)
    network_speed_gbps: float = 0.0

    @classmethod
    def from_dict(cls, data: Dict[str, Any]) -> "InstanceSpecs":
        disks = [DiskInfo.from_dict(d) for d in data.get("disks", [])]
        return cls(
            cpu_model=data.get("cpu_model"),
            cpu_cores=data.get("cpu_cores", 0),
            memory_gb=data.get("memory_gb", 0),
            disks=disks,
            network_speed_gbps=data.get("network_speed_gbps", 0.0),
        )


@dataclass
class Instance:
    """A compute instance (bare metal or cloud)."""

    uuid: str
    name: str
    type: str
    status: str
    auto_renew: bool
    created_at: str
    region: Optional[str] = None
    hostname: Optional[str] = None
    pay_cycle: Optional[str] = None
    next_pay_at: Optional[str] = None
    primary_ipv4: Optional[str] = None
    ipv4_list: List[str] = field(default_factory=list)
    ipv6_list: List[str] = field(default_factory=list)
    specs: Optional[InstanceSpecs] = None
    os_name: Optional[str] = None

    @classmethod
    def from_dict(cls, data: Dict[str, Any]) -> "Instance":
        specs = None
        if data.get("specs"):
            specs = InstanceSpecs.from_dict(data["specs"])
        return cls(
            uuid=data["uuid"],
            name=data["name"],
            type=data["type"],
            status=data["status"],
            auto_renew=data.get("auto_renew", False),
            created_at=data["created_at"],
            region=data.get("region"),
            hostname=data.get("hostname"),
            pay_cycle=data.get("pay_cycle"),
            next_pay_at=data.get("next_pay_at"),
            primary_ipv4=data.get("primary_ipv4"),
            ipv4_list=data.get("ipv4_list", []),
            ipv6_list=data.get("ipv6_list", []),
            specs=specs,
            os_name=data.get("os_name"),
        )


@dataclass
class PowerStatus:
    """Power status of an instance."""

    status: str
    is_installing: bool

    @classmethod
    def from_dict(cls, data: Dict[str, Any]) -> "PowerStatus":
        return cls(
            status=data["status"],
            is_installing=data.get("is_installing", False),
        )


@dataclass
class TrafficPackage:
    """A single traffic package attached to an instance."""

    name: str
    capacity_gb: int
    used_gb: float
    reset_cycle: Optional[str] = None

    @classmethod
    def from_dict(cls, data: Dict[str, Any]) -> "TrafficPackage":
        return cls(
            name=data["name"],
            capacity_gb=data["capacity_gb"],
            used_gb=data["used_gb"],
            reset_cycle=data.get("reset_cycle"),
        )


@dataclass
class Traffic:
    """Traffic information for an instance."""

    packages: List[TrafficPackage] = field(default_factory=list)

    @classmethod
    def from_dict(cls, data: Dict[str, Any]) -> "Traffic":
        packages = [
            TrafficPackage.from_dict(p) for p in data.get("packages", [])
        ]
        return cls(packages=packages)


@dataclass
class OSProfileFeatures:
    """Features supported by an OS profile."""

    allow_ssh_keys: bool = False
    allow_set_hostname: bool = False

    @classmethod
    def from_dict(cls, data: Dict[str, Any]) -> "OSProfileFeatures":
        return cls(
            allow_ssh_keys=data.get("allow_ssh_keys", False),
            allow_set_hostname=data.get("allow_set_hostname", False),
        )


@dataclass
class OSProfile:
    """An available OS profile for rebuild or rescue."""

    id: int
    name: str
    category: str
    is_rescue: bool
    features: OSProfileFeatures = field(default_factory=OSProfileFeatures)

    @classmethod
    def from_dict(cls, data: Dict[str, Any]) -> "OSProfile":
        features = OSProfileFeatures.from_dict(data.get("features", {}))
        return cls(
            id=data["id"],
            name=data["name"],
            category=data["category"],
            is_rescue=data.get("is_rescue", False),
            features=features,
        )


@dataclass
class ResourceUsage:
    """Resource usage (memory, disk)."""

    limit: float = 0.0
    usage: float = 0.0
    free: float = 0.0
    percentage: float = 0.0
    unit: str = ""

    @classmethod
    def from_dict(cls, data: Dict[str, Any]) -> "ResourceUsage":
        return cls(
            limit=data.get("limit", 0.0),
            usage=data.get("usage", 0.0),
            free=data.get("free", 0.0),
            percentage=data.get("percentage", 0.0),
            unit=data.get("unit", ""),
        )


@dataclass
class BandwidthUsage:
    """Bandwidth quota usage."""

    limit: float = 0.0
    allowance: float = 0.0
    usage: float = 0.0
    inbound: float = 0.0
    outbound: float = 0.0
    free: float = 0.0
    percentage: float = 0.0
    usage_unit: str = ""
    limit_unit: str = ""
    started_time: str = ""
    end_time: str = ""

    @classmethod
    def from_dict(cls, data: Dict[str, Any]) -> "BandwidthUsage":
        return cls(**{k: data.get(k, getattr(cls, k, "")) for k in cls.__dataclass_fields__})


@dataclass
class NetworkUsage:
    """Current network throughput."""

    inbound: float = 0.0
    outbound: float = 0.0
    unit: str = ""

    @classmethod
    def from_dict(cls, data: Dict[str, Any]) -> "NetworkUsage":
        return cls(
            inbound=data.get("inbound", 0.0),
            outbound=data.get("outbound", 0.0),
            unit=data.get("unit", ""),
        )


@dataclass
class CPUUsage:
    """CPU usage."""

    percentage: float = 0.0

    @classmethod
    def from_dict(cls, data: Dict[str, Any]) -> "CPUUsage":
        return cls(percentage=data.get("percentage", 0.0))


@dataclass
class Metrics:
    """Instance performance metrics."""

    cpu: CPUUsage = field(default_factory=CPUUsage)
    memory: ResourceUsage = field(default_factory=ResourceUsage)
    disk: ResourceUsage = field(default_factory=ResourceUsage)
    bandwidth: BandwidthUsage = field(default_factory=BandwidthUsage)
    network: NetworkUsage = field(default_factory=NetworkUsage)

    @classmethod
    def from_dict(cls, data: Dict[str, Any]) -> "Metrics":
        return cls(
            cpu=CPUUsage.from_dict(data.get("cpu", {})),
            memory=ResourceUsage.from_dict(data.get("memory", {})),
            disk=ResourceUsage.from_dict(data.get("disk", {})),
            bandwidth=BandwidthUsage.from_dict(data.get("bandwidth", {})),
            network=NetworkUsage.from_dict(data.get("network", {})),
        )


@dataclass
class PaginatedResult:
    """Paginated result wrapper."""

    items: List[Any]
    total: int
    page: int
    page_size: int

    @classmethod
    def from_dict(cls, data: Dict[str, Any], item_cls: type) -> "PaginatedResult":
        items = [item_cls.from_dict(i) for i in data.get("items", [])]
        return cls(
            items=items,
            total=data.get("total", 0),
            page=data.get("page", 1),
            page_size=data.get("page_size", 20),
        )


@dataclass
class Balance:
    """Account balance."""

    available: float
    locked: float
    currency: str

    @classmethod
    def from_dict(cls, data: Dict[str, Any]) -> "Balance":
        return cls(
            available=data["available"],
            locked=data["locked"],
            currency=data["currency"],
        )


@dataclass
class Invoice:
    """A billing invoice."""

    uuid: str
    amount: float
    status: str
    category: str
    created_at: str
    due_at: Optional[str] = None

    @classmethod
    def from_dict(cls, data: Dict[str, Any]) -> "Invoice":
        return cls(
            uuid=data["uuid"],
            amount=data["amount"],
            status=data["status"],
            category=data["category"],
            created_at=data["created_at"],
            due_at=data.get("due_at"),
        )
