package tencentcloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"testing"
	"time"
)

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
					resource.TestCheckResourceAttrSet("tencentcloud_cls_data_transform.data_transform", "dst_resources"),
					resource.TestCheckResourceAttrSet("tencentcloud_cls_data_transform.data_transform.dst_resources", "topic_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_cls_data_transform.data_transform.dst_resources", "alias")),
			},
			{
				ResourceName:      "tencentcloud_cls_data_transform.data_transform",
				ImportState:       true,
				ImportStateVerify: true,
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

resource "tencentcloud_cls_data_transform" "data_transform" {
  func_type = 1
  src_topic_id = "ff2a8b6e-c9ee-4774-9bb5-bdfb4c92a16c"
  name = "iac-test-src"
  etl_content = "ext_sep(\"content\", \"f1, f2, f3\", sep=\",\", quote=\"\", restrict=False, mode=\"overwrite\")fields_drop(\"content\")"
  task_type = 3
  enable_flag = 1
  dst_resources {
    topic_id = "b0e8d897-6e2e-4ef0-bca2-c79228bc0b45"
    alias = "iac-test-dst"

  }
}

`
