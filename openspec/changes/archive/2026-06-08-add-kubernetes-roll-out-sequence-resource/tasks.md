## 1. Resource Implementation

- [x] 1.1 Create resource file `tencentcloud/services/tke/resource_tc_kubernetes_roll_out_sequence.go` with schema definition (name, sequence_flows, enabled), CRUD functions (Create, Read, Update, Delete), retry logic using `tccommon.ReadRetryTimeout`, and proper error handling
- [x] 1.2 Register `tencentcloud_kubernetes_roll_out_sequence` resource in `tencentcloud/provider.go` resource map
- [x] 1.3 Add `tencentcloud_kubernetes_roll_out_sequence` entry in `tencentcloud/provider.md`

## 2. Documentation

- [x] 2.1 Create example usage documentation file `tencentcloud/services/tke/resource_tc_kubernetes_roll_out_sequence.md` with Example Usage and Import sections

## 3. Unit Tests

- [x] 3.1 Create unit test file `tencentcloud/services/tke/resource_tc_kubernetes_roll_out_sequence_test.go` with gomonkey-mocked tests covering Create, Read, Update, and Delete operations
- [x] 3.2 Run unit tests with `go test -gcflags=all=-l` to verify they pass
