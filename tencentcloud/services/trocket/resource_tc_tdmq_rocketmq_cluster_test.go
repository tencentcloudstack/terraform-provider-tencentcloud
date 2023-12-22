package trocket_test

import (
	"context"
	"fmt"
	"log"
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svctrocket "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/trocket"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	sdkErrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
)

func TestAccTencentCloudTdmqRocketmqClusterResource_basic(t *testing.T) {
	t.Parallel()
	terraformId := "tencentcloud_tdmq_rocketmq_cluster.example"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheck(t) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckTdmqRocketmqClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTdmqRocketmqCluster,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTdmqRocketmqClusterExists(terraformId),
					resource.TestCheckResourceAttrSet(terraformId, "id"),
					resource.TestCheckResourceAttr(terraformId, "cluster_name", "tf_example"),
					resource.TestCheckResourceAttr(terraformId, "remark", "remark."),
				),
			},
			{
				Config: testAccTdmqRocketmqClusterUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTdmqRocketmqClusterExists(terraformId),
					resource.TestCheckResourceAttrSet(terraformId, "id"),
					resource.TestCheckResourceAttr(terraformId, "cluster_name", "tf_example_update"),
					resource.TestCheckResourceAttr(terraformId, "remark", "remark."),
				),
			},
			{
				ResourceName:      "tencentcloud_tdmq_rocketmq_cluster.example",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckTdmqRocketmqClusterDestroy(s *terraform.State) error {
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := svctrocket.NewTdmqRocketmqService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
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
		logId := tccommon.GetLogId(tccommon.ContextNil)
		ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

		rs, ok := s.RootModule().Resources[re]
		if !ok {
			return fmt.Errorf("Rocketmq instance  %s is not found", re)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("Rocketmq instance id is not set")
		}
		log.Printf("exist id: %v", rs.Primary.ID)

		service := svctrocket.NewTdmqRocketmqService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
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
resource "tencentcloud_tdmq_rocketmq_cluster" "example" {
  cluster_name = "tf_example"
  remark       = "remark."
}
`
const testAccTdmqRocketmqClusterUpdate = `
resource "tencentcloud_tdmq_rocketmq_cluster" "example" {
  cluster_name = "tf_example_update"
  remark       = "remark."
}
`
