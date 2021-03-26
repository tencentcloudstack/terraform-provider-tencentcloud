package tencentcloud

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"

	sdkError "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	ssl "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ssl/v20191205"
)

func TestAccTencentCloudSslCertificate_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSslCertificateDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSslCertificate("CA", testAccSslCertificateCA, "CA", ""),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSslCertificateExists("tencentcloud_ssl_certificate.foo"),
					resource.TestCheckResourceAttr("tencentcloud_ssl_certificate.foo", "name", "CA"),
					resource.TestCheckResourceAttr("tencentcloud_ssl_certificate.foo", "type", "CA"),
					resource.TestCheckResourceAttr("tencentcloud_ssl_certificate.foo", "project_id", "0"),
					resource.TestCheckResourceAttrSet("tencentcloud_ssl_certificate.foo", "cert"),
					resource.TestCheckNoResourceAttr("tencentcloud_ssl_certificate.foo", "key"),
					resource.TestCheckResourceAttrSet("tencentcloud_ssl_certificate.foo", "product_zh_name"),
					resource.TestCheckResourceAttrSet("tencentcloud_ssl_certificate.foo", "domain"),
					resource.TestCheckResourceAttrSet("tencentcloud_ssl_certificate.foo", "status"),
					resource.TestCheckResourceAttrSet("tencentcloud_ssl_certificate.foo", "begin_time"),
					resource.TestCheckResourceAttrSet("tencentcloud_ssl_certificate.foo", "end_time"),
					resource.TestCheckResourceAttrSet("tencentcloud_ssl_certificate.foo", "create_time"),
					resource.TestCheckResourceAttr("tencentcloud_ssl_certificate.foo", "subject_names.#", "0"),
				),
			},
			{
				ResourceName:      "tencentcloud_ssl_certificate.foo",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccTencentCloudSslCertificate_svr(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSslCertificateDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSslCertificate("SVR", testAccSslCertificateCA, "server", testAccSslCertificateKey),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSslCertificateExists("tencentcloud_ssl_certificate.foo"),
					resource.TestCheckResourceAttr("tencentcloud_ssl_certificate.foo", "name", "server"),
					resource.TestCheckResourceAttr("tencentcloud_ssl_certificate.foo", "type", "SVR"),
					resource.TestCheckResourceAttr("tencentcloud_ssl_certificate.foo", "project_id", "0"),
					resource.TestCheckResourceAttrSet("tencentcloud_ssl_certificate.foo", "cert"),
					resource.TestCheckResourceAttrSet("tencentcloud_ssl_certificate.foo", "key"),
					resource.TestCheckResourceAttrSet("tencentcloud_ssl_certificate.foo", "product_zh_name"),
					resource.TestCheckResourceAttrSet("tencentcloud_ssl_certificate.foo", "domain"),
					resource.TestCheckResourceAttrSet("tencentcloud_ssl_certificate.foo", "status"),
					resource.TestCheckResourceAttrSet("tencentcloud_ssl_certificate.foo", "begin_time"),
					resource.TestCheckResourceAttrSet("tencentcloud_ssl_certificate.foo", "end_time"),
					resource.TestCheckResourceAttrSet("tencentcloud_ssl_certificate.foo", "create_time"),
					resource.TestCheckResourceAttr("tencentcloud_ssl_certificate.foo", "subject_names.#", "0"),
				),
			},
		},
	})
}

func testAccCheckSslCertificateDestroy(s *terraform.State) error {
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	sslService := SSLService{
		client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn,
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_ssl_certificate" {
			continue
		}
		resourceId := rs.Primary.ID
		describeRequest := ssl.NewDescribeCertificateDetailRequest()
		describeRequest.CertificateId = helper.String(resourceId)
		var (
			describeResponse *ssl.DescribeCertificateDetailResponse
			outErr, inErr    error
		)
		outErr = resource.Retry(readRetryTimeout, func() *resource.RetryError {
			describeResponse, inErr = sslService.DescribeCertificateDetail(ctx, describeRequest)
			if inErr != nil {
				if sdkErr, ok := inErr.(*sdkError.TencentCloudSDKError); ok {
					if sdkErr.Code == CertificateNotFound {
						return nil
					}
				}
				return retryError(inErr)
			}
			return nil
		})
		if outErr != nil {
			return outErr
		}

		if describeResponse != nil && describeResponse.Response != nil && describeResponse.Response.CertificateId != nil {
			return errors.New("certificate still exists")
		}
	}
	return nil
}

func testAccCheckSslCertificateExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)
		sslService := SSLService{
			client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn,
		}

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("[TECENT_TERRAFORM_CHECK][SSL certificate][Exists] check: SSL certificate %s is not found", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("[TECENT_TERRAFORM_CHECK][SSL certificate][Exists] check: SSL certificate certificateId is not set")
		}
		resourceId := rs.Primary.ID
		describeRequest := ssl.NewDescribeCertificateDetailRequest()
		describeRequest.CertificateId = helper.String(resourceId)
		var (
			describeResponse *ssl.DescribeCertificateDetailResponse
			outErr, inErr    error
		)
		outErr = resource.Retry(readRetryTimeout, func() *resource.RetryError {
			describeResponse, inErr = sslService.DescribeCertificateDetail(ctx, describeRequest)
			if inErr != nil {
				if sdkErr, ok := inErr.(*sdkError.TencentCloudSDKError); ok {
					if sdkErr.Code == CertificateNotFound {
						return nil
					}
				}
				return retryError(inErr)
			}
			return nil
		})
		if outErr != nil {
			return outErr
		}

		if describeResponse == nil || describeResponse.Response == nil || describeResponse.Response.CertificateId == nil {
			return fmt.Errorf("certificateId %s does not exist", resourceId)
		}

		return nil
	}
}

func testAccSslCertificate(certificateType string, cert, name, key string) string {
	const str = `
resource "tencentcloud_ssl_certificate" "foo" {
  type = "%s"
  cert = "%s"
  %s
  %s
}

`

	if name != "" {
		name = "name = \"" + name + "\""
	}

	if key != "" {
		key = fmt.Sprintf("key = \"%s\"", key)
	}

	return fmt.Sprintf(str, certificateType, cert, name, key)
}

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

const testAccSslCertificateKey = "Public Key Info:\\n	Public Key Algorithm: RSA\\n	Key Security Level:" +
	" High (3072 bits)\\nmodulus:\\n	00:d5:1c:a9:f8:32:8a:48:d1:4a:66:c3:ed:" +
	"c1:f4:7e\\n	8f:33:a7:43:c6:05:89:e2:0a:62:06:a4:57:af:3c:77\\n	c6:20:a" +
	"6:89:5e:9d:b9:9e:10:fa:f3:4f:ce:ed:ba:71\\n	43:ec:6f:5f:e8:02:c6:34:3f:" +
	"2b:8f:b9:e9:8f:40:28\\n	c3:18:4e:2f:f8:7f:bf:50:fd:bb:45:2c:21:a8:f9:8" +
	"b\\n	61:96:26:ba:8d:02:11:8b:6b:15:a2:10:1f:e1:cb:95\\n	8e:8b:0c:43" +
	":ac:67:a2:81:0a:89:90:53:39:d8:4a:7b\\n	1b:56:31:ec:f1:c9:98:a1:46:c0:df" +
	":45:84:8b:05:8b\\n	57:0d:65:f5:51:a2:96:6e:09:ef:d1:30:f3:6e:b0:d5\\n	" +
	"59:68:88:be:a9:c1:1f:7e:c1:35:72:67:e6:b9:29:46\\n	dc:0a:74:44:a3:88:c9" +
	":1d:1a:8c:31:33:a9:99:fd:ff\\n	44:7b:57:05:82:17:de:ea:cb:c4:cd:62:4f:7" +
	"8:2b:fd\\n	6d:ab:7c:f2:85:5c:d1:67:cd:f1:39:27:eb:4e:a1:1f\\n	6f:7b:28" +
	":dc:05:36:8f:22:a0:61:dc:2b:fc:56:52:ef\\n	65:42:1d:2f:08:b2:58:27:a8:9" +
	"5:b8:8b:32:17:66:a8\\n	08:2a:61:da:45:9b:4c:37:88:d4:96:88:ba:97:fe:a" +
	"5\\n	92:db:7a:2e:c8:21:f7:11:02:79:cc:73:89:22:6d:3c\\n	c6:a8:83:cf:" +
	"1a:0b:c1:e7:32:a2:49:de:44:9c:18:22\\n	f7:e8:66:06:f8:a6:b0:e4:f9:c3:95" +
	":65:c1:69:e8:83\\n	7b:cc:ff:79:62:bd:f9:1a:f6:33:ea:e1:4c:dc:68:b0\\n	" +
	"05:46:7c:32:72:c4:05:f2:89:7e:3d:2d:0d:2f:ed:ff\\n	6a:cf:07:3f:d0:98:83" +
	":6e:fb:14:bf:57:26:e8:e8:29\\n	a0:43:0f:f5:25:ea:a8:17:70:1a:a7:8a:6c:3" +
	"d:bf:da\\n	39:4f:e3:cc:65:60:fe:a5:d5:2a:09:1b:71:a0:67:c7\\n	5b:\\npub" +
	"lic exponent:\\n	01:00:01:\\nprivate exponent:\\n	00:b1:56:d0:fa:00:d4" +
	":a2:13:c7:5e:0c:dc:e4:f1:97\\n	ff:82:74:46:29:9a:a2:4a:bf:69:23:2d:ce:e" +
	"9:bb:df\\n	cf:b7:8b:dd:f4:26:3c:38:14:d9:3f:6f:c2:3a:81:53\\n	8f:ba:48" +
	":53:fe:b5:90:4a:19:e7:1e:0b:0f:18:6d:c3\\n	7d:d5:d3:fa:87:47:86:e4:d6:b" +
	"f:e7:a7:f9:ba:ab:2e\\n	19:5e:e1:8b:8b:9b:95:0d:f7:66:61:1e:19:e9:c3:8" +
	"8\\n	08:be:1c:ce:93:c1:09:b1:68:1b:61:46:60:74:64:46\\n	5d:51:34:ea:7" +
	"f:a9:ca:a1:2a:47:85:84:4b:ef:84:05\\n	97:c3:46:7d:06:19:ce:24:73:90:64:" +
	"fb:df:16:d5:80\\n	34:8e:90:7c:58:b6:a4:86:ce:30:b3:ab:52:8b:f2:95\\n	4" +
	"c:b6:46:5a:77:db:73:c0:0c:3f:6d:12:18:a8:54:7c\\n	ff:77:c3:ca:89:9f:63:" +
	"98:ef:48:2d:c1:09:70:6e:ea\\n	cb:bb:78:91:42:8a:22:3e:21:ef:a5:bf:16:ee" +
	":66:45\\n	e5:f0:26:6a:85:8e:e1:69:62:ac:05:00:a6:44:ba:c8\\n	ac:10:00:" +
	"97:f5:51:65:7f:9a:1f:7b:99:9d:02:d4:87\\n	50:ce:74:06:51:67:fa:fb:90:e4" +
	":33:79:f2:a8:61:ee\\n	45:1d:87:ca:22:5b:ac:e7:32:38:f8:2c:fd:55:92:1" +
	"e\\n	3d:60:1e:7c:4b:fd:28:ff:e5:b1:02:6a:aa:22:f7:ae\\n	a8:36:90:7b:a" +
	"6:f7:29:05:14:3a:21:da:36:05:f9:b0\\n	9d:f7:fb:10:75:d7:2c:21:32:95:e7:" +
	"f7:17:be:09:cb\\n	66:fe:f1:69:71:df:a4:5e:3f:fc:67:6c:37:65:b8:51\\n	c" +
	"6:22:38:fb:07:ce:89:54:50:43:71:44:3d:c3:51:5a\\n	bd:e5:c7:87:b2:ea:7b:" +
	"64:0f:5d:34:9c:a1:52:b3:ce\\n	06:86:ba:7a:05:80:48:b3:6c:1b:79:74:9b:49" +
	":f2:30\\n	c1:\\nprime1:\\n	00:e4:31:46:59:3d:24:f7:31:d9:22:26:af:c0:3e:" +
	"f5\\n	c1:6d:be:ba:d3:9e:3f:b9:2c:43:a0:d0:47:09:e4:35\\n	63:19:a4:33:8" +
	"2:af:f9:76:3c:11:c2:cb:34:f9:a6:ab\\n	dd:ab:64:5a:6b:9c:c1:2a:52:89:64:" +
	"7e:b5:a7:f0:4d\\n	29:13:a4:cf:17:f4:f2:0d:a0:6e:b9:5d:95:41:10:df\\n	a" +
	"e:f3:7a:13:49:21:66:73:2a:b7:e2:8d:7c:c2:34:e5\\n	3f:bd:78:ca:fc:64:c5:" +
	"1c:3a:66:7a:12:53:96:bd:b0\\n	c3:7a:0c:ec:5e:55:c0:c3:3f:7f:25:72:f4:e2" +
	":19:94\\n	9d:65:15:be:c8:82:20:57:12:97:b2:a8:4d:3d:e0:8f\\n	e2:1f:d0:" +
	"c8:49:aa:f4:34:fa:91:d1:d1:cc:98:bc:3d\\n	8b:b1:9b:8f:fd:ef:03:dd:92:fb" +
	":ca:99:45:af:cc:83\\n	58:4c:bb:ba:73:9e:23:84:f9:7e:4f:40:fe:00:b5:b" +
	"f\\n	6f:\\nprime2:\\n	00:ef:14:ef:73:fc:0c:fc:e3:87:d9:7f:a6:f8:55:8" +
	"6\\n	57:63:8a:86:87:f5:ef:63:20:1f:b2:ae:28:dc:ab:59\\n	80:8f:15:64:4" +
	"4:e2:bc:a5:7b:d7:69:ef:30:b1:83:b3\\n	bd:09:fd:4a:0c:c6:31:5b:a4:79:d0:" +
	"e5:d3:a8:31:fd\\n	59:ea:52:63:cf:17:a7:c1:54:bf:a8:11:9b:b1:85:47\\n	5" +
	"a:08:a6:9c:2f:47:9d:ac:5d:e8:7c:e4:31:6c:99:71\\n	04:7d:20:98:be:8b:60:" +
	"07:66:2d:b9:41:10:ea:dd:5b\\n	87:20:65:62:ea:75:a7:a6:04:a2:18:66:6b:db" +
	":5b:a4\\n	9f:12:97:cb:7c:8c:d2:e0:ce:02:ef:1e:df:a1:9d:6a\\n	bc:00:38:" +
	"18:36:a1:c5:97:16:be:7a:df:5f:4f:4f:de\\n	a3:cb:25:fe:f6:67:0d:31:aa:0a" +
	":d4:1b:be:df:91:2c\\n	05:14:20:37:cc:4f:50:33:a6:50:1b:90:f9:b2:08:8" +
	"0\\n	d5:\\ncoefficient:\\n	47:d1:7f:ca:93:6a:14:9b:fe:85:8d:c2:15:11:52:" +
	"a2\\n	a5:bc:f5:6b:a2:69:76:49:1e:09:79:f1:15:bf:39:48\\n	41:ff:92:78:a" +
	"f:bc:7d:6f:76:3b:32:9e:08:d2:42:06\\n	04:5f:36:e0:be:a8:1d:21:5c:ec:39:" +
	"09:e0:77:c5:86\\n	06:e6:ce:98:16:fc:0f:30:de:a7:69:7a:8f:dd:01:42\\n	2" +
	"a:22:f5:b7:c2:fc:c8:90:5a:78:dc:b3:e0:4d:e7:2d\\n	98:6c:e3:34:1b:d7:e8:" +
	"f8:90:57:7e:4d:41:d6:4a:29\\n	81:92:eb:89:5b:45:85:dd:b9:16:20:63:cb:59" +
	":f6:06\\n	59:c1:dd:3b:6b:92:0a:5e:5e:63:4a:f1:a7:d5:16:b9\\n	8b:6c:d8:" +
	"ad:76:0e:2d:3c:e0:b3:73:e0:6d:af:d4:a2\\n	bc:4b:fd:6c:2d:d7:5d:4d:cd:28" +
	":03:64:b2:ef:9a:1d\\n	82:8d:53:40:c5:f8:fb:f3:63:de:8e:1a:21:b6:35:1" +
	"4\\n	\\nexp1:\\n	00:9c:a5:8a:d2:65:dc:03:69:8f:d2:16:d6:9d:55:5b\\n	2" +
	"5:4e:ae:18:d8:7e:90:e6:10:11:d8:ca:41:89:f3:c4\\n	06:64:aa:c8:c5:95:01:" +
	"dd:fd:7c:7f:c9:39:73:8b:cb\\n	fd:9e:d3:84:12:cd:87:f9:02:b1:d8:6f:f7:49" +
	":f2:f7\\n	35:14:8c:15:b2:2f:6f:1e:95:9c:8c:d9:46:45:65:4c\\n	f8:6f:a1:" +
	"c4:ad:76:25:3b:37:ff:05:a1:f5:1b:e8:6d\\n	db:64:b9:10:37:55:01:ce:cf:f4" +
	":5b:26:4b:85:76:70\\n	6a:b0:55:40:c9:bd:7a:57:4e:36:7d:41:be:03:9c:6" +
	"5\\n	dd:ea:6f:94:09:56:f2:d6:73:27:f9:f7:f9:16:5a:1a\\n	cb:b2:e5:83:2" +
	"8:b7:17:6f:6a:f7:41:1f:11:a1:63:cf\\n	a8:1e:e3:58:64:8c:78:8d:d9:81:c9:" +
	"e1:8f:ea:0f:ad\\n	b6:a6:ee:54:1f:5c:56:ab:c9:0d:c1:60:2f:3d:d3:86\\n	3" +
	"7:\\nexp2:\\n	64:12:b7:48:2d:30:a2:89:fa:cb:27:8b:94:56:f8:2c\\n	8c:15" +
	":e7:c9:f1:3f:8a:96:5f:ae:43:08:07:96:11:98\\n	a6:4b:a5:f4:cf:93:77:11:2" +
	"7:51:c8:34:f1:98:d7:1b\\n	41:9b:2b:eb:bc:e9:dc:1a:34:83:24:30:3c:2e:f0:" +
	"85\\n	3a:77:d2:1f:55:1f:7a:e5:26:74:0b:2a:c8:5b:a9:4a\\n	1e:64:de:eb:4" +
	"b:66:cc:47:62:91:24:53:2b:c9:ee:6c\\n	9a:93:92:5b:ef:aa:fa:6d:e2:a5:b0:" +
	"7e:8c:50:ab:87\\n	1c:20:54:0f:1f:c0:54:d5:8b:a3:fa:fb:1a:8e:79:91\\n	b" +
	"c:0e:9d:b6:3c:9b:e8:4d:53:1d:14:27:37:56:d4:de\\n	6c:99:0e:49:8f:dd:4d:" +
	"28:d0:02:4e:8d:6e:7d:58:0b\\n	e7:74:b8:0c:1b:86:82:4b:52:cd:05:f0:17:54" +
	":84:c0\\n	7b:74:20:e6:fc:2b:ed:f2:a7:85:62:61:a2:0b:bd:21\\n	\\nPublic " +
	"Key PIN:\\n	pin-sha256:t5OXXC5gYqMNtUMsTqRs3A3vhfK2BiXVOgYzIEYv7Y8=\\nPubl" +
	"ic Key ID:\\n	sha256:b793975c2e6062a30db5432c4ea46cdc0def85f2b60625d53a" +
	"063320462fed8f\\n	sha1:2f058a05034b2fdb3771bded4e7c8f560fa6a423\\n-----B" +
	"EGIN RSA PRIVATE KEY-----\\nMIIG5AIBAAKCAYEA1Ryp+DKKSNFKZsPtwfR+jzOnQ8YFi" +
	"eIKYgakV688d8Ygpole\\nnbmeEPrzT87tunFD7G9f6ALGND8rj7npj0AowxhOL/h/v1D9u0U" +
	"sIaj5i2GWJrqN\\nAhGLaxWiEB/hy5WOiwxDrGeigQqJkFM52Ep7G1Yx7PHJmKFGwN9FhIsFi" +
	"1cNZfVR\\nopZuCe/RMPNusNVZaIi+qcEffsE1cmfmuSlG3Ap0RKOIyR0ajDEzqZn9/0R7VwW" +
	"C\\nF97qy8TNYk94K/1tq3zyhVzRZ83xOSfrTqEfb3so3AU2jyKgYdwr/FZS72VCHS8I\\nslg" +
	"nqJW4izIXZqgIKmHaRZtMN4jUloi6l/6lktt6Lsgh9xECecxziSJtPMaog88a\\nC8HnMqJJ3" +
	"kScGCL36GYG+Kaw5PnDlWXBaeiDe8z/eWK9+Rr2M+rhTNxosAVGfDJy\\nxAXyiX49LQ0v7f9" +
	"qzwc/0JiDbvsUv1cm6OgpoEMP9SXqqBdwGqeKbD2/2jlP48xl\\nYP6l1SoJG3GgZ8dbAgMBA" +
	"AECggGBALFW0PoA1KITx14M3OTxl/+CdEYpmqJKv2kj\\nLc7pu9/Pt4vd9CY8OBTZP2/COoF" +
	"Tj7pIU/61kEoZ5x4LDxhtw33V0/qHR4bk1r/n\\np/m6qy4ZXuGLi5uVDfdmYR4Z6cOICL4cz" +
	"pPBCbFoG2FGYHRkRl1RNOp/qcqhKkeF\\nhEvvhAWXw0Z9BhnOJHOQZPvfFtWANI6QfFi2pIb" +
	"OMLOrUovylUy2Rlp323PADD9t\\nEhioVHz/d8PKiZ9jmO9ILcEJcG7qy7t4kUKKIj4h76W/F" +
	"u5mReXwJmqFjuFpYqwF\\nAKZEusisEACX9VFlf5ofe5mdAtSHUM50BlFn+vuQ5DN58qhh7kU" +
	"dh8oiW6znMjj4\\nLP1Vkh49YB58S/0o/+WxAmqqIveuqDaQe6b3KQUUOiHaNgX5sJ33+xB11" +
	"ywhMpXn\\n9xe+Cctm/vFpcd+kXj/8Z2w3ZbhRxiI4+wfOiVRQQ3FEPcNRWr3lx4ey6ntkD10" +
	"0\\nnKFSs84Ghrp6BYBIs2wbeXSbSfIwwQKBwQDkMUZZPST3MdkiJq/APvXBbb66054/\\nuSx" +
	"DoNBHCeQ1YxmkM4Kv+XY8EcLLNPmmq92rZFprnMEqUolkfrWn8E0pE6TPF/Ty\\nDaBuuV2VQ" +
	"RDfrvN6E0khZnMqt+KNfMI05T+9eMr8ZMUcOmZ6ElOWvbDDegzsXlXA\\nwz9/JXL04hmUnWU" +
	"VvsiCIFcSl7KoTT3gj+If0MhJqvQ0+pHR0cyYvD2LsZuP/e8D\\n3ZL7yplFr8yDWEy7unOeI" +
	"4T5fk9A/gC1v28CgcEA7xTvc/wM/OOH2X+m+FWGV2OK\\nhof172MgH7KuKNyrWYCPFWRE4ry" +
	"le9dp7zCxg7O9Cf1KDMYxW6R50OXTqDH9WepS\\nY88Xp8FUv6gRm7GFR1oIppwvR52sXeh85" +
	"DFsmXEEfSCYvotgB2YtuUEQ6t1bhyBl\\nYup1p6YEohhma9tbpJ8Sl8t8jNLgzgLvHt+hnWq" +
	"8ADgYNqHFlxa+et9fT0/eo8sl\\n/vZnDTGqCtQbvt+RLAUUIDfMT1AzplAbkPmyCIDVAoHBA" +
	"JylitJl3ANpj9IW1p1V\\nWyVOrhjYfpDmEBHYykGJ88QGZKrIxZUB3f18f8k5c4vL/Z7ThBL" +
	"Nh/kCsdhv90ny\\n9zUUjBWyL28elZyM2UZFZUz4b6HErXYlOzf/BaH1G+ht22S5EDdVAc7P9" +
	"FsmS4V2\\ncGqwVUDJvXpXTjZ9Qb4DnGXd6m+UCVby1nMn+ff5Floay7Llgyi3F29q90EfEaF" +
	"j\\nz6ge41hkjHiN2YHJ4Y/qD622pu5UH1xWq8kNwWAvPdOGNwKBwGQSt0gtMKKJ+ssn\\ni5R" +
	"W+CyMFefJ8T+Kll+uQwgHlhGYpkul9M+TdxEnUcg08ZjXG0GbK+u86dwaNIMk\\nMDwu8IU6d" +
	"9IfVR965SZ0CyrIW6lKHmTe60tmzEdikSRTK8nubJqTklvvqvpt4qWw\\nfoxQq4ccIFQPH8B" +
	"U1Yuj+vsajnmRvA6dtjyb6E1THRQnN1bU3myZDkmP3U0o0AJO\\njW59WAvndLgMG4aCS1LNB" +
	"fAXVITAe3Qg5vwr7fKnhWJhogu9IQKBwEfRf8qTahSb\\n/oWNwhURUqKlvPVroml2SR4JefE" +
	"VvzlIQf+SeK+8fW92OzKeCNJCBgRfNuC+qB0h\\nXOw5CeB3xYYG5s6YFvwPMN6naXqP3QFCK" +
	"iL1t8L8yJBaeNyz4E3nLZhs4zQb1+j4\\nkFd+TUHWSimBkuuJW0WF3bkWIGPLWfYGWcHdO2u" +
	"SCl5eY0rxp9UWuYts2K12Di08\\n4LNz4G2v1KK8S/1sLdddTc0oA2Sy75odgo1TQMX4+/Nj3" +
	"o4aIbY1FA==\\n-----END RSA PRIVATE KEY-----"
