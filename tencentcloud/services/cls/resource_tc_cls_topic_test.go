package cls_test

import (
	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	localcls "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/cls"

	"context"
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func init() {
	resource.AddTestSweepers("tencentcloud_cls_topic", &resource.Sweeper{
		Name: "tencentcloud_cls_topic",
		F:    testSweepClsTopic,
	})
}

func testSweepClsTopic(region string) error {
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	sharedClient, err := tcacctest.SharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("getting tencentcloud client error: %s", err.Error())
	}
	client := sharedClient.(tccommon.ProviderMeta)

	clsService := localcls.NewClsService(client.GetAPIV3Conn())

	instances, err := clsService.DescribeClsTopicByFilter(ctx, nil)
	if err != nil {
		return fmt.Errorf("get instance list error: %s", err.Error())
	}

	for _, v := range instances {
		instanceId := v.TopicId
		instanceName := v.TopicName

		now := time.Now()

		createTime := tccommon.StringToTime(*v.CreateTime)
		interval := now.Sub(createTime).Minutes()
		if strings.HasPrefix(*instanceName, tcacctest.KeepResource) || strings.HasPrefix(*instanceName, tcacctest.DefaultResource) {
			continue
		}
		// less than 30 minute, not delete
		if tccommon.NeedProtect == 1 && int64(interval) < 30 {
			continue
		}

		if err = clsService.DeleteClsTopic(ctx, *instanceId); err != nil {
			log.Printf("[ERROR] sweep instance %s error: %s", *instanceId, err.Error())
		}
	}
	return nil
}

func TestAccTencentCloudClsTopic_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheck(t) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccClsTopic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClsTopicExists("tencentcloud_cls_topic.topic"),
					resource.TestCheckResourceAttr("tencentcloud_cls_topic.topic", "topic_name", "tf-topic-test"),
				),
			},
			{
				ResourceName:      "tencentcloud_cls_topic.topic",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckClsTopicExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := tccommon.GetLogId(tccommon.ContextNil)
		ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("[CHECK][CLS topic][Exists] check: CLB topic %s is not found", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("[CHECK][CLS topic][Exists] check: CLB topic id is not set")
		}
		clsService := localcls.NewClsService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
		instance, err := clsService.DescribeClsTopicById(ctx, rs.Primary.ID)
		if err != nil {
			return err
		}

		if instance == nil {
			return fmt.Errorf("[CHECK][CLS topic][Exists] id %s is not exist", rs.Primary.ID)
		}
		return nil
	}
}

const testAccClsTopic = `
resource "tencentcloud_cls_logset" "logset" {
  logset_name = "tf-topic-test"
  tags        = {
    "test" = "test"
  }
}

resource "tencentcloud_cls_topic" "topic" {
  auto_split           = true
  logset_id            = tencentcloud_cls_logset.logset.id
  max_split_partitions = 20
  partition_count      = 1
  period               = 10
  storage_type         = "hot"
  tags                 = {
    "test" = "test"
  }
  topic_name           = "tf-topic-test"
}
`
