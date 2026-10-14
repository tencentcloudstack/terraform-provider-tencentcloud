---
subcategory: "TencentCloud Lighthouse(Lighthouse)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_lighthouse_blueprint"
sidebar_current: "docs-tencentcloud-resource-lighthouse_blueprint"
description: |-
  Provides a resource to create a lighthouse blueprint
---

# tencentcloud_lighthouse_blueprint

Provides a resource to create a lighthouse blueprint

## Example Usage

```hcl
resource "tencentcloud_lighthouse_blueprint" "blueprint" {
  blueprint_name = "blueprint_name_test"
  description    = "blueprint_description_test"
  instance_id    = "lhins-xxxxxx"
}
```

## Argument Reference

The following arguments are supported:

* `blueprint_name` - (Required, String) Blueprint name, which can contain up to 60 characters.
* `description` - (Optional, String) Blueprint description, which can contain up to 60 characters.
* `instance_id` - (Optional, String) ID of the instance for which to make a blueprint.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.


## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [CreateBlueprint](https://cloud.tencent.com/document/api/1207/54385) - This API is used to create an image.
* Create: [DescribeBlueprints](https://cloud.tencent.com/document/api/1207/47689) - This API is used to query the information of an image.
* Read: [DescribeBlueprints](https://cloud.tencent.com/document/api/1207/47689) - This API is used to query the information of an image.
* Update: [ModifyBlueprintAttribute](https://cloud.tencent.com/document/api/1207/54383) - This API is used to modify an image attribute.
* Delete: [DeleteBlueprints](https://cloud.tencent.com/document/api/1207/54384) - This API is used to delete an image.

## Import

lighthouse blueprint can be imported using the id, e.g.

```
terraform import tencentcloud_lighthouse_blueprint.blueprint blueprint_id
```

