Provides a resource to create a TEO bind security template

~> **NOTE:** If the domain name you input has been bound to a policy template (including site-level protection policies), the default value is to replace the template currently bound to the domain name.

~> **NOTE:** The current resource can only bind/unbind the template and domain name belonging to the same site.

Example Usage

```hcl
resource "tencentcloud_teo_bind_security_template" "example" {
  zone_id     = "zone-3skoch6ingbw"
  template_id = "temp-3s1pzyam2nxp"
  entity      = "tf.example.com"
  operate     = "unbind-use-default"
}

output "bind_message" {
  value = tencentcloud_teo_bind_security_template.example.message
}
```
Import

TEO bind security template can be imported using the zoneId#templateId#entity, e.g.

```
terraform import tencentcloud_teo_bind_security_template.teo_bind_security_template zone-3skoch6ingbw#temp-3s1pzyam2nxp#tf.example.com
```