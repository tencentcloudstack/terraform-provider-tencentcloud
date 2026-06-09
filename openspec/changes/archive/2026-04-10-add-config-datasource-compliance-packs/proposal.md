# Add Config Compliance Pack Data Sources

## What

Add two new Terraform data sources for the Tencent Cloud Config (配置审计) service:

1. `tencentcloud_config_compliance_packs` — query user-created compliance packs via `ListCompliancePacks`
2. `tencentcloud_system_config_compliance_packs` — query system-provided compliance packs via `ListSystemCompliancePacks`

## Why

The existing `tencentcloud_config_compliance_pack` resource covers CRUD for user-defined compliance packs, but there are no data sources to list or query compliance packs. Users need to:

- Enumerate existing compliance packs (e.g., to reference IDs in other resources)
- Query system-provided compliance pack templates to understand available options
- Filter packs by name, risk level, status, or compliance result for operational reporting

## APIs Used

| Data Source | API | Notes |
|---|---|---|
| `tencentcloud_config_compliance_packs` | `ListCompliancePacks` | Paginated; supports filter by name, risk level, status, compliance result, order |
| `tencentcloud_system_config_compliance_packs` | `ListSystemCompliancePacks` | Paginated; no filters |

## Resource IDs

Data sources use `helper.BuildToken()` as computed ID (read-only, no real resource ID needed).
