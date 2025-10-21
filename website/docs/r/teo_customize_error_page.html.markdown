---
subcategory: "TencentCloud EdgeOne(TEO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_teo_customize_error_page"
sidebar_current: "docs-tencentcloud-resource-teo_customize_error_page"
description: |-
  Provides a resource to create a TEO customize error page
---

# tencentcloud_teo_customize_error_page

Provides a resource to create a TEO customize error page

## Example Usage

### If content_type is application/json

```hcl
resource "tencentcloud_teo_customize_error_page" "example" {
  zone_id      = "zone-3edjdliiw3he"
  name         = "tf-example"
  content_type = "application/json"
  description  = "description."
  content = jsonencode({
    "key" : "value",
  })
}
```

### If content_type is text/html

```hcl
resource "tencentcloud_teo_customize_error_page" "example" {
  zone_id      = "zone-3edjdliiw3he"
  name         = "tf-example"
  content_type = "text/html"
  description  = "description."
  content      = <<-EOF
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>Error Page</title>
</head>
<body>
    customize error page
</body>
</html>
  EOF
}
```

### If content_type is text/plain

```hcl
resource "tencentcloud_teo_customize_error_page" "example" {
  zone_id      = "zone-3edjdliiw3he"
  name         = "tf-example"
  content_type = "text/plain"
  description  = "description."
  content      = "customize error page"
}
```

### If content_type is text/xml

```hcl
resource "tencentcloud_teo_customize_error_page" "example" {
  zone_id      = "zone-3edjdliiw3he"
  name         = "tf-example"
  content_type = "text/xml"
  description  = "description."
  content      = <<-EOF
<?xml version="1.0" encoding="UTF-8"?>
<?xml-stylesheet type="text/css" href="#internalStyle"?>
<error-page>
  <message>customize error page</message>
</error-page>
  EOF
}
```

## Argument Reference

The following arguments are supported:

* `content_type` - (Required, String) Custom error page type, with values:<li>text/html; </li><li>application/json;</li><li>text/plain;</li><li>text/xml.</li>.
* `name` - (Required, String) Custom error page name. The name must be 2-30 characters long.
* `zone_id` - (Required, String, ForceNew) Zone ID.
* `content` - (Optional, String) Custom error page content, not exceeding 2 KB.
* `description` - (Optional, String) Custom error page description, not exceeding 60 characters.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `page_id` - Page ID.


## Import

TEO customize error page can be imported using the id, e.g.

```
terraform import tencentcloud_teo_customize_error_page.example zone-3edjdliiw3he#p-3egexy9b4426
```

