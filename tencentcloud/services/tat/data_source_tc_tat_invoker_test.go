package tat_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudTatInvokerDataSource -v
func TestAccTencentCloudTatInvokerDataSource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheck(t) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceTatInvoker,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_tat_invoker.invoker"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tat_invoker.invoker", "invoker_set.0.command_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tat_invoker.invoker", "invoker_set.0.created_time"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tat_invoker.invoker", "invoker_set.0.enable"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tat_invoker.invoker", "invoker_set.0.instance_ids.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tat_invoker.invoker", "invoker_set.0.invoker_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tat_invoker.invoker", "invoker_set.0.name"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tat_invoker.invoker", "invoker_set.0.schedule_settings.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tat_invoker.invoker", "invoker_set.0.type"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tat_invoker.invoker", "invoker_set.0.updated_time"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tat_invoker.invoker", "invoker_set.0.username"),
				),
			},
		},
	})
}

const testAccDataSourceTatInvoker = `

data "tencentcloud_tat_invoker" "invoker" {
	# invoker_id = ""
	# command_id = ""
	# type = ""
}

`
