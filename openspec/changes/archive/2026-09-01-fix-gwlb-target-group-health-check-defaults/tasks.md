## 1. CRUD Fix

- [x] 1.1 In `resourceTencentCloudGwlbTargetGroupCreate`, read `timeout`/`interval_time`/`health_num`/`un_health_num` via `d.GetOk("health_check.0.<field>")` instead of the flattened `health_check` map, so omitted fields stay `nil` and are not sent to the API
- [x] 1.2 Apply the same change in `resourceTencentCloudGwlbTargetGroupUpdate`
- [x] 1.3 Keep the schema unchanged (fields remain `Optional: true, Computed: true`)

## 2. Test Update

- [x] 2.1 Add a test case in `resource_tc_gwlb_target_group_test.go` that creates a target group with a `health_check` block but omits `timeout`, `interval_time`, `health_num`, and `un_health_num`, and verifies the resource is created successfully with the API defaults reflected in state

## 3. Documentation

- [x] 3.1 No doc change required (schema is unchanged)
