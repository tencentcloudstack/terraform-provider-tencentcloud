package rum_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudRumTawInstanceDataSource -v
func TestAccTencentCloudRumTawInstanceDataSource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheck(t) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceRumTawInstance,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_rum_taw_instance.taw_instance"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_rum_taw_instance.taw_instance", "instance_set.#"),
					resource.TestCheckResourceAttr("data.tencentcloud_rum_taw_instance.taw_instance", "instance_set.0.area_id", "1"),
					resource.TestCheckResourceAttr("data.tencentcloud_rum_taw_instance.taw_instance", "instance_set.0.charge_status", "1"),
					resource.TestCheckResourceAttr("data.tencentcloud_rum_taw_instance.taw_instance", "instance_set.0.charge_type", "1"),
					resource.TestCheckResourceAttr("data.tencentcloud_rum_taw_instance.taw_instance", "instance_set.0.cluster_id", "0"),
					resource.TestCheckResourceAttr("data.tencentcloud_rum_taw_instance.taw_instance", "instance_set.0.data_retention_days", "30"),
					resource.TestCheckResourceAttr("data.tencentcloud_rum_taw_instance.taw_instance", "instance_set.0.instance_desc", "Automated testing, do not delete"),
					resource.TestCheckResourceAttr("data.tencentcloud_rum_taw_instance.taw_instance", "instance_set.0.instance_id", "rum-pasZKEI3RLgakj"),
					resource.TestCheckResourceAttr("data.tencentcloud_rum_taw_instance.taw_instance", "instance_set.0.instance_name", "keep-rum"),
					resource.TestCheckResourceAttr("data.tencentcloud_rum_taw_instance.taw_instance", "instance_set.0.instance_status", "2"),
				),
			},
		},
	})
}

const testAccDataSourceRumTawInstance = `

data "tencentcloud_rum_taw_instance" "taw_instance" {
	charge_statuses = [1,]
	charge_types = [1,]
	area_ids = [1,]
	instance_statuses = [2,]
	instance_ids = ["rum-pasZKEI3RLgakj",]
}

`
