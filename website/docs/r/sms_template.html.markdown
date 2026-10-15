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

* Create: [AddSmsTemplate](https://cloud.tencent.com/document/api/382/55974) - 添加短信模板
* Read: [DescribeSmsTemplateList](https://cloud.tencent.com/document/api/382/52067) - 短信模板状态查询
* Update: [ModifySmsTemplate](https://cloud.tencent.com/document/api/382/55972) - 修改短信模板
* Delete: [DeleteSmsTemplate](https://cloud.tencent.com/document/api/382/55973) - 删除短信模板

