package tencentcloud

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

//go test -i; go test -test.run TestAccTencentCloudResourceTag_basic -v
func TestAccTencentCloudResourceTag_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckResourceTagDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTagResourceTag,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceTagExists("tencentcloud_resource_tag.resource_tag"),
					resource.TestCheckResourceAttr("tencentcloud_resource_tag.resource_tag", "tag_key", "test_terraform_key"),
					resource.TestCheckResourceAttr("tencentcloud_resource_tag.resource_tag", "tag_value", "Terraform_value"),
					resource.TestCheckResourceAttr("tencentcloud_resource_tag.resource_tag", "resource", "qcs::cvm:ap-guangzhou:uin/100020512675:instance/ins-kfrlvcp4")),
			},
			{
				ResourceName:      "tencentcloud_resource_tag.resource_tag",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
func testAccCheckResourceTagDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_resource_tag" {
			continue
		}
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)
		service := TagService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}

		tags, err := service.DescribeResourceTagById(ctx, rs.Primary.Attributes["tag_key"],
			rs.Primary.Attributes["tag_value"], rs.Primary.Attributes["resource"])
		if err != nil {
			return err
		}
		if tags == nil {
			return nil
		}
		return fmt.Errorf("delete resourceTag key %s fail, still on server", rs.Primary.Attributes["tag_key"])
	}
	return nil
}

func testAccCheckResourceTagExists(r string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)

		rs, ok := s.RootModule().Resources[r]
		if !ok {
			return fmt.Errorf("resource %s is not found", r)
		}

		service := TagService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
		res, err := service.DescribeResourceTagById(ctx, rs.Primary.Attributes["tag_key"],
			rs.Primary.Attributes["tag_value"], rs.Primary.Attributes["resource"])
		if err != nil {
			return err
		}
		if res != nil && res.Resource != nil && res.Tags != nil {
			return nil
		}

		return fmt.Errorf("resourceTag %s not found on server", rs.Primary.Attributes["tag_key"])
	}
}

const testAccTagResourceTag = `

resource "tencentcloud_resource_tag" "resource_tag" {
  tag_key = "test_terraform_key"
  tag_value = "Terraform_value"
  resource = "qcs::cvm:ap-guangzhou:uin/100020512675:instance/ins-kfrlvcp4"
}

`
