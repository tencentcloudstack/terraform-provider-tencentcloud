package cls_test

import (
	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	localcls "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/cls"

	"context"
	"fmt"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/stretchr/testify/assert"
	cls "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cls/v20201016"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
)

// go test -i; go test -test.run TestAccTencentCloudClsScheduledSqlResource_basic -v
func TestAccTencentCloudClsScheduledSqlResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		CheckDestroy: testAccCheckClsScheduledSqlDestroy,
		Providers:    tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccClsScheduledSql,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClsScheduledSqlExists("tencentcloud_cls_scheduled_sql.scheduled_sql"),
					resource.TestCheckResourceAttrSet("tencentcloud_cls_scheduled_sql.scheduled_sql", "src_topic_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_cls_scheduled_sql.scheduled_sql", "name"),
					resource.TestCheckResourceAttrSet("tencentcloud_cls_scheduled_sql.scheduled_sql", "enable_flag"),
					resource.TestCheckResourceAttrSet("tencentcloud_cls_scheduled_sql.scheduled_sql", "dst_resource.#"),
					resource.TestCheckResourceAttrSet("tencentcloud_cls_scheduled_sql.scheduled_sql", "scheduled_sql_content"),
					resource.TestCheckResourceAttrSet("tencentcloud_cls_scheduled_sql.scheduled_sql", "process_start_time"),
					resource.TestCheckResourceAttrSet("tencentcloud_cls_scheduled_sql.scheduled_sql", "process_type"),
					resource.TestCheckResourceAttrSet("tencentcloud_cls_scheduled_sql.scheduled_sql", "process_time_window"),
					resource.TestCheckResourceAttrSet("tencentcloud_cls_scheduled_sql.scheduled_sql", "process_delay")),
			},
			{
				ResourceName:      "tencentcloud_cls_scheduled_sql.scheduled_sql",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckClsScheduledSqlDestroy(s *terraform.State) error {
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	clsService := localcls.NewClsService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_cls_scheduled_sql" {
			continue
		}
		instance, err := clsService.DescribeClsScheduledSqlById(ctx, rs.Primary.ID)
		if err != nil {
			continue
		}
		if instance != nil {
			return fmt.Errorf("[CHECK][CLS ScheduledSql][Destroy] check: CLS ScheduledSql still exists: %s", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckClsScheduledSqlExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := tccommon.GetLogId(tccommon.ContextNil)
		ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("[CHECK][CLS ScheduledSql][Exists] check: CLS ScheduledSql %s is not found", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("[CHECK][CLS ScheduledSql][Create] check: CLS ScheduledSql id is not set")
		}
		clsService := localcls.NewClsService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
		taskRes, err := clsService.DescribeClsScheduledSqlById(ctx, rs.Primary.ID)
		if err != nil {
			return err
		}
		if taskRes == nil {
			return fmt.Errorf("[CHECK][CLS ScheduledSql][Exists] id %s is not exist", rs.Primary.ID)
		}
		return nil
	}
}

const testAccClsScheduledSql = `

resource "tencentcloud_cls_logset" "logset" {
  logset_name = "tf-example"
  tags = {
    "createdBy" = "terraform"
  }
}
resource "tencentcloud_cls_topic" "topic" {
  topic_name           = "tf-example"
  logset_id            = tencentcloud_cls_logset.logset.id
  auto_split           = false
  max_split_partitions = 20
  partition_count      = 1
  period               = 10
  storage_type         = "hot"
  tags                 = {
    "test" = "test",
  }
}
resource "tencentcloud_cls_scheduled_sql" "scheduled_sql" {
  src_topic_id = tencentcloud_cls_topic.topic.id
  name = "tf-example"
  enable_flag = 1
  dst_resource {
    topic_id = tencentcloud_cls_topic.topic.id
    region = "ap-guangzhou"
    biz_type = 0
    metric_name = "test"

  }
  scheduled_sql_content = "xxx"
  process_start_time = 1723117637000
  process_type = 1
  process_period = 10
  process_time_window = "@m-15m,@m"
  process_delay = 5
  src_topic_region = "ap-guangzhou"
  syntax_rule = 0
}
`

// mockMetaForScheduledSql is a mock ProviderMeta for scheduled_sql unit tests
type mockMetaForScheduledSql struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMetaForScheduledSql) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockMetaForScheduledSql{}

func newMockMetaForScheduledSql() *mockMetaForScheduledSql {
	return &mockMetaForScheduledSql{client: &connectivity.TencentCloudClient{}}
}

func ptrStrScheduledSql(s string) *string {
	return &s
}

func ptrInt64ScheduledSql(v int64) *int64 {
	return &v
}

func ptrUint64ScheduledSql(v uint64) *uint64 {
	return &v
}

// go test ./tencentcloud/services/cls/ -run "TestScheduledSqlDstResourceParams" -v -count=1 -gcflags="all=-l"

// TestScheduledSqlDstResourceParams_Create tests that Create correctly populates the new dst_resource fields in the request
func TestScheduledSqlDstResourceParams_Create(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	clsClient := &cls.Client{}
	patches.ApplyMethodReturn(newMockMetaForScheduledSql().client, "UseClsClient", clsClient)

	var capturedRequest *cls.CreateScheduledSqlRequest
	patches.ApplyMethodFunc(clsClient, "CreateScheduledSql", func(request *cls.CreateScheduledSqlRequest) (*cls.CreateScheduledSqlResponse, error) {
		capturedRequest = request
		resp := cls.NewCreateScheduledSqlResponse()
		resp.Response = &cls.CreateScheduledSqlResponseParams{
			TaskId:    ptrStrScheduledSql("task-test-metric-params"),
			RequestId: ptrStrScheduledSql("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(clsClient, "DescribeScheduledSqlInfo", func(request *cls.DescribeScheduledSqlInfoRequest) (*cls.DescribeScheduledSqlInfoResponse, error) {
		resp := cls.NewDescribeScheduledSqlInfoResponse()
		resp.Response = &cls.DescribeScheduledSqlInfoResponseParams{
			ScheduledSqlTaskInfos: []*cls.ScheduledSqlTaskInfo{
				{
					TaskId:     ptrStrScheduledSql("task-test-metric-params"),
					Name:       ptrStrScheduledSql("tf-example-metric-params"),
					SrcTopicId: ptrStrScheduledSql("topic-src-123"),
					EnableFlag: ptrInt64ScheduledSql(1),
					DstResource: &cls.ScheduledSqlResouceInfo{
						TopicId:    ptrStrScheduledSql("topic-dst-123"),
						Region:     ptrStrScheduledSql("ap-guangzhou"),
						BizType:    ptrInt64ScheduledSql(1),
						MetricName: ptrStrScheduledSql("metric1"),
						MetricNames: []*string{
							ptrStrScheduledSql("metric1"),
							ptrStrScheduledSql("metric2"),
						},
						MetricLabels: []*string{
							ptrStrScheduledSql("label1"),
							ptrStrScheduledSql("label2"),
						},
						CustomTime: ptrStrScheduledSql("timestamp"),
						CustomMetricLabels: []*cls.MetricLabel{
							{Key: ptrStrScheduledSql("env"), Value: ptrStrScheduledSql("production")},
							{Key: ptrStrScheduledSql("app"), Value: ptrStrScheduledSql("myapp")},
						},
					},
					ScheduledSqlContent: ptrStrScheduledSql("select * from log"),
					ProcessStartTime:    ptrStrScheduledSql("2024-01-01 00:00:00"),
					ProcessType:         ptrInt64ScheduledSql(1),
					ProcessPeriod:       ptrInt64ScheduledSql(10),
					ProcessTimeWindow:   ptrStrScheduledSql("@m-15m,@m"),
					ProcessDelay:        ptrInt64ScheduledSql(5),
					SrcTopicRegion:      ptrStrScheduledSql("ap-guangzhou"),
					SyntaxRule:          ptrUint64ScheduledSql(0),
				},
			},
			TotalCount: ptrUint64ScheduledSql(1),
			RequestId:  ptrStrScheduledSql("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaForScheduledSql()
	res := localcls.ResourceTencentCloudClsScheduledSql()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"src_topic_id":          "topic-src-123",
		"name":                  "tf-example-metric-params",
		"enable_flag":           1,
		"scheduled_sql_content": "select * from log",
		"process_start_time":    1704067200000,
		"process_type":          1,
		"process_period":        10,
		"process_time_window":   "@m-15m,@m",
		"process_delay":         5,
		"src_topic_region":      "ap-guangzhou",
		"syntax_rule":           0,
		"dst_resource": []interface{}{
			map[string]interface{}{
				"topic_id":    "topic-dst-123",
				"region":      "ap-guangzhou",
				"biz_type":    1,
				"metric_name": "metric1",
				"metric_names": []interface{}{
					"metric1",
					"metric2",
				},
				"metric_labels": []interface{}{
					"label1",
					"label2",
				},
				"custom_time": "timestamp",
				"custom_metric_labels": []interface{}{
					map[string]interface{}{
						"key":   "env",
						"value": "production",
					},
					map[string]interface{}{
						"key":   "app",
						"value": "myapp",
					},
				},
			},
		},
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "task-test-metric-params", d.Id())
	assert.NotNil(t, capturedRequest)
	assert.NotNil(t, capturedRequest.DstResource)

	dst := capturedRequest.DstResource
	assert.NotNil(t, dst.MetricNames)
	assert.Len(t, dst.MetricNames, 2)
	assert.Equal(t, "metric1", *dst.MetricNames[0])
	assert.Equal(t, "metric2", *dst.MetricNames[1])

	assert.NotNil(t, dst.MetricLabels)
	assert.Len(t, dst.MetricLabels, 2)
	assert.Equal(t, "label1", *dst.MetricLabels[0])
	assert.Equal(t, "label2", *dst.MetricLabels[1])

	assert.NotNil(t, dst.CustomTime)
	assert.Equal(t, "timestamp", *dst.CustomTime)

	assert.NotNil(t, dst.CustomMetricLabels)
	assert.Len(t, dst.CustomMetricLabels, 2)
	assert.Equal(t, "env", *dst.CustomMetricLabels[0].Key)
	assert.Equal(t, "production", *dst.CustomMetricLabels[0].Value)
	assert.Equal(t, "app", *dst.CustomMetricLabels[1].Key)
	assert.Equal(t, "myapp", *dst.CustomMetricLabels[1].Value)
}

// TestScheduledSqlDstResourceParams_Read tests that Read populates the new dst_resource fields from the Describe response
func TestScheduledSqlDstResourceParams_Read(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	clsClient := &cls.Client{}
	patches.ApplyMethodReturn(newMockMetaForScheduledSql().client, "UseClsClient", clsClient)

	patches.ApplyMethodFunc(clsClient, "DescribeScheduledSqlInfo", func(request *cls.DescribeScheduledSqlInfoRequest) (*cls.DescribeScheduledSqlInfoResponse, error) {
		resp := cls.NewDescribeScheduledSqlInfoResponse()
		resp.Response = &cls.DescribeScheduledSqlInfoResponseParams{
			ScheduledSqlTaskInfos: []*cls.ScheduledSqlTaskInfo{
				{
					TaskId:     ptrStrScheduledSql("task-test-metric-params"),
					Name:       ptrStrScheduledSql("tf-example-metric-params"),
					SrcTopicId: ptrStrScheduledSql("topic-src-123"),
					EnableFlag: ptrInt64ScheduledSql(1),
					DstResource: &cls.ScheduledSqlResouceInfo{
						TopicId:    ptrStrScheduledSql("topic-dst-123"),
						Region:     ptrStrScheduledSql("ap-guangzhou"),
						BizType:    ptrInt64ScheduledSql(1),
						MetricName: ptrStrScheduledSql("metric1"),
						MetricNames: []*string{
							ptrStrScheduledSql("metric1"),
							ptrStrScheduledSql("metric2"),
						},
						MetricLabels: []*string{
							ptrStrScheduledSql("label1"),
							ptrStrScheduledSql("label2"),
						},
						CustomTime: ptrStrScheduledSql("timestamp"),
						CustomMetricLabels: []*cls.MetricLabel{
							{Key: ptrStrScheduledSql("env"), Value: ptrStrScheduledSql("production")},
							{Key: ptrStrScheduledSql("app"), Value: ptrStrScheduledSql("myapp")},
						},
					},
					ScheduledSqlContent: ptrStrScheduledSql("select * from log"),
					ProcessStartTime:    ptrStrScheduledSql("2024-01-01 00:00:00"),
					ProcessType:         ptrInt64ScheduledSql(1),
					ProcessPeriod:       ptrInt64ScheduledSql(10),
					ProcessTimeWindow:   ptrStrScheduledSql("@m-15m,@m"),
					ProcessDelay:        ptrInt64ScheduledSql(5),
					SrcTopicRegion:      ptrStrScheduledSql("ap-guangzhou"),
					SyntaxRule:          ptrUint64ScheduledSql(0),
				},
			},
			TotalCount: ptrUint64ScheduledSql(1),
			RequestId:  ptrStrScheduledSql("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaForScheduledSql()
	res := localcls.ResourceTencentCloudClsScheduledSql()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"src_topic_id":          "topic-src-123",
		"name":                  "tf-example-metric-params",
		"enable_flag":           1,
		"scheduled_sql_content": "select * from log",
		"process_start_time":    1704067200000,
		"process_type":          1,
		"process_period":        10,
		"process_time_window":   "@m-15m,@m",
		"process_delay":         5,
		"src_topic_region":      "ap-guangzhou",
		"syntax_rule":           0,
	})
	d.SetId("task-test-metric-params")

	err := res.Read(d, meta)
	assert.NoError(t, err)

	dstResources := d.Get("dst_resource").([]interface{})
	assert.Len(t, dstResources, 1)

	dstMap := dstResources[0].(map[string]interface{})
	assert.Equal(t, "topic-dst-123", dstMap["topic_id"])
	assert.Equal(t, "ap-guangzhou", dstMap["region"])
	assert.Equal(t, 1, dstMap["biz_type"])
	assert.Equal(t, "metric1", dstMap["metric_name"])

	metricNames := dstMap["metric_names"].([]interface{})
	assert.Len(t, metricNames, 2)
	assert.Equal(t, "metric1", metricNames[0])
	assert.Equal(t, "metric2", metricNames[1])

	metricLabels := dstMap["metric_labels"].([]interface{})
	assert.Len(t, metricLabels, 2)
	assert.Equal(t, "label1", metricLabels[0])
	assert.Equal(t, "label2", metricLabels[1])

	assert.Equal(t, "timestamp", dstMap["custom_time"])

	customMetricLabels := dstMap["custom_metric_labels"].([]interface{})
	assert.Len(t, customMetricLabels, 2)
	label0 := customMetricLabels[0].(map[string]interface{})
	assert.Equal(t, "env", label0["key"])
	assert.Equal(t, "production", label0["value"])
	label1 := customMetricLabels[1].(map[string]interface{})
	assert.Equal(t, "app", label1["key"])
	assert.Equal(t, "myapp", label1["value"])
}

// TestScheduledSqlDstResourceParams_Read_NilFields tests that Read handles nil new dst_resource fields without panic
func TestScheduledSqlDstResourceParams_Read_NilFields(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	clsClient := &cls.Client{}
	patches.ApplyMethodReturn(newMockMetaForScheduledSql().client, "UseClsClient", clsClient)

	patches.ApplyMethodFunc(clsClient, "DescribeScheduledSqlInfo", func(request *cls.DescribeScheduledSqlInfoRequest) (*cls.DescribeScheduledSqlInfoResponse, error) {
		resp := cls.NewDescribeScheduledSqlInfoResponse()
		resp.Response = &cls.DescribeScheduledSqlInfoResponseParams{
			ScheduledSqlTaskInfos: []*cls.ScheduledSqlTaskInfo{
				{
					TaskId:     ptrStrScheduledSql("task-test-nil-fields"),
					Name:       ptrStrScheduledSql("tf-example-nil-fields"),
					SrcTopicId: ptrStrScheduledSql("topic-src-123"),
					EnableFlag: ptrInt64ScheduledSql(1),
					DstResource: &cls.ScheduledSqlResouceInfo{
						TopicId: ptrStrScheduledSql("topic-dst-123"),
						Region:  ptrStrScheduledSql("ap-guangzhou"),
						BizType: ptrInt64ScheduledSql(0),
					},
					ScheduledSqlContent: ptrStrScheduledSql("select * from log"),
					ProcessStartTime:    ptrStrScheduledSql("2024-01-01 00:00:00"),
					ProcessType:         ptrInt64ScheduledSql(1),
					ProcessPeriod:       ptrInt64ScheduledSql(10),
					ProcessTimeWindow:   ptrStrScheduledSql("@m-15m,@m"),
					ProcessDelay:        ptrInt64ScheduledSql(5),
					SrcTopicRegion:      ptrStrScheduledSql("ap-guangzhou"),
					SyntaxRule:          ptrUint64ScheduledSql(0),
				},
			},
			TotalCount: ptrUint64ScheduledSql(1),
			RequestId:  ptrStrScheduledSql("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaForScheduledSql()
	res := localcls.ResourceTencentCloudClsScheduledSql()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"src_topic_id":          "topic-src-123",
		"name":                  "tf-example-nil-fields",
		"enable_flag":           1,
		"scheduled_sql_content": "select * from log",
		"process_start_time":    1704067200000,
		"process_type":          1,
		"process_period":        10,
		"process_time_window":   "@m-15m,@m",
		"process_delay":         5,
		"src_topic_region":      "ap-guangzhou",
		"syntax_rule":           0,
	})
	d.SetId("task-test-nil-fields")

	err := res.Read(d, meta)
	assert.NoError(t, err)

	dstResources := d.Get("dst_resource").([]interface{})
	assert.Len(t, dstResources, 1)

	dstMap := dstResources[0].(map[string]interface{})
	assert.Equal(t, "topic-dst-123", dstMap["topic_id"])
	assert.Equal(t, "ap-guangzhou", dstMap["region"])
	assert.Equal(t, 0, dstMap["biz_type"])
	_, exists := dstMap["metric_names"]
	assert.False(t, exists)
	_, exists = dstMap["metric_labels"]
	assert.False(t, exists)
	_, exists = dstMap["custom_time"]
	assert.False(t, exists)
	_, exists = dstMap["custom_metric_labels"]
	assert.False(t, exists)
}

// TestScheduledSqlDstResourceParams_Update tests that Update correctly populates the new dst_resource fields in the modify request
func TestScheduledSqlDstResourceParams_Update(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	clsClient := &cls.Client{}
	patches.ApplyMethodReturn(newMockMetaForScheduledSql().client, "UseClsClient", clsClient)

	var capturedModifyRequest *cls.ModifyScheduledSqlRequest
	patches.ApplyMethodFunc(clsClient, "ModifyScheduledSql", func(request *cls.ModifyScheduledSqlRequest) (*cls.ModifyScheduledSqlResponse, error) {
		capturedModifyRequest = request
		resp := cls.NewModifyScheduledSqlResponse()
		resp.Response = &cls.ModifyScheduledSqlResponseParams{
			RequestId: ptrStrScheduledSql("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(clsClient, "DescribeScheduledSqlInfo", func(request *cls.DescribeScheduledSqlInfoRequest) (*cls.DescribeScheduledSqlInfoResponse, error) {
		resp := cls.NewDescribeScheduledSqlInfoResponse()
		resp.Response = &cls.DescribeScheduledSqlInfoResponseParams{
			ScheduledSqlTaskInfos: []*cls.ScheduledSqlTaskInfo{
				{
					TaskId:     ptrStrScheduledSql("task-test-metric-params"),
					Name:       ptrStrScheduledSql("tf-example-metric-params"),
					SrcTopicId: ptrStrScheduledSql("topic-src-123"),
					EnableFlag: ptrInt64ScheduledSql(1),
					DstResource: &cls.ScheduledSqlResouceInfo{
						TopicId:    ptrStrScheduledSql("topic-dst-123"),
						Region:     ptrStrScheduledSql("ap-guangzhou"),
						BizType:    ptrInt64ScheduledSql(1),
						MetricName: ptrStrScheduledSql("metric1"),
						MetricNames: []*string{
							ptrStrScheduledSql("metric1"),
							ptrStrScheduledSql("metric2"),
						},
						MetricLabels: []*string{
							ptrStrScheduledSql("label1"),
							ptrStrScheduledSql("label2"),
						},
						CustomTime: ptrStrScheduledSql("timestamp"),
						CustomMetricLabels: []*cls.MetricLabel{
							{Key: ptrStrScheduledSql("env"), Value: ptrStrScheduledSql("production")},
							{Key: ptrStrScheduledSql("app"), Value: ptrStrScheduledSql("myapp")},
						},
					},
					ScheduledSqlContent: ptrStrScheduledSql("select * from log"),
					ProcessStartTime:    ptrStrScheduledSql("2024-01-01 00:00:00"),
					ProcessType:         ptrInt64ScheduledSql(1),
					ProcessPeriod:       ptrInt64ScheduledSql(10),
					ProcessTimeWindow:   ptrStrScheduledSql("@m-15m,@m"),
					ProcessDelay:        ptrInt64ScheduledSql(5),
					SrcTopicRegion:      ptrStrScheduledSql("ap-guangzhou"),
					SyntaxRule:          ptrUint64ScheduledSql(0),
				},
			},
			TotalCount: ptrUint64ScheduledSql(1),
			RequestId:  ptrStrScheduledSql("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaForScheduledSql()
	res := localcls.ResourceTencentCloudClsScheduledSql()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"src_topic_id":          "topic-src-123",
		"name":                  "tf-example-metric-params",
		"enable_flag":           1,
		"scheduled_sql_content": "select * from log",
		"process_start_time":    1704067200000,
		"process_type":          1,
		"process_period":        10,
		"process_time_window":   "@m-15m,@m",
		"process_delay":         5,
		"src_topic_region":      "ap-guangzhou",
		"syntax_rule":           0,
		"dst_resource": []interface{}{
			map[string]interface{}{
				"topic_id":    "topic-dst-123",
				"region":      "ap-guangzhou",
				"biz_type":    1,
				"metric_name": "metric1",
				"metric_names": []interface{}{
					"metric1",
					"metric2",
				},
				"metric_labels": []interface{}{
					"label1",
					"label2",
				},
				"custom_time": "timestamp",
				"custom_metric_labels": []interface{}{
					map[string]interface{}{
						"key":   "env",
						"value": "production",
					},
					map[string]interface{}{
						"key":   "app",
						"value": "myapp",
					},
				},
			},
		},
	})
	d.SetId("task-test-metric-params")

	err := res.Update(d, meta)
	assert.NoError(t, err)
	assert.NotNil(t, capturedModifyRequest)
	assert.NotNil(t, capturedModifyRequest.DstResource)

	dst := capturedModifyRequest.DstResource
	assert.NotNil(t, dst.MetricNames)
	assert.Len(t, dst.MetricNames, 2)
	assert.Equal(t, "metric1", *dst.MetricNames[0])
	assert.Equal(t, "metric2", *dst.MetricNames[1])

	assert.NotNil(t, dst.MetricLabels)
	assert.Len(t, dst.MetricLabels, 2)
	assert.Equal(t, "label1", *dst.MetricLabels[0])
	assert.Equal(t, "label2", *dst.MetricLabels[1])

	assert.NotNil(t, dst.CustomTime)
	assert.Equal(t, "timestamp", *dst.CustomTime)

	assert.NotNil(t, dst.CustomMetricLabels)
	assert.Len(t, dst.CustomMetricLabels, 2)
	assert.Equal(t, "env", *dst.CustomMetricLabels[0].Key)
	assert.Equal(t, "production", *dst.CustomMetricLabels[0].Value)
	assert.Equal(t, "app", *dst.CustomMetricLabels[1].Key)
	assert.Equal(t, "myapp", *dst.CustomMetricLabels[1].Value)
}

// TestScheduledSqlDstResourceParams_Schema tests the new schema definitions in the dst_resource block
func TestScheduledSqlDstResourceParams_Schema(t *testing.T) {
	res := localcls.ResourceTencentCloudClsScheduledSql()

	assert.NotNil(t, res)
	assert.Contains(t, res.Schema, "dst_resource")

	dstResourceSchema := res.Schema["dst_resource"]
	assert.NotNil(t, dstResourceSchema.Elem)

	elemResource, ok := dstResourceSchema.Elem.(*schema.Resource)
	assert.True(t, ok)
	assert.NotNil(t, elemResource.Schema)

	metricNamesSchema, ok := elemResource.Schema["metric_names"]
	assert.True(t, ok)
	assert.Equal(t, schema.TypeList, metricNamesSchema.Type)
	assert.True(t, metricNamesSchema.Optional)
	assert.False(t, metricNamesSchema.Required)

	metricLabelsSchema, ok := elemResource.Schema["metric_labels"]
	assert.True(t, ok)
	assert.Equal(t, schema.TypeList, metricLabelsSchema.Type)
	assert.True(t, metricLabelsSchema.Optional)
	assert.False(t, metricLabelsSchema.Required)

	customTimeSchema, ok := elemResource.Schema["custom_time"]
	assert.True(t, ok)
	assert.Equal(t, schema.TypeString, customTimeSchema.Type)
	assert.True(t, customTimeSchema.Optional)
	assert.False(t, customTimeSchema.Required)

	customMetricLabelsSchema, ok := elemResource.Schema["custom_metric_labels"]
	assert.True(t, ok)
	assert.Equal(t, schema.TypeList, customMetricLabelsSchema.Type)
	assert.True(t, customMetricLabelsSchema.Optional)
	assert.False(t, customMetricLabelsSchema.Required)

	labelElem, ok := customMetricLabelsSchema.Elem.(*schema.Resource)
	assert.True(t, ok)
	assert.Contains(t, labelElem.Schema, "key")
	assert.Contains(t, labelElem.Schema, "value")
	assert.Equal(t, schema.TypeString, labelElem.Schema["key"].Type)
	assert.Equal(t, schema.TypeString, labelElem.Schema["value"].Type)
}
