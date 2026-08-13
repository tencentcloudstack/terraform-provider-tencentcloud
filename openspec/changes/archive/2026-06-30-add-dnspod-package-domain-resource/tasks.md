## 1. Resource Implementation

- [x] 1.1 Create `tencentcloud/services/dnspod/resource_tc_dnspod_package_domain.go` with schema definition and CRUD functions (Create/Read/Update/Delete/Import)
- [x] 1.2 Create `tencentcloud/services/dnspod/resource_tc_dnspod_package_domain_test.go` with unit tests using gomonkey mock
- [x] 1.3 Create `tencentcloud/services/dnspod/resource_tc_dnspod_package_domain.md` with resource documentation

## 2. Provider Registration

- [x] 2.1 Register `tencentcloud_dnspod_package_domain` resource in `tencentcloud/provider.go`
- [x] 2.2 Add resource entry in `tencentcloud/provider.md`

## 3. Validation

- [x] 3.1 Run `go test -gcflags=all=-l` to execute unit tests and ensure all pass
- [x] 3.2 Run `openspec status --change "add-dnspod-package-domain-resource"` to confirm all artifacts complete
