---
subcategory: "Simple Email Service(SES)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_ses_batch_send_email"
sidebar_current: "docs-tencentcloud-resource-ses_batch_send_email"
description: |-
  Provides a resource to create a ses batch_send_email
---

# tencentcloud_ses_batch_send_email

Provides a resource to create a ses batch_send_email

## Example Usage

```hcl
resource "tencentcloud_ses_batch_send_email" "batch_send_email" {
  from_email_address = "aaa@iac-tf.cloud"
  receiver_id        = 1063742
  subject            = "terraform test"
  task_type          = 1
  reply_to_addresses = "reply@mail.qcloud.com"
  template {
    template_id   = 99629
    template_data = "{\"name\":\"xxx\",\"age\":\"xx\"}"

  }

  cycle_param {
    begin_time    = "2023-09-07 15:10:00"
    interval_time = 1
  }
  timed_param {
    begin_time = "2023-09-07 15:20:00"
  }
  unsubscribe = "0"
  ad_location = 0
}
```

## Argument Reference

The following arguments are supported:

* `from_email_address` - (Required, String, ForceNew) Sender address. Enter a sender address such as noreply@mail.qcloud.com. To display the sender name, enter the address in the following format:sender &amp;amp;lt;email address&amp;amp;gt;. For example:Tencent Cloud team &amp;amp;lt;noreply@mail.qcloud.com&amp;amp;gt;.
* `receiver_id` - (Required, Int, ForceNew) Recipient group ID.
* `subject` - (Required, String, ForceNew) Email subject.
* `task_type` - (Required, Int, ForceNew) Task type. 1: immediate; 2: scheduled; 3: recurring.
* `ad_location` - (Optional, Int, ForceNew) Whether to add an ad tag. 0: Add no tag; 1: Add before the subject; 2: Add after the subject.
* `attachments` - (Optional, List, ForceNew) Attachment parameters to set when you need to send attachments. This parameter is currently unavailable.
* `cycle_param` - (Optional, List, ForceNew) Parameter required for a recurring sending task.
* `reply_to_addresses` - (Optional, String, ForceNew) Reply-to address. You can enter a valid personal email address that can receive emails. If this parameter is left empty, reply emails will fail to be sent.
* `template` - (Optional, List, ForceNew) Template when emails are sent using a template.
* `timed_param` - (Optional, List, ForceNew) Parameter required for a scheduled sending task.
* `unsubscribe` - (Optional, String, ForceNew) Unsubscribe link option.  0: Do not add unsubscribe link; 1: English 2: Simplified Chinese;  3: Traditional Chinese; 4: Spanish; 5: French;  6: German; 7: Japanese; 8: Korean;  9: Arabic; 10: Thai.

The `attachments` object supports the following:

* `content` - (Required, String) Base64-encoded attachment content. You can send attachments of up to 4 MB in the total size.Note: The TencentCloud API supports a request packet of up to 8 MB in size, and the size of the attachmentcontent will increase by 1.5 times after Base64 encoding. Therefore, you need to keep the total size of allattachments below 4 MB. If the entire request exceeds 8 MB, the API will return an error.
* `file_name` - (Required, String) Attachment name, which cannot exceed 255 characters. Some attachment types are not supported. For details, see [Attachment Types.](https://www.tencentcloud.com/document/product/1084/42373?has_map=1).

The `cycle_param` object supports the following:

* `begin_time` - (Required, String) Start time of the task.
* `interval_time` - (Required, Int) Task recurrence in hours.
* `term_cycle` - (Optional, Int) Specifies whether to end the cycle. This parameter is used to update the task. Valid values: 0: No; 1: Yes.

The `template` object supports the following:

* `template_data` - (Required, String) Variable parameters in the template. Please use json.dump to format the JSON object into a string type.The object is a set of key-value pairs. Each key denotes a variable, which is represented by {{key}}. The key will be replaced with the correspondingvalue (represented by {{value}}) when sending the email.Note: The parameter value cannot be data of a complex type such as HTML.Example: {name:xxx,age:xx}.
* `template_id` - (Required, Int) Template ID. If you do not have any template, please create one.

The `timed_param` object supports the following:

* `begin_time` - (Required, String) Start time of a scheduled sending task.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



