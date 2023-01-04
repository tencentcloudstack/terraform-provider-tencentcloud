---
subcategory: "TencentCloud Elastic Microservice(TEM)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_tem_application"
sidebar_current: "docs-tencentcloud-resource-tem_application"
description: |-
  Provides a resource to create a tem application
---

# tencentcloud_tem_application

Provides a resource to create a tem application

## Example Usage

```hcl
resource "tencentcloud_tem_application" "application" {
  application_name          = "demo"
  description               = "demo for test"
  coding_language           = "JAVA"
  use_default_image_service = 0
  repo_type                 = 2
  repo_name                 = "qcloud/nginx"
  repo_server               = "ccr.ccs.tencentyun.com"
  tag {
    tag_key   = "createdBy"
    tag_value = "terraform"
  }
}
```

## Argument Reference

The following arguments are supported:

* `application_name` - (Required, String) application name.
* `coding_language` - (Required, String) program language, like JAVA.
* `description` - (Optional, String) application description.
* `instance_id` - (Optional, String) tcr instance id.
* `repo_name` - (Optional, String) repository name.
* `repo_server` - (Optional, String) registry address.
* `repo_type` - (Optional, Int) repo type, 0: tcr personal, 1: tcr enterprise, 2: public repository, 3: tcr hosted by tem, 4: demo image.
* `tag` - (Optional, List) application tag list.
* `use_default_image_service` - (Optional, Int) create image repo or not.

The `tag` object supports the following:

* `tag_key` - (Optional, String) tag key.
* `tag_value` - (Optional, String) tag value.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



