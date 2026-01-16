Use this data source to query detailed information of CFW ccn instance region status

Example Usage

```hcl
data "tencentcloud_cfw_ccn_instance_region_status" "example" {
  ccn_id = "ccn-fkb9bo2v"
  instance_ids = [
    "vpc-axbsvrrg"
  ]
  routing_mode = 1
}
```
