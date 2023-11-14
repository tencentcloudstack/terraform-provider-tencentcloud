package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudApigatewayAPIResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccApigatewayAPI,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_apigateway_a_p_i.a_p_i", "id")),
			},
			{
				ResourceName:      "tencentcloud_apigateway_a_p_i.a_p_i",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccApigatewayAPI = `

resource "tencentcloud_apigateway_a_p_i" "a_p_i" {
  service_id = ""
  service_type = ""
  service_timeout = 
  protocol = ""
  request_config {
		path = ""
		method = ""

  }
  api_name = ""
  api_desc = ""
  api_type = ""
  auth_type = ""
  enable_c_o_r_s = 
  constant_parameters {
		name = ""
		desc = ""
		position = ""
		default_value = ""

  }
  request_parameters {
		name = ""
		desc = ""
		position = ""
		type = ""
		default_value = ""
		required = 

  }
  api_business_type = ""
  service_mock_return_message = ""
  micro_services {
		cluster_id = ""
		namespace_id = ""
		micro_service_name = ""

  }
  service_tsf_load_balance_conf {
		is_load_balance = 
		method = ""
		session_stick_required = 
		session_stick_timeout = 

  }
  service_tsf_health_check_conf {
		is_health_check = 
		request_volume_threshold = 
		sleep_window_in_milliseconds = 
		error_threshold_percentage = 

  }
  target_services {
		vm_ip = ""
		vpc_id = ""
		vm_port = 
		host_ip = ""
		docker_ip = ""

  }
  target_services_load_balance_conf = 
  target_services_health_check_conf {
		is_health_check = 
		request_volume_threshold = 
		sleep_window_in_milliseconds = 
		error_threshold_percentage = 

  }
  service_scf_function_name = ""
  service_websocket_register_function_name = ""
  service_websocket_cleanup_function_name = ""
  service_websocket_transport_function_name = ""
  service_scf_function_namespace = ""
  service_scf_function_qualifier = ""
  service_websocket_register_function_namespace = ""
  service_websocket_register_function_qualifier = ""
  service_websocket_transport_function_namespace = ""
  service_websocket_transport_function_qualifier = ""
  service_websocket_cleanup_function_namespace = ""
  service_websocket_cleanup_function_qualifier = ""
  service_scf_is_integrated_response = 
  is_debug_after_charge = 
  is_delete_response_error_codes = 
  response_type = ""
  response_success_example = ""
  response_fail_example = ""
  service_config {
		product = ""
		uniq_vpc_id = ""
		url = ""
		path = ""
		method = ""
		upstream_id = ""
		cos_config {
			action = ""
			bucket_name = ""
			authorization = 
			path_match_mode = ""
		}

  }
  auth_relation_api_id = ""
  service_parameters {
		name = ""
		position = ""
		relevant_request_parameter_position = ""
		relevant_request_parameter_name = ""
		default_value = ""
		relevant_request_parameter_desc = ""
		relevant_request_parameter_type = ""

  }
  oauth_config {
		public_key = ""
		token_location = ""
		login_redirect_url = ""

  }
  response_error_codes {
		code = 
		msg = ""
		desc = ""
		converted_code = 
		need_convert = 

  }
  target_namespace_id = ""
  user_type = ""
  is_base64_encoded = 
  event_bus_id = ""
  service_scf_function_type = ""
  e_i_a_m_app_type = ""
  e_i_a_m_auth_type = ""
  token_timeout = 
  e_i_a_m_app_id = ""
  owner = ""
  tags = {
    "createdBy" = "terraform"
  }
}

`
