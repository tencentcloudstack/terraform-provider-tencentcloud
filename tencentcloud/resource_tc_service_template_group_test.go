package tencentcloud

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccTencentCloudServiceTemplateGroup_basic_and_update(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckServiceTemplateGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccServiceTemplateGroup_basic,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("tencentcloud_service_template_group.group", "name", "test"),
					resource.TestCheckResourceAttr("tencentcloud_service_template_group.group", "template_ids.#", "1"),
				),
			},
			{
				ResourceName:      "tencentcloud_service_template_group.group",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccServiceTemplateGroup_basic_update_remark,
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckServiceTemplateGroupExists("tencentcloud_service_template_group.group"),
					resource.TestCheckResourceAttr("tencentcloud_service_template_group.group", "name", "test_update"),
					resource.TestCheckResourceAttr("tencentcloud_service_template_group.group", "template_ids.#", "1"),
				),
			},
		},
	})
}

func testAccCheckServiceTemplateGroupDestroy(s *terraform.State) error {
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	vpcService := VpcService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_service_template_group" {
			continue
		}

		_, has, err := vpcService.DescribeServiceTemplateGroupById(ctx, rs.Primary.ID)
		if has {
			return fmt.Errorf("Service template group still exists")
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func testAccCheckServiceTemplateGroupExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Service template group %s is not found", n)
		}

		vpcService := VpcService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
		_, has, err := vpcService.DescribeServiceTemplateGroupById(ctx, rs.Primary.ID)
		if !has {
			return fmt.Errorf("Service template group %s is not found", rs.Primary.ID)
		}
		if err != nil {
			return err
		}

		return nil
	}
}

const testAccServiceTemplateGroup_basic = `
resource "tencentcloud_service_template" "template" {
  name = "test"
  services = ["tcp:all"]
}

resource "tencentcloud_service_template_group" "group"{
	name = "test"
	template_ids = [tencentcloud_service_template.template.id]
}
`

const testAccServiceTemplateGroup_basic_update_remark = `
resource "tencentcloud_service_template" "template" {
  name = "test"
  services = ["tcp:all"]
}

resource "tencentcloud_service_template" "templateB" {
  name = "testB"
  services = ["tcp:80", "udp:90,111"]
}

resource "tencentcloud_service_template_group" "group"{
	name = "test_update"
	template_ids = [tencentcloud_service_template.templateB.id]
}
`
