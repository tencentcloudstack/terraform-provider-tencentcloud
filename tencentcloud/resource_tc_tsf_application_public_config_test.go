package tencentcloud

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// go test -i; go test -test.run TestAccTencentCloudTsfApplicationPublicConfigResource_basic -v
func TestAccTencentCloudTsfApplicationPublicConfigResource_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckCommon(t, ACCOUNT_TYPE_TSF) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckTsfApplicationPublicConfigDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTsfApplicationPublicConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTsfApplicationPublicConfigExists("tencentcloud_tsf_application_public_config.application_public_config"),
					resource.TestCheckResourceAttrSet("tencentcloud_tsf_application_public_config.application_public_config", "id"),
					resource.TestCheckResourceAttr("tencentcloud_tsf_application_public_config.application_public_config", "config_name", "terraform-test"),
					resource.TestCheckResourceAttr("tencentcloud_tsf_application_public_config.application_public_config", "config_type", "P"),
					resource.TestCheckResourceAttr("tencentcloud_tsf_application_public_config.application_public_config", "config_value", "test: 1"),
					resource.TestCheckResourceAttr("tencentcloud_tsf_application_public_config.application_public_config", "config_version", "1.0"),
					resource.TestCheckResourceAttr("tencentcloud_tsf_application_public_config.application_public_config", "config_version_desc", "product version"),
					resource.TestCheckResourceAttr("tencentcloud_tsf_application_public_config.application_public_config", "encode_with_base64", "true"),
				),
			},
			// {
			// 	ResourceName:      "tencentcloud_tsf_application_public_config.application_public_config",
			// 	ImportState:       true,
			// 	ImportStateVerify: true,
			// },
		},
	})
}

func testAccCheckTsfApplicationPublicConfigDestroy(s *terraform.State) error {
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	service := TsfService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_tsf_application_release_config" {
			continue
		}

		res, err := service.DescribeTsfApplicationPublicConfigById(ctx, rs.Primary.ID)
		if err != nil {
			return err
		}

		if res != nil {
			return fmt.Errorf("tsf ApplicationPublicConfig %s still exists", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckTsfApplicationPublicConfigExists(r string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)

		rs, ok := s.RootModule().Resources[r]
		if !ok {
			return fmt.Errorf("resource %s is not found", r)
		}

		service := TsfService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
		res, err := service.DescribeTsfApplicationPublicConfigById(ctx, rs.Primary.ID)
		if err != nil {
			return err
		}

		if res == nil {
			return fmt.Errorf("tsf ApplicationPublicConfig %s is not found", rs.Primary.ID)
		}

		return nil
	}
}

const testAccTsfApplicationPublicConfig = `

resource "tencentcloud_tsf_application_public_config" "application_public_config" {
	config_name = "terraform-test"
	config_version = "1.0"
	config_value = "test: 1"
	config_version_desc = "product version"
	config_type = "P"
	encode_with_base64 = true
	# program_id_list =
}

`
