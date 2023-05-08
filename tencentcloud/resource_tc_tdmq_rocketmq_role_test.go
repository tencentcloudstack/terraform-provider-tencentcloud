package tencentcloud

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	sdkErrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
)

func TestAccTencentCloudTdmqRocketmqRoleResource_basic(t *testing.T) {
	t.Parallel()
	terraformId := "tencentcloud_tdmq_rocketmq_role.role"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckTdmqRocketmqRoleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTdmqRocketmqRole,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTdmqRocketmqRoleExists(terraformId),
					resource.TestCheckResourceAttr(terraformId, "role_name", "test_rocketmq_role"),
					resource.TestCheckResourceAttr(terraformId, "remark", "test rocketmq role"),
				),
			},
			{
				Config: testAccTdmqRocketmqRoleUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTdmqRocketmqRoleExists(terraformId),
					resource.TestCheckResourceAttr(terraformId, "role_name", "test_rocketmq_role"),
					resource.TestCheckResourceAttr(terraformId, "remark", "test rocketmq role update"),
				),
			},
			{
				ResourceName:      "tencentcloud_tdmq_rocketmq_role.role",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckTdmqRocketmqRoleDestroy(s *terraform.State) error {
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TdmqRocketmqService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_tdmq_rocketmq_role" {
			continue
		}
		idSplit := strings.Split(rs.Primary.ID, FILED_SP)
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
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)

		rs, ok := s.RootModule().Resources[re]
		if !ok {
			return fmt.Errorf("Rocketmq role  %s is not found", re)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("Rocketmq role id is not set")
		}

		idSplit := strings.Split(rs.Primary.ID, FILED_SP)
		if len(idSplit) != 2 {
			return fmt.Errorf("id is broken,%s", rs.Primary.ID)
		}
		clusterId := idSplit[0]
		roleName := idSplit[1]
		service := TdmqRocketmqService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
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
resource "tencentcloud_tdmq_rocketmq_cluster" "cluster" {
	cluster_name = "test_rocketmq"
	remark = "test recket mq"
}

resource "tencentcloud_tdmq_rocketmq_role" "role" {
  role_name = "test_rocketmq_role"
  remark = "test rocketmq role"
  cluster_id = tencentcloud_tdmq_rocketmq_cluster.cluster.cluster_id
}
`

const testAccTdmqRocketmqRoleUpdate = `
resource "tencentcloud_tdmq_rocketmq_cluster" "cluster" {
	cluster_name = "test_rocketmq"
	remark = "test recket mq"
}

resource "tencentcloud_tdmq_rocketmq_role" "role" {
  role_name = "test_rocketmq_role"
  remark = "test rocketmq role update"
  cluster_id = tencentcloud_tdmq_rocketmq_cluster.cluster.cluster_id
}
`
