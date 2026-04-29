---
subcategory: "TencentCloud EdgeOne(TEO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_teo_web_security_templates"
sidebar_current: "docs-tencentcloud-datasource-teo_web_security_templates"
description: |-
  Use this data source to query detailed information of TEO web security policy templates
---

# tencentcloud_teo_web_security_templates

Use this data source to query detailed information of TEO web security policy templates

## Example Usage

```hcl
data "tencentcloud_teo_web_security_templates" "templates" {
  zone_ids = ["zone-3fkff38fyw8s"]
}
```

## Argument Reference

The following arguments are supported:

* `zone_ids` - (Required, List: [`String`]) Zone ID list. Up to 100 zone IDs per query.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `security_policy_templates` - Security policy template list.
  * `bind_domains` - Domain binding information for the policy template.
    * `domain` - Domain name.
    * `status` - Binding status. Valid values: process (binding in progress), online (binding succeeded), fail (binding failed).
    * `zone_id` - Zone ID that the domain belongs to.
  * `template_id` - Policy template ID.
  * `template_name` - Policy template name.
  * `zone_id` - Zone ID that the policy template belongs to.


