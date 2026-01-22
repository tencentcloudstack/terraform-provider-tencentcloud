package cvm_test

import (
	"testing"

	acctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudInstanceTypesDataSource_basic -v
func TestAccTencentCloudCvmInstanceTypesDataSource_Basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.AccPreCheck(t)
		},
		Providers: acctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCvmInstanceTypesDataSource_BasicCreate,
				Check:  resource.ComposeTestCheckFunc(acctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_instance_types.example"), resource.TestCheckResourceAttr("data.tencentcloud_instance_types.example", "instance_types.0.cpu_core_count", "4"), resource.TestCheckResourceAttr("data.tencentcloud_instance_types.example", "instance_types.0.memory_size", "8"), resource.TestCheckResourceAttr("data.tencentcloud_instance_types.example", "instance_types.0.availability_zone", "ap-guangzhou-3")),
			},
		},
	})
}

const testAccCvmInstanceTypesDataSource_BasicCreate = `

data "tencentcloud_instance_types" "example" {
    availability_zone = "ap-guangzhou-3"
    cpu_core_count = 4
    memory_size = 8
}

`

func TestAccTencentCloudCvmInstanceTypesDataSource_Sell(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.AccPreCheck(t)
		},
		Providers: acctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCvmInstanceTypesDataSource_SellCreate,
				Check:  resource.ComposeTestCheckFunc(acctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_instance_types.example"), resource.TestCheckResourceAttr("data.tencentcloud_instance_types.example", "instance_types.0.cpu_core_count", "2"), resource.TestCheckResourceAttr("data.tencentcloud_instance_types.example", "instance_types.0.memory_size", "2"), resource.TestCheckResourceAttr("data.tencentcloud_instance_types.example", "instance_types.0.availability_zone", "ap-guangzhou-3"), resource.TestCheckResourceAttr("data.tencentcloud_instance_types.example", "instance_types.0.family", "SA2")),
			},
		},
	})
}
func TestAccTencentCloudCvmInstanceTypesDataSource_WithCbsFilter(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.AccPreCheck(t)
		},
		Providers: acctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCvmInstanceTypesDataSource_WithCbsFilter,
				Check: resource.ComposeTestCheckFunc(
					acctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_instance_types.with_cbs_filter"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_instance_types.with_cbs_filter", "instance_types.0.cbs_configs.#"),
					resource.TestCheckResourceAttr("data.tencentcloud_instance_types.with_cbs_filter", "instance_types.0.cbs_configs.0.disk_type", "CLOUD_SSD"),
					resource.TestCheckResourceAttr("data.tencentcloud_instance_types.with_cbs_filter", "instance_types.0.cbs_configs.0.disk_charge_type", "PREPAID"),
					resource.TestCheckResourceAttr("data.tencentcloud_instance_types.with_cbs_filter", "instance_types.0.cbs_configs.0.disk_usage", "SYSTEM_DISK"),
				),
			},
		},
	})
}

const testAccCvmInstanceTypesDataSource_SellCreate = `

data "tencentcloud_instance_types" "example" {
    cpu_core_count = 2
    memory_size = 2
    exclude_sold_out = true
    
    filter {
        name = "instance-family"
        values = ["SA2"]
    }
    filter {
        name = "zone"
        values = ["ap-guangzhou-3"]
    }
}

`

const testAccCvmInstanceTypesDataSource_WithCbsFilter = `

data "tencentcloud_instance_types" "with_cbs_filter" {
    cpu_core_count = 2
    memory_size = 2
    exclude_sold_out = true
    
    filter {
        name = "instance-family"
        values = ["S6"]
    }
    filter {
        name = "zone"
        values = ["ap-guangzhou-6"]
    }
	cbs_filter {
        disk_types = ["CLOUD_SSD"]
        disk_charge_type = "PREPAID"
        disk_usage = "SYSTEM_DISK"
    }
}
`
