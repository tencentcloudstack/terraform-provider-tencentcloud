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

// go test -i; go test -test.run TestAccTencentCloudTdmqRocketmqNamespaceResource_basic -v
func TestAccTencentCloudTdmqRocketmqNamespaceResource_basic(t *testing.T) {
	t.Parallel()
	terraformId := "tencentcloud_tdmq_rocketmq_namespace.example"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheck(t) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckTdmqRocketmqNamespaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTdmqRocketmqNamespace,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTdmqRocketmqNamespaceExists(terraformId),
					resource.TestCheckResourceAttr(terraformId, "remark", "remark."),
				),
			},
			{
				Config: testAccTdmqRocketmqNamespaceUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTdmqRocketmqNamespaceExists(terraformId),
					resource.TestCheckResourceAttr(terraformId, "remark", "remark update."),
				),
			},
			{
				ResourceName:      "tencentcloud_tdmq_rocketmq_namespace.example",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckTdmqRocketmqNamespaceDestroy(s *terraform.State) error {
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := svctrocket.NewTdmqRocketmqService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_tdmq_rocketmq_namespace" {
			continue
		}
		idSplit := strings.Split(rs.Primary.ID, tccommon.FILED_SP)
		if len(idSplit) != 2 {
			return fmt.Errorf("id is broken,%s", rs.Primary.ID)
		}
		clusterId := idSplit[0]
		namespaceName := idSplit[1]
		namespaces, err := service.DescribeTdmqRocketmqNamespace(ctx, namespaceName, clusterId)
		if err != nil {
			if sdkerr, ok := err.(*sdkErrors.TencentCloudSDKError); ok {
				if sdkerr.Code == "ResourceUnavailable" || sdkerr.Code == "ResourceNotFound.Cluster" {
					return nil
				}
			}
			return err
		}

		if len(namespaces) != 0 {
			return fmt.Errorf("Rocketmq namespace still exist, id: %v", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckTdmqRocketmqNamespaceExists(re string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := tccommon.GetLogId(tccommon.ContextNil)
		ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

		rs, ok := s.RootModule().Resources[re]
		if !ok {
			return fmt.Errorf("Rocketmq namespace  %s is not found", re)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("Rocketmq namespace id is not set")
		}

		idSplit := strings.Split(rs.Primary.ID, tccommon.FILED_SP)
		if len(idSplit) != 2 {
			return fmt.Errorf("id is broken,%s", rs.Primary.ID)
		}
		clusterId := idSplit[0]
		namespaceName := idSplit[1]
		service := svctrocket.NewTdmqRocketmqService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
		namespaces, err := service.DescribeTdmqRocketmqNamespace(ctx, namespaceName, clusterId)
		if err != nil {
			return err
		}

		if len(namespaces) == 0 {
			return fmt.Errorf("Rocketmq namespace not found, id: %v", rs.Primary.ID)
		}

		return nil
	}
}

const testAccTdmqRocketmqNamespace = `
resource "tencentcloud_tdmq_rocketmq_cluster" "example" {
  cluster_name = "tf_example"
  remark       = "remark."
}

resource "tencentcloud_tdmq_rocketmq_namespace" "example" {
  cluster_id     = tencentcloud_tdmq_rocketmq_cluster.example.cluster_id
  namespace_name = "tf_example"
  remark         = "remark."
}
`

const testAccTdmqRocketmqNamespaceUpdate = `
resource "tencentcloud_tdmq_rocketmq_cluster" "example" {
  cluster_name = "tf_example"
  remark       = "remark."
}

resource "tencentcloud_tdmq_rocketmq_namespace" "example" {
  cluster_id     = tencentcloud_tdmq_rocketmq_cluster.example.cluster_id
  namespace_name = "tf_example"
  remark         = "remark update."
}
`
