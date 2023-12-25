package tcmq_test

import (
	"context"
	"fmt"
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svctcmq "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/tcmq"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
)

func TestAccTencentCloudTcmqTopicResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		CheckDestroy: testAccCheckTcmqTopicDestroy,
		Providers:    tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTcmqTopic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTcmqTopicExists("tencentcloud_tcmq_topic.topic"),
					resource.TestCheckResourceAttrSet("tencentcloud_tcmq_topic.topic", "id"),
					resource.TestCheckResourceAttr("tencentcloud_tcmq_topic.topic", "topic_name", "test_topic"),
				),
			},
			{
				ResourceName:            "tencentcloud_tcmq_topic.topic",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"trace"},
			},
		},
	})
}

func testAccCheckTcmqTopicDestroy(s *terraform.State) error {
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := svctcmq.NewTcmqService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_tcmq_topic" {
			continue
		}
		topic, err := service.DescribeTcmqTopicById(ctx, rs.Primary.ID)
		if topic != nil {
			return fmt.Errorf("TcmqTopic instance still exists")
		}
		if err != nil {
			if e, ok := err.(*errors.TencentCloudSDKError); ok {
				if e.GetCode() == "ResourceNotFound" {
					return nil
				}
			}
			return err
		}
	}
	return nil
}

func testAccCheckTcmqTopicExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := tccommon.GetLogId(tccommon.ContextNil)
		ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("TcmqTopic %s is not found", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("TcmqTopic id is not set")
		}

		service := svctcmq.NewTcmqService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
		topic, err := service.DescribeTcmqTopicById(ctx, rs.Primary.ID)
		if topic == nil {
			return fmt.Errorf("TcmqTopic %s is not found", rs.Primary.ID)
		}
		if err != nil {
			return err
		}

		return nil
	}
}

const testAccTcmqTopic = `
resource "tencentcloud_tcmq_topic" "topic" {
	topic_name = "test_topic"
}
`
