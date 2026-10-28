## 1. Data Source Implementation

- [x] 1.1 Create `tencentcloud/services/tco/data_source_tc_organization_permission_policies_in_role_configuration.go` with schema definition and Read function that calls `ListPermissionPoliciesInRoleConfiguration` API with retry logic
- [x] 1.2 Register the data source `tencentcloud_organization_permission_policies_in_role_configuration` in `tencentcloud/provider.go`

## 2. Documentation

- [x] 2.1 Create `tencentcloud/services/tco/data_source_tc_organization_permission_policies_in_role_configuration.md` with example usage
- [x] 2.2 Add the data source entry to `tencentcloud/provider.md`

## 3. Testing

- [x] 3.1 Create `tencentcloud/services/tco/data_source_tc_organization_permission_policies_in_role_configuration_test.go` with unit tests using gomonkey to mock the API call and verify the Read function logic
