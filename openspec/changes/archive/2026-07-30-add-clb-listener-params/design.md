## Context

The `tencentcloud_clb_listener` resource is a Terraform resource of type RESOURCE_KIND_GENERAL for managing CLB listeners. It supports CRUD operations via CreateListener, DescribeListeners, ModifyListener, and DeleteListener APIs.

The CLB SDK (`github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/clb/v20180317`) already contains the following four parameters across the relevant API structs:

- **MaxConn** (`*int64`): Listener-level maximum concurrent connections limit. Available in `CreateListenerRequest`, `ModifyListenerRequest`, and `Listener` (DescribeListeners response).
- **MaxCps** (`*int64`): Listener-level maximum new connections per second limit. Available in `CreateListenerRequest`, `ModifyListenerRequest`, and `Listener`.
- **ProxyProtocol** (`*bool`): Enable proxy protocol for TCP_SSL/QUIC listeners. Available in `CreateListenerRequest` and `ModifyListenerRequest`. **Not present in `Listener` struct** (not returned by DescribeListeners).
- **DataCompressMode** (`*string`): Data compression mode (`transparent` or `compatibility`). Available in `CreateListenerRequest`, `ModifyListenerRequest`, and `Listener`.

These parameters are validated against the vendor SDK (`vendor/github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/clb/v20180317/models.go`).

## Goals / Non-Goals

**Goals:**
- Add `MaxConn`, `MaxCps`, `ProxyProtocol`, and `DataCompressMode` as Optional schema fields to `tencentcloud_clb_listener`
- Wire these fields into Create, Read, and Update functions following existing patterns
- Ensure backward compatibility (all new fields are Optional)

**Non-Goals:**
- No validation logic beyond what the API enforces (no protocol-specific checks in Terraform)
- No changes to Delete behavior
- No changes to existing schema fields

## Decisions

### Schema Design

| Field | Type | Optional | Computed | Notes |
|-------|------|----------|----------|-------|
| `max_conn` | TypeInt | Yes | Yes | -1 means unlimited |
| `max_cps` | TypeInt | Yes | Yes | -1 means unlimited |
| `proxy_protocol` | TypeBool | Yes | No | Write-only: not in Listener struct |
| `data_compress_mode` | TypeString | Yes | Yes | `transparent` or `compatibility` |

**Rationale**: `MaxConn`, `MaxCps`, `DataCompressMode` are marked Computed because they are returned by DescribeListeners. `ProxyProtocol` is NOT Computed because it is absent from the `Listener` struct in the DescribeListeners response — it is a write-only parameter used in Create/Update only.

### Create Implementation

Following the existing pattern (lines 435-453), each new parameter is checked with `d.GetOkExists` and set on the request:

```go
if v, ok := d.GetOkExists("max_conn"); ok {
    request.MaxConn = helper.IntInt64(v.(int))
}
```

### Read Implementation

Following the existing pattern (lines 606-630), each parameter is set with nil-check before calling `d.Set`:

```go
if instance.MaxConn != nil {
    _ = d.Set("max_conn", instance.MaxConn)
}
```

`ProxyProtocol` is NOT set in Read because the `Listener` struct does not have this field.

### Update Implementation

Following the existing pattern (lines 756-784), each parameter is checked with `d.HasChange`:

```go
if d.HasChange("max_conn") {
    changeFlag = true
    request.MaxConn = helper.IntInt64(d.Get("max_conn").(int))
}
```

## Risks / Trade-offs

- **ProxyProtocol is write-only**: If a user sets `proxy_protocol = true` and later runs `terraform apply` again, Terraform will see the state contains `proxy_protocol = true` but the Read will not set it, potentially flagging a drift. This is an inherent limitation of the API (DescribeListeners does not return this field). Users should be aware via documentation.
- **No protocol validation**: The Terraform code does not validate that these parameters are only set for compatible protocols (e.g., MaxConn/MaxCps only for performance capacity-type instances with TCP/UDP/TCP_SSL/QUIC listeners). Invalid configurations will be rejected by the API at apply time.
