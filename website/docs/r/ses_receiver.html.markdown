---
subcategory: "Simple Email Service(SES)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_ses_receiver"
sidebar_current: "docs-tencentcloud-resource-ses_receiver"
description: |-
  Provides a resource to create a ses receiver
---

# tencentcloud_ses_receiver

Provides a resource to create a ses receiver

## Example Usage

```hcl
resource "tencentcloud_ses_receiver" "receiver" {
  receivers_name = "terraform_test"
  desc           = "description"

  data {
    email = "abc@abc.com"
  }

  data {
    email = "abcd@abcd.com"
  }
}
```



```hcl
resource "tencentcloud_ses_receiver" "receiver" {
  receivers_name = "terraform_test"
  desc           = "description"

  data {
    email         = "abc@abc.com"
    template_data = "{\"name\":\"xxx\",\"age\":\"xx\"}"
  }

  data {
    email         = "abcd@abcd.com"
    template_data = "{\"name\":\"xxx\",\"age\":\"xx\"}"
  }
}
```

## Argument Reference

The following arguments are supported:

* `data` - (Required, Set, ForceNew) Recipient email and template parameters in array format. The number of recipients is limited to within 20,000. If there is an object in the `data` list that inputs `template_data`, then other objects are also required.
* `receivers_name` - (Required, String, ForceNew) Recipient group name.
* `desc` - (Optional, String, ForceNew) Recipient group description.

The `data` object supports the following:

* `email` - (Required, String, ForceNew) Recipient email addresses.
* `template_data` - (Optional, String, ForceNew) Variable parameters in the template, please use json.dump to format the JSON object as a string type. The object is a set of key-value pairs, where each key represents a variable in the template, and the variables in the template are represented by {{key}}, and the corresponding values will be replaced with {{value}} when sent.Note: Parameter values cannot be complex data such as HTML. The total length of TemplateData (the entire JSON structure) should be less than 800 bytes.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.


## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [CreateReceiver](https://cloud.tencent.com/document/api/1288/68360) - This API is used to create a recipient group, which is the list of target email addresses for batch sending emails. After creating a group, you need to upload recipient email addresses.
* Create: [CreateReceiverDetail](https://cloud.tencent.com/document/api/1288/68359) - This API is used to add recipient email addresses (up to 20,000 at a time) to a recipient group. This will be processed asynchronously. If the data volume is large, it may take some time to upload.
* Create: [CreateReceiverDetailWithData](https://cloud.tencent.com/document/api/1288/70776) - Add recipient addresses with Template parameters. Use this API to import Template parameters while adding recipient addresses, ensuring each recipient address uses Template variables with different values when sending emails.
* Read: [ListReceiverDetails](https://cloud.tencent.com/document/api/1288/97935) - Query ALL recipient email addresses in the recipient list based on the recipient list id with paging query. Filter queries can be based on the recipient email address.
* Read: [ListReceivers](https://cloud.tencent.com/document/api/1288/68358) - This API is used to query recipient groups. It supports pagination, fuzzy query, and query by status.
* Delete: [DeleteReceiver](https://cloud.tencent.com/document/api/1288/71169) - This API is used to delete a recipient group and all recipient email addresses in the group based on the recipient group ID.

## Import

ses email_address can be imported using the id, e.g.
```
$ terraform import tencentcloud_ses_receiver.receiver receiverId
```

