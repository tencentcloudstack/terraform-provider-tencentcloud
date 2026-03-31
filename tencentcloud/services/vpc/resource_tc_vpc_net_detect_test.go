package vpc_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudVpcNetDetectResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVpcNetDetect,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_vpc_net_detect.net_detect", "id")),
			},
			{
				Config: testAccVpcNetDetectUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("tencentcloud_vpc_net_detect.net_detect", "net_detect_name", "terraform-for-test"),
				),
			},
			{
				ResourceName:      "tencentcloud_vpc_net_detect.net_detect",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccVpcNetDetect = `

resource "tencentcloud_vpc_net_detect" "net_detect" {
  net_detect_name       = "terrform-test"
  vpc_id                = "vpc-4owdpnwr"
  subnet_id             = "subnet-c1l35990"
  next_hop_destination  = "172.16.128.57"
  next_hop_type         = "NORMAL_CVM"
  detect_destination_ip = [
    "10.0.0.1",
    "10.0.0.2",
  ]
}

`

const testAccVpcNetDetectUpdate = `

resource "tencentcloud_vpc_net_detect" "net_detect" {
  net_detect_name       = "terraform-for-test"
  vpc_id                = "vpc-4owdpnwr"
  subnet_id             = "subnet-c1l35990"
  next_hop_destination  = "172.16.128.57"
  next_hop_type         = "NORMAL_CVM"
  detect_destination_ip = [
    "10.0.0.1",
    "10.0.0.2",
  ]
}

`

func TestAccTencentCloudVpcNetDetectResource_tags(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVpcNetDetectWithTags,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("tencentcloud_vpc_net_detect.net_detect_tags", "tags.%", "2"),
					resource.TestCheckResourceAttr("tencentcloud_vpc_net_detect.net_detect_tags", "tags.Environment", "test"),
					resource.TestCheckResourceAttr("tencentcloud_vpc_net_detect.net_detect_tags", "tags.Owner", "devops"),
				),
			},
			{
				Config: testAccVpcNetDetectWithTagsUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("tencentcloud_vpc_net_detect.net_detect_tags", "tags.%", "3"),
					resource.TestCheckResourceAttr("tencentcloud_vpc_net_detect.net_detect_tags", "tags.Environment", "production"),
					resource.TestCheckResourceAttr("tencentcloud_vpc_net_detect.net_detect_tags", "tags.Owner", "devops"),
					resource.TestCheckResourceAttr("tencentcloud_vpc_net_detect.net_detect_tags", "tags.Project", "network"),
				),
			},
			{
				Config: testAccVpcNetDetectWithNoTags,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("tencentcloud_vpc_net_detect.net_detect_tags", "tags.%", "0"),
				),
			},
			{
				ResourceName:      "tencentcloud_vpc_net_detect.net_detect_tags",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccVpcNetDetectWithTags = `
resource "tencentcloud_vpc_net_detect" "net_detect_tags" {
  net_detect_name       = "terraform-test-with-tags"
  vpc_id                = "vpc-4owdpnwr"
  subnet_id             = "subnet-c1l35990"
  next_hop_destination  = "172.16.128.57"
  next_hop_type         = "NORMAL_CVM"
  detect_destination_ip = [
    "10.0.0.1",
    "10.0.0.2",
  ]

  tags = {
    Environment = "test"
    Owner      = "devops"
  }
}
`

const testAccVpcNetDetectWithTagsUpdate = `
resource "tencentcloud_vpc_net_detect" "net_detect_tags" {
  net_detect_name       = "terraform-test-with-tags"
  vpc_id                = "vpc-4owdpnwr"
  subnet_id             = "subnet-c1l35990"
  next_hop_destination  = "172.16.128.57"
  next_hop_type         = "NORMAL_CVM"
  detect_destination_ip = [
    "10.0.0.1",
    "10.0.0.2",
  ]

  tags = {
    Environment = "production"
    Owner      = "devops"
    Project    = "network"
  }
}
`

const testAccVpcNetDetectWithNoTags = `
resource "tencentcloud_vpc_net_detect" "net_detect_tags" {
  net_detect_name       = "terraform-test-with-tags"
  vpc_id                = "vpc-4owdpnwr"
  subnet_id             = "subnet-c1l35990"
  next_hop_destination  = "172.16.128.57"
  next_hop_type         = "NORMAL_CVM"
  detect_destination_ip = [
    "10.0.0.1",
    "10.0.0.2",
  ]
}
`
