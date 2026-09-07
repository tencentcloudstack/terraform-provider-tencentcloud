Use this data source to query historical instant tasks of CAT (Cat).

Example Usage

```hcl
data "tencentcloud_cat_instant_tasks" "example" {
}

output "tasks" {
  value = data.tencentcloud_cat_instant_tasks.example.tasks
}

output "total" {
  value = data.tencentcloud_cat_instant_tasks.example.total
}
```
