package vpc_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudVpcRoutePolicyEntriesResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVpcRoutePolicyEntries,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_vpc_route_policy_entries.example", "id"),
				),
			},
			{
				Config: testAccVpcRoutePolicyEntriesUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_vpc_route_policy_entries.example", "id"),
				),
			},
			{
				ResourceName:      "tencentcloud_vpc_route_policy_entries.example",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccVpcRoutePolicyEntries = `
resource "tencentcloud_vpc_route_policy" "example" {
  route_policy_name        = "tf-example"
  route_policy_description = "remark."
}

resource "tencentcloud_vpc_route_policy_entries" "example" {
  route_policy_id = tencentcloud_vpc_route_policy.example.id
  route_policy_entry_set {
    cidr_block   = "10.10.0.0/16"
    route_type   = "ANY"
    gateway_type = "VPN"
    gateway_id   = "vpngw-may3cb0m"
    action       = "ACCEPT"
  }

  route_policy_entry_set {
    cidr_block   = "172.16.0.0/16"
    description  = "remark"
    route_type   = "ANY"
    gateway_type = "EIP"
    priority     = 10
    action       = "ACCEPT"
  }

  route_policy_entry_set {
    cidr_block   = "192.168.0.0/16"
    description  = "remark"
    route_type   = "ANY"
    gateway_type = "HAVIP"
    gateway_id   = "havip-r3ar5p86"
    priority     = 1
    action       = "ACCEPT"
  }
}
`

const testAccVpcRoutePolicyEntriesUpdate = `
resource "tencentcloud_vpc_route_policy" "example" {
  route_policy_name        = "tf-example"
  route_policy_description = "remark."
}

resource "tencentcloud_vpc_route_policy_entries" "example" {
  route_policy_id = tencentcloud_vpc_route_policy.example.id
  route_policy_entry_set {
    cidr_block   = "10.10.0.0/16"
    route_type   = "ANY"
    gateway_type = "VPN"
    gateway_id   = "vpngw-may3cb0m"
    action       = "ACCEPT"
  }
}
`
