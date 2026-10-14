---
subcategory: "Tencent Container Registry(TCR)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_tcr_token"
sidebar_current: "docs-tencentcloud-resource-tcr_token"
description: |-
  Use this resource to create tcr long term token.
---

# tencentcloud_tcr_token

Use this resource to create tcr long term token.

## Example Usage

### Create a token for tcr instance

```hcl
resource "tencentcloud_tcr_instance" "example" {
  name          = "tf-example-tcr"
  instance_type = "basic"
  delete_bucket = true
  tags = {
    "createdBy" = "terraform"
  }
}

resource "tencentcloud_tcr_token" "example" {
  instance_id = tencentcloud_tcr_instance.example.id
  description = "example for the tcr token"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String, ForceNew) ID of the TCR instance.
* `description` - (Optional, String, ForceNew) Description of the token. Valid length is [0~255].
* `enable` - (Optional, Bool) Indicate to enable this token or not.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `create_time` - Create time.
* `token_id` - Sub ID of the TCR token. The full ID of token format like `instance_id#token_id`.
* `token` - The content of the token.
* `user_name` - User name of the token.

## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [CreateInstanceToken](https://cloud.tencent.com/document/api/1141/41571) - This API is used to create a temporary or long-term instance access credential.
* Create: [ModifyInstanceToken](https://cloud.tencent.com/document/api/1141/42910) - This API is used to update the status of the specified long-term access credential in an instance.
* Read: [DescribeInstanceToken](https://cloud.tencent.com/document/api/1141/42911) - This API is used to query the information of long-term access credentials.
* Update: [ModifyInstanceToken](https://cloud.tencent.com/document/api/1141/42910) - This API is used to update the status of the specified long-term access credential in an instance.
* Delete: [DeleteInstanceToken](https://cloud.tencent.com/document/api/1141/42912) - This API is used to delete a long-term access credential.
* Delete: [DescribeInstanceToken](https://cloud.tencent.com/document/api/1141/42911) - This API is used to query the information of long-term access credentials.

## Import

tcr token can be imported using the id, e.g.

```
$ terraform import tencentcloud_tcr_token.example instance_id#token_id
```

