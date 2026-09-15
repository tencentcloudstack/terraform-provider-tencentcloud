package ioa_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	ioav20220601 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ioa/v20220601"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	svcioa "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/ioa"
)

type mockMetaIoaCompanyDirectoryConfig struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMetaIoaCompanyDirectoryConfig) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockMetaIoaCompanyDirectoryConfig{}

func newMockMetaIoaCompanyDirectoryConfig() *mockMetaIoaCompanyDirectoryConfig {
	return &mockMetaIoaCompanyDirectoryConfig{client: &connectivity.TencentCloudClient{}}
}

func ptrStrIoa(s string) *string { return &s }

func ptrBoolIoa(b bool) *bool { return &b }

func ptrInt64Ioa(i int64) *int64 { return &i }

// TestIoaCompanyDirectoryConfigCreate tests the Create business logic with mocked SDK calls.
func TestIoaCompanyDirectoryConfigCreate(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	ioaClient := &ioav20220601.Client{}
	patches.ApplyMethodReturn(newMockMetaIoaCompanyDirectoryConfig().client, "UseIoaV20220601Client", ioaClient)

	var capturedRequest *ioav20220601.CreateCompanyDirectoryConfigRequest
	patches.ApplyMethodFunc(ioaClient, "CreateCompanyDirectoryConfigWithContext", func(ctx context.Context, request *ioav20220601.CreateCompanyDirectoryConfigRequest) (*ioav20220601.CreateCompanyDirectoryConfigResponse, error) {
		capturedRequest = request
		resp := ioav20220601.NewCreateCompanyDirectoryConfigResponse()
		resp.Response = &ioav20220601.CreateCompanyDirectoryConfigResponseParams{
			Data: &ioav20220601.DirectoryConfigResultData{
				Id:                   ptrInt64Ioa(123),
				IdentifySourceId:     ptrStrIoa("identify-1"),
				AuthSourceId:         ptrStrIoa("auth-src-1"),
				AuthConfigId:         ptrInt64Ioa(456),
				AuthPolicyId:         ptrInt64Ioa(789),
				AuthSupportPlatforms: []*string{ptrStrIoa("PC"), ptrStrIoa("Mobile")},
				AuthMethods:          []*string{ptrStrIoa("授权认证")},
			},
			RequestId: ptrStrIoa("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(ioaClient, "DescribeCompanyDirectoryConfig", func(request *ioav20220601.DescribeCompanyDirectoryConfigRequest) (*ioav20220601.DescribeCompanyDirectoryConfigResponse, error) {
		resp := ioav20220601.NewDescribeCompanyDirectoryConfigResponse()
		resp.Response = &ioav20220601.DescribeCompanyDirectoryConfigResponseParams{
			Data: &ioav20220601.DirectoryConfigData{
				Id:                 ptrInt64Ioa(123),
				Type:               ptrStrIoa("WeCom"),
				Name:               ptrStrIoa("tf-example"),
				Config:             ptrStrIoa("encrypted-hex-config-data"),
				SyncEnable:         ptrBoolIoa(true),
				SyncPolicy:         ptrStrIoa("daily"),
				SyncPolicyParams:   ptrStrIoa("{\"hour\":2}"),
				CreateAuthConfig:   ptrBoolIoa(true),
				DisplayOnLoginPage: ptrBoolIoa(true),
				Description:        ptrStrIoa("tf example description"),
				SourceId:           ptrStrIoa("src-1"),
				NameI18n: []*ioav20220601.I18nString{
					{Lang: ptrStrIoa("zh-CN"), Value: ptrStrIoa("示例目录")},
				},
			},
			RequestId: ptrStrIoa("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaIoaCompanyDirectoryConfig()
	res := svcioa.ResourceTencentCloudIoaCompanyDirectoryConfig()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"type":                  "WeCom",
		"name":                  "tf-example",
		"config":                "encrypted-hex-config-data",
		"sync_enable":           true,
		"sync_policy":           "daily",
		"sync_policy_params":    "{\"hour\":2}",
		"create_auth_config":    true,
		"display_on_login_page": true,
		"description":           "tf example description",
		"scene":                 "API",
		"name_i18n": []interface{}{
			map[string]interface{}{
				"lang":  "zh-CN",
				"value": "示例目录",
			},
		},
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "123", d.Id())

	assert.NotNil(t, capturedRequest)
	assert.Equal(t, "WeCom", *capturedRequest.Type)
	assert.Equal(t, "tf-example", *capturedRequest.Name)
	assert.Equal(t, "encrypted-hex-config-data", *capturedRequest.Config)
	assert.True(t, *capturedRequest.SyncEnable)
	assert.Equal(t, "daily", *capturedRequest.SyncPolicy)
	assert.True(t, *capturedRequest.CreateAuthConfig)
	assert.True(t, *capturedRequest.DisplayOnLoginPage)
	assert.Equal(t, "API", *capturedRequest.Scene)
	assert.Len(t, capturedRequest.NameI18n, 1)
	assert.Equal(t, "zh-CN", *capturedRequest.NameI18n[0].Lang)

	assert.Equal(t, "identify-1", d.Get("identify_source_id"))
	assert.Equal(t, "auth-src-1", d.Get("auth_source_id"))
	assert.Equal(t, 456, d.Get("auth_config_id"))
	assert.Equal(t, 789, d.Get("auth_policy_id"))
}

// TestIoaCompanyDirectoryConfigCreate_EmptyResponse verifies that an empty create response returns an error without setting id.
func TestIoaCompanyDirectoryConfigCreate_EmptyResponse(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	ioaClient := &ioav20220601.Client{}
	patches.ApplyMethodReturn(newMockMetaIoaCompanyDirectoryConfig().client, "UseIoaV20220601Client", ioaClient)

	patches.ApplyMethodFunc(ioaClient, "CreateCompanyDirectoryConfigWithContext", func(ctx context.Context, request *ioav20220601.CreateCompanyDirectoryConfigRequest) (*ioav20220601.CreateCompanyDirectoryConfigResponse, error) {
		resp := ioav20220601.NewCreateCompanyDirectoryConfigResponse()
		resp.Response = &ioav20220601.CreateCompanyDirectoryConfigResponseParams{
			Data: &ioav20220601.DirectoryConfigResultData{},
		}
		return resp, nil
	})

	meta := newMockMetaIoaCompanyDirectoryConfig()
	res := svcioa.ResourceTencentCloudIoaCompanyDirectoryConfig()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"type":                  "WeCom",
		"name":                  "tf-example",
		"config":                "encrypted-hex-config-data",
		"sync_enable":           true,
		"sync_policy":           "daily",
		"sync_policy_params":    "{\"hour\":2}",
		"create_auth_config":    true,
		"display_on_login_page": true,
	})

	err := res.Create(d, meta)
	assert.Error(t, err)
	assert.Equal(t, "", d.Id())
}

// TestIoaCompanyDirectoryConfigRead tests the Read business logic with a populated describe response.
func TestIoaCompanyDirectoryConfigRead(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	ioaClient := &ioav20220601.Client{}
	patches.ApplyMethodReturn(newMockMetaIoaCompanyDirectoryConfig().client, "UseIoaV20220601Client", ioaClient)

	patches.ApplyMethodFunc(ioaClient, "DescribeCompanyDirectoryConfig", func(request *ioav20220601.DescribeCompanyDirectoryConfigRequest) (*ioav20220601.DescribeCompanyDirectoryConfigResponse, error) {
		assert.NotNil(t, request.Id)
		resp := ioav20220601.NewDescribeCompanyDirectoryConfigResponse()
		resp.Response = &ioav20220601.DescribeCompanyDirectoryConfigResponseParams{
			Data: &ioav20220601.DirectoryConfigData{
				Id:                 ptrInt64Ioa(123),
				Type:               ptrStrIoa("WeCom"),
				Name:               ptrStrIoa("tf-example"),
				Config:             ptrStrIoa("encrypted-hex-config-data"),
				SyncEnable:         ptrBoolIoa(true),
				SyncPolicy:         ptrStrIoa("daily"),
				SyncPolicyParams:   ptrStrIoa("{\"hour\":2}"),
				CreateAuthConfig:   ptrBoolIoa(true),
				DisplayOnLoginPage: ptrBoolIoa(true),
				Description:        ptrStrIoa("updated-desc"),
				SourceId:           ptrStrIoa("src-1"),
				NameI18n: []*ioav20220601.I18nString{
					{Lang: ptrStrIoa("zh-CN"), Value: ptrStrIoa("示例目录")},
				},
			},
			RequestId: ptrStrIoa("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaIoaCompanyDirectoryConfig()
	res := svcioa.ResourceTencentCloudIoaCompanyDirectoryConfig()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
	d.SetId("123")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "123", d.Id())
	assert.Equal(t, "WeCom", d.Get("type"))
	assert.Equal(t, "tf-example", d.Get("name"))
	assert.Equal(t, "encrypted-hex-config-data", d.Get("config"))
	assert.True(t, d.Get("sync_enable").(bool))
	assert.Equal(t, "daily", d.Get("sync_policy"))
	assert.True(t, d.Get("create_auth_config").(bool))
	assert.True(t, d.Get("display_on_login_page").(bool))
	assert.Equal(t, "updated-desc", d.Get("description"))
	assert.Equal(t, "src-1", d.Get("source_id"))
}

// TestIoaCompanyDirectoryConfigRead_NotFound verifies that a nil Data clears the id.
func TestIoaCompanyDirectoryConfigRead_NotFound(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	ioaClient := &ioav20220601.Client{}
	patches.ApplyMethodReturn(newMockMetaIoaCompanyDirectoryConfig().client, "UseIoaV20220601Client", ioaClient)

	patches.ApplyMethodFunc(ioaClient, "DescribeCompanyDirectoryConfig", func(request *ioav20220601.DescribeCompanyDirectoryConfigRequest) (*ioav20220601.DescribeCompanyDirectoryConfigResponse, error) {
		resp := ioav20220601.NewDescribeCompanyDirectoryConfigResponse()
		resp.Response = &ioav20220601.DescribeCompanyDirectoryConfigResponseParams{}
		return resp, nil
	})

	meta := newMockMetaIoaCompanyDirectoryConfig()
	res := svcioa.ResourceTencentCloudIoaCompanyDirectoryConfig()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
	d.SetId("123")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "", d.Id())
}

// TestIoaCompanyDirectoryConfigUpdate tests the Update business logic with a mocked Modify call.
func TestIoaCompanyDirectoryConfigUpdate(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	ioaClient := &ioav20220601.Client{}
	patches.ApplyMethodReturn(newMockMetaIoaCompanyDirectoryConfig().client, "UseIoaV20220601Client", ioaClient)

	var capturedRequest *ioav20220601.ModifyCompanyDirectoryConfigRequest
	patches.ApplyMethodFunc(ioaClient, "ModifyCompanyDirectoryConfigWithContext", func(ctx context.Context, request *ioav20220601.ModifyCompanyDirectoryConfigRequest) (*ioav20220601.ModifyCompanyDirectoryConfigResponse, error) {
		capturedRequest = request
		resp := ioav20220601.NewModifyCompanyDirectoryConfigResponse()
		resp.Response = &ioav20220601.ModifyCompanyDirectoryConfigResponseParams{
			Data: &ioav20220601.DirectoryConfigResultData{
				Id:               ptrInt64Ioa(123),
				IdentifySourceId: ptrStrIoa("identify-2"),
			},
			RequestId: ptrStrIoa("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(ioaClient, "DescribeCompanyDirectoryConfig", func(request *ioav20220601.DescribeCompanyDirectoryConfigRequest) (*ioav20220601.DescribeCompanyDirectoryConfigResponse, error) {
		resp := ioav20220601.NewDescribeCompanyDirectoryConfigResponse()
		resp.Response = &ioav20220601.DescribeCompanyDirectoryConfigResponseParams{
			Data: &ioav20220601.DirectoryConfigData{
				Id:                 ptrInt64Ioa(123),
				Type:               ptrStrIoa("WeCom"),
				Name:               ptrStrIoa("tf-example-updated"),
				Config:             ptrStrIoa("encrypted-hex-config-data"),
				SyncEnable:         ptrBoolIoa(true),
				SyncPolicy:         ptrStrIoa("weekly"),
				SyncPolicyParams:   ptrStrIoa("{\"day\":1}"),
				CreateAuthConfig:   ptrBoolIoa(true),
				DisplayOnLoginPage: ptrBoolIoa(false),
				Description:        ptrStrIoa("updated-desc"),
				SourceId:           ptrStrIoa("src-1"),
			},
			RequestId: ptrStrIoa("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaIoaCompanyDirectoryConfig()
	res := svcioa.ResourceTencentCloudIoaCompanyDirectoryConfig()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"type":                  "WeCom",
		"name":                  "tf-example-updated",
		"config":                "encrypted-hex-config-data",
		"sync_enable":           true,
		"sync_policy":           "weekly",
		"sync_policy_params":    "{\"day\":1}",
		"create_auth_config":    true,
		"display_on_login_page": false,
		"description":           "updated-desc",
	})
	d.SetId("123")

	err := res.Update(d, meta)
	assert.NoError(t, err)

	assert.NotNil(t, capturedRequest)
	assert.Equal(t, int64(123), *capturedRequest.Id)
	assert.Equal(t, "tf-example-updated", *capturedRequest.Name)
	assert.Equal(t, "weekly", *capturedRequest.SyncPolicy)
	assert.False(t, *capturedRequest.DisplayOnLoginPage)
	assert.Equal(t, "identify-2", d.Get("identify_source_id"))
}

// TestIoaCompanyDirectoryConfigDelete tests the Delete business logic with a mocked DeleteAccountGroup call.
func TestIoaCompanyDirectoryConfigDelete(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	ioaClient := &ioav20220601.Client{}
	patches.ApplyMethodReturn(newMockMetaIoaCompanyDirectoryConfig().client, "UseIoaV20220601Client", ioaClient)

	var capturedRequest *ioav20220601.DeleteAccountGroupRequest
	patches.ApplyMethodFunc(ioaClient, "DeleteAccountGroupWithContext", func(ctx context.Context, request *ioav20220601.DeleteAccountGroupRequest) (*ioav20220601.DeleteAccountGroupResponse, error) {
		capturedRequest = request
		resp := ioav20220601.NewDeleteAccountGroupResponse()
		resp.Response = &ioav20220601.DeleteAccountGroupResponseParams{
			RequestId: ptrStrIoa("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaIoaCompanyDirectoryConfig()
	res := svcioa.ResourceTencentCloudIoaCompanyDirectoryConfig()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
	d.SetId("123")

	err := res.Delete(d, meta)
	assert.NoError(t, err)

	assert.NotNil(t, capturedRequest)
	assert.Equal(t, uint64(123), *capturedRequest.AccountGroupId)
	assert.Nil(t, capturedRequest.DomainInstanceId)
}

// TestIoaCompanyDirectoryConfigDelete_APIError verifies that delete API errors are propagated.
func TestIoaCompanyDirectoryConfigDelete_APIError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	ioaClient := &ioav20220601.Client{}
	patches.ApplyMethodReturn(newMockMetaIoaCompanyDirectoryConfig().client, "UseIoaV20220601Client", ioaClient)

	patches.ApplyMethodFunc(ioaClient, "DeleteAccountGroupWithContext", func(ctx context.Context, request *ioav20220601.DeleteAccountGroupRequest) (*ioav20220601.DeleteAccountGroupResponse, error) {
		return nil, fmt.Errorf("[TencentCloudSDKError] Code=InternalError, Message=Internal error")
	})

	meta := newMockMetaIoaCompanyDirectoryConfig()
	res := svcioa.ResourceTencentCloudIoaCompanyDirectoryConfig()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
	d.SetId("123")

	err := res.Delete(d, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "InternalError")
}
