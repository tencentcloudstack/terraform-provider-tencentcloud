---
subcategory: "TencentCloud EdgeOne(TEO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_teo_bind_security_template"
sidebar_current: "docs-tencentcloud-resource-teo_bind_security_template"
description: |-
  Provides a resource to create a TEO bind security template
---

# tencentcloud_teo_bind_security_template

Provides a resource to create a TEO bind security template

~> **NOTE:** If the domain name you input has been bound to a policy template (including site-level protection policies), the default value is to replace the template currently bound to the domain name.

~> **NOTE:** The current resource can only bind/unbind the template and domain name belonging to the same site.

## Example Usage

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

## Argument Reference

The following arguments are supported:

* `entity` - (Required, String, ForceNew) The domain name to be bound to the policy template (or unbound from the policy template).
* `template_id` - (Required, String, ForceNew) Specifies the ID of the policy template or the site global policy to be bound or unbound.
<li>To bind to a policy template, or unbind from it, specify the policy template ID.</li>
<li>To bind to the site global policy, or unbind from it, use the `@ZoneLevel@domain` parameter value.</li>
Note: After unbinding, the domain name will use an independent policy and the rule quota will be calculated separately. Please make sure the plan rule quota is sufficient before unbinding.
* `zone_id` - (Required, String, ForceNew) The site ID to which the policy template to be bound or unbound belongs.
* `operate` - (Optional, String, ForceNew) Bind or unbind operation option. Valid values:
<li>`unbind-keep-policy`: unbind the domain name from the policy template while retaining the current policy.</li>
<li>`unbind-use-default`: unbind the domain name from the policy template and use the default blank policy.</li>
Default value: `unbind-keep-policy`. Note: The unbind operation currently only supports unbinding a single domain name. That is, when the value of `operate` is `unbind-keep-policy` or `unbind-use-default`, only one domain name can be unbound.
* `over_write` - (Optional, Bool, ForceNew) If the passed-in domain name is already bound to a policy template (including site-level protection policies), this parameter indicates whether to replace the template. Default value is `true`. Supported values:
<li>`true`: replace the template currently bound to the domain name.</li>
<li>`false`: do not replace the template currently bound to the domain name.</li>
Note: When set to `false`, if the passed-in domain name is already bound to a policy template, the API will return an error; the site-level protection policy is also a type of policy template.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `message` - Instance configuration delivery message. This field provides human-readable information about the delivery status (e.g., failure reasons).
* `status` - Instance configuration delivery status, the possible values are: `online`: the configuration has taken effect; `fail`: the configuration failed; `process`: the configuration is being delivered.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to `10m`) Used when creating the resource.
* `delete` - (Defaults to `10m`) Used when deleting the resource.

## Import

TEO bind security template can be imported using the zoneId#templateId#entity, e.g.

```
terraform import tencentcloud_teo_bind_security_template.teo_bind_security_template zone-3skoch6ingbw#temp-3s1pzyam2nxp#tf.example.com
```

