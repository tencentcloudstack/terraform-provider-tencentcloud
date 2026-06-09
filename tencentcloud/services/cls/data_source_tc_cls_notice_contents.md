Use this data source to query detailed information of CLS notice content templates.

Example Usage

Query all notice content templates

```hcl
data "tencentcloud_cls_notice_contents" "example" {}
```

Query by template name

```hcl
data "tencentcloud_cls_notice_contents" "example" {
  filters {
    key    = "name"
    values = ["DefaultTemplate(English)"]
  }
}
```

Query by template ID

```hcl
data "tencentcloud_cls_notice_contents" "example" {
  filters {
    key    = "noticeContentId"
    values = ["Default-en"]
  }
}
```
