package css_test

import (
	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svccss "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/css"

	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccTencentCloudCssWatermarkRuleAttachmentResource_basic(t *testing.T) {
	t.Parallel()
	baseTime := time.Now().UTC().Add(10 * time.Hour)
	startTime := baseTime.Format(time.RFC3339)
	endTime := baseTime.Add(1 * time.Hour).Format(time.RFC3339)
	liveUrl := "rtmp://5000.liveplay.myqcloud.com/live/stream1"
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckWatermarkRuleAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCssWatermarkRuleAttachment, tcacctest.DefaultCSSLiveType, liveUrl, tcacctest.DefaultCSSDomainName, tcacctest.DefaultCSSAppName, startTime, endTime, tcacctest.DefaultCSSOperator, tcacctest.DefaultCSSPrefix),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCssWatermarkRuleAttachmentExists("tencentcloud_css_watermark_rule_attachment.watermark_rule_attachment"),
					resource.TestCheckResourceAttrSet("tencentcloud_css_watermark_rule_attachment.watermark_rule_attachment", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_css_watermark_rule_attachment.watermark_rule_attachment", "domain_name"),
					resource.TestCheckResourceAttrSet("tencentcloud_css_watermark_rule_attachment.watermark_rule_attachment", "app_name"),
					resource.TestCheckResourceAttrSet("tencentcloud_css_watermark_rule_attachment.watermark_rule_attachment", "stream_name"),
					resource.TestCheckResourceAttrSet("tencentcloud_css_watermark_rule_attachment.watermark_rule_attachment", "template_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_css_watermark_rule_attachment.watermark_rule_attachment", "create_time"),
					resource.TestCheckResourceAttrSet("tencentcloud_css_watermark_rule_attachment.watermark_rule_attachment", "update_time"),
				),
			},
			{
				ResourceName:      "tencentcloud_css_watermark_rule_attachment.watermark_rule_attachment",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckWatermarkRuleAttachmentDestroy(s *terraform.State) error {
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	cssService := svccss.NewCssService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_css_watermark_rule_attachment" {
			continue
		}

		idSplit := strings.Split(rs.Primary.ID, tccommon.FILED_SP)
		domainName := idSplit[0]
		appName := idSplit[1]
		streamName := idSplit[2]
		templateId := idSplit[3]

		rule, err := cssService.DescribeCssWatermarkRuleAttachment(ctx, domainName, appName, streamName, templateId)
		if err != nil {
			return nil
		}

		if rule != nil {
			return fmt.Errorf("css watermark rule attachment still exist, Id: %v", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckCssWatermarkRuleAttachmentExists(re string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := tccommon.GetLogId(tccommon.ContextNil)
		ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

		rs, ok := s.RootModule().Resources[re]
		if !ok {
			return fmt.Errorf("css watermark rule attachment %s is not found", re)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("css watermark rule attachment id is not set")
		}

		cssService := svccss.NewCssService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
		idSplit := strings.Split(rs.Primary.ID, tccommon.FILED_SP)
		domainName := idSplit[0]
		appName := idSplit[1]
		streamName := idSplit[2]
		templateId := idSplit[3]

		rule, err := cssService.DescribeCssWatermarkRuleAttachment(ctx, domainName, appName, streamName, templateId)
		if err != nil {
			return err
		}

		if rule == nil {
			return fmt.Errorf("css watermark rule attachment not found, Id: %v", rs.Primary.ID)
		}
		return nil
	}
}

const testAccCssWatermarkRuleAttachment = `
resource "tencentcloud_css_pull_stream_task" "example" {
	stream_name = "tf_example_stream_name"
	source_type = "%s"
	source_urls = ["%s"]
	domain_name = "%s"
	app_name    = "%s"
	start_time  = "%s"
	end_time    = "%s"
	operator    = "%s"
	comment     = "This is a e2e test case."
  }
  
  resource "tencentcloud_css_watermark" "example" {
	picture_url    = "https://main.qcloudimg.com/raw/c3e0cf113a5c5346b776ecbcfbdcfc72.svg"
	watermark_name = "%swm_rule"
	x_position     = 0
	y_position     = 0
	width          = 0
	height         = 0
  }
  
  resource "tencentcloud_css_watermark_rule_attachment" "watermark_rule_attachment" {
	domain_name = tencentcloud_css_pull_stream_task.example.domain_name
	app_name    = tencentcloud_css_pull_stream_task.example.app_name
	stream_name = tencentcloud_css_pull_stream_task.example.stream_name
	template_id = tencentcloud_css_watermark.example.id
  }

`
