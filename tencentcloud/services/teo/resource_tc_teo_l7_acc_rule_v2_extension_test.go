package teo_test

import (
	"testing"

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
