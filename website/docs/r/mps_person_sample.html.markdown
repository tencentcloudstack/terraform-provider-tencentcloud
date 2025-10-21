---
subcategory: "Media Processing Service(MPS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_mps_person_sample"
sidebar_current: "docs-tencentcloud-resource-mps_person_sample"
description: |-
  Provides a resource to create a mps person_sample
---

# tencentcloud_mps_person_sample

Provides a resource to create a mps person_sample

## Example Usage

```hcl
resource "tencentcloud_mps_person_sample" "person_sample" {
  name = "test"
  usages = [
    "Review.Face"
  ]
  description = "test"
  face_contents = [
    filebase64("./person.png")
  ]
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, String) Material name, length limit: 20 characters.
* `usages` - (Required, Set: [`String`]) Material application scene, optional value:1. Recognition.Face: used for content recognition 2. Review.Face: used for inappropriate content identification 3. All: contains all of the above, equivalent to 1+2.
* `description` - (Optional, String) Material description, length limit: 1024 characters.
* `face_contents` - (Optional, Set: [`String`]) Material image [Base64](https://tools.ietf.org/html/rfc4648) encoded string only supports jpeg and png image formats. Array length limit: 5 images.Note: The picture must be a single portrait with clearer facial features, with a pixel size of not less than 200*200.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

mps person_sample can be imported using the id, e.g.

```
terraform import tencentcloud_mps_person_sample.person_sample person_sample_id
```

