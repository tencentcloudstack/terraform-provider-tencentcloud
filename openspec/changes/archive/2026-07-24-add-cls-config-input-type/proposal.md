## Why

The `tencentcloud_cls_config` resource is missing the `input_type` parameter, which is already supported by the CLS cloud API across `CreateConfig`, `ModifyConfig`, and `DescribeConfigs` (via `ConfigInfo`). This parameter is essential for users who need to configure Windows event log collection or syslog collection, where `input_type` is required to specify the input source type. Without this parameter, Terraform users cannot properly manage CLS configs for Windows event or syslog scenarios.

## What Changes

- Add a new optional `input_type` field of type `TypeString` to the `tencentcloud_cls_config` resource schema
- Support the `input_type` parameter in the Create (`CreateConfig`), Read (`DescribeConfigs`), and Update (`ModifyConfig`) operations
- Update the resource documentation (`.md`) to include the new parameter

## Capabilities

### New Capabilities
- `cls-config-input-type`: Support configuring the log input type (file, windows_event, syslog) for CLS collection configurations

### Modified Capabilities
<!-- No existing capabilities are being modified at the spec level -->

## Impact

- **Affected code**: `tencentcloud/services/cls/resource_tc_cls_config.go` (schema definition, Create, Read, Update methods)
- **Affected docs**: `tencentcloud/services/cls/resource_tc_cls_config.md`
- **API dependencies**: `CreateConfig`, `ModifyConfig`, `DescribeConfigs` (ConfigInfo response) from `tencentcloud-sdk-go/tencentcloud/cls/v20201016`
- **Backward compatibility**: Fully backward compatible — the new parameter is Optional with no default value