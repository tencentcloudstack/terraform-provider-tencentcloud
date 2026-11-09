Use this data source to query detailed information of a CloudAudit audit track set.

Example Usage

```hcl
data "tencentcloud_cloudaudit_audit" "audit" {
  audit_name = "tf-audit"
}
```