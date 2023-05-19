package tencentcloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const testAccDbbrainTopSpaceSchemaObject = "data.tencentcloud_dbbrain_top_space_schemas.top_space_schemas"

func TestAccTencentCloudDbbrainTopSpaceSchemasDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccDbbrainTopSpaceSchemasDataSource, defaultDbBrainInstanceId),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID(testAccDbbrainTopSpaceSchemaObject),
					resource.TestCheckResourceAttr(testAccDbbrainTopSpaceSchemaObject, "instance_id", defaultDbBrainInstanceId),
					resource.TestCheckResourceAttr(testAccDbbrainTopSpaceSchemaObject, "sort_by", "DataLength"),
					resource.TestCheckResourceAttr(testAccDbbrainTopSpaceSchemaObject, "product", "mysql"),
					resource.TestCheckResourceAttrSet(testAccDbbrainTopSpaceSchemaObject, "top_space_schemas.#"),
					resource.TestCheckResourceAttrSet(testAccDbbrainTopSpaceSchemaObject, "top_space_schemas.0.table_schema"),
					resource.TestCheckResourceAttrSet(testAccDbbrainTopSpaceSchemaObject, "top_space_schemas.0.data_length"),
					resource.TestCheckResourceAttrSet(testAccDbbrainTopSpaceSchemaObject, "top_space_schemas.0.index_length"),
					resource.TestCheckResourceAttrSet(testAccDbbrainTopSpaceSchemaObject, "top_space_schemas.0.data_free"),
					resource.TestCheckResourceAttrSet(testAccDbbrainTopSpaceSchemaObject, "top_space_schemas.0.total_length"),
					resource.TestCheckResourceAttrSet(testAccDbbrainTopSpaceSchemaObject, "top_space_schemas.0.frag_ratio"),
					resource.TestCheckResourceAttrSet(testAccDbbrainTopSpaceSchemaObject, "top_space_schemas.0.table_rows"),
					resource.TestCheckResourceAttrSet(testAccDbbrainTopSpaceSchemaObject, "top_space_schemas.0.physical_file_size"),
					resource.TestCheckResourceAttrSet(testAccDbbrainTopSpaceSchemaObject, "timestamp"),
				),
			},
		},
	})
}

const testAccDbbrainTopSpaceSchemasDataSource = `

data "tencentcloud_dbbrain_top_space_schemas" "top_space_schemas" {
  instance_id = "%s"
  sort_by = "DataLength"
  product = "mysql"
}

`
