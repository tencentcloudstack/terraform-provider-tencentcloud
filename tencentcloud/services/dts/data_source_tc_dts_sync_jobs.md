Use this data source to query detailed information of dts syncJobs

Example Usage

```hcl
resource "tencentcloud_dts_sync_job" "job" {
	job_name = "tf_dts_test"
	pay_mode = "PostPay"
	src_database_type = "mysql"
	src_region = "ap-guangzhou"
	dst_database_type = "cynosdbmysql"
	dst_region = "ap-guangzhou"
	tags {
	  tag_key = "aaa"
	  tag_value = "bbb"
	}
	auto_renew = 0
   instance_class = "micro"
  }

data "tencentcloud_dts_sync_jobs" "sync_jobs" {
  job_id = tencentcloud_dts_sync_job.job.id
  job_name = "tf_dts_test"
}
```