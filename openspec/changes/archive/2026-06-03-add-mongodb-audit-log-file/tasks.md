## 1. Resource Implementation

- [x] 1.1 Create resource file `tencentcloud/services/mongodb/resource_tc_mongodb_audit_log_file.go` with schema definition, CRUD functions (Create, Read, Update, Delete), including retry logic, composite ID handling, immutableArgs pattern in Update, and proper error checking
- [x] 1.2 Register the resource `tencentcloud_mongodb_audit_log_file` in `tencentcloud/provider.go`
- [x] 1.3 Add the resource entry to `tencentcloud/provider.md`

## 2. Testing

- [x] 2.1 Create unit test file `tencentcloud/services/mongodb/resource_tc_mongodb_audit_log_file_test.go` with gomonkey-based mock tests covering Create, Read, Update (immutable error), and Delete scenarios
- [x] 2.2 Run unit tests with `go test -gcflags=all=-l` to verify all tests pass

## 3. Documentation

- [x] 3.1 Create resource example documentation file `tencentcloud/services/mongodb/resource_tc_mongodb_audit_log_file.md` with Example Usage and Import sections
