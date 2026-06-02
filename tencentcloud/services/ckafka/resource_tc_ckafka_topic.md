Use this resource to create ckafka topic.

Example Usage

```hcl
resource "tencentcloud_ckafka_topic" "example" {
  instance_id                    = "ckafka-bzmja3mb"
  topic_name                     = "tf-example"
  note                           = "topic note"
  replica_num                    = 2
  partition_num                  = 3
  clean_up_policy                = "delete"
  sync_replica_min_num           = 2
  unclean_leader_election_enable = false
  segment                        = 86400000
  retention                      = 60000
  max_message_bytes              = 4096
}
```

Import

ckafka topic can be imported using the instance_id#topic_name, e.g.

```
terraform import tencentcloud_ckafka_topic.example ckafka-bzmja3mb#tf-example
```