---
subcategory: "Cloud Load Balancer(CLB)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_clb_snat_ip"
sidebar_current: "docs-tencentcloud-resource-clb_snat_ip"
description: |-
  Provide a resource to create a SnatIp of CLB instance.
---

# tencentcloud_clb_snat_ip

Provide a resource to create a SnatIp of CLB instance.

~> **NOTE:** Target CLB instance must enable `snat_pro` before creating snat ips.

~> **NOTE:** Dynamic allocate IP doesn't support for now.

## Example Usage

```hcl
resource "tencentcloud_clb_snat_ip" "example" {
  clb_id = "lb-jnx618r2"
  ips {
    subnet_id = "subnet-hhi88a58"
    ip        = "10.0.30.10"
  }

  ips {
    subnet_id = "subnet-d4umunpy"
  }
}
```

## Argument Reference

The following arguments are supported:

* `clb_id` - (Required, String) CLB instance ID.
* `ips` - (Optional, Set) Snat IP address config.

The `ips` object supports the following:

* `subnet_id` - (Required, String) Subnet ID.
* `ip` - (Optional, String) Snat IP.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.


## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [CreateLoadBalancerSnatIps](https://cloud.tencent.com/document/api/214/41505) - 添加SnatIp
* Read: [DescribeLoadBalancers](https://cloud.tencent.com/document/api/214/30685) - 查询负载均衡实例列表
* Update: [CreateLoadBalancerSnatIps](https://cloud.tencent.com/document/api/214/41505) - 添加SnatIp
* Update: [DeleteLoadBalancerSnatIps](https://cloud.tencent.com/document/api/214/41503) - 删除SnatIp
* Delete: [DeleteLoadBalancerSnatIps](https://cloud.tencent.com/document/api/214/41503) - 删除SnatIp
* Delete: [DescribeLoadBalancers](https://cloud.tencent.com/document/api/214/30685) - 查询负载均衡实例列表

## Import

Clb instance snat ip can be imported by clb instance id, e.g.

```
terraform import tencentcloud_clb_snat_ip.example lb-jnx618r2
```

