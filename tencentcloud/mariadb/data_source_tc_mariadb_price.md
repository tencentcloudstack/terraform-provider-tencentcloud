Use this data source to query detailed information of mariadb price

Example Usage

```hcl
data "tencentcloud_mariadb_price" "price" {
  zone       = "ap-guangzhou-3"
  node_count = 2
  memory     = 2
  storage    = 20
  buy_count  = 1
  period     = 1
  paymode    = "prepaid"
}
```