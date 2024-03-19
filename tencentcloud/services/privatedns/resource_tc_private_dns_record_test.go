package privatedns_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudPrivateDnsRecord_basic -v
func TestAccTencentCloudPrivateDnsRecord_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheck(t) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccPrivateDnsRecord,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_private_dns_record.example", "zone_id"),
					resource.TestCheckResourceAttr("tencentcloud_private_dns_record.example", "record_type", "A"),
					resource.TestCheckResourceAttr("tencentcloud_private_dns_record.example", "record_value", "192.168.1.2"),
					resource.TestCheckResourceAttr("tencentcloud_private_dns_record.example", "sub_domain", "www"),
					resource.TestCheckResourceAttrSet("tencentcloud_private_dns_record.example", "ttl"),
					resource.TestCheckResourceAttrSet("tencentcloud_private_dns_record.example", "weight"),
					resource.TestCheckResourceAttrSet("tencentcloud_private_dns_record.example", "mx"),
				),
			},
			{
				ResourceName:      "tencentcloud_private_dns_record.example",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccPrivateDnsRecordUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_private_dns_record.example", "zone_id"),
					resource.TestCheckResourceAttr("tencentcloud_private_dns_record.example", "record_type", "A"),
					resource.TestCheckResourceAttr("tencentcloud_private_dns_record.example", "record_value", "192.168.1.3"),
					resource.TestCheckResourceAttr("tencentcloud_private_dns_record.example", "sub_domain", "www"),
					resource.TestCheckResourceAttrSet("tencentcloud_private_dns_record.example", "ttl"),
					resource.TestCheckResourceAttrSet("tencentcloud_private_dns_record.example", "weight"),
					resource.TestCheckResourceAttrSet("tencentcloud_private_dns_record.example", "mx"),
				),
			},
		},
	})
}

const testAccPrivateDnsRecord = `
resource "tencentcloud_vpc" "vpc" {
  name       = "vpc-example"
  cidr_block = "10.0.0.0/16"
}

resource "tencentcloud_private_dns_zone" "example" {
  domain = "domain.com"
  remark = "remark."

  vpc_set {
    region      = "ap-guangzhou"
    uniq_vpc_id = tencentcloud_vpc.vpc.id
  }

  dns_forward_status   = "DISABLED"
  cname_speedup_status = "ENABLED"

  tags = {
    createdBy : "terraform"
  }
}

resource "tencentcloud_private_dns_record" "example" {
  zone_id      = tencentcloud_private_dns_zone.example.id
  record_type  = "A"
  record_value = "192.168.1.2"
  sub_domain   = "www"
  ttl          = 300
  weight       = 1
  mx           = 0
}
`

const testAccPrivateDnsRecordUpdate = `
resource "tencentcloud_vpc" "vpc" {
  name       = "vpc-example"
  cidr_block = "10.0.0.0/16"
}

resource "tencentcloud_private_dns_zone" "example" {
  domain = "domain.com"
  remark = "remark."

  vpc_set {
    region      = "ap-guangzhou"
    uniq_vpc_id = tencentcloud_vpc.vpc.id
  }

  dns_forward_status   = "DISABLED"
  cname_speedup_status = "ENABLED"

  tags = {
    createdBy : "terraform"
  }
}

resource "tencentcloud_private_dns_record" "example" {
  zone_id      = tencentcloud_private_dns_zone.example.id
  record_type  = "A"
  record_value = "192.168.1.3"
  sub_domain   = "www"
  ttl          = 300
  weight       = 1
  mx           = 0
}
`
