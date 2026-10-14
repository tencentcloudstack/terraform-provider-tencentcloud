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
  tags = {
    "created" = "terraform"
  }
}
```

## Argument Reference

The following arguments are supported:

* `application_name` - (Required, String) application name.
* `coding_language` - (Required, String) program language, like JAVA.
* `description` - (Required, String) application description.
* `instance_id` - (Optional, String) tcr instance id.
* `repo_name` - (Optional, String) repository name.
* `repo_server` - (Optional, String) registry address.
* `repo_type` - (Optional, Int) repo type, 0: tcr personal, 1: tcr enterprise, 2: public repository, 3: tcr hosted by tem, 4: demo image.
* `tags` - (Optional, Map) application tag list.
* `use_default_image_service` - (Optional, Int) create image repo or not.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.


## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [CreateApplication](https://cloud.tencent.com/document/api/1371/60162) - This API is used to create an application.
* Create: [ModifyResourceTags](https://cloud.tencent.com/document/api/651/35322) - This API is used to modify all tags associated with a resource.
* Read: [DescribeApplications](https://cloud.tencent.com/document/api/1371/78169) - This API is to query the list of running applications.
* Read: [DescribeResourceTagsByResourceIds](https://cloud.tencent.com/document/api/651/35319) - This API is used to query the tag key-value pairs associated with an existing resource.
* Update: [ModifyApplicationInfo](https://cloud.tencent.com/document/api/1371/60160) - This API is used to modify the basic information of an application.
* Update: [ModifyResourceTags](https://cloud.tencent.com/document/api/651/35322) - This API is used to modify all tags associated with a resource.
* Delete: [DeleteApplication](https://cloud.tencent.com/document/api/1371/66645) - This API is used to delete an application.

