package clb_test

import (
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	clb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/clb/v20180317"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	localclb "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/clb"
)

// mockMetaTargetGroup implements tccommon.ProviderMeta
type mockMetaTargetGroup struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMetaTargetGroup) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockMetaTargetGroup{}

func newMockMetaTargetGroup() *mockMetaTargetGroup {
	return &mockMetaTargetGroup{client: &connectivity.TencentCloudClient{}}
}

func ptrStringTargetGroup(s string) *string {
	return &s
}

func ptrBoolTargetGroup(b bool) *bool {
	return &b
}

func ptrUint64TargetGroup(v uint64) *uint64 {
	return &v
}

// go test ./tencentcloud/services/clb/ -run "TestClbTargetGroupSnatEnable" -v -count=1 -gcflags="all=-l"

// TestClbTargetGroupSnatEnable_CreateWithSnatEnableTrue tests Create with snat_enable = true
func TestClbTargetGroupSnatEnable_CreateWithSnatEnableTrue(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	clbClient := &clb.Client{}
	patches.ApplyMethodReturn(newMockMetaTargetGroup().client, "UseClbClient", clbClient)

	patches.ApplyMethodFunc(clbClient, "CreateTargetGroup", func(request *clb.CreateTargetGroupRequest) (*clb.CreateTargetGroupResponse, error) {
		assert.NotNil(t, request.SnatEnable)
		assert.Equal(t, true, *request.SnatEnable)
		assert.Equal(t, "test-snat-tg", *request.TargetGroupName)

		resp := clb.NewCreateTargetGroupResponse()
		resp.Response = &clb.CreateTargetGroupResponseParams{
			TargetGroupId: ptrStringTargetGroup("lbtg-12345678"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(clbClient, "DescribeTargetGroupList", func(request *clb.DescribeTargetGroupListRequest) (*clb.DescribeTargetGroupListResponse, error) {
		resp := clb.NewDescribeTargetGroupListResponse()
		resp.Response = &clb.DescribeTargetGroupListResponseParams{
			TotalCount: ptrUint64TargetGroup(1),
			TargetGroupSet: []*clb.TargetGroupInfo{
				{
					TargetGroupId:   ptrStringTargetGroup("lbtg-12345678"),
					TargetGroupName: ptrStringTargetGroup("test-snat-tg"),
					VpcId:           ptrStringTargetGroup("vpc-xxxxxx"),
					Port:            ptrUint64TargetGroup(80),
					TargetGroupType: ptrStringTargetGroup("v2"),
					Protocol:        ptrStringTargetGroup("TCP"),
					SnatEnable:      ptrBoolTargetGroup(true),
				},
			},
		}
		return resp, nil
	})

	meta := newMockMetaTargetGroup()
	res := localclb.ResourceTencentCloudClbTargetGroup()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"target_group_name": "test-snat-tg",
		"vpc_id":            "vpc-xxxxxx",
		"port":              80,
		"type":              "v2",
		"protocol":          "TCP",
		"snat_enable":       true,
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "lbtg-12345678", d.Id())
	assert.Equal(t, true, d.Get("snat_enable"))
}

// TestClbTargetGroupSnatEnable_CreateWithSnatEnableFalse tests Create with snat_enable = false
func TestClbTargetGroupSnatEnable_CreateWithSnatEnableFalse(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	clbClient := &clb.Client{}
	patches.ApplyMethodReturn(newMockMetaTargetGroup().client, "UseClbClient", clbClient)

	patches.ApplyMethodFunc(clbClient, "CreateTargetGroup", func(request *clb.CreateTargetGroupRequest) (*clb.CreateTargetGroupResponse, error) {
		assert.NotNil(t, request.SnatEnable)
		assert.Equal(t, false, *request.SnatEnable)

		resp := clb.NewCreateTargetGroupResponse()
		resp.Response = &clb.CreateTargetGroupResponseParams{
			TargetGroupId: ptrStringTargetGroup("lbtg-22345678"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(clbClient, "DescribeTargetGroupList", func(request *clb.DescribeTargetGroupListRequest) (*clb.DescribeTargetGroupListResponse, error) {
		resp := clb.NewDescribeTargetGroupListResponse()
		resp.Response = &clb.DescribeTargetGroupListResponseParams{
			TotalCount: ptrUint64TargetGroup(1),
			TargetGroupSet: []*clb.TargetGroupInfo{
				{
					TargetGroupId:   ptrStringTargetGroup("lbtg-22345678"),
					TargetGroupName: ptrStringTargetGroup("test-snat-false"),
					VpcId:           ptrStringTargetGroup("vpc-xxxxxx"),
					Port:            ptrUint64TargetGroup(80),
					TargetGroupType: ptrStringTargetGroup("v2"),
					Protocol:        ptrStringTargetGroup("TCP"),
					SnatEnable:      ptrBoolTargetGroup(false),
				},
			},
		}
		return resp, nil
	})

	meta := newMockMetaTargetGroup()
	res := localclb.ResourceTencentCloudClbTargetGroup()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"target_group_name": "test-snat-false",
		"vpc_id":            "vpc-xxxxxx",
		"port":              80,
		"type":              "v2",
		"protocol":          "TCP",
		"snat_enable":       false,
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "lbtg-22345678", d.Id())
	assert.Equal(t, false, d.Get("snat_enable"))
}

// TestClbTargetGroupSnatEnable_CreateWithoutSnatEnable tests Create without snat_enable specified
func TestClbTargetGroupSnatEnable_CreateWithoutSnatEnable(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	clbClient := &clb.Client{}
	patches.ApplyMethodReturn(newMockMetaTargetGroup().client, "UseClbClient", clbClient)

	patches.ApplyMethodFunc(clbClient, "CreateTargetGroup", func(request *clb.CreateTargetGroupRequest) (*clb.CreateTargetGroupResponse, error) {
		assert.Nil(t, request.SnatEnable)

		resp := clb.NewCreateTargetGroupResponse()
		resp.Response = &clb.CreateTargetGroupResponseParams{
			TargetGroupId: ptrStringTargetGroup("lbtg-32345678"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(clbClient, "DescribeTargetGroupList", func(request *clb.DescribeTargetGroupListRequest) (*clb.DescribeTargetGroupListResponse, error) {
		resp := clb.NewDescribeTargetGroupListResponse()
		resp.Response = &clb.DescribeTargetGroupListResponseParams{
			TotalCount: ptrUint64TargetGroup(1),
			TargetGroupSet: []*clb.TargetGroupInfo{
				{
					TargetGroupId:   ptrStringTargetGroup("lbtg-32345678"),
					TargetGroupName: ptrStringTargetGroup("test-no-snat"),
					VpcId:           ptrStringTargetGroup("vpc-xxxxxx"),
					Port:            ptrUint64TargetGroup(80),
					TargetGroupType: ptrStringTargetGroup("v2"),
					Protocol:        ptrStringTargetGroup("TCP"),
				},
			},
		}
		return resp, nil
	})

	meta := newMockMetaTargetGroup()
	res := localclb.ResourceTencentCloudClbTargetGroup()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"target_group_name": "test-no-snat",
		"vpc_id":            "vpc-xxxxxx",
		"port":              80,
		"type":              "v2",
		"protocol":          "TCP",
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "lbtg-32345678", d.Id())
}

// TestClbTargetGroupSnatEnable_ReadWithSnatEnable tests Read with SnatEnable in response
func TestClbTargetGroupSnatEnable_ReadWithSnatEnable(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	clbClient := &clb.Client{}
	patches.ApplyMethodReturn(newMockMetaTargetGroup().client, "UseClbClient", clbClient)

	patches.ApplyMethodFunc(clbClient, "DescribeTargetGroupList", func(request *clb.DescribeTargetGroupListRequest) (*clb.DescribeTargetGroupListResponse, error) {
		resp := clb.NewDescribeTargetGroupListResponse()
		resp.Response = &clb.DescribeTargetGroupListResponseParams{
			TotalCount: ptrUint64TargetGroup(1),
			TargetGroupSet: []*clb.TargetGroupInfo{
				{
					TargetGroupId:   ptrStringTargetGroup("lbtg-12345678"),
					TargetGroupName: ptrStringTargetGroup("test-read-snat"),
					VpcId:           ptrStringTargetGroup("vpc-xxxxxx"),
					Port:            ptrUint64TargetGroup(80),
					TargetGroupType: ptrStringTargetGroup("v2"),
					Protocol:        ptrStringTargetGroup("TCP"),
					SnatEnable:      ptrBoolTargetGroup(true),
				},
			},
		}
		return resp, nil
	})

	meta := newMockMetaTargetGroup()
	res := localclb.ResourceTencentCloudClbTargetGroup()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"target_group_name": "test-read-snat",
	})
	d.SetId("lbtg-12345678")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, true, d.Get("snat_enable"))
}

// TestClbTargetGroupSnatEnable_ReadWithSnatEnableNil tests Read when SnatEnable is nil
func TestClbTargetGroupSnatEnable_ReadWithSnatEnableNil(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	clbClient := &clb.Client{}
	patches.ApplyMethodReturn(newMockMetaTargetGroup().client, "UseClbClient", clbClient)

	patches.ApplyMethodFunc(clbClient, "DescribeTargetGroupList", func(request *clb.DescribeTargetGroupListRequest) (*clb.DescribeTargetGroupListResponse, error) {
		resp := clb.NewDescribeTargetGroupListResponse()
		resp.Response = &clb.DescribeTargetGroupListResponseParams{
			TotalCount: ptrUint64TargetGroup(1),
			TargetGroupSet: []*clb.TargetGroupInfo{
				{
					TargetGroupId:   ptrStringTargetGroup("lbtg-12345678"),
					TargetGroupName: ptrStringTargetGroup("test-read-no-snat"),
					VpcId:           ptrStringTargetGroup("vpc-xxxxxx"),
					Port:            ptrUint64TargetGroup(80),
					TargetGroupType: ptrStringTargetGroup("v2"),
					Protocol:        ptrStringTargetGroup("TCP"),
				},
			},
		}
		return resp, nil
	})

	meta := newMockMetaTargetGroup()
	res := localclb.ResourceTencentCloudClbTargetGroup()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"target_group_name": "test-read-no-snat",
	})
	d.SetId("lbtg-12345678")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, false, d.Get("snat_enable"))
}

// TestClbTargetGroupSnatEnable_UpdateSnatEnable tests Update with snat_enable change
func TestClbTargetGroupSnatEnable_UpdateSnatEnable(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	clbClient := &clb.Client{}
	patches.ApplyMethodReturn(newMockMetaTargetGroup().client, "UseClbClient", clbClient)

	modifyCalled := false
	patches.ApplyMethodFunc(clbClient, "ModifyTargetGroupAttribute", func(request *clb.ModifyTargetGroupAttributeRequest) (*clb.ModifyTargetGroupAttributeResponse, error) {
		assert.NotNil(t, request.TargetGroupId)
		assert.Equal(t, "lbtg-12345678", *request.TargetGroupId)
		assert.NotNil(t, request.SnatEnable)
		assert.Equal(t, true, *request.SnatEnable)
		modifyCalled = true

		resp := clb.NewModifyTargetGroupAttributeResponse()
		resp.Response = &clb.ModifyTargetGroupAttributeResponseParams{
			RequestId: ptrStringTargetGroup("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(clbClient, "DescribeTargetGroupList", func(request *clb.DescribeTargetGroupListRequest) (*clb.DescribeTargetGroupListResponse, error) {
		resp := clb.NewDescribeTargetGroupListResponse()
		resp.Response = &clb.DescribeTargetGroupListResponseParams{
			TotalCount: ptrUint64TargetGroup(1),
			TargetGroupSet: []*clb.TargetGroupInfo{
				{
					TargetGroupId:   ptrStringTargetGroup("lbtg-12345678"),
					TargetGroupName: ptrStringTargetGroup("test-update-snat"),
					VpcId:           ptrStringTargetGroup("vpc-xxxxxx"),
					Port:            ptrUint64TargetGroup(80),
					TargetGroupType: ptrStringTargetGroup("v2"),
					Protocol:        ptrStringTargetGroup("TCP"),
					SnatEnable:      ptrBoolTargetGroup(true),
				},
			},
		}
		return resp, nil
	})

	meta := newMockMetaTargetGroup()
	res := localclb.ResourceTencentCloudClbTargetGroup()

	// Use TestResourceDataRaw which marks all fields as changed (simulating a diff)
	// Exclude immutable fields (vpc_id, full_listen_switch, ip_version) to avoid immutable field error
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"target_group_name": "test-update-snat",
		"port":              80,
		"snat_enable":       true,
	})
	d.SetId("lbtg-12345678")

	err := res.Update(d, meta)
	assert.NoError(t, err)
	assert.True(t, modifyCalled, "ModifyTargetGroupAttribute should have been called")
	assert.Equal(t, true, d.Get("snat_enable"))
}
