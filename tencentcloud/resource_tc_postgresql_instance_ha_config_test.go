package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudPostgresqlInstanceHAConfigResource_basic -v
func TestAccTencentCloudPostgresqlInstanceHAConfigResource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccPostgresqlInstanceHAConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_postgresql_instance_ha_config.example", "id"),
					resource.TestCheckResourceAttr("tencentcloud_postgresql_instance_ha_config.example", "sync_mode", "Semi-sync"),
					resource.TestCheckResourceAttrSet("tencentcloud_postgresql_instance_ha_config.example", "max_standby_latency"),
					resource.TestCheckResourceAttrSet("tencentcloud_postgresql_instance_ha_config.example", "max_standby_lag"),
					resource.TestCheckResourceAttrSet("tencentcloud_postgresql_instance_ha_config.example", "max_sync_standby_latency"),
					resource.TestCheckResourceAttrSet("tencentcloud_postgresql_instance_ha_config.example", "max_sync_standby_lag"),
				),
			},
			{
				ResourceName:      "tencentcloud_postgresql_instance_ha_config.example",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccPostgresqlInstanceHAConfigAsync,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_postgresql_instance_ha_config.example", "id"),
					resource.TestCheckResourceAttr("tencentcloud_postgresql_instance_ha_config.example", "sync_mode", "Async"),
					resource.TestCheckResourceAttrSet("tencentcloud_postgresql_instance_ha_config.example", "max_standby_latency"),
					resource.TestCheckResourceAttrSet("tencentcloud_postgresql_instance_ha_config.example", "max_standby_lag"),
				),
			},
		},
	})
}

const testAccPostgresqlInstanceHAConfig = `
resource "tencentcloud_postgresql_instance_ha_config" "example" {
  instance_id              = "postgres-gzg9jb2n"
  sync_mode                = "Semi-sync"
  max_standby_latency      = 10737418240
  max_standby_lag          = 10
  max_sync_standby_latency = 52428800
  max_sync_standby_lag     = 5
}
`

const testAccPostgresqlInstanceHAConfigAsync = `
resource "tencentcloud_postgresql_instance_ha_config" "example" {
  instance_id              = "postgres-gzg9jb2n"
  sync_mode                = "Async"
  max_standby_latency      = 10737418240
  max_standby_lag          = 10
}
`
