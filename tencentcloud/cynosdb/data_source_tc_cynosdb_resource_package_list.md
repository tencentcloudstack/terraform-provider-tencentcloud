Use this data source to query detailed information of cynosdb resource_package_list

Example Usage

```hcl
data "tencentcloud_cynosdb_resource_package_list" "resource_package_list" {
  package_id      = ["package-hy4d2ppl"]
  package_name    = ["keep-package-disk"]
  package_type    = ["DISK"]
  package_region  = ["china"]
  status          = ["using"]
  order_by        = ["startTime"]
  order_direction = "DESC"
}
```