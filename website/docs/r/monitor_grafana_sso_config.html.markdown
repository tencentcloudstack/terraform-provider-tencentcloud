---
subcategory: "TencentCloud Managed Service for Grafana(TCMG)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_monitor_grafana_sso_config"
sidebar_current: "docs-tencentcloud-resource-monitor_grafana_sso_config"
description: |-
  Provides a resource to create a monitor grafana_sso_config
---

# tencentcloud_monitor_grafana_sso_config

Provides a resource to create a monitor grafana_sso_config

## Example Usage

```hcl
resource "tencentcloud_monitor_grafana_sso_config" "grafana_sso_config" {
  instance_id = "grafana-dp2hnnfa"
  enable_sso  = false
}
```

## Argument Reference

The following arguments are supported:

* `enable_sso` - (Required, Bool) Whether to enable SSO: `true` for enabling; `false` for disabling.
* `instance_id` - (Required, String, ForceNew) Grafana instance ID.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.


## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [DescribeGrafanaInstances](https://www.tencentcloud.com/document/product/248/49759) - This API is used to list all Grafana instances under a user account.
* Create: [EnableGrafanaSSO](https://www.tencentcloud.com/document/product/248/49752) - This API is used to set the Grafana SSO through a Tencent Cloud account.
* Read: [DescribeGrafanaInstances](https://www.tencentcloud.com/document/product/248/49759) - This API is used to list all Grafana instances under a user account.
* Update: [DescribeGrafanaInstances](https://www.tencentcloud.com/document/product/248/49759) - This API is used to list all Grafana instances under a user account.
* Update: [EnableGrafanaSSO](https://www.tencentcloud.com/document/product/248/49752) - This API is used to set the Grafana SSO through a Tencent Cloud account.

## Import

monitor grafana_sso_config can be imported using the id, e.g.

```
terraform import tencentcloud_monitor_grafana_sso_config.grafana_sso_config instance_id
```

