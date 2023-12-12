Use this data source to query detailed information of dcdb renewal_price

Example Usage

```hcl
data "tencentcloud_dcdb_renewal_price" "renewal_price" {
	instance_id = local.dcdb_id
	period      = 1
	amount_unit = "pent"
}
```