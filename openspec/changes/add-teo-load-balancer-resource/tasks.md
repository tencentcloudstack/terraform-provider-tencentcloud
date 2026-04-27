## 1. Service Layer

- [x] 1.1 Append `DescribeTeoLoadBalancerById(ctx, zoneId, instanceId)` to `tencentcloud/services/teo/service_tencentcloud_teo.go` — wraps `DescribeLoadBalancerList` with `Filters=[{Name:"InstanceId",Values:[instanceId]}]` and Retry; returns `*teo.LoadBalancer` or nil

## 2. Resource Implementation

- [x] 2.1 Create `tencentcloud/services/teo/resource_tc_teo_load_balancer.go` with full schema following `tencentcloud_igtm_strategy` style:
  - Schema fields: `zone_id` (Required, ForceNew), `name` (Required), `type` (Required, ForceNew), `origin_groups` (Required, List), `health_checker` (Optional, List MaxItems:1), `steering_policy` (Optional), `failover_policy` (Optional), `instance_id` (Computed)
  - `health_checker` sub-fields: `type`, `port`, `interval`, `timeout`, `health_threshold`, `critical_threshold`, `path`, `method`, `expected_codes`, `follow_redirect`, `headers` (List: `key`,`value`), `send_context`, `recv_context`
  - `origin_groups` sub-fields: `priority`, `origin_group_id`

- [x] 2.2 Implement Create:
  - Build `CreateLoadBalancerRequest` from schema fields
  - Call `CreateLoadBalancerWithContext`; extract `InstanceId` from response
  - Set resource ID to `strings.Join([]string{zoneId, instanceId}, tccommon.FILED_SP)`
  - Call Read

- [x] 2.3 Implement Read:
  - Split ID into `zoneId` and `instanceId`
  - Call `DescribeTeoLoadBalancerById`; if nil → `d.SetId("")`
  - Populate all schema fields from `LoadBalancer` struct: `name`, `type`, `steering_policy`, `failover_policy`, `instance_id`, `origin_groups`, `health_checker`

- [x] 2.4 Implement Update:
  - Build `ModifyLoadBalancerRequest` with `ZoneId` + `InstanceId`
  - Set all updatable fields (name, origin_groups, health_checker, steering_policy, failover_policy) from d
  - Call `ModifyLoadBalancerWithContext`
  - Call Read

- [x] 2.5 Implement Delete:
  - Build `DeleteLoadBalancerRequest` with `ZoneId` + `InstanceId`
  - Call `DeleteLoadBalancerWithContext` with Retry

## 3. Provider Registration

- [x] 3.1 Register `tencentcloud_teo_load_balancer` in `tencentcloud/provider.go` ResourcesMap, pointing to `teo.ResourceTencentCloudTeoLoadBalancer()`

## 4. Documentation & Tests

- [x] 4.1 Create `tencentcloud/services/teo/resource_tc_teo_load_balancer.md` — document all arguments, attributes, and import syntax with example HCL
- [x] 4.2 Create `tencentcloud/services/teo/resource_tc_teo_load_balancer_test.go` — basic acceptance test covering create/update/import/delete following `resource_tc_igtm_strategy_test.go` style
