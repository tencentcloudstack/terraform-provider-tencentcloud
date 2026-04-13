## 1. Schema Definition

- [x] 1.1 Add `remark` field to Resource Schema as Optional String with description
- [x] 1.2 Add `enable_deletion_protection` field to Resource Schema as Optional Bool with description
- [x] 1.3 Add `enable_risk_warning` field to Resource Schema as Optional Bool with description

## 2. Read Function Implementation

- [x] 2.1 Add remark field reading logic in resourceTencentCloudTdmqRabbitmqVipInstanceRead function
- [x] 2.2 Add enable_deletion_protection field reading logic in resourceTencentCloudTdmqRabbitmqVipInstanceRead function
- [x] 2.3 Add enable_risk_warning field reading logic in resourceTencentCloudTdmqRabbitmqVipInstanceRead function (Note: Not available in cloud API response, field is only for update operations)
- [x] 2.4 Ensure nil value handling for all new fields in read operation

## 3. Update Function Implementation

- [x] 3.1 Add remark update logic in resourceTencentCloudTdmqRabbitmqVipInstanceUpdate function with d.HasChange() check
- [x] 3.2 Add enable_deletion_protection update logic in resourceTencentCloudTdmqRabbitmqVipInstanceUpdate function with d.HasChange() check
- [x] 3.3 Add enable_risk_warning update logic in resourceTencentCloudTdmqRabbitmqVipInstanceUpdate function with d.HasChange() check
- [x] 3.4 Ensure proper API parameter mapping using helper.String() and helper.Bool() functions
- [x] 3.5 Set needUpdate flag when any of the new fields change

## 4. Unit Tests

- [x] 4.1 Create unit test for remark field update in resource_tc_tdmq_rabbitmq_vip_instance_test.go using mock cloud API
- [x] 4.2 Create unit test for enable_deletion_protection field update in resource_tc_tdmq_rabbitmq_vip_instance_test.go using mock cloud API
- [x] 4.3 Create unit test for enable_risk_warning field update in resource_tc_tdmq_rabbitmq_vip_instance_test.go using mock cloud API
- [x] 4.4 Create unit test for multiple field updates in resource_tc_tdmq_rabbitmq_vip_instance_test.go using mock cloud API
- [x] 4.5 Create unit test for read operations with new fields in resource_tc_tdmq_rabbitmq_vip_instance_test.go using mock cloud API

## 5. Documentation

- [x] 5.1 Update resource_tc_tdmq_rabbitmq_vip_instance.md example file to include usage examples for the new fields
- [x] 5.2 Run make doc to generate website/docs/ documentation files (this will be executed by tfpacer-finalize skill)
