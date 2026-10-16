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

* Create: [CreateBlueprint](https://www.tencentcloud.com/search/?wd=CreateBlueprint&type=document&lang=en)
* Create: [DescribeBlueprints](https://www.tencentcloud.com/search/?wd=DescribeBlueprints&type=document&lang=en)
* Read: [DescribeBlueprints](https://www.tencentcloud.com/search/?wd=DescribeBlueprints&type=document&lang=en)
* Update: [ModifyBlueprintAttribute](https://www.tencentcloud.com/search/?wd=ModifyBlueprintAttribute&type=document&lang=en)
* Delete: [DeleteBlueprints](https://www.tencentcloud.com/search/?wd=DeleteBlueprints&type=document&lang=en)

## Import

lighthouse blueprint can be imported using the id, e.g.

```
terraform import tencentcloud_lighthouse_blueprint.blueprint blueprint_id
```

