package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccTencentCloudElasticsearchInstancesDataSource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckElasticsearchInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccElasticsearchInstancesDataSource,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckElasticsearchInstanceExists("tencentcloud_elasticsearch_instance.foo"),
					resource.TestCheckResourceAttr("data.tencentcloud_elasticsearch_instances.foo", "instance_list.#", "1"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_elasticsearch_instances.foo", "instance_list.0.instance_id"),
					resource.TestCheckResourceAttr("data.tencentcloud_elasticsearch_instances.foo", "instance_list.0.instance_name", "tf-ci-test"),
					resource.TestCheckResourceAttr("data.tencentcloud_elasticsearch_instances.foo", "instance_list.0.availability_zone", defaultAZone),
					resource.TestCheckResourceAttr("data.tencentcloud_elasticsearch_instances.foo", "instance_list.0.vpc_id", defaultVpcId),
					resource.TestCheckResourceAttr("data.tencentcloud_elasticsearch_instances.foo", "instance_list.0.subnet_id", defaultSubnetId),
					resource.TestCheckResourceAttr("data.tencentcloud_elasticsearch_instances.foo", "instance_list.0.version", "7.5.1"),
					resource.TestCheckResourceAttr("data.tencentcloud_elasticsearch_instances.foo", "instance_list.0.license_type", "oss"),
					resource.TestCheckResourceAttr("data.tencentcloud_elasticsearch_instances.foo", "instance_list.0.charge_type", "POSTPAID_BY_HOUR"),
					resource.TestCheckResourceAttr("data.tencentcloud_elasticsearch_instances.foo", "instance_list.0.tags.test", "terraform"),
					resource.TestCheckResourceAttr("data.tencentcloud_elasticsearch_instances.foo", "instance_list.0.node_info_list.#", "1"),
					resource.TestCheckResourceAttr("data.tencentcloud_elasticsearch_instances.foo", "instance_list.0.node_info_list.0.node_num", "2"),
					resource.TestCheckResourceAttr("data.tencentcloud_elasticsearch_instances.foo", "instance_list.0.node_info_list.0.node_type", "ES.S1.MEDIUM4"),
					resource.TestCheckResourceAttr("data.tencentcloud_elasticsearch_instances.foo", "instance_list.0.node_info_list.0.encrypt", "false"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_elasticsearch_instances.foo", "instance_list.0.create_time"),
				),
			},
		},
	})
}

const testAccElasticsearchInstancesDataSource = defaultVpcVariable + `
resource "tencentcloud_elasticsearch_instance" "foo" {
	instance_name     = "tf-ci-test"
	availability_zone = var.availability_zone
	version           = "7.5.1"
	vpc_id            = var.vpc_id
	subnet_id         = var.subnet_id
	password          = "Test1234"
	license_type      = "oss"
  
	node_info_list {
	  node_num  = 2
	  node_type = "ES.S1.MEDIUM4"
	}
  
	tags = {
	  test = "terraform"
	}
  }
  
  data "tencentcloud_elasticsearch_instances" "foo" {
	instance_id = tencentcloud_elasticsearch_instance.foo.id
  }
`
