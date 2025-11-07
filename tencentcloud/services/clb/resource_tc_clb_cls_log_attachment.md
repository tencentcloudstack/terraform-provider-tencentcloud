Provides a resource to create a CLB cls log attachment

Example Usage

```hcl
resource "tencentcloud_clb_log_topic" "example" {
  log_set_id = "2ed70190-bf06-4777-980d-2d8a327a2554"
  topic_name = "tf-example"
  status     = true
}

resource "tencentcloud_clb_cls_log_attachment" "example" {
  load_balancer_id = "lb-n26tx0bm"
  log_set_id       = "2ed70190-bf06-4777-980d-2d8a327a2554"
  log_topic_id     = tencentcloud_clb_log_topic.example.id
}
```

Import

CLB cls log attachment can be imported using the loadBalancerId#logSetId#logTopicId, e.g.

```
terraform import tencentcloud_clb_cls_log_attachment.example lb-n26tx0bm#2ed70190-bf06-4777-980d-2d8a327a2554#ac2fda28-3e79-4b51-b193-bfcf1aeece24
```
