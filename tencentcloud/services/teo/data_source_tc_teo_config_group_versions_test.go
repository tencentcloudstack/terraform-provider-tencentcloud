package teo_test

import (
	"context"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	teov20220901 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/teo"
)

type mockMetaTeoConfigGroupVersionsDS struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMetaTeoConfigGroupVersionsDS) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockMetaTeoConfigGroupVersionsDS{}

func newMockMetaTeoConfigGroupVersionsDS() *mockMetaTeoConfigGroupVersionsDS {
	return &mockMetaTeoConfigGroupVersionsDS{client: &connectivity.TencentCloudClient{}}
}

func ptrStrTeoCfgGrpVerDS(s string) *string {
	return &s
}

func TestAccTencentCloudTeoConfigGroupVersionsDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{{
			Config: testAccTeoConfigGroupVersionsDataSource,
			Check: resource.ComposeTestCheckFunc(
				tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_teo_config_group_versions.teo_config_group_versions"),
				resource.TestCheckResourceAttr("data.tencentcloud_teo_config_group_versions.teo_config_group_versions", "zone_id", "zone-2xkazzl8yf6k"),
				resource.TestCheckResourceAttr("data.tencentcloud_teo_config_group_versions.teo_config_group_versions", "group_id", "cg-3lchxitnb5pb"),
				resource.TestCheckResourceAttrSet("data.tencentcloud_teo_config_group_versions.teo_config_group_versions", "config_group_version_infos.#"),
			),
		}},
	})
}

const testAccTeoConfigGroupVersionsDataSource = `

data "tencentcloud_teo_config_group_versions" "teo_config_group_versions" {
  zone_id = "zone-2xkazzl8yf6k"
  group_id = "cg-3lchxitnb5pb"
}
`

// go test ./tencentcloud/services/teo/ -run "TestTeoConfigGroupVersionsDataSource" -v -count=1 -gcflags="all=-l"

// TestTeoConfigGroupVersionsDataSource_Read_SourceVersionPopulated tests that source_version is
// populated in config_group_version_infos when the API returns a non-nil SourceVersion.
func TestTeoConfigGroupVersionsDataSource_Read_SourceVersionPopulated(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	patches.ApplyMethodFunc(&teo.TeoService{}, "DescribeTeoConfigGroupVersionsByFilter", func(_ context.Context, _ map[string]interface{}) ([]*teov20220901.ConfigGroupVersionInfo, error) {
		return []*teov20220901.ConfigGroupVersionInfo{
			{
				VersionId:     ptrStrTeoCfgGrpVerDS("ver-2kplomhisdcb"),
				VersionNumber: ptrStrTeoCfgGrpVerDS("1"),
				SourceVersion: ptrStrTeoCfgGrpVerDS("ver-aaaaaaaaaaaa"),
				GroupId:       ptrStrTeoCfgGrpVerDS("cg-3lchxitnb5pb"),
				GroupType:     ptrStrTeoCfgGrpVerDS("l7_acceleration"),
				Description:   ptrStrTeoCfgGrpVerDS("test version"),
				Status:        ptrStrTeoCfgGrpVerDS("active"),
			},
		}, nil
	})

	meta := newMockMetaTeoConfigGroupVersionsDS()
	res := teo.DataSourceTencentCloudTeoConfigGroupVersions()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":  "zone-2qtuhspy7cr6",
		"group_id": "cg-3lchxitnb5pb",
	})

	err := res.Read(d, meta)
	assert.NoError(t, err)

	configGroupVersionInfos := d.Get("config_group_version_infos").([]interface{})
	assert.Equal(t, 1, len(configGroupVersionInfos))
	info := configGroupVersionInfos[0].(map[string]interface{})
	assert.Equal(t, "ver-2kplomhisdcb", info["version_id"])
	assert.Equal(t, "ver-aaaaaaaaaaaa", info["source_version"])
}

// TestTeoConfigGroupVersionsDataSource_Read_SourceVersionNil tests that source_version is safely
// omitted when the API returns a nil SourceVersion and no error is raised.
func TestTeoConfigGroupVersionsDataSource_Read_SourceVersionNil(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	patches.ApplyMethodFunc(&teo.TeoService{}, "DescribeTeoConfigGroupVersionsByFilter", func(_ context.Context, _ map[string]interface{}) ([]*teov20220901.ConfigGroupVersionInfo, error) {
		return []*teov20220901.ConfigGroupVersionInfo{
			{
				VersionId:     ptrStrTeoCfgGrpVerDS("ver-2kplomhisdcb"),
				VersionNumber: ptrStrTeoCfgGrpVerDS("0"),
				SourceVersion: nil,
				GroupId:       ptrStrTeoCfgGrpVerDS("cg-3lchxitnb5pb"),
				GroupType:     ptrStrTeoCfgGrpVerDS("l7_acceleration"),
				Description:   ptrStrTeoCfgGrpVerDS("initial version"),
				Status:        ptrStrTeoCfgGrpVerDS("active"),
			},
		}, nil
	})

	meta := newMockMetaTeoConfigGroupVersionsDS()
	res := teo.DataSourceTencentCloudTeoConfigGroupVersions()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":  "zone-2qtuhspy7cr6",
		"group_id": "cg-3lchxitnb5pb",
	})

	err := res.Read(d, meta)
	assert.NoError(t, err)

	configGroupVersionInfos := d.Get("config_group_version_infos").([]interface{})
	assert.Equal(t, 1, len(configGroupVersionInfos))
	info := configGroupVersionInfos[0].(map[string]interface{})
	assert.Equal(t, "ver-2kplomhisdcb", info["version_id"])
	// When SourceVersion is nil, the Read loop skips setting the key, so the value defaults to empty string.
	assert.Equal(t, "", info["source_version"])
}

// TestTeoConfigGroupVersionsDataSource_Schema tests the schema definition includes source_version.
func TestTeoConfigGroupVersionsDataSource_Schema(t *testing.T) {
	res := teo.DataSourceTencentCloudTeoConfigGroupVersions()

	assert.NotNil(t, res)
	assert.NotNil(t, res.Read)

	configGroupVersionInfos := res.Schema["config_group_version_infos"]
	assert.NotNil(t, configGroupVersionInfos)
	elem := configGroupVersionInfos.Elem.(*schema.Resource)
	sourceVersionSchema, ok := elem.Schema["source_version"]
	assert.True(t, ok)
	assert.NotNil(t, sourceVersionSchema)
	assert.Equal(t, schema.TypeString, sourceVersionSchema.Type)
	assert.True(t, sourceVersionSchema.Optional)
}
