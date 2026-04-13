# Design: tencentcloud_config_compliance_pack

## Directory Structure

```
tencentcloud/services/cfg/
├── service_tencentcloud_config.go
├── resource_tc_config_compliance_pack.go
├── resource_tc_config_compliance_pack.md
└── resource_tc_config_compliance_pack_test.go
```

## SDK Package

`github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/config/v20220802`

Import alias: `configv20220802`

## Resource Schema

| Field | Type | Required | Computed | Description |
|---|---|---|---|---|
| `compliance_pack_name` | String | Yes | No | Compliance pack name |
| `risk_level` | Int | Yes | No | Risk level: 1=high, 2=medium, 3=low |
| `config_rules` | List | Yes | No | List of compliance pack rules |
| `description` | String | No | No | Description |
| `status` | String | No | Yes | Compliance pack status: ACTIVE/UN_ACTIVE |
| `compliance_pack_id` | String | Computed | Yes | The unique ID returned by AddCompliancePack |
| `create_time` | String | Computed | Yes | Creation time |

### config_rules block

| Field | Type | Required | Description |
|---|---|---|---|
| `identifier` | String | Yes | Rule identifier (managed rule name or cloud function ARN) |
| `rule_name` | String | No | Rule name |
| `description` | String | No | Rule description |
| `risk_level` | Int | No | Rule risk level |
| `managed_rule_identifier` | String | No | Managed rule identifier |
| `input_parameter` | List | No | Rule input parameters |

### input_parameter block

| Field | Type | Required | Description |
|---|---|---|---|
| `parameter_key` | String | Yes | Parameter key |
| `type` | String | No | Parameter type (Require/Optional) |
| `value` | String | No | Parameter value |

## Key Design Decisions

1. **Resource ID**: Use `CompliancePackId` directly (not composite key).
2. **Status Update**: `UpdateCompliancePackStatus` is called separately when `status` changes.
3. **Delete Pre-condition**: Before deleting, set status to `UN_ACTIVE` via `UpdateCompliancePackStatus`.
4. **Style**: Follow `tencentcloud_igtm_strategy` style for retry logic, service layer pattern, and CRUD structure.

## Client Registration

Add to `tencentcloud/connectivity/client.go`:
- Import `configv20220802 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/config/v20220802"`
- Add `configv20220802Conn *configv20220802.Client` field
- Add `UseConfigV20220802Client()` method
