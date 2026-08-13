package cfw_test

import (
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	cfwv20190904 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cfw/v20190904"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	svccfw "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/cfw"
)

// mockMetaCfwClusterNatFwSwitch implements tccommon.ProviderMeta
type mockMetaCfwClusterNatFwSwitch struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMetaCfwClusterNatFwSwitch) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockMetaCfwClusterNatFwSwitch{}

func ptrCfwNatFwSwitchString(s string) *string {
	return &s
}

func ptrCfwNatFwSwitchInt64(i int64) *int64 {
	return &i
}

// go test ./tencentcloud/services/cfw/ -run "TestCfwClusterNatFwSwitch" -v -count=1 -gcflags="all=-l"

// TestCfwClusterNatFwSwitch_Create tests the Create function
func TestCfwClusterNatFwSwitch_Create(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	cfwClient := &cfwv20190904.Client{}
	mockClient := &connectivity.TencentCloudClient{}
	patches.ApplyMethodReturn(mockClient, "UseCfwV20190904Client", cfwClient)

	patches.ApplyMethodFunc(cfwClient, "OpenClusterNatFwSwitchWithContext", func(_ interface{}, request *cfwv20190904.OpenClusterNatFwSwitchRequest) (*cfwv20190904.OpenClusterNatFwSwitchResponse, error) {
		resp := cfwv20190904.NewOpenClusterNatFwSwitchResponse()
		resp.Response = &cfwv20190904.OpenClusterNatFwSwitchResponseParams{
			RequestId: ptrCfwNatFwSwitchString("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(cfwClient, "DescribeNatCcnFwSwitch", func(request *cfwv20190904.DescribeNatCcnFwSwitchRequest) (*cfwv20190904.DescribeNatCcnFwSwitchResponse, error) {
		resp := cfwv20190904.NewDescribeNatCcnFwSwitchResponse()
		resp.Response = &cfwv20190904.DescribeNatCcnFwSwitchResponseParams{
			SwitchMode:  ptrCfwNatFwSwitchInt64(1),
			RoutingMode: ptrCfwNatFwSwitchInt64(1),
			Bypass:      ptrCfwNatFwSwitchInt64(0),
			CcnId:       ptrCfwNatFwSwitchString("ccn-test001"),
			RequestId:   ptrCfwNatFwSwitchString("fake-request-id"),
		}
		return resp, nil
	})

	meta := &mockMetaCfwClusterNatFwSwitch{client: mockClient}
	res := svccfw.ResourceTencentCloudCfwClusterNatFwSwitch()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"nat_ccn_switch": []interface{}{
			map[string]interface{}{
				"nat_ins_id":           "cfwnat-test001",
				"ccn_id":               "ccn-test001",
				"switch_mode":          1,
				"routing_mode":         1,
				"access_instance_list": []interface{}{},
				"lead_vpc_cidr":        "",
			},
		},
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "cfwnat-test001"+tccommon.FILED_SP+"ccn-test001", d.Id())
}

// TestCfwClusterNatFwSwitch_Read tests the Read function
func TestCfwClusterNatFwSwitch_Read(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	cfwClient := &cfwv20190904.Client{}
	mockClient := &connectivity.TencentCloudClient{}
	patches.ApplyMethodReturn(mockClient, "UseCfwV20190904Client", cfwClient)

	patches.ApplyMethodFunc(cfwClient, "DescribeNatCcnFwSwitch", func(request *cfwv20190904.DescribeNatCcnFwSwitchRequest) (*cfwv20190904.DescribeNatCcnFwSwitchResponse, error) {
		resp := cfwv20190904.NewDescribeNatCcnFwSwitchResponse()
		resp.Response = &cfwv20190904.DescribeNatCcnFwSwitchResponseParams{
			SwitchMode:  ptrCfwNatFwSwitchInt64(1),
			RoutingMode: ptrCfwNatFwSwitchInt64(1),
			Bypass:      ptrCfwNatFwSwitchInt64(0),
			CcnId:       ptrCfwNatFwSwitchString("ccn-test001"),
			RequestId:   ptrCfwNatFwSwitchString("fake-request-id"),
		}
		return resp, nil
	})

	meta := &mockMetaCfwClusterNatFwSwitch{client: mockClient}
	res := svccfw.ResourceTencentCloudCfwClusterNatFwSwitch()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
	d.SetId("cfwnat-test001" + tccommon.FILED_SP + "ccn-test001")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), int64(d.Get("switch_mode").(int)))
	assert.Equal(t, int64(1), int64(d.Get("routing_mode").(int)))
	assert.Equal(t, int64(0), int64(d.Get("bypass").(int)))
	assert.Equal(t, "ccn-test001", d.Get("ccn_id").(string))
}

// TestCfwClusterNatFwSwitch_Update tests the Update function
func TestCfwClusterNatFwSwitch_Update(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	cfwClient := &cfwv20190904.Client{}
	mockClient := &connectivity.TencentCloudClient{}
	patches.ApplyMethodReturn(mockClient, "UseCfwV20190904Client", cfwClient)

	patches.ApplyMethodFunc(cfwClient, "ModifyClusterNatFwSwitchWithContext", func(_ interface{}, request *cfwv20190904.ModifyClusterNatFwSwitchRequest) (*cfwv20190904.ModifyClusterNatFwSwitchResponse, error) {
		resp := cfwv20190904.NewModifyClusterNatFwSwitchResponse()
		resp.Response = &cfwv20190904.ModifyClusterNatFwSwitchResponseParams{
			RequestId: ptrCfwNatFwSwitchString("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(cfwClient, "DescribeNatCcnFwSwitch", func(request *cfwv20190904.DescribeNatCcnFwSwitchRequest) (*cfwv20190904.DescribeNatCcnFwSwitchResponse, error) {
		resp := cfwv20190904.NewDescribeNatCcnFwSwitchResponse()
		resp.Response = &cfwv20190904.DescribeNatCcnFwSwitchResponseParams{
			SwitchMode:  ptrCfwNatFwSwitchInt64(2),
			RoutingMode: ptrCfwNatFwSwitchInt64(0),
			Bypass:      ptrCfwNatFwSwitchInt64(0),
			CcnId:       ptrCfwNatFwSwitchString("ccn-test001"),
			RequestId:   ptrCfwNatFwSwitchString("fake-request-id"),
		}
		return resp, nil
	})

	meta := &mockMetaCfwClusterNatFwSwitch{client: mockClient}
	res := svccfw.ResourceTencentCloudCfwClusterNatFwSwitch()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"nat_ccn_switch": []interface{}{
			map[string]interface{}{
				"nat_ins_id":           "cfwnat-test001",
				"ccn_id":               "ccn-test001",
				"switch_mode":          2,
				"routing_mode":         0,
				"access_instance_list": []interface{}{},
				"lead_vpc_cidr":        "",
			},
		},
	})
	d.SetId("cfwnat-test001" + tccommon.FILED_SP + "ccn-test001")

	err := res.Update(d, meta)
	assert.NoError(t, err)
}

// TestCfwClusterNatFwSwitch_Delete tests the Delete function
func TestCfwClusterNatFwSwitch_Delete(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	cfwClient := &cfwv20190904.Client{}
	mockClient := &connectivity.TencentCloudClient{}
	patches.ApplyMethodReturn(mockClient, "UseCfwV20190904Client", cfwClient)

	patches.ApplyMethodFunc(cfwClient, "CloseClusterNatFwSwitchWithContext", func(_ interface{}, request *cfwv20190904.CloseClusterNatFwSwitchRequest) (*cfwv20190904.CloseClusterNatFwSwitchResponse, error) {
		resp := cfwv20190904.NewCloseClusterNatFwSwitchResponse()
		resp.Response = &cfwv20190904.CloseClusterNatFwSwitchResponseParams{
			RequestId: ptrCfwNatFwSwitchString("fake-request-id"),
		}
		return resp, nil
	})

	meta := &mockMetaCfwClusterNatFwSwitch{client: mockClient}
	res := svccfw.ResourceTencentCloudCfwClusterNatFwSwitch()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
	d.SetId("cfwnat-test001" + tccommon.FILED_SP + "ccn-test001")

	err := res.Delete(d, meta)
	assert.NoError(t, err)
}

// TestCfwClusterNatFwSwitch_Schema tests the schema definition
func TestCfwClusterNatFwSwitch_Schema(t *testing.T) {
	res := svccfw.ResourceTencentCloudCfwClusterNatFwSwitch()

	assert.NotNil(t, res)
	assert.Contains(t, res.Schema, "nat_ccn_switch")
	assert.Contains(t, res.Schema, "switch_mode")
	assert.Contains(t, res.Schema, "routing_mode")
	assert.Contains(t, res.Schema, "bypass")
	assert.Contains(t, res.Schema, "ccn_id")

	natCcnSwitch := res.Schema["nat_ccn_switch"]
	assert.Equal(t, schema.TypeList, natCcnSwitch.Type)
	assert.True(t, natCcnSwitch.Required)
}
