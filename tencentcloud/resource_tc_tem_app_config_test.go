package tencentcloud

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
)

// go test -i; go test -test.run TestAccTencentCloudTemAppConfigResource_basic -v
func TestAccTencentCloudTemAppConfigResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckTemAppConfigDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTemAppConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTemAppConfigExists("tencentcloud_tem_app_config.appConfig"),
					resource.TestCheckResourceAttrSet("tencentcloud_tem_app_config.appConfig", "id"),
					resource.TestCheckResourceAttr("tencentcloud_tem_app_config.appConfig", "environment_id", defaultEnvironmentId),
					resource.TestCheckResourceAttr("tencentcloud_tem_app_config.appConfig", "name", "demo"),
					resource.TestCheckResourceAttr("tencentcloud_tem_app_config.appConfig", "config_data.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_tem_app_config.appConfig", "config_data.0.key", "key"),
					resource.TestCheckResourceAttr("tencentcloud_tem_app_config.appConfig", "config_data.0.value", "value"),
				),
			},
			{
				ResourceName:      "tencentcloud_tem_app_config.appConfig",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckTemAppConfigDestroy(s *terraform.State) error {
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	service := TemService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_tem_app_config" {
			continue
		}
		idSplit := strings.Split(rs.Primary.ID, FILED_SP)
		if len(idSplit) != 2 {
			return fmt.Errorf("id is broken,%s", rs.Primary.ID)
		}
		environmentId := idSplit[0]
		name := idSplit[1]

		res, err := service.DescribeTemAppConfig(ctx, environmentId, name)
		if err != nil {
			if sdkErr, ok := err.(*errors.TencentCloudSDKError); ok {
				if sdkErr.Code == "InternalError.DescribeConfigDataError" {
					return nil
				}
			}
			return err
		}

		if res != nil {
			return fmt.Errorf("tem app config %s still exists", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckTemAppConfigExists(r string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)

		rs, ok := s.RootModule().Resources[r]
		if !ok {
			return fmt.Errorf("resource %s is not found", r)
		}

		idSplit := strings.Split(rs.Primary.ID, FILED_SP)
		if len(idSplit) != 2 {
			return fmt.Errorf("id is broken,%s", rs.Primary.ID)
		}
		environmentId := idSplit[0]
		name := idSplit[1]

		service := TemService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
		res, err := service.DescribeTemAppConfig(ctx, environmentId, name)
		if err != nil {
			return err
		}

		if res == nil {
			return fmt.Errorf("tem app config %s is not found", rs.Primary.ID)
		}

		return nil
	}
}

const testAccTemAppConfigVar = `
variable "environment_id" {
  default = "` + defaultEnvironmentId + `"
}
`

const testAccTemAppConfig = testAccTemAppConfigVar + `

resource "tencentcloud_tem_app_config" "appConfig" {
  environment_id = var.environment_id
  name = "demo"
  config_data {
    key = "key"
    value = "value"
  }
}

`
