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

// go test -i; go test -test.run TestAccTencentCloudTdmqRocketmqTopicResource_basic -v
func TestAccTencentCloudTdmqRocketmqTopicResource_basic(t *testing.T) {
	t.Parallel()
	terraformId := "tencentcloud_tdmq_rocketmq_topic.example"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheck(t) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckTdmqRocketmqTopicDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTdmqRocketmqTopic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTdmqRocketmqTopicExists(terraformId),
					resource.TestCheckResourceAttrSet(terraformId, "id"),
					resource.TestCheckResourceAttrSet(terraformId, "cluster_id"),
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

func testAccCheckTdmqRocketmqTopicDestroy(s *terraform.State) error {
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := svctrocket.NewTdmqRocketmqService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_tdmq_rocketmq_topic" {
			continue
		}
		idSplit := strings.Split(rs.Primary.ID, tccommon.FILED_SP)
		if len(idSplit) != 4 {
			return fmt.Errorf("id is broken,%s", rs.Primary.ID)
		}
		clusterId := idSplit[0]
		namespaceName := idSplit[1]
		topicName := idSplit[3]

		topicList, err := service.DescribeTdmqRocketmqTopic(ctx, clusterId, namespaceName, topicName)

		if err != nil {
			if sdkerr, ok := err.(*sdkErrors.TencentCloudSDKError); ok {
				if sdkerr.Code == "ResourceUnavailable" || sdkerr.Code == "ResourceNotFound.Cluster" {
					return nil
				}
			}
			return err
		}

		if len(topicList) != 0 {
			return fmt.Errorf("Rocketmq topic still exist, id: %v", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckTdmqRocketmqTopicExists(re string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := tccommon.GetLogId(tccommon.ContextNil)
		ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

		rs, ok := s.RootModule().Resources[re]
		if !ok {
			return fmt.Errorf("Rocketmq topic  %s is not found", re)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("Rocketmq topic id is not set")
		}

		idSplit := strings.Split(rs.Primary.ID, tccommon.FILED_SP)
		if len(idSplit) != 4 {
			return fmt.Errorf("id is broken,%s", rs.Primary.ID)
		}
		clusterId := idSplit[0]
		namespaceName := idSplit[1]
		topicName := idSplit[3]

		service := svctrocket.NewTdmqRocketmqService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
		topicList, err := service.DescribeTdmqRocketmqTopic(ctx, clusterId, namespaceName, topicName)

		if err != nil {
			return err
		}

		if len(topicList) == 0 {
			return fmt.Errorf("Rocketmq topic not found, id: %v", rs.Primary.ID)
		}

		return nil
	}
}

const testAccTdmqRocketmqTopic = `
resource "tencentcloud_tdmq_rocketmq_cluster" "example" {
  cluster_name = "tf_example"
  remark       = "remark."
}

resource "tencentcloud_tdmq_rocketmq_namespace" "example" {
  cluster_id     = tencentcloud_tdmq_rocketmq_cluster.example.cluster_id
  namespace_name = "tf_example_namespace"
  remark         = "remark."
}

resource "tencentcloud_tdmq_rocketmq_topic" "example" {
  topic_name     = "tf_example"
  namespace_name = tencentcloud_tdmq_rocketmq_namespace.example.namespace_name
  cluster_id     = tencentcloud_tdmq_rocketmq_cluster.example.cluster_id
  type           = "Normal"
  remark         = "remark."
}
`
