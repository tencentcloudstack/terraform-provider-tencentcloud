Use this data source to query detailed information of clb target_health

Example Usage

```hcl
data "tencentcloud_clb_target_health" "target_health" {
  load_balancer_ids = ["lb-5dnrkgry"]
}
```