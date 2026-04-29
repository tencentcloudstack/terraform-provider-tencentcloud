package teo_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	teov20220901 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	svcteo "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/teo"

	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudTeoOriginGroup_basic -v
func TestAccTencentCloudTeoOriginGroup_basic(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_PRIVATE) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckOriginGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTeoOriginGroup,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOriginGroupExists("tencentcloud_teo_origin_group.basic"),
					resource.TestCheckResourceAttrSet("tencentcloud_teo_origin_group.basic", "zone_id"),
					resource.TestCheckResourceAttr("tencentcloud_teo_origin_group.basic", "name", "keep-group-1"),
					resource.TestCheckResourceAttr("tencentcloud_teo_origin_group.basic", "type", "GENERAL"),
					resource.TestCheckResourceAttr("tencentcloud_teo_origin_group.basic", "records.#", "3"),
					resource.TestCheckResourceAttrSet("tencentcloud_teo_origin_group.basic", "records.0.record"),
					resource.TestCheckResourceAttrSet("tencentcloud_teo_origin_group.basic", "records.0.type"),
					resource.TestCheckResourceAttrSet("tencentcloud_teo_origin_group.basic", "records.0.weight"),
					resource.TestCheckResourceAttrSet("tencentcloud_teo_origin_group.basic", "records.0.private"),
					resource.TestCheckResourceAttrSet("tencentcloud_teo_origin_group.basic", "records.1.record"),
					resource.TestCheckResourceAttrSet("tencentcloud_teo_origin_group.basic", "records.1.type"),
					resource.TestCheckResourceAttrSet("tencentcloud_teo_origin_group.basic", "records.1.weight"),
					resource.TestCheckResourceAttrSet("tencentcloud_teo_origin_group.basic", "records.1.private"),
					resource.TestCheckResourceAttrSet("tencentcloud_teo_origin_group.basic", "records.2.record"),
					resource.TestCheckResourceAttrSet("tencentcloud_teo_origin_group.basic", "records.2.type"),
					resource.TestCheckResourceAttrSet("tencentcloud_teo_origin_group.basic", "records.2.weight"),
					resource.TestCheckResourceAttrSet("tencentcloud_teo_origin_group.basic", "records.2.private"),
				),
			},
			{
				ResourceName:      "tencentcloud_teo_origin_group.basic",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccTeoOriginGroupUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOriginGroupExists("tencentcloud_teo_origin_group.basic"),
					resource.TestCheckResourceAttrSet("tencentcloud_teo_origin_group.basic", "zone_id"),
					resource.TestCheckResourceAttr("tencentcloud_teo_origin_group.basic", "records.0.private", "true"),
					resource.TestCheckResourceAttr("tencentcloud_teo_origin_group.basic", "records.0.private_parameters.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_teo_origin_group.basic", "records.0.private_parameters.0.name", "SecretAccessKey"),
					resource.TestCheckResourceAttr("tencentcloud_teo_origin_group.basic", "records.0.private_parameters.0.value", "test"),
					resource.TestCheckResourceAttrSet("tencentcloud_teo_origin_group.basic", "create_time"),
					resource.TestCheckResourceAttrSet("tencentcloud_teo_origin_group.basic", "update_time"),
				),
			},
		},
	})
}

func testAccCheckOriginGroupDestroy(s *terraform.State) error {
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	service := svcteo.NewTeoService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_teo_origin_group" {
			continue
		}
		idSplit := strings.Split(rs.Primary.ID, tccommon.FILED_SP)
		if len(idSplit) != 2 {
			return fmt.Errorf("id is broken,%s", rs.Primary.ID)
		}
		zoneId := idSplit[0]
		originGroupId := idSplit[1]

		originGroup, err := service.DescribeTeoOriginGroup(ctx, zoneId, originGroupId)
		if originGroup != nil {
			return fmt.Errorf("zone originGroup %s still exists", rs.Primary.ID)
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func testAccCheckOriginGroupExists(r string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := tccommon.GetLogId(tccommon.ContextNil)
		ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

		rs, ok := s.RootModule().Resources[r]
		if !ok {
			return fmt.Errorf("resource %s is not found", r)
		}
		idSplit := strings.Split(rs.Primary.ID, tccommon.FILED_SP)
		if len(idSplit) != 2 {
			return fmt.Errorf("id is broken,%s", rs.Primary.ID)
		}
		zoneId := idSplit[0]
		originGroupId := idSplit[1]

		service := svcteo.NewTeoService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
		originGroup, err := service.DescribeTeoOriginGroup(ctx, zoneId, originGroupId)
		if originGroup == nil {
			return fmt.Errorf("zone originGroup %s is not found", rs.Primary.ID)
		}
		if err != nil {
			return err
		}

		return nil
	}
}

const testAccTeoOriginGroup = testAccTeoZone + `

resource "tencentcloud_teo_origin_group" "basic" {
  name    = "keep-group-1"
  type    = "GENERAL"
  zone_id = tencentcloud_teo_zone.basic.id

  records {
    record  = var.zone_name
    type    = "IP_DOMAIN"
    weight  = 100
    private = false
  }
  records {
    private   = false
    record    = "21.1.1.1"
    type      = "IP_DOMAIN"
    weight    = 100
  }
  records {
    private   = false
    record    = "21.1.1.2"
    type      = "IP_DOMAIN"
    weight    = 11
  }
}

`

const testAccTeoOriginGroupUpdate = testAccTeoZone + `

resource "tencentcloud_teo_origin_group" "basic" {
  name    = "keep-group-1"
  type    = "GENERAL"
  zone_id = tencentcloud_teo_zone.basic.id

  records {
    record  = var.zone_name
    type    = "IP_DOMAIN"
    weight  = 100
    private = true

    private_parameters {
      name = "SecretAccessKey"
      value = "test"
    }
  }
}

`

// go test ./tencentcloud/services/teo/ -run "TestTeoOriginGroup" -v -count=1 -gcflags="all=-l"

// mockMetaForOriginGroup implements tccommon.ProviderMeta
type mockMetaForOriginGroup struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMetaForOriginGroup) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockMetaForOriginGroup{}

func newMockMetaForOriginGroup() *mockMetaForOriginGroup {
	return &mockMetaForOriginGroup{client: &connectivity.TencentCloudClient{}}
}

func ptrStrOriginGroup(s string) *string {
	return &s
}

func ptrUint64OriginGroup(u uint64) *uint64 {
	return &u
}

func ptrBoolOriginGroup(b bool) *bool {
	return &b
}

// TestTeoOriginGroup_Read_WithZoneFields tests Read populates zone_id, zone_name, alias_zone_name in references
func TestTeoOriginGroup_Read_WithZoneFields(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	mockClient := newMockMetaForOriginGroup().client
	patches.ApplyMethodReturn(mockClient, "UseTeoClient", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeOriginGroup", func(request *teov20220901.DescribeOriginGroupRequest) (*teov20220901.DescribeOriginGroupResponse, error) {
		resp := teov20220901.NewDescribeOriginGroupResponse()
		resp.Response = &teov20220901.DescribeOriginGroupResponseParams{
			OriginGroups: []*teov20220901.OriginGroup{
				{
					GroupId: ptrStrOriginGroup("og-test123"),
					Name:    ptrStrOriginGroup("test-origin-group"),
					Type:    ptrStrOriginGroup("GENERAL"),
					Records: []*teov20220901.OriginRecord{
						{
							Record: ptrStrOriginGroup("1.1.1.1"),
							Type:   ptrStrOriginGroup("IP_DOMAIN"),
							Weight: ptrUint64OriginGroup(100),
						},
					},
					References: []*teov20220901.OriginGroupReference{
						{
							InstanceType:  ptrStrOriginGroup("acceleration-domain"),
							InstanceId:    ptrStrOriginGroup("acc-domain-001"),
							InstanceName:  ptrStrOriginGroup("test-domain.com"),
							ZoneId:        ptrStrOriginGroup("zone-abc123"),
							ZoneName:      ptrStrOriginGroup("test-zone.com"),
							AliasZoneName: ptrStrOriginGroup("alias-test-zone.com"),
						},
						{
							InstanceType:  ptrStrOriginGroup("rule-engine"),
							InstanceId:    ptrStrOriginGroup("rule-002"),
							InstanceName:  ptrStrOriginGroup("test-rule"),
							ZoneId:        ptrStrOriginGroup("zone-def456"),
							ZoneName:      ptrStrOriginGroup("another-zone.com"),
							AliasZoneName: ptrStrOriginGroup("alias-another-zone.com"),
						},
					},
					CreateTime: ptrStrOriginGroup("2024-01-01T00:00:00Z"),
					UpdateTime: ptrStrOriginGroup("2024-01-02T00:00:00Z"),
				},
			},
			RequestId: ptrStrOriginGroup("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaForOriginGroup()
	res := svcteo.ResourceTencentCloudTeoOriginGroup()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-test123",
		"type":    "GENERAL",
	})
	d.SetId("zone-test123#og-test123")

	err := res.Read(d, meta)
	assert.NoError(t, err)

	references := d.Get("references").([]interface{})
	assert.Len(t, references, 2)

	ref0 := references[0].(map[string]interface{})
	assert.Equal(t, "acceleration-domain", ref0["instance_type"])
	assert.Equal(t, "acc-domain-001", ref0["instance_id"])
	assert.Equal(t, "test-domain.com", ref0["instance_name"])
	assert.Equal(t, "zone-abc123", ref0["zone_id"])
	assert.Equal(t, "test-zone.com", ref0["zone_name"])
	assert.Equal(t, "alias-test-zone.com", ref0["alias_zone_name"])

	ref1 := references[1].(map[string]interface{})
	assert.Equal(t, "rule-engine", ref1["instance_type"])
	assert.Equal(t, "rule-002", ref1["instance_id"])
	assert.Equal(t, "test-rule", ref1["instance_name"])
	assert.Equal(t, "zone-def456", ref1["zone_id"])
	assert.Equal(t, "another-zone.com", ref1["zone_name"])
	assert.Equal(t, "alias-another-zone.com", ref1["alias_zone_name"])
}

// TestTeoOriginGroup_Read_NilZoneFields tests Read handles nil zone fields gracefully
func TestTeoOriginGroup_Read_NilZoneFields(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	mockClient := newMockMetaForOriginGroup().client
	patches.ApplyMethodReturn(mockClient, "UseTeoClient", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeOriginGroup", func(request *teov20220901.DescribeOriginGroupRequest) (*teov20220901.DescribeOriginGroupResponse, error) {
		resp := teov20220901.NewDescribeOriginGroupResponse()
		resp.Response = &teov20220901.DescribeOriginGroupResponseParams{
			OriginGroups: []*teov20220901.OriginGroup{
				{
					GroupId: ptrStrOriginGroup("og-test456"),
					Name:    ptrStrOriginGroup("test-origin-group-2"),
					Type:    ptrStrOriginGroup("HTTP"),
					Records: []*teov20220901.OriginRecord{
						{
							Record: ptrStrOriginGroup("2.2.2.2"),
							Type:   ptrStrOriginGroup("IP_DOMAIN"),
							Weight: ptrUint64OriginGroup(50),
						},
					},
					References: []*teov20220901.OriginGroupReference{
						{
							InstanceType: ptrStrOriginGroup("loadbalance"),
							InstanceId:   ptrStrOriginGroup("lb-001"),
							InstanceName: ptrStrOriginGroup("test-lb"),
							// ZoneId, ZoneName, AliasZoneName are nil
						},
					},
					CreateTime: ptrStrOriginGroup("2024-01-01T00:00:00Z"),
					UpdateTime: ptrStrOriginGroup("2024-01-02T00:00:00Z"),
				},
			},
			RequestId: ptrStrOriginGroup("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaForOriginGroup()
	res := svcteo.ResourceTencentCloudTeoOriginGroup()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-test456",
		"type":    "HTTP",
	})
	d.SetId("zone-test456#og-test456")

	err := res.Read(d, meta)
	assert.NoError(t, err)

	references := d.Get("references").([]interface{})
	assert.Len(t, references, 1)

	ref0 := references[0].(map[string]interface{})
	assert.Equal(t, "loadbalance", ref0["instance_type"])
	assert.Equal(t, "lb-001", ref0["instance_id"])
	assert.Equal(t, "test-lb", ref0["instance_name"])
	// zone fields should be empty when API returns nil
	assert.Equal(t, "", ref0["zone_id"])
	assert.Equal(t, "", ref0["zone_name"])
	assert.Equal(t, "", ref0["alias_zone_name"])
}

// TestTeoOriginGroup_ReferencesSchema tests the references schema includes zone fields
func TestTeoOriginGroup_ReferencesSchema(t *testing.T) {
	res := svcteo.ResourceTencentCloudTeoOriginGroup()

	assert.NotNil(t, res)
	assert.Contains(t, res.Schema, "references")

	referencesSchema := res.Schema["references"]
	assert.Equal(t, schema.TypeList, referencesSchema.Type)
	assert.True(t, referencesSchema.Computed)

	referencesElem, ok := referencesSchema.Elem.(*schema.Resource)
	assert.True(t, ok)

	assert.Contains(t, referencesElem.Schema, "instance_type")
	assert.Contains(t, referencesElem.Schema, "instance_id")
	assert.Contains(t, referencesElem.Schema, "instance_name")
	assert.Contains(t, referencesElem.Schema, "zone_id")
	assert.Contains(t, referencesElem.Schema, "zone_name")
	assert.Contains(t, referencesElem.Schema, "alias_zone_name")

	// Verify zone fields are computed-only
	zoneIdSchema := referencesElem.Schema["zone_id"]
	assert.Equal(t, schema.TypeString, zoneIdSchema.Type)
	assert.True(t, zoneIdSchema.Computed)
	assert.False(t, zoneIdSchema.Optional)
	assert.False(t, zoneIdSchema.Required)

	zoneNameSchema := referencesElem.Schema["zone_name"]
	assert.Equal(t, schema.TypeString, zoneNameSchema.Type)
	assert.True(t, zoneNameSchema.Computed)
	assert.False(t, zoneNameSchema.Optional)
	assert.False(t, zoneNameSchema.Required)

	aliasZoneNameSchema := referencesElem.Schema["alias_zone_name"]
	assert.Equal(t, schema.TypeString, aliasZoneNameSchema.Type)
	assert.True(t, aliasZoneNameSchema.Computed)
	assert.False(t, aliasZoneNameSchema.Optional)
	assert.False(t, aliasZoneNameSchema.Required)
}
