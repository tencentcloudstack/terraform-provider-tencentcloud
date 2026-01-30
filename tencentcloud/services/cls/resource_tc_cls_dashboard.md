Provides a resource to create a CLS Dashboard.

Example Usage

```hcl
resource "tencentcloud_cls_dashboard" "dashboard" {
  dashboard_name = "my-dashboard"
}
```

With configuration data

```hcl
resource "tencentcloud_cls_dashboard" "dashboard" {
  dashboard_name = "production-dashboard"
  data           = jsonencode({
    timezone = "browser"
    subType  = "CLS_Host"
  })
  
  tags = {
    "team"        = "ops"
    "environment" = "production"
  }
}
```

Import

CLS dashboard can be imported using the id, e.g.

```
$ terraform import tencentcloud_cls_dashboard.dashboard dashboard-xxxx-xxxx-xxxx-xxxx
```
