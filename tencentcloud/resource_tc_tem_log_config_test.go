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

// go test -i; go test -test.run TestAccTencentCloudTemLogConfigResource_basic -v
func TestAccTencentCloudTemLogConfigResource_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckTemLogConfigDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTemLogConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTemLogConfigExists("tencentcloud_tem_log_config.logConfig"),
					resource.TestCheckResourceAttrSet("tencentcloud_tem_log_config.logConfig", "id"),
					resource.TestCheckResourceAttr("tencentcloud_tem_log_config.logConfig", "environment_id", defaultEnvironmentId),
					resource.TestCheckResourceAttr("tencentcloud_tem_log_config.logConfig", "application_id", defaultApplicationId),
					resource.TestCheckResourceAttr("tencentcloud_tem_log_config.logConfig", "name", "terraform-test"),
					resource.TestCheckResourceAttr("tencentcloud_tem_log_config.logConfig", "logset_id", defaultLogsetId),
					resource.TestCheckResourceAttr("tencentcloud_tem_log_config.logConfig", "topic_id", defaultTopicId),
					resource.TestCheckResourceAttr("tencentcloud_tem_log_config.logConfig", "input_type", "container_stdout"),
					resource.TestCheckResourceAttr("tencentcloud_tem_log_config.logConfig", "log_type", "minimalist_log"),
				),
			},
			{
				ResourceName:      "tencentcloud_tem_log_config.logConfig",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckTemLogConfigDestroy(s *terraform.State) error {
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	service := TemService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_tem_log_config" {
			continue
		}

		idSplit := strings.Split(rs.Primary.ID, FILED_SP)
		if len(idSplit) != 3 {
			return fmt.Errorf("id is broken,%s", rs.Primary.ID)
		}
		environmentId := idSplit[0]
		applicationId := idSplit[1]
		name := idSplit[2]

		res, err := service.DescribeTemLogConfig(ctx, environmentId, applicationId, name)
		if err != nil {
			ee, ok := err.(*sdkErrors.TencentCloudSDKError)
			if !ok {
				return err
			}
			if ee.Code == "ResourceNotFound.LogConfigNotFound" {
				return nil
			}
			return err
		}

		if res != nil {
			return fmt.Errorf("tem log config %s still exists", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckTemLogConfigExists(r string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)

		rs, ok := s.RootModule().Resources[r]
		if !ok {
			return fmt.Errorf("resource %s is not found", r)
		}
		idSplit := strings.Split(rs.Primary.ID, FILED_SP)
		if len(idSplit) != 3 {
			return fmt.Errorf("id is broken,%s", rs.Primary.ID)
		}
		environmentId := idSplit[0]
		applicationId := idSplit[1]
		name := idSplit[2]

		service := TemService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
		res, err := service.DescribeTemLogConfig(ctx, environmentId, applicationId, name)
		if err != nil {
			return err
		}

		if res == nil {
			return fmt.Errorf("tem log config %s is not found", rs.Primary.ID)
		}

		return nil
	}
}

const testAccTemLogConfigVar = `
variable "environment_id" {
	default = "` + defaultEnvironmentId + `"
}

variable "application_id" {
	default = "` + defaultApplicationId + `"
}

variable "logset_id" {
	default = "` + defaultLogsetId + `"
}

variable "topic_id" {
	default = "` + defaultTopicId + `"
}

variable "workload_id" {
	default = "` + defaultEnvironmentId + "#" + defaultApplicationId + `"
}
`

const testAccTemLogConfig = testAccTemLogConfigVar + `

resource "tencentcloud_tem_log_config" "logConfig" {
	environment_id = var.environment_id
	application_id = var.application_id
	workload_id = var.workload_id
	name = "terraform-test"
	logset_id = var.logset_id
	topic_id = var.topic_id
	input_type = "container_stdout"
	log_type = "minimalist_log"
}

`
