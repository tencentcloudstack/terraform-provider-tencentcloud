package tencentcloud

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccTencentCloudTCRToken_basic_and_update(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckTCRTokenDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTCRToken_basic,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("tencentcloud_tcr_token.mytcr_token", "description", "test"),
					resource.TestCheckResourceAttr("tencentcloud_tcr_token.mytcr_token", "enable", "true"),
					resource.TestCheckResourceAttrSet("tencentcloud_tcr_token.mytcr_token", "token_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_tcr_token.mytcr_token", "create_time"),
					resource.TestCheckResourceAttrSet("tencentcloud_tcr_token.mytcr_token", "token"),
					resource.TestCheckResourceAttrSet("tencentcloud_tcr_token.mytcr_token", "user_name"),
				),
				Destroy: false,
			},
			{
				ResourceName:            "tencentcloud_tcr_token.mytcr_token",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"token", "user_name"},
			},
			{
				Config: testAccTCRToken_basic_update_remark,
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckTCRTokenExists("tencentcloud_tcr_token.mytcr_token"),
					resource.TestCheckResourceAttr("tencentcloud_tcr_token.mytcr_token", "enable", "false"),
				),
			},
		},
	})
}

func testAccCheckTCRTokenDestroy(s *terraform.State) error {
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	tcrService := TCRService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_tcr_token" {
			continue
		}
		items := strings.Split(rs.Primary.ID, FILED_SP)
		if len(items) != 2 {
			return fmt.Errorf("invalid ID %s", rs.Primary.ID)
		}

		instanceId := items[0]
		tokenId := items[1]
		_, has, err := tcrService.DescribeTCRLongTermTokenById(ctx, instanceId, tokenId)
		if has {
			return fmt.Errorf("TCR token still exists")
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func testAccCheckTCRTokenExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("TCR token %s is not found", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("TCR token id is not set")
		}
		items := strings.Split(rs.Primary.ID, FILED_SP)
		if len(items) != 2 {
			return fmt.Errorf("invalid ID %s", rs.Primary.ID)
		}

		instanceId := items[0]
		tokenId := items[1]
		tcrService := TCRService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
		_, has, err := tcrService.DescribeTCRLongTermTokenById(ctx, instanceId, tokenId)
		if !has {
			return fmt.Errorf("TCR token %s is not found", rs.Primary.ID)
		}
		if err != nil {
			return err
		}

		return nil
	}
}

const testAccTCRToken_basic = `
resource "tencentcloud_tcr_instance" "mytcr_instance" {
  name        = "testacctcrinstance"
  instance_type = "basic"
  delete_bucket = true

  tags ={
	test = "test"
  }
}

resource "tencentcloud_tcr_token" "mytcr_token" {
  instance_id = tencentcloud_tcr_instance.mytcr_instance.id
  description       = "test"
}`

const testAccTCRToken_basic_update_remark = `
resource "tencentcloud_tcr_instance" "mytcr_instance" {
  name        = "testacctcrinstance"
  instance_type = "basic"
  delete_bucket = true

  tags ={
	test = "test"
  }
}

resource "tencentcloud_tcr_token" "mytcr_token" {
  instance_id = tencentcloud_tcr_instance.mytcr_instance.id
  description       = "test"
  enable   = false
}`
