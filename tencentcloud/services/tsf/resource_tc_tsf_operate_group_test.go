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

// go test -i; go test -test.run TestAccTencentCloudTsfOperateGroupResource_basic -v
func TestAccTencentCloudTsfOperateGroupResource_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_TSF) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckTsfUnitNamespaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTsfOperateGroup,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTsfOperateGroupExists("tencentcloud_tsf_operate_group.operate_group"),
					resource.TestCheckResourceAttrSet("tencentcloud_tsf_operate_group.operate_group", "id"),
				),
			},
			{
				Config: testAccTsfOperateGroupUp,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTsfOperateGroupExists("tencentcloud_tsf_operate_group.operate_group"),
					resource.TestCheckResourceAttrSet("tencentcloud_tsf_operate_group.operate_group", "id"),
				),
			},
		},
	})
}

func testAccCheckTsfOperateGroupExists(r string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := tccommon.GetLogId(tccommon.ContextNil)
		ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

		rs, ok := s.RootModule().Resources[r]
		if !ok {
			return fmt.Errorf("resource %s is not found", r)
		}

		service := svctsf.NewTsfService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
		res, err := service.DescribeTsfStartGroupById(ctx, rs.Primary.ID)
		if err != nil {
			return err
		}

		if res == nil {
			return fmt.Errorf("tsf group %s is not found", rs.Primary.ID)
		}

		if *res.GroupStatus != "Running" && *res.GroupStatus != "Paused" {
			return fmt.Errorf("tsf group %s start or stop operation failed", rs.Primary.ID)
		}

		return nil
	}
}

const testAccTsfOperateGroup = `

resource "tencentcloud_tsf_operate_group" "operate_group" {
	group_id = "group-yrjkln9v"
	operate  = "stop"
}

`

const testAccTsfOperateGroupUp = `

resource "tencentcloud_tsf_operate_group" "operate_group" {
	group_id = "group-yrjkln9v"
	operate  = "start"
}

`
