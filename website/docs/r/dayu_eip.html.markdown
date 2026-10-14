---
subcategory: "Anti-DDoS(DayuV2)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dayu_eip"
sidebar_current: "docs-tencentcloud-resource-dayu_eip"
description: |-
  Use this resource to create dayu eip rule
---

# tencentcloud_dayu_eip

Use this resource to create dayu eip rule

## Example Usage

```hcl
resource "tencentcloud_dayu_eip" "test" {
  resource_id          = "bgpip-000004xg"
  eip                  = "162.62.163.50"
  bind_resource_id     = "ins-4m0jvxic"
  bind_resource_region = "hk"
  bind_resource_type   = "cvm"
}
```

## Argument Reference

The following arguments are supported:

* `bind_resource_id` - (Required, String, ForceNew) Resource id to bind.
* `bind_resource_region` - (Required, String, ForceNew) Resource region to bind.
* `bind_resource_type` - (Required, String, ForceNew) Resource type to bind, value range [`clb`, `cvm`].
* `eip` - (Required, String, ForceNew) Eip of the resource.
* `resource_id` - (Required, String, ForceNew) ID of the resource.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `created_time` - Created time of the resource instance.
* `eip_address_status` - Eip address status of the resource instance.
* `eip_bound_rsc_eni` - Eip bound rsc eni of the resource instance.
* `eip_bound_rsc_ins` - Eip bound rsc ins of the resource instance.
* `eip_bound_rsc_vip` - Eip bound rsc vip of the resource instance.
* `expired_time` - Expired time of the resource instance.
* `modify_time` - Modify time of the resource instance.
* `protection_status` - Protection status of the resource instance.
* `resource_region` - Region of the resource instance.

## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [AssociateDDoSEipAddress](https://cloud.tencent.com/document/api/297/95297) - This API is used to bind an EIP to an Anti-DDoS Advanced instance or a specified private IP of an ENI.
* Create: [AssociateDDoSEipLoadBalancer](https://cloud.tencent.com/document/api/297/95296) - This API is used to bind an Anti-DDoS EIP to the specified private IP of a CLB instance.
* Create: [DescribeListBGPIPInstances](https://cloud.tencent.com/document/api/297/95347) - This API is used to get a list of Anti-DDoS Advanced instances.
* Read: [DescribeListBGPIPInstances](https://cloud.tencent.com/document/api/297/95347) - This API is used to get a list of Anti-DDoS Advanced instances.
* Delete: [DisassociateDDoSEipAddress](https://cloud.tencent.com/document/api/297/95295) - This API is used to unbind an Anti-DDoS EIP.

