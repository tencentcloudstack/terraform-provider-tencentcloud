## 1. Schema Definition

- [x] 1.1 Add `kms_key_id` (TypeString, Optional, ForceNew) to `tencentcloud_instance` resource schema in `resource_tc_instance.go`

## 2. Create Function (RunInstances)

- [x] 2.1 Add `kms_key_id` mapping to `SystemDisk.KmsKeyId` in the RunInstances request builder (around line 891)

## 3. Update Function (ResetInstance)

- [x] 3.1 Add `kms_key_id` mapping to `SystemDisk.KmsKeyId` in the ResetInstance request builder (around line 2089)

## 4. Documentation

- [x] 4.1 Update `resource_tc_instance.md` with `kms_key_id` usage example

## 5. Unit Tests

- [x] 5.1 Add unit test cases in `resource_tc_instance_test.go` for the new `kms_key_id` parameter using gomonkey mocks

## 6. Verification

- [x] 6.1 Run `go test -gcflags=all=-l` on the test file to verify unit tests pass
- [x] 6.2 Verify the code compiles by running `go build ./...`