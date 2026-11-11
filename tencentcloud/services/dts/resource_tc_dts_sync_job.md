Provides a resource to create a DTS sync job

~> **NOTE:** Import function does not support field `existed_job_id`.

Example Usage

```hcl
resource "tencentcloud_dts_sync_job" "example" {
  pay_mode          = "PostPay"
  src_database_type = "mysql"
  src_region        = "ap-guangzhou"
  dst_database_type = "cynosdbmysql"
  dst_region        = "ap-guangzhou"
  auto_renew        = 0
  instance_class    = "micro"
  tags {
    tag_key   = "createBy"
    tag_value = "Terraform"
  }
}
```

Import

DTS sync job can be imported using the id, e.g.

```
terraform import tencentcloud_dts_sync_job.example sync-hpb214ua
```
