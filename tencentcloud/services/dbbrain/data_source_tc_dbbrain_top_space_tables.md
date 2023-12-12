Use this data source to query detailed information of dbbrain top_space_tables

Example Usage

Sort by PhysicalFileSize
```hcl
data "tencentcloud_dbbrain_top_space_tables" "top_space_tables" {
  instance_id = "%s"
  sort_by = "PhysicalFileSize"
  product = "mysql"
}
```

Sort by TotalLength
```hcl
data "tencentcloud_dbbrain_top_space_tables" "top_space_tables" {
  instance_id = "%s"
  sort_by = "PhysicalFileSize"
  product = "mysql"
}
```