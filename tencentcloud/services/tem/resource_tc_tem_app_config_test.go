package tem_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svctem "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/tem"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
)

// go test -i; go test -test.run TestAccTencentCloudTemAppConfigResource_basic -v
func TestAccTencentCloudTemAppConfigResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckTemAppConfigDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTemAppConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTemAppConfigExists("tencentcloud_tem_app_config.appConfig"),
					resource.TestCheckResourceAttrSet("tencentcloud_tem_app_config.appConfig", "id"),
					resource.TestCheckResourceAttr("tencentcloud_tem_app_config.appConfig", "environment_id", tcacctest.DefaultEnvironmentId),
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
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	service := svctem.NewTemService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_tem_app_config" {
			continue
		}
		idSplit := strings.Split(rs.Primary.ID, tccommon.FILED_SP)
		if len(idSplit) != 2 {
			return fmt.Errorf("id is broken,%s", rs.Primary.ID)
		}
		environmentId := idSplit[0]
		name := idSplit[1]

		res, err := service.DescribeTemAppConfig(ctx, environmentId, name)
		if err != nil {
			if sdkErr, ok := err.(*errors.TencentCloudSDKError); ok {
				if sdkErr.Code == "ResourceNotFound.ConfigDataNotFound" {
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
		logId := tccommon.GetLogId(tccommon.ContextNil)
		ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

		rs, ok := s.RootModule().Resources[r]
		if !ok {
			return fmt.Errorf("resource %s is not found", r)
		}

		idSplit := strings.Split(rs.Primary.ID, tccommon.FILED_SP)
		if len(idSplit) != 2 {
			return fmt.Errorf("id is broken,%s", rs.Primary.ID)
		}
		environmentId := idSplit[0]
		name := idSplit[1]

		service := svctem.NewTemService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
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
  default = "` + tcacctest.DefaultEnvironmentId + `"
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
