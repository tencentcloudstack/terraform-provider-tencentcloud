package pls_test

import (
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

func TestAccTencentCloudVpcEndPointResource_withSecurityGroupId(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVpcEndPointWithSecurityGroupId,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_vpc_end_point.end_point", "id"),
					resource.TestCheckResourceAttr("tencentcloud_vpc_end_point.end_point", "end_point_name", "terraform_test"),
					resource.TestCheckResourceAttr("tencentcloud_vpc_end_point.end_point", "security_group_id", "sg-ghvp9djf"),
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

func TestAccTencentCloudVpcEndPointResource_withTags(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVpcEndPointWithTags,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_vpc_end_point.end_point", "id"),
					resource.TestCheckResourceAttr("tencentcloud_vpc_end_point.end_point", "end_point_name", "terraform_test"),
					resource.TestCheckResourceAttr("tencentcloud_vpc_end_point.end_point", "tags.0.key", "env"),
					resource.TestCheckResourceAttr("tencentcloud_vpc_end_point.end_point", "tags.0.value", "test"),
				),
			},
			{
				Config: testAccVpcEndPointWithTagsUpdate,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("tencentcloud_vpc_end_point.end_point", "tags.0.key", "env"),
					resource.TestCheckResourceAttr("tencentcloud_vpc_end_point.end_point", "tags.0.value", "prod"),
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

func TestAccTencentCloudVpcEndPointResource_withIpAddressType(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVpcEndPointWithIpAddressType,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_vpc_end_point.end_point", "id"),
					resource.TestCheckResourceAttr("tencentcloud_vpc_end_point.end_point", "end_point_name", "terraform_test"),
					resource.TestCheckResourceAttr("tencentcloud_vpc_end_point.end_point", "ip_address_type", "Ipv4"),
				),
			},
			{
				Config: testAccVpcEndPointWithIpAddressTypeUpdate,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("tencentcloud_vpc_end_point.end_point", "ip_address_type", "Ipv6"),
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

func TestAccTencentCloudVpcEndPointResource_withoutNewFields(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVpcEndPointWithoutNewFields,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_vpc_end_point.end_point", "id"),
					resource.TestCheckResourceAttr("tencentcloud_vpc_end_point.end_point", "end_point_name", "terraform_test"),
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

const testAccVpcEndPointWithSecurityGroupId = `

resource "tencentcloud_vpc_end_point" "end_point" {
  end_point_name       = "terraform_test"
  end_point_service_id = "vpcsvc-5y4yb85d"
  end_point_vip        = "10.0.0.58"
  security_group_id    = "sg-ghvp9djf"
  subnet_id            = "subnet-cpknsqgo"
  vpc_id               = "vpc-gmq0mxoj"
}

`

const testAccVpcEndPointWithTags = `

resource "tencentcloud_vpc_end_point" "end_point" {
  end_point_name       = "terraform_test"
  end_point_service_id = "vpcsvc-5y4yb85d"
  end_point_vip        = "10.0.0.58"
  subnet_id            = "subnet-cpknsqgo"
  vpc_id               = "vpc-gmq0mxoj"

  tags {
    key   = "env"
    value = "test"
  }
  tags {
    key   = "owner"
    value = "terraform"
  }
}

`

const testAccVpcEndPointWithTagsUpdate = `

resource "tencentcloud_vpc_end_point" "end_point" {
  end_point_name       = "terraform_test"
  end_point_service_id = "vpcsvc-5y4yb85d"
  end_point_vip        = "10.0.0.58"
  subnet_id            = "subnet-cpknsqgo"
  vpc_id               = "vpc-gmq0mxoj"

  tags {
    key   = "env"
    value = "prod"
  }
  tags {
    key   = "owner"
    value = "terraform"
  }
}

`

const testAccVpcEndPointWithIpAddressType = `

resource "tencentcloud_vpc_end_point" "end_point" {
  end_point_name       = "terraform_test"
  end_point_service_id = "vpcsvc-5y4yb85d"
  end_point_vip        = "10.0.0.58"
  subnet_id            = "subnet-cpknsqgo"
  vpc_id               = "vpc-gmq0mxoj"
  ip_address_type       = "Ipv4"
}

`

const testAccVpcEndPointWithIpAddressTypeUpdate = `

resource "tencentcloud_vpc_end_point" "end_point" {
  end_point_name       = "terraform_test"
  end_point_service_id = "vpcsvc-5y4yb85d"
  end_point_vip        = "10.0.0.58"
  subnet_id            = "subnet-cpknsqgo"
  vpc_id               = "vpc-gmq0mxoj"
  ip_address_type       = "Ipv6"
}

`

const testAccVpcEndPointWithoutNewFields = `

resource "tencentcloud_vpc_end_point" "end_point" {
  end_point_name       = "terraform_test"
  end_point_service_id = "vpcsvc-5y4yb85d"
  end_point_vip        = "10.0.0.58"
  subnet_id            = "subnet-cpknsqgo"
  vpc_id               = "vpc-gmq0mxoj"
}

`
