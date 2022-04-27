package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccTencentCloudCkafkaAclsDataSource(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckCommon(t, ACCOUNT_TYPE_PREPAY) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCkafkaAclDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTencentCloudDataSourceCkafkaAcl,
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckCkafkaAclExists("tencentcloud_ckafka_acl.foo"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_ckafka_acls.foo", "acl_list.0.operation_type"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_ckafka_acls.foo", "acl_list.0.permission_type"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_ckafka_acls.foo", "acl_list.0.resource_name"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_ckafka_acls.foo", "acl_list.0.resource_type"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_ckafka_acls.foo", "acl_list.0.host"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_ckafka_acls.foo", "acl_list.0.principal"),
				),
			},
		},
	})
}

const testAccTencentCloudDataSourceCkafkaAcl = defaultKafkaVariable + `
resource "tencentcloud_ckafka_user" "foo" {
	instance_id  = var.instance_id
	account_name = "tf-test-acl-data"
	password     = "test1234"
}

resource "tencentcloud_ckafka_topic" "kafka_topic_acl" {
	instance_id                     = var.instance_id
	topic_name                      = "ckafka-topic-acl-data-test"
	replica_num                     = 2
	partition_num                   = 1
	note                            = "test topic"
	enable_white_list               = true
	ip_white_list                   = ["192.168.1.1"]
	clean_up_policy                 = "delete"
	sync_replica_min_num            = 1
	unclean_leader_election_enable  = false
	segment                         = 86400000
	retention                       = 60000
}

resource "tencentcloud_ckafka_acl" foo {
  instance_id     = var.instance_id
  resource_type   = "TOPIC"
  resource_name   = tencentcloud_ckafka_topic.kafka_topic_acl.topic_name
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
`
