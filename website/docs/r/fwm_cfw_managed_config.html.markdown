---
subcategory: "Firewall Manager(FWM)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_fwm_cfw_managed_config"
sidebar_current: "docs-tencentcloud-resource-fwm_cfw_managed_config"
description: |-
  Provides a resource to manage the cloud firewall managed switch of an organization member for the fwm (Cloud Firewall Managed) product.
---

# tencentcloud_fwm_cfw_managed_config

Provides a resource to manage the cloud firewall managed switch of an organization member for the fwm (Cloud Firewall Managed) product.

## Example Usage

```hcl
resource "tencentcloud_fwm_cfw_managed_config" "example" {
  member_id   = "mem-tencent-84035535ae473b76"
  cfw_managed = 1
}
```

## Argument Reference

The following arguments are supported:

* `cfw_managed` - (Required, Int) Cloud firewall managed switch, 1: enable managed, 0: disable managed.
* `member_id` - (Required, String, ForceNew) Target member id of the organization.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.


## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [ModifyCfwManaged](https://www.tencentcloud.com/search/?wd=ModifyCfwManaged&type=document&lang=en)
* Read: [DescribeOrganMembers](https://www.tencentcloud.com/search/?wd=DescribeOrganMembers&type=document&lang=en)
* Update: [ModifyCfwManaged](https://www.tencentcloud.com/search/?wd=ModifyCfwManaged&type=document&lang=en)

## Import

fwm cfw managed can be imported using the composite id, e.g.

```
terraform import tencentcloud_fwm_cfw_managed_config.example mem-tencent-84035535ae473b76
```

