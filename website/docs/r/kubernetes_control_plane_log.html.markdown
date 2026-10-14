---
subcategory: "Tencent Kubernetes Engine(TKE)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_kubernetes_control_plane_log"
sidebar_current: "docs-tencentcloud-resource-kubernetes_control_plane_log"
description: |-
  Provides a resource to create a TKE kubernetes control plane log
---

# tencentcloud_kubernetes_control_plane_log

Provides a resource to create a TKE kubernetes control plane log

## Example Usage

### Use automatic creation of log_set_id and topic_id

```hcl
resource "tencentcloud_kubernetes_control_plane_log" "example" {
  cluster_id   = "cls-rng1h5ei"
  cluster_type = "tke"
  components {
    name         = "karpenter"
    topic_region = "ap-guangzhou"
  }

  delete_log_set_and_topic = true
}
```

### Use custom log_set_id and topic_id

```hcl
resource "tencentcloud_kubernetes_control_plane_log" "example" {
  cluster_id   = "cls-rng1h5ei"
  cluster_type = "tke"
  components {
    name         = "cluster-autoscaler"
    log_set_id   = "40eed846-0f43-44b1-b216-c786a8970b1f"
    topic_id     = "21918a54-9ab4-40bc-90cd-c600cff00695"
    topic_region = "ap-guangzhou"
  }
}
```

## Argument Reference

The following arguments are supported:

* `cluster_id` - (Required, String, ForceNew) Cluster ID.
* `cluster_type` - (Required, String, ForceNew) Cluster type. currently only support tke.
* `components` - (Required, List, ForceNew) Component name list. currently supports cluster-autoscaler, kapenter.
* `delete_log_set_and_topic` - (Optional, Bool) Whether to simultaneously delete the log set and topic. If the log set and topic are used by other collection rules, they will not be deleted. Default is false.

The `components` object supports the following:

* `name` - (Required, String) Component name.
* `log_level` - (Optional, Int, **Deprecated**) It has been deprecated from version 1.82.54. Log level. for components that support dynamic adjustment, you can specify this parameter when enabling logs.
* `log_set_id` - (Optional, String) Logset ID. if not specified, auto-create.
* `topic_id` - (Optional, String) Log topic ID. if not specified, auto-create.
* `topic_region` - (Optional, String) topic region. this parameter enables cross-region shipping of logs.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.


## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [DescribeControlPlaneLogs](https://www.tencentcloud.com/document/product/457/75764) - This API is used to query plug-in log collection settings.
* Create: [EnableControlPlaneLogs](https://www.tencentcloud.com/document/product/457/75760) - This API is used to create log collection settings for plugins.
* Read: [DescribeControlPlaneLogs](https://www.tencentcloud.com/document/product/457/75764) - This API is used to query plug-in log collection settings.
* Delete: [DescribeControlPlaneLogs](https://www.tencentcloud.com/document/product/457/75764) - This API is used to query plug-in log collection settings.
* Delete: [DisableControlPlaneLogs](https://www.tencentcloud.com/document/product/457/75762) - This API is used to delete the log collection configuration of a plugin.

## Import

TKE kubernetes control plane log can be imported using the clusterId#clusterType#componentName, e.g.

```
terraform import tencentcloud_kubernetes_control_plane_log.example cls-rng1h5ei#tke#cluster-autoscaler
```

