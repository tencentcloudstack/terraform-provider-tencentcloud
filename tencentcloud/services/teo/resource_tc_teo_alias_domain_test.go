package teo_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	teov20220901 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/teo"
)

type mockMetaForAliasDomain struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMetaForAliasDomain) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockMetaForAliasDomain{}

func newMockMetaForAliasDomain() *mockMetaForAliasDomain {
	return &mockMetaForAliasDomain{client: &connectivity.TencentCloudClient{}}
}

func ptrStrAliasDomain(s string) *string {
	return &s
}

func ptrInt64AliasDomain(i int64) *int64 {
	return &i
}

func makeDescribeAliasDomainsResponse() *teov20220901.DescribeAliasDomainsResponse {
	resp := teov20220901.NewDescribeAliasDomainsResponse()
	resp.Response = &teov20220901.DescribeAliasDomainsResponseParams{
		TotalCount: ptrInt64AliasDomain(1),
		AliasDomains: []*teov20220901.AliasDomain{
			{
				AliasName:  ptrStrAliasDomain("alias.example.com"),
				ZoneId:     ptrStrAliasDomain("zone-test123"),
				TargetName: ptrStrAliasDomain("target.example.com"),
				Status:     ptrStrAliasDomain("active"),
				ForbidMode: ptrInt64AliasDomain(0),
				CreatedOn:  ptrStrAliasDomain("2024-01-01T00:00:00Z"),
				ModifiedOn: ptrStrAliasDomain("2024-01-01T00:00:00Z"),
			},
		},
		RequestId: ptrStrAliasDomain("fake-request-id"),
	}
	return resp
}

func TestAliasDomainReadSuccess(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaForAliasDomain().client, "UseTeoClient", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeAliasDomains", func(request *teov20220901.DescribeAliasDomainsRequest) (*teov20220901.DescribeAliasDomainsResponse, error) {
		return makeDescribeAliasDomainsResponse(), nil
	})

	meta := newMockMetaForAliasDomain()
	res := teo.ResourceTencentCloudTeoAliasDomain()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":     "zone-test123",
		"alias_name":  "alias.example.com",
		"target_name": "target.example.com",
	})
	d.SetId("zone-test123#alias.example.com")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "target.example.com", d.Get("target_name"))
	assert.Equal(t, "alias.example.com", d.Get("alias_name"))
	assert.Equal(t, "active", d.Get("status"))
	assert.Equal(t, 0, d.Get("forbid_mode"))
	assert.Equal(t, "2024-01-01T00:00:00Z", d.Get("created_on"))
	assert.Equal(t, "2024-01-01T00:00:00Z", d.Get("modified_on"))
}

func TestAliasDomainReadNotFound(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaForAliasDomain().client, "UseTeoClient", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeAliasDomains", func(request *teov20220901.DescribeAliasDomainsRequest) (*teov20220901.DescribeAliasDomainsResponse, error) {
		resp := teov20220901.NewDescribeAliasDomainsResponse()
		resp.Response = &teov20220901.DescribeAliasDomainsResponseParams{
			TotalCount:  ptrInt64AliasDomain(0),
			AliasDomains: []*teov20220901.AliasDomain{},
			RequestId:   ptrStrAliasDomain("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaForAliasDomain()
	res := teo.ResourceTencentCloudTeoAliasDomain()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":     "zone-test123",
		"alias_name":  "alias.example.com",
		"target_name": "target.example.com",
	})
	d.SetId("zone-test123#alias.example.com")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "", d.Id())
}

func TestAliasDomainReadApiError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaForAliasDomain().client, "UseTeoClient", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeAliasDomains", func(request *teov20220901.DescribeAliasDomainsRequest) (*teov20220901.DescribeAliasDomainsResponse, error) {
		return nil, fmt.Errorf("[TencentCloudSDKError] Code=InternalError, Message=internal error")
	})

	meta := newMockMetaForAliasDomain()
	res := teo.ResourceTencentCloudTeoAliasDomain()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":     "zone-test123",
		"alias_name":  "alias.example.com",
		"target_name": "target.example.com",
	})
	d.SetId("zone-test123#alias.example.com")

	err := res.Read(d, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "InternalError")
}

func TestAliasDomainCreateSuccess(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaForAliasDomain().client, "UseTeoClient", teoClient)

	patches.ApplyMethodFunc(teoClient, "CreateAliasDomainWithContext", func(_ context.Context, request *teov20220901.CreateAliasDomainRequest) (*teov20220901.CreateAliasDomainResponse, error) {
		resp := teov20220901.NewCreateAliasDomainResponse()
		resp.Response = &teov20220901.CreateAliasDomainResponseParams{
			RequestId: ptrStrAliasDomain("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(teoClient, "DescribeAliasDomains", func(request *teov20220901.DescribeAliasDomainsRequest) (*teov20220901.DescribeAliasDomainsResponse, error) {
		return makeDescribeAliasDomainsResponse(), nil
	})

	meta := newMockMetaForAliasDomain()
	res := teo.ResourceTencentCloudTeoAliasDomain()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":     "zone-test123",
		"alias_name":  "alias.example.com",
		"target_name": "target.example.com",
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "zone-test123#alias.example.com", d.Id())
}

func TestAliasDomainCreateWithCert(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaForAliasDomain().client, "UseTeoClient", teoClient)

	patches.ApplyMethodFunc(teoClient, "CreateAliasDomainWithContext", func(_ context.Context, request *teov20220901.CreateAliasDomainRequest) (*teov20220901.CreateAliasDomainResponse, error) {
		resp := teov20220901.NewCreateAliasDomainResponse()
		resp.Response = &teov20220901.CreateAliasDomainResponseParams{
			RequestId: ptrStrAliasDomain("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(teoClient, "DescribeAliasDomains", func(request *teov20220901.DescribeAliasDomainsRequest) (*teov20220901.DescribeAliasDomainsResponse, error) {
		return makeDescribeAliasDomainsResponse(), nil
	})

	meta := newMockMetaForAliasDomain()
	res := teo.ResourceTencentCloudTeoAliasDomain()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":     "zone-test123",
		"alias_name":  "alias.example.com",
		"target_name": "target.example.com",
		"cert_type":   "hosting",
		"cert_id":     []interface{}{"cert-abc123"},
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "zone-test123#alias.example.com", d.Id())
}

func TestAliasDomainUpdateSuccess(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaForAliasDomain().client, "UseTeoClient", teoClient)

	patches.ApplyMethodFunc(teoClient, "ModifyAliasDomainWithContext", func(_ context.Context, request *teov20220901.ModifyAliasDomainRequest) (*teov20220901.ModifyAliasDomainResponse, error) {
		resp := teov20220901.NewModifyAliasDomainResponse()
		resp.Response = &teov20220901.ModifyAliasDomainResponseParams{
			RequestId: ptrStrAliasDomain("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(teoClient, "DescribeAliasDomains", func(request *teov20220901.DescribeAliasDomainsRequest) (*teov20220901.DescribeAliasDomainsResponse, error) {
		resp := teov20220901.NewDescribeAliasDomainsResponse()
		resp.Response = &teov20220901.DescribeAliasDomainsResponseParams{
			TotalCount: ptrInt64AliasDomain(1),
			AliasDomains: []*teov20220901.AliasDomain{
				{
					AliasName:  ptrStrAliasDomain("alias.example.com"),
					ZoneId:     ptrStrAliasDomain("zone-test123"),
					TargetName: ptrStrAliasDomain("target2.example.com"),
					Status:     ptrStrAliasDomain("active"),
					ForbidMode: ptrInt64AliasDomain(0),
					CreatedOn:  ptrStrAliasDomain("2024-01-01T00:00:00Z"),
					ModifiedOn: ptrStrAliasDomain("2024-01-03T00:00:00Z"),
				},
			},
			RequestId: ptrStrAliasDomain("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaForAliasDomain()
	res := teo.ResourceTencentCloudTeoAliasDomain()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":     "zone-test123",
		"alias_name":  "alias.example.com",
		"target_name": "target2.example.com",
	})
	d.SetId("zone-test123#alias.example.com")

	err := res.Update(d, meta)
	assert.NoError(t, err)
}

func TestAliasDomainDeleteSuccess(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaForAliasDomain().client, "UseTeoClient", teoClient)

	patches.ApplyMethodFunc(teoClient, "DeleteAliasDomainWithContext", func(_ context.Context, request *teov20220901.DeleteAliasDomainRequest) (*teov20220901.DeleteAliasDomainResponse, error) {
		resp := teov20220901.NewDeleteAliasDomainResponse()
		resp.Response = &teov20220901.DeleteAliasDomainResponseParams{
			RequestId: ptrStrAliasDomain("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaForAliasDomain()
	res := teo.ResourceTencentCloudTeoAliasDomain()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":     "zone-test123",
		"alias_name":  "alias.example.com",
		"target_name": "target.example.com",
	})
	d.SetId("zone-test123#alias.example.com")

	err := res.Delete(d, meta)
	assert.NoError(t, err)
}

func TestAliasDomainSchema(t *testing.T) {
	res := teo.ResourceTencentCloudTeoAliasDomain()

	assert.NotNil(t, res)
	assert.Contains(t, res.Schema, "zone_id")
	assert.Contains(t, res.Schema, "alias_name")
	assert.Contains(t, res.Schema, "target_name")
	assert.Contains(t, res.Schema, "cert_type")
	assert.Contains(t, res.Schema, "cert_id")
	assert.Contains(t, res.Schema, "status")
	assert.Contains(t, res.Schema, "forbid_mode")
	assert.Contains(t, res.Schema, "created_on")
	assert.Contains(t, res.Schema, "modified_on")

	assert.True(t, res.Schema["zone_id"].ForceNew)
	assert.True(t, res.Schema["alias_name"].ForceNew)

	assert.True(t, res.Schema["status"].Computed)
	assert.True(t, res.Schema["forbid_mode"].Computed)
	assert.True(t, res.Schema["created_on"].Computed)
	assert.True(t, res.Schema["modified_on"].Computed)

	assert.Equal(t, schema.TypeList, res.Schema["cert_id"].Type)
	assert.NotNil(t, res.Importer)
}
