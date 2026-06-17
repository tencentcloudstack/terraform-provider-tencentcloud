package tpulsar_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/tpulsar"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	tdmq "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tdmq/v20200217"
)

// go test -i; go test -test.run TestAccTencentCloudTdmqTopicResource_basic -v
func TestAccTencentCloudTdmqTopicResource_basic(t *testing.T) {
	terraformId := "tencentcloud_tdmq_topic.example"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheck(t) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTdmqTopic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(terraformId, "environ_id"),
					resource.TestCheckResourceAttrSet(terraformId, "cluster_id"),
					resource.TestCheckResourceAttr(terraformId, "topic_name", "tf-example-topic"),
					resource.TestCheckResourceAttr(terraformId, "partitions", "6"),
					resource.TestCheckResourceAttr(terraformId, "pulsar_topic_type", "3"),
					resource.TestCheckResourceAttr(terraformId, "remark", "remark."),
				),
			},
			{
				Config: testAccTdmqTopicUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(terraformId, "environ_id"),
					resource.TestCheckResourceAttrSet(terraformId, "cluster_id"),
					resource.TestCheckResourceAttr(terraformId, "topic_name", "tf-example-topic"),
					resource.TestCheckResourceAttr(terraformId, "partitions", "8"),
					resource.TestCheckResourceAttr(terraformId, "pulsar_topic_type", "3"),
					resource.TestCheckResourceAttr(terraformId, "remark", "remark update."),
				),
			},
		},
	})
}

const testAccTdmqTopic = `
resource "tencentcloud_tdmq_instance" "example" {
  cluster_name = "tf_example"
  remark       = "remark."
  tags         = {
    "createdBy" = "terraform"
  }
}

resource "tencentcloud_tdmq_namespace" "example" {
  environ_name = "tf_example"
  msg_ttl      = 300
  cluster_id   = tencentcloud_tdmq_instance.example.id
  retention_policy {
    time_in_minutes = 60
    size_in_mb      = 10
  }
  remark = "remark."
}

resource "tencentcloud_tdmq_topic" "example" {
  environ_id        = tencentcloud_tdmq_namespace.example.environ_name
  cluster_id        = tencentcloud_tdmq_instance.example.id
  topic_name        = "tf-example-topic"
  partitions        = 6
  pulsar_topic_type = 3
  remark            = "remark."
}
`

const testAccTdmqTopicUpdate = `
resource "tencentcloud_tdmq_instance" "example" {
  cluster_name = "tf_example"
  remark       = "remark."
  tags         = {
    "createdBy" = "terraform"
  }
}

resource "tencentcloud_tdmq_namespace" "example" {
  environ_name = "tf_example"
  msg_ttl      = 300
  cluster_id   = tencentcloud_tdmq_instance.example.id
  retention_policy {
    time_in_minutes = 60
    size_in_mb      = 10
  }
  remark = "remark."
}

resource "tencentcloud_tdmq_topic" "example" {
  environ_id        = tencentcloud_tdmq_namespace.example.environ_name
  cluster_id        = tencentcloud_tdmq_instance.example.id
  topic_name        = "tf-example-topic"
  partitions        = 8
  pulsar_topic_type = 3
  remark            = "remark update."
}
`

// mockMetaTdmqTopic implements tccommon.ProviderMeta
type mockMetaTdmqTopic struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMetaTdmqTopic) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockMetaTdmqTopic{}

func ptrTdmqTopicString(s string) *string {
	return &s
}

func ptrTdmqTopicInt64(i int64) *int64 {
	return &i
}

func ptrTdmqTopicUint64(i uint64) *uint64 {
	return &i
}

// go test ./tencentcloud/services/tpulsar/ -run "TestTdmqTopicTags" -v -count=1 -gcflags="all=-l"

// TestTdmqTopicTags_CreateWithTags tests that tags are passed to CreateTdmqTopic when specified
func TestTdmqTopicTags_CreateWithTags(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	mockClient := &connectivity.TencentCloudClient{}
	tdmqClient := &tdmq.Client{}
	patches.ApplyMethodReturn(mockClient, "UseTdmqClient", tdmqClient)

	// Capture the tags passed to CreateTopic
	var capturedTags []*tdmq.Tag
	patches.ApplyMethodFunc(tdmqClient, "CreateTopic", func(request *tdmq.CreateTopicRequest) (*tdmq.CreateTopicResponse, error) {
		capturedTags = request.Tags
		resp := tdmq.NewCreateTopicResponse()
		resp.Response = &tdmq.CreateTopicResponseParams{
			EnvironmentId: ptrTdmqTopicString("test-env"),
			TopicName:     ptrTdmqTopicString("test-topic"),
			Partitions:    ptrTdmqTopicUint64(6),
			RequestId:     ptrTdmqTopicString("fake-request-id"),
		}
		return resp, nil
	})

	// Patch DescribeTopics for the read after create
	patches.ApplyMethodFunc(tdmqClient, "DescribeTopics", func(request *tdmq.DescribeTopicsRequest) (*tdmq.DescribeTopicsResponse, error) {
		resp := tdmq.NewDescribeTopicsResponse()
		resp.Response = &tdmq.DescribeTopicsResponseParams{
			TopicSets: []*tdmq.Topic{
				{
					TopicName:       ptrTdmqTopicString("test-topic"),
					Partitions:      ptrTdmqTopicInt64(6),
					TopicType:       ptrTdmqTopicUint64(0),
					PulsarTopicType: ptrTdmqTopicInt64(3),
					Remark:          ptrTdmqTopicString("test remark"),
					CreateTime:      ptrTdmqTopicString("2024-01-01 00:00:00"),
					Tags: []*tdmq.Tag{
						{
							TagKey:   ptrTdmqTopicString("env"),
							TagValue: ptrTdmqTopicString("test"),
						},
					},
				},
			},
			TotalCount: ptrTdmqTopicUint64(1),
			RequestId:  ptrTdmqTopicString("fake-request-id"),
		}
		return resp, nil
	})

	meta := &mockMetaTdmqTopic{client: mockClient}
	res := tpulsar.ResourceTencentCloudTdmqTopic()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"environ_id":        "test-env",
		"topic_name":        "test-topic",
		"partitions":        6,
		"cluster_id":        "test-cluster",
		"pulsar_topic_type": 3,
		"remark":            "test remark",
		"tags": map[string]interface{}{
			"env": "test",
		},
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "test-topic", d.Id())

	// Verify tags were passed to CreateTopic
	assert.NotNil(t, capturedTags)
	assert.Len(t, capturedTags, 1)
	assert.Equal(t, "env", *capturedTags[0].TagKey)
	assert.Equal(t, "test", *capturedTags[0].TagValue)
}

// TestTdmqTopicTags_CreateWithoutTags tests that tags are not passed when not specified
func TestTdmqTopicTags_CreateWithoutTags(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	mockClient := &connectivity.TencentCloudClient{}
	tdmqClient := &tdmq.Client{}
	patches.ApplyMethodReturn(mockClient, "UseTdmqClient", tdmqClient)

	var capturedTags []*tdmq.Tag
	patches.ApplyMethodFunc(tdmqClient, "CreateTopic", func(request *tdmq.CreateTopicRequest) (*tdmq.CreateTopicResponse, error) {
		capturedTags = request.Tags
		resp := tdmq.NewCreateTopicResponse()
		resp.Response = &tdmq.CreateTopicResponseParams{
			EnvironmentId: ptrTdmqTopicString("test-env"),
			TopicName:     ptrTdmqTopicString("test-topic"),
			Partitions:    ptrTdmqTopicUint64(6),
			RequestId:     ptrTdmqTopicString("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(tdmqClient, "DescribeTopics", func(request *tdmq.DescribeTopicsRequest) (*tdmq.DescribeTopicsResponse, error) {
		resp := tdmq.NewDescribeTopicsResponse()
		resp.Response = &tdmq.DescribeTopicsResponseParams{
			TopicSets: []*tdmq.Topic{
				{
					TopicName:       ptrTdmqTopicString("test-topic"),
					Partitions:      ptrTdmqTopicInt64(6),
					TopicType:       ptrTdmqTopicUint64(0),
					PulsarTopicType: ptrTdmqTopicInt64(3),
					Remark:          ptrTdmqTopicString("test remark"),
					CreateTime:      ptrTdmqTopicString("2024-01-01 00:00:00"),
				},
			},
			TotalCount: ptrTdmqTopicUint64(1),
			RequestId:  ptrTdmqTopicString("fake-request-id"),
		}
		return resp, nil
	})

	meta := &mockMetaTdmqTopic{client: mockClient}
	res := tpulsar.ResourceTencentCloudTdmqTopic()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"environ_id":        "test-env",
		"topic_name":        "test-topic",
		"partitions":        6,
		"cluster_id":        "test-cluster",
		"pulsar_topic_type": 3,
		"remark":            "test remark",
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "test-topic", d.Id())

	// Verify tags were NOT passed to CreateTopic
	assert.Nil(t, capturedTags)
}

// TestTdmqTopicTags_ReadWithTags tests that tags are correctly read from DescribeTopics response
func TestTdmqTopicTags_ReadWithTags(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	mockClient := &connectivity.TencentCloudClient{}
	tdmqClient := &tdmq.Client{}
	patches.ApplyMethodReturn(mockClient, "UseTdmqClient", tdmqClient)

	patches.ApplyMethodFunc(tdmqClient, "DescribeTopics", func(request *tdmq.DescribeTopicsRequest) (*tdmq.DescribeTopicsResponse, error) {
		resp := tdmq.NewDescribeTopicsResponse()
		resp.Response = &tdmq.DescribeTopicsResponseParams{
			TopicSets: []*tdmq.Topic{
				{
					TopicName:       ptrTdmqTopicString("test-topic"),
					Partitions:      ptrTdmqTopicInt64(6),
					TopicType:       ptrTdmqTopicUint64(0),
					PulsarTopicType: ptrTdmqTopicInt64(3),
					Remark:          ptrTdmqTopicString("test remark"),
					CreateTime:      ptrTdmqTopicString("2024-01-01 00:00:00"),
					Tags: []*tdmq.Tag{
						{
							TagKey:   ptrTdmqTopicString("env"),
							TagValue: ptrTdmqTopicString("production"),
						},
						{
							TagKey:   ptrTdmqTopicString("team"),
							TagValue: ptrTdmqTopicString("backend"),
						},
					},
				},
			},
			TotalCount: ptrTdmqTopicUint64(1),
			RequestId:  ptrTdmqTopicString("fake-request-id"),
		}
		return resp, nil
	})

	meta := &mockMetaTdmqTopic{client: mockClient}
	res := tpulsar.ResourceTencentCloudTdmqTopic()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"environ_id":        "test-env",
		"topic_name":        "test-topic",
		"partitions":        6,
		"cluster_id":        "test-cluster",
		"pulsar_topic_type": 3,
		"remark":            "test remark",
	})
	d.SetId("test-topic")

	err := res.Read(d, meta)
	assert.NoError(t, err)

	tags := d.Get("tags").(map[string]interface{})
	assert.Len(t, tags, 2)
	assert.Equal(t, "production", tags["env"])
	assert.Equal(t, "backend", tags["team"])
}

// TestTdmqTopicTags_Schema tests the schema definition of tags
func TestTdmqTopicTags_Schema(t *testing.T) {
	res := tpulsar.ResourceTencentCloudTdmqTopic()

	assert.NotNil(t, res)
	assert.Contains(t, res.Schema, "tags")

	tagsSchema := res.Schema["tags"]
	assert.Equal(t, schema.TypeMap, tagsSchema.Type)
	assert.True(t, tagsSchema.Optional)
	assert.True(t, tagsSchema.ForceNew)
}

// TestTdmqTopicTags_ReadWithNilTags tests that tags are not set when response Tags field is nil
func TestTdmqTopicTags_ReadWithNilTags(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	mockClient := &connectivity.TencentCloudClient{}
	tdmqClient := &tdmq.Client{}
	patches.ApplyMethodReturn(mockClient, "UseTdmqClient", tdmqClient)

	patches.ApplyMethodFunc(tdmqClient, "DescribeTopics", func(request *tdmq.DescribeTopicsRequest) (*tdmq.DescribeTopicsResponse, error) {
		resp := tdmq.NewDescribeTopicsResponse()
		resp.Response = &tdmq.DescribeTopicsResponseParams{
			TopicSets: []*tdmq.Topic{
				{
					TopicName:       ptrTdmqTopicString("test-topic"),
					Partitions:      ptrTdmqTopicInt64(6),
					TopicType:       ptrTdmqTopicUint64(0),
					PulsarTopicType: ptrTdmqTopicInt64(3),
					Remark:          ptrTdmqTopicString("test remark"),
					CreateTime:      ptrTdmqTopicString("2024-01-01 00:00:00"),
					Tags:            nil,
				},
			},
			TotalCount: ptrTdmqTopicUint64(1),
			RequestId:  ptrTdmqTopicString("fake-request-id"),
		}
		return resp, nil
	})

	meta := &mockMetaTdmqTopic{client: mockClient}
	res := tpulsar.ResourceTencentCloudTdmqTopic()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"environ_id":        "test-env",
		"topic_name":        "test-topic",
		"partitions":        6,
		"cluster_id":        "test-cluster",
		"pulsar_topic_type": 3,
		"remark":            "test remark",
	})
	d.SetId("test-topic")

	err := res.Read(d, meta)
	assert.NoError(t, err)

	tags := d.Get("tags").(map[string]interface{})
	assert.Len(t, tags, 0)
}
