Provides a resource to create a privatedns forward rule

Example Usage

```hcl
resource "tencentcloud_private_dns_end_point" "example" {
  end_point_name       = "tf-example"
  end_point_service_id = "vpcsvc-61wcwmar"
  end_point_region     = "ap-guangzhou"
  ip_num               = 1
}

resource "tencentcloud_private_dns_forward_rule" "example" {
  rule_name    = "tf-example"
  rule_type    = "DOWN"
  zone_id      = "zone-cmmbvaq8"
  end_point_id = tencentcloud_private_dns_end_point.example.id
}
```

Import

private dns forward rule can be imported using the id, e.g.

```
terraform import tencentcloud_private_dns_forward_rule.example fid-dbc2c0a97c
```
