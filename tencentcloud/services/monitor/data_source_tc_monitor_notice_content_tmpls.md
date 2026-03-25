Use this data source to query monitor notice content templates.

Example Usage

Query all templates

```hcl
data "tencentcloud_monitor_notice_content_tmpls" "example" {}
```

Query by filter

```hcl
data "tencentcloud_monitor_notice_content_tmpls" "example" {
  tmpl_ids      = ["ntpl-plu46bk5"]
  tmpl_name     = "tf-example"
  notice_id     = "notice-c2af43ee-1a4b-4c4a-ae3e-f81481280101"
  tmpl_language = "en"
  monitor_type  = "MT_QCE"
}
```
