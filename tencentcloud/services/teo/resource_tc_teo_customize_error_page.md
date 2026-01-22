Provides a resource to create a TEO customize error page

Example Usage

If content_type is application/json

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

If content_type is text/html

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

If content_type is text/plain

```hcl
resource "tencentcloud_teo_customize_error_page" "example" {
  zone_id      = "zone-3edjdliiw3he"
  name         = "tf-example"
  content_type = "text/plain"
  description  = "description."
  content      = "customize error page"
}
```

If content_type is text/xml

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

Import

TEO customize error page can be imported using the id, e.g.

```
terraform import tencentcloud_teo_customize_error_page.example zone-3edjdliiw3he#p-3egexy9b4426
```
