Provides a resource to create a TCSS refresh task operation

Example Usage

```hcl
resource "tencentcloud_tcss_refresh_task_operation" "example" {}
```

Or

```hcl
resource "tencentcloud_tcss_refresh_task_operation" "example" {
  cluster_ids = [
    "cls-fdy7hm1q"
  ]
  is_sync_list_only = false
}
```

