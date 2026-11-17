## Context

The `tencentcloud_teo_zone_setting` resource manages zone-level settings for TencentCloud EdgeOne (TEO). It currently supports 16 configurable parameters (cache, cache_key, max_age, offline_cache, quic, post_max_size, compression, upstream_http2, force_redirect, https, origin, smart_routing, web_socket, client_ip_header, cache_prefresh, ipv6). The cloud API (`DescribeZoneSetting`/`ModifyZoneSetting`) also exposes a `JITVideoProcess` field that is not yet mapped in the Terraform resource.

The `JITVideoProcess` struct in the SDK (`github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901`) contains a single `Switch` field (string, values: "on"/"off") controlling whether JIT video processing is enabled for the zone.

## Goals / Non-Goals

**Goals:**
- Add `jit_video_process` as an Optional+Computed parameter to `tencentcloud_teo_zone_setting`
- Support reading the field from `DescribeZoneSetting` response (`ZoneSetting.JITVideoProcess`)
- Support writing the field via `ModifyZoneSetting` request (`Request.JITVideoProcess`)
- Maintain backward compatibility (existing configurations without this field continue to work)
- Add unit tests using gomonkey mock approach
- Update resource documentation (.md file)

**Non-Goals:**
- Adding other missing parameters (client_ip_country, grpc, network_error_logging, image_optimize, standard_debug) — those are separate changes
- Modifying the resource's Create or Delete methods (this resource uses CreateWithoutTimeout pattern and Delete is a no-op)

## Decisions

1. **Schema type: TypeList with MaxItems=1** — Consistent with all other switch-based parameters in this resource (ipv6, upstream_http2, quic, etc.). The inner element has a single `switch` field of type string.

2. **Optional + Computed** — The field is Optional so existing configurations don't break, and Computed so the API-returned value is stored in state when not explicitly set. This matches the pattern used by `ipv6`, `upstream_http2`, and other similar fields.

3. **Placement in mutableArgs array** — The field name `jit_video_process` will be appended to the `mutableArgs` slice in the Update method to trigger the ModifyZoneSetting API call when this field changes.

4. **Read pattern** — Follow the same nil-check pattern as `ipv6`: check if `respData.JITVideoProcess != nil`, then check if `respData.JITVideoProcess.Switch != nil` before setting the value.

5. **Update pattern** — Follow the same pattern as `ipv6`: use `helper.InterfacesHeadMap(d, "jit_video_process")` to extract the map, create a `teo.JITVideoProcess{}` struct, set the Switch field, and assign to `request.JITVideoProcess`.

## Risks / Trade-offs

- [Risk] Users who import existing zone settings will see a diff if JIT video processing was previously configured via console → Mitigation: Using Computed flag ensures the API value is read into state on first refresh, preventing spurious diffs.
- [Risk] The field description in the SDK says "不填写表示保持原有配置" (not filling means keeping existing config) → Mitigation: Only include in the request when the field has changed (controlled by mutableArgs check), matching existing behavior.
