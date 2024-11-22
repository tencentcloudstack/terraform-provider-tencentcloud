---
subcategory: "Tencent Kubernetes Engine(TKE)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_kubernetes_cluster_common_names"
sidebar_current: "docs-tencentcloud-datasource-kubernetes_cluster_common_names"
description: |-
  Provide a datasource to query cluster CommonNames.
---

# tencentcloud_kubernetes_cluster_common_names

Provide a datasource to query cluster CommonNames.

## Example Usage

### Query common names by subaccount uins

```hcl
data "tencentcloud_kubernetes_cluster_common_names" "example" {
  cluster_id      = "cls-fdy7hm1q"
  subaccount_uins = ["100037718139", "100031340176"]
}
```

### Query common names by role ids

```hcl
data "tencentcloud_kubernetes_cluster_common_names" "example" {
  cluster_id = "cls-fdy7hm1q"
  role_ids   = ["4611686018441060141"]
}
```

## Argument Reference

The following arguments are supported:

* `cluster_id` - (Optional, String) Cluster ID.
* `result_output_file` - (Optional, String) Used to save results.
* `role_ids` - (Optional, List: [`String`]) List of Role ID. Up to 50 sub-accounts can be passed in at a time.
* `subaccount_uins` - (Optional, List: [`String`]) List of sub-account. Up to 50 sub-accounts can be passed in at a time.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `list` - List of the CommonName in the certificate of the client corresponding to the sub-account UIN.
  * `common_name` - The CommonName in the certificate of the client corresponding to the sub-account.
  * `common_names` - (**Deprecated**) It has been deprecated from version 1.81.140. Please use `common_name`. The CommonName in the certificate of the client corresponding to the sub-account.
  * `subaccount_uin` - User UIN.


