---
subcategory: "Virtual Private Cloud(VPC)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_vpc_bandwidth_package"
sidebar_current: "docs-tencentcloud-resource-vpc_bandwidth_package"
description: |-
  Provides a resource to create a VPC bandwidth package
---

# tencentcloud_vpc_bandwidth_package

Provides a resource to create a VPC bandwidth package

## Example Usage

```hcl
resource "tencentcloud_vpc_bandwidth_package" "example" {
  network_type           = "BGP"
  charge_type            = "TOP5_POSTPAID_BY_MONTH"
  bandwidth_package_name = "tf-example"
  tag = {
    createdBy = "Terraform"
  }
}
```

### PrePaid Bandwidth Package

```hcl
resource "tencentcloud_vpc_bandwidth_package" "example" {
  network_type           = "BGP"
  charge_type            = "FIXED_PREPAID_BY_MONTH"
  bandwidth_package_name = "tf-example"
  time_span              = 3
  internet_max_bandwidth = 100
  tag = {
    createdBy = "Terraform"
  }
}
```

### Bandwidth Package With Egress

```hcl
resource "tencentcloud_vpc_bandwidth_package" "example" {
  network_type           = "SINGLEISP_CMCC"
  charge_type            = "ENHANCED95_POSTPAID_BY_MONTH"
  bandwidth_package_name = "tf-example"
  internet_max_bandwidth = 400
  egress                 = "center_egress2"
  tag = {
    createdBy = "Terraform"
  }
}
```

## Argument Reference

The following arguments are supported:

* `bandwidth_package_name` - (Optional, String) Bandwidth package name.
* `charge_type` - (Optional, String) Bandwidth package billing type, default: `TOP5_POSTPAID_BY_MONTH`. Optional value: `TOP5_POSTPAID_BY_MONTH`: TOP5 billed by monthly postpaid; `PERCENT95_POSTPAID_BY_MONTH`: 95 billed monthly postpaid; `FIXED_PREPAID_BY_MONTH`: Monthly prepaid billing (Type FIXED_PREPAID_BY_MONTH product API capability is under construction); `BANDWIDTH_POSTPAID_BY_DAY`: bandwidth billed by daily postpaid; `ENHANCED95_POSTPAID_BY_MONTH`: enhanced 95 billed monthly postpaid.
* `egress` - (Optional, String) Network egress. It defaults to `center_egress1`. If you want to try the egress feature, please [submit a ticket](https://console.cloud.tencent.com/workorder/category).
* `internet_max_bandwidth` - (Optional, Int) Bandwidth packet speed limit size. Unit: Mbps, -1 means no speed limit.
* `network_type` - (Optional, String) Bandwidth packet type, default: `BGP`. Optional value: `BGP`: common BGP shared bandwidth package; `HIGH_QUALITY_BGP`: High Quality BGP Shared Bandwidth Package; `SINGLEISP_CMCC`: CMCC shared bandwidth package; `SINGLEISP_CTCC:`: CTCC shared bandwidth package; `SINGLEISP_CUCC`: CUCC shared bandwidth package.
* `tag` - (Optional, Map) Tag description list.
* `tags` - (Optional, Map, **Deprecated**) Use `tag` instead. Tag description list.
* `time_span` - (Optional, Int, ForceNew) The purchase duration of the prepaid monthly bandwidth package, unit: month, value range: 1~60.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.


## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [CreateBandwidthPackage](https://www.tencentcloud.com/document/product/215/36922) - This API is used to create a [device bandwidth package](https://intl.cloud.tencent.com/document/product/684/15245?from_cn_redirect=1#bwptype) or an [IP bandwidth package](https://intl.cloud.tencent.com/document/product/684/15245?
* Create: [DescribeBandwidthPackages](https://www.tencentcloud.com/document/product/215/36919) - This API is used to query bandwidth package information, including the unique ID of the bandwidth package, the type, the billing mode, the name, and the resource information.
* Create: [ModifyResourceTags](https://www.tencentcloud.com/document/product/651/33012) - This API is used to modify all tags associated with a resource.
* Read: [DescribeBandwidthPackages](https://www.tencentcloud.com/document/product/215/36919) - This API is used to query bandwidth package information, including the unique ID of the bandwidth package, the type, the billing mode, the name, and the resource information.
* Read: [DescribeResourceTagsByResourceIds](https://www.tencentcloud.com/document/product/651/33022) - This API is used to query the tag key-value pairs associated with an existing resource.
* Update: [ModifyBandwidthPackageAttribute](https://www.tencentcloud.com/document/product/215/36918) - This API is used to modify the attributes of a bandwidth package, including the bandwidth package name, and so on.
* Update: [ModifyBandwidthPackageBandwidth](https://www.tencentcloud.com/document/product/215/65948) - This API is used to adjust the bandwidth of a [bandwidth package](https://www.tencentcloud.com/document/product/684/15245).
* Update: [ModifyResourceTags](https://www.tencentcloud.com/document/product/651/33012) - This API is used to modify all tags associated with a resource.
* Delete: [DeleteBandwidthPackage](https://www.tencentcloud.com/document/product/215/36921) - This API is used to delete bandwidth packages, including [device bandwidth packages](https://intl.cloud.tencent.com/document/product/684/15246?from_cn_redirect=1#.E8.AE.BE.E5.A4.87.E5.B8.A6.E5.AE.BD.E5.8C.

## Import

vpc bandwidth_package can be imported using the id, e.g.
```
$ terraform import tencentcloud_vpc_bandwidth_package.example bwp-hq8h7qpy
```

