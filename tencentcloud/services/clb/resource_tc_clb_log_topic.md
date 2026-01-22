Provides a resource to create a CLB log topic.

Example Usage

```hcl
resource "tencentcloud_clb_log_topic" "example" {
  log_set_id = "2ed70190-bf06-4777-980d-2d8a327a2554"
  topic_name = "tf-example"
  status     = true
}
```

Import

CLB log topic can be imported using the id, e.g.

```
terraform import tencentcloud_clb_log_topic.example be1a83dd-04b4-4807-89bf-8daddce0df71
```
