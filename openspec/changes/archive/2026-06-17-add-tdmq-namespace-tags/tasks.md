## 1. Service Layer Updates

- [x] 1.1 Update `CreateTdmqNamespace` function in `tencentcloud/services/tdmq/service_tencentcloud_tdmq.go` to accept a `tags []*tdmq.Tag` parameter and assign it to `request.Tags`
- [x] 1.2 Verify that the `DescribeTdmqNamespaceById` function returns the `Environment` struct which already includes the `Tags` field (no change needed, just confirm)

## 2. Resource Schema and CRUD Updates

- [x] 2.1 Add `Tags` schema field to `tencentcloud_tdmq_namespace` resource in `tencentcloud/services/tpulsar/resource_tc_tdmq_namespace.go` as `TypeList`, `Optional`, with `tag_key` (Required TypeString) and `tag_value` (Required TypeString) sub-fields
- [x] 2.2 Update `resourceTencentCloudTdmqNamespaceCreate` to read Tags from schema and pass to `CreateTdmqNamespace` service function
- [x] 2.3 Update `resourceTencentCloudTdmqNamespaceRead` to flatten `info.Tags` into the schema (check for nil before setting)
- [x] 2.4 Update `resourceTencentCloudTdmqNamespaceUpdate` to detect tag changes and use `svctag.ModifyTags` (via `TagResources`/`UnTagResources` APIs) to apply tag updates in-place
- [x] 2.5 Verify the `Delete` method does not need changes (Tags are destroyed with the namespace)

## 3. Documentation

- [x] 3.1 Update `tencentcloud/services/tpulsar/resource_tc_tdmq_namespace.md` to include Tags parameter in example usage

## 4. Unit Tests

- [x] 4.1 Add unit tests in `tencentcloud/services/tpulsar/resource_tc_tdmq_namespace_test.go` for the Tags parameter using gomonkey mock approach (mock the TDMQ service methods)
- [x] 4.2 Run `go test -gcflags=all=-l` on the test file to verify tests pass
