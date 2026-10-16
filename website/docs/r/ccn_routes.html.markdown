---
subcategory: "Cloud Connect Network(CCN)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_ccn_routes"
sidebar_current: "docs-tencentcloud-resource-ccn_routes"
description: |-
  Provides a resource to create a vpc ccn_routes switch
---

# tencentcloud_ccn_routes

Provides a resource to create a vpc ccn_routes switch

## Example Usage

```hcl
resource "tencentcloud_ccn_routes" "example" {
  ccn_id   = "ccn-gr7nynbd"
  route_id = "ccnr-5uhewx1s"
  switch   = "off"
}
```

## Argument Reference

The following arguments are supported:

* `ccn_id` - (Required, String, ForceNew) CCN Instance ID.
* `route_id` - (Required, String, ForceNew) CCN Route Id List.
* `switch` - (Required, String) `on`: Enable, `off`: Disable.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.


## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [DisableCcnRoutes](https://www.tencentcloud.com/document/api/215/34785)
* Create: [EnableCcnRoutes](https://www.tencentcloud.com/document/api/215/34784)
* Read: [DescribeCcnRoutes](https://www.tencentcloud.com/document/api/215/34788)
* Update: [DisableCcnRoutes](https://www.tencentcloud.com/document/api/215/34785)
* Update: [EnableCcnRoutes](https://www.tencentcloud.com/document/api/215/34784)

## Import

vpc ccn_routes can be imported using the id, e.g.

```
terraform import tencentcloud_ccn_routes.ccn_routes ccn-gr7nynbd#ccnr-5uhewx1s
```

