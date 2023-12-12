Provides a resource to create a vpc enable_end_point_connect

Example Usage

```hcl
resource "tencentcloud_vpc_enable_end_point_connect" "enable_end_point_connect" {
  end_point_service_id = "vpcsvc-98jddhcz"
  end_point_id         = ["vpce-6q0ftmke"]
  accept_flag          = true
}
```