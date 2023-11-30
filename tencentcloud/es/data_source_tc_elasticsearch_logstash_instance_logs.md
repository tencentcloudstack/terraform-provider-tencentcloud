Use this data source to query detailed information of elasticsearch logstash_instance_logs

Example Usage

```hcl
data "tencentcloud_elasticsearch_logstash_instance_logs" "logstash_instance_logs" {
	instance_id = "ls-xxxxxx"
	log_type = 1
	start_time = "2023-10-31 10:30:00"
	end_time = "2023-10-31 10:30:10"
}
```