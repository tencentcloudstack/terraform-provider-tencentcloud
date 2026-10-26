package tke_test

import (
	"context"
	"errors"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	tkev20180525 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tke/v20180525"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	tkesvc "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/tke"
)

// go test ./tencentcloud/services/tke/ -run "TestTkeAddonValuesDataSource" -v -count=1 -gcflags="all=-l"

type mockMetaForTkeAddonValuesDS struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMetaForTkeAddonValuesDS) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockMetaForTkeAddonValuesDS{}

func newMockMetaForTkeAddonValuesDS() *mockMetaForTkeAddonValuesDS {
	return &mockMetaForTkeAddonValuesDS{client: &connectivity.TencentCloudClient{}}
}

func ptrStrTkeAddonValuesDS(s string) *string { return &s }

// TestTkeAddonValuesDataSource_ReadBasic verifies that a populated response is
// flattened into values/default_values and that the id is cluster_id#addon_name.
func TestTkeAddonValuesDataSource_ReadBasic(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	patches.ApplyMethodFunc(&tkesvc.TkeService{}, "DescribeAddonValuesList", func(_ context.Context, clusterId, addonName string) (*tkev20180525.DescribeAddonValuesResponseParams, error) {
		assert.Equal(t, "cls-5yezvaxo", clusterId)
		assert.Equal(t, "ip-masq-agent", addonName)
		return &tkev20180525.DescribeAddonValuesResponseParams{
			Values:        ptrStrTkeAddonValuesDS(`{"nonMasqueradeCIDRs":["10.0.0.0/8"]}`),
			DefaultValues: ptrStrTkeAddonValuesDS(`{"nonMasqueradeCIDRs":["10.0.0.0/8"],"masqLinkLocal":false}`),
		}, nil
	})

	meta := newMockMetaForTkeAddonValuesDS()
	res := tkesvc.DataSourceTencentCloudTkeAddonValues()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"cluster_id": "cls-5yezvaxo",
		"addon_name": "ip-masq-agent",
	})

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, `{"nonMasqueradeCIDRs":["10.0.0.0/8"]}`, d.Get("values").(string))
	assert.Equal(t, `{"nonMasqueradeCIDRs":["10.0.0.0/8"],"masqLinkLocal":false}`, d.Get("default_values").(string))
	assert.Equal(t, "cls-5yezvaxo#ip-masq-agent", d.Id())
}

// TestTkeAddonValuesDataSource_ReadWithNilValues verifies that a nil Values field
// (addon not installed) is skipped without error while default_values is still set.
func TestTkeAddonValuesDataSource_ReadWithNilValues(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	patches.ApplyMethodFunc(&tkesvc.TkeService{}, "DescribeAddonValuesList", func(_ context.Context, _ string, _ string) (*tkev20180525.DescribeAddonValuesResponseParams, error) {
		return &tkev20180525.DescribeAddonValuesResponseParams{
			DefaultValues: ptrStrTkeAddonValuesDS(`{"replicas":1}`),
		}, nil
	})

	meta := newMockMetaForTkeAddonValuesDS()
	res := tkesvc.DataSourceTencentCloudTkeAddonValues()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"cluster_id": "cls-5yezvaxo",
		"addon_name": "ip-masq-agent",
	})

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "", d.Get("values").(string))
	assert.Equal(t, `{"replicas":1}`, d.Get("default_values").(string))
	assert.Equal(t, "cls-5yezvaxo#ip-masq-agent", d.Id())
}

// TestTkeAddonValuesDataSource_ReadWithNilDefaultValues verifies that a nil
// DefaultValues field is skipped without error.
func TestTkeAddonValuesDataSource_ReadWithNilDefaultValues(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	patches.ApplyMethodFunc(&tkesvc.TkeService{}, "DescribeAddonValuesList", func(_ context.Context, _ string, _ string) (*tkev20180525.DescribeAddonValuesResponseParams, error) {
		return &tkev20180525.DescribeAddonValuesResponseParams{
			Values: ptrStrTkeAddonValuesDS(`{"replicas":2}`),
		}, nil
	})

	meta := newMockMetaForTkeAddonValuesDS()
	res := tkesvc.DataSourceTencentCloudTkeAddonValues()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"cluster_id": "cls-5yezvaxo",
		"addon_name": "ip-masq-agent",
	})

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, `{"replicas":2}`, d.Get("values").(string))
	assert.Equal(t, "", d.Get("default_values").(string))
	assert.Equal(t, "cls-5yezvaxo#ip-masq-agent", d.Id())
}

// TestTkeAddonValuesDataSource_ReadError verifies that a service layer error is propagated.
func TestTkeAddonValuesDataSource_ReadError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	patches.ApplyMethodFunc(&tkesvc.TkeService{}, "DescribeAddonValuesList", func(_ context.Context, _ string, _ string) (*tkev20180525.DescribeAddonValuesResponseParams, error) {
		return nil, errors.New("ResourceUnavailable")
	})

	meta := newMockMetaForTkeAddonValuesDS()
	res := tkesvc.DataSourceTencentCloudTkeAddonValues()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"cluster_id": "cls-5yezvaxo",
		"addon_name": "ip-masq-agent",
	})

	err := res.Read(d, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "ResourceUnavailable")
}

// TestTkeAddonValuesDataSource_ReadNilResponse verifies that a nil response from
// the service layer is reported as an error instead of silently clearing the id.
func TestTkeAddonValuesDataSource_ReadNilResponse(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	patches.ApplyMethodFunc(&tkesvc.TkeService{}, "DescribeAddonValuesList", func(_ context.Context, _ string, _ string) (*tkev20180525.DescribeAddonValuesResponseParams, error) {
		return nil, nil
	})

	meta := newMockMetaForTkeAddonValuesDS()
	res := tkesvc.DataSourceTencentCloudTkeAddonValues()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"cluster_id": "cls-5yezvaxo",
		"addon_name": "ip-masq-agent",
	})

	err := res.Read(d, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "DescribeAddonValuesResponse is nil")
	assert.Equal(t, "", d.Id())
}

// TestTkeAddonValuesDataSource_Schema verifies the schema definition of the data source.
func TestTkeAddonValuesDataSource_Schema(t *testing.T) {
	res := tkesvc.DataSourceTencentCloudTkeAddonValues()
	assert.NotNil(t, res)
	assert.NotNil(t, res.Read)

	expected := map[string]struct {
		fieldType schema.ValueType
		required  bool
		optional  bool
		computed  bool
	}{
		"cluster_id":         {schema.TypeString, true, false, false},
		"addon_name":         {schema.TypeString, true, false, false},
		"values":             {schema.TypeString, false, false, true},
		"default_values":     {schema.TypeString, false, false, true},
		"result_output_file": {schema.TypeString, false, true, false},
	}

	assert.Len(t, res.Schema, len(expected))
	for name, want := range expected {
		got := res.Schema[name]
		assert.NotNil(t, got, "schema field %s is missing", name)
		if got == nil {
			continue
		}

		assert.Equal(t, want.fieldType, got.Type, "field %s type mismatch", name)
		assert.Equal(t, want.required, got.Required, "field %s Required mismatch", name)
		assert.Equal(t, want.optional, got.Optional, "field %s Optional mismatch", name)
		assert.Equal(t, want.computed, got.Computed, "field %s Computed mismatch", name)
	}
}
