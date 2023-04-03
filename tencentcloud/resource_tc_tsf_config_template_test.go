package tencentcloud

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	sdkErrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
)

// go test -i; go test -test.run TestAccTencentCloudTsfConfigTemplateResource_basic -v
func TestAccTencentCloudTsfConfigTemplateResource_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckCommon(t, ACCOUNT_TYPE_TSF) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckTsfConfigTemplateDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTsfConfigTemplate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTsfConfigTemplateExists("tencentcloud_tsf_config_template.config_template"),
					resource.TestCheckResourceAttrSet("tencentcloud_tsf_config_template.config_template", "id"),
					resource.TestCheckResourceAttr("tencentcloud_tsf_config_template.config_template", "config_template_name", "terraform-template-name"),
					resource.TestCheckResourceAttr("tencentcloud_tsf_config_template.config_template", "config_template_type", "Ribbon"),
					resource.TestCheckResourceAttrSet("tencentcloud_tsf_config_template.config_template", "config_template_value"),
					resource.TestCheckResourceAttr("tencentcloud_tsf_config_template.config_template", "config_template_desc", "terraform-test"),
				),
			},
			// {
			// 	ResourceName:      "tencentcloud_tsf_config_template.config_template",
			// 	ImportState:       true,
			// 	ImportStateVerify: true,
			// },
			{
				Config: testAccTsfConfigTemplateUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTsfConfigTemplateExists("tencentcloud_tsf_config_template.config_template"),
					resource.TestCheckResourceAttrSet("tencentcloud_tsf_config_template.config_template", "id"),
					resource.TestCheckResourceAttr("tencentcloud_tsf_config_template.config_template", "config_template_name", "terraform-template-name"),
					resource.TestCheckResourceAttr("tencentcloud_tsf_config_template.config_template", "config_template_type", "Ribbon"),
					resource.TestCheckResourceAttrSet("tencentcloud_tsf_config_template.config_template", "config_template_value"),
					resource.TestCheckResourceAttr("tencentcloud_tsf_config_template.config_template", "config_template_desc", "terraform-test"),
				),
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
			code := err.(*sdkErrors.TencentCloudSDKError).Code
			if code == "FailedOperation.ConfigTemplateImportFailed" {
				return nil
			}
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
	config_template_name = "terraform-template-name"
	config_template_type = "Ribbon"
	config_template_value = <<-EOT
	  ribbon.ReadTimeout: 5000
	  ribbon.ConnectTimeout: 2000
	  ribbon.MaxAutoRetries: 0
	  ribbon.MaxAutoRetriesNextServer: 1
	  ribbon.OkToRetryOnAllOperations: true
	EOT
	config_template_desc = "terraform-test"
}

`

const testAccTsfConfigTemplateUpdate = `

resource "tencentcloud_tsf_config_template" "config_template" {
	config_template_name = "terraform-template-name"
	config_template_type = "Ribbon"
	config_template_value = <<-EOT
	  ribbon.ReadTimeout: 5000
	  ribbon.ConnectTimeout: 2000
	  ribbon.MaxAutoRetries: 0
	  ribbon.MaxAutoRetriesNextServer: 1
	  ribbon.OkToRetryOnAllOperations: false
	EOT
	config_template_desc = "terraform-test"
}

`
