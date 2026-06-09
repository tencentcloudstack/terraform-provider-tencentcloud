---
subcategory: "TencentCloud EdgeOne(TEO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_teo_default_certificate"
sidebar_current: "docs-tencentcloud-datasource-teo_default_certificate"
description: |-
  Use this data source to query detailed information of TEO default certificates
---

# tencentcloud_teo_default_certificate

Use this data source to query detailed information of TEO default certificates

## Example Usage

```hcl
data "tencentcloud_teo_default_certificate" "example" {
  zone_id = "zone-2qtuhspy7cr6"
}
```

### Query with filters

```hcl
data "tencentcloud_teo_default_certificate" "example" {
  filters {
    name = "zone-id"
    values = [
      "zone-2qtuhspy7cr6"
    ]
  }
}
```

## Argument Reference

The following arguments are supported:

* `filters` - (Optional, List) Filter conditions, the upper limit of Filters.Values is 5. The detailed filtering conditions are as follows: zone-id - Filter by zone ID. At least one of `zone_id` or `filters` must be specified.
* `result_output_file` - (Optional, String) Used to save results.
* `zone_id` - (Optional, String) Zone ID. At least one of `zone_id` or `filters` must be specified.

The `filters` object supports the following:

* `name` - (Required, String) Filter name.
* `values` - (Required, Set) Filter value.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `default_server_cert_info` - Default certificate list.
  * `alias` - Certificate alias.
  * `cert_id` - Server certificate ID.
  * `common_name` - Certificate common name.
  * `effective_time` - Certificate effective time.
  * `expire_time` - Certificate expiration time.
  * `message` - Failure reason when Status is failed.
  * `sign_algo` - Certificate signing algorithm.
  * `status` - Deploy status. Valid values: processing (deploying), deployed (deployed), failed (deploy failed).
  * `subject_alt_name` - Certificate SAN domains.
  * `type` - Certificate type. Valid values: default (default certificate), upload (user uploaded), managed (Tencent Cloud managed).


