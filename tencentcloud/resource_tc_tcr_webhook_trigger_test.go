package tencentcloud

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	sdkErrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func TestAccTencentCloudTcrWebhookTriggerResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckTCRWebhookTriggerDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccTCRWebhookTrigger, "webhooktrigger", "webhooktrigger", "webhooktrigger"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTCRWebhookTriggerExists("tencentcloud_tcr_webhook_trigger.my_trigger"),
					resource.TestCheckResourceAttrSet("tencentcloud_tcr_webhook_trigger.my_trigger", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_tcr_webhook_trigger.my_trigger", "registry_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_tcr_webhook_trigger.my_trigger", "namespace"),
					resource.TestCheckResourceAttr("tencentcloud_tcr_webhook_trigger.my_trigger", "trigger.0.name", "trigger-webhooktrigger"),
					resource.TestCheckResourceAttr("tencentcloud_tcr_webhook_trigger.my_trigger", "trigger.0.event_types.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_tcr_webhook_trigger.my_trigger", "trigger.0.condition", ".*"),
					resource.TestCheckResourceAttr("tencentcloud_tcr_webhook_trigger.my_trigger", "trigger.0.enabled", "true"),
					resource.TestCheckResourceAttr("tencentcloud_tcr_webhook_trigger.my_trigger", "trigger.0.targets.0.address", "http://example.org/post"),
					resource.TestCheckResourceAttr("tencentcloud_tcr_webhook_trigger.my_trigger", "trigger.0.targets.0.headers.0.key", "X-Custom-Header"),
					resource.TestCheckResourceAttr("tencentcloud_tcr_webhook_trigger.my_trigger", "trigger.0.targets.0.headers.0.values.#", "1"),
					resource.TestCheckResourceAttrSet("tencentcloud_tcr_webhook_trigger.my_trigger", "trigger.0.description"),
					// resource.TestCheckResourceAttrSet("tencentcloud_tcr_webhook_trigger.my_trigger", "trigger.0.namespace_id"),
				),
			},
			{
				Config: fmt.Sprintf(testAccTCRWebhookTrigger_update, "webhooktrigger", "webhooktrigger", "webhooktrigger"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTCRWebhookTriggerExists("tencentcloud_tcr_webhook_trigger.my_trigger"),
					resource.TestCheckResourceAttrSet("tencentcloud_tcr_webhook_trigger.my_trigger", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_tcr_webhook_trigger.my_trigger", "registry_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_tcr_webhook_trigger.my_trigger", "namespace"),
					resource.TestCheckResourceAttr("tencentcloud_tcr_webhook_trigger.my_trigger", "trigger.0.name", "trigger-webhooktrigger"),
					resource.TestCheckResourceAttr("tencentcloud_tcr_webhook_trigger.my_trigger", "trigger.0.event_types.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_tcr_webhook_trigger.my_trigger", "trigger.0.condition", ".*test"),
					resource.TestCheckResourceAttr("tencentcloud_tcr_webhook_trigger.my_trigger", "trigger.0.enabled", "false"),
					resource.TestCheckResourceAttr("tencentcloud_tcr_webhook_trigger.my_trigger", "trigger.0.targets.0.address", "http://example.org/post"),
					resource.TestCheckResourceAttr("tencentcloud_tcr_webhook_trigger.my_trigger", "trigger.0.targets.0.headers.0.key", "X-Custom-Header"),
					resource.TestCheckResourceAttr("tencentcloud_tcr_webhook_trigger.my_trigger", "trigger.0.targets.0.headers.0.values.#", "1"),
					resource.TestCheckResourceAttrSet("tencentcloud_tcr_webhook_trigger.my_trigger", "trigger.0.description"),
					// resource.TestCheckResourceAttrSet("tencentcloud_tcr_webhook_trigger.my_trigger", "trigger.0.namespace_id"),
				),
			},
			{
				ResourceName:      "tencentcloud_tcr_webhook_trigger.my_trigger",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckTCRWebhookTriggerDestroy(s *terraform.State) error {
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TCRService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_tcr_webhook_trigger" {
			continue
		}

		idSplit := strings.Split(rs.Primary.ID, FILED_SP)
		if len(idSplit) != 3 {
			return fmt.Errorf("id is broken,%s", rs.Primary.ID)
		}
		registryId := idSplit[0]
		namespaceName := idSplit[1]
		triggerId := helper.StrToInt64(idSplit[2])

		trigger, err := service.DescribeTcrWebhookTriggerById(ctx, registryId, triggerId, namespaceName)
		if err != nil {
			if ee, ok := err.(*sdkErrors.TencentCloudSDKError); ok {
				if ee.Code == "ResourceNotFound" {
					return nil
				}
			}
			return err
		}

		if trigger != nil {
			return fmt.Errorf("Tcr web hook trigger still exist, Id: %v", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckTCRWebhookTriggerExists(re string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)
		service := TCRService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}

		rs, ok := s.RootModule().Resources[re]
		if !ok {
			return fmt.Errorf("Tcr web hook trigger %s is not found", re)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("Tcr web hook trigger id is not set")
		}

		idSplit := strings.Split(rs.Primary.ID, FILED_SP)
		if len(idSplit) != 3 {
			return fmt.Errorf("id is broken,%s", rs.Primary.ID)
		}
		registryId := idSplit[0]
		namespaceName := idSplit[1]
		triggerId := helper.StrToInt64(idSplit[2])

		trigger, err := service.DescribeTcrWebhookTriggerById(ctx, registryId, triggerId, namespaceName)
		if err != nil {
			if ee, ok := err.(*sdkErrors.TencentCloudSDKError); ok {
				if ee.Code == "ResourceNotFound" {
					return fmt.Errorf("Tcr web hook trigger not found[ResourceNotFound], Id: %v", rs.Primary.ID)
				}
			}
			return err
		}

		if trigger == nil {
			return fmt.Errorf("Tcr web hook trigger not found, Id: %v", rs.Primary.ID)
		}
		return nil
	}
}

const testAccTCRInstance_webhooktrigger = `
resource "tencentcloud_tcr_instance" "mytcr_webhooktrigger" {
  name        = "tf-test-tcr-%s"
  instance_type = "basic"
  delete_bucket = true

  tags ={
	test = "test"
  }
}

resource "tencentcloud_tcr_namespace" "my_ns" {
	instance_id 	 = tencentcloud_tcr_instance.mytcr_webhooktrigger.id
	name			 = "tf_test_ns_%s"
	is_public		 = true
	is_auto_scan	 = true
	is_prevent_vul = true
	severity		 = "medium"
	cve_whitelist_items	{
	  cve_id = "cve-xxxxx"
	}
  }

  data "tencentcloud_tcr_namespaces" "id_test" {
	instance_id = tencentcloud_tcr_namespace.my_ns.instance_id
  }

  locals {
    ns_id = data.tencentcloud_tcr_namespaces.id_test.namespace_list.0.id
  }

`

const testAccTCRWebhookTrigger = testAccTCRInstance_webhooktrigger + `

resource "tencentcloud_tcr_webhook_trigger" "my_trigger" {
  registry_id = tencentcloud_tcr_instance.mytcr_webhooktrigger.id
  namespace = tencentcloud_tcr_namespace.my_ns.name
  trigger {
		name = "trigger-%s"
		targets {
			address = "http://example.org/post"
			headers {
				key = "X-Custom-Header"
				values = ["a"]
			}
		}
		event_types = ["pushImage"]
		condition = ".*"
		enabled = true
		description = "this is trigger description"
		namespace_id = local.ns_id

  }
  tags = {
    "createdBy" = "terraform"
  }
}

`

const testAccTCRWebhookTrigger_update = testAccTCRInstance_webhooktrigger + `

resource "tencentcloud_tcr_webhook_trigger" "my_trigger" {
  registry_id = tencentcloud_tcr_instance.mytcr_webhooktrigger.id
  namespace = tencentcloud_tcr_namespace.my_ns.name
  trigger {
		name = "trigger-%s"
		targets {
			address = "http://example.org/post"
			headers {
				key = "X-Custom-Header"
				values = ["abc"]
			}
		}
		targets {
			address = "http://example.org/get"
			headers {
				key = "X-Custom-Header"
				values = ["xxx"]
			}
		}
		event_types = ["deleteImage"]
		condition = ".*test"
		enabled = false
		description = "this is trigger description deleted"
		namespace_id = local.ns_id

  }
  tags = {
    "createdBy" = "terraform"
  }
}

`
