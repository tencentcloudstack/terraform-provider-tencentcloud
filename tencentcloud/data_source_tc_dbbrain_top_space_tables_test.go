package tencentcloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const testAccDbbrainTopSpaceTablesObject = "data.tencentcloud_dbbrain_top_space_tables.top_space_tables"

func TestAccTencentCloudDbbrainTopSpaceTablesDataSource_fromDataLength(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccDbbrainTopSpaceTablesDataSource_dataLen, defaultDbBrainInstanceId),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID(testAccDbbrainTopSpaceTablesObject),
					resource.TestCheckResourceAttr(testAccDbbrainTopSpaceTablesObject, "instance_id", defaultDbBrainInstanceId),
					resource.TestCheckResourceAttr(testAccDbbrainTopSpaceTablesObject, "sort_by", "PhysicalFileSize"),
					resource.TestCheckResourceAttr(testAccDbbrainTopSpaceTablesObject, "product", "mysql"),
					resource.TestCheckResourceAttrSet(testAccDbbrainTopSpaceTablesObject, "top_space_tables.#"),
					resource.TestCheckResourceAttrSet(testAccDbbrainTopSpaceTablesObject, "top_space_tables.0.table_name"),
					resource.TestCheckResourceAttrSet(testAccDbbrainTopSpaceTablesObject, "top_space_tables.0.table_schema"),
					resource.TestCheckResourceAttrSet(testAccDbbrainTopSpaceTablesObject, "top_space_tables.0.engine"),
					resource.TestCheckResourceAttrSet(testAccDbbrainTopSpaceTablesObject, "top_space_tables.0.data_length"),
					resource.TestCheckResourceAttrSet(testAccDbbrainTopSpaceTablesObject, "top_space_tables.0.index_length"),
					resource.TestCheckResourceAttrSet(testAccDbbrainTopSpaceTablesObject, "top_space_tables.0.data_free"),
					resource.TestCheckResourceAttrSet(testAccDbbrainTopSpaceTablesObject, "top_space_tables.0.total_length"),
					resource.TestCheckResourceAttrSet(testAccDbbrainTopSpaceTablesObject, "top_space_tables.0.frag_ratio"),
					resource.TestCheckResourceAttrSet(testAccDbbrainTopSpaceTablesObject, "top_space_tables.0.table_rows"),
					resource.TestCheckResourceAttrSet(testAccDbbrainTopSpaceTablesObject, "top_space_tables.0.physical_file_size"),
					resource.TestCheckResourceAttrSet(testAccDbbrainTopSpaceTablesObject, "timestamp"),
				),
			},
		},
	})
}

func TestAccTencentCloudDbbrainTopSpaceTablesDataSource_fromIndexLength(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccDbbrainTopSpaceTablesDataSource_indexLen, defaultDbBrainInstanceId),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID(testAccDbbrainTopSpaceTablesObject),
					resource.TestCheckResourceAttr(testAccDbbrainTopSpaceTablesObject, "instance_id", defaultDbBrainInstanceId),
					resource.TestCheckResourceAttr(testAccDbbrainTopSpaceTablesObject, "sort_by", "TotalLength"),
					resource.TestCheckResourceAttr(testAccDbbrainTopSpaceTablesObject, "product", "mysql"),
					resource.TestCheckResourceAttrSet(testAccDbbrainTopSpaceTablesObject, "top_space_tables.#"),
					resource.TestCheckResourceAttrSet(testAccDbbrainTopSpaceTablesObject, "top_space_tables.0.table_name"),
					resource.TestCheckResourceAttrSet(testAccDbbrainTopSpaceTablesObject, "top_space_tables.0.table_schema"),
					resource.TestCheckResourceAttrSet(testAccDbbrainTopSpaceTablesObject, "top_space_tables.0.engine"),
					resource.TestCheckResourceAttrSet(testAccDbbrainTopSpaceTablesObject, "top_space_tables.0.data_length"),
					resource.TestCheckResourceAttrSet(testAccDbbrainTopSpaceTablesObject, "top_space_tables.0.index_length"),
					resource.TestCheckResourceAttrSet(testAccDbbrainTopSpaceTablesObject, "top_space_tables.0.data_free"),
					resource.TestCheckResourceAttrSet(testAccDbbrainTopSpaceTablesObject, "top_space_tables.0.total_length"),
					resource.TestCheckResourceAttrSet(testAccDbbrainTopSpaceTablesObject, "top_space_tables.0.frag_ratio"),
					resource.TestCheckResourceAttrSet(testAccDbbrainTopSpaceTablesObject, "top_space_tables.0.table_rows"),
					resource.TestCheckResourceAttrSet(testAccDbbrainTopSpaceTablesObject, "top_space_tables.0.physical_file_size"),
					resource.TestCheckResourceAttrSet(testAccDbbrainTopSpaceTablesObject, "timestamp"),
				),
			},
		},
	})
}

const testAccDbbrainTopSpaceTablesDataSource_dataLen = `

data "tencentcloud_dbbrain_top_space_tables" "top_space_tables" {
  instance_id = "%s"
  sort_by = "PhysicalFileSize" 
  product = "mysql"
}

`

const testAccDbbrainTopSpaceTablesDataSource_indexLen = `

data "tencentcloud_dbbrain_top_space_tables" "top_space_tables" {
  instance_id = "%s"
  sort_by = "TotalLength" 
  product = "mysql"
}

`
