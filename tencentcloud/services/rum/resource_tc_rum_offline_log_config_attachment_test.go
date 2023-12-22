package rum_test

import (
	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svcrum "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/rum"

	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// go test -i; go test -test.run TestAccTencentCloudRumOfflineLogConfigAttachmentResource_basic -v
func TestAccTencentCloudRumOfflineLogConfigAttachmentResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheck(t) },
		Providers:    tcacctest.AccProviders,
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
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	service := svcrum.NewRumService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_rum_offline_log_config_attachment" {
			continue
		}
		idSplit := strings.Split(rs.Primary.ID, tccommon.FILED_SP)
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
		logId := tccommon.GetLogId(tccommon.ContextNil)
		ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

		rs, ok := s.RootModule().Resources[r]
		if !ok {
			return fmt.Errorf("resource %s is not found", r)
		}

		idSplit := strings.Split(rs.Primary.ID, tccommon.FILED_SP)
		if len(idSplit) != 2 {
			return fmt.Errorf("id is broken,%s", rs.Primary.ID)
		}
		projectKey := idSplit[0]
		uniqueId := idSplit[1]

		service := svcrum.NewRumService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
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
