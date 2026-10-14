---
subcategory: "Virtual Private Cloud(VPC)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_reserve_ip_address"
sidebar_current: "docs-tencentcloud-resource-reserve_ip_address"
description: |-
  Provides a resource to create a VPC reserve ip address
---

# tencentcloud_reserve_ip_address

Provides a resource to create a VPC reserve ip address

## Example Usage

```hcl
resource "tencentcloud_reserve_ip_address" "example" {
  vpc_id      = "vpc-i5yyodl9"
  subnet_id   = "subnet-hhi88a58"
  ip_address  = "10.0.30.27"
  name        = "tf-example"
  description = "remark."
  tags = {
    createBy = "Terraform"
  }
}
```

## Argument Reference

The following arguments are supported:

* `vpc_id` - (Required, String) VPC unique ID.
* `description` - (Optional, String) The IP description is retained on the intranet.
* `ip_address` - (Optional, String) Specify the reserved IP address of the intranet for which the IP application is requested.
* `name` - (Optional, String) The IP name is reserved for the intranet.
* `subnet_id` - (Optional, String) Subnet ID.
* `tags` - (Optional, Map) Tags.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `created_time` - Created time.
* `ip_type` - Ip type for product application.
* `reserve_ip_id` - Reserve ip ID.
* `resource_id` - The intranet retains the resource instance ID bound to the IPs.
* `state` - Binding status.

## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [CreateReserveIpAddresses](https://cloud.tencent.com/document/api/215/112867) - This API is used to create a reserved private IP address.
* Create: [ModifyResourceTags](https://cloud.tencent.com/document/api/651/35322) - This API is used to modify all tags associated with a resource.
* Read: [DescribeReserveIpAddresses](https://cloud.tencent.com/document/api/215/112865) - This API is used to query reserved private IP addresses.
* Read: [DescribeResourceTagsByResourceIds](https://cloud.tencent.com/document/api/651/35319) - This API is used to query the tag key-value pairs associated with an existing resource.
* Update: [ModifyReserveIpAddress](https://cloud.tencent.com/document/api/215/112864) - This API is used to modify a reserved private IP address.
* Update: [ModifyResourceTags](https://cloud.tencent.com/document/api/651/35322) - This API is used to modify all tags associated with a resource.
* Delete: [DeleteReserveIpAddresses](https://cloud.tencent.com/document/api/215/112866) - This API is used to delete a reserved private IP address.

## Import

VPC reserve ip address can be imported using the vpcId#reserveIpId, e.g.

```
terraform import tencentcloud_reserve_ip_address.example vpc-i5yyodl9#rsvip-cz3ces44
```

