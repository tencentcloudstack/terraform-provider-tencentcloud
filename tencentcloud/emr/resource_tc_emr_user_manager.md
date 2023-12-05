Provides a resource to create a emr user

Example Usage

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

Import

emr user_manager can be imported using the id, e.g.

```
terraform import tencentcloud_emr_user_manager.user_manager instanceId#userName
```