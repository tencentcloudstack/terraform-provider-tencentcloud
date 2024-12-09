Provides a resource to create a emr serverless_hbase_instance

Example Usage

```hcl
resource "tencentcloud_serverless_hbase_instance" "serverless_hbase_instance" {
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

Import

emr serverless_hbase_instance can be imported using the id, e.g.

```
terraform import tencentcloud_serverless_hbase_instance.serverless_hbase_instance serverless_hbase_instance_id
```
