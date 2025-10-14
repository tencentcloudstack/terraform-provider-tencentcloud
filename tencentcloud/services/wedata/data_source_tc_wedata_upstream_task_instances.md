Use this data source to query detailed information of wedata upstream task instances

Example Usage

```hcl
data "tencentcloud_wedata_task_instances" "wedata_task_instances" {
  project_id = "1859317240494305280"
}

locals {
  instance_keys = data.tencentcloud_wedata_task_instances.wedata_task_instances.data[0].items[*].instance_key
}

data "tencentcloud_wedata_upstream_task_instances" "wedata_upstream_task_instances" {
  for_each = toset(local.instance_keys)

  project_id  = "1859317240494305280"
  instance_key = each.value
}
```