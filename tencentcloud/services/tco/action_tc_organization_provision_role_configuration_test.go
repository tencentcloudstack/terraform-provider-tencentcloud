package tco_test

import (
	"context"
	"fmt"
	"math/big"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-framework/action"
	action_schema "github.com/hashicorp/terraform-plugin-framework/action/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/stretchr/testify/assert"
	organization "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/organization/v20210331"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/sharedmeta"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/tco"
)

func ptrStringOrganizationProvisionRoleConfiguration(s string) *string {
	return &s
}

// organizationProvisionRoleConfigurationSchema returns the framework action
// schema used to build the raw config value for InvokeRequest.
func organizationProvisionRoleConfigurationSchema(t *testing.T) action_schema.Schema {
	a := tco.NewTcoProvisionRoleConfiguration()
	schemaResp := &action.SchemaResponse{}
	a.Schema(context.Background(), action.SchemaRequest{}, schemaResp)
	if schemaResp.Diagnostics.HasError() {
		t.Fatalf("failed to build action schema: %v", schemaResp.Diagnostics)
	}
	return schemaResp.Schema
}

// provisionTimeoutsObjType builds the tftypes.Object type for the timeouts block.
func provisionTimeoutsObjType() tftypes.Object {
	return tftypes.Object{
		AttributeTypes: map[string]tftypes.Type{
			"invoke": tftypes.String,
		},
	}
}

// newOrganizationProvisionRoleConfigurationInvokeRequest builds an
// action.InvokeRequest carrying the given raw attribute values. The timeouts
// block is left unconfigured so the built-in default applies.
func newOrganizationProvisionRoleConfigurationInvokeRequest(t *testing.T, zoneId, roleConfigurationId, targetType string, targetUin int64) action.InvokeRequest {
	s := organizationProvisionRoleConfigurationSchema(t)

	raw := tftypes.NewValue(tftypes.Object{
		AttributeTypes: map[string]tftypes.Type{
			"zone_id":               tftypes.String,
			"role_configuration_id": tftypes.String,
			"target_type":           tftypes.String,
			"target_uin":            tftypes.Number,
			"timeouts":              provisionTimeoutsObjType(),
		},
	}, map[string]tftypes.Value{
		"zone_id":               tftypes.NewValue(tftypes.String, zoneId),
		"role_configuration_id": tftypes.NewValue(tftypes.String, roleConfigurationId),
		"target_type":           tftypes.NewValue(tftypes.String, targetType),
		"target_uin":            tftypes.NewValue(tftypes.Number, big.NewFloat(float64(targetUin))),
		"timeouts":              tftypes.NewValue(provisionTimeoutsObjType(), nil),
	})

	req := action.InvokeRequest{}
	req.Config.Raw = raw
	req.Config.Schema = s
	return req
}

// setOrganizationProvisionRoleConfigurationClient injects a mock connectivity
// client into the action instance via the Configure method.
func setOrganizationProvisionRoleConfigurationClient(a *tco.TcoProvisionRoleConfiguration, client *connectivity.TencentCloudClient) {
	meta := &sharedmeta.ProviderMeta{Client: client}
	configureReq := action.ConfigureRequest{ProviderData: meta}
	configureResp := &action.ConfigureResponse{}
	a.Configure(context.Background(), configureReq, configureResp)
}

// go test ./tencentcloud/services/tco/ -run "TestOrganizationProvisionRoleConfiguration" -v -count=1 -gcflags="all=-l"

// TestOrganizationProvisionRoleConfiguration_Invoke_Success tests a successful
// invoke, including waiting for the asynchronous task to succeed.
func TestOrganizationProvisionRoleConfiguration_Invoke_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	organizationClient := &organization.Client{}
	patches.ApplyMethodReturn(client, "UseOrganizationClient", organizationClient)

	var gotRequest *organization.ProvisionRoleConfigurationRequest
	patches.ApplyMethodFunc(organizationClient, "ProvisionRoleConfigurationWithContext", func(_ context.Context, request *organization.ProvisionRoleConfigurationRequest) (*organization.ProvisionRoleConfigurationResponse, error) {
		gotRequest = request
		resp := organization.NewProvisionRoleConfigurationResponse()
		resp.Response = &organization.ProvisionRoleConfigurationResponseParams{
			Task: &organization.RoleProvisioningsTask{
				TaskId:              ptrStringOrganizationProvisionRoleConfiguration("task-1111222233334444"),
				RoleConfigurationId: ptrStringOrganizationProvisionRoleConfiguration("rc-1111222233334444"),
			},
			RequestId: ptrStringOrganizationProvisionRoleConfiguration("fake-request-id"),
		}
		return resp, nil
	})

	var gotTaskRequest *organization.GetTaskStatusRequest
	patches.ApplyMethodFunc(organizationClient, "GetTaskStatusWithContext", func(_ context.Context, request *organization.GetTaskStatusRequest) (*organization.GetTaskStatusResponse, error) {
		gotTaskRequest = request
		resp := organization.NewGetTaskStatusResponse()
		resp.Response = &organization.GetTaskStatusResponseParams{
			TaskStatus: &organization.TaskStatus{
				Status: ptrStringOrganizationProvisionRoleConfiguration("Success"),
				TaskId: ptrStringOrganizationProvisionRoleConfiguration("task-1111222233334444"),
			},
			RequestId: ptrStringOrganizationProvisionRoleConfiguration("fake-request-id"),
		}
		return resp, nil
	})

	a := tco.NewTcoProvisionRoleConfiguration()
	setOrganizationProvisionRoleConfigurationClient(a.(*tco.TcoProvisionRoleConfiguration), client)

	req := newOrganizationProvisionRoleConfigurationInvokeRequest(t, "z-1111222233334444", "rc-1111222233334444", "MemberUin", 100000000001)
	resp := &action.InvokeResponse{}
	a.Invoke(context.Background(), req, resp)

	assert.False(t, resp.Diagnostics.HasError(), "expected no error diagnostics, got: %v", resp.Diagnostics)

	assert.NotNil(t, gotRequest)
	assert.Equal(t, "z-1111222233334444", *gotRequest.ZoneId)
	assert.Equal(t, "rc-1111222233334444", *gotRequest.RoleConfigurationId)
	assert.Equal(t, "MemberUin", *gotRequest.TargetType)
	assert.Equal(t, int64(100000000001), *gotRequest.TargetUin)

	assert.NotNil(t, gotTaskRequest)
	assert.Equal(t, "z-1111222233334444", *gotTaskRequest.ZoneId)
	assert.Equal(t, "task-1111222233334444", *gotTaskRequest.TaskId)
}

// TestOrganizationProvisionRoleConfiguration_Invoke_TaskFailed tests that a
// failed asynchronous task is reported as an error.
func TestOrganizationProvisionRoleConfiguration_Invoke_TaskFailed(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	organizationClient := &organization.Client{}
	patches.ApplyMethodReturn(client, "UseOrganizationClient", organizationClient)

	patches.ApplyMethodFunc(organizationClient, "ProvisionRoleConfigurationWithContext", func(_ context.Context, request *organization.ProvisionRoleConfigurationRequest) (*organization.ProvisionRoleConfigurationResponse, error) {
		resp := organization.NewProvisionRoleConfigurationResponse()
		resp.Response = &organization.ProvisionRoleConfigurationResponseParams{
			Task: &organization.RoleProvisioningsTask{
				TaskId: ptrStringOrganizationProvisionRoleConfiguration("task-1111222233334444"),
			},
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(organizationClient, "GetTaskStatusWithContext", func(_ context.Context, request *organization.GetTaskStatusRequest) (*organization.GetTaskStatusResponse, error) {
		resp := organization.NewGetTaskStatusResponse()
		resp.Response = &organization.GetTaskStatusResponseParams{
			TaskStatus: &organization.TaskStatus{
				Status:        ptrStringOrganizationProvisionRoleConfiguration("Failed"),
				TaskId:        ptrStringOrganizationProvisionRoleConfiguration("task-1111222233334444"),
				FailureReason: ptrStringOrganizationProvisionRoleConfiguration("target account not exist"),
			},
		}
		return resp, nil
	})

	a := tco.NewTcoProvisionRoleConfiguration()
	setOrganizationProvisionRoleConfigurationClient(a.(*tco.TcoProvisionRoleConfiguration), client)

	req := newOrganizationProvisionRoleConfigurationInvokeRequest(t, "z-1111222233334444", "rc-1111222233334444", "MemberUin", 100000000001)
	resp := &action.InvokeResponse{}
	a.Invoke(context.Background(), req, resp)

	assert.True(t, resp.Diagnostics.HasError())
	assert.Contains(t, resp.Diagnostics.Errors()[0].Summary(), "waiting for the role configuration provisioning task")
	assert.Contains(t, resp.Diagnostics.Errors()[0].Detail(), "target account not exist")
}

// TestOrganizationProvisionRoleConfiguration_Invoke_MissingTaskId tests that a
// response without a task id is reported as an error instead of polling forever.
func TestOrganizationProvisionRoleConfiguration_Invoke_MissingTaskId(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	organizationClient := &organization.Client{}
	patches.ApplyMethodReturn(client, "UseOrganizationClient", organizationClient)

	patches.ApplyMethodFunc(organizationClient, "ProvisionRoleConfigurationWithContext", func(_ context.Context, request *organization.ProvisionRoleConfigurationRequest) (*organization.ProvisionRoleConfigurationResponse, error) {
		resp := organization.NewProvisionRoleConfigurationResponse()
		resp.Response = &organization.ProvisionRoleConfigurationResponseParams{}
		return resp, nil
	})

	a := tco.NewTcoProvisionRoleConfiguration()
	setOrganizationProvisionRoleConfigurationClient(a.(*tco.TcoProvisionRoleConfiguration), client)

	req := newOrganizationProvisionRoleConfigurationInvokeRequest(t, "z-1111222233334444", "rc-1111222233334444", "MemberUin", 100000000001)
	resp := &action.InvokeResponse{}
	a.Invoke(context.Background(), req, resp)

	assert.True(t, resp.Diagnostics.HasError())
	assert.Contains(t, resp.Diagnostics.Errors()[0].Summary(), "provisioning role configuration")
	assert.Contains(t, resp.Diagnostics.Errors()[0].Detail(), "does not contain a task id")
}

// TestOrganizationProvisionRoleConfiguration_Invoke_APIError tests API error handling
func TestOrganizationProvisionRoleConfiguration_Invoke_APIError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	organizationClient := &organization.Client{}
	patches.ApplyMethodReturn(client, "UseOrganizationClient", organizationClient)

	patches.ApplyMethodFunc(organizationClient, "ProvisionRoleConfigurationWithContext", func(_ context.Context, request *organization.ProvisionRoleConfigurationRequest) (*organization.ProvisionRoleConfigurationResponse, error) {
		return nil, fmt.Errorf("[TencentCloudSDKError] Code=ResourceNotFound.RoleConfigurationNotExist, Message=role configuration not exist")
	})

	a := tco.NewTcoProvisionRoleConfiguration()
	setOrganizationProvisionRoleConfigurationClient(a.(*tco.TcoProvisionRoleConfiguration), client)

	req := newOrganizationProvisionRoleConfigurationInvokeRequest(t, "z-1111222233334444", "rc-invalid", "MemberUin", 100000000001)
	resp := &action.InvokeResponse{}
	a.Invoke(context.Background(), req, resp)

	assert.True(t, resp.Diagnostics.HasError())
	assert.Contains(t, resp.Diagnostics.Errors()[0].Summary(), "provisioning role configuration")
	assert.Contains(t, resp.Diagnostics.Errors()[0].Detail(), "ResourceNotFound.RoleConfigurationNotExist")
}

// TestOrganizationProvisionRoleConfiguration_Invoke_MissingInput tests missing
// required inputs are rejected before any API call
func TestOrganizationProvisionRoleConfiguration_Invoke_MissingInput(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	organizationClient := &organization.Client{}
	patches.ApplyMethodReturn(client, "UseOrganizationClient", organizationClient)

	apiCalled := false
	patches.ApplyMethodFunc(organizationClient, "ProvisionRoleConfigurationWithContext", func(_ context.Context, request *organization.ProvisionRoleConfigurationRequest) (*organization.ProvisionRoleConfigurationResponse, error) {
		apiCalled = true
		return organization.NewProvisionRoleConfigurationResponse(), nil
	})

	a := tco.NewTcoProvisionRoleConfiguration()
	setOrganizationProvisionRoleConfigurationClient(a.(*tco.TcoProvisionRoleConfiguration), client)

	cases := []struct {
		name                string
		zoneId              string
		roleConfigurationId string
		targetType          string
		targetUin           int64
		expectedSummary     string
	}{
		{
			name:                "empty zone_id",
			zoneId:              "",
			roleConfigurationId: "rc-1111222233334444",
			targetType:          "MemberUin",
			targetUin:           100000000001,
			expectedSummary:     "Missing zone_id",
		},
		{
			name:                "empty role_configuration_id",
			zoneId:              "z-1111222233334444",
			roleConfigurationId: "",
			targetType:          "MemberUin",
			targetUin:           100000000001,
			expectedSummary:     "Missing role_configuration_id",
		},
		{
			name:                "empty target_type",
			zoneId:              "z-1111222233334444",
			roleConfigurationId: "rc-1111222233334444",
			targetType:          "",
			targetUin:           100000000001,
			expectedSummary:     "Missing target_type",
		},
		{
			name:                "zero target_uin",
			zoneId:              "z-1111222233334444",
			roleConfigurationId: "rc-1111222233334444",
			targetType:          "MemberUin",
			targetUin:           0,
			expectedSummary:     "Missing target_uin",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			req := newOrganizationProvisionRoleConfigurationInvokeRequest(t, c.zoneId, c.roleConfigurationId, c.targetType, c.targetUin)
			resp := &action.InvokeResponse{}
			a.Invoke(context.Background(), req, resp)

			assert.True(t, resp.Diagnostics.HasError())
			assert.Contains(t, resp.Diagnostics.Errors()[0].Summary(), c.expectedSummary)
			assert.False(t, apiCalled)
		})
	}
}

// TestOrganizationProvisionRoleConfiguration_Invoke_ClientNotConfigured tests the nil-client guard
func TestOrganizationProvisionRoleConfiguration_Invoke_ClientNotConfigured(t *testing.T) {
	a := tco.NewTcoProvisionRoleConfiguration()

	req := newOrganizationProvisionRoleConfigurationInvokeRequest(t, "z-1111222233334444", "rc-1111222233334444", "MemberUin", 100000000001)
	resp := &action.InvokeResponse{}
	a.Invoke(context.Background(), req, resp)

	assert.True(t, resp.Diagnostics.HasError())
	assert.Contains(t, resp.Diagnostics.Errors()[0].Summary(), "Provider not configured")
}

// TestOrganizationProvisionRoleConfiguration_MetadataAndSchema validates metadata and schema definition
func TestOrganizationProvisionRoleConfiguration_MetadataAndSchema(t *testing.T) {
	a := tco.NewTcoProvisionRoleConfiguration()

	metaResp := &action.MetadataResponse{}
	a.Metadata(context.Background(), action.MetadataRequest{}, metaResp)
	assert.Equal(t, "tencentcloud_organization_provision_role_configuration", metaResp.TypeName)

	schemaResp := &action.SchemaResponse{}
	a.Schema(context.Background(), action.SchemaRequest{}, schemaResp)
	assert.False(t, schemaResp.Diagnostics.HasError())

	attrs := schemaResp.Schema.Attributes
	for _, name := range []string{"zone_id", "role_configuration_id", "target_type", "target_uin"} {
		assert.Contains(t, attrs, name)
		assert.True(t, attrs[name].IsRequired())
	}

	assert.Contains(t, schemaResp.Schema.Blocks, "timeouts")

	uinAttr, ok := attrs["target_uin"].(action_schema.Int64Attribute)
	assert.True(t, ok)
	assert.Equal(t, types.Int64Type, uinAttr.GetType())
}
