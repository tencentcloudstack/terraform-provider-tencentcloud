## 1. Schema Definition

- [x] 1.1 Add `remark` field to resource schema with TypeString, Optional: true, Computed: true, and appropriate description
- [x] 1.2 Add `enable_deletion_protection` field to resource schema with TypeBool, Optional: true, Computed: true, and appropriate description
- [x] 1.3 Add `enable_risk_warning` field to resource schema with TypeBool, Optional: true, Computed: true, and appropriate description

## 2. Create Logic Implementation

- [x] 2.1 Add `enable_deletion_protection` parameter to CreateRabbitMQVipInstance API request in resourceTencentCloudTdmqRabbitmqVipInstanceCreate function

## 3. Read Logic Implementation

- [x] 3.1 Read `remark` field from DescribeRabbitMQVipInstanceResponseParams and set to schema in resourceTencentCloudTdmqRabbitmqVipInstanceRead function
- [x] 3.2 Read `enable_deletion_protection` field from DescribeRabbitMQVipInstanceResponseParams and set to schema in resourceTencentCloudTdmqRabbitmqVipInstanceRead function
- [x] 3.3 Handle nil values for new fields gracefully in read operation

## 4. Update Logic Implementation

- [x] 4.1 Remove `remark`, `enable_deletion_protection`, `enable_risk_warning` from immutableArgs list in resourceTencentCloudTdmqRabbitmqVipInstanceUpdate function
- [x] 4.2 Add `remark` field update logic to resourceTencentCloudTdmqRabbitmqVipInstanceUpdate function using d.HasChange() and ModifyRabbitMQVipInstance API
- [x] 4.3 Add `enable_deletion_protection` field update logic to resourceTencentCloudTdmqRabbitmqVipInstanceUpdate function using d.HasChange() and ModifyRabbitMQVipInstance API
- [x] 4.4 Add `enable_risk_warning` field update logic to resourceTencentCloudTdmqRabbitmqVipInstanceUpdate function using d.HasChange() and ModifyRabbitMQVipInstance API

## 5. Unit Testing

- [x] 5.1 Add unit tests for new schema fields in resource_tc_tdmq_rabbitmq_vip_instance_test.go
- [x] 5.2 Add unit tests for Create API integration with new fields using mock cloud API
- [x] 5.3 Add unit tests for Read API integration with new fields using mock cloud API
- [x] 5.4 Add unit tests for Update API integration with new fields using mock cloud API

## 6. Documentation

- [x] 6.1 Run `make doc` to generate website documentation for updated resource schema
- [x] 6.2 Verify generated documentation includes `remark`, `enable_deletion_protection`, `enable_risk_warning` fields with correct descriptions

## 7. Code Quality

- [x] 7.1 Run `go fmt` to format all modified Go code files
- [x] 7.2 Verify no linting errors in modified code
