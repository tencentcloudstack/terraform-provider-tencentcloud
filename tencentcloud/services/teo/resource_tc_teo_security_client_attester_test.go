package teo_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudTeoSecurityClientAttesterResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTeoSecurityClientAttester,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_teo_security_client_attester.example", "id"),
					resource.TestCheckResourceAttr("tencentcloud_teo_security_client_attester.example", "client_attesters.0.name", "tf-example-attester"),
					resource.TestCheckResourceAttr("tencentcloud_teo_security_client_attester.example", "client_attesters.0.attester_source", "TC-RCE"),
					resource.TestCheckResourceAttr("tencentcloud_teo_security_client_attester.example", "client_attesters.0.attester_duration", "300s"),
					resource.TestCheckResourceAttrSet("tencentcloud_teo_security_client_attester.example", "client_attesters.0.id"),
					resource.TestCheckResourceAttrSet("tencentcloud_teo_security_client_attester.example", "client_attesters.0.type"),
				),
			},
			{
				Config: testAccTeoSecurityClientAttesterUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("tencentcloud_teo_security_client_attester.example", "client_attesters.0.name", "tf-example-attester-updated"),
					resource.TestCheckResourceAttr("tencentcloud_teo_security_client_attester.example", "client_attesters.0.attester_duration", "600s"),
				),
			},
			{
				ResourceName:      "tencentcloud_teo_security_client_attester.example",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccTeoSecurityClientAttester = `
resource "tencentcloud_teo_security_client_attester" "example" {
  zone_id = "zone-3fkff38fyw8s"

  client_attesters {
    name              = "tf-example"
    attester_source   = "TC-RCE"
    attester_duration = "300s"

    tc_rce_option {
      channel = "12399223"
      region  = "ap-beijing"
    }
  }
}
`

const testAccTeoSecurityClientAttesterUpdate = `
resource "tencentcloud_teo_security_client_attester" "example" {
  zone_id = "zone-3fkff38fyw8s"

  client_attesters {
    name              = "tf-example0-update"
    attester_source   = "TC-RCE"
    attester_duration = "600s"

    tc_rce_option {
      channel = "12399223"
      region  = "ap-beijing"
    }
  }
}
`
