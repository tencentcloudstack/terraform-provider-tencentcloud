package cls_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/stretchr/testify/assert"
	cls "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cls/v20201016"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	localcls "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/cls"
)

// mockMetaForDataTransform is a mock ProviderMeta for data transform cross-account tests
type mockMetaForDataTransform struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMetaForDataTransform) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockMetaForDataTransform{}

func newMockMetaForDataTransform() *mockMetaForDataTransform {
	return &mockMetaForDataTransform{client: &connectivity.TencentCloudClient{}}
}

func ptrStrDataTransform(s string) *string {
	return &s
}

func ptrBoolDataTransform(b bool) *bool {
	return &b
}

func ptrInt64DataTransform(v int64) *int64 {
	return &v
}

func ptrUint64DataTransform(v uint64) *uint64 {
	return &v
}

// go test -i; go test -test.run TestAccTencentCloudClsDataTransformResource_basic -v
func TestAccTencentCloudClsDataTransformResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		CheckDestroy: testAccCheckClsDataTransformDestroy,
		Providers:    tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccClsDataTransform,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClsDataTransformExists("tencentcloud_cls_data_transform.data_transform"),
					resource.TestCheckResourceAttrSet("tencentcloud_cls_data_transform.data_transform", "func_type"),
					resource.TestCheckResourceAttrSet("tencentcloud_cls_data_transform.data_transform", "src_topic_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_cls_data_transform.data_transform", "name"),
					resource.TestCheckResourceAttrSet("tencentcloud_cls_data_transform.data_transform", "etl_content"),
					resource.TestCheckResourceAttrSet("tencentcloud_cls_data_transform.data_transform", "task_type"),
					resource.TestCheckResourceAttrSet("tencentcloud_cls_data_transform.data_transform", "enable_flag"),
					resource.TestCheckResourceAttrSet("tencentcloud_cls_data_transform.data_transform", "dst_resources.#"),
					resource.TestCheckResourceAttrSet("tencentcloud_cls_data_transform.data_transform", "dst_resources.0.topic_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_cls_data_transform.data_transform", "dst_resources.0.alias")),
			},
		},
	})
}

func testAccCheckClsDataTransformDestroy(s *terraform.State) error {
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	clsService := localcls.NewClsService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_cls_data_transform" {
			continue
		}
		time.Sleep(5 * time.Second)
		instance, err := clsService.DescribeClsDataTransformById(ctx, rs.Primary.ID)
		if err != nil {
			continue
		}
		if instance != nil && instance.TaskId != nil && *instance.TaskId == rs.Primary.ID {
			return fmt.Errorf("[CHECK][CLS dataTransform][Destroy] check: CLS dataTransform still exists: %s", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckClsDataTransformExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := tccommon.GetLogId(tccommon.ContextNil)
		ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("[CHECK][CLS dataTransform][Exists] check: CLS dataTransform %s is not found", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("[CHECK][CLS dataTransform][Create] check: CLS dataTransform id is not set")
		}
		clsService := localcls.NewClsService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
		taskRes, err := clsService.DescribeClsDataTransformById(ctx, rs.Primary.ID)
		if err != nil {
			return err
		}
		if taskRes == nil {
			return fmt.Errorf("[CHECK][CLS redirection][Exists] id %s is not exist", rs.Primary.ID)
		}
		return nil
	}
}

const testAccClsDataTransform = `
resource "tencentcloud_cls_logset" "logset_src" {
  logset_name = "tf-example-src"
  tags = {
    "createdBy" = "terraform"
  }
}
resource "tencentcloud_cls_topic" "topic_src" {
  topic_name           = "tf-example_src"
  logset_id            = tencentcloud_cls_logset.logset_src.id
  auto_split           = false
  max_split_partitions = 20
  partition_count      = 1
  period               = 10
  storage_type         = "hot"
  tags                 = {
    "test" = "test",
  }
}
resource "tencentcloud_cls_logset" "logset_dst" {
  logset_name = "tf-example-dst"
  tags = {
    "createdBy" = "terraform"
  }
}
resource "tencentcloud_cls_topic" "topic_dst" {
  topic_name           = "tf-example-dst"
  logset_id            = tencentcloud_cls_logset.logset_dst.id
  auto_split           = false
  max_split_partitions = 20
  partition_count      = 1
  period               = 10
  storage_type         = "hot"
  tags                 = {
    "test" = "test",
  }
}
resource "tencentcloud_cls_data_transform" "data_transform" {
  func_type = 1
  src_topic_id = tencentcloud_cls_topic.topic_src.id
  name = "iac-test-src"
  etl_content = "ext_sep(\"content\", \"f1, f2, f3\", sep=\",\", quote=\"\", restrict=False, mode=\"overwrite\")fields_drop(\"content\")"
  task_type = 3
  enable_flag = 1
  dst_resources {
    topic_id = tencentcloud_cls_topic.topic_dst.id
    alias = "iac-test-dst"

  }
}

`

// go test ./tencentcloud/services/cls/ -run "TestDataTransformCrossAccount" -v -count=1 -gcflags="all=-l"

// TestDataTransformCrossAccount_Create tests that Create correctly populates cross-account fields in the request
func TestDataTransformCrossAccount_Create(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	clsClient := &cls.Client{}
	patches.ApplyMethodReturn(newMockMetaForDataTransform().client, "UseClsClient", clsClient)

	var capturedRequest *cls.CreateDataTransformRequest
	patches.ApplyMethodFunc(clsClient, "CreateDataTransform", func(request *cls.CreateDataTransformRequest) (*cls.CreateDataTransformResponse, error) {
		capturedRequest = request
		resp := cls.NewCreateDataTransformResponse()
		resp.Response = &cls.CreateDataTransformResponseParams{
			TaskId:    ptrStrDataTransform("task-test-cross-account"),
			RequestId: ptrStrDataTransform("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(clsClient, "DescribeDataTransformInfo", func(request *cls.DescribeDataTransformInfoRequest) (*cls.DescribeDataTransformInfoResponse, error) {
		resp := cls.NewDescribeDataTransformInfoResponse()
		resp.Response = &cls.DescribeDataTransformInfoResponseParams{
			DataTransformTaskInfos: []*cls.DataTransformTaskInfo{
				{
					Name:       ptrStrDataTransform("tf-example-cross-account"),
					TaskId:     ptrStrDataTransform("task-test-cross-account"),
					SrcTopicId: ptrStrDataTransform("topic-src-123"),
					EtlContent: ptrStrDataTransform("ext_sep(\"content\", \"f1, f2, f3\", sep=\",\", quote=\"\", restrict=False, mode=\"overwrite\")fields_drop(\"content\")"),
					EnableFlag: ptrInt64DataTransform(1),
					DstResources: []*cls.DataTransformResouceInfo{
						{
							TopicId:        ptrStrDataTransform("topic-id-in-target-account"),
							Alias:          ptrStrDataTransform("cross-account-dst"),
							IsCrossAccount: ptrBoolDataTransform(true),
							RoleARN:        ptrStrDataTransform("qcs::cam::uin/123456789:roleName/cls-cross-account-role"),
							ExternalId:     ptrStrDataTransform("external-id-value"),
							TopicName:      ptrStrDataTransform("target-topic-name"),
							LogsetName:     ptrStrDataTransform("target-logset-name"),
						},
					},
				},
			},
			TotalCount: ptrUint64DataTransform(1),
			RequestId:  ptrStrDataTransform("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaForDataTransform()
	res := localcls.ResourceTencentCloudClsDataTransform()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"func_type":    1,
		"src_topic_id": "topic-src-123",
		"name":         "tf-example-cross-account",
		"etl_content":  "ext_sep(\"content\", \"f1, f2, f3\", sep=\",\", quote=\"\", restrict=False, mode=\"overwrite\")fields_drop(\"content\")",
		"task_type":    3,
		"enable_flag":  1,
		"dst_resources": []interface{}{
			map[string]interface{}{
				"topic_id":         "topic-id-in-target-account",
				"alias":            "cross-account-dst",
				"is_cross_account": true,
				"role_arn":         "qcs::cam::uin/123456789:roleName/cls-cross-account-role",
				"external_id":      "external-id-value",
				"topic_name":       "target-topic-name",
				"logset_name":      "target-logset-name",
			},
		},
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "task-test-cross-account", d.Id())
	assert.NotNil(t, capturedRequest)
	assert.Len(t, capturedRequest.DstResources, 1)

	dst := capturedRequest.DstResources[0]
	assert.NotNil(t, dst.IsCrossAccount)
	assert.True(t, *dst.IsCrossAccount)
	assert.NotNil(t, dst.RoleARN)
	assert.Equal(t, "qcs::cam::uin/123456789:roleName/cls-cross-account-role", *dst.RoleARN)
	assert.NotNil(t, dst.ExternalId)
	assert.Equal(t, "external-id-value", *dst.ExternalId)
	assert.NotNil(t, dst.TopicName)
	assert.Equal(t, "target-topic-name", *dst.TopicName)
	assert.NotNil(t, dst.LogsetName)
	assert.Equal(t, "target-logset-name", *dst.LogsetName)
}

// TestDataTransformCrossAccount_Create_WithoutCrossAccount tests that Create does not set cross-account fields when not specified
func TestDataTransformCrossAccount_Create_WithoutCrossAccount(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	clsClient := &cls.Client{}
	patches.ApplyMethodReturn(newMockMetaForDataTransform().client, "UseClsClient", clsClient)

	var capturedRequest *cls.CreateDataTransformRequest
	patches.ApplyMethodFunc(clsClient, "CreateDataTransform", func(request *cls.CreateDataTransformRequest) (*cls.CreateDataTransformResponse, error) {
		capturedRequest = request
		resp := cls.NewCreateDataTransformResponse()
		resp.Response = &cls.CreateDataTransformResponseParams{
			TaskId:    ptrStrDataTransform("task-test-same-account"),
			RequestId: ptrStrDataTransform("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(clsClient, "DescribeDataTransformInfo", func(request *cls.DescribeDataTransformInfoRequest) (*cls.DescribeDataTransformInfoResponse, error) {
		resp := cls.NewDescribeDataTransformInfoResponse()
		resp.Response = &cls.DescribeDataTransformInfoResponseParams{
			DataTransformTaskInfos: []*cls.DataTransformTaskInfo{
				{
					Name:       ptrStrDataTransform("tf-example-same-account"),
					TaskId:     ptrStrDataTransform("task-test-same-account"),
					SrcTopicId: ptrStrDataTransform("topic-src-123"),
					EtlContent: ptrStrDataTransform("ext_sep(\"content\", \"f1, f2, f3\", sep=\",\", quote=\"\", restrict=False, mode=\"overwrite\")fields_drop(\"content\")"),
					EnableFlag: ptrInt64DataTransform(1),
					DstResources: []*cls.DataTransformResouceInfo{
						{
							TopicId: ptrStrDataTransform("topic-dst-123"),
							Alias:   ptrStrDataTransform("same-account-dst"),
						},
					},
				},
			},
			TotalCount: ptrUint64DataTransform(1),
			RequestId:  ptrStrDataTransform("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaForDataTransform()
	res := localcls.ResourceTencentCloudClsDataTransform()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"func_type":    1,
		"src_topic_id": "topic-src-123",
		"name":         "tf-example-same-account",
		"etl_content":  "ext_sep(\"content\", \"f1, f2, f3\", sep=\",\", quote=\"\", restrict=False, mode=\"overwrite\")fields_drop(\"content\")",
		"task_type":    3,
		"enable_flag":  1,
		"dst_resources": []interface{}{
			map[string]interface{}{
				"topic_id": "topic-dst-123",
				"alias":    "same-account-dst",
			},
		},
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "task-test-same-account", d.Id())
	assert.NotNil(t, capturedRequest)
	assert.Len(t, capturedRequest.DstResources, 1)

	dst := capturedRequest.DstResources[0]
	assert.Nil(t, dst.IsCrossAccount)
	assert.Nil(t, dst.RoleARN)
	assert.Nil(t, dst.ExternalId)
	assert.Nil(t, dst.TopicName)
	assert.Nil(t, dst.LogsetName)
}

// TestDataTransformCrossAccount_Read_WithCrossAccount tests that Read populates cross-account fields from API response
func TestDataTransformCrossAccount_Read_WithCrossAccount(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	clsClient := &cls.Client{}
	patches.ApplyMethodReturn(newMockMetaForDataTransform().client, "UseClsClient", clsClient)

	patches.ApplyMethodFunc(clsClient, "DescribeDataTransformInfo", func(request *cls.DescribeDataTransformInfoRequest) (*cls.DescribeDataTransformInfoResponse, error) {
		resp := cls.NewDescribeDataTransformInfoResponse()
		resp.Response = &cls.DescribeDataTransformInfoResponseParams{
			DataTransformTaskInfos: []*cls.DataTransformTaskInfo{
				{
					Name:       ptrStrDataTransform("tf-example-cross-account"),
					TaskId:     ptrStrDataTransform("task-test-cross-account"),
					SrcTopicId: ptrStrDataTransform("topic-src-123"),
					EtlContent: ptrStrDataTransform("ext_sep(\"content\", \"f1, f2, f3\", sep=\",\", quote=\"\", restrict=False, mode=\"overwrite\")fields_drop(\"content\")"),
					EnableFlag: ptrInt64DataTransform(1),
					DstResources: []*cls.DataTransformResouceInfo{
						{
							TopicId:        ptrStrDataTransform("topic-id-in-target-account"),
							Alias:          ptrStrDataTransform("cross-account-dst"),
							IsCrossAccount: ptrBoolDataTransform(true),
							RoleARN:        ptrStrDataTransform("qcs::cam::uin/123456789:roleName/cls-cross-account-role"),
							ExternalId:     ptrStrDataTransform("external-id-value"),
							TopicName:      ptrStrDataTransform("target-topic-name"),
							LogsetName:     ptrStrDataTransform("target-logset-name"),
						},
					},
				},
			},
			TotalCount: ptrUint64DataTransform(1),
			RequestId:  ptrStrDataTransform("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaForDataTransform()
	res := localcls.ResourceTencentCloudClsDataTransform()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"func_type":    1,
		"src_topic_id": "topic-src-123",
		"name":         "tf-example-cross-account",
		"etl_content":  "ext_sep(\"content\", \"f1, f2, f3\", sep=\",\", quote=\"\", restrict=False, mode=\"overwrite\")fields_drop(\"content\")",
		"task_type":    3,
		"enable_flag":  1,
	})
	d.SetId("task-test-cross-account")

	err := res.Read(d, meta)
	assert.NoError(t, err)

	dstResources := d.Get("dst_resources").([]interface{})
	assert.Len(t, dstResources, 1)

	dstMap := dstResources[0].(map[string]interface{})
	assert.Equal(t, "topic-id-in-target-account", dstMap["topic_id"])
	assert.Equal(t, "cross-account-dst", dstMap["alias"])
	assert.Equal(t, true, dstMap["is_cross_account"])
	assert.Equal(t, "qcs::cam::uin/123456789:roleName/cls-cross-account-role", dstMap["role_arn"])
	assert.Equal(t, "external-id-value", dstMap["external_id"])
	assert.Equal(t, "target-topic-name", dstMap["topic_name"])
	assert.Equal(t, "target-logset-name", dstMap["logset_name"])
}

// TestDataTransformCrossAccount_Read_NilCrossAccount tests that Read handles nil cross-account fields
func TestDataTransformCrossAccount_Read_NilCrossAccount(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	clsClient := &cls.Client{}
	patches.ApplyMethodReturn(newMockMetaForDataTransform().client, "UseClsClient", clsClient)

	patches.ApplyMethodFunc(clsClient, "DescribeDataTransformInfo", func(request *cls.DescribeDataTransformInfoRequest) (*cls.DescribeDataTransformInfoResponse, error) {
		resp := cls.NewDescribeDataTransformInfoResponse()
		resp.Response = &cls.DescribeDataTransformInfoResponseParams{
			DataTransformTaskInfos: []*cls.DataTransformTaskInfo{
				{
					Name:       ptrStrDataTransform("tf-example-same-account"),
					TaskId:     ptrStrDataTransform("task-test-same-account"),
					SrcTopicId: ptrStrDataTransform("topic-src-123"),
					EtlContent: ptrStrDataTransform("ext_sep(\"content\", \"f1, f2, f3\", sep=\",\", quote=\"\", restrict=False, mode=\"overwrite\")fields_drop(\"content\")"),
					EnableFlag: ptrInt64DataTransform(1),
					DstResources: []*cls.DataTransformResouceInfo{
						{
							TopicId: ptrStrDataTransform("topic-dst-123"),
							Alias:   ptrStrDataTransform("same-account-dst"),
						},
					},
				},
			},
			TotalCount: ptrUint64DataTransform(1),
			RequestId:  ptrStrDataTransform("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaForDataTransform()
	res := localcls.ResourceTencentCloudClsDataTransform()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"func_type":    1,
		"src_topic_id": "topic-src-123",
		"name":         "tf-example-same-account",
		"etl_content":  "ext_sep(\"content\", \"f1, f2, f3\", sep=\",\", quote=\"\", restrict=False, mode=\"overwrite\")fields_drop(\"content\")",
		"task_type":    3,
		"enable_flag":  1,
	})
	d.SetId("task-test-same-account")

	err := res.Read(d, meta)
	assert.NoError(t, err)

	dstResources := d.Get("dst_resources").([]interface{})
	assert.Len(t, dstResources, 1)

	dstMap := dstResources[0].(map[string]interface{})
	assert.Equal(t, "topic-dst-123", dstMap["topic_id"])
	assert.Equal(t, "same-account-dst", dstMap["alias"])
	_, exists := dstMap["is_cross_account"]
	assert.False(t, exists)
	_, exists = dstMap["role_arn"]
	assert.False(t, exists)
}

// TestDataTransformCrossAccount_Update tests that Update correctly populates cross-account fields in the modify request
func TestDataTransformCrossAccount_Update(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	clsClient := &cls.Client{}
	patches.ApplyMethodReturn(newMockMetaForDataTransform().client, "UseClsClient", clsClient)

	var capturedModifyRequest *cls.ModifyDataTransformRequest
	patches.ApplyMethodFunc(clsClient, "ModifyDataTransform", func(request *cls.ModifyDataTransformRequest) (*cls.ModifyDataTransformResponse, error) {
		capturedModifyRequest = request
		resp := cls.NewModifyDataTransformResponse()
		resp.Response = &cls.ModifyDataTransformResponseParams{
			RequestId: ptrStrDataTransform("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(clsClient, "DescribeDataTransformInfo", func(request *cls.DescribeDataTransformInfoRequest) (*cls.DescribeDataTransformInfoResponse, error) {
		resp := cls.NewDescribeDataTransformInfoResponse()
		resp.Response = &cls.DescribeDataTransformInfoResponseParams{
			DataTransformTaskInfos: []*cls.DataTransformTaskInfo{
				{
					Name:       ptrStrDataTransform("tf-example-cross-account-updated"),
					TaskId:     ptrStrDataTransform("task-test-cross-account"),
					SrcTopicId: ptrStrDataTransform("topic-src-123"),
					EtlContent: ptrStrDataTransform("ext_sep(\"content\", \"f1, f2, f3\", sep=\",\", quote=\"\", restrict=False, mode=\"overwrite\")fields_drop(\"content\")"),
					EnableFlag: ptrInt64DataTransform(1),
					DstResources: []*cls.DataTransformResouceInfo{
						{
							TopicId:        ptrStrDataTransform("topic-id-in-target-account"),
							Alias:          ptrStrDataTransform("cross-account-dst-updated"),
							IsCrossAccount: ptrBoolDataTransform(true),
							RoleARN:        ptrStrDataTransform("qcs::cam::uin/123456789:roleName/cls-cross-account-role"),
							ExternalId:     ptrStrDataTransform("external-id-value"),
							TopicName:      ptrStrDataTransform("target-topic-name"),
							LogsetName:     ptrStrDataTransform("target-logset-name"),
						},
					},
				},
			},
			TotalCount: ptrUint64DataTransform(1),
			RequestId:  ptrStrDataTransform("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaForDataTransform()
	res := localcls.ResourceTencentCloudClsDataTransform()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"func_type":    1,
		"src_topic_id": "topic-src-123",
		"name":         "tf-example-cross-account-updated",
		"etl_content":  "ext_sep(\"content\", \"f1, f2, f3\", sep=\",\", quote=\"\", restrict=False, mode=\"overwrite\")fields_drop(\"content\")",
		"task_type":    3,
		"enable_flag":  1,
		"dst_resources": []interface{}{
			map[string]interface{}{
				"topic_id":         "topic-id-in-target-account",
				"alias":            "cross-account-dst-updated",
				"is_cross_account": true,
				"role_arn":         "qcs::cam::uin/123456789:roleName/cls-cross-account-role",
				"external_id":      "external-id-value",
				"topic_name":       "target-topic-name",
				"logset_name":      "target-logset-name",
			},
		},
	})
	d.SetId("task-test-cross-account")

	err := res.Update(d, meta)
	assert.NoError(t, err)
	assert.NotNil(t, capturedModifyRequest)
	assert.Len(t, capturedModifyRequest.DstResources, 1)

	dst := capturedModifyRequest.DstResources[0]
	assert.NotNil(t, dst.IsCrossAccount)
	assert.True(t, *dst.IsCrossAccount)
	assert.NotNil(t, dst.RoleARN)
	assert.Equal(t, "qcs::cam::uin/123456789:roleName/cls-cross-account-role", *dst.RoleARN)
	assert.NotNil(t, dst.ExternalId)
	assert.Equal(t, "external-id-value", *dst.ExternalId)
	assert.NotNil(t, dst.TopicName)
	assert.Equal(t, "target-topic-name", *dst.TopicName)
	assert.NotNil(t, dst.LogsetName)
	assert.Equal(t, "target-logset-name", *dst.LogsetName)
}

// TestDataTransformCrossAccount_Schema tests the cross-account schema definitions
func TestDataTransformCrossAccount_Schema(t *testing.T) {
	res := localcls.ResourceTencentCloudClsDataTransform()

	assert.NotNil(t, res)
	assert.Contains(t, res.Schema, "dst_resources")

	dstResourcesSchema := res.Schema["dst_resources"]
	assert.NotNil(t, dstResourcesSchema.Elem)

	elemResource, ok := dstResourcesSchema.Elem.(*schema.Resource)
	assert.True(t, ok)
	assert.NotNil(t, elemResource.Schema)

	isCrossAccountSchema, ok := elemResource.Schema["is_cross_account"]
	assert.True(t, ok)
	assert.Equal(t, schema.TypeBool, isCrossAccountSchema.Type)
	assert.True(t, isCrossAccountSchema.Optional)
	assert.False(t, isCrossAccountSchema.Required)

	roleArnSchema, ok := elemResource.Schema["role_arn"]
	assert.True(t, ok)
	assert.Equal(t, schema.TypeString, roleArnSchema.Type)
	assert.True(t, roleArnSchema.Optional)
	assert.False(t, roleArnSchema.Required)

	externalIdSchema, ok := elemResource.Schema["external_id"]
	assert.True(t, ok)
	assert.Equal(t, schema.TypeString, externalIdSchema.Type)
	assert.True(t, externalIdSchema.Optional)
	assert.False(t, externalIdSchema.Required)

	topicNameSchema, ok := elemResource.Schema["topic_name"]
	assert.True(t, ok)
	assert.Equal(t, schema.TypeString, topicNameSchema.Type)
	assert.True(t, topicNameSchema.Optional)
	assert.False(t, topicNameSchema.Required)

	logsetNameSchema, ok := elemResource.Schema["logset_name"]
	assert.True(t, ok)
	assert.Equal(t, schema.TypeString, logsetNameSchema.Type)
	assert.True(t, logsetNameSchema.Optional)
	assert.False(t, logsetNameSchema.Required)
}
