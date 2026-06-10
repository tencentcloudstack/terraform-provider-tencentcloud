package cfw_test

import (
	"context"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	cfwv20190904 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cfw/v20190904"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/cfw"
)

type mockMetaClusterFwBypass struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMetaClusterFwBypass) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockMetaClusterFwBypass{}

func newMockMetaClusterFwBypass() *mockMetaClusterFwBypass {
	return &mockMetaClusterFwBypass{client: &connectivity.TencentCloudClient{}}
}

func ptrStringCFB(s string) *string {
	return &s
}

func ptrInt64CFB(i int64) *int64 {
	return &i
}

func TestCfwClusterFwBypassConfig_Read_Normal(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	cfwClient := &cfwv20190904.Client{}
	patches.ApplyMethodReturn(newMockMetaClusterFwBypass().client, "UseCfwV20190904Client", cfwClient)

	patches.ApplyMethodFunc(cfwClient, "DescribeClusterNatCcnFwSwitchListWithContext", func(_ context.Context, request *cfwv20190904.DescribeClusterNatCcnFwSwitchListRequest) (*cfwv20190904.DescribeClusterNatCcnFwSwitchListResponse, error) {
		resp := cfwv20190904.NewDescribeClusterNatCcnFwSwitchListResponse()
		total := int64(1)
		bypass := int64(0)
		status := int64(1)
		resp.Response = &cfwv20190904.DescribeClusterNatCcnFwSwitchListResponseParams{
			RequestId: ptrStringCFB("fake-request-id"),
			Total:     &total,
			Data: []*cfwv20190904.NatFwSwitchDetailS{
				{
					InsObj:     ptrStringCFB("nat-xxxxxxxx"),
					ObjName:    ptrStringCFB("test-nat"),
					FwType:     ptrStringCFB("NAT_FW"),
					AssetType:  ptrStringCFB("nat_ccn"),
					Region:     ptrStringCFB("ap-guangzhou"),
					Status:     &status,
					Bypass:     &bypass,
					AttachId:   ptrStringCFB("ccn-xxxxxxxx"),
					AttachName: ptrStringCFB("test-ccn"),
				},
			},
		}
		return resp, nil
	})

	meta := newMockMetaClusterFwBypass()
	res := cfw.ResourceTencentCloudCfwClusterFwBypassConfig()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"fw_type":    "NAT_FW",
		"ccn_id":     "ccn-xxxxxxxx",
		"nat_ins_id": "nat-xxxxxxxx",
		"enable":     false,
	})
	d.SetId("NAT_FW#ccn-xxxxxxxx#nat-xxxxxxxx")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "NAT_FW#ccn-xxxxxxxx#nat-xxxxxxxx", d.Id())
	assert.Equal(t, false, d.Get("enable"))
}

func TestCfwClusterFwBypassConfig_Read_BypassEnabled(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	cfwClient := &cfwv20190904.Client{}
	patches.ApplyMethodReturn(newMockMetaClusterFwBypass().client, "UseCfwV20190904Client", cfwClient)

	patches.ApplyMethodFunc(cfwClient, "DescribeClusterNatCcnFwSwitchListWithContext", func(_ context.Context, request *cfwv20190904.DescribeClusterNatCcnFwSwitchListRequest) (*cfwv20190904.DescribeClusterNatCcnFwSwitchListResponse, error) {
		resp := cfwv20190904.NewDescribeClusterNatCcnFwSwitchListResponse()
		total := int64(1)
		bypass := int64(1)
		status := int64(1)
		resp.Response = &cfwv20190904.DescribeClusterNatCcnFwSwitchListResponseParams{
			RequestId: ptrStringCFB("fake-request-id"),
			Total:     &total,
			Data: []*cfwv20190904.NatFwSwitchDetailS{
				{
					InsObj:     ptrStringCFB("nat-xxxxxxxx"),
					ObjName:    ptrStringCFB("test-nat"),
					FwType:     ptrStringCFB("NAT_FW"),
					AssetType:  ptrStringCFB("nat_ccn"),
					Region:     ptrStringCFB("ap-guangzhou"),
					Status:     &status,
					Bypass:     &bypass,
					AttachId:   ptrStringCFB("ccn-xxxxxxxx"),
					AttachName: ptrStringCFB("test-ccn"),
				},
			},
		}
		return resp, nil
	})

	meta := newMockMetaClusterFwBypass()
	res := cfw.ResourceTencentCloudCfwClusterFwBypassConfig()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"fw_type":    "NAT_FW",
		"ccn_id":     "ccn-xxxxxxxx",
		"nat_ins_id": "nat-xxxxxxxx",
		"enable":     true,
	})
	d.SetId("NAT_FW#ccn-xxxxxxxx#nat-xxxxxxxx")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "NAT_FW#ccn-xxxxxxxx#nat-xxxxxxxx", d.Id())
	assert.Equal(t, true, d.Get("enable"))
}

func TestCfwClusterFwBypassConfig_Read_NotFound(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	cfwClient := &cfwv20190904.Client{}
	patches.ApplyMethodReturn(newMockMetaClusterFwBypass().client, "UseCfwV20190904Client", cfwClient)

	patches.ApplyMethodFunc(cfwClient, "DescribeClusterNatCcnFwSwitchListWithContext", func(_ context.Context, request *cfwv20190904.DescribeClusterNatCcnFwSwitchListRequest) (*cfwv20190904.DescribeClusterNatCcnFwSwitchListResponse, error) {
		resp := cfwv20190904.NewDescribeClusterNatCcnFwSwitchListResponse()
		total := int64(0)
		resp.Response = &cfwv20190904.DescribeClusterNatCcnFwSwitchListResponseParams{
			RequestId: ptrStringCFB("fake-request-id"),
			Total:     &total,
			Data:      []*cfwv20190904.NatFwSwitchDetailS{},
		}
		return resp, nil
	})

	meta := newMockMetaClusterFwBypass()
	res := cfw.ResourceTencentCloudCfwClusterFwBypassConfig()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"fw_type": "VPC_FW",
		"ccn_id":  "ccn-notexist",
		"enable":  false,
	})
	d.SetId("VPC_FW#ccn-notexist")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "", d.Id())
}

func TestCfwClusterFwBypassConfig_Read_VpcFw_Normal(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	cfwClient := &cfwv20190904.Client{}
	patches.ApplyMethodReturn(newMockMetaClusterFwBypass().client, "UseCfwV20190904Client", cfwClient)

	patches.ApplyMethodFunc(cfwClient, "DescribeClusterNatCcnFwSwitchListWithContext", func(_ context.Context, request *cfwv20190904.DescribeClusterNatCcnFwSwitchListRequest) (*cfwv20190904.DescribeClusterNatCcnFwSwitchListResponse, error) {
		resp := cfwv20190904.NewDescribeClusterNatCcnFwSwitchListResponse()
		total := int64(1)
		bypass := int64(1)
		status := int64(1)
		resp.Response = &cfwv20190904.DescribeClusterNatCcnFwSwitchListResponseParams{
			RequestId: ptrStringCFB("fake-request-id"),
			Total:     &total,
			Data: []*cfwv20190904.NatFwSwitchDetailS{
				{
					InsObj:     ptrStringCFB("vpc-xxxxxxxx"),
					ObjName:    ptrStringCFB("test-vpc"),
					FwType:     ptrStringCFB("VPC_FW"),
					AssetType:  ptrStringCFB("nat_ccn"),
					Region:     ptrStringCFB("ap-beijing"),
					Status:     &status,
					Bypass:     &bypass,
					AttachId:   ptrStringCFB("ccn-xxxxxxxx"),
					AttachName: ptrStringCFB("test-ccn"),
				},
			},
		}
		return resp, nil
	})

	meta := newMockMetaClusterFwBypass()
	res := cfw.ResourceTencentCloudCfwClusterFwBypassConfig()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"fw_type": "VPC_FW",
		"ccn_id":  "ccn-xxxxxxxx",
		"enable":  true,
	})
	d.SetId("VPC_FW#ccn-xxxxxxxx")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "VPC_FW#ccn-xxxxxxxx", d.Id())
	assert.Equal(t, true, d.Get("enable"))
}

func TestCfwClusterFwBypassConfig_Update_Normal(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	cfwClient := &cfwv20190904.Client{}
	patches.ApplyMethodReturn(newMockMetaClusterFwBypass().client, "UseCfwV20190904Client", cfwClient)

	patches.ApplyMethodFunc(cfwClient, "ModifyClusterFwBypassWithContext", func(_ context.Context, request *cfwv20190904.ModifyClusterFwBypassRequest) (*cfwv20190904.ModifyClusterFwBypassResponse, error) {
		assert.NotNil(t, request.FwType)
		assert.Equal(t, "VPC_FW", *request.FwType)
		assert.NotNil(t, request.CcnId)
		assert.Equal(t, "ccn-xxxxxxxx", *request.CcnId)
		assert.NotNil(t, request.Enable)
		assert.Equal(t, false, *request.Enable)

		resp := cfwv20190904.NewModifyClusterFwBypassResponse()
		resp.Response = &cfwv20190904.ModifyClusterFwBypassResponseParams{
			RequestId: ptrStringCFB("fake-request-id-update"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(cfwClient, "DescribeClusterNatCcnFwSwitchListWithContext", func(_ context.Context, request *cfwv20190904.DescribeClusterNatCcnFwSwitchListRequest) (*cfwv20190904.DescribeClusterNatCcnFwSwitchListResponse, error) {
		resp := cfwv20190904.NewDescribeClusterNatCcnFwSwitchListResponse()
		total := int64(1)
		bypass := int64(0)
		status := int64(1)
		resp.Response = &cfwv20190904.DescribeClusterNatCcnFwSwitchListResponseParams{
			RequestId: ptrStringCFB("fake-request-id"),
			Total:     &total,
			Data: []*cfwv20190904.NatFwSwitchDetailS{
				{
					InsObj:     ptrStringCFB("vpc-xxxxxxxx"),
					ObjName:    ptrStringCFB("test-vpc"),
					FwType:     ptrStringCFB("VPC_FW"),
					AssetType:  ptrStringCFB("nat_ccn"),
					Region:     ptrStringCFB("ap-guangzhou"),
					Status:     &status,
					Bypass:     &bypass,
					AttachId:   ptrStringCFB("ccn-xxxxxxxx"),
					AttachName: ptrStringCFB("test-ccn"),
				},
			},
		}
		return resp, nil
	})

	meta := newMockMetaClusterFwBypass()
	res := cfw.ResourceTencentCloudCfwClusterFwBypassConfig()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"fw_type": "VPC_FW",
		"ccn_id":  "ccn-xxxxxxxx",
		"enable":  false,
	})
	d.SetId("VPC_FW#ccn-xxxxxxxx")

	err := res.Update(d, meta)
	assert.NoError(t, err)
}

func TestCfwClusterFwBypassConfig_Update_NatFw(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	cfwClient := &cfwv20190904.Client{}
	patches.ApplyMethodReturn(newMockMetaClusterFwBypass().client, "UseCfwV20190904Client", cfwClient)

	patches.ApplyMethodFunc(cfwClient, "ModifyClusterFwBypassWithContext", func(_ context.Context, request *cfwv20190904.ModifyClusterFwBypassRequest) (*cfwv20190904.ModifyClusterFwBypassResponse, error) {
		assert.NotNil(t, request.FwType)
		assert.Equal(t, "NAT_FW", *request.FwType)
		assert.NotNil(t, request.CcnId)
		assert.Equal(t, "ccn-xxxxxxxx", *request.CcnId)
		assert.NotNil(t, request.NatInsId)
		assert.Equal(t, "nat-xxxxxxxx", *request.NatInsId)
		assert.NotNil(t, request.Enable)
		assert.Equal(t, true, *request.Enable)

		resp := cfwv20190904.NewModifyClusterFwBypassResponse()
		resp.Response = &cfwv20190904.ModifyClusterFwBypassResponseParams{
			RequestId: ptrStringCFB("fake-request-id-nat-update"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(cfwClient, "DescribeClusterNatCcnFwSwitchListWithContext", func(_ context.Context, request *cfwv20190904.DescribeClusterNatCcnFwSwitchListRequest) (*cfwv20190904.DescribeClusterNatCcnFwSwitchListResponse, error) {
		resp := cfwv20190904.NewDescribeClusterNatCcnFwSwitchListResponse()
		total := int64(1)
		bypass := int64(1)
		status := int64(1)
		resp.Response = &cfwv20190904.DescribeClusterNatCcnFwSwitchListResponseParams{
			RequestId: ptrStringCFB("fake-request-id"),
			Total:     &total,
			Data: []*cfwv20190904.NatFwSwitchDetailS{
				{
					InsObj:     ptrStringCFB("nat-xxxxxxxx"),
					ObjName:    ptrStringCFB("test-nat"),
					FwType:     ptrStringCFB("NAT_FW"),
					AssetType:  ptrStringCFB("nat_ccn"),
					Region:     ptrStringCFB("ap-guangzhou"),
					Status:     &status,
					Bypass:     &bypass,
					AttachId:   ptrStringCFB("ccn-xxxxxxxx"),
					AttachName: ptrStringCFB("test-ccn"),
				},
			},
		}
		return resp, nil
	})

	meta := newMockMetaClusterFwBypass()
	res := cfw.ResourceTencentCloudCfwClusterFwBypassConfig()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"fw_type":    "NAT_FW",
		"ccn_id":     "ccn-xxxxxxxx",
		"nat_ins_id": "nat-xxxxxxxx",
		"enable":     true,
	})
	d.SetId("NAT_FW#ccn-xxxxxxxx#nat-xxxxxxxx")

	err := res.Update(d, meta)
	assert.NoError(t, err)
}
