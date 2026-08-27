## Context

Redis instances support password complexity policies (密码复杂度策略) that enforce security requirements on password creation and reset operations. The TencentCloud Redis API provides two endpoints for this:

- `DescribeInstancePasswordPolicy` — reads the current password policy for an instance
- `ModifyInstancePasswordPolicy` — updates the password policy for an instance

This resource is a **RESOURCE_KIND_CONFIG** type, meaning the configuration exists as long as the Redis instance exists. There is no separate creation or deletion of the policy itself — only reading and updating.

## Goals / Non-Goals

**Goals:**
- Provide a Terraform resource `tencentcloud_redis_instance_password_policy` to manage Redis instance password complexity policies
- Support reading the current policy via `DescribeInstancePasswordPolicy`
- Support updating the policy via `ModifyInstancePasswordPolicy`
- Follow the standard CONFIG resource pattern: Create sets the ID and calls Update, Read queries the policy, Update modifies the policy, Delete is a no-op

**Non-Goals:**
- This resource does NOT create or delete Redis instances
- This resource does NOT manage the Redis instance password itself (only the policy rules)

## Decisions

### 1. Resource ID: InstanceId

The resource ID is simply the Redis instance ID (`InstanceId`). Since the password policy is a configuration of the instance (not a separate entity), using the instance ID as the resource ID is clean and consistent with other CONFIG resources.

### 2. Schema Structure: Nested `password_policy` Block

The `password_policy` is modeled as a `TypeList` with `MaxItems: 1`, containing nested fields that mirror the `PasswordPolicy` struct from the SDK:

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `instance_id` | TypeString | Required, ForceNew | Redis instance ID |
| `password_policy` | TypeList | Required | Nested block (MaxItems: 1) |
| `password_policy.enabled` | TypeBool | Required | Whether password complexity is enabled |
| `password_policy.min_letter_count` | TypeInt | Optional | Minimum letter count (range: 1-16) |
| `password_policy.min_digit_count` | TypeInt | Optional | Minimum digit count (range: 1-16) |
| `password_policy.min_special_count` | TypeInt | Optional | Minimum special character count (range: 1-16) |
| `password_policy.min_length` | TypeInt | Optional | Minimum total password length (range: 8-64) |

### 3. CONFIG Pattern: Create → Update, Delete → No-op

Following the standard config resource pattern:
- **Create**: Set `d.SetId(instanceId)`, then call Update to apply the desired configuration
- **Read**: Call `DescribeInstancePasswordPolicy`, populate all fields from the response
- **Update**: Call `ModifyInstancePasswordPolicy` with the current configuration
- **Delete**: No-op (return nil); the policy configuration persists as long as the instance exists

### 4. API Calls with Retry

All API calls use `resource.Retry` with `tccommon.ReadRetryTimeout` / `tccommon.WriteRetryTimeout` for resilience. Errors are wrapped with `tccommon.RetryError()`.

### 5. Service Layer

A new `service_tencentcloud_redis.go` file will be created in `tencentcloud/services/redis/` to provide the service helper function `DescribeRedisInstancePasswordPolicyById`.

## Risks / Trade-offs

- **[Risk] Instance not found**: If the Redis instance is deleted outside of Terraform, the Read will return an error. Mitigation: Check for `ResourceNotFound.InstanceNotExists` error and call `d.SetId("")` to allow Terraform to detect the drift.
- **[Risk] API version compatibility**: The `PasswordPolicy` struct exists in the current SDK version. Mitigation: The SDK is already vendored and verified.