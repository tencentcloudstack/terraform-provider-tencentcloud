package tencentcloud

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
)

func TestAccTencentCloudTcmqQueueResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		CheckDestroy: testAccCheckTcmqQueueDestroy,
		Providers:    testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTcmqQueue,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTcmqQueueExists("tencentcloud_tcmq_queue.queue"),
					resource.TestCheckResourceAttrSet("tencentcloud_tcmq_queue.queue", "id"),
					resource.TestCheckResourceAttr("tencentcloud_tcmq_queue.queue", "queue_name", "test_queue"),
				),
			},
			{
				ResourceName:            "tencentcloud_tcmq_queue.queue",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"max_receive_count", "max_time_to_live"},
			},
		},
	})
}

func testAccCheckTcmqQueueDestroy(s *terraform.State) error {
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TcmqService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_tcmq_queue" {
			continue
		}

		queue, err := service.DescribeTcmqQueueById(ctx, rs.Primary.ID)
		if queue != nil {
			return fmt.Errorf("TcmqQueue instance still exists")
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

func testAccCheckTcmqQueueExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("TcmqQueue %s is not found", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("TcmqQueue id is not set")
		}

		service := TcmqService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
		queue, err := service.DescribeTcmqQueueById(ctx, rs.Primary.ID)
		if queue == nil {
			return fmt.Errorf("TcmqQueue %s is not found", rs.Primary.ID)
		}
		if err != nil {
			return err
		}

		return nil
	}
}

const testAccTcmqQueue = `
resource "tencentcloud_tcmq_queue" "queue" {
	queue_name="test_queue"
}
`
