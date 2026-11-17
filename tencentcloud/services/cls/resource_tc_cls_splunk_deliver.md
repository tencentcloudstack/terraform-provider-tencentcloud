Provides a resource to create a CLS splunk deliver.

Example Usage

```hcl
resource "tencentcloud_cls_splunk_deliver" "example" {
  topic_id  = "eb417b1d-fc00-46f2-85f7-47ca0848a6b3"
  name      = "tf-example"
  index_ack = 1

  net_info {
    host     = "110.11.22.106"
    port     = 8088
    token    = "e27274fb-****-****-****-****00206282"
    net_type = 2
    is_ssl   = true
  }

  metadata_info {
    format = "json"
    meta_fields = [
      "__HOSTNAME__",
      "__FILENAME__",
      "__TIMESTAMP__",
    ]
  }
}
```

Import

cls splunk deliver can be imported using the topicId#taskId, e.g.

```
terraform import tencentcloud_cls_splunk_deliver.example eb417b1d-fc00-46f2-85f7-47ca0848a6b3#716dae34-7eba-4b09-a3d9-6fe8e510e236
```