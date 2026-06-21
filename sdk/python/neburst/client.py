"""Core HTTP client with HMAC-SHA256 request signing."""

import base64
import hashlib
import hmac
import json
import time
import uuid
from typing import Any, Dict, List, Optional
from urllib.parse import quote_plus

import requests

from .types import APIError

_EMPTY_BODY_HASH = "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"


def _parse_combined_key(combined: str) -> tuple:
    try:
        data = json.loads(base64.b64decode(combined))
        if "key_id" in data and "secret" in data:
            return data["key_id"], data["secret"]
    except Exception:
        pass
    return combined, ""


class NeburstClient:
    """Neburst OpenAPI client.

    Supports two calling styles::

        NeburstClient("https://api.neburst.com", "nb_key_...", "nb_secret_...")
        NeburstClient("https://api.neburst.com", "eyJrZXlfaWQ...base64...")
    """

    def __init__(
        self,
        base_url: str,
        key_id_or_combined: str,
        secret: Optional[str] = None,
        timeout: int = 30,
    ) -> None:
        self.base_url = base_url.rstrip("/")
        if secret is not None:
            self.key_id = key_id_or_combined
            self.secret = secret
        else:
            self.key_id, self.secret = _parse_combined_key(key_id_or_combined)
        self.timeout = timeout
        self._session = requests.Session()

    def _sign(self, method: str, path: str, query_string: str, body: bytes) -> Dict[str, str]:
        timestamp = str(int(time.time()))
        nonce = str(uuid.uuid4())
        body_hash = hashlib.sha256(body).hexdigest() if body else _EMPTY_BODY_HASH
        string_to_sign = timestamp + "\n" + method + "\n" + path + "\n" + query_string + "\n" + body_hash
        signature = hmac.new(self.secret.encode(), string_to_sign.encode(), hashlib.sha256).hexdigest()
        return {
            "X-Nb-Key": self.key_id,
            "X-Nb-Timestamp": timestamp,
            "X-Nb-Nonce": nonce,
            "X-Nb-Signature": signature,
        }

    @staticmethod
    def _sorted_query_string(params: Optional[Dict[str, Any]]) -> str:
        if not params:
            return ""
        parts = []  # type: List[str]
        for key in sorted(params.keys()):
            value = params[key]
            if value is None:
                continue
            parts.append(quote_plus(str(key)) + "=" + quote_plus(str(value)))
        return "&".join(parts)

    def _request(self, method: str, path: str, params: Optional[Dict[str, Any]] = None, json_body: Optional[Any] = None) -> Any:
        full_path = "/open/v1" + path
        body = json.dumps(json_body, separators=(",", ":")).encode() if json_body is not None else b""
        query_string = self._sorted_query_string(params)
        headers = self._sign(method.upper(), full_path, query_string, body)
        if json_body is not None:
            headers["Content-Type"] = "application/json"
        url = self.base_url + full_path
        if query_string:
            url += "?" + query_string
        resp = self._session.request(method=method.upper(), url=url, headers=headers, data=body if body else None, timeout=self.timeout)
        resp.raise_for_status()
        envelope = resp.json()
        if envelope.get("code", 0) != 0:
            raise APIError(code=envelope["code"], msg=envelope.get("msg", ""))
        return envelope.get("data")

    # ── Cloud Instance (/compute/instance/*) ──

    def list_instances(self, page: int = 1, page_size: int = 20) -> "PaginatedResult":
        from .types import Instance, PaginatedResult
        data = self._request("GET", "/compute/instance/list", params={"page": page, "page_size": page_size})
        return PaginatedResult.from_dict(data, Instance)

    def get_instance(self, instance_id: str) -> "Instance":
        from .types import Instance
        return Instance.from_dict(self._request("GET", "/compute/instance/{}".format(instance_id)))

    def get_instance_status(self, instance_id: str) -> "PowerStatus":
        from .types import PowerStatus
        return PowerStatus.from_dict(self._request("GET", "/compute/instance/{}/status".format(instance_id)))

    def get_instance_traffic(self, instance_id: str) -> "Traffic":
        from .types import Traffic
        return Traffic.from_dict(self._request("GET", "/compute/instance/{}/traffic".format(instance_id)))

    def cloud_power_action(self, instance_id: str, action: str) -> None:
        self._request("POST", "/compute/instance/{}/power".format(instance_id), json_body={"action": action})

    def get_cloud_metrics(self, instance_id: str) -> "Metrics":
        from .types import Metrics
        return Metrics.from_dict(self._request("GET", "/compute/instance/{}/metrics".format(instance_id)))

    # ── Bare Metal (/compute/bare-metal/*) ──

    def list_bare_metal_instances(self, page: int = 1, page_size: int = 20) -> "PaginatedResult":
        from .types import Instance, PaginatedResult
        data = self._request("GET", "/compute/bare-metal/list", params={"page": page, "page_size": page_size})
        return PaginatedResult.from_dict(data, Instance)

    def get_bare_metal_instance(self, instance_id: str) -> "Instance":
        from .types import Instance
        return Instance.from_dict(self._request("GET", "/compute/bare-metal/{}".format(instance_id)))

    def get_bare_metal_status(self, instance_id: str) -> "PowerStatus":
        from .types import PowerStatus
        return PowerStatus.from_dict(self._request("GET", "/compute/bare-metal/{}/status".format(instance_id)))

    def get_bare_metal_traffic(self, instance_id: str) -> "Traffic":
        from .types import Traffic
        return Traffic.from_dict(self._request("GET", "/compute/bare-metal/{}/traffic".format(instance_id)))

    def bare_metal_power_action(self, instance_id: str, action: str) -> None:
        self._request("POST", "/compute/bare-metal/{}/power".format(instance_id), json_body={"action": action})

    def get_bare_metal_metrics(self, instance_id: str) -> "Metrics":
        from .types import Metrics
        return Metrics.from_dict(self._request("GET", "/compute/bare-metal/{}/metrics".format(instance_id)))

    def get_reinstall_profiles(self, instance_id: str) -> List["OSProfile"]:
        from .types import OSProfile
        data = self._request("GET", "/compute/bare-metal/{}/profiles".format(instance_id))
        return [OSProfile.from_dict(p) for p in data] if data else []

    def get_rescue_profiles(self, instance_id: str) -> List["OSProfile"]:
        from .types import OSProfile
        data = self._request("GET", "/compute/bare-metal/{}/rescue-profiles".format(instance_id))
        return [OSProfile.from_dict(p) for p in data] if data else []

    def rebuild_instance(self, instance_id: str, profile_id: int, hostname: Optional[str] = None, public_keys: Optional[List[str]] = None) -> None:
        body = {"profile_id": profile_id}  # type: Dict[str, Any]
        if hostname is not None:
            body["hostname"] = hostname
        if public_keys is not None:
            body["public_keys"] = public_keys
        self._request("POST", "/compute/bare-metal/{}/rebuild".format(instance_id), json_body=body)

    def rescue_instance(self, instance_id: str, profile_id: int) -> None:
        self._request("POST", "/compute/bare-metal/{}/rescue".format(instance_id), json_body={"profile_id": profile_id})

    # ── Billing (/billing/*) ──

    def get_balance(self) -> "Balance":
        from .types import Balance
        return Balance.from_dict(self._request("GET", "/billing/balance"))

    def list_invoices(self, page: int = 1, page_size: int = 20) -> "PaginatedResult":
        from .types import Invoice, PaginatedResult
        data = self._request("GET", "/billing/invoices", params={"page": page, "page_size": page_size})
        return PaginatedResult.from_dict(data, Invoice)

    def get_invoice(self, invoice_id: str) -> "Invoice":
        from .types import Invoice
        return Invoice.from_dict(self._request("GET", "/billing/invoices/{}".format(invoice_id)))

    # ── User (/user/*) ──

    def get_user_info(self) -> dict:
        return self._request("GET", "/user/info")
