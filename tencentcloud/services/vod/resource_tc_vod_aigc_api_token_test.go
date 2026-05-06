package vod_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	svcvod "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/vod"
)

func TestAccTencentCloudVodAigcApiTokenResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheck(t) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckVodAigcApiTokenDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVodAigcApiTokenBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVodAigcApiTokenExists("tencentcloud_vod_aigc_api_token.example"),
					resource.TestCheckResourceAttrSet("tencentcloud_vod_aigc_api_token.example", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_vod_aigc_api_token.example", "sub_app_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_vod_aigc_api_token.example", "api_token"),
				),
			},
			{
				ResourceName:      "tencentcloud_vod_aigc_api_token.example",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckVodAigcApiTokenDestroy(s *terraform.State) error {
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	vodService := svcvod.NewVodService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_vod_aigc_api_token" {
			continue
		}
		idSplit := strings.Split(rs.Primary.ID, tccommon.FILED_SP)
		if len(idSplit) != 2 {
			return fmt.Errorf("vod aigc api token id is malformed: %s", rs.Primary.ID)
		}
		subAppId := uint64(helper.StrToInt(idSplit[0]))
		apiToken := idSplit[1]

		exists, err := vodService.DescribeVodAigcApiTokenById(ctx, subAppId, apiToken)
		if err != nil {
			return err
		}
		if exists {
			return fmt.Errorf("vod aigc api token still exists: %s", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckVodAigcApiTokenExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := tccommon.GetLogId(tccommon.ContextNil)
		ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("vod aigc api token %s is not found", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("vod aigc api token id is not set")
		}
		idSplit := strings.Split(rs.Primary.ID, tccommon.FILED_SP)
		if len(idSplit) != 2 {
			return fmt.Errorf("vod aigc api token id is malformed: %s", rs.Primary.ID)
		}
		subAppId := uint64(helper.StrToInt(idSplit[0]))
		apiToken := idSplit[1]

		vodService := svcvod.NewVodService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
		exists, err := vodService.DescribeVodAigcApiTokenById(ctx, subAppId, apiToken)
		if err != nil {
			return err
		}
		if !exists {
			return fmt.Errorf("vod aigc api token doesn't exist: %s", rs.Primary.ID)
		}
		return nil
	}
}

const testAccVodAigcApiTokenBasic = `
resource "tencentcloud_vod_sub_application" "sub_application" {
  name        = "aigc-api-token-subapp"
  status      = "On"
  description = "sub application for tencentcloud_vod_aigc_api_token acceptance test"
}

resource "tencentcloud_vod_aigc_api_token" "example" {
  sub_app_id = tonumber(split("#", tencentcloud_vod_sub_application.sub_application.id)[1])
}
`
