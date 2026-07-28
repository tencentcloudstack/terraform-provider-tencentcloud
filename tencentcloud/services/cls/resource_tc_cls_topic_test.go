package cls_test

import (
	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	localcls "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/cls"

	"context"
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/stretchr/testify/assert"
	clsv20201016 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cls/v20201016"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
)

func init() {
	resource.AddTestSweepers("tencentcloud_cls_topic", &resource.Sweeper{
		Name: "tencentcloud_cls_topic",
		F:    testSweepClsTopic,
	})
}

func testSweepClsTopic(region string) error {
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	sharedClient, err := tcacctest.SharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("getting tencentcloud client error: %s", err.Error())
	}
	client := sharedClient.(tccommon.ProviderMeta).GetAPIV3Conn()

	clsService := localcls.NewClsService(client)

	instances, err := clsService.DescribeClsTopicByFilter(ctx, nil)
	if err != nil {
		return fmt.Errorf("get instance list error: %s", err.Error())
	}

	// add scanning resources
	var resources, nonKeepResources []*tccommon.ResourceInstance
	for _, v := range instances {
		if !tccommon.CheckResourcePersist(*v.TopicName, *v.CreateTime) {
			nonKeepResources = append(nonKeepResources, &tccommon.ResourceInstance{
				Id:   *v.TopicId,
				Name: *v.TopicName,
			})
		}
		resources = append(resources, &tccommon.ResourceInstance{
			Id:         *v.TopicId,
			Name:       *v.TopicName,
			CreateTime: *v.CreateTime,
		})
	}
	tccommon.ProcessScanCloudResources(client, resources, nonKeepResources, "CreateTopic")

	for _, v := range instances {
		instanceId := v.TopicId
		instanceName := v.TopicName

		now := time.Now()

		createTime := tccommon.StringToTime(*v.CreateTime)
		interval := now.Sub(createTime).Minutes()
		if strings.HasPrefix(*instanceName, tcacctest.KeepResource) || strings.HasPrefix(*instanceName, tcacctest.DefaultResource) {
			continue
		}
		// less than 30 minute, not delete
		if tccommon.NeedProtect == 1 && int64(interval) < 30 {
			continue
		}

		if err = clsService.DeleteClsTopic(ctx, *instanceId); err != nil {
			log.Printf("[ERROR] sweep instance %s error: %s", *instanceId, err.Error())
		}
	}
	return nil
}

// go test -i; go test -test.run TestAccTencentCloudClsTopic_basic -v
func TestAccTencentCloudClsTopic_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheck(t) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccClsTopic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClsTopicExists("tencentcloud_cls_topic.example"),
					resource.TestCheckResourceAttr("tencentcloud_cls_topic.example", "topic_name", "tf_example"),
					resource.TestCheckResourceAttr("tencentcloud_cls_topic.example", "storage_type", "hot"),
					resource.TestCheckResourceAttr("tencentcloud_cls_topic.example", "describes", "Test Demo."),
				),
			},
			{
				ResourceName:      "tencentcloud_cls_topic.example",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccClsTopicUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClsTopicExists("tencentcloud_cls_topic.example"),
					resource.TestCheckResourceAttr("tencentcloud_cls_topic.example", "topic_name", "tf_example_update"),
					resource.TestCheckResourceAttr("tencentcloud_cls_topic.example", "storage_type", "hot"),
					resource.TestCheckResourceAttr("tencentcloud_cls_topic.example", "describes", "Test Demo Update."),
				),
			},
		},
	})
}

func testAccCheckClsTopicExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := tccommon.GetLogId(tccommon.ContextNil)
		ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("[CHECK][CLS topic][Exists] check: CLB topic %s is not found", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("[CHECK][CLS topic][Exists] check: CLB topic id is not set")
		}
		clsService := localcls.NewClsService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
		instance, err := clsService.DescribeClsTopicById(ctx, rs.Primary.ID)
		if err != nil {
			return err
		}

		if instance == nil {
			return fmt.Errorf("[CHECK][CLS topic][Exists] id %s is not exist", rs.Primary.ID)
		}
		return nil
	}
}

const testAccClsTopic = `
resource "tencentcloud_cls_logset" "example" {
  logset_name = "tf_example"
  tags        = {
    "demo" = "test"
  }
}

resource "tencentcloud_cls_topic" "example" {
  topic_name           = "tf_example"
  logset_id            = tencentcloud_cls_logset.example.id
  auto_split           = false
  max_split_partitions = 20
  partition_count      = 1
  period               = 30
  storage_type         = "hot"
  describes            = "Test Demo."
  hot_period           = 10
  tags                 = {
    "test" = "test",
  }
}
`

const testAccClsTopicUpdate = `
resource "tencentcloud_cls_logset" "example" {
  logset_name = "tf_example"
  tags        = {
    "demo" = "test"
  }
}

resource "tencentcloud_cls_topic" "example" {
  topic_name           = "tf_example_update"
  logset_id            = tencentcloud_cls_logset.example.id
  auto_split           = false
  max_split_partitions = 20
  partition_count      = 1
  period               = 30
  storage_type         = "hot"
  describes            = "Test Demo Update."
  hot_period           = 10
}
`

// --- gomonkey mock unit tests for biz_type parameter ---

type mockMetaForClsTopic struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMetaForClsTopic) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockMetaForClsTopic{}

func newMockMetaForClsTopic() *mockMetaForClsTopic {
	return &mockMetaForClsTopic{client: &connectivity.TencentCloudClient{}}
}

func ptrStringCT(s string) *string { return &s }
func ptrBoolCT(v bool) *bool       { return &v }
func ptrUint64CT(v uint64) *uint64 { return &v }
func ptrInt64CT(v int64) *int64    { return &v }

// TestClsTopic_Create_WithBizType verifies that when biz_type is set,
// the CreateTopic request carries the correct BizType value.
func TestClsTopic_Create_WithBizType(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	clsClient := &clsv20201016.Client{}
	patches.ApplyMethodReturn(newMockMetaForClsTopic().client, "UseClsClient", clsClient)

	var capturedRequest *clsv20201016.CreateTopicRequest
	patches.ApplyMethodFunc(clsClient, "CreateTopic", func(request *clsv20201016.CreateTopicRequest) (*clsv20201016.CreateTopicResponse, error) {
		capturedRequest = request
		resp := clsv20201016.NewCreateTopicResponse()
		resp.Response = &clsv20201016.CreateTopicResponseParams{
			TopicId:   ptrStringCT("fake-topic-id"),
			RequestId: ptrStringCT("fake-request-id"),
		}
		return resp, nil
	})

	// Mock DescribeClsTopicById to return a topic with BizType=1 for the read-back after create
	topicInfo := &clsv20201016.TopicInfo{
		LogsetId:           ptrStringCT("fake-logset-id"),
		TopicId:            ptrStringCT("fake-topic-id"),
		TopicName:          ptrStringCT("tf-test-topic"),
		PartitionCount:     ptrInt64CT(1),
		AutoSplit:          ptrBoolCT(false),
		MaxSplitPartitions: ptrInt64CT(20),
		StorageType:        ptrStringCT("hot"),
		Period:             ptrInt64CT(30),
		HotPeriod:          ptrUint64CT(10),
		Describes:          ptrStringCT("Test Demo."),
		IsWebTracking:      ptrBoolCT(false),
		BizType:            ptrUint64CT(1),
		Tags: []*clsv20201016.Tag{
			{Key: ptrStringCT("test"), Value: ptrStringCT("test")},
		},
	}
	patches.ApplyMethodFunc(&localcls.ClsService{}, "DescribeClsTopicById", func(_ context.Context, _ string) (*clsv20201016.TopicInfo, error) {
		return topicInfo, nil
	})

	meta := newMockMetaForClsTopic()
	res := localcls.ResourceTencentCloudClsTopic()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"logset_id":            "fake-logset-id",
		"topic_name":           "tf-test-topic",
		"partition_count":      1,
		"auto_split":           false,
		"max_split_partitions": 20,
		"period":               30,
		"storage_type":         "hot",
		"describes":            "Test Demo.",
		"hot_period":           10,
		"biz_type":             1,
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.NotEmpty(t, d.Id())

	// Verify BizType was set on the create request
	assert.NotNil(t, capturedRequest.BizType)
	assert.Equal(t, uint64(1), *capturedRequest.BizType)

	// Verify biz_type is read back into state
	assert.Equal(t, 1, d.Get("biz_type"))
}

// TestClsTopic_Create_WithoutBizType verifies that when biz_type is not set,
// the CreateTopic request does not carry BizType and the default value (0)
// is read back from the API.
func TestClsTopic_Create_WithoutBizType(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	clsClient := &clsv20201016.Client{}
	patches.ApplyMethodReturn(newMockMetaForClsTopic().client, "UseClsClient", clsClient)

	var capturedRequest *clsv20201016.CreateTopicRequest
	patches.ApplyMethodFunc(clsClient, "CreateTopic", func(request *clsv20201016.CreateTopicRequest) (*clsv20201016.CreateTopicResponse, error) {
		capturedRequest = request
		resp := clsv20201016.NewCreateTopicResponse()
		resp.Response = &clsv20201016.CreateTopicResponseParams{
			TopicId:   ptrStringCT("fake-topic-id-2"),
			RequestId: ptrStringCT("fake-request-id-2"),
		}
		return resp, nil
	})

	// Mock DescribeClsTopicById to return a topic with BizType=0 (default log topic)
	topicInfo := &clsv20201016.TopicInfo{
		LogsetId:           ptrStringCT("fake-logset-id"),
		TopicId:            ptrStringCT("fake-topic-id-2"),
		TopicName:          ptrStringCT("tf-test-topic-2"),
		PartitionCount:     ptrInt64CT(1),
		AutoSplit:          ptrBoolCT(false),
		MaxSplitPartitions: ptrInt64CT(20),
		StorageType:        ptrStringCT("hot"),
		Period:             ptrInt64CT(30),
		HotPeriod:          ptrUint64CT(10),
		Describes:          ptrStringCT("Test Demo."),
		IsWebTracking:      ptrBoolCT(false),
		BizType:            ptrUint64CT(0),
	}
	patches.ApplyMethodFunc(&localcls.ClsService{}, "DescribeClsTopicById", func(_ context.Context, _ string) (*clsv20201016.TopicInfo, error) {
		return topicInfo, nil
	})

	meta := newMockMetaForClsTopic()
	res := localcls.ResourceTencentCloudClsTopic()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"logset_id":            "fake-logset-id",
		"topic_name":           "tf-test-topic-2",
		"partition_count":      1,
		"auto_split":           false,
		"max_split_partitions": 20,
		"period":               30,
		"storage_type":         "hot",
		"describes":            "Test Demo.",
		"hot_period":           10,
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.NotEmpty(t, d.Id())

	// Verify BizType was NOT set on the create request (not specified by user)
	assert.Nil(t, capturedRequest.BizType)

	// Verify biz_type is read back as 0 (default log topic)
	assert.Equal(t, 0, d.Get("biz_type"))
}

// TestClsTopic_Read_BizType verifies that biz_type is correctly read from
// the TopicInfo response and set in resource data.
func TestClsTopic_Read_BizType(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	clsClient := &clsv20201016.Client{}
	patches.ApplyMethodReturn(newMockMetaForClsTopic().client, "UseClsClient", clsClient)

	// Mock DescribeClsTopicById to return a topic with BizType=1 (metric topic)
	topicInfo := &clsv20201016.TopicInfo{
		LogsetId:           ptrStringCT("fake-logset-id"),
		TopicId:            ptrStringCT("fake-topic-id-read"),
		TopicName:          ptrStringCT("tf-test-topic-read"),
		PartitionCount:     ptrInt64CT(1),
		AutoSplit:          ptrBoolCT(false),
		MaxSplitPartitions: ptrInt64CT(20),
		StorageType:        ptrStringCT("hot"),
		Period:             ptrInt64CT(30),
		HotPeriod:          ptrUint64CT(10),
		Describes:          ptrStringCT("Test Demo."),
		IsWebTracking:      ptrBoolCT(false),
		BizType:            ptrUint64CT(1),
	}
	patches.ApplyMethodFunc(&localcls.ClsService{}, "DescribeClsTopicById", func(_ context.Context, _ string) (*clsv20201016.TopicInfo, error) {
		return topicInfo, nil
	})

	meta := newMockMetaForClsTopic()
	res := localcls.ResourceTencentCloudClsTopic()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"logset_id":            "fake-logset-id",
		"topic_name":           "tf-test-topic-read",
		"partition_count":      1,
		"auto_split":           false,
		"max_split_partitions": 20,
		"period":               30,
		"storage_type":         "hot",
		"describes":            "Test Demo.",
		"hot_period":           10,
	})
	d.SetId("fake-topic-id-read")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "fake-topic-id-read", d.Id())

	// Verify biz_type is correctly read from the API response
	assert.Equal(t, 1, d.Get("biz_type"))
}

// TestClsTopic_Read_BizTypeNil verifies that when BizType is nil in the
// TopicInfo response, biz_type is not set in resource data (no crash).
func TestClsTopic_Read_BizTypeNil(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	clsClient := &clsv20201016.Client{}
	patches.ApplyMethodReturn(newMockMetaForClsTopic().client, "UseClsClient", clsClient)

	// Mock DescribeClsTopicById to return a topic with BizType=nil
	topicInfo := &clsv20201016.TopicInfo{
		LogsetId:           ptrStringCT("fake-logset-id"),
		TopicId:            ptrStringCT("fake-topic-id-nil"),
		TopicName:          ptrStringCT("tf-test-topic-nil"),
		PartitionCount:     ptrInt64CT(1),
		AutoSplit:          ptrBoolCT(false),
		MaxSplitPartitions: ptrInt64CT(20),
		StorageType:        ptrStringCT("hot"),
		Period:             ptrInt64CT(30),
		HotPeriod:          ptrUint64CT(10),
		Describes:          ptrStringCT("Test Demo."),
		IsWebTracking:      ptrBoolCT(false),
		BizType:            nil,
	}
	patches.ApplyMethodFunc(&localcls.ClsService{}, "DescribeClsTopicById", func(_ context.Context, _ string) (*clsv20201016.TopicInfo, error) {
		return topicInfo, nil
	})

	meta := newMockMetaForClsTopic()
	res := localcls.ResourceTencentCloudClsTopic()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"logset_id":            "fake-logset-id",
		"topic_name":           "tf-test-topic-nil",
		"partition_count":      1,
		"auto_split":           false,
		"max_split_partitions": 20,
		"period":               30,
		"storage_type":         "hot",
		"describes":            "Test Demo.",
		"hot_period":           10,
	})
	d.SetId("fake-topic-id-nil")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "fake-topic-id-nil", d.Id())
}

// TestClsTopic_Update_BizTypeImmutable verifies that attempting to update
// biz_type returns an error since it is in the immutableArgs array.
func TestClsTopic_Update_BizTypeImmutable(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	clsClient := &clsv20201016.Client{}
	patches.ApplyMethodReturn(newMockMetaForClsTopic().client, "UseClsClient", clsClient)

	meta := newMockMetaForClsTopic()
	res := localcls.ResourceTencentCloudClsTopic()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"logset_id":            "fake-logset-id",
		"topic_name":           "tf-test-topic-immutable",
		"partition_count":      1,
		"auto_split":           false,
		"max_split_partitions": 20,
		"period":               30,
		"storage_type":         "hot",
		"describes":            "Test Demo.",
		"hot_period":           10,
		"biz_type":             1,
	})
	d.SetId("fake-topic-id-immutable")

	// Force biz_type to be detected as changed
	patches.ApplyMethodFunc(d, "HasChange", func(key string) bool {
		return key == "biz_type"
	})

	err := res.Update(d, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "biz_type")
}
