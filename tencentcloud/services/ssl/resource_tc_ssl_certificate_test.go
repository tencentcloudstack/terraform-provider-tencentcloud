package ssl_test

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svcssl "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/ssl"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"

	sdkError "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	ssl "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ssl/v20191205"
)

func init() {
	resource.AddTestSweepers("tencentcloud_ssl_certificate", &resource.Sweeper{
		Name: "tencentcloud_ssl_certificate",
		F:    testSweepSslCertificate,
	})
}

func testSweepSslCertificate(region string) error {
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	sharedClient, err := tcacctest.SharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("getting tencentcloud client error: %s", err.Error())
	}
	client := sharedClient.(tccommon.ProviderMeta)

	sslService := svcssl.NewSSLService(client.GetAPIV3Conn())
	describeRequest := ssl.NewDescribeCertificatesRequest()
	instances, err := sslService.DescribeCertificates(ctx, describeRequest)
	if err != nil {
		return fmt.Errorf("get instance list error: %s", err.Error())
	}

	for _, v := range instances {

		instanceId := *v.CertificateId
		instanceName := *v.Alias
		now := time.Now()
		createTime := tccommon.StringToTime(*v.CertBeginTime)
		interval := now.Sub(createTime).Minutes()

		if strings.HasPrefix(instanceName, tcacctest.KeepResource) || strings.HasPrefix(instanceName, tcacctest.DefaultResource) {
			continue
		}

		if tccommon.NeedProtect == 1 && int64(interval) < 30 {
			continue
		}

		request := ssl.NewDeleteCertificateRequest()
		request.CertificateId = helper.String(instanceId)
		if _, err = sslService.DeleteCertificate(ctx, request); err != nil {
			log.Printf("[ERROR] sweep instance %s error: %s", instanceId, err.Error())
		}
	}

	return nil
}

func TestAccTencentCloudSslCertificate_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheck(t) },
		Providers:    tcacctest.AccProviders,
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

func TestAccTencentCloudSslCertificate_tags(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheck(t) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckSslCertificateDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSslCertificateWithTags(`{
					tagKey1="tagValue1"
				}`),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSslCertificateExists("tencentcloud_ssl_certificate.test-ssl-certificate-tag"),
					resource.TestCheckResourceAttr("tencentcloud_ssl_certificate.test-ssl-certificate-tag", "tags.tagKey1", "tagValue1"),
				),
			},
			{
				Config: testAccSslCertificateWithTags(`{
					tagKey1="tagValue1"
					tagKey2="tagValue2"
				}`),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSslCertificateExists("tencentcloud_ssl_certificate.test-ssl-certificate-tag"),
					resource.TestCheckResourceAttr("tencentcloud_ssl_certificate.test-ssl-certificate-tag", "tags.tagKey1", "tagValue1"),
					resource.TestCheckResourceAttr("tencentcloud_ssl_certificate.test-ssl-certificate-tag", "tags.tagKey2", "tagValue2"),
				),
			},
		},
	})
}

func TestAccTencentCloudSslCertificate_svr(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheck(t) },
		Providers:    tcacctest.AccProviders,
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
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	sslService := svcssl.NewSSLService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
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
		outErr = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			describeResponse, inErr = sslService.DescribeCertificateDetail(ctx, describeRequest)
			if inErr != nil {
				if sdkErr, ok := inErr.(*sdkError.TencentCloudSDKError); ok {
					if sdkErr.Code == svcssl.CertificateNotFound {
						return nil
					}
				}
				return tccommon.RetryError(inErr)
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
		logId := tccommon.GetLogId(tccommon.ContextNil)
		ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		sslService := svcssl.NewSSLService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())

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
		outErr = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			describeResponse, inErr = sslService.DescribeCertificateDetail(ctx, describeRequest)
			if inErr != nil {
				if sdkErr, ok := inErr.(*sdkError.TencentCloudSDKError); ok {
					if sdkErr.Code == svcssl.CertificateNotFound {
						return nil
					}
				}
				return tccommon.RetryError(inErr)
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

func testAccSslCertificateWithTags(tags string) string {
	const str = `
resource "tencentcloud_ssl_certificate" "test-ssl-certificate-tag" {
  type = "CA"
  cert = "%s"
  name = "test-ssl-certificate-tag"
  tags =%s
}`
	return fmt.Sprintf(str, testAccSslCertificateCA, tags)
}

const testAccSslCertificateCA = "-----BEGIN CERTIFICATE-----\\n" +
	"MIIDzjCCAragAwIBAgIUN9h/fEqxX/xxyLcAVUT67S/zN3kwDQYJKoZIhvcNAQEL\\n" +
	"BQAweDELMAkGA1UEBhMCQ04xCzAJBgNVBAgMAlNYMQ0wCwYDVQQHDARYSUFOMQsw\\n" +
	"CQYDVQQKDAJUQzELMAkGA1UECwwCTU0xETAPBgNVBAMMCHRlc3QtaWFjMSAwHgYJ\\n" +
	"KoZIhvcNAQkBFhExMjU0Njg4NTU5QHFxLmNvbTAeFw0yMzA4MjkwODA0NTNaFw0z\\n" +
	"MzA4MjYwODA0NTNaMHgxCzAJBgNVBAYTAkNOMQswCQYDVQQIDAJTWDENMAsGA1UE\\n" +
	"BwwEWElBTjELMAkGA1UECgwCVEMxCzAJBgNVBAsMAk1NMREwDwYDVQQDDAh0ZXN0\\n" +
	"LWlhYzEgMB4GCSqGSIb3DQEJARYRMTI1NDY4ODU1OUBxcS5jb20wggEiMA0GCSqG\\n" +
	"SIb3DQEBAQUAA4IBDwAwggEKAoIBAQDbXEnXjfHZUtyTlF2BDZxfuxkJ3CMIBbrT\\n" +
	"Tp2YXoj/jMKRCUgAg5AGD5/uuJ9yF9HMyUr0VCDrf7JJLxk5hWbpbUahATjotGd7\\n" +
	"6z/5e9IM+IPTP//wa/7I3tL3fTq9nj1cmVxpmiC5wGwFSAssxkLVQDiQVKHW+6Oi\\n" +
	"Qde3ENHQj+IwMq4bohAThKxCXqeDqtfAzJLPfMmLG6HTEwh10isbS1BDhXeYSoGN\\n" +
	"HtkzTg07B0p76DnQEztJjKNY1/zVFwDSHpIbdsRAKsEsHEV0hzrjrBSCSwqfbuuG\\n" +
	"LzWEVKWvmh+R6h2xhw9BBws7kYcw3cXUN+PiMB9aBSnRupFyASNpAgMBAAGjUDBO\\n" +
	"MB0GA1UdDgQWBBR+qkSE6nATEINfR3xEZ9QKQfZZ+DAfBgNVHSMEGDAWgBR+qkSE\\n" +
	"6nATEINfR3xEZ9QKQfZZ+DAMBgNVHRMEBTADAQH/MA0GCSqGSIb3DQEBCwUAA4IB\\n" +
	"AQAbgRo+4gcTGfXaXLv5KgkhPOyzpyaoAbfWtk8PkAlZb+8j6bxb4qEsaVwXrH8U\\n" +
	"rRIe26YsoRCP3dYvJITaOO8KT0VZnO/KZeIdoHnguhjcynx60zt4hLQ+83N/JMJr\\n" +
	"lLX+JIg2nhwnS97aXdZL6sdcraNKXD63Gyjhp4VMeSDi9juIYfheB6STbM/os0Aj\\n" +
	"Xn12QxOg/nbAJah3aqYlZ38Js6lJ2vkd9AYptdPOw8SyeavrSiQg14iihLydh0mN\\n" +
	"GfiXelsRMOgB46B/Xvk/ir3wes36rkrN+ixUiVcLYDME77O3Z4Mre2PzueOIWNzR\\n" +
	"Axxb5inrStx6colrfzQThuwF\\n-----END CERTIFICATE-----"

const testAccSslCertificateKey = "-----BEGIN PRIVATE KEY-----\\n" +
	"MIIEvQIBADANBgkqhkiG9w0BAQEFAASCBKcwggSjAgEAAoIBAQDbXEnXjfHZUtyT\\n" +
	"lF2BDZxfuxkJ3CMIBbrTTp2YXoj/jMKRCUgAg5AGD5/uuJ9yF9HMyUr0VCDrf7JJ\\n" +
	"Lxk5hWbpbUahATjotGd76z/5e9IM+IPTP//wa/7I3tL3fTq9nj1cmVxpmiC5wGwF\\n" +
	"SAssxkLVQDiQVKHW+6OiQde3ENHQj+IwMq4bohAThKxCXqeDqtfAzJLPfMmLG6HT\\n" +
	"Ewh10isbS1BDhXeYSoGNHtkzTg07B0p76DnQEztJjKNY1/zVFwDSHpIbdsRAKsEs\\n" +
	"HEV0hzrjrBSCSwqfbuuGLzWEVKWvmh+R6h2xhw9BBws7kYcw3cXUN+PiMB9aBSnR\\n" +
	"upFyASNpAgMBAAECggEAIpBnqDnbBlHapLxngVq6LZFnEBkqQezZM8t65JPcxVuS\\n" +
	"GtVaDY6tZm8W1cAsi4c6TSjYkSgiackcuBBeSqR9A0HvM5ZkN7KZbbqzQWXjwpxz\\n" +
	"9RjsBJ+XrWIC3vFSDKe+5nTZzV/2UR6DRs/DxwHUbRKp9wAG4j+TWJFEYrmZPeHr\\n" +
	"RpkqR59RjODCLjjqB1tLS4dejvNNJlljPDjmKgWyZWNRxtlXF/+75IK3CykEmEpu\\n" +
	"ttSNjydb+VC99iKvd+XSvj3Km7FWECy6XH2LDGxGYkD5SD+1JBcJdRMq7QYyWmbz\\n" +
	"+SbbVclxdIp6QJD7l10A4dkSJT22V9omRn2Yn8UwmwKBgQD62FOqFMNjsIOX3gcz\\n" +
	"OEIEqkSytG/ndEDmKX7rQARzxKCb2n6NQtGPKSY6yjVQultUyuWlopQZzzTMg2Lr\\n" +
	"6iJfbpZLXRs1TclYVE0XV/Alk055qvTTPBnNzo9AHFT4io6j0MEcI1JbIy9DEDb7\\n" +
	"jxsFsCBgj43UYjq+0NqU1JndGwKBgQDf3lMJ6FAfi+nWR5UKTt9FLb46A2D4aNLV\\n" +
	"zSZU8P0ExRQo6NHQuqwBeDKXGbQH2JDpv4427DzGjdkT+5s6htgOR6ReKJ+2uRg6\\n" +
	"EfU+vzuykoKto6k6+aur1sVnn8ZG2zX9bnaGzSRsLLgK7eeGjbpEHgtKmBz0X3C4\\n" +
	"9DXL6sRdywKBgQDnonevqTi8h7UcuhRgAeVEtY52jxR+4OVFJLBkwFrcJIhDI0KV\\n" +
	"Y0xsLI124F7XSx8nb60chMLKCoMxD2p7e1t+UHpM4Y9Ma6YwALing7bom9xtkaY+\\n" +
	"oVMar1Gs2/zC/f+12gFY4G0eZ6EvBnwfVAiZ+ggL4sQPiR3CMs6FfMUQXQKBgAT6\\n" +
	"+znzMyUghblAqm4qRwlQ9TRxMs0T9+zNvZaSLe7XO5WVaGWOYZk+xVFbPwgVp1Or\\n" +
	"8UwDgW6hZTzukguBSHk42s1FdhgokgNott4Ifxl/7OxUAcXQHCOciZO+mDinU9Ip\\n" +
	"jPV+xtqpPAbyN/5kVMqDKJkmPS6qmOVkeXXp6Sh1AoGATVLo1V61XmwyWHn57CeK\\n" +
	"ofwhPYgQumIpUad51hm2ONflFa6dvkO3n0XDNwQ4PN6nJ0Qja+QyDrZ9zRxcd3YJ\\n" +
	"NqZ/TEDT5H8L1M/ovdgt2H2M8AL4QZS4eBa0JgOfQUokgDmfpqZYagVBqasO3YFN\\n" +
	"a31BbeeDG10TcTUM1JxNFYg=\\n-----END PRIVATE KEY-----"
