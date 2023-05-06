package tencentcloud

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccTencentCloudMonitorTempAttachment_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckCommon(t, ACCOUNT_TYPE_COMMON) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckTempAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testTempAttachment_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTempAttachmentExists("tencentcloud_monitor_tmp_tke_template_attachment.basic"),
					resource.TestCheckResourceAttr("tencentcloud_monitor_tmp_tke_template_attachment.basic", "template_id", "temp-gqunlvo1"),
					resource.TestCheckResourceAttr("tencentcloud_monitor_tmp_tke_template_attachment.basic", "targets.0.instance_id", "prom-1lspn8sw"),
					resource.TestCheckResourceAttr("tencentcloud_monitor_tmp_tke_template_attachment.basic", "targets.0.region", "ap-guangzhou"),
				),
			},
		},
	})
}

func testAccCheckTempAttachmentDestroy(s *terraform.State) error {
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	recordService := MonitorService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_monitor_tmp_tke_template_attachment" {
			continue
		}
		items := strings.Split(rs.Primary.ID, FILED_SP)
		if len(items) != 3 {
			return fmt.Errorf("invalid ID %s", rs.Primary.ID)
		}

		templateId := items[0]
		instanceId := items[1]
		region := items[2]
		targets, err := recordService.DescribePrometheusTempSync(ctx, templateId)
		if err != nil {
			return err
		}

		if len(targets) > 0 {
			for _, v := range targets {
				if *v.InstanceId == instanceId && *v.Region == region {
					return fmt.Errorf("associated instance information %s still exists", rs.Primary.ID)
				}
			}
		}
	}

	return nil
}

func testAccCheckTempAttachmentExists(r string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)

		rs, ok := s.RootModule().Resources[r]
		if !ok {
			return fmt.Errorf("resource %s is not found", r)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("resource id is not set")
		}
		items := strings.Split(rs.Primary.ID, FILED_SP)
		if len(items) != 3 {
			return fmt.Errorf("invalid ID %s", rs.Primary.ID)
		}

		templateId := items[0]
		instanceId := items[1]
		region := items[2]
		service := MonitorService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
		targets, err := service.DescribePrometheusTempSync(ctx, templateId)
		if err != nil {
			return err
		}

		if len(targets) < 1 {
			return fmt.Errorf("associated instance information %s is not found", rs.Primary.ID)
		}
		for i, v := range targets {
			if *v.InstanceId == instanceId && *v.Region == region {
				return nil
			}
			if i == len(targets)-1 {
				return fmt.Errorf("associated instance information %s is not found", rs.Primary.ID)
			}
		}

		return nil
	}
}

const testTempAttachmentVar = `
variable "prometheus_id" {
  default = "` + defaultPrometheusId + `"
}
variable "template_id" {
  default = "` + defaultTemplateId + `"
}
variable "region" {
  default = "ap-guangzhou"
}`

const testTempAttachment_basic = testTempAttachmentVar + `
resource "tencentcloud_monitor_tmp_tke_template_attachment" "basic" {
  template_id  = var.template_id

  targets {
    region      = var.region
    instance_id = var.prometheus_id
  }
}`
