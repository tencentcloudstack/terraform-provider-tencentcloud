Provides a resource to create a bi embed_token

Example Usage

```hcl
resource "tencentcloud_bi_embed_token_apply" "embed_token" {
  project_id   = 11015030
  page_id      = 10520483
  scope        = "page"
  expire_time  = "240"
  user_corp_id = "100022975249"
  user_id      = "100024664626"
}
```