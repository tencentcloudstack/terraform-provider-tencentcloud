package tencentcloud

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
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

var TestAccDataSourceTencentCloudGaapCertificatesBasic = fmt.Sprintf(`
resource tencentcloud_gaap_certificate "foo" {
  type    = "SERVER"
  name    = "ci-server-ca"
  content = %s
  key     = %s
}

data "tencentcloud_gaap_certificates" "foo" {
  id = "${tencentcloud_gaap_certificate.foo.id}"
}
`, "<<EOF"+testAccGaapCertificateServerCert+"EOF", "<<EOF"+testAccGaapCertificateServerKey+"EOF")

// fuzzy search
var TestAccDataSourceTencentCloudGaapCertificatesName = fmt.Sprintf(`
resource tencentcloud_gaap_certificate "foo" {
  type    = "SERVER"
  name    = "ci-server-ca"
  content = %s
  key     = %s
}

data "tencentcloud_gaap_certificates" "foo" {
  name = "${tencentcloud_gaap_certificate.foo.name}"
}
`, "<<EOF"+testAccGaapCertificateServerCert+"EOF", "<<EOF"+testAccGaapCertificateServerKey+"EOF")

var TestAccDataSourceTencentCloudGaapCertificatesType = fmt.Sprintf(
	`
resource tencentcloud_gaap_certificate "foo" {
  type    = "SERVER"
  name    = "ci-server-ca"
  content = %s
  key     = %s
}

data "tencentcloud_gaap_certificates" "foo" {
  type = "${tencentcloud_gaap_certificate.foo.type}"
}
`, "<<EOF"+testAccGaapCertificateServerCert+"EOF", "<<EOF"+testAccGaapCertificateServerKey+"EOF")
