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
)

// go test -i; go test -test.run TestAccTencentCloudTsfApplicationConfigResource_basic -v
func TestAccTencentCloudTsfApplicationConfigResource_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_TSF) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckTsfApplicationConfigDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTsfApplicationConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTsfApplicationConfigExists("tencentcloud_tsf_application_config.application_config"),
					resource.TestCheckResourceAttrSet("tencentcloud_tsf_application_config.application_config", "id"),
					resource.TestCheckResourceAttr("tencentcloud_tsf_application_config.application_config", "config_name", "tf-test-config"),
					resource.TestCheckResourceAttr("tencentcloud_tsf_application_config.application_config", "config_version", "1.0"),
					resource.TestCheckResourceAttr("tencentcloud_tsf_application_config.application_config", "config_value", "name: \"name\""),
					resource.TestCheckResourceAttr("tencentcloud_tsf_application_config.application_config", "config_version_desc", "version desc"),
					// resource.TestCheckResourceAttr("tencentcloud_tsf_application_config.application_config", "encode_with_base64", "false"),
				),
			},
			// {
			// 	ResourceName:      "tencentcloud_tsf_application_config.application_config",
			// 	ImportState:       true,
			// 	ImportStateVerify: true,
			// },
		},
	})
}

func testAccCheckTsfApplicationConfigDestroy(s *terraform.State) error {
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	service := svctsf.NewTsfService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_tsf_application_config" {
			continue
		}

		res, err := service.DescribeTsfApplicationConfigById(ctx, rs.Primary.ID, "")
		if err != nil {
			return err
		}

		if res != nil {
			return fmt.Errorf("tsf application config %s still exists", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckTsfApplicationConfigExists(r string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := tccommon.GetLogId(tccommon.ContextNil)
		ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

		rs, ok := s.RootModule().Resources[r]
		if !ok {
			return fmt.Errorf("resource %s is not found", r)
		}

		service := svctsf.NewTsfService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
		res, err := service.DescribeTsfApplicationConfigById(ctx, rs.Primary.ID, "")
		if err != nil {
			return err
		}

		if res == nil {
			return fmt.Errorf("tsf application config %s is not found", rs.Primary.ID)
		}

		return nil
	}
}

const testAccTsfApplicationConfigVar = `
variable "application_id" {
	default = "` + tcacctest.DefaultTsfApplicationId + `"
}
`

const testAccTsfApplicationConfig = testAccTsfApplicationConfigVar + `

resource "tencentcloud_tsf_application_config" "application_config" {
	config_name = "tf-test-config"
	config_version = "1.0"
	config_value = "name: \"name\""
	application_id = var.application_id
	config_version_desc = "version desc"
	# config_type = ""
	# encode_with_base64 = false
	# program_id_list =
}

`
