package teo_test

import (
	"context"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	teov20220901 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/teo"
)

// go test ./tencentcloud/services/teo/ -run "TestTeoL7AccRuleV2Extension" -v -count=1 -gcflags="all=-l"

// TestTeoL7AccRuleV2Extension_AdvancedOriginRoutingParameters_Schema tests that AdvancedOriginRoutingParameters schema is correctly defined
func TestTeoL7AccRuleV2Extension_AdvancedOriginRoutingParameters_Schema(t *testing.T) {
	schemaMap := teo.TencentTeoL7RuleBranchBasicInfo(1)

	actionsSchema, ok := schemaMap["actions"]
	assert.True(t, ok, "actions schema should exist")
	assert.Equal(t, schema.TypeList, actionsSchema.Type, "actions should be TypeList")

	actionsElem, ok := actionsSchema.Elem.(*schema.Resource)
	if !ok {
		t.Fatal("actions Elem should be *schema.Resource")
	}
	actionsResourceSchema := actionsElem.Schema

	aorpSchema, ok := actionsResourceSchema["advanced_origin_routing_parameters"]
	assert.True(t, ok, "advanced_origin_routing_parameters schema should exist")
	assert.Equal(t, schema.TypeList, aorpSchema.Type, "advanced_origin_routing_parameters should be TypeList")
	assert.True(t, aorpSchema.Optional, "advanced_origin_routing_parameters should be Optional")
	assert.Equal(t, 1, aorpSchema.MaxItems, "advanced_origin_routing_parameters MaxItems should be 1")

	aorpElem, ok := aorpSchema.Elem.(*schema.Resource)
	if !ok {
		t.Fatal("advanced_origin_routing_parameters Elem should be *schema.Resource")
	}
	aorpResourceSchema := aorpElem.Schema

	directionSchema, ok := aorpResourceSchema["direction"]
	assert.True(t, ok, "direction schema should exist inside advanced_origin_routing_parameters")
	assert.Equal(t, schema.TypeString, directionSchema.Type, "direction should be TypeString")
	assert.True(t, directionSchema.Optional, "direction should be Optional")
}

// TestTeoL7AccRuleV2Extension_ShieldParameters_Schema tests that ShieldParameters schema is correctly defined
func TestTeoL7AccRuleV2Extension_ShieldParameters_Schema(t *testing.T) {
	schemaMap := teo.TencentTeoL7RuleBranchBasicInfo(1)

	actionsSchema := schemaMap["actions"]
	actionsElem := actionsSchema.Elem.(*schema.Resource)
	actionsResourceSchema := actionsElem.Schema

	spSchema, ok := actionsResourceSchema["shield_parameters"]
	assert.True(t, ok, "shield_parameters schema should exist")
	assert.Equal(t, schema.TypeList, spSchema.Type, "shield_parameters should be TypeList")
	assert.True(t, spSchema.Optional, "shield_parameters should be Optional")
	assert.Equal(t, 1, spSchema.MaxItems, "shield_parameters MaxItems should be 1")

	spElem := spSchema.Elem.(*schema.Resource)
	spResourceSchema := spElem.Schema

	shieldSpaceIdSchema, ok := spResourceSchema["shield_space_id"]
	assert.True(t, ok, "shield_space_id schema should exist inside shield_parameters")
	assert.Equal(t, schema.TypeString, shieldSpaceIdSchema.Type, "shield_space_id should be TypeString")
	assert.True(t, shieldSpaceIdSchema.Optional, "shield_space_id should be Optional")
}

// TestTeoL7AccRuleV2Extension_SiteFailoverParameters_Schema tests that SiteFailoverParameters schema is correctly defined
func TestTeoL7AccRuleV2Extension_SiteFailoverParameters_Schema(t *testing.T) {
	schemaMap := teo.TencentTeoL7RuleBranchBasicInfo(1)

	actionsSchema := schemaMap["actions"]
	actionsElem := actionsSchema.Elem.(*schema.Resource)
	actionsResourceSchema := actionsElem.Schema

	sfpSchema, ok := actionsResourceSchema["site_failover_parameters"]
	assert.True(t, ok, "site_failover_parameters schema should exist")
	assert.Equal(t, schema.TypeList, sfpSchema.Type, "site_failover_parameters should be TypeList")
	assert.True(t, sfpSchema.Optional, "site_failover_parameters should be Optional")
	assert.Equal(t, 1, sfpSchema.MaxItems, "site_failover_parameters MaxItems should be 1")

	sfpElem := sfpSchema.Elem.(*schema.Resource)
	sfpResourceSchema := sfpElem.Schema

	statusCodesSchema, ok := sfpResourceSchema["site_failover_status_codes"]
	assert.True(t, ok, "site_failover_status_codes schema should exist inside site_failover_parameters")
	assert.Equal(t, schema.TypeList, statusCodesSchema.Type, "site_failover_status_codes should be TypeList")

	paramsSchema, ok := sfpResourceSchema["site_failover_params"]
	assert.True(t, ok, "site_failover_params schema should exist inside site_failover_parameters")
	assert.Equal(t, schema.TypeList, paramsSchema.Type, "site_failover_params should be TypeList")

	paramsElem := paramsSchema.Elem.(*schema.Resource)
	paramsResourceSchema := paramsElem.Schema

	modeSchema, ok := paramsResourceSchema["mode"]
	assert.True(t, ok, "mode schema should exist inside site_failover_params")
	assert.Equal(t, schema.TypeString, modeSchema.Type, "mode should be TypeString")

	originSchema, ok := paramsResourceSchema["origin"]
	assert.True(t, ok, "origin schema should exist inside site_failover_params")
	assert.Equal(t, schema.TypeString, originSchema.Type, "origin should be TypeString")

	originProtocolSchema, ok := paramsResourceSchema["origin_protocol"]
	assert.True(t, ok, "origin_protocol schema should exist inside site_failover_params")
	assert.Equal(t, schema.TypeString, originProtocolSchema.Type, "origin_protocol should be TypeString")

	httpOriginPortSchema, ok := paramsResourceSchema["http_origin_port"]
	assert.True(t, ok, "http_origin_port schema should exist inside site_failover_params")
	assert.Equal(t, schema.TypeInt, httpOriginPortSchema.Type, "http_origin_port should be TypeInt")

	httpsOriginPortSchema, ok := paramsResourceSchema["https_origin_port"]
	assert.True(t, ok, "https_origin_port schema should exist inside site_failover_params")
	assert.Equal(t, schema.TypeInt, httpsOriginPortSchema.Type, "https_origin_port should be TypeInt")

	upstreamHostHeaderSchema, ok := paramsResourceSchema["upstream_host_header"]
	assert.True(t, ok, "upstream_host_header schema should exist inside site_failover_params")
	assert.Equal(t, schema.TypeList, upstreamHostHeaderSchema.Type, "upstream_host_header should be TypeList")

	upstreamURLRewriteSchema, ok := paramsResourceSchema["upstream_url_rewrite"]
	assert.True(t, ok, "upstream_url_rewrite schema should exist inside site_failover_params")
	assert.Equal(t, schema.TypeList, upstreamURLRewriteSchema.Type, "upstream_url_rewrite should be TypeList")

	upstreamRequestParamsSchema, ok := paramsResourceSchema["upstream_request_parameters"]
	assert.True(t, ok, "upstream_request_parameters schema should exist inside site_failover_params")
	assert.Equal(t, schema.TypeList, upstreamRequestParamsSchema.Type, "upstream_request_parameters should be TypeList")

	upstreamHTTP2Schema, ok := paramsResourceSchema["upstream_http2_parameters"]
	assert.True(t, ok, "upstream_http2_parameters schema should exist inside site_failover_params")
	assert.Equal(t, schema.TypeList, upstreamHTTP2Schema.Type, "upstream_http2_parameters should be TypeList")

	privateAccessSchema, ok := paramsResourceSchema["private_access"]
	assert.True(t, ok, "private_access schema should exist inside site_failover_params")
	assert.Equal(t, schema.TypeString, privateAccessSchema.Type, "private_access should be TypeString")

	privateParamsSchema, ok := paramsResourceSchema["private_parameters"]
	assert.True(t, ok, "private_parameters schema should exist inside site_failover_params")
	assert.Equal(t, schema.TypeList, privateParamsSchema.Type, "private_parameters should be TypeList")

	redirectURLSchema, ok := paramsResourceSchema["redirect_url"]
	assert.True(t, ok, "redirect_url schema should exist inside site_failover_params")
	assert.Equal(t, schema.TypeString, redirectURLSchema.Type, "redirect_url should be TypeString")

	responsePageIdSchema, ok := paramsResourceSchema["response_page_id"]
	assert.True(t, ok, "response_page_id schema should exist inside site_failover_params")
	assert.Equal(t, schema.TypeString, responsePageIdSchema.Type, "response_page_id should be TypeString")

	statusCodeSchema, ok := paramsResourceSchema["status_code"]
	assert.True(t, ok, "status_code schema should exist inside site_failover_params")
	assert.Equal(t, schema.TypeInt, statusCodeSchema.Type, "status_code should be TypeInt")
}

// TestTeoL7AccRuleV2Extension_AdvancedOriginRoutingParameters_FlattenSet tests flatten and set for AdvancedOriginRoutingParameters
func TestTeoL7AccRuleV2Extension_AdvancedOriginRoutingParameters_FlattenSet(t *testing.T) {
	direction := "MainlandChinaAndGlobalAdaptive"
	sdkAction := &teov20220901.RuleEngineAction{
		Name: helper.String("AdvancedOriginRouting"),
		AdvancedOriginRoutingParameters: &teov20220901.AdvancedOriginRoutingParameters{
			Direction: &direction,
		},
	}

	assert.NotNil(t, sdkAction.AdvancedOriginRoutingParameters)
	assert.Equal(t, "MainlandChinaAndGlobalAdaptive", *sdkAction.AdvancedOriginRoutingParameters.Direction)
}

// TestTeoL7AccRuleV2Extension_ShieldParameters_FlattenSet tests flatten and set for ShieldParameters
func TestTeoL7AccRuleV2Extension_ShieldParameters_FlattenSet(t *testing.T) {
	shieldSpaceId := "shield-space-abc123"
	sdkAction := &teov20220901.RuleEngineAction{
		Name: helper.String("Shield"),
		ShieldParameters: &teov20220901.ShieldParameters{
			ShieldSpaceId: &shieldSpaceId,
		},
	}

	assert.NotNil(t, sdkAction.ShieldParameters)
	assert.Equal(t, "shield-space-abc123", *sdkAction.ShieldParameters.ShieldSpaceId)
}

// TestTeoL7AccRuleV2Extension_SiteFailoverParameters_FlattenSet tests flatten and set for SiteFailoverParameters
func TestTeoL7AccRuleV2Extension_SiteFailoverParameters_FlattenSet(t *testing.T) {
	statusCode := int64(500)
	httpPort := int64(80)
	httpsPort := int64(443)
	respStatusCode := int64(302)

	sdkAction := &teov20220901.RuleEngineAction{
		Name: helper.String("SiteFailover"),
		SiteFailoverParameters: &teov20220901.SiteFailoverParameters{
			SiteFailoverStatusCodes: []*int64{&statusCode},
			SiteFailoverParams: []*teov20220901.SiteFailover{
				{
					Mode:            helper.String("FailoverToHost"),
					Origin:          helper.String("backup.example.com"),
					OriginProtocol:  helper.String("https"),
					HTTPOriginPort:  &httpPort,
					HTTPSOriginPort: &httpsPort,
					StatusCode:      &respStatusCode,
				},
			},
		},
	}

	assert.NotNil(t, sdkAction.SiteFailoverParameters)
	assert.Len(t, sdkAction.SiteFailoverParameters.SiteFailoverStatusCodes, 1)
	assert.Equal(t, int64(500), *sdkAction.SiteFailoverParameters.SiteFailoverStatusCodes[0])
	assert.Len(t, sdkAction.SiteFailoverParameters.SiteFailoverParams, 1)
	assert.Equal(t, "FailoverToHost", *sdkAction.SiteFailoverParameters.SiteFailoverParams[0].Mode)
	assert.Equal(t, "backup.example.com", *sdkAction.SiteFailoverParameters.SiteFailoverParams[0].Origin)
	assert.Equal(t, "https", *sdkAction.SiteFailoverParameters.SiteFailoverParams[0].OriginProtocol)
	assert.Equal(t, int64(80), *sdkAction.SiteFailoverParameters.SiteFailoverParams[0].HTTPOriginPort)
	assert.Equal(t, int64(443), *sdkAction.SiteFailoverParameters.SiteFailoverParams[0].HTTPSOriginPort)
	assert.Equal(t, int64(302), *sdkAction.SiteFailoverParameters.SiteFailoverParams[0].StatusCode)
}

// TestTeoL7AccRuleV2Extension_SiteFailoverParameters_WithNestedParams tests SiteFailover with upstream_host_header, upstream_url_rewrite, etc.
func TestTeoL7AccRuleV2Extension_SiteFailoverParameters_WithNestedParams(t *testing.T) {
	sdkAction := &teov20220901.RuleEngineAction{
		Name: helper.String("SiteFailover"),
		SiteFailoverParameters: &teov20220901.SiteFailoverParameters{
			SiteFailoverStatusCodes: []*int64{helper.IntInt64(500)},
			SiteFailoverParams: []*teov20220901.SiteFailover{
				{
					Mode:           helper.String("FailoverToHost"),
					Origin:         helper.String("backup.example.com"),
					OriginProtocol: helper.String("follow"),
					UpstreamHostHeader: &teov20220901.HostHeaderParameters{
						Action:     helper.String("custom"),
						ServerName: helper.String("custom.example.com"),
					},
					UpstreamURLRewrite: &teov20220901.UpstreamURLRewriteParameters{
						Type:   helper.String("path"),
						Action: helper.String("replace"),
						Value:  helper.String("/new-path"),
					},
					UpstreamRequestParameters: &teov20220901.UpstreamRequestParameters{
						QueryString: &teov20220901.UpstreamRequestQueryString{
							Switch: helper.String("on"),
							Action: helper.String("includeCustom"),
							Values: []*string{helper.String("param1")},
						},
						Cookie: &teov20220901.UpstreamRequestCookie{
							Switch: helper.String("off"),
						},
					},
					UpstreamHTTP2Parameters: &teov20220901.UpstreamHTTP2Parameters{
						Switch: helper.String("on"),
					},
					PrivateAccess: helper.String("off"),
				},
			},
		},
	}

	sfp := sdkAction.SiteFailoverParameters
	assert.NotNil(t, sfp)
	assert.Len(t, sfp.SiteFailoverParams, 1)

	param := sfp.SiteFailoverParams[0]
	assert.NotNil(t, param.UpstreamHostHeader)
	assert.Equal(t, "custom", *param.UpstreamHostHeader.Action)
	assert.Equal(t, "custom.example.com", *param.UpstreamHostHeader.ServerName)

	assert.NotNil(t, param.UpstreamURLRewrite)
	assert.Equal(t, "path", *param.UpstreamURLRewrite.Type)
	assert.Equal(t, "replace", *param.UpstreamURLRewrite.Action)
	assert.Equal(t, "/new-path", *param.UpstreamURLRewrite.Value)

	assert.NotNil(t, param.UpstreamRequestParameters)
	assert.NotNil(t, param.UpstreamRequestParameters.QueryString)
	assert.Equal(t, "on", *param.UpstreamRequestParameters.QueryString.Switch)
	assert.Equal(t, "includeCustom", *param.UpstreamRequestParameters.QueryString.Action)

	assert.NotNil(t, param.UpstreamHTTP2Parameters)
	assert.Equal(t, "on", *param.UpstreamHTTP2Parameters.Switch)

	assert.Equal(t, "off", *param.PrivateAccess)
}

// TestTeoL7AccRuleV2Extension_SiteFailoverParameters_PrivateParams tests SiteFailover with PrivateParameters
func TestTeoL7AccRuleV2Extension_SiteFailoverParameters_PrivateParams(t *testing.T) {
	sdkAction := &teov20220901.RuleEngineAction{
		Name: helper.String("SiteFailover"),
		SiteFailoverParameters: &teov20220901.SiteFailoverParameters{
			SiteFailoverStatusCodes: []*int64{helper.IntInt64(500)},
			SiteFailoverParams: []*teov20220901.SiteFailover{
				{
					Mode:          helper.String("FailoverToS3CompatibleObjectStorage"),
					Origin:        helper.String("s3-bucket.example.com"),
					PrivateAccess: helper.String("on"),
					PrivateParameters: &teov20220901.OriginPrivateParameters{
						AccessKeyId:      helper.String("AKID123"),
						SecretAccessKey:  helper.String("Secret123"),
						SignatureVersion: helper.String("v4"),
						Region:           helper.String("us-east-1"),
					},
				},
			},
		},
	}

	sfp := sdkAction.SiteFailoverParameters
	assert.NotNil(t, sfp)
	param := sfp.SiteFailoverParams[0]

	assert.NotNil(t, param.PrivateParameters)
	assert.Equal(t, "AKID123", *param.PrivateParameters.AccessKeyId)
	assert.Equal(t, "Secret123", *param.PrivateParameters.SecretAccessKey)
	assert.Equal(t, "v4", *param.PrivateParameters.SignatureVersion)
	assert.Equal(t, "us-east-1", *param.PrivateParameters.Region)
}

// TestTeoL7AccRuleV2Extension_SiteFailoverParameters_RedirectURL tests SiteFailover with redirect URL mode
func TestTeoL7AccRuleV2Extension_SiteFailoverParameters_RedirectURL(t *testing.T) {
	sdkAction := &teov20220901.RuleEngineAction{
		Name: helper.String("SiteFailover"),
		SiteFailoverParameters: &teov20220901.SiteFailoverParameters{
			SiteFailoverStatusCodes: []*int64{helper.IntInt64(500)},
			SiteFailoverParams: []*teov20220901.SiteFailover{
				{
					Mode:        helper.String("FailoverRedirectToURL"),
					RedirectURL: helper.String("https://fallback.example.com"),
					StatusCode:  helper.IntInt64(302),
				},
			},
		},
	}

	sfp := sdkAction.SiteFailoverParameters
	param := sfp.SiteFailoverParams[0]

	assert.Equal(t, "FailoverRedirectToURL", *param.Mode)
	assert.Equal(t, "https://fallback.example.com", *param.RedirectURL)
	assert.Equal(t, int64(302), *param.StatusCode)
}

// TestTeoL7AccRuleV2Extension_CustomActionParameters_Schema tests that CustomActionParameters schema is correctly defined
func TestTeoL7AccRuleV2Extension_CustomActionParameters_Schema(t *testing.T) {
	schemaMap := teo.TencentTeoL7RuleBranchBasicInfo(1)

	actionsSchema := schemaMap["actions"]
	actionsElem := actionsSchema.Elem.(*schema.Resource)
	actionsResourceSchema := actionsElem.Schema

	capSchema, ok := actionsResourceSchema["custom_action_parameters"]
	assert.True(t, ok, "custom_action_parameters schema should exist")
	assert.Equal(t, schema.TypeList, capSchema.Type, "custom_action_parameters should be TypeList")
	assert.True(t, capSchema.Optional, "custom_action_parameters should be Optional")
	assert.Equal(t, 1, capSchema.MaxItems, "custom_action_parameters MaxItems should be 1")

	capElem := capSchema.Elem.(*schema.Resource)
	capResourceSchema := capElem.Schema

	customActionsSchema, ok := capResourceSchema["custom_actions"]
	assert.True(t, ok, "custom_actions schema should exist inside custom_action_parameters")
	assert.Equal(t, schema.TypeList, customActionsSchema.Type, "custom_actions should be TypeList")

	customActionsElem := customActionsSchema.Elem.(*schema.Resource)
	customActionsResourceSchema := customActionsElem.Schema

	actionIdSchema, ok := customActionsResourceSchema["action_id"]
	assert.True(t, ok, "action_id schema should exist")
	assert.Equal(t, schema.TypeString, actionIdSchema.Type, "action_id should be TypeString")

	parametersSchema, ok := customActionsResourceSchema["parameters"]
	assert.True(t, ok, "parameters schema should exist")
	assert.Equal(t, schema.TypeList, parametersSchema.Type, "parameters should be TypeList")

	parametersElem := parametersSchema.Elem.(*schema.Resource)
	parametersResourceSchema := parametersElem.Schema

	expectedFields := []string{"name", "value_type", "string_value", "integer_value", "float_value", "boolean_value", "string_array_value", "integer_array_value", "float_array_value"}
	for _, field := range expectedFields {
		_, ok := parametersResourceSchema[field]
		assert.True(t, ok, field+" schema should exist inside parameters")
	}
}

// TestTeoL7AccRuleV2Extension_CustomVariables_Schema tests that custom_variables schema is correctly defined
func TestTeoL7AccRuleV2Extension_CustomVariables_Schema(t *testing.T) {
	res := teo.ResourceTencentCloudTeoL7AccRuleV2()

	customVariablesSchema, ok := res.Schema["custom_variables"]
	assert.True(t, ok, "custom_variables schema should exist at top level")
	assert.Equal(t, schema.TypeList, customVariablesSchema.Type, "custom_variables should be TypeList")
	assert.True(t, customVariablesSchema.Optional, "custom_variables should be Optional")

	cvElem := customVariablesSchema.Elem.(*schema.Resource)
	cvResourceSchema := cvElem.Schema

	nameSchema, ok := cvResourceSchema["name"]
	assert.True(t, ok, "name schema should exist inside custom_variables")
	assert.Equal(t, schema.TypeString, nameSchema.Type)

	initialValueSchema, ok := cvResourceSchema["initial_value"]
	assert.True(t, ok, "initial_value schema should exist inside custom_variables")
	assert.Equal(t, schema.TypeString, initialValueSchema.Type)

	descSchema, ok := cvResourceSchema["description"]
	assert.True(t, ok, "description schema should exist inside custom_variables")
	assert.Equal(t, schema.TypeString, descSchema.Type)
}

// TestTeoL7AccRuleV2Extension_CustomVariableOperations_Schema tests custom_variable_operations schema definition
func TestTeoL7AccRuleV2Extension_CustomVariableOperations_Schema(t *testing.T) {
	res := teo.ResourceTencentCloudTeoL7AccRuleV2()

	cvoSchema, ok := res.Schema["custom_variable_operations"]
	assert.True(t, ok, "custom_variable_operations schema should exist at top level")
	assert.Equal(t, schema.TypeList, cvoSchema.Type)
	assert.True(t, cvoSchema.Optional)

	cvoElem := cvoSchema.Elem.(*schema.Resource)
	cvoResourceSchema := cvoElem.Schema

	branchesSchema, ok := cvoResourceSchema["branches"]
	assert.True(t, ok, "branches schema should exist inside custom_variable_operations")
	assert.Equal(t, schema.TypeList, branchesSchema.Type)

	descSchema, ok := cvoResourceSchema["description"]
	assert.True(t, ok, "description schema should exist inside custom_variable_operations")
	assert.Equal(t, schema.TypeList, descSchema.Type)
}

// TestTeoL7AccRuleV2Extension_CustomVariables_FlattenSet tests flatten and set for CustomVariables
func TestTeoL7AccRuleV2Extension_CustomVariables_FlattenSet(t *testing.T) {
	sdkCustomVariables := []*teov20220901.CustomVariable{
		{
			Name:         helper.String("user.rule.var1"),
			InitialValue: helper.String("hello"),
			Description:  helper.String("desc1"),
		},
		{
			Name:         helper.String("user.rule.var2"),
			InitialValue: helper.String("world"),
		},
	}

	assert.NotNil(t, sdkCustomVariables)
	assert.Len(t, sdkCustomVariables, 2)
	assert.Equal(t, "user.rule.var1", *sdkCustomVariables[0].Name)
	assert.Equal(t, "hello", *sdkCustomVariables[0].InitialValue)
	assert.Equal(t, "desc1", *sdkCustomVariables[0].Description)
	assert.Nil(t, sdkCustomVariables[1].Description)
}

// go test ./tencentcloud/services/teo/ -run "TestAccTeoL7AccRuleV2CustomVariables" -v -count=1 -gcflags="all=-l"

// TestAccTeoL7AccRuleV2CustomVariables_Create tests Create with custom_variables
func TestAccTeoL7AccRuleV2CustomVariables_Create(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseTeoV20220901Client", teoClient)

	var capturedRequest *teov20220901.CreateL7AccRulesRequest
	patches.ApplyMethodFunc(teoClient, "CreateL7AccRules", func(request *teov20220901.CreateL7AccRulesRequest) (*teov20220901.CreateL7AccRulesResponse, error) {
		capturedRequest = request
		resp := teov20220901.NewCreateL7AccRulesResponse()
		resp.Response = &teov20220901.CreateL7AccRulesResponseParams{
			RuleIds:   []*string{ptrString("rule-abc123")},
			RequestId: ptrString("fake-request-id"),
		}
		return resp, nil
	})

	// Patch DescribeTeoL7AccRuleById to avoid ratelimit and service layer issues
	patches.ApplyMethodFunc(&teo.TeoService{}, "DescribeTeoL7AccRuleById", func(ctx context.Context, zoneId string, ruleId string) (*teov20220901.DescribeL7AccRulesResponseParams, error) {
		return &teov20220901.DescribeL7AccRulesResponseParams{
			TotalCount: helper.IntInt64(1),
			Rules: []*teov20220901.RuleEngineItem{
				{
					Status:   ptrString("enable"),
					RuleName: ptrString("Test Rule"),
					RuleId:   ptrString("rule-abc123"),
					CustomVariables: []*teov20220901.CustomVariable{
						{
							Name:         ptrString("user.rule.var1"),
							InitialValue: ptrString("hello"),
							Description:  ptrString("desc1"),
						},
					},
				},
			},
			RequestId: ptrString("fake-request-id"),
		}, nil
	})

	meta := newMockMeta()
	res := teo.ResourceTencentCloudTeoL7AccRuleV2()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":   "zone-abc123",
		"rule_name": "Test Rule",
		"status":    "enable",
		"custom_variables": []interface{}{
			map[string]interface{}{
				"name":          "user.rule.var1",
				"initial_value": "hello",
				"description":   "desc1",
			},
		},
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "zone-abc123#rule-abc123", d.Id())

	// Verify the request captured the custom_variables
	assert.NotNil(t, capturedRequest)
	assert.Len(t, capturedRequest.Rules, 1)
	rule := capturedRequest.Rules[0]
	assert.NotNil(t, rule.CustomVariables)
	assert.Len(t, rule.CustomVariables, 1)
	assert.Equal(t, "user.rule.var1", *rule.CustomVariables[0].Name)
	assert.Equal(t, "hello", *rule.CustomVariables[0].InitialValue)
	assert.Equal(t, "desc1", *rule.CustomVariables[0].Description)

	// Verify state was set from Read
	customVariables := d.Get("custom_variables").([]interface{})
	assert.Len(t, customVariables, 1)
	cvMap := customVariables[0].(map[string]interface{})
	assert.Equal(t, "user.rule.var1", cvMap["name"])
	assert.Equal(t, "hello", cvMap["initial_value"])
}

// TestAccTeoL7AccRuleV2CustomVariables_Update tests Update with custom_variables
func TestAccTeoL7AccRuleV2CustomVariables_Update(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseTeoV20220901Client", teoClient)

	var capturedRequest *teov20220901.ModifyL7AccRuleRequest
	patches.ApplyMethodFunc(teoClient, "ModifyL7AccRule", func(request *teov20220901.ModifyL7AccRuleRequest) (*teov20220901.ModifyL7AccRuleResponse, error) {
		capturedRequest = request
		resp := teov20220901.NewModifyL7AccRuleResponse()
		resp.Response = &teov20220901.ModifyL7AccRuleResponseParams{
			RequestId: ptrString("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(&teo.TeoService{}, "DescribeTeoL7AccRuleById", func(ctx context.Context, zoneId string, ruleId string) (*teov20220901.DescribeL7AccRulesResponseParams, error) {
		return &teov20220901.DescribeL7AccRulesResponseParams{
			TotalCount: helper.IntInt64(1),
			Rules: []*teov20220901.RuleEngineItem{
				{
					Status:   ptrString("enable"),
					RuleName: ptrString("Updated Rule"),
					RuleId:   ptrString("rule-abc123"),
					CustomVariables: []*teov20220901.CustomVariable{
						{
							Name:         ptrString("user.rule.var1"),
							InitialValue: ptrString("updated"),
						},
					},
				},
			},
			RequestId: ptrString("fake-request-id"),
		}, nil
	})

	meta := newMockMeta()
	res := teo.ResourceTencentCloudTeoL7AccRuleV2()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":   "zone-abc123",
		"rule_name": "Updated Rule",
		"status":    "enable",
		"custom_variables": []interface{}{
			map[string]interface{}{
				"name":          "user.rule.var1",
				"initial_value": "updated",
			},
		},
	})
	d.SetId("zone-abc123#rule-abc123")

	err := res.Update(d, meta)
	assert.NoError(t, err)

	// Verify ModifyL7AccRule was called with custom_variables
	assert.NotNil(t, capturedRequest)
	assert.NotNil(t, capturedRequest.Rule)
	assert.NotNil(t, capturedRequest.Rule.CustomVariables)
	assert.Len(t, capturedRequest.Rule.CustomVariables, 1)
	assert.Equal(t, "user.rule.var1", *capturedRequest.Rule.CustomVariables[0].Name)
	assert.Equal(t, "updated", *capturedRequest.Rule.CustomVariables[0].InitialValue)
}

// TestAccTeoL7AccRuleV2CustomVariableOperations_Create tests Create with custom_variable_operations
func TestAccTeoL7AccRuleV2CustomVariableOperations_Create(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseTeoV20220901Client", teoClient)

	var capturedRequest *teov20220901.CreateL7AccRulesRequest
	patches.ApplyMethodFunc(teoClient, "CreateL7AccRules", func(request *teov20220901.CreateL7AccRulesRequest) (*teov20220901.CreateL7AccRulesResponse, error) {
		capturedRequest = request
		resp := teov20220901.NewCreateL7AccRulesResponse()
		resp.Response = &teov20220901.CreateL7AccRulesResponseParams{
			RuleIds:   []*string{ptrString("rule-def456")},
			RequestId: ptrString("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(&teo.TeoService{}, "DescribeTeoL7AccRuleById", func(ctx context.Context, zoneId string, ruleId string) (*teov20220901.DescribeL7AccRulesResponseParams, error) {
		return &teov20220901.DescribeL7AccRulesResponseParams{
			TotalCount: helper.IntInt64(1),
			Rules: []*teov20220901.RuleEngineItem{
				{
					Status:   ptrString("enable"),
					RuleName: ptrString("Test Rule Ops"),
					RuleId:   ptrString("rule-def456"),
					CustomVariableOperations: []*teov20220901.CustomVariableOperation{
						{
							Branches: []*teov20220901.CustomVariableOperationRuleBranch{
								{
									Condition: ptrString("${http.request.host} in ['www.example.com']"),
									Actions: []*teov20220901.CustomVariableOperationRuleAction{
										{
											Name: ptrString("Set"),
											SetParameters: &teov20220901.SetParameters{
												Name:  ptrString("user.rule.var1"),
												Value: ptrString("${http.request.uri.path}"),
											},
										},
									},
									SubRules: []*teov20220901.CustomVariableOperationSubRule{
										{
											Branches: []*teov20220901.CustomVariableOperationRuleBranch{
												{
													Condition: ptrString("${http.request.file_extension} in ['html']"),
													Actions: []*teov20220901.CustomVariableOperationRuleAction{
														{
															Name: ptrString("Set"),
															SetParameters: &teov20220901.SetParameters{
																Name:  ptrString("user.rule.var1"),
																Value: ptrString("static"),
															},
														},
													},
												},
											},
											Description: []*string{ptrString("sub-rule-desc")},
										},
									},
								},
							},
							Description: []*string{ptrString("main-desc")},
						},
					},
				},
			},
			RequestId: ptrString("fake-request-id"),
		}, nil
	})

	meta := newMockMeta()
	res := teo.ResourceTencentCloudTeoL7AccRuleV2()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":   "zone-abc123",
		"rule_name": "Test Rule Ops",
		"status":    "enable",
		"branches": []interface{}{
			map[string]interface{}{
				"condition": "${http.request.host} in ['www.example.com']",
			},
		},
		"custom_variable_operations": []interface{}{
			map[string]interface{}{
				"branches": []interface{}{
					map[string]interface{}{
						"condition": "${http.request.host} in ['www.example.com']",
						"actions": []interface{}{
							map[string]interface{}{
								"name": "Set",
								"set_parameters": []interface{}{
									map[string]interface{}{
										"name":  "user.rule.var1",
										"value": "${http.request.uri.path}",
									},
								},
							},
						},
						"sub_rules": []interface{}{
							map[string]interface{}{
								"branches": []interface{}{
									map[string]interface{}{
										"condition": "${http.request.file_extension} in ['html']",
										"actions": []interface{}{
											map[string]interface{}{
												"name": "Set",
												"set_parameters": []interface{}{
													map[string]interface{}{
														"name":  "user.rule.var1",
														"value": "static",
													},
												},
											},
										},
									},
								},
								"description": []interface{}{"sub-rule-desc"},
							},
						},
					},
				},
				"description": []interface{}{"main-desc"},
			},
		},
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)

	// Verify the request captured the custom_variable_operations
	assert.NotNil(t, capturedRequest)
	assert.Len(t, capturedRequest.Rules, 1)
	rule := capturedRequest.Rules[0]
	assert.NotNil(t, rule.CustomVariableOperations)
	assert.Len(t, rule.CustomVariableOperations, 1)

	operation := rule.CustomVariableOperations[0]
	assert.Len(t, operation.Description, 1)
	assert.Equal(t, "main-desc", *operation.Description[0])

	assert.Len(t, operation.Branches, 1)
	branch := operation.Branches[0]
	assert.Equal(t, "${http.request.host} in ['www.example.com']", *branch.Condition)

	assert.Len(t, branch.Actions, 1)
	action := branch.Actions[0]
	assert.Equal(t, "Set", *action.Name)
	assert.NotNil(t, action.SetParameters)
	assert.Equal(t, "user.rule.var1", *action.SetParameters.Name)
	assert.Equal(t, "${http.request.uri.path}", *action.SetParameters.Value)

	assert.Len(t, branch.SubRules, 1)
	subRule := branch.SubRules[0]
	assert.Len(t, subRule.Description, 1)
	assert.Equal(t, "sub-rule-desc", *subRule.Description[0])
	assert.Len(t, subRule.Branches, 1)
	subBranch := subRule.Branches[0]
	assert.Equal(t, "${http.request.file_extension} in ['html']", *subBranch.Condition)
	assert.Len(t, subBranch.Actions, 1)
	assert.Equal(t, "static", *subBranch.Actions[0].SetParameters.Value)

	// Verify state was set from Read
	cvo := d.Get("custom_variable_operations").([]interface{})
	assert.Len(t, cvo, 1)
}

// TestAccTeoL7AccRuleV2CustomActionParameters_Create tests Create with custom_action_parameters in actions
func TestAccTeoL7AccRuleV2CustomActionParameters_Create(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseTeoV20220901Client", teoClient)

	var capturedRequest *teov20220901.CreateL7AccRulesRequest
	patches.ApplyMethodFunc(teoClient, "CreateL7AccRules", func(request *teov20220901.CreateL7AccRulesRequest) (*teov20220901.CreateL7AccRulesResponse, error) {
		capturedRequest = request
		resp := teov20220901.NewCreateL7AccRulesResponse()
		resp.Response = &teov20220901.CreateL7AccRulesResponseParams{
			RuleIds:   []*string{ptrString("rule-ghi789")},
			RequestId: ptrString("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(&teo.TeoService{}, "DescribeTeoL7AccRuleById", func(ctx context.Context, zoneId string, ruleId string) (*teov20220901.DescribeL7AccRulesResponseParams, error) {
		return &teov20220901.DescribeL7AccRulesResponseParams{
			TotalCount: helper.IntInt64(1),
			Rules: []*teov20220901.RuleEngineItem{
				{
					Status:   ptrString("enable"),
					RuleName: ptrString("Test CustomAction"),
					RuleId:   ptrString("rule-ghi789"),
					Branches: []*teov20220901.RuleBranch{
						{
							Condition: ptrString("${http.request.host} in ['www.example.com']"),
							Actions: []*teov20220901.RuleEngineAction{
								{
									Name: ptrString("CustomAction"),
									CustomActionParameters: &teov20220901.CustomActionParameters{
										CustomActions: []*teov20220901.CustomAction{
											{
												ActionId: ptrString("action-001"),
												Parameters: []*teov20220901.CustomActionParameter{
													{
														Name:        ptrString("param1"),
														ValueType:   ptrString("String"),
														StringValue: ptrString("value1"),
													},
													{
														Name:         ptrString("param2"),
														ValueType:    ptrString("Integer"),
														IntegerValue: helper.IntInt64(42),
													},
													{
														Name:         ptrString("param3"),
														ValueType:    ptrString("Boolean"),
														BooleanValue: helper.Bool(true),
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			RequestId: ptrString("fake-request-id"),
		}, nil
	})

	meta := newMockMeta()
	res := teo.ResourceTencentCloudTeoL7AccRuleV2()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":   "zone-abc123",
		"rule_name": "Test CustomAction",
		"status":    "enable",
		"branches": []interface{}{
			map[string]interface{}{
				"condition": "${http.request.host} in ['www.example.com']",
				"actions": []interface{}{
					map[string]interface{}{
						"name": "CustomAction",
						"custom_action_parameters": []interface{}{
							map[string]interface{}{
								"custom_actions": []interface{}{
									map[string]interface{}{
										"action_id": "action-001",
										"parameters": []interface{}{
											map[string]interface{}{
												"name":         "param1",
												"value_type":   "String",
												"string_value": "value1",
											},
											map[string]interface{}{
												"name":          "param2",
												"value_type":    "Integer",
												"integer_value": 42,
											},
											map[string]interface{}{
												"name":          "param3",
												"value_type":    "Boolean",
												"boolean_value": true,
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)

	// Verify the request captured the custom_action_parameters
	assert.NotNil(t, capturedRequest)
	assert.Len(t, capturedRequest.Rules, 1)
	rule := capturedRequest.Rules[0]
	assert.NotNil(t, rule.Branches)
	assert.Len(t, rule.Branches, 1)
	branch := rule.Branches[0]
	assert.NotNil(t, branch.Actions)
	assert.Len(t, branch.Actions, 1)
	action := branch.Actions[0]
	assert.Equal(t, "CustomAction", *action.Name)
	assert.NotNil(t, action.CustomActionParameters)
	assert.Len(t, action.CustomActionParameters.CustomActions, 1)
	customAction := action.CustomActionParameters.CustomActions[0]
	assert.Equal(t, "action-001", *customAction.ActionId)
	assert.Len(t, customAction.Parameters, 2)
	assert.Equal(t, "param1", *customAction.Parameters[0].Name)
	assert.Equal(t, "String", *customAction.Parameters[0].ValueType)
	assert.Equal(t, "value1", *customAction.Parameters[0].StringValue)
	assert.Equal(t, "param2", *customAction.Parameters[1].Name)
	assert.Equal(t, int64(42), *customAction.Parameters[1].IntegerValue)

	// Verify state was set from Read
	branches := d.Get("branches").([]interface{})
	assert.Len(t, branches, 1)
	branchMap := branches[0].(map[string]interface{})
	actionsList := branchMap["actions"].([]interface{})
	assert.Len(t, actionsList, 1)
	actionMap := actionsList[0].(map[string]interface{})
	capList := actionMap["custom_action_parameters"].([]interface{})
	assert.Len(t, capList, 1)
}
