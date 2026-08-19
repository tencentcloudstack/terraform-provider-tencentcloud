## 1. Schema Definition

- [x] 1.1 Add `is_ssc` parameter (TypeBool, Optional, ForceNew, Description: "Whether it is storage-compute separation.") to the resource schema in `resource_tc_cdwdoris_instance.go`

## 2. Create Function

- [x] 2.1 In `resourceTencentCloudCdwdorisInstanceCreate`, add code to read `is_ssc` from the resource data and set `request.IsSSC` using `helper.Bool()` when the value is present

## 3. Update Function

- [x] 3.1 Add `"is_ssc"` to the `immutableArgs` slice in `resourceTencentCloudCdwdorisInstanceUpdate` to prevent modification after creation

## 4. Tests

- [x] 4.1 Add unit test cases in `resource_tc_cdwdoris_instance_test.go` for the `is_ssc` parameter using gomonkey mock, covering: create with is_ssc=true, create without is_ssc, and update attempt (should error)

## 5. Documentation

- [x] 5.1 Update `resource_tc_cdwdoris_instance.md` to include `is_ssc` in the example usage section

## 6. Validation

- [x] 6.1 Verify the code compiles correctly by running `go vet` against the modified files
- [x] 6.2 Run unit tests with `go test -gcflags=all=-l` to ensure all test cases pass