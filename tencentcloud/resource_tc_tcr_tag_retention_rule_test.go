package tencentcloud

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	sdkErrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
)

func TestAccTencentCloudTcrTagRetentionRuleResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckTCRTagRetentionRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTcrTagRetentionRule,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTCRTagRetentionRuleExists("tencentcloud_tcr_tag_retention_rule.my_rule"),
					resource.TestCheckResourceAttrSet("tencentcloud_tcr_tag_retention_rule.my_rule", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_tcr_tag_retention_rule.my_rule", "registry_id"),
					resource.TestCheckResourceAttr("tencentcloud_tcr_tag_retention_rule.my_rule", "namespace_name", "tf_test_ns_retention"),
					resource.TestCheckResourceAttrSet("tencentcloud_tcr_tag_retention_rule.my_rule", "retention_rule.#"),
					resource.TestCheckResourceAttr("tencentcloud_tcr_tag_retention_rule.my_rule", "retention_rule.0.key", "nDaysSinceLastPush"),
					resource.TestCheckResourceAttr("tencentcloud_tcr_tag_retention_rule.my_rule", "retention_rule.0.value", "1"),
					resource.TestCheckResourceAttr("tencentcloud_tcr_tag_retention_rule.my_rule", "cron_setting", "daily"),
					resource.TestCheckResourceAttr("tencentcloud_tcr_tag_retention_rule.my_rule", "disabled", "false"),
				),
			},
			{
				Config: testAccTcrTagRetentionRule_update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTCRTagRetentionRuleExists("tencentcloud_tcr_tag_retention_rule.my_rule"),
					resource.TestCheckResourceAttrSet("tencentcloud_tcr_tag_retention_rule.my_rule", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_tcr_tag_retention_rule.my_rule", "registry_id"),
					resource.TestCheckResourceAttr("tencentcloud_tcr_tag_retention_rule.my_rule", "namespace_name", "tf_test_ns_retention"),
					resource.TestCheckResourceAttrSet("tencentcloud_tcr_tag_retention_rule.my_rule", "retention_rule.#"),
					resource.TestCheckResourceAttr("tencentcloud_tcr_tag_retention_rule.my_rule", "retention_rule.0.key", "nDaysSinceLastPush"),
					resource.TestCheckResourceAttr("tencentcloud_tcr_tag_retention_rule.my_rule", "retention_rule.0.value", "2"),
					resource.TestCheckResourceAttr("tencentcloud_tcr_tag_retention_rule.my_rule", "cron_setting", "daily"),
					resource.TestCheckResourceAttr("tencentcloud_tcr_tag_retention_rule.my_rule", "disabled", "true"),
				),
			},
			{
				ResourceName:      "tencentcloud_tcr_tag_retention_rule.my_rule",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckTCRTagRetentionRuleDestroy(s *terraform.State) error {
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TCRService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_tcr_tag_retention_rule" {
			continue
		}

		idSplit := strings.Split(rs.Primary.ID, FILED_SP)
		if len(idSplit) != 3 {
			return fmt.Errorf("id is broken,%s", rs.Primary.ID)
		}
		registryId := idSplit[0]
		namespaceName := idSplit[1]
		retentionId := idSplit[2]

		rule, err := service.DescribeTcrTagRetentionRuleById(ctx, registryId, namespaceName, &retentionId)
		if err != nil {
			if ee, ok := err.(*sdkErrors.TencentCloudSDKError); ok {
				if ee.Code == "ResourceNotFound" {
					return nil
				}
			}
			return err
		}

		if rule != nil {
			return fmt.Errorf("Tcr Tag Retention Rule still exist, Id: %v", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckTCRTagRetentionRuleExists(re string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)
		service := TCRService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}

		rs, ok := s.RootModule().Resources[re]
		if !ok {
			return fmt.Errorf("Tcr Tag Retention Rule  %s is not found", re)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("Tcr Tag Retention Rule id is not set")
		}

		idSplit := strings.Split(rs.Primary.ID, FILED_SP)
		if len(idSplit) != 3 {
			return fmt.Errorf("id is broken,%s", rs.Primary.ID)
		}
		registryId := idSplit[0]
		namespaceName := idSplit[1]
		retentionId := idSplit[2]

		rule, err := service.DescribeTcrTagRetentionRuleById(ctx, registryId, namespaceName, &retentionId)
		if err != nil {
			if ee, ok := err.(*sdkErrors.TencentCloudSDKError); ok {
				if ee.Code == "ResourceNotFound" {
					return fmt.Errorf("Tcr Tag Retention Rule not found[ResourceNotFound], Id: %v", rs.Primary.ID)
				}
			}
			return err
		}

		if rule == nil {
			return fmt.Errorf("Tcr Tag Retention Rule not found, Id: %v", rs.Primary.ID)
		}
		return nil
	}
}

const testAccTCRInstance_retention = `
resource "tencentcloud_tcr_instance" "mytcr_retention" {
  name        = "tf-test-tcr-retention"
  instance_type = "basic"
  delete_bucket = true

  tags ={
	test = "test"
  }
}`

const testAccTcrTagRetentionRule = testAccTCRInstance_retention + `

resource "tencentcloud_tcr_namespace" "my_ns" {
  instance_id 	 = tencentcloud_tcr_instance.mytcr_retention.id
  name			 = "tf_test_ns_retention"
  is_public		 = true
  is_auto_scan	 = true
  is_prevent_vul = true
  severity		 = "medium"
  cve_whitelist_items	{
    cve_id = "cve-xxxxx"
  }
}

resource "tencentcloud_tcr_tag_retention_rule" "my_rule" {
  registry_id = tencentcloud_tcr_instance.mytcr_retention.id
  namespace_name = tencentcloud_tcr_namespace.my_ns.name
  retention_rule {
		key = "nDaysSinceLastPush"
		value = 1
  }
  cron_setting = "daily"
  disabled = false
}

`

const testAccTcrTagRetentionRule_update = testAccTCRInstance_retention + `

resource "tencentcloud_tcr_namespace" "my_ns" {
  instance_id 	 = tencentcloud_tcr_instance.mytcr_retention.id
  name			 = "tf_test_ns_retention"
  is_public		 = true
  is_auto_scan	 = true
  is_prevent_vul = true
  severity		 = "medium"
  cve_whitelist_items	{
    cve_id = "cve-xxxxx"
  }
}

resource "tencentcloud_tcr_tag_retention_rule" "my_rule" {
  registry_id = tencentcloud_tcr_instance.mytcr_retention.id
  namespace_name = tencentcloud_tcr_namespace.my_ns.name
  retention_rule {
		key = "nDaysSinceLastPush"
		value = 2
  }
  cron_setting = "daily"
  disabled = true
}

`
