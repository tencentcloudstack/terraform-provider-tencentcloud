package tencentcloud

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	sdkErrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
)

func TestAccTencentCloudTdmqRocketmqTopicResource_basic(t *testing.T) {
	t.Parallel()
	terraformId := "tencentcloud_tdmq_rocketmq_topic.topic"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
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
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TdmqRocketmqService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_tdmq_rocketmq_topic" {
			continue
		}
		idSplit := strings.Split(rs.Primary.ID, FILED_SP)
		if len(idSplit) != 4 {
			return fmt.Errorf("id is broken,%s", rs.Primary.ID)
		}
		clusterId := idSplit[0]
		namespaceName := idSplit[1]
		topicName := idSplit[3]

		topicList, err := service.DescribeTdmqRocketmqTopic(ctx, clusterId, namespaceName, topicName)

		if err != nil {
			if sdkerr, ok := err.(*sdkErrors.TencentCloudSDKError); ok {
				if sdkerr.Code == "ResourceUnavailable" {
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
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)

		rs, ok := s.RootModule().Resources[re]
		if !ok {
			return fmt.Errorf("Rocketmq topic  %s is not found", re)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("Rocketmq topic id is not set")
		}

		idSplit := strings.Split(rs.Primary.ID, FILED_SP)
		if len(idSplit) != 4 {
			return fmt.Errorf("id is broken,%s", rs.Primary.ID)
		}
		clusterId := idSplit[0]
		namespaceName := idSplit[1]
		topicName := idSplit[3]

		service := TdmqRocketmqService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
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
resource "tencentcloud_tdmq_rocketmq_cluster" "cluster" {
	cluster_name = "test_rocketmq"
	remark = "test recket mq"
}

resource "tencentcloud_tdmq_rocketmq_namespace" "namespace" {
  cluster_id = tencentcloud_tdmq_rocketmq_cluster.cluster.cluster_id
  namespace_name = "test_namespace"
  ttl = 65000
  retention_time = 65000
  remark = "test namespace"
}

resource "tencentcloud_tdmq_rocketmq_topic" "topic" {
  topic_name = "test_rocketmq_topic"
  namespace_name = tencentcloud_tdmq_rocketmq_namespace.namespace.namespace_name
  type = "Normal"
  cluster_id = tencentcloud_tdmq_rocketmq_cluster.cluster.cluster_id
  remark = "test rocketmq topic"
}
`
