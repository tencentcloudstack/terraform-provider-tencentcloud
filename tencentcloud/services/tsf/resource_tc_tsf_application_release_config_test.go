package tsf_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svctsf "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/tsf"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// go test -i; go test -test.run TestAccTencentCloudTsfApplicationReleaseConfigResource_basic -v
func TestAccTencentCloudTsfApplicationReleaseConfigResource_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_TSF) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckTsfApplicationReleaseConfigDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTsfApplicationReleaseConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTsfApplicationReleaseConfigExists("tencentcloud_tsf_application_release_config.application_release_config"),
					resource.TestCheckResourceAttrSet("tencentcloud_tsf_application_release_config.application_release_config", "id"),
					resource.TestCheckResourceAttr("tencentcloud_tsf_application_release_config.application_release_config", "config_id", tcacctest.DefaultTsfConfigId),
					resource.TestCheckResourceAttr("tencentcloud_tsf_application_release_config.application_release_config", "group_id", tcacctest.DefaultTsfGroupId),
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
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	service := svctsf.NewTsfService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_tsf_application_release_config" {
			continue
		}

		idSplit := strings.Split(rs.Primary.ID, tccommon.FILED_SP)
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
		logId := tccommon.GetLogId(tccommon.ContextNil)
		ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

		rs, ok := s.RootModule().Resources[r]
		if !ok {
			return fmt.Errorf("resource %s is not found", r)
		}
		idSplit := strings.Split(rs.Primary.ID, tccommon.FILED_SP)
		if len(idSplit) != 3 {
			return fmt.Errorf("invalid ID %s", rs.Primary.ID)
		}
		configId := idSplit[0]
		groupId := idSplit[1]

		service := svctsf.NewTsfService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
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
	default = "` + tcacctest.DefaultTsfGroupId + `"
}
variable "config_id" {
	default = "` + tcacctest.DefaultTsfConfigId + `"
}
`

const testAccTsfApplicationReleaseConfig = testAccTsfApplicationReleaseConfigVar + `

resource "tencentcloud_tsf_application_release_config" "application_release_config" {
  config_id = var.config_id
  group_id = var.group_id
  release_desc = "terraform_release_desc"
}

`
