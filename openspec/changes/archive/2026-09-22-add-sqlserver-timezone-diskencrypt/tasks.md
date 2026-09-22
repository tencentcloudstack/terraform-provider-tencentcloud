# Tasks: Add TimeZone and DiskEncryptFlag Support to SqlServer Basic Instance

**Change ID**: `add-sqlserver-timezone-diskencrypt`

---

## Phase 1: Schema Definition

- [x] Task 1.1: Add `time_zone` schema field (Optional, Computed, ForceNew, TypeString) to `resource_tc_sqlserver_basic_instance.go`
- [x] Task 1.2: Add `disk_encrypt_flag` schema field (Optional, Computed, ForceNew, TypeInt, ValidateFunc 0-1) to `resource_tc_sqlserver_basic_instance.go`
- [x] Task 1.3: Add both fields to `immutableArgs` list in the update function
- [x] Task 1.4: Verify schema compiles correctly

## Phase 2: Service Layer

- [x] Task 2.1: Create `DescribeSqlserverInstanceAttributeById` method in `service_tencentcloud_sqlserver.go`
- [x] Task 2.2: Implement API request construction with InstanceId, ratelimit.Check, error handling
- [x] Task 2.3: Extract response and return `DescribeDBInstancesAttributeResponseParams` with nil-safe check
- [x] Task 2.4: Modify `CreateSqlserverBasicInstance` to read `time_zone` and `disk_encrypt_flag` from paramMap
- [x] Task 2.5: Pass `TimeZone` and `DiskEncryptFlag` to `CreateBasicDBInstancesRequest`

## Phase 3: Create Operation

- [x] Task 3.1: Read `time_zone` from schema using `GetOk` in Create function, add to paramMap
- [x] Task 3.2: Read `disk_encrypt_flag` from schema using `GetOkExists` in Create function, add to paramMap
- [x] Task 3.3: Verify paramMap passed to service layer `CreateSqlserverBasicInstance`

## Phase 4: Read Operation

- [x] Task 4.1: Read `TimeZone` from existing `DescribeDBInstances` response with nil check
- [x] Task 4.2: Call `DescribeSqlserverInstanceAttributeById` with retry logic in Read function
- [x] Task 4.3: Read `IsDiskEncryptFlag` from attribute response with double nil check, convert *int64 to int
- [x] Task 4.4: Graceful error handling - warn log on attribute API failure without failing whole read
- [x] Task 4.5: Verify all nil pointer checks in place

## Phase 5: Testing

- [x] Task 5.1: Add acceptance test `TestAccTencentCloudSqlserverBasicInstanceTimeZone` with custom timezone
- [x] Task 5.2: Add acceptance test `TestAccTencentCloudSqlserverBasicInstanceDiskEncrypt` with encryption enabled
- [x] Task 5.3: Add acceptance test `TestAccTencentCloudSqlserverBasicInstanceTimeZoneAndEncrypt` with both parameters

## Phase 6: Documentation

- [x] Task 6.1: Add usage examples with timezone and disk encryption to `resource_tc_sqlserver_basic_instance.md`
- [x] Task 6.2: Update website documentation via `make doc`

## Phase 7: Validation

- [x] Task 7.1: Verify provider registration (already exists, resource only adds params)
- [x] Task 7.2: Verify vendor SDK fields match implementation (TimeZone, DiskEncryptFlag, IsDiskEncryptFlag)
- [x] Task 7.3: Final review - all schema fields, create, read, immutableArgs, tests, docs complete