Provides a resource to create a eip address_transform

Example Usage

```hcl
resource "tencentcloud_eip_address_transform" "address_transform" {
  instance_id = ""
}
```

Import

eip address_transform can be imported using the id, e.g.

```
terraform import tencentcloud_eip_address_transform.address_transform address_transform_id
```