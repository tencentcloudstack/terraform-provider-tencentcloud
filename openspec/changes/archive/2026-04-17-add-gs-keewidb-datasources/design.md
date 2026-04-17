# Design: GS Android Instances and KeeWiDB Instances Data Sources

## File Layout

| File | Action |
|---|---|
| `tencentcloud/connectivity/client.go` | Modified — add `UseGsV20190118Client` and `UseKeewidbV20220308Client` |
| `tencentcloud/services/gs/data_source_tc_gs_android_instances.go` | New |
| `tencentcloud/services/gs/service_tencentcloud_gs.go` | New |
| `tencentcloud/services/gs/data_source_tc_gs_android_instances.md` | New |
| `tencentcloud/services/gs/data_source_tc_gs_android_instances_test.go` | New |
| `tencentcloud/services/keewidb/data_source_tc_keewidb_instances.go` | New |
| `tencentcloud/services/keewidb/service_tencentcloud_keewidb.go` | New |
| `tencentcloud/services/keewidb/data_source_tc_keewidb_instances.md` | New |
| `tencentcloud/services/keewidb/data_source_tc_keewidb_instances_test.go` | New |
| `tencentcloud/provider.go` | Modified — register both data sources |
| `go.mod` | Modified — add gs and keewidb SDK deps |
| `vendor/` | Modified — `go mod vendor` to pull in new SDK packages |

## SDK

- GS: `gsv20190118 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/gs/v20190118"`
- KeeWiDB: `keewidbv20220308 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/keewidb/v20220308"`

## Code Style

Strictly follow `data_source_tc_igtm_instance_list.go`:
- `paramMap` built from d.GetOk
- Service method receives `paramMap`
- For-loop pagination with `resource.Retry` wrapping each API call inside the loop
- Flatten results into list of maps, `d.Set`

## Pagination

### GS DescribeAndroidInstances
- `Limit = 100` (max per docs)
- `Offset` increments by 100 each iteration
- Stop when `len(page) < 100` or `total reached`

### KeeWiDB DescribeInstances
- `Limit = 1000` (max per docs)
- `Offset` increments by 1000 each iteration
- Stop when `len(page) < 1000` or `total reached`

## Schema: tencentcloud_gs_android_instances

**Optional filters:**
- `android_instance_ids` (List of String)
- `android_instance_region` (String)
- `android_instance_zone` (String)
- `android_instance_group_ids` (List of String)

**Computed:** `android_instance_list` (List)

Fields per item: `android_instance_id`, `android_instance_region`, `android_instance_zone`, `state`, `android_instance_type`, `android_instance_image_id`, `width`, `height`, `host_serial_number`, `android_instance_group_id`, `name`, `user_id`, `private_ip`, `create_time`, `host_server_serial_number`, `service_status`, `android_instance_model`

## Schema: tencentcloud_keewidb_instances

**Optional filters:**
- `instance_id` (String)
- `instance_name` (String)
- `search_key` (String)
- `uniq_vpc_ids` (List of String)
- `uniq_subnet_ids` (List of String)
- `project_ids` (List of Int)
- `status` (List of Int)
- `billing_mode` (String)

**Computed:** `instance_list` (List)

Key fields per item: `instance_id`, `instance_name`, `status`, `region_id`, `zone_id`, `uniq_vpc_id`, `uniq_subnet_id`, `wan_ip`, `port`, `createtime`, `size`, `type`, `auto_renew_flag`, `deadline_time`, `engine`, `product_type`, `billing_mode`, `project_id`, `project_name`, `no_auth`, `disk_size`, `region`, `machine_memory`, `compression`
