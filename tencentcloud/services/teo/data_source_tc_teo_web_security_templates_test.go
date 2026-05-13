package teo_test

import (
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

// go test ./tencentcloud/services/teo/ -run "TestTeoWebSecurityTemplates" -v -count=1 -gcflags="all=-l"

// TestTeoWebSecurityTemplatesDataSource_Success tests Read calls API and sets security_policy_templates
func TestTeoWebSecurityTemplatesDataSource_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	// Patch UseTeoV20220901Client to return a non-nil client
	teoClient := &teov20220901.Client{}
	mockClient := &connectivity.TencentCloudClient{}
	patches.ApplyMethodReturn(mockClient, "UseTeoV20220901Client", teoClient)

	// Patch DescribeWebSecurityTemplates to return success
	patches.ApplyMethodFunc(teoClient, "DescribeWebSecurityTemplates", func(request *teov20220901.DescribeWebSecurityTemplatesRequest) (*teov20220901.DescribeWebSecurityTemplatesResponse, error) {
		resp := teov20220901.NewDescribeWebSecurityTemplatesResponse()
		resp.Response = &teov20220901.DescribeWebSecurityTemplatesResponseParams{
			TotalCount: ptrWebSecTmplInt64(2),
			SecurityPolicyTemplates: []*teov20220901.SecurityPolicyTemplateInfo{
				{
					ZoneId:       ptrWebSecTmplString("zone-abc123"),
					TemplateId:   ptrWebSecTmplString("template-001"),
					TemplateName: ptrWebSecTmplString("Default Template"),
					BindDomains: []*teov20220901.BindDomainInfo{
						{
							Domain: ptrWebSecTmplString("example.com"),
							ZoneId: ptrWebSecTmplString("zone-abc123"),
							Status: ptrWebSecTmplString("online"),
						},
						{
							Domain: ptrWebSecTmplString("test.example.com"),
							ZoneId: ptrWebSecTmplString("zone-abc123"),
							Status: ptrWebSecTmplString("process"),
						},
					},
				},
				{
					ZoneId:       ptrWebSecTmplString("zone-def456"),
					TemplateId:   ptrWebSecTmplString("template-002"),
					TemplateName: ptrWebSecTmplString("Custom Template"),
					BindDomains:  []*teov20220901.BindDomainInfo{},
				},
			},
			RequestId: ptrWebSecTmplString("fake-request-id"),
		}
		return resp, nil
	})

	meta := &webSecTmplMockMeta{client: mockClient}
	res := teo.DataSourceTencentCloudTeoWebSecurityTemplates()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_ids": []interface{}{"zone-abc123", "zone-def456"},
	})

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.NotEmpty(t, d.Id())

	templates := d.Get("security_policy_templates").([]interface{})
	assert.Len(t, templates, 2)

	tmpl0 := templates[0].(map[string]interface{})
	assert.Equal(t, "zone-abc123", tmpl0["zone_id"])
	assert.Equal(t, "template-001", tmpl0["template_id"])
	assert.Equal(t, "Default Template", tmpl0["template_name"])

	bindDomains0 := tmpl0["bind_domains"].([]interface{})
	assert.Len(t, bindDomains0, 2)

	bd0 := bindDomains0[0].(map[string]interface{})
	assert.Equal(t, "example.com", bd0["domain"])
	assert.Equal(t, "zone-abc123", bd0["zone_id"])
	assert.Equal(t, "online", bd0["status"])

	bd1 := bindDomains0[1].(map[string]interface{})
	assert.Equal(t, "test.example.com", bd1["domain"])
	assert.Equal(t, "zone-abc123", bd1["zone_id"])
	assert.Equal(t, "process", bd1["status"])

	tmpl1 := templates[1].(map[string]interface{})
	assert.Equal(t, "zone-def456", tmpl1["zone_id"])
	assert.Equal(t, "template-002", tmpl1["template_id"])
	assert.Equal(t, "Custom Template", tmpl1["template_name"])
}

// TestTeoWebSecurityTemplatesDataSource_EmptyResult tests Read with empty API response
func TestTeoWebSecurityTemplatesDataSource_EmptyResult(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	mockClient := &connectivity.TencentCloudClient{}
	patches.ApplyMethodReturn(mockClient, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeWebSecurityTemplates", func(request *teov20220901.DescribeWebSecurityTemplatesRequest) (*teov20220901.DescribeWebSecurityTemplatesResponse, error) {
		resp := teov20220901.NewDescribeWebSecurityTemplatesResponse()
		resp.Response = &teov20220901.DescribeWebSecurityTemplatesResponseParams{
			TotalCount:              ptrWebSecTmplInt64(0),
			SecurityPolicyTemplates: []*teov20220901.SecurityPolicyTemplateInfo{},
			RequestId:               ptrWebSecTmplString("fake-request-id"),
		}
		return resp, nil
	})

	meta := &webSecTmplMockMeta{client: mockClient}
	res := teo.DataSourceTencentCloudTeoWebSecurityTemplates()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_ids": []interface{}{"zone-nonexistent"},
	})

	err := res.Read(d, meta)
	assert.NoError(t, err)

	templates := d.Get("security_policy_templates").([]interface{})
	assert.Len(t, templates, 0)
}

// TestTeoWebSecurityTemplatesDataSource_APIError tests Read handles API error
func TestTeoWebSecurityTemplatesDataSource_APIError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	mockClient := &connectivity.TencentCloudClient{}
	patches.ApplyMethodReturn(mockClient, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeWebSecurityTemplates", func(request *teov20220901.DescribeWebSecurityTemplatesRequest) (*teov20220901.DescribeWebSecurityTemplatesResponse, error) {
		return nil, fmt.Errorf("[TencentCloudSDKError] Code=InvalidParameter, Message=Invalid ZoneIds")
	})

	meta := &webSecTmplMockMeta{client: mockClient}
	res := teo.DataSourceTencentCloudTeoWebSecurityTemplates()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_ids": []interface{}{"invalid-zone"},
	})

	err := res.Read(d, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "InvalidParameter")
}

// TestTeoWebSecurityTemplatesDataSource_Schema validates schema definition
func TestTeoWebSecurityTemplatesDataSource_Schema(t *testing.T) {
	res := teo.DataSourceTencentCloudTeoWebSecurityTemplates()

	assert.NotNil(t, res)
	assert.NotNil(t, res.Read)

	assert.Contains(t, res.Schema, "zone_ids")
	assert.Contains(t, res.Schema, "security_policy_templates")
	assert.Contains(t, res.Schema, "result_output_file")

	zoneIds := res.Schema["zone_ids"]
	assert.Equal(t, schema.TypeList, zoneIds.Type)
	assert.True(t, zoneIds.Required)
	assert.Equal(t, 100, zoneIds.MaxItems)

	templates := res.Schema["security_policy_templates"]
	assert.Equal(t, schema.TypeList, templates.Type)
	assert.True(t, templates.Computed)

	resultOutputFile := res.Schema["result_output_file"]
	assert.Equal(t, schema.TypeString, resultOutputFile.Type)
	assert.True(t, resultOutputFile.Optional)
}

// webSecTmplMockMeta implements tccommon.ProviderMeta
type webSecTmplMockMeta struct {
	client *connectivity.TencentCloudClient
}

func (m *webSecTmplMockMeta) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &webSecTmplMockMeta{}

func ptrWebSecTmplString(s string) *string {
	return &s
}

func ptrWebSecTmplInt64(n int64) *int64 {
	return &n
}
