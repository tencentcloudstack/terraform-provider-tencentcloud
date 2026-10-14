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

* Create: [CreateEmailAddress](https://cloud.tencent.com/document/api/1288/51052) - After the sender domain is verified, you need a sender address to send emails. For example, if your sender domain is mail.qcloud.com, your sender address can be service@mail.qcloud.com. If you want to display your name (such as "Tencent Cloud") in the inbox list of the recipients, the sender address should be in the format of `Tencent Cloud <email address>`. Please note that there must be a space between your name and the first angle bracket.
* Create: [UpdateEmailSmtpPassWord](https://cloud.tencent.com/document/api/1288/84274) - This API is used to set the SMTP password. Initially, no SMTP password is set for your email address, so emails cannot be sent over SMTP. To send emails over SMTP, you must set the SMTP password. The set password can be changed subsequently.
* Read: [ListEmailAddress](https://cloud.tencent.com/document/api/1288/51050) - This API is used to get the list of sender addresses.
* Update: [UpdateEmailSmtpPassWord](https://cloud.tencent.com/document/api/1288/84274) - This API is used to set the SMTP password. Initially, no SMTP password is set for your email address, so emails cannot be sent over SMTP. To send emails over SMTP, you must set the SMTP password. The set password can be changed subsequently.
* Delete: [DeleteEmailAddress](https://cloud.tencent.com/document/api/1288/51051) - This API is used to delete a sender address.

## Import

ses email_address can be imported using the id, e.g.
```
$ terraform import tencentcloud_ses_email_address.example demo@iac-terraform.cloud
```

