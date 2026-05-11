package cdb_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudMysqlProxyAddressConfigResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMysqlProxyAddressConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_mysql_proxy_address_config.example", "id"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_proxy_address_config.example", "weight_mode", "system"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_proxy_address_config.example", "is_kick_out", "true"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_proxy_address_config.example", "fail_over", "true"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_proxy_address_config.example", "auto_add_ro", "true"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_proxy_address_config.example", "read_only", "false"),
				),
			},
			{
				Config: testAccMysqlProxyAddressConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("tencentcloud_mysql_proxy_address_config.example", "weight_mode", "system"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_proxy_address_config.example", "is_kick_out", "true"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_proxy_address_config.example", "min_count", "0"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_proxy_address_config.example", "max_delay", "10"),
				),
			},
			{
				ResourceName:      "tencentcloud_mysql_proxy_address_config.example",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"proxy_address_id",
				},
			},
		},
	})
}

const testAccMysqlProxyAddressConfig = `
resource "tencentcloud_mysql_proxy_address_config" "example" {
  instance_id       = "cdb-o2t7gmjl"
  proxy_group_id    = "proxy-ov7dqp8n"
  proxy_address_id  = "proxyaddr-y8dnlfs0"
  weight_mode       = "system"
  is_kick_out       = true
  min_count         = 0
  max_delay         = 10
  fail_over         = true
  auto_add_ro       = true
  read_only         = false
  trans_split       = false
  connection_pool   = true
  auto_load_balance = true
  access_mode       = "nearby"
  proxy_allocation {
    region = "ap-guangzhou"
    zone   = "ap-guangzhou-6"

    proxy_instance {
      instance_id = "cdb-o2t7gmjl"
      weight      = 0
    }
  }

  proxy_allocation {
    region = "ap-guangzhou"
    zone   = "ap-guangzhou-7"

    proxy_instance {
      instance_id = "cdb-o2t7gmjl"
      weight      = 0
    }
  }
}
`

const testAccMysqlProxyAddressConfigUpdate = `
resource "tencentcloud_mysql_proxy_address_config" "example" {
  instance_id       = "cdb-o2t7gmjl"
  proxy_group_id    = "proxy-ov7dqp8n"
  proxy_address_id  = "proxyaddr-y8dnlfs0"
  weight_mode       = "system"
  is_kick_out       = true
  min_count         = 0
  max_delay         = 10
  fail_over         = true
  auto_add_ro       = true
  read_only         = false
  trans_split       = false
  connection_pool   = false
  auto_load_balance = true
  access_mode       = "nearby"
  proxy_allocation {
    region = "ap-guangzhou"
    zone   = "ap-guangzhou-6"

    proxy_instance {
      instance_id = "cdb-o2t7gmjl"
      weight      = 0
    }
  }

  proxy_allocation {
    region = "ap-guangzhou"
    zone   = "ap-guangzhou-7"

    proxy_instance {
      instance_id = "cdb-o2t7gmjl"
      weight      = 0
    }
  }
}
`
