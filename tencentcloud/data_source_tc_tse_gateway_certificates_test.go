package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -test.run TestAccTencentCloudTseGatewayCertificatesDataSource_basic -v
func TestAccTencentCloudTseGatewayCertificatesDataSource_basic(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTseGatewayCertificatesDataSource,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_tse_gateway_certificates.gateway_certificates"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tse_gateway_certificates.gateway_certificates", "result.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tse_gateway_certificates.gateway_certificates", "result.0.certificates_list.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tse_gateway_certificates.gateway_certificates", "result.0.certificates_list.0.bind_domains.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tse_gateway_certificates.gateway_certificates", "result.0.certificates_list.0.cert_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tse_gateway_certificates.gateway_certificates", "result.0.certificates_list.0.id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tse_gateway_certificates.gateway_certificates", "result.0.certificates_list.0.cert_source"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tse_gateway_certificates.gateway_certificates", "result.0.certificates_list.0.create_time"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tse_gateway_certificates.gateway_certificates", "result.0.certificates_list.0.crt"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tse_gateway_certificates.gateway_certificates", "result.0.certificates_list.0.expire_time"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tse_gateway_certificates.gateway_certificates", "result.0.certificates_list.0.issue_time"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tse_gateway_certificates.gateway_certificates", "result.0.certificates_list.0.key"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tse_gateway_certificates.gateway_certificates", "result.0.certificates_list.0.name"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tse_gateway_certificates.gateway_certificates", "result.0.certificates_list.0.status"),
				),
			},
		},
	})
}

const testAccTseGatewayCertificatesDataSource = testAccTseCngwCertificate + `

data "tencentcloud_tse_gateway_certificates" "gateway_certificates" {
  gateway_id = var.gateway_id
  filters {
    key = "BindDomain"
    value = "example.com"
  }
  depends_on = [ tencentcloud_tse_cngw_certificate.cngw_certificate ]
}

`
