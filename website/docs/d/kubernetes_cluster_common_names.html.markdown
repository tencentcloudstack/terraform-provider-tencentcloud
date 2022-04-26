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

```hcl
data "tencentcloud_kubernetes_cluster_common_names" "foo" {
  cluster_id      = "cls-12345678"
  subaccount_uins = ["1234567890", "0987654321"]
}
```

## Argument Reference

The following arguments are supported:

* `cluster_id` - (Optional) Cluster ID.
* `result_output_file` - (Optional) Used for save result.
* `role_ids` - (Optional) List of Role ID. Up to 50 sub-accounts can be passed in at a time.
* `subaccount_uins` - (Optional) List of sub-account. Up to 50 sub-accounts can be passed in at a time.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `list` - List of the CommonName in the certificate of the client corresponding to the sub-account UIN.
  * `common_names` - The CommonName in the certificate of the client corresponding to the sub-account.
  * `subaccount_uin` - User UIN.


