package tencentcloud

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

// go test -i; go test -test.run TestAccTencentCloudRumOfflineLogConfigAttachmentResource_basic -v
func TestAccTencentCloudRumOfflineLogConfigAttachmentResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRumOfflineLogConfigDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRumOfflineLogConfigAttachment,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRumOfflineLogConfigExists("tencentcloud_rum_offline_log_config_attachment.offlineLogConfigAttachment"),
					resource.TestCheckResourceAttr("tencentcloud_rum_offline_log_config_attachment.offlineLogConfigAttachment", "project_key", "e72G2Ulmvyk507X8x5"),
					resource.TestCheckResourceAttr("tencentcloud_rum_offline_log_config_attachment.offlineLogConfigAttachment", "unique_id", "100027012456"),
					resource.TestCheckResourceAttr("tencentcloud_rum_offline_log_config_attachment.offlineLogConfigAttachment", "msg", "success"),
				),
			},
			{
				ResourceName:      "tencentcloud_rum_offline_log_config_attachment.offlineLogConfigAttachment",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckRumOfflineLogConfigDestroy(s *terraform.State) error {
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	service := RumService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_rum_offline_log_config_attachment" {
			continue
		}
		idSplit := strings.Split(rs.Primary.ID, FILED_SP)
		if len(idSplit) != 2 {
			return fmt.Errorf("id is broken,%s", rs.Primary.ID)
		}
		projectKey := idSplit[0]
		uniqueId := idSplit[1]

		logConfig, err := service.DescribeRumOfflineLogConfigAttachment(ctx, projectKey, uniqueId)
		if logConfig != nil && len(logConfig.UniqueIDSet) > 0 {
			return fmt.Errorf("rum logConfig %s still exists", rs.Primary.ID)
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func testAccCheckRumOfflineLogConfigExists(r string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)

		rs, ok := s.RootModule().Resources[r]
		if !ok {
			return fmt.Errorf("resource %s is not found", r)
		}

		idSplit := strings.Split(rs.Primary.ID, FILED_SP)
		if len(idSplit) != 2 {
			return fmt.Errorf("id is broken,%s", rs.Primary.ID)
		}
		projectKey := idSplit[0]
		uniqueId := idSplit[1]

		service := RumService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
		logConfig, err := service.DescribeRumOfflineLogConfigAttachment(ctx, projectKey, uniqueId)
		if logConfig == nil || len(logConfig.UniqueIDSet) < 1 {
			return fmt.Errorf("rum logConfig %s is not found", rs.Primary.ID)
		}
		if err != nil {
			return err
		}

		return nil
	}
}

const testAccRumOfflineLogConfigAttachment = `

resource "tencentcloud_rum_offline_log_config_attachment" "offlineLogConfigAttachment" {
	project_key = "e72G2Ulmvyk507X8x5"
	unique_id = "100027012456"
}

`
