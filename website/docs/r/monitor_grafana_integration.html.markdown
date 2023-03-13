---
subcategory: "TencentCloud Managed Service for Grafana(TCMG)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_monitor_grafana_integration"
sidebar_current: "docs-tencentcloud-resource-monitor_grafana_integration"
description: |-
  Provides a resource to create a monitor grafanaIntegration
---

# tencentcloud_monitor_grafana_integration

Provides a resource to create a monitor grafanaIntegration

## Example Usage

```hcl
resource "tencentcloud_monitor_grafana_integration" "grafanaIntegration" {
  instance_id = "grafana-50nj6v00"
  kind        = "tencentcloud-monitor-app"
  content     = "{\"kind\":\"tencentcloud-monitor-app\",\"spec\":{\"dataSourceSpec\":{\"authProvider\":{\"__anyOf\":\"使用密钥\",\"useRole\":true,\"secretId\":\"arunma@tencent.com\",\"secretKey\":\"12345678\"},\"name\":\"uint-test\"},\"grafanaSpec\":{\"organizationIds\":[]}}}"
}

resource "tencentcloud_monitor_grafana_integration" "grafanaIntegration_update" {
  content     = "{\"id\":\"integration-9st6kqz6\",\"kind\":\"tencentcloud-monitor-app\",\"spec\":{\"dataSourceSpec\":{\"name\":\"uint-test3\",\"authProvider\":{\"secretId\":\"ROLE\",\"useRole\":true,\"__anyOf\":\"使用服务角色\"}},\"grafanaSpec\":{\"organizationIds\":[]}}}"
  instance_id = "grafana-50nj6v00"
  kind        = "tencentcloud-monitor-app"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String) grafana instance id.
* `content` - (Optional, String) generated json string of given integration json schema.
* `description` - (Optional, String) integration desc.
* `kind` - (Optional, String) integration json schema kind.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `integration_id` - integration id.


