Use this data source to query detailed information of gaap proxy groups

Example Usage

```hcl
data "tencentcloud_gaap_proxy_groups" "proxy_groups" {
	project_id = 0
	filters {
		name = "GroupId"
		values = ["lg-5anbbou5"]
	}
}
```