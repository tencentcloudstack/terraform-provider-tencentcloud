Provides a resource to create a Private Dns inbound endpoint

Example Usage

```hcl
resource "tencentcloud_private_dns_inbound_endpoint" "example" {
  endpoint_name   = "tf-example"
  endpoint_region = "ap-guangzhou"
  endpoint_vpc    = "vpc-i5yyodl9"
  subnet_ip {
    subnet_id  = "subnet-hhi88a58"
    subnet_vip = "10.0.30.2"
  }

  subnet_ip {
    subnet_id  = "subnet-5rrirqyc"
    subnet_vip = "10.0.0.11"
  }

  subnet_ip {
    subnet_id  = "subnet-60ut6n10"
  }
}
```
