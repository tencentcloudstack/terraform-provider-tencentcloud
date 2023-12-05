Use this data source to query detailed information of gaap proxy group statistics

Example Usage

```hcl
data "tencentcloud_gaap_proxy_group_statistics" "proxy_group_statistics" {
	group_id = "link-8lpyo88p"
	start_time = "2023-10-09 00:00:00"
	end_time = "2023-10-09 23:59:59"
	metric_names = ["InBandwidth", "OutBandwidth", "InFlow", "OutFlow"]
	granularity = 300
}
```