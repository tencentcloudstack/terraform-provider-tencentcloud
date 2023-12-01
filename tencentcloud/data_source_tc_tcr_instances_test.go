package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

var testDataTCRInstancesName = "data.tencentcloud_tcr_instances.tcr"

func TestAccTencentCloudTcrInstancesData(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckCommon(t, ACCOUNT_TYPE_COMMON) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckTCRInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTencentCloudDataTCRInstancesBasic,
				PreConfig: func() {
					testAccStepSetRegion(t, "ap-shanghai")
					testAccPreCheckCommon(t, ACCOUNT_TYPE_COMMON)
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(testDataTCRInstancesName, "instance_list.0.id"),
					resource.TestCheckResourceAttrSet(testDataTCRInstancesName, "instance_list.0.instance_type"),
					resource.TestCheckResourceAttrSet(testDataTCRInstancesName, "instance_list.0.internal_end_point"),
					resource.TestCheckResourceAttrSet(testDataTCRInstancesName, "instance_list.0.public_domain"),
					resource.TestCheckResourceAttrSet(testDataTCRInstancesName, "instance_list.0.status"),
				),
			},
		},
	})
}

const testAccTencentCloudDataTCRInstancesBasic = defaultTCRInstanceData
