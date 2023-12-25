Use this data source to query detailed information of waf waf_infos

Example Usage

```hcl
data "tencentcloud_waf_waf_infos" "example" {
  params {
    load_balancer_id = "lb-A8VF445"
  }
}
```

Or

```hcl
data "tencentcloud_waf_waf_infos" "example" {
  params {
    load_balancer_id = "lb-A8VF445"
    listener_id      = "lbl-nonkgvc2"
    domain_id        = "waf-MPtWPK5Q"
  }
}
```
