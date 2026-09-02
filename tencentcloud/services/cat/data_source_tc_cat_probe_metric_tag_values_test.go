package cat_test

import (
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	cat "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cat/v20180409"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/cat"
)

type mockMetaForCatProbeMetricTagValuesDS struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMetaForCatProbeMetricTagValuesDS) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockMetaForCatProbeMetricTagValuesDS{}

func newMockMetaForCatProbeMetricTagValuesDS() *mockMetaForCatProbeMetricTagValuesDS {
	return &mockMetaForCatProbeMetricTagValuesDS{client: &connectivity.TencentCloudClient{}}
}

func ptrStringCatPMDS(s string) *string { return &s }

func TestCatProbeMetricTagValuesDS_ReadBasic(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	catClient := &cat.Client{}
	patches.ApplyMethodReturn(newMockMetaForCatProbeMetricTagValuesDS().client, "UseCatClient", catClient)

	patches.ApplyMethodFunc(catClient, "DescribeProbeMetricTagValues", func(request *cat.DescribeProbeMetricTagValuesRequest) (*cat.DescribeProbeMetricTagValuesResponse, error) {
		resp := cat.NewDescribeProbeMetricTagValuesResponse()
		resp.Response = &cat.DescribeProbeMetricTagValuesResponseParams{
			TagValueSet: ptrStringCatPMDS("[\"www.qq.com\",\"www.baidu.com\"]"),
			RequestId:   ptrStringCatPMDS("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaForCatProbeMetricTagValuesDS()
	res := cat.DataSourceTencentCloudCatProbeMetricTagValues()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"analyze_task_type": "AnalyzeTaskType_Network",
		"key":               "host",
		"filter":            "www.qq.com",
		"time_range":        "1h",
	})

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.NotEmpty(t, d.Id())

	tagValueSet := d.Get("tag_value_set").(string)
	assert.Equal(t, "[\"www.qq.com\",\"www.baidu.com\"]", tagValueSet)
}

func TestCatProbeMetricTagValuesDS_ReadWithFilters(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	catClient := &cat.Client{}
	patches.ApplyMethodReturn(newMockMetaForCatProbeMetricTagValuesDS().client, "UseCatClient", catClient)

	patches.ApplyMethodFunc(catClient, "DescribeProbeMetricTagValues", func(request *cat.DescribeProbeMetricTagValuesRequest) (*cat.DescribeProbeMetricTagValuesResponse, error) {
		assert.Equal(t, "AnalyzeTaskType_Network", *request.AnalyzeTaskType)
		assert.Equal(t, "area", *request.Key)
		assert.Equal(t, int(2), len(request.Filters))
		resp := cat.NewDescribeProbeMetricTagValuesResponse()
		resp.Response = &cat.DescribeProbeMetricTagValuesResponseParams{
			TagValueSet: ptrStringCatPMDS("[\"beijing\",\"shanghai\"]"),
			RequestId:   ptrStringCatPMDS("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaForCatProbeMetricTagValuesDS()
	res := cat.DataSourceTencentCloudCatProbeMetricTagValues()
	filters := schema.NewSet(schema.HashString, []interface{}{})
	filters.Add("\"host\" = 'www.qq.com'")
	filters.Add("time >= now()-1h")
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"analyze_task_type": "AnalyzeTaskType_Network",
		"key":               "area",
		"filters":           filters,
	})

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.NotEmpty(t, d.Id())

	tagValueSet := d.Get("tag_value_set").(string)
	assert.Equal(t, "[\"beijing\",\"shanghai\"]", tagValueSet)
}

func TestCatProbeMetricTagValuesDS_ReadEmptyResponse(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	catClient := &cat.Client{}
	patches.ApplyMethodReturn(newMockMetaForCatProbeMetricTagValuesDS().client, "UseCatClient", catClient)

	patches.ApplyMethodFunc(catClient, "DescribeProbeMetricTagValues", func(request *cat.DescribeProbeMetricTagValuesRequest) (*cat.DescribeProbeMetricTagValuesResponse, error) {
		resp := cat.NewDescribeProbeMetricTagValuesResponse()
		resp.Response = &cat.DescribeProbeMetricTagValuesResponseParams{}
		return resp, nil
	})

	meta := newMockMetaForCatProbeMetricTagValuesDS()
	res := cat.DataSourceTencentCloudCatProbeMetricTagValues()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"analyze_task_type": "AnalyzeTaskType_Network",
		"key":               "host",
	})

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.NotEmpty(t, d.Id())

	tagValueSet := d.Get("tag_value_set").(string)
	assert.Empty(t, tagValueSet)
}

func TestCatProbeMetricTagValuesDS_Schema(t *testing.T) {
	res := cat.DataSourceTencentCloudCatProbeMetricTagValues()

	assert.NotNil(t, res)
	assert.Contains(t, res.Schema, "analyze_task_type")
	assert.Contains(t, res.Schema, "key")
	assert.Contains(t, res.Schema, "filter")
	assert.Contains(t, res.Schema, "filters")
	assert.Contains(t, res.Schema, "time_range")
	assert.Contains(t, res.Schema, "tag_value_set")
	assert.Contains(t, res.Schema, "result_output_file")

	analyzeTaskTypeSchema := res.Schema["analyze_task_type"]
	assert.Equal(t, schema.TypeString, analyzeTaskTypeSchema.Type)
	assert.True(t, analyzeTaskTypeSchema.Optional)

	keySchema := res.Schema["key"]
	assert.Equal(t, schema.TypeString, keySchema.Type)
	assert.True(t, keySchema.Optional)

	filtersSchema := res.Schema["filters"]
	assert.Equal(t, schema.TypeSet, filtersSchema.Type)
	assert.True(t, filtersSchema.Optional)

	tagValueSetSchema := res.Schema["tag_value_set"]
	assert.Equal(t, schema.TypeString, tagValueSetSchema.Type)
	assert.True(t, tagValueSetSchema.Computed)
}
