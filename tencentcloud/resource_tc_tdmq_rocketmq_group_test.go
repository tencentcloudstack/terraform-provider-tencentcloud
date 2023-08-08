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

// go test -i; go test -test.run TestAccTencentCloudTdmqRocketmqGroupResource_basic -v
func TestAccTencentCloudTdmqRocketmqGroupResource_basic(t *testing.T) {
	t.Parallel()
	terraformId := "tencentcloud_tdmq_rocketmq_group.example"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckTdmqRocketmqGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTdmqRocketmqGroup,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTdmqRocketmqGroupExists(terraformId),
					resource.TestCheckResourceAttrSet(terraformId, "cluster_id"),
					resource.TestCheckResourceAttrSet(terraformId, "namespace"),
					resource.TestCheckResourceAttrSet(terraformId, "group_name"),
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

func testAccCheckTdmqRocketmqGroupDestroy(s *terraform.State) error {
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TdmqRocketmqService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_tdmq_rocketmq_group" {
			continue
		}
		idSplit := strings.Split(rs.Primary.ID, FILED_SP)
		if len(idSplit) != 3 {
			return fmt.Errorf("id is broken,%s", rs.Primary.ID)
		}
		clusterId := idSplit[0]
		namespaceName := idSplit[1]
		groupName := idSplit[2]

		groupList, err := service.DescribeTdmqRocketmqGroup(ctx, clusterId, namespaceName, groupName)

		if err != nil {
			if sdkerr, ok := err.(*sdkErrors.TencentCloudSDKError); ok {
				if sdkerr.Code == "ResourceUnavailable" || sdkerr.Code == "ResourceNotFound.Cluster" {
					return nil
				}
			}
			return err
		}

		if len(groupList) != 0 {
			return fmt.Errorf("Rocketmq group still exist, id: %v", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckTdmqRocketmqGroupExists(re string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)

		rs, ok := s.RootModule().Resources[re]
		if !ok {
			return fmt.Errorf("Rocketmq group  %s is not found", re)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("Rocketmq group id is not set")
		}

		idSplit := strings.Split(rs.Primary.ID, FILED_SP)
		if len(idSplit) != 3 {
			return fmt.Errorf("id is broken,%s", rs.Primary.ID)
		}
		clusterId := idSplit[0]
		namespaceName := idSplit[1]
		groupName := idSplit[2]

		service := TdmqRocketmqService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
		groupList, err := service.DescribeTdmqRocketmqGroup(ctx, clusterId, namespaceName, groupName)

		if err != nil {
			return err
		}

		if len(groupList) == 0 {
			return fmt.Errorf("Rocketmq group not found, id: %v", rs.Primary.ID)
		}

		return nil
	}
}

const testAccTdmqRocketmqGroup = `
resource "tencentcloud_tdmq_rocketmq_cluster" "example" {
  cluster_name = "tf_example"
  remark       = "remark."
}

resource "tencentcloud_tdmq_rocketmq_namespace" "example" {
  cluster_id     = tencentcloud_tdmq_rocketmq_cluster.example.cluster_id
  namespace_name = "tf_example"
  remark         = "remark."
}

resource "tencentcloud_tdmq_rocketmq_group" "example" {
  group_name       = "tf_example"
  cluster_id       = tencentcloud_tdmq_rocketmq_cluster.example.cluster_id
  namespace        = tencentcloud_tdmq_rocketmq_namespace.example.namespace_name
  read_enable      = true
  broadcast_enable = true
  remark           = "remark."
}
`
