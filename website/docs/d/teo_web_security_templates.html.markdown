---
subcategory: "TencentCloud EdgeOne(TEO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_teo_web_security_templates"
sidebar_current: "docs-tencentcloud-datasource-teo_web_security_templates"
description: |-
  Use this data source to query detailed information of TEO web security templates
---

# tencentcloud_teo_web_security_templates

Use this data source to query detailed information of TEO web security templates

## Example Usage

```hcl
data "tencentcloud_teo_web_security_templates" "example" {
  zone_ids = [
    "zone-3fkff38fyw8s",
  ]
}
```

## Argument Reference

The following arguments are supported:

* `zone_ids` - (Required, Set: [`String`]) List of zone IDs. A maximum of 100 zones can be queried in a single request.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `security_policy_templates` - List of policy templates.
  * `bind_domains` - Information about domains bound to the policy template.
    * `domain` - Domain name.
    * `status` - Binding status. valid values:. 
<li>`process`: binding in progress</li>
<li>`online`: binding succeeded.</li>
<Li>`fail`: binding failed.</li>.
    * `zone_id` - Zone ID to which the domain belongs.
  * `template_id` - Policy template ID.
  * `template_name` - The name of the policy template.
  * `zone_id` - The zone ID to which the policy template belongs.


