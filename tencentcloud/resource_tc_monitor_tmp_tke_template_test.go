package tencentcloud

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccTencentCloudMonitorTemplate_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckCommon(t, ACCOUNT_TYPE_COMMON) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckTemplateDestroy,
		Steps: []resource.TestStep{
			{
				Config: testTemplate_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTemplateExists("tencentcloud_monitor_tmp_tke_template.basic"),
					resource.TestCheckResourceAttr("tencentcloud_monitor_tmp_tke_template.basic", "template.0.name", "test-template"),
					resource.TestCheckResourceAttr("tencentcloud_monitor_tmp_tke_template.basic", "template.0.level", "instance"),
				),
			},
			{
				Config: testTemplate_update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTemplateExists("tencentcloud_monitor_tmp_tke_template.basic"),
					resource.TestCheckResourceAttr("tencentcloud_monitor_tmp_tke_template.basic", "template.0.name", "test-template_update"),
					resource.TestCheckResourceAttr("tencentcloud_monitor_tmp_tke_template.basic", "template.0.level", "instance"),
				),
			},
			{
				ResourceName:      "tencentcloud_monitor_tmp_tke_template.basic",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckTemplateDestroy(s *terraform.State) error {
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	service := MonitorService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_monitor_tmp_tke_template" {
			continue
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("resource id is not set")
		}

		template, err := service.DescribeTmpTkeTemplateById(ctx, rs.Primary.ID)
		if err != nil {
			return err
		}

		if template != nil {
			return fmt.Errorf("template %s still exists", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckTemplateExists(r string) resource.TestCheckFunc {
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

		service := MonitorService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
		template, err := service.DescribeTmpTkeTemplateById(ctx, rs.Primary.ID)
		if err != nil {
			return err
		}

		if template == nil {
			return fmt.Errorf("template %s is not found", rs.Primary.ID)
		}

		return nil
	}
}

const testTemplate_basic = `
resource "tencentcloud_monitor_tmp_tke_template" "basic" {
  template {
	name	= "test-template"
	level	= "instance"
  }
}`

const testTemplate_update = `
resource "tencentcloud_monitor_tmp_tke_template" "basic" {
  template {
	name	= "test-template_update"
	level	= "instance"
  }
}`
