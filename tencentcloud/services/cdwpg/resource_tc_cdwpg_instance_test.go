package cdwpg_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudCdwpgInstanceResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_PREPAY) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCdwpgInstance,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_cdwpg_instance.instance", "id"),
					resource.TestCheckResourceAttr("tencentcloud_cdwpg_instance.instance", "instance_name", "test_pg"),
					resource.TestCheckResourceAttr("tencentcloud_cdwpg_instance.instance", "zone", "ap-guangzhou-6"),
					resource.TestCheckResourceAttr("tencentcloud_cdwpg_instance.instance", "user_vpc_id", "vpc-axrsmmrv"),
					resource.TestCheckResourceAttr("tencentcloud_cdwpg_instance.instance", "user_subnet_id", "subnet-kxaxknmg"),
					resource.TestCheckResourceAttr("tencentcloud_cdwpg_instance.instance", "charge_properties.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_cdwpg_instance.instance", "resources.#", "2"),
					resource.TestCheckResourceAttr("tencentcloud_cdwpg_instance.instance", "tags.tagKey", "tagValue"),
				),
			},
			{
				Config: testAccCdwpgInstanceUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_cdwpg_instance.instance", "id"),
					resource.TestCheckResourceAttr("tencentcloud_cdwpg_instance.instance", "instance_name", "test_pg_update"),
					resource.TestCheckResourceAttr("tencentcloud_cdwpg_instance.instance", "tags.tagKey", "tagValueUpdate"),
				),
			},
			{
				Config: testAccCdwpgInstanceScaleOut,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_cdwpg_instance.instance", "id"),
					resource.TestCheckResourceAttr("tencentcloud_cdwpg_instance.instance", "resources.0.count", "4"),
				),
			},
			{
				ResourceName:            "tencentcloud_cdwpg_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"admin_password"},
			},
		},
	})
}

func TestAccTencentCloudCdwpgInstanceResource_withVersion(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_PREPAY) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCdwpgInstanceWithVersion,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_cdwpg_instance.instance", "id"),
					resource.TestCheckResourceAttr("tencentcloud_cdwpg_instance.instance", "product_version", "3.16.9.3"),
				),
			},
			{
				Config: testAccCdwpgInstanceWithVersionUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_cdwpg_instance.instance", "id"),
					resource.TestCheckResourceAttr("tencentcloud_cdwpg_instance.instance", "product_version", "3.16.9.4"),
				),
			},
			{
				ResourceName:            "tencentcloud_cdwpg_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"admin_password"},
			},
		},
	})
}

const testAccCdwpgInstance = `
resource "tencentcloud_cdwpg_instance" "instance" {
	instance_name  = "test_pg"
	zone           = "ap-guangzhou-6"
	user_vpc_id    = "vpc-axrsmmrv"
	user_subnet_id = "subnet-kxaxknmg"
	charge_properties {
	  renew_flag  = 0
	  time_span   = 1
	  time_unit   = "h"
	  charge_type = "POSTPAID_BY_HOUR"
  
	}
	admin_password = "bWJSZDVtVmZkNExJ"
	resources {
	  spec_name = "S_4_16_H_CN"
	  count     = 2
	  disk_spec {
		disk_type  = "CLOUD_HSSD"
		disk_size  = 200
		disk_count = 1
	  }
	  type = "cn"
  
	}
	resources {
	  spec_name = "S_4_16_H_CN"
	  count     = 2
	  disk_spec {
		disk_type  = "CLOUD_HSSD"
		disk_size  = 20
		disk_count = 10
	  }
	  type = "dn"
  
	}
	tags = {
	  "tagKey" = "tagValue"
	}
}
`

const testAccCdwpgInstanceUpdate = `
resource "tencentcloud_cdwpg_instance" "instance" {
	instance_name  = "test_pg_update"
	zone           = "ap-guangzhou-6"
	user_vpc_id    = "vpc-axrsmmrv"
	user_subnet_id = "subnet-kxaxknmg"
	charge_properties {
	  renew_flag  = 0
	  time_span   = 1
	  time_unit   = "h"
	  charge_type = "POSTPAID_BY_HOUR"
  
	}
	admin_password = "bWJSZDVtVmZkNExJ"
	resources {
	  spec_name = "S_4_16_H_CN"
	  count     = 2
	  disk_spec {
		disk_type  = "CLOUD_HSSD"
		disk_size  = 200
		disk_count = 1
	  }
	  type = "cn"
  
	}
	resources {
	  spec_name = "S_4_16_H_CN"
	  count     = 2
	  disk_spec {
		disk_type  = "CLOUD_HSSD"
		disk_size  = 20
		disk_count = 10
	  }
	  type = "dn"
  
	}
	tags = {
	  "tagKey" = "tagValueUpdate"
	}
}
`

const testAccCdwpgInstanceScaleOut = `
resource "tencentcloud_cdwpg_instance" "instance" {
	instance_name  = "test_pg_update"
	zone           = "ap-guangzhou-6"
	user_vpc_id    = "vpc-axrsmmrv"
	user_subnet_id = "subnet-kxaxknmg"
	charge_properties {
	  renew_flag  = 0
	  time_span   = 1
	  time_unit   = "h"
	  charge_type = "POSTPAID_BY_HOUR"
  
	}
	admin_password = "bWJSZDVtVmZkNExJ"
	resources {
	  spec_name = "S_4_16_H_CN"
	  count     = 4
	  disk_spec {
		disk_type  = "CLOUD_HSSD"
		disk_size  = 200
		disk_count = 1
	  }
	  type = "cn"
  
	}
	resources {
	  spec_name = "S_4_16_H_CN"
	  count     = 2
	  disk_spec {
		disk_type  = "CLOUD_HSSD"
		disk_size  = 20
		disk_count = 10
	  }
	  type = "dn"
  
	}
	tags = {
	  "tagKey" = "tagValueUpdate"
	}
}
`

const testAccCdwpgInstanceWithVersion = `
resource "tencentcloud_cdwpg_instance" "instance" {
	instance_name  = "test_pg"
	zone           = "ap-guangzhou-6"
	user_vpc_id    = "vpc-axrsmmrv"
	user_subnet_id = "subnet-kxaxknmg"
	charge_properties {
	  renew_flag  = 0
	  time_span   = 1
	  time_unit   = "h"
	  charge_type = "POSTPAID_BY_HOUR"
  
	}
	admin_password = "bWJSZDVtVmZkNExJ"
	resources {
	  spec_name = "S_4_16_H_CN"
	  count     = 2
	  disk_spec {
		disk_type  = "CLOUD_HSSD"
		disk_size  = 200
		disk_count = 1
	  }
	  type = "cn"
  
	}
	resources {
	  spec_name = "S_4_16_H_CN"
	  count     = 2
	  disk_spec {
		disk_type  = "CLOUD_HSSD"
		disk_size  = 20
		disk_count = 10
	  }
	  type = "dn"
  
	}
	tags = {
	  "tagKey" = "tagValue"
	}
	product_version = "3.16.9.3"
}
`

const testAccCdwpgInstanceWithVersionUpdate = `
resource "tencentcloud_cdwpg_instance" "instance" {
	instance_name  = "test_pg"
	zone           = "ap-guangzhou-6"
	user_vpc_id    = "vpc-axrsmmrv"
	user_subnet_id = "subnet-kxaxknmg"
	charge_properties {
	  renew_flag  = 0
	  time_span   = 1
	  time_unit   = "h"
	  charge_type = "POSTPAID_BY_HOUR"
  
	}
	admin_password = "bWJSZDVtVmZkNExJ"
	resources {
	  spec_name = "S_4_16_H_CN"
	  count     = 2
	  disk_spec {
		disk_type  = "CLOUD_HSSD"
		disk_size  = 200
		disk_count = 1
	  }
	  type = "cn"
  
	}
	resources {
	  spec_name = "S_4_16_H_CN"
	  count     = 2
	  disk_spec {
		disk_type  = "CLOUD_HSSD"
		disk_size  = 20
		disk_count = 10
	  }
	  type = "dn"
  
	}
	tags = {
	  "tagKey" = "tagValue"
	}
	product_version = "3.16.9.4"
}
`
