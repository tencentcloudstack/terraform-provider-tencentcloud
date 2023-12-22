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

// go test -i; go test -test.run TestAccTencentCloudTsfLaneResource_basic -v
func TestAccTencentCloudTsfLaneResource_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_TSF) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckTsfLaneDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTsfLane,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTsfLaneExists("tencentcloud_tsf_lane.lane"),
					resource.TestCheckResourceAttrSet("tencentcloud_tsf_lane.lane", "id"),
					resource.TestCheckResourceAttr("tencentcloud_tsf_lane.lane", "lane_name", "terraform-lane"),
					resource.TestCheckResourceAttr("tencentcloud_tsf_lane.lane", "remark", "lane desc"),
					resource.TestCheckResourceAttr("tencentcloud_tsf_lane.lane", "lane_group_list.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_tsf_lane.lane", "lane_group_list.0.group_id", tcacctest.DefaultTsfGroupId),
					resource.TestCheckResourceAttr("tencentcloud_tsf_lane.lane", "lane_group_list.0.entrance", "true"),
				),
			},
			// {
			// 	ResourceName:      "tencentcloud_tsf_lane.lane",
			// 	ImportState:       true,
			// 	ImportStateVerify: true,
			// },
		},
	})
}

func testAccCheckTsfLaneDestroy(s *terraform.State) error {
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	service := svctsf.NewTsfService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_tsf_lane" {
			continue
		}

		res, err := service.DescribeTsfLaneById(ctx, rs.Primary.ID)
		if err != nil {
			return err
		}

		if res != nil {
			return fmt.Errorf("tsf lane %s still exists", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckTsfLaneExists(r string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := tccommon.GetLogId(tccommon.ContextNil)
		ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

		rs, ok := s.RootModule().Resources[r]
		if !ok {
			return fmt.Errorf("resource %s is not found", r)
		}

		service := svctsf.NewTsfService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
		res, err := service.DescribeTsfLaneById(ctx, rs.Primary.ID)
		if err != nil {
			return err
		}

		if res == nil {
			return fmt.Errorf("tsf lane %s is not found", rs.Primary.ID)
		}

		return nil
	}
}

const testAccTsfLaneVar = `
variable "group_id" {
	default = "` + tcacctest.DefaultTsfGroupId + `"
}
`

const testAccTsfLane = testAccTsfLaneVar + `

resource "tencentcloud_tsf_lane" "lane" {
	lane_name = "terraform-lane"
	remark = "lane desc"
	lane_group_list {
		  group_id = var.group_id
		  entrance = true
	}
}

`
