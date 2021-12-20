package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

var testDataTCRInstancesNameAll = "data.tencentcloud_tcr_instances.all_test"

func TestAccTencentCloudDataTCRInstances(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckTCRInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTencentCloudDataTCRInstancesBasic,
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckTCRInstanceExists("tencentcloud_tcr_instance.mytcr_instance"),
					resource.TestCheckResourceAttrSet(testDataTCRInstancesNameAll, "instance_list.0.id"),
					resource.TestCheckResourceAttr(testDataTCRInstancesNameAll, "instance_list.0.instance_type", "standard"),
					resource.TestCheckResourceAttrSet(testDataTCRInstancesNameAll, "instance_list.0.internal_end_point"),
					resource.TestCheckResourceAttrSet(testDataTCRInstancesNameAll, "instance_list.0.public_domain"),
					resource.TestCheckResourceAttrSet(testDataTCRInstancesNameAll, "instance_list.0.status"),
					resource.TestCheckResourceAttr(testDataTCRInstancesNameAll, "instance_list.0.tags.test", "test"),
				),
			},
		},
	})
}

const testAccTencentCloudDataTCRInstancesBasic = `
resource "tencentcloud_tcr_instance" "mytcr_instance" {
  name        = "testacctcrinstance"
  instance_type = "standard"
  delete_bucket = true

  tags ={
	test = "test"
  }
}

data "tencentcloud_tcr_instances" "all_test" {
  name = tencentcloud_tcr_instance.mytcr_instance.name
}

`
