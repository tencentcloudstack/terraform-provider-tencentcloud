package cls_test

import (
	"strconv"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	clsv20201016 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cls/v20201016"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/cls"
)

// itoa wraps strconv.Itoa for readable pagination id construction.
func itoa(i int) string { return strconv.Itoa(i) }

type mockMetaForClsMachineGroupsDS struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMetaForClsMachineGroupsDS) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockMetaForClsMachineGroupsDS{}

func newMockMetaForClsMachineGroupsDS() *mockMetaForClsMachineGroupsDS {
	return &mockMetaForClsMachineGroupsDS{client: &connectivity.TencentCloudClient{}}
}

func ptrStringMGDS(s string) *string { return &s }
func ptrInt64MGDS(v int64) *int64    { return &v }
func ptrUint64MGDS(v uint64) *uint64 { return &v }
func ptrBoolMGDS(v bool) *bool       { return &v }
func ptrStringSliceMGDS(ss []string) []*string {
	out := make([]*string, len(ss))
	for i := range ss {
		out[i] = &ss[i]
	}
	return out
}

// buildFullMachineGroupInfo returns a MachineGroupInfo with every field populated,
// including nested machine_group_type, tags, and meta_tags.
func buildFullMachineGroupInfo() *clsv20201016.MachineGroupInfo {
	return &clsv20201016.MachineGroupInfo{
		GroupId:   ptrStringMGDS("mg-aaaa1111"),
		GroupName: ptrStringMGDS("tf-machine-group-1"),
		MachineGroupType: &clsv20201016.MachineGroupTypeInfo{
			Type:   ptrStringMGDS("ip"),
			Values: ptrStringSliceMGDS([]string{"192.168.1.1", "192.168.1.2"}),
		},
		CreateTime: ptrStringMGDS("2024-01-01 10:00:00"),
		Tags: []*clsv20201016.Tag{
			{Key: ptrStringMGDS("env"), Value: ptrStringMGDS("prod")},
			{Key: ptrStringMGDS("team"), Value: ptrStringMGDS("infra")},
		},
		AutoUpdate:       ptrStringMGDS("true"),
		UpdateStartTime:  ptrStringMGDS("17:05:00"),
		UpdateEndTime:    ptrStringMGDS("19:05:00"),
		ServiceLogging:   ptrBoolMGDS(true),
		DelayCleanupTime: ptrInt64MGDS(30),
		MetaTags: []*clsv20201016.MetaTagInfo{
			{Key: ptrStringMGDS("k1"), Value: ptrStringMGDS("v1")},
			{Key: ptrStringMGDS("k2"), Value: ptrStringMGDS("v2")},
		},
		OSType: ptrUint64MGDS(0),
	}
}

func buildMachineGroupInfoWithNilNested() *clsv20201016.MachineGroupInfo {
	return &clsv20201016.MachineGroupInfo{
		GroupId:          ptrStringMGDS("mg-bbbb2222"),
		GroupName:        ptrStringMGDS("tf-machine-group-2"),
		MachineGroupType: nil,
		CreateTime:       ptrStringMGDS("2024-02-02 11:00:00"),
		Tags:             nil,
		AutoUpdate:       ptrStringMGDS("false"),
		UpdateStartTime:  nil,
		UpdateEndTime:    nil,
		ServiceLogging:   ptrBoolMGDS(false),
		DelayCleanupTime: ptrInt64MGDS(15),
		MetaTags:         nil,
		OSType:           ptrUint64MGDS(1),
	}
}

func TestClsMachineGroupsDS_ReadBasic(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	clsClient := &clsv20201016.Client{}
	patches.ApplyMethodReturn(newMockMetaForClsMachineGroupsDS().client, "UseClsClient", clsClient)

	fullInfo := buildFullMachineGroupInfo()
	patches.ApplyMethodFunc(clsClient, "DescribeMachineGroups", func(request *clsv20201016.DescribeMachineGroupsRequest) (*clsv20201016.DescribeMachineGroupsResponse, error) {
		resp := clsv20201016.NewDescribeMachineGroupsResponse()
		resp.Response = &clsv20201016.DescribeMachineGroupsResponseParams{
			MachineGroups: []*clsv20201016.MachineGroupInfo{fullInfo},
			TotalCount:    ptrInt64MGDS(1),
			RequestId:     ptrStringMGDS("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaForClsMachineGroupsDS()
	res := cls.DataSourceTencentCloudClsMachineGroups()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.NotEmpty(t, d.Id())

	machineGroups := d.Get("machine_groups").([]interface{})
	assert.Len(t, machineGroups, 1)

	group0 := machineGroups[0].(map[string]interface{})
	assert.Equal(t, "mg-aaaa1111", group0["group_id"].(string))
	assert.Equal(t, "tf-machine-group-1", group0["group_name"].(string))
	assert.Equal(t, "2024-01-01 10:00:00", group0["create_time"].(string))
	assert.Equal(t, "true", group0["auto_update"].(string))
	assert.Equal(t, "17:05:00", group0["update_start_time"].(string))
	assert.Equal(t, "19:05:00", group0["update_end_time"].(string))
	assert.Equal(t, true, group0["service_logging"].(bool))
	assert.Equal(t, 30, group0["delay_cleanup_time"].(int))
	assert.Equal(t, 0, group0["os_type"].(int))

	// machine_group_type (list of one)
	mgtList := group0["machine_group_type"].([]interface{})
	assert.Len(t, mgtList, 1)
	mgt := mgtList[0].(map[string]interface{})
	assert.Equal(t, "ip", mgt["type"].(string))
	mgtValues := mgt["values"].(*schema.Set).List()
	assert.Len(t, mgtValues, 2)

	// tags
	tagsList := group0["tags"].([]interface{})
	assert.Len(t, tagsList, 2)
	tag0 := tagsList[0].(map[string]interface{})
	assert.Equal(t, "env", tag0["key"].(string))
	assert.Equal(t, "prod", tag0["value"].(string))

	// meta_tags
	metaTagsList := group0["meta_tags"].([]interface{})
	assert.Len(t, metaTagsList, 2)
	metaTag0 := metaTagsList[0].(map[string]interface{})
	assert.Equal(t, "k1", metaTag0["key"].(string))
	assert.Equal(t, "v1", metaTag0["value"].(string))

	totalCount := d.Get("total_count").(int)
	assert.Equal(t, 1, totalCount)
}

func TestClsMachineGroupsDS_ReadWithNilNested(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	clsClient := &clsv20201016.Client{}
	patches.ApplyMethodReturn(newMockMetaForClsMachineGroupsDS().client, "UseClsClient", clsClient)

	nilNestedInfo := buildMachineGroupInfoWithNilNested()
	patches.ApplyMethodFunc(clsClient, "DescribeMachineGroups", func(request *clsv20201016.DescribeMachineGroupsRequest) (*clsv20201016.DescribeMachineGroupsResponse, error) {
		resp := clsv20201016.NewDescribeMachineGroupsResponse()
		resp.Response = &clsv20201016.DescribeMachineGroupsResponseParams{
			MachineGroups: []*clsv20201016.MachineGroupInfo{nilNestedInfo},
			TotalCount:    ptrInt64MGDS(1),
			RequestId:     ptrStringMGDS("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaForClsMachineGroupsDS()
	res := cls.DataSourceTencentCloudClsMachineGroups()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.NotEmpty(t, d.Id())

	machineGroups := d.Get("machine_groups").([]interface{})
	assert.Len(t, machineGroups, 1)

	group0 := machineGroups[0].(map[string]interface{})
	assert.Equal(t, "mg-bbbb2222", group0["group_id"].(string))
	assert.Equal(t, "tf-machine-group-2", group0["group_name"].(string))
	assert.Equal(t, "2024-02-02 11:00:00", group0["create_time"].(string))
	assert.Equal(t, "false", group0["auto_update"].(string))
	assert.Equal(t, false, group0["service_logging"].(bool))
	assert.Equal(t, 15, group0["delay_cleanup_time"].(int))
	assert.Equal(t, 1, group0["os_type"].(int))

	// nil nested fields should not panic; machine_group_type should be absent/empty
	mgtList := group0["machine_group_type"].([]interface{})
	assert.Empty(t, mgtList)

	tagsList := group0["tags"].([]interface{})
	assert.Empty(t, tagsList)

	metaTagsList := group0["meta_tags"].([]interface{})
	assert.Empty(t, metaTagsList)

	totalCount := d.Get("total_count").(int)
	assert.Equal(t, 1, totalCount)
}

func TestClsMachineGroupsDS_ReadWithPagination(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	clsClient := &clsv20201016.Client{}
	patches.ApplyMethodReturn(newMockMetaForClsMachineGroupsDS().client, "UseClsClient", clsClient)

	// Build 3 pages of results: page1 (100 items), page2 (100 items), page3 (3 items, < limit -> stop).
	// Use offset to distinguish the page returned by the mock.
	limit := int64(100)
	page3Count := 3

	patches.ApplyMethodFunc(clsClient, "DescribeMachineGroups", func(request *clsv20201016.DescribeMachineGroupsRequest) (*clsv20201016.DescribeMachineGroupsResponse, error) {
		offset := 0
		if request.Offset != nil {
			offset = int(*request.Offset)
		}

		resp := clsv20201016.NewDescribeMachineGroupsResponse()
		var groups []*clsv20201016.MachineGroupInfo
		switch offset {
		case 0:
			// page 1: full page of 100
			groups = make([]*clsv20201016.MachineGroupInfo, 0, int(limit))
			for i := 0; i < int(limit); i++ {
				groups = append(groups, &clsv20201016.MachineGroupInfo{
					GroupId:   ptrStringMGDS("mg-page1-" + itoa(i)),
					GroupName: ptrStringMGDS("tf-page1-" + itoa(i)),
				})
			}
		case int(limit):
			// page 2: full page of 100
			groups = make([]*clsv20201016.MachineGroupInfo, 0, int(limit))
			for i := 0; i < int(limit); i++ {
				groups = append(groups, &clsv20201016.MachineGroupInfo{
					GroupId:   ptrStringMGDS("mg-page2-" + itoa(i)),
					GroupName: ptrStringMGDS("tf-page2-" + itoa(i)),
				})
			}
		case 2 * int(limit):
			// page 3: partial page of 3 (< limit -> stop)
			groups = make([]*clsv20201016.MachineGroupInfo, 0, page3Count)
			for i := 0; i < page3Count; i++ {
				groups = append(groups, &clsv20201016.MachineGroupInfo{
					GroupId:   ptrStringMGDS("mg-page3-" + itoa(i)),
					GroupName: ptrStringMGDS("tf-page3-" + itoa(i)),
				})
			}
		default:
			groups = []*clsv20201016.MachineGroupInfo{}
		}

		resp.Response = &clsv20201016.DescribeMachineGroupsResponseParams{
			MachineGroups: groups,
			TotalCount:    ptrInt64MGDS(int64(2*int(limit) + page3Count)),
			RequestId:     ptrStringMGDS("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaForClsMachineGroupsDS()
	res := cls.DataSourceTencentCloudClsMachineGroups()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.NotEmpty(t, d.Id())

	machineGroups := d.Get("machine_groups").([]interface{})
	expectedLen := 2*int(limit) + page3Count
	assert.Len(t, machineGroups, expectedLen)

	// spot-check first and last
	first := machineGroups[0].(map[string]interface{})
	assert.Equal(t, "mg-page1-0", first["group_id"].(string))

	last := machineGroups[expectedLen-1].(map[string]interface{})
	assert.Equal(t, "mg-page3-2", last["group_id"].(string))

	totalCount := d.Get("total_count").(int)
	assert.Equal(t, expectedLen, totalCount)
}

func TestClsMachineGroupsDS_Schema(t *testing.T) {
	res := cls.DataSourceTencentCloudClsMachineGroups()

	assert.NotNil(t, res)
	assert.Contains(t, res.Schema, "filters")
	assert.Contains(t, res.Schema, "result_output_file")
	assert.Contains(t, res.Schema, "machine_groups")
	assert.Contains(t, res.Schema, "total_count")

	filtersSchema := res.Schema["filters"]
	assert.Equal(t, schema.TypeList, filtersSchema.Type)
	assert.True(t, filtersSchema.Optional)

	machineGroupsSchema := res.Schema["machine_groups"]
	assert.Equal(t, schema.TypeList, machineGroupsSchema.Type)
	assert.True(t, machineGroupsSchema.Computed)

	elemRes := machineGroupsSchema.Elem.(*schema.Resource)
	assert.Contains(t, elemRes.Schema, "group_id")
	assert.Contains(t, elemRes.Schema, "group_name")
	assert.Contains(t, elemRes.Schema, "machine_group_type")
	assert.Contains(t, elemRes.Schema, "create_time")
	assert.Contains(t, elemRes.Schema, "tags")
	assert.Contains(t, elemRes.Schema, "auto_update")
	assert.Contains(t, elemRes.Schema, "update_start_time")
	assert.Contains(t, elemRes.Schema, "update_end_time")
	assert.Contains(t, elemRes.Schema, "service_logging")
	assert.Contains(t, elemRes.Schema, "delay_cleanup_time")
	assert.Contains(t, elemRes.Schema, "meta_tags")
	assert.Contains(t, elemRes.Schema, "os_type")

	totalCountSchema := res.Schema["total_count"]
	assert.Equal(t, schema.TypeInt, totalCountSchema.Type)
	assert.True(t, totalCountSchema.Computed)
}
