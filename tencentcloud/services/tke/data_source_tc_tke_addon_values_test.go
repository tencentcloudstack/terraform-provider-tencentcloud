package tke_test

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	sdkErrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
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

// TestTkeAddonValuesDataSource_ReadNotInstalled verifies the real-world semantics of
// an addon whose chart exists but which is not installed in the cluster: the cloud API
// still returns a rendered (chart default) payload, so `values` is NOT empty and is
// exactly equal to `default_values`. The data source must pass both through unchanged
// and must still set the id.
func TestTkeAddonValuesDataSource_ReadNotInstalled(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	rendered := `{"global":{"cluster":{"id":"","uin":""}},"agent":{"resources":{"limits":{"memory":"500Mi"}}}}`
	patches.ApplyMethodFunc(&tkesvc.TkeService{}, "DescribeAddonValuesList", func(_ context.Context, _ string, _ string) (*tkev20180525.DescribeAddonValuesResponseParams, error) {
		return &tkev20180525.DescribeAddonValuesResponseParams{
			Values:        ptrStrTkeAddonValuesDS(rendered),
			DefaultValues: ptrStrTkeAddonValuesDS(rendered),
		}, nil
	})

	meta := newMockMetaForTkeAddonValuesDS()
	res := tkesvc.DataSourceTencentCloudTkeAddonValues()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"cluster_id": "cls-dhhfi9n8",
		"addon_name": "tke-log-agent",
	})

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, rendered, d.Get("values").(string), "values must not be cleared for a not-installed addon")
	assert.Equal(t, rendered, d.Get("default_values").(string))

	// The point of this case: the two fields are indistinguishable, and the payload
	// carries no cluster information -> "values is not empty" cannot be used to tell
	// whether the addon is installed.
	var decoded struct {
		Global struct {
			Cluster struct {
				Id *string `json:"id"`
			} `json:"cluster"`
		} `json:"global"`
	}
	assert.NoError(t, json.Unmarshal([]byte(d.Get("values").(string)), &decoded))
	assert.NotNil(t, decoded.Global.Cluster.Id)
	assert.Equal(t, "", *decoded.Global.Cluster.Id, "not-installed payload must not carry a cluster id")
	assert.Equal(t, "cls-dhhfi9n8#tke-log-agent", d.Id())
}

// TestTkeAddonValuesDataSource_ReadNotFoundChartNonRetryable verifies that a
// deterministic cloud error ("not found chart") fails fast instead of exhausting the
// 3-minute read timeout. `ResourceUnavailable` is a generic retryable code, so the
// data source must special-case these messages.
func TestTkeAddonValuesDataSource_ReadNotFoundChartNonRetryable(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	var calls int
	patches.ApplyMethodFunc(&tkesvc.TkeService{}, "DescribeAddonValuesList", func(_ context.Context, _ string, addonName string) (*tkev20180525.DescribeAddonValuesResponseParams, error) {
		calls++
		return nil, sdkErrors.NewTencentCloudSDKError("ResourceUnavailable", "not found chart", "req-1")
	})

	meta := newMockMetaForTkeAddonValuesDS()
	res := tkesvc.DataSourceTencentCloudTkeAddonValues()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"cluster_id": "cls-dhhfi9n8",
		"addon_name": "no-such-addon",
	})

	err := res.Read(d, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found chart")
	assert.Equal(t, 1, calls, "deterministic error must not be retried")
	assert.Equal(t, "", d.Id())
}

// TestTkeAddonValuesDataSource_ReadEmptyResourceNameNonRetryable covers the other
// deterministic message: an empty addon name is rejected by the cloud API with
// `resource name may not be empty`, which must also fail fast.
func TestTkeAddonValuesDataSource_ReadEmptyResourceNameNonRetryable(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	var calls int
	patches.ApplyMethodFunc(&tkesvc.TkeService{}, "DescribeAddonValuesList", func(_ context.Context, _ string, _ string) (*tkev20180525.DescribeAddonValuesResponseParams, error) {
		calls++
		return nil, sdkErrors.NewTencentCloudSDKError("ResourceUnavailable", "resource name may not be empty", "req-2")
	})

	meta := newMockMetaForTkeAddonValuesDS()
	res := tkesvc.DataSourceTencentCloudTkeAddonValues()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"cluster_id": "cls-dhhfi9n8",
		"addon_name": "x",
	})

	err := res.Read(d, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "resource name may not be empty")
	assert.Equal(t, 1, calls, "deterministic error must not be retried")
}

// TestTkeAddonValuesDataSource_ReadRetryableErrorStillRetried makes sure the
// deterministic-error special case does not swallow genuinely retryable errors
// (e.g. throttling): they must still be retried rather than failing immediately.
func TestTkeAddonValuesDataSource_ReadRetryableErrorStillRetried(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	var calls int
	patches.ApplyMethodFunc(&tkesvc.TkeService{}, "DescribeAddonValuesList", func(_ context.Context, _ string, _ string) (*tkev20180525.DescribeAddonValuesResponseParams, error) {
		calls++
		if calls < 3 {
			return nil, sdkErrors.NewTencentCloudSDKError("RequestLimitExceeded", "too many requests", "req-3")
		}
		return &tkev20180525.DescribeAddonValuesResponseParams{
			Values:        ptrStrTkeAddonValuesDS(`{"a":1}`),
			DefaultValues: ptrStrTkeAddonValuesDS(`{"a":1}`),
		}, nil
	})

	meta := newMockMetaForTkeAddonValuesDS()
	res := tkesvc.DataSourceTencentCloudTkeAddonValues()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"cluster_id": "cls-dhhfi9n8",
		"addon_name": "ip-masq-agent",
	})

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, calls, 3, "retryable errors must still be retried")
	assert.Equal(t, `{"a":1}`, d.Get("values").(string))
	assert.Equal(t, "cls-dhhfi9n8#ip-masq-agent", d.Id())
}

// TestTkeAddonValuesDataSource_ResultOutputFile verifies that `result_output_file`
// contains the real payload. Passing ResourceData to WriteToFile used to produce `{}`.
func TestTkeAddonValuesDataSource_ResultOutputFile(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	patches.ApplyMethodFunc(&tkesvc.TkeService{}, "DescribeAddonValuesList", func(_ context.Context, _ string, _ string) (*tkev20180525.DescribeAddonValuesResponseParams, error) {
		return &tkev20180525.DescribeAddonValuesResponseParams{
			Values:        ptrStrTkeAddonValuesDS(`{"limits":{"memory":"500Mi"}}`),
			DefaultValues: ptrStrTkeAddonValuesDS(`{"limits":{"memory":"500Mi"}}`),
		}, nil
	})

	outFile := filepath.Join(t.TempDir(), "addon-values.json")
	meta := newMockMetaForTkeAddonValuesDS()
	res := tkesvc.DataSourceTencentCloudTkeAddonValues()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"cluster_id":         "cls-dhhfi9n8",
		"addon_name":         "rdma-agent",
		"result_output_file": outFile,
	})

	err := res.Read(d, meta)
	assert.NoError(t, err)

	raw, readErr := os.ReadFile(outFile)
	assert.NoError(t, readErr)
	assert.NotEqual(t, "{}", strings.TrimSpace(string(raw)), "result_output_file must not be an empty JSON object")

	var payload map[string]string
	assert.NoError(t, json.Unmarshal(raw, &payload))
	assert.Equal(t, `{"limits":{"memory":"500Mi"}}`, payload["values"])
	assert.Equal(t, `{"limits":{"memory":"500Mi"}}`, payload["default_values"])
}

// TestTkeAddonValuesDataSource_EmptyArgumentsRejected verifies that empty
// `cluster_id` / `addon_name` are rejected at plan time, so the request never
// reaches the cloud API (where it would be retried for 3 minutes).
func TestTkeAddonValuesDataSource_EmptyArgumentsRejected(t *testing.T) {
	res := tkesvc.DataSourceTencentCloudTkeAddonValues()

	for _, field := range []string{"cluster_id", "addon_name"} {
		s := res.Schema[field]
		assert.NotNil(t, s, "schema field %s is missing", field)
		assert.NotNil(t, s.ValidateFunc, "schema field %s must have a ValidateFunc", field)

		_, errs := s.ValidateFunc("", field)
		assert.NotEmpty(t, errs, "empty %s must be rejected", field)

		_, errs = s.ValidateFunc("cls-dhhfi9n8", field)
		assert.Empty(t, errs, "non-empty %s must be accepted", field)
	}
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
