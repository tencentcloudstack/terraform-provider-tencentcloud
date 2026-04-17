## 0. SDK Dependencies

- [x] 0.1 Add `go get github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/gs/...` and `go get github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/keewidb/...`, update `go.mod` and run `go mod vendor`
- [x] 0.2 Add `UseGsV20191118Client()` to `tencentcloud/connectivity/client.go`
- [x] 0.3 Add `UseKeewidbV20220308Client()` to `tencentcloud/connectivity/client.go`

## 1. GS Android Instances Data Source

- [x] 1.1 Create `tencentcloud/services/gs/service_tencentcloud_gs.go` with `DescribeGsAndroidInstancesByFilter()` (pagination Limit=100)
- [x] 1.2 Create `tencentcloud/services/gs/data_source_tc_gs_android_instances.go` with schema and read handler
- [x] 1.3 Create `tencentcloud/services/gs/data_source_tc_gs_android_instances.md`
- [x] 1.4 Create `tencentcloud/services/gs/data_source_tc_gs_android_instances_test.go`

## 2. KeeWiDB Instances Data Source

- [x] 2.1 Create `tencentcloud/services/keewidb/service_tencentcloud_keewidb.go` with `DescribeKeewidbInstancesByFilter()` (pagination Limit=1000)
- [x] 2.2 Create `tencentcloud/services/keewidb/data_source_tc_keewidb_instances.go` with schema and read handler
- [x] 2.3 Create `tencentcloud/services/keewidb/data_source_tc_keewidb_instances.md`
- [x] 2.4 Create `tencentcloud/services/keewidb/data_source_tc_keewidb_instances_test.go`

## 3. Provider Registration

- [x] 3.1 Register `tencentcloud_gs_android_instances` in `provider.go`
- [x] 3.2 Register `tencentcloud_keewidb_instances` in `provider.go`

## 4. Refinements

- [x] 4.1 Add missing input fields to `tencentcloud_gs_android_instances`: `label_selector` (LabelRequirement list: key/operator/values) and `filters` (Filter list: name/values)
- [x] 4.2 Add missing output field `android_instance_labels` (list of key/value) to `android_instance_list`
- [x] 4.3 Update `service_tencentcloud_gs.go` to pass `LabelSelector` and `Filters` from paramMap to request

## 5. KeeWiDB Schema Refinements

- [x] 5.1 Add missing input fields to `data_source_tc_keewidb_instances.go`: `order_by`, `order_type`, `type`, `auto_renew`, `vpc_ids`, `subnet_ids`, `search_keys`, `tag_keys`, `tag_list` (skip internal params TypeList/MonitorVersion)
- [x] 5.2 Update `service_tencentcloud_keewidb.go` to pass new filter fields from paramMap to request
