Provides a resource to create a emr emr_yarn

Example Usage

```hcl
resource "tencentcloud_emr_yarn" "emr_yarn" {
  instance_id = "emr-rzrochgp"
  enable_resource_schedule = true
  scheduler = "fair"
  fair_global_config {
    user_max_apps_default = 1000
  }
}
```

Import

emr emr_yarn can be imported using the id, e.g.

```
terraform import tencentcloud_emr_yarn.emr_yarn emr_instance_id
```
