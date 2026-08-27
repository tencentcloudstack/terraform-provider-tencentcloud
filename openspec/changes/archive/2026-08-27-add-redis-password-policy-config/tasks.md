## 1. Service Layer

- [x] 1.1 Add `DescribeRedisInstancePasswordPolicyById` function to `tencentcloud/services/crs/service_tencentcloud_redis.go` that calls `DescribeInstancePasswordPolicy` API

## 2. Resource Implementation

- [x] 2.1 Create `resource_tc_redis_instance_password_policy_config.go` with schema definition (instance_id, password_policy nested block with enabled, min_letter_count, min_digit_count, min_special_count, min_length)
- [x] 2.2 Implement Create: set `d.SetId(instanceId)`, then call Update to apply configuration
- [x] 2.3 Implement Read: call `DescribeInstancePasswordPolicy` API, populate state from response, handle nil PasswordPolicy by clearing ID
- [x] 2.4 Implement Update: call `ModifyInstancePasswordPolicy` API with retry, then call Read to refresh state
- [x] 2.5 Implement Delete: no-op (return nil), CONFIG resource persists as long as instance exists

## 3. Provider Registration

- [x] 3.1 Register `tencentcloud_redis_instance_password_policy` in `tencentcloud/provider.go` ResourcesMap

## 4. Documentation & Tests

- [x] 4.1 Create `resource_tc_redis_instance_password_policy_config.md` with example usage
- [x] 4.2 Create `resource_tc_redis_instance_password_policy_config_test.go` with unit tests using gomonkey mocks