Use this data source to query detailed information of organization org_financial_by_product

Example Usage

```hcl

data "tencentcloud_organization_org_financial_by_product" "org_financial_by_product" {
  month = "2023-05"
  end_month = "2023-09"
  product_codes = ["p_eip"]
  }
```