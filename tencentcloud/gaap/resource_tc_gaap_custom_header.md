Provides a resource to create a gaap custom_header

Example Usage

```hcl
resource "tencentcloud_gaap_custom_header" "custom_header" {
  rule_id = "rule-xxxxxx"
  headers {
    header_name  = "HeaderName1"
    header_value = "HeaderValue1"
  }
  headers {
    header_name  = "HeaderName2"
    header_value = "HeaderValue2"
  }
}
```

Import

gaap custom_header can be imported using the id, e.g.

```
terraform import tencentcloud_gaap_custom_header.custom_header ruleId
```