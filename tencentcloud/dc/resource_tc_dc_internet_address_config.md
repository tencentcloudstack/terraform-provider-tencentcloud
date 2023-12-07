Provides a resource to create a dc internet_address_config

Example Usage

```hcl
resource "tencentcloud_dc_internet_address" "internet_address" {
  mask_len = 30
  addr_type = 2
  addr_proto = 0
}

resource "tencentcloud_dc_internet_address_config" "internet_address_config" {
  instance_id = tencentcloud_dc_internet_address.internet_address.id
  enable = false
}
```

Import

dc internet_address_config can be imported using the id, e.g.

```
terraform import tencentcloud_dc_internet_address_config.internet_address_config internet_address_id
```