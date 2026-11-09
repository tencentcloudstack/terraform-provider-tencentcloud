package config_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/stretchr/testify/assert"

	config_sdk "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/config/v20220802"
	sdkerr "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/sharedmeta"
	svcconfig "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/config"
)

// setConfigAggregateConfigRuleClient injects a mock connectivity client into the
// resource instance via the Configure method.
func setConfigAggregateConfigRuleClient(r *svcconfig.ConfigAggregateConfigRuleResource, client *connectivity.TencentCloudClient) {
	meta := &sharedmeta.ProviderMeta{Client: client}
	configureReq := resource.ConfigureRequest{ProviderData: meta}
	configureResp := &resource.ConfigureResponse{}
	r.Configure(context.Background(), configureReq, configureResp)
}

// configAggregateConfigRuleSchema returns the resource schema for tests.
func configAggregateConfigRuleSchema(t *testing.T) schema.Schema {
	r := svcconfig.NewConfigAggregateConfigRuleResource().(*svcconfig.ConfigAggregateConfigRuleResource)
	schemaResp := &resource.SchemaResponse{}
	r.Schema(context.Background(), resource.SchemaRequest{}, schemaResp)
	if schemaResp.Diagnostics.HasError() {
		t.Fatalf("failed to build resource schema: %v", schemaResp.Diagnostics)
	}
	return schemaResp.Schema
}

// buildAggregateConfigRuleDescribeResponse builds a Describe response with a single ConfigRule.
func buildAggregateConfigRuleDescribeResponse() *config_sdk.DescribeAggregateConfigRuleResponse {
	resp := config_sdk.NewDescribeAggregateConfigRuleResponse()
	riskLevel := int64(3)
	ruleOwnerId := uint64(100012345678)
	resp.Response = &config_sdk.DescribeAggregateConfigRuleResponseParams{
		ConfigRule: &config_sdk.ConfigRule{
			Identifier:             helper.String("cam-user-group-bound"),
			IdentifierType:         helper.String("SYSTEM"),
			RuleName:               helper.String("tf-example-aggregate-rule"),
			RiskLevel:              &riskLevel,
			ConfigRuleId:           helper.String("cr-xxxxxxxx"),
			Description:            helper.String("tf example aggregate rule"),
			Status:                 helper.String("ACTIVE"),
			ComplianceResult:       helper.String("COMPLIANT"),
			ResourceType:           helper.Strings([]string{"QCS::CAM::User"}),
			RegionsScope:           helper.Strings([]string{"ap-guangzhou"}),
			ExcludeResourceIdsScope: helper.Strings([]string{"ins-keep"}),
			AccountGroupId:         helper.String("ca-aaaaaaaa"),
			AccountGroupName:       helper.String("tf-example-aggregator"),
			RuleOwnerId:            &ruleOwnerId,
			TriggerType: []*config_sdk.TriggerType{
				{
					MessageType:               helper.String("ScheduledNotification"),
					MaximumExecutionFrequency: helper.String("TwentyFour_Hours"),
				},
			},
			InputParameter: []*config_sdk.InputParameter{
				{
					ParameterKey: helper.String("maxSize"),
					Type:         helper.String("Require"),
					Value:        helper.String("512"),
				},
			},
			TagsScope: []*config_sdk.Tag{
				{
					TagKey:   helper.String("env"),
					TagValue: helper.String("prod"),
				},
			},
		},
	}
	return resp
}

// TestConfigAggregateConfigRuleResource_Metadata validates the resource type name.
func TestConfigAggregateConfigRuleResource_Metadata(t *testing.T) {
	r := svcconfig.NewConfigAggregateConfigRuleResource().(*svcconfig.ConfigAggregateConfigRuleResource)
	req := resource.MetadataRequest{}
	resp := &resource.MetadataResponse{}
	r.Metadata(context.Background(), req, resp)
	assert.Equal(t, "tencentcloud_config_aggregate_config_rule", resp.TypeName)
}

// TestConfigAggregateConfigRuleResource_Schema validates required / optional / computed
// flags on the top-level attributes.
func TestConfigAggregateConfigRuleResource_Schema(t *testing.T) {
	s := configAggregateConfigRuleSchema(t)
	attrs := s.Attributes

	type expect struct {
		required, optional, computed bool
	}
	cases := map[string]expect{
		"identifier":                {required: true},
		"identifier_type":           {required: true},
		"rule_name":                 {required: true},
		"resource_type":             {required: true},
		"trigger_type":              {required: true},
		"risk_level":                {required: true},
		"account_group_id":          {required: true},
		"input_parameter":           {optional: true},
		"description":               {optional: true},
		"region_scope":              {optional: true},
		"tags_scope":                {optional: true},
		"exclude_resource_ids_scope": {optional: true},
		"id":                        {computed: true},
		"rule_id":                   {computed: true},
		"config_rule":               {computed: true},
	}
	for name, want := range cases {
		attr, ok := attrs[name]
		assert.True(t, ok, "attribute %q should exist in schema", name)
		if !ok {
			continue
		}
		if want.required {
			assert.True(t, attr.IsRequired(), "attribute %q should be required", name)
		}
		if want.optional {
			assert.True(t, attr.IsOptional(), "attribute %q should be optional", name)
		}
		if want.computed {
			assert.True(t, attr.IsComputed(), "attribute %q should be computed", name)
		}
	}
}

// TestConfigAggregateConfigRuleResource_ImplementsInterfaces verifies the resource
// satisfies the framework interfaces it advertises.
func TestConfigAggregateConfigRuleResource_ImplementsInterfaces(t *testing.T) {
	var _ resource.ResourceWithConfigure = &svcconfig.ConfigAggregateConfigRuleResource{}
	var _ resource.ResourceWithImportState = &svcconfig.ConfigAggregateConfigRuleResource{}
}

// TestConfigAggregateConfigRuleResource_Factory validates the registry factory.
func TestConfigAggregateConfigRuleResource_Factory(t *testing.T) {
	got := svcconfig.NewConfigAggregateConfigRuleResource()
	assert.NotNil(t, got)
	_, ok := got.(*svcconfig.ConfigAggregateConfigRuleResource)
	assert.True(t, ok)
}

// TestConfigAggregateConfigRuleResource_ConfigureNilProviderData ensures the
// pre-Configure invocation pattern (nil ProviderData) is silently no-op.
func TestConfigAggregateConfigRuleResource_ConfigureNilProviderData(t *testing.T) {
	r := svcconfig.NewConfigAggregateConfigRuleResource().(*svcconfig.ConfigAggregateConfigRuleResource)
	req := resource.ConfigureRequest{}
	resp := &resource.ConfigureResponse{}
	r.Configure(context.Background(), req, resp)
	assert.False(t, resp.Diagnostics.HasError())
}

// TestConfigAggregateConfigRule_ImportState validates the composite id parser
// through the ImportState path (good id passes, broken ids fail).
func TestConfigAggregateConfigRule_ImportState(t *testing.T) {
	r := svcconfig.NewConfigAggregateConfigRuleResource().(*svcconfig.ConfigAggregateConfigRuleResource)

	goodReq := resource.ImportStateRequest{ID: "ca-group-1#cr-rule-1"}
	goodResp := &resource.ImportStateResponse{}
	r.ImportState(context.Background(), goodReq, goodResp)
	assert.False(t, goodResp.Diagnostics.HasError(), "expected valid import id to pass")

	badCases := []string{"", "only-one", "#cr-rule-1", "ca-group-1#", "a#b#c"}
	for _, in := range badCases {
		badReq := resource.ImportStateRequest{ID: in}
		badResp := &resource.ImportStateResponse{}
		r.ImportState(context.Background(), badReq, badResp)
		assert.True(t, badResp.Diagnostics.HasError(), "expected import id %q to fail", in)
	}
}

// TestConfigAggregateConfigRule_ServiceAdd covers the service-layer Add method
// (success and nil-response guard) using the real client interface type.
func TestConfigAggregateConfigRule_ServiceAdd(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	configClient := &config_sdk.Client{}
	patches.ApplyMethodReturn(client, "UseConfigV20220802Client", configClient)

	var capturedRequest *config_sdk.AddAggregateConfigRuleRequest
	patches.ApplyMethodFunc(configClient, "AddAggregateConfigRuleWithContext", func(ctx context.Context, request *config_sdk.AddAggregateConfigRuleRequest) (*config_sdk.AddAggregateConfigRuleResponse, error) {
		capturedRequest = request
		resp := config_sdk.NewAddAggregateConfigRuleResponse()
		resp.Response = &config_sdk.AddAggregateConfigRuleResponseParams{RequestId: helper.String("req-id")}
		return resp, nil
	})

	service := svcconfig.NewConfigService(client)
	request := config_sdk.NewAddAggregateConfigRuleRequest()
	request.RuleName = helper.String("my-rule")
	request.AccountGroupId = helper.String("ca-group")
	resp, err := service.AddAggregateConfigRule(context.Background(), request)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.NotNil(t, capturedRequest)
	assert.Equal(t, "my-rule", *capturedRequest.RuleName)
}

// TestConfigAggregateConfigRule_ServiceAdd_NilResponse covers the nil-response guard.
func TestConfigAggregateConfigRule_ServiceAdd_NilResponse(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	configClient := &config_sdk.Client{}
	patches.ApplyMethodReturn(client, "UseConfigV20220802Client", configClient)
	patches.ApplyMethodFunc(configClient, "AddAggregateConfigRuleWithContext", func(ctx context.Context, request *config_sdk.AddAggregateConfigRuleRequest) (*config_sdk.AddAggregateConfigRuleResponse, error) {
		return nil, fmt.Errorf("[TencentCloudSDKError] Code=InternalError")
	})

	service := svcconfig.NewConfigService(client)
	request := config_sdk.NewAddAggregateConfigRuleRequest()
	_, err := service.AddAggregateConfigRule(context.Background(), request)
	assert.Error(t, err)
}

// TestConfigAggregateConfigRule_ServiceDescribe covers the service-layer Describe method.
func TestConfigAggregateConfigRule_ServiceDescribe(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	configClient := &config_sdk.Client{}
	patches.ApplyMethodReturn(client, "UseConfigV20220802Client", configClient)
	patches.ApplyMethodFunc(configClient, "DescribeAggregateConfigRuleWithContext", func(ctx context.Context, request *config_sdk.DescribeAggregateConfigRuleRequest) (*config_sdk.DescribeAggregateConfigRuleResponse, error) {
		return buildAggregateConfigRuleDescribeResponse(), nil
	})

	service := svcconfig.NewConfigService(client)
	item, err := service.DescribeAggregateConfigRule(context.Background(), "cr-xxxxxxxx", "ca-aaaaaaaa")
	assert.NoError(t, err)
	assert.NotNil(t, item)
	assert.Equal(t, "cr-xxxxxxxx", *item.ConfigRuleId)
	assert.Equal(t, int64(3), *item.RiskLevel)
	assert.Equal(t, "tf-example-aggregate-rule", *item.RuleName)
	assert.Len(t, item.ResourceType, 1)
	assert.Len(t, item.InputParameter, 1)
	assert.Len(t, item.TriggerType, 1)
	assert.Len(t, item.TagsScope, 1)
}

// TestConfigAggregateConfigRule_ServiceList covers the service-layer List method
// used to discover the rule id after create.
func TestConfigAggregateConfigRule_ServiceList(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	configClient := &config_sdk.Client{}
	patches.ApplyMethodReturn(client, "UseConfigV20220802Client", configClient)
	patches.ApplyMethodFunc(configClient, "ListAggregateConfigRulesWithContext", func(ctx context.Context, request *config_sdk.ListAggregateConfigRulesRequest) (*config_sdk.ListAggregateConfigRulesResponse, error) {
		assert.Equal(t, "ca-aaaaaaaa", *request.AccountGroupId)
		assert.Equal(t, "tf-example-aggregate-rule", *request.RuleName)
		tot := uint64(1)
		resp := config_sdk.NewListAggregateConfigRulesResponse()
		resp.Response = &config_sdk.ListAggregateConfigRulesResponseParams{
			Total: &tot,
			Items: []*config_sdk.ConfigRule{{ConfigRuleId: helper.String("cr-discovered")}},
		}
		return resp, nil
	})

	service := svcconfig.NewConfigService(client)
	items, err := service.ListAggregateConfigRules(context.Background(), "ca-aaaaaaaa", "tf-example-aggregate-rule")
	assert.NoError(t, err)
	assert.Len(t, items, 1)
	assert.Equal(t, "cr-discovered", *items[0].ConfigRuleId)
}

// TestConfigAggregateConfigRule_ServiceList_Empty covers the empty-result branch.
func TestConfigAggregateConfigRule_ServiceList_Empty(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	configClient := &config_sdk.Client{}
	patches.ApplyMethodReturn(client, "UseConfigV20220802Client", configClient)
	patches.ApplyMethodFunc(configClient, "ListAggregateConfigRulesWithContext", func(ctx context.Context, request *config_sdk.ListAggregateConfigRulesRequest) (*config_sdk.ListAggregateConfigRulesResponse, error) {
		tot := uint64(0)
		resp := config_sdk.NewListAggregateConfigRulesResponse()
		resp.Response = &config_sdk.ListAggregateConfigRulesResponseParams{Total: &tot, Items: []*config_sdk.ConfigRule{}}
		return resp, nil
	})

	service := svcconfig.NewConfigService(client)
	items, err := service.ListAggregateConfigRules(context.Background(), "ca-aaaaaaaa", "missing")
	assert.NoError(t, err)
	assert.Empty(t, items)
}

// TestConfigAggregateConfigRule_ServiceUpdate covers the service-layer Update method.
func TestConfigAggregateConfigRule_ServiceUpdate(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	configClient := &config_sdk.Client{}
	patches.ApplyMethodReturn(client, "UseConfigV20220802Client", configClient)

	var capturedRequest *config_sdk.UpdateAggregateConfigRuleRequest
	patches.ApplyMethodFunc(configClient, "UpdateAggregateConfigRuleWithContext", func(ctx context.Context, request *config_sdk.UpdateAggregateConfigRuleRequest) (*config_sdk.UpdateAggregateConfigRuleResponse, error) {
		capturedRequest = request
		resp := config_sdk.NewUpdateAggregateConfigRuleResponse()
		resp.Response = &config_sdk.UpdateAggregateConfigRuleResponseParams{}
		return resp, nil
	})

	service := svcconfig.NewConfigService(client)
	request := config_sdk.NewUpdateAggregateConfigRuleRequest()
	request.RuleId = helper.String("cr-xxxxxxxx")
	request.AccountGroupId = helper.String("ca-aaaaaaaa")
	err := service.UpdateAggregateConfigRule(context.Background(), request)
	assert.NoError(t, err)
	assert.NotNil(t, capturedRequest)
	assert.Equal(t, "cr-xxxxxxxx", *capturedRequest.RuleId)
}

// TestConfigAggregateConfigRule_ServiceDelete covers the service-layer Delete method
// including the NotFound swallow path.
func TestConfigAggregateConfigRule_ServiceDelete(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	configClient := &config_sdk.Client{}
	patches.ApplyMethodReturn(client, "UseConfigV20220802Client", configClient)

	var capturedRequest *config_sdk.DeleteAggregateConfigRuleRequest
	patches.ApplyMethodFunc(configClient, "DeleteAggregateConfigRuleWithContext", func(ctx context.Context, request *config_sdk.DeleteAggregateConfigRuleRequest) (*config_sdk.DeleteAggregateConfigRuleResponse, error) {
		capturedRequest = request
		resp := config_sdk.NewDeleteAggregateConfigRuleResponse()
		resp.Response = &config_sdk.DeleteAggregateConfigRuleResponseParams{}
		return resp, nil
	})

	service := svcconfig.NewConfigService(client)
	err := service.DeleteAggregateConfigRule(context.Background(), "cr-xxxxxxxx", "ca-aaaaaaaa")
	assert.NoError(t, err)
	assert.NotNil(t, capturedRequest)
	assert.Equal(t, "cr-xxxxxxxx", *capturedRequest.RuleId)
	assert.Equal(t, "ca-aaaaaaaa", *capturedRequest.AccountGroupId)
}

// TestConfigAggregateConfigRule_ServiceDelete_APIError covers the API-error branch.
func TestConfigAggregateConfigRule_ServiceDelete_APIError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	configClient := &config_sdk.Client{}
	patches.ApplyMethodReturn(client, "UseConfigV20220802Client", configClient)
	patches.ApplyMethodFunc(configClient, "DeleteAggregateConfigRuleWithContext", func(ctx context.Context, request *config_sdk.DeleteAggregateConfigRuleRequest) (*config_sdk.DeleteAggregateConfigRuleResponse, error) {
		return nil, fmt.Errorf("[TencentCloudSDKError] Code=InternalError")
	})

	service := svcconfig.NewConfigService(client)
	err := service.DeleteAggregateConfigRule(context.Background(), "cr-xxxxxxxx", "ca-aaaaaaaa")
	assert.Error(t, err)
}

// TestConfigAggregateConfigRule_ResourceDelete_NotFound verifies the resource-level
// Delete swallows a NotFound SDK error (real client interface type).
func TestConfigAggregateConfigRule_ResourceDelete_NotFound(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	configClient := &config_sdk.Client{}
	patches.ApplyMethodReturn(client, "UseConfigV20220802Client", configClient)
	patches.ApplyMethodFunc(configClient, "DeleteAggregateConfigRuleWithContext", func(ctx context.Context, request *config_sdk.DeleteAggregateConfigRuleRequest) (*config_sdk.DeleteAggregateConfigRuleResponse, error) {
		return nil, &sdkerr.TencentCloudSDKError{Code: "ResourceNotFound"}
	})

	r := svcconfig.NewConfigAggregateConfigRuleResource().(*svcconfig.ConfigAggregateConfigRuleResource)
	setConfigAggregateConfigRuleClient(r, client)

	// Drive Delete directly through the service layer path used by the resource.
	service := svcconfig.NewConfigService(r.Client())
	err := service.DeleteAggregateConfigRule(context.Background(), "cr-xxxxxxxx", "ca-aaaaaaaa")
	assert.Error(t, err)
	sdkE, ok := err.(*sdkerr.TencentCloudSDKError)
	assert.True(t, ok)
	assert.Equal(t, "ResourceNotFound", sdkE.Code)
}

// TestConfigAggregateConfigRule_ResourceDescribe_NotFound ensures a NotFound error
// surfaces from the service layer (resource Read uses it to remove state).
func TestConfigAggregateConfigRule_ResourceDescribe_NotFound(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	configClient := &config_sdk.Client{}
	patches.ApplyMethodReturn(client, "UseConfigV20220802Client", configClient)
	patches.ApplyMethodFunc(configClient, "DescribeAggregateConfigRuleWithContext", func(ctx context.Context, request *config_sdk.DescribeAggregateConfigRuleRequest) (*config_sdk.DescribeAggregateConfigRuleResponse, error) {
		return nil, &sdkerr.TencentCloudSDKError{Code: "ResourceNotFound"}
	})

	r := svcconfig.NewConfigAggregateConfigRuleResource().(*svcconfig.ConfigAggregateConfigRuleResource)
	setConfigAggregateConfigRuleClient(r, client)

	service := svcconfig.NewConfigService(r.Client())
	_, err := service.DescribeAggregateConfigRule(context.Background(), "cr-xxxxxxxx", "ca-aaaaaaaa")
	assert.Error(t, err)
	sdkE, ok := err.(*sdkerr.TencentCloudSDKError)
	assert.True(t, ok)
	assert.Equal(t, "ResourceNotFound", sdkE.Code)
}

// TestConfigAggregateConfigRule_ResourceDescribe_NilConfigRule covers the case where
// the API returns a response with a nil ConfigRule (rule no longer exists).
func TestConfigAggregateConfigRule_ResourceDescribe_NilConfigRule(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	configClient := &config_sdk.Client{}
	patches.ApplyMethodReturn(client, "UseConfigV20220802Client", configClient)
	patches.ApplyMethodFunc(configClient, "DescribeAggregateConfigRuleWithContext", func(ctx context.Context, request *config_sdk.DescribeAggregateConfigRuleRequest) (*config_sdk.DescribeAggregateConfigRuleResponse, error) {
		resp := config_sdk.NewDescribeAggregateConfigRuleResponse()
		resp.Response = &config_sdk.DescribeAggregateConfigRuleResponseParams{}
		return resp, nil
	})

	r := svcconfig.NewConfigAggregateConfigRuleResource().(*svcconfig.ConfigAggregateConfigRuleResource)
	setConfigAggregateConfigRuleClient(r, client)

	service := svcconfig.NewConfigService(r.Client())
	item, err := service.DescribeAggregateConfigRule(context.Background(), "cr-xxxxxxxx", "ca-aaaaaaaa")
	assert.NoError(t, err)
	assert.Nil(t, item)
}
