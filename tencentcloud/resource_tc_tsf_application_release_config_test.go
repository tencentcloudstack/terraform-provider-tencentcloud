package tencentcloud

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// go test -i; go test -test.run TestAccTencentCloudTsfApplicationReleaseConfigResource_basic -v
func TestAccTencentCloudTsfApplicationReleaseConfigResource_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckCommon(t, ACCOUNT_TYPE_TSF) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckTsfApplicationReleaseConfigDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTsfApplicationReleaseConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTsfApplicationReleaseConfigExists("tencentcloud_tsf_application_release_config.application_release_config"),
					resource.TestCheckResourceAttrSet("tencentcloud_tsf_application_release_config.application_release_config", "id"),
					resource.TestCheckResourceAttr("tencentcloud_tsf_application_release_config.application_release_config", "config_id", defaultTsfConfigId),
					resource.TestCheckResourceAttr("tencentcloud_tsf_application_release_config.application_release_config", "group_id", defaultTsfGroupId),
					resource.TestCheckResourceAttr("tencentcloud_tsf_application_release_config.application_release_config", "release_desc", "terraform_release_desc")),
			},
			{
				ResourceName:      "tencentcloud_tsf_application_release_config.application_release_config",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckTsfApplicationReleaseConfigDestroy(s *terraform.State) error {
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	service := TsfService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_tsf_application_release_config" {
			continue
		}

		idSplit := strings.Split(rs.Primary.ID, FILED_SP)
		if len(idSplit) != 3 {
			return fmt.Errorf("invalid ID %s", rs.Primary.ID)
		}
		configId := idSplit[0]
		groupId := idSplit[1]

		res, err := service.DescribeTsfApplicationReleaseConfigById(ctx, configId, groupId)
		if err != nil {
			return err
		}

		if res != nil {
			return fmt.Errorf("tsf ApplicationReleaseConfig %s still exists", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckTsfApplicationReleaseConfigExists(r string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)

		rs, ok := s.RootModule().Resources[r]
		if !ok {
			return fmt.Errorf("resource %s is not found", r)
		}
		idSplit := strings.Split(rs.Primary.ID, FILED_SP)
		if len(idSplit) != 3 {
			return fmt.Errorf("invalid ID %s", rs.Primary.ID)
		}
		configId := idSplit[0]
		groupId := idSplit[1]

		service := TsfService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
		res, err := service.DescribeTsfApplicationReleaseConfigById(ctx, configId, groupId)
		if err != nil {
			return err
		}

		if res == nil {
			return fmt.Errorf("tsf ApplicationReleaseConfig %s is not found", rs.Primary.ID)
		}

		return nil
	}
}

const testAccTsfApplicationReleaseConfigVar = `
variable "group_id" {
	default = "` + defaultTsfGroupId + `"
}
variable "config_id" {
	default = "` + defaultTsfConfigId + `"
}
`

const testAccTsfApplicationReleaseConfig = testAccTsfApplicationReleaseConfigVar + `

resource "tencentcloud_tsf_application_release_config" "application_release_config" {
  config_id = var.config_id
  group_id = var.group_id
  release_desc = "terraform_release_desc"
}

`
