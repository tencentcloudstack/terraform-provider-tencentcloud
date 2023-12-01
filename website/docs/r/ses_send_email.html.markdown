---
subcategory: "Simple Email Service(SES)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_ses_send_email"
sidebar_current: "docs-tencentcloud-resource-ses_send_email"
description: |-
  Provides a resource to create a ses send_email
---

# tencentcloud_ses_send_email

Provides a resource to create a ses send_email

## Example Usage

```hcl
resource "tencentcloud_ses_send_email" "send_email" {
  from_email_address = "aaa@iac-tf.cloud"
  destination        = ["1055482519@qq.com"]
  subject            = "test subject"
  reply_to_addresses = "aaa@iac-tf.cloud"

  template {
    template_id   = 99629
    template_data = "{\"name\":\"xxx\",\"age\":\"xx\"}"
  }

  unsubscribe  = "1"
  trigger_type = 1
}
```

## Argument Reference

The following arguments are supported:

* `destination` - (Required, Set: [`String`], ForceNew) Recipient email addresses. You can send an email to up to 50 recipients at a time. Note: the email content will display all recipient addresses. To send one-to-one emails to several recipients, please call the API multiple times to send the emails.
* `from_email_address` - (Required, String, ForceNew) Sender address. Enter a sender address, for example, noreply@mail.qcloud.com.To display the sender name, enter the address in the following format:Sender.
* `subject` - (Required, String, ForceNew) Email subject.
* `attachments` - (Optional, List, ForceNew) Parameters for the attachments to be sent. The TencentCloud API supports a request packet of up to 8 MB in size,and the size of the attachment content will increase by 1.5 times after Base64 encoding. Therefore,you need to keep the total size of all attachments below 4 MB. If the entire request exceeds 8 MB,the API will return an error.
* `bcc` - (Optional, Set: [`String`], ForceNew) The email address of the cc recipient can support up to 20 cc recipients.
* `cc` - (Optional, Set: [`String`], ForceNew) Cc recipient email address, up to 20 people can be copied.
* `reply_to_addresses` - (Optional, String, ForceNew) Reply-to address. You can enter a valid personal email address that can receive emails. If this parameter is left empty, reply emails will fail to be sent.
* `template` - (Optional, List, ForceNew) Template parameters for template-based sending. As Simple has been disused, Template is required.
* `trigger_type` - (Optional, Int, ForceNew) Email triggering type. 0 (default): non-trigger-based, suitable for marketing emails and non-immediate emails;1: trigger-based, suitable for immediate emails such as emails containing verification codes.If the size of an email exceeds a specified value,the system will automatically choose the non-trigger-based type.
* `unsubscribe` - (Optional, String, ForceNew) Unsubscribe link option.  0: Do not add unsubscribe link; 1: English 2: Simplified Chinese;  3: Traditional Chinese; 4: Spanish; 5: French;  6: German; 7: Japanese; 8: Korean;  9: Arabic; 10: Thai.

The `attachments` object supports the following:

* `content` - (Required, String) Base64-encoded attachment content. You can send attachments of up to 4 MB in the total size.Note: The TencentCloud API supports a request packet of up to 8 MB in size, and the size of the attachmentcontent will increase by 1.5 times after Base64 encoding. Therefore, you need to keep the total size of allattachments below 4 MB. If the entire request exceeds 8 MB, the API will return an error.
* `file_name` - (Required, String) Attachment name, which cannot exceed 255 characters. Some attachment types are not supported. For details, see [Attachment Types.](https://www.tencentcloud.com/document/product/1084/42373?has_map=1).

The `template` object supports the following:

* `template_data` - (Required, String) Variable parameters in the template. Please use json.dump to format the JSON object into a string type.The object is a set of key-value pairs. Each key denotes a variable, which is represented by {{key}}. The key will be replaced with the correspondingvalue (represented by {{value}}) when sending the email.Note: The parameter value cannot be data of a complex type such as HTML.Example: {name:xxx,age:xx}.
* `template_id` - (Required, Int) Template ID. If you do not have any template, please create one.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



