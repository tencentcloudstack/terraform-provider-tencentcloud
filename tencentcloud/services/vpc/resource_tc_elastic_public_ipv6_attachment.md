Provides a resource to create a vpc elastic_public_ipv6_attachment

Example Usage

```hcl
resource "tencentcloud_elastic_public_ipv6_attachment" "elastic_public_ipv6_attachment" {
  ipv6_address_id = "eipv6-xxxxxx"
  network_interface_id = "eni-xxxxxx"
  private_ipv6_address = "xxxxxx"
}
```

Import

vpc elastic_public_ipv6_attachment can be imported using the id, e.g.

```
terraform import tencentcloud_elastic_public_ipv6_attachment.elastic_public_ipv6_attachment elastic_public_ipv6_attachment_id
```
