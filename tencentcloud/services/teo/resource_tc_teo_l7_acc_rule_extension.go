package teo

import (
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	teo "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"
	teov20220901 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func TencentTeoL7RuleBranchBasicInfo(depth int) map[string]*schema.Schema {
	schemaMap := map[string]*schema.Schema{
		"condition": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "",
		},
		"actions": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"name": {
						Type:        schema.TypeString,
						Required:    true,
						Description: "",
					},
					"cache_parameters": {
						Type:        schema.TypeList,
						Optional:    true,
						MaxItems:    1,
						Description: "",
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"follow_origin": {
									Type:        schema.TypeList,
									Optional:    true,
									MaxItems:    1,
									Description: "",
									Elem: &schema.Resource{
										Schema: map[string]*schema.Schema{
											"switch": {
												Type:        schema.TypeString,
												Required:    true,
												Description: "",
											},
											"default_cache": {
												Type:        schema.TypeString,
												Optional:    true,
												Description: "",
											},
											"default_cache_strategy": {
												Type:        schema.TypeString,
												Optional:    true,
												Description: "",
											},
											"default_cache_time": {
												Type:        schema.TypeInt,
												Optional:    true,
												Description: "",
											},
										},
									},
								},
								"no_cache": {
									Type:        schema.TypeList,
									Optional:    true,
									MaxItems:    1,
									Description: "",
									Elem: &schema.Resource{
										Schema: map[string]*schema.Schema{
											"switch": {
												Type:        schema.TypeString,
												Required:    true,
												Description: "",
											},
										},
									},
								},
								"custom_time": {
									Type:        schema.TypeList,
									Optional:    true,
									MaxItems:    1,
									Description: "",
									Elem: &schema.Resource{
										Schema: map[string]*schema.Schema{
											"switch": {
												Type:        schema.TypeString,
												Optional:    true,
												Description: "",
											},
											"ignore_cache_control": {
												Type:        schema.TypeString,
												Optional:    true,
												Description: "",
											},
											"cache_time": {
												Type:        schema.TypeInt,
												Optional:    true,
												Description: "",
											},
										},
									},
								},
							},
						},
					},
					"cache_key_parameters": {
						Type:        schema.TypeList,
						Optional:    true,
						MaxItems:    1,
						Description: "",
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"full_url_cache": {
									Type:        schema.TypeString,
									Optional:    true,
									Description: "",
								},
								"query_string": {
									Type:        schema.TypeList,
									Optional:    true,
									MaxItems:    1,
									Description: "",
									Elem: &schema.Resource{
										Schema: map[string]*schema.Schema{
											"switch": {
												Type:        schema.TypeString,
												Optional:    true,
												Description: "",
											},
											"action": {
												Type:        schema.TypeString,
												Optional:    true,
												Description: "",
											},
											"values": {
												Type:        schema.TypeList,
												Optional:    true,
												Description: "",
												Elem: &schema.Schema{
													Type: schema.TypeString,
												},
											},
										},
									},
								},
								"ignore_case": {
									Type:        schema.TypeString,
									Optional:    true,
									Description: "",
								},
								"header": {
									Type:        schema.TypeList,
									Optional:    true,
									MaxItems:    1,
									Description: "",
									Elem: &schema.Resource{
										Schema: map[string]*schema.Schema{
											"switch": {
												Type:        schema.TypeString,
												Optional:    true,
												Description: "",
											},
											"values": {
												Type:        schema.TypeList,
												Optional:    true,
												Description: "",
												Elem: &schema.Schema{
													Type: schema.TypeString,
												},
											},
										},
									},
								},
								"scheme": {
									Type:        schema.TypeString,
									Optional:    true,
									Description: "",
								},
								"cookie": {
									Type:        schema.TypeList,
									Optional:    true,
									MaxItems:    1,
									Description: "",
									Elem: &schema.Resource{
										Schema: map[string]*schema.Schema{
											"switch": {
												Type:        schema.TypeString,
												Optional:    true,
												Description: "",
											},
											"action": {
												Type:        schema.TypeString,
												Optional:    true,
												Description: "",
											},
											"values": {
												Type:        schema.TypeList,
												Optional:    true,
												Description: "",
												Elem: &schema.Schema{
													Type: schema.TypeString,
												},
											},
										},
									},
								},
							},
						},
					},
					"cache_prefresh_parameters": {
						Type:        schema.TypeList,
						Optional:    true,
						MaxItems:    1,
						Description: "",
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"switch": {
									Type:        schema.TypeString,
									Optional:    true,
									Description: "",
								},
								"cache_time_percent": {
									Type:        schema.TypeInt,
									Optional:    true,
									Description: "",
								},
							},
						},
					},
					"access_url_redirect_parameters": {
						Type:        schema.TypeList,
						Optional:    true,
						MaxItems:    1,
						Description: "",
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"status_code": {
									Type:        schema.TypeInt,
									Optional:    true,
									Description: "",
								},
								"protocol": {
									Type:        schema.TypeString,
									Optional:    true,
									Description: "",
								},
								"host_name": {
									Type:        schema.TypeList,
									Optional:    true,
									MaxItems:    1,
									Description: "",
									Elem: &schema.Resource{
										Schema: map[string]*schema.Schema{
											"action": {
												Type:        schema.TypeString,
												Optional:    true,
												Description: "",
											},
											"value": {
												Type:        schema.TypeString,
												Optional:    true,
												Description: "",
											},
										},
									},
								},
								"url_path": {
									Type:        schema.TypeList,
									Optional:    true,
									MaxItems:    1,
									Description: "",
									Elem: &schema.Resource{
										Schema: map[string]*schema.Schema{
											"action": {
												Type:        schema.TypeString,
												Optional:    true,
												Description: "",
											},
											"regex": {
												Type:        schema.TypeString,
												Optional:    true,
												Description: "",
											},
											"value": {
												Type:        schema.TypeString,
												Optional:    true,
												Description: "",
											},
										},
									},
								},
								"query_string": {
									Type:        schema.TypeList,
									Optional:    true,
									MaxItems:    1,
									Description: "",
									Elem: &schema.Resource{
										Schema: map[string]*schema.Schema{
											"action": {
												Type:        schema.TypeString,
												Optional:    true,
												Description: "",
											},
										},
									},
								},
							},
						},
					},
					"upstream_url_rewrite_parameters": {
						Type:        schema.TypeList,
						Optional:    true,
						MaxItems:    1,
						Description: "",
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"type": {
									Type:        schema.TypeString,
									Optional:    true,
									Description: "",
								},
								"action": {
									Type:        schema.TypeString,
									Optional:    true,
									Description: "",
								},
								"value": {
									Type:        schema.TypeString,
									Optional:    true,
									Description: "",
								},
								"regex": {
									Type:        schema.TypeString,
									Optional:    true,
									Description: "",
								},
							},
						},
					},
					"quic_parameters": {
						Type:        schema.TypeList,
						Optional:    true,
						MaxItems:    1,
						Description: "",
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"switch": {
									Type:        schema.TypeString,
									Optional:    true,
									Description: "",
								},
							},
						},
					},
					"web_socket_parameters": {
						Type:        schema.TypeList,
						Optional:    true,
						MaxItems:    1,
						Description: "",
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"switch": {
									Type:        schema.TypeString,
									Optional:    true,
									Description: "",
								},
								"timeout": {
									Type:        schema.TypeInt,
									Optional:    true,
									Description: "",
								},
							},
						},
					},
					"authentication_parameters": {
						Type:        schema.TypeList,
						Optional:    true,
						MaxItems:    1,
						Description: "",
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"auth_type": {
									Type:        schema.TypeString,
									Optional:    true,
									Description: "",
								},
								"secret_key": {
									Type:        schema.TypeString,
									Optional:    true,
									Description: "",
								},
								"timeout": {
									Type:        schema.TypeInt,
									Optional:    true,
									Description: "",
								},
								"backup_secret_key": {
									Type:        schema.TypeString,
									Optional:    true,
									Description: "",
								},
								"auth_param": {
									Type:        schema.TypeString,
									Optional:    true,
									Description: "",
								},
								"time_param": {
									Type:        schema.TypeString,
									Optional:    true,
									Description: "",
								},
								"time_format": {
									Type:        schema.TypeString,
									Optional:    true,
									Description: "",
								},
							},
						},
					},
					"max_age_parameters": {
						Type:        schema.TypeList,
						Optional:    true,
						MaxItems:    1,
						Description: "",
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"follow_origin": {
									Type:        schema.TypeString,
									Optional:    true,
									Description: "",
								},
								"cache_time": {
									Type:        schema.TypeInt,
									Optional:    true,
									Description: "",
								},
							},
						},
					},
					"status_code_cache_parameters": {
						Type:        schema.TypeList,
						Optional:    true,
						MaxItems:    1,
						Description: "",
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"status_code_cache_params": {
									Type:        schema.TypeList,
									Optional:    true,
									Description: "",
									Elem: &schema.Resource{
										Schema: map[string]*schema.Schema{
											"status_code": {
												Type:        schema.TypeInt,
												Optional:    true,
												Description: "",
											},
											"cache_time": {
												Type:        schema.TypeInt,
												Optional:    true,
												Description: "",
											},
										},
									},
								},
							},
						},
					},
					"offline_cache_parameters": {
						Type:        schema.TypeList,
						Optional:    true,
						MaxItems:    1,
						Description: "",
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"switch": {
									Type:        schema.TypeString,
									Optional:    true,
									Description: "",
								},
							},
						},
					},
					"smart_routing_parameters": {
						Type:        schema.TypeList,
						Optional:    true,
						MaxItems:    1,
						Description: "",
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"switch": {
									Type:        schema.TypeString,
									Optional:    true,
									Description: "",
								},
							},
						},
					},
					"range_origin_pull_parameters": {
						Type:        schema.TypeList,
						Optional:    true,
						MaxItems:    1,
						Description: "",
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"switch": {
									Type:        schema.TypeString,
									Optional:    true,
									Description: "",
								},
							},
						},
					},
					"upstream_http2_parameters": {
						Type:        schema.TypeList,
						Optional:    true,
						MaxItems:    1,
						Description: "",
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"switch": {
									Type:        schema.TypeString,
									Optional:    true,
									Description: "",
								},
							},
						},
					},
					"host_header_parameters": {
						Type:        schema.TypeList,
						Optional:    true,
						MaxItems:    1,
						Description: "",
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"action": {
									Type:        schema.TypeString,
									Optional:    true,
									Description: "",
								},
								"server_name": {
									Type:        schema.TypeString,
									Optional:    true,
									Description: "",
								},
							},
						},
					},
					"force_redirect_https_parameters": {
						Type:        schema.TypeList,
						Optional:    true,
						MaxItems:    1,
						Description: "",
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"switch": {
									Type:        schema.TypeString,
									Optional:    true,
									Description: "",
								},
								"redirect_status_code": {
									Type:        schema.TypeInt,
									Optional:    true,
									Description: "",
								},
							},
						},
					},
					"origin_pull_protocol_parameters": {
						Type:        schema.TypeList,
						Optional:    true,
						MaxItems:    1,
						Description: "",
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"protocol": {
									Type:        schema.TypeString,
									Optional:    true,
									Description: "",
								},
							},
						},
					},
					"compression_parameters": {
						Type:        schema.TypeList,
						Optional:    true,
						MaxItems:    1,
						Description: "",
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"switch": {
									Type:        schema.TypeString,
									Optional:    true,
									Description: "",
								},
								"algorithms": {
									Type:        schema.TypeList,
									Optional:    true,
									Description: "",
									Elem: &schema.Schema{
										Type: schema.TypeString,
									},
								},
							},
						},
					},
					"hsts_parameters": {
						Type:        schema.TypeList,
						Optional:    true,
						MaxItems:    1,
						Description: "",
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"switch": {
									Type:        schema.TypeString,
									Optional:    true,
									Description: "",
								},
								"timeout": {
									Type:        schema.TypeInt,
									Optional:    true,
									Description: "",
								},
								"include_sub_domains": {
									Type:        schema.TypeString,
									Optional:    true,
									Description: "",
								},
								"preload": {
									Type:        schema.TypeString,
									Optional:    true,
									Description: "",
								},
							},
						},
					},
					"client_ip_header_parameters": {
						Type:        schema.TypeList,
						Optional:    true,
						MaxItems:    1,
						Description: "",
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"switch": {
									Type:        schema.TypeString,
									Optional:    true,
									Description: "",
								},
								"header_name": {
									Type:        schema.TypeString,
									Optional:    true,
									Description: "",
								},
							},
						},
					},
					"ocsp_stapling_parameters": {
						Type:        schema.TypeList,
						Optional:    true,
						MaxItems:    1,
						Description: "",
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"switch": {
									Type:        schema.TypeString,
									Optional:    true,
									Description: "",
								},
							},
						},
					},
					"http2_parameters": {
						Type:        schema.TypeList,
						Optional:    true,
						MaxItems:    1,
						Description: "",
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"switch": {
									Type:        schema.TypeString,
									Optional:    true,
									Description: "",
								},
							},
						},
					},
					"post_max_size_parameters": {
						Type:        schema.TypeList,
						Optional:    true,
						MaxItems:    1,
						Description: "",
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"switch": {
									Type:        schema.TypeString,
									Optional:    true,
									Description: "",
								},
								"max_size": {
									Type:        schema.TypeInt,
									Optional:    true,
									Description: "",
								},
							},
						},
					},
					"client_ip_country_parameters": {
						Type:        schema.TypeList,
						Optional:    true,
						MaxItems:    1,
						Description: "",
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"switch": {
									Type:        schema.TypeString,
									Optional:    true,
									Description: "",
								},
								"header_name": {
									Type:        schema.TypeString,
									Optional:    true,
									Description: "",
								},
							},
						},
					},
					"upstream_follow_redirect_parameters": {
						Type:        schema.TypeList,
						Optional:    true,
						MaxItems:    1,
						Description: "",
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"switch": {
									Type:        schema.TypeString,
									Optional:    true,
									Description: "",
								},
								"max_times": {
									Type:        schema.TypeInt,
									Optional:    true,
									Description: "",
								},
							},
						},
					},
					"upstream_request_parameters": {
						Type:        schema.TypeList,
						Optional:    true,
						MaxItems:    1,
						Description: "",
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"query_string": {
									Type:        schema.TypeList,
									Optional:    true,
									MaxItems:    1,
									Description: "",
									Elem: &schema.Resource{
										Schema: map[string]*schema.Schema{
											"switch": {
												Type:        schema.TypeString,
												Optional:    true,
												Description: "",
											},
											"action": {
												Type:        schema.TypeString,
												Optional:    true,
												Description: "",
											},
											"values": {
												Type:        schema.TypeList,
												Optional:    true,
												Description: "",
												Elem: &schema.Schema{
													Type: schema.TypeString,
												},
											},
										},
									},
								},
								"cookie": {
									Type:        schema.TypeList,
									Optional:    true,
									MaxItems:    1,
									Description: "",
									Elem: &schema.Resource{
										Schema: map[string]*schema.Schema{
											"switch": {
												Type:        schema.TypeString,
												Optional:    true,
												Description: "",
											},
											"action": {
												Type:        schema.TypeString,
												Optional:    true,
												Description: "",
											},
											"values": {
												Type:        schema.TypeList,
												Optional:    true,
												Description: "",
												Elem: &schema.Schema{
													Type: schema.TypeString,
												},
											},
										},
									},
								},
							},
						},
					},
					"tls_config_parameters": {
						Type:        schema.TypeList,
						Optional:    true,
						MaxItems:    1,
						Description: "",
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"version": {
									Type:        schema.TypeList,
									Optional:    true,
									Description: "",
									Elem: &schema.Schema{
										Type: schema.TypeString,
									},
								},
								"cipher_suite": {
									Type:        schema.TypeString,
									Optional:    true,
									Description: "",
								},
							},
						},
					},
					"modify_origin_parameters": {
						Type:        schema.TypeList,
						Optional:    true,
						MaxItems:    1,
						Description: "",
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"origin_type": {
									Type:        schema.TypeString,
									Optional:    true,
									Description: "",
								},
								"origin": {
									Type:        schema.TypeString,
									Optional:    true,
									Description: "",
								},
								"origin_protocol": {
									Type:        schema.TypeString,
									Optional:    true,
									Description: "",
								},
								"http_origin_port": {
									Type:         schema.TypeInt,
									Optional:     true,
									ValidateFunc: tccommon.ValidateIntegerInRange(1, 65535),
									Description:  "",
								},
								"https_origin_port": {
									Type:         schema.TypeInt,
									Optional:     true,
									ValidateFunc: tccommon.ValidateIntegerInRange(1, 65535),
									Description:  "",
								},
								"private_access": {
									Type:        schema.TypeString,
									Optional:    true,
									Description: "",
								},
								"private_parameters": {
									Type:        schema.TypeList,
									Optional:    true,
									MaxItems:    1,
									Description: "",
									Elem: &schema.Resource{
										Schema: map[string]*schema.Schema{
											"access_key_id": {
												Type:        schema.TypeString,
												Required:    true,
												Description: "",
											},
											"secret_access_key": {
												Type:        schema.TypeString,
												Required:    true,
												Description: "",
											},
											"signature_version": {
												Type:        schema.TypeString,
												Required:    true,
												Description: "",
											},
											"region": {
												Type:        schema.TypeString,
												Optional:    true,
												Description: "",
											},
										},
									},
								},
							},
						},
					},
					"http_upstream_timeout_parameters": {
						Type:        schema.TypeList,
						Optional:    true,
						MaxItems:    1,
						Description: "",
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"response_timeout": {
									Type:        schema.TypeInt,
									Optional:    true,
									Description: "",
								},
							},
						},
					},
					"http_response_parameters": {
						Type:        schema.TypeList,
						Optional:    true,
						MaxItems:    1,
						Description: "",
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"status_code": {
									Type:        schema.TypeInt,
									Optional:    true,
									Description: "",
								},
								"response_page": {
									Type:        schema.TypeString,
									Optional:    true,
									Description: "",
								},
							},
						},
					},
					"error_page_parameters": {
						Type:        schema.TypeList,
						Optional:    true,
						MaxItems:    1,
						Description: "",
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"error_page_params": {
									Type:        schema.TypeList,
									Optional:    true,
									Description: "",
									Elem: &schema.Resource{
										Schema: map[string]*schema.Schema{
											"status_code": {
												Type:        schema.TypeInt,
												Required:    true,
												Description: "",
											},
											"redirect_url": {
												Type:        schema.TypeString,
												Required:    true,
												Description: "",
											},
										},
									},
								},
							},
						},
					},
					"modify_response_header_parameters": {
						Type:        schema.TypeList,
						Optional:    true,
						MaxItems:    1,
						Description: "",
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"header_actions": {
									Type:        schema.TypeList,
									Optional:    true,
									Description: "",
									Elem: &schema.Resource{
										Schema: map[string]*schema.Schema{
											"action": {
												Type:        schema.TypeString,
												Required:    true,
												Description: "",
											},
											"name": {
												Type:        schema.TypeString,
												Required:    true,
												Description: "",
											},
											"value": {
												Type:        schema.TypeString,
												Optional:    true,
												Description: "",
											},
										},
									},
								},
							},
						},
					},
					"modify_request_header_parameters": {
						Type:        schema.TypeList,
						Optional:    true,
						MaxItems:    1,
						Description: "",
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"header_actions": {
									Type:        schema.TypeList,
									Optional:    true,
									Description: "",
									Elem: &schema.Resource{
										Schema: map[string]*schema.Schema{
											"action": {
												Type:        schema.TypeString,
												Required:    true,
												Description: "",
											},
											"name": {
												Type:        schema.TypeString,
												Required:    true,
												Description: "",
											},
											"value": {
												Type:        schema.TypeString,
												Optional:    true,
												Description: "",
											},
										},
									},
								},
							},
						},
					},
					"response_speed_limit_parameters": {
						Type:        schema.TypeList,
						Optional:    true,
						MaxItems:    1,
						Description: "",
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"mode": {
									Type:        schema.TypeString,
									Required:    true,
									Description: "",
								},
								"max_speed": {
									Type:        schema.TypeString,
									Required:    true,
									Description: "",
								},
								"start_at": {
									Type:        schema.TypeString,
									Optional:    true,
									Description: "",
								},
							},
						},
					},
					"set_content_identifier_parameters": {
						Type:        schema.TypeList,
						Optional:    true,
						MaxItems:    1,
						Description: "",
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"content_identifier": {
									Type:        schema.TypeString,
									Optional:    true,
									Description: "",
								},
							},
						},
					},
					"content_compression_parameters": {
						Type:        schema.TypeList,
						Optional:    true,
						MaxItems:    1,
						Description: "",
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"switch": {
									Type:        schema.TypeString,
									Required:    true,
									Description: "",
								},
							},
						},
					},
					"vary_parameters": {
						Type:        schema.TypeList,
						Optional:    true,
						MaxItems:    1,
						Description: "",
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"switch": {
									Type:        schema.TypeString,
									Required:    true,
									Description: "",
								},
							},
						},
					},
					"origin_authentication_parameters": {
						Type:        schema.TypeList,
						Optional:    true,
						MaxItems:    1,
						Description: "",
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"request_properties": {
									Type:        schema.TypeList,
									Required:    true,
									Description: "",
									Elem: &schema.Resource{
										Schema: map[string]*schema.Schema{
											"type": {
												Type:        schema.TypeString,
												Required:    true,
												Description: "",
											},
											"name": {
												Type:        schema.TypeString,
												Required:    true,
												Description: "",
											},
											"value": {
												Type:        schema.TypeString,
												Required:    true,
												Description: "",
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
	}

	if depth < 8 {
		schemaMap["sub_rules"] = &schema.Schema{
			Type:        schema.TypeList,
			Optional:    true,
			Description: "",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"branches": {
						Type:        schema.TypeList,
						Optional:    true,
						Description: "",
						Elem: &schema.Resource{
							Schema: TencentTeoL7RuleBranchBasicInfo(depth + 1),
						},
					},
					"description": {
						Type:        schema.TypeList,
						Optional:    true,
						Description: "",
						Elem: &schema.Schema{
							Type: schema.TypeString,
						},
					},
				},
			},
		}
	}

	return schemaMap
}

func resourceTencentCloudTeoL7AccRuleGetBranchs(rulesMap map[string]interface{}) []*teo.RuleBranch {
	ruleBranchs := []*teo.RuleBranch{}
	if v, ok := rulesMap["branches"]; ok {
		for _, item := range v.([]interface{}) {
			branchesMap := item.(map[string]interface{})
			ruleBranch := teov20220901.RuleBranch{}
			if v, ok := branchesMap["condition"].(string); ok && v != "" {
				ruleBranch.Condition = helper.String(v)
			}
			if v, ok := branchesMap["actions"]; ok {
				for _, item := range v.([]interface{}) {
					actionsMap := item.(map[string]interface{})
					ruleEngineAction := teov20220901.RuleEngineAction{}
					if v, ok := actionsMap["name"].(string); ok && v != "" {
						ruleEngineAction.Name = helper.String(v)
					}
					if cacheParametersMap, ok := helper.ConvertInterfacesHeadToMap(actionsMap["cache_parameters"]); ok {
						cacheParameters := teov20220901.CacheParameters{}
						if followOriginMap, ok := helper.ConvertInterfacesHeadToMap(cacheParametersMap["follow_origin"]); ok {
							followOrigin := teov20220901.FollowOrigin{}
							if v, ok := followOriginMap["switch"].(string); ok && v != "" {
								followOrigin.Switch = helper.String(v)
							}
							if v, ok := followOriginMap["default_cache"].(string); ok && v != "" {
								followOrigin.DefaultCache = helper.String(v)
							}
							if v, ok := followOriginMap["default_cache_strategy"].(string); ok && v != "" {
								followOrigin.DefaultCacheStrategy = helper.String(v)
							}
							if v, ok := followOriginMap["default_cache_time"].(int); ok {
								followOrigin.DefaultCacheTime = helper.IntInt64(v)
							}
							cacheParameters.FollowOrigin = &followOrigin
						}
						if noCacheMap, ok := helper.ConvertInterfacesHeadToMap(cacheParametersMap["no_cache"]); ok {
							noCache := teov20220901.NoCache{}
							if v, ok := noCacheMap["switch"].(string); ok && v != "" {
								noCache.Switch = helper.String(v)
							}
							cacheParameters.NoCache = &noCache
						}
						if customTimeMap, ok := helper.ConvertInterfacesHeadToMap(cacheParametersMap["custom_time"]); ok {
							customTime := teov20220901.CustomTime{}
							if v, ok := customTimeMap["switch"].(string); ok && v != "" {
								customTime.Switch = helper.String(v)
							}
							if v, ok := customTimeMap["ignore_cache_control"].(string); ok && v != "" {
								customTime.IgnoreCacheControl = helper.String(v)
							}
							if v, ok := customTimeMap["cache_time"].(int); ok {
								customTime.CacheTime = helper.IntInt64(v)
							}
							cacheParameters.CustomTime = &customTime
						}
						ruleEngineAction.CacheParameters = &cacheParameters
					}
					if cacheKeyParametersMap, ok := helper.ConvertInterfacesHeadToMap(actionsMap["cache_key_parameters"]); ok {
						cacheKeyParameters := teov20220901.CacheKeyParameters{}
						if v, ok := cacheKeyParametersMap["full_url_cache"].(string); ok && v != "" {
							cacheKeyParameters.FullURLCache = helper.String(v)
						}
						if queryStringMap, ok := helper.ConvertInterfacesHeadToMap(cacheKeyParametersMap["query_string"]); ok {
							cacheKeyQueryString := teov20220901.CacheKeyQueryString{}
							if v, ok := queryStringMap["switch"].(string); ok && v != "" {
								cacheKeyQueryString.Switch = helper.String(v)
							}
							if v, ok := queryStringMap["action"].(string); ok && v != "" {
								cacheKeyQueryString.Action = helper.String(v)
							}
							if v, ok := queryStringMap["values"]; ok {
								valuesSet := v.([]interface{})
								for i := range valuesSet {
									values := valuesSet[i].(string)
									cacheKeyQueryString.Values = append(cacheKeyQueryString.Values, helper.String(values))
								}
							}
							cacheKeyParameters.QueryString = &cacheKeyQueryString
						}
						if v, ok := cacheKeyParametersMap["ignore_case"].(string); ok && v != "" {
							cacheKeyParameters.IgnoreCase = helper.String(v)
						}
						if headerMap, ok := helper.ConvertInterfacesHeadToMap(cacheKeyParametersMap["header"]); ok {
							cacheKeyHeader := teov20220901.CacheKeyHeader{}
							if v, ok := headerMap["switch"].(string); ok && v != "" {
								cacheKeyHeader.Switch = helper.String(v)
							}
							if v, ok := headerMap["values"]; ok {
								valuesSet := v.([]interface{})
								for i := range valuesSet {
									values := valuesSet[i].(string)
									cacheKeyHeader.Values = append(cacheKeyHeader.Values, helper.String(values))
								}
							}
							cacheKeyParameters.Header = &cacheKeyHeader
						}
						if v, ok := cacheKeyParametersMap["scheme"].(string); ok && v != "" {
							cacheKeyParameters.Scheme = helper.String(v)
						}
						if cookieMap, ok := helper.ConvertInterfacesHeadToMap(cacheKeyParametersMap["cookie"]); ok {
							cacheKeyCookie := teov20220901.CacheKeyCookie{}
							if v, ok := cookieMap["switch"].(string); ok && v != "" {
								cacheKeyCookie.Switch = helper.String(v)
							}
							if v, ok := cookieMap["action"].(string); ok && v != "" {
								cacheKeyCookie.Action = helper.String(v)
							}
							if v, ok := cookieMap["values"]; ok {
								valuesSet := v.([]interface{})
								for i := range valuesSet {
									values := valuesSet[i].(string)
									cacheKeyCookie.Values = append(cacheKeyCookie.Values, helper.String(values))
								}
							}
							cacheKeyParameters.Cookie = &cacheKeyCookie
						}
						ruleEngineAction.CacheKeyParameters = &cacheKeyParameters
					}
					if cachePrefreshParametersMap, ok := helper.ConvertInterfacesHeadToMap(actionsMap["cache_prefresh_parameters"]); ok {
						cachePrefreshParameters := teov20220901.CachePrefreshParameters{}
						if v, ok := cachePrefreshParametersMap["switch"].(string); ok && v != "" {
							cachePrefreshParameters.Switch = helper.String(v)
						}
						if v, ok := cachePrefreshParametersMap["cache_time_percent"].(int); ok {
							cachePrefreshParameters.CacheTimePercent = helper.IntInt64(v)
						}
						ruleEngineAction.CachePrefreshParameters = &cachePrefreshParameters
					}
					if accessURLRedirectParametersMap, ok := helper.ConvertInterfacesHeadToMap(actionsMap["access_url_redirect_parameters"]); ok {
						accessURLRedirectParameters := teov20220901.AccessURLRedirectParameters{}
						if v, ok := accessURLRedirectParametersMap["status_code"].(int); ok {
							accessURLRedirectParameters.StatusCode = helper.IntInt64(v)
						}
						if v, ok := accessURLRedirectParametersMap["protocol"].(string); ok && v != "" {
							accessURLRedirectParameters.Protocol = helper.String(v)
						}
						if hostNameMap, ok := helper.ConvertInterfacesHeadToMap(accessURLRedirectParametersMap["host_name"]); ok {
							hostName := teov20220901.HostName{}
							if v, ok := hostNameMap["action"].(string); ok && v != "" {
								hostName.Action = helper.String(v)
							}
							if v, ok := hostNameMap["value"].(string); ok && v != "" {
								hostName.Value = helper.String(v)
							}
							accessURLRedirectParameters.HostName = &hostName
						}
						if uRLPathMap, ok := helper.ConvertInterfacesHeadToMap(accessURLRedirectParametersMap["url_path"]); ok {
							uRLPath := teov20220901.URLPath{}
							if v, ok := uRLPathMap["action"].(string); ok && v != "" {
								uRLPath.Action = helper.String(v)
							}
							if v, ok := uRLPathMap["regex"].(string); ok && v != "" {
								uRLPath.Regex = helper.String(v)
							}
							if v, ok := uRLPathMap["value"].(string); ok && v != "" {
								uRLPath.Value = helper.String(v)
							}
							accessURLRedirectParameters.URLPath = &uRLPath
						}
						if queryStringMap, ok := helper.ConvertInterfacesHeadToMap(accessURLRedirectParametersMap["query_string"]); ok {
							accessURLRedirectQueryString := teov20220901.AccessURLRedirectQueryString{}
							if v, ok := queryStringMap["action"].(string); ok && v != "" {
								accessURLRedirectQueryString.Action = helper.String(v)
							}
							accessURLRedirectParameters.QueryString = &accessURLRedirectQueryString
						}
						ruleEngineAction.AccessURLRedirectParameters = &accessURLRedirectParameters
					}
					if upstreamURLRewriteParametersMap, ok := helper.ConvertInterfacesHeadToMap(actionsMap["upstream_url_rewrite_parameters"]); ok {
						upstreamURLRewriteParameters := teov20220901.UpstreamURLRewriteParameters{}
						if v, ok := upstreamURLRewriteParametersMap["type"].(string); ok && v != "" {
							upstreamURLRewriteParameters.Type = helper.String(v)
						}
						if v, ok := upstreamURLRewriteParametersMap["action"].(string); ok && v != "" {
							upstreamURLRewriteParameters.Action = helper.String(v)
						}
						if v, ok := upstreamURLRewriteParametersMap["value"].(string); ok && v != "" {
							upstreamURLRewriteParameters.Value = helper.String(v)
						}
						if v, ok := upstreamURLRewriteParametersMap["regex"].(string); ok && v != "" {
							upstreamURLRewriteParameters.Regex = helper.String(v)
						}
						ruleEngineAction.UpstreamURLRewriteParameters = &upstreamURLRewriteParameters
					}
					if qUICParametersMap, ok := helper.ConvertInterfacesHeadToMap(actionsMap["quic_parameters"]); ok {
						qUICParameters := teov20220901.QUICParameters{}
						if v, ok := qUICParametersMap["switch"].(string); ok && v != "" {
							qUICParameters.Switch = helper.String(v)
						}
						ruleEngineAction.QUICParameters = &qUICParameters
					}
					if webSocketParametersMap, ok := helper.ConvertInterfacesHeadToMap(actionsMap["web_socket_parameters"]); ok {
						webSocketParameters := teov20220901.WebSocketParameters{}
						if v, ok := webSocketParametersMap["switch"].(string); ok && v != "" {
							webSocketParameters.Switch = helper.String(v)
						}
						if v, ok := webSocketParametersMap["timeout"].(int); ok {
							webSocketParameters.Timeout = helper.IntInt64(v)
						}
						ruleEngineAction.WebSocketParameters = &webSocketParameters
					}
					if authenticationParametersMap, ok := helper.ConvertInterfacesHeadToMap(actionsMap["authentication_parameters"]); ok {
						authenticationParameters := teov20220901.AuthenticationParameters{}
						if v, ok := authenticationParametersMap["auth_type"].(string); ok && v != "" {
							authenticationParameters.AuthType = helper.String(v)
						}
						if v, ok := authenticationParametersMap["secret_key"].(string); ok && v != "" {
							authenticationParameters.SecretKey = helper.String(v)
						}
						if v, ok := authenticationParametersMap["timeout"].(int); ok {
							authenticationParameters.Timeout = helper.IntInt64(v)
						}
						if v, ok := authenticationParametersMap["backup_secret_key"].(string); ok && v != "" {
							authenticationParameters.BackupSecretKey = helper.String(v)
						}
						if v, ok := authenticationParametersMap["auth_param"].(string); ok && v != "" {
							authenticationParameters.AuthParam = helper.String(v)
						}
						if v, ok := authenticationParametersMap["time_param"].(string); ok && v != "" {
							authenticationParameters.TimeParam = helper.String(v)
						}
						if v, ok := authenticationParametersMap["time_format"].(string); ok && v != "" {
							authenticationParameters.TimeFormat = helper.String(v)
						}
						ruleEngineAction.AuthenticationParameters = &authenticationParameters
					}
					if maxAgeParametersMap, ok := helper.ConvertInterfacesHeadToMap(actionsMap["max_age_parameters"]); ok {
						maxAgeParameters := teov20220901.MaxAgeParameters{}
						if v, ok := maxAgeParametersMap["follow_origin"].(string); ok && v != "" {
							maxAgeParameters.FollowOrigin = helper.String(v)
						}
						if v, ok := maxAgeParametersMap["cache_time"].(int); ok {
							maxAgeParameters.CacheTime = helper.IntInt64(v)
						}
						ruleEngineAction.MaxAgeParameters = &maxAgeParameters
					}
					if statusCodeCacheParametersMap, ok := helper.ConvertInterfacesHeadToMap(actionsMap["status_code_cache_parameters"]); ok {
						statusCodeCacheParameters := teov20220901.StatusCodeCacheParameters{}
						if v, ok := statusCodeCacheParametersMap["status_code_cache_params"]; ok {
							for _, item := range v.([]interface{}) {
								statusCodeCacheParamsMap := item.(map[string]interface{})
								statusCodeCacheParam := teov20220901.StatusCodeCacheParam{}
								if v, ok := statusCodeCacheParamsMap["status_code"].(int); ok {
									statusCodeCacheParam.StatusCode = helper.IntInt64(v)
								}
								if v, ok := statusCodeCacheParamsMap["cache_time"].(int); ok {
									statusCodeCacheParam.CacheTime = helper.IntInt64(v)
								}
								statusCodeCacheParameters.StatusCodeCacheParams = append(statusCodeCacheParameters.StatusCodeCacheParams, &statusCodeCacheParam)
							}
						}
						ruleEngineAction.StatusCodeCacheParameters = &statusCodeCacheParameters
					}
					if offlineCacheParametersMap, ok := helper.ConvertInterfacesHeadToMap(actionsMap["offline_cache_parameters"]); ok {
						offlineCacheParameters := teov20220901.OfflineCacheParameters{}
						if v, ok := offlineCacheParametersMap["switch"].(string); ok && v != "" {
							offlineCacheParameters.Switch = helper.String(v)
						}
						ruleEngineAction.OfflineCacheParameters = &offlineCacheParameters
					}
					if smartRoutingParametersMap, ok := helper.ConvertInterfacesHeadToMap(actionsMap["smart_routing_parameters"]); ok {
						smartRoutingParameters := teov20220901.SmartRoutingParameters{}
						if v, ok := smartRoutingParametersMap["switch"].(string); ok && v != "" {
							smartRoutingParameters.Switch = helper.String(v)
						}
						ruleEngineAction.SmartRoutingParameters = &smartRoutingParameters
					}
					if rangeOriginPullParametersMap, ok := helper.ConvertInterfacesHeadToMap(actionsMap["range_origin_pull_parameters"]); ok {
						rangeOriginPullParameters := teov20220901.RangeOriginPullParameters{}
						if v, ok := rangeOriginPullParametersMap["switch"].(string); ok && v != "" {
							rangeOriginPullParameters.Switch = helper.String(v)
						}
						ruleEngineAction.RangeOriginPullParameters = &rangeOriginPullParameters
					}
					if upstreamHTTP2ParametersMap, ok := helper.ConvertInterfacesHeadToMap(actionsMap["upstream_http2_parameters"]); ok {
						upstreamHTTP2Parameters := teov20220901.UpstreamHTTP2Parameters{}
						if v, ok := upstreamHTTP2ParametersMap["switch"].(string); ok && v != "" {
							upstreamHTTP2Parameters.Switch = helper.String(v)
						}
						ruleEngineAction.UpstreamHTTP2Parameters = &upstreamHTTP2Parameters
					}
					if hostHeaderParametersMap, ok := helper.ConvertInterfacesHeadToMap(actionsMap["host_header_parameters"]); ok {
						hostHeaderParameters := teov20220901.HostHeaderParameters{}
						if v, ok := hostHeaderParametersMap["action"].(string); ok && v != "" {
							hostHeaderParameters.Action = helper.String(v)
						}
						if v, ok := hostHeaderParametersMap["server_name"].(string); ok && v != "" {
							hostHeaderParameters.ServerName = helper.String(v)
						}
						ruleEngineAction.HostHeaderParameters = &hostHeaderParameters
					}
					if forceRedirectHTTPSParametersMap, ok := helper.ConvertInterfacesHeadToMap(actionsMap["force_redirect_https_parameters"]); ok {
						forceRedirectHTTPSParameters := teov20220901.ForceRedirectHTTPSParameters{}
						if v, ok := forceRedirectHTTPSParametersMap["switch"].(string); ok && v != "" {
							forceRedirectHTTPSParameters.Switch = helper.String(v)
						}
						if v, ok := forceRedirectHTTPSParametersMap["redirect_status_code"].(int); ok {
							forceRedirectHTTPSParameters.RedirectStatusCode = helper.IntInt64(v)
						}
						ruleEngineAction.ForceRedirectHTTPSParameters = &forceRedirectHTTPSParameters
					}
					if originPullProtocolParametersMap, ok := helper.ConvertInterfacesHeadToMap(actionsMap["origin_pull_protocol_parameters"]); ok {
						originPullProtocolParameters := teov20220901.OriginPullProtocolParameters{}
						if v, ok := originPullProtocolParametersMap["protocol"].(string); ok && v != "" {
							originPullProtocolParameters.Protocol = helper.String(v)
						}
						ruleEngineAction.OriginPullProtocolParameters = &originPullProtocolParameters
					}
					if compressionParametersMap, ok := helper.ConvertInterfacesHeadToMap(actionsMap["compression_parameters"]); ok {
						compressionParameters := teov20220901.CompressionParameters{}
						if v, ok := compressionParametersMap["switch"].(string); ok && v != "" {
							compressionParameters.Switch = helper.String(v)
						}
						if v, ok := compressionParametersMap["algorithms"]; ok {
							algorithmsSet := v.([]interface{})
							for i := range algorithmsSet {
								algorithms := algorithmsSet[i].(string)
								compressionParameters.Algorithms = append(compressionParameters.Algorithms, helper.String(algorithms))
							}
						}
						ruleEngineAction.CompressionParameters = &compressionParameters
					}
					if hSTSParametersMap, ok := helper.ConvertInterfacesHeadToMap(actionsMap["hsts_parameters"]); ok {
						hSTSParameters := teov20220901.HSTSParameters{}
						if v, ok := hSTSParametersMap["switch"].(string); ok && v != "" {
							hSTSParameters.Switch = helper.String(v)
						}
						if v, ok := hSTSParametersMap["timeout"].(int); ok {
							hSTSParameters.Timeout = helper.IntInt64(v)
						}
						if v, ok := hSTSParametersMap["include_sub_domains"].(string); ok && v != "" {
							hSTSParameters.IncludeSubDomains = helper.String(v)
						}
						if v, ok := hSTSParametersMap["preload"].(string); ok && v != "" {
							hSTSParameters.Preload = helper.String(v)
						}
						ruleEngineAction.HSTSParameters = &hSTSParameters
					}
					if clientIPHeaderParametersMap, ok := helper.ConvertInterfacesHeadToMap(actionsMap["client_ip_header_parameters"]); ok {
						clientIPHeaderParameters := teov20220901.ClientIPHeaderParameters{}
						if v, ok := clientIPHeaderParametersMap["switch"].(string); ok && v != "" {
							clientIPHeaderParameters.Switch = helper.String(v)
						}
						if v, ok := clientIPHeaderParametersMap["header_name"].(string); ok && v != "" {
							clientIPHeaderParameters.HeaderName = helper.String(v)
						}
						ruleEngineAction.ClientIPHeaderParameters = &clientIPHeaderParameters
					}
					if oCSPStaplingParametersMap, ok := helper.ConvertInterfacesHeadToMap(actionsMap["ocsp_stapling_parameters"]); ok {
						oCSPStaplingParameters := teov20220901.OCSPStaplingParameters{}
						if v, ok := oCSPStaplingParametersMap["switch"].(string); ok && v != "" {
							oCSPStaplingParameters.Switch = helper.String(v)
						}
						ruleEngineAction.OCSPStaplingParameters = &oCSPStaplingParameters
					}
					if hTTP2ParametersMap, ok := helper.ConvertInterfacesHeadToMap(actionsMap["http2_parameters"]); ok {
						hTTP2Parameters := teov20220901.HTTP2Parameters{}
						if v, ok := hTTP2ParametersMap["switch"].(string); ok && v != "" {
							hTTP2Parameters.Switch = helper.String(v)
						}
						ruleEngineAction.HTTP2Parameters = &hTTP2Parameters
					}
					if postMaxSizeParametersMap, ok := helper.ConvertInterfacesHeadToMap(actionsMap["post_max_size_parameters"]); ok {
						postMaxSizeParameters := teov20220901.PostMaxSizeParameters{}
						if v, ok := postMaxSizeParametersMap["switch"].(string); ok && v != "" {
							postMaxSizeParameters.Switch = helper.String(v)
						}
						if v, ok := postMaxSizeParametersMap["max_size"].(int); ok {
							postMaxSizeParameters.MaxSize = helper.IntInt64(v)
						}
						ruleEngineAction.PostMaxSizeParameters = &postMaxSizeParameters
					}
					if clientIPCountryParametersMap, ok := helper.ConvertInterfacesHeadToMap(actionsMap["client_ip_country_parameters"]); ok {
						clientIPCountryParameters := teov20220901.ClientIPCountryParameters{}
						if v, ok := clientIPCountryParametersMap["switch"].(string); ok && v != "" {
							clientIPCountryParameters.Switch = helper.String(v)
						}
						if v, ok := clientIPCountryParametersMap["header_name"].(string); ok && v != "" {
							clientIPCountryParameters.HeaderName = helper.String(v)
						}
						ruleEngineAction.ClientIPCountryParameters = &clientIPCountryParameters
					}
					if upstreamFollowRedirectParametersMap, ok := helper.ConvertInterfacesHeadToMap(actionsMap["upstream_follow_redirect_parameters"]); ok {
						upstreamFollowRedirectParameters := teov20220901.UpstreamFollowRedirectParameters{}
						if v, ok := upstreamFollowRedirectParametersMap["switch"].(string); ok && v != "" {
							upstreamFollowRedirectParameters.Switch = helper.String(v)
						}
						if v, ok := upstreamFollowRedirectParametersMap["max_times"].(int); ok {
							upstreamFollowRedirectParameters.MaxTimes = helper.IntInt64(v)
						}
						ruleEngineAction.UpstreamFollowRedirectParameters = &upstreamFollowRedirectParameters
					}
					if upstreamRequestParametersMap, ok := helper.ConvertInterfacesHeadToMap(actionsMap["upstream_request_parameters"]); ok {
						upstreamRequestParameters := teov20220901.UpstreamRequestParameters{}
						if queryStringMap, ok := helper.ConvertInterfacesHeadToMap(upstreamRequestParametersMap["query_string"]); ok {
							upstreamRequestQueryString := teov20220901.UpstreamRequestQueryString{}
							if v, ok := queryStringMap["switch"].(string); ok && v != "" {
								upstreamRequestQueryString.Switch = helper.String(v)
							}
							if v, ok := queryStringMap["action"].(string); ok && v != "" {
								upstreamRequestQueryString.Action = helper.String(v)
							}
							if v, ok := queryStringMap["values"]; ok {
								valuesSet := v.([]interface{})
								for i := range valuesSet {
									values := valuesSet[i].(string)
									upstreamRequestQueryString.Values = append(upstreamRequestQueryString.Values, helper.String(values))
								}
							}
							upstreamRequestParameters.QueryString = &upstreamRequestQueryString
						}
						if cookieMap, ok := helper.ConvertInterfacesHeadToMap(upstreamRequestParametersMap["cookie"]); ok {
							upstreamRequestCookie := teov20220901.UpstreamRequestCookie{}
							if v, ok := cookieMap["switch"].(string); ok && v != "" {
								upstreamRequestCookie.Switch = helper.String(v)
							}
							if v, ok := cookieMap["action"].(string); ok && v != "" {
								upstreamRequestCookie.Action = helper.String(v)
							}
							if v, ok := cookieMap["values"]; ok {
								valuesSet := v.([]interface{})
								for i := range valuesSet {
									values := valuesSet[i].(string)
									upstreamRequestCookie.Values = append(upstreamRequestCookie.Values, helper.String(values))
								}
							}
							upstreamRequestParameters.Cookie = &upstreamRequestCookie
						}
						ruleEngineAction.UpstreamRequestParameters = &upstreamRequestParameters
					}
					if tLSConfigParametersMap, ok := helper.ConvertInterfacesHeadToMap(actionsMap["tls_config_parameters"]); ok {
						tLSConfigParameters := teov20220901.TLSConfigParameters{}
						if v, ok := tLSConfigParametersMap["version"]; ok {
							versionSet := v.([]interface{})
							for i := range versionSet {
								version := versionSet[i].(string)
								tLSConfigParameters.Version = append(tLSConfigParameters.Version, helper.String(version))
							}
						}
						if v, ok := tLSConfigParametersMap["cipher_suite"].(string); ok && v != "" {
							tLSConfigParameters.CipherSuite = helper.String(v)
						}
						ruleEngineAction.TLSConfigParameters = &tLSConfigParameters
					}
					if modifyOriginParametersMap, ok := helper.ConvertInterfacesHeadToMap(actionsMap["modify_origin_parameters"]); ok {
						modifyOriginParameters := teov20220901.ModifyOriginParameters{}
						if v, ok := modifyOriginParametersMap["origin_type"].(string); ok && v != "" {
							modifyOriginParameters.OriginType = helper.String(v)
						}
						if v, ok := modifyOriginParametersMap["origin"].(string); ok && v != "" {
							modifyOriginParameters.Origin = helper.String(v)
						}
						if v, ok := modifyOriginParametersMap["origin_protocol"].(string); ok && v != "" {
							modifyOriginParameters.OriginProtocol = helper.String(v)
						}
						if v, ok := modifyOriginParametersMap["http_origin_port"].(int); ok && v != 0 {
							modifyOriginParameters.HTTPOriginPort = helper.IntInt64(v)
						}
						if v, ok := modifyOriginParametersMap["https_origin_port"].(int); ok && v != 0 {
							modifyOriginParameters.HTTPSOriginPort = helper.IntInt64(v)
						}
						if v, ok := modifyOriginParametersMap["private_access"].(string); ok && v != "" {
							modifyOriginParameters.PrivateAccess = helper.String(v)
						}
						if privateParametersMap, ok := helper.ConvertInterfacesHeadToMap(modifyOriginParametersMap["private_parameters"]); ok {
							originPrivateParameters := teov20220901.OriginPrivateParameters{}
							if v, ok := privateParametersMap["access_key_id"].(string); ok && v != "" {
								originPrivateParameters.AccessKeyId = helper.String(v)
							}
							if v, ok := privateParametersMap["secret_access_key"].(string); ok && v != "" {
								originPrivateParameters.SecretAccessKey = helper.String(v)
							}
							if v, ok := privateParametersMap["signature_version"].(string); ok && v != "" {
								originPrivateParameters.SignatureVersion = helper.String(v)
							}
							if v, ok := privateParametersMap["region"].(string); ok && v != "" {
								originPrivateParameters.Region = helper.String(v)
							}
							modifyOriginParameters.PrivateParameters = &originPrivateParameters
						}
						ruleEngineAction.ModifyOriginParameters = &modifyOriginParameters
					}
					if hTTPUpstreamTimeoutParametersMap, ok := helper.ConvertInterfacesHeadToMap(actionsMap["http_upstream_timeout_parameters"]); ok {
						hTTPUpstreamTimeoutParameters := teov20220901.HTTPUpstreamTimeoutParameters{}
						if v, ok := hTTPUpstreamTimeoutParametersMap["response_timeout"].(int); ok {
							hTTPUpstreamTimeoutParameters.ResponseTimeout = helper.IntInt64(v)
						}
						ruleEngineAction.HTTPUpstreamTimeoutParameters = &hTTPUpstreamTimeoutParameters
					}
					if httpResponseParametersMap, ok := helper.ConvertInterfacesHeadToMap(actionsMap["http_response_parameters"]); ok {
						hTTPResponseParameters := teov20220901.HTTPResponseParameters{}
						if v, ok := httpResponseParametersMap["status_code"].(int); ok {
							hTTPResponseParameters.StatusCode = helper.IntInt64(v)
						}
						if v, ok := httpResponseParametersMap["response_page"].(string); ok && v != "" {
							hTTPResponseParameters.ResponsePage = helper.String(v)
						}
						ruleEngineAction.HttpResponseParameters = &hTTPResponseParameters
					}
					if errorPageParametersMap, ok := helper.ConvertInterfacesHeadToMap(actionsMap["error_page_parameters"]); ok {
						errorPageParameters := teov20220901.ErrorPageParameters{}
						if v, ok := errorPageParametersMap["error_page_params"]; ok {
							for _, item := range v.([]interface{}) {
								errorPageParamsMap := item.(map[string]interface{})
								errorPage := teov20220901.ErrorPage{}
								if v, ok := errorPageParamsMap["status_code"].(int); ok {
									errorPage.StatusCode = helper.IntInt64(v)
								}
								if v, ok := errorPageParamsMap["redirect_url"].(string); ok && v != "" {
									errorPage.RedirectURL = helper.String(v)
								}
								errorPageParameters.ErrorPageParams = append(errorPageParameters.ErrorPageParams, &errorPage)
							}
						}
						ruleEngineAction.ErrorPageParameters = &errorPageParameters
					}
					if modifyResponseHeaderParametersMap, ok := helper.ConvertInterfacesHeadToMap(actionsMap["modify_response_header_parameters"]); ok {
						modifyResponseHeaderParameters := teov20220901.ModifyResponseHeaderParameters{}
						if v, ok := modifyResponseHeaderParametersMap["header_actions"]; ok {
							for _, item := range v.([]interface{}) {
								headerActionsMap := item.(map[string]interface{})
								headerAction := teov20220901.HeaderAction{}
								if v, ok := headerActionsMap["action"].(string); ok && v != "" {
									headerAction.Action = helper.String(v)
								}
								if v, ok := headerActionsMap["name"].(string); ok && v != "" {
									headerAction.Name = helper.String(v)
								}
								if v, ok := headerActionsMap["value"].(string); ok && v != "" {
									headerAction.Value = helper.String(v)
								}
								modifyResponseHeaderParameters.HeaderActions = append(modifyResponseHeaderParameters.HeaderActions, &headerAction)
							}
						}
						ruleEngineAction.ModifyResponseHeaderParameters = &modifyResponseHeaderParameters
					}
					if modifyRequestHeaderParametersMap, ok := helper.ConvertInterfacesHeadToMap(actionsMap["modify_request_header_parameters"]); ok {
						modifyRequestHeaderParameters := teov20220901.ModifyRequestHeaderParameters{}
						if v, ok := modifyRequestHeaderParametersMap["header_actions"]; ok {
							for _, item := range v.([]interface{}) {
								headerActionsMap := item.(map[string]interface{})
								headerAction := teov20220901.HeaderAction{}
								if v, ok := headerActionsMap["action"].(string); ok && v != "" {
									headerAction.Action = helper.String(v)
								}
								if v, ok := headerActionsMap["name"].(string); ok && v != "" {
									headerAction.Name = helper.String(v)
								}
								if v, ok := headerActionsMap["value"].(string); ok && v != "" {
									headerAction.Value = helper.String(v)
								}
								modifyRequestHeaderParameters.HeaderActions = append(modifyRequestHeaderParameters.HeaderActions, &headerAction)
							}
						}
						ruleEngineAction.ModifyRequestHeaderParameters = &modifyRequestHeaderParameters
					}
					if responseSpeedLimitParametersMap, ok := helper.ConvertInterfacesHeadToMap(actionsMap["response_speed_limit_parameters"]); ok {
						responseSpeedLimitParameters := teov20220901.ResponseSpeedLimitParameters{}
						if v, ok := responseSpeedLimitParametersMap["mode"].(string); ok && v != "" {
							responseSpeedLimitParameters.Mode = helper.String(v)
						}
						if v, ok := responseSpeedLimitParametersMap["max_speed"].(string); ok && v != "" {
							responseSpeedLimitParameters.MaxSpeed = helper.String(v)
						}
						if v, ok := responseSpeedLimitParametersMap["start_at"].(string); ok && v != "" {
							responseSpeedLimitParameters.StartAt = helper.String(v)
						}
						ruleEngineAction.ResponseSpeedLimitParameters = &responseSpeedLimitParameters
					}
					if setContentIdentifierParametersMap, ok := helper.ConvertInterfacesHeadToMap(actionsMap["set_content_identifier_parameters"]); ok {
						setContentIdentifierParameters := teov20220901.SetContentIdentifierParameters{}
						if v, ok := setContentIdentifierParametersMap["content_identifier"].(string); ok && v != "" {
							setContentIdentifierParameters.ContentIdentifier = helper.String(v)
						}
						ruleEngineAction.SetContentIdentifierParameters = &setContentIdentifierParameters
					}
					if contentCompressionParametersMap, ok := helper.ConvertInterfacesHeadToMap(actionsMap["content_compression_parameters"]); ok {
						contentCompressionParameters := teov20220901.ContentCompressionParameters{}
						if v, ok := contentCompressionParametersMap["switch"].(string); ok && v != "" {
							contentCompressionParameters.Switch = helper.String(v)
						}
						ruleEngineAction.ContentCompressionParameters = &contentCompressionParameters
					}
					if varyParametersMap, ok := helper.ConvertInterfacesHeadToMap(actionsMap["vary_parameters"]); ok {
						varyParameters := teov20220901.VaryParameters{}
						if v, ok := varyParametersMap["switch"].(string); ok && v != "" {
							varyParameters.Switch = helper.String(v)
						}
						ruleEngineAction.VaryParameters = &varyParameters
					}
					if originAuthenticationParametersMap, ok := helper.ConvertInterfacesHeadToMap(actionsMap["origin_authentication_parameters"]); ok {
						originAuthenticationParameters := teov20220901.OriginAuthenticationParameters{}
						if v, ok := originAuthenticationParametersMap["request_properties"]; ok {
							for _, item := range v.([]interface{}) {
								propMap := item.(map[string]interface{})
								prop := teov20220901.OriginAuthenticationRequestProperties{}
								if v, ok := propMap["type"].(string); ok && v != "" {
									prop.Type = helper.String(v)
								}
								if v, ok := propMap["name"].(string); ok && v != "" {
									prop.Name = helper.String(v)
								}
								if v, ok := propMap["value"].(string); ok && v != "" {
									prop.Value = helper.String(v)
								}
								originAuthenticationParameters.RequestProperties = append(originAuthenticationParameters.RequestProperties, &prop)
							}
						}
						ruleEngineAction.OriginAuthenticationParameters = &originAuthenticationParameters
					}
					ruleBranch.Actions = append(ruleBranch.Actions, &ruleEngineAction)
				}
			}

			if v, ok := branchesMap["sub_rules"]; ok {
				for _, item := range v.([]interface{}) {
					subRulesMap := item.(map[string]interface{})
					ruleEngineSubRule := teov20220901.RuleEngineSubRule{}
					if _, ok := subRulesMap["branches"]; ok {
						branchs := resourceTencentCloudTeoL7AccRuleGetBranchs(subRulesMap)
						ruleEngineSubRule.Branches = branchs
					}
					if v, ok := subRulesMap["description"]; ok {
						descriptionSet := v.([]interface{})
						for i := range descriptionSet {
							description := descriptionSet[i].(string)
							ruleEngineSubRule.Description = append(ruleEngineSubRule.Description, helper.String(description))
						}
					}
					ruleBranch.SubRules = append(ruleBranch.SubRules, &ruleEngineSubRule)
				}
			}
			ruleBranchs = append(ruleBranchs, &ruleBranch)
		}
	}

	return ruleBranchs
}

func resourceTencentCloudTeoL7AccRuleSetBranchs(ruleBranches []*teo.RuleBranch) []map[string]interface{} {
	branchesList := make([]map[string]interface{}, 0, len(ruleBranches))
	if len(ruleBranches) > 0 {
		for _, branches := range ruleBranches {
			branchesMap := map[string]interface{}{}

			if branches.Condition != nil {
				branchesMap["condition"] = branches.Condition
			}

			actionsList := make([]map[string]interface{}, 0, len(branches.Actions))
			if branches.Actions != nil {
				for _, actions := range branches.Actions {
					actionsMap := map[string]interface{}{}

					if actions.Name != nil {
						actionsMap["name"] = actions.Name
					}

					cacheParametersMap := map[string]interface{}{}

					if actions.CacheParameters != nil {
						followOriginMap := map[string]interface{}{}

						if actions.CacheParameters.FollowOrigin != nil {
							if actions.CacheParameters.FollowOrigin.Switch != nil {
								followOriginMap["switch"] = actions.CacheParameters.FollowOrigin.Switch
							}

							if actions.CacheParameters.FollowOrigin.DefaultCache != nil {
								followOriginMap["default_cache"] = actions.CacheParameters.FollowOrigin.DefaultCache
							}

							if actions.CacheParameters.FollowOrigin.DefaultCacheStrategy != nil {
								followOriginMap["default_cache_strategy"] = actions.CacheParameters.FollowOrigin.DefaultCacheStrategy
							}

							if actions.CacheParameters.FollowOrigin.DefaultCacheTime != nil {
								followOriginMap["default_cache_time"] = actions.CacheParameters.FollowOrigin.DefaultCacheTime
							}

							cacheParametersMap["follow_origin"] = []interface{}{followOriginMap}
						}

						noCacheMap := map[string]interface{}{}

						if actions.CacheParameters.NoCache != nil {
							if actions.CacheParameters.NoCache.Switch != nil {
								noCacheMap["switch"] = actions.CacheParameters.NoCache.Switch
							}

							cacheParametersMap["no_cache"] = []interface{}{noCacheMap}
						}

						customTimeMap := map[string]interface{}{}

						if actions.CacheParameters.CustomTime != nil {
							if actions.CacheParameters.CustomTime.Switch != nil {
								customTimeMap["switch"] = actions.CacheParameters.CustomTime.Switch
							}

							if actions.CacheParameters.CustomTime.IgnoreCacheControl != nil {
								customTimeMap["ignore_cache_control"] = actions.CacheParameters.CustomTime.IgnoreCacheControl
							}

							if actions.CacheParameters.CustomTime.CacheTime != nil {
								customTimeMap["cache_time"] = actions.CacheParameters.CustomTime.CacheTime
							}

							cacheParametersMap["custom_time"] = []interface{}{customTimeMap}
						}

						actionsMap["cache_parameters"] = []interface{}{cacheParametersMap}
					}

					cacheKeyParametersMap := map[string]interface{}{}

					if actions.CacheKeyParameters != nil {
						if actions.CacheKeyParameters.FullURLCache != nil {
							cacheKeyParametersMap["full_url_cache"] = actions.CacheKeyParameters.FullURLCache
						}

						queryStringMap := map[string]interface{}{}

						if actions.CacheKeyParameters.QueryString != nil {
							if actions.CacheKeyParameters.QueryString.Switch != nil {
								queryStringMap["switch"] = actions.CacheKeyParameters.QueryString.Switch
							}

							if actions.CacheKeyParameters.QueryString.Action != nil {
								queryStringMap["action"] = actions.CacheKeyParameters.QueryString.Action
							}

							if actions.CacheKeyParameters.QueryString.Values != nil {
								queryStringMap["values"] = actions.CacheKeyParameters.QueryString.Values
							}

							cacheKeyParametersMap["query_string"] = []interface{}{queryStringMap}
						}

						if actions.CacheKeyParameters.IgnoreCase != nil {
							cacheKeyParametersMap["ignore_case"] = actions.CacheKeyParameters.IgnoreCase
						}

						headerMap := map[string]interface{}{}

						if actions.CacheKeyParameters.Header != nil {
							if actions.CacheKeyParameters.Header.Switch != nil {
								headerMap["switch"] = actions.CacheKeyParameters.Header.Switch
							}

							if actions.CacheKeyParameters.Header.Values != nil {
								headerMap["values"] = actions.CacheKeyParameters.Header.Values
							}

							cacheKeyParametersMap["header"] = []interface{}{headerMap}
						}

						if actions.CacheKeyParameters.Scheme != nil {
							cacheKeyParametersMap["scheme"] = actions.CacheKeyParameters.Scheme
						}

						cookieMap := map[string]interface{}{}

						if actions.CacheKeyParameters.Cookie != nil {
							if actions.CacheKeyParameters.Cookie.Switch != nil {
								cookieMap["switch"] = actions.CacheKeyParameters.Cookie.Switch
							}

							if actions.CacheKeyParameters.Cookie.Action != nil {
								cookieMap["action"] = actions.CacheKeyParameters.Cookie.Action
							}

							if actions.CacheKeyParameters.Cookie.Values != nil {
								cookieMap["values"] = actions.CacheKeyParameters.Cookie.Values
							}

							cacheKeyParametersMap["cookie"] = []interface{}{cookieMap}
						}

						actionsMap["cache_key_parameters"] = []interface{}{cacheKeyParametersMap}
					}

					cachePrefreshParametersMap := map[string]interface{}{}

					if actions.CachePrefreshParameters != nil {
						if actions.CachePrefreshParameters.Switch != nil {
							cachePrefreshParametersMap["switch"] = actions.CachePrefreshParameters.Switch
						}

						if actions.CachePrefreshParameters.CacheTimePercent != nil {
							cachePrefreshParametersMap["cache_time_percent"] = actions.CachePrefreshParameters.CacheTimePercent
						}

						actionsMap["cache_prefresh_parameters"] = []interface{}{cachePrefreshParametersMap}
					}

					accessURLRedirectParametersMap := map[string]interface{}{}

					if actions.AccessURLRedirectParameters != nil {
						if actions.AccessURLRedirectParameters.StatusCode != nil {
							accessURLRedirectParametersMap["status_code"] = actions.AccessURLRedirectParameters.StatusCode
						}

						if actions.AccessURLRedirectParameters.Protocol != nil {
							accessURLRedirectParametersMap["protocol"] = actions.AccessURLRedirectParameters.Protocol
						}

						hostNameMap := map[string]interface{}{}

						if actions.AccessURLRedirectParameters.HostName != nil {
							if actions.AccessURLRedirectParameters.HostName.Action != nil {
								hostNameMap["action"] = actions.AccessURLRedirectParameters.HostName.Action
							}

							if actions.AccessURLRedirectParameters.HostName.Value != nil {
								hostNameMap["value"] = actions.AccessURLRedirectParameters.HostName.Value
							}

							accessURLRedirectParametersMap["host_name"] = []interface{}{hostNameMap}
						}

						uRLPathMap := map[string]interface{}{}

						if actions.AccessURLRedirectParameters.URLPath != nil {
							if actions.AccessURLRedirectParameters.URLPath.Action != nil {
								uRLPathMap["action"] = actions.AccessURLRedirectParameters.URLPath.Action
							}

							if actions.AccessURLRedirectParameters.URLPath.Regex != nil {
								uRLPathMap["regex"] = actions.AccessURLRedirectParameters.URLPath.Regex
							}

							if actions.AccessURLRedirectParameters.URLPath.Value != nil {
								uRLPathMap["value"] = actions.AccessURLRedirectParameters.URLPath.Value
							}

							accessURLRedirectParametersMap["url_path"] = []interface{}{uRLPathMap}
						}

						queryStringMap := map[string]interface{}{}

						if actions.AccessURLRedirectParameters.QueryString != nil {
							if actions.AccessURLRedirectParameters.QueryString.Action != nil {
								queryStringMap["action"] = actions.AccessURLRedirectParameters.QueryString.Action
							}

							accessURLRedirectParametersMap["query_string"] = []interface{}{queryStringMap}
						}

						actionsMap["access_url_redirect_parameters"] = []interface{}{accessURLRedirectParametersMap}
					}

					upstreamURLRewriteParametersMap := map[string]interface{}{}

					if actions.UpstreamURLRewriteParameters != nil {
						if actions.UpstreamURLRewriteParameters.Type != nil {
							upstreamURLRewriteParametersMap["type"] = actions.UpstreamURLRewriteParameters.Type
						}

						if actions.UpstreamURLRewriteParameters.Action != nil {
							upstreamURLRewriteParametersMap["action"] = actions.UpstreamURLRewriteParameters.Action
						}

						if actions.UpstreamURLRewriteParameters.Value != nil {
							upstreamURLRewriteParametersMap["value"] = actions.UpstreamURLRewriteParameters.Value
						}

						if actions.UpstreamURLRewriteParameters.Regex != nil {
							upstreamURLRewriteParametersMap["regex"] = actions.UpstreamURLRewriteParameters.Regex

						}

						actionsMap["upstream_url_rewrite_parameters"] = []interface{}{upstreamURLRewriteParametersMap}
					}

					qUICParametersMap := map[string]interface{}{}

					if actions.QUICParameters != nil {
						if actions.QUICParameters.Switch != nil {
							qUICParametersMap["switch"] = actions.QUICParameters.Switch
						}

						actionsMap["quic_parameters"] = []interface{}{qUICParametersMap}
					}

					webSocketParametersMap := map[string]interface{}{}

					if actions.WebSocketParameters != nil {
						if actions.WebSocketParameters.Switch != nil {
							webSocketParametersMap["switch"] = actions.WebSocketParameters.Switch
						}

						if actions.WebSocketParameters.Timeout != nil {
							webSocketParametersMap["timeout"] = actions.WebSocketParameters.Timeout
						}

						actionsMap["web_socket_parameters"] = []interface{}{webSocketParametersMap}
					}

					authenticationParametersMap := map[string]interface{}{}

					if actions.AuthenticationParameters != nil {
						if actions.AuthenticationParameters.AuthType != nil {
							authenticationParametersMap["auth_type"] = actions.AuthenticationParameters.AuthType
						}

						if actions.AuthenticationParameters.SecretKey != nil {
							authenticationParametersMap["secret_key"] = actions.AuthenticationParameters.SecretKey
						}

						if actions.AuthenticationParameters.Timeout != nil {
							authenticationParametersMap["timeout"] = actions.AuthenticationParameters.Timeout
						}

						if actions.AuthenticationParameters.BackupSecretKey != nil {
							authenticationParametersMap["backup_secret_key"] = actions.AuthenticationParameters.BackupSecretKey
						}

						if actions.AuthenticationParameters.AuthParam != nil {
							authenticationParametersMap["auth_param"] = actions.AuthenticationParameters.AuthParam
						}

						if actions.AuthenticationParameters.TimeParam != nil {
							authenticationParametersMap["time_param"] = actions.AuthenticationParameters.TimeParam
						}

						if actions.AuthenticationParameters.TimeFormat != nil {
							authenticationParametersMap["time_format"] = actions.AuthenticationParameters.TimeFormat
						}

						actionsMap["authentication_parameters"] = []interface{}{authenticationParametersMap}
					}

					maxAgeParametersMap := map[string]interface{}{}

					if actions.MaxAgeParameters != nil {
						if actions.MaxAgeParameters.FollowOrigin != nil {
							maxAgeParametersMap["follow_origin"] = actions.MaxAgeParameters.FollowOrigin
						}

						if actions.MaxAgeParameters.CacheTime != nil {
							maxAgeParametersMap["cache_time"] = actions.MaxAgeParameters.CacheTime
						}

						actionsMap["max_age_parameters"] = []interface{}{maxAgeParametersMap}
					}

					statusCodeCacheParametersMap := map[string]interface{}{}

					if actions.StatusCodeCacheParameters != nil {
						statusCodeCacheParamsList := make([]map[string]interface{}, 0, len(actions.StatusCodeCacheParameters.StatusCodeCacheParams))
						if actions.StatusCodeCacheParameters.StatusCodeCacheParams != nil {
							for _, statusCodeCacheParams := range actions.StatusCodeCacheParameters.StatusCodeCacheParams {
								statusCodeCacheParamsMap := map[string]interface{}{}

								if statusCodeCacheParams.StatusCode != nil {
									statusCodeCacheParamsMap["status_code"] = statusCodeCacheParams.StatusCode
								}

								if statusCodeCacheParams.CacheTime != nil {
									statusCodeCacheParamsMap["cache_time"] = statusCodeCacheParams.CacheTime
								}

								statusCodeCacheParamsList = append(statusCodeCacheParamsList, statusCodeCacheParamsMap)
							}

							statusCodeCacheParametersMap["status_code_cache_params"] = statusCodeCacheParamsList
						}
						actionsMap["status_code_cache_parameters"] = []interface{}{statusCodeCacheParametersMap}
					}

					offlineCacheParametersMap := map[string]interface{}{}

					if actions.OfflineCacheParameters != nil {
						if actions.OfflineCacheParameters.Switch != nil {
							offlineCacheParametersMap["switch"] = actions.OfflineCacheParameters.Switch
						}

						actionsMap["offline_cache_parameters"] = []interface{}{offlineCacheParametersMap}
					}

					smartRoutingParametersMap := map[string]interface{}{}

					if actions.SmartRoutingParameters != nil {
						if actions.SmartRoutingParameters.Switch != nil {
							smartRoutingParametersMap["switch"] = actions.SmartRoutingParameters.Switch
						}

						actionsMap["smart_routing_parameters"] = []interface{}{smartRoutingParametersMap}
					}

					rangeOriginPullParametersMap := map[string]interface{}{}

					if actions.RangeOriginPullParameters != nil {
						if actions.RangeOriginPullParameters.Switch != nil {
							rangeOriginPullParametersMap["switch"] = actions.RangeOriginPullParameters.Switch
						}

						actionsMap["range_origin_pull_parameters"] = []interface{}{rangeOriginPullParametersMap}
					}

					upstreamHTTP2ParametersMap := map[string]interface{}{}

					if actions.UpstreamHTTP2Parameters != nil {
						if actions.UpstreamHTTP2Parameters.Switch != nil {
							upstreamHTTP2ParametersMap["switch"] = actions.UpstreamHTTP2Parameters.Switch
						}

						actionsMap["upstream_http2_parameters"] = []interface{}{upstreamHTTP2ParametersMap}
					}

					hostHeaderParametersMap := map[string]interface{}{}

					if actions.HostHeaderParameters != nil {
						if actions.HostHeaderParameters.Action != nil {
							hostHeaderParametersMap["action"] = actions.HostHeaderParameters.Action
						}

						if actions.HostHeaderParameters.ServerName != nil {
							hostHeaderParametersMap["server_name"] = actions.HostHeaderParameters.ServerName
						}

						actionsMap["host_header_parameters"] = []interface{}{hostHeaderParametersMap}
					}

					forceRedirectHTTPSParametersMap := map[string]interface{}{}

					if actions.ForceRedirectHTTPSParameters != nil {
						if actions.ForceRedirectHTTPSParameters.Switch != nil {
							forceRedirectHTTPSParametersMap["switch"] = actions.ForceRedirectHTTPSParameters.Switch
						}

						if actions.ForceRedirectHTTPSParameters.RedirectStatusCode != nil {
							forceRedirectHTTPSParametersMap["redirect_status_code"] = actions.ForceRedirectHTTPSParameters.RedirectStatusCode
						}

						actionsMap["force_redirect_https_parameters"] = []interface{}{forceRedirectHTTPSParametersMap}
					}

					originPullProtocolParametersMap := map[string]interface{}{}

					if actions.OriginPullProtocolParameters != nil {
						if actions.OriginPullProtocolParameters.Protocol != nil {
							originPullProtocolParametersMap["protocol"] = actions.OriginPullProtocolParameters.Protocol
						}

						actionsMap["origin_pull_protocol_parameters"] = []interface{}{originPullProtocolParametersMap}
					}

					compressionParametersMap := map[string]interface{}{}

					if actions.CompressionParameters != nil {
						if actions.CompressionParameters.Switch != nil {
							compressionParametersMap["switch"] = actions.CompressionParameters.Switch
						}

						if actions.CompressionParameters.Algorithms != nil {
							compressionParametersMap["algorithms"] = actions.CompressionParameters.Algorithms
						}

						actionsMap["compression_parameters"] = []interface{}{compressionParametersMap}
					}

					hSTSParametersMap := map[string]interface{}{}

					if actions.HSTSParameters != nil {
						if actions.HSTSParameters.Switch != nil {
							hSTSParametersMap["switch"] = actions.HSTSParameters.Switch
						}

						if actions.HSTSParameters.Timeout != nil {
							hSTSParametersMap["timeout"] = actions.HSTSParameters.Timeout
						}

						if actions.HSTSParameters.IncludeSubDomains != nil {
							hSTSParametersMap["include_sub_domains"] = actions.HSTSParameters.IncludeSubDomains
						}

						if actions.HSTSParameters.Preload != nil {
							hSTSParametersMap["preload"] = actions.HSTSParameters.Preload
						}

						actionsMap["hsts_parameters"] = []interface{}{hSTSParametersMap}
					}

					clientIPHeaderParametersMap := map[string]interface{}{}

					if actions.ClientIPHeaderParameters != nil {
						if actions.ClientIPHeaderParameters.Switch != nil {
							clientIPHeaderParametersMap["switch"] = actions.ClientIPHeaderParameters.Switch
						}

						if actions.ClientIPHeaderParameters.HeaderName != nil {
							clientIPHeaderParametersMap["header_name"] = actions.ClientIPHeaderParameters.HeaderName
						}

						actionsMap["client_ip_header_parameters"] = []interface{}{clientIPHeaderParametersMap}
					}

					oCSPStaplingParametersMap := map[string]interface{}{}

					if actions.OCSPStaplingParameters != nil {
						if actions.OCSPStaplingParameters.Switch != nil {
							oCSPStaplingParametersMap["switch"] = actions.OCSPStaplingParameters.Switch
						}

						actionsMap["ocsp_stapling_parameters"] = []interface{}{oCSPStaplingParametersMap}
					}

					hTTP2ParametersMap := map[string]interface{}{}

					if actions.HTTP2Parameters != nil {
						if actions.HTTP2Parameters.Switch != nil {
							hTTP2ParametersMap["switch"] = actions.HTTP2Parameters.Switch
						}

						actionsMap["http2_parameters"] = []interface{}{hTTP2ParametersMap}
					}

					postMaxSizeParametersMap := map[string]interface{}{}

					if actions.PostMaxSizeParameters != nil {
						if actions.PostMaxSizeParameters.Switch != nil {
							postMaxSizeParametersMap["switch"] = actions.PostMaxSizeParameters.Switch
						}

						if actions.PostMaxSizeParameters.MaxSize != nil {
							postMaxSizeParametersMap["max_size"] = actions.PostMaxSizeParameters.MaxSize
						}

						actionsMap["post_max_size_parameters"] = []interface{}{postMaxSizeParametersMap}
					}

					clientIPCountryParametersMap := map[string]interface{}{}

					if actions.ClientIPCountryParameters != nil {
						if actions.ClientIPCountryParameters.Switch != nil {
							clientIPCountryParametersMap["switch"] = actions.ClientIPCountryParameters.Switch
						}

						if actions.ClientIPCountryParameters.HeaderName != nil {
							clientIPCountryParametersMap["header_name"] = actions.ClientIPCountryParameters.HeaderName
						}

						actionsMap["client_ip_country_parameters"] = []interface{}{clientIPCountryParametersMap}
					}

					upstreamFollowRedirectParametersMap := map[string]interface{}{}

					if actions.UpstreamFollowRedirectParameters != nil {
						if actions.UpstreamFollowRedirectParameters.Switch != nil {
							upstreamFollowRedirectParametersMap["switch"] = actions.UpstreamFollowRedirectParameters.Switch
						}

						if actions.UpstreamFollowRedirectParameters.MaxTimes != nil {
							upstreamFollowRedirectParametersMap["max_times"] = actions.UpstreamFollowRedirectParameters.MaxTimes
						}

						actionsMap["upstream_follow_redirect_parameters"] = []interface{}{upstreamFollowRedirectParametersMap}
					}

					upstreamRequestParametersMap := map[string]interface{}{}

					if actions.UpstreamRequestParameters != nil {
						queryStringMap := map[string]interface{}{}

						if actions.UpstreamRequestParameters.QueryString != nil {
							if actions.UpstreamRequestParameters.QueryString.Switch != nil {
								queryStringMap["switch"] = actions.UpstreamRequestParameters.QueryString.Switch
							}

							if actions.UpstreamRequestParameters.QueryString.Action != nil {
								queryStringMap["action"] = actions.UpstreamRequestParameters.QueryString.Action
							}

							if actions.UpstreamRequestParameters.QueryString.Values != nil {
								queryStringMap["values"] = actions.UpstreamRequestParameters.QueryString.Values
							}

							upstreamRequestParametersMap["query_string"] = []interface{}{queryStringMap}
						}

						cookieMap := map[string]interface{}{}

						if actions.UpstreamRequestParameters.Cookie != nil {
							if actions.UpstreamRequestParameters.Cookie.Switch != nil {
								cookieMap["switch"] = actions.UpstreamRequestParameters.Cookie.Switch
							}

							if actions.UpstreamRequestParameters.Cookie.Action != nil {
								cookieMap["action"] = actions.UpstreamRequestParameters.Cookie.Action
							}

							if actions.UpstreamRequestParameters.Cookie.Values != nil {
								cookieMap["values"] = actions.UpstreamRequestParameters.Cookie.Values
							}

							upstreamRequestParametersMap["cookie"] = []interface{}{cookieMap}
						}

						actionsMap["upstream_request_parameters"] = []interface{}{upstreamRequestParametersMap}
					}

					tLSConfigParametersMap := map[string]interface{}{}

					if actions.TLSConfigParameters != nil {
						if actions.TLSConfigParameters.Version != nil {
							tLSConfigParametersMap["version"] = actions.TLSConfigParameters.Version
						}

						if actions.TLSConfigParameters.CipherSuite != nil {
							tLSConfigParametersMap["cipher_suite"] = actions.TLSConfigParameters.CipherSuite
						}

						actionsMap["tls_config_parameters"] = []interface{}{tLSConfigParametersMap}
					}

					modifyOriginParametersMap := map[string]interface{}{}

					if actions.ModifyOriginParameters != nil {
						if actions.ModifyOriginParameters.OriginType != nil {
							modifyOriginParametersMap["origin_type"] = actions.ModifyOriginParameters.OriginType
						}

						if actions.ModifyOriginParameters.Origin != nil {
							modifyOriginParametersMap["origin"] = actions.ModifyOriginParameters.Origin
						}

						if actions.ModifyOriginParameters.OriginProtocol != nil {
							modifyOriginParametersMap["origin_protocol"] = actions.ModifyOriginParameters.OriginProtocol
						}

						if actions.ModifyOriginParameters.HTTPOriginPort != nil {
							modifyOriginParametersMap["http_origin_port"] = actions.ModifyOriginParameters.HTTPOriginPort
						}

						if actions.ModifyOriginParameters.HTTPSOriginPort != nil {
							modifyOriginParametersMap["https_origin_port"] = actions.ModifyOriginParameters.HTTPSOriginPort
						}

						if actions.ModifyOriginParameters.PrivateAccess != nil {
							modifyOriginParametersMap["private_access"] = actions.ModifyOriginParameters.PrivateAccess
						}

						privateParametersMap := map[string]interface{}{}

						if actions.ModifyOriginParameters.PrivateParameters != nil {
							if actions.ModifyOriginParameters.PrivateParameters.AccessKeyId != nil {
								privateParametersMap["access_key_id"] = actions.ModifyOriginParameters.PrivateParameters.AccessKeyId
							}

							if actions.ModifyOriginParameters.PrivateParameters.SecretAccessKey != nil {
								privateParametersMap["secret_access_key"] = actions.ModifyOriginParameters.PrivateParameters.SecretAccessKey
							}

							if actions.ModifyOriginParameters.PrivateParameters.SignatureVersion != nil {
								privateParametersMap["signature_version"] = actions.ModifyOriginParameters.PrivateParameters.SignatureVersion
							}

							if actions.ModifyOriginParameters.PrivateParameters.Region != nil {
								privateParametersMap["region"] = actions.ModifyOriginParameters.PrivateParameters.Region
							}

							modifyOriginParametersMap["private_parameters"] = []interface{}{privateParametersMap}
						}

						actionsMap["modify_origin_parameters"] = []interface{}{modifyOriginParametersMap}
					}

					hTTPUpstreamTimeoutParametersMap := map[string]interface{}{}

					if actions.HTTPUpstreamTimeoutParameters != nil {
						if actions.HTTPUpstreamTimeoutParameters.ResponseTimeout != nil {
							hTTPUpstreamTimeoutParametersMap["response_timeout"] = actions.HTTPUpstreamTimeoutParameters.ResponseTimeout
						}

						actionsMap["http_upstream_timeout_parameters"] = []interface{}{hTTPUpstreamTimeoutParametersMap}
					}

					httpResponseParametersMap := map[string]interface{}{}

					if actions.HttpResponseParameters != nil {
						if actions.HttpResponseParameters.StatusCode != nil {
							httpResponseParametersMap["status_code"] = actions.HttpResponseParameters.StatusCode
						}

						if actions.HttpResponseParameters.ResponsePage != nil {
							httpResponseParametersMap["response_page"] = actions.HttpResponseParameters.ResponsePage
						}

						actionsMap["http_response_parameters"] = []interface{}{httpResponseParametersMap}
					}

					errorPageParametersMap := map[string]interface{}{}

					if actions.ErrorPageParameters != nil {
						errorPageParamsList := make([]map[string]interface{}, 0, len(actions.ErrorPageParameters.ErrorPageParams))
						if actions.ErrorPageParameters.ErrorPageParams != nil {
							for _, errorPageParams := range actions.ErrorPageParameters.ErrorPageParams {
								errorPageParamsMap := map[string]interface{}{}

								if errorPageParams.StatusCode != nil {
									errorPageParamsMap["status_code"] = errorPageParams.StatusCode
								}

								if errorPageParams.RedirectURL != nil {
									errorPageParamsMap["redirect_url"] = errorPageParams.RedirectURL
								}

								errorPageParamsList = append(errorPageParamsList, errorPageParamsMap)
							}

							errorPageParametersMap["error_page_params"] = errorPageParamsList
						}
						actionsMap["error_page_parameters"] = []interface{}{errorPageParametersMap}
					}

					modifyResponseHeaderParametersMap := map[string]interface{}{}

					if actions.ModifyResponseHeaderParameters != nil {
						headerActionsList := make([]map[string]interface{}, 0, len(actions.ModifyResponseHeaderParameters.HeaderActions))
						if actions.ModifyResponseHeaderParameters.HeaderActions != nil {
							for _, headerActions := range actions.ModifyResponseHeaderParameters.HeaderActions {
								headerActionsMap := map[string]interface{}{}

								if headerActions.Action != nil {
									headerActionsMap["action"] = headerActions.Action
								}

								if headerActions.Name != nil {
									headerActionsMap["name"] = headerActions.Name
								}

								if headerActions.Value != nil {
									headerActionsMap["value"] = headerActions.Value
								}

								headerActionsList = append(headerActionsList, headerActionsMap)
							}

							modifyResponseHeaderParametersMap["header_actions"] = headerActionsList
						}
						actionsMap["modify_response_header_parameters"] = []interface{}{modifyResponseHeaderParametersMap}
					}

					modifyRequestHeaderParametersMap := map[string]interface{}{}

					if actions.ModifyRequestHeaderParameters != nil {
						headerActionsList := make([]map[string]interface{}, 0, len(actions.ModifyRequestHeaderParameters.HeaderActions))
						if actions.ModifyRequestHeaderParameters.HeaderActions != nil {
							for _, headerActions := range actions.ModifyRequestHeaderParameters.HeaderActions {
								headerActionsMap := map[string]interface{}{}

								if headerActions.Action != nil {
									headerActionsMap["action"] = headerActions.Action
								}

								if headerActions.Name != nil {
									headerActionsMap["name"] = headerActions.Name
								}

								if headerActions.Value != nil {
									headerActionsMap["value"] = headerActions.Value
								}

								headerActionsList = append(headerActionsList, headerActionsMap)
							}

							modifyRequestHeaderParametersMap["header_actions"] = headerActionsList
						}
						actionsMap["modify_request_header_parameters"] = []interface{}{modifyRequestHeaderParametersMap}
					}

					responseSpeedLimitParametersMap := map[string]interface{}{}

					if actions.ResponseSpeedLimitParameters != nil {
						if actions.ResponseSpeedLimitParameters.Mode != nil {
							responseSpeedLimitParametersMap["mode"] = actions.ResponseSpeedLimitParameters.Mode
						}

						if actions.ResponseSpeedLimitParameters.MaxSpeed != nil {
							responseSpeedLimitParametersMap["max_speed"] = actions.ResponseSpeedLimitParameters.MaxSpeed
						}

						if actions.ResponseSpeedLimitParameters.StartAt != nil {
							responseSpeedLimitParametersMap["start_at"] = actions.ResponseSpeedLimitParameters.StartAt
						}

						actionsMap["response_speed_limit_parameters"] = []interface{}{responseSpeedLimitParametersMap}
					}

					setContentIdentifierParametersMap := map[string]interface{}{}

					if actions.SetContentIdentifierParameters != nil {
						if actions.SetContentIdentifierParameters.ContentIdentifier != nil {
							setContentIdentifierParametersMap["content_identifier"] = actions.SetContentIdentifierParameters.ContentIdentifier
						}

						actionsMap["set_content_identifier_parameters"] = []interface{}{setContentIdentifierParametersMap}
					}

					contentCompressionParametersMap := map[string]interface{}{}
					if actions.ContentCompressionParameters != nil {
						if actions.ContentCompressionParameters.Switch != nil {
							contentCompressionParametersMap["switch"] = actions.ContentCompressionParameters.Switch
						}

						actionsMap["content_compression_parameters"] = []interface{}{contentCompressionParametersMap}
					}

					varyParametersMap := map[string]interface{}{}
					if actions.VaryParameters != nil {
						if actions.VaryParameters.Switch != nil {
							varyParametersMap["switch"] = actions.VaryParameters.Switch
						}

						actionsMap["vary_parameters"] = []interface{}{varyParametersMap}
					}

					originAuthenticationParametersMap := map[string]interface{}{}
					if actions.OriginAuthenticationParameters != nil {
						requestPropertiesList := make([]map[string]interface{}, 0, len(actions.OriginAuthenticationParameters.RequestProperties))
						for _, prop := range actions.OriginAuthenticationParameters.RequestProperties {
							propMap := map[string]interface{}{}
							if prop.Type != nil {
								propMap["type"] = prop.Type
							}
							if prop.Name != nil {
								propMap["name"] = prop.Name
							}
							if prop.Value != nil {
								propMap["value"] = prop.Value
							}
							requestPropertiesList = append(requestPropertiesList, propMap)
						}
						originAuthenticationParametersMap["request_properties"] = requestPropertiesList

						actionsMap["origin_authentication_parameters"] = []interface{}{originAuthenticationParametersMap}
					}

					actionsList = append(actionsList, actionsMap)
				}

				branchesMap["actions"] = actionsList
			}

			subRulesList := make([]map[string]interface{}, 0, len(branches.SubRules))
			if branches.SubRules != nil {
				for _, subRules := range branches.SubRules {
					subRulesMap := map[string]interface{}{}

					if subRules.Branches != nil {
						subRulesMap["branches"] = resourceTencentCloudTeoL7AccRuleSetBranchs(subRules.Branches)
					}

					if subRules.Description != nil {
						subRulesMap["description"] = subRules.Description
					}

					subRulesList = append(subRulesList, subRulesMap)
				}

				branchesMap["sub_rules"] = subRulesList
			}

			branchesList = append(branchesList, branchesMap)
		}
	}
	return branchesList
}

func resourceTencentCloudTeoL7AccRuleContent(rules []*teo.RuleEngineItem) (string, error) {
	type Content struct {
		FormatVersion string                `json:"FormatVersion,omitempty"`
		Rules         []*teo.RuleEngineItem `json:"Rules,omitempty"`
	}
	content := Content{
		FormatVersion: "1.0",
		Rules:         rules,
	}
	contentBytes, err := json.Marshal(content)
	if err != nil {
		return "", err
	}
	return string(contentBytes), nil
}
