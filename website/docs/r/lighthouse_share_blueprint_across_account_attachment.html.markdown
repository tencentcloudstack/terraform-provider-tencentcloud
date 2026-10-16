---
subcategory: "TencentCloud Lighthouse(Lighthouse)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_lighthouse_share_blueprint_across_account_attachment"
sidebar_current: "docs-tencentcloud-resource-lighthouse_share_blueprint_across_account_attachment"
description: |-
  Provides a resource to create a lighthouse share blueprint across account attachment share
---

# tencentcloud_lighthouse_share_blueprint_across_account_attachment

Provides a resource to create a lighthouse share blueprint across account attachment share

## Example Usage

```hcl
resource "tencentcloud_lighthouse_share_blueprint_across_account_attachment" "share_blueprint_across_account_attachment" {
  blueprint_id = "lhbp-xxxxxx"
  account_ids  = ["100012345678"]
}
```

## Argument Reference

The following arguments are supported:

* `account_ids` - (Required, Set: [`String`]) List of target TencentCloud account IDs to share the blueprint with.
* `blueprint_id` - (Required, String, ForceNew) Lighthouse blueprint ID.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.


## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [DescribeBlueprintsShareAcrossAccountInfos](https://www.tencentcloud.com/search/?wd=DescribeBlueprintsShareAcrossAccountInfos&type=document&lang=en)
* Create: [ShareBlueprintAcrossAccounts](https://www.tencentcloud.com/document/api/1103/75693)
* Read: [DescribeBlueprintsShareAcrossAccountInfos](https://www.tencentcloud.com/search/?wd=DescribeBlueprintsShareAcrossAccountInfos&type=document&lang=en)
* Update: [CancelShareBlueprintAcrossAccounts](https://www.tencentcloud.com/document/api/1103/75613)
* Update: [DescribeBlueprintsShareAcrossAccountInfos](https://www.tencentcloud.com/search/?wd=DescribeBlueprintsShareAcrossAccountInfos&type=document&lang=en)
* Update: [ShareBlueprintAcrossAccounts](https://www.tencentcloud.com/document/api/1103/75693)
* Delete: [CancelShareBlueprintAcrossAccounts](https://www.tencentcloud.com/document/api/1103/75613)

## Import

tencentcloud_lighthouse_share_blueprint_across_account_attachment can be imported using the blueprint_id, e.g.

```
terraform import tencentcloud_lighthouse_share_blueprint_across_account_attachment.share_blueprint_across_account_attachment lhbp-xxxxxx
```

