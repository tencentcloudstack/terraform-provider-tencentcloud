# Design: Config Resources Data Sources

## File Layout

| File | Action |
|---|---|
| `data_source_tc_config_discovered_resources.go` | New |
| `data_source_tc_config_discovered_resources.md` | New |
| `data_source_tc_config_discovered_resources_test.go` | New |
| `data_source_tc_config_resource_types.go` | New |
| `data_source_tc_config_resource_types.md` | New |
| `data_source_tc_config_resource_types_test.go` | New |
| `service_tencentcloud_config.go` | Modified |
| `provider.go` | Modified |

## Pagination Strategy

### ListDiscoveredResources
Uses `NextToken` cursor pagination (not Offset):
- Set `MaxResults=200` per page
- Loop: each iteration wrapped in `resource.Retry`; after success, set `request.NextToken = response.NextToken`; break when `NextToken` is nil or empty

### ListResourceTypes
No pagination — single call, no loop needed. Wrapped in `resource.Retry`.

## Schema: tencentcloud_config_discovered_resources

**Optional filters:**
- `filters` (List) — name/value filter (filter on `resourceName` or `resourceId`)
- `tags` (List) — tag key/value filter
- `order_type` (String) — `asc` or `desc`

**Computed:** `resource_list` (List of `ResourceListInfo`)
- `resource_type`, `resource_name`, `resource_id`, `resource_region`, `resource_status`, `resource_delete` (Int), `resource_create_time`, `resource_zone`, `compliance_result`, `tags` (List of key/value)

## Schema: tencentcloud_config_resource_types

**No filters** (API accepts no params)

**Computed:** `resource_type_list` (List of `ConfigResource`)
- `product`, `product_name`, `resource_type`, `resource_type_name`
