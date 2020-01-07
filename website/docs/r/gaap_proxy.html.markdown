---
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_gaap_proxy"
sidebar_current: "docs-tencentcloud-resource-gaap_proxy"
description: |-
  Provides a resource to create a GAAP proxy.
---

# tencentcloud_gaap_proxy

Provides a resource to create a GAAP proxy.

## Example Usage

```hcl
resource "tencentcloud_gaap_proxy" "foo" {
  name              = "ci-test-gaap-proxy"
  bandwidth         = 10
  concurrent        = 2
  access_region     = "SouthChina"
  realserver_region = "NorthChina"

  tags = {
    test = "test"
  }
}
```

## Argument Reference

The following arguments are supported:

* `access_region` - (Required, ForceNew) Access region of the GAAP proxy, the available values include `NorthChina`, `EastChina`, `SouthChina`, `SouthwestChina`, `Hongkong`, `SL_TAIWAN`, `SoutheastAsia`, `Korea`, `SL_India`, `SL_Australia`, `Europe`, `SL_UK`, `SL_SouthAmerica`, `NorthAmerica`, `SL_MiddleUSA`, `Canada`, `SL_VIET`, `WestIndia`, `Thailand`, `Virginia`, `Russia`, `Japan` and `SL_Indonesia`.
* `bandwidth` - (Required) Maximum bandwidth of the GAAP proxy, unit is Mbps, the available values include `10`, `20`, `50`, `100`, `200`, `500` and `1000`.
* `concurrent` - (Required) Maximum concurrency of the GAAP proxy, unit is 10k, the available values include `2`, `5`, `10`, `20`, `30`, `40`, `50`, `60`, `70`, `80`, `90` and `100`.
* `name` - (Required) Name of the GAAP proxy, the maximum length is 30.
* `realserver_region` - (Required, ForceNew) Region of the GAAP realserver, the available values include `NorthChina`, `EastChina`, `SouthChina`, `SouthwestChina`, `Hongkong`, `SL_TAIWAN`, `SoutheastAsia`, `Korea`, `SL_India`, `SL_Australia`, `Europe`, `SL_UK`, `SL_SouthAmerica`, `NorthAmerica`, `SL_MiddleUSA`, `Canada`, `SL_VIET`, `WestIndia`, `Thailand`, `Virginia`, `Russia`, `Japan` and `SL_Indonesia`.
* `enable` - (Optional) Indicates whether GAAP proxy is enabled, default value is `true`.
* `project_id` - (Optional) ID of the project within the GAAP proxy, '0' means is default project.
* `tags` - (Optional) Tags of the GAAP proxy.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `create_time` - Creation time of the GAAP proxy.
* `domain` - Access domain of the GAAP proxy.
* `forward_ip` - Forwarding IP of the GAAP proxy.
* `ip` - Access IP of the GAAP proxy.
* `scalable` - Indicates whether GAAP proxy can scalable.
* `status` - Status of the GAAP proxy.
* `support_protocols` - Supported protocols of the GAAP proxy.


## Import

GAAP proxy can be imported using the id, e.g.

```
  $ terraform import tencentcloud_gaap_proxy.foo link-11112222
```

