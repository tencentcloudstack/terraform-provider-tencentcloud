package tencentcloud

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccTencentCloudGaapCertificate_basic(t *testing.T) {
	id := new(string)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGaapCertificateDestroy(id),
		Steps: []resource.TestStep{
			{
				Config: testAccGaapCertificate(0, "\"test:tx2KGdo3zJg/.\"", "", ""),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGaapCertificateExists("tencentcloud_gaap_certificate.foo", id),
					resource.TestCheckResourceAttr("tencentcloud_gaap_certificate.foo", "type", "0"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_certificate.foo", "content", "test:tx2KGdo3zJg/."),
					resource.TestCheckNoResourceAttr("tencentcloud_gaap_certificate.foo", "key"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_certificate.foo", "create_time"),
					resource.TestCheckNoResourceAttr("tencentcloud_gaap_certificate.foo", "begin_time"),
					resource.TestCheckNoResourceAttr("tencentcloud_gaap_certificate.foo", "end_time"),
					resource.TestCheckNoResourceAttr("tencentcloud_gaap_certificate.foo", "issuer_cn"),
					resource.TestCheckNoResourceAttr("tencentcloud_gaap_certificate.foo", "subject_cn"),
				),
			},
		},
	})
}

func TestAccTencentCloudGaapCertificate_clientCA(t *testing.T) {
	id := new(string)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGaapCertificateDestroy(id),
		Steps: []resource.TestStep{
			{
				Config: testAccGaapCertificate(1, "<<EOF\n"+testAccGaapCertificateClientCA+"EOF", "", "<<EOF\n"+testAccGaapCertificateClientCAKey+"EOF"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGaapCertificateExists("tencentcloud_gaap_certificate.foo", id),
					resource.TestCheckResourceAttr("tencentcloud_gaap_certificate.foo", "type", "1"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_certificate.foo", "content", testAccGaapCertificateClientCA),
					resource.TestCheckResourceAttr("tencentcloud_gaap_certificate.foo", "key", testAccGaapCertificateClientCAKey),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_certificate.foo", "create_time"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_certificate.foo", "begin_time"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_certificate.foo", "end_time"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_certificate.foo", "issuer_cn"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_certificate.foo", "subject_cn"),
				),
			},
		},
	})
}

func TestAccTencentCloudGaapCertificate_ServerSSL(t *testing.T) {
	id := new(string)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGaapCertificateDestroy(id),
		Steps: []resource.TestStep{
			{
				Config: testAccGaapCertificate(2, "<<EOF\n"+testAccGaapCertificateServerCert+"EOF", "", "<<EOF\n"+testAccGaapCertificateServerKey+"EOF"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGaapCertificateExists("tencentcloud_gaap_certificate.foo", id),
					resource.TestCheckResourceAttr("tencentcloud_gaap_certificate.foo", "type", "2"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_certificate.foo", "content", testAccGaapCertificateServerCert),
					resource.TestCheckResourceAttr("tencentcloud_gaap_certificate.foo", "key", testAccGaapCertificateServerKey),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_certificate.foo", "create_time"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_certificate.foo", "begin_time"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_certificate.foo", "end_time"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_certificate.foo", "issuer_cn"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_certificate.foo", "subject_cn"),
				),
			},
		},
	})
}

func TestAccTencentCloudGaapCertificate_realserverCA(t *testing.T) {
	id := new(string)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGaapCertificateDestroy(id),
		Steps: []resource.TestStep{
			{
				Config: testAccGaapCertificate(3, "<<EOF\n"+testAccGaapCertificateClientCA+"EOF", "", "<<EOF\n"+testAccGaapCertificateClientCAKey+"EOF"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGaapCertificateExists("tencentcloud_gaap_certificate.foo", id),
					resource.TestCheckResourceAttr("tencentcloud_gaap_certificate.foo", "type", "3"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_certificate.foo", "content", testAccGaapCertificateClientCA),
					resource.TestCheckResourceAttr("tencentcloud_gaap_certificate.foo", "key", testAccGaapCertificateClientCAKey),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_certificate.foo", "create_time"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_certificate.foo", "begin_time"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_certificate.foo", "end_time"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_certificate.foo", "issuer_cn"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_certificate.foo", "subject_cn"),
				),
			},
		},
	})
}

func TestAccTencentCloudGaapCertificate_ProxySSL(t *testing.T) {
	id := new(string)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGaapCertificateDestroy(id),
		Steps: []resource.TestStep{
			{
				Config: testAccGaapCertificate(4, "<<EOF\n"+testAccGaapCertificateServerCert+"EOF", "", "<<EOF\n"+testAccGaapCertificateServerKey+"EOF"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGaapCertificateExists("tencentcloud_gaap_certificate.foo", id),
					resource.TestCheckResourceAttr("tencentcloud_gaap_certificate.foo", "type", "4"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_certificate.foo", "content", testAccGaapCertificateServerCert),
					resource.TestCheckResourceAttr("tencentcloud_gaap_certificate.foo", "key", testAccGaapCertificateServerKey),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_certificate.foo", "create_time"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_certificate.foo", "begin_time"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_certificate.foo", "end_time"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_certificate.foo", "issuer_cn"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_certificate.foo", "subject_cn"),
				),
			},
		},
	})
}

func TestAccTencentCloudGaapCertificate_updateName(t *testing.T) {
	id := new(string)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGaapCertificateDestroy(id),
		Steps: []resource.TestStep{
			{
				Config: testAccGaapCertificate(0, "\"test:tx2KGdo3zJg/.\"", "", ""),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGaapCertificateExists("tencentcloud_gaap_certificate.foo", id),
					resource.TestCheckResourceAttr("tencentcloud_gaap_certificate.foo", "type", "0"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_certificate.foo", "content", "test:tx2KGdo3zJg/."),
					resource.TestCheckNoResourceAttr("tencentcloud_gaap_certificate.foo", "key"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_certificate.foo", "create_time"),
					resource.TestCheckNoResourceAttr("tencentcloud_gaap_certificate.foo", "begin_time"),
					resource.TestCheckNoResourceAttr("tencentcloud_gaap_certificate.foo", "end_time"),
					resource.TestCheckNoResourceAttr("tencentcloud_gaap_certificate.foo", "issuer_cn"),
					resource.TestCheckNoResourceAttr("tencentcloud_gaap_certificate.foo", "subject_cn"),
				),
			},
			{
				Config: testAccGaapCertificate(0, "\"test:tx2KGdo3zJg/.\"", "new-name", ""),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("tencentcloud_gaap_certificate.foo"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_certificate.foo", "name", "new-name"),
				),
			},
		},
	})
}

func testAccCheckGaapCertificateExists(n string, id *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no certicicate id is set")
		}

		service := GaapService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}

		certificate, err := service.DescribeCertificateById(context.TODO(), rs.Primary.ID)
		if err != nil {
			return err
		}

		if certificate == nil {
			return fmt.Errorf("certificate not found: %s", rs.Primary.ID)
		}

		if certificate.CertificateId == nil {
			return errors.New("certificate id is nil")
		}

		*id = *certificate.CertificateId

		if *id == "" {
			return fmt.Errorf("certificate not found: %s", rs.Primary.ID)
		}

		return nil
	}
}

func testAccCheckGaapCertificateDestroy(id *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*TencentCloudClient).apiV3Conn
		service := GaapService{client: client}

		certificate, err := service.DescribeCertificateById(context.TODO(), *id)
		if err != nil {
			return err
		}

		if certificate != nil {
			return fmt.Errorf("certificate still exists")
		}

		return nil
	}
}

func testAccGaapCertificate(certificateType int, content, name, key string) string {
	const str = `
resource tencentcloud_gaap_certificate "foo" {
  type = %d
  content = %s
  %s
  %s
}
`

	if name != "" {
		name = "name = \"" + name + "\""
	}

	if key != "" {
		key = "key = " + key
	}

	return fmt.Sprintf(str, certificateType, content, name, key)
}

const testAccGaapCertificateClientCA = `
-----BEGIN CERTIFICATE-----
MIIEDjCCAnagAwIBAgIBATANBgkqhkiG9w0BAQsFADAoMQ0wCwYDVQQDEwR0ZXN0
MRcwFQYDVQQKEw50ZXJyYWZvcm0gdGVzdDAeFw0xOTA4MTMwMzA4MjBaFw0yOTA4
MTAwMzA4MjBaMCgxDTALBgNVBAMTBHRlc3QxFzAVBgNVBAoTDnRlcnJhZm9ybSB0
ZXN0MIIBojANBgkqhkiG9w0BAQEFAAOCAY8AMIIBigKCAYEA0k2vqg/GHtFP5P7r
dbzswfx1jSHeK9r4StV4mGOAoKyzvAJA5BvYbAHpSrL2ZAd6ShjHgRVU1qEroeFn
8fwTrAVQttMltBFABx7G4iN4Zf6EUXzhhFN6vVVbWaqhYhrdMoPvZxgGSA/4hG4W
GIr8MXZzXbKLoRoz4Bvq1Ymg5eO14KLJFSTahvIkG60egGN5pmi4czxWy2U7ycA5
Q5TuQBnF0rKQJW5XKIV3kr5YrzDdJK7up9E6Od4T5jz+qY97KAjIpWD/pTAsc7+6
fPBpY7NHT9Bw0fDmvsWO/PtswY4hW02n86b5eWA9sfKJGphhsBxgpuuhmxYHS6pA
B+C7IkyxcADNT5u9tEo2JGOj+/veXKrEhZin7inKsQLD0WOobcg1Rh/3NSWD7geF
TJBRnzgplaN7cK6c/utEAAnngS38q4DGBR/jHmkWjAeQPZj1eLLBk686HEEbKeU+
9yAVcPRhA9tuL7wMeSX32VunWZunoA/f8iuGZYJlZsNBqyJbAgMBAAGjQzBBMA8G
A1UdEwEB/wQFMAMBAf8wDwYDVR0PAQH/BAUDAweGADAdBgNVHQ4EFgQUKwfrmq79
1mY831S6UHARHtgYnlgwDQYJKoZIhvcNAQELBQADggGBAInM+aeaHoZdw9B9nAH2
HscEoOulF+RxnysSXTTRLd2VQph4+ynlfRZT4evLBBj/ppmqjp8F7/OcRiiZwSXl
namyP/UUINtHfgDM0kll/5Za0aYzMhrORNw+3ythIv2yPJX8t4LmsG1L4PMO8ZU8
N0K9XyKRaL/tq6rw1gQM152OmNgTzfAQoKYxrvbftOZz4J0ZACctuBmwtp5upKvJ
36aQ4wJLUzOt69mnW+AaL5EPA37mwtzdnzTTxd3SBfOYXjsflc3l2raljJznnqU2
ySynjb6L3D3L/pObL1Uu7nQBy8CazJBsBsVFK/pr61vcllm8lG7vOhHOUSFUeezq
FWukAolm9/cagmD6IhNishM3Uzng+UYyCC8uQq3Z7FGqJpXSI79wZYjudnCLPVCg
OIfJHQeJFLryn6GxiSYmYs6dgUJiiTV+I/2Y5X7ZFdb5FC1J/WmvoCv6yO7NiirY
BSgfV0lp5CuV8SfiSClpYfrM28NbNgxveUqET642BJOPLQ==
-----END CERTIFICATE-----
`

const testAccGaapCertificateClientCAKey = `
Public Key Info:
	Public Key Algorithm: RSA
	Key Security Level: High (3072 bits)

modulus:
	00:d2:4d:af:aa:0f:c6:1e:d1:4f:e4:fe:eb:75:bc:ec
	c1:fc:75:8d:21:de:2b:da:f8:4a:d5:78:98:63:80:a0
	ac:b3:bc:02:40:e4:1b:d8:6c:01:e9:4a:b2:f6:64:07
	7a:4a:18:c7:81:15:54:d6:a1:2b:a1:e1:67:f1:fc:13
	ac:05:50:b6:d3:25:b4:11:40:07:1e:c6:e2:23:78:65
	fe:84:51:7c:e1:84:53:7a:bd:55:5b:59:aa:a1:62:1a
	dd:32:83:ef:67:18:06:48:0f:f8:84:6e:16:18:8a:fc
	31:76:73:5d:b2:8b:a1:1a:33:e0:1b:ea:d5:89:a0:e5
	e3:b5:e0:a2:c9:15:24:da:86:f2:24:1b:ad:1e:80:63
	79:a6:68:b8:73:3c:56:cb:65:3b:c9:c0:39:43:94:ee
	40:19:c5:d2:b2:90:25:6e:57:28:85:77:92:be:58:af
	30:dd:24:ae:ee:a7:d1:3a:39:de:13:e6:3c:fe:a9:8f
	7b:28:08:c8:a5:60:ff:a5:30:2c:73:bf:ba:7c:f0:69
	63:b3:47:4f:d0:70:d1:f0:e6:be:c5:8e:fc:fb:6c:c1
	8e:21:5b:4d:a7:f3:a6:f9:79:60:3d:b1:f2:89:1a:98
	61:b0:1c:60:a6:eb:a1:9b:16:07:4b:aa:40:07:e0:bb
	22:4c:b1:70:00:cd:4f:9b:bd:b4:4a:36:24:63:a3:fb
	fb:de:5c:aa:c4:85:98:a7:ee:29:ca:b1:02:c3:d1:63
	a8:6d:c8:35:46:1f:f7:35:25:83:ee:07:85:4c:90:51
	9f:38:29:95:a3:7b:70:ae:9c:fe:eb:44:00:09:e7:81
	2d:fc:ab:80:c6:05:1f:e3:1e:69:16:8c:07:90:3d:98
	f5:78:b2:c1:93:af:3a:1c:41:1b:29:e5:3e:f7:20:15
	70:f4:61:03:db:6e:2f:bc:0c:79:25:f7:d9:5b:a7:59
	9b:a7:a0:0f:df:f2:2b:86:65:82:65:66:c3:41:ab:22
	5b:

public exponent:
	01:00:01:

private exponent:
	5a:f7:77:90:9c:1a:1a:a2:77:68:9a:4b:c7:35:dd:43
	5b:ac:8d:4b:a5:0a:5b:41:23:3d:8b:58:7f:51:d8:2e
	5b:e0:6b:29:1e:82:5c:ee:fb:34:aa:37:17:14:d5:97
	34:0d:db:de:1e:18:00:6e:de:ac:bb:0f:77:40:8e:51
	ce:4a:c7:8a:35:b8:d9:ed:54:27:1f:e8:19:67:ae:d6
	94:ed:9a:93:01:e6:0a:25:73:92:7c:0a:ae:9b:fc:fa
	c9:2b:00:97:1b:71:3c:22:8c:60:dc:2d:7a:98:43:d6
	31:62:5d:99:29:84:9a:0c:ee:57:a5:10:90:e3:a4:0d
	07:53:0f:96:e3:2a:79:cb:fd:59:59:0f:5d:2e:33:d9
	1e:fe:15:2a:e3:62:b7:c0:26:48:72:79:52:9e:4d:20
	35:05:b8:c9:bc:48:34:9e:46:cc:d4:98:08:f5:db:71
	cb:76:5e:a4:a4:ba:7f:f1:1f:fa:83:90:c6:a7:19:84
	67:5c:b9:99:2c:d2:24:4b:40:4f:0f:da:21:fa:1d:77
	7f:bf:19:2e:30:45:bb:0e:70:81:66:a8:9e:45:55:27
	83:a2:cd:56:22:2e:fd:37:49:09:39:e4:c0:24:eb:fb
	c6:e4:44:bb:5a:54:fb:fc:7d:f8:5f:90:4f:a7:e4:f8
	f8:0a:a9:71:ec:4e:7f:fc:3c:5a:5c:e2:50:73:a2:1e
	35:ef:be:1a:b4:b3:f9:ea:c8:0d:6c:eb:87:fe:20:2c
	90:2b:fc:6a:04:e2:a1:35:f2:d3:d8:a1:df:c5:00:fb
	9b:d8:f7:32:2b:25:b9:18:e3:32:d0:bb:92:c1:64:e1
	44:64:e5:13:54:a9:14:c5:16:81:84:cd:0a:4d:20:9a
	7c:ef:3a:13:ca:d6:fc:3c:e8:5b:b8:53:ac:13:4e:45
	cc:67:4e:5f:6a:c2:9b:98:9c:9d:5a:e7:db:49:7b:69
	5e:d9:e1:ae:59:87:e7:02:5b:ac:31:78:55:a2:24:61
	

prime1:
	00:e3:c8:9d:d9:46:6d:29:3a:22:be:46:01:a2:a7:9d
	18:70:6d:6b:37:6a:fd:70:d4:16:dc:58:4b:70:80:d7
	d5:df:70:7a:ea:ed:23:0f:13:a9:f1:a1:f7:2c:11:8a
	96:d3:0f:9f:33:31:0d:e4:8a:06:44:f0:9f:51:92:4a
	08:4d:ba:96:34:21:13:30:5b:14:95:29:b8:ea:55:ee
	a1:1f:5a:42:1b:df:1e:06:8e:d4:cd:02:0c:58:07:65
	03:69:5d:73:a3:01:dc:75:86:a3:b0:bb:a8:5d:61:66
	bf:c6:8a:c7:1f:ca:17:1a:01:a0:fa:ad:11:16:4b:c3
	65:75:6d:b4:a3:19:74:54:0b:9a:4e:ca:f7:98:15:63
	a2:cb:c6:0f:e9:cb:98:91:f9:5f:45:28:4b:e2:9a:a7
	5b:ab:66:e6:08:ba:2d:af:67:40:59:8f:08:e5:0a:f1
	0b:71:85:b4:d3:70:17:cf:7f:1f:4c:77:bd:6e:ba:c8
	57:

prime2:
	00:ec:5a:be:a1:21:e0:07:74:19:94:f6:bb:2f:22:42
	a7:1c:11:7a:22:18:2e:57:cb:76:63:75:ac:42:a8:b1
	b8:e6:8a:e8:c3:bf:a2:02:d4:e9:b1:dd:80:3d:57:42
	04:ec:26:0b:78:6e:e5:b9:ed:d7:96:02:56:0d:d4:4e
	58:0b:93:07:ec:31:20:75:42:13:0f:17:36:de:11:ab
	15:ad:f1:24:54:3e:2b:51:5f:50:70:09:d6:8d:83:31
	71:d4:1d:c6:2b:a6:19:89:9c:be:8e:9c:e8:05:31:26
	78:29:83:c2:75:14:00:5b:13:8b:34:36:23:c4:6c:75
	d6:e4:be:33:15:5b:e6:22:f9:b5:8b:13:1c:0a:d9:53
	a2:18:4d:1d:23:09:23:1e:ac:1c:e9:ae:40:cf:e5:d1
	d7:ec:c8:50:8c:f0:bf:6b:13:9c:39:d4:28:cc:92:04
	55:74:ac:80:9d:38:cc:df:96:ed:f7:1c:db:7f:1c:c3
	9d:

coefficient:
	5f:d6:59:d1:aa:a7:45:e3:27:f9:68:4b:a1:4d:ab:e6
	2c:31:d9:69:60:43:ad:3e:ed:0a:3e:10:20:41:e9:44
	fc:9c:c4:a9:d7:80:bc:ad:48:e3:20:d8:6a:49:b9:9a
	75:40:a3:8c:db:2f:1f:86:a0:33:04:04:4c:d8:1b:22
	77:68:ae:e5:a7:40:dc:11:50:92:46:85:0d:f5:72:1b
	f5:29:6b:72:c3:8c:d4:e5:95:71:a3:db:b3:e0:f3:5c
	95:c2:33:df:61:1a:d3:2b:64:14:31:1f:22:80:78:3a
	0f:c9:d2:c3:53:c1:65:fb:a8:b4:b7:a9:2f:21:46:be
	e1:a3:22:b9:ce:1d:34:93:ef:6d:64:97:3d:1f:29:b7
	db:07:57:02:2d:9e:60:42:00:13:ca:68:25:aa:ea:e3
	86:1c:5a:5c:8c:13:87:42:4e:b0:99:cf:fe:c9:35:4e
	60:16:ce:ee:5d:c4:94:87:86:7c:bc:9c:1a:6a:de:74
	

exp1:
	2f:b4:0b:02:be:d2:0b:a9:46:2f:6c:ff:d8:ad:9b:a3
	cb:9f:ce:ad:6b:75:aa:54:70:79:32:f0:91:9e:1a:15
	8b:56:c6:17:3f:14:71:8a:df:b3:60:05:20:b0:87:c2
	b0:6e:fc:1b:3f:71:b6:64:05:8e:18:8a:75:0d:da:fd
	44:32:08:54:e0:7c:61:4b:21:d3:5f:4a:7f:a1:01:79
	b2:d4:37:36:19:12:f3:b2:a3:f4:4f:32:80:99:03:d9
	a4:0e:53:32:57:28:71:60:82:15:78:27:79:6e:f5:92
	c2:24:bc:30:f7:24:c1:68:87:eb:17:a5:95:72:c6:78
	10:10:aa:9c:e1:ae:d0:0f:22:00:d0:56:eb:fd:c6:c1
	08:45:f1:7d:38:0b:3a:be:eb:e9:d2:d7:99:9c:63:e2
	39:d2:31:e8:af:fc:f3:57:ff:4d:02:3c:8d:a3:2d:fa
	04:8b:48:e4:62:fc:49:93:48:dd:bd:d1:e9:47:aa:bf
	

exp2:
	00:aa:fa:26:26:c8:2b:99:18:ba:9c:d1:33:ad:b0:1a
	09:6b:3b:95:5f:6e:a0:af:b0:26:cf:62:43:9f:e1:0a
	d7:98:26:ea:a5:18:4e:a9:bf:bd:e2:86:3a:8b:a5:40
	c3:f0:d4:c0:bd:79:73:c3:e1:b3:6f:f7:f6:aa:80:67
	c0:37:77:40:66:f3:4e:e8:a4:48:c7:44:e2:d3:18:72
	eb:f6:ed:97:7b:3d:91:f6:86:7d:d8:de:b5:0f:9f:d5
	37:e3:db:3a:0d:3f:55:ff:ff:9c:b5:eb:f9:c8:2f:bb
	05:77:fa:b1:c4:22:18:f3:c9:9a:8c:c5:91:70:39:89
	b1:4f:bd:eb:94:4a:eb:9b:67:8a:95:b5:d8:36:5c:ea
	ac:3a:ea:25:e6:bf:68:61:7a:de:5b:89:4b:a7:59:38
	48:fc:28:18:65:2c:8f:3c:75:6e:31:43:31:d4:72:d8
	a1:bd:0a:40:eb:fb:37:f1:92:ec:48:33:de:f1:00:2c
	dd:


Public Key PIN:
	pin-sha256:HDQpHefsN5rr+X/cLabz5OF3OlA1vTQF92x7f1VT86g=
Public Key ID:
	sha256:1c34291de7ec379aebf97fdc2da6f3e4e1773a5035bd3405f76c7b7f5553f3a8
	sha1:2b07eb9aaefdd6663cdf54ba5070111ed8189e58

-----BEGIN RSA PRIVATE KEY-----
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
-----END RSA PRIVATE KEY-----
`

const testAccGaapCertificateServerCert = `
-----BEGIN CERTIFICATE-----
MIIERzCCAq+gAwIBAgIBAjANBgkqhkiG9w0BAQsFADAoMQ0wCwYDVQQDEwR0ZXN0
MRcwFQYDVQQKEw50ZXJyYWZvcm0gdGVzdDAeFw0xOTA4MTMwMzE5MzlaFw0yOTA4
MTAwMzE5MzlaMC4xEzARBgNVBAMTCnNlcnZlciBzc2wxFzAVBgNVBAoTDnRlcnJh
Zm9ybS10ZXN0MIIBojANBgkqhkiG9w0BAQEFAAOCAY8AMIIBigKCAYEA1Ryp+DKK
SNFKZsPtwfR+jzOnQ8YFieIKYgakV688d8YgpolenbmeEPrzT87tunFD7G9f6ALG
ND8rj7npj0AowxhOL/h/v1D9u0UsIaj5i2GWJrqNAhGLaxWiEB/hy5WOiwxDrGei
gQqJkFM52Ep7G1Yx7PHJmKFGwN9FhIsFi1cNZfVRopZuCe/RMPNusNVZaIi+qcEf
fsE1cmfmuSlG3Ap0RKOIyR0ajDEzqZn9/0R7VwWCF97qy8TNYk94K/1tq3zyhVzR
Z83xOSfrTqEfb3so3AU2jyKgYdwr/FZS72VCHS8IslgnqJW4izIXZqgIKmHaRZtM
N4jUloi6l/6lktt6Lsgh9xECecxziSJtPMaog88aC8HnMqJJ3kScGCL36GYG+Kaw
5PnDlWXBaeiDe8z/eWK9+Rr2M+rhTNxosAVGfDJyxAXyiX49LQ0v7f9qzwc/0JiD
bvsUv1cm6OgpoEMP9SXqqBdwGqeKbD2/2jlP48xlYP6l1SoJG3GgZ8dbAgMBAAGj
djB0MAwGA1UdEwEB/wQCMAAwEwYDVR0lBAwwCgYIKwYBBQUHAwEwDwYDVR0PAQH/
BAUDAweAADAdBgNVHQ4EFgQULwWKBQNLL9s3cb3tTnyPVg+mpCMwHwYDVR0jBBgw
FoAUKwfrmq791mY831S6UHARHtgYnlgwDQYJKoZIhvcNAQELBQADggGBAMo5RglS
AHdPgaicWJvmvjjexjF/42b7Rz4pPfMjYw6uYO8He/f4UZWv5CZLrbEe7MywaK3y
0OsfH8AhyN29pv2x8g9wbmq7omZIOZ0oCAGduEXs/A/qY/hFaCohdkz/IN8qi6JW
VXreGli3SrpcHFchSwHTyJEXgkutcGAsOvdsOuVSmplOyrkLHc8uUe8SG4j8kGyg
EzaszFjHkR7g1dVyDVUedc588mjkQxYeAamJgfkgIhljWKMa2XzkVMcVfQHfNpM1
n+bu8SmqRt9Wma2bMijKRG/Blm756LoI+skY+WRZmlDnq8zj95TT0vceGP0FUWh5
hKyiocABmpQs9OK9HMi8vgSWISP+fYgkm/bKtKup2NbZBoO5/VL2vCEPInYzUhBO
jCbLMjNjtM5KriCaR7wDARgHiG0gBEPOEW1PIjZ9UOH+LtIxbNZ4eEIIINLHnBHf
L+doVeZtS/gJc4G4Adr5HYuaS9ZxJ0W2uy0eQlOHzjyxR6Mf/rpnilJlcQ==
-----END CERTIFICATE-----
`

const testAccGaapCertificateServerKey = `
Public Key Info:
	Public Key Algorithm: RSA
	Key Security Level: High (3072 bits)

modulus:
	00:d5:1c:a9:f8:32:8a:48:d1:4a:66:c3:ed:c1:f4:7e
	8f:33:a7:43:c6:05:89:e2:0a:62:06:a4:57:af:3c:77
	c6:20:a6:89:5e:9d:b9:9e:10:fa:f3:4f:ce:ed:ba:71
	43:ec:6f:5f:e8:02:c6:34:3f:2b:8f:b9:e9:8f:40:28
	c3:18:4e:2f:f8:7f:bf:50:fd:bb:45:2c:21:a8:f9:8b
	61:96:26:ba:8d:02:11:8b:6b:15:a2:10:1f:e1:cb:95
	8e:8b:0c:43:ac:67:a2:81:0a:89:90:53:39:d8:4a:7b
	1b:56:31:ec:f1:c9:98:a1:46:c0:df:45:84:8b:05:8b
	57:0d:65:f5:51:a2:96:6e:09:ef:d1:30:f3:6e:b0:d5
	59:68:88:be:a9:c1:1f:7e:c1:35:72:67:e6:b9:29:46
	dc:0a:74:44:a3:88:c9:1d:1a:8c:31:33:a9:99:fd:ff
	44:7b:57:05:82:17:de:ea:cb:c4:cd:62:4f:78:2b:fd
	6d:ab:7c:f2:85:5c:d1:67:cd:f1:39:27:eb:4e:a1:1f
	6f:7b:28:dc:05:36:8f:22:a0:61:dc:2b:fc:56:52:ef
	65:42:1d:2f:08:b2:58:27:a8:95:b8:8b:32:17:66:a8
	08:2a:61:da:45:9b:4c:37:88:d4:96:88:ba:97:fe:a5
	92:db:7a:2e:c8:21:f7:11:02:79:cc:73:89:22:6d:3c
	c6:a8:83:cf:1a:0b:c1:e7:32:a2:49:de:44:9c:18:22
	f7:e8:66:06:f8:a6:b0:e4:f9:c3:95:65:c1:69:e8:83
	7b:cc:ff:79:62:bd:f9:1a:f6:33:ea:e1:4c:dc:68:b0
	05:46:7c:32:72:c4:05:f2:89:7e:3d:2d:0d:2f:ed:ff
	6a:cf:07:3f:d0:98:83:6e:fb:14:bf:57:26:e8:e8:29
	a0:43:0f:f5:25:ea:a8:17:70:1a:a7:8a:6c:3d:bf:da
	39:4f:e3:cc:65:60:fe:a5:d5:2a:09:1b:71:a0:67:c7
	5b:

public exponent:
	01:00:01:

private exponent:
	00:b1:56:d0:fa:00:d4:a2:13:c7:5e:0c:dc:e4:f1:97
	ff:82:74:46:29:9a:a2:4a:bf:69:23:2d:ce:e9:bb:df
	cf:b7:8b:dd:f4:26:3c:38:14:d9:3f:6f:c2:3a:81:53
	8f:ba:48:53:fe:b5:90:4a:19:e7:1e:0b:0f:18:6d:c3
	7d:d5:d3:fa:87:47:86:e4:d6:bf:e7:a7:f9:ba:ab:2e
	19:5e:e1:8b:8b:9b:95:0d:f7:66:61:1e:19:e9:c3:88
	08:be:1c:ce:93:c1:09:b1:68:1b:61:46:60:74:64:46
	5d:51:34:ea:7f:a9:ca:a1:2a:47:85:84:4b:ef:84:05
	97:c3:46:7d:06:19:ce:24:73:90:64:fb:df:16:d5:80
	34:8e:90:7c:58:b6:a4:86:ce:30:b3:ab:52:8b:f2:95
	4c:b6:46:5a:77:db:73:c0:0c:3f:6d:12:18:a8:54:7c
	ff:77:c3:ca:89:9f:63:98:ef:48:2d:c1:09:70:6e:ea
	cb:bb:78:91:42:8a:22:3e:21:ef:a5:bf:16:ee:66:45
	e5:f0:26:6a:85:8e:e1:69:62:ac:05:00:a6:44:ba:c8
	ac:10:00:97:f5:51:65:7f:9a:1f:7b:99:9d:02:d4:87
	50:ce:74:06:51:67:fa:fb:90:e4:33:79:f2:a8:61:ee
	45:1d:87:ca:22:5b:ac:e7:32:38:f8:2c:fd:55:92:1e
	3d:60:1e:7c:4b:fd:28:ff:e5:b1:02:6a:aa:22:f7:ae
	a8:36:90:7b:a6:f7:29:05:14:3a:21:da:36:05:f9:b0
	9d:f7:fb:10:75:d7:2c:21:32:95:e7:f7:17:be:09:cb
	66:fe:f1:69:71:df:a4:5e:3f:fc:67:6c:37:65:b8:51
	c6:22:38:fb:07:ce:89:54:50:43:71:44:3d:c3:51:5a
	bd:e5:c7:87:b2:ea:7b:64:0f:5d:34:9c:a1:52:b3:ce
	06:86:ba:7a:05:80:48:b3:6c:1b:79:74:9b:49:f2:30
	c1:

prime1:
	00:e4:31:46:59:3d:24:f7:31:d9:22:26:af:c0:3e:f5
	c1:6d:be:ba:d3:9e:3f:b9:2c:43:a0:d0:47:09:e4:35
	63:19:a4:33:82:af:f9:76:3c:11:c2:cb:34:f9:a6:ab
	dd:ab:64:5a:6b:9c:c1:2a:52:89:64:7e:b5:a7:f0:4d
	29:13:a4:cf:17:f4:f2:0d:a0:6e:b9:5d:95:41:10:df
	ae:f3:7a:13:49:21:66:73:2a:b7:e2:8d:7c:c2:34:e5
	3f:bd:78:ca:fc:64:c5:1c:3a:66:7a:12:53:96:bd:b0
	c3:7a:0c:ec:5e:55:c0:c3:3f:7f:25:72:f4:e2:19:94
	9d:65:15:be:c8:82:20:57:12:97:b2:a8:4d:3d:e0:8f
	e2:1f:d0:c8:49:aa:f4:34:fa:91:d1:d1:cc:98:bc:3d
	8b:b1:9b:8f:fd:ef:03:dd:92:fb:ca:99:45:af:cc:83
	58:4c:bb:ba:73:9e:23:84:f9:7e:4f:40:fe:00:b5:bf
	6f:

prime2:
	00:ef:14:ef:73:fc:0c:fc:e3:87:d9:7f:a6:f8:55:86
	57:63:8a:86:87:f5:ef:63:20:1f:b2:ae:28:dc:ab:59
	80:8f:15:64:44:e2:bc:a5:7b:d7:69:ef:30:b1:83:b3
	bd:09:fd:4a:0c:c6:31:5b:a4:79:d0:e5:d3:a8:31:fd
	59:ea:52:63:cf:17:a7:c1:54:bf:a8:11:9b:b1:85:47
	5a:08:a6:9c:2f:47:9d:ac:5d:e8:7c:e4:31:6c:99:71
	04:7d:20:98:be:8b:60:07:66:2d:b9:41:10:ea:dd:5b
	87:20:65:62:ea:75:a7:a6:04:a2:18:66:6b:db:5b:a4
	9f:12:97:cb:7c:8c:d2:e0:ce:02:ef:1e:df:a1:9d:6a
	bc:00:38:18:36:a1:c5:97:16:be:7a:df:5f:4f:4f:de
	a3:cb:25:fe:f6:67:0d:31:aa:0a:d4:1b:be:df:91:2c
	05:14:20:37:cc:4f:50:33:a6:50:1b:90:f9:b2:08:80
	d5:

coefficient:
	47:d1:7f:ca:93:6a:14:9b:fe:85:8d:c2:15:11:52:a2
	a5:bc:f5:6b:a2:69:76:49:1e:09:79:f1:15:bf:39:48
	41:ff:92:78:af:bc:7d:6f:76:3b:32:9e:08:d2:42:06
	04:5f:36:e0:be:a8:1d:21:5c:ec:39:09:e0:77:c5:86
	06:e6:ce:98:16:fc:0f:30:de:a7:69:7a:8f:dd:01:42
	2a:22:f5:b7:c2:fc:c8:90:5a:78:dc:b3:e0:4d:e7:2d
	98:6c:e3:34:1b:d7:e8:f8:90:57:7e:4d:41:d6:4a:29
	81:92:eb:89:5b:45:85:dd:b9:16:20:63:cb:59:f6:06
	59:c1:dd:3b:6b:92:0a:5e:5e:63:4a:f1:a7:d5:16:b9
	8b:6c:d8:ad:76:0e:2d:3c:e0:b3:73:e0:6d:af:d4:a2
	bc:4b:fd:6c:2d:d7:5d:4d:cd:28:03:64:b2:ef:9a:1d
	82:8d:53:40:c5:f8:fb:f3:63:de:8e:1a:21:b6:35:14
	

exp1:
	00:9c:a5:8a:d2:65:dc:03:69:8f:d2:16:d6:9d:55:5b
	25:4e:ae:18:d8:7e:90:e6:10:11:d8:ca:41:89:f3:c4
	06:64:aa:c8:c5:95:01:dd:fd:7c:7f:c9:39:73:8b:cb
	fd:9e:d3:84:12:cd:87:f9:02:b1:d8:6f:f7:49:f2:f7
	35:14:8c:15:b2:2f:6f:1e:95:9c:8c:d9:46:45:65:4c
	f8:6f:a1:c4:ad:76:25:3b:37:ff:05:a1:f5:1b:e8:6d
	db:64:b9:10:37:55:01:ce:cf:f4:5b:26:4b:85:76:70
	6a:b0:55:40:c9:bd:7a:57:4e:36:7d:41:be:03:9c:65
	dd:ea:6f:94:09:56:f2:d6:73:27:f9:f7:f9:16:5a:1a
	cb:b2:e5:83:28:b7:17:6f:6a:f7:41:1f:11:a1:63:cf
	a8:1e:e3:58:64:8c:78:8d:d9:81:c9:e1:8f:ea:0f:ad
	b6:a6:ee:54:1f:5c:56:ab:c9:0d:c1:60:2f:3d:d3:86
	37:

exp2:
	64:12:b7:48:2d:30:a2:89:fa:cb:27:8b:94:56:f8:2c
	8c:15:e7:c9:f1:3f:8a:96:5f:ae:43:08:07:96:11:98
	a6:4b:a5:f4:cf:93:77:11:27:51:c8:34:f1:98:d7:1b
	41:9b:2b:eb:bc:e9:dc:1a:34:83:24:30:3c:2e:f0:85
	3a:77:d2:1f:55:1f:7a:e5:26:74:0b:2a:c8:5b:a9:4a
	1e:64:de:eb:4b:66:cc:47:62:91:24:53:2b:c9:ee:6c
	9a:93:92:5b:ef:aa:fa:6d:e2:a5:b0:7e:8c:50:ab:87
	1c:20:54:0f:1f:c0:54:d5:8b:a3:fa:fb:1a:8e:79:91
	bc:0e:9d:b6:3c:9b:e8:4d:53:1d:14:27:37:56:d4:de
	6c:99:0e:49:8f:dd:4d:28:d0:02:4e:8d:6e:7d:58:0b
	e7:74:b8:0c:1b:86:82:4b:52:cd:05:f0:17:54:84:c0
	7b:74:20:e6:fc:2b:ed:f2:a7:85:62:61:a2:0b:bd:21
	


Public Key PIN:
	pin-sha256:t5OXXC5gYqMNtUMsTqRs3A3vhfK2BiXVOgYzIEYv7Y8=
Public Key ID:
	sha256:b793975c2e6062a30db5432c4ea46cdc0def85f2b60625d53a063320462fed8f
	sha1:2f058a05034b2fdb3771bded4e7c8f560fa6a423

-----BEGIN RSA PRIVATE KEY-----
MIIG5AIBAAKCAYEA1Ryp+DKKSNFKZsPtwfR+jzOnQ8YFieIKYgakV688d8Ygpole
nbmeEPrzT87tunFD7G9f6ALGND8rj7npj0AowxhOL/h/v1D9u0UsIaj5i2GWJrqN
AhGLaxWiEB/hy5WOiwxDrGeigQqJkFM52Ep7G1Yx7PHJmKFGwN9FhIsFi1cNZfVR
opZuCe/RMPNusNVZaIi+qcEffsE1cmfmuSlG3Ap0RKOIyR0ajDEzqZn9/0R7VwWC
F97qy8TNYk94K/1tq3zyhVzRZ83xOSfrTqEfb3so3AU2jyKgYdwr/FZS72VCHS8I
slgnqJW4izIXZqgIKmHaRZtMN4jUloi6l/6lktt6Lsgh9xECecxziSJtPMaog88a
C8HnMqJJ3kScGCL36GYG+Kaw5PnDlWXBaeiDe8z/eWK9+Rr2M+rhTNxosAVGfDJy
xAXyiX49LQ0v7f9qzwc/0JiDbvsUv1cm6OgpoEMP9SXqqBdwGqeKbD2/2jlP48xl
YP6l1SoJG3GgZ8dbAgMBAAECggGBALFW0PoA1KITx14M3OTxl/+CdEYpmqJKv2kj
Lc7pu9/Pt4vd9CY8OBTZP2/COoFTj7pIU/61kEoZ5x4LDxhtw33V0/qHR4bk1r/n
p/m6qy4ZXuGLi5uVDfdmYR4Z6cOICL4czpPBCbFoG2FGYHRkRl1RNOp/qcqhKkeF
hEvvhAWXw0Z9BhnOJHOQZPvfFtWANI6QfFi2pIbOMLOrUovylUy2Rlp323PADD9t
EhioVHz/d8PKiZ9jmO9ILcEJcG7qy7t4kUKKIj4h76W/Fu5mReXwJmqFjuFpYqwF
AKZEusisEACX9VFlf5ofe5mdAtSHUM50BlFn+vuQ5DN58qhh7kUdh8oiW6znMjj4
LP1Vkh49YB58S/0o/+WxAmqqIveuqDaQe6b3KQUUOiHaNgX5sJ33+xB11ywhMpXn
9xe+Cctm/vFpcd+kXj/8Z2w3ZbhRxiI4+wfOiVRQQ3FEPcNRWr3lx4ey6ntkD100
nKFSs84Ghrp6BYBIs2wbeXSbSfIwwQKBwQDkMUZZPST3MdkiJq/APvXBbb66054/
uSxDoNBHCeQ1YxmkM4Kv+XY8EcLLNPmmq92rZFprnMEqUolkfrWn8E0pE6TPF/Ty
DaBuuV2VQRDfrvN6E0khZnMqt+KNfMI05T+9eMr8ZMUcOmZ6ElOWvbDDegzsXlXA
wz9/JXL04hmUnWUVvsiCIFcSl7KoTT3gj+If0MhJqvQ0+pHR0cyYvD2LsZuP/e8D
3ZL7yplFr8yDWEy7unOeI4T5fk9A/gC1v28CgcEA7xTvc/wM/OOH2X+m+FWGV2OK
hof172MgH7KuKNyrWYCPFWRE4ryle9dp7zCxg7O9Cf1KDMYxW6R50OXTqDH9WepS
Y88Xp8FUv6gRm7GFR1oIppwvR52sXeh85DFsmXEEfSCYvotgB2YtuUEQ6t1bhyBl
Yup1p6YEohhma9tbpJ8Sl8t8jNLgzgLvHt+hnWq8ADgYNqHFlxa+et9fT0/eo8sl
/vZnDTGqCtQbvt+RLAUUIDfMT1AzplAbkPmyCIDVAoHBAJylitJl3ANpj9IW1p1V
WyVOrhjYfpDmEBHYykGJ88QGZKrIxZUB3f18f8k5c4vL/Z7ThBLNh/kCsdhv90ny
9zUUjBWyL28elZyM2UZFZUz4b6HErXYlOzf/BaH1G+ht22S5EDdVAc7P9FsmS4V2
cGqwVUDJvXpXTjZ9Qb4DnGXd6m+UCVby1nMn+ff5Floay7Llgyi3F29q90EfEaFj
z6ge41hkjHiN2YHJ4Y/qD622pu5UH1xWq8kNwWAvPdOGNwKBwGQSt0gtMKKJ+ssn
i5RW+CyMFefJ8T+Kll+uQwgHlhGYpkul9M+TdxEnUcg08ZjXG0GbK+u86dwaNIMk
MDwu8IU6d9IfVR965SZ0CyrIW6lKHmTe60tmzEdikSRTK8nubJqTklvvqvpt4qWw
foxQq4ccIFQPH8BU1Yuj+vsajnmRvA6dtjyb6E1THRQnN1bU3myZDkmP3U0o0AJO
jW59WAvndLgMG4aCS1LNBfAXVITAe3Qg5vwr7fKnhWJhogu9IQKBwEfRf8qTahSb
/oWNwhURUqKlvPVroml2SR4JefEVvzlIQf+SeK+8fW92OzKeCNJCBgRfNuC+qB0h
XOw5CeB3xYYG5s6YFvwPMN6naXqP3QFCKiL1t8L8yJBaeNyz4E3nLZhs4zQb1+j4
kFd+TUHWSimBkuuJW0WF3bkWIGPLWfYGWcHdO2uSCl5eY0rxp9UWuYts2K12Di08
4LNz4G2v1KK8S/1sLdddTc0oA2Sy75odgo1TQMX4+/Nj3o4aIbY1FA==
-----END RSA PRIVATE KEY-----
`
