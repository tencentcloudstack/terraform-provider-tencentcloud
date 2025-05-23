Provides a resource to create a Private Dns extend end point

Example Usage

If access_type is CLB

```hcl
resource "tencentcloud_private_dns_extend_end_point" "example" {
  end_point_name   = "tf-example"
  end_point_region = "ap-jakarta"
  forward_ip {
    access_type       = "CLB"
    host              = "10.0.1.12"
    port              = 9000
    vpc_id            = "vpc-1v2i79fc"
  }
}
```

If access_type is CCN

```hcl
resource "tencentcloud_private_dns_extend_end_point" "example" {
  end_point_name   = "tf-example"
  end_point_region = "ap-jakarta"
  forward_ip {
    access_type       = "CCN"
    host              = "1.1.1.1"
    port              = 8080
    vpc_id            = "vpc-2qjckjg2"
    access_gateway_id = "ccn-eo13f8ub"
  }
}
```

Import

Private Dns extend end point can be imported using the id, e.g.

```
terraform import tencentcloud_private_dns_extend_end_point.example eid-960fb0ee9677
```
