package tencentcloud

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

// go test -i; go test -test.run TestAccTencentCloudTsfApplicationConfigResource_basic -v
func TestAccTencentCloudTsfApplicationConfigResource_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckCommon(t, ACCOUNT_TYPE_TSF) },
		Providers:    testAccProviders,
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
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	service := TsfService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
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
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)

		rs, ok := s.RootModule().Resources[r]
		if !ok {
			return fmt.Errorf("resource %s is not found", r)
		}

		service := TsfService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
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
	default = "` + defaultTsfApplicationId + `"
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
