package cynosdb_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudCynosdbInstanceParamResource_basic -v
func TestAccTencentCloudCynosdbInstanceParamResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCynosdbInstanceParam,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_cynosdb_instance_param.instance_param", "id"),
					resource.TestCheckResourceAttr("tencentcloud_cynosdb_instance_param.instance_param", "cluster_id", tcacctest.DefaultCynosdbClusterId),
					resource.TestCheckResourceAttr("tencentcloud_cynosdb_instance_param.instance_param", "instance_id", "cynosdbmysql-ins-rikr6z4o"),
					resource.TestCheckResourceAttr("tencentcloud_cynosdb_instance_param.instance_param", "is_in_maintain_period", "no"),
					resource.TestCheckResourceAttr("tencentcloud_cynosdb_instance_param.instance_param", "instance_param_list.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_cynosdb_instance_param.instance_param", "instance_param_list.0.current_value", "0"),
					resource.TestCheckResourceAttr("tencentcloud_cynosdb_instance_param.instance_param", "instance_param_list.0.param_name", "init_connect"),
				),
			},
		},
	})
}

const testAccCynosdbInstanceParam = tcacctest.CommonCynosdb + `

resource "tencentcloud_cynosdb_instance_param" "instance_param" {
	cluster_id            = var.cynosdb_cluster_id
	instance_id           = "cynosdbmysql-ins-rikr6z4o"
	is_in_maintain_period = "no"
  
	instance_param_list {
	  current_value = "0"
	  param_name    = "init_connect"
	}
}

`
