Provides a resource to create a tcmq queue

Example Usage

```hcl
resource "tencentcloud_tcmq_queue" "queue" {
  queue_name = "queue_name"
}
```

Import

tcmq queue can be imported using the id, e.g.

```
terraform import tencentcloud_tcmq_queue.queue queue_id
```