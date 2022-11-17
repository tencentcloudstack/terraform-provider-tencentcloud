package tencentcloud

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func init() {
	resource.AddTestSweepers("tencentcloud_tdcpg_instance", &resource.Sweeper{
		Name: "tencentcloud_tdcpg_instance",
		F:    testSweepTdcpgInstance,
	})
}

// go test -v ./tencentcloud -sweep=ap-guangzhou -sweep-run=tencentcloud_tdcpg_instance
func testSweepTdcpgInstance(r string) error {
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	cli, _ := sharedClientForRegion(r)
	tdcpgService := TdcpgService{client: cli.(*TencentCloudClient).apiV3Conn}

	instances, err := tdcpgService.DescribeTdcpgInstancesByFilter(ctx, helper.String(defaultTdcpgClusterId), nil)
	if err != nil {
		return err
	}
	if instances == nil {
		return fmt.Errorf("tdcpg instances not exists. clusterId:[%s]", defaultTdcpgClusterId)
	}

	// delete all instances which has specified prefix under the default cluster
	for _, v := range instances {
		delId := v.InstanceId
		delName := v.InstanceName
		status := *v.Status

		if status == "running" && strings.HasPrefix(*delName, defaultTdcpgTestNamePrefix) {
			err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
				err := tdcpgService.DeleteTdcpgInstanceById(ctx, helper.String(defaultTdcpgClusterId), delId)
				if err != nil {
					return retryError(err)
				}
				return nil
			})
			if err != nil {
				return fmt.Errorf("[ERROR] delete tdcpg instance %s failed. reason:[%s]", *delId, err.Error())
			}
		}
	}
	return nil
}

func TestAccTencentCloudTdcpgInstanceResource_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckTdcpgInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccTdcpgInstance_basic, defaultTdcpgClusterId, defaultTdcpgTestNamePrefix),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTdcpgInstanceExists("tencentcloud_tdcpg_instance.instance"),
					resource.TestCheckResourceAttr("tencentcloud_tdcpg_instance.instance", "cluster_id", defaultTdcpgClusterId),
					resource.TestCheckResourceAttrSet("tencentcloud_tdcpg_instance.instance", "cpu"),
					resource.TestCheckResourceAttrSet("tencentcloud_tdcpg_instance.instance", "memory"),
					resource.TestMatchResourceAttr("tencentcloud_tdcpg_instance.instance", "instance_name", regexp.MustCompile(defaultTdcpgTestNamePrefix)),
				),
			},
			{
				Config: fmt.Sprintf(testAccTdcpgInstance_update, defaultTdcpgClusterId, defaultTdcpgTestNamePrefix),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTdcpgInstanceExists("tencentcloud_tdcpg_instance.instance"),
					resource.TestCheckResourceAttr("tencentcloud_tdcpg_instance.instance", "cluster_id", defaultTdcpgClusterId),
					resource.TestCheckResourceAttr("tencentcloud_tdcpg_instance.instance", "cpu", "2"),
					resource.TestCheckResourceAttr("tencentcloud_tdcpg_instance.instance", "memory", "4"),
					resource.TestMatchResourceAttr("tencentcloud_tdcpg_instance.instance", "instance_name", regexp.MustCompile(defaultTdcpgTestNamePrefix)),
					resource.TestCheckResourceAttr("tencentcloud_tdcpg_instance.instance", "operation_timing", "IMMEDIATE"),
				),
			},
			{
				ResourceName:            "tencentcloud_tdcpg_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"operation_timing"},
			},
		},
	})
}

func testAccCheckTdcpgInstanceDestroy(s *terraform.State) error {
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	tdcpgService := TdcpgService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_tdcpg_instance" {
			continue
		}
		ids := strings.Split(rs.Primary.ID, FILED_SP)

		ret, err := tdcpgService.DescribeTdcpgInstance(ctx, &ids[0], &ids[1])
		if err != nil {
			return err
		}

		if ret != nil && len(ret.InstanceSet) > 0 {
			status := *ret.InstanceSet[0].Status
			if status == "deleting" || status == "deleted" || status == "isolated" || status == "isolating" {
				return nil
			}
			return fmt.Errorf("tdcpg instance still exist, id: %v, status: %v", rs.Primary.ID, status)
		}
	}
	return nil
}

func testAccCheckTdcpgInstanceExists(re string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)

		rs, ok := s.RootModule().Resources[re]
		if !ok {
			return fmt.Errorf("tdcpg instance  %s is not found", re)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("tdcpg instance id is not set")
		}

		tdcpgService := TdcpgService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
		ids := strings.Split(rs.Primary.ID, FILED_SP)
		ret, err := tdcpgService.DescribeTdcpgInstance(ctx, &ids[0], &ids[1])
		if err != nil {
			return err
		}

		if ret == nil || len(ret.InstanceSet) == 0 {
			return fmt.Errorf("tdcpg instance not found, id: %v", rs.Primary.ID)
		}

		return nil
	}
}

const testAccTdcpgInstance_basic = `

resource "tencentcloud_tdcpg_instance" "instance" {
  cluster_id = "%s"
  cpu = 1
  memory = 1
  instance_name = "%sinstance"
}

`

const testAccTdcpgInstance_update = `

resource "tencentcloud_tdcpg_instance" "instance" {
  cluster_id = "%s"
  cpu = 2
  memory = 4
  instance_name = "%sinstance"
  operation_timing = "IMMEDIATE"
}

`
