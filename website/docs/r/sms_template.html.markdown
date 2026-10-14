---
subcategory: "Short Message Service(SMS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_sms_template"
sidebar_current: "docs-tencentcloud-resource-sms_template"
description: |-
  Provides a resource to create a sms template
---

# tencentcloud_sms_template

Provides a resource to create a sms template

## Example Usage

### Create a sms template

```hcl
resource "tencentcloud_sms_template" "template" {
  template_name    = "tf_example_sms_template"
  template_content = "example for sms template"
  international    = 0 # Mainland China SMS
  sms_type         = 0 # regular SMS
  remark           = "terraform example"
}
```

## Argument Reference

The following arguments are supported:

* `international` - (Required, Int) Whether it is Global SMS: 0: Mainland China SMS; 1: Global SMS.
* `remark` - (Required, String) Template remarks, such as reason for application and use case.
* `sms_type` - (Required, Int) SMS type. 0: regular SMS, 1: marketing SMS.
* `template_content` - (Required, String) Message Template Content.
* `template_name` - (Required, String) Message Template name, which must be unique.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.


## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [AddSmsTemplate](https://cloud.tencent.com/document/api/382/55974) - This API is used to add an SMS template. Please read the [Tencent Cloud SMS Body Template Review Standards](https://intl.cloud.tencent.com/document/product/382/39023?from_cn_redirect=1) before starting an application.
* Read: [DescribeSmsTemplateList](https://cloud.tencent.com/document/api/382/52067)
* Update: [ModifySmsTemplate](https://cloud.tencent.com/document/api/382/55972) - This API is used to modify an SMS body template. Please read the [Tencent Cloud SMS Body Template Review Standards](https://intl.cloud.tencent.com/document/product/382/39023?from_cn_redirect=1) before making a modification.
* Delete: [DeleteSmsTemplate](https://cloud.tencent.com/document/api/382/55973)

