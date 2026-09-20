package teo_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	teov20220901 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/teo"
)

// go test ./tencentcloud/services/teo/ -run "TestTeoInferenceAPIToken" -v -count=1 -gcflags="all=-l"

// TestTeoInferenceAPIToken_Create_Success tests Create calls API and sets composite ID
func TestTeoInferenceAPIToken_Create_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "CreateInferenceAPITokenWithContext", func(_ context.Context, request *teov20220901.CreateInferenceAPITokenRequest) (*teov20220901.CreateInferenceAPITokenResponse, error) {
		assert.Equal(t, "zone-1234567890", *request.ZoneId)
		assert.Equal(t, "my-inference-token", *request.Name)
		resp := teov20220901.NewCreateInferenceAPITokenResponse()
		resp.Response = &teov20220901.CreateInferenceAPITokenResponseParams{
			TokenId:   ptrString("token-abcdefghij"),
			Content:   ptrString("token-content-secret"),
			RequestId: ptrString("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(teoClient, "DescribeInferenceAPITokensWithContext", func(_ context.Context, request *teov20220901.DescribeInferenceAPITokensRequest) (*teov20220901.DescribeInferenceAPITokensResponse, error) {
		resp := teov20220901.NewDescribeInferenceAPITokensResponse()
		resp.Response = &teov20220901.DescribeInferenceAPITokensResponseParams{
			TotalCount: ptrInt64(1),
			Tokens: []*teov20220901.InferenceAPIToken{
				{
					TokenId:    ptrString("token-abcdefghij"),
					Name:       ptrString("my-inference-token"),
					Content:    ptrString("token-content-secret"),
					CreateTime: ptrString("2025-01-01T00:00:00Z"),
				},
			},
			RequestId: ptrString("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMeta()
	res := teo.ResourceTencentCloudTeoInferenceAPIToken()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-1234567890",
		"name":    "my-inference-token",
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "zone-1234567890#token-abcdefghij", d.Id())
}

// TestTeoInferenceAPIToken_Create_EmptyTokenId tests Create handles empty TokenId
func TestTeoInferenceAPIToken_Create_EmptyTokenId(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "CreateInferenceAPITokenWithContext", func(_ context.Context, request *teov20220901.CreateInferenceAPITokenRequest) (*teov20220901.CreateInferenceAPITokenResponse, error) {
		resp := teov20220901.NewCreateInferenceAPITokenResponse()
		resp.Response = &teov20220901.CreateInferenceAPITokenResponseParams{
			TokenId:   ptrString(""),
			RequestId: ptrString("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMeta()
	res := teo.ResourceTencentCloudTeoInferenceAPIToken()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-1234567890",
		"name":    "my-inference-token",
	})

	err := res.Create(d, meta)
	assert.Error(t, err)
}

// TestTeoInferenceAPIToken_Create_APIError tests Create handles API error
func TestTeoInferenceAPIToken_Create_APIError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "CreateInferenceAPITokenWithContext", func(_ context.Context, request *teov20220901.CreateInferenceAPITokenRequest) (*teov20220901.CreateInferenceAPITokenResponse, error) {
		return nil, fmt.Errorf("[TencentCloudSDKError] Code=InvalidParameter, Message=Invalid zone_id")
	})

	meta := newMockMeta()
	res := teo.ResourceTencentCloudTeoInferenceAPIToken()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-invalid",
		"name":    "my-inference-token",
	})

	err := res.Create(d, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "InvalidParameter")
}

// TestTeoInferenceAPIToken_Read_Success tests Read retrieves token data
func TestTeoInferenceAPIToken_Read_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeInferenceAPITokensWithContext", func(_ context.Context, request *teov20220901.DescribeInferenceAPITokensRequest) (*teov20220901.DescribeInferenceAPITokensResponse, error) {
		resp := teov20220901.NewDescribeInferenceAPITokensResponse()
		resp.Response = &teov20220901.DescribeInferenceAPITokensResponseParams{
			TotalCount: ptrInt64(1),
			Tokens: []*teov20220901.InferenceAPIToken{
				{
					TokenId:    ptrString("token-abcdefghij"),
					Name:       ptrString("my-inference-token"),
					Content:    ptrString("token-content-secret"),
					CreateTime: ptrString("2025-01-01T00:00:00Z"),
				},
			},
			RequestId: ptrString("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMeta()
	res := teo.ResourceTencentCloudTeoInferenceAPIToken()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-1234567890",
		"name":    "my-inference-token",
	})
	d.SetId("zone-1234567890#token-abcdefghij")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "zone-1234567890", d.Get("zone_id"))
	assert.Equal(t, "token-abcdefghij", d.Get("token_id"))
	assert.Equal(t, "my-inference-token", d.Get("name"))
	assert.Equal(t, "token-content-secret", d.Get("content"))
	assert.Equal(t, "2025-01-01T00:00:00Z", d.Get("create_time"))
}

// TestTeoInferenceAPIToken_Read_NotFound tests Read handles token not found
func TestTeoInferenceAPIToken_Read_NotFound(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeInferenceAPITokensWithContext", func(_ context.Context, request *teov20220901.DescribeInferenceAPITokensRequest) (*teov20220901.DescribeInferenceAPITokensResponse, error) {
		resp := teov20220901.NewDescribeInferenceAPITokensResponse()
		resp.Response = &teov20220901.DescribeInferenceAPITokensResponseParams{
			TotalCount: ptrInt64(0),
			Tokens:     []*teov20220901.InferenceAPIToken{},
			RequestId:  ptrString("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMeta()
	res := teo.ResourceTencentCloudTeoInferenceAPIToken()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":  "zone-1234567890",
		"name":     "my-inference-token",
		"token_id": "token-abcdefghij",
	})
	d.SetId("zone-1234567890#token-abcdefghij")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "", d.Id())
}

// TestTeoInferenceAPIToken_Update_ImmutableError tests Update detects immutable field change
func TestTeoInferenceAPIToken_Update_ImmutableError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseTeoV20220901Client", teoClient)

	meta := newMockMeta()
	res := teo.ResourceTencentCloudTeoInferenceAPIToken()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-1234567890",
		"name":    "new-inference-token",
	})
	d.SetId("zone-1234567890#token-abcdefghij")

	patches.ApplyMethodFunc(d, "HasChange", func(key string) bool {
		return key == "name"
	})

	err := res.Update(d, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "name")
	assert.Contains(t, err.Error(), "cannot be changed")
}

// TestTeoInferenceAPIToken_Delete_Success tests Delete removes token
func TestTeoInferenceAPIToken_Delete_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DeleteInferenceAPITokenWithContext", func(_ context.Context, request *teov20220901.DeleteInferenceAPITokenRequest) (*teov20220901.DeleteInferenceAPITokenResponse, error) {
		assert.Equal(t, "zone-1234567890", *request.ZoneId)
		assert.Equal(t, "token-abcdefghij", *request.TokenId)
		resp := teov20220901.NewDeleteInferenceAPITokenResponse()
		resp.Response = &teov20220901.DeleteInferenceAPITokenResponseParams{
			RequestId: ptrString("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMeta()
	res := teo.ResourceTencentCloudTeoInferenceAPIToken()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":  "zone-1234567890",
		"name":     "my-inference-token",
		"token_id": "token-abcdefghij",
	})
	d.SetId("zone-1234567890#token-abcdefghij")

	err := res.Delete(d, meta)
	assert.NoError(t, err)
}

// TestTeoInferenceAPIToken_Delete_APIError tests Delete handles API error
func TestTeoInferenceAPIToken_Delete_APIError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DeleteInferenceAPITokenWithContext", func(_ context.Context, request *teov20220901.DeleteInferenceAPITokenRequest) (*teov20220901.DeleteInferenceAPITokenResponse, error) {
		return nil, fmt.Errorf("[TencentCloudSDKError] Code=ResourceNotFound, Message=Token not found")
	})

	meta := newMockMeta()
	res := teo.ResourceTencentCloudTeoInferenceAPIToken()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":  "zone-1234567890",
		"name":     "my-inference-token",
		"token_id": "token-abcdefghij",
	})
	d.SetId("zone-1234567890#token-abcdefghij")

	err := res.Delete(d, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "ResourceNotFound")
}

// TestTeoInferenceAPIToken_Schema validates schema definition
func TestTeoInferenceAPIToken_Schema(t *testing.T) {
	res := teo.ResourceTencentCloudTeoInferenceAPIToken()

	assert.NotNil(t, res)
	assert.NotNil(t, res.Create)
	assert.NotNil(t, res.Read)
	assert.NotNil(t, res.Update)
	assert.NotNil(t, res.Delete)
	assert.NotNil(t, res.Importer)

	assert.Contains(t, res.Schema, "zone_id")
	zoneId := res.Schema["zone_id"]
	assert.Equal(t, schema.TypeString, zoneId.Type)
	assert.True(t, zoneId.Required)
	assert.True(t, zoneId.ForceNew)

	assert.Contains(t, res.Schema, "name")
	name := res.Schema["name"]
	assert.Equal(t, schema.TypeString, name.Type)
	assert.True(t, name.Required)
	assert.True(t, name.ForceNew)

	assert.Contains(t, res.Schema, "token_id")
	tokenId := res.Schema["token_id"]
	assert.Equal(t, schema.TypeString, tokenId.Type)
	assert.True(t, tokenId.Computed)

	assert.Contains(t, res.Schema, "content")
	content := res.Schema["content"]
	assert.Equal(t, schema.TypeString, content.Type)
	assert.True(t, content.Computed)
	assert.True(t, content.Sensitive)

	assert.Contains(t, res.Schema, "create_time")
	createTime := res.Schema["create_time"]
	assert.Equal(t, schema.TypeString, createTime.Type)
	assert.True(t, createTime.Computed)
}

func ptrInt64(n int64) *int64 {
	return &n
}
