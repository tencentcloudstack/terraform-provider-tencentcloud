package cls_test

import (
	"context"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	clsv20201016 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cls/v20201016"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/cls"
)

type mockMetaForClsCloudProductLogTaskV2 struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMetaForClsCloudProductLogTaskV2) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockMetaForClsCloudProductLogTaskV2{}

func newMockMetaForClsCloudProductLogTaskV2() *mockMetaForClsCloudProductLogTaskV2 {
	return &mockMetaForClsCloudProductLogTaskV2{client: &connectivity.TencentCloudClient{}}
}

func ptrStringCPLT2(s string) *string { return &s }
func ptrInt64CPLT2(v int64) *int64    { return &v }
func ptrUint64CPLT2(v uint64) *uint64 { return &v }

// buildCloudProductLogTaskResp builds a DescribeCloudProductLogTasks response with the
// given task info, used both by the state-refresh polling in Create and by Read.
func buildCloudProductLogTaskResp(task *clsv20201016.CloudProductLogTaskInfo) *clsv20201016.DescribeCloudProductLogTasksResponse {
	resp := clsv20201016.NewDescribeCloudProductLogTasksResponse()
	resp.Response = &clsv20201016.DescribeCloudProductLogTasksResponseParams{
		Tasks:      []*clsv20201016.CloudProductLogTaskInfo{task},
		TotalCount: ptrUint64CPLT2(1),
		RequestId:  ptrStringCPLT2("fake-request-id"),
	}
	return resp
}

// patchDescribeCloudProductLogTasks mocks the read-back API used by both the
// Create state-refresh poller and the Read function.
func patchDescribeCloudProductLogTasks(patches *gomonkey.Patches, clsClient *clsv20201016.Client, task *clsv20201016.CloudProductLogTaskInfo) {
	patches.ApplyMethodFunc(clsClient, "DescribeCloudProductLogTasks", func(request *clsv20201016.DescribeCloudProductLogTasksRequest) (*clsv20201016.DescribeCloudProductLogTasksResponse, error) {
		return buildCloudProductLogTaskResp(task), nil
	})
}

// TestClsCloudProductLogTaskV2_Create_WithTags verifies that when the `tags`
// parameter is set, the create request carries the populated Tags slice.
func TestClsCloudProductLogTaskV2_Create_WithTags(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	clsClient := &clsv20201016.Client{}
	patches.ApplyMethodReturn(newMockMetaForClsCloudProductLogTaskV2().client, "UseClsV20201016Client", clsClient)
	patches.ApplyMethodReturn(newMockMetaForClsCloudProductLogTaskV2().client, "UseClsClient", clsClient)

	var capturedRequest *clsv20201016.CreateCloudProductLogCollectionRequest
	patches.ApplyMethodFunc(clsClient, "CreateCloudProductLogCollectionWithContext", func(_ context.Context, request *clsv20201016.CreateCloudProductLogCollectionRequest) (*clsv20201016.CreateCloudProductLogCollectionResponse, error) {
		capturedRequest = request
		resp := clsv20201016.NewCreateCloudProductLogCollectionResponse()
		resp.Response = &clsv20201016.CreateCloudProductLogCollectionResponseParams{
			RequestId: ptrStringCPLT2("fake-request-id"),
		}
		return resp, nil
	})

	task := &clsv20201016.CloudProductLogTaskInfo{
		InstanceId: ptrStringCPLT2("postgres-0an6hpv3"),
		ClsRegion:  ptrStringCPLT2("ap-guangzhou"),
		LogType:    ptrStringCPLT2("PostgreSQL-SLOW"),
		Status:     ptrInt64CPLT2(1),
		TopicTags: []*clsv20201016.Tag{
			{Key: ptrStringCPLT2("Environment"), Value: ptrStringCPLT2("production")},
			{Key: ptrStringCPLT2("Team"), Value: ptrStringCPLT2("backend")},
		},
	}
	patchDescribeCloudProductLogTasks(patches, clsClient, task)

	meta := newMockMetaForClsCloudProductLogTaskV2()
	res := cls.ResourceTencentCloudClsCloudProductLogTaskV2()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"instance_id":          "postgres-0an6hpv3",
		"assumer_name":         "PostgreSQL",
		"log_type":             "PostgreSQL-SLOW",
		"cloud_product_region": "gz",
		"cls_region":           "ap-guangzhou",
		"logset_name":          "tf-example",
		"topic_name":           "tf-example",
		"tags": map[string]interface{}{
			"Environment": "production",
			"Team":        "backend",
		},
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.NotEmpty(t, d.Id())

	// Verify the Tags field was populated on the create request.
	assert.NotNil(t, capturedRequest.Tags)
	assert.Equal(t, 2, len(capturedRequest.Tags))
	tagMap := make(map[string]string)
	for _, tag := range capturedRequest.Tags {
		tagMap[*tag.Key] = *tag.Value
	}
	assert.Equal(t, "production", tagMap["Environment"])
	assert.Equal(t, "backend", tagMap["Team"])

	// Verify tags are read back into state from the (mocked) read API.
	tagsInState := d.Get("tags").(map[string]interface{})
	assert.Equal(t, "production", tagsInState["Environment"])
	assert.Equal(t, "backend", tagsInState["Team"])
}

// TestClsCloudProductLogTaskV2_Create_WithoutTags verifies that when the
// `tags` parameter is not set, the create request Tags field is left nil.
func TestClsCloudProductLogTaskV2_Create_WithoutTags(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	clsClient := &clsv20201016.Client{}
	patches.ApplyMethodReturn(newMockMetaForClsCloudProductLogTaskV2().client, "UseClsV20201016Client", clsClient)
	patches.ApplyMethodReturn(newMockMetaForClsCloudProductLogTaskV2().client, "UseClsClient", clsClient)

	var capturedRequest *clsv20201016.CreateCloudProductLogCollectionRequest
	patches.ApplyMethodFunc(clsClient, "CreateCloudProductLogCollectionWithContext", func(_ context.Context, request *clsv20201016.CreateCloudProductLogCollectionRequest) (*clsv20201016.CreateCloudProductLogCollectionResponse, error) {
		capturedRequest = request
		resp := clsv20201016.NewCreateCloudProductLogCollectionResponse()
		resp.Response = &clsv20201016.CreateCloudProductLogCollectionResponseParams{
			RequestId: ptrStringCPLT2("fake-request-id"),
		}
		return resp, nil
	})

	task := &clsv20201016.CloudProductLogTaskInfo{
		InstanceId: ptrStringCPLT2("postgres-0an6hpv3"),
		ClsRegion:  ptrStringCPLT2("ap-guangzhou"),
		LogType:    ptrStringCPLT2("PostgreSQL-SLOW"),
		Status:     ptrInt64CPLT2(1),
	}
	patchDescribeCloudProductLogTasks(patches, clsClient, task)

	meta := newMockMetaForClsCloudProductLogTaskV2()
	res := cls.ResourceTencentCloudClsCloudProductLogTaskV2()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"instance_id":          "postgres-0an6hpv3",
		"assumer_name":         "PostgreSQL",
		"log_type":             "PostgreSQL-SLOW",
		"cloud_product_region": "gz",
		"cls_region":           "ap-guangzhou",
		"logset_name":          "tf-example",
		"topic_name":           "tf-example",
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.NotEmpty(t, d.Id())
	// Tags should not be populated on the request when not configured.
	assert.Nil(t, capturedRequest.Tags)
}

// TestClsCloudProductLogTaskV2_Update_Tags verifies that when the `tags`
// parameter changes, ModifyCloudProductLogCollection is called with the
// updated Tags value and the resource is not recreated.
func TestClsCloudProductLogTaskV2_Update_Tags(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	clsClient := &clsv20201016.Client{}
	patches.ApplyMethodReturn(newMockMetaForClsCloudProductLogTaskV2().client, "UseClsV20201016Client", clsClient)
	patches.ApplyMethodReturn(newMockMetaForClsCloudProductLogTaskV2().client, "UseClsClient", clsClient)

	var capturedRequest *clsv20201016.ModifyCloudProductLogCollectionRequest
	patches.ApplyMethodFunc(clsClient, "ModifyCloudProductLogCollectionWithContext", func(_ context.Context, request *clsv20201016.ModifyCloudProductLogCollectionRequest) (*clsv20201016.ModifyCloudProductLogCollectionResponse, error) {
		capturedRequest = request
		resp := clsv20201016.NewModifyCloudProductLogCollectionResponse()
		resp.Response = &clsv20201016.ModifyCloudProductLogCollectionResponseParams{
			RequestId: ptrStringCPLT2("fake-request-id-update"),
		}
		return resp, nil
	})

	task := &clsv20201016.CloudProductLogTaskInfo{
		InstanceId: ptrStringCPLT2("postgres-0an6hpv3"),
		ClsRegion:  ptrStringCPLT2("ap-guangzhou"),
		LogType:    ptrStringCPLT2("PostgreSQL-SLOW"),
		Status:     ptrInt64CPLT2(1),
		TopicTags: []*clsv20201016.Tag{
			{Key: ptrStringCPLT2("Environment"), Value: ptrStringCPLT2("staging")},
		},
	}
	patchDescribeCloudProductLogTasks(patches, clsClient, task)

	meta := newMockMetaForClsCloudProductLogTaskV2()
	res := cls.ResourceTencentCloudClsCloudProductLogTaskV2()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"instance_id":          "postgres-0an6hpv3",
		"assumer_name":         "PostgreSQL",
		"log_type":             "PostgreSQL-SLOW",
		"cloud_product_region": "gz",
		"cls_region":           "ap-guangzhou",
		"tags": map[string]interface{}{
			"Environment": "staging",
		},
	})
	d.SetId("postgres-0an6hpv3#PostgreSQL#PostgreSQL-SLOW#gz")

	// Force only `tags` to be detected as changed so immutableArgs (cls_region)
	// and other update branches are not triggered.
	patches.ApplyMethodFunc(d, "HasChange", func(key string) bool {
		return key == "tags"
	})

	err := res.Update(d, meta)
	assert.NoError(t, err)

	// Verify the modify request carries the updated Tags.
	assert.NotNil(t, capturedRequest)
	assert.NotNil(t, capturedRequest.Tags)
	assert.Equal(t, 1, len(capturedRequest.Tags))
	assert.Equal(t, "Environment", *capturedRequest.Tags[0].Key)
	assert.Equal(t, "staging", *capturedRequest.Tags[0].Value)
	// The id must be preserved (no recreation).
	assert.Equal(t, "postgres-0an6hpv3#PostgreSQL#PostgreSQL-SLOW#gz", d.Id())
}

// TestClsCloudProductLogTaskV2_Update_RemoveTags verifies that when all tags
// are removed, ModifyCloudProductLogCollection is called with an empty (non-nil)
// Tags slice so the API clears them.
func TestClsCloudProductLogTaskV2_Update_RemoveTags(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	clsClient := &clsv20201016.Client{}
	patches.ApplyMethodReturn(newMockMetaForClsCloudProductLogTaskV2().client, "UseClsV20201016Client", clsClient)
	patches.ApplyMethodReturn(newMockMetaForClsCloudProductLogTaskV2().client, "UseClsClient", clsClient)

	var capturedRequest *clsv20201016.ModifyCloudProductLogCollectionRequest
	patches.ApplyMethodFunc(clsClient, "ModifyCloudProductLogCollectionWithContext", func(_ context.Context, request *clsv20201016.ModifyCloudProductLogCollectionRequest) (*clsv20201016.ModifyCloudProductLogCollectionResponse, error) {
		capturedRequest = request
		resp := clsv20201016.NewModifyCloudProductLogCollectionResponse()
		resp.Response = &clsv20201016.ModifyCloudProductLogCollectionResponseParams{
			RequestId: ptrStringCPLT2("fake-request-id-update"),
		}
		return resp, nil
	})

	task := &clsv20201016.CloudProductLogTaskInfo{
		InstanceId: ptrStringCPLT2("postgres-0an6hpv3"),
		ClsRegion:  ptrStringCPLT2("ap-guangzhou"),
		LogType:    ptrStringCPLT2("PostgreSQL-SLOW"),
		Status:     ptrInt64CPLT2(1),
	}
	patchDescribeCloudProductLogTasks(patches, clsClient, task)

	meta := newMockMetaForClsCloudProductLogTaskV2()
	res := cls.ResourceTencentCloudClsCloudProductLogTaskV2()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"instance_id":          "postgres-0an6hpv3",
		"assumer_name":         "PostgreSQL",
		"log_type":             "PostgreSQL-SLOW",
		"cloud_product_region": "gz",
		"cls_region":           "ap-guangzhou",
	})
	d.SetId("postgres-0an6hpv3#PostgreSQL#PostgreSQL-SLOW#gz")

	// Force only `tags` to be detected as changed so immutableArgs (cls_region)
	// and other update branches are not triggered. Even though the user removed
	// tags, we simulate the change detection so the modify path runs.
	patches.ApplyMethodFunc(d, "HasChange", func(key string) bool {
		return key == "tags"
	})

	err := res.Update(d, meta)
	assert.NoError(t, err)

	// Verify the modify request carries an empty (non-nil) Tags slice.
	assert.NotNil(t, capturedRequest)
	assert.NotNil(t, capturedRequest.Tags)
	assert.Equal(t, 0, len(capturedRequest.Tags))
}

// TestClsCloudProductLogTaskV2_Read_Tags verifies that tags from the API
// response (TopicTags) are flattened into the Terraform state.
func TestClsCloudProductLogTaskV2_Read_Tags(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	clsClient := &clsv20201016.Client{}
	patches.ApplyMethodReturn(newMockMetaForClsCloudProductLogTaskV2().client, "UseClsV20201016Client", clsClient)
	patches.ApplyMethodReturn(newMockMetaForClsCloudProductLogTaskV2().client, "UseClsClient", clsClient)

	task := &clsv20201016.CloudProductLogTaskInfo{
		InstanceId: ptrStringCPLT2("postgres-0an6hpv3"),
		ClsRegion:  ptrStringCPLT2("ap-guangzhou"),
		LogType:    ptrStringCPLT2("PostgreSQL-SLOW"),
		Status:     ptrInt64CPLT2(1),
		TopicTags: []*clsv20201016.Tag{
			{Key: ptrStringCPLT2("Environment"), Value: ptrStringCPLT2("production")},
			{Key: ptrStringCPLT2("Team"), Value: ptrStringCPLT2("backend")},
		},
	}
	patchDescribeCloudProductLogTasks(patches, clsClient, task)

	meta := newMockMetaForClsCloudProductLogTaskV2()
	res := cls.ResourceTencentCloudClsCloudProductLogTaskV2()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"instance_id":          "postgres-0an6hpv3",
		"assumer_name":         "PostgreSQL",
		"log_type":             "PostgreSQL-SLOW",
		"cloud_product_region": "gz",
		"cls_region":           "ap-guangzhou",
	})
	d.SetId("postgres-0an6hpv3#PostgreSQL#PostgreSQL-SLOW#gz")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "postgres-0an6hpv3#PostgreSQL#PostgreSQL-SLOW#gz", d.Id())

	tags := d.Get("tags").(map[string]interface{})
	assert.Equal(t, "production", tags["Environment"])
	assert.Equal(t, "backend", tags["Team"])
}

// TestClsCloudProductLogTaskV2_Read_Tags_FromLogsetTags verifies that when
// TopicTags is empty, tags are read from LogsetTags as a fallback.
func TestClsCloudProductLogTaskV2_Read_Tags_FromLogsetTags(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	clsClient := &clsv20201016.Client{}
	patches.ApplyMethodReturn(newMockMetaForClsCloudProductLogTaskV2().client, "UseClsV20201016Client", clsClient)
	patches.ApplyMethodReturn(newMockMetaForClsCloudProductLogTaskV2().client, "UseClsClient", clsClient)

	task := &clsv20201016.CloudProductLogTaskInfo{
		InstanceId: ptrStringCPLT2("postgres-0an6hpv3"),
		ClsRegion:  ptrStringCPLT2("ap-guangzhou"),
		LogType:    ptrStringCPLT2("PostgreSQL-SLOW"),
		Status:     ptrInt64CPLT2(1),
		LogsetTags: []*clsv20201016.Tag{
			{Key: ptrStringCPLT2("Owner"), Value: ptrStringCPLT2("ops")},
		},
	}
	patchDescribeCloudProductLogTasks(patches, clsClient, task)

	meta := newMockMetaForClsCloudProductLogTaskV2()
	res := cls.ResourceTencentCloudClsCloudProductLogTaskV2()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"instance_id":          "postgres-0an6hpv3",
		"assumer_name":         "PostgreSQL",
		"log_type":             "PostgreSQL-SLOW",
		"cloud_product_region": "gz",
		"cls_region":           "ap-guangzhou",
	})
	d.SetId("postgres-0an6hpv3#PostgreSQL#PostgreSQL-SLOW#gz")

	err := res.Read(d, meta)
	assert.NoError(t, err)

	tags := d.Get("tags").(map[string]interface{})
	assert.Equal(t, "ops", tags["Owner"])
}
