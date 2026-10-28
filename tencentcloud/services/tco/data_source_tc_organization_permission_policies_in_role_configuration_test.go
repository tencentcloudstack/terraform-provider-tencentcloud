package tco_test

import (
	"context"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	organization "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/organization/v20210331"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/tco"
)

type mockMetaPermissionPoliciesDS struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMetaPermissionPoliciesDS) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockMetaPermissionPoliciesDS{}

func newMockMetaPermissionPoliciesDS() *mockMetaPermissionPoliciesDS {
	return &mockMetaPermissionPoliciesDS{client: &connectivity.TencentCloudClient{}}
}

func ptrStringPermissionPoliciesDS(s string) *string {
	return &s
}

func ptrInt64PermissionPoliciesDS(i int64) *int64 {
	return &i
}

// TestOrganizationPermissionPoliciesInRoleConfigurationDS_ReadSuccess tests data source Read with successful response
func TestOrganizationPermissionPoliciesInRoleConfigurationDS_ReadSuccess(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	orgClient := &organization.Client{}
	patches.ApplyMethodReturn(newMockMetaPermissionPoliciesDS().client, "UseOrganizationClient", orgClient)

	patches.ApplyMethodFunc(orgClient, "ListPermissionPoliciesInRoleConfigurationWithContext", func(_ context.Context, request *organization.ListPermissionPoliciesInRoleConfigurationRequest) (*organization.ListPermissionPoliciesInRoleConfigurationResponse, error) {
		resp := organization.NewListPermissionPoliciesInRoleConfigurationResponse()
		resp.Response = &organization.ListPermissionPoliciesInRoleConfigurationResponseParams{
			TotalCounts: ptrInt64PermissionPoliciesDS(2),
			RolePolicies: []*organization.RolePolicie{
				{
					RolePolicyId:   ptrInt64PermissionPoliciesDS(1),
					RolePolicyName: ptrStringPermissionPoliciesDS("AdministratorAccess"),
					RolePolicyType: ptrStringPermissionPoliciesDS("System"),
					AddTime:        ptrStringPermissionPoliciesDS("2024-01-01 00:00:00"),
				},
				{
					RolePolicyId:       ptrInt64PermissionPoliciesDS(2),
					RolePolicyName:     ptrStringPermissionPoliciesDS("CustomPolicy"),
					RolePolicyType:     ptrStringPermissionPoliciesDS("Custom"),
					RolePolicyDocument: ptrStringPermissionPoliciesDS(`{"version":"2.0","statement":[{"effect":"allow","action":["*"],"resource":["*"]}]}`),
					AddTime:            ptrStringPermissionPoliciesDS("2024-01-02 00:00:00"),
				},
			},
			RequestId: ptrStringPermissionPoliciesDS("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaPermissionPoliciesDS()
	res := tco.DataSourceTencentCloudOrganizationPermissionPoliciesInRoleConfiguration()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":               "z-xxxxxx",
		"role_configuration_id": "rc-xxxxxx",
	})

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "z-xxxxxx"+tccommon.FILED_SP+"rc-xxxxxx", d.Id())
	assert.Equal(t, 2, d.Get("total_counts"))

	rolePolicies := d.Get("role_policies").([]interface{})
	assert.Equal(t, 2, len(rolePolicies))

	policy0 := rolePolicies[0].(map[string]interface{})
	assert.Equal(t, 1, policy0["role_policy_id"])
	assert.Equal(t, "AdministratorAccess", policy0["role_policy_name"])
	assert.Equal(t, "System", policy0["role_policy_type"])
	assert.Equal(t, "2024-01-01 00:00:00", policy0["add_time"])
	assert.Equal(t, "", policy0["role_policy_document"])

	policy1 := rolePolicies[1].(map[string]interface{})
	assert.Equal(t, 2, policy1["role_policy_id"])
	assert.Equal(t, "CustomPolicy", policy1["role_policy_name"])
	assert.Equal(t, "Custom", policy1["role_policy_type"])
	assert.Equal(t, `{"version":"2.0","statement":[{"effect":"allow","action":["*"],"resource":["*"]}]}`, policy1["role_policy_document"])
	assert.Equal(t, "2024-01-02 00:00:00", policy1["add_time"])
}

// TestOrganizationPermissionPoliciesInRoleConfigurationDS_ReadEmpty tests data source Read with nil response
func TestOrganizationPermissionPoliciesInRoleConfigurationDS_ReadEmpty(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	orgClient := &organization.Client{}
	patches.ApplyMethodReturn(newMockMetaPermissionPoliciesDS().client, "UseOrganizationClient", orgClient)

	patches.ApplyMethodFunc(orgClient, "ListPermissionPoliciesInRoleConfigurationWithContext", func(_ context.Context, request *organization.ListPermissionPoliciesInRoleConfigurationRequest) (*organization.ListPermissionPoliciesInRoleConfigurationResponse, error) {
		resp := organization.NewListPermissionPoliciesInRoleConfigurationResponse()
		resp.Response = nil
		return resp, nil
	})

	meta := newMockMetaPermissionPoliciesDS()
	res := tco.DataSourceTencentCloudOrganizationPermissionPoliciesInRoleConfiguration()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":               "z-xxxxxx",
		"role_configuration_id": "rc-xxxxxx",
	})

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "", d.Id())
}
