package tke_test

import (
	"context"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	tkev20180525 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tke/v20180525"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/tke"
)

type mockMetaRollOutSequence struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMetaRollOutSequence) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockMetaRollOutSequence{}

func newMockMetaRollOutSequence() *mockMetaRollOutSequence {
	return &mockMetaRollOutSequence{client: &connectivity.TencentCloudClient{}}
}

func ptrStringRollOut(s string) *string {
	return &s
}

func ptrInt64RollOut(v int64) *int64 {
	return &v
}

func ptrBoolRollOut(v bool) *bool {
	return &v
}

// go test ./tencentcloud/services/tke/ -run "TestRollOutSequence" -v -count=1 -gcflags="all=-l"

// TestRollOutSequence_Create tests the Create operation
func TestRollOutSequence_Create(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	tkeClient := &tkev20180525.Client{}
	patches.ApplyMethodReturn(newMockMetaRollOutSequence().client, "UseTkeClient", tkeClient)

	patches.ApplyMethodFunc(tkeClient, "CreateRollOutSequenceWithContext", func(ctx context.Context, request *tkev20180525.CreateRollOutSequenceRequest) (*tkev20180525.CreateRollOutSequenceResponse, error) {
		assert.Equal(t, "test-sequence", *request.Name)
		assert.True(t, *request.Enabled)
		assert.Len(t, request.SequenceFlows, 1)
		assert.Len(t, request.SequenceFlows[0].Tags, 1)
		assert.Equal(t, "env", *request.SequenceFlows[0].Tags[0].Key)
		assert.Equal(t, "production", *request.SequenceFlows[0].Tags[0].Value[0])
		assert.Equal(t, int64(3600), *request.SequenceFlows[0].SoakTime)

		resp := tkev20180525.NewCreateRollOutSequenceResponse()
		resp.Response = &tkev20180525.CreateRollOutSequenceResponseParams{
			ID:        ptrInt64RollOut(123),
			RequestId: ptrStringRollOut("fake-request-id"),
		}
		return resp, nil
	})

	// Mock DescribeRollOutSequencesWithContext for the Read call after Create
	patches.ApplyMethodFunc(tkeClient, "DescribeRollOutSequencesWithContext", func(ctx context.Context, request *tkev20180525.DescribeRollOutSequencesRequest) (*tkev20180525.DescribeRollOutSequencesResponse, error) {
		resp := tkev20180525.NewDescribeRollOutSequencesResponse()
		resp.Response = &tkev20180525.DescribeRollOutSequencesResponseParams{
			Sequences: []*tkev20180525.RollOutSequence{
				{
					ID:      ptrInt64RollOut(123),
					Name:    ptrStringRollOut("test-sequence"),
					Enabled: ptrBoolRollOut(true),
					SequenceFlows: []*tkev20180525.SequenceFlow{
						{
							Tags: []*tkev20180525.SequenceTag{
								{
									Key:   ptrStringRollOut("env"),
									Value: []*string{ptrStringRollOut("production")},
								},
							},
							SoakTime: ptrInt64RollOut(3600),
						},
					},
				},
			},
			TotalCount: ptrInt64RollOut(1),
			RequestId:  ptrStringRollOut("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaRollOutSequence()
	res := tke.ResourceTencentCloudKubernetesRollOutSequence()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"name":    "test-sequence",
		"enabled": true,
		"sequence_flows": []interface{}{
			map[string]interface{}{
				"tags": []interface{}{
					map[string]interface{}{
						"key":   "env",
						"value": []interface{}{"production"},
					},
				},
				"soak_time": 3600,
			},
		},
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "123", d.Id())
}

// TestRollOutSequence_Read tests the Read operation
func TestRollOutSequence_Read(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	tkeClient := &tkev20180525.Client{}
	patches.ApplyMethodReturn(newMockMetaRollOutSequence().client, "UseTkeClient", tkeClient)

	patches.ApplyMethodFunc(tkeClient, "DescribeRollOutSequencesWithContext", func(ctx context.Context, request *tkev20180525.DescribeRollOutSequencesRequest) (*tkev20180525.DescribeRollOutSequencesResponse, error) {
		resp := tkev20180525.NewDescribeRollOutSequencesResponse()
		resp.Response = &tkev20180525.DescribeRollOutSequencesResponseParams{
			Sequences: []*tkev20180525.RollOutSequence{
				{
					ID:      ptrInt64RollOut(456),
					Name:    ptrStringRollOut("read-sequence"),
					Enabled: ptrBoolRollOut(false),
					SequenceFlows: []*tkev20180525.SequenceFlow{
						{
							Tags: []*tkev20180525.SequenceTag{
								{
									Key:   ptrStringRollOut("cluster"),
									Value: []*string{ptrStringRollOut("cls-abc"), ptrStringRollOut("cls-def")},
								},
							},
							SoakTime: ptrInt64RollOut(1800),
						},
					},
				},
			},
			TotalCount: ptrInt64RollOut(1),
			RequestId:  ptrStringRollOut("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaRollOutSequence()
	res := tke.ResourceTencentCloudKubernetesRollOutSequence()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"name":           "",
		"enabled":        false,
		"sequence_flows": []interface{}{},
	})
	d.SetId("456")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "456", d.Id())
	assert.Equal(t, "read-sequence", d.Get("name").(string))
	assert.Equal(t, false, d.Get("enabled").(bool))
}

// TestRollOutSequence_Read_NotFound tests Read when sequence is not found
func TestRollOutSequence_Read_NotFound(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	tkeClient := &tkev20180525.Client{}
	patches.ApplyMethodReturn(newMockMetaRollOutSequence().client, "UseTkeClient", tkeClient)

	patches.ApplyMethodFunc(tkeClient, "DescribeRollOutSequencesWithContext", func(ctx context.Context, request *tkev20180525.DescribeRollOutSequencesRequest) (*tkev20180525.DescribeRollOutSequencesResponse, error) {
		resp := tkev20180525.NewDescribeRollOutSequencesResponse()
		resp.Response = &tkev20180525.DescribeRollOutSequencesResponseParams{
			Sequences:  []*tkev20180525.RollOutSequence{},
			TotalCount: ptrInt64RollOut(0),
			RequestId:  ptrStringRollOut("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaRollOutSequence()
	res := tke.ResourceTencentCloudKubernetesRollOutSequence()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"name":           "",
		"enabled":        false,
		"sequence_flows": []interface{}{},
	})
	d.SetId("999")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "", d.Id())
}

// TestRollOutSequence_Update tests the Update operation
func TestRollOutSequence_Update(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	tkeClient := &tkev20180525.Client{}
	patches.ApplyMethodReturn(newMockMetaRollOutSequence().client, "UseTkeClient", tkeClient)

	patches.ApplyMethodFunc(tkeClient, "ModifyRollOutSequenceWithContext", func(ctx context.Context, request *tkev20180525.ModifyRollOutSequenceRequest) (*tkev20180525.ModifyRollOutSequenceResponse, error) {
		assert.Equal(t, int64(789), *request.ID)
		assert.Equal(t, "updated-sequence", *request.Name)
		assert.False(t, *request.Enabled)

		resp := tkev20180525.NewModifyRollOutSequenceResponse()
		resp.Response = &tkev20180525.ModifyRollOutSequenceResponseParams{
			RequestId: ptrStringRollOut("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(tkeClient, "DescribeRollOutSequencesWithContext", func(ctx context.Context, request *tkev20180525.DescribeRollOutSequencesRequest) (*tkev20180525.DescribeRollOutSequencesResponse, error) {
		resp := tkev20180525.NewDescribeRollOutSequencesResponse()
		resp.Response = &tkev20180525.DescribeRollOutSequencesResponseParams{
			Sequences: []*tkev20180525.RollOutSequence{
				{
					ID:      ptrInt64RollOut(789),
					Name:    ptrStringRollOut("updated-sequence"),
					Enabled: ptrBoolRollOut(false),
					SequenceFlows: []*tkev20180525.SequenceFlow{
						{
							Tags: []*tkev20180525.SequenceTag{
								{
									Key:   ptrStringRollOut("env"),
									Value: []*string{ptrStringRollOut("staging")},
								},
							},
							SoakTime: ptrInt64RollOut(900),
						},
					},
				},
			},
			TotalCount: ptrInt64RollOut(1),
			RequestId:  ptrStringRollOut("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaRollOutSequence()
	res := tke.ResourceTencentCloudKubernetesRollOutSequence()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"name":    "updated-sequence",
		"enabled": false,
		"sequence_flows": []interface{}{
			map[string]interface{}{
				"tags": []interface{}{
					map[string]interface{}{
						"key":   "env",
						"value": []interface{}{"staging"},
					},
				},
				"soak_time": 900,
			},
		},
	})
	d.SetId("789")

	err := res.Update(d, meta)
	assert.NoError(t, err)
}

// TestRollOutSequence_Delete tests the Delete operation
func TestRollOutSequence_Delete(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	tkeClient := &tkev20180525.Client{}
	patches.ApplyMethodReturn(newMockMetaRollOutSequence().client, "UseTkeClient", tkeClient)

	patches.ApplyMethodFunc(tkeClient, "DeleteRollOutSequenceWithContext", func(ctx context.Context, request *tkev20180525.DeleteRollOutSequenceRequest) (*tkev20180525.DeleteRollOutSequenceResponse, error) {
		assert.Equal(t, int64(321), *request.ID)

		resp := tkev20180525.NewDeleteRollOutSequenceResponse()
		resp.Response = &tkev20180525.DeleteRollOutSequenceResponseParams{
			RequestId: ptrStringRollOut("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaRollOutSequence()
	res := tke.ResourceTencentCloudKubernetesRollOutSequence()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"name":           "delete-sequence",
		"enabled":        true,
		"sequence_flows": []interface{}{},
	})
	d.SetId("321")

	err := res.Delete(d, meta)
	assert.NoError(t, err)
}
