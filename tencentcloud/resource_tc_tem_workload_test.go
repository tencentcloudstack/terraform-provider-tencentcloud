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

// go test -i; go test -test.run TestAccTencentCloudTemWorkloadResource_basic -v
func TestAccTencentCloudTemWorkloadResource_basic(t *testing.T) {
	// t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckTemWorkloadDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTemWorkload,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTemWorkloadExists("tencentcloud_tem_workload.workload"),
					resource.TestCheckResourceAttrSet("tencentcloud_tem_workload.workload", "id"),
					resource.TestCheckResourceAttr("tencentcloud_tem_workload.workload", "environment_id", defaultEnvironmentId),
					resource.TestCheckResourceAttrSet("tencentcloud_tem_workload.workload", "application_id"),
					resource.TestCheckResourceAttr("tencentcloud_tem_workload.workload", "deploy_version", "hello-world"),
					resource.TestCheckResourceAttr("tencentcloud_tem_workload.workload", "deploy_mode", "IMAGE"),
					resource.TestCheckResourceAttr("tencentcloud_tem_workload.workload", "img_repo", "tem_demo/tem_demo"),
					resource.TestCheckResourceAttr("tencentcloud_tem_workload.workload", "repo_server", "ccr.ccs.tencentyun.com"),
					resource.TestCheckResourceAttr("tencentcloud_tem_workload.workload", "init_pod_num", "1"),
					resource.TestCheckResourceAttr("tencentcloud_tem_workload.workload", "cpu_spec", "1"),
					resource.TestCheckResourceAttr("tencentcloud_tem_workload.workload", "memory_spec", "1"),
				),
			},
			{
				ResourceName:      "tencentcloud_tem_workload.workload",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckTemWorkloadDestroy(s *terraform.State) error {
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	service := TemService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_tem_workload" {
			continue
		}

		idSplit := strings.Split(rs.Primary.ID, FILED_SP)
		if len(idSplit) != 2 {
			return fmt.Errorf("id is broken,%s", rs.Primary.ID)
		}
		environmentId := idSplit[0]
		applicationId := idSplit[1]

		res, err := service.DescribeTemWorkload(ctx, environmentId, applicationId)
		if err != nil {
			ee, ok := err.(*sdkErrors.TencentCloudSDKError)
			if !ok {
				return err
			}
			if ee.Code == "ResourceNotFound.ServiceRunningVersionNotFound" || ee.Code == "ResourceNotFound.ServiceNotFound" {
				return nil
			}
			return err
		}

		if res != nil {
			return fmt.Errorf("tem workload %s still exists", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckTemWorkloadExists(r string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)

		rs, ok := s.RootModule().Resources[r]
		if !ok {
			return fmt.Errorf("resource %s is not found", r)
		}
		idSplit := strings.Split(rs.Primary.ID, FILED_SP)
		if len(idSplit) != 2 {
			return fmt.Errorf("id is broken,%s", rs.Primary.ID)
		}
		environmentId := idSplit[0]
		applicationId := idSplit[1]

		service := TemService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
		res, err := service.DescribeTemWorkload(ctx, environmentId, applicationId)
		if err != nil {
			return err
		}

		if res == nil {
			return fmt.Errorf("tem workload %s is not found", rs.Primary.ID)
		}

		return nil
	}
}

const testAccTemWorkloadVar = `
variable "environment_id" {
	default = "` + defaultEnvironmentId + `"
}
`

const testAccTemWorkload = testAccTemApplication + testAccTemWorkloadVar + `

resource "tencentcloud_tem_workload" "workload" {
  application_id     = tencentcloud_tem_application.application.id
  environment_id     = var.environment_id
  deploy_version     = "hello-world"
  deploy_mode        = "IMAGE"
  img_repo           = "tem_demo/tem_demo"
  repo_server        = "ccr.ccs.tencentyun.com"
  init_pod_num       = 1
  cpu_spec           = 1
  memory_spec        = 1
}

`
