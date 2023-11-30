Provides a resource to create a dlc suspend_resume_data_engine

Example Usage

```hcl
resource "tencentcloud_dlc_suspend_resume_data_engine" "suspend_resume_data_engine" {
  data_engine_name = "example-iac"
  operate = "suspend"
}
```

Import

dlc suspend_resume_data_engine can be imported using the id, e.g.

```
terraform import tencentcloud_dlc_suspend_resume_data_engine.suspend_resume_data_engine suspend_resume_data_engine_id
```