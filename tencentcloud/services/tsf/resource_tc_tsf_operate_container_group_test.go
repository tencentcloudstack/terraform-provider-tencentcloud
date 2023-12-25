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

// go test -i; go test -test.run TestAccTencentCloudTsfOperateContainerGroupResource_basic -v
func TestAccTencentCloudTsfOperateContainerGroupResource_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_TSF) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTsfOperateContainerGroup,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTsfOperateContainerGroupExists("tencentcloud_tsf_operate_container_group.operate_container_group"),
					resource.TestCheckResourceAttrSet("tencentcloud_tsf_operate_container_group.operate_container_group", "id"),
				),
			},
			{
				Config: testAccTsfOperateContainerGroupUp,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTsfOperateContainerGroupExists("tencentcloud_tsf_operate_container_group.operate_container_group"),
					resource.TestCheckResourceAttrSet("tencentcloud_tsf_operate_container_group.operate_container_group", "id"),
				),
			},
		},
	})
}

func testAccCheckTsfOperateContainerGroupExists(r string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := tccommon.GetLogId(tccommon.ContextNil)
		ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

		rs, ok := s.RootModule().Resources[r]
		if !ok {
			return fmt.Errorf("resource %s is not found", r)
		}

		service := svctsf.NewTsfService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
		res, err := service.DescribeTsfStartContainerGroupById(ctx, rs.Primary.ID)
		if err != nil {
			return err
		}

		if res == nil {
			return fmt.Errorf("tsf container group %s is not found", rs.Primary.ID)
		}

		if *res.Status != "Running" && *res.Status != "Paused" {
			return fmt.Errorf("tsf container group %s start or stop operation failed", rs.Primary.ID)
		}

		return nil
	}
}

const testAccTsfOperateContainerGroup = `

resource "tencentcloud_tsf_operate_container_group" "operate_container_group" {
  group_id = "group-ympdpdzy"
  operate = "stop"
}

`
const testAccTsfOperateContainerGroupUp = `

resource "tencentcloud_tsf_operate_container_group" "operate_container_group" {
  group_id = "group-ympdpdzy"
  operate = "start"
}

`
