## 1. Code Implementation

- [x] 1.1 Add HostHeader parameter handling in resourceTencentCloudTeoOriginGroupCreate function
- [x] 1.2 Verify the implementation matches the pattern used in resourceTencentCloudTeoOriginGroupUpdate function

## 2. Code Verification

- [x] 2.1 Run `go build` to ensure code compiles without errors
- [x] 2.2 Run `go fmt` to ensure code formatting is correct
- [x] 2.3 Run `golint` to check for linting issues (if applicable)

## 3. Testing

- [x] 3.1 Run acceptance tests for tencentcloud_teo_origin_group resource with TF_ACC=1
- [x] 3.2 Verify that creating an origin group with host_header parameter works correctly
- [x] 3.3 Verify that creating an origin group without host_header parameter still works

## 4. Documentation

- [x] 4.1 Verify that the resource example file (resource_tc_teo_origin_group.md) is up to date
- [x] 4.2 Run `make doc` to regenerate website documentation if needed
- [x] 4.3 Verify that the generated documentation includes the host_header parameter
