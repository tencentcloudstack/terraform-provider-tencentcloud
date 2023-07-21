---
subcategory: "Simple Email Service(SES)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_ses_template"
sidebar_current: "docs-tencentcloud-resource-ses_template"
description: |-
  Provides a resource to create a ses template.
---

# tencentcloud_ses_template

Provides a resource to create a ses template.

## Example Usage

### Create a ses template instance

```hcl
resource "tencentcloud_ses_template" "example" {
  template_name = "tf_example_ses_temp" "
  template_content {
    text = " example for the ses template "
  }
}
```

## Argument Reference

The following arguments are supported:

* `template_content` - (Required, List) Sms Template Content.
* `template_name` - (Required, String) smsTemplateName, which must be required.

The `template_content` object supports the following:

* `html` - (Optional, String) Html code after base64.
* `text` - (Optional, String) Text content after base64.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

ses template can be imported using the id, e.g.
```
$ terraform import tencentcloud_ses_template.example template_id
```

