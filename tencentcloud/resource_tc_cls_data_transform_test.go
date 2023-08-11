package tencentcloud

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// go test -i; go test -test.run TestAccTencentCloudClsDataTransformResource_basic -v
func TestAccTencentCloudClsDataTransformResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		CheckDestroy: testAccCheckClsDataTransformDestroy,
		Providers:    testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccClsDataTransform,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClsDataTransformExists("tencentcloud_cls_data_transform.data_transform"),
					resource.TestCheckResourceAttrSet("tencentcloud_cls_data_transform.data_transform", "func_type"),
					resource.TestCheckResourceAttrSet("tencentcloud_cls_data_transform.data_transform", "src_topic_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_cls_data_transform.data_transform", "name"),
					resource.TestCheckResourceAttrSet("tencentcloud_cls_data_transform.data_transform", "etl_content"),
					resource.TestCheckResourceAttrSet("tencentcloud_cls_data_transform.data_transform", "task_type"),
					resource.TestCheckResourceAttrSet("tencentcloud_cls_data_transform.data_transform", "enable_flag"),
					resource.TestCheckResourceAttrSet("tencentcloud_cls_data_transform.data_transform", "dst_resources.#"),
					resource.TestCheckResourceAttrSet("tencentcloud_cls_data_transform.data_transform", "dst_resources.0.topic_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_cls_data_transform.data_transform", "dst_resources.0.alias")),
			},
		},
	})
}

func testAccCheckClsDataTransformDestroy(s *terraform.State) error {
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	clsService := ClsService{
		client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn,
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_cls_data_transform" {
			continue
		}
		time.Sleep(5 * time.Second)
		instance, err := clsService.DescribeClsDataTransformById(ctx, rs.Primary.ID)
		if err != nil {
			continue
		}
		if instance != nil && instance.TaskId != nil && *instance.TaskId == rs.Primary.ID {
			return fmt.Errorf("[CHECK][CLS dataTransform][Destroy] check: CLS dataTransform still exists: %s", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckClsDataTransformExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("[CHECK][CLS dataTransform][Exists] check: CLS dataTransform %s is not found", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("[CHECK][CLS dataTransform][Create] check: CLS dataTransform id is not set")
		}
		clsService := ClsService{
			client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn,
		}
		taskRes, err := clsService.DescribeClsDataTransformById(ctx, rs.Primary.ID)
		if err != nil {
			return err
		}
		if taskRes == nil {
			return fmt.Errorf("[CHECK][CLS redirection][Exists] id %s is not exist", rs.Primary.ID)
		}
		return nil
	}
}

const testAccClsDataTransform = `
resource "tencentcloud_cls_logset" "logset_src" {
  logset_name = "tf-example-src"
  tags = {
    "createdBy" = "terraform"
  }
}
resource "tencentcloud_cls_topic" "topic_src" {
  topic_name           = "tf-example_src"
  logset_id            = tencentcloud_cls_logset.logset_src.id
  auto_split           = false
  max_split_partitions = 20
  partition_count      = 1
  period               = 10
  storage_type         = "hot"
  tags                 = {
    "test" = "test",
  }
}
resource "tencentcloud_cls_logset" "logset_dst" {
  logset_name = "tf-example-dst"
  tags = {
    "createdBy" = "terraform"
  }
}
resource "tencentcloud_cls_topic" "topic_dst" {
  topic_name           = "tf-example-dst"
  logset_id            = tencentcloud_cls_logset.logset_dst.id
  auto_split           = false
  max_split_partitions = 20
  partition_count      = 1
  period               = 10
  storage_type         = "hot"
  tags                 = {
    "test" = "test",
  }
}
resource "tencentcloud_cls_data_transform" "data_transform" {
  func_type = 1
  src_topic_id = tencentcloud_cls_topic.topic_src.id
  name = "iac-test-src"
  etl_content = "ext_sep(\"content\", \"f1, f2, f3\", sep=\",\", quote=\"\", restrict=False, mode=\"overwrite\")fields_drop(\"content\")"
  task_type = 3
  enable_flag = 1
  dst_resources {
    topic_id = tencentcloud_cls_topic.topic_dst.id
    alias = "iac-test-dst"

  }
}

`
