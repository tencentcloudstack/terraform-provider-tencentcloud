package tpulsar_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

// go test -i; go test -test.run TestAccTencentCloudTdmqSubscriptionResource_basic -v
func TestAccTencentCloudTdmqSubscriptionResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTdmqSubscription,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_tdmq_subscription.example", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_tdmq_subscription.example", "cluster_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_tdmq_subscription.example", "environment_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_tdmq_subscription.example", "topic_name"),
					resource.TestCheckResourceAttr("tencentcloud_tdmq_subscription.example", "subscription_name", "tf-example-subscription"),
					resource.TestCheckResourceAttr("tencentcloud_tdmq_subscription.example", "remark", "remark."),
				),
			},
			{
				ResourceName:      "tencentcloud_tdmq_subscription.example",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccTdmqSubscription = `
resource "tencentcloud_tdmq_instance" "example" {
  cluster_name = "tf_example"
  remark       = "remark."
  tags         = {
    "createdBy" = "terraform"
  }
}

resource "tencentcloud_tdmq_namespace" "example" {
  environ_name = "tf_example"
  msg_ttl      = 300
  cluster_id   = tencentcloud_tdmq_instance.example.id
  retention_policy {
    time_in_minutes = 60
    size_in_mb      = 10
  }
  remark = "remark."
}

resource "tencentcloud_tdmq_topic" "example" {
  cluster_id        = tencentcloud_tdmq_instance.example.id
  environ_id        = tencentcloud_tdmq_namespace.example.environ_name
  topic_name        = "tf-example-topic"
  partitions        = 1
  pulsar_topic_type = 3
  remark            = "remark."
}

resource "tencentcloud_tdmq_subscription" "example" {
  cluster_id               = tencentcloud_tdmq_instance.example.id
  environment_id           = tencentcloud_tdmq_namespace.example.environ_name
  topic_name               = tencentcloud_tdmq_topic.example.topic_name
  subscription_name        = "tf-example-subscription"
  remark                   = "remark."
  auto_create_policy_topic = true
  auto_delete_policy_topic = true
}
`
