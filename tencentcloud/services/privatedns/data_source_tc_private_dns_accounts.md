Use this data source to query detailed information of privatedns accounts

Example Usage

Query all accounts

```hcl
data "tencentcloud_private_dns_accounts" "example" {}
```

Query accounts by filters

```hcl
data "tencentcloud_private_dns_accounts" "example" {
  filters {
    name   = "AccountUin"
    values = ["100022770160"]
  }
}
```