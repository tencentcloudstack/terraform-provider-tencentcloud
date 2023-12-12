Use this data source to query detailed information of mariadb renewal_price

Example Usage

```hcl
data "tencentcloud_mariadb_renewal_price" "renewal_price" {
  instance_id = "tdsql-9vqvls95"
  period      = 2
}
```