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

// go test -i; go test -test.run TestAccTencentCloudTdmqRocketmqEnvironmentRoleResource_basic -v
func TestAccTencentCloudTdmqRocketmqEnvironmentRoleResource_basic(t *testing.T) {
	t.Parallel()
	terraformId := "tencentcloud_tdmq_rocketmq_environment_role.example"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckTdmqRocketmqEnvironmentRoleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTdmqRocketmqEnvironmentRole,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTdmqRocketmqEnvironmentRoleExists(terraformId),
					resource.TestCheckResourceAttr(terraformId, "role_name", "tf_example_role"),
					resource.TestCheckResourceAttr(terraformId, "permissions.#", "2"),
				),
			},
			{
				Config: testAccTdmqRocketmqEnvironmentUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTdmqRocketmqEnvironmentRoleExists(terraformId),
					resource.TestCheckResourceAttr(terraformId, "permissions.#", "1"),
				),
			},
			{
				ResourceName:      terraformId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckTdmqRocketmqEnvironmentRoleDestroy(s *terraform.State) error {
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TdmqRocketmqService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_tdmq_rocketmq_environment_role" {
			continue
		}
		idSplit := strings.Split(rs.Primary.ID, FILED_SP)
		if len(idSplit) != 3 {
			return fmt.Errorf("id is broken,%s", rs.Primary.ID)
		}
		clusterId := idSplit[0]
		roleName := idSplit[1]
		environmentName := idSplit[2]

		environmentRoles, err := service.DescribeTdmqRocketmqEnvironmentRole(ctx, clusterId, roleName, environmentName)

		if err != nil {
			if sdkerr, ok := err.(*sdkErrors.TencentCloudSDKError); ok {
				if sdkerr.Code == "ResourceNotFound.Cluster" {
					return nil
				}
			}
			return err
		}

		if len(environmentRoles) != 0 {
			return fmt.Errorf("Rocketmq environment role still exist, id: %v", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckTdmqRocketmqEnvironmentRoleExists(re string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)

		rs, ok := s.RootModule().Resources[re]
		if !ok {
			return fmt.Errorf("Rocketmq environment role  %s is not found", re)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("Rocketmq environment role id is not set")
		}

		idSplit := strings.Split(rs.Primary.ID, FILED_SP)
		if len(idSplit) != 3 {
			return fmt.Errorf("id is broken,%s", rs.Primary.ID)
		}
		clusterId := idSplit[0]
		roleName := idSplit[1]
		environmentName := idSplit[2]

		service := TdmqRocketmqService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
		environmentRoles, err := service.DescribeTdmqRocketmqEnvironmentRole(ctx, clusterId, roleName, environmentName)

		if err != nil {
			return err
		}

		if len(environmentRoles) == 0 {
			return fmt.Errorf("Rocketmq environment role not found, id: %v", rs.Primary.ID)
		}

		return nil
	}
}

const testAccTdmqRocketmqEnvironmentRole = `
resource "tencentcloud_tdmq_rocketmq_cluster" "example" {
  cluster_name = "tf_example"
  remark       = "remark."
}

resource "tencentcloud_tdmq_rocketmq_role" "example" {
  role_name  = "tf_example_role"
  remark     = "remark."
  cluster_id = tencentcloud_tdmq_rocketmq_cluster.example.cluster_id
}

resource "tencentcloud_tdmq_rocketmq_namespace" "example" {
  cluster_id     = tencentcloud_tdmq_rocketmq_cluster.example.cluster_id
  namespace_name = "tf_example_namespace"
  remark         = "remark."
}

resource "tencentcloud_tdmq_rocketmq_environment_role" "example" {
  environment_name = tencentcloud_tdmq_rocketmq_namespace.example.namespace_name
  role_name        = tencentcloud_tdmq_rocketmq_role.example.role_name
  permissions      = ["produce", "consume"]
  cluster_id       = tencentcloud_tdmq_rocketmq_cluster.example.cluster_id
}
`

const testAccTdmqRocketmqEnvironmentUpdate = `
resource "tencentcloud_tdmq_rocketmq_cluster" "example" {
  cluster_name = "tf_example"
  remark       = "remark."
}

resource "tencentcloud_tdmq_rocketmq_role" "example" {
  role_name  = "tf_example_role"
  remark     = "remark."
  cluster_id = tencentcloud_tdmq_rocketmq_cluster.example.cluster_id
}

resource "tencentcloud_tdmq_rocketmq_namespace" "example" {
  cluster_id     = tencentcloud_tdmq_rocketmq_cluster.example.cluster_id
  namespace_name = "tf_example_namespace"
  remark         = "remark."
}

resource "tencentcloud_tdmq_rocketmq_environment_role" "example" {
  environment_name = tencentcloud_tdmq_rocketmq_namespace.example.namespace_name
  role_name        = tencentcloud_tdmq_rocketmq_role.example.role_name
  permissions      = ["produce"]
  cluster_id       = tencentcloud_tdmq_rocketmq_cluster.example.cluster_id
}
`
