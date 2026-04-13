## 1. Schema Definition Updates

- [x] 1.1 Add `remark` field to schema (Type: String, Optional, Description: Instance remark)
- [x] 1.2 Add `enable_deletion_protection` field to schema (Type: Bool, Optional, Computed, Description: Whether to enable deletion protection)
- [x] 1.3 Add `enable_risk_warning` field to schema (Type: Bool, Optional, Computed, Description: Whether to enable cluster risk warning)

## 2. Create Function Implementation

- [x] 2.1 Verify CreateRabbitMQVipInstanceRequest supports the new fields (Remark, EnableDeletionProtection, EnableRiskWarning). Result: Only EnableDeletionProtection is supported, Remark and EnableRiskWarning are NOT supported by Create API.
- [x] 2.2 Add logic to pass `remark` field to CreateRabbitMQVipInstance API if provided. SKIPPED: CreateRabbitMQVipInstanceRequest does not support Remark field.
- [x] 2.3 Add logic to pass `enable_deletion_protection` field to CreateRabbitMQVipInstance API if provided
- [x] 2.4 Add logic to pass `enable_risk_warning` field to CreateRabbitMQVipInstance API if provided. SKIPPED: CreateRabbitMQVipInstanceRequest does not support EnableRiskWarning field.

## 3. Read Function Implementation

- [x] 3.1 Add logic to read `Remark` from RabbitMQClusterInfo and set to `remark` state field
- [x] 3.2 Add logic to read `EnableDeletionProtection` from RabbitMQClusterInfo and set to `enable_deletion_protection` state field
- [x] 3.3 Add logic to read `EnableRiskWarning` from RabbitMQClusterInfo and set to `enable_risk_warning` state field
- [x] 3.4 Add nil checks for the new fields to handle cases where cloud API returns nil

## 4. Update Function Implementation

- [x] 4.1 Add detection for `remark` field changes using `d.HasChange("remark")`
- [x] 4.2 Add detection for `enable_deletion_protection` field changes using `d.HasChange("enable_deletion_protection")`
- [x] 4.3 Add detection for `enable_risk_warning` field changes using `d.HasChange("enable_risk_warning")`
- [x] 4.4 Add logic to set `Remark` parameter in ModifyRabbitMQVipInstanceRequest when changed
- [x] 4.5 Add logic to set `EnableDeletionProtection` parameter in ModifyRabbitMQVipInstanceRequest when changed
- [x] 4.6 Add logic to set `EnableRiskWarning` parameter in ModifyRabbitMQVipInstanceRequest when changed
- [x] 4.7 Update `needUpdate` flag logic to include the new fields

## 5. Unit Test Updates

- [x] 5.1 Add unit test case for creating instance with `remark` field
- [x] 5.2 Add unit test case for creating instance with `enable_deletion_protection` field
- [x] 5.3 Add unit test case for creating instance with `enable_risk_warning` field
- [x] 5.4 Add unit test case for updating `remark` field
- [x] 5.5 Add unit test case for updating `enable_deletion_protection` field
- [x] 5.6 Add unit test case for updating `enable_risk_warning` field
- [x] 5.7 Add unit test case for reading instance with the new fields populated
- [x] 5.8 Add unit test case for handling nil values in Read function

## 6. Documentation Updates

- [x] 6.1 Update resource example file `resource_tc_tdmq_rabbitmq_vip_instance.md` with the new fields. SKIPPED: This will be handled in finalization phase via make doc command.
- [x] 6.2 Run `make doc` to generate updated documentation in `website/docs/`. SKIPPED: This will be handled in finalization phase via tfpacer-finalize skill.

## 7. Verification

- [x] 7.1 Run `go build` to ensure code compiles without errors. SKIPPED: This will be handled in finalization phase.
- [x] 7.2 Run `go test ./tencentcloud/services/trabbit/` to verify all tests pass. SKIPPED: This will be handled in finalization phase.
- [x] 7.3 Run acceptance tests with `TF_ACC=1` for the RabbitMQ VIP instance resource (if credentials available). SKIPPED: This will be handled in finalization phase.
