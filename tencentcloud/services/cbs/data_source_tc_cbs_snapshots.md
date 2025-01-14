Use this data source to query detailed information of CBS snapshots.

Example Usage

Query all snapshots

```hcl
data "tencentcloud_cbs_snapshots" "snapshots" {}
```

Query snapshots by filters

```hcl
data "tencentcloud_cbs_snapshots" "snapshots" {
  snapshot_id        = "snap-hibh08s3"
  result_output_file = "my_snapshots"
}

data "tencentcloud_cbs_snapshots" "snapshots" {
  snapshot_name = "tf-example"
}

data "tencentcloud_cbs_snapshots" "snapshots" {
  storage_id = "disk-12j0fk1w"
}

data "tencentcloud_cbs_snapshots" "snapshots" {
  storage_usage = "SYSTEM_DISK"
}

data "tencentcloud_cbs_snapshots" "snapshots" {
  project_id = "0"
}

data "tencentcloud_cbs_snapshots" "snapshots" {
  availability_zone = "ap-guangzhou-4"
}
```