Provides a resource to create a CynosDB (TDSQL-C) LibraDB read-only analytics engine instance attachment

Example Usage

```hcl
resource "tencentcloud_cynosdb_libra_db_instance_attachment" "example" {
  cluster_id       = "cynosdbmysql-xxxxxxxx"
  zone             = "ap-guangzhou-3"
  cpu              = 4
  mem              = 8
  storage_size     = 100
  pay_mode         = 0
  instance_name    = "tf-example"
  instance_type    = "Common"
  storage_type     = "CLOUD_SSD"
  vpc_id           = "vpc-xxxxxxxx"
  subnet_id        = "subnet-xxxxxxxx"
  libra_db_version = "3.1.2"
  src_instance_id  = "cynosdbmysql-ins-xxxxxxxx"
}
```

Import

CynosDB LibraDB instance attachment can be imported using the cluster_id#instance_id, e.g.

```
terraform import tencentcloud_cynosdb_libra_db_instance_attachment.example cynosdbmysql-xxxxxxxx#cynosdbmysql-ins-yyyyyyyy
```
