---
subcategory: "Global Application Acceleration(GAAP)"
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

* `access_region` - (Required, String, ForceNew) Access region of the GAAP proxy. Valid value: `Hongkong`, `SoutheastAsia`, `Korea`, `Europe`, `NorthAmerica`, `Canada`, `WestIndia`, `Thailand`, `Virginia`, `Japan`, `Taipei`, `SL_AZURE_NorthUAE`, `SL_AZURE_EastAUS`, `SL_AZURE_NorthCentralUSA`, `SL_AZURE_SouthIndia`, `SL_AZURE_SouthBrazil`, `SL_AZURE_NorthZAF`, `SL_AZURE_SoutheastAsia`, `SL_AZURE_CentralFrance`, `SL_AZURE_SouthEngland`, `SL_AZURE_EastUS`, `SL_AZURE_WestUS`, `SL_AZURE_SouthCentralUSA`, `Jakarta`, `Beijing`, `Shanghai`, `Guangzhou`, `Chengdu`, `SL_AZURE_NorwayEast`, `Chongqing`, `Nanjing`, `SaoPaulo`, `SL_AZURE_JapanEast`, `Changsha`, `Xian`, `Wuhan`, `Fuzhou`, `Shenyang`, `Zhengzhou`, `Jinan`, `Hangzhou`, `Shijiazhuang`, `Hefei`.
* `bandwidth` - (Required, Int) Maximum bandwidth of the GAAP proxy, unit is Mbps. Valid value: `10`, `20`, `50`, `100`, `200`, `500`, `1000`, `2000`, `5000` and `10000`. To set `2000`, `5000` or `10000`, you need to apply for a whitelist from Tencent Cloud.
* `concurrent` - (Required, Int) Maximum concurrency of the GAAP proxy, unit is 10k. Valid value: `2`, `5`, `10`, `20`, `30`, `40`, `50`, `60`, `70`, `80`, `90`, `100`, `150`, `200`, `250` and `300`. To set `150`, `200`, `250` or `300`, you need to apply for a whitelist from Tencent Cloud.
* `name` - (Required, String) Name of the GAAP proxy, the maximum length is 30.
* `realserver_region` - (Required, String, ForceNew) Region of the GAAP realserver. Valid value: `Hongkong`, `SoutheastAsia`, `Korea`, `Europe`, `NorthAmerica`, `Canada`, `WestIndia`, `Thailand`, `Virginia`, `Japan`, `Taipei`, `SL_AZURE_NorthUAE`, `SL_AZURE_EastAUS`, `SL_AZURE_NorthCentralUSA`, `SL_AZURE_SouthIndia`, `SL_AZURE_SouthBrazil`, `SL_AZURE_NorthZAF`, `SL_AZURE_SoutheastAsia`, `SL_AZURE_CentralFrance`, `SL_AZURE_SouthEngland`, `SL_AZURE_EastUS`, `SL_AZURE_WestUS`, `SL_AZURE_SouthCentralUSA`, `Jakarta`, `Beijing`, `Shanghai`, `Guangzhou`, `Chengdu`, `SL_AZURE_NorwayEast`, `Chongqing`, `Nanjing`, `SaoPaulo`, `SL_AZURE_JapanEast`.
* `enable` - (Optional, Bool) Indicates whether GAAP proxy is enabled, default value is `true`.
* `network_type` - (Optional, String, ForceNew) Network type. `normal`: regular BGP, `cn2`: boutique BGP, `triple`: triple play.
* `project_id` - (Optional, Int) ID of the project within the GAAP proxy, `0` means is default project.
* `tags` - (Optional, Map) Tags of the GAAP proxy. Tags that do not exist are not created automatically.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `create_time` - Creation time of the GAAP proxy.
* `domain` - Access domain of the GAAP proxy.
* `forward_ip` - Forwarding IP of the GAAP proxy.
* `ip` - Access IP of the GAAP proxy.
* `scalable` - Indicates whether GAAP proxy can scalable.
* `status` - Status of the GAAP proxy.
* `support_protocols` - Supported protocols of the GAAP proxy.

## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [CloseProxies](https://cloud.tencent.com/document/api/608/36973) - This API (CloseProxies) is used to disable connections. If disabled, no traffic will be generated, but the basic configuration fee will still be incurred.
* Create: [CreateProxy](https://cloud.tencent.com/document/api/608/36971) - This API is used to create/replicate an acceleration connection with the specified configuration. To replicate a connection, the basic configuration parameters need to be set for the new connection, and `ClonedProxyId` is needed to identify the replicated connection.
* Create: [DescribeProxies](https://cloud.tencent.com/document/api/608/36963) - This API (DescribeProxies) is used to query the connection instance list.
* Read: [DescribeProxies](https://cloud.tencent.com/document/api/608/36963) - This API (DescribeProxies) is used to query the connection instance list.
* Update: [CloseProxies](https://cloud.tencent.com/document/api/608/36973) - This API (CloseProxies) is used to disable connections. If disabled, no traffic will be generated, but the basic configuration fee will still be incurred.
* Update: [DescribeProxies](https://cloud.tencent.com/document/api/608/36963) - This API (DescribeProxies) is used to query the connection instance list.
* Update: [ModifyProxiesAttribute](https://cloud.tencent.com/document/api/608/36956) - This API (ModifyProxiesAttribute) is used to modify instance attributes (currently only supports modifying connection name).
* Update: [ModifyProxiesProject](https://cloud.tencent.com/document/api/608/36955) - This API (ModifyProxiesProject) is used to modify the project to which a connection belongs.
* Update: [ModifyProxyConfiguration](https://cloud.tencent.com/document/api/608/36954) - This API (ModifyProxyConfiguration) is used to modify connection configurations. You can expand or reduce the capacity based on current business requirements. It only supports connections with a Scalarable of 1. Scalarable can be obtained via DescribeProxies API.
* Update: [ModifyResourceTags](https://cloud.tencent.com/document/api/651/35322) - This API is used to modify all tags associated with a resource.
* Update: [OpenProxies](https://cloud.tencent.com/document/api/608/36952) - This API (OpenProxies) is used to enable one or more connections.
* Delete: [DescribeProxies](https://cloud.tencent.com/document/api/608/36963) - This API (DescribeProxies) is used to query the connection instance list.
* Delete: [DestroyProxies](https://cloud.tencent.com/document/api/608/36958) - This API (DestroyProxies) is used to terminate a connection. If terminated, no fees will be incurred.

## Import

GAAP proxy can be imported using the id, e.g.

```
  $ terraform import tencentcloud_gaap_proxy.foo link-11112222
```

