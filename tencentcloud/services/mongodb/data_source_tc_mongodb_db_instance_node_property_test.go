package mongodb_test

import (
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	mongodbv20190725 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/mongodb/v20190725"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/mongodb"
)

type mockMetaForNodeProperty struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMetaForNodeProperty) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockMetaForNodeProperty{}

func newMockMetaForNodeProperty() *mockMetaForNodeProperty {
	return &mockMetaForNodeProperty{client: &connectivity.TencentCloudClient{}}
}

func ptrStrNodeProp(s string) *string {
	return &s
}

func ptrInt64NodeProp(i int64) *int64 {
	return &i
}

func ptrBoolNodeProp(b bool) *bool {
	return &b
}

// go test ./tencentcloud/services/mongodb/ -run "TestMongodbDbInstanceNodeProperty" -v -count=1 -gcflags="all=-l"

// TestMongodbDbInstanceNodeProperty_Read_Success tests successful read with mongos and replicate_sets
func TestMongodbDbInstanceNodeProperty_Read_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	mongodbClient := &mongodbv20190725.Client{}
	patches.ApplyMethodReturn(newMockMetaForNodeProperty().client, "UseMongodbClient", mongodbClient)

	patches.ApplyMethodFunc(mongodbClient, "DescribeDBInstanceNodeProperty", func(request *mongodbv20190725.DescribeDBInstanceNodePropertyRequest) (*mongodbv20190725.DescribeDBInstanceNodePropertyResponse, error) {
		resp := mongodbv20190725.NewDescribeDBInstanceNodePropertyResponse()
		resp.Response = &mongodbv20190725.DescribeDBInstanceNodePropertyResponseParams{
			Mongos: []*mongodbv20190725.NodeProperty{
				{
					Zone:              ptrStrNodeProp("ap-guangzhou-3"),
					NodeName:          ptrStrNodeProp("cmgo-test1234_0"),
					Address:           ptrStrNodeProp("10.0.0.1:27017"),
					WanServiceAddress: ptrStrNodeProp(""),
					Role:              ptrStrNodeProp("PRIMARY"),
					Hidden:            ptrBoolNodeProp(false),
					Status:            ptrStrNodeProp("NORMAL"),
					SlaveDelay:        ptrInt64NodeProp(0),
					Priority:          ptrInt64NodeProp(1),
					Votes:             ptrInt64NodeProp(1),
					Tags:              []*mongodbv20190725.NodeTag{},
					ReplicateSetId:    ptrStrNodeProp("rs0"),
				},
			},
			ReplicateSets: []*mongodbv20190725.ReplicateSetInfo{
				{
					Nodes: []*mongodbv20190725.NodeProperty{
						{
							Zone:           ptrStrNodeProp("ap-guangzhou-3"),
							NodeName:       ptrStrNodeProp("cmgo-test1234_1"),
							Address:        ptrStrNodeProp("10.0.0.2:27017"),
							Role:           ptrStrNodeProp("SECONDARY"),
							Hidden:         ptrBoolNodeProp(false),
							Status:         ptrStrNodeProp("NORMAL"),
							SlaveDelay:     ptrInt64NodeProp(0),
							Priority:       ptrInt64NodeProp(1),
							Votes:          ptrInt64NodeProp(1),
							ReplicateSetId: ptrStrNodeProp("rs0"),
						},
					},
				},
			},
			RequestId: ptrStrNodeProp("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaForNodeProperty()
	res := mongodb.DataSourceTencentCloudMongodbDbInstanceNodeProperty()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"instance_id": "cmgo-test1234",
	})

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "cmgo-test1234", d.Id())

	mongos := d.Get("mongos").([]interface{})
	assert.Equal(t, 1, len(mongos))
	mongosNode := mongos[0].(map[string]interface{})
	assert.Equal(t, "ap-guangzhou-3", mongosNode["zone"])
	assert.Equal(t, "PRIMARY", mongosNode["role"])

	replicateSets := d.Get("replicate_sets").([]interface{})
	assert.Equal(t, 1, len(replicateSets))
	rsItem := replicateSets[0].(map[string]interface{})
	nodes := rsItem["nodes"].([]interface{})
	assert.Equal(t, 1, len(nodes))
	node := nodes[0].(map[string]interface{})
	assert.Equal(t, "SECONDARY", node["role"])
}

// TestMongodbDbInstanceNodeProperty_Read_EmptyMongos tests read when mongos is nil
func TestMongodbDbInstanceNodeProperty_Read_EmptyMongos(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	mongodbClient := &mongodbv20190725.Client{}
	patches.ApplyMethodReturn(newMockMetaForNodeProperty().client, "UseMongodbClient", mongodbClient)

	patches.ApplyMethodFunc(mongodbClient, "DescribeDBInstanceNodeProperty", func(request *mongodbv20190725.DescribeDBInstanceNodePropertyRequest) (*mongodbv20190725.DescribeDBInstanceNodePropertyResponse, error) {
		resp := mongodbv20190725.NewDescribeDBInstanceNodePropertyResponse()
		resp.Response = &mongodbv20190725.DescribeDBInstanceNodePropertyResponseParams{
			Mongos: nil,
			ReplicateSets: []*mongodbv20190725.ReplicateSetInfo{
				{
					Nodes: []*mongodbv20190725.NodeProperty{
						{
							Zone:   ptrStrNodeProp("ap-guangzhou-3"),
							Role:   ptrStrNodeProp("PRIMARY"),
							Status: ptrStrNodeProp("NORMAL"),
						},
					},
				},
			},
			RequestId: ptrStrNodeProp("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaForNodeProperty()
	res := mongodb.DataSourceTencentCloudMongodbDbInstanceNodeProperty()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"instance_id": "cmgo-test1234",
	})

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "cmgo-test1234", d.Id())

	mongos := d.Get("mongos").([]interface{})
	assert.Equal(t, 0, len(mongos))

	replicateSets := d.Get("replicate_sets").([]interface{})
	assert.Equal(t, 1, len(replicateSets))
}

// TestMongodbDbInstanceNodeProperty_Read_APIError tests read when API returns error
func TestMongodbDbInstanceNodeProperty_Read_APIError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	mongodbClient := &mongodbv20190725.Client{}
	patches.ApplyMethodReturn(newMockMetaForNodeProperty().client, "UseMongodbClient", mongodbClient)

	patches.ApplyMethodFunc(mongodbClient, "DescribeDBInstanceNodeProperty", func(request *mongodbv20190725.DescribeDBInstanceNodePropertyRequest) (*mongodbv20190725.DescribeDBInstanceNodePropertyResponse, error) {
		return nil, assert.AnError
	})

	meta := newMockMetaForNodeProperty()
	res := mongodb.DataSourceTencentCloudMongodbDbInstanceNodeProperty()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"instance_id": "cmgo-test1234",
	})

	err := res.Read(d, meta)
	assert.Error(t, err)
}

// TestMongodbDbInstanceNodeProperty_Read_WithTags tests read with node tags
func TestMongodbDbInstanceNodeProperty_Read_WithTags(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	mongodbClient := &mongodbv20190725.Client{}
	patches.ApplyMethodReturn(newMockMetaForNodeProperty().client, "UseMongodbClient", mongodbClient)

	patches.ApplyMethodFunc(mongodbClient, "DescribeDBInstanceNodeProperty", func(request *mongodbv20190725.DescribeDBInstanceNodePropertyRequest) (*mongodbv20190725.DescribeDBInstanceNodePropertyResponse, error) {
		resp := mongodbv20190725.NewDescribeDBInstanceNodePropertyResponse()
		resp.Response = &mongodbv20190725.DescribeDBInstanceNodePropertyResponseParams{
			Mongos: []*mongodbv20190725.NodeProperty{
				{
					Zone:   ptrStrNodeProp("ap-guangzhou-3"),
					Role:   ptrStrNodeProp("PRIMARY"),
					Status: ptrStrNodeProp("NORMAL"),
					Tags: []*mongodbv20190725.NodeTag{
						{
							TagKey:   ptrStrNodeProp("env"),
							TagValue: ptrStrNodeProp("prod"),
						},
					},
				},
			},
			ReplicateSets: nil,
			RequestId:     ptrStrNodeProp("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaForNodeProperty()
	res := mongodb.DataSourceTencentCloudMongodbDbInstanceNodeProperty()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"instance_id": "cmgo-test1234",
	})

	err := res.Read(d, meta)
	assert.NoError(t, err)

	mongos := d.Get("mongos").([]interface{})
	assert.Equal(t, 1, len(mongos))
	mongosNode := mongos[0].(map[string]interface{})
	tags := mongosNode["tags"].([]interface{})
	assert.Equal(t, 1, len(tags))
	tag := tags[0].(map[string]interface{})
	assert.Equal(t, "env", tag["tag_key"])
	assert.Equal(t, "prod", tag["tag_value"])
}

// TestMongodbDbInstanceNodeProperty_Schema tests the schema definition
func TestMongodbDbInstanceNodeProperty_Schema(t *testing.T) {
	res := mongodb.DataSourceTencentCloudMongodbDbInstanceNodeProperty()

	assert.NotNil(t, res)
	assert.NotNil(t, res.Read)

	// Check required fields
	assert.True(t, res.Schema["instance_id"].Required)

	// Check optional fields
	assert.True(t, res.Schema["node_ids"].Optional)
	assert.True(t, res.Schema["roles"].Optional)
	assert.True(t, res.Schema["only_hidden"].Optional)
	assert.True(t, res.Schema["priority"].Optional)
	assert.True(t, res.Schema["votes"].Optional)
	assert.True(t, res.Schema["tags"].Optional)

	// Check computed fields
	assert.True(t, res.Schema["mongos"].Computed)
	assert.True(t, res.Schema["replicate_sets"].Computed)
}
