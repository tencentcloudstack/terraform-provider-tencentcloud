Use this data source to query detailed information of elasticsearch logstash_instance_operations

Example Usage

```hcl
data "tencentcloud_elasticsearch_logstash_instance_operations" "logstash_instance_operations" {
	instance_id = "ls-xxxxxx"
	start_time = "2018-01-01 00:00:00"
	end_time = "2023-10-31 10:12:45"
}
```