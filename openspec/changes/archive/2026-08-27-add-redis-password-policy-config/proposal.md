## Why

Redis instances support password complexity policies (密码复杂度策略) to enforce security requirements on password changes. Currently, the Terraform provider has no resource to manage this configuration, forcing users to configure it manually via the console or API. This change adds a `tencentcloud_redis_instance_password_policy` resource to manage Redis instance password policies declaratively through Terraform.

## What Changes

- Add new Terraform resource `tencentcloud_redis_instance_password_policy` (RESOURCE_KIND_CONFIG) for managing Redis instance password complexity policies
- Resource uses `DescribeInstancePasswordPolicy` API for reading the current password policy configuration
- Resource uses `ModifyInstancePasswordPolicy` API for updating the password policy configuration
- No Create/Delete operations — as a CONFIG resource, the configuration exists as long as the instance exists; Create sets the ID and calls Update, Delete is a no-op

## Capabilities

### New Capabilities
- `redis-instance-password-policy`: Manage Redis instance password complexity policy configuration, including enabled status, minimum letter count, minimum digit count, minimum special character count, and minimum password length.

### Modified Capabilities
<!-- None -->

## Impact

- **Affected code**: `tencentcloud/services/redis/` (new service package), `tencentcloud/provider.go` (registration)
- **New files**: `resource_tc_redis_instance_password_policy_config.go`, `resource_tc_redis_instance_password_policy_config_test.go`, `resource_tc_redis_instance_password_policy_config.md`, `service_tencentcloud_redis.go`
- **SDK dependency**: `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/redis/v20180412` (already in vendor)
- **No breaking changes** to existing resources or configurations