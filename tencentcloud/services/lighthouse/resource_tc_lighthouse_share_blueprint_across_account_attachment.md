Provides a resource to create a lighthouse share blueprint across account attachment share

Example Usage

```hcl
resource "tencentcloud_lighthouse_share_blueprint_across_account_attachment" "share_blueprint_across_account_attachment" {
  blueprint_id = "lhbp-xxxxxx"
  account_ids  = ["100012345678"]
}
```

Import

lighthouse share_blueprint_across_account_attachment can be imported using the id, e.g.

```
terraform import tencentcloud_lighthouse_share_blueprint_across_account_attachment.share_blueprint_across_account_attachment lhbp-xxxxxx#100012345678
```
