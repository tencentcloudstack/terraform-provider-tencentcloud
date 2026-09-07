package teo_test

import (
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

func TestAccTencentCloudTeoDeployConfigGroupVersionResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTeoDeployConfigGroupVersion,
				Check: resource.ComposeTestCheckFunc(
					// Basic resource attributes
					resource.TestCheckResourceAttrSet("tencentcloud_teo_deploy_config_group_version.teo_deploy_config_group_version", "id"),

					// Input parameters validation
					resource.TestCheckResourceAttr("tencentcloud_teo_deploy_config_group_version.teo_deploy_config_group_version", "zone_id", "zone-2xkazzl8yf6k"),
					resource.TestCheckResourceAttr("tencentcloud_teo_deploy_config_group_version.teo_deploy_config_group_version", "env_id", "env-3lchxiq1h855"),
					resource.TestCheckResourceAttr("tencentcloud_teo_deploy_config_group_version.teo_deploy_config_group_version", "description", "Deploy config group version for production"),

					// Config group version infos validation
					resource.TestCheckResourceAttr("tencentcloud_teo_deploy_config_group_version.teo_deploy_config_group_version", "config_group_version_infos.#", "2"),
					resource.TestCheckResourceAttr("tencentcloud_teo_deploy_config_group_version.teo_deploy_config_group_version", "config_group_version_infos.0.version_id", "ver-3lchxizh2mqn"),
					resource.TestCheckResourceAttr("tencentcloud_teo_deploy_config_group_version.teo_deploy_config_group_version", "config_group_version_infos.1.version_id", "ver-3lchxjdciuzx"),

					// Computed attributes validation
					resource.TestCheckResourceAttrSet("tencentcloud_teo_deploy_config_group_version.teo_deploy_config_group_version", "record_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_teo_deploy_config_group_version.teo_deploy_config_group_version", "deploy_time"),
					resource.TestCheckResourceAttrSet("tencentcloud_teo_deploy_config_group_version.teo_deploy_config_group_version", "status"),

					// Optional computed attributes (may or may not be present)
					resource.TestCheckResourceAttr("tencentcloud_teo_deploy_config_group_version.teo_deploy_config_group_version", "status", "success"),
				),
			},
		},
	})
}

const testAccTeoDeployConfigGroupVersion = `

resource "tencentcloud_teo_deploy_config_group_version" "teo_deploy_config_group_version" {
  zone_id = "zone-2xkazzl8yf6k"
  env_id = "env-3lchxiq1h855"
  description = "Deploy config group version for production"
  # l7_acceleration
  config_group_version_infos {
    version_id = "ver-3lchxizh2mqn"
  }
  # edge_functions
  config_group_version_infos {
    version_id = "ver-3lchxjdciuzx"
  }
}
`

// mockDeployConfigGroupVersionMeta implements tccommon.ProviderMeta for unit tests of deploy_config_group_version.
type mockDeployConfigGroupVersionMeta struct {
	client *connectivity.TencentCloudClient
}

func (m *mockDeployConfigGroupVersionMeta) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockDeployConfigGroupVersionMeta{}

func newMockDeployConfigGroupVersionMeta() *mockDeployConfigGroupVersionMeta {
	return &mockDeployConfigGroupVersionMeta{client: &connectivity.TencentCloudClient{}}
}

func ptrStringDeployCfgGrpVer(s string) *string {
	return &s
}

// TestTeoDeployConfigGroupVersion_ReadSourceVersionPopulated tests Read fills source_version from DescribeDeployHistory when SourceVersion is non-nil.
//
// go test ./tencentcloud/services/teo/ -run "TestTeoDeployConfigGroupVersion_ReadSourceVersionPopulated" -v -count=1 -gcflags="all=-l"
func TestTeoDeployConfigGroupVersion_ReadSourceVersionPopulated(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockDeployConfigGroupVersionMeta().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeDeployHistory", func(request *teov20220901.DescribeDeployHistoryRequest) (*teov20220901.DescribeDeployHistoryResponse, error) {
		resp := teov20220901.NewDescribeDeployHistoryResponse()
		resp.Response = &teov20220901.DescribeDeployHistoryResponseParams{
			TotalCount: ptrUint64DeployCfgGrpVer(1),
			Records: []*teov20220901.DeployRecord{
				{
					RecordId:   ptrStringDeployCfgGrpVer("rec-12345678"),
					DeployTime: ptrStringDeployCfgGrpVer("2026-08-11T10:00:00Z"),
					Status:     ptrStringDeployCfgGrpVer("success"),
					Message:    ptrStringDeployCfgGrpVer("deploy ok"),
					ConfigGroupVersionInfos: []*teov20220901.ConfigGroupVersionInfo{
						{
							VersionId:     ptrStringDeployCfgGrpVer("ver-aaaaaaaaaaaa"),
							VersionNumber: ptrStringDeployCfgGrpVer("1"),
							GroupId:       ptrStringDeployCfgGrpVer("cg-bbbbbbbbbbbb"),
							GroupType:     ptrStringDeployCfgGrpVer("l7_acceleration"),
							Description:   ptrStringDeployCfgGrpVer("version derived from source"),
							Status:        ptrStringDeployCfgGrpVer("active"),
							CreateTime:    ptrStringDeployCfgGrpVer("2026-08-10T09:00:00Z"),
							SourceVersion: ptrStringDeployCfgGrpVer("ver-cccccccccccc"),
						},
					},
				},
			},
			RequestId: ptrStringDeployCfgGrpVer("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockDeployConfigGroupVersionMeta()
	res := teo.ResourceTencentCloudTeoDeployConfigGroupVersion()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-2xkazzl8yf6k",
		"env_id":  "env-3lchxiq1h855",
	})
	d.SetId("zone-2xkazzl8yf6k#env-3lchxiq1h855#rec-12345678")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "zone-2xkazzl8yf6k#env-3lchxiq1h855#rec-12345678", d.Id())

	infos := d.Get("config_group_version_infos").(*schema.Set).List()
	assert.Equal(t, 1, len(infos))
	infoMap := infos[0].(map[string]interface{})
	assert.Equal(t, "ver-aaaaaaaaaaaa", infoMap["version_id"])
	assert.Equal(t, "ver-cccccccccccc", infoMap["source_version"])
}

// TestTeoDeployConfigGroupVersion_ReadSourceVersionNil tests Read skips source_version when SourceVersion is nil and does not error.
//
// go test ./tencentcloud/services/teo/ -run "TestTeoDeployConfigGroupVersion_ReadSourceVersionNil" -v -count=1 -gcflags="all=-l"
func TestTeoDeployConfigGroupVersion_ReadSourceVersionNil(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockDeployConfigGroupVersionMeta().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeDeployHistory", func(request *teov20220901.DescribeDeployHistoryRequest) (*teov20220901.DescribeDeployHistoryResponse, error) {
		resp := teov20220901.NewDescribeDeployHistoryResponse()
		resp.Response = &teov20220901.DescribeDeployHistoryResponseParams{
			TotalCount: ptrUint64DeployCfgGrpVer(1),
			Records: []*teov20220901.DeployRecord{
				{
					RecordId:   ptrStringDeployCfgGrpVer("rec-87654321"),
					DeployTime: ptrStringDeployCfgGrpVer("2026-08-11T11:00:00Z"),
					Status:     ptrStringDeployCfgGrpVer("success"),
					Message:    ptrStringDeployCfgGrpVer("deploy ok"),
					ConfigGroupVersionInfos: []*teov20220901.ConfigGroupVersionInfo{
						{
							VersionId:     ptrStringDeployCfgGrpVer("ver-dddddddddddd"),
							VersionNumber: ptrStringDeployCfgGrpVer("0"),
							GroupId:       ptrStringDeployCfgGrpVer("cg-eeeeeeeeeeee"),
							GroupType:     ptrStringDeployCfgGrpVer("edge_functions"),
							Description:   ptrStringDeployCfgGrpVer("initial version without source"),
							Status:        ptrStringDeployCfgGrpVer("active"),
							CreateTime:    ptrStringDeployCfgGrpVer("2026-08-10T08:00:00Z"),
							SourceVersion: nil,
						},
					},
				},
			},
			RequestId: ptrStringDeployCfgGrpVer("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockDeployConfigGroupVersionMeta()
	res := teo.ResourceTencentCloudTeoDeployConfigGroupVersion()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-2xkazzl8yf6k",
		"env_id":  "env-3lchxiq1h855",
	})
	d.SetId("zone-2xkazzl8yf6k#env-3lchxiq1h855#rec-87654321")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "zone-2xkazzl8yf6k#env-3lchxiq1h855#rec-87654321", d.Id())

	infos := d.Get("config_group_version_infos").(*schema.Set).List()
	assert.Equal(t, 1, len(infos))
	infoMap := infos[0].(map[string]interface{})
	assert.Equal(t, "ver-dddddddddddd", infoMap["version_id"])
	// source_version should be empty string since SourceVersion is nil and skipped.
	assert.Equal(t, "", infoMap["source_version"])
}

// TestTeoDeployConfigGroupVersion_ReadNotFound tests Read clears id when no deploy record returned.
//
// go test ./tencentcloud/services/teo/ -run "TestTeoDeployConfigGroupVersion_ReadNotFound" -v -count=1 -gcflags="all=-l"
func TestTeoDeployConfigGroupVersion_ReadNotFound(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockDeployConfigGroupVersionMeta().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeDeployHistory", func(request *teov20220901.DescribeDeployHistoryRequest) (*teov20220901.DescribeDeployHistoryResponse, error) {
		resp := teov20220901.NewDescribeDeployHistoryResponse()
		resp.Response = &teov20220901.DescribeDeployHistoryResponseParams{
			TotalCount: ptrUint64DeployCfgGrpVer(0),
			Records:    []*teov20220901.DeployRecord{},
			RequestId:  ptrStringDeployCfgGrpVer("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockDeployConfigGroupVersionMeta()
	res := teo.ResourceTencentCloudTeoDeployConfigGroupVersion()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-2xkazzl8yf6k",
		"env_id":  "env-3lchxiq1h855",
	})
	d.SetId("zone-2xkazzl8yf6k#env-3lchxiq1h855#rec-not-exist")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "", d.Id())
}

// TestTeoDeployConfigGroupVersion_Schema validates schema definition for source_version.
//
// go test ./tencentcloud/services/teo/ -run "TestTeoDeployConfigGroupVersion_Schema" -v -count=1 -gcflags="all=-l"
func TestTeoDeployConfigGroupVersion_Schema(t *testing.T) {
	res := teo.ResourceTencentCloudTeoDeployConfigGroupVersion()
	assert.NotNil(t, res)

	infosSchema, ok := res.Schema["config_group_version_infos"]
	assert.True(t, ok)
	assert.NotNil(t, infosSchema.Elem)

	elemRes, ok := infosSchema.Elem.(*schema.Resource)
	assert.True(t, ok)
	assert.NotNil(t, elemRes.Schema)

	sourceVersionSchema, ok := elemRes.Schema["source_version"]
	assert.True(t, ok)
	assert.Equal(t, schema.TypeString, sourceVersionSchema.Type)
	assert.True(t, sourceVersionSchema.Computed)
	assert.False(t, sourceVersionSchema.Required)
	assert.False(t, sourceVersionSchema.Optional)
	assert.False(t, sourceVersionSchema.ForceNew)
}

func ptrUint64DeployCfgGrpVer(n uint64) *uint64 {
	return &n
}
