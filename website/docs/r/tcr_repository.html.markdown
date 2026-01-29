---
subcategory: "Tencent Container Registry(TCR)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_tcr_repository"
sidebar_current: "docs-tencentcloud-resource-tcr_repository"
description: |-
  Use this resource to create TCR repository.
---

# tencentcloud_tcr_repository

Use this resource to create TCR repository.

## Example Usage

### Create a tcr repository instance

```hcl
resource "tencentcloud_tcr_instance" "example" {
  name          = "tf-example"
  instance_type = "standard"
  delete_bucket = true
  tags = {
    "createdBy" = "Terraform"
  }
}

resource "tencentcloud_tcr_namespace" "example" {
  instance_id = tencentcloud_tcr_instance.example.id
  name        = "tf_example"
  severity    = "medium"
}

resource "tencentcloud_tcr_repository" "example" {
  instance_id    = tencentcloud_tcr_instance.example.id
  namespace_name = tencentcloud_tcr_namespace.example.name
  name           = "tf-example"
  brief_desc     = "desc."
  description    = "description."
  force_delete   = true
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String, ForceNew) ID of the TCR instance.
* `name` - (Required, String, ForceNew) Name of the TCR repository. Valid length is [2~200]. It can only contain lowercase letters, numbers and separators (`.`, `_`, `-`, `/`), and cannot start, end or continue with separators. Support the use of multi-level address formats, such as `sub1/sub2/repo`.
* `namespace_name` - (Required, String, ForceNew) Name of the TCR namespace.
* `brief_desc` - (Optional, String) Brief description of the repository. Valid length is [1~100].
* `description` - (Optional, String) Description of the repository. Valid length is [1~1000].
* `force_delete` - (Optional, Bool) The default value is true, meaning that the repository will be deleted directly regardless of whether it contains any images; false means that the existence of images will be checked before deleting the repository.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `create_time` - Create time.
* `is_public` - Indicate the repository is public or not.
* `update_time` - Last updated time.
* `url` - URL of the repository.


## Import

TCR repository can be imported using the instanceId#nameSpaceName#name, e.g.

```
terraform import tencentcloud_tcr_repository.example tcr-s1jud21h#tf_example#tf-example
```

