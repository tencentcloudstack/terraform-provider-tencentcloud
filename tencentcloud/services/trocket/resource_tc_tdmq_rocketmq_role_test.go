package trocket_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svctrocket "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/trocket"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	sdkErrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
)

// go test -i; go test -test.run TestAccTencentCloudTdmqRocketmqRoleResource_basic -v
func TestAccTencentCloudTdmqRocketmqRoleResource_basic(t *testing.T) {
	t.Parallel()
	terraformId := "tencentcloud_tdmq_rocketmq_role.example"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheck(t) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckTdmqRocketmqRoleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTdmqRocketmqRole,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTdmqRocketmqRoleExists(terraformId),
					resource.TestCheckResourceAttr(terraformId, "role_name", "tf_example"),
					resource.TestCheckResourceAttr(terraformId, "remark", "remark."),
				),
			},
			{
				Config: testAccTdmqRocketmqRoleUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTdmqRocketmqRoleExists(terraformId),
					resource.TestCheckResourceAttr(terraformId, "role_name", "tf_example"),
					resource.TestCheckResourceAttr(terraformId, "remark", "remark update."),
				),
			},
			{
				ResourceName:      "tencentcloud_tdmq_rocketmq_role.example",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckTdmqRocketmqRoleDestroy(s *terraform.State) error {
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := svctrocket.NewTdmqRocketmqService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_tdmq_rocketmq_role" {
			continue
		}
		idSplit := strings.Split(rs.Primary.ID, tccommon.FILED_SP)
		if len(idSplit) != 2 {
			return fmt.Errorf("id is broken,%s", rs.Primary.ID)
		}
		clusterId := idSplit[0]
		roleName := idSplit[1]
		role, err := service.DescribeTdmqRocketmqRole(ctx, clusterId, roleName)
		if err != nil {
			if sdkerr, ok := err.(*sdkErrors.TencentCloudSDKError); ok {
				if sdkerr.Code == "ResourceNotFound.Cluster" {
					return nil
				}
			}
			return err
		}

		if role != nil {
			return fmt.Errorf("Rocketmq role still exist, id: %v", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckTdmqRocketmqRoleExists(re string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := tccommon.GetLogId(tccommon.ContextNil)
		ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

		rs, ok := s.RootModule().Resources[re]
		if !ok {
			return fmt.Errorf("Rocketmq role  %s is not found", re)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("Rocketmq role id is not set")
		}

		idSplit := strings.Split(rs.Primary.ID, tccommon.FILED_SP)
		if len(idSplit) != 2 {
			return fmt.Errorf("id is broken,%s", rs.Primary.ID)
		}
		clusterId := idSplit[0]
		roleName := idSplit[1]
		service := svctrocket.NewTdmqRocketmqService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
		role, err := service.DescribeTdmqRocketmqRole(ctx, clusterId, roleName)

		if err != nil {
			return err
		}

		if role == nil {
			return fmt.Errorf("Rocketmq role not found, id: %v", rs.Primary.ID)
		}

		return nil
	}
}

const testAccTdmqRocketmqRole = `
resource "tencentcloud_tdmq_rocketmq_cluster" "example" {
  cluster_name = "tf_example"
  remark       = "remark."
}

resource "tencentcloud_tdmq_rocketmq_role" "example" {
  cluster_id = tencentcloud_tdmq_rocketmq_cluster.example.cluster_id
  role_name  = "tf_example"
  remark     = "remark."
}
`

const testAccTdmqRocketmqRoleUpdate = `
resource "tencentcloud_tdmq_rocketmq_cluster" "example" {
  cluster_name = "tf_example"
  remark       = "remark."
}

resource "tencentcloud_tdmq_rocketmq_role" "example" {
  cluster_id = tencentcloud_tdmq_rocketmq_cluster.example.cluster_id
  role_name  = "tf_example"
  remark     = "remark update."
}
`
