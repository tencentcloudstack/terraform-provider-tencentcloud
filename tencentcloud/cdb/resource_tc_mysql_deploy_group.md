Provides a resource to create a mysql deploy_group

Example Usage

```hcl
resource "tencentcloud_mysql_deploy_group" "example" {
  deploy_group_name = "tf-example"
  description       = "desc."
  limit_num         = 1
  dev_class         = ["TS85"]
}
```

Import

mysql deploy_group can be imported using the id, e.g.

```
terraform import tencentcloud_mysql_deploy_group.deploy_group deploy_group_id
```