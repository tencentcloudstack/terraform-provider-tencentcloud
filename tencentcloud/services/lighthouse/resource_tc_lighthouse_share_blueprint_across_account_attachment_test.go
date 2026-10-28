package lighthouse_test

import (
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	lighthouse "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/lighthouse/v20200324"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	svclighthouse "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/lighthouse"
)

// mockMetaForShareBlueprint implements tccommon.ProviderMeta
type mockMetaForShareBlueprint struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMetaForShareBlueprint) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockMetaForShareBlueprint{}

func newMockMetaForShareBlueprint() *mockMetaForShareBlueprint {
	return &mockMetaForShareBlueprint{client: &connectivity.TencentCloudClient{}}
}

// go test ./tencentcloud/services/lighthouse/ -run "TestShareBlueprintAcrossAccountAttachment" -v -count=1 -gcflags="all=-l"

func TestShareBlueprintAcrossAccountAttachment_Create(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	lighthouseClient := &lighthouse.Client{}
	patches.ApplyMethodReturn(newMockMetaForShareBlueprint().client, "UseLighthouseClient", lighthouseClient)

	patches.ApplyMethodFunc(lighthouseClient, "ShareBlueprintAcrossAccounts", func(request *lighthouse.ShareBlueprintAcrossAccountsRequest) (*lighthouse.ShareBlueprintAcrossAccountsResponse, error) {
		resp := lighthouse.NewShareBlueprintAcrossAccountsResponse()
		resp.Response = &lighthouse.ShareBlueprintAcrossAccountsResponseParams{
			RequestId: ptrStr("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(lighthouseClient, "DescribeBlueprintsShareAcrossAccountInfos", func(request *lighthouse.DescribeBlueprintsShareAcrossAccountInfosRequest) (*lighthouse.DescribeBlueprintsShareAcrossAccountInfosResponse, error) {
		resp := lighthouse.NewDescribeBlueprintsShareAcrossAccountInfosResponse()
		resp.Response = &lighthouse.DescribeBlueprintsShareAcrossAccountInfosResponseParams{
			TotalCount: ptrInt64(2),
			BlueprintShareAcrossAccountInfoSet: []*lighthouse.BlueprintShareAcrossAccountInfo{
				{BlueprintId: ptrStr("lhbp-xxxx"), AccountId: ptrStr("100012345678")},
				{BlueprintId: ptrStr("lhbp-xxxx"), AccountId: ptrStr("100087654321")},
			},
			RequestId: ptrStr("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaForShareBlueprint()
	res := svclighthouse.ResourceTencentCloudLighthouseShareBlueprintAcrossAccountAttachment()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"blueprint_id": "lhbp-xxxx",
		"account_ids":  []interface{}{"100012345678", "100087654321"},
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "lhbp-xxxx", d.Id())
	assert.Equal(t, "lhbp-xxxx", d.Get("blueprint_id"))
	// Verify account_ids are sorted after Create
	actualAccountIds := d.Get("account_ids").([]interface{})
	assert.Equal(t, []interface{}{"100012345678", "100087654321"}, actualAccountIds)
}

func TestShareBlueprintAcrossAccountAttachment_Create_SortsAccountIds(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	lighthouseClient := &lighthouse.Client{}
	patches.ApplyMethodReturn(newMockMetaForShareBlueprint().client, "UseLighthouseClient", lighthouseClient)

	patches.ApplyMethodFunc(lighthouseClient, "ShareBlueprintAcrossAccounts", func(request *lighthouse.ShareBlueprintAcrossAccountsRequest) (*lighthouse.ShareBlueprintAcrossAccountsResponse, error) {
		resp := lighthouse.NewShareBlueprintAcrossAccountsResponse()
		resp.Response = &lighthouse.ShareBlueprintAcrossAccountsResponseParams{
			RequestId: ptrStr("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(lighthouseClient, "DescribeBlueprintsShareAcrossAccountInfos", func(request *lighthouse.DescribeBlueprintsShareAcrossAccountInfosRequest) (*lighthouse.DescribeBlueprintsShareAcrossAccountInfosResponse, error) {
		resp := lighthouse.NewDescribeBlueprintsShareAcrossAccountInfosResponse()
		resp.Response = &lighthouse.DescribeBlueprintsShareAcrossAccountInfosResponseParams{
			TotalCount: ptrInt64(2),
			BlueprintShareAcrossAccountInfoSet: []*lighthouse.BlueprintShareAcrossAccountInfo{
				{BlueprintId: ptrStr("lhbp-xxxx"), AccountId: ptrStr("100012345678")},
				{BlueprintId: ptrStr("lhbp-xxxx"), AccountId: ptrStr("100087654321")},
			},
			RequestId: ptrStr("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaForShareBlueprint()
	res := svclighthouse.ResourceTencentCloudLighthouseShareBlueprintAcrossAccountAttachment()
	// Provide account IDs in reverse order to test sorting
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"blueprint_id": "lhbp-xxxx",
		"account_ids":  []interface{}{"100087654321", "100012345678"},
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	// Verify that account_ids are sorted after Create (should be ascending order)
	actualAccountIds := d.Get("account_ids").([]interface{})
	assert.Equal(t, []interface{}{"100012345678", "100087654321"}, actualAccountIds, "account_ids should be sorted in ascending order after Create")
}

func TestShareBlueprintAcrossAccountAttachment_Read(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	lighthouseClient := &lighthouse.Client{}
	patches.ApplyMethodReturn(newMockMetaForShareBlueprint().client, "UseLighthouseClient", lighthouseClient)

	patches.ApplyMethodFunc(lighthouseClient, "DescribeBlueprintsShareAcrossAccountInfos", func(request *lighthouse.DescribeBlueprintsShareAcrossAccountInfosRequest) (*lighthouse.DescribeBlueprintsShareAcrossAccountInfosResponse, error) {
		resp := lighthouse.NewDescribeBlueprintsShareAcrossAccountInfosResponse()
		resp.Response = &lighthouse.DescribeBlueprintsShareAcrossAccountInfosResponseParams{
			TotalCount: ptrInt64(2),
			BlueprintShareAcrossAccountInfoSet: []*lighthouse.BlueprintShareAcrossAccountInfo{
				{BlueprintId: ptrStr("lhbp-xxxx"), AccountId: ptrStr("100012345678")},
				{BlueprintId: ptrStr("lhbp-xxxx"), AccountId: ptrStr("100087654321")},
			},
			RequestId: ptrStr("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaForShareBlueprint()
	res := svclighthouse.ResourceTencentCloudLighthouseShareBlueprintAcrossAccountAttachment()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"blueprint_id": "lhbp-xxxx",
		"account_ids":  []interface{}{},
	})
	d.SetId("lhbp-xxxx")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "lhbp-xxxx", d.Get("blueprint_id"))
}

func TestShareBlueprintAcrossAccountAttachment_Read_NotFound(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	lighthouseClient := &lighthouse.Client{}
	patches.ApplyMethodReturn(newMockMetaForShareBlueprint().client, "UseLighthouseClient", lighthouseClient)

	patches.ApplyMethodFunc(lighthouseClient, "DescribeBlueprintsShareAcrossAccountInfos", func(request *lighthouse.DescribeBlueprintsShareAcrossAccountInfosRequest) (*lighthouse.DescribeBlueprintsShareAcrossAccountInfosResponse, error) {
		resp := lighthouse.NewDescribeBlueprintsShareAcrossAccountInfosResponse()
		resp.Response = &lighthouse.DescribeBlueprintsShareAcrossAccountInfosResponseParams{
			TotalCount:                         ptrInt64(0),
			BlueprintShareAcrossAccountInfoSet: []*lighthouse.BlueprintShareAcrossAccountInfo{},
			RequestId:                          ptrStr("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaForShareBlueprint()
	res := svclighthouse.ResourceTencentCloudLighthouseShareBlueprintAcrossAccountAttachment()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"blueprint_id": "lhbp-xxxx",
		"account_ids":  []interface{}{},
	})
	d.SetId("lhbp-xxxx")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "", d.Id())
}

func TestShareBlueprintAcrossAccountAttachment_Delete(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	lighthouseClient := &lighthouse.Client{}
	patches.ApplyMethodReturn(newMockMetaForShareBlueprint().client, "UseLighthouseClient", lighthouseClient)

	patches.ApplyMethodFunc(lighthouseClient, "CancelShareBlueprintAcrossAccounts", func(request *lighthouse.CancelShareBlueprintAcrossAccountsRequest) (*lighthouse.CancelShareBlueprintAcrossAccountsResponse, error) {
		resp := lighthouse.NewCancelShareBlueprintAcrossAccountsResponse()
		resp.Response = &lighthouse.CancelShareBlueprintAcrossAccountsResponseParams{
			RequestId: ptrStr("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaForShareBlueprint()
	res := svclighthouse.ResourceTencentCloudLighthouseShareBlueprintAcrossAccountAttachment()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"blueprint_id": "lhbp-xxxx",
		"account_ids":  []interface{}{"100012345678", "100087654321"},
	})
	d.SetId("lhbp-xxxx")

	err := res.Delete(d, meta)
	assert.NoError(t, err)
}

func TestShareBlueprintAcrossAccountAttachment_Import(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	lighthouseClient := &lighthouse.Client{}
	patches.ApplyMethodReturn(newMockMetaForShareBlueprint().client, "UseLighthouseClient", lighthouseClient)

	patches.ApplyMethodFunc(lighthouseClient, "DescribeBlueprintsShareAcrossAccountInfos", func(request *lighthouse.DescribeBlueprintsShareAcrossAccountInfosRequest) (*lighthouse.DescribeBlueprintsShareAcrossAccountInfosResponse, error) {
		resp := lighthouse.NewDescribeBlueprintsShareAcrossAccountInfosResponse()
		resp.Response = &lighthouse.DescribeBlueprintsShareAcrossAccountInfosResponseParams{
			TotalCount: ptrInt64(2),
			BlueprintShareAcrossAccountInfoSet: []*lighthouse.BlueprintShareAcrossAccountInfo{
				{BlueprintId: ptrStr("lhbp-xxxx"), AccountId: ptrStr("100012345678")},
				{BlueprintId: ptrStr("lhbp-xxxx"), AccountId: ptrStr("100087654321")},
			},
			RequestId: ptrStr("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaForShareBlueprint()
	res := svclighthouse.ResourceTencentCloudLighthouseShareBlueprintAcrossAccountAttachment()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"blueprint_id": "",
		"account_ids":  []interface{}{},
	})
	d.SetId("lhbp-xxxx")

	// Import uses ImportStatePassthrough which calls Read
	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "lhbp-xxxx", d.Get("blueprint_id"))
	assert.Equal(t, "lhbp-xxxx", d.Id())
}

func TestShareBlueprintAcrossAccountAttachment_Update(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	lighthouseClient := &lighthouse.Client{}
	patches.ApplyMethodReturn(newMockMetaForShareBlueprint().client, "UseLighthouseClient", lighthouseClient)

	patches.ApplyMethodFunc(lighthouseClient, "CancelShareBlueprintAcrossAccounts", func(request *lighthouse.CancelShareBlueprintAcrossAccountsRequest) (*lighthouse.CancelShareBlueprintAcrossAccountsResponse, error) {
		resp := lighthouse.NewCancelShareBlueprintAcrossAccountsResponse()
		resp.Response = &lighthouse.CancelShareBlueprintAcrossAccountsResponseParams{
			RequestId: ptrStr("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(lighthouseClient, "ShareBlueprintAcrossAccounts", func(request *lighthouse.ShareBlueprintAcrossAccountsRequest) (*lighthouse.ShareBlueprintAcrossAccountsResponse, error) {
		resp := lighthouse.NewShareBlueprintAcrossAccountsResponse()
		resp.Response = &lighthouse.ShareBlueprintAcrossAccountsResponseParams{
			RequestId: ptrStr("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(lighthouseClient, "DescribeBlueprintsShareAcrossAccountInfos", func(request *lighthouse.DescribeBlueprintsShareAcrossAccountInfosRequest) (*lighthouse.DescribeBlueprintsShareAcrossAccountInfosResponse, error) {
		resp := lighthouse.NewDescribeBlueprintsShareAcrossAccountInfosResponse()
		resp.Response = &lighthouse.DescribeBlueprintsShareAcrossAccountInfosResponseParams{
			TotalCount: ptrInt64(3),
			BlueprintShareAcrossAccountInfoSet: []*lighthouse.BlueprintShareAcrossAccountInfo{
				{BlueprintId: ptrStr("lhbp-xxxx"), AccountId: ptrStr("100012345678")},
				{BlueprintId: ptrStr("lhbp-xxxx"), AccountId: ptrStr("100087654321")},
				{BlueprintId: ptrStr("lhbp-xxxx"), AccountId: ptrStr("100099999999")},
			},
			RequestId: ptrStr("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaForShareBlueprint()
	res := svclighthouse.ResourceTencentCloudLighthouseShareBlueprintAcrossAccountAttachment()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"blueprint_id": "lhbp-xxxx",
		"account_ids":  []interface{}{"100012345678", "100087654321"},
	})
	d.SetId("lhbp-xxxx")

	// Simulate update: change account_ids to add one and remove none
	// Old: ["100012345678", "100087654321"]
	// New: ["100012345678", "100087654321", "100099999999"]
	err := d.Set("account_ids", []interface{}{"100012345678", "100087654321", "100099999999"})
	assert.NoError(t, err)

	err = res.Update(d, meta)
	assert.NoError(t, err)
}

func TestShareBlueprintAcrossAccountAttachment_Update_RemoveAndAdd(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	lighthouseClient := &lighthouse.Client{}
	patches.ApplyMethodReturn(newMockMetaForShareBlueprint().client, "UseLighthouseClient", lighthouseClient)

	removeCalled := false
	addCalled := false

	patches.ApplyMethodFunc(lighthouseClient, "CancelShareBlueprintAcrossAccounts", func(request *lighthouse.CancelShareBlueprintAcrossAccountsRequest) (*lighthouse.CancelShareBlueprintAcrossAccountsResponse, error) {
		removeCalled = true
		resp := lighthouse.NewCancelShareBlueprintAcrossAccountsResponse()
		resp.Response = &lighthouse.CancelShareBlueprintAcrossAccountsResponseParams{
			RequestId: ptrStr("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(lighthouseClient, "ShareBlueprintAcrossAccounts", func(request *lighthouse.ShareBlueprintAcrossAccountsRequest) (*lighthouse.ShareBlueprintAcrossAccountsResponse, error) {
		addCalled = true
		resp := lighthouse.NewShareBlueprintAcrossAccountsResponse()
		resp.Response = &lighthouse.ShareBlueprintAcrossAccountsResponseParams{
			RequestId: ptrStr("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(lighthouseClient, "DescribeBlueprintsShareAcrossAccountInfos", func(request *lighthouse.DescribeBlueprintsShareAcrossAccountInfosRequest) (*lighthouse.DescribeBlueprintsShareAcrossAccountInfosResponse, error) {
		resp := lighthouse.NewDescribeBlueprintsShareAcrossAccountInfosResponse()
		resp.Response = &lighthouse.DescribeBlueprintsShareAcrossAccountInfosResponseParams{
			TotalCount: ptrInt64(2),
			BlueprintShareAcrossAccountInfoSet: []*lighthouse.BlueprintShareAcrossAccountInfo{
				{BlueprintId: ptrStr("lhbp-xxxx"), AccountId: ptrStr("100012345678")},
				{BlueprintId: ptrStr("lhbp-xxxx"), AccountId: ptrStr("100099999999")},
			},
			RequestId: ptrStr("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaForShareBlueprint()
	res := svclighthouse.ResourceTencentCloudLighthouseShareBlueprintAcrossAccountAttachment()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"blueprint_id": "lhbp-xxxx",
		"account_ids":  []interface{}{"100012345678", "100087654321"},
	})
	d.SetId("lhbp-xxxx")

	// Simulate update: remove 100087654321, add 100099999999
	// Old: ["100012345678", "100087654321"]
	// New: ["100012345678", "100099999999"]
	err := d.Set("account_ids", []interface{}{"100012345678", "100099999999"})
	assert.NoError(t, err)

	err = res.Update(d, meta)
	assert.NoError(t, err)
	assert.True(t, removeCalled, "CancelShareBlueprintAcrossAccounts should be called for removed accounts")
	assert.True(t, addCalled, "ShareBlueprintAcrossAccounts should be called for added accounts")
}

func ptrStr(s string) *string {
	return &s
}

func ptrInt64(v int64) *int64 {
	return &v
}

// TestAccTencentCloudLighthouseShareBlueprintAcrossAccountAttachmentResource_basic is kept for e2e testing
func TestAccTencentCloudLighthouseShareBlueprintAcrossAccountAttachmentResource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { /*tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_PREPAY)*/ },
		Providers: nil, //testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccLighthouseShareBlueprintAcrossAccountAttachment,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_lighthouse_share_blueprint_across_account_attachment.attachment", "id")),
			},
			{
				ResourceName:      "tencentcloud_lighthouse_share_blueprint_across_account_attachment.attachment",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccLighthouseShareBlueprintAcrossAccountAttachment = `
resource "tencentcloud_lighthouse_share_blueprint_across_account_attachment" "attachment" {
  blueprint_id = "lhbp-xxxx"
  account_ids  = ["100012345678"]
}
`
