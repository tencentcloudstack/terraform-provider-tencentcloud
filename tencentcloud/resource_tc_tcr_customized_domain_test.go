package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudTcrCustomizedDomainResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTcrCustomizedDomain,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_tcr_customized_domain.customized_domain", "id")),
			},
			{
				ResourceName:      "tencentcloud_tcr_customized_domain.customized_domain",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccTcrCustomizedDomain = `

resource "tencentcloud_tcr_customized_domain" "customized_domain" {
  registry_id = "tcr-xxx"
  domain_name = "xxx.test.com"
  certificate_id = "kWGTVuU3"
  tags = {
    "createdBy" = "terraform"
  }
}

`
