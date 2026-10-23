---
subcategory: "Cloud Firewall Managed(FWM)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_fwm_cfw_managed"
sidebar_current: "docs-tencentcloud-resource-fwm_cfw_managed"
description: |-
  Provides a resource to manage the cloud firewall managed switch of an organization member for the fwm (Cloud Firewall Managed) product.
---

# tencentcloud_fwm_cfw_managed

Provides a resource to manage the cloud firewall managed switch of an organization member for the fwm (Cloud Firewall Managed) product.

## Example Usage

```hcl
resource "tencentcloud_fwm_cfw_managed" "example" {
  app_id      = "1234567890"
  member_id   = "100012345678"
  cfw_managed = 1
}
```

## Argument Reference

The following arguments are supported:

* `app_id` - (Required, String, ForceNew) App id of the target member, used to locate the member by AppId filter.
* `cfw_managed` - (Required, Int) Cloud firewall managed switch, 1: enable managed, 0: disable managed.
* `member_id` - (Required, String, ForceNew) Target member id of the organization.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.


## API Reference

This resource interacts with the following TencentCloud APIs:

* Read: [DescribeOrganMembers](https://www.tencentcloud.com/search/?wd=DescribeOrganMembers&type=document&lang=en)
* Update:
  - [DescribeOrganMembers](https://www.tencentcloud.com/search/?wd=DescribeOrganMembers&type=document&lang=en)
  - [ModifyCfwManaged](https://www.tencentcloud.com/search/?wd=ModifyCfwManaged&type=document&lang=en)

## Import

fwm cfw managed can be imported using the composite id, e.g. `app_id#member_id`.

```
terraform import tencentcloud_fwm_cfw_managed.example 1234567890#100012345678
```

