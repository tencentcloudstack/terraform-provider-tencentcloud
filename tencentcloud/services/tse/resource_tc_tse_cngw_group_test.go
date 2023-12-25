package tse_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svctse "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/tse"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	sdkErrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
)

// go test -i; go test -test.run TestAccTencentCloudTseCngwGroupResource_basic -v
func TestAccTencentCloudNeedFixTseCngwGroupResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckTseCngwGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTseCngwGroup,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTseCngwGroupExists("tencentcloud_tse_cngw_group.cngw_group"),
					resource.TestCheckResourceAttrSet("tencentcloud_tse_cngw_group.cngw_group", "id"),
					resource.TestCheckResourceAttr("tencentcloud_tse_cngw_group.cngw_group", "description", "terraform test"),
					resource.TestCheckResourceAttr("tencentcloud_tse_cngw_group.cngw_group", "name", "terraform-group"),
					resource.TestCheckResourceAttr("tencentcloud_tse_cngw_group.cngw_group", "node_config.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_tse_cngw_group.cngw_group", "node_config.0.number", "2"),
					resource.TestCheckResourceAttr("tencentcloud_tse_cngw_group.cngw_group", "node_config.0.specification", "1c2g"),
				),
			},
			{
				ResourceName:      "tencentcloud_tse_cngw_group.cngw_group",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckTseCngwGroupDestroy(s *terraform.State) error {
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	service := svctse.NewTseService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_tse_cngw_gateway" {
			continue
		}

		idSplit := strings.Split(rs.Primary.ID, tccommon.FILED_SP)
		if len(idSplit) != 2 {
			return fmt.Errorf("invalid ID %s", rs.Primary.ID)
		}
		gatewayId := idSplit[0]
		groupId := idSplit[1]

		res, err := service.DescribeTseCngwGroupById(ctx, gatewayId, groupId)
		if err != nil {
			if sdkErr, ok := err.(*sdkErrors.TencentCloudSDKError); ok {
				if sdkErr.Code == "ResourceNotFound.InstanceNotFound" {
					return nil
				}
			}
			return err
		}

		if res != nil {
			return fmt.Errorf("tse gateway %s still exists", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckTseCngwGroupExists(r string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := tccommon.GetLogId(tccommon.ContextNil)
		ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

		rs, ok := s.RootModule().Resources[r]
		if !ok {
			return fmt.Errorf("resource %s is not found", r)
		}

		idSplit := strings.Split(rs.Primary.ID, tccommon.FILED_SP)
		if len(idSplit) != 2 {
			return fmt.Errorf("invalid ID %s", rs.Primary.ID)
		}
		gatewayId := idSplit[0]
		groupId := idSplit[1]

		service := svctse.NewTseService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
		res, err := service.DescribeTseCngwGroupById(ctx, gatewayId, groupId)
		if err != nil {
			return err
		}

		if res == nil {
			return fmt.Errorf("tse gateway %s is not found", rs.Primary.ID)
		}

		return nil
	}
}

const testAccTseCngwGroup = tcacctest.DefaultTseVar + `

resource "tencentcloud_tse_cngw_group" "cngw_group" {
  description = "terraform desc"
  gateway_id  = var.gateway_id
  name        = "terraform-group"
  subnet_id   = "subnet-dwj7ipnc"

  node_config {
    number        = 2
    specification = "1c2g"
  }
}
`
