---
subcategory: "ControlCenter"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_controlcenter_account_factory_baseline_items"
sidebar_current: "docs-tencentcloud-datasource-controlcenter_account_factory_baseline_items"
description: |-
  Use this data source to query detailed information of Controlcenter account factory baseline items
---

# tencentcloud_controlcenter_account_factory_baseline_items

Use this data source to query detailed information of Controlcenter account factory baseline items

## Example Usage

```hcl
data "tencentcloud_controlcenter_account_factory_baseline_items" "example" {}
```

## Argument Reference

The following arguments are supported:

* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `baseline_items` - Account factory baseline list.
  * `classify_en` - Baseline english classification, with a length of 2-64 english characters. cannot be empty.
  * `classify` - Baseline classification. length: 2-32 english or chinese characters. values cannot be empty.
  * `depends_on` - Baseline item dependency. value range of N depends on the count of other baseline items it relies on.
    * `identifier` - Specifies the unique identifier for the feature item, can only contain `english letters`, `digits`, and `@,._[]-:()()[]+=.`, with a length of 2-128 characters.
    * `type` - Dependency type. valid values: LandingZoneSetUp or AccountFactorySetUp. LandingZoneSetUp refers to the dependency of landingZone. AccountFactorySetUp refers to the dependency of account factory.
  * `description_en` - Baseline item english description, with a length of 2 to 1024 english characters. it is empty by default.
  * `description` - Baseline description, with a length of 2 to 256 english or chinese characters. it is empty by default.
  * `identifier` - Specifies the unique identifier for account factory baseline item, can only contain `english letters`, `digits`, and `@,._[]-:()()[]+=.`, with a length of 2-128 characters.
  * `name_en` - Baseline item english name. specifies a unique name for the baseline item. supports a combination of english letters, digits, spaces, and symbols @, &, _, [], -. valid values: 1-64 english characters.
  * `name` - Baseline item name. specifies a unique name for the feature item. supports a combination of english letters, numbers, chinese characters, and symbols @, &, _, [, ], -. valid values: 1-25 chinese or english characters.
  * `required` - Specifies whether the baseline item is required (1: required; 0: optional).
  * `weight` - Baseline item weight. the smaller the value, the higher the weight. value range equal to or greater than 0.


