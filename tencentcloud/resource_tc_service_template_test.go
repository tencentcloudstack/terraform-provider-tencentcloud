package tencentcloud

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccTencentCloudServiceTemplate_basic_and_update(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckServiceTemplateDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccServiceTemplate_basic,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("tencentcloud_service_template.template", "name", "test"),
					resource.TestCheckResourceAttr("tencentcloud_service_template.template", "services.#", "1"),
				),
			},
			{
				ResourceName:      "tencentcloud_service_template.template",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccServiceTemplate_basic_update_remark,
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckServiceTemplateExists("tencentcloud_service_template.template"),
					resource.TestCheckResourceAttr("tencentcloud_service_template.template", "name", "test_update"),
					resource.TestCheckResourceAttr("tencentcloud_service_template.template", "services.#", "2"),
				),
			},
		},
	})
}

func testAccCheckServiceTemplateDestroy(s *terraform.State) error {
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	vpcService := VpcService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_service_template" {
			continue
		}

		_, has, err := vpcService.DescribeServiceTemplateById(ctx, rs.Primary.ID)
		if has {
			return fmt.Errorf("service template still exists")
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func testAccCheckServiceTemplateExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Service template %s is not found", n)
		}

		vpcService := VpcService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
		_, has, err := vpcService.DescribeServiceTemplateById(ctx, rs.Primary.ID)
		if !has {
			return fmt.Errorf("Service template %s is not found", rs.Primary.ID)
		}
		if err != nil {
			return err
		}

		return nil
	}
}

const testAccServiceTemplate_basic = `
resource "tencentcloud_service_template" "template" {
  name = "test"
  services = ["tcp:80"]
}`

const testAccServiceTemplate_basic_update_remark = `
resource "tencentcloud_service_template" "template" {
  name = "test_update"
  services = ["udp:all", "tcp:80,90"]
}`
