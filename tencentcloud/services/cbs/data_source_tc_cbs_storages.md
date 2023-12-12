Use this data source to query detailed information of CBS storages.

Example Usage

```hcl
data "tencentcloud_cbs_storages" "storages" {
  storage_id         = "disk-kdt0sq6m"
  result_output_file = "mytestpath"
}
```

The following snippet shows the new supported query params

```hcl
data "tencentcloud_cbs_storages" "whats_new" {
  charge_type = ["POSTPAID_BY_HOUR", "PREPAID"]
  portable = true
  storage_state = ["ATTACHED"]
  instance_ips = ["10.0.0.2"]
  instance_name = ["my-instance"]
  tag_keys = ["foo"]
  tag_values = ["bar", "baz"]
}
```