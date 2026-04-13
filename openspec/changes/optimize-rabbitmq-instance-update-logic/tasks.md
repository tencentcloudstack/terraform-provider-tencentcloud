## 1. Schema Definition

- [x] 1.1 Add `remark` parameter to resource schema in `tencentcloud/services/trabbit/resource_tc_tdmq_rabbitmq_vip_instance.go`
  - Type: `TypeString`
  - Optional: `true`
  - Computed: `true`
  - Description: "Instance remark or description information"
- [x] 1.2 Add `enable_deletion_protection` parameter to resource schema in `tencentcloud/services/trabbit/resource_tc_tdmq_rabbitmq_vip_instance.go`
  - Type: `TypeBool`
  - Optional: `true`
  - Computed: `true`
  - Description: "Whether to enable deletion protection for the instance"
- [x] 1.3 Add `enable_risk_warning` parameter to resource schema in `tencentcloud/services/trabbit/resource_tc_tdmq_rabbitmq_vip_instance.go`
  - Type: `TypeBool`
  - Optional: `true`
  - Computed: `true`
  - Description: "Whether to enable cluster risk warning for the instance"

## 2. Update Function Implementation

- [x] 2.1 Add `remark` update logic to `resourceTencentCloudTdmqRabbitmqVipInstanceUpdate` function
  - Check if `remark` has changed using `d.HasChange("remark")`
  - If changed, set `request.Remark` with the new value using `helper.String()`
  - Set `needUpdate = true`
- [x] 2.2 Add `enable_deletion_protection` update logic to `resourceTencentCloudTdmqRabbitmqVipInstanceUpdate` function
  - Check if `enable_deletion_protection` has changed using `d.HasChange("enable_deletion_protection")`
  - If changed, set `request.EnableDeletionProtection` with the new value using `helper.Bool()`
  - Set `needUpdate = true`
- [x] 2.3 Add `enable_risk_warning` update logic to `resourceTencentCloudTdmqRabbitmqVipInstanceUpdate` function
  - Check if `enable_risk_warning` has changed using `d.HasChange("enable_risk_warning")`
  - If changed, set `request.EnableRiskWarning` with the new value using `helper.Bool()`
  - Set `needUpdate = true`

## 3. Read Function Implementation

- [x] 3.1 Add `remark` read logic to `resourceTencentCloudTdmqRabbitmqVipInstanceRead` function
  - Read `rabbitmqVipInstance.ClusterInfo.Remark` from API response
  - Set `remark` in resource state using `d.Set("remark", ...)`
- [x] 3.2 Add `enable_deletion_protection` read logic to `resourceTencentCloudTdmqRabbitmqVipInstanceRead` function
  - Read `rabbitmqVipInstance.ClusterInfo.EnableDeletionProtection` from API response
  - Set `enable_deletion_protection` in resource state using `d.Set("enable_deletion_protection", ...)`
- [x] 3.3 Add `enable_risk_warning` read logic to `resourceTencentCloudTdmqRabbitmqVipInstanceRead` function
  - Read `rabbitmqVipInstance.ClusterInfo.EnableRiskWarning` from API response
  - Set `enable_risk_warning` in resource state using `d.Set("enable_risk_warning", ...)`

## 4. Unit Tests

- [x] 4.1 Add unit test for updating `remark` parameter in `tencentcloud/services/trabbit/resource_tc_tdmq_rabbitmq_vip_instance_test.go`
  - Mock `ModifyRabbitMQVipInstance` API to verify `Remark` field is set
  - Verify the update is called only when `remark` changes
  - Use mock cloud API, do not call real cloud API
- [x] 4.2 Add unit test for updating `enable_deletion_protection` parameter in `tencentcloud/services/trabbit/resource_tc_tdmq_rabbitmq_vip_instance_test.go`
  - Mock `ModifyRabbitMQVipInstance` API to verify `EnableDeletionProtection` field is set
  - Verify both enabling and disabling scenarios
  - Use mock cloud API, do not call real cloud API
- [x] 4.3 Add unit test for updating `enable_risk_warning` parameter in `tencentcloud/services/trabbit/resource_tc_tdmq_rabbitmq_vip_instance_test.go`
  - Mock `ModifyRabbitMQVipInstance` API to verify `EnableRiskWarning` field is set
  - Verify both enabling and disabling scenarios
  - Use mock cloud API, do not call real cloud API
- [x] 4.4 Add unit test for reading new parameters in `tencentcloud/services/trabbit/resource_tc_tdmq_rabbitmq_vip_instance_test.go`
  - Mock `DescribeTdmqRabbitmqVipInstanceById` API response with all three new parameters set
  - Verify that the Read function correctly sets all three parameters in resource state
  - Use mock cloud API, do not call real cloud API
- [x] 4.5 Add unit test for backward compatibility in `tencentcloud/services/trabbit/resource_tc_tdmq_rabbitmq_vip_instance_test.go`
  - Test updating an instance without setting the new parameters
  - Verify that existing configurations continue to work without errors
  - Use mock cloud API, do not call real cloud API

## 5. Documentation

- [x] 5.1 Update resource example file `tencentcloud/services/trabbit/resource_tc_tdmq_rabbitmq_vip_instance.md`
  - Add example usage of `remark` parameter with sample value
  - Add example usage of `enable_deletion_protection` parameter with sample value
  - Add example usage of `enable_risk_warning` parameter with sample value
  - Update argument reference section with descriptions of new parameters
