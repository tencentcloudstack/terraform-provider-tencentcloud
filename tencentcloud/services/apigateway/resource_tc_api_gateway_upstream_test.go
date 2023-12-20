package apigateway_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudApigatewayUpstreamResource_basic -v
func TestAccTencentCloudApigatewayUpstreamResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccApigatewayUpstreamCvm,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_api_gateway_upstream.example_cvm", "id"),
				),
			},
			{
				ResourceName:      "tencentcloud_api_gateway_upstream.example_cvm",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccApigatewayUpstreamCvmUpdateCvm,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_api_gateway_upstream.example_cvm", "id"),
				),
			},
		},
	})
}

const testAccApigatewayUpstreamCvm = `
data "tencentcloud_availability_zones_by_product" "zones" {
  product = "cvm"
}

data "tencentcloud_images" "images" {
  image_type       = ["PUBLIC_IMAGE"]
  image_name_regex = "Final"
}

data "tencentcloud_instance_types" "instance_types" {
  filter {
    name   = "instance-family"
    values = ["S5"]
  }

  cpu_core_count   = 2
  exclude_sold_out = true
}

resource "tencentcloud_vpc" "vpc" {
  name       = "vpc-example"
  cidr_block = "10.0.0.0/16"
}

resource "tencentcloud_subnet" "subnet" {
  availability_zone = data.tencentcloud_availability_zones_by_product.zones.zones.3.name
  name              = "subnet-example"
  vpc_id            = tencentcloud_vpc.vpc.id
  cidr_block        = "10.0.0.0/16"
  is_multicast      = false
}

resource "tencentcloud_instance" "example" {
  instance_name     = "tf_example"
  availability_zone = data.tencentcloud_availability_zones_by_product.zones.zones.3.name
  image_id          = data.tencentcloud_images.images.images.0.image_id
  instance_type     = data.tencentcloud_instance_types.instance_types.instance_types.0.instance_type
  system_disk_type  = "CLOUD_PREMIUM"
  system_disk_size  = 50
  hostname          = "terraform"
  project_id        = 0
  vpc_id            = tencentcloud_vpc.vpc.id
  subnet_id         = tencentcloud_subnet.subnet.id

  data_disks {
    data_disk_type = "CLOUD_PREMIUM"
    data_disk_size = 50
    encrypt        = false
  }

  tags = {
    tagKey = "tagValue"
  }
}

resource "tencentcloud_api_gateway_upstream" "example_cvm" {
  scheme               = "HTTP"
  algorithm            = "ROUND-ROBIN"
  uniq_vpc_id          = tencentcloud_vpc.vpc.id
  upstream_name        = "tf_example"
  upstream_description = "desc."
  upstream_type        = "IP_PORT"
  retries              = 5

  nodes {
    host           = "1.1.1.1"
    port           = 9090
    weight         = 10
    vm_instance_id = tencentcloud_instance.example.id
    tags           = ["tags"]
  }

  tags = {
    "createdBy" = "terraform"
  }
}
`

const testAccApigatewayUpstreamCvmUpdateCvm = `
data "tencentcloud_availability_zones_by_product" "zones" {
  product = "cvm"
}

data "tencentcloud_images" "images" {
  image_type       = ["PUBLIC_IMAGE"]
  image_name_regex = "Final"
}

data "tencentcloud_instance_types" "instance_types" {
  filter {
    name   = "instance-family"
    values = ["S5"]
  }

  cpu_core_count   = 2
  exclude_sold_out = true
}

resource "tencentcloud_vpc" "vpc" {
  name       = "vpc-example"
  cidr_block = "10.0.0.0/16"
}

resource "tencentcloud_subnet" "subnet" {
  availability_zone = data.tencentcloud_availability_zones_by_product.zones.zones.3.name
  name              = "subnet-example"
  vpc_id            = tencentcloud_vpc.vpc.id
  cidr_block        = "10.0.0.0/16"
  is_multicast      = false
}

resource "tencentcloud_instance" "example" {
  instance_name     = "tf_example"
  availability_zone = data.tencentcloud_availability_zones_by_product.zones.zones.3.name
  image_id          = data.tencentcloud_images.images.images.0.image_id
  instance_type     = data.tencentcloud_instance_types.instance_types.instance_types.0.instance_type
  system_disk_type  = "CLOUD_PREMIUM"
  system_disk_size  = 50
  hostname          = "terraform"
  project_id        = 0
  vpc_id            = tencentcloud_vpc.vpc.id
  subnet_id         = tencentcloud_subnet.subnet.id

  data_disks {
    data_disk_type = "CLOUD_PREMIUM"
    data_disk_size = 50
    encrypt        = false
  }

  tags = {
    tagKey = "tagValue"
  }
}

resource "tencentcloud_api_gateway_upstream" "example_cvm" {
  scheme               = "HTTPS"
  algorithm            = "ROUND-ROBIN"
  uniq_vpc_id          = tencentcloud_vpc.vpc.id
  upstream_name        = "tf_example_update"
  upstream_description = "update desc."
  upstream_type        = "IP_PORT"
  retries              = 10

  nodes {
    host           = "2.2.2.2"
    port           = 8080
    weight         = 20
    vm_instance_id = tencentcloud_instance.example.id
    tags           = ["tags"]
  }

  tags = {
    "createdBy" = "terraform"
  }
}
`
