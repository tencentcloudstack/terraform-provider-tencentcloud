Provides a resource to create a emr lite_hbase_instance. The `disk_size` parameter supports online expansion (expansion only, not scaling in) by calling the ModifySLInstance API of the EMR service.

Example Usage

```hcl
resource "tencentcloud_lite_hbase_instance" "lite_hbase_instance" {
  instance_name = "tf-test"
  pay_mode = 0
  disk_type = "CLOUD_HSSD"
  disk_size = 100
  node_type = "8C32G"
  zone_settings {
    zone = "ap-shanghai-2"
    vpc_settings {
      vpc_id = "vpc-xxxxxx"
      subnet_id = "subnet-xxxxxx"
    }
    node_num = 3
  }
  tags {
    tag_key = "test"
    tag_value = "test"
  }
}
```

The `disk_size` can be expanded by updating the value, for example from 100 to 200:

```hcl
resource "tencentcloud_lite_hbase_instance" "lite_hbase_instance" {
  instance_name = "tf-test"
  pay_mode = 0
  disk_type = "CLOUD_HSSD"
  disk_size = 200
  node_type = "8C32G"
  zone_settings {
    zone = "ap-shanghai-2"
    vpc_settings {
      vpc_id = "vpc-xxxxxx"
      subnet_id = "subnet-xxxxxx"
    }
    node_num = 3
  }
  tags {
    tag_key = "test"
    tag_value = "test"
  }
}
```

Import

emr lite_hbase_instance can be imported using the id, e.g.

```
terraform import tencentcloud_lite_hbase_instance.lite_hbase_instance lite_hbase_instance_id
```
