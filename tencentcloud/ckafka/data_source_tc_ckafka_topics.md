Use this data source to query detailed information of ckafka topic.

Example Usage

```hcl
resource "tencentcloud_ckafka_topic" "foo" {
	instance_id                     = "ckafka-f9ife4zz"
	topic_name                      = "example"
	note                            = "topic note"
	replica_num                     = 2
	partition_num                   = 1
	enable_white_list               = true
	ip_white_list                   = ["ip1","ip2"]
	clean_up_policy                 = "delete"
	sync_replica_min_num            = 1
	unclean_leader_election_enable  = false
	segment                         = 3600000
	retention                       = 60000
	max_message_bytes               = 1024
}
```