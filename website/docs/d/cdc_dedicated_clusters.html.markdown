---
subcategory: "CDC"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cdc_dedicated_clusters"
sidebar_current: "docs-tencentcloud-datasource-cdc_dedicated_clusters"
description: |-
  Use this data source to query detailed information of CDC dedicated clusters
---

# tencentcloud_cdc_dedicated_clusters

Use this data source to query detailed information of CDC dedicated clusters

## Example Usage

### Query all dedicated clusters

```hcl
data "tencentcloud_cdc_dedicated_clusters" "example" {}
```

### Query dedicated clusters by filters

```hcl
data "tencentcloud_cdc_dedicated_clusters" "example" {
  name = "tf-example"
}

data "tencentcloud_cdc_dedicated_clusters" "example" {
  dedicated_cluster_ids = [
    "cluster-aiaui7ei",
    "cluster-262n63e8"
  ]
}

data "tencentcloud_cdc_dedicated_clusters" "example" {
  zones = [
    "ap-guangzhou-2",
    "ap-guangzhou-6"
  ]
}

data "tencentcloud_cdc_dedicated_clusters" "example" {
  site_ids = [
    "site-2qu42ele",
    "site-fp8gp962"
  ]
}

data "tencentcloud_cdc_dedicated_clusters" "example" {
  lifecycle_statuses = [
    "PENDING",
    "RUNNING"
  ]
}
```

## Argument Reference

The following arguments are supported:

* `dedicated_cluster_ids` - (Optional, Set: [`String`]) Query by one or more instance IDs. Example of instance ID: cluster-xxxxxxxx.
* `lifecycle_statuses` - (Optional, Set: [`String`]) Filter by CDC life cycle.
* `name` - (Optional, String) Name of fuzzy matching CDC.
* `result_output_file` - (Optional, String) Used to save results.
* `site_ids` - (Optional, Set: [`String`]) Filter by site id.
* `zones` - (Optional, Set: [`String`]) Filter by AZ name.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `dedicated_cluster_set` - List of CDCs that meet the conditions.


