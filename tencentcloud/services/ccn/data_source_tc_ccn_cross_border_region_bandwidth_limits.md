Use this data source to query detailed information of ccn_cross_border_region_bandwidth_limits

-> **NOTE:** This resource is dedicated to Unicom.

Example Usage

```hcl
data "tencentcloud_ccn_cross_border_region_bandwidth_limits" "ccn_region_bandwidth_limits" {
  filters {
    name   = "source-region"
    values = ["ap-guangzhou"]
  }

  filters {
    name   = "destination-region"
    values = ["ap-shanghai"]
  }
}
```