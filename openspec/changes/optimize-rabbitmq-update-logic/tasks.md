## 1. Schema Updates

- [x] 1.1 Add `remark` field to resource schema
- [x] 1.2 Add `enable_deletion_protection` field to resource schema
- [x] 1.3 Add `enable_risk_warning` field to resource schema

## 2. Update Function Implementation

- [x] 2.1 Add `remark` field change detection and API call in Update function
- [x] 2.2 Add `enable_deletion_protection` field change detection and API call in Update function
- [x] 2.3 Add `enable_risk_warning` field change detection and API call in Update function

## 3. Read Function Implementation

- [x] 3.1 Read `remark` field from `DescribeRabbitMQVipInstance` API response
- [x] 3.2 Read `enable_deletion_protection` field from `DescribeRabbitMQVipInstances` API response
- [x] 3.3 Read `enable_risk_warning` field from `DescribeRabbitMQVipInstances` API response

## 4. Code Formatting

- [x] 4.1 Run `go fmt` on resource_tc_tdmq_rabbitmq_vip_instance.go
- [x] 4.2 Verify no formatting errors exist

## 5. Documentation Updates (Optional)

- [x] 5.1 Update resource_tc_tdmq_rabbitmq_vip_instance.md example file with new fields
- [x] 5.2 Run `make doc` to regenerate website documentation

## 6. Verification

- [x] 6.1 Run unit tests for RabbitMQ resource
- [x] 6.2 Run acceptance tests for RabbitMQ resource (if environment available)
- [x] 6.3 Verify backward compatibility with existing resources
