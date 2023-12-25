package tcr_test

import (
	"fmt"
	"testing"
	"time"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudTcrServiceAccountResource_basic(t *testing.T) {
	t.Parallel()
	loc, _ := time.LoadLocation("Asia/Chongqing")
	expireTime := time.Now().AddDate(0, 0, 10).In(loc).UnixMilli()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_COMMON) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTcrServiceAccount,
				PreConfig: func() {
					tcacctest.AccStepSetRegion(t, "ap-shanghai")
					tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_COMMON)
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_tcr_service_account.example", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_tcr_service_account.example", "registry_id"),
					resource.TestCheckResourceAttr("tencentcloud_tcr_service_account.example", "name", "tf_example_account"),
					resource.TestCheckResourceAttrSet("tencentcloud_tcr_service_account.example", "permissions.#"),
					resource.TestCheckResourceAttr("tencentcloud_tcr_service_account.example", "permissions.0.resource", "tf_test_tcr_namespace"),
					resource.TestCheckResourceAttrSet("tencentcloud_tcr_service_account.example", "permissions.0.actions.#"),
					resource.TestCheckTypeSetElemAttr("tencentcloud_tcr_service_account.example", "permissions.0.actions.*", "tcr:PushRepository"),
					resource.TestCheckTypeSetElemAttr("tencentcloud_tcr_service_account.example", "permissions.0.actions.*", "tcr:PullRepository"),
					resource.TestCheckResourceAttr("tencentcloud_tcr_service_account.example", "description", "tf example for tcr service account"),
					resource.TestCheckResourceAttr("tencentcloud_tcr_service_account.example", "duration", "10"),
					resource.TestCheckResourceAttr("tencentcloud_tcr_service_account.example", "disable", "false"),
				),
			},
			{
				Config: fmt.Sprintf(testAccTcrServiceAccount_Update, expireTime),
				PreConfig: func() {
					tcacctest.AccStepSetRegion(t, "ap-shanghai")
					tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_COMMON)
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_tcr_service_account.example", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_tcr_service_account.example", "permissions.#"),
					resource.TestCheckResourceAttr("tencentcloud_tcr_service_account.example", "permissions.0.resource", "tf_test_tcr_namespace"),
					resource.TestCheckResourceAttrSet("tencentcloud_tcr_service_account.example", "permissions.0.actions.#"),
					resource.TestCheckTypeSetElemAttr("tencentcloud_tcr_service_account.example", "permissions.0.actions.*", "tcr:PushRepository"),
					resource.TestCheckResourceAttr("tencentcloud_tcr_service_account.example", "description", "CHANGED tf example for tcr service account"),
					resource.TestCheckResourceAttrSet("tencentcloud_tcr_service_account.example", "expires_at"),
					// resource.TestCheckResourceAttr("tencentcloud_tcr_service_account.example", "disable", "true"),
				),
			},
			{
				ResourceName:            "tencentcloud_tcr_service_account.example",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"duration", "password"},
			},
		},
	})
}

const testAccTcrServiceAccount = `

resource "tencentcloud_tcr_instance" "example" {
	name          = "tf-example-tcr-instance"
	instance_type = "premium"
	delete_bucket = true
	tags = {
	  "createdBy" = "terraform"
	}
  }
  
  resource "tencentcloud_tcr_namespace" "example" {
	instance_id    = tencentcloud_tcr_instance.example.id
	name           = "tf_test_tcr_namespace"
	is_public      = true
	is_auto_scan   = true
	is_prevent_vul = true
	severity       = "medium"
	cve_whitelist_items {
	  cve_id = "tf_example_cve_id"
	}
  }
  
  resource "tencentcloud_tcr_service_account" "example" {
	registry_id = tencentcloud_tcr_instance.example.id
	name        = "tf_example_account"
	permissions {
	  resource = tencentcloud_tcr_namespace.example.name
	  actions  = ["tcr:PushRepository", "tcr:PullRepository"]
	}
	description = "tf example for tcr service account"
	duration    = 10
	disable     = false
	tags = {
	  "createdBy" = "terraform"
	}
  }

`

const testAccTcrServiceAccount_Update = `

resource "tencentcloud_tcr_instance" "example" {
	name          = "tf-example-tcr-instance"
	instance_type = "premium"
	delete_bucket = true
	tags = {
	  "createdBy" = "terraform"
	}
  }
  
  resource "tencentcloud_tcr_namespace" "example" {
	instance_id    = tencentcloud_tcr_instance.example.id
	name           = "tf_test_tcr_namespace"
	is_public      = true
	is_auto_scan   = true
	is_prevent_vul = true
	severity       = "medium"
	cve_whitelist_items {
	  cve_id = "tf_example_cve_id"
	}
  }
  
  resource "tencentcloud_tcr_service_account" "example" {
	registry_id = tencentcloud_tcr_instance.example.id
	name        = "tf_example_account"
	permissions {
	  resource = tencentcloud_tcr_namespace.example.name
	  actions  = ["tcr:PushRepository"]
	}
	description = "CHANGED tf example for tcr service account"
	expires_at  = %d
	disable     = false
	tags = {
	  "createdBy" = "terraform"
	}
  }

`
