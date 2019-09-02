package tencentcloud

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccDataSourceTencentCloudGaapCertificates_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: TestAccDataSourceTencentCloudGaapCertificatesBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_gaap_certificates.foo"),
					resource.TestCheckResourceAttr("data.tencentcloud_gaap_certificates.foo", "certificates.#", "1"),
					resource.TestCheckResourceAttr("data.tencentcloud_gaap_certificates.foo", "certificates.0.name", "ci-server-ca"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_certificates.foo", "certificates.0.type"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_certificates.foo", "certificates.0.create_time"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_certificates.foo", "certificates.0.begin_time"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_certificates.foo", "certificates.0.end_time"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_certificates.foo", "certificates.0.issuer_cn"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_certificates.foo", "certificates.0.subject_cn"),
				),
			},
		},
	})
}

func TestAccDataSourceTencentCloudGaapCertificates_name(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: TestAccDataSourceTencentCloudGaapCertificatesName,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_gaap_certificates.foo"),
					resource.TestMatchResourceAttr("data.tencentcloud_gaap_certificates.foo", "certificates.#", regexp.MustCompile(`^[1-9]\d*$`)),
					resource.TestMatchResourceAttr("data.tencentcloud_gaap_certificates.foo", "certificates.0.name", regexp.MustCompile("ci-server-ca")),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_certificates.foo", "certificates.0.type"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_certificates.foo", "certificates.0.create_time"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_certificates.foo", "certificates.0.begin_time"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_certificates.foo", "certificates.0.end_time"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_certificates.foo", "certificates.0.issuer_cn"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_certificates.foo", "certificates.0.subject_cn"),
				),
			},
		},
	})
}

func TestAccDataSourceTencentCloudGaapCertificates_type(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: TestAccDataSourceTencentCloudGaapCertificatesType,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_gaap_certificates.foo"),
					resource.TestMatchResourceAttr("data.tencentcloud_gaap_certificates.foo", "certificates.#", regexp.MustCompile(`^[1-9]\d*$`)),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_certificates.foo", "certificates.0.name"),
					resource.TestCheckResourceAttr("data.tencentcloud_gaap_certificates.foo", "certificates.0.type", "SERVER"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_certificates.foo", "certificates.0.create_time"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_certificates.foo", "certificates.0.begin_time"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_certificates.foo", "certificates.0.end_time"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_certificates.foo", "certificates.0.issuer_cn"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_certificates.foo", "certificates.0.subject_cn"),
				),
			},
		},
	})
}

var TestAccDataSourceTencentCloudGaapCertificatesBasic = testAccGaapCertificate("SERVER", "<<EOF\n"+testAccGaapCertificateServerCert+"EOF", "ci-server-ca", "<<EOF\n"+testAccGaapCertificateServerKey+"EOF") + `
data "tencentcloud_gaap_certificates" "foo" {
  id = "${tencentcloud_gaap_certificate.foo.id}"
}
`

// fuzzy search
var TestAccDataSourceTencentCloudGaapCertificatesName = testAccGaapCertificate("SERVER", "<<EOF\n"+testAccGaapCertificateServerCert+"EOF", "test-ci-server-ca-test", "<<EOF\n"+testAccGaapCertificateServerKey+"EOF") + `
data "tencentcloud_gaap_certificates" "foo" {
  name = "${tencentcloud_gaap_certificate.foo.name}"
}
`

var TestAccDataSourceTencentCloudGaapCertificatesType = testAccGaapCertificate("SERVER", "<<EOF\n"+testAccGaapCertificateServerCert+"EOF", "ci-server-ca", "<<EOF\n"+testAccGaapCertificateServerKey+"EOF") + `
data "tencentcloud_gaap_certificates" "foo" {
  type = "${tencentcloud_gaap_certificate.foo.type}"
}
`
