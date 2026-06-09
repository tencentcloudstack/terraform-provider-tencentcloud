# Design: Config Compliance Pack Data Sources

## Architecture

Follows the standard terraform-provider-tencentcloud pattern for data sources (reference: `tencentcloud_igtm_instance_list`):

```
provider.go (registration)
    └─ data_source_tc_config_compliance_packs.go        (DataSource schema + Read handler)
    └─ data_source_tc_system_config_compliance_packs.go (DataSource schema + Read handler)
           └─ service_tencentcloud_config.go            (service methods — appended)
                  └─ config SDK v20220802               (ListCompliancePacks / ListSystemCompliancePacks)
```

## File Layout

| File | Action |
|---|---|
| `tencentcloud/services/config/data_source_tc_config_compliance_packs.go` | New — data source 1 |
| `tencentcloud/services/config/data_source_tc_config_compliance_packs.md` | New — inline doc |
| `tencentcloud/services/config/data_source_tc_config_compliance_packs_test.go` | New — acceptance test |
| `tencentcloud/services/config/data_source_tc_system_config_compliance_packs.go` | New — data source 2 |
| `tencentcloud/services/config/data_source_tc_system_config_compliance_packs.md` | New — inline doc |
| `tencentcloud/services/config/data_source_tc_system_config_compliance_packs_test.go` | New — acceptance test |
| `tencentcloud/services/config/service_tencentcloud_config.go` | Modified — add two service methods |
| `tencentcloud/provider.go` | Modified — register two new data sources |

## Pagination Strategy

Both APIs use `Limit` + `Offset` pagination with a `Total` field in the response.
- Loop until all pages fetched: `offset += limit` until `offset >= total`
- Each iteration is wrapped in `resource.Retry(tccommon.ReadRetryTimeout, ...)` for retry on transient errors
- Limit value: use **100** (the practical maximum for config APIs; no explicit max documented, using a safe large value)

## Schema Design

### `tencentcloud_config_compliance_packs`

**Filters (Optional):**
- `compliance_pack_name` (String) — filter by name
- `risk_level` (List of Int) — filter by risk level (1/2/3)
- `status` (String) — filter by status (ACTIVE / NO_ACTIVE)
- `compliance_result` (List of String) — filter by compliance result
- `order_type` (String) — sort direction (asc / desc)

**Result:**
- `compliance_pack_list` (List) — list of `ConfigCompliancePack` objects
- `result_output_file` (String, Optional) — output file path

**ConfigCompliancePack fields:**
`compliance_pack_id`, `compliance_pack_name`, `status`, `risk_level`, `compliance_result`, `create_time`, `description`, `rule_count`, `no_compliant_names`

### `tencentcloud_system_config_compliance_packs`

**No filters** (API has no filter params)

**Result:**
- `compliance_pack_list` (List) — list of `SystemCompliancePack` objects
- `result_output_file` (String, Optional) — output file path

**SystemCompliancePack fields:**
`compliance_pack_id`, `compliance_pack_name`, `description`, `risk_level`, `config_rules` (nested list with `identifier`, `rule_name`, `description`, `risk_level`, `create_time`, `update_time`)

## Service Methods

```go
// Appended to service_tencentcloud_config.go

func (me *ConfigService) DescribeConfigCompliancePacksByFilter(ctx, paramMap) ([]*ConfigCompliancePack, error)
func (me *ConfigService) DescribeSystemConfigCompliancePacks(ctx) ([]*SystemCompliancePack, error)
```

Both methods implement paged fetching with retry inside the loop.
