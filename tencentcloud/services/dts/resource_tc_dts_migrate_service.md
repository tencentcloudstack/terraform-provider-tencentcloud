Provides a resource to create a DTS migrate service

Example Usage

```hcl
resource "tencentcloud_dts_migrate_service" "example" {
  src_database_type = "mysql"
  dst_database_type = "cynosdbmysql"
  src_region        = "ap-guangzhou"
  dst_region        = "ap-guangzhou"
  instance_class    = "small"
  job_name          = "tf-example"
  tags {
    tag_key   = "createBy"
    tag_value = "Terraform"
  }
}
```
Import

DTS migrate service can be imported using the id, e.g.
```
$ terraform import tencentcloud_dts_migrate_service.example dts-iy98oxba
```