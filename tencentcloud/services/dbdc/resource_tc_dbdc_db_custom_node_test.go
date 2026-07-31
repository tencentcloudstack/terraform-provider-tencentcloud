package dbdc_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudDbdcDbCustomNodeResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDbdcDbCustomNode,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_dbdc_db_custom_node.example", "id"),
					resource.TestCheckResourceAttr("tencentcloud_dbdc_db_custom_node.example", "zone", "ap-shanghai-5"),
					resource.TestCheckResourceAttr("tencentcloud_dbdc_db_custom_node.example", "node_type", "DB.AT5.8XLARGE128"),
					resource.TestCheckResourceAttr("tencentcloud_dbdc_db_custom_node.example", "node_name", "tf-example"),
					resource.TestCheckResourceAttr("tencentcloud_dbdc_db_custom_node.example", "tags.createBy", "Terraform"),
					resource.TestCheckResourceAttrSet("tencentcloud_dbdc_db_custom_node.example", "node_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_dbdc_db_custom_node.example", "status"),
				),
			},
			{
				Config: testAccDbdcDbCustomNodeUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("tencentcloud_dbdc_db_custom_node.example", "tags.createBy", "TerraformUpdate"),
				),
			},
			{
				ResourceName:            "tencentcloud_dbdc_db_custom_node.example",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"login_settings", "period", "node_count", "auto_voucher", "voucher_ids", "host_name"},
			},
		},
	})
}

func TestAccTencentCloudDbdcDbCustomNodeResource_createParams(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDbdcDbCustomNodeCreateParams,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_dbdc_db_custom_node.example", "id"),
					resource.TestCheckResourceAttr("tencentcloud_dbdc_db_custom_node.example", "charge_type", "PREPAID"),
					resource.TestCheckResourceAttr("tencentcloud_dbdc_db_custom_node.example", "network_mode", "privatelink"),
					resource.TestCheckResourceAttr("tencentcloud_dbdc_db_custom_node.example", "host_name", "tf-example-host"),
					resource.TestCheckResourceAttr("tencentcloud_dbdc_db_custom_node.example", "security_group_ids.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_dbdc_db_custom_node.example", "system_disk.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_dbdc_db_custom_node.example", "data_disks.#", "1"),
				),
			},
			{
				Config: testAccDbdcDbCustomNodeCreateParamsUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("tencentcloud_dbdc_db_custom_node.example", "security_group_ids.#", "2"),
				),
			},
			{
				ResourceName:            "tencentcloud_dbdc_db_custom_node.example",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"login_settings", "period", "node_count", "auto_voucher", "voucher_ids", "host_name"},
			},
		},
	})
}

const testAccDbdcDbCustomNode = `
resource "tencentcloud_dbdc_db_custom_node" "example" {
  zone      = "ap-shanghai-5"
  image_id  = "img-xxxxxxxx"
  vpc_id    = "vpc-xxxxxxxx"
  subnet_id = "subnet-xxxxxxxx"
  node_type = "DB.AT5.8XLARGE128"
  period    = 1
  node_name = "tf-example"

  login_settings {
    password = "Passw0rd@2024"
  }

  auto_renew = 0

  tags = {
    createBy = "Terraform"
  }
}
`

const testAccDbdcDbCustomNodeUpdate = `
resource "tencentcloud_dbdc_db_custom_node" "example" {
  zone      = "ap-shanghai-5"
  image_id  = "img-xxxxxxxx"
  vpc_id    = "vpc-xxxxxxxx"
  subnet_id = "subnet-xxxxxxxx"
  node_type = "DB.AT5.8XLARGE128"
  period    = 1
  node_name = "tf-example"

  login_settings {
    password = "Passw0rd@2024"
  }

  auto_renew = 1

  tags = {
    createBy = "TerraformUpdate"
  }
}
`

const testAccDbdcDbCustomNodeCreateParams = `
resource "tencentcloud_dbdc_db_custom_node" "example" {
  zone       = "ap-shanghai-5"
  image_id   = "img-xxxxxxxx"
  vpc_id     = "vpc-xxxxxxxx"
  subnet_id  = "subnet-xxxxxxxx"
  node_type  = "DB.SA5.8XLARGE128"
  period     = 1
  node_name  = "tf-example"

  charge_type   = "PREPAID"
  network_mode  = "privatelink"
  host_name     = "tf-example-host"
  security_group_ids = ["sg-xxxxxxxx"]

  system_disk {
    disk_type = "CLOUD_HSSD"
    disk_size = 100
  }

  data_disks {
    disk_type = "CLOUD_HSSD"
    disk_size = 500
    disk_name = "tf-example-data"
  }

  login_settings {
    password = "Passw0rd@2024"
  }

  auto_renew = 0

  tags = {
    createBy = "Terraform"
  }
}
`

const testAccDbdcDbCustomNodeCreateParamsUpdate = `
resource "tencentcloud_dbdc_db_custom_node" "example" {
  zone       = "ap-shanghai-5"
  image_id   = "img-xxxxxxxx"
  vpc_id     = "vpc-xxxxxxxx"
  subnet_id  = "subnet-xxxxxxxx"
  node_type  = "DB.SA5.8XLARGE128"
  period     = 1
  node_name  = "tf-example"

  charge_type   = "PREPAID"
  network_mode  = "privatelink"
  host_name     = "tf-example-host"
  security_group_ids = ["sg-xxxxxxxx", "sg-yyyyyyyy"]

  system_disk {
    disk_type = "CLOUD_HSSD"
    disk_size = 100
  }

  data_disks {
    disk_type = "CLOUD_HSSD"
    disk_size = 500
    disk_name = "tf-example-data"
  }

  login_settings {
    password = "Passw0rd@2024"
  }

  auto_renew = 0

  tags = {
    createBy = "Terraform"
  }
}
`
