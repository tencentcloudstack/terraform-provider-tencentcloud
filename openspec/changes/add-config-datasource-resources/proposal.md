# Add Config Resources Data Sources

## What

Add two new Terraform data sources for the Tencent Cloud Config service:

1. `tencentcloud_config_discovered_resources` — query discovered resources via `ListDiscoveredResources`
2. `tencentcloud_config_resource_types` — query supported resource types via `ListResourceTypes`

## Why

Users managing cloud compliance need to:
- List all discovered resources tracked by Config (with optional name/ID filters and tag filters)
- Enumerate all resource types supported by Config for reference when building rules

## APIs Used

| Data Source | API | Pagination |
|---|---|---|
| `tencentcloud_config_discovered_resources` | `ListDiscoveredResources` | NextToken-based, MaxResults=200 |
| `tencentcloud_config_resource_types` | `ListResourceTypes` | No pagination (returns all at once) |
