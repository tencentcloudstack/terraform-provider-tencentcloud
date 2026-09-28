## 1. Schema Definition

- [x] 1.1 Add `custom_image` field (Type=String, Optional, no ForceNew) to the `native` block schema in `resource_tc_kubernetes_native_node_pool.go`, placed alongside the existing `machine_type` field, with a description explaining it specifies the custom image ID

## 2. Create Function Implementation

- [x] 2.1 In `resourceTencentCloudKubernetesNativeNodePoolCreate()`, read `custom_image` from the `nativeMap` and assign it to `createNativeNodePoolParam.CustomImage` via `helper.String()`, following the existing pattern used for `machine_type`

## 3. Read Function Implementation

- [x] 3.1 In `resourceTencentCloudKubernetesNativeNodePoolRead()`, when `respData.Native` is not nil, nil-check `respData.Native.CustomImage` and set `nativeMap["custom_image"]` when present, following the existing pattern used for `machine_type`

## 4. Update Function Implementation

- [x] 4.1 In `resourceTencentCloudKubernetesNativeNodePoolUpdate()`, read `custom_image` from the `nativeMap` and assign it to `updateNativeNodePoolParam.CustomImage` via `helper.String()`, following the existing pattern used for `key_ids`

## 5. Documentation

- [x] 5.1 Update `resource_tc_kubernetes_native_node_pool.md` to document the new `custom_image` parameter in the `native` block

## 6. Unit Test

- [x] 6.1 Add unit test cases in `resource_tc_kubernetes_native_node_pool_test.go` covering the `custom_image` parameter in Create, Read, and Update paths using gomonkey mocks

## 7. Verification

- [x] 7.1 Verify the generated code compiles and the new field is correctly wired to `Native.CustomImage` in Create/Update and flattened from `NativeNodePoolInfo.CustomImage` in Read
