package tencentcloud

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	sdkErrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
)

func TestAccTencentCloudTsfApplicationFileConfigResource_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckTsfApplicationFileConfigDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTsfApplicationFileConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTsfApplicationFileConfigExists("tencentcloud_tsf_application_file_config.application_file_config"),
					resource.TestCheckResourceAttrSet("tencentcloud_tsf_application_file_config.application_file_config", "id"),
					resource.TestCheckResourceAttr("tencentcloud_tsf_application_file_config.application_file_config", "config_name", "terraform-test"),
					resource.TestCheckResourceAttr("tencentcloud_tsf_application_file_config.application_file_config", "config_version", "v1"),
					resource.TestCheckResourceAttr("tencentcloud_tsf_application_file_config.application_file_config", "config_file_name", "terraform-config-name"),
					resource.TestCheckResourceAttr("tencentcloud_tsf_application_file_config.application_file_config", "config_file_value", "ZWNobyAidGVzdCI="),
					resource.TestCheckResourceAttr("tencentcloud_tsf_application_file_config.application_file_config", "application_id", "application-ym9mxmza"),
					resource.TestCheckResourceAttr("tencentcloud_tsf_application_file_config.application_file_config", "config_file_path", "/etc/test/"),
					resource.TestCheckResourceAttr("tencentcloud_tsf_application_file_config.application_file_config", "config_version_desc", "terraform test version"),
					resource.TestCheckResourceAttr("tencentcloud_tsf_application_file_config.application_file_config", "config_file_code", "utf-8"),
					resource.TestCheckResourceAttr("tencentcloud_tsf_application_file_config.application_file_config", "config_post_cmd", "echo \"test1\""),
					resource.TestCheckResourceAttr("tencentcloud_tsf_application_file_config.application_file_config", "encode_with_base64", "true"),
					resource.TestCheckResourceAttr("tencentcloud_tsf_application_file_config.application_file_config", "tags.createdBy", "terraform"),
				),
			},
			{
				ResourceName:      "tencentcloud_tsf_application_file_config.application_file_config",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckTsfApplicationFileConfigDestroy(s *terraform.State) error {
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	service := TsfService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_tsf_application_file_config" {
			continue
		}

		res, err := service.DescribeTsfGroupById(ctx, rs.Primary.ID)
		if err != nil {
			ee, ok := err.(*sdkErrors.TencentCloudSDKError)
			if !ok {
				return err
			}
			if ee.Code == "ResourceNotFound.GroupNotExist" {
				return nil
			}
			return err
		}

		if res != nil {
			return fmt.Errorf("tsf application file config %s still exists", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckTsfApplicationFileConfigExists(r string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)

		rs, ok := s.RootModule().Resources[r]
		if !ok {
			return fmt.Errorf("resource %s is not found", r)
		}

		service := TsfService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
		res, err := service.DescribeTsfGroupById(ctx, rs.Primary.ID)
		if err != nil {
			return err
		}

		if res == nil {
			return fmt.Errorf("tsf application file config %s is not found", rs.Primary.ID)
		}

		return nil
	}
}

const testAccTsfApplicationFileConfig = `

resource "tencentcloud_tsf_application_file_config" "application_file_config" {
	config_name = "terraform-test"
	config_version = "v1"
	config_file_name = "terraform-config-name"
	config_file_value = "ZWNobyAidGVzdCI="
	application_id = "application-ym9mxmza"
	config_file_path = "/etc/test/"
	config_version_desc = "terraform test version"
	config_file_code = "utf-8"
	config_post_cmd = "echo \"test1\""
	encode_with_base64 = true
	# program_id_list =
	tags = {
	  "createdBy" = "terraform"
	}
}

`
