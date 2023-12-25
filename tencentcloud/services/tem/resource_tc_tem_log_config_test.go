package tem_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svctem "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/tem"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	sdkErrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
)

// go test -i; go test -test.run TestAccTencentCloudTemLogConfigResource_basic -v
func TestAccTencentCloudTemLogConfigResource_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckTemLogConfigDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTemLogConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTemLogConfigExists("tencentcloud_tem_log_config.logConfig"),
					resource.TestCheckResourceAttrSet("tencentcloud_tem_log_config.logConfig", "id"),
					resource.TestCheckResourceAttr("tencentcloud_tem_log_config.logConfig", "environment_id", tcacctest.DefaultEnvironmentId),
					resource.TestCheckResourceAttr("tencentcloud_tem_log_config.logConfig", "application_id", tcacctest.DefaultApplicationId),
					resource.TestCheckResourceAttr("tencentcloud_tem_log_config.logConfig", "name", "terraform-test"),
					resource.TestCheckResourceAttr("tencentcloud_tem_log_config.logConfig", "logset_id", tcacctest.DefaultLogsetId),
					resource.TestCheckResourceAttr("tencentcloud_tem_log_config.logConfig", "topic_id", tcacctest.DefaultTopicId),
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
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	service := svctem.NewTemService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_tem_log_config" {
			continue
		}

		idSplit := strings.Split(rs.Primary.ID, tccommon.FILED_SP)
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
		logId := tccommon.GetLogId(tccommon.ContextNil)
		ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

		rs, ok := s.RootModule().Resources[r]
		if !ok {
			return fmt.Errorf("resource %s is not found", r)
		}
		idSplit := strings.Split(rs.Primary.ID, tccommon.FILED_SP)
		if len(idSplit) != 3 {
			return fmt.Errorf("id is broken,%s", rs.Primary.ID)
		}
		environmentId := idSplit[0]
		applicationId := idSplit[1]
		name := idSplit[2]

		service := svctem.NewTemService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
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
	default = "` + tcacctest.DefaultEnvironmentId + `"
}

variable "application_id" {
	default = "` + tcacctest.DefaultApplicationId + `"
}

variable "logset_id" {
	default = "` + tcacctest.DefaultLogsetId + `"
}

variable "topic_id" {
	default = "` + tcacctest.DefaultTopicId + `"
}

variable "workload_id" {
	default = "` + tcacctest.DefaultEnvironmentId + "#" + tcacctest.DefaultApplicationId + `"
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
