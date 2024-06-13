package tpulsar_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudTdmqTopicWithFullIdResource_basic(t *testing.T) {
	t.Parallel()
	terraformId := "tencentcloud_tdmq_topic_with_full_id.example"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheck(t) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTdmqTopicWithFullId,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(terraformId, "cluster_id"),
					resource.TestCheckResourceAttrSet(terraformId, "environ_id"),
					resource.TestCheckResourceAttrSet(terraformId, "topic_name"),
					resource.TestCheckResourceAttrSet(terraformId, "topic_type"),
					resource.TestCheckResourceAttrSet(terraformId, "create_time"),
					resource.TestCheckResourceAttr(terraformId, "partitions", "6"),
					resource.TestCheckResourceAttr(terraformId, "remark", "remark."),
				),
			},
			{
				Config: testAccTdmqTopicWithFullIdUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(terraformId, "partitions", "7"),
					resource.TestCheckResourceAttr(terraformId, "remark", "remark update."),
				),
			},
			{
				ResourceName:      terraformId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccTdmqTopicWithFullId = `
resource "tencentcloud_tdmq_instance" "example" {
	cluster_name = "tf_example_topic_full_id"
	remark       = "remark."
	tags         = {
	  "createdBy" = "terraform"
	}
  }
  
  resource "tencentcloud_tdmq_namespace" "example" {
	environ_name = "tf_example_topic_full_id"
	msg_ttl      = 300
	cluster_id   = tencentcloud_tdmq_instance.example.id
	retention_policy {
	  time_in_minutes = 60
	  size_in_mb      = 10
	}
	remark = "remark."
  }
  
  resource "tencentcloud_tdmq_topic_with_full_id" "example" {
	environ_id        = tencentcloud_tdmq_namespace.example.environ_name
	cluster_id        = tencentcloud_tdmq_instance.example.id
	topic_name        = "tf-example-topic-with-full-id"
	partitions        = 6
	pulsar_topic_type = 3
	remark            = "remark."
  }
`

const testAccTdmqTopicWithFullIdUpdate = `
resource "tencentcloud_tdmq_instance" "example" {
	cluster_name = "tf_example_topic_full_id"
	remark       = "remark."
	tags         = {
	  "createdBy" = "terraform"
	}
  }
  
  resource "tencentcloud_tdmq_namespace" "example" {
	environ_name = "tf_example_topic_full_id"
	msg_ttl      = 300
	cluster_id   = tencentcloud_tdmq_instance.example.id
	retention_policy {
	  time_in_minutes = 60
	  size_in_mb      = 10
	}
	remark = "remark."
  }
  
  resource "tencentcloud_tdmq_topic_with_full_id" "example" {
	environ_id        = tencentcloud_tdmq_namespace.example.environ_name
	cluster_id        = tencentcloud_tdmq_instance.example.id
	topic_name        = "tf-example-topic-with-full-id"
	partitions        = 7
	pulsar_topic_type = 3
	remark            = "remark update."
  }
`
