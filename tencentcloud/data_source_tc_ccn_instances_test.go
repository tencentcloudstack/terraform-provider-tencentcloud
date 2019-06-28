package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccDataSourceTencentCloudCcnV3Instances_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: TestAccDataSourceTencentCloudCcnInstances,

				Check: resource.ComposeTestCheckFunc(

					//id filter
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_ccn_instances.id_instances"),
					resource.TestCheckResourceAttr("data.tencentcloud_ccn_instances.id_instances", "instance_list.#", "1"),
					resource.TestCheckResourceAttr("data.tencentcloud_ccn_instances.id_instances", "instance_list.0.name", "ci-temp-test-ccn"),
					resource.TestCheckResourceAttr("data.tencentcloud_ccn_instances.id_instances", "instance_list.0.description", "ci-temp-test-ccn-des"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_ccn_instances.id_instances", "instance_list.0.ccn_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_ccn_instances.id_instances", "instance_list.0.qos"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_ccn_instances.id_instances", "instance_list.0.state"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_ccn_instances.id_instances", "instance_list.0.attachment_list.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_ccn_instances.id_instances", "instance_list.0.create_time"),

					//name filter ,Every VPC with a "guagua_vpc_instance_test" name will be found
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_ccn_instances.name_instances"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_ccn_instances.id_instances", "instance_list.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_ccn_instances.id_instances", "instance_list.0.name"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_ccn_instances.id_instances", "instance_list.0.description"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_ccn_instances.id_instances", "instance_list.0.ccn_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_ccn_instances.id_instances", "instance_list.0.qos"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_ccn_instances.id_instances", "instance_list.0.state"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_ccn_instances.id_instances", "instance_list.0.attachment_list.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_ccn_instances.id_instances", "instance_list.0.create_time"),
				),
			},
		},
	})
}

const TestAccDataSourceTencentCloudCcnInstances = `
resource tencentcloud_ccn main{
	name ="ci-temp-test-ccn"
	description="ci-temp-test-ccn-des"
	qos ="AG"
}

data tencentcloud_ccn_instances id_instances{
	ccn_id = "${tencentcloud_ccn.main.id}"
}

data tencentcloud_ccn_instances name_instances{
	name = "${tencentcloud_ccn.main.name}"
}

`
