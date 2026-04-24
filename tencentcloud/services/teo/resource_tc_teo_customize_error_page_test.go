package teo_test

import (
	"fmt"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	teov20220901 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/teo"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

// mockMetaForCustomizeErrorPage implements tccommon.ProviderMeta
type mockMetaForCustomizeErrorPage struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMetaForCustomizeErrorPage) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockMetaForCustomizeErrorPage{}

func newMockMetaForCustomizeErrorPage() *mockMetaForCustomizeErrorPage {
	return &mockMetaForCustomizeErrorPage{client: &connectivity.TencentCloudClient{}}
}

func ptrStrCustomizeErrorPage(s string) *string {
	return &s
}

func TestAccTencentCloudTeoCustomizeErrorPageResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTeoCustomizeErrorPage,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_teo_customize_error_page.example", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_teo_customize_error_page.example", "zone_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_teo_customize_error_page.example", "name"),
					resource.TestCheckResourceAttrSet("tencentcloud_teo_customize_error_page.example", "content_type"),
					resource.TestCheckResourceAttrSet("tencentcloud_teo_customize_error_page.example", "description"),
					resource.TestCheckResourceAttrSet("tencentcloud_teo_customize_error_page.example", "content"),
				),
			},
			{
				Config: testAccTeoCustomizeErrorPageUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_teo_customize_error_page.example", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_teo_customize_error_page.example", "zone_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_teo_customize_error_page.example", "name"),
					resource.TestCheckResourceAttrSet("tencentcloud_teo_customize_error_page.example", "content_type"),
					resource.TestCheckResourceAttrSet("tencentcloud_teo_customize_error_page.example", "description"),
					resource.TestCheckResourceAttrSet("tencentcloud_teo_customize_error_page.example", "content"),
				),
			},
			{
				ResourceName:      "tencentcloud_teo_customize_error_page.example",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccTeoCustomizeErrorPage = `
resource "tencentcloud_teo_customize_error_page" "example" {
  zone_id      = "zone-3edjdliiw3he"
  name         = "tf-example"
  content_type = "text/plain"
  description  = "description."
  content      = "customize error page"
}
`

const testAccTeoCustomizeErrorPageUpdate = `
resource "tencentcloud_teo_customize_error_page" "example" {
  zone_id      = "zone-3edjdliiw3he"
  name         = "tf-example-update"
  content_type = "application/json"
  description  = "description update."
  content = jsonencode({
    "key" : "value",
  })
}
`

// go test ./tencentcloud/services/teo/ -run "TestCustomizeErrorPageErrorPages" -v -count=1 -gcflags="all=-l"

// TestCustomizeErrorPageErrorPages_Read_Success tests Read populates references attribute
func TestCustomizeErrorPageErrorPages_Read_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaForCustomizeErrorPage().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeCustomErrorPages", func(request *teov20220901.DescribeCustomErrorPagesRequest) (*teov20220901.DescribeCustomErrorPagesResponse, error) {
		resp := teov20220901.NewDescribeCustomErrorPagesResponse()
		resp.Response = &teov20220901.DescribeCustomErrorPagesResponseParams{
			TotalCount: ptrUint64CustomizeErrorPage(1),
			ErrorPages: []*teov20220901.CustomErrorPage{
				{
					PageId:      ptrStrCustomizeErrorPage("page-abc123"),
					ZoneId:      ptrStrCustomizeErrorPage("zone-test123"),
					Name:        ptrStrCustomizeErrorPage("test-error-page"),
					ContentType: ptrStrCustomizeErrorPage("text/html"),
					Description: ptrStrCustomizeErrorPage("test description"),
					Content:     ptrStrCustomizeErrorPage("<html>error</html>"),
					References: []*teov20220901.ErrorPageReference{
						{
							BusinessId: ptrStrCustomizeErrorPage("rule-001"),
						},
						{
							BusinessId: ptrStrCustomizeErrorPage("rule-002"),
						},
					},
				},
			},
			RequestId: ptrStrCustomizeErrorPage("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaForCustomizeErrorPage()
	res := teo.ResourceTencentCloudTeoCustomizeErrorPage()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":      "zone-test123",
		"name":         "test-error-page",
		"content_type": "text/html",
	})
	d.SetId("zone-test123#page-abc123")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "test-error-page", d.Get("name"))
	assert.Equal(t, "text/html", d.Get("content_type"))
	assert.Equal(t, "test description", d.Get("description"))
	assert.Equal(t, "<html>error</html>", d.Get("content"))

	// Verify top-level references computed attribute
	references := d.Get("references").([]interface{})
	assert.Equal(t, 2, len(references))
	assert.Equal(t, "rule-001", references[0])
	assert.Equal(t, "rule-002", references[1])
}

// TestCustomizeErrorPageErrorPages_Read_EmptyReferences tests Read with no references
func TestCustomizeErrorPageErrorPages_Read_EmptyReferences(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaForCustomizeErrorPage().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeCustomErrorPages", func(request *teov20220901.DescribeCustomErrorPagesRequest) (*teov20220901.DescribeCustomErrorPagesResponse, error) {
		resp := teov20220901.NewDescribeCustomErrorPagesResponse()
		resp.Response = &teov20220901.DescribeCustomErrorPagesResponseParams{
			TotalCount: ptrUint64CustomizeErrorPage(1),
			ErrorPages: []*teov20220901.CustomErrorPage{
				{
					PageId:      ptrStrCustomizeErrorPage("page-abc123"),
					ZoneId:      ptrStrCustomizeErrorPage("zone-test123"),
					Name:        ptrStrCustomizeErrorPage("test-error-page"),
					ContentType: ptrStrCustomizeErrorPage("text/plain"),
					Description: ptrStrCustomizeErrorPage("test description"),
					Content:     ptrStrCustomizeErrorPage("error page content"),
				},
			},
			RequestId: ptrStrCustomizeErrorPage("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaForCustomizeErrorPage()
	res := teo.ResourceTencentCloudTeoCustomizeErrorPage()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":      "zone-test123",
		"name":         "test-error-page",
		"content_type": "text/plain",
	})
	d.SetId("zone-test123#page-abc123")

	err := res.Read(d, meta)
	assert.NoError(t, err)

	// Verify top-level references is empty
	references := d.Get("references").([]interface{})
	assert.Equal(t, 0, len(references))
}

// TestCustomizeErrorPageErrorPages_Read_NotFound tests Read handles resource not found
func TestCustomizeErrorPageErrorPages_Read_NotFound(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaForCustomizeErrorPage().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeCustomErrorPages", func(request *teov20220901.DescribeCustomErrorPagesRequest) (*teov20220901.DescribeCustomErrorPagesResponse, error) {
		return nil, fmt.Errorf("[TencentCloudSDKError] Code=ResourceNotFound, Message=Error page not found")
	})

	meta := newMockMetaForCustomizeErrorPage()
	res := teo.ResourceTencentCloudTeoCustomizeErrorPage()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":      "zone-test123",
		"name":         "test-error-page",
		"content_type": "text/html",
	})
	d.SetId("zone-test123#page-abc123")

	err := res.Read(d, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "ResourceNotFound")
}

// TestCustomizeErrorPageErrorPages_Schema tests the references schema definition
func TestCustomizeErrorPageErrorPages_Schema(t *testing.T) {
	res := teo.ResourceTencentCloudTeoCustomizeErrorPage()

	assert.NotNil(t, res)

	// Check references field at top level
	assert.Contains(t, res.Schema, "references")
	referencesSchema := res.Schema["references"]
	assert.Equal(t, schema.TypeList, referencesSchema.Type)
	assert.True(t, referencesSchema.Computed)
	assert.False(t, referencesSchema.Optional)
	assert.False(t, referencesSchema.Required)

	// Verify references is TypeList of TypeString
	referencesElem, ok := referencesSchema.Elem.(*schema.Schema)
	assert.True(t, ok)
	assert.Equal(t, schema.TypeString, referencesElem.Type)
}

func ptrUint64CustomizeErrorPage(u uint64) *uint64 {
	return &u
}
