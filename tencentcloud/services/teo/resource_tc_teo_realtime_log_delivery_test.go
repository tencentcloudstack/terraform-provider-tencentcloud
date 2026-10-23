package teo_test

import (
	"context"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	teov20220901 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/teo"
)

func TestAccTencentCloudTeoRealtimeLogDeliveryResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTeoRealtimeLogDelivery,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_teo_realtime_log_delivery.teo_realtime_log_delivery", "id"),
					resource.TestCheckResourceAttr("tencentcloud_teo_realtime_log_delivery.teo_realtime_log_delivery", "area", "overseas"),
					resource.TestCheckResourceAttr("tencentcloud_teo_realtime_log_delivery.teo_realtime_log_delivery", "delivery_status", "enabled"),
					resource.TestCheckResourceAttr("tencentcloud_teo_realtime_log_delivery.teo_realtime_log_delivery", "entity_list.#", "1"),
					resource.TestCheckResourceAttrSet("tencentcloud_teo_realtime_log_delivery.teo_realtime_log_delivery", "fields.#"),
					resource.TestCheckResourceAttr("tencentcloud_teo_realtime_log_delivery.teo_realtime_log_delivery", "log_type", "application"),
					resource.TestCheckResourceAttr("tencentcloud_teo_realtime_log_delivery.teo_realtime_log_delivery", "sample", "0"),
					resource.TestCheckResourceAttr("tencentcloud_teo_realtime_log_delivery.teo_realtime_log_delivery", "task_name", "test"),
					resource.TestCheckResourceAttr("tencentcloud_teo_realtime_log_delivery.teo_realtime_log_delivery", "task_type", "s3"),
					resource.TestCheckResourceAttr("tencentcloud_teo_realtime_log_delivery.teo_realtime_log_delivery", "log_format.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_teo_realtime_log_delivery.teo_realtime_log_delivery", "s3.#", "1"),
					resource.TestCheckResourceAttrSet("tencentcloud_teo_realtime_log_delivery.teo_realtime_log_delivery", "s3.0.access_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_teo_realtime_log_delivery.teo_realtime_log_delivery", "s3.0.access_key"),
					resource.TestCheckResourceAttrSet("tencentcloud_teo_realtime_log_delivery.teo_realtime_log_delivery", "s3.0.bucket"),
					resource.TestCheckResourceAttrSet("tencentcloud_teo_realtime_log_delivery.teo_realtime_log_delivery", "s3.0.compress_type"),
					resource.TestCheckResourceAttrSet("tencentcloud_teo_realtime_log_delivery.teo_realtime_log_delivery", "s3.0.endpoint"),
					resource.TestCheckResourceAttrSet("tencentcloud_teo_realtime_log_delivery.teo_realtime_log_delivery", "s3.0.region"),
				),
			},
			{
				ResourceName:      "tencentcloud_teo_realtime_log_delivery.teo_realtime_log_delivery",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccTeoRealtimeLogDeliveryUp,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_teo_realtime_log_delivery.teo_realtime_log_delivery", "id"),
					resource.TestCheckResourceAttr("tencentcloud_teo_realtime_log_delivery.teo_realtime_log_delivery", "area", "overseas"),
					resource.TestCheckResourceAttr("tencentcloud_teo_realtime_log_delivery.teo_realtime_log_delivery", "delivery_status", "disabled"),
					resource.TestCheckResourceAttr("tencentcloud_teo_realtime_log_delivery.teo_realtime_log_delivery", "entity_list.#", "1"),
					resource.TestCheckResourceAttrSet("tencentcloud_teo_realtime_log_delivery.teo_realtime_log_delivery", "fields.#"),
					resource.TestCheckResourceAttr("tencentcloud_teo_realtime_log_delivery.teo_realtime_log_delivery", "log_type", "application"),
					resource.TestCheckResourceAttr("tencentcloud_teo_realtime_log_delivery.teo_realtime_log_delivery", "sample", "0"),
					resource.TestCheckResourceAttr("tencentcloud_teo_realtime_log_delivery.teo_realtime_log_delivery", "task_name", "test"),
					resource.TestCheckResourceAttr("tencentcloud_teo_realtime_log_delivery.teo_realtime_log_delivery", "task_type", "s3"),
					resource.TestCheckResourceAttr("tencentcloud_teo_realtime_log_delivery.teo_realtime_log_delivery", "log_format.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_teo_realtime_log_delivery.teo_realtime_log_delivery", "s3.#", "1"),
					resource.TestCheckResourceAttrSet("tencentcloud_teo_realtime_log_delivery.teo_realtime_log_delivery", "s3.0.access_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_teo_realtime_log_delivery.teo_realtime_log_delivery", "s3.0.access_key"),
					resource.TestCheckResourceAttrSet("tencentcloud_teo_realtime_log_delivery.teo_realtime_log_delivery", "s3.0.bucket"),
					resource.TestCheckResourceAttrSet("tencentcloud_teo_realtime_log_delivery.teo_realtime_log_delivery", "s3.0.compress_type"),
					resource.TestCheckResourceAttrSet("tencentcloud_teo_realtime_log_delivery.teo_realtime_log_delivery", "s3.0.endpoint"),
					resource.TestCheckResourceAttrSet("tencentcloud_teo_realtime_log_delivery.teo_realtime_log_delivery", "s3.0.region"),
				),
			},
		},
	})
}

const testAccTeoRealtimeLogDelivery = `

resource "tencentcloud_teo_realtime_log_delivery" "teo_realtime_log_delivery" {
    area            = "overseas"
    delivery_status = "enabled"
    entity_list     = [
        "sid-2yvhjw98uaco",
    ]
    fields          = [
        "ServiceID",
        "ConnectTimeStamp",
        "DisconnetTimeStamp",
        "DisconnetReason",
        "ClientRealIP",
        "ClientRegion",
        "EdgeIP",
        "ForwardProtocol",
        "ForwardPort",
        "SentBytes",
        "ReceivedBytes",
        "LogTimeStamp",
    ]
    log_type        = "application"
    sample          = 0
    task_name       = "test"
    task_type       = "s3"
    zone_id         = "zone-2qtuhspy7cr6"

    log_format {
        field_delimiter  = ","
        format_type      = "json"
        record_delimiter = "\n"
        record_prefix    = "{"
        record_suffix    = "}"
    }

    s3 {
        access_id     = "xxxxxxxxxx"
        access_key    = "xxxxxxxxxx"
        bucket        = "test-1253833068"
        compress_type = "gzip"
        endpoint      = "https://test-1253833068.cos.ap-nanjing.myqcloud.com"
        region        = "ap-nanjing"
    }
}
`

const testAccTeoRealtimeLogDeliveryUp = `

resource "tencentcloud_teo_realtime_log_delivery" "teo_realtime_log_delivery" {
    area            = "overseas"
    delivery_status = "disabled"
    entity_list     = [
        "sid-2yvhjw98uaco",
    ]
    fields          = [
        "ServiceID",
        "ConnectTimeStamp",
        "DisconnetTimeStamp",
        "DisconnetReason",
        "ClientRealIP",
        "ClientRegion",
        "EdgeIP",
        "ForwardProtocol",
        "ForwardPort",
        "SentBytes",
        "ReceivedBytes",
        "LogTimeStamp",
    ]
    log_type        = "application"
    sample          = 0
    task_name       = "test"
    task_type       = "s3"
    zone_id         = "zone-2qtuhspy7cr6"

    log_format {
        field_delimiter  = ","
        format_type      = "json"
        record_delimiter = "\n"
        record_prefix    = "{"
        record_suffix    = "}"
    }

    s3 {
        access_id     = "xxxxxxxxxx"
        access_key    = "xxxxxxxxxx"
        bucket        = "test-1253833068"
        compress_type = "gzip"
        endpoint      = "https://test-1253833068.cos.ap-nanjing.myqcloud.com"
        region        = "ap-nanjing"
    }
}
`

// mockMetaForRealtimeLogDelivery implements tccommon.ProviderMeta
type mockMetaForRealtimeLogDelivery struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMetaForRealtimeLogDelivery) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockMetaForRealtimeLogDelivery{}

func newMockMetaForRealtimeLogDelivery() *mockMetaForRealtimeLogDelivery {
	return &mockMetaForRealtimeLogDelivery{client: &connectivity.TencentCloudClient{}}
}

func ptrStrRTLD(s string) *string { return &s }
func ptrBoolRTLD(b bool) *bool    { return &b }

// go test ./tencentcloud/services/teo/ -run "TestRealtimeLogDelivery" -v -count=1 -gcflags="all=-l"

// TestRealtimeLogDelivery_Create_WithNewFields tests Create with custom_expression_fields, aws_s3, and log_format.record_template
func TestRealtimeLogDelivery_Create_WithNewFields(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaForRealtimeLogDelivery().client, "UseTeoClient", teoClient)

	var capturedCreateReq *teov20220901.CreateRealtimeLogDeliveryTaskRequest
	patches.ApplyMethodFunc(teoClient, "CreateRealtimeLogDeliveryTaskWithContext", func(_ context.Context, request *teov20220901.CreateRealtimeLogDeliveryTaskRequest) (*teov20220901.CreateRealtimeLogDeliveryTaskResponse, error) {
		capturedCreateReq = request
		assert.Equal(t, "zone-12345678", *request.ZoneId)
		assert.Equal(t, "aws_s3", *request.TaskType)
		// verify custom_expression_fields
		assert.NotNil(t, request.CustomExpressionFields)
		assert.Equal(t, 1, len(request.CustomExpressionFields))
		if len(request.CustomExpressionFields) > 0 {
			assert.Equal(t, "expr_field_1", *request.CustomExpressionFields[0].Name)
			assert.Equal(t, "$request_id", *request.CustomExpressionFields[0].Expression)
			assert.True(t, *request.CustomExpressionFields[0].Enabled)
		}
		// verify aws_s3
		assert.NotNil(t, request.AWSS3)
		if request.AWSS3 != nil {
			assert.Equal(t, "s3.us-west-2.amazonaws.com", *request.AWSS3.Endpoint)
			assert.Equal(t, "my-bucket", *request.AWSS3.BucketName)
			assert.Equal(t, "eo-logs/l7-access-logs", *request.AWSS3.Path)
			assert.Equal(t, "arn:aws:iam::123456789012:role/eo-logs", *request.AWSS3.RoleARN)
			assert.Equal(t, "gzip", *request.AWSS3.CompressType)
		}
		// verify log_format.record_template
		assert.NotNil(t, request.LogFormat)
		if request.LogFormat != nil {
			assert.Equal(t, "template", *request.LogFormat.FormatType)
			assert.NotNil(t, request.LogFormat.RecordTemplate)
			if request.LogFormat.RecordTemplate != nil {
				assert.Equal(t, "$time $host $request_uri", *request.LogFormat.RecordTemplate)
			}
		}
		resp := teov20220901.NewCreateRealtimeLogDeliveryTaskResponse()
		resp.Response = &teov20220901.CreateRealtimeLogDeliveryTaskResponseParams{
			TaskId:    ptrStrRTLD("task-abcdef"),
			RequestId: ptrStrRTLD("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(teoClient, "DescribeRealtimeLogDeliveryTasks", func(request *teov20220901.DescribeRealtimeLogDeliveryTasksRequest) (*teov20220901.DescribeRealtimeLogDeliveryTasksResponse, error) {
		resp := teov20220901.NewDescribeRealtimeLogDeliveryTasksResponse()
		resp.Response = &teov20220901.DescribeRealtimeLogDeliveryTasksResponseParams{
			TotalCount: uint64PtrRTLD(1),
			RealtimeLogDeliveryTasks: []*teov20220901.RealtimeLogDeliveryTask{
				{
					TaskId:         ptrStrRTLD("task-abcdef"),
					TaskName:       ptrStrRTLD("test-task"),
					TaskType:       ptrStrRTLD("aws_s3"),
					DeliveryStatus: ptrStrRTLD("enabled"),
					ZoneId:         ptrStrRTLD("zone-12345678"),
					LogType:        ptrStrRTLD("domain"),
					Area:           ptrStrRTLD("overseas"),
					EntityList:     []*string{ptrStrRTLD("example.com")},
					Fields:         []*string{ptrStrRTLD("ClientIP")},
					Sample:         uint64PtrRTLD(1000),
					CustomExpressionFields: []*teov20220901.CustomExpressionField{
						{
							Name:       ptrStrRTLD("expr_field_1"),
							Expression: ptrStrRTLD("$request_id"),
							Enabled:    ptrBoolRTLD(true),
						},
					},
					LogFormat: &teov20220901.LogFormat{
						FormatType:     ptrStrRTLD("template"),
						RecordTemplate: ptrStrRTLD("$time $host $request_uri"),
					},
					AWSS3: &teov20220901.AWSS3{
						Endpoint:     ptrStrRTLD("s3.us-west-2.amazonaws.com"),
						BucketName:   ptrStrRTLD("my-bucket"),
						Path:         ptrStrRTLD("eo-logs/l7-access-logs"),
						RoleARN:      ptrStrRTLD("arn:aws:iam::123456789012:role/eo-logs"),
						CompressType: ptrStrRTLD("gzip"),
					},
				},
			},
			RequestId: ptrStrRTLD("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaForRealtimeLogDelivery()
	res := teo.ResourceTencentCloudTeoRealtimeLogDelivery()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":     "zone-12345678",
		"task_name":   "test-task",
		"task_type":   "aws_s3",
		"log_type":    "domain",
		"area":        "overseas",
		"entity_list": []interface{}{"example.com"},
		"fields":      []interface{}{"ClientIP"},
		"sample":      1000,
		"custom_expression_fields": []interface{}{
			map[string]interface{}{
				"name":       "expr_field_1",
				"expression": "$request_id",
				"enabled":    true,
			},
		},
		"log_format": []interface{}{
			map[string]interface{}{
				"format_type":     "template",
				"record_template": "$time $host $request_uri",
			},
		},
		"aws_s3": []interface{}{
			map[string]interface{}{
				"endpoint":      "s3.us-west-2.amazonaws.com",
				"bucket_name":   "my-bucket",
				"path":          "eo-logs/l7-access-logs",
				"role_arn":      "arn:aws:iam::123456789012:role/eo-logs",
				"compress_type": "gzip",
			},
		},
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "zone-12345678#task-abcdef", d.Id())
	assert.NotNil(t, capturedCreateReq)
}

// TestRealtimeLogDelivery_Read_WithNewFields tests Read backfills custom_expression_fields, aws_s3, record_template
func TestRealtimeLogDelivery_Read_WithNewFields(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaForRealtimeLogDelivery().client, "UseTeoClient", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeRealtimeLogDeliveryTasks", func(request *teov20220901.DescribeRealtimeLogDeliveryTasksRequest) (*teov20220901.DescribeRealtimeLogDeliveryTasksResponse, error) {
		assert.Equal(t, "zone-12345678", *request.ZoneId)
		resp := teov20220901.NewDescribeRealtimeLogDeliveryTasksResponse()
		resp.Response = &teov20220901.DescribeRealtimeLogDeliveryTasksResponseParams{
			TotalCount: uint64PtrRTLD(1),
			RealtimeLogDeliveryTasks: []*teov20220901.RealtimeLogDeliveryTask{
				{
					TaskId:         ptrStrRTLD("task-abcdef"),
					TaskName:       ptrStrRTLD("test-task"),
					TaskType:       ptrStrRTLD("aws_s3"),
					DeliveryStatus: ptrStrRTLD("enabled"),
					ZoneId:         ptrStrRTLD("zone-12345678"),
					LogType:        ptrStrRTLD("domain"),
					Area:           ptrStrRTLD("overseas"),
					EntityList:     []*string{ptrStrRTLD("example.com")},
					Fields:         []*string{ptrStrRTLD("ClientIP")},
					Sample:         uint64PtrRTLD(1000),
					CustomExpressionFields: []*teov20220901.CustomExpressionField{
						{
							Name:       ptrStrRTLD("expr_field_1"),
							Expression: ptrStrRTLD("$request_id"),
							Enabled:    ptrBoolRTLD(true),
						},
					},
					LogFormat: &teov20220901.LogFormat{
						FormatType:     ptrStrRTLD("template"),
						RecordTemplate: ptrStrRTLD("$time $host $request_uri"),
					},
					AWSS3: &teov20220901.AWSS3{
						Endpoint:     ptrStrRTLD("s3.us-west-2.amazonaws.com"),
						BucketName:   ptrStrRTLD("my-bucket"),
						Path:         ptrStrRTLD("eo-logs/l7-access-logs"),
						RoleARN:      ptrStrRTLD("arn:aws:iam::123456789012:role/eo-logs"),
						CompressType: ptrStrRTLD("gzip"),
					},
				},
			},
			RequestId: ptrStrRTLD("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaForRealtimeLogDelivery()
	res := teo.ResourceTencentCloudTeoRealtimeLogDelivery()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
	d.SetId("zone-12345678#task-abcdef")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	// verify custom_expression_fields backfilled
	assert.Equal(t, 1, d.Get("custom_expression_fields.#"))
	assert.Equal(t, "expr_field_1", d.Get("custom_expression_fields.0.name"))
	assert.Equal(t, "$request_id", d.Get("custom_expression_fields.0.expression"))
	assert.Equal(t, true, d.Get("custom_expression_fields.0.enabled"))
	// verify aws_s3 backfilled
	assert.Equal(t, 1, d.Get("aws_s3.#"))
	assert.Equal(t, "s3.us-west-2.amazonaws.com", d.Get("aws_s3.0.endpoint"))
	assert.Equal(t, "my-bucket", d.Get("aws_s3.0.bucket_name"))
	assert.Equal(t, "eo-logs/l7-access-logs", d.Get("aws_s3.0.path"))
	assert.Equal(t, "arn:aws:iam::123456789012:role/eo-logs", d.Get("aws_s3.0.role_arn"))
	assert.Equal(t, "gzip", d.Get("aws_s3.0.compress_type"))
	// verify record_template backfilled
	assert.Equal(t, "template", d.Get("log_format.0.format_type"))
	assert.Equal(t, "$time $host $request_uri", d.Get("log_format.0.record_template"))
}

// TestRealtimeLogDelivery_Read_NullAWSS3 tests Read handles null aws_s3 without error
func TestRealtimeLogDelivery_Read_NullAWSS3(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaForRealtimeLogDelivery().client, "UseTeoClient", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeRealtimeLogDeliveryTasks", func(request *teov20220901.DescribeRealtimeLogDeliveryTasksRequest) (*teov20220901.DescribeRealtimeLogDeliveryTasksResponse, error) {
		resp := teov20220901.NewDescribeRealtimeLogDeliveryTasksResponse()
		resp.Response = &teov20220901.DescribeRealtimeLogDeliveryTasksResponseParams{
			TotalCount: uint64PtrRTLD(1),
			RealtimeLogDeliveryTasks: []*teov20220901.RealtimeLogDeliveryTask{
				{
					TaskId:         ptrStrRTLD("task-abcdef"),
					TaskName:       ptrStrRTLD("test-task"),
					TaskType:       ptrStrRTLD("cls"),
					DeliveryStatus: ptrStrRTLD("enabled"),
					ZoneId:         ptrStrRTLD("zone-12345678"),
					LogType:        ptrStrRTLD("domain"),
					Area:           ptrStrRTLD("overseas"),
					EntityList:     []*string{ptrStrRTLD("example.com")},
					Fields:         []*string{ptrStrRTLD("ClientIP")},
					Sample:         uint64PtrRTLD(1000),
					AWSS3:          nil,
				},
			},
			RequestId: ptrStrRTLD("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaForRealtimeLogDelivery()
	res := teo.ResourceTencentCloudTeoRealtimeLogDelivery()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
	d.SetId("zone-12345678#task-abcdef")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, 0, d.Get("aws_s3.#"))
}

// TestRealtimeLogDelivery_Update_WithNewFields tests Update with custom_expression_fields, aws_s3, log_format.record_template
func TestRealtimeLogDelivery_Update_WithNewFields(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaForRealtimeLogDelivery().client, "UseTeoClient", teoClient)

	var capturedModifyReq *teov20220901.ModifyRealtimeLogDeliveryTaskRequest
	patches.ApplyMethodFunc(teoClient, "ModifyRealtimeLogDeliveryTaskWithContext", func(_ context.Context, request *teov20220901.ModifyRealtimeLogDeliveryTaskRequest) (*teov20220901.ModifyRealtimeLogDeliveryTaskResponse, error) {
		capturedModifyReq = request
		assert.Equal(t, "zone-12345678", *request.ZoneId)
		assert.Equal(t, "task-abcdef", *request.TaskId)
		// verify custom_expression_fields
		assert.NotNil(t, request.CustomExpressionFields)
		if len(request.CustomExpressionFields) > 0 {
			assert.Equal(t, "expr_field_1", *request.CustomExpressionFields[0].Name)
			assert.Equal(t, "$request_id", *request.CustomExpressionFields[0].Expression)
			assert.True(t, *request.CustomExpressionFields[0].Enabled)
		}
		// verify aws_s3
		assert.NotNil(t, request.AWSS3)
		if request.AWSS3 != nil {
			assert.Equal(t, "s3.us-west-2.amazonaws.com", *request.AWSS3.Endpoint)
			assert.Equal(t, "my-bucket-updated", *request.AWSS3.BucketName)
		}
		// verify log_format.record_template
		assert.NotNil(t, request.LogFormat)
		if request.LogFormat != nil {
			assert.Equal(t, "template", *request.LogFormat.FormatType)
			assert.NotNil(t, request.LogFormat.RecordTemplate)
		}
		resp := teov20220901.NewModifyRealtimeLogDeliveryTaskResponse()
		resp.Response = &teov20220901.ModifyRealtimeLogDeliveryTaskResponseParams{
			RequestId: ptrStrRTLD("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(teoClient, "DescribeRealtimeLogDeliveryTasks", func(request *teov20220901.DescribeRealtimeLogDeliveryTasksRequest) (*teov20220901.DescribeRealtimeLogDeliveryTasksResponse, error) {
		resp := teov20220901.NewDescribeRealtimeLogDeliveryTasksResponse()
		resp.Response = &teov20220901.DescribeRealtimeLogDeliveryTasksResponseParams{
			TotalCount: uint64PtrRTLD(1),
			RealtimeLogDeliveryTasks: []*teov20220901.RealtimeLogDeliveryTask{
				{
					TaskId:         ptrStrRTLD("task-abcdef"),
					TaskName:       ptrStrRTLD("test-task"),
					TaskType:       ptrStrRTLD("aws_s3"),
					DeliveryStatus: ptrStrRTLD("enabled"),
					ZoneId:         ptrStrRTLD("zone-12345678"),
					LogType:        ptrStrRTLD("domain"),
					Area:           ptrStrRTLD("overseas"),
					EntityList:     []*string{ptrStrRTLD("example.com")},
					Fields:         []*string{ptrStrRTLD("ClientIP")},
					Sample:         uint64PtrRTLD(1000),
					CustomExpressionFields: []*teov20220901.CustomExpressionField{
						{
							Name:       ptrStrRTLD("expr_field_1"),
							Expression: ptrStrRTLD("$request_id"),
							Enabled:    ptrBoolRTLD(true),
						},
					},
					LogFormat: &teov20220901.LogFormat{
						FormatType:     ptrStrRTLD("template"),
						RecordTemplate: ptrStrRTLD("$time $host $request_uri"),
					},
					AWSS3: &teov20220901.AWSS3{
						Endpoint:     ptrStrRTLD("s3.us-west-2.amazonaws.com"),
						BucketName:   ptrStrRTLD("my-bucket-updated"),
						Path:         ptrStrRTLD("eo-logs/l7-access-logs"),
						RoleARN:      ptrStrRTLD("arn:aws:iam::123456789012:role/eo-logs"),
						CompressType: ptrStrRTLD("gzip"),
					},
				},
			},
			RequestId: ptrStrRTLD("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaForRealtimeLogDelivery()
	res := teo.ResourceTencentCloudTeoRealtimeLogDelivery()

	// Set up state with new values to trigger Update
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":         "zone-12345678",
		"task_name":       "test-task",
		"task_type":       "aws_s3",
		"log_type":        "domain",
		"area":            "overseas",
		"entity_list":     []interface{}{"example.com"},
		"fields":          []interface{}{"ClientIP"},
		"sample":          1000,
		"task_id":         "task-abcdef",
		"delivery_status": "enabled",
		"custom_expression_fields": []interface{}{
			map[string]interface{}{
				"name":       "expr_field_1",
				"expression": "$request_id",
				"enabled":    true,
			},
		},
		"log_format": []interface{}{
			map[string]interface{}{
				"format_type":     "template",
				"record_template": "$time $host $request_uri",
			},
		},
		"aws_s3": []interface{}{
			map[string]interface{}{
				"endpoint":      "s3.us-west-2.amazonaws.com",
				"bucket_name":   "my-bucket-updated",
				"path":          "eo-logs/l7-access-logs",
				"role_arn":      "arn:aws:iam::123456789012:role/eo-logs",
				"compress_type": "gzip",
			},
		},
	})
	d.SetId("zone-12345678#task-abcdef")

	// Manually trigger HasChange by simulating the new state differs from old
	// Since TestResourceDataRaw sets the current state, we need to mark changes
	// Use d.Set to update values which will register as changes from the "old" state
	_ = d.Set("custom_expression_fields", []interface{}{
		map[string]interface{}{
			"name":       "expr_field_1",
			"expression": "$request_id_updated",
			"enabled":    true,
		},
	})

	err := res.Update(d, meta)
	assert.NoError(t, err)
	assert.NotNil(t, capturedModifyReq)
}

func uint64PtrRTLD(v uint64) *uint64 {
	return &v
}
