package vpc_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudElasticPublicIpv6AttachmentResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccElasticPublicIpv6Attachment,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_elastic_public_ipv6_attachment.elastic_public_ipv6_attachment", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_elastic_public_ipv6_attachment.elastic_public_ipv6_attachment", "ipv6_address_id"),
					resource.TestCheckResourceAttr("tencentcloud_elastic_public_ipv6_attachment.elastic_public_ipv6_attachment", "network_interface_id", "eni-k6pjc7nr"),
					resource.TestCheckResourceAttr("tencentcloud_elastic_public_ipv6_attachment.elastic_public_ipv6_attachment", "private_ipv6_address", "fd76:3600:700:1900:0:9dbc:f353:871"),
				),
			},
			{
				ResourceName:            "tencentcloud_elastic_public_ipv6_attachment.elastic_public_ipv6_attachment",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"keep_bind_with_eni"},
			},
		},
	})
}

const testAccElasticPublicIpv6Attachment = `
resource "tencentcloud_elastic_public_ipv6" "elastic_public_ipv6" {
    address_name = "test"
    internet_max_bandwidth_out = 1
}
resource "tencentcloud_elastic_public_ipv6_attachment" "elastic_public_ipv6_attachment" {
  ipv6_address_id = tencentcloud_elastic_public_ipv6.elastic_public_ipv6.id
  network_interface_id = "eni-k6pjc7nr"
  private_ipv6_address = "fd76:3600:700:1900:0:9dbc:f353:871"
}
`
