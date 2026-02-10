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
    access_type = "CCN"
    hosts = [
      "1.1.1.1:8080",
      "2.2.2.2:9090",
    ]
    port              = 8080
    vpc_id            = "vpc-h70u60bi"
    access_gateway_id = "ccn-4s3g3yg5"
  }
}
```

Import

Private Dns extend end point can be imported using the id, e.g.

```
terraform import tencentcloud_private_dns_extend_end_point.example eid-960fb0ee9677
```
