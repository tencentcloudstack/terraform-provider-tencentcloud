package ssl_test

import (
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	ssl "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ssl/v20191205"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	svcssl "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/ssl"
)

// go test ./tencentcloud/services/ssl/ -run "TestSslCertificateBindResourceTaskDetailDataSource" -v -count=1 -gcflags="all=-l"

type mockMetaBindResourceTaskDetail struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMetaBindResourceTaskDetail) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockMetaBindResourceTaskDetail{}

func newMockMetaBindResourceTaskDetail() *mockMetaBindResourceTaskDetail {
	return &mockMetaBindResourceTaskDetail{client: &connectivity.TencentCloudClient{}}
}

func ptrStrBind(s string) *string {
	return &s
}

func ptrUint64Bind(n uint64) *uint64 {
	return &n
}

func ptrInt64Bind(n int64) *int64 {
	return &n
}

// TestSslCertificateBindResourceTaskDetailDataSource_ReadSuccess tests successful read with status=1
func TestSslCertificateBindResourceTaskDetailDataSource_ReadSuccess(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	sslClient := &ssl.Client{}
	patches.ApplyMethodReturn(newMockMetaBindResourceTaskDetail().client, "UseSSLCertificateClient", sslClient)

	patches.ApplyMethodFunc(sslClient, "DescribeCertificateBindResourceTaskDetail", func(request *ssl.DescribeCertificateBindResourceTaskDetailRequest) (*ssl.DescribeCertificateBindResourceTaskDetailResponse, error) {
		resp := ssl.NewDescribeCertificateBindResourceTaskDetailResponse()
		resp.Response = &ssl.DescribeCertificateBindResourceTaskDetailResponseParams{
			Status:     ptrUint64Bind(1),
			CacheTime:  ptrStrBind("2025-07-09 12:00:00"),
			RequestId:  ptrStrBind("fake-request-id"),
			CLB:        []*ssl.ClbInstanceList{},
			CDN:        []*ssl.CdnInstanceList{},
			WAF:        []*ssl.WafInstanceList{},
			DDOS:       []*ssl.DdosInstanceList{},
			LIVE:       []*ssl.LiveInstanceList{},
			VOD:        []*ssl.VODInstanceList{},
			TKE:        []*ssl.TkeInstanceList{},
			APIGATEWAY: []*ssl.ApiGatewayInstanceList{},
			TCB:        []*ssl.TCBInstanceList{},
			TEO:        []*ssl.TeoInstanceList{},
			TSE:        []*ssl.TSEInstanceList{},
			COS:        []*ssl.COSInstanceList{},
			TDMQ:       []*ssl.TDMQInstanceList{},
			MQTT:       []*ssl.MQTTInstanceList{},
			GAAP:       []*ssl.GAAPInstanceList{},
			SCF:        []*ssl.SCFInstanceList{},
		}
		return resp, nil
	})

	meta := newMockMetaBindResourceTaskDetail()
	res := svcssl.DataSourceTencentCloudSslCertificateBindResourceTaskDetail()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"task_id": "task-12345",
	})

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.NotEmpty(t, d.Id())

	assert.Equal(t, 1, d.Get("status"))
	assert.Equal(t, "2025-07-09 12:00:00", d.Get("cache_time"))
}

// TestSslCertificateBindResourceTaskDetailDataSource_ReadWithData tests read with actual resource data
func TestSslCertificateBindResourceTaskDetailDataSource_ReadWithData(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	sslClient := &ssl.Client{}
	patches.ApplyMethodReturn(newMockMetaBindResourceTaskDetail().client, "UseSSLCertificateClient", sslClient)

	patches.ApplyMethodFunc(sslClient, "DescribeCertificateBindResourceTaskDetail", func(request *ssl.DescribeCertificateBindResourceTaskDetailRequest) (*ssl.DescribeCertificateBindResourceTaskDetailResponse, error) {
		resp := ssl.NewDescribeCertificateBindResourceTaskDetailResponse()
		resp.Response = &ssl.DescribeCertificateBindResourceTaskDetailResponseParams{
			Status:    ptrUint64Bind(1),
			CacheTime: ptrStrBind("2025-07-09 12:00:00"),
			RequestId: ptrStrBind("fake-request-id"),
			CLB: []*ssl.ClbInstanceList{
				{
					Region:     ptrStrBind("ap-guangzhou"),
					TotalCount: ptrUint64Bind(1),
					Error:      ptrStrBind(""),
					InstanceList: []*ssl.ClbInstanceDetail{
						{
							LoadBalancerId:   ptrStrBind("lb-12345"),
							LoadBalancerName: ptrStrBind("test-lb"),
							Forward:          ptrInt64Bind(1),
							Listeners: []*ssl.ClbListener{
								{
									ListenerId:   ptrStrBind("lbl-12345"),
									ListenerName: ptrStrBind("test-listener"),
									SniSwitch:    ptrUint64Bind(1),
									Protocol:     ptrStrBind("HTTPS"),
									Certificate: &ssl.Certificate{
										CertId:   ptrStrBind("cert-12345"),
										SSLMode:  ptrStrBind("UNIDIRECTIONAL"),
										DnsNames: []*string{ptrStrBind("example.com")},
									},
								},
							},
						},
					},
				},
			},
			CDN:        []*ssl.CdnInstanceList{},
			WAF:        []*ssl.WafInstanceList{},
			DDOS:       []*ssl.DdosInstanceList{},
			LIVE:       []*ssl.LiveInstanceList{},
			VOD:        []*ssl.VODInstanceList{},
			TKE:        []*ssl.TkeInstanceList{},
			APIGATEWAY: []*ssl.ApiGatewayInstanceList{},
			TCB:        []*ssl.TCBInstanceList{},
			TEO:        []*ssl.TeoInstanceList{},
			TSE:        []*ssl.TSEInstanceList{},
			COS:        []*ssl.COSInstanceList{},
			TDMQ:       []*ssl.TDMQInstanceList{},
			MQTT:       []*ssl.MQTTInstanceList{},
			GAAP:       []*ssl.GAAPInstanceList{},
			SCF:        []*ssl.SCFInstanceList{},
		}
		return resp, nil
	})

	meta := newMockMetaBindResourceTaskDetail()
	res := svcssl.DataSourceTencentCloudSslCertificateBindResourceTaskDetail()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"task_id": "task-12345",
	})

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.NotEmpty(t, d.Id())

	assert.Equal(t, 1, d.Get("status"))

	clbList := d.Get("clb").([]interface{})
	assert.Len(t, clbList, 1)
	clbMap := clbList[0].(map[string]interface{})
	assert.Equal(t, "ap-guangzhou", clbMap["region"])
	assert.Equal(t, 1, clbMap["total_count"])

	instanceList := clbMap["instance_list"].([]interface{})
	assert.Len(t, instanceList, 1)
	instanceMap := instanceList[0].(map[string]interface{})
	assert.Equal(t, "lb-12345", instanceMap["load_balancer_id"])
	assert.Equal(t, "test-lb", instanceMap["load_balancer_name"])
	assert.Equal(t, 1, instanceMap["forward"])

	listeners := instanceMap["listeners"].([]interface{})
	assert.Len(t, listeners, 1)
	listenerMap := listeners[0].(map[string]interface{})
	assert.Equal(t, "lbl-12345", listenerMap["listener_id"])
	assert.Equal(t, "HTTPS", listenerMap["protocol"])
}

// TestSslCertificateBindResourceTaskDetailDataSource_EmptyResponse tests nil response returns error
func TestSslCertificateBindResourceTaskDetailDataSource_EmptyResponse(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	sslClient := &ssl.Client{}
	patches.ApplyMethodReturn(newMockMetaBindResourceTaskDetail().client, "UseSSLCertificateClient", sslClient)

	patches.ApplyMethodFunc(sslClient, "DescribeCertificateBindResourceTaskDetail", func(request *ssl.DescribeCertificateBindResourceTaskDetailRequest) (*ssl.DescribeCertificateBindResourceTaskDetailResponse, error) {
		resp := ssl.NewDescribeCertificateBindResourceTaskDetailResponse()
		resp.Response = nil
		return resp, nil
	})

	meta := newMockMetaBindResourceTaskDetail()
	res := svcssl.DataSourceTencentCloudSslCertificateBindResourceTaskDetail()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"task_id": "task-12345",
	})

	err := res.Read(d, meta)
	assert.Error(t, err)
}

// TestSslCertificateBindResourceTaskDetailDataSource_Schema validates schema definition
func TestSslCertificateBindResourceTaskDetailDataSource_Schema(t *testing.T) {
	res := svcssl.DataSourceTencentCloudSslCertificateBindResourceTaskDetail()

	assert.NotNil(t, res)
	assert.NotNil(t, res.Read)

	assert.Contains(t, res.Schema, "task_id")
	assert.Contains(t, res.Schema, "resource_types")
	assert.Contains(t, res.Schema, "regions")
	assert.Contains(t, res.Schema, "result_output_file")
	assert.Contains(t, res.Schema, "status")
	assert.Contains(t, res.Schema, "cache_time")
	assert.Contains(t, res.Schema, "clb")
	assert.Contains(t, res.Schema, "cdn")
	assert.Contains(t, res.Schema, "waf")
	assert.Contains(t, res.Schema, "ddos")
	assert.Contains(t, res.Schema, "live")
	assert.Contains(t, res.Schema, "vod")
	assert.Contains(t, res.Schema, "tke")
	assert.Contains(t, res.Schema, "apigateway")
	assert.Contains(t, res.Schema, "tcb")
	assert.Contains(t, res.Schema, "teo")
	assert.Contains(t, res.Schema, "tse")
	assert.Contains(t, res.Schema, "cos")
	assert.Contains(t, res.Schema, "tdmq")
	assert.Contains(t, res.Schema, "mqtt")
	assert.Contains(t, res.Schema, "gaap")
	assert.Contains(t, res.Schema, "scf")

	taskId := res.Schema["task_id"]
	assert.Equal(t, schema.TypeString, taskId.Type)
	assert.True(t, taskId.Required)

	resourceTypes := res.Schema["resource_types"]
	assert.Equal(t, schema.TypeSet, resourceTypes.Type)
	assert.True(t, resourceTypes.Optional)

	status := res.Schema["status"]
	assert.Equal(t, schema.TypeInt, status.Type)
	assert.True(t, status.Computed)

	clb := res.Schema["clb"]
	assert.Equal(t, schema.TypeList, clb.Type)
	assert.True(t, clb.Computed)
}
