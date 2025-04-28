Use this data source to query detailed information of CBS storages.

Example Usage

Query all CBS storages

```hcl
data "tencentcloud_cbs_storages" "example" {}
```

Query CBS by storage id

```hcl
data "tencentcloud_cbs_storages" "example" {
  storage_id         = "disk-6goq404g"
  result_output_file = "my-test-path"
}
```

Query CBS by dedicated cluster id

```hcl
data "tencentcloud_cbs_storages" "example" {
  dedicated_cluster_id = "cluster-262n63e8"
}
```

The following snippet shows the new supported query params

```hcl
data "tencentcloud_cbs_storages" "example" {
  charge_type   = ["POSTPAID_BY_HOUR", "PREPAID", "CDCPAID", "DEDICATED_CLUSTER_PAID"]
  storage_state = ["ATTACHED"]
  instance_ips  = ["10.0.0.2"]
  instance_name = ["my-instance"]
  tag_keys      = ["example"]
  tag_values    = ["bar", "baz"]
  portable      = true
}
```