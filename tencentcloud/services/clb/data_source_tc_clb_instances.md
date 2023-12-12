Use this data source to query detailed information of CLB

Example Usage

```hcl
data "tencentcloud_clb_instances" "foo" {
  clb_id             = "lb-k2zjp9lv"
  network_type       = "OPEN"
  clb_name           = "myclb"
  project_id         = 0
  result_output_file = "mytestpath"
}
```