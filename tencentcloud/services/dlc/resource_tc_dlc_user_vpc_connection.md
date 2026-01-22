Provides a resource to create a DLC user vpc connection

Example Usage

```hcl
resource "tencentcloud_dlc_user_vpc_connection" "example" {
  user_vpc_id            = "vpc-f7fa1fu5"
  user_subnet_id         = "subnet-ds2t3udw"
  user_vpc_endpoint_name = "tf-example"
  engine_network_id      = "DataEngine-Network-2mfg9icb"
}
```

Or

```hcl
resource "tencentcloud_dlc_user_vpc_connection" "example" {
  user_vpc_id            = "vpc-f7fa1fu5"
  user_subnet_id         = "subnet-ds2t3udw"
  user_vpc_endpoint_name = "tf-example"
  engine_network_id      = "DataEngine-Network-2mfg9icb"
  user_vpc_endpoint_vip  = "10.0.1.10"
}
```
