---
subcategory: "Anti-DDoS(DayuV2)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_antiddos_ddos_black_white_ip"
sidebar_current: "docs-tencentcloud-resource-antiddos_ddos_black_white_ip"
description: |-
  Provides a resource to create a antiddos ddos black white ip
---

# tencentcloud_antiddos_ddos_black_white_ip

Provides a resource to create a antiddos ddos black white ip

## Example Usage

```hcl
resource "tencentcloud_antiddos_ddos_black_white_ip" "ddos_black_white_ip" {
  instance_id = "bgp-xxxxxx"
  ip          = "1.2.3.5"
  mask        = 0
  type        = "black"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String) instance id.
* `ip` - (Required, String) ip list.
* `mask` - (Required, Int) ip mask.
* `type` - (Required, String) ip type, black: black ip list, white: white ip list.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.


## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [CreateDDoSBlackWhiteIpList](https://www.tencentcloud.com/search/?wd=CreateDDoSBlackWhiteIpList&type=document&lang=en)
* Read: [DescribeDDoSBlackWhiteIpList](https://www.tencentcloud.com/search/?wd=DescribeDDoSBlackWhiteIpList&type=document&lang=en)
* Update: [ModifyDDoSBlackWhiteIpList](https://www.tencentcloud.com/search/?wd=ModifyDDoSBlackWhiteIpList&type=document&lang=en)
* Delete: [DeleteDDoSBlackWhiteIpList](https://www.tencentcloud.com/search/?wd=DeleteDDoSBlackWhiteIpList&type=document&lang=en)

## Import

antiddos ddos_black_white_ip can be imported using the id, e.g.

```
terraform import tencentcloud_antiddos_ddos_black_white_ip.ddos_black_white_ip ${instanceId}#${ip}
```

