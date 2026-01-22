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
	client := sharedClient.(tccommon.ProviderMeta).GetAPIV3Conn()

	clsService := localcls.NewClsService(client)

	instances, err := clsService.DescribeClsTopicByFilter(ctx, nil)
	if err != nil {
		return fmt.Errorf("get instance list error: %s", err.Error())
	}

	// add scanning resources
	var resources, nonKeepResources []*tccommon.ResourceInstance
	for _, v := range instances {
		if !tccommon.CheckResourcePersist(*v.TopicName, *v.CreateTime) {
			nonKeepResources = append(nonKeepResources, &tccommon.ResourceInstance{
				Id:   *v.TopicId,
				Name: *v.TopicName,
			})
		}
		resources = append(resources, &tccommon.ResourceInstance{
			Id:         *v.TopicId,
			Name:       *v.TopicName,
			CreateTime: *v.CreateTime,
		})
	}
	tccommon.ProcessScanCloudResources(client, resources, nonKeepResources, "CreateTopic")

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

// go test -i; go test -test.run TestAccTencentCloudClsTopic_basic -v
func TestAccTencentCloudClsTopic_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheck(t) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccClsTopic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClsTopicExists("tencentcloud_cls_topic.example"),
					resource.TestCheckResourceAttr("tencentcloud_cls_topic.example", "topic_name", "tf_example"),
					resource.TestCheckResourceAttr("tencentcloud_cls_topic.example", "storage_type", "hot"),
					resource.TestCheckResourceAttr("tencentcloud_cls_topic.example", "describes", "Test Demo."),
				),
			},
			{
				ResourceName:      "tencentcloud_cls_topic.example",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccClsTopicUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClsTopicExists("tencentcloud_cls_topic.example"),
					resource.TestCheckResourceAttr("tencentcloud_cls_topic.example", "topic_name", "tf_example_update"),
					resource.TestCheckResourceAttr("tencentcloud_cls_topic.example", "storage_type", "hot"),
					resource.TestCheckResourceAttr("tencentcloud_cls_topic.example", "describes", "Test Demo Update."),
				),
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
resource "tencentcloud_cls_logset" "example" {
  logset_name = "tf_example"
  tags        = {
    "demo" = "test"
  }
}

resource "tencentcloud_cls_topic" "example" {
  topic_name           = "tf_example"
  logset_id            = tencentcloud_cls_logset.example.id
  auto_split           = false
  max_split_partitions = 20
  partition_count      = 1
  period               = 30
  storage_type         = "hot"
  describes            = "Test Demo."
  hot_period           = 10
  tags                 = {
    "test" = "test",
  }
}
`

const testAccClsTopicUpdate = `
resource "tencentcloud_cls_logset" "example" {
  logset_name = "tf_example"
  tags        = {
    "demo" = "test"
  }
}

resource "tencentcloud_cls_topic" "example" {
  topic_name           = "tf_example_update"
  logset_id            = tencentcloud_cls_logset.example.id
  auto_split           = false
  max_split_partitions = 20
  partition_count      = 1
  period               = 30
  storage_type         = "hot"
  describes            = "Test Demo Update."
  hot_period           = 10
}
`
