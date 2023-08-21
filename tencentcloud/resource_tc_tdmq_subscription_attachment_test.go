package tencentcloud

import (
	"context"
	"fmt"
	"strings"
	"testing"

	sdkErrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// go test -i; go test -test.run TestAccTencentCloudTdmqSubscriptionAttachmentResource_basic -v
func TestAccTencentCloudTdmqSubscriptionAttachmentResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		CheckDestroy: testAccCheckTdmqSubscriptionAttachmentDestroy,
		Providers:    testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTdmqSubscriptionAttachment,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_tdmq_subscription_attachment.example", "id"),
				),
			},
			{
				ResourceName:      "tencentcloud_tdmq_subscription_attachment.example",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckTdmqSubscriptionAttachmentDestroy(s *terraform.State) error {
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	service := TdmqService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_tdmq_subscription_attachment" {
			continue
		}

		idSplit := strings.Split(rs.Primary.ID, FILED_SP)
		if len(idSplit) != 4 {
			return fmt.Errorf("id is broken,%s", rs.Primary.ID)
		}

		environmentId := idSplit[0]
		Topic := idSplit[1]
		subscriptionName := idSplit[2]
		clusterId := idSplit[3]

		response, err := service.DescribeTdmqSubscriptionAttachmentById(ctx, environmentId, Topic, subscriptionName, clusterId)

		if err != nil {
			if sdkerr, ok := err.(*sdkErrors.TencentCloudSDKError); ok {
				if sdkerr.Code == "ResourceNotFound.Cluster" {
					return nil
				}
			}
			return err
		}

		if response != nil {
			return fmt.Errorf("tdmq subscription attachment still exist, id: %v", rs.Primary.ID)
		}
	}

	return nil
}

const testAccTdmqSubscriptionAttachment = `
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
  partitions        = 1
  pulsar_topic_type = 3
  remark            = "remark."
}

resource "tencentcloud_tdmq_subscription_attachment" "example" {
  environment_id           = tencentcloud_tdmq_namespace.example.environ_name
  cluster_id               = tencentcloud_tdmq_instance.example.id
  topic_name               = tencentcloud_tdmq_topic.example.topic_name
  subscription_name        = "tf-example-subcription"
  remark                   = "remark."
  auto_create_policy_topic = true
}
`
