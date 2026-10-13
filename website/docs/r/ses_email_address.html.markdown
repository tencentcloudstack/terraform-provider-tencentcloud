---
subcategory: "Simple Email Service(SES)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_ses_email_address"
sidebar_current: "docs-tencentcloud-resource-ses_email_address"
description: |-
  Provides a resource to create a ses email address
---

# tencentcloud_ses_email_address

Provides a resource to create a ses email address

## Example Usage

### Create ses email address

```hcl
resource "tencentcloud_ses_email_address" "example" {
  email_address     = "demo@iac-terraform.cloud"
  email_sender_name = "root"
}
```

### Set smtp password

```hcl
resource "tencentcloud_ses_email_address" "example" {
  email_address     = "demo@iac-terraform.cloud"
  email_sender_name = "root"
  smtp_password     = "Password@123"
}
```

## Argument Reference

The following arguments are supported:

* `email_address` - (Required, String, ForceNew) Your sender address(You can create up to 10 sender addresses for each domain).
* `email_sender_name` - (Optional, String, ForceNew) Sender name.
* `smtp_password` - (Optional, String) Password for SMTP, Length limit 64.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.


## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [CreateEmailAddress](https://cloud.tencent.com/document/api/1288/51052) - 新建发信地址
* Create: [UpdateEmailSmtpPassWord](https://cloud.tencent.com/document/api/1288/84274) - 设置邮箱的smtp密码
* Read: [ListEmailAddress](https://cloud.tencent.com/document/api/1288/51050) - 获取发信地址列表
* Update: [UpdateEmailSmtpPassWord](https://cloud.tencent.com/document/api/1288/84274) - 设置邮箱的smtp密码
* Delete: [DeleteEmailAddress](https://cloud.tencent.com/document/api/1288/51051) - 删除发信地址

## Import

ses email_address can be imported using the id, e.g.
```
$ terraform import tencentcloud_ses_email_address.example demo@iac-terraform.cloud
```

