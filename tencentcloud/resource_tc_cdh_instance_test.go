package tencentcloud

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccTencentCloudCdhInstance_basic(t *testing.T) {
	t.Parallel()
	resourceName := "tencentcloud_cdh_instance.foo"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCdhInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: TestAccTencentCloudCdhInstance_basicConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCdhInstanceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "availability_zone", "ap-guangzhou-3"),
					resource.TestCheckResourceAttr(resourceName, "project_id", "0"),
					resource.TestCheckResourceAttr(resourceName, "host_type", "HC20"),
					resource.TestCheckResourceAttr(resourceName, "host_name", "unit-test"),
					resource.TestCheckResourceAttr(resourceName, "charge_type", "PREPAID"),
					resource.TestCheckResourceAttr(resourceName, "prepaid_renew_flag", "NOTIFY_AND_MANUAL_RENEW"),
					resource.TestCheckResourceAttr(resourceName, "host_state", "RUNNING"),
					resource.TestCheckResourceAttr(resourceName, "host_resource.#", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "create_time"),
					resource.TestCheckResourceAttrSet(resourceName, "expired_time"),
				),
			},
			{
				Config: TestAccTencentCloudCdhInstance_modify,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCdhInstanceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "host_name", "unit-test-modify"),
					resource.TestCheckResourceAttr(resourceName, "prepaid_renew_flag", "DISABLE_NOTIFY_AND_MANUAL_RENEW"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"prepaid_period"},
			},
		},
	})
}

func testAccCheckCdhInstanceDestroy(s *terraform.State) error {
	return nil
}

func testAccCheckCdhInstanceExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)

		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("[CHECK][CDH instance][Exists] check: CDH instance %s is not found", name)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("[CHECK][CDH instance][Exists] check:CDH instance id is not set")
		}
		cdhService := CdhService{
			client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn,
		}
		instance, err := cdhService.DescribeCdhInstanceById(ctx, rs.Primary.ID)
		if err != nil {
			return err
		}
		if instance == nil {
			return fmt.Errorf("[CHECK][CDH instance][Exists] id %s is not exist", rs.Primary.ID)
		}
		return nil
	}
}

const TestAccTencentCloudCdhInstance_basicConfig = `
variable "availability_zone" {
  default = "ap-guangzhou-3"
}

resource "tencentcloud_cdh_instance" "foo" {
  availability_zone = var.availability_zone
  host_type = "HC20"
  charge_type = "PREPAID"
  prepaid_period = 1
  host_name = "unit-test"
  prepaid_renew_flag = "NOTIFY_AND_MANUAL_RENEW"
}
`

const TestAccTencentCloudCdhInstance_modify = `
variable "availability_zone" {
  default = "ap-guangzhou-3"
}

resource "tencentcloud_cdh_instance" "foo" {
  availability_zone = var.availability_zone
  host_type = "HC20"
  charge_type = "PREPAID"
  prepaid_period = 1
  host_name = "unit-test-modify"
  project_id = ` + defaultProjectId + `
  prepaid_renew_flag = "DISABLE_NOTIFY_AND_MANUAL_RENEW"
}
`
