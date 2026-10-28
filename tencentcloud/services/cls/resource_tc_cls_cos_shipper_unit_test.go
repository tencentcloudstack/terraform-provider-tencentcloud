package cls_test

import (
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	cls "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cls/v20201016"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	localcls "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/cls"
)

type mockMetaForCosShipper struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMetaForCosShipper) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockMetaForCosShipper{}

func newMockMetaForCosShipper() *mockMetaForCosShipper {
	return &mockMetaForCosShipper{client: &connectivity.TencentCloudClient{}}
}

func ptrStrCosShipper(s string) *string {
	return &s
}

func ptrUint64CosShipper(v uint64) *uint64 {
	return &v
}

// go test ./tencentcloud/services/cls/ -run "TestCosShipperTimeZone" -v -count=1 -gcflags="all=-l"

// TestCosShipperTimeZone_Create_WithTimeZone tests Create sets TimeZone in request
func TestCosShipperTimeZone_Create_WithTimeZone(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	clsClient := &cls.Client{}
	patches.ApplyMethodReturn(newMockMetaForCosShipper().client, "UseClsClient", clsClient)

	var capturedRequest *cls.CreateShipperRequest
	patches.ApplyMethodFunc(clsClient, "CreateShipper", func(request *cls.CreateShipperRequest) (*cls.CreateShipperResponse, error) {
		capturedRequest = request
		resp := cls.NewCreateShipperResponse()
		resp.Response = &cls.CreateShipperResponseParams{
			ShipperId: ptrStrCosShipper("shipper-test-123"),
			RequestId: ptrStrCosShipper("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(clsClient, "DescribeShippers", func(request *cls.DescribeShippersRequest) (*cls.DescribeShippersResponse, error) {
		resp := cls.NewDescribeShippersResponse()
		resp.Response = &cls.DescribeShippersResponseParams{
			Shippers: []*cls.ShipperInfo{
				{
					ShipperId:   ptrStrCosShipper("shipper-test-123"),
					TopicId:     ptrStrCosShipper("topic-test-123"),
					Bucket:      ptrStrCosShipper("test-bucket-1234567890"),
					Prefix:      ptrStrCosShipper("logs/"),
					ShipperName: ptrStrCosShipper("test-shipper"),
					Interval:    ptrUint64CosShipper(300),
					MaxSize:     ptrUint64CosShipper(256),
					TimeZone:    ptrStrCosShipper("GMT+08:00"),
				},
			},
			TotalCount: ptrUint64CosShipper(1),
			RequestId:  ptrStrCosShipper("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaForCosShipper()
	res := localcls.ResourceTencentCloudClsCosShipper()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"topic_id":     "topic-test-123",
		"bucket":       "test-bucket-1234567890",
		"prefix":       "logs/",
		"shipper_name": "test-shipper",
		"interval":     300,
		"max_size":     256,
		"time_zone":    "GMT+08:00",
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "shipper-test-123", d.Id())
	assert.NotNil(t, capturedRequest)
	assert.NotNil(t, capturedRequest.TimeZone)
	assert.Equal(t, "GMT+08:00", *capturedRequest.TimeZone)
}

// TestCosShipperTimeZone_Create_WithoutTimeZone tests Create does not set TimeZone when not specified
func TestCosShipperTimeZone_Create_WithoutTimeZone(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	clsClient := &cls.Client{}
	patches.ApplyMethodReturn(newMockMetaForCosShipper().client, "UseClsClient", clsClient)

	var capturedRequest *cls.CreateShipperRequest
	patches.ApplyMethodFunc(clsClient, "CreateShipper", func(request *cls.CreateShipperRequest) (*cls.CreateShipperResponse, error) {
		capturedRequest = request
		resp := cls.NewCreateShipperResponse()
		resp.Response = &cls.CreateShipperResponseParams{
			ShipperId: ptrStrCosShipper("shipper-test-456"),
			RequestId: ptrStrCosShipper("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(clsClient, "DescribeShippers", func(request *cls.DescribeShippersRequest) (*cls.DescribeShippersResponse, error) {
		resp := cls.NewDescribeShippersResponse()
		resp.Response = &cls.DescribeShippersResponseParams{
			Shippers: []*cls.ShipperInfo{
				{
					ShipperId:   ptrStrCosShipper("shipper-test-456"),
					TopicId:     ptrStrCosShipper("topic-test-123"),
					Bucket:      ptrStrCosShipper("test-bucket-1234567890"),
					Prefix:      ptrStrCosShipper("logs/"),
					ShipperName: ptrStrCosShipper("test-shipper"),
					Interval:    ptrUint64CosShipper(300),
					MaxSize:     ptrUint64CosShipper(256),
				},
			},
			TotalCount: ptrUint64CosShipper(1),
			RequestId:  ptrStrCosShipper("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaForCosShipper()
	res := localcls.ResourceTencentCloudClsCosShipper()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"topic_id":     "topic-test-123",
		"bucket":       "test-bucket-1234567890",
		"prefix":       "logs/",
		"shipper_name": "test-shipper",
		"interval":     300,
		"max_size":     256,
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "shipper-test-456", d.Id())
	assert.NotNil(t, capturedRequest)
	assert.Nil(t, capturedRequest.TimeZone)
}

// TestCosShipperTimeZone_Read_WithTimeZone tests Read populates time_zone from API response
func TestCosShipperTimeZone_Read_WithTimeZone(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	clsClient := &cls.Client{}
	patches.ApplyMethodReturn(newMockMetaForCosShipper().client, "UseClsClient", clsClient)

	patches.ApplyMethodFunc(clsClient, "DescribeShippers", func(request *cls.DescribeShippersRequest) (*cls.DescribeShippersResponse, error) {
		resp := cls.NewDescribeShippersResponse()
		resp.Response = &cls.DescribeShippersResponseParams{
			Shippers: []*cls.ShipperInfo{
				{
					ShipperId:   ptrStrCosShipper("shipper-test-789"),
					TopicId:     ptrStrCosShipper("topic-test-123"),
					Bucket:      ptrStrCosShipper("test-bucket-1234567890"),
					Prefix:      ptrStrCosShipper("logs/"),
					ShipperName: ptrStrCosShipper("test-shipper"),
					Interval:    ptrUint64CosShipper(300),
					MaxSize:     ptrUint64CosShipper(256),
					TimeZone:    ptrStrCosShipper("UTC+08:00"),
				},
			},
			TotalCount: ptrUint64CosShipper(1),
			RequestId:  ptrStrCosShipper("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaForCosShipper()
	res := localcls.ResourceTencentCloudClsCosShipper()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"topic_id":     "topic-test-123",
		"bucket":       "test-bucket-1234567890",
		"prefix":       "logs/",
		"shipper_name": "test-shipper",
	})
	d.SetId("shipper-test-789")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "UTC+08:00", d.Get("time_zone"))
}

// TestCosShipperTimeZone_Read_NilTimeZone tests Read handles nil TimeZone
func TestCosShipperTimeZone_Read_NilTimeZone(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	clsClient := &cls.Client{}
	patches.ApplyMethodReturn(newMockMetaForCosShipper().client, "UseClsClient", clsClient)

	patches.ApplyMethodFunc(clsClient, "DescribeShippers", func(request *cls.DescribeShippersRequest) (*cls.DescribeShippersResponse, error) {
		resp := cls.NewDescribeShippersResponse()
		resp.Response = &cls.DescribeShippersResponseParams{
			Shippers: []*cls.ShipperInfo{
				{
					ShipperId:   ptrStrCosShipper("shipper-test-789"),
					TopicId:     ptrStrCosShipper("topic-test-123"),
					Bucket:      ptrStrCosShipper("test-bucket-1234567890"),
					Prefix:      ptrStrCosShipper("logs/"),
					ShipperName: ptrStrCosShipper("test-shipper"),
					Interval:    ptrUint64CosShipper(300),
					MaxSize:     ptrUint64CosShipper(256),
					TimeZone:    nil,
				},
			},
			TotalCount: ptrUint64CosShipper(1),
			RequestId:  ptrStrCosShipper("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaForCosShipper()
	res := localcls.ResourceTencentCloudClsCosShipper()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"topic_id":     "topic-test-123",
		"bucket":       "test-bucket-1234567890",
		"prefix":       "logs/",
		"shipper_name": "test-shipper",
	})
	d.SetId("shipper-test-789")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "", d.Get("time_zone"))
}

// TestCosShipperTimeZone_Update_WithTimeZone tests Update sets TimeZone in request
func TestCosShipperTimeZone_Update_WithTimeZone(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	clsClient := &cls.Client{}
	patches.ApplyMethodReturn(newMockMetaForCosShipper().client, "UseClsClient", clsClient)

	var capturedModifyRequest *cls.ModifyShipperRequest
	patches.ApplyMethodFunc(clsClient, "ModifyShipper", func(request *cls.ModifyShipperRequest) (*cls.ModifyShipperResponse, error) {
		capturedModifyRequest = request
		resp := cls.NewModifyShipperResponse()
		resp.Response = &cls.ModifyShipperResponseParams{
			RequestId: ptrStrCosShipper("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(clsClient, "DescribeShippers", func(request *cls.DescribeShippersRequest) (*cls.DescribeShippersResponse, error) {
		resp := cls.NewDescribeShippersResponse()
		resp.Response = &cls.DescribeShippersResponseParams{
			Shippers: []*cls.ShipperInfo{
				{
					ShipperId:   ptrStrCosShipper("shipper-test-update"),
					TopicId:     ptrStrCosShipper("topic-test-123"),
					Bucket:      ptrStrCosShipper("test-bucket-1234567890"),
					Prefix:      ptrStrCosShipper("logs/"),
					ShipperName: ptrStrCosShipper("test-shipper"),
					Interval:    ptrUint64CosShipper(300),
					MaxSize:     ptrUint64CosShipper(256),
					TimeZone:    ptrStrCosShipper("GMT+09:00"),
				},
			},
			TotalCount: ptrUint64CosShipper(1),
			RequestId:  ptrStrCosShipper("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaForCosShipper()
	res := localcls.ResourceTencentCloudClsCosShipper()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"topic_id":     "topic-test-123",
		"bucket":       "test-bucket-1234567890",
		"prefix":       "logs/",
		"shipper_name": "test-shipper",
		"time_zone":    "GMT+09:00",
	})
	d.SetId("shipper-test-update")

	// Simulate HasChange by marking the key as changed
	_ = d.Set("time_zone", "GMT+08:00")

	err := res.Update(d, meta)
	assert.NoError(t, err)
	assert.NotNil(t, capturedModifyRequest)
	assert.NotNil(t, capturedModifyRequest.TimeZone)
	assert.Equal(t, "GMT+08:00", *capturedModifyRequest.TimeZone)
}

// TestCosShipperTimeZone_Schema tests the time_zone schema definition
func TestCosShipperTimeZone_Schema(t *testing.T) {
	res := localcls.ResourceTencentCloudClsCosShipper()

	assert.NotNil(t, res)
	assert.Contains(t, res.Schema, "time_zone")

	timeZoneSchema := res.Schema["time_zone"]
	assert.Equal(t, schema.TypeString, timeZoneSchema.Type)
	assert.True(t, timeZoneSchema.Optional)
	assert.True(t, timeZoneSchema.Computed)
	assert.False(t, timeZoneSchema.Required)
}
