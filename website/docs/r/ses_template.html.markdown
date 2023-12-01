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

### Create a ses text template

```hcl
resource "tencentcloud_ses_template" "example" {
  template_name = "tf_example_ses_temp" "
  template_content {
    text = " example for the ses template "
  }
}
```

### Create a ses html template

```hcl
resource "tencentcloud_ses_template" "example" {
  template_name = "tf_example_ses_temp"
  template_content {
    html = <<-EOT
<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>mail title</title>
</head>
<body>
<div class="container">
  <h1>Welcome to our service! </h1>
  <p>Dear user,</p>
  <p>Thank you for using Tencent Cloud:</p>
  <p><a href="https://cloud.tencent.com/document/product/1653">https://cloud.tencent.com/document/product/1653</a></p>
  <p>If you did not request this email, please ignore it. </p>
  <p><strong>from the iac team</strong></p>
</div>
</body>
</html>
    EOT
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

