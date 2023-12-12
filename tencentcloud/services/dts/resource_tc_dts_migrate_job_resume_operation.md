Provides a resource to create a dts migrate_job_resume_operation

Example Usage

```hcl
resource "tencentcloud_dts_migrate_job_resume_operation" "resume" {
	job_id = "job_id"
	resume_option = "normal"
}
```