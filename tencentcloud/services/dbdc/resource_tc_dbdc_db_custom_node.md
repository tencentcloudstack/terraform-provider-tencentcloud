Provides a resource to create a DBDC db custom node.

Example Usage

```hcl
resource "tencentcloud_dbdc_db_custom_node" "example" {
  zone       = "ap-shanghai-5"
  image_id   = "img-rm13akp3"
  vpc_id     = "vpc-py7mlxqm"
  subnet_id  = "subnet-qd4upp83"
  node_type  = "DB.SA5.8XLARGE128"
  period     = 1
  auto_renew = 1
  node_name  = "tf-example"

  charge_type   = "PREPAID"
  network_mode  = "privatelink"
  host_name     = "tf-example-node"
  security_group_ids = ["sg-xxxxxxxx"]

  system_disk {
    disk_type = "CLOUD_HSSD"
    disk_size = 100
  }

  data_disks {
    disk_type = "CLOUD_HSSD"
    disk_size = 500
    disk_name = "tf-example-data"
  }

  login_settings {
    password = "Password@2026"
  }

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
