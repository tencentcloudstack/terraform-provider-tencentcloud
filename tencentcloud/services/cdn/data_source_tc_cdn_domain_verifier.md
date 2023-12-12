Provides a resource to check or create a cdn Domain Verify Record

~> **NOTE:**

Example Usage

```hcl
data "tencentcloud_cdn_domain_verifier" "vr" {
  domain = "www.examplexxx123.com"
  auto_verify = true # auto create record if not verified
  freeze_record = true # once been freeze and verified, it will never be changed again
}

locals {
  recordValue = data.tencentcloud_cdn_domain_verifier.record
  recordType = data.tencentcloud_cdn_domain_verifier.record_type
}
```