package tencentcloud

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccTencentCloudNeedFixTsfConfigTemplateResource_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckTsfConfigTemplateDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTsfConfigTemplate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTsfConfigTemplateExists("tencentcloud_tsf_config_template.config_template"),
					resource.TestCheckResourceAttrSet("tencentcloud_tsf_config_template.config_template", "id"),
					resource.TestCheckResourceAttr("tencentcloud_tsf_config_template.config_template", "config_template_name", ""),
					resource.TestCheckResourceAttr("tencentcloud_tsf_config_template.config_template", "config_template_type", ""),
				),
			},
			{
				ResourceName:      "tencentcloud_tsf_config_template.config_template",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckTsfConfigTemplateDestroy(s *terraform.State) error {
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	service := TsfService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_tsf_config_template" {
			continue
		}

		res, err := service.DescribeTsfConfigTemplateById(ctx, rs.Primary.ID)
		if err != nil {
			return err
		}

		if res != nil {
			return fmt.Errorf("tsf ConfigTemplate %s still exists", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckTsfConfigTemplateExists(r string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)

		rs, ok := s.RootModule().Resources[r]
		if !ok {
			return fmt.Errorf("resource %s is not found", r)
		}

		service := TsfService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
		res, err := service.DescribeTsfConfigTemplateById(ctx, rs.Primary.ID)
		if err != nil {
			return err
		}

		if res == nil {
			return fmt.Errorf("tsf ConfigTemplate %s is not found", rs.Primary.ID)
		}

		return nil
	}
}

const testAccTsfConfigTemplate = `

resource "tencentcloud_tsf_config_template" "config_template" {
  config_template_name = ""
  config_template_type = ""
  config_template_value = ""
  config_template_desc = ""
  program_id_list = 
      }

`
