package tencentcloud

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccTencentCloudClbInstanceTopic(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckClbListenerRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccClbInstanceTopic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClbInstanceTopicExists("tencentcloud_clb_log_topic.topic"),
					resource.TestCheckResourceAttr("tencentcloud_clb_log_topic.topic", "topic_name", "clb-topic-test"),
				),
			},
		},
	})
}

func testAccCheckClbInstanceTopicExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("[CHECK][CLB topic][Exists] check: CLB topic %s is not found", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("[CHECK][CLB topic][Exists] check: CLB topic id is not set")
		}
		clsService := ClsService{
			client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn,
		}
		instance, err := clsService.DescribeClsTopicById(ctx, rs.Primary.ID)
		if err != nil {
			return err
		}

		if instance == nil {
			return fmt.Errorf("[CHECK][CLB topic][Exists] id %s is not exist", rs.Primary.ID)
		}
		return nil
	}
}

const testAccClbInstanceTopic = `
resource "tencentcloud_clb_log_set" "set1" {
    period = 7
}

resource "tencentcloud_clb_log_topic" "topic" {
    log_set_id = tencentcloud_clb_log_set.set1.id
    topic_name="clb-topic-test"
}
`
