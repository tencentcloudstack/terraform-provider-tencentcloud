Use this data source to query detailed information of waf ports

Example Usage

```hcl
data "tencentcloud_waf_ports" "example" {}
```

Or

```hcl
data "tencentcloud_waf_ports" "example" {
  edition     = "clb-waf"
  instance_id = "waf_2kxtlbky00b2v1fn"
}
```

Or

```hcl
data "tencentcloud_waf_ports" "example" {
  edition     = "sparta-waf"
  instance_id = "waf_2ka80zly0702e8j3"
}
```
