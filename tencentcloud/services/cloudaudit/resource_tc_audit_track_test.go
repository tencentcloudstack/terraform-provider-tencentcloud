package cloudaudit_test

import (
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	audit "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cloudaudit/v20190319"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	auditsvc "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/cloudaudit"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudAuditTrackResource_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheck(t) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccAuditTrack,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_audit_track.track", "id"),
					resource.TestCheckResourceAttr("tencentcloud_audit_track.track", "name", "terraform_track"),
					resource.TestCheckResourceAttr("tencentcloud_audit_track.track", "action_type", "Read"),
					resource.TestCheckResourceAttr("tencentcloud_audit_track.track", "status", "1"),
					resource.TestCheckResourceAttr("tencentcloud_audit_track.track", "track_for_all_members", "0"),
				),
			},
			{
				ResourceName:      "tencentcloud_audit_track.track",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccAuditTrack = `

resource "tencentcloud_audit_track" "track" {
  action_type           = "Read"
  event_names           = [
    "*",
  ]
  name                  = "terraform_track"
  resource_type         = "*"
  status                = 1
  track_for_all_members = 0

  storage {
    storage_name   = "keep-bucket"
    storage_prefix = "cloudaudit"
    storage_region = "ap-guangzhou"
    storage_type   = "cos"
  }
}
`

// ---- gomonkey-based unit tests for storage.compress ----
// Run with: go test ./tencentcloud/services/cloudaudit/ -run "TestAuditTrackCompress" -v -count=1 -gcflags="all=-l"

type mockMetaAuditTrack struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMetaAuditTrack) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockMetaAuditTrack{}

func newMockMetaAuditTrack() *mockMetaAuditTrack {
	return &mockMetaAuditTrack{client: &connectivity.TencentCloudClient{}}
}

func ptrUint64AuditTrack(v uint64) *uint64 {
	return &v
}

func ptrStringAuditTrack(s string) *string {
	return &s
}

func ptrStrListAuditTrack(values ...string) []*string {
	ret := make([]*string, 0, len(values))
	for i := range values {
		ret = append(ret, &values[i])
	}
	return ret
}

func auditTrackStorageRaw(compress int) []interface{} {
	return []interface{}{
		map[string]interface{}{
			"storage_type":   "cos",
			"storage_region": "ap-guangzhou",
			"storage_name":   "keep-bucket",
			"storage_prefix": "cloudaudit",
			"compress":       compress,
		},
	}
}

// TestAuditTrackCompress_Create verifies compress is passed to CreateAuditTrack.
func TestAuditTrackCompress_Create(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	auditClient := &audit.Client{}
	patches.ApplyMethodReturn(newMockMetaAuditTrack().client, "UseAuditClient", auditClient)

	var capturedRequest *audit.CreateAuditTrackRequest
	patches.ApplyMethodFunc(auditClient, "CreateAuditTrack", func(request *audit.CreateAuditTrackRequest) (*audit.CreateAuditTrackResponse, error) {
		capturedRequest = request
		resp := audit.NewCreateAuditTrackResponse()
		resp.Response = &audit.CreateAuditTrackResponseParams{
			TrackId:   ptrUint64AuditTrack(1001),
			RequestId: ptrStringAuditTrack("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(auditClient, "DescribeAuditTrack", func(request *audit.DescribeAuditTrackRequest) (*audit.DescribeAuditTrackResponse, error) {
		resp := audit.NewDescribeAuditTrackResponse()
		resp.Response = &audit.DescribeAuditTrackResponseParams{
			Name:         ptrStringAuditTrack("tf_audit_track"),
			ActionType:   ptrStringAuditTrack("Read"),
			ResourceType: ptrStringAuditTrack("*"),
			Status:       ptrUint64AuditTrack(1),
			EventNames:   ptrStrListAuditTrack("*"),
			Storage: &audit.Storage{
				StorageType:   ptrStringAuditTrack("cos"),
				StorageRegion: ptrStringAuditTrack("ap-guangzhou"),
				StorageName:   ptrStringAuditTrack("keep-bucket"),
				StoragePrefix: ptrStringAuditTrack("cloudaudit"),
				Compress:      ptrUint64AuditTrack(1),
			},
			TrackForAllMembers: ptrUint64AuditTrack(0),
			RequestId:          ptrStringAuditTrack("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaAuditTrack()
	res := auditsvc.ResourceTencentCloudAuditTrack()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"name":                  "tf_audit_track",
		"action_type":           "Read",
		"resource_type":         "*",
		"status":                1,
		"event_names":           []interface{}{"*"},
		"storage":               auditTrackStorageRaw(1),
		"track_for_all_members": 0,
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "1001", d.Id())

	assert.NotNil(t, capturedRequest.Storage)
	assert.NotNil(t, capturedRequest.Storage.Compress)
	assert.Equal(t, uint64(1), *capturedRequest.Storage.Compress)

	storageList := d.Get("storage").([]interface{})
	assert.Len(t, storageList, 1)
	storageMap := storageList[0].(map[string]interface{})
	assert.Equal(t, 1, storageMap["compress"].(int))
}

// TestAuditTrackCompress_CreateDefault verifies compress is not sent when not specified.
func TestAuditTrackCompress_CreateDefault(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	auditClient := &audit.Client{}
	patches.ApplyMethodReturn(newMockMetaAuditTrack().client, "UseAuditClient", auditClient)

	var capturedRequest *audit.CreateAuditTrackRequest
	patches.ApplyMethodFunc(auditClient, "CreateAuditTrack", func(request *audit.CreateAuditTrackRequest) (*audit.CreateAuditTrackResponse, error) {
		capturedRequest = request
		resp := audit.NewCreateAuditTrackResponse()
		resp.Response = &audit.CreateAuditTrackResponseParams{
			TrackId:   ptrUint64AuditTrack(1002),
			RequestId: ptrStringAuditTrack("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(auditClient, "DescribeAuditTrack", func(request *audit.DescribeAuditTrackRequest) (*audit.DescribeAuditTrackResponse, error) {
		resp := audit.NewDescribeAuditTrackResponse()
		resp.Response = &audit.DescribeAuditTrackResponseParams{
			Name:         ptrStringAuditTrack("tf_audit_track"),
			ActionType:   ptrStringAuditTrack("Read"),
			ResourceType: ptrStringAuditTrack("*"),
			Status:       ptrUint64AuditTrack(1),
			EventNames:   ptrStrListAuditTrack("*"),
			Storage: &audit.Storage{
				StorageType:   ptrStringAuditTrack("cos"),
				StorageRegion: ptrStringAuditTrack("ap-guangzhou"),
				StorageName:   ptrStringAuditTrack("keep-bucket"),
				StoragePrefix: ptrStringAuditTrack("cloudaudit"),
			},
			TrackForAllMembers: ptrUint64AuditTrack(0),
			RequestId:          ptrStringAuditTrack("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaAuditTrack()
	res := auditsvc.ResourceTencentCloudAuditTrack()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"name":          "tf_audit_track",
		"action_type":   "Read",
		"resource_type": "*",
		"status":        1,
		"event_names":   []interface{}{"*"},
		"storage": []interface{}{
			map[string]interface{}{
				"storage_type":   "cos",
				"storage_region": "ap-guangzhou",
				"storage_name":   "keep-bucket",
				"storage_prefix": "cloudaudit",
			},
		},
		"track_for_all_members": 0,
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)

	assert.NotNil(t, capturedRequest.Storage)
	assert.Nil(t, capturedRequest.Storage.Compress)
}

// TestAuditTrackCompress_Update verifies compress change triggers ModifyAuditTrack.
func TestAuditTrackCompress_Update(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	auditClient := &audit.Client{}
	patches.ApplyMethodReturn(newMockMetaAuditTrack().client, "UseAuditClient", auditClient)

	var capturedRequest *audit.ModifyAuditTrackRequest
	patches.ApplyMethodFunc(auditClient, "ModifyAuditTrack", func(request *audit.ModifyAuditTrackRequest) (*audit.ModifyAuditTrackResponse, error) {
		capturedRequest = request
		resp := audit.NewModifyAuditTrackResponse()
		resp.Response = &audit.ModifyAuditTrackResponseParams{
			RequestId: ptrStringAuditTrack("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(auditClient, "DescribeAuditTrack", func(request *audit.DescribeAuditTrackRequest) (*audit.DescribeAuditTrackResponse, error) {
		resp := audit.NewDescribeAuditTrackResponse()
		resp.Response = &audit.DescribeAuditTrackResponseParams{
			Name:         ptrStringAuditTrack("tf_audit_track"),
			ActionType:   ptrStringAuditTrack("Read"),
			ResourceType: ptrStringAuditTrack("*"),
			Status:       ptrUint64AuditTrack(1),
			EventNames:   ptrStrListAuditTrack("*"),
			Storage: &audit.Storage{
				StorageType:   ptrStringAuditTrack("cos"),
				StorageRegion: ptrStringAuditTrack("ap-guangzhou"),
				StorageName:   ptrStringAuditTrack("keep-bucket"),
				StoragePrefix: ptrStringAuditTrack("cloudaudit"),
				Compress:      ptrUint64AuditTrack(2),
			},
			TrackForAllMembers: ptrUint64AuditTrack(0),
			RequestId:          ptrStringAuditTrack("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaAuditTrack()
	res := auditsvc.ResourceTencentCloudAuditTrack()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"name":                  "tf_audit_track",
		"action_type":           "Read",
		"resource_type":         "*",
		"status":                1,
		"event_names":           []interface{}{"*"},
		"storage":               auditTrackStorageRaw(2),
		"track_for_all_members": 0,
	})
	d.SetId("1003")

	patches.ApplyMethodFunc(d, "HasChange", func(key string) bool {
		return key == "storage"
	})

	err := res.Update(d, meta)
	assert.NoError(t, err)

	assert.NotNil(t, capturedRequest)
	assert.NotNil(t, capturedRequest.Storage)
	assert.NotNil(t, capturedRequest.Storage.Compress)
	assert.Equal(t, uint64(2), *capturedRequest.Storage.Compress)

	storageList := d.Get("storage").([]interface{})
	assert.Len(t, storageList, 1)
	storageMap := storageList[0].(map[string]interface{})
	assert.Equal(t, 2, storageMap["compress"].(int))
}

// TestAuditTrackCompress_ReadNil verifies Read does not panic and skips
// compress when the API returns nil for Compress. Since the whole storage
// block is overwritten by d.Set, compress falls back to its zero value (0).
func TestAuditTrackCompress_ReadNil(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	auditClient := &audit.Client{}
	patches.ApplyMethodReturn(newMockMetaAuditTrack().client, "UseAuditClient", auditClient)

	patches.ApplyMethodFunc(auditClient, "DescribeAuditTrack", func(request *audit.DescribeAuditTrackRequest) (*audit.DescribeAuditTrackResponse, error) {
		resp := audit.NewDescribeAuditTrackResponse()
		resp.Response = &audit.DescribeAuditTrackResponseParams{
			Name:         ptrStringAuditTrack("tf_audit_track"),
			ActionType:   ptrStringAuditTrack("Read"),
			ResourceType: ptrStringAuditTrack("*"),
			Status:       ptrUint64AuditTrack(1),
			EventNames:   ptrStrListAuditTrack("*"),
			Storage: &audit.Storage{
				StorageType:   ptrStringAuditTrack("cos"),
				StorageRegion: ptrStringAuditTrack("ap-guangzhou"),
				StorageName:   ptrStringAuditTrack("keep-bucket"),
				StoragePrefix: ptrStringAuditTrack("cloudaudit"),
			},
			TrackForAllMembers: ptrUint64AuditTrack(0),
			RequestId:          ptrStringAuditTrack("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaAuditTrack()
	res := auditsvc.ResourceTencentCloudAuditTrack()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"name":                  "tf_audit_track",
		"action_type":           "Read",
		"resource_type":         "*",
		"status":                1,
		"event_names":           []interface{}{"*"},
		"storage":               auditTrackStorageRaw(1),
		"track_for_all_members": 0,
	})
	d.SetId("1004")

	err := res.Read(d, meta)
	assert.NoError(t, err)

	storageList := d.Get("storage").([]interface{})
	assert.Len(t, storageList, 1)
	storageMap := storageList[0].(map[string]interface{})
	// Compress is nil in the API response, so it is skipped and falls back to 0.
	assert.Equal(t, 0, storageMap["compress"].(int))
	// Other storage fields are still populated correctly.
	assert.Equal(t, "cos", storageMap["storage_type"].(string))
	assert.Equal(t, "keep-bucket", storageMap["storage_name"].(string))
}

// TestAuditTrackCompress_ReadValue verifies Read sets compress when API returns a value.
func TestAuditTrackCompress_ReadValue(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	auditClient := &audit.Client{}
	patches.ApplyMethodReturn(newMockMetaAuditTrack().client, "UseAuditClient", auditClient)

	patches.ApplyMethodFunc(auditClient, "DescribeAuditTrack", func(request *audit.DescribeAuditTrackRequest) (*audit.DescribeAuditTrackResponse, error) {
		resp := audit.NewDescribeAuditTrackResponse()
		resp.Response = &audit.DescribeAuditTrackResponseParams{
			Name:         ptrStringAuditTrack("tf_audit_track"),
			ActionType:   ptrStringAuditTrack("Read"),
			ResourceType: ptrStringAuditTrack("*"),
			Status:       ptrUint64AuditTrack(1),
			EventNames:   ptrStrListAuditTrack("*"),
			Storage: &audit.Storage{
				StorageType:   ptrStringAuditTrack("cos"),
				StorageRegion: ptrStringAuditTrack("ap-guangzhou"),
				StorageName:   ptrStringAuditTrack("keep-bucket"),
				StoragePrefix: ptrStringAuditTrack("cloudaudit"),
				Compress:      ptrUint64AuditTrack(2),
			},
			TrackForAllMembers: ptrUint64AuditTrack(0),
			RequestId:          ptrStringAuditTrack("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaAuditTrack()
	res := auditsvc.ResourceTencentCloudAuditTrack()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"name":                  "tf_audit_track",
		"action_type":           "Read",
		"resource_type":         "*",
		"status":                1,
		"event_names":           []interface{}{"*"},
		"storage":               auditTrackStorageRaw(0),
		"track_for_all_members": 0,
	})
	d.SetId("1005")

	err := res.Read(d, meta)
	assert.NoError(t, err)

	storageList := d.Get("storage").([]interface{})
	assert.Len(t, storageList, 1)
	storageMap := storageList[0].(map[string]interface{})
	assert.Equal(t, 2, storageMap["compress"].(int))
}

// TestAuditTrackCompress_Schema validates the compress schema definition.
func TestAuditTrackCompress_Schema(t *testing.T) {
	res := auditsvc.ResourceTencentCloudAuditTrack()

	storageField, ok := res.Schema["storage"]
	assert.True(t, ok, "storage should exist in schema")

	storageSchema := storageField.Elem.(*schema.Resource).Schema
	field, ok := storageSchema["compress"]
	assert.True(t, ok, "compress should exist in storage schema")
	assert.Equal(t, schema.TypeInt, field.Type)
	assert.True(t, field.Optional)
	assert.False(t, field.ForceNew, "compress should not be ForceNew")
}
