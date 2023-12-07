Provides a resource to create a dc internet_address

Example Usage

```hcl
resource "tencentcloud_dc_internet_address" "internet_address" {
  mask_len = 30
  addr_type = 2
  addr_proto = 0
}
```

Import

dc internet_address can be imported using the id, e.g.

```
terraform import tencentcloud_dc_internet_address.internet_address internet_address_id
```