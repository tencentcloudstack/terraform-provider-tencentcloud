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

func TestAccTencentCloudTeoZonesDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTeoZonesDataSource,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_teo_zones.teo_zones"),
				),
			},
		},
	})
}

const testAccTeoZonesDataSource = `

data "tencentcloud_teo_zones" "teo_zones" {
  filters {
    name = "tag-value"
    values = ["terraform"]
  }
}
`

// go test ./tencentcloud/services/teo/ -run "TestTeoZonesDataSource_workModeInfos" -v -count=1 -gcflags="all=-l"

type mockMetaTeoZones struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMetaTeoZones) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockMetaTeoZones{}

func newMockMetaTeoZones() *mockMetaTeoZones {
	return &mockMetaTeoZones{client: &connectivity.TencentCloudClient{}}
}

func ptrStringTeoZones(s string) *string {
	return &s
}

// TestTeoZonesDataSource_workModeInfos tests that WorkModeInfos is flattened into work_mode_infos
func TestTeoZonesDataSource_workModeInfos(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaTeoZones().client, "UseTeoClient", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeZones", func(request *teov20220901.DescribeZonesRequest) (*teov20220901.DescribeZonesResponse, error) {
		resp := teov20220901.NewDescribeZonesResponse()
		resp.Response = &teov20220901.DescribeZonesResponseParams{
			Zones: []*teov20220901.Zone{
				{
					ZoneId:   ptrStringTeoZones("zone-2noz78a8ev6k"),
					ZoneName: ptrStringTeoZones("example.com"),
					Status:   ptrStringTeoZones("active"),
					Type:     ptrStringTeoZones("full"),
					WorkModeInfos: []*teov20220901.ConfigGroupWorkModeInfo{
						{
							ConfigGroupType: ptrStringTeoZones("l7_acceleration"),
							WorkMode:        ptrStringTeoZones("immediate_effect"),
						},
						{
							ConfigGroupType: ptrStringTeoZones("web_security"),
							WorkMode:        ptrStringTeoZones("version_control"),
						},
					},
				},
			},
			RequestId: ptrStringTeoZones("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaTeoZones()
	res := teo.DataSourceTencentCloudTeoZones()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.NotEmpty(t, d.Id())

	zones := d.Get("zones").([]interface{})
	assert.Len(t, zones, 1)

	zoneMap := zones[0].(map[string]interface{})
	assert.Equal(t, "zone-2noz78a8ev6k", zoneMap["zone_id"])
	assert.Equal(t, "example.com", zoneMap["zone_name"])

	workModeInfos := zoneMap["work_mode_infos"].([]interface{})
	assert.Len(t, workModeInfos, 2)

	firstInfo := workModeInfos[0].(map[string]interface{})
	assert.Equal(t, "l7_acceleration", firstInfo["config_group_type"])
	assert.Equal(t, "immediate_effect", firstInfo["work_mode"])

	secondInfo := workModeInfos[1].(map[string]interface{})
	assert.Equal(t, "web_security", secondInfo["config_group_type"])
	assert.Equal(t, "version_control", secondInfo["work_mode"])
}

// TestTeoZonesDataSource_workModeInfosNil tests that a nil WorkModeInfos does not error and other fields are populated
func TestTeoZonesDataSource_workModeInfosNil(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaTeoZones().client, "UseTeoClient", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeZones", func(request *teov20220901.DescribeZonesRequest) (*teov20220901.DescribeZonesResponse, error) {
		resp := teov20220901.NewDescribeZonesResponse()
		resp.Response = &teov20220901.DescribeZonesResponseParams{
			Zones: []*teov20220901.Zone{
				{
					ZoneId:        ptrStringTeoZones("zone-2noz78a8ev6k"),
					ZoneName:      ptrStringTeoZones("example.com"),
					Status:        ptrStringTeoZones("active"),
					Type:          ptrStringTeoZones("full"),
					WorkModeInfos: nil,
				},
			},
			RequestId: ptrStringTeoZones("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaTeoZones()
	res := teo.DataSourceTencentCloudTeoZones()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.NotEmpty(t, d.Id())

	zones := d.Get("zones").([]interface{})
	assert.Len(t, zones, 1)

	zoneMap := zones[0].(map[string]interface{})
	assert.Equal(t, "zone-2noz78a8ev6k", zoneMap["zone_id"])
	assert.Equal(t, "example.com", zoneMap["zone_name"])
	assert.Equal(t, "active", zoneMap["status"])

	// work_mode_infos should be empty when WorkModeInfos is nil
	workModeInfos := zoneMap["work_mode_infos"].([]interface{})
	assert.Len(t, workModeInfos, 0)
}

// TestTeoZonesDataSource_workModeInfosSchema validates the schema definition of work_mode_infos
func TestTeoZonesDataSource_workModeInfosSchema(t *testing.T) {
	res := teo.DataSourceTencentCloudTeoZones()

	assert.NotNil(t, res)
	assert.NotNil(t, res.Read)

	zonesSchema := res.Schema["zones"]
	assert.Equal(t, schema.TypeList, zonesSchema.Type)
	assert.True(t, zonesSchema.Computed)

	zonesElem := zonesSchema.Elem.(*schema.Resource)
	assert.Contains(t, zonesElem.Schema, "work_mode_infos")

	workModeInfosSchema := zonesElem.Schema["work_mode_infos"]
	assert.Equal(t, schema.TypeList, workModeInfosSchema.Type)
	assert.True(t, workModeInfosSchema.Computed)

	workModeInfosElem := workModeInfosSchema.Elem.(*schema.Resource)
	assert.Contains(t, workModeInfosElem.Schema, "config_group_type")
	assert.Contains(t, workModeInfosElem.Schema, "work_mode")

	configGroupTypeSchema := workModeInfosElem.Schema["config_group_type"]
	assert.Equal(t, schema.TypeString, configGroupTypeSchema.Type)
	assert.True(t, configGroupTypeSchema.Computed)

	workModeSchema := workModeInfosElem.Schema["work_mode"]
	assert.Equal(t, schema.TypeString, workModeSchema.Type)
	assert.True(t, workModeSchema.Computed)
}
