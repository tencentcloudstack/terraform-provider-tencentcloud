package privatedns_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudPrivateDnsInboundEndpointResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccPrivateDnsInboundEndpoint,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_private_dns_inbound_endpoint.example", "id")),
			},
			{
				Config: testAccPrivateDnsInboundEndpointUpdate,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_private_dns_inbound_endpoint.example", "id")),
			},
		},
	})
}

const testAccPrivateDnsInboundEndpoint = `
resource "tencentcloud_private_dns_inbound_endpoint" "example" {
  endpoint_name   = "tf-example"
  endpoint_region = "ap-guangzhou"
  endpoint_vpc    = "vpc-i5yyodl9"
  subnet_ip {
    subnet_id  = "subnet-hhi88a58"
    subnet_vip = "10.0.30.2"
  }

  subnet_ip {
    subnet_id  = "subnet-5rrirqyc"
    subnet_vip = "10.0.0.11"
  }

  subnet_ip {
    subnet_id  = "subnet-60ut6n10"
  }
}
`

const testAccPrivateDnsInboundEndpointUpdate = `
resource "tencentcloud_private_dns_inbound_endpoint" "example" {
  endpoint_name   = "tf-example-update"
  endpoint_region = "ap-guangzhou"
  endpoint_vpc    = "vpc-i5yyodl9"
  subnet_ip {
    subnet_id  = "subnet-hhi88a58"
    subnet_vip = "10.0.30.2"
  }

  subnet_ip {
    subnet_id  = "subnet-5rrirqyc"
    subnet_vip = "10.0.0.11"
  }

  subnet_ip {
    subnet_id  = "subnet-60ut6n10"
  }
}
`
