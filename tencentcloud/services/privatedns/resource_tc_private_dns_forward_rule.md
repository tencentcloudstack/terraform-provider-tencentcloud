Provides a resource to create a Private Dns forward rule

Example Usage

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

resource "tencentcloud_private_dns_forward_rule" "example" {
  rule_name    = "tf-example"
  rule_type    = "DOWN"
  zone_id      = "zone-cmmbvaq8"
  end_point_id = tencentcloud_private_dns_extend_end_point.example.id
}
```

Import

Private Dns forward rule can be imported using the id, e.g.

```
terraform import tencentcloud_private_dns_forward_rule.example fid-dbc2c0a97c
```
