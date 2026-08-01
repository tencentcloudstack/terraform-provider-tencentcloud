Provides a resource to create a DBDC db custom node.

~> **NOTE:** Not all `zone` support configuration `system_disk` and `data_disks`. You can check the supported zone list with `tencentcloud_dbdc_db_custom_node_types`.

Example Usage

Create a PREPAID DBDC db custom node

```hcl
resource "tencentcloud_dbdc_db_custom_node" "example" {
  zone       = "ap-shanghai-5"
  image_id   = "img-rm13akp3"
  vpc_id     = "vpc-py7mlxqm"
  subnet_id  = "subnet-qd4upp83"
  node_type  = "DB.AT5.8XLARGE128"
  period     = 1
  auto_renew = 1
  node_name  = "tf-example"

  login_settings {
    password = "Password@2026"
  }

  tags = {
    createBy = "Terraform"
  }
}
```

Create a POSTPAID DBDC db custom node

```hcl
resource "tencentcloud_dbdc_db_custom_node" "example" {
  zone         = "ap-bangkok-2"
  image_id     = "img-rm13akp3"
  vpc_id       = "vpc-doprnsrq"
  subnet_id    = "subnet-kzp25bun"
  node_type    = "DB.SA5.8XLARGE128"
  node_name    = "tf-example"
  host_name    = "hostName"
  network_mode = "cross_tenant_eni"
  charge_type  = "POSTPAID"

  login_settings {
    password = "Password@2026"
  }

  system_disk {
    disk_size = 100
    disk_type = "CLOUD_HSSD"
  }

  data_disks {
    disk_size = 200
    disk_type = "CLOUD_HSSD"
  }

  security_group_ids = [
    "sg-avup6l04",
  ]

  tags = {
    createBy = "Terraform"
  }
}
```

Import

DBDC db custom node can be imported using the id, e.g.

```
terraform import tencentcloud_dbdc_db_custom_node.example dbcn-ttiyh58n
```
