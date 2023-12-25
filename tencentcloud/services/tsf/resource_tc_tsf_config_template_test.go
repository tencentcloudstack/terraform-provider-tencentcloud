package tsf_test

import (
	"context"
	"fmt"
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svctsf "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/tsf"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	sdkErrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
)

// go test -i; go test -test.run TestAccTencentCloudTsfConfigTemplateResource_basic -v
func TestAccTencentCloudTsfConfigTemplateResource_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_TSF) },
		Providers:    tcacctest.AccProviders,
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
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	service := svctsf.NewTsfService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
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
		logId := tccommon.GetLogId(tccommon.ContextNil)
		ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

		rs, ok := s.RootModule().Resources[r]
		if !ok {
			return fmt.Errorf("resource %s is not found", r)
		}

		service := svctsf.NewTsfService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
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
