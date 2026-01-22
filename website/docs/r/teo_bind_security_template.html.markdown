---
subcategory: "TencentCloud EdgeOne(TEO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_teo_bind_security_template"
sidebar_current: "docs-tencentcloud-resource-teo_bind_security_template"
description: |-
  Provides a resource to create a teo bind_security_template
---

# tencentcloud_teo_bind_security_template

Provides a resource to create a teo bind_security_template

~> **NOTE:** If the domain name you input has been bound to a policy template (including site-level protection policies), the default value is to replace the template currently bound to the domain name.
~> **NOTE:** The current resource can only bind/unbind the template and domain name belonging to the same site.

## Example Usage

```hcl
resource "tencentcloud_teo_bind_security_template" "teo_bind_security_template" {
  operate     = "unbind-use-default"
  template_id = "temp-7dr7dm78"
  zone_id     = "zone-39quuimqg8r6"
  entity      = "aaa.makn.cn"
}
```

## Argument Reference

The following arguments are supported:

* `entity` - (Required, String, ForceNew) List of domain names to bind to/unbind from a policy template.
* `template_id` - (Required, String, ForceNew) Specifies the ID of the policy template or the site global policy to be bound or unbound.
<li>To bind to a policy template, or unbind from it, specify the policy template ID.</li>.
<li>To bind to the site's global policy, or unbind from it, use the @ZoneLevel@domain parameter value.</li>.

Note: After unbinding, the domain name will use an independent policy and rule quota will be calculated separately. Please make sure there is sufficient rule quota before unbinding.
* `zone_id` - (Required, String, ForceNew) Site ID of the policy template to be bound to or unbound from.
* `operate` - (Optional, String, ForceNew) Unbind operation option. valid values: `unbind-keep-policy`: unbind a domain name from the policy template while retaining the current policy. `unbind-use-default`: unbind a domain name from the policy template and use the default blank policy. default value: `unbind-keep-policy`.
* `over_write` - (Optional, Bool, ForceNew) If the passed-in domain is already bound to a policy template (including site-level protection policies), setting this parameter indicates whether to replace that template. The default value is true. Supported values are: `true`: Replace the currently bound template for the domain. `false`: Do not replace the currently bound template for the domain. Note: When set to false, if the passed-in domain is already bound to a policy template, the API will return an error; site-level protection policies are also a type of policy template.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `status` - Instance configuration delivery status, the possible values are: `online`: the configuration has taken effect; `fail`: the configuration failed; `process`: the configuration is being delivered.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to `3m`) Used when creating the resource.

## Import

teo application_proxy_rule can be imported using the zoneId#templateId#entity, e.g.
```
terraform import tencentcloud_teo_bind_security_template.teo_bind_security_template zone-39quuimqg8r6#temp-7dr7dm78#aaa.makn.cn
```

