Use this data source to query detailed information of monitor tmp instances

Example Usage

```hcl
data "tencentcloud_monitor_tmp_instances" "tmp_instances" {
	instance_ids = ["prom-xxxxxx"]
}
```