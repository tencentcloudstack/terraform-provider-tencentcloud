package pls_test

import (
	"regexp"
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudVpcEndPointResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVpcEndPoint,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_vpc_end_point.end_point", "id"),
					resource.TestCheckResourceAttr("tencentcloud_vpc_end_point.end_point", "end_point_name", "terraform_test"),
					resource.TestCheckResourceAttr("tencentcloud_vpc_end_point.end_point", "security_groups_ids.0", "sg-ghvp9djf"),
				),
			},

			{
				Config: testAccVpcEndPointUpdate,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("tencentcloud_vpc_end_point.end_point", "end_point_name", "terraform_test_for"),
					resource.TestCheckResourceAttr("tencentcloud_vpc_end_point.end_point", "security_groups_ids.0", "sg-3k7vtgf7"),
				),
			},
			{
				ResourceName:      "tencentcloud_vpc_end_point.end_point",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccVpcEndPoint = `

resource "tencentcloud_vpc_end_point" "end_point" {
  end_point_name       = "terraform_test"
  end_point_service_id = "vpcsvc-5y4yb85d"
  end_point_vip        = "10.0.0.58"
  security_groups_ids  = [
    "sg-ghvp9djf",
    "sg-if748odn",
    "sg-3k7vtgf7",
  ]
  subnet_id = "subnet-cpknsqgo"
  vpc_id    = "vpc-gmq0mxoj"
}

`

const testAccVpcEndPointUpdate = `

resource "tencentcloud_vpc_end_point" "end_point" {
  end_point_name       = "terraform_test_for"
  end_point_service_id = "vpcsvc-5y4yb85d"
  end_point_vip        = "10.0.0.58"
  security_groups_ids  = [
	"sg-3k7vtgf7",
    "sg-ghvp9djf",
    "sg-if748odn",
  ]
  subnet_id = "subnet-cpknsqgo"
  vpc_id    = "vpc-gmq0mxoj"
}

`

func TestAccTencentCloudVpcEndPoint_SecurityGroupId(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVpcEndPointSecurityGroupId,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_vpc_end_point.security_group_test", "id"),
					resource.TestCheckResourceAttr("tencentcloud_vpc_end_point.security_group_test", "security_group_id", "sg-ghvp9djf"),
				),
			},
			{
				Config: testAccVpcEndPointSecurityGroupIdUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("tencentcloud_vpc_end_point.security_group_test", "security_group_id", "sg-3k7vtgf7"),
				),
			},
		},
	})
}

const testAccVpcEndPointSecurityGroupId = `
resource "tencentcloud_vpc_end_point" "security_group_test" {
  end_point_name        = "terraform_sg_test"
  end_point_service_id  = "vpcsvc-5y4yb85d"
  security_group_id     = "sg-ghvp9djf"
  subnet_id            = "subnet-cpknsqgo"
  vpc_id               = "vpc-gmq0mxoj"
}
`

const testAccVpcEndPointSecurityGroupIdUpdate = `
resource "tencentcloud_vpc_end_point" "security_group_test" {
  end_point_name        = "terraform_sg_test"
  end_point_service_id  = "vpcsvc-5y4yb85d"
  security_group_id     = "sg-3k7vtgf7"
  subnet_id            = "subnet-cpknsqgo"
  vpc_id               = "vpc-gmq0mxoj"
}
`

func TestAccTencentCloudVpcEndPoint_Tags(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVpcEndPointTags,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_vpc_end_point.tags_test", "id"),
					resource.TestCheckResourceAttr("tencentcloud_vpc_end_point.tags_test", "tags.env", "test"),
					resource.TestCheckResourceAttr("tencentcloud_vpc_end_point.tags_test", "tags.project", "terraform"),
				),
			},
			{
				Config: testAccVpcEndPointTagsUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("tencentcloud_vpc_end_point.tags_test", "tags.env", "prod"),
					resource.TestCheckResourceAttr("tencentcloud_vpc_end_point.tags_test", "tags.project", "terraform"),
				),
			},
		},
	})
}

const testAccVpcEndPointTags = `
resource "tencentcloud_vpc_end_point" "tags_test" {
  end_point_name       = "terraform_tags_test"
  end_point_service_id = "vpcsvc-5y4yb85d"
  tags                = {
    env     = "test"
    project = "terraform"
  }
  subnet_id = "subnet-cpknsqgo"
  vpc_id    = "vpc-gmq0mxoj"
}
`

const testAccVpcEndPointTagsUpdate = `
resource "tencentcloud_vpc_end_point" "tags_test" {
  end_point_name       = "terraform_tags_test"
  end_point_service_id = "vpcsvc-5y4yb85d"
  tags                = {
    env     = "prod"
    project = "terraform"
  }
  subnet_id = "subnet-cpknsqgo"
  vpc_id    = "vpc-gmq0mxoj"
}
`

func TestAccTencentCloudVpcEndPoint_IpAddressType(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVpcEndPointIpAddressTypeDefault,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_vpc_end_point.ip_test", "id"),
					resource.TestCheckResourceAttr("tencentcloud_vpc_end_point.ip_test", "ip_address_type", "Ipv4"),
				),
			},
			{
				Config: testAccVpcEndPointIpAddressTypeIpv6,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("tencentcloud_vpc_end_point.ip_test", "ip_address_type", "Ipv6"),
				),
			},
			{
				Config:      testAccVpcEndPointIpAddressTypeInvalid,
				ExpectError: regexp.MustCompile(`Invalid ip_address_type value`),
			},
		},
	})
}

const testAccVpcEndPointIpAddressTypeDefault = `
resource "tencentcloud_vpc_end_point" "ip_test" {
  end_point_name       = "terraform_ip_test"
  end_point_service_id = "vpcsvc-5y4yb85d"
  subnet_id           = "subnet-cpknsqgo"
  vpc_id              = "vpc-gmq0mxoj"
}
`

const testAccVpcEndPointIpAddressTypeIpv6 = `
resource "tencentcloud_vpc_end_point" "ip_test" {
  end_point_name       = "terraform_ip_test"
  end_point_service_id = "vpcsvc-5y4yb85d"
  ip_address_type      = "Ipv6"
  subnet_id           = "subnet-cpknsqgo"
  vpc_id              = "vpc-gmq0mxoj"
}
`

const testAccVpcEndPointIpAddressTypeInvalid = `
resource "tencentcloud_vpc_end_point" "ip_test" {
  end_point_name       = "terraform_ip_test"
  end_point_service_id = "vpcsvc-5y4yb85d"
  ip_address_type      = "InvalidType"
  subnet_id           = "subnet-cpknsqgo"
  vpc_id              = "vpc-gmq0mxoj"
}
`

func TestAccTencentCloudVpcEndPoint_AllFields(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVpcEndPointAllFields,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_vpc_end_point.all_fields_test", "id"),
					resource.TestCheckResourceAttr("tencentcloud_vpc_end_point.all_fields_test", "security_group_id", "sg-ghvp9djf"),
					resource.TestCheckResourceAttr("tencentcloud_vpc_end_point.all_fields_test", "tags.env", "test"),
					resource.TestCheckResourceAttr("tencentcloud_vpc_end_point.all_fields_test", "ip_address_type", "Ipv4"),
				),
			},
		},
	})
}

const testAccVpcEndPointAllFields = `
resource "tencentcloud_vpc_end_point" "all_fields_test" {
  end_point_name       = "terraform_all_fields_test"
  end_point_service_id = "vpcsvc-5y4yb85d"
  security_group_id    = "sg-ghvp9djf"
  tags                = {
    env     = "test"
    project = "terraform"
  }
  ip_address_type      = "Ipv4"
  subnet_id           = "subnet-cpknsqgo"
  vpc_id              = "vpc-gmq0mxoj"
}
`

