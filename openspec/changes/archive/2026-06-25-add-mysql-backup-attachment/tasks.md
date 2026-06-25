## 1. Resource Implementation

- [x] 1.1 Create `resource_tc_mysql_backup_attachment.go` in `tencentcloud/services/cdb/` with schema definition (instance_id, backup_method, backup_db_table_list, manual_backup_name, encryption_flag, backup_id) and CRD lifecycle functions
- [x] 1.2 Implement Create function: call `CreateBackup` API with retry, validate BackupId is not empty, set ID as `backup_id#instance_id`, and call Read function
- [x] 1.3 Implement Read function: call `DescribeBackups` API with instance_id, iterate Items to find matching backup_id, set fields or clear ID if not found
- [x] 1.4 Implement Delete function: parse composite ID, call `DeleteBackup` API with instance_id and backup_id
- [x] 1.5 Implement Import support via `schema.ImportStatePassthrough`

## 2. Provider Registration

- [x] 2.1 Register `tencentcloud_mysql_backup` resource in `tencentcloud/provider.go` under the CDB resources section

## 3. Documentation

- [x] 3.1 Create `resource_tc_mysql_backup_attachment.md` in `tencentcloud/services/cdb/` with usage examples, argument reference, and import instructions

## 4. Testing

- [x] 4.1 Create `resource_tc_mysql_backup_attachment_test.go` in `tencentcloud/services/cdb/` with unit tests using gomonkey to mock API calls
- [x] 4.2 Run unit tests with `go test -gcflags=all=-l` to verify all tests pass

## 5. Validation

- [ ] 5.1 Run `make doc` to generate website documentation
- [ ] 5.2 Run `gofmt` on all modified Go files
- [x] 5.3 Verify the code compiles without errors
