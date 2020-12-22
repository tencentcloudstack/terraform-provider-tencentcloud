package tencentcloud

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccTencentCloudCamSAMLProvider_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCamSAMLProviderDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCamSAMLProvider_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCamSAMLProviderExists("tencentcloud_cam_saml_provider.saml_provider_basic"),
					resource.TestCheckResourceAttr("tencentcloud_cam_saml_provider.saml_provider_basic", "name", "cam-saml-provider-test"),
					resource.TestCheckResourceAttrSet("tencentcloud_cam_saml_provider.saml_provider_basic", "meta_data"),
					resource.TestCheckResourceAttr("tencentcloud_cam_saml_provider.saml_provider_basic", "description", "test"),
				),
			}, {
				Config: testAccCamSAMLProvider_update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCamSAMLProviderExists("tencentcloud_cam_saml_provider.saml_provider_basic"),
					resource.TestCheckResourceAttr("tencentcloud_cam_saml_provider.saml_provider_basic", "name", "cam-saml-provider-test"),
					resource.TestCheckResourceAttrSet("tencentcloud_cam_saml_provider.saml_provider_basic", "meta_data"), resource.TestCheckResourceAttr("tencentcloud_cam_saml_provider.saml_provider_basic", "description", "test2"),
				),
			},
			{
				ResourceName:            "tencentcloud_cam_saml_provider.saml_provider_basic",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"provider_arn"},
			},
		},
	})
}

func testAccCheckCamSAMLProviderDestroy(s *terraform.State) error {
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	camService := CamService{
		client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn,
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_cam_saml_provider" {
			continue
		}

		instance, err := camService.DescribeSAMLProviderById(ctx, rs.Primary.ID)
		if err == nil && instance != nil && instance.Response != nil && instance.Response.Name != nil {
			return fmt.Errorf("[CHECK][CAM SAML provider][Destroy] check: CAM SAML provider still exists: %s", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckCamSAMLProviderExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("[CHECK][CAM SAML provider][Exists] check: CAM SAML provider %s is not found", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("[CHECK][CAM SAML provider][Exists] check: CAM SAML provider id is not set")
		}
		camService := CamService{
			client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn,
		}
		instance, err := camService.DescribeSAMLProviderById(ctx, rs.Primary.ID)
		if err != nil {
			return err
		}
		if instance == nil || instance.Response == nil || instance.Response.Name == nil {
			return fmt.Errorf("[CHECK][CAM SAML provider][Exists] check: CAM SAML provider %s is not exist", rs.Primary.ID)
		}
		return nil
	}
}

const testAccCamSAMLProvider_basic = `
resource "tencentcloud_cam_saml_provider" "saml_provider_basic" {
  name        = "cam-saml-provider-test"
  meta_data   = "PD94bWwgdmVyc2lvbj0iMS4wIiBlbmNvZGluZz0iVVRGLTgiPz48bWQ6RW50aXR5RGVzY3JpcHRvciBlbnRpdHlJRD0iaHR0cDovL3d3dy5va3RhLmNvbS9leGsxa3F4bWNqUW1HQURNeTM1NyIgeG1sbnM6bWQ9InVybjpvYXNpczpuYW1lczp0YzpTQU1MOjIuMDptZXRhZGF0YSI+PG1kOklEUFNTT0Rlc2NyaXB0b3IgV2FudEF1dGhuUmVxdWVzdHNTaWduZWQ9ImZhbHNlIiBwcm90b2NvbFN1cHBvcnRFbnVtZXJhdGlvbj0idXJuOm9hc2lzOm5hbWVzOnRjOlNBTUw6Mi4wOnByb3RvY29sIj48bWQ6S2V5RGVzY3JpcHRvciB1c2U9InNpZ25pbmciPjxkczpLZXlJbmZvIHhtbG5zOmRzPSJodHRwOi8vd3d3LnczLm9yZy8yMDAwLzA5L3htbGRzaWcjIj48ZHM6WDUwOURhdGE+PGRzOlg1MDlDZXJ0aWZpY2F0ZT5NSUlEb0RDQ0FvaWdBd0lCQWdJR0FXM0lTcExvTUEwR0NTcUdTSWIzRFFFQkN3VUFNSUdRTVFzd0NRWURWUVFHRXdKVlV6RVRNQkVHDQpBMVVFQ0F3S1EyRnNhV1p2Y201cFlURVdNQlFHQTFVRUJ3d05VMkZ1SUVaeVlXNWphWE5qYnpFTk1Bc0dBMVVFQ2d3RVQydDBZVEVVDQpNQklHQTFVRUN3d0xVMU5QVUhKdmRtbGtaWEl4RVRBUEJnTlZCQU1NQ0dsa2VIVmxkblJoTVJ3d0dnWUpLb1pJaHZjTkFRa0JGZzFwDQpibVp2UUc5cmRHRXVZMjl0TUI0WERURTVNVEF4TkRBek1qSXhNMW9YRFRJNU1UQXhOREF6TWpNeE0xb3dnWkF4Q3pBSkJnTlZCQVlUDQpBbFZUTVJNd0VRWURWUVFJREFwRFlXeHBabTl5Ym1saE1SWXdGQVlEVlFRSERBMVRZVzRnUm5KaGJtTnBjMk52TVEwd0N3WURWUVFLDQpEQVJQYTNSaE1SUXdFZ1lEVlFRTERBdFRVMDlRY205MmFXUmxjakVSTUE4R0ExVUVBd3dJYVdSNGRXVjJkR0V4SERBYUJna3Foa2lHDQo5dzBCQ1FFV0RXbHVabTlBYjJ0MFlTNWpiMjB3Z2dFaU1BMEdDU3FHU0liM0RRRUJBUVVBQTRJQkR3QXdnZ0VLQW9JQkFRQ2g4b3dqDQpZK2dQSUM3blQvNTduLzdmeXJzcDlHMXdxa2UxdXhjMHVrTndnQXozOVNpelY3QVhLMWRReTFLaThXWjJJMzFEczJkT0FNQ1FKR2pWDQpUWWNNbnA3KzhqUzNLdmxNUkRJamk5cmxuUi9vcnBvMll1RHVWby9jVzdidlRIS2h2REo1QWZRaWxzYlNPTXdUOWM2TVlYZGhBNVBwDQpzelFsK1UrdHJmcXUrdUorSER4SVQxdlhWaVI5YlY2SUFRSzZpbWZoc2wxWmVSUytjbVFVNEpjQWlYT0xtTnFVVWM2UkpxUzhrMW1mDQpBLzhmb2VyMGc3SG4xZDVXclpCc2gyUlR2Vzh1ZVdadHQ3dmh4QTlGdE5kSVlEcXJ0eElmMlZXcXhrSHM3WFZDSm5wTnJITVovT1BRDQpGY21YSGVxNlJJMlB3Q1RlOW8zZHZpM0hqeXBaOEl4dkFnTUJBQUV3RFFZSktvWklodmNOQVFFTEJRQURnZ0VCQUFHaHk1bG9nbGtTDQoyVHg2YS90MnF5VEx0YVV5cEwrNGhySGJoMVAweVVMc0NrSnFsM2wrWG9VZDZCY2FJaFNSVGFPQk95ODViL0UzelJ4K3JzQXJwTjVVDQp5ZThuUEM4a05PYW5vTk9wWnZvYmhpTzFlMFIvYmxEcnRBL0o5UlBwMWtmdlhmS2NSTTU3TlRCWXppTURlbnFQUTRFOWN1U2lGdFFxDQpJYmpIbThaM1B1YXgwRitldkZ3U1pJMDNCWXNISGw1d1EraEJBS3hTdTJINEZRdU93Zmpnb2EveEN6Z1NKYjJ2UXdEc1MxMk9mSkNiDQpSRm1ZL1VYZXQramFhdEVORktLZStZSUJpU0J2WG1adTN0MHN5NDZTNzlPVzBacXJ0NUh2bElsT2lpTFpaN1FZamxjM1kxeG1LZ1luDQpXM2M2WGZkdmhGWHo0ZDdkbWYvTUdpNGY0enM9PC9kczpYNTA5Q2VydGlmaWNhdGU+PC9kczpYNTA5RGF0YT48L2RzOktleUluZm8+PC9tZDpLZXlEZXNjcmlwdG9yPjxtZDpOYW1lSURGb3JtYXQ+dXJuOm9hc2lzOm5hbWVzOnRjOlNBTUw6MS4xOm5hbWVpZC1mb3JtYXQ6dW5zcGVjaWZpZWQ8L21kOk5hbWVJREZvcm1hdD48bWQ6TmFtZUlERm9ybWF0PnVybjpvYXNpczpuYW1lczp0YzpTQU1MOjEuMTpuYW1laWQtZm9ybWF0OmVtYWlsQWRkcmVzczwvbWQ6TmFtZUlERm9ybWF0PjxtZDpTaW5nbGVTaWduT25TZXJ2aWNlIEJpbmRpbmc9InVybjpvYXNpczpuYW1lczp0YzpTQU1MOjIuMDpiaW5kaW5nczpIVFRQLVBPU1QiIExvY2F0aW9uPSJodHRwczovL2lkeHVldnRhLm9rdGEuY29tL2FwcC9pZHh1ZW9yZzYzNzM1OF90ZXN0XzEvZXhrMWtxeG1jalFtR0FETXkzNTcvc3NvL3NhbWwiLz48bWQ6U2luZ2xlU2lnbk9uU2VydmljZSBCaW5kaW5nPSJ1cm46b2FzaXM6bmFtZXM6dGM6U0FNTDoyLjA6YmluZGluZ3M6SFRUUC1SZWRpcmVjdCIgTG9jYXRpb249Imh0dHBzOi8vaWR4dWV2dGEub2t0YS5jb20vYXBwL2lkeHVlb3JnNjM3MzU4X3Rlc3RfMS9leGsxa3F4bWNqUW1HQURNeTM1Ny9zc28vc2FtbCIvPjwvbWQ6SURQU1NPRGVzY3JpcHRvcj48L21kOkVudGl0eURlc2NyaXB0b3I+"
  description = "test"
}
`

const testAccCamSAMLProvider_update = `
resource "tencentcloud_cam_saml_provider" "saml_provider_basic" {
  name        = "cam-saml-provider-test"
  meta_data   = "PD94bWwgdmVyc2lvbj0iMS4wIiBlbmNvZGluZz0iVVRGLTgiPz48bWQ6RW50aXR5RGVzY3JpcHRvciBlbnRpdHlJRD0iaHR0cDovL3d3dy5va3RhLmNvbS9leGsxa3F4bWNqUW1HQURNeTM1NyIgeG1sbnM6bWQ9InVybjpvYXNpczpuYW1lczp0YzpTQU1MOjIuMDptZXRhZGF0YSI+PG1kOklEUFNTT0Rlc2NyaXB0b3IgV2FudEF1dGhuUmVxdWVzdHNTaWduZWQ9ImZhbHNlIiBwcm90b2NvbFN1cHBvcnRFbnVtZXJhdGlvbj0idXJuOm9hc2lzOm5hbWVzOnRjOlNBTUw6Mi4wOnByb3RvY29sIj48bWQ6S2V5RGVzY3JpcHRvciB1c2U9InNpZ25pbmciPjxkczpLZXlJbmZvIHhtbG5zOmRzPSJodHRwOi8vd3d3LnczLm9yZy8yMDAwLzA5L3htbGRzaWcjIj48ZHM6WDUwOURhdGE+PGRzOlg1MDlDZXJ0aWZpY2F0ZT5NSUlEb0RDQ0FvaWdBd0lCQWdJR0FXM0lTcExvTUEwR0NTcUdTSWIzRFFFQkN3VUFNSUdRTVFzd0NRWURWUVFHRXdKVlV6RVRNQkVHDQpBMVVFQ0F3S1EyRnNhV1p2Y201cFlURVdNQlFHQTFVRUJ3d05VMkZ1SUVaeVlXNWphWE5qYnpFTk1Bc0dBMVVFQ2d3RVQydDBZVEVVDQpNQklHQTFVRUN3d0xVMU5QVUhKdmRtbGtaWEl4RVRBUEJnTlZCQU1NQ0dsa2VIVmxkblJoTVJ3d0dnWUpLb1pJaHZjTkFRa0JGZzFwDQpibVp2UUc5cmRHRXVZMjl0TUI0WERURTVNVEF4TkRBek1qSXhNMW9YRFRJNU1UQXhOREF6TWpNeE0xb3dnWkF4Q3pBSkJnTlZCQVlUDQpBbFZUTVJNd0VRWURWUVFJREFwRFlXeHBabTl5Ym1saE1SWXdGQVlEVlFRSERBMVRZVzRnUm5KaGJtTnBjMk52TVEwd0N3WURWUVFLDQpEQVJQYTNSaE1SUXdFZ1lEVlFRTERBdFRVMDlRY205MmFXUmxjakVSTUE4R0ExVUVBd3dJYVdSNGRXVjJkR0V4SERBYUJna3Foa2lHDQo5dzBCQ1FFV0RXbHVabTlBYjJ0MFlTNWpiMjB3Z2dFaU1BMEdDU3FHU0liM0RRRUJBUVVBQTRJQkR3QXdnZ0VLQW9JQkFRQ2g4b3dqDQpZK2dQSUM3blQvNTduLzdmeXJzcDlHMXdxa2UxdXhjMHVrTndnQXozOVNpelY3QVhLMWRReTFLaThXWjJJMzFEczJkT0FNQ1FKR2pWDQpUWWNNbnA3KzhqUzNLdmxNUkRJamk5cmxuUi9vcnBvMll1RHVWby9jVzdidlRIS2h2REo1QWZRaWxzYlNPTXdUOWM2TVlYZGhBNVBwDQpzelFsK1UrdHJmcXUrdUorSER4SVQxdlhWaVI5YlY2SUFRSzZpbWZoc2wxWmVSUytjbVFVNEpjQWlYT0xtTnFVVWM2UkpxUzhrMW1mDQpBLzhmb2VyMGc3SG4xZDVXclpCc2gyUlR2Vzh1ZVdadHQ3dmh4QTlGdE5kSVlEcXJ0eElmMlZXcXhrSHM3WFZDSm5wTnJITVovT1BRDQpGY21YSGVxNlJJMlB3Q1RlOW8zZHZpM0hqeXBaOEl4dkFnTUJBQUV3RFFZSktvWklodmNOQVFFTEJRQURnZ0VCQUFHaHk1bG9nbGtTDQoyVHg2YS90MnF5VEx0YVV5cEwrNGhySGJoMVAweVVMc0NrSnFsM2wrWG9VZDZCY2FJaFNSVGFPQk95ODViL0UzelJ4K3JzQXJwTjVVDQp5ZThuUEM4a05PYW5vTk9wWnZvYmhpTzFlMFIvYmxEcnRBL0o5UlBwMWtmdlhmS2NSTTU3TlRCWXppTURlbnFQUTRFOWN1U2lGdFFxDQpJYmpIbThaM1B1YXgwRitldkZ3U1pJMDNCWXNISGw1d1EraEJBS3hTdTJINEZRdU93Zmpnb2EveEN6Z1NKYjJ2UXdEc1MxMk9mSkNiDQpSRm1ZL1VYZXQramFhdEVORktLZStZSUJpU0J2WG1adTN0MHN5NDZTNzlPVzBacXJ0NUh2bElsT2lpTFpaN1FZamxjM1kxeG1LZ1luDQpXM2M2WGZkdmhGWHo0ZDdkbWYvTUdpNGY0enM9PC9kczpYNTA5Q2VydGlmaWNhdGU+PC9kczpYNTA5RGF0YT48L2RzOktleUluZm8+PC9tZDpLZXlEZXNjcmlwdG9yPjxtZDpOYW1lSURGb3JtYXQ+dXJuOm9hc2lzOm5hbWVzOnRjOlNBTUw6MS4xOm5hbWVpZC1mb3JtYXQ6dW5zcGVjaWZpZWQ8L21kOk5hbWVJREZvcm1hdD48bWQ6TmFtZUlERm9ybWF0PnVybjpvYXNpczpuYW1lczp0YzpTQU1MOjEuMTpuYW1laWQtZm9ybWF0OmVtYWlsQWRkcmVzczwvbWQ6TmFtZUlERm9ybWF0PjxtZDpTaW5nbGVTaWduT25TZXJ2aWNlIEJpbmRpbmc9InVybjpvYXNpczpuYW1lczp0YzpTQU1MOjIuMDpiaW5kaW5nczpIVFRQLVBPU1QiIExvY2F0aW9uPSJodHRwczovL2lkeHVldnRhLm9rdGEuY29tL2FwcC9pZHh1ZW9yZzYzNzM1OF90ZXN0XzEvZXhrMWtxeG1jalFtR0FETXkzNTcvc3NvL3NhbWwiLz48bWQ6U2luZ2xlU2lnbk9uU2VydmljZSBCaW5kaW5nPSJ1cm46b2FzaXM6bmFtZXM6dGM6U0FNTDoyLjA6YmluZGluZ3M6SFRUUC1SZWRpcmVjdCIgTG9jYXRpb249Imh0dHBzOi8vaWR4dWV2dGEub2t0YS5jb20vYXBwL2lkeHVlb3JnNjM3MzU4X3Rlc3RfMS9leGsxa3F4bWNqUW1HQURNeTM1Ny9zc28vc2FtbCIvPjwvbWQ6SURQU1NPRGVzY3JpcHRvcj48L21kOkVudGl0eURlc2NyaXB0b3I+"
  description = "test2"
}
`
