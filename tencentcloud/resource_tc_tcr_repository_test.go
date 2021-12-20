package tencentcloud

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccTencentCloudTCRRepository_basic_and_update(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckTCRRepositoryDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTCRRepository_basic,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("tencentcloud_tcr_repository.mytcr_repository", "name", "test"),
					resource.TestCheckResourceAttr("tencentcloud_tcr_repository.mytcr_repository", "brief_desc", "111"),
					resource.TestCheckResourceAttr("tencentcloud_tcr_repository.mytcr_repository", "description", "111111111111111111111111111111111111"),
					resource.TestCheckResourceAttrSet("tencentcloud_tcr_repository.mytcr_repository", "create_time"),
					resource.TestCheckResourceAttrSet("tencentcloud_tcr_repository.mytcr_repository", "update_time"),
					resource.TestCheckResourceAttrSet("tencentcloud_tcr_repository.mytcr_repository", "is_public"),
				),
				Destroy: false,
			},
			{
				ResourceName:      "tencentcloud_tcr_repository.mytcr_repository",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccTCRRepository_basic_update_remark,
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckTCRRepositoryExists("tencentcloud_tcr_repository.mytcr_repository"),
					resource.TestCheckResourceAttr("tencentcloud_tcr_repository.mytcr_repository", "brief_desc", "2222"),
					resource.TestCheckResourceAttr("tencentcloud_tcr_repository.mytcr_repository", "description", "211111111111111111111111111111111111"),
				),
			},
		},
	})
}

func testAccCheckTCRRepositoryDestroy(s *terraform.State) error {
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	tcrService := TCRService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_tcr_repository" {
			continue
		}
		items := strings.Split(rs.Primary.ID, FILED_SP)
		if len(items) != 3 {
			return fmt.Errorf("invalid ID %s", rs.Primary.ID)
		}

		instanceId := items[0]
		namespaceName := items[1]
		repositoryName := items[2]
		_, has, err := tcrService.DescribeTCRRepositoryById(ctx, instanceId, namespaceName, repositoryName)
		if has {
			return fmt.Errorf("TCR repository still exists")
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func testAccCheckTCRRepositoryExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("TCR repository %s is not found", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("TCR repository id is not set")
		}
		items := strings.Split(rs.Primary.ID, FILED_SP)
		if len(items) != 3 {
			return fmt.Errorf("invalid ID %s", rs.Primary.ID)
		}

		instanceId := items[0]
		namespaceName := items[1]
		repositoryName := items[2]
		tcrService := TCRService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
		_, has, err := tcrService.DescribeTCRRepositoryById(ctx, instanceId, namespaceName, repositoryName)
		if !has {
			return fmt.Errorf("TCR repository %s is not found", rs.Primary.ID)
		}
		if err != nil {
			return err
		}

		return nil
	}
}

const testAccTCRRepository_basic = `
resource "tencentcloud_tcr_instance" "mytcr_instance" {
  name        = "testacctcrinstance"
  instance_type = "basic"
  delete_bucket = true

  tags ={
	test = "test"
  }
}

resource "tencentcloud_tcr_namespace" "mytcr_namespace" {
  instance_id = tencentcloud_tcr_instance.mytcr_instance.id
  name        = "test"
  is_public   = false
}
resource "tencentcloud_tcr_repository" "mytcr_repository" {
  instance_id = tencentcloud_tcr_instance.mytcr_instance.id
  namespace_name        = tencentcloud_tcr_namespace.mytcr_namespace.name
  name = "test"
  brief_desc = "111"
  description = "111111111111111111111111111111111111"
}`

const testAccTCRRepository_basic_update_remark = `
resource "tencentcloud_tcr_instance" "mytcr_instance" {
  name        = "testacctcrinstance"
  instance_type = "basic"
  delete_bucket = true

  tags ={
	test = "test"
  }
}

resource "tencentcloud_tcr_namespace" "mytcr_namespace" {
  instance_id = tencentcloud_tcr_instance.mytcr_instance.id
  name        = "test"
  is_public   = false
}
resource "tencentcloud_tcr_repository" "mytcr_repository" {
  instance_id = tencentcloud_tcr_instance.mytcr_instance.id
  namespace_name        = tencentcloud_tcr_namespace.mytcr_namespace.name
  name = "test"
  brief_desc = "2222"
  description = "211111111111111111111111111111111111"
}`
