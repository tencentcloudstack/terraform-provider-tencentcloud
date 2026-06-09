Use this data source to query CVM repair tasks.

Example Usage

```hcl
data "tencentcloud_cvm_repair_tasks" "tasks" {
  task_status = [1, 4]
}
```

Query with multiple filters

```hcl
data "tencentcloud_cvm_repair_tasks" "filtered" {
  product         = "CVM"
  task_status     = [1, 2]
  instance_ids    = ["ins-xxxxxxxx"]
  start_date      = "2023-01-01 00:00:00"
  end_date        = "2023-12-31 23:59:59"
  order_field     = "CreateTime"
  order           = 1
}
```
