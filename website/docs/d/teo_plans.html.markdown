---
subcategory: "TencentCloud EdgeOne(TEO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_teo_plans"
sidebar_current: "docs-tencentcloud-datasource-teo_plans"
description: |-
  Use this data source to query detailed information of TEO plans
---

# tencentcloud_teo_plans

Use this data source to query detailed information of TEO plans

## Example Usage

### Query all plans

```hcl
data "tencentcloud_teo_plans" "example" {}
```

### Query plans by filters

```hcl
data "tencentcloud_teo_plans" "example" {
  order     = "expire-time"
  direction = "desc"
  filters {
    name = "plan-id"
    values = [
      "edgeone-2o1xvpmq7nn",
      "edgeone-2mezmk9s2xdx"
    ]
  }

  filters {
    name = "plan-type"
    values = [
      "plan-trial",
      "plan-personal",
      "plan-basic",
      "plan-standard",
      "plan-enterprise"
    ]
  }

  filters {
    name = "area"
    values = [
      "overseas",
      "mainland",
      "global"
    ]
  }

  filters {
    name = "status"
    values = [
      "normal",
      "expiring-soon",
      "expired",
      "isolated"
    ]
  }
}
```

## Argument Reference

The following arguments are supported:

* `direction` - (Optional, String) Sorting direction, the possible values are: <li>asc: sort from small to large; </li><li>desc: sort from large to small. </li>If not filled in, the default value desc will be used.
* `filters` - (Optional, List) Filter conditions, the upper limit of Filters. Values is 20. The detailed filtering conditions are as follows: <li>plan-type<br>Filter according to [<strong>Package Type</strong>]. <br>Optional types are: <br>plan-trial: Trial Package; <br>plan-personal: Personal Package; <br>plan-basic: Basic Package; <br>plan-standard: Standard Package; <br>plan-enterprise: Enterprise Package. </li><li>plan-id<br>Filter according to [<strong>Package ID</strong>]. The package ID is in the form of: edgeone-268z103ob0sx.</li><li>area<br>Filter according to [<strong>Package Acceleration Region</strong>]. </li>Service area, optional types are: <br>mainland: Mainland China; <br>overseas: Global (excluding Mainland China); <br>global: Global (including Mainland China).<br><li>status<br>Filter by [<strong>Package Status</strong>].<br>The available statuses are:<br>normal: normal status;<br>expiring-soon: about to expire;<br>expired: expired;<br>isolated: isolated.</li>.
* `order` - (Optional, String) Sorting field, the values are: <li> enable-time: effective time; </li><li> expire-time: expiration time. </li> If not filled in, the default value enable-time will be used.
* `result_output_file` - (Optional, String) Used to save results.

The `filters` object supports the following:

* `name` - (Required, String) Filter name.
* `values` - (Required, Set) Filter value.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `plans` - Plan list.
  * `acc_traffic_capacity` - The content acceleration traffic specifications in the package, unit: byte.
  * `area` - Service area, the values are: <li>mainland: Mainland China; </li><li>overseas: Worldwide (excluding Mainland China); </li><li>global: Worldwide (including Mainland China).</li>.
  * `bindable` - Whether the package allows binding of new sites, the values are: <li>true: allows binding of new sites; </li><li>false: does not allow binding of new sites.</li>.
  * `cross_mlc_traffic_capacity` - The optimized traffic specifications of the Chinese mainland network in the package, unit: bytes.
  * `ddos_traffic_capacity` - DDoS protection traffic specifications within the package, unit: bytes.
  * `enabled_time` - The package effective time.
  * `expired_time` - The expiration date of the package.
  * `features` - The functions supported by the package have the following values: <li>ContentAcceleration: content acceleration function; </li><li>SmartAcceleration: smart acceleration function; </li><li>L4: four-layer acceleration function; </li><li>Waf: advanced web protection; </li><li>QUIC: QUIC function; </li><li>CrossMLC: Chinese mainland network optimization function; </li><li>ProcessMedia: media processing function; </li><li>L4DDoS: four-layer DDoS protection function; </li>L7DDoS function will only have one of the following specifications <li>L7DDoS.CM30G; seven-layer DDoS protection function - Chinese mainland 30G minimum bandwidth specification; </li><li>L7DDoS.CM60G; seven-layer DDoS protection function - Chinese mainland 60G minimum bandwidth specification; </li> <li>L7DDoS.CM100G; Layer 7 DDoS protection function - 100G guaranteed bandwidth for mainland China;</li><li>L7DDoS.Anycast300G; Layer 7 DDoS protection function - 300G guaranteed bandwidth for Anycast outside mainland China;</li><li>L7DDoS.AnycastUnlimited; Layer 7 DDoS protection function - unlimited full protection for Anycast outside mainland China;</li><li>L7DDoS.CM30G_Anycast300G; Layer 7 DDoS protection function - 30G guaranteed bandwidth for mainland China </li><li>L7DDoS.CM60G_Anycast300G; Layer 7 DDoS protection function - 60G guaranteed bandwidth in mainland China, 300G guaranteed bandwidth in anycast outside mainland China; </li><li>L7DDoS.CM100G_Anycast300G; Layer 7 DDoS protection function - 100G guaranteed bandwidth in mainland China, 300G guaranteed bandwidth in anycast outside mainland China; </li><li>L7DDoS.CM30G_AnycastUnlimited d; Layer 7 DDoS protection function - 30G guaranteed bandwidth in mainland China, unlimited Anycast protection outside mainland China; </li><li>L7DDoS.CM60G_AnycastUnlimited; Layer 7 DDoS protection function - 60G guaranteed bandwidth in mainland China, unlimited Anycast protection outside mainland China; </li><li>L7DDoS.CM100G_AnycastUnlimited; Layer 7 DDoS protection function - 100G guaranteed bandwidth in mainland China, unlimited Anycast protection outside mainland China; </li>.
  * `l4_traffic_capacity` - Layer 4 acceleration traffic specifications within the package, unit: byte.
  * `pay_mode` - Payment type, possible values: <li>0: post-payment; </li><li>1: pre-payment.</li>.
  * `plan_id` - Plan ID.
  * `plan_type` - Plan type. Possible values are: <li>plan-trial: Trial plan; </li><li>plan-personal: Personal plan; </li><li>plan-basic: Basic plan; </li><li>plan-standard: Standard plan; </li><li>plan-enterprise-v2: Enterprise plan; </li><li>plan-enterprise-model-a: Enterprise Model A plan. </li><li>plan-enterprise: Old Enterprise plan. </li>.
  * `sec_request_capacity` - The number of secure requests in the package, unit: times.
  * `sec_traffic_capacity` - The security flow specification within the package, unit: byte.
  * `smart_request_capacity` - The number of intelligent acceleration requests in the package, unit: times.
  * `smart_traffic_capacity` - Smart acceleration traffic specifications within the package, unit: byte.
  * `status` - Package status, the values are: <li>normal: normal status; </li><li>expiring-soon: about to expire; </li><li>expired: expired; </li><li>isolated: isolated; </li><li>overdue-isolated: overdue isolated.</li>.
  * `vau_capacity` - VAU specifications in the package, unit: piece.
  * `zones_info` - Site information bound to the package, including site ID, site name, and site status.
    * `paused` - Whether the site is disabled. The possible values are: <li>false: not disabled; </li><li>true: disabled.</li>.
    * `zone_id` - Zone ID.
    * `zone_name` - Zone name.


