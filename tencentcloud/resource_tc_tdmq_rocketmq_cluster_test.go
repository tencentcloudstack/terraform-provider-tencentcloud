package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	sdkErrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
)

func TestAccTencentCloudTdmqRocketmqClusterResource_basic(t *testing.T) {
	t.Parallel()
	terraformId := "tencentcloud_tdmq_rocketmq_cluster.cluster"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckTdmqRocketmqClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTdmqRocketmqCluster,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTdmqRocketmqClusterExists(terraformId),
					resource.TestCheckResourceAttrSet(terraformId, "id"),
					resource.TestCheckResourceAttr(terraformId, "cluster_name", "test_rocketmq"),
					resource.TestCheckResourceAttr(terraformId, "remark", "test rocket mq"),
				),
			},
			{
				Config: testAccTdmqRocketmqClusterUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTdmqRocketmqClusterExists(terraformId),
					resource.TestCheckResourceAttrSet(terraformId, "id"),
					resource.TestCheckResourceAttr(terraformId, "cluster_name", "test_rocketmq_update"),
					resource.TestCheckResourceAttr(terraformId, "remark", "test rocket update"),
				),
			},
			{
				ResourceName:      "tencentcloud_tdmq_rocketmq_cluster.cluster",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckTdmqRocketmqClusterDestroy(s *terraform.State) error {
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TdmqRocketmqService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_tdmq_rocketmq_cluster" {
			continue
		}
		log.Printf("destroy id: %v", rs.Primary.ID)
		cluster, err := service.DescribeTdmqRocketmqCluster(ctx, rs.Primary.ID)

		if err != nil {
			if sdkerr, ok := err.(*sdkErrors.TencentCloudSDKError); ok {
				if sdkerr.Code == "ResourceNotFound.Cluster" {
					return nil
				}
			}
			return err
		}

		if cluster != nil {
			return fmt.Errorf("Rocketmq instance still exist, id: %v", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckTdmqRocketmqClusterExists(re string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)

		rs, ok := s.RootModule().Resources[re]
		if !ok {
			return fmt.Errorf("Rocketmq instance  %s is not found", re)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("Rocketmq instance id is not set")
		}
		log.Printf("exist id: %v", rs.Primary.ID)

		service := TdmqRocketmqService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
		cluster, err := service.DescribeTdmqRocketmqCluster(ctx, rs.Primary.ID)
		if err != nil {
			return err
		}

		if cluster == nil {
			return fmt.Errorf("Rocketmq instance not found, id: %v", rs.Primary.ID)
		}

		return nil
	}
}

const testAccTdmqRocketmqCluster = `
resource "tencentcloud_tdmq_rocketmq_cluster" "cluster" {
	cluster_name = "test_rocketmq"
	remark = "test rocket mq"
}
`
const testAccTdmqRocketmqClusterUpdate = `
resource "tencentcloud_tdmq_rocketmq_cluster" "cluster" {
	cluster_name = "test_rocketmq_update"
	remark = "test rocket update"
}
`
