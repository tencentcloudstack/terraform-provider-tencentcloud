package lighthouse_test

import (
	"errors"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	lighthousev20200324 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/lighthouse/v20200324"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/lighthouse"
)

type mockMetaShareBlueprint struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMetaShareBlueprint) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockMetaShareBlueprint{}

func newMockMetaShareBlueprint() *mockMetaShareBlueprint {
	return &mockMetaShareBlueprint{client: &connectivity.TencentCloudClient{}}
}

func ptrStringShareBlueprint(s string) *string {
	return &s
}

// go test ./tencentcloud/services/lighthouse/ -run "TestLighthouseShareBlueprintAcrossAccount" -v -count=1 -gcflags="all=-l"

func TestLighthouseShareBlueprintAcrossAccount_CreateSuccess(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	lighthouseClient := &lighthousev20200324.Client{}
	patches.ApplyMethodReturn(newMockMetaShareBlueprint().client, "UseLighthouseClient", lighthouseClient)

	patches.ApplyMethodFunc(lighthouseClient, "ShareBlueprintAcrossAccounts", func(request *lighthousev20200324.ShareBlueprintAcrossAccountsRequest) (*lighthousev20200324.ShareBlueprintAcrossAccountsResponse, error) {
		assert.NotNil(t, request.BlueprintId)
		assert.Equal(t, "lhbp-xxxxxx", *request.BlueprintId)
		assert.NotNil(t, request.AccountIds)
		assert.Equal(t, 2, len(request.AccountIds))
		assert.Equal(t, "100000000001", *request.AccountIds[0])
		assert.Equal(t, "100000000002", *request.AccountIds[1])

		resp := lighthousev20200324.NewShareBlueprintAcrossAccountsResponse()
		resp.Response = &lighthousev20200324.ShareBlueprintAcrossAccountsResponseParams{
			RequestId: ptrStringShareBlueprint("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaShareBlueprint()
	res := lighthouse.ResourceTencentCloudLighthouseShareBlueprintAcrossAccountOperation()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"blueprint_id": "lhbp-xxxxxx",
		"account_ids":  []interface{}{"100000000001", "100000000002"},
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.NotEmpty(t, d.Id())
	assert.Equal(t, "lhbp-xxxxxx", d.Get("blueprint_id"))
}

func TestLighthouseShareBlueprintAcrossAccount_CreateApiError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	lighthouseClient := &lighthousev20200324.Client{}
	patches.ApplyMethodReturn(newMockMetaShareBlueprint().client, "UseLighthouseClient", lighthouseClient)

	patches.ApplyMethodFunc(lighthouseClient, "ShareBlueprintAcrossAccounts", func(request *lighthousev20200324.ShareBlueprintAcrossAccountsRequest) (*lighthousev20200324.ShareBlueprintAcrossAccountsResponse, error) {
		return nil, errors.New("InternalError")
	})

	meta := newMockMetaShareBlueprint()
	res := lighthouse.ResourceTencentCloudLighthouseShareBlueprintAcrossAccountOperation()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"blueprint_id": "lhbp-xxxxxx",
		"account_ids":  []interface{}{"100000000001"},
	})

	err := res.Create(d, meta)
	assert.Error(t, err)
}

func TestLighthouseShareBlueprintAcrossAccount_ReadNoOp(t *testing.T) {
	meta := newMockMetaShareBlueprint()
	res := lighthouse.ResourceTencentCloudLighthouseShareBlueprintAcrossAccountOperation()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"blueprint_id": "lhbp-xxxxxx",
		"account_ids":  []interface{}{"100000000001"},
	})
	d.SetId("test-id")

	err := res.Read(d, meta)
	assert.NoError(t, err)
}

func TestLighthouseShareBlueprintAcrossAccount_DeleteNoOp(t *testing.T) {
	meta := newMockMetaShareBlueprint()
	res := lighthouse.ResourceTencentCloudLighthouseShareBlueprintAcrossAccountOperation()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"blueprint_id": "lhbp-xxxxxx",
		"account_ids":  []interface{}{"100000000001"},
	})
	d.SetId("test-id")

	err := res.Delete(d, meta)
	assert.NoError(t, err)
}
