Provides a resource to create a bi datasource_cloud

Example Usage

```hcl
resource "tencentcloud_bi_datasource_cloud" "datasource_cloud" {
  charset    = "utf8"
  db_name    = "bi_dev"
  db_type    = "MYSQL"
  db_user    = "root"
  project_id = "11015056"
  db_pwd     = "xxxxxx"
  service_type {
    instance_id = "cdb-12viotu5"
    region     = "ap-guangzhou"
    type       = "Cloud"
  }
  source_name = "tf-test1"
  vip         = "10.0.0.4"
  vport       = "3306"
  region_id   = "gz"
  vpc_id      = 5292713
}
```