package clb_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/stretchr/testify/assert"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	localclb "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/clb"
	localcls "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/cls"

	clbv20180317 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/clb/v20180317"
	clsv20201016 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cls/v20201016"
)

func TestAccTencentCloudClbInstanceTopic(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheck(t) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckClbListenerRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccClbInstanceTopic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClbInstanceTopicExists("tencentcloud_clb_log_topic.topic"),
					resource.TestCheckResourceAttr("tencentcloud_clb_log_topic.topic", "topic_name", "clb-topic-test"),
				),
			},
			{
				Config: testAccClbInstanceTopicWithTags,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClbInstanceTopicExists("tencentcloud_clb_log_topic.topic"),
					resource.TestCheckResourceAttr("tencentcloud_clb_log_topic.topic", "topic_name", "clb-topic-test"),
					resource.TestCheckResourceAttr("tencentcloud_clb_log_topic.topic", "tags.env", "prod"),
				),
			},
			{
				Config: testAccClbInstanceTopicWithTagsUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClbInstanceTopicExists("tencentcloud_clb_log_topic.topic"),
					resource.TestCheckResourceAttr("tencentcloud_clb_log_topic.topic", "topic_name", "clb-topic-test"),
					resource.TestCheckResourceAttr("tencentcloud_clb_log_topic.topic", "tags.team", "dev"),
				),
			},
			{
				Config: testAccClbInstanceTopic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClbInstanceTopicExists("tencentcloud_clb_log_topic.topic"),
					resource.TestCheckResourceAttr("tencentcloud_clb_log_topic.topic", "topic_name", "clb-topic-test"),
					resource.TestCheckResourceAttr("tencentcloud_clb_log_topic.topic", "tags.%", "0"),
				),
			},
		},
	})
}

func testAccCheckClbInstanceTopicExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := tccommon.GetLogId(tccommon.ContextNil)
		ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("[CHECK][CLB topic][Exists] check: CLB topic %s is not found", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("[CHECK][CLB topic][Exists] check: CLB topic id is not set")
		}
		clsService := localcls.NewClsService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
		instance, err := clsService.DescribeClsTopicById(ctx, rs.Primary.ID, nil)
		if err != nil {
			return err
		}

		if instance == nil {
			return fmt.Errorf("[CHECK][CLB topic][Exists] id %s is not exist", rs.Primary.ID)
		}
		return nil
	}
}

const testAccClbInstanceTopic = `
resource "tencentcloud_clb_log_set" "set1" {
    period = 7
}

resource "tencentcloud_clb_log_topic" "topic" {
    log_set_id = tencentcloud_clb_log_set.set1.id
    topic_name="clb-topic-test"
}
`

const testAccClbInstanceTopicWithTags = `
resource "tencentcloud_clb_log_set" "set1" {
    period = 7
}

resource "tencentcloud_clb_log_topic" "topic" {
    log_set_id = tencentcloud_clb_log_set.set1.id
    topic_name = "clb-topic-test"
    tags = {
      env = "prod"
    }
}
`

const testAccClbInstanceTopicWithTagsUpdate = `
resource "tencentcloud_clb_log_set" "set1" {
    period = 7
}

resource "tencentcloud_clb_log_topic" "topic" {
    log_set_id = tencentcloud_clb_log_set.set1.id
    topic_name = "clb-topic-test"
    tags = {
      team = "dev"
    }
}
`

// --- gomonkey mock unit tests for period parameter ---

type mockMetaForClbLogTopic struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMetaForClbLogTopic) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockMetaForClbLogTopic{}

func newMockMetaForClbLogTopic() *mockMetaForClbLogTopic {
	return &mockMetaForClbLogTopic{client: &connectivity.TencentCloudClient{}}
}

func ptrStringCLT(s string) *string { return &s }
func ptrInt64CLT(i int64) *int64    { return &i }
func ptrBoolCLT(b bool) *bool       { return &b }

// TestClbLogTopic_Create_WithPeriod verifies that when period is set, the
// CreateTopic request carries the correct Period value (cast to *uint64).
func TestClbLogTopic_Create_WithPeriod(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	clbClient := &clbv20180317.Client{}
	patches.ApplyMethodReturn(newMockMetaForClbLogTopic().client, "UseClbClient", clbClient)

	clsClient := &clsv20201016.Client{}
	patches.ApplyMethodReturn(newMockMetaForClbLogTopic().client, "UseClsClient", clsClient)

	var capturedRequest *clbv20180317.CreateTopicRequest
	patches.ApplyMethodFunc(clbClient, "CreateTopic", func(request *clbv20180317.CreateTopicRequest) (*clbv20180317.CreateTopicResponse, error) {
		capturedRequest = request
		resp := clbv20180317.NewCreateTopicResponse()
		resp.Response = &clbv20180317.CreateTopicResponseParams{
			TopicId:   ptrStringCLT("fake-topic-id"),
			RequestId: ptrStringCLT("fake-request-id"),
		}
		return resp, nil
	})

	// Mock DescribeClsLogset so log_set_id validation passes
	patches.ApplyMethodFunc(&localcls.ClsService{}, "DescribeClsLogset", func(_ context.Context, _ string) (*clsv20201016.LogsetInfo, error) {
		return &clsv20201016.LogsetInfo{
			LogsetId:   ptrStringCLT("fake-logset-id"),
			LogsetName: ptrStringCLT("clb_logset"),
		}, nil
	})

	// Mock DescribeClsTopicById for the read-back after create
	topicInfo := &clsv20201016.TopicInfo{
		LogsetId:   ptrStringCLT("fake-logset-id"),
		TopicId:    ptrStringCLT("fake-topic-id"),
		TopicName:  ptrStringCLT("tf-test-topic"),
		Status:     ptrBoolCLT(true),
		Period:     ptrInt64CLT(30),
		CreateTime: ptrStringCLT("2024-01-01 00:00:00"),
	}
	patches.ApplyMethodFunc(&localcls.ClsService{}, "DescribeClsTopicById", func(_ context.Context, _ string, _ *uint64) (*clsv20201016.TopicInfo, error) {
		return topicInfo, nil
	})

	meta := newMockMetaForClbLogTopic()
	res := localclb.ResourceTencentCloudClbLogTopic()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"log_set_id": "fake-logset-id",
		"topic_name": "tf-test-topic",
		"period":     30,
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "fake-topic-id", d.Id())

	// Verify Period was set on the create request as *uint64
	assert.NotNil(t, capturedRequest.Period)
	assert.Equal(t, uint64(30), *capturedRequest.Period)

	// Verify period is read back into state
	assert.Equal(t, 30, d.Get("period"))
}

// TestClbLogTopic_Create_WithoutPeriod verifies that when period is not set,
// the CreateTopic request does not carry Period and the API default applies.
func TestClbLogTopic_Create_WithoutPeriod(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	clbClient := &clbv20180317.Client{}
	patches.ApplyMethodReturn(newMockMetaForClbLogTopic().client, "UseClbClient", clbClient)

	clsClient := &clsv20201016.Client{}
	patches.ApplyMethodReturn(newMockMetaForClbLogTopic().client, "UseClsClient", clsClient)

	var capturedRequest *clbv20180317.CreateTopicRequest
	patches.ApplyMethodFunc(clbClient, "CreateTopic", func(request *clbv20180317.CreateTopicRequest) (*clbv20180317.CreateTopicResponse, error) {
		capturedRequest = request
		resp := clbv20180317.NewCreateTopicResponse()
		resp.Response = &clbv20180317.CreateTopicResponseParams{
			TopicId:   ptrStringCLT("fake-topic-id-2"),
			RequestId: ptrStringCLT("fake-request-id-2"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(&localcls.ClsService{}, "DescribeClsLogset", func(_ context.Context, _ string) (*clsv20201016.LogsetInfo, error) {
		return &clsv20201016.LogsetInfo{
			LogsetId:   ptrStringCLT("fake-logset-id"),
			LogsetName: ptrStringCLT("clb_logset"),
		}, nil
	})

	topicInfo := &clsv20201016.TopicInfo{
		LogsetId:   ptrStringCLT("fake-logset-id"),
		TopicId:    ptrStringCLT("fake-topic-id-2"),
		TopicName:  ptrStringCLT("tf-test-topic-2"),
		Status:     ptrBoolCLT(true),
		Period:     ptrInt64CLT(30),
		CreateTime: ptrStringCLT("2024-01-01 00:00:00"),
	}
	patches.ApplyMethodFunc(&localcls.ClsService{}, "DescribeClsTopicById", func(_ context.Context, _ string, _ *uint64) (*clsv20201016.TopicInfo, error) {
		return topicInfo, nil
	})

	meta := newMockMetaForClbLogTopic()
	res := localclb.ResourceTencentCloudClbLogTopic()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"log_set_id": "fake-logset-id",
		"topic_name": "tf-test-topic-2",
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "fake-topic-id-2", d.Id())

	// Verify Period was NOT set on the create request (user did not specify it)
	assert.Nil(t, capturedRequest.Period)
}

// TestClbLogTopic_Read_PeriodPopulated verifies that period is correctly read
// from the TopicInfo.Period response and set in resource data.
func TestClbLogTopic_Read_PeriodPopulated(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	clsClient := &clsv20201016.Client{}
	patches.ApplyMethodReturn(newMockMetaForClbLogTopic().client, "UseClsClient", clsClient)

	topicInfo := &clsv20201016.TopicInfo{
		LogsetId:   ptrStringCLT("fake-logset-id"),
		TopicId:    ptrStringCLT("fake-topic-id"),
		TopicName:  ptrStringCLT("tf-test-topic"),
		Status:     ptrBoolCLT(true),
		Period:     ptrInt64CLT(60),
		CreateTime: ptrStringCLT("2024-01-01 00:00:00"),
	}
	patches.ApplyMethodFunc(&localcls.ClsService{}, "DescribeClsTopicById", func(_ context.Context, _ string, _ *uint64) (*clsv20201016.TopicInfo, error) {
		return topicInfo, nil
	})

	meta := newMockMetaForClbLogTopic()
	res := localclb.ResourceTencentCloudClbLogTopic()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"log_set_id": "fake-logset-id",
		"topic_name": "tf-test-topic",
	})
	d.SetId("fake-topic-id")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "fake-topic-id", d.Id())

	// Verify period is populated from TopicInfo.Period
	assert.Equal(t, 60, d.Get("period"))
}

// TestClbLogTopic_Read_PeriodNil verifies that when Period is nil in the
// TopicInfo response, period is not set in resource data (no crash).
func TestClbLogTopic_Read_PeriodNil(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	clsClient := &clsv20201016.Client{}
	patches.ApplyMethodReturn(newMockMetaForClbLogTopic().client, "UseClsClient", clsClient)

	topicInfo := &clsv20201016.TopicInfo{
		LogsetId:   ptrStringCLT("fake-logset-id"),
		TopicId:    ptrStringCLT("fake-topic-id-nil"),
		TopicName:  ptrStringCLT("tf-test-topic-nil"),
		Status:     ptrBoolCLT(true),
		Period:     nil,
		CreateTime: ptrStringCLT("2024-01-01 00:00:00"),
	}
	patches.ApplyMethodFunc(&localcls.ClsService{}, "DescribeClsTopicById", func(_ context.Context, _ string, _ *uint64) (*clsv20201016.TopicInfo, error) {
		return topicInfo, nil
	})

	meta := newMockMetaForClbLogTopic()
	res := localclb.ResourceTencentCloudClbLogTopic()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"log_set_id": "fake-logset-id",
		"topic_name": "tf-test-topic-nil",
	})
	d.SetId("fake-topic-id-nil")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "fake-topic-id-nil", d.Id())

	// Verify period defaults to 0 (not set) when TopicInfo.Period is nil
	assert.Equal(t, 0, d.Get("period"))
}

// TestClbLogTopic_Update_Period verifies that updating the period parameter
// calls the ModifyTopic API with the new Period value (cast to *int64).
func TestClbLogTopic_Update_Period(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	clsClient := &clsv20201016.Client{}
	patches.ApplyMethodReturn(newMockMetaForClbLogTopic().client, "UseClsClient", clsClient)

	var capturedRequest *clsv20201016.ModifyTopicRequest
	patches.ApplyMethodFunc(clsClient, "ModifyTopic", func(request *clsv20201016.ModifyTopicRequest) (*clsv20201016.ModifyTopicResponse, error) {
		capturedRequest = request
		resp := clsv20201016.NewModifyTopicResponse()
		resp.Response = &clsv20201016.ModifyTopicResponseParams{
			RequestId: ptrStringCLT("fake-request-id"),
		}
		return resp, nil
	})

	// Mock DescribeClsTopicById for the read-back after update
	topicInfo := &clsv20201016.TopicInfo{
		LogsetId:   ptrStringCLT("fake-logset-id"),
		TopicId:    ptrStringCLT("fake-topic-id"),
		TopicName:  ptrStringCLT("tf-test-topic"),
		Status:     ptrBoolCLT(true),
		Period:     ptrInt64CLT(60),
		CreateTime: ptrStringCLT("2024-01-01 00:00:00"),
	}
	patches.ApplyMethodFunc(&localcls.ClsService{}, "DescribeClsTopicById", func(_ context.Context, _ string, _ *uint64) (*clsv20201016.TopicInfo, error) {
		return topicInfo, nil
	})

	meta := newMockMetaForClbLogTopic()
	res := localclb.ResourceTencentCloudClbLogTopic()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"log_set_id": "fake-logset-id",
		"topic_name": "tf-test-topic",
		"period":     60,
	})
	d.SetId("fake-topic-id")

	// Force period to be detected as changed
	patches.ApplyMethodFunc(d, "HasChange", func(key string) bool {
		return key == "period"
	})

	err := res.Update(d, meta)
	assert.NoError(t, err)

	// Verify ModifyTopic was called with Period as *int64
	assert.NotNil(t, capturedRequest)
	assert.NotNil(t, capturedRequest.Period)
	assert.Equal(t, int64(60), *capturedRequest.Period)
	assert.NotNil(t, capturedRequest.TopicId)
	assert.Equal(t, "fake-topic-id", *capturedRequest.TopicId)

	// Verify the updated value is reflected in state after read-back
	assert.Equal(t, 60, d.Get("period"))
}
