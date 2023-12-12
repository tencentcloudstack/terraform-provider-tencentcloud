Use this data source to query detailed information of ccn cross_border_compliance

Example Usage

```hcl
data "tencentcloud_ccn_cross_border_compliance" "cross_border_compliance" {
  service_provider = "UNICOM"
  compliance_id = 10002
  email = "test@tencent.com"
  service_start_date = "2020-07-29"
  service_end_date = "2021-07-29"
  state = "APPROVED"
}
```