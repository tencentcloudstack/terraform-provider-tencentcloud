---
subcategory: "Cloud Virtual Machine(CVM)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_eip_association"
sidebar_current: "docs-tencentcloud-resource-eip_association"
description: |-
  Provides an eip resource associated with other resource like CVM, ENI and CLB.
---

# tencentcloud_eip_association

Provides an eip resource associated with other resource like CVM, ENI and CLB.

~> **NOTE:** Please DO NOT define `allocate_public_ip` in `tencentcloud_instance` resource when using `tencentcloud_eip_association`.

## Example Usage

### Bind elastic public IP By Instance ID

```hcl
data "tencentcloud_availability_zones" "zones" {}

data "tencentcloud_images" "image" {
  image_type       = ["PUBLIC_IMAGE"]
  image_name_regex = "Final"
}

data "tencentcloud_instance_types" "instance_types" {
  filter {
    name   = "zone"
    values = [data.tencentcloud_availability_zones.zones.zones.0.name]
  }

  filter {
    name   = "instance-family"
    values = ["S5"]
  }

  cpu_core_count   = 2
  exclude_sold_out = true
}

resource "tencentcloud_vpc" "vpc" {
  name       = "example-vpc"
  cidr_block = "10.0.0.0/16"
}

resource "tencentcloud_subnet" "subnet" {
  availability_zone = data.tencentcloud_availability_zones.zones.zones.0.name
  name              = "example-vpc"
  vpc_id            = tencentcloud_vpc.vpc.id
  cidr_block        = "10.0.0.0/16"
  is_multicast      = false
}

resource "tencentcloud_eip" "eip" {
  name                 = "example-eip"
  internet_charge_type = "TRAFFIC_POSTPAID_BY_HOUR"
  type                 = "EIP"
}

resource "tencentcloud_instance" "example" {
  instance_name            = "example-cvm"
  availability_zone        = data.tencentcloud_availability_zones.zones.zones.0.name
  image_id                 = data.tencentcloud_images.image.images.0.image_id
  instance_type            = data.tencentcloud_instance_types.instance_types.instance_types.0.instance_type
  system_disk_type         = "CLOUD_PREMIUM"
  disable_security_service = true
  disable_monitor_service  = true
  vpc_id                   = tencentcloud_vpc.vpc.id
  subnet_id                = tencentcloud_subnet.subnet.id
}

resource "tencentcloud_eip_association" "example" {
  eip_id      = tencentcloud_eip.eip.id
  instance_id = tencentcloud_instance.example.id
}
```

### Bind elastic public IP By elastic network card

```hcl
data "tencentcloud_availability_zones" "zones" {}

resource "tencentcloud_vpc" "vpc" {
  name       = "example-vpc"
  cidr_block = "10.0.0.0/16"
}

resource "tencentcloud_subnet" "subnet" {
  availability_zone = data.tencentcloud_availability_zones.zones.zones.0.name
  name              = "example-vpc"
  vpc_id            = tencentcloud_vpc.vpc.id
  cidr_block        = "10.0.0.0/16"
  is_multicast      = false
}

resource "tencentcloud_eni" "eni" {
  name        = "example-eni"
  vpc_id      = tencentcloud_vpc.vpc.id
  subnet_id   = tencentcloud_subnet.subnet.id
  description = "eni desc"
  ipv4_count  = 1
}

resource "tencentcloud_eip" "eip" {
  name                 = "example-eip"
  internet_charge_type = "TRAFFIC_POSTPAID_BY_HOUR"
  type                 = "EIP"
}

resource "tencentcloud_eip_association" "example" {
  eip_id               = tencentcloud_eip.eip.id
  network_interface_id = tencentcloud_eni.eni.id
  private_ip           = tencentcloud_eni.eni.ipv4_info[0].ip
}
```

### Bind elastic public IP By WAF(SaaS)

```hcl
resource "tencentcloud_eip_association" "example" {
  eip_id      = "eip-4stgtfb8"
  instance_id = "waf_2opxlbky67yub3gm"
}
```

## Argument Reference

The following arguments are supported:

* `eip_id` - (Required, String, ForceNew) The ID of EIP.
* `instance_id` - (Optional, String, ForceNew) The ID of the target resource to associate with the Elastic IP (EIP). Supported targets include a CVM instance, SaaS WAF instance, CLB instance, or a VPC endpoint.
Limitation (GWLB VPC endpoint): Only an EIP in the bound state can be associated with a GWLB-type VPC endpoint through this field, enabling more advanced networking scenarios.
Mutual exclusivity: This field conflicts with `network_interface_id` and `private_ip`. Only one association target can be specified per request.
* `network_interface_id` - (Optional, String, ForceNew) Indicates the network interface id like `eni-xxxxxx`. This field is conflict with `instance_id`.
* `private_ip` - (Optional, String, ForceNew) Indicates an IP belongs to the `network_interface_id`. This field is conflict with `instance_id`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `cdc_id` - ID of the dedicated cluster.

## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [AssociateAddress](https://www.tencentcloud.com/document/product/215/16700) - This API is used to bind an [Elastic IP](https://intl.cloud.tencent.com/document/product/213/1941?from_cn_redirect=1) (EIP for short) to the specified private IP of an instance or ENI.
* Create: [DescribeAddresses](https://www.tencentcloud.com/document/product/215/16702) - This API (DescribeAddresses) is used to query the information of one or multiple [Elastic IPs](https://intl.cloud.tencent.com/document/product/213/1941?from_cn_redirect=1).
* Read: [DescribeAddresses](https://www.tencentcloud.com/document/product/215/16702) - This API (DescribeAddresses) is used to query the information of one or multiple [Elastic IPs](https://intl.cloud.tencent.com/document/product/213/1941?from_cn_redirect=1).
* Delete: [DescribeAddresses](https://www.tencentcloud.com/document/product/215/16702) - This API (DescribeAddresses) is used to query the information of one or multiple [Elastic IPs](https://intl.cloud.tencent.com/document/product/213/1941?from_cn_redirect=1).
* Delete: [DescribeNatGateways](https://www.tencentcloud.com/document/product/215/34752) - This API is used to query NAT gateways.
* Delete: [DescribeTaskResult](https://www.tencentcloud.com/document/product/215/34772) - This API is used to query the EIP async job execution results.
* Delete: [DisassociateAddress](https://www.tencentcloud.com/document/product/215/16703) - This API is used to unbind an [Elastic IP](https://intl.cloud.tencent.com/document/product/213/1941?from_cn_redirect=1) (EIP for short).
* Delete: [DisassociateNatGatewayAddress](https://www.tencentcloud.com/document/product/215/34751) - This API is used to unbind an EIP from a NAT gateway.

## Import

Eip association can be imported using the id, e.g.

```
$ terraform import tencentcloud_eip_association.example eip-41s6jwy4::ins-0cxjwrog
```

