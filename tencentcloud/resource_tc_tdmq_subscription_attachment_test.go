package tencentcloud

import (
	"context"
	"fmt"
	"strings"
	"testing"

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
					resource.TestCheckResourceAttrSet("tencentcloud_tdmq_subscription_attachment.subscription_attachment", "id"),
				),
			},
			{
				ResourceName:      "tencentcloud_tdmq_subscription_attachment.subscription_attachment",
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
			return err
		}

		if response != nil {
			return fmt.Errorf("tdmq subscription attachment still exist, id: %v", rs.Primary.ID)
		}
	}

	return nil
}

const testAccTdmqSubscriptionAttachment = `
resource "tencentcloud_tdmq_subscription_attachment" "subscription_attachment" {
  environment_id           = "keep-ns"
  topic_name               = "keep-topic"
  subscription_name        = "test-subcription"
  remark                   = "test"
  cluster_id               = "pulsar-9n95ax58b9vn"
  auto_create_policy_topic = true
}
`
