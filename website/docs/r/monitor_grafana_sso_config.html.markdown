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

* Create: [DescribeGrafanaInstances](https://cloud.tencent.com/document/api/248/79208) - 列出用户所有的 Grafana 服务
* Create: [EnableGrafanaSSO](https://cloud.tencent.com/document/api/248/79200) - 设置 Grafana 单点登录
* Read: [DescribeGrafanaInstances](https://cloud.tencent.com/document/api/248/79208) - 列出用户所有的 Grafana 服务
* Update: [DescribeGrafanaInstances](https://cloud.tencent.com/document/api/248/79208) - 列出用户所有的 Grafana 服务
* Update: [EnableGrafanaSSO](https://cloud.tencent.com/document/api/248/79200) - 设置 Grafana 单点登录

## Import

monitor grafana_sso_config can be imported using the id, e.g.

```
terraform import tencentcloud_monitor_grafana_sso_config.grafana_sso_config instance_id
```

