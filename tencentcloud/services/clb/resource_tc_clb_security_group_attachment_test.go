package clb_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudClbSecurityGroupAttachmentResource_basic -v
func TestAccTencentCloudClbSecurityGroupAttachmentResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccClbSecurityGroupAttachment,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_clb_security_group_attachment.example", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_clb_security_group_attachment.example", "security_group"),
					resource.TestCheckResourceAttrSet("tencentcloud_clb_security_group_attachment.example", "load_balancer_ids.#"),
				),
			},
			{
				ResourceName:      "tencentcloud_clb_security_group_attachment.example",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccClbSecurityGroupAttachment = `
# create vpc
resource "tencentcloud_vpc" "vpc" {
  name       = "vpc"
  cidr_block = "10.0.0.0/16"
}

# create vpc subnet
resource "tencentcloud_subnet" "subnet" {
  name              = "subnet"
  vpc_id            = tencentcloud_vpc.vpc.id
  availability_zone = "ap-guangzhou-6"
  cidr_block        = "10.0.20.0/28"
  is_multicast      = false
}

# create clb
resource "tencentcloud_clb_instance" "example" {
  network_type = "INTERNAL"
  clb_name     = "tf-example"
  project_id   = 0
  vpc_id       = tencentcloud_vpc.vpc.id
  subnet_id    = tencentcloud_subnet.subnet.id

  tags = {
    "example" = "test"
  }
}

# attachment
resource "tencentcloud_clb_security_group_attachment" "example" {
  security_group    = "sg-5275dorp"
  load_balancer_ids = [tencentcloud_clb_instance.example.id]
}
`
