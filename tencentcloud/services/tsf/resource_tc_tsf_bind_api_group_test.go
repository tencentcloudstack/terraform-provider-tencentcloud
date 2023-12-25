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

// go test -i; go test -test.run TestAccTencentCloudTsfBindApiGroupResource_basic -v
func TestAccTencentCloudTsfBindApiGroupResource_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_TSF) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckTsfBindApiGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTsfBindApiGroup,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTsfBindApiGroupExists("tencentcloud_tsf_bind_api_group.bind_api_group"),
					resource.TestCheckResourceAttrSet("tencentcloud_tsf_bind_api_group.bind_api_group", "id"),
					resource.TestCheckResourceAttr("tencentcloud_tsf_bind_api_group.bind_api_group", "gateway_deploy_group_id", "group-vzd97zpy"),
					resource.TestCheckResourceAttr("tencentcloud_tsf_bind_api_group.bind_api_group", "group_id", "grp-qp0rj3zi"),
				),
			},
			{
				ResourceName:      "tencentcloud_tsf_bind_api_group.bind_api_group",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckTsfBindApiGroupDestroy(s *terraform.State) error {
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	service := svctsf.NewTsfService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_tsf_bind_api_group" {
			continue
		}
		idSplit := strings.Split(rs.Primary.ID, tccommon.FILED_SP)
		if len(idSplit) != 2 {
			return fmt.Errorf("id is broken,%s", rs.Primary.ID)
		}
		groupId := idSplit[0]
		gatewayDeployGroupId := idSplit[1]

		res, err := service.DescribeTsfBindApiGroupById(ctx, groupId, gatewayDeployGroupId)
		if err != nil {
			return err
		}

		if res != nil {
			return fmt.Errorf("tsf bindApiGroup %s still exists", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckTsfBindApiGroupExists(r string) resource.TestCheckFunc {
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
		groupId := idSplit[0]
		gatewayDeployGroupId := idSplit[1]

		service := svctsf.NewTsfService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
		res, err := service.DescribeTsfBindApiGroupById(ctx, groupId, gatewayDeployGroupId)
		if err != nil {
			return err
		}

		if res == nil {
			return fmt.Errorf("tsf bindApiGroup %s is not found", rs.Primary.ID)
		}

		return nil
	}
}

const testAccTsfBindApiGroup = `

resource "tencentcloud_tsf_bind_api_group" "bind_api_group" {
  gateway_deploy_group_id = "group-vzd97zpy"
  group_id = "grp-qp0rj3zi"
}

`
