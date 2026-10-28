package ckafka_test

import (
	"context"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	ckafka "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ckafka/v20190819"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	localckafka "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/ckafka"
)

// -----------------------------------------------------------------------------
// Mock-based tests for tencentcloud_ckafka_routes data source Read path.
//
// go test ./tencentcloud/services/ckafka/ -run "TestCkafkaRoutesDataSource" -v -count=1 -gcflags="all=-l"
// -----------------------------------------------------------------------------

type mockMetaCkafkaRoutesDataSource struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMetaCkafkaRoutesDataSource) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockMetaCkafkaRoutesDataSource{}

func newMockMetaCkafkaRoutesDataSource() *mockMetaCkafkaRoutesDataSource {
	return &mockMetaCkafkaRoutesDataSource{client: &connectivity.TencentCloudClient{}}
}

func routesPtrString(s string) *string { return &s }
func routesPtrInt64(v int64) *int64    { return &v }

// buildCkafkaRoute builds a Route with all fields populated for the Read path.
func buildCkafkaRoute(routeId int64, status int64) *ckafka.Route {
	return &ckafka.Route{
		AccessType: routesPtrInt64(0),
		RouteId:    routesPtrInt64(routeId),
		VipType:    routesPtrInt64(3),
		VipList: []*ckafka.VipEntity{
			{
				Vip:   routesPtrString("10.0.0.1"),
				Vport: routesPtrString("9092"),
			},
		},
		Domain:          routesPtrString("ckafka-route.example.com"),
		DomainPort:      routesPtrInt64(9092),
		DeleteTimestamp: routesPtrString(""),
		Subnet:          routesPtrString("subnet-j5vja918"),
		BrokerVipList: []*ckafka.VipEntity{
			{
				Vip:   routesPtrString("10.0.0.2"),
				Vport: routesPtrString("9093"),
			},
		},
		VpcId:  routesPtrString("vpc-axrsmmrv"),
		Note:   routesPtrString("test route"),
		Status: routesPtrInt64(status),
	}
}

// TestCkafkaRoutesDataSource_Read_Basic verifies that the Read function maps
// the DescribeRoute response into the routers list, including nested
// vip_list/broker_vip_list expansion.
func TestCkafkaRoutesDataSource_Read_Basic(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	ckafkaClient := &ckafka.Client{}
	patches.ApplyMethodReturn(newMockMetaCkafkaRoutesDataSource().client, "UseCkafkaClient", ckafkaClient)

	instanceId := "ckafka-routes-test"
	patches.ApplyMethodFunc(ckafkaClient, "DescribeRoute", func(request *ckafka.DescribeRouteRequest) (*ckafka.DescribeRouteResponse, error) {
		resp := ckafka.NewDescribeRouteResponse()
		resp.Response = &ckafka.DescribeRouteResponseParams{
			Result: &ckafka.RouteResponse{
				Routers: []*ckafka.Route{
					buildCkafkaRoute(135912, 2),
				},
			},
			RequestId: routesPtrString("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaCkafkaRoutesDataSource()
	res := localckafka.DataSourceTencentCloudCkafkaRoutes()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"instance_id": instanceId,
	})

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.NotEmpty(t, d.Id())

	routers := d.Get("routers").([]interface{})
	assert.Equal(t, 1, len(routers))

	route := routers[0].(map[string]interface{})
	assert.Equal(t, int64(0), route["access_type"])
	assert.Equal(t, int64(135912), route["route_id"])
	assert.Equal(t, int64(3), route["vip_type"])
	assert.Equal(t, "ckafka-route.example.com", route["domain"])
	assert.Equal(t, int64(9092), route["domain_port"])
	assert.Equal(t, "subnet-j5vja918", route["subnet"])
	assert.Equal(t, "vpc-axrsmmrv", route["vpc_id"])
	assert.Equal(t, "test route", route["note"])
	assert.Equal(t, int64(2), route["status"])

	vipList := route["vip_list"].([]interface{})
	assert.Equal(t, 1, len(vipList))
	vip := vipList[0].(map[string]interface{})
	assert.Equal(t, "10.0.0.1", vip["vip"])
	assert.Equal(t, "9092", vip["vport"])

	brokerVipList := route["broker_vip_list"].([]interface{})
	assert.Equal(t, 1, len(brokerVipList))
	brokerVip := brokerVipList[0].(map[string]interface{})
	assert.Equal(t, "10.0.0.2", brokerVip["vip"])
	assert.Equal(t, "9093", brokerVip["vport"])
}

// TestCkafkaRoutesDataSource_Read_MultipleRoutes verifies that multiple routes
// in the response are all mapped into the routers list.
func TestCkafkaRoutesDataSource_Read_MultipleRoutes(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	ckafkaClient := &ckafka.Client{}
	patches.ApplyMethodReturn(newMockMetaCkafkaRoutesDataSource().client, "UseCkafkaClient", ckafkaClient)

	instanceId := "ckafka-routes-multi"
	patches.ApplyMethodFunc(ckafkaClient, "DescribeRoute", func(request *ckafka.DescribeRouteRequest) (*ckafka.DescribeRouteResponse, error) {
		resp := ckafka.NewDescribeRouteResponse()
		resp.Response = &ckafka.DescribeRouteResponseParams{
			Result: &ckafka.RouteResponse{
				Routers: []*ckafka.Route{
					buildCkafkaRoute(1, 2),
					buildCkafkaRoute(2, 2),
				},
			},
			RequestId: routesPtrString("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaCkafkaRoutesDataSource()
	res := localckafka.DataSourceTencentCloudCkafkaRoutes()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"instance_id": instanceId,
	})

	err := res.Read(d, meta)
	assert.NoError(t, err)

	routers := d.Get("routers").([]interface{})
	assert.Equal(t, 2, len(routers))
}

// TestCkafkaRoutesDataSource_Read_WithRouteId verifies that the route_id
// optional parameter is passed through to the request.
func TestCkafkaRoutesDataSource_Read_WithRouteId(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	ckafkaClient := &ckafka.Client{}
	patches.ApplyMethodReturn(newMockMetaCkafkaRoutesDataSource().client, "UseCkafkaClient", ckafkaClient)

	instanceId := "ckafka-routes-routeid"
	var capturedRouteId *int64
	var capturedMainRouteFlag *bool
	patches.ApplyMethodFunc(ckafkaClient, "DescribeRoute", func(request *ckafka.DescribeRouteRequest) (*ckafka.DescribeRouteResponse, error) {
		capturedRouteId = request.RouteId
		capturedMainRouteFlag = request.MainRouteFlag
		resp := ckafka.NewDescribeRouteResponse()
		resp.Response = &ckafka.DescribeRouteResponseParams{
			Result: &ckafka.RouteResponse{
				Routers: []*ckafka.Route{
					buildCkafkaRoute(135912, 2),
				},
			},
			RequestId: routesPtrString("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaCkafkaRoutesDataSource()
	res := localckafka.DataSourceTencentCloudCkafkaRoutes()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"instance_id":     instanceId,
		"route_id":        135912,
		"main_route_flag": true,
	})

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.NotNil(t, capturedRouteId)
	assert.Equal(t, int64(135912), *capturedRouteId)
	assert.NotNil(t, capturedMainRouteFlag)
	assert.True(t, *capturedMainRouteFlag)
}

// TestCkafkaRoutesDataSource_Read_NilFieldsSkipped verifies that nil fields in
// the response route are skipped (not set) without causing panics.
func TestCkafkaRoutesDataSource_Read_NilFieldsSkipped(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	ckafkaClient := &ckafka.Client{}
	patches.ApplyMethodReturn(newMockMetaCkafkaRoutesDataSource().client, "UseCkafkaClient", ckafkaClient)

	instanceId := "ckafka-routes-nil"
	patches.ApplyMethodFunc(ckafkaClient, "DescribeRoute", func(request *ckafka.DescribeRouteRequest) (*ckafka.DescribeRouteResponse, error) {
		resp := ckafka.NewDescribeRouteResponse()
		resp.Response = &ckafka.DescribeRouteResponseParams{
			Result: &ckafka.RouteResponse{
				Routers: []*ckafka.Route{
					{
						// Only RouteId is set; all other fields are nil.
						RouteId: routesPtrInt64(1),
					},
				},
			},
			RequestId: routesPtrString("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaCkafkaRoutesDataSource()
	res := localckafka.DataSourceTencentCloudCkafkaRoutes()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"instance_id": instanceId,
	})

	err := res.Read(d, meta)
	assert.NoError(t, err)

	routers := d.Get("routers").([]interface{})
	assert.Equal(t, 1, len(routers))
	route := routers[0].(map[string]interface{})
	assert.Equal(t, int64(1), route["route_id"])
	// nil fields should yield zero values
	assert.Equal(t, 0, route["access_type"])
	assert.Equal(t, "", route["domain"])
}

// TestCkafkaRoutesDataSource_Read_EmptyResponse verifies that when the
// DescribeRoute response returns an empty routers list, the Read function
// returns an error (NonRetryableError) and does not clear the id.
func TestCkafkaRoutesDataSource_Read_EmptyResponse(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	ckafkaClient := &ckafka.Client{}
	patches.ApplyMethodReturn(newMockMetaCkafkaRoutesDataSource().client, "UseCkafkaClient", ckafkaClient)

	instanceId := "ckafka-routes-empty"
	patches.ApplyMethodFunc(ckafkaClient, "DescribeRoute", func(request *ckafka.DescribeRouteRequest) (*ckafka.DescribeRouteResponse, error) {
		resp := ckafka.NewDescribeRouteResponse()
		resp.Response = &ckafka.DescribeRouteResponseParams{
			Result: &ckafka.RouteResponse{
				Routers: []*ckafka.Route{},
			},
			RequestId: routesPtrString("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaCkafkaRoutesDataSource()
	res := localckafka.DataSourceTencentCloudCkafkaRoutes()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"instance_id": instanceId,
	})

	err := res.Read(d, meta)
	assert.Error(t, err)
}

// TestCkafkaRoutesDataSource_Read_NilResult verifies that when the
// DescribeRoute response returns nil Result, the Read function returns an
// error (NonRetryableError).
func TestCkafkaRoutesDataSource_Read_NilResult(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	ckafkaClient := &ckafka.Client{}
	patches.ApplyMethodReturn(newMockMetaCkafkaRoutesDataSource().client, "UseCkafkaClient", ckafkaClient)

	instanceId := "ckafka-routes-nil-result"
	patches.ApplyMethodFunc(ckafkaClient, "DescribeRoute", func(request *ckafka.DescribeRouteRequest) (*ckafka.DescribeRouteResponse, error) {
		resp := ckafka.NewDescribeRouteResponse()
		resp.Response = &ckafka.DescribeRouteResponseParams{
			Result:    nil,
			RequestId: routesPtrString("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaCkafkaRoutesDataSource()
	res := localckafka.DataSourceTencentCloudCkafkaRoutes()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"instance_id": instanceId,
	})

	err := res.Read(d, meta)
	assert.Error(t, err)
}

// TestCkafkaRoutesDataSource_Read_APIError verifies that when the DescribeRoute
// API returns an error, the Read function returns the error.
func TestCkafkaRoutesDataSource_Read_APIError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	ckafkaClient := &ckafka.Client{}
	patches.ApplyMethodReturn(newMockMetaCkafkaRoutesDataSource().client, "UseCkafkaClient", ckafkaClient)

	instanceId := "ckafka-routes-api-error"
	patches.ApplyMethodFunc(ckafkaClient, "DescribeRoute", func(request *ckafka.DescribeRouteRequest) (*ckafka.DescribeRouteResponse, error) {
		return nil, context.DeadlineExceeded
	})

	meta := newMockMetaCkafkaRoutesDataSource()
	res := localckafka.DataSourceTencentCloudCkafkaRoutes()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"instance_id": instanceId,
	})

	err := res.Read(d, meta)
	assert.Error(t, err)
}
