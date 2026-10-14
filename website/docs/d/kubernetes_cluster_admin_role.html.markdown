---
subcategory: "Tencent Kubernetes Engine(TKE)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_kubernetes_cluster_admin_role"
sidebar_current: "docs-tencentcloud-datasource-kubernetes_cluster_admin_role"
description: |-
  Provide a datasource to acquire TKE cluster admin role.
---

# tencentcloud_kubernetes_cluster_admin_role

Provide a datasource to acquire TKE cluster admin role.

Use this data source to grant the current user (or sub-account) the `tke:admin` ClusterRole in the specified Kubernetes cluster. This is typically used when a CAM sub-account needs to be granted cluster administrator permissions through a CAM policy.

## Example Usage

```hcl
data "tencentcloud_kubernetes_cluster_admin_role" "foo" {
  cluster_id = "cls-xxxxxxxx"
}

output "request_id" {
  value = data.tencentcloud_kubernetes_cluster_admin_role.foo.request_id
}
```

## Argument Reference

The following arguments are supported:

* `cluster_id` - (Required, String) Cluster ID.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `request_id` - The request ID returned by the API.

## API Reference

This data source queries the following TencentCloud APIs:

* Read: [AcquireClusterAdminRole](https://cloud.tencent.com/document/api/457/49014) - This API can be called to acquire the ClusterRole tke:admin. By setting a CAM policy, you can grant permission of this API to a sub-account that has higher permission in CAM.

