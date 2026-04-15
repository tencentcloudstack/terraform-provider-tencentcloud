package teo_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudTeoDnsRecordV9Resource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTeoDnsRecordV9_basic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_teo_dns_record_v9.test", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_teo_dns_record_v9.test", "zone_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_teo_dns_record_v9.test", "name"),
					resource.TestCheckResourceAttr("tencentcloud_teo_dns_record_v9.test", "type", "A"),
					resource.TestCheckResourceAttrSet("tencentcloud_teo_dns_record_v9.test", "content"),
					resource.TestCheckResourceAttr("tencentcloud_teo_dns_record_v9.test", "ttl", "300"),
					resource.TestCheckResourceAttrSet("tencentcloud_teo_dns_record_v9.test", "status"),
					resource.TestCheckResourceAttrSet("tencentcloud_teo_dns_record_v9.test", "created_on"),
					resource.TestCheckResourceAttrSet("tencentcloud_teo_dns_record_v9.test", "modified_on"),
				),
			},
			{
				ResourceName:      "tencentcloud_teo_dns_record_v9.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccTeoDnsRecordV9_update,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_teo_dns_record_v9.test", "id"),
					resource.TestCheckResourceAttr("tencentcloud_teo_dns_record_v9.test", "content", "5.6.7.8"),
					resource.TestCheckResourceAttr("tencentcloud_teo_dns_record_v9.test", "ttl", "600"),
				),
			},
			{
				Config: testAccTeoDnsRecordV9_delete,
			},
		},
	})
}

func TestAccTencentCloudTeoDnsRecordV9Resource_cname(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTeoDnsRecordV9_cname,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_teo_dns_record_v9.test", "id"),
					resource.TestCheckResourceAttr("tencentcloud_teo_dns_record_v9.test", "type", "CNAME"),
					resource.TestCheckResourceAttr("tencentcloud_teo_dns_record_v9.test", "location", "Asia"),
					resource.TestCheckResourceAttr("tencentcloud_teo_dns_record_v9.test", "weight", "10"),
				),
			},
			{
				ResourceName:      "tencentcloud_teo_dns_record_v9.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccTencentCloudTeoDnsRecordV9Resource_mx(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTeoDnsRecordV9_mx,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_teo_dns_record_v9.test", "id"),
					resource.TestCheckResourceAttr("tencentcloud_teo_dns_record_v9.test", "type", "MX"),
					resource.TestCheckResourceAttr("tencentcloud_teo_dns_record_v9.test", "priority", "10"),
				),
			},
			{
				ResourceName:      "tencentcloud_teo_dns_record_v9.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccTencentCloudTeoDnsRecordV9Resource_locationWeight(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTeoDnsRecordV9_locationWeight,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_teo_dns_record_v9.test", "id"),
					resource.TestCheckResourceAttr("tencentcloud_teo_dns_record_v9.test", "type", "A"),
					resource.TestCheckResourceAttr("tencentcloud_teo_dns_record_v9.test", "location", "Asia"),
					resource.TestCheckResourceAttr("tencentcloud_teo_dns_record_v9.test", "weight", "20"),
				),
			},
			{
				ResourceName:      "tencentcloud_teo_dns_record_v9.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccTeoDnsRecordV9_basic = `
resource "tencentcloud_teo_zone" "test" {
  zone_name = "tf-test-dns-record-v9.example.com"
  area      = "mainland"
  type      = "full"
}

resource "tencentcloud_teo_dns_record_v9" "test" {
  zone_id = tencentcloud_teo_zone.test.id
  name     = "www"
  type     = "A"
  content  = "1.2.3.4"
  ttl      = 300
}
`

const testAccTeoDnsRecordV9_update = `
resource "tencentcloud_teo_zone" "test" {
  zone_name = "tf-test-dns-record-v9.example.com"
  area      = "mainland"
  type      = "full"
}

resource "tencentcloud_teo_dns_record_v9" "test" {
  zone_id = tencentcloud_teo_zone.test.id
  name     = "www"
  type     = "A"
  content  = "5.6.7.8"
  ttl      = 600
}
`

const testAccTeoDnsRecordV9_delete = `
resource "tencentcloud_teo_zone" "test" {
  zone_name = "tf-test-dns-record-v9.example.com"
  area      = "mainland"
  type      = "full"
}
`

const testAccTeoDnsRecordV9_cname = `
resource "tencentcloud_teo_zone" "test" {
  zone_name = "tf-test-dns-record-v9-cname.example.com"
  area      = "mainland"
  type      = "full"
}

resource "tencentcloud_teo_dns_record_v9" "test" {
  zone_id  = tencentcloud_teo_zone.test.id
  name     = "www"
  type     = "CNAME"
  content  = "alias.example.com"
  location = "Asia"
  weight   = 10
  ttl      = 600
}
`

const testAccTeoDnsRecordV9_mx = `
resource "tencentcloud_teo_zone" "test" {
  zone_name = "tf-test-dns-record-v9-mx.example.com"
  area      = "mainland"
  type      = "full"
}

resource "tencentcloud_teo_dns_record_v9" "test" {
  zone_id  = tencentcloud_teo_zone.test.id
  name     = "@"
  type     = "MX"
  content  = "mail.example.com"
  priority = 10
}
`

const testAccTeoDnsRecordV9_locationWeight = `
resource "tencentcloud_teo_zone" "test" {
  zone_name = "tf-test-dns-record-v9-location.example.com"
  area      = "mainland"
  type      = "full"
}

resource "tencentcloud_teo_dns_record_v9" "test" {
  zone_id  = tencentcloud_teo_zone.test.id
  name     = "www"
  type     = "A"
  content  = "10.0.0.1"
  location = "Asia"
  weight   = 20
  ttl      = 600
}
`
