package dbdc_test

import (
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	dbdcv20201029 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dbdc/v20201029"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/dbdc"
)

// go test ./tencentcloud/services/dbdc/ -run "TestDbdcDbCustomImagesDS" -v -count=1 -gcflags="all=-l"

func TestDbdcDbCustomImagesDS_ReadBasic(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	dbdcClient := &dbdcv20201029.Client{}
	patches.ApplyMethodReturn(newMockMetaDbdcDS().client, "UseDbdcV20201029Client", dbdcClient)

	patches.ApplyMethodFunc(dbdcClient, "DescribeDBCustomImages", func(request *dbdcv20201029.DescribeDBCustomImagesRequest) (*dbdcv20201029.DescribeDBCustomImagesResponse, error) {
		resp := dbdcv20201029.NewDescribeDBCustomImagesResponse()
		resp.Response = &dbdcv20201029.DescribeDBCustomImagesResponseParams{
			TotalCount: ptrInt64(2),
			ImageSet: []*dbdcv20201029.DBCustomImage{
				{
					ImageId:      ptrStr("img-abc123"),
					OsName:       ptrStr("CentOS 7.6"),
					ImageType:    ptrStr("PUBLIC_IMAGE"),
					Architecture: ptrStr("x86_64"),
				},
				{
					ImageId:      ptrStr("img-def456"),
					OsName:       ptrStr("Ubuntu 20.04"),
					ImageType:    ptrStr("PRIVATE_IMAGE"),
					Architecture: ptrStr("arm64"),
				},
			},
		}
		return resp, nil
	})

	meta := newMockMetaDbdcDS()
	res := dbdc.DataSourceTencentCloudDbdcDbCustomImages()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.NotEmpty(t, d.Id())

	imageSet := d.Get("image_set").([]interface{})
	assert.Len(t, imageSet, 2)

	image0 := imageSet[0].(map[string]interface{})
	assert.Equal(t, "img-abc123", image0["image_id"].(string))
	assert.Equal(t, "CentOS 7.6", image0["os_name"].(string))
	assert.Equal(t, "PUBLIC_IMAGE", image0["image_type"].(string))
	assert.Equal(t, "x86_64", image0["architecture"].(string))

	image1 := imageSet[1].(map[string]interface{})
	assert.Equal(t, "img-def456", image1["image_id"].(string))
	assert.Equal(t, "Ubuntu 20.04", image1["os_name"].(string))
	assert.Equal(t, "PRIVATE_IMAGE", image1["image_type"].(string))
	assert.Equal(t, "arm64", image1["architecture"].(string))
}

func TestDbdcDbCustomImagesDS_ReadWithNilFields(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	dbdcClient := &dbdcv20201029.Client{}
	patches.ApplyMethodReturn(newMockMetaDbdcDS().client, "UseDbdcV20201029Client", dbdcClient)

	patches.ApplyMethodFunc(dbdcClient, "DescribeDBCustomImages", func(request *dbdcv20201029.DescribeDBCustomImagesRequest) (*dbdcv20201029.DescribeDBCustomImagesResponse, error) {
		resp := dbdcv20201029.NewDescribeDBCustomImagesResponse()
		resp.Response = &dbdcv20201029.DescribeDBCustomImagesResponseParams{
			TotalCount: ptrInt64(1),
			ImageSet: []*dbdcv20201029.DBCustomImage{
				{
					ImageId:   ptrStr("img-nil-fields"),
					ImageType: ptrStr("PUBLIC_IMAGE"),
				},
			},
		}
		return resp, nil
	})

	meta := newMockMetaDbdcDS()
	res := dbdc.DataSourceTencentCloudDbdcDbCustomImages()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})

	err := res.Read(d, meta)
	assert.NoError(t, err)

	imageSet := d.Get("image_set").([]interface{})
	assert.Len(t, imageSet, 1)

	image0 := imageSet[0].(map[string]interface{})
	assert.Equal(t, "img-nil-fields", image0["image_id"].(string))
	assert.Equal(t, "PUBLIC_IMAGE", image0["image_type"].(string))
	// OsName and Architecture are nil in the API response, Terraform SDK defaults them to empty strings
	assert.Equal(t, "", image0["os_name"].(string))
	assert.Equal(t, "", image0["architecture"].(string))
}

func TestDbdcDbCustomImagesDS_ReadWithEmptyResponse(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	dbdcClient := &dbdcv20201029.Client{}
	patches.ApplyMethodReturn(newMockMetaDbdcDS().client, "UseDbdcV20201029Client", dbdcClient)

	patches.ApplyMethodFunc(dbdcClient, "DescribeDBCustomImages", func(request *dbdcv20201029.DescribeDBCustomImagesRequest) (*dbdcv20201029.DescribeDBCustomImagesResponse, error) {
		resp := dbdcv20201029.NewDescribeDBCustomImagesResponse()
		resp.Response = &dbdcv20201029.DescribeDBCustomImagesResponseParams{
			TotalCount: ptrInt64(0),
			ImageSet:   []*dbdcv20201029.DBCustomImage{},
		}
		return resp, nil
	})

	meta := newMockMetaDbdcDS()
	res := dbdc.DataSourceTencentCloudDbdcDbCustomImages()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})

	err := res.Read(d, meta)
	// When response is empty, the datasource should return an error (NonRetryableError)
	assert.Error(t, err)
}

func TestDbdcDbCustomImagesDS_Schema(t *testing.T) {
	res := dbdc.DataSourceTencentCloudDbdcDbCustomImages()

	assert.NotNil(t, res)
	assert.Contains(t, res.Schema, "image_set")
	assert.Contains(t, res.Schema, "result_output_file")

	imageSetSchema := res.Schema["image_set"]
	assert.Equal(t, schema.TypeList, imageSetSchema.Type)
	assert.True(t, imageSetSchema.Computed)

	elemRes := imageSetSchema.Elem.(*schema.Resource)
	assert.Contains(t, elemRes.Schema, "image_id")
	assert.Contains(t, elemRes.Schema, "os_name")
	assert.Contains(t, elemRes.Schema, "image_type")
	assert.Contains(t, elemRes.Schema, "architecture")

	imageIdSchema := elemRes.Schema["image_id"]
	assert.Equal(t, schema.TypeString, imageIdSchema.Type)
	assert.True(t, imageIdSchema.Computed)

	osNameSchema := elemRes.Schema["os_name"]
	assert.Equal(t, schema.TypeString, osNameSchema.Type)
	assert.True(t, osNameSchema.Computed)

	imageTypeSchema := elemRes.Schema["image_type"]
	assert.Equal(t, schema.TypeString, imageTypeSchema.Type)
	assert.True(t, imageTypeSchema.Computed)

	architectureSchema := elemRes.Schema["architecture"]
	assert.Equal(t, schema.TypeString, architectureSchema.Type)
	assert.True(t, architectureSchema.Computed)
}
