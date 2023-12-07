Use this data source to query detailed information of css domains

Example Usage

```hcl
data "tencentcloud_css_domains" "domains" {
  domain_type = 0
  play_type = 1
  is_delay_live = 0
}
```