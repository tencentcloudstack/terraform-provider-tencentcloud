---
subcategory: "Direct Connect Gateway(DCG)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dc_gateway_attachment"
sidebar_current: "docs-tencentcloud-resource-dc_gateway_attachment"
description: |-
  Provides a resource to create a dc_gateway_attachment
---

# tencentcloud_dc_gateway_attachment

Provides a resource to create a dc_gateway_attachment

## Example Usage

```hcl
resource "tencentcloud_dc_gateway_attachment" "dc_gateway_attachment" {
  vpc_id                    = "vpc-4h9v4mo3"
  nat_gateway_id            = "nat-7kanjc6y"
  direct_connect_gateway_id = "dcg-dmbhf7jf"
}
```

## Argument Reference

The following arguments are supported:

* `direct_connect_gateway_id` - (Required, String, ForceNew) DirectConnectGatewayId.
* `nat_gateway_id` - (Required, String, ForceNew) NatGatewayId.
* `vpc_id` - (Required, String, ForceNew) vpc id.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.


## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [AssociateDirectConnectGatewayNatGateway](https://cloud.tencent.com/document/api/215/47109) - 专线网关绑定NAT网关
* Read: [DescribeDirectConnectGateways](https://cloud.tencent.com/document/api/215/30644) - 查询专线网关
* Delete: [DisassociateDirectConnectGatewayNatGateway](https://cloud.tencent.com/document/api/215/47108) - 专线网关解绑NAT网关

## Import

dc_gateway_attachment can be imported using the id, e.g.

```
terraform import tencentcloud_vpc_dc_gateway_attachment.dc_gateway_attachment vpcId#dcgId#ngId
```

