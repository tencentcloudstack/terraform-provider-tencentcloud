Use this data source to query detailed information of CDWDoris cluster instance nodes.

Example Usage

Query cdwdoris instance nodes by instance id

```hcl
data "tencentcloud_cdwdoris_instance_nodes" "example" {
  instance_id = "cdwdoris-rhbflamd"
}
```

Query cdwdoris instance nodes with optional filters

```hcl
data "tencentcloud_cdwdoris_instance_nodes" "example" {
  instance_id    = "cdwdoris-rhbflamd"
  node_role      = "data"
  display_policy = "All"
}
```