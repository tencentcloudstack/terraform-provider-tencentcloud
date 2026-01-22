Provides a resource to create a redis parameter template

Example Usage

```hcl
resource "tencentcloud_redis_param_template" "example" {
  name         = "tf_example"
  description  = "This is an example redis param template."
  product_type = 6
  params_override {
    key   = "timeout"
    value = "7200"
  }
}
```

Copy from another template

```hcl
resource "tencentcloud_redis_param_template" "example" {
  name         = "tf-template"
  description  = "This is an example redis param template."
  product_type = 6
  params_override {
    key   = "timeout"
    value = "7200"
  }
}

resource "tencentcloud_redis_param_template" "example_copy" {
  name        = "tf-template-copied"
  description = "This is an copied redis param template from tf-template."
  template_id = tencentcloud_redis_param_template.example.id
}
```

Import

redis param_template can be imported using the id, e.g.
```
$ terraform import tencentcloud_redis_param_template.example crs-cfg-oyyon8f6
```
