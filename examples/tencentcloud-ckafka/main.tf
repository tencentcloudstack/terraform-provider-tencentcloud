resource "tencentcloud_ckafka_user" "foo" {
  instance_id  = "ckafka-f9ife4zz"
  account_name = "test"
  password     = "test1234"
}

data "tencentcloud_ckafka_users" "foo" {
  instance_id  = tencentcloud_ckafka_user.foo.instance_id
  account_name = tencentcloud_ckafka_user.foo.account_name
}

resource "tencentcloud_ckafka_acl" foo {
  instance_id     = "ckafka-f9ife4zz"
  resource_type   = "TOPIC"
  resource_name   = "topic-tf-test"
  operation_type  = "WRITE"
  permission_type = "ALLOW"
  host            = "10.10.10.0"
  principal       = tencentcloud_ckafka_user.foo.account_name
}

data "tencentcloud_ckafka_acls" "foo" {
  instance_id   = tencentcloud_ckafka_acl.foo.instance_id
  resource_type = tencentcloud_ckafka_acl.foo.resource_type
  resource_name = tencentcloud_ckafka_acl.foo.resource_name
}

resource "tencentcloud_ckafka_topic" "kafka_topic" {
	instance_id                         = "ckafka-f9ife4zz"
	topic_name                          = "ckafka-topic-tf-test"
	note                                = "this is test ckafka topic"
	replica_num                         = 2
	partition_num                       = 2
	enable_white_list                   = true
	ip_white_list                       = ["192.168.1.1"]
	clean_up_policy                     = "delete"
	sync_replica_min_num                = 1
	unclean_leader_election_enable      = false
	segment                             = 3600000
	retention                           = 60000
	max_message_bytes                   = 0
}

data "tencentcloud_ckafka_topics" "kafka_topics" {
	instance_id						= tencentcloud_ckafka_topic.kafka_topic.instance_id
	topic_name						= tencentcloud_ckafka_topic.kafka_topic.topic_name
}