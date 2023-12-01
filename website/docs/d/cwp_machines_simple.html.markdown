---
subcategory: "Cwp"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cwp_machines_simple"
sidebar_current: "docs-tencentcloud-datasource-cwp_machines_simple"
description: |-
  Use this data source to query detailed information of cwp machines_simple
---

# tencentcloud_cwp_machines_simple

Use this data source to query detailed information of cwp machines_simple

## Example Usage

```hcl
data "tencentcloud_cwp_machines_simple" "example" {
  machine_type   = "CVM"
  machine_region = "ap-guangzhou"
  project_ids    = [1210293, 1157652]
}
```

### Query by Keyword filter

```hcl
data "tencentcloud_cwp_machines_simple" "example" {
  machine_type   = "CVM"
  machine_region = "ap-guangzhou"
  project_ids    = [0]

  filters {
    name        = "Keywords"
    values      = ["tf_example"]
    exact_match = true
  }
}
```

### Query by Version filter

```hcl
data "tencentcloud_cwp_machines_simple" "example" {
  machine_type   = "CVM"
  machine_region = "ap-guangzhou"
  project_ids    = [0]

  filters {
    name        = "Version"
    values      = ["BASIC_VERSION"]
    exact_match = true
  }
}
```

### Query by TagId filter

```hcl
data "tencentcloud_cwp_machines_simple" "example" {
  machine_type   = "ALL"
  machine_region = "all-regions"

  filters {
    name        = "TagId"
    values      = ["13771"]
    exact_match = true
  }
}
```

## Argument Reference

The following arguments are supported:

* `machine_region` - (Required, String) The area where the machine belongs,Such as: ap-guangzhou, ap-shanghai, all-regions: All server region types.
* `machine_type` - (Required, String) Service types. -CVM: Cloud Virtual Machine; -ECM: Edge Computing Machine; -LH: Lighthouse; -Other: Mixed cloud; -ALL: All server types.
* `filters` - (Optional, List) filter list.
* `project_ids` - (Optional, Set: [`Int`]) Project id list.
* `result_output_file` - (Optional, String) Used to save results.

The `filters` object supports the following:

* `name` - (Required, String) Only supported Keywords, Version and TagId.
* `values` - (Required, Set) If `name` is `Keywords`: enter keyword query; If `name` is `Version`: enter PRO_VERSION: Professional Edition | BASIC_VERSION: Basic | Flagship: Flagship | ProtectedMachines: Professional+Flagship | UnFlagship: Non Flagship | PRO_POST_PAY: Professional Edition Pay by Volume | PRO_PRE_PAY: Professional Edition Monthly Package query; If `name` is `TagId`: enter tag ID query.
* `exact_match` - (Optional, Bool) exact match. true or false.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `machines` - Machine list.
  * `cloud_tags` - Cloud tags detailNote: This field may return null, indicating that a valid value cannot be obtained.
    * `tag_key` - Tag key.
    * `tag_value` - Tag value.
  * `instance_id` - Instance IDNote: This field may return null, indicating that a valid value cannot be obtained.
  * `instance_state` - Instance status.
  * `is_pro_version` - Paid version or not. true: yes; false: no.
  * `kernel_version` - Core Version.
  * `license_order` - License Order ObjectNote: This field may return null, indicating that a valid value cannot be obtained.
    * `license_id` - License ID.
    * `license_type` - License Types.
    * `resource_id` - Resource ID.
    * `source_type` - Order types.
    * `status` - License Order Status.
  * `machine_ip` - Machine Internal net IP.
  * `machine_name` - Machine name.
  * `machine_os` - Machine OS System.
  * `machine_type` - Service types. -CVM: Cloud Virtual Machine; -ECM: Edge Computing Machine -LH: Lighthouse; -Other: Mixed cloud; -ALL: All server types.
  * `machine_wan_ip` - Machine Outer net IP.
  * `pay_mode` - Payment model. POSTPAY: Pay as you go; PREPAY: Monthly subscription.
  * `project_id` - Project ID.
  * `protect_type` - Protection Version. -BASIC_VERSION: Basic Version; -PRO_VERSION: Pro Version -Flagship: Flagship Version; -GENERAL_DISCOUNT: CWP-LH Version.
  * `quuid` - Cloud server sole UUID.
  * `region_info` - Region detail.
    * `region_code` - Region Code.
    * `region_id` - Region ID.
    * `region_name_en` - Regional English name.
    * `region_name` - Regional Chinese name.
    * `region` - Region, Such as ap-guangzhou, ap-shanghai, ap-beijing.
  * `tag` - Tag.
    * `name` - Tag name.
    * `rid` - Relevance tag id.
    * `tag_id` - Tag ID.
  * `uuid` - Cwp client sole UUID.


