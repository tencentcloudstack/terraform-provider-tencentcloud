package tcmg_test

import (
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	monitor "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/monitor/v20180724"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/tcmg"
)

type mockMetaForMonitorGrafanaVersionsDS struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMetaForMonitorGrafanaVersionsDS) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockMetaForMonitorGrafanaVersionsDS{}

func newMockMetaForMonitorGrafanaVersionsDS() *mockMetaForMonitorGrafanaVersionsDS {
	return &mockMetaForMonitorGrafanaVersionsDS{client: &connectivity.TencentCloudClient{}}
}

func ptrStringGVDS(s string) *string { return &s }

func buildGrafanaVersion(alias, version string) *monitor.GrafanaVersion {
	return &monitor.GrafanaVersion{
		Alias:   ptrStringGVDS(alias),
		Version: ptrStringGVDS(version),
	}
}

func TestMonitorGrafanaVersionsDS_ReadBasic(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	monitorClient := &monitor.Client{}
	patches.ApplyMethodReturn(newMockMetaForMonitorGrafanaVersionsDS().client, "UseMonitorClient", monitorClient)

	patches.ApplyMethodFunc(monitorClient, "DescribeGrafanaVersions", func(request *monitor.DescribeGrafanaVersionsRequest) (*monitor.DescribeGrafanaVersionsResponse, error) {
		resp := monitor.NewDescribeGrafanaVersionsResponse()
		resp.Response = &monitor.DescribeGrafanaVersionsResponseParams{
			Versions: []*monitor.GrafanaVersion{
				buildGrafanaVersion("stable", "10.2.5"),
				buildGrafanaVersion("beta", "10.3.0"),
			},
			RequestId: ptrStringGVDS("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaForMonitorGrafanaVersionsDS()
	res := tcmg.DataSourceTencentCloudMonitorGrafanaVersions()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.NotEmpty(t, d.Id())

	versions := d.Get("versions").([]interface{})
	assert.Len(t, versions, 2)

	version0 := versions[0].(map[string]interface{})
	assert.Equal(t, "stable", version0["alias"].(string))
	assert.Equal(t, "10.2.5", version0["version"].(string))

	version1 := versions[1].(map[string]interface{})
	assert.Equal(t, "beta", version1["alias"].(string))
	assert.Equal(t, "10.3.0", version1["version"].(string))
}

func TestMonitorGrafanaVersionsDS_ReadEmpty(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	monitorClient := &monitor.Client{}
	patches.ApplyMethodReturn(newMockMetaForMonitorGrafanaVersionsDS().client, "UseMonitorClient", monitorClient)

	patches.ApplyMethodFunc(monitorClient, "DescribeGrafanaVersions", func(request *monitor.DescribeGrafanaVersionsRequest) (*monitor.DescribeGrafanaVersionsResponse, error) {
		resp := monitor.NewDescribeGrafanaVersionsResponse()
		resp.Response = &monitor.DescribeGrafanaVersionsResponseParams{
			Versions:  []*monitor.GrafanaVersion{},
			RequestId: ptrStringGVDS("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaForMonitorGrafanaVersionsDS()
	res := tcmg.DataSourceTencentCloudMonitorGrafanaVersions()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})

	err := res.Read(d, meta)
	assert.Error(t, err)
}

func TestMonitorGrafanaVersionsDS_Schema(t *testing.T) {
	res := tcmg.DataSourceTencentCloudMonitorGrafanaVersions()

	assert.NotNil(t, res)
	assert.Contains(t, res.Schema, "versions")
	assert.Contains(t, res.Schema, "result_output_file")

	versionsSchema := res.Schema["versions"]
	assert.Equal(t, schema.TypeList, versionsSchema.Type)
	assert.True(t, versionsSchema.Computed)

	elemRes := versionsSchema.Elem.(*schema.Resource)
	assert.Contains(t, elemRes.Schema, "alias")
	assert.Contains(t, elemRes.Schema, "version")

	outputSchema := res.Schema["result_output_file"]
	assert.Equal(t, schema.TypeString, outputSchema.Type)
	assert.True(t, outputSchema.Optional)
}
