Provides a resource to create a Private Dns end point

Example Usage

```hcl
resource "tencentcloud_private_dns_end_point" "example" {
  end_point_name       = "tf-example"
  end_point_service_id = "vpcsvc-61wcwmar"
  end_point_region     = "ap-guangzhou"
  ip_num               = 1
}
```

Import

Private Dns end point can be imported using the id, e.g.

```
terraform import tencentcloud_private_dns_end_point.example eid-77a246c867
```
