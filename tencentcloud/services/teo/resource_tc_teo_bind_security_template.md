Provides a resource to create a teo bind_security_template

~> **NOTE:** If the domain name you input has been bound to a policy template (including site-level protection policies), the default value is to replace the template currently bound to the domain name.
~> **NOTE:** The current resource can only bind/unbind the template and domain name belonging to the same site.

Example Usage

```hcl
resource "tencentcloud_teo_bind_security_template" "teo_bind_security_template" {
}

```
Import

teo application_proxy_rule can be imported using the zoneId#templateId, e.g.
```
terraform import tencentcloud_teo_bind_security_template.teo_bind_security_template zone-2983wizgxqvm#templateId
```