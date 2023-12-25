package teo_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -test.run TestAccTencentCloudTeoOwnershipVerifyResource_basic -v
func TestAccTencentCloudTeoOwnershipVerifyResource_basic(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_PRIVATE) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckZoneDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTeoOwnershipVerify,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_teo_ownership_verify.ownership_verify", "id"),
					resource.TestCheckResourceAttr("tencentcloud_teo_ownership_verify.ownership_verify", "domain", "tf-teo.xyz"),
					resource.TestCheckResourceAttr("tencentcloud_teo_ownership_verify.ownership_verify", "result", "ok"),
					resource.TestCheckResourceAttr("tencentcloud_teo_ownership_verify.ownership_verify", "status", "success"),
				),
			},
		},
	})
}

const testAccTeoOwnershipVerify = testAccTeoZone + `

resource "tencentcloud_teo_ownership_verify" "ownership_verify" {
  domain = var.zone_name

  depends_on = [ tencentcloud_teo_zone.basic ]
}

`
