package dbdc_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudDbdcDbCustomClusterResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDbdcDbCustomCluster,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_dbdc_db_custom_cluster.example", "id"),
					resource.TestCheckResourceAttr("tencentcloud_dbdc_db_custom_cluster.example", "cluster_name", "tf-example"),
					resource.TestCheckResourceAttr("tencentcloud_dbdc_db_custom_cluster.example", "cluster_description", "tf example cluster"),
					resource.TestCheckResourceAttr("tencentcloud_dbdc_db_custom_cluster.example", "tags.createBy", "Terraform"),
					resource.TestCheckResourceAttrSet("tencentcloud_dbdc_db_custom_cluster.example", "cluster_status"),
				),
			},
			{
				Config: testAccDbdcDbCustomClusterUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("tencentcloud_dbdc_db_custom_cluster.example", "tags.createBy", "TerraformUpdate"),
				),
			},
			{
				ResourceName:      "tencentcloud_dbdc_db_custom_cluster.example",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccDbdcDbCustomCluster = `
resource "tencentcloud_dbdc_db_custom_cluster" "example" {
  cluster_name        = "tf-example"
  cluster_description = "tf example cluster"

  container_network {
    vpc_id     = "vpc-xxxxxxxx"
    subnet_ids = ["subnet-xxxxxxxx"]
  }

  api_server_network {
    vpc_id    = "vpc-xxxxxxxx"
    subnet_id = "subnet-xxxxxxxx"
  }

  tags = {
    createBy = "Terraform"
  }
}
`

const testAccDbdcDbCustomClusterUpdate = `
resource "tencentcloud_dbdc_db_custom_cluster" "example" {
  cluster_name        = "tf-example"
  cluster_description = "tf example cluster"

  container_network {
    vpc_id     = "vpc-xxxxxxxx"
    subnet_ids = ["subnet-xxxxxxxx"]
  }

  api_server_network {
    vpc_id    = "vpc-xxxxxxxx"
    subnet_id = "subnet-xxxxxxxx"
  }

  tags = {
    createBy = "TerraformUpdate"
  }
}
`
