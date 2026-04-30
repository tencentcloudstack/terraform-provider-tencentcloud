## Context

Terraform Provider for TencentCloud currently supports TEO (EdgeOne) multi-path gateway resources (e.g., `tencentcloud_teo_multi_path_gateway_line`, `tencentcloud_teo_multi_path_gateway_secret_key`) but lacks a data source to query the multi-path gateway origin ACL details. The `DescribeMultiPathGatewayOriginACL` API is available in the SDK and provides origin ACL information including current and next version ACLs with IPv4/IPv6 CIDR lists.

The existing `tencentcloud_teo_origin_acl` data source queries the regular (L4/L7) origin ACL, which has a different API (`DescribeOriginACL`) and different response structure. The new data source is specific to multi-path gateways and requires a `gateway_id` parameter.

## Goals / Non-Goals

**Goals:**
- Add `tencentcloud_teo_multi_path_gateway_origin_acl` data source using `DescribeMultiPathGatewayOriginACL` API
- Follow existing TEO datasource patterns (e.g., `data_source_tc_teo_origin_acl.go`)
- Properly handle nested response structures (MultiPathGatewayCurrentOriginACL, MultiPathGatewayNextOriginACL, Addresses)
- Register the data source in provider.go and provider.md
- Add unit tests using gomonkey mock approach
- Add documentation .md file

**Non-Goals:**
- Modifying the existing `tencentcloud_teo_origin_acl` data source or resource
- Adding CRUD operations (this is a read-only data source)
- Adding filters or pagination (the API returns a single result per zone+gateway)

## Decisions

1. **Schema design follows cloud API structure**: The `multi_path_gateway_origin_acl_info` output block mirrors the `MultiPathGatewayOriginACLInfo` SDK struct with nested blocks for `multi_path_gateway_current_origin_acl` and `multi_path_gateway_next_origin_acl`. This matches the pattern used in `data_source_tc_teo_origin_acl.go`.

2. **Address blocks use TypeSet for ipv4/ipv6**: Following the existing `tencentcloud_teo_origin_acl` pattern, IPv4 and IPv6 lists use `schema.TypeSet` since order is not significant.

3. **Single object output (not list)**: The API returns a single `MultiPathGatewayOriginACLInfo` object, not an array. The output block is defined as `TypeList` with `MaxItems: 1` to match the Terraform SDK pattern for single nested objects, consistent with `data_source_tc_teo_origin_acl.go`.

4. **Resource ID uses composite key**: Since both `zone_id` and `gateway_id` are required to identify the resource, the data source ID will be set as `zone_id + FILED_SP + gateway_id`.

5. **Service method pattern**: Add `DescribeTeoMultiPathGatewayOriginAclByFilter` to `service_tencentcloud_teo.go`, following the existing service method pattern with `paramMap` for input parameters.

6. **No pagination needed**: The `DescribeMultiPathGatewayOriginACL` API returns a single result (not paginated), so no pagination logic is required.

## Risks / Trade-offs

- **API field nil safety**: The response has deeply nested structs where any level could be nil. The Read function must check nil at each level before accessing child fields to avoid panics. → Mitigation: Follow the nil-check pattern from `data_source_tc_teo_origin_acl.go`.
- **No deprecated field mapping needed**: Unlike the existing `tencentcloud_teo_origin_acl` which has deprecated `i_pv4`/`i_pv6` fields, the new API only has `IPv4`/`IPv6` in the `Addresses` struct. The new data source will only use `ipv4`/`ipv6` schema names without deprecated aliases.
