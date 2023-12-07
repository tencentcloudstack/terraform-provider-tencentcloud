Provides a resource to create a HA VIP EIP attachment.

Example Usage

```hcl
resource "tencentcloud_ha_vip_eip_attachment" "foo" {
  havip_id   = "havip-kjqwe4ba"
  address_ip = "1.1.1.1"
}
```

Import

HA VIP EIP attachment can be imported using the id, e.g.

```
$ terraform import tencentcloud_ha_vip_eip_attachment.foo havip-kjqwe4ba#1.1.1.1
```