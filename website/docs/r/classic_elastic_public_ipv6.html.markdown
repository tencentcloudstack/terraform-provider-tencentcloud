---
subcategory: "Virtual Private Cloud(VPC)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_classic_elastic_public_ipv6"
sidebar_current: "docs-tencentcloud-resource-classic_elastic_public_ipv6"
description: |-
  Provides a resource to create a vpc classic_elastic_public_ipv6
---

# tencentcloud_classic_elastic_public_ipv6

Provides a resource to create a vpc classic_elastic_public_ipv6

## Example Usage

```hcl
resource "tencentcloud_classic_elastic_public_ipv6" "classic_elastic_public_ipv6" {
  ip6_address                = "xxxxxx"
  internet_max_bandwidth_out = 2
  tags = {
    "testkey" = "testvalue"
  }
}
```

## Argument Reference

The following arguments are supported:

* `ip6_address` - (Required, String) IPV6 addresses that require public network access.
* `bandwidth_package_id` - (Optional, String) Bandwidth package id, move the account up, and you need to pass in the ipv6 address to apply for bandwidth package charging mode.
* `internet_charge_type` - (Optional, String) Network billing model. IPV6 currently supports `TRAFFIC_POSTPAID_BY_HOUR` and `BANDWIDTH_PACKAGE`. The default network charging mode is `TRAFFIC_POSTPAID_BY_HOUR`.
* `internet_max_bandwidth_out` - (Optional, Int) Bandwidth in Mbps. Default is 1Mbps.
* `tags` - (Optional, Map) Tags.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.


## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [AllocateIp6AddressesBandwidth](https://cloud.tencent.com/document/api/215/40090) - 开通IPv6公网带宽
* Create: [DescribeTaskResult](https://cloud.tencent.com/document/api/215/36271) - 查询异步任务执行结果
* Read: [DescribeIp6Addresses](https://cloud.tencent.com/document/api/215/40089) - 查询传统弹性公网IPv6列表
* Read: [DescribeResourceTagsByResourceIds](https://cloud.tencent.com/document/api/651/35319) - 批量查看资源关联的标签
* Update: [DescribeTaskResult](https://cloud.tencent.com/document/api/215/36271) - 查询异步任务执行结果
* Update: [ModifyIp6AddressesBandwidth](https://cloud.tencent.com/document/api/215/40088) - 调整传统弹性公网IPv6带宽
* Update: [ModifyResourceTags](https://cloud.tencent.com/document/api/651/35322) - 批量修改资源关联的标签
* Delete: [DescribeTaskResult](https://cloud.tencent.com/document/api/215/36271) - 查询异步任务执行结果
* Delete: [ReleaseIp6AddressesBandwidth](https://cloud.tencent.com/document/api/215/40087) - 关闭IPv6公网带宽

## Import

vpc classic_elastic_public_ipv6 can be imported using the id, e.g.

```
terraform import tencentcloud_classic_elastic_public_ipv6.classic_elastic_public_ipv6 classic_elastic_public_ipv6_id
```

