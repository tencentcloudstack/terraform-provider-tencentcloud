Use this data source to query detailed information of organization org_financial_by_month

Example Usage

```hcl

data "tencentcloud_organization_org_financial_by_month" "org_financial_by_month" {
  end_month = "2023-05"
  member_uins = [100026517717]
}
```