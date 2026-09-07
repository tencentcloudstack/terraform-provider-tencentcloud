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

type mockMetaForSplunkDeliver struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMetaForSplunkDeliver) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockMetaForSplunkDeliver{}

func newMockMetaForSplunkDeliver() *mockMetaForSplunkDeliver {
	return &mockMetaForSplunkDeliver{client: &connectivity.TencentCloudClient{}}
}

func ptrStrSplunkDeliver(s string) *string {
	return &s
}

func ptrUint64SplunkDeliver(v uint64) *uint64 {
	return &v
}

func ptrInt64SplunkDeliver(v int64) *int64 {
	return &v
}

func ptrBoolSplunkDeliver(v bool) *bool {
	return &v
}

// TestSplunkDeliver_Create tests the Create operation
func TestSplunkDeliver_Create(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	clsClient := &cls.Client{}
	patches.ApplyMethodReturn(newMockMetaForSplunkDeliver().client, "UseClsClient", clsClient)

	var capturedRequest *cls.CreateSplunkDeliverRequest
	patches.ApplyMethodFunc(clsClient, "CreateSplunkDeliver", func(request *cls.CreateSplunkDeliverRequest) (*cls.CreateSplunkDeliverResponse, error) {
		capturedRequest = request
		resp := cls.NewCreateSplunkDeliverResponse()
		resp.Response = &cls.CreateSplunkDeliverResponseParams{
			TaskId:    ptrStrSplunkDeliver("task-test-123"),
			RequestId: ptrStrSplunkDeliver("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(clsClient, "DescribeSplunkDelivers", func(request *cls.DescribeSplunkDeliversRequest) (*cls.DescribeSplunkDeliversResponse, error) {
		resp := cls.NewDescribeSplunkDeliversResponse()
		resp.Response = &cls.DescribeSplunkDeliversResponseParams{
			Infos: []*cls.SplunkDeliverInfo{
				{
					TaskId:  ptrStrSplunkDeliver("task-test-123"),
					TopicId: ptrStrSplunkDeliver("topic-test-123"),
					Name:    ptrStrSplunkDeliver("test-deliver"),
					Enable:  ptrInt64SplunkDeliver(1),
					NetInfo: &cls.NetInfo{
						Host:    ptrStrSplunkDeliver("10.0.0.1"),
						Port:    ptrUint64SplunkDeliver(8088),
						Token:   ptrStrSplunkDeliver("test-token"),
						NetType: ptrUint64SplunkDeliver(2),
						IsSSL:   ptrBoolSplunkDeliver(true),
					},
					Metadata: &cls.MetadataInfo{
						Format:     ptrStrSplunkDeliver("json"),
						MetaFields: []*string{ptrStrSplunkDeliver("__SOURCE__"), ptrStrSplunkDeliver("__FILENAME__")},
						EnableTag:  ptrBoolSplunkDeliver(true),
					},
					IndexAck: ptrInt64SplunkDeliver(1),
				},
			},
			Total:     ptrUint64SplunkDeliver(1),
			RequestId: ptrStrSplunkDeliver("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaForSplunkDeliver()
	res := localcls.ResourceTencentCloudClsSplunkDeliver()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"topic_id": "topic-test-123",
		"name":     "test-deliver",
		"net_info": []interface{}{
			map[string]interface{}{
				"host":     "10.0.0.1",
				"port":     8088,
				"token":    "test-token",
				"net_type": 2,
				"is_ssl":   true,
			},
		},
		"metadata_info": []interface{}{
			map[string]interface{}{
				"format": "json",
				"meta_fields": schema.NewSet(schema.HashString, []interface{}{
					"__SOURCE__",
					"__FILENAME__",
				}),
				"enable_tag": true,
			},
		},
		"index_ack": 1,
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "task-test-123#topic-test-123", d.Id())
	assert.NotNil(t, capturedRequest)
	assert.NotNil(t, capturedRequest.Name)
	assert.Equal(t, "test-deliver", *capturedRequest.Name)
	assert.NotNil(t, capturedRequest.NetInfo)
	assert.Equal(t, "10.0.0.1", *capturedRequest.NetInfo.Host)
}

// TestSplunkDeliver_Create_NilTaskId tests Create with nil TaskId
func TestSplunkDeliver_Create_NilTaskId(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	clsClient := &cls.Client{}
	patches.ApplyMethodReturn(newMockMetaForSplunkDeliver().client, "UseClsClient", clsClient)

	patches.ApplyMethodFunc(clsClient, "CreateSplunkDeliver", func(request *cls.CreateSplunkDeliverRequest) (*cls.CreateSplunkDeliverResponse, error) {
		resp := cls.NewCreateSplunkDeliverResponse()
		resp.Response = &cls.CreateSplunkDeliverResponseParams{
			TaskId:    nil,
			RequestId: ptrStrSplunkDeliver("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaForSplunkDeliver()
	res := localcls.ResourceTencentCloudClsSplunkDeliver()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"topic_id": "topic-test-123",
		"name":     "test-deliver",
		"net_info": []interface{}{
			map[string]interface{}{
				"host":     "10.0.0.1",
				"port":     8088,
				"token":    "test-token",
				"net_type": 2,
			},
		},
		"metadata_info": []interface{}{
			map[string]interface{}{
				"format": "json",
				"meta_fields": schema.NewSet(schema.HashString, []interface{}{
					"__SOURCE__",
				}),
				"enable_tag": true,
			},
		},
	})

	err := res.Create(d, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "TaskId")
}

// TestSplunkDeliver_Read tests the Read operation
func TestSplunkDeliver_Read(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	clsClient := &cls.Client{}
	patches.ApplyMethodReturn(newMockMetaForSplunkDeliver().client, "UseClsClient", clsClient)

	patches.ApplyMethodFunc(clsClient, "DescribeSplunkDelivers", func(request *cls.DescribeSplunkDeliversRequest) (*cls.DescribeSplunkDeliversResponse, error) {
		resp := cls.NewDescribeSplunkDeliversResponse()
		resp.Response = &cls.DescribeSplunkDeliversResponseParams{
			Infos: []*cls.SplunkDeliverInfo{
				{
					TaskId:  ptrStrSplunkDeliver("task-test-123"),
					TopicId: ptrStrSplunkDeliver("topic-test-123"),
					Name:    ptrStrSplunkDeliver("test-deliver"),
					Enable:  ptrInt64SplunkDeliver(1),
					NetInfo: &cls.NetInfo{
						Host:    ptrStrSplunkDeliver("10.0.0.1"),
						Port:    ptrUint64SplunkDeliver(8088),
						Token:   ptrStrSplunkDeliver("test-token"),
						NetType: ptrUint64SplunkDeliver(2),
						IsSSL:   ptrBoolSplunkDeliver(true),
					},
					Metadata: &cls.MetadataInfo{
						Format:     ptrStrSplunkDeliver("json"),
						MetaFields: []*string{ptrStrSplunkDeliver("__SOURCE__"), ptrStrSplunkDeliver("__FILENAME__")},
						EnableTag:  ptrBoolSplunkDeliver(true),
					},
					HasServiceLog: ptrInt64SplunkDeliver(2),
					IndexAck:      ptrInt64SplunkDeliver(1),
					Source:        ptrStrSplunkDeliver("test-source"),
					SourceType:    ptrStrSplunkDeliver("test-source-type"),
					Index:         ptrStrSplunkDeliver("main"),
					Channel:       ptrStrSplunkDeliver("test-channel"),
					DSLFilter:     ptrStrSplunkDeliver(""),
				},
			},
			Total:     ptrUint64SplunkDeliver(1),
			RequestId: ptrStrSplunkDeliver("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaForSplunkDeliver()
	res := localcls.ResourceTencentCloudClsSplunkDeliver()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"topic_id": "topic-test-123",
		"name":     "test-deliver",
	})
	d.SetId("task-test-123#topic-test-123")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "test-deliver", d.Get("name"))
	assert.Equal(t, "main", d.Get("index"))
	assert.Equal(t, "test-source", d.Get("source"))
	assert.Equal(t, "test-source-type", d.Get("source_type"))
	assert.Equal(t, "test-channel", d.Get("channel"))
}

// TestSplunkDeliver_Read_EmptyInfos tests Read with empty Infos
func TestSplunkDeliver_Read_EmptyInfos(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	clsClient := &cls.Client{}
	patches.ApplyMethodReturn(newMockMetaForSplunkDeliver().client, "UseClsClient", clsClient)

	patches.ApplyMethodFunc(clsClient, "DescribeSplunkDelivers", func(request *cls.DescribeSplunkDeliversRequest) (*cls.DescribeSplunkDeliversResponse, error) {
		resp := cls.NewDescribeSplunkDeliversResponse()
		resp.Response = &cls.DescribeSplunkDeliversResponseParams{
			Infos:     []*cls.SplunkDeliverInfo{},
			Total:     ptrUint64SplunkDeliver(0),
			RequestId: ptrStrSplunkDeliver("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaForSplunkDeliver()
	res := localcls.ResourceTencentCloudClsSplunkDeliver()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"topic_id": "topic-test-123",
		"name":     "test-deliver",
	})
	d.SetId("task-test-123#topic-test-123")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "", d.Id())
}

// TestSplunkDeliver_Update tests the Update operation
func TestSplunkDeliver_Update(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	clsClient := &cls.Client{}
	patches.ApplyMethodReturn(newMockMetaForSplunkDeliver().client, "UseClsClient", clsClient)

	var capturedModifyRequest *cls.ModifySplunkDeliverRequest
	patches.ApplyMethodFunc(clsClient, "ModifySplunkDeliver", func(request *cls.ModifySplunkDeliverRequest) (*cls.ModifySplunkDeliverResponse, error) {
		capturedModifyRequest = request
		resp := cls.NewModifySplunkDeliverResponse()
		resp.Response = &cls.ModifySplunkDeliverResponseParams{
			RequestId: ptrStrSplunkDeliver("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(clsClient, "DescribeSplunkDelivers", func(request *cls.DescribeSplunkDeliversRequest) (*cls.DescribeSplunkDeliversResponse, error) {
		resp := cls.NewDescribeSplunkDeliversResponse()
		resp.Response = &cls.DescribeSplunkDeliversResponseParams{
			Infos: []*cls.SplunkDeliverInfo{
				{
					TaskId:  ptrStrSplunkDeliver("task-test-123"),
					TopicId: ptrStrSplunkDeliver("topic-test-123"),
					Name:    ptrStrSplunkDeliver("updated-name"),
					Enable:  ptrInt64SplunkDeliver(0),
					NetInfo: &cls.NetInfo{
						Host:    ptrStrSplunkDeliver("10.0.0.2"),
						Port:    ptrUint64SplunkDeliver(8089),
						Token:   ptrStrSplunkDeliver("new-token"),
						NetType: ptrUint64SplunkDeliver(2),
					},
					Metadata: &cls.MetadataInfo{
						Format:     ptrStrSplunkDeliver("rawlog"),
						MetaFields: []*string{ptrStrSplunkDeliver("__SOURCE__")},
						EnableTag:  ptrBoolSplunkDeliver(false),
					},
				},
			},
			Total:     ptrUint64SplunkDeliver(1),
			RequestId: ptrStrSplunkDeliver("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaForSplunkDeliver()
	res := localcls.ResourceTencentCloudClsSplunkDeliver()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"topic_id": "topic-test-123",
		"name":     "test-deliver",
		"net_info": []interface{}{
			map[string]interface{}{
				"host":     "10.0.0.1",
				"port":     8088,
				"token":    "test-token",
				"net_type": 2,
			},
		},
		"metadata_info": []interface{}{
			map[string]interface{}{
				"format": "json",
				"meta_fields": schema.NewSet(schema.HashString, []interface{}{
					"__SOURCE__",
				}),
				"enable_tag": true,
			},
		},
	})
	d.SetId("task-test-123#topic-test-123")

	// Simulate name change
	_ = d.Set("name", "updated-name")

	err := res.Update(d, meta)
	assert.NoError(t, err)
	assert.NotNil(t, capturedModifyRequest)
	assert.NotNil(t, capturedModifyRequest.Name)
	assert.Equal(t, "updated-name", *capturedModifyRequest.Name)
}

// TestSplunkDeliver_Delete tests the Delete operation
func TestSplunkDeliver_Delete(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	clsClient := &cls.Client{}
	patches.ApplyMethodReturn(newMockMetaForSplunkDeliver().client, "UseClsClient", clsClient)

	var capturedDeleteRequest *cls.DeleteSplunkDeliverRequest
	patches.ApplyMethodFunc(clsClient, "DeleteSplunkDeliver", func(request *cls.DeleteSplunkDeliverRequest) (*cls.DeleteSplunkDeliverResponse, error) {
		capturedDeleteRequest = request
		resp := cls.NewDeleteSplunkDeliverResponse()
		resp.Response = &cls.DeleteSplunkDeliverResponseParams{
			RequestId: ptrStrSplunkDeliver("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaForSplunkDeliver()
	res := localcls.ResourceTencentCloudClsSplunkDeliver()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"topic_id": "topic-test-123",
		"name":     "test-deliver",
	})
	d.SetId("task-test-123#topic-test-123")

	err := res.Delete(d, meta)
	assert.NoError(t, err)
	assert.NotNil(t, capturedDeleteRequest)
	assert.Equal(t, "task-test-123", *capturedDeleteRequest.TaskId)
	assert.Equal(t, "topic-test-123", *capturedDeleteRequest.TopicId)
}

// TestSplunkDeliver_Schema tests the schema definition
func TestSplunkDeliver_Schema(t *testing.T) {
	res := localcls.ResourceTencentCloudClsSplunkDeliver()

	assert.NotNil(t, res)
	assert.Contains(t, res.Schema, "topic_id")
	assert.Contains(t, res.Schema, "name")
	assert.Contains(t, res.Schema, "net_info")
	assert.Contains(t, res.Schema, "metadata_info")
	assert.Contains(t, res.Schema, "external_role")
	assert.Contains(t, res.Schema, "enable")
	assert.Contains(t, res.Schema, "task_id")

	topicIdSchema := res.Schema["topic_id"]
	assert.True(t, topicIdSchema.Required)
	assert.True(t, topicIdSchema.ForceNew)

	nameSchema := res.Schema["name"]
	assert.True(t, nameSchema.Required)

	enableSchema := res.Schema["enable"]
	assert.True(t, enableSchema.Optional)
	assert.True(t, enableSchema.Computed)

	taskIdSchema := res.Schema["task_id"]
	assert.True(t, taskIdSchema.Computed)
}

// TestSplunkDeliver_Import tests the import configuration
func TestSplunkDeliver_Import(t *testing.T) {
	res := localcls.ResourceTencentCloudClsSplunkDeliver()

	assert.NotNil(t, res.Importer)
	assert.NotNil(t, res.Importer.State)
}
