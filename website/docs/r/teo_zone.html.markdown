---
subcategory: "TencentCloud EdgeOne(TEO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_teo_zone"
sidebar_current: "docs-tencentcloud-resource-teo_zone"
description: |-
  Provides a resource to create a teo zone
---

# tencentcloud_teo_zone

Provides a resource to create a teo zone

## Example Usage

### Basic Usage

```hcl
resource "tencentcloud_teo_zone" "zone" {
  zone_name       = "tf-teo.com"
  type            = "partial"
  area            = "overseas"
  alias_zone_name = "teo-test"
  paused          = false
  plan_id         = "edgeone-2kfv1h391n6w"
  tags = {
    "createdBy" = "terraform"
  }
}
```

### Enable Version Control Mode

```hcl
resource "tencentcloud_teo_zone" "zone_with_version_control" {
  zone_name       = "tf-teo-version.com"
  type            = "partial"
  area            = "overseas"
  alias_zone_name = "teo-version-test"
  paused          = false
  plan_id         = "edgeone-2kfv1h391n6w"

  work_mode_infos {
    config_group_type = "l7_acceleration"
    work_mode         = "version_control"
  }

  tags = {
    "createdBy" = "terraform"
  }
}
```

## Argument Reference

The following arguments are supported:

* `area` - (Required, String) When the Type value is partial/full, the acceleration region of the L7 domain name. The following are the values of this parameter, and the default value is overseas if not filled in. When the Type value is noDomainAccess, please leave this value empty:
  - global: Global availability zone.
  - mainland: Chinese mainland availability zone.
  - overseas: Global availability zone (excluding Chinese mainland).
* `plan_id` - (Required, String, ForceNew) The target Plan ID to be bound. When you have an existing Plan in your account, you can fill in this parameter to directly bind the site to the Plan. If you do not have a Plan that can be bound at the moment, please go to the console to purchase a Plan to complete the site creation.
* `type` - (Required, String) Site access type. The value of this parameter is as follows, and the default is partial if not filled in:partial: CNAME access; full: NS access; noDomainAccess: No domain access.
* `zone_name` - (Required, String, ForceNew) Site name. When accessing CNAME/NS, please pass the second-level domain (example.com) as the site name; when accessing without a domain name, please leave this value empty.
* `alias_zone_name` - (Optional, String) Alias site identifier. Limit the input to a combination of numbers, English, - and _, within 20 characters. For details, refer to the alias site identifier. If there is no such usage scenario, leave this field empty.
* `paused` - (Optional, Bool) Indicates whether the site is disabled.
* `tags` - (Optional, Map) Tag description list.
* `work_mode_infos` - (Optional, List) Configuration group work mode. Each configuration module of the site can enable version control mode or immediate effect mode according to the configuration group dimension. For details, please refer to [Version Management](https://cloud.tencent.com/document/product/1552/113690).

The `work_mode_infos` object supports the following:

* `config_group_type` - (Required, String) Configuration group type. Valid values: `l7_acceleration`: L7 acceleration configuration group; `edge_functions`: Edge functions configuration group.
* `work_mode` - (Required, String) Work mode. Valid values: `immediate_effect`: Immediate effect mode; `version_control`: Version control mode.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `name_servers` - NS list allocated by Tencent Cloud.
* `ownership_verification` - Ownership verification information. Note: This field may return null, indicating that no valid value can be obtained.
  * `dns_verification` - CNAME access, using DNS to resolve the information required for authentication. For details, please refer to [Site/Domain Name Ownership Verification ](https://cloud.tencent.com/document/product/1552/70789#7af6ecf8-afca-4e35-8811-b5797ed1bde5). Note: This field may return null, indicating that no valid value can be obtained.
    * `record_type` - Record type.
    * `record_value` - Record the value.
    * `subdomain` - Host record.
* `status` - Site status. Valid values: `active`: NS is switched; `pending`: NS is not switched; `moved`: NS is moved; `deactivated`: this site is blocked.


## Import

teo zone can be imported using the id, e.g.
```
terraform import tencentcloud_teo_zone.zone zone_id
```

