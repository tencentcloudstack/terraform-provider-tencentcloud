---
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_eip"
sidebar_current: "docs-tencentcloud-datasource-eip"
description: |-
  Provides an available EIP for the user.
---

# tencentcloud_eip

The EIP data source fetch proper EIP from user's EIP pool.

~> **NOTE:** It has been deprecated and replaced by tencentcloud_eips.

## Example Usage

```hcl
data "tencentcloud_eip" "my_eip" {
  filter {
    name   = "address-status"
    values = ["UNBIND"]
  }
}
```

## Argument Reference

 * `filter` - (Optional) One or more name/value pairs to filter off of. There are several valid keys:  `address-id`,`address-name`,`address-ip`. For a full reference, check out [DescribeImages in the TencentCloud API reference](https://intl.cloud.tencent.com/document/api/213/9451#filter).

## Attributes Reference

 * `id` - An EIP id indicate the uniqueness of a certain EIP,  which can be used for instance binding or network interface binding.
 * `public_ip` - An public IP address for the EIP.
 * `status` - The status of the EIP, there are several status like `BIND`, `UNBIND`, and `BIND_ENI`. For a full reference, check out [DescribeImages in the TencentCloud API reference](https://intl.cloud.tencent.com/document/api/213/9452#eip_state).