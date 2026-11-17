package teo_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/assert"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/teo"
)

func TestAccTencentCloudTeoFunctionResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTeoFunction,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_teo_function.teo_function", "id"),
					resource.TestCheckResourceAttr("tencentcloud_teo_function.teo_function", "name", "aaa"),
					resource.TestCheckResourceAttr("tencentcloud_teo_function.teo_function", "remark", "test"),
					resource.TestCheckResourceAttr("tencentcloud_teo_function.teo_function", "content", `addEventListener('fetch', e => {
  const response = new Response('Hello World!!');
  e.respondWith(response);
});
`),
				),
			},
			{
				ResourceName:      "tencentcloud_teo_function.teo_function",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccTeoFunctionUp,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_teo_function.teo_function", "id"),
					resource.TestCheckResourceAttr("tencentcloud_teo_function.teo_function", "name", "aaa"),
					resource.TestCheckResourceAttr("tencentcloud_teo_function.teo_function", "remark", "test-update"),
					resource.TestCheckResourceAttr("tencentcloud_teo_function.teo_function", "content", `addEventListener('fetch', e => {
  const response = new Response('Hello World');
  e.respondWith(response);
});
`),
				),
			},
		},
	})
}

const testAccTeoFunction = `

resource "tencentcloud_teo_function" "teo_function" {
    content     = <<-EOT
        addEventListener('fetch', e => {
          const response = new Response('Hello World!!');
          e.respondWith(response);
        });
    EOT
    name        = "aaa"
    remark      = "test"
    zone_id     = "zone-2qtuhspy7cr6"
}
`
const testAccTeoFunctionUp = `

resource "tencentcloud_teo_function" "teo_function" {
    content     = <<-EOT
        addEventListener('fetch', e => {
          const response = new Response('Hello World');
          e.respondWith(response);
        });
    EOT
    name        = "aaa"
    remark      = "test-update"
    zone_id     = "zone-2qtuhspy7cr6"
}
`

func TestParseTeoFunctionOriginalName(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		input    string
		zoneId   string
		expected string
	}{
		{
			name:     "standard concatenated name",
			input:    "my-func-zone-2qtuhspy7cr6-1310708577",
			zoneId:   "zone-2qtuhspy7cr6-1310708577",
			expected: "my-func",
		},
		{
			name:     "original name without hyphens",
			input:    "myfunc-zone-2qtuhspy7cr6-1310708577",
			zoneId:   "zone-2qtuhspy7cr6-1310708577",
			expected: "myfunc",
		},
		{
			name:     "original name with multiple hyphens",
			input:    "my-test-func-v2-zone-2qtuhspy7cr6-1310708577",
			zoneId:   "zone-2qtuhspy7cr6-1310708577",
			expected: "my-test-func-v2",
		},
		{
			name:     "name without -zone- substring",
			input:    "myfunc",
			zoneId:   "zone-2qtuhspy7cr6-1310708577",
			expected: "myfunc",
		},
		{
			name:     "empty string",
			input:    "",
			zoneId:   "zone-2qtuhspy7cr6-1310708577",
			expected: "",
		},
		{
			name:     "name starts with zone",
			input:    "zone-2qtuhspy7cr6-1310708577",
			zoneId:   "zone-2qtuhspy7cr6-1310708577",
			expected: "zone-2qtuhspy7cr6-1310708577",
		},
		{
			name:     "name with zone appearing before -zone-",
			input:    "zone-proxy-zone-2qtuhspy7cr6-1310708577",
			zoneId:   "zone-2qtuhspy7cr6-1310708577",
			expected: "zone-proxy",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := teo.ParseTeoFunctionOriginalName(tt.input, tt.zoneId)
			assert.Equal(t, tt.expected, result)
		})
	}
}
