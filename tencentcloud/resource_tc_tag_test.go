package tencentcloud

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

//go test -i; go test -test.run TestAccTencentCloudTagResource_basic -v
func TestAccTencentCloudTagResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckTagDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTag,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTagExists("tencentcloud_tag.tag"),
					resource.TestCheckResourceAttr("tencentcloud_tag.tag", "tag_key", "test_terraform_tag_key"),
					resource.TestCheckResourceAttr("tencentcloud_tag.tag", "tag_value", "Terraform_tag_value")),
			},
			{
				ResourceName:      "tencentcloud_tag.tag",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckTagDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_tag" {
			continue
		}
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)
		service := TagService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}

		tags, err := service.DescribeTagResourceById(ctx, rs.Primary.Attributes["tag_key"], rs.Primary.Attributes["tag_value"])
		if err != nil {
			return err
		}
		if tags == nil {
			return nil
		}
		return fmt.Errorf("delete tag key %s fail, still on server", rs.Primary.Attributes["tag_key"])
	}
	return nil
}

func testAccCheckTagExists(r string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)

		rs, ok := s.RootModule().Resources[r]
		if !ok {
			return fmt.Errorf("resource %s is not found", r)
		}

		service := TagService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
		tags, err := service.DescribeTagResourceById(ctx, rs.Primary.Attributes["tag_key"], rs.Primary.Attributes["tag_value"])
		if err != nil {
			return err
		}
		if tags != nil && tags.TagKey != nil && tags.TagValue != nil {
			return nil
		}

		return fmt.Errorf("tag %s not found on server", rs.Primary.Attributes["tag_key"])
	}
}

const testAccTag = `

resource "tencentcloud_tag" "tag" {
  tag_key = "test_terraform_tag_key"
  tag_value = "Terraform_tag_value"
}

`
