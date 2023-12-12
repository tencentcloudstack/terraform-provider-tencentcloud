Provides a resource to create a dc instance

Example Usage

```hcl
resource "tencentcloud_dc_instance" "instance" {
  access_point_id         = "ap-shenzhen-b-ft"
  bandwidth               = 10
  customer_contact_number = "0"
  direct_connect_name     = "terraform-for-test"
  line_operator           = "In-houseWiring"
  port_type               = "10GBase-LR"
  sign_law                = true
  vlan                    = -1
}
```

Import

dc instance can be imported using the id, e.g.

```
terraform import tencentcloud_dc_instance.instance dc_id
```