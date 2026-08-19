## 1. Resource Implementation

- [x] 1.1 Create resource file `tencentcloud/services/mongodb/resource_tc_mongodb_audit_service.go` with Schema definition (instance_id, log_expire_day, audit_all, rule_filters, computed fields) and CRUD functions (resourceTencentCloudMongodbAuditServiceCreate, Read, Update, Delete) including async polling for Create/Delete, retry logic with tccommon.ReadRetryTimeout, and Timeouts block
- [x] 1.2 Register `tencentcloud_mongodb_audit_service` resource in `tencentcloud/provider.go` and `tencentcloud/provider.md`

## 2. Documentation

- [x] 2.1 Create resource documentation file `tencentcloud/services/mongodb/resource_tc_mongodb_audit_service.md` with one-line description, Example Usage (full audit and rule-based audit), and Import section

## 3. Unit Tests

- [x] 3.1 Create unit test file `tencentcloud/services/mongodb/resource_tc_mongodb_audit_service_test.go` using gomonkey mock to test Create, Read, Update, Delete flows, and verify tests pass with `go test -gcflags=all=-l`
