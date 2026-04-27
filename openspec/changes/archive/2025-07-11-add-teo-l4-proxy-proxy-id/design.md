## Context

The `tencentcloud_teo_l4_proxy` resource manages TencentCloud EdgeOne L4 proxy instances. Currently, the resource stores `proxy_id` as part of the composite resource ID (format: `zone_id#proxy_id`), but does not expose it as a separate schema attribute. This means users cannot directly reference `proxy_id` in their Terraform configurations (e.g., for outputs or as input to other resources).

The `CreateL4Proxy` API already returns `ProxyId` in its response, and the `DescribeL4Proxy` API returns `ProxyId` in the `L4Proxy` struct. No API changes are needed.

## Goals / Non-Goals

**Goals:**
- Add `proxy_id` as a computed string attribute to the `tencentcloud_teo_l4_proxy` resource schema
- Populate `proxy_id` in the Create function from the `CreateL4Proxy` API response
- Populate `proxy_id` in the Read function from the `DescribeL4Proxy` API response
- Maintain full backward compatibility with existing configurations

**Non-Goals:**
- Changing the composite ID format (zone_id#proxy_id)
- Modifying any existing schema attributes
- Adding any other new parameters beyond `proxy_id`

## Decisions

1. **Schema attribute type**: `proxy_id` will be a `Computed: true` string field, since it is read-only and set by the cloud API. Users cannot set it directly.

2. **Value source in Create**: After the `CreateL4Proxy` API call succeeds, `proxy_id` will be set from `response.Response.ProxyId`. This value is already available in the current code (line 170 of `resource_tc_teo_l4_proxy.go`).

3. **Value source in Read**: After the `DescribeL4Proxy` API call succeeds, `proxy_id` will be set from `respData.ProxyId`. The `L4Proxy` struct already contains this field.

4. **No changes to Update/Delete**: `proxy_id` is a computed field and does not need to be handled in update or delete operations.

## Risks / Trade-offs

- **[Minimal risk]** Adding a computed attribute is backward compatible. Existing state files will have the new attribute populated on the next `terraform refresh` or `terraform plan`. → No mitigation needed.
