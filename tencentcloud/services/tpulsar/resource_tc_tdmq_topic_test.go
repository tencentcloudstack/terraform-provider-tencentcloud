package tpulsar_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudTdmqTopicResource_basic -v
func TestAccTencentCloudTdmqTopicResource_basic(t *testing.T) {
	terraformId := "tencentcloud_tdmq_topic.example"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheck(t) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTdmqTopic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(terraformId, "environ_id"),
					resource.TestCheckResourceAttrSet(terraformId, "cluster_id"),
					resource.TestCheckResourceAttr(terraformId, "topic_name", "tf-example-topic"),
					resource.TestCheckResourceAttr(terraformId, "partitions", "6"),
					resource.TestCheckResourceAttr(terraformId, "pulsar_topic_type", "3"),
					resource.TestCheckResourceAttr(terraformId, "remark", "remark."),
				),
			},
			{
				Config: testAccTdmqTopicUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(terraformId, "environ_id"),
					resource.TestCheckResourceAttrSet(terraformId, "cluster_id"),
					resource.TestCheckResourceAttr(terraformId, "topic_name", "tf-example-topic"),
					resource.TestCheckResourceAttr(terraformId, "partitions", "8"),
					resource.TestCheckResourceAttr(terraformId, "pulsar_topic_type", "3"),
					resource.TestCheckResourceAttr(terraformId, "remark", "remark update."),
				),
			},
		},
	})
}

const testAccTdmqTopic = `
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
  environ_id        = tencentcloud_tdmq_namespace.example.environ_name
  cluster_id        = tencentcloud_tdmq_instance.example.id
  topic_name        = "tf-example-topic"
  partitions        = 6
  pulsar_topic_type = 3
  remark            = "remark."
}
`

const testAccTdmqTopicUpdate = `
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
  environ_id        = tencentcloud_tdmq_namespace.example.environ_name
  cluster_id        = tencentcloud_tdmq_instance.example.id
  topic_name        = "tf-example-topic"
  partitions        = 8
  pulsar_topic_type = 3
  remark            = "remark update."
}
`
