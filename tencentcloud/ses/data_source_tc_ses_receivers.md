Use this data source to query detailed information of ses receivers

Example Usage

```hcl
data "tencentcloud_ses_receivers" "receivers" {
  status   = 3
  key_word = "keep"
}
```