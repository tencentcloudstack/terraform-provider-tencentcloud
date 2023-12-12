Provides a resource to create a lighthouse key_pair_attachment

Example Usage

```hcl
resource "tencentcloud_lighthouse_key_pair_attachment" "key_pair_attachment" {
  key_id = "lhkp-xxxxxx"
  instance_id = "lhins-xxxxxx"
}
```

Import

lighthouse key_pair_attachment can be imported using the id, e.g.

```
terraform import tencentcloud_lighthouse_key_pair_attachment.key_pair_attachment key_pair_attachment_id
```