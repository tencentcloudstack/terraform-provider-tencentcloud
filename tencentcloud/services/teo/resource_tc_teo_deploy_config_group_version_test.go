package teo_test

import (
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	teov20220901 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/teo"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

// mockMetaForDeployConfigGroupVersion implements tccommon.ProviderMeta
type mockMetaForDeployConfigGroupVersion struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMetaForDeployConfigGroupVersion) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockMetaForDeployConfigGroupVersion{}

func newMockMetaForDeployConfigGroupVersion() *mockMetaForDeployConfigGroupVersion {
	return &mockMetaForDeployConfigGroupVersion{client: &connectivity.TencentCloudClient{}}
}

func ptrStrDeployConfigGroupVersion(s string) *string {
	return &s
}

func ptrUint64DeployConfigGroupVersion(u uint64) *uint64 {
	return &u
}

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

// mockDeployHistoryResponseForEnvInfo builds a DescribeDeployHistory response that contains
// one deploy record, so the existing deploy-record read path succeeds and reaches the new
// environment-info read logic.
func mockDeployHistoryResponseForEnvInfo() *teov20220901.DescribeDeployHistoryResponse {
	resp := teov20220901.NewDescribeDeployHistoryResponse()
	resp.Response = &teov20220901.DescribeDeployHistoryResponseParams{
		TotalCount: ptrUint64DeployConfigGroupVersion(1),
		Records: []*teov20220901.DeployRecord{
			{
				RecordId:   ptrStrDeployConfigGroupVersion("rec-001"),
				DeployTime: ptrStrDeployConfigGroupVersion("2026-08-11T00:00:00Z"),
				Status:     ptrStrDeployConfigGroupVersion("success"),
				Message:    ptrStrDeployConfigGroupVersion("deploy ok"),
			},
		},
		RequestId: ptrStrDeployConfigGroupVersion("fake-request-id-history"),
	}
	return resp
}

// mockEnvironmentsResponseForEnvInfo builds a DescribeEnvironments response containing the
// target env id with all env-related fields populated.
func mockEnvironmentsResponseForEnvInfo(envId string) *teov20220901.DescribeEnvironmentsResponse {
	resp := teov20220901.NewDescribeEnvironmentsResponse()
	resp.Response = &teov20220901.DescribeEnvironmentsResponseParams{
		TotalCount: ptrUint64DeployConfigGroupVersion(2),
		EnvInfos: []*teov20220901.EnvInfo{
			{
				EnvId:   ptrStrDeployConfigGroupVersion("env-other"),
				EnvType: ptrStrDeployConfigGroupVersion("staging"),
				Status:  ptrStrDeployConfigGroupVersion("running"),
				Scope:   []*string{ptrStrDeployConfigGroupVersion("1.2.3.4")},
			},
			{
				EnvId:      ptrStrDeployConfigGroupVersion(envId),
				EnvType:    ptrStrDeployConfigGroupVersion("production"),
				Status:     ptrStrDeployConfigGroupVersion("running"),
				Scope:      []*string{ptrStrDeployConfigGroupVersion("ALL")},
				CreateTime: ptrStrDeployConfigGroupVersion("2026-08-10T00:00:00Z"),
				UpdateTime: ptrStrDeployConfigGroupVersion("2026-08-11T00:00:00Z"),
				CurrentConfigGroupVersionInfos: []*teov20220901.ConfigGroupVersionInfo{
					{
						VersionId:     ptrStrDeployConfigGroupVersion("ver-aaa"),
						VersionNumber: ptrStrDeployConfigGroupVersion("3"),
						SourceVersion: ptrStrDeployConfigGroupVersion("ver-bbb"),
						GroupType:     ptrStrDeployConfigGroupVersion("l7_acceleration"),
						GroupId:       ptrStrDeployConfigGroupVersion("cg-aaa"),
						Description:   ptrStrDeployConfigGroupVersion("env effective version"),
						Status:        ptrStrDeployConfigGroupVersion("active"),
						CreateTime:    ptrStrDeployConfigGroupVersion("2026-08-09T00:00:00Z"),
					},
				},
			},
		},
		RequestId: ptrStrDeployConfigGroupVersion("fake-request-id-env"),
	}
	return resp
}

// TestDeployConfigGroupVersionEnvInfo_Read_Success tests Read populates all the new env-related
// computed fields when the target env is present in DescribeEnvironments.
func TestDeployConfigGroupVersionEnvInfo_Read_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaForDeployConfigGroupVersion().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeDeployHistory", func(request *teov20220901.DescribeDeployHistoryRequest) (*teov20220901.DescribeDeployHistoryResponse, error) {
		return mockDeployHistoryResponseForEnvInfo(), nil
	})

	patches.ApplyMethodFunc(teoClient, "DescribeEnvironments", func(request *teov20220901.DescribeEnvironmentsRequest) (*teov20220901.DescribeEnvironmentsResponse, error) {
		return mockEnvironmentsResponseForEnvInfo("env-target"), nil
	})

	meta := newMockMetaForDeployConfigGroupVersion()
	res := teo.ResourceTencentCloudTeoDeployConfigGroupVersion()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":     "zone-test123",
		"env_id":      "env-target",
		"description": "deploy desc",
	})
	d.SetId("zone-test123#env-target#rec-001")

	err := res.Read(d, meta)
	assert.NoError(t, err)

	// Existing deploy-record fields
	assert.Equal(t, "rec-001", d.Get("record_id"))
	assert.Equal(t, "success", d.Get("status"))

	// New env-related computed fields
	assert.Equal(t, 2, d.Get("total_count"))
	assert.Equal(t, "production", d.Get("env_type"))
	assert.Equal(t, "ALL", d.Get("scope").([]interface{})[0])

	// env create/update time
	assert.Equal(t, "2026-08-10T00:00:00Z", d.Get("env_create_time"))
	assert.Equal(t, "2026-08-11T00:00:00Z", d.Get("env_update_time"))

	// current_config_group_version_infos set
	currentSet := d.Get("current_config_group_version_infos").(*schema.Set)
	assert.Equal(t, 1, currentSet.Len())
	list := currentSet.List()
	item := list[0].(map[string]interface{})
	assert.Equal(t, "ver-aaa", item["version_id"])
	assert.Equal(t, "3", item["version_number"])
	assert.Equal(t, "ver-bbb", item["source_version"])
	assert.Equal(t, "l7_acceleration", item["group_type"])
	assert.Equal(t, "cg-aaa", item["group_id"])
	assert.Equal(t, "env effective version", item["description"])
	assert.Equal(t, "active", item["status"])
	assert.Equal(t, "2026-08-09T00:00:00Z", item["create_time"])

	// ID must NOT be cleared
	assert.Equal(t, "zone-test123#env-target#rec-001", d.Id())
}

// TestDeployConfigGroupVersionEnvInfo_Read_EnvNotFound tests Read when DescribeEnvironments
// does not contain the target env id: no error, ID not cleared, new fields not set.
func TestDeployConfigGroupVersionEnvInfo_Read_EnvNotFound(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaForDeployConfigGroupVersion().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeDeployHistory", func(request *teov20220901.DescribeDeployHistoryRequest) (*teov20220901.DescribeDeployHistoryResponse, error) {
		return mockDeployHistoryResponseForEnvInfo(), nil
	})

	patches.ApplyMethodFunc(teoClient, "DescribeEnvironments", func(request *teov20220901.DescribeEnvironmentsRequest) (*teov20220901.DescribeEnvironmentsResponse, error) {
		resp := teov20220901.NewDescribeEnvironmentsResponse()
		resp.Response = &teov20220901.DescribeEnvironmentsResponseParams{
			TotalCount: ptrUint64DeployConfigGroupVersion(1),
			EnvInfos: []*teov20220901.EnvInfo{
				{
					EnvId:   ptrStrDeployConfigGroupVersion("env-other"),
					EnvType: ptrStrDeployConfigGroupVersion("staging"),
					Status:  ptrStrDeployConfigGroupVersion("running"),
				},
			},
			RequestId: ptrStrDeployConfigGroupVersion("fake-request-id-env"),
		}
		return resp, nil
	})

	meta := newMockMetaForDeployConfigGroupVersion()
	res := teo.ResourceTencentCloudTeoDeployConfigGroupVersion()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":     "zone-test123",
		"env_id":      "env-target",
		"description": "deploy desc",
	})
	d.SetId("zone-test123#env-target#rec-001")

	err := res.Read(d, meta)
	assert.NoError(t, err)

	// Existing deploy-record fields still populated
	assert.Equal(t, "rec-001", d.Get("record_id"))

	// New env fields remain unset (zero values)
	assert.Equal(t, 0, d.Get("total_count"))
	assert.Equal(t, "", d.Get("env_type"))
	assert.Equal(t, 0, len(d.Get("scope").([]interface{})))
	assert.Equal(t, 0, d.Get("current_config_group_version_infos").(*schema.Set).Len())

	// ID must NOT be cleared
	assert.Equal(t, "zone-test123#env-target#rec-001", d.Id())
}

// TestDeployConfigGroupVersionEnvInfo_Read_PartialNilFields tests Read when the target EnvInfo
// has some nil fields: nil fields are skipped (not set) without error.
func TestDeployConfigGroupVersionEnvInfo_Read_PartialNilFields(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaForDeployConfigGroupVersion().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeDeployHistory", func(request *teov20220901.DescribeDeployHistoryRequest) (*teov20220901.DescribeDeployHistoryResponse, error) {
		return mockDeployHistoryResponseForEnvInfo(), nil
	})

	patches.ApplyMethodFunc(teoClient, "DescribeEnvironments", func(request *teov20220901.DescribeEnvironmentsRequest) (*teov20220901.DescribeEnvironmentsResponse, error) {
		resp := teov20220901.NewDescribeEnvironmentsResponse()
		resp.Response = &teov20220901.DescribeEnvironmentsResponseParams{
			TotalCount: ptrUint64DeployConfigGroupVersion(1),
			EnvInfos: []*teov20220901.EnvInfo{
				{
					EnvId:      ptrStrDeployConfigGroupVersion("env-target"),
					EnvType:    ptrStrDeployConfigGroupVersion("production"),
					Status:     ptrStrDeployConfigGroupVersion("running"),
					Scope:      nil,
					CreateTime: nil,
					UpdateTime: nil,
					CurrentConfigGroupVersionInfos: []*teov20220901.ConfigGroupVersionInfo{
						{
							VersionId:     ptrStrDeployConfigGroupVersion("ver-aaa"),
							VersionNumber: nil,
							SourceVersion: nil,
							GroupType:     nil,
							GroupId:       ptrStrDeployConfigGroupVersion("cg-aaa"),
							Description:   nil,
							Status:        ptrStrDeployConfigGroupVersion("active"),
							CreateTime:    nil,
						},
					},
				},
			},
			RequestId: ptrStrDeployConfigGroupVersion("fake-request-id-env"),
		}
		return resp, nil
	})

	meta := newMockMetaForDeployConfigGroupVersion()
	res := teo.ResourceTencentCloudTeoDeployConfigGroupVersion()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":     "zone-test123",
		"env_id":      "env-target",
		"description": "deploy desc",
	})
	d.SetId("zone-test123#env-target#rec-001")

	err := res.Read(d, meta)
	assert.NoError(t, err)

	// env_type set, but env_create_time/env_update_time/scope skipped (nil)
	assert.Equal(t, "production", d.Get("env_type"))
	assert.Equal(t, "", d.Get("env_create_time"))
	assert.Equal(t, "", d.Get("env_update_time"))
	assert.Equal(t, 0, len(d.Get("scope").([]interface{})))

	// current_config_group_version_infos: only non-nil sub-fields populated
	currentSet := d.Get("current_config_group_version_infos").(*schema.Set)
	assert.Equal(t, 1, currentSet.Len())
	item := currentSet.List()[0].(map[string]interface{})
	assert.Equal(t, "ver-aaa", item["version_id"])
	assert.Equal(t, "cg-aaa", item["group_id"])
	assert.Equal(t, "active", item["status"])
	assert.Equal(t, "", item["version_number"])
	assert.Equal(t, "", item["source_version"])
	assert.Equal(t, "", item["group_type"])
	assert.Equal(t, "", item["description"])
	assert.Equal(t, "", item["create_time"])

	// ID must NOT be cleared
	assert.Equal(t, "zone-test123#env-target#rec-001", d.Id())
}

// TestDeployConfigGroupVersionEnvInfo_Schema verifies the new computed fields exist and are
// Computed-only (no Required/Optional/ForceNew).
func TestDeployConfigGroupVersionEnvInfo_Schema(t *testing.T) {
	res := teo.ResourceTencentCloudTeoDeployConfigGroupVersion()
	assert.NotNil(t, res)

	checks := map[string]schema.ValueType{
		"total_count":     schema.TypeInt,
		"env_type":        schema.TypeString,
		"scope":           schema.TypeList,
		"env_create_time": schema.TypeString,
		"env_update_time": schema.TypeString,
	}
	for name, typ := range checks {
		s, ok := res.Schema[name]
		assert.True(t, ok, "field %s should exist", name)
		assert.Equal(t, typ, s.Type, "field %s type", name)
		assert.True(t, s.Computed, "field %s should be Computed", name)
		assert.False(t, s.Required, "field %s should not be Required", name)
		assert.False(t, s.Optional, "field %s should not be Optional", name)
		assert.False(t, s.ForceNew, "field %s should not be ForceNew", name)
	}

	// current_config_group_version_infos is a TypeSet of schema.Resource
	cgv, ok := res.Schema["current_config_group_version_infos"]
	assert.True(t, ok)
	assert.Equal(t, schema.TypeSet, cgv.Type)
	assert.True(t, cgv.Computed)
	assert.False(t, cgv.Required)
	assert.False(t, cgv.Optional)
	assert.False(t, cgv.ForceNew)

	r, ok := cgv.Elem.(*schema.Resource)
	assert.True(t, ok)
	subFields := []string{
		"version_id", "version_number", "source_version", "group_type",
		"group_id", "description", "status", "create_time",
	}
	for _, f := range subFields {
		s, ok := r.Schema[f]
		assert.True(t, ok, "sub field %s should exist", f)
		assert.Equal(t, schema.TypeString, s.Type)
		assert.True(t, s.Computed)
		assert.False(t, s.Required)
		assert.False(t, s.Optional)
		assert.False(t, s.ForceNew)
	}
}
