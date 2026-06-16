Provides a resource to create a TDMQ professional cluster instance

Example Usage

```hcl
resource "tencentcloud_tdmq_pro_instance" "example" {
  zone_ids        = [200002, 200003, 200004]
  product_name    = "PULSAR.P1.MINI2"
  storage_size    = 200
  auto_renew_flag = 1
  time_span       = 1
  cluster_name    = "tf-example-pro-instance"
  auto_voucher    = 0

  vpc {
    vpc_id    = "vpc-xxxx"
    subnet_id = "subnet-xxxx"
  }
}
```

Import

tdmq pro_instance can be imported using the id (cluster_id), e.g.

```
terraform import tencentcloud_tdmq_pro_instance.example pulsar-xxxxxxxxxxxx
```
