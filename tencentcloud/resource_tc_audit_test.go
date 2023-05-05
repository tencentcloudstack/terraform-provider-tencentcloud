package tencentcloud

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccTencentCloudNeedFixAudit_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAuditDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAudit_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAuditExists("tencentcloud_audit.audit_basic"),
					resource.TestCheckResourceAttr("tencentcloud_audit.audit_basic", "name", "audittest"),
					resource.TestCheckResourceAttr("tencentcloud_audit.audit_basic", "read_write_attribute", "3"),
					resource.TestCheckResourceAttr("tencentcloud_audit.audit_basic", "cos_bucket", "test"),
					resource.TestCheckResourceAttr("tencentcloud_audit.audit_basic", "cos_region", "ap-hongkong"),
					resource.TestCheckResourceAttr("tencentcloud_audit.audit_basic", "log_file_prefix", "test"),
					resource.TestCheckResourceAttr("tencentcloud_audit.audit_basic", "audit_switch", "true"),
				),
			},
			{
				ResourceName:      "tencentcloud_audit.audit_basic",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccAudit_update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAuditExists("tencentcloud_audit.audit_basic"),
					resource.TestCheckResourceAttr("tencentcloud_audit.audit_basic", "name", "audittest"),
					resource.TestCheckResourceAttr("tencentcloud_audit.audit_basic", "read_write_attribute", "2"),
					resource.TestCheckResourceAttr("tencentcloud_audit.audit_basic", "cos_bucket", "test1"),
					resource.TestCheckResourceAttr("tencentcloud_audit.audit_basic", "cos_region", "ap-shanghai"),
					resource.TestCheckResourceAttr("tencentcloud_audit.audit_basic", "log_file_prefix", "test11"),
					resource.TestCheckResourceAttr("tencentcloud_audit.audit_basic", "audit_switch", "false"),
				),
			},
		},
	})
}

func TestAccTencentCloudNeedFixAudit_kms(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAuditDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAudit_kms,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAuditExists("tencentcloud_audit.audit_kms"),
					resource.TestCheckResourceAttr("tencentcloud_audit.audit_kms", "name", "audittest"),
					resource.TestCheckResourceAttr("tencentcloud_audit.audit_kms", "enable_kms_encry", "true"),
					resource.TestCheckResourceAttrSet("tencentcloud_audit.audit_kms", "key_id"),
					resource.TestCheckResourceAttr("tencentcloud_audit.audit_kms", "read_write_attribute", "3"),
					resource.TestCheckResourceAttr("tencentcloud_audit.audit_kms", "cos_bucket", "test"),
					resource.TestCheckResourceAttr("tencentcloud_audit.audit_kms", "cos_region", "ap-hongkong"),
					resource.TestCheckResourceAttr("tencentcloud_audit.audit_kms", "log_file_prefix", "test"),
					resource.TestCheckResourceAttr("tencentcloud_audit.audit_kms", "audit_switch", "true"),
				),
			},
			{
				ResourceName:      "tencentcloud_audit.audit_kms",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccAudit_kms_update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAuditExists("tencentcloud_audit.audit_kms"),
					resource.TestCheckResourceAttr("tencentcloud_audit.audit_kms", "name", "audittest"),
					resource.TestCheckResourceAttr("tencentcloud_audit.audit_kms", "enable_kms_encry", "false"),
					resource.TestCheckResourceAttr("tencentcloud_audit.audit_kms", "read_write_attribute", "2"),
					resource.TestCheckResourceAttr("tencentcloud_audit.audit_kms", "cos_bucket", "test1"),
					resource.TestCheckResourceAttr("tencentcloud_audit.audit_kms", "cos_region", "ap-shanghai"),
					resource.TestCheckResourceAttr("tencentcloud_audit.audit_kms", "log_file_prefix", "test11"),
					resource.TestCheckResourceAttr("tencentcloud_audit.audit_kms", "audit_switch", "false"),
				),
			},
		},
	})
}

func testAccCheckAuditDestroy(s *terraform.State) error {
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	auditService := AuditService{
		client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn,
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_audit" {
			continue
		}
		time.Sleep(5 * time.Second)
		_, has, err := auditService.DescribeAuditById(ctx, rs.Primary.ID)
		if err != nil {
			return err
		}
		if has {
			return fmt.Errorf("[CHECK][Audit][Exists] id %s still exist", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckAuditExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("[CHECK][Audit][Exists] check: Audit %s is not found", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("[CHECK][Audit][Create] check: Audit id is not set")
		}
		auditService := AuditService{
			client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn,
		}
		_, has, err := auditService.DescribeAuditById(ctx, rs.Primary.ID)
		if err != nil {
			return err
		}
		if !has {
			return fmt.Errorf("[CHECK][Audit][Exists] id %s is not exist", rs.Primary.ID)
		}
		return nil
	}
}

const testAccAudit_basic = `
resource "tencentcloud_audit" "audit_basic" {
  name        = "audittest"
  cos_bucket	= "test"
  cos_region = "ap-hongkong"
  log_file_prefix = "test"
  audit_switch = true
  read_write_attribute = 3
}
`

const testAccAudit_update = `
resource "tencentcloud_audit" "audit_basic" {
  name        = "audittest"
  cos_bucket	= "test1"
  cos_region = "ap-shanghai"
  log_file_prefix = "test11"
  audit_switch = false
  read_write_attribute = 2
}
`

const testAccAudit_kms = `
data "tencentcloud_audit_key_alias" "all" {
	region = "ap-hongkong"
}

resource "tencentcloud_audit" "audit_kms" {
  name        = "audittest"
  cos_bucket	= "test"
  cos_region = "ap-hongkong"
  enable_kms_encry = true
  log_file_prefix = "test"
  key_id = data.tencentcloud_audit_key_alias.all.audit_key_alias_list.0.key_id
  audit_switch = true
  read_write_attribute = 3
}
`

const testAccAudit_kms_update = `
resource "tencentcloud_audit" "audit_kms" {
  name        = "audittest"
  cos_bucket	= "test1"
  cos_region = "ap-shanghai"
  enable_kms_encry = false
  log_file_prefix = "test11"
  audit_switch = false
  read_write_attribute = 2
}
`
