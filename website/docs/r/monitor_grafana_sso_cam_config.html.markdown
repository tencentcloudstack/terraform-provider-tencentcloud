---
subcategory: "TencentCloud Managed Service for Grafana(TCMG)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_monitor_grafana_sso_cam_config"
sidebar_current: "docs-tencentcloud-resource-monitor_grafana_sso_cam_config"
description: |-
  Provides a resource to create a monitor grafana_sso_cam_config
---

# tencentcloud_monitor_grafana_sso_cam_config

Provides a resource to create a monitor grafana_sso_cam_config

## Example Usage

```hcl
resource "tencentcloud_monitor_grafana_sso_cam_config" "grafana_sso_cam_config" {
  instance_id          = "grafana-dp2hnnfa"
  enable_sso_cam_check = false
}
```

## Argument Reference

The following arguments are supported:

* `enable_sso_cam_check` - (Required, Bool) Whether to enable the CAM authorization: `true` for enabling; `false` for disabling.
* `instance_id` - (Required, String, ForceNew) Grafana instance ID.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.


## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [EnableSSOCamCheck](https://www.tencentcloud.com/document/product/248/49751) - This API is used to set whether to enable CAM authentication during SSO.
* Read: [DescribeGrafanaInstances](https://www.tencentcloud.com/document/product/248/49759) - This API is used to list all Grafana instances under a user account.
* Update: [EnableSSOCamCheck](https://www.tencentcloud.com/document/product/248/49751) - This API is used to set whether to enable CAM authentication during SSO.

## Import

monitor grafana_sso_cam_config can be imported using the id, e.g.

```
terraform import tencentcloud_monitor_grafana_sso_cam_config.grafana_sso_cam_config instance_id
```

