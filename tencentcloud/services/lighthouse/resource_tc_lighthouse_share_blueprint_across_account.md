Provides a resource to share a Lighthouse blueprint across accounts.

Example Usage

```hcl
resource "tencentcloud_lighthouse_share_blueprint_across_account" "example" {
  blueprint_id = "lhbp-xxxxxx"
  account_ids  = ["100000000001", "100000000002"]
}
```
