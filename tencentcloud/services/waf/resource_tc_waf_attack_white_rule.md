Provides a resource to create a WAF attack white rule

Example Usage

Using type_ids

```hcl
resource "tencentcloud_waf_attack_white_rule" "example" {
  domain = "www.demo.com"
  name   = "tf-example"
  status = 1
  mode   = 0
  rules {
    match_field   = "IP"
    match_method  = "ipmatch"
    match_content = "1.1.1.1"
  }

  rules {
    match_field   = "Referer"
    match_method  = "eq"
    match_content = "referer content"
  }

  rules {
    match_field   = "URL"
    match_method  = "contains"
    match_content = "/prefix"
  }

  rules {
    match_field   = "HTTP_METHOD"
    match_method  = "neq"
    match_content = "POST"
  }

  rules {
    match_field   = "GET"
    match_method  = "ncontains"
    match_content = "value"
    match_params  = "key"
  }

  type_ids = [
    "010000000",
    "020000000",
    "030000000",
    "040000000",
    "050000000",
    "060000000",
    "090000000",
    "110000000"
  ]
}
```

Using signature_ids

```hcl
resource "tencentcloud_waf_attack_white_rule" "example" {
  domain = "www.demo.com"
  name   = "tf-example"
  status = 0
  mode   = 1
  rules {
    match_field   = "IP"
    match_method  = "ipmatch"
    match_content = "1.1.1.1"
  }

  signature_ids = [
    "60270036",
    "10000047"
  ]
}
```

Import

WAF attack white rule can be imported using the id, e.g.

```
terraform import tencentcloud_waf_attack_white_rule.example www.demo.com#38562
```
