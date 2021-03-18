package tencentcloud

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccTencentCloudTCRInstance_basic_and_update(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckTCRInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTCRInstance_basic,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("tencentcloud_tcr_instance.mytcr_instance", "name", "testacctcrinstance"),
					resource.TestCheckResourceAttr("tencentcloud_tcr_instance.mytcr_instance", "instance_type", "basic"),
					resource.TestCheckResourceAttr("tencentcloud_tcr_instance.mytcr_instance", "tags.test", "test"),
					resource.TestCheckResourceAttr("tencentcloud_tcr_instance.mytcr_instance", "delete_bucket", "true"),
					resource.TestCheckResourceAttrSet("tencentcloud_tcr_instance.mytcr_instance", "internal_end_point"),
					resource.TestCheckResourceAttrSet("tencentcloud_tcr_instance.mytcr_instance", "status"),
					resource.TestCheckResourceAttrSet("tencentcloud_tcr_instance.mytcr_instance", "public_domain"),
				),
				Destroy: false,
			},
			{
				ResourceName:            "tencentcloud_tcr_instance.mytcr_instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"delete_bucket"},
			},
			{
				Config: testAccTCRInstance_basic_update_remark,
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckTCRInstanceExists("tencentcloud_tcr_instance.mytcr_instance"),
					resource.TestCheckResourceAttr("tencentcloud_tcr_instance.mytcr_instance", "tags.tf", "tf"),
					resource.TestCheckResourceAttr("tencentcloud_tcr_instance.mytcr_instance", "delete_bucket", "true"),
				),
			},
		},
	})
}

func testAccCheckTCRInstanceDestroy(s *terraform.State) error {
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	tcrService := TCRService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_tcr_instance" {
			continue
		}
		_, has, err := tcrService.DescribeTCRInstanceById(ctx, rs.Primary.ID)
		if has {
			return fmt.Errorf("TCR instance still exists")
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func testAccCheckTCRInstanceExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("TCR instance %s is not found", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("TCR instance id is not set")
		}

		tcrService := TCRService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
		_, has, err := tcrService.DescribeTCRInstanceById(ctx, rs.Primary.ID)
		if !has {
			return fmt.Errorf("TCR instance %s is not found", rs.Primary.ID)
		}
		if err != nil {
			return err
		}

		return nil
	}
}

const testAccTCRInstance_basic = `
resource "tencentcloud_tcr_instance" "mytcr_instance" {
  name        = "testacctcrinstance"
  instance_type = "basic"
  delete_bucket = true

  tags ={
	test = "test"
  }
}`

const testAccTCRInstance_basic_update_remark = `
resource "tencentcloud_tcr_instance" "mytcr_instance" {
  name        = "testacctcrinstance"
  instance_type = "basic"
  delete_bucket = true

  tags ={
	tf = "tf"
  }
}`
