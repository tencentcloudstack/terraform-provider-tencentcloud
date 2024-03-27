Provides a resource to create a vpc eni_ipv4_address

Example Usage

```hcl
data "tencentcloud_enis" "eni" {
  name = "Primary ENI"
}

resource "tencentcloud_eni_ipv4_address" "eni_ipv4_address" {
  network_interface_id = data.tencentcloud_enis.eni.enis.0.id
  secondary_private_ip_address_count = 3
}
```

Import

vpc eni_ipv4_address can be imported using the id, e.g.

```
terraform import tencentcloud_eni_ipv4_address.eni_ipv4_address eni_id
```