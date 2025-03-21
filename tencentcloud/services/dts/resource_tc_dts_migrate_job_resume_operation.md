Provides a resource to create a DTS migrate job resume operation

Example Usage

```hcl
resource "tencentcloud_dts_migrate_job_resume_operation" "example" {
  job_id        = "dts-puwyj5uy"
  resume_option = "normal"
}
```