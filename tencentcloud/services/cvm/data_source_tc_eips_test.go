package cvm_test

import (
	"testing"

	acctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudCvmEipsDataSource_Basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.AccPreCheck(t)
		},
		Providers: acctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCvmEipsDataSource_BasicCreate,
				Check:  resource.ComposeTestCheckFunc(acctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_eips.example")),
			},
		},
	})
}

const testAccCvmEipsDataSource_BasicCreate = `

data "tencentcloud_eips" "example" {
}

`

func TestAccTencentCloudCvmEipsDataSource_ById(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.AccPreCheck(t)
		},
		Providers: acctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCvmEipsDataSource_ByIdCreate,
				Check:  resource.ComposeTestCheckFunc(acctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_eips.example_by_id"), resource.TestCheckResourceAttrSet("data.tencentcloud_eips.example_by_id", "eip_list.0.public_ip"), resource.TestCheckResourceAttrSet("data.tencentcloud_eips.example_by_id", "eip_list.0.create_time"), resource.TestCheckResourceAttr("data.tencentcloud_eips.example_by_id", "eip_list.#", "1"), resource.TestCheckResourceAttrSet("data.tencentcloud_eips.example_by_id", "eip_list.0.eip_id"), resource.TestCheckResourceAttr("data.tencentcloud_eips.example_by_id", "eip_list.0.eip_name", "tf-example"), resource.TestCheckResourceAttrSet("data.tencentcloud_eips.example_by_id", "eip_list.0.eip_type"), resource.TestCheckResourceAttrSet("data.tencentcloud_eips.example_by_id", "eip_list.0.status")),
			},
		},
	})
}

const testAccCvmEipsDataSource_ByIdCreate = `

data "tencentcloud_eips" "example_by_id" {
    eip_id = tencentcloud_eip.example.id
}
resource "tencentcloud_eip" "example" {
    
    tags = {
        test = "test"
    }
    name = "tf-example"
}

`

func TestAccTencentCloudCvmEipsDataSource_ByName(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.AccPreCheck(t)
		},
		Providers: acctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCvmEipsDataSource_ByNameCreate,
				Check:  resource.ComposeTestCheckFunc(acctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_eips.example_by_name"), resource.TestCheckResourceAttr("data.tencentcloud_eips.example_by_name", "eip_list.0.eip_name", "tf-example")),
			},
		},
	})
}

const testAccCvmEipsDataSource_ByNameCreate = `

data "tencentcloud_eips" "example_by_name" {
    eip_name = tencentcloud_eip.example.name
}
resource "tencentcloud_eip" "example" {
    name = "tf-example"
    
    tags = {
        test = "test"
    }
}

`

func TestAccTencentCloudCvmEipsDataSource_ByTags(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.AccPreCheck(t)
		},
		Providers: acctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCvmEipsDataSource_ByTagsCreate,
				Check:  resource.ComposeTestCheckFunc(acctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_eips.example_by_tags"), resource.TestCheckResourceAttr("data.tencentcloud_eips.example_by_tags", "eip_list.0.tags.test", "test")),
			},
		},
	})
}

const testAccCvmEipsDataSource_ByTagsCreate = `

data "tencentcloud_eips" "example_by_tags" {
    tags = tencentcloud_eip.example.tags
}
resource "tencentcloud_eip" "example" {
    name = "tf-example"
    
    tags = {
        test = "test"
    }
}

`
