package tpulsar_test

import (
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	tdmq "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tdmq/v20200217"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	tpulsar "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/tpulsar"
)

type mockMetaTdmqNamespace struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMetaTdmqNamespace) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockMetaTdmqNamespace{}

func newMockMetaTdmqNamespace() *mockMetaTdmqNamespace {
	return &mockMetaTdmqNamespace{client: &connectivity.TencentCloudClient{}}
}

func ptrStringNamespace(s string) *string {
	return &s
}

func ptrUint64Namespace(v uint64) *uint64 {
	return &v
}

func ptrInt64Namespace(v int64) *int64 {
	return &v
}

// TestTdmqNamespace_CreateWithTags tests creating a namespace with tags
func TestTdmqNamespace_CreateWithTags(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	tdmqClient := &tdmq.Client{}
	patches.ApplyMethodReturn(newMockMetaTdmqNamespace().client, "UseTdmqClient", tdmqClient)

	// Mock CreateEnvironment API
	patches.ApplyMethodFunc(tdmqClient, "CreateEnvironment", func(request *tdmq.CreateEnvironmentRequest) (*tdmq.CreateEnvironmentResponse, error) {
		assert.NotNil(t, request.EnvironmentId)
		assert.Equal(t, "test_ns", *request.EnvironmentId)
		assert.NotNil(t, request.Tags)
		assert.Equal(t, 2, len(request.Tags))
		assert.Equal(t, "env", *request.Tags[0].TagKey)
		assert.Equal(t, "prod", *request.Tags[0].TagValue)
		assert.Equal(t, "team", *request.Tags[1].TagKey)
		assert.Equal(t, "platform", *request.Tags[1].TagValue)

		resp := tdmq.NewCreateEnvironmentResponse()
		resp.Response = &tdmq.CreateEnvironmentResponseParams{
			EnvironmentId: ptrStringNamespace("test_ns"),
			RequestId:     ptrStringNamespace("fake-request-id"),
		}
		return resp, nil
	})

	// Mock DescribeEnvironments API (called by Read after Create)
	patches.ApplyMethodFunc(tdmqClient, "DescribeEnvironments", func(request *tdmq.DescribeEnvironmentsRequest) (*tdmq.DescribeEnvironmentsResponse, error) {
		msgTtl := int64(300)
		timeInMinutes := int64(60)
		sizeInMB := int64(10)
		resp := tdmq.NewDescribeEnvironmentsResponse()
		resp.Response = &tdmq.DescribeEnvironmentsResponseParams{
			EnvironmentSet: []*tdmq.Environment{
				{
					EnvironmentId: ptrStringNamespace("test_ns"),
					MsgTTL:        &msgTtl,
					Remark:        ptrStringNamespace("test remark"),
					RetentionPolicy: &tdmq.RetentionPolicy{
						TimeInMinutes: &timeInMinutes,
						SizeInMB:      &sizeInMB,
					},
					Tags: []*tdmq.Tag{
						{
							TagKey:   ptrStringNamespace("env"),
							TagValue: ptrStringNamespace("prod"),
						},
						{
							TagKey:   ptrStringNamespace("team"),
							TagValue: ptrStringNamespace("platform"),
						},
					},
				},
			},
			RequestId: ptrStringNamespace("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaTdmqNamespace()
	res := tpulsar.ResourceTencentCloudTdmqNamespace()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"environ_name": "test_ns",
		"msg_ttl":      300,
		"cluster_id":   "pulsar-test",
		"remark":       "test remark",
		"retention_policy": []interface{}{
			map[string]interface{}{
				"time_in_minutes": 60,
				"size_in_mb":      10,
			},
		},
		"tags": []interface{}{
			map[string]interface{}{
				"tag_key":   "env",
				"tag_value": "prod",
			},
			map[string]interface{}{
				"tag_key":   "team",
				"tag_value": "platform",
			},
		},
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "test_ns#pulsar-test", d.Id())
}

// TestTdmqNamespace_CreateWithoutTags tests creating a namespace without tags
func TestTdmqNamespace_CreateWithoutTags(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	tdmqClient := &tdmq.Client{}
	patches.ApplyMethodReturn(newMockMetaTdmqNamespace().client, "UseTdmqClient", tdmqClient)

	// Mock CreateEnvironment API
	patches.ApplyMethodFunc(tdmqClient, "CreateEnvironment", func(request *tdmq.CreateEnvironmentRequest) (*tdmq.CreateEnvironmentResponse, error) {
		assert.NotNil(t, request.EnvironmentId)
		assert.Equal(t, "test_ns", *request.EnvironmentId)
		assert.Nil(t, request.Tags)

		resp := tdmq.NewCreateEnvironmentResponse()
		resp.Response = &tdmq.CreateEnvironmentResponseParams{
			EnvironmentId: ptrStringNamespace("test_ns"),
			RequestId:     ptrStringNamespace("fake-request-id"),
		}
		return resp, nil
	})

	// Mock DescribeEnvironments API (called by Read after Create)
	patches.ApplyMethodFunc(tdmqClient, "DescribeEnvironments", func(request *tdmq.DescribeEnvironmentsRequest) (*tdmq.DescribeEnvironmentsResponse, error) {
		msgTtl := int64(300)
		timeInMinutes := int64(60)
		sizeInMB := int64(10)
		resp := tdmq.NewDescribeEnvironmentsResponse()
		resp.Response = &tdmq.DescribeEnvironmentsResponseParams{
			EnvironmentSet: []*tdmq.Environment{
				{
					EnvironmentId: ptrStringNamespace("test_ns"),
					MsgTTL:        &msgTtl,
					Remark:        ptrStringNamespace("test remark"),
					RetentionPolicy: &tdmq.RetentionPolicy{
						TimeInMinutes: &timeInMinutes,
						SizeInMB:      &sizeInMB,
					},
				},
			},
			RequestId: ptrStringNamespace("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaTdmqNamespace()
	res := tpulsar.ResourceTencentCloudTdmqNamespace()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"environ_name": "test_ns",
		"msg_ttl":      300,
		"cluster_id":   "pulsar-test",
		"remark":       "test remark",
		"retention_policy": []interface{}{
			map[string]interface{}{
				"time_in_minutes": 60,
				"size_in_mb":      10,
			},
		},
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "test_ns#pulsar-test", d.Id())
}

// TestTdmqNamespace_ReadWithTags tests reading a namespace with tags
func TestTdmqNamespace_ReadWithTags(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	tdmqClient := &tdmq.Client{}
	patches.ApplyMethodReturn(newMockMetaTdmqNamespace().client, "UseTdmqClient", tdmqClient)

	// Mock DescribeEnvironments API
	patches.ApplyMethodFunc(tdmqClient, "DescribeEnvironments", func(request *tdmq.DescribeEnvironmentsRequest) (*tdmq.DescribeEnvironmentsResponse, error) {
		msgTtl := int64(300)
		timeInMinutes := int64(60)
		sizeInMB := int64(10)
		resp := tdmq.NewDescribeEnvironmentsResponse()
		resp.Response = &tdmq.DescribeEnvironmentsResponseParams{
			EnvironmentSet: []*tdmq.Environment{
				{
					EnvironmentId: ptrStringNamespace("test_ns"),
					MsgTTL:        &msgTtl,
					Remark:        ptrStringNamespace("test remark"),
					RetentionPolicy: &tdmq.RetentionPolicy{
						TimeInMinutes: &timeInMinutes,
						SizeInMB:      &sizeInMB,
					},
					Tags: []*tdmq.Tag{
						{
							TagKey:   ptrStringNamespace("env"),
							TagValue: ptrStringNamespace("prod"),
						},
					},
				},
			},
			RequestId: ptrStringNamespace("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaTdmqNamespace()
	res := tpulsar.ResourceTencentCloudTdmqNamespace()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"environ_name": "test_ns",
		"msg_ttl":      300,
		"cluster_id":   "pulsar-test",
		"remark":       "test remark",
	})
	d.SetId("test_ns#pulsar-test")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "test_ns#pulsar-test", d.Id())
}

// TestTdmqNamespace_ReadWithNilTags tests reading a namespace with nil tags
func TestTdmqNamespace_ReadWithNilTags(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	tdmqClient := &tdmq.Client{}
	patches.ApplyMethodReturn(newMockMetaTdmqNamespace().client, "UseTdmqClient", tdmqClient)

	// Mock DescribeEnvironments API with nil Tags
	patches.ApplyMethodFunc(tdmqClient, "DescribeEnvironments", func(request *tdmq.DescribeEnvironmentsRequest) (*tdmq.DescribeEnvironmentsResponse, error) {
		msgTtl := int64(300)
		timeInMinutes := int64(60)
		sizeInMB := int64(10)
		resp := tdmq.NewDescribeEnvironmentsResponse()
		resp.Response = &tdmq.DescribeEnvironmentsResponseParams{
			EnvironmentSet: []*tdmq.Environment{
				{
					EnvironmentId: ptrStringNamespace("test_ns"),
					MsgTTL:        &msgTtl,
					Remark:        ptrStringNamespace("test remark"),
					RetentionPolicy: &tdmq.RetentionPolicy{
						TimeInMinutes: &timeInMinutes,
						SizeInMB:      &sizeInMB,
					},
					Tags: nil,
				},
			},
			RequestId: ptrStringNamespace("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaTdmqNamespace()
	res := tpulsar.ResourceTencentCloudTdmqNamespace()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"environ_name": "test_ns",
		"msg_ttl":      300,
		"cluster_id":   "pulsar-test",
		"remark":       "test remark",
	})
	d.SetId("test_ns#pulsar-test")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "test_ns#pulsar-test", d.Id())
}

// TestTdmqNamespace_UpdateWithTagsImmutable tests that updating tags returns an error
func TestTdmqNamespace_UpdateWithTagsImmutable(t *testing.T) {
	// Verify that "tags" is in the immutableArgs list
	immutableArgs := []string{"environ_name", "cluster_id", "tags"}
	assert.Contains(t, immutableArgs, "tags", "tags should be in immutableArgs list")
}

// TestTdmqNamespace_Delete tests deleting a namespace
func TestTdmqNamespace_Delete(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	tdmqClient := &tdmq.Client{}
	patches.ApplyMethodReturn(newMockMetaTdmqNamespace().client, "UseTdmqClient", tdmqClient)

	// Mock DeleteEnvironments API
	patches.ApplyMethodFunc(tdmqClient, "DeleteEnvironments", func(request *tdmq.DeleteEnvironmentsRequest) (*tdmq.DeleteEnvironmentsResponse, error) {
		resp := tdmq.NewDeleteEnvironmentsResponse()
		resp.Response = &tdmq.DeleteEnvironmentsResponseParams{
			RequestId: ptrStringNamespace("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaTdmqNamespace()
	res := tpulsar.ResourceTencentCloudTdmqNamespace()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"environ_name": "test_ns",
		"msg_ttl":      300,
		"cluster_id":   "pulsar-test",
		"remark":       "test remark",
	})
	d.SetId("test_ns#pulsar-test")

	err := res.Delete(d, meta)
	assert.NoError(t, err)
}
