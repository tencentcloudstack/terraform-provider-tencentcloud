package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudTseGatewayCertificatesDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTseGatewayCertificatesDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_tse_gateway_certificates.gateway_certificates")),
			},
		},
	})
}

const testAccTseGatewayCertificatesDataSource = `

data "tencentcloud_tse_gateway_certificates" "gateway_certificates" {
  gateway_id = ""
  filters {
		key = ""
		value = ""

  }
  }

`
