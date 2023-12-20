package clb_test

import (
	"fmt"
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudNeedFixClbReplaceCertForLbsResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccClbReplaceCertForLbs,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_clb_replace_cert_for_lbs.replace_cert_for_lbs", "id")),
			},
		},
	})
}

var testAccClbReplaceCertForLbs = fmt.Sprintf(`

resource "tencentcloud_clb_replace_cert_for_lbs" "replace_cert_for_lbs" {
  old_certificate_id = "zjUMifFK"
  certificate {
    cert_ca_name = "test"
	cert_ca_content = "%s"
  }
}
`, testAccSslCertificateCA)

const testAccSslCertificateCA = "-----BEGIN CERTIFICATE-----\\nMIIERzCCAq+gAwIBAgIBAjANBgkqhkiG9w0BAQsF" +
	"ADAoMQ0wCwYDVQQDEwR0ZXN0\\nMRcwFQYDVQQKEw50ZXJyYWZvcm0gdGVzdDAeFw0xOTA4MTM" +
	"wMzE5MzlaFw0yOTA4\\nMTAwMzE5MzlaMC4xEzARBgNVBAMTCnNlcnZlciBzc2wxFzAVBgNVBA" +
	"oTDnRlcnJh\\nZm9ybS10ZXN0MIIBojANBgkqhkiG9w0BAQEFAAOCAY8AMIIBigKCAYEA1Ryp+" +
	"DKK\\nSNFKZsPtwfR+jzOnQ8YFieIKYgakV688d8YgpolenbmeEPrzT87tunFD7G9f6ALG\\nND" +
	"8rj7npj0AowxhOL/h/v1D9u0UsIaj5i2GWJrqNAhGLaxWiEB/hy5WOiwxDrGei\\ngQqJkFM52" +
	"Ep7G1Yx7PHJmKFGwN9FhIsFi1cNZfVRopZuCe/RMPNusNVZaIi+qcEf\\nfsE1cmfmuSlG3Ap0" +
	"RKOIyR0ajDEzqZn9/0R7VwWCF97qy8TNYk94K/1tq3zyhVzR\\nZ83xOSfrTqEfb3so3AU2jyK" +
	"gYdwr/FZS72VCHS8IslgnqJW4izIXZqgIKmHaRZtM\\nN4jUloi6l/6lktt6Lsgh9xECecxziS" +
	"JtPMaog88aC8HnMqJJ3kScGCL36GYG+Kaw\\n5PnDlWXBaeiDe8z/eWK9+Rr2M+rhTNxosAVGf" +
	"DJyxAXyiX49LQ0v7f9qzwc/0JiD\\nbvsUv1cm6OgpoEMP9SXqqBdwGqeKbD2/2jlP48xlYP6l" +
	"1SoJG3GgZ8dbAgMBAAGj\\ndjB0MAwGA1UdEwEB/wQCMAAwEwYDVR0lBAwwCgYIKwYBBQUHAwE" +
	"wDwYDVR0PAQH/\\nBAUDAweAADAdBgNVHQ4EFgQULwWKBQNLL9s3cb3tTnyPVg+mpCMwHwYDVR" +
	"0jBBgw\\nFoAUKwfrmq791mY831S6UHARHtgYnlgwDQYJKoZIhvcNAQELBQADggGBAMo5RglS\\nA" +
	"HdPgaicWJvmvjjexjF/42b7Rz4pPfMjYw6uYO8He/f4UZWv5CZLrbEe7MywaK3y\\n0OsfH8Ah" +
	"yN29pv2x8g9wbmq7omZIOZ0oCAGduEXs/A/qY/hFaCohdkz/IN8qi6JW\\nVXreGli3SrpcHFc" +
	"hSwHTyJEXgkutcGAsOvdsOuVSmplOyrkLHc8uUe8SG4j8kGyg\\nEzaszFjHkR7g1dVyDVUedc" +
	"588mjkQxYeAamJgfkgIhljWKMa2XzkVMcVfQHfNpM1\\nn+bu8SmqRt9Wma2bMijKRG/Blm756" +
	"LoI+skY+WRZmlDnq8zj95TT0vceGP0FUWh5\\nhKyiocABmpQs9OK9HMi8vgSWISP+fYgkm/bK" +
	"tKup2NbZBoO5/VL2vCEPInYzUhBO\\njCbLMjNjtM5KriCaR7wDARgHiG0gBEPOEW1PIjZ9UOH" +
	"+LtIxbNZ4eEIIINLHnBHf\\nL+doVeZtS/gJc4G4Adr5HYuaS9ZxJ0W2uy0eQlOHzjyxR6Mf/r" +
	"pnilJlcQ==\\n-----END CERTIFICATE-----"
