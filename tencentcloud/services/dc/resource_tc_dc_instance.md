Provides a resource to create a DC instance

Example Usage

Create direct connect instance

```hcl
resource "tencentcloud_dc_instance" "example" {
  direct_connect_name = "tf-example"
  access_point_id     = "ap-shenzhen-b-ft"
  line_operator       = "In-houseWiring"
  port_type           = "10GBase-LR"
}
```

Or

```hcl
resource "tencentcloud_dc_instance" "example" {
  direct_connect_name     = "tf-example"
  access_point_id         = "ap-shenzhen-b-ft"
  line_operator           = "In-houseWiring"
  port_type               = "10GBase-LR"
  bandwidth               = 100
  vlan                    = 1
  customer_contact_number = "0"
  sign_law                = true
}
```

Import

DC instance can be imported using the id, e.g.

```
terraform import tencentcloud_dc_instance.example dc-ovxsm3u5
```