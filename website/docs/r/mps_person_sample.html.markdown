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


## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [CreatePersonSample](https://cloud.tencent.com/document/api/862/39445) - This API is used to create image samples for video processing operations such as content recognition and inappropriate information detection with the help of technologies such as facial feature positioning.
* Read: [DescribePersonSamples](https://cloud.tencent.com/document/api/862/39441) - This API is used to query the information of image samples. It supports paginated queries by image ID, name, and tag.
* Update: [ModifyPersonSample](https://cloud.tencent.com/document/api/862/39439) - This API is used to modify image samples by image ID. You can use it to modify the name and description of an image sample and add/delete/reset facial features or tags. There must be at least one image left after the deletion of facial features; otherwise, please reset instead of delete the facial features.
* Delete: [DeletePersonSample](https://cloud.tencent.com/document/api/862/39443) - This API is used to delete image samples by image ID.

## Import

mps person_sample can be imported using the id, e.g.

```
terraform import tencentcloud_mps_person_sample.person_sample person_sample_id
```

