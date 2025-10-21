---
subcategory: "PrivateDNS"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_private_dns_private_zone_list"
sidebar_current: "docs-tencentcloud-datasource-private_dns_private_zone_list"
description: |-
  Use this data source to query detailed information of Private Dns private zone list
---

# tencentcloud_private_dns_private_zone_list

Use this data source to query detailed information of Private Dns private zone list

## Example Usage

### Query All private zones:

```hcl
data "tencentcloud_private_dns_private_zone_list" "example" {}
```

### Query private zones by ZoneId:

```hcl
data "tencentcloud_private_dns_private_zone_list" "example" {
  filters {
    name   = "ZoneId"
    values = ["zone-6xg5xgky1"]
  }
}
```

### Query private zones by Domain:

```hcl
data "tencentcloud_private_dns_private_zone_list" "example" {
  filters {
    name   = "Domain"
    values = ["domain.com"]
  }
}
```

### Query private zones by Vpc:

```hcl
data "tencentcloud_private_dns_private_zone_list" "example" {
  filters {
    name   = "Vpc"
    values = ["vpc-axrsmmrv"]
  }
}
```

## Argument Reference

The following arguments are supported:

* `filters` - (Optional, List) filters.
* `result_output_file` - (Optional, String) Used to save results.

The `filters` object supports the following:

* `name` - (Required, String) name.
* `values` - (Required, Set) values.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `private_zone_set` - Private Zone Set.
  * `account_vpc_set` - VPC list of bound associated accounts.
    * `region` - Region.
    * `uin` - uin.
    * `uniq_vpc_id` - Vpc Id.
  * `cname_speedup_status` - CNAME acceleration status: enabled: ENABLED, off, DISABLED.
  * `created_on` - Create time.
  * `deleted_vpc_set` - List of deleted VPCs.
    * `region` - Region.
    * `uniq_vpc_id` - Vpc Id.
  * `dns_forward_status` - Domain name recursive resolution status: enabled: ENABLED, disabled, DISABLED.
  * `domain` - Domain.
  * `end_point_name` - End point name.
  * `forward_address` - Forwarded address.
  * `forward_rule_name` - Forwarding rule name.
  * `forward_rule_type` - Forwarding rule type: from cloud to cloud, DOWN; From cloud to cloud, UP, currently only supports DOWN.
  * `is_custom_tld` - Custom TLD.
  * `owner_uin` - Owner Uin.
  * `record_count` - Record count.
  * `remark` - Remark.
  * `status` - Private domain bound VPC status, not associated with vpc: SUSPEND, associated with VPC: ENABLED, associated with VPC failed: FAILED.
  * `tags` - tags.
    * `tag_key` - tag key.
    * `tag_value` - tag value.
  * `updated_on` - Update time.
  * `vpc_set` - Vpc list.
    * `region` - Region.
    * `uniq_vpc_id` - Vpc Id.
  * `zone_id` - PrivateZone ID.


