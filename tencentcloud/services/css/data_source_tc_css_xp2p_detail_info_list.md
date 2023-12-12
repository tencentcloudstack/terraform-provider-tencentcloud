Use this data source to query detailed information of css xp2p_detail_info_list

Example Usage

```hcl
data "tencentcloud_css_xp2p_detail_info_list" "xp2p_detail_info_list" {
  query_time   = "2023-11-01T14:55:01+08:00"
  type         = ["live"]
}
```