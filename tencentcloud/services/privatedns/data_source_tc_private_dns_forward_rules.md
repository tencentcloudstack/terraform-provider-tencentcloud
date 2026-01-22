Use this data source to query detailed information of Private Dns forward rules

Example Usage

Query all private dns forward rules

```hcl
data "tencentcloud_private_dns_forward_rules" "example" {}
```

Query all private dns forward rules by filters

```hcl
data "tencentcloud_private_dns_forward_rules" "example" {
  filters {
    name   = "RuleId"
    values = ["fid-2ece6ca305"]
  }

  filters {
    name   = "RuleName"
    values = ["tf-example"]
  }

  filters {
    name   = "RuleType"
    values = ["DOWN"]
  }

  filters {
    name   = "ZoneId"
    values = ["zone-04jlawty"]
  }

  filters {
    name   = "EndPointId"
    values = ["eid-e9d5880672"]
  }

  filters {
    name   = "EndPointName"
    values = ["tf-example"]
  }
}
```
