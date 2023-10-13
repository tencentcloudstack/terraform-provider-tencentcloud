---
subcategory: "MapReduce(EMR)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_emr_user_manager"
sidebar_current: "docs-tencentcloud-resource-emr_user_manager"
description: |-
  Provides a resource to create a emr user
---

# tencentcloud_emr_user_manager

Provides a resource to create a emr user

## Example Usage

```hcl
data "tencentcloud_emr" "my_emr" {
  display_strategy = "clusterList"
}

resource "tencentcloud_emr_user_manager" "user_manager" {
  instance_id = data.tencentcloud_emr.my_emr.clusters.0.cluster_id
  user_name   = "tf-test"
  user_group  = "group1"
  password    = "tf@123456"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String, ForceNew) Cluster string ID.
* `password` - (Required, String) PassWord.
* `user_group` - (Required, String, ForceNew) User group membership.
* `user_name` - (Required, String, ForceNew) Username.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `create_time` - Create time.
* `download_keytab_url` - Download keytab url.
* `support_download_keytab` - If support download keytab.
* `user_type` - User type.


## Import

emr user_manager can be imported using the id, e.g.

```
terraform import tencentcloud_emr_user_manager.user_manager instanceId#userName
```

