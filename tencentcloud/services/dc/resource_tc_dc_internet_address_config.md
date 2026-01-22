Provides a resource to create a dc internet address config

Example Usage

```hcl
resource "tencentcloud_dc_internet_address" "example" {
  mask_len   = 30
  addr_type  = 2
  addr_proto = 0
}

resource "tencentcloud_dc_internet_address_config" "example" {
  instance_id = tencentcloud_dc_internet_address.example.id
  enable      = true
}
```

Import

dc internet address config can be imported using the id, e.g.

```
terraform import tencentcloud_dc_internet_address_config.example ipv4-5091pc5v
```