---
subcategory: "Serverless Cloud Function(SCF)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_scf_layer_versions"
sidebar_current: "docs-tencentcloud-datasource-scf_layer_versions"
description: |-
  Use this data source to query detailed information of scf layer_versions
---

# tencentcloud_scf_layer_versions

Use this data source to query detailed information of scf layer_versions

## Example Usage

```hcl
data "tencentcloud_scf_layer_versions" "layer_versions" {
  layer_name = "tf-test"
}
```

## Argument Reference

The following arguments are supported:

* `layer_name` - (Required, String) Layer name.
* `compatible_runtime` - (Optional, Set: [`String`]) Compatible runtimes.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `layer_versions` - Layer version list.
  * `add_time` - Creation time.
  * `compatible_runtimes` - Runtime applicable to a versionNote: This field may return null, indicating that no valid values can be obtained.
  * `description` - Version descriptionNote: This field may return null, indicating that no valid values can be obtained.
  * `layer_name` - Layer name.
  * `layer_version` - Version number.
  * `license_info` - License informationNote: This field may return null, indicating that no valid values can be obtained.
  * `stamp` - StampNote: This field may return null, indicating that no valid values can be obtained.
  * `status` - Current status of specific layer version. For valid values, please see [here](https://intl.cloud.tencent.com/document/product/583/47175?from_cn_redirect=1#.E5.B1.82.EF.BC.88layer.EF.BC.89.E7.8A.B6.E6.80.81).


