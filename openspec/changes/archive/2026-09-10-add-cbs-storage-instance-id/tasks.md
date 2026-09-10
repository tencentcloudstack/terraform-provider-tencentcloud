## 1. Schema Definition

- [x] 1.1 Add `instance_id` field to the schema in `tencentcloud/services/cbs/resource_tc_cbs_storage.go` with Type=String, Optional=true, and a description explaining it specifies the CVM instance ID for auto-mounting the data disk during creation
- [x] 1.2 Verify the new field is placed in the Optional parameters section (alongside other optional fields like `snapshot_id`, `project_id`, etc.) and does not conflict with existing schema fields

## 2. Create Function Implementation

- [x] 2.1 In `resourceTencentCloudCbsStorageCreate()`, read `instance_id` from resource data
- [x] 2.2 If `instance_id` is set, populate `request.AutoMountConfiguration` with `&cbs.AutoMountConfiguration{InstanceId: []*string{helper.String(instanceId)}}` before the `CreateDisks` API call
- [x] 2.3 Ensure `AutoMountConfiguration` is only set when `instance_id` is provided (backward compatibility - do not set it for existing configurations)
- [x] 2.4 Verify the Create function follows existing patterns (retry logic, nil checks on response, error handling) and does not break the existing disk-id extraction logic

## 3. Read Function Implementation

- [x] 3.1 In `resourceTencentCloudCbsStorageRead()`, after the existing `storage` object is retrieved, check if `storage.InstanceId` is non-nil
- [x] 3.2 If non-nil, set `instance_id` in resource data via `d.Set("instance_id", storage.InstanceId)`
- [x] 3.3 Ensure the nil check guard is in place so that an empty/unmounted disk does not overwrite state with an empty string

## 4. Documentation

- [x] 4.1 Update `tencentcloud/services/cbs/resource_tc_cbs_storage.md` to add a new example showing CBS storage creation with `instance_id` for auto-mounting
- [x] 4.2 In the example, demonstrate `instance_id = "ins-xxxxxxxx"` alongside existing required fields
- [x] 4.3 Verify the .md file format follows the gendoc requirements (one-line description, Example Usage, Import section) and does not manually add Argument Reference / Attribute Reference sections

## 5. Testing

- [x] 5.1 Add a test case for `instance_id` in `tencentcloud/services/cbs/resource_tc_cbs_storage_test.go` using the existing Terraform acceptance test suite pattern (the resource already uses TF test suite, so extend it for the modified resource)
- [x] 5.2 The test should create a CBS storage with `instance_id` specified and verify the field is correctly set in state
- [x] 5.3 Verify backward compatibility: ensure existing test cases continue to pass without `instance_id`

## 6. Validation

- [x] 6.1 Manually review code changes to ensure compilation correctness (all function returns checked, no unused variables)
- [x] 6.2 Verify the new parameter maps correctly: `instance_id` → `AutoMountConfiguration.InstanceId` (Create) and `Disk.InstanceId` → `instance_id` (Read)
- [x] 6.3 Confirm only ONE new parameter (`instance_id`) was added and no other unrelated changes were made
- [x] 6.4 Verify backward compatibility: no existing schema fields were modified or removed
