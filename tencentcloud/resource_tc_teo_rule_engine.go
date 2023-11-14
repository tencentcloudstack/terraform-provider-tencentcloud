/*
Provides a resource to create a teo rule_engine

Example Usage

```hcl
resource "tencentcloud_teo_rule_engine" "rule_engine" {
  zone_id = ""
  rule_name = ""
  status = ""
  rules {
		actions {
			normal_action {
				action = ""
				parameters {
					name = ""
					values =
				}
			}
			rewrite_action {
				action = ""
				parameters {
					action = ""
					name = ""
					values =
				}
			}
			code_action {
				action = ""
				parameters {
					status_code =
					name = ""
					values =
				}
			}
		}
		conditions {
			conditions {
				operator = ""
				target = ""
				values =
				ignore_case =
				name = ""
				ignore_name_case =
			}
		}
		sub_rules {
			rules {
				conditions {
					conditions {
						operator = ""
						target = ""
						values =
						ignore_case =
						name = ""
						ignore_name_case =
					}
				}
				actions {
					normal_action {
						action = ""
						parameters {
							name = ""
							values =
						}
					}
					rewrite_action {
						action = ""
						parameters {
							action = ""
							name = ""
							values =
						}
					}
					code_action {
						action = ""
						parameters {
							status_code =
							name = ""
							values =
						}
					}
				}
			}
			tags =
		}

  }
}
```

Import

teo rule_engine can be imported using the id, e.g.

```
terraform import tencentcloud_teo_rule_engine.rule_engine rule_engine_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	teo "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
	"strings"
)

func resourceTencentCloudTeoRule_engine() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTeoRule_engineCreate,
		Read:   resourceTencentCloudTeoRule_engineRead,
		Update: resourceTencentCloudTeoRule_engineUpdate,
		Delete: resourceTencentCloudTeoRule_engineDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"zone_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "ID of the site.",
			},

			"rule_name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "The rule name (1 to 255 characters).",
			},

			"status": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Rule status. Values:&amp;amp;lt;li&amp;amp;gt;`enable`: Enabled&amp;amp;lt;/li&amp;amp;gt;&amp;amp;lt;li&amp;amp;gt;`disable`: Disabled&amp;amp;lt;/li&amp;amp;gt;.",
			},

			"rules": {
				Required:    true,
				Type:        schema.TypeList,
				Description: "The rule content.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"actions": {
							Type:        schema.TypeList,
							Required:    true,
							Description: "Feature to be executed.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"normal_action": {
										Type:        schema.TypeList,
										MaxItems:    1,
										Optional:    true,
										Description: "Common operation. Values:&amp;lt;li&amp;gt;`AccessUrlRedirect`: Access URL rewrite&amp;lt;/li&amp;gt;&amp;lt;li&amp;gt;`UpstreamUrlRedirect`: Origin-pull URL rewrite&amp;lt;/li&amp;gt;&amp;lt;li&amp;gt;`QUIC`: QUIC&amp;lt;/li&amp;gt;&amp;lt;li&amp;gt;`WebSocket`: WebSocket&amp;lt;/li&amp;gt;&amp;lt;li&amp;gt;`VideoSeek`: Video dragging&amp;lt;/li&amp;gt;&amp;lt;li&amp;gt;`Authentication`: Token authentication&amp;lt;/li&amp;gt;&amp;lt;li&amp;gt;`CacheKey`: Custom cache key&amp;lt;/li&amp;gt;&amp;lt;li&amp;gt;`Cache`: Node cache TTL&amp;lt;/li&amp;gt;&amp;lt;li&amp;gt;`MaxAge`: Browser cache TTL&amp;lt;/li&amp;gt;&amp;lt;li&amp;gt;`OfflineCache`: Offline cache&amp;lt;/li&amp;gt;&amp;lt;li&amp;gt;`SmartRouting`: Smart acceleration&amp;lt;/li&amp;gt;&amp;lt;li&amp;gt;`RangeOriginPull`: Range GETs&amp;lt;/li&amp;gt;&amp;lt;li&amp;gt;`UpstreamHttp2`: HTTP/2 forwarding&amp;lt;/li&amp;gt;&amp;lt;li&amp;gt;`HostHeader`: Host header rewrite&amp;lt;/li&amp;gt;&amp;lt;li&amp;gt;`ForceRedirect`: Force HTTPS&amp;lt;/li&amp;gt;&amp;lt;li&amp;gt;`OriginPullProtocol`: Origin-pull HTTPS&amp;lt;/li&amp;gt;&amp;lt;li&amp;gt;`CachePrefresh`: Cache prefresh&amp;lt;/li&amp;gt;&amp;lt;li&amp;gt;`Compression`: Smart compression&amp;lt;/li&amp;gt;&amp;lt;li&amp;gt;`Hsts`&amp;lt;/li&amp;gt;&amp;lt;li&amp;gt;`ClientIpHeader`&amp;lt;/li&amp;gt;&amp;lt;li&amp;gt;`SslTlsSecureConf`&amp;lt;/li&amp;gt;&amp;lt;li&amp;gt;`OcspStapling`&amp;lt;/li&amp;gt;&amp;lt;li&amp;gt;`Http2`: HTTP/2 access&amp;lt;/li&amp;gt;&amp;lt;li&amp;gt;`UpstreamFollowRedirect`: Follow origin redirect&amp;lt;/li&amp;gt;&amp;lt;li&amp;gt;`Origin`: Origin&amp;lt;/li&amp;gt;Note: This field may return `null`, indicating that no valid value can be obtained.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"action": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "Feature name. You can call the [DescribeRulesSetting](https://tcloud4api.woa.com/document/product/1657/79433?!preview&amp;amp;!document=1) API to view the requirements for entering the feature name.",
												},
												"parameters": {
													Type:        schema.TypeList,
													Required:    true,
													Description: "Parameter.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"name": {
																Type:        schema.TypeString,
																Required:    true,
																Description: "Parameter name. You can call the [DescribeRulesSetting](https://tcloud4api.woa.com/document/product/1657/79433?!preview&amp;amp;!document=1) API to view the requirements for entering the parameter name.",
															},
															"values": {
																Type: schema.TypeSet,
																Elem: &schema.Schema{
																	Type: schema.TypeString,
																},
																Required:    true,
																Description: "The parameter value.",
															},
														},
													},
												},
											},
										},
									},
									"rewrite_action": {
										Type:        schema.TypeList,
										MaxItems:    1,
										Optional:    true,
										Description: "Feature operation with a request/response header. Features of this type include:&amp;lt;li&amp;gt;`RequestHeader`: HTTP request header modification.&amp;lt;/li&amp;gt;&amp;lt;li&amp;gt;`ResponseHeader`: HTTP response header modification.&amp;lt;/li&amp;gt;Note: This field may return null, indicating that no valid values can be obtained.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"action": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "Feature name. You can call the [DescribeRulesSetting](https://tcloud4api.woa.com/document/product/1657/79433?!preview&amp;amp;!document=1) API to view the requirements for entering the feature name.",
												},
												"parameters": {
													Type:        schema.TypeList,
													Required:    true,
													Description: "Parameter.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"action": {
																Type:        schema.TypeString,
																Required:    true,
																Description: "Feature parameter name. You can call the [DescribeRulesSetting](https://tcloud4api.woa.com/document/product/1657/79433?!preview&amp;amp;!document=1) API to view the requirements for entering the parameter name, which has three values:&amp;lt;li&amp;gt;add: Add the HTTP header.&amp;lt;/li&amp;gt;&amp;lt;li&amp;gt;set: Rewrite the HTTP header.&amp;lt;/li&amp;gt;&amp;lt;li&amp;gt;del: Delete the HTTP header.&amp;lt;/li&amp;gt;.",
															},
															"name": {
																Type:        schema.TypeString,
																Required:    true,
																Description: "Parameter name.",
															},
															"values": {
																Type: schema.TypeSet,
																Elem: &schema.Schema{
																	Type: schema.TypeString,
																},
																Required:    true,
																Description: "Parameter value.",
															},
														},
													},
												},
											},
										},
									},
									"code_action": {
										Type:        schema.TypeList,
										MaxItems:    1,
										Optional:    true,
										Description: "Feature operation with a status code. Features of this type include:&amp;lt;li&amp;gt;`ErrorPage`: Custom error page.&amp;lt;/li&amp;gt;&amp;lt;li&amp;gt;`StatusCodeCache`: Status code cache TTL.&amp;lt;/li&amp;gt;Note: This field may return null, indicating that no valid values can be obtained.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"action": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "Feature name. You can call the [DescribeRulesSetting](https://tcloud4api.woa.com/document/product/1657/79433?!preview&amp;amp;!document=1) API to view the requirements for entering the feature name.",
												},
												"parameters": {
													Type:        schema.TypeList,
													Required:    true,
													Description: "Operation parameter.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"status_code": {
																Type:        schema.TypeInt,
																Required:    true,
																Description: "The status code.",
															},
															"name": {
																Type:        schema.TypeString,
																Required:    true,
																Description: "The parameter name. You can call the [DescribeRulesSetting](https://tcloud4api.woa.com/document/product/1657/79433?!preview&amp;amp;!document=1) API to view the requirements for entering the parameter name.",
															},
															"values": {
																Type: schema.TypeSet,
																Elem: &schema.Schema{
																	Type: schema.TypeString,
																},
																Required:    true,
																Description: "The parameter value.",
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
						"conditions": {
							Type:        schema.TypeList,
							Required:    true,
							Description: "Feature execution conditions.Note: If any condition in the array is met, the feature will run.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"conditions": {
										Type:        schema.TypeList,
										Required:    true,
										Description: "Rule engine condition. This condition will be considered met if all items in the array are met.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"operator": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "Operator. Valid values:&amp;lt;li&amp;gt;`equals`: Equals&amp;lt;/li&amp;gt;&amp;lt;li&amp;gt;`notEquals`: Does not equal&amp;lt;/li&amp;gt;&amp;lt;li&amp;gt;`exist`: Exists&amp;lt;/li&amp;gt;&amp;lt;li&amp;gt;`notexist`: Does not exist&amp;lt;/li&amp;gt;.",
												},
												"target": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "The match type. Values:&amp;lt;li&amp;gt;`filename`: File name&amp;lt;/li&amp;gt;&amp;lt;li&amp;gt;`extension`: File extension&amp;lt;/li&amp;gt;&amp;lt;li&amp;gt;`host`: Host&amp;lt;/li&amp;gt;&amp;lt;li&amp;gt;`full_url`: Full URL, which indicates the complete URL path under the current site and must contain the HTTP protocol, host, and path.&amp;lt;/li&amp;gt;&amp;lt;li&amp;gt;`url`: Partial URL under the current site&amp;lt;/li&amp;gt;&amp;lt;li&amp;gt;`client_country`: Country/Region of the client&amp;lt;/li&amp;gt;&amp;lt;li&amp;gt;`query_string`: Query string in the request URL&amp;lt;/li&amp;gt;&amp;lt;li&amp;gt;`request_header`: HTTP request header&amp;lt;/li&amp;gt;.",
												},
												"values": {
													Type: schema.TypeSet,
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
													Optional:    true,
													Description: "The parameter value of the match type. It can be an empty string only when `Target=query string/request header` and `Operator=exist/notexist`.&amp;lt;li&amp;gt;When `Target=extension`, enter the file extension, such as jpg and txt.&amp;lt;/li&amp;gt;&amp;lt;li&amp;gt;When `Target=filename`, enter the file name, such as foo in foo.jpg.&amp;lt;/li&amp;gt;&amp;lt;li&amp;gt;When `Target=all`, it indicates any site request.&amp;lt;/li&amp;gt;&amp;lt;li&amp;gt;When `Target=host`, enter the host under the current site, such as www.maxx55.com.&amp;lt;/li&amp;gt;&amp;lt;li&amp;gt;When `Target=url`, enter the partial URL path under the current site, such as /example.&amp;lt;/li&amp;gt;&amp;lt;li&amp;gt;When `Target=full_url`, enter the complete URL under the current site. It must contain the HTTP protocol, host, and path, such as https://www.maxx55.cn/example.&amp;lt;/li&amp;gt;&amp;lt;li&amp;gt;When `Target=client_country`, enter the ISO-3166 country/region code.&amp;lt;/li&amp;gt;&amp;lt;li&amp;gt;When `Target=query_string`, enter the value of the query string, such as cn and 1 in lang=cn&amp;amp;version=1.&amp;lt;/li&amp;gt;&amp;lt;li&amp;gt;When `Target=request_header`, enter the HTTP request header value, such as zh-CN,zh;q=0.9 in the Accept-Language:zh-CN,zh;q=0.9 header.&amp;lt;/li&amp;gt;.",
												},
												"ignore_case": {
													Type:        schema.TypeBool,
													Optional:    true,
													Description: "Whether the parameter value is case insensitive. Default value: false.",
												},
												"name": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "The parameter name of the match type. This field is required only when `Target=query_string/request_header`.&amp;lt;li&amp;gt;`query_string`: Name of the query string, such as lang and version in lang=cn&amp;amp;version=1.&amp;lt;/li&amp;gt;&amp;lt;li&amp;gt;`request_header`: Name of the HTTP request header, such as Accept-Language in the Accept-Language:zh-CN,zh;q=0.9 header.&amp;lt;/li&amp;gt;.",
												},
												"ignore_name_case": {
													Type:        schema.TypeBool,
													Optional:    true,
													Description: "Whether the parameter name is case insensitive. Default value: `false`.Note: This field may return null, indicating that no valid values can be obtained.",
												},
											},
										},
									},
								},
							},
						},
						"sub_rules": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "The nested rule.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"rules": {
										Type:        schema.TypeList,
										Required:    true,
										Description: "Nested rule settings.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"conditions": {
													Type:        schema.TypeList,
													Required:    true,
													Description: "The condition that determines if a feature should run.Note: If any condition in the array is met, the feature will run.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"conditions": {
																Type:        schema.TypeList,
																Required:    true,
																Description: "Rule engine condition. This condition will be considered met if all items in the array are met.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"operator": {
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "Operator. Valid values:&amp;lt;li&amp;gt;`equals`: Equals&amp;lt;/li&amp;gt;&amp;lt;li&amp;gt;`notEquals`: Does not equal&amp;lt;/li&amp;gt;&amp;lt;li&amp;gt;`exist`: Exists&amp;lt;/li&amp;gt;&amp;lt;li&amp;gt;`notexist`: Does not exist&amp;lt;/li&amp;gt;.",
																		},
																		"target": {
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "The match type. Values:&amp;lt;li&amp;gt;`filename`: File name&amp;lt;/li&amp;gt;&amp;lt;li&amp;gt;`extension`: File extension&amp;lt;/li&amp;gt;&amp;lt;li&amp;gt;`host`: Host&amp;lt;/li&amp;gt;&amp;lt;li&amp;gt;`full_url`: Full URL, which indicates the complete URL path under the current site and must contain the HTTP protocol, host, and path.&amp;lt;/li&amp;gt;&amp;lt;li&amp;gt;`url`: Partial URL under the current site&amp;lt;/li&amp;gt;&amp;lt;li&amp;gt;`client_country`: Country/Region of the client&amp;lt;/li&amp;gt;&amp;lt;li&amp;gt;`query_string`: Query string in the request URL&amp;lt;/li&amp;gt;&amp;lt;li&amp;gt;`request_header`: HTTP request header&amp;lt;/li&amp;gt;.",
																		},
																		"values": {
																			Type: schema.TypeSet,
																			Elem: &schema.Schema{
																				Type: schema.TypeString,
																			},
																			Optional:    true,
																			Description: "The parameter value of the match type. It can be an empty string only when `Target=query string/request header` and `Operator=exist/notexist`.&amp;lt;li&amp;gt;When `Target=extension`, enter the file extension, such as jpg and txt.&amp;lt;/li&amp;gt;&amp;lt;li&amp;gt;When `Target=filename`, enter the file name, such as foo in foo.jpg.&amp;lt;/li&amp;gt;&amp;lt;li&amp;gt;When `Target=all`, it indicates any site request.&amp;lt;/li&amp;gt;&amp;lt;li&amp;gt;When `Target=host`, enter the host under the current site, such as www.maxx55.com.&amp;lt;/li&amp;gt;&amp;lt;li&amp;gt;When `Target=url`, enter the partial URL path under the current site, such as /example.&amp;lt;/li&amp;gt;&amp;lt;li&amp;gt;When `Target=full_url`, enter the complete URL under the current site. It must contain the HTTP protocol, host, and path, such as https://www.maxx55.cn/example.&amp;lt;/li&amp;gt;&amp;lt;li&amp;gt;When `Target=client_country`, enter the ISO-3166 country/region code.&amp;lt;/li&amp;gt;&amp;lt;li&amp;gt;When `Target=query_string`, enter the value of the query string, such as cn and 1 in lang=cn&amp;amp;version=1.&amp;lt;/li&amp;gt;&amp;lt;li&amp;gt;When `Target=request_header`, enter the HTTP request header value, such as zh-CN,zh;q=0.9 in the Accept-Language:zh-CN,zh;q=0.9 header.&amp;lt;/li&amp;gt;.",
																		},
																		"ignore_case": {
																			Type:        schema.TypeBool,
																			Optional:    true,
																			Description: "Whether the parameter value is case insensitive. Default value: false.",
																		},
																		"name": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "The parameter name of the match type. This field is required only when `Target=query_string/request_header`.&amp;lt;li&amp;gt;`query_string`: Name of the query string, such as lang and version in lang=cn&amp;amp;version=1.&amp;lt;/li&amp;gt;&amp;lt;li&amp;gt;`request_header`: Name of the HTTP request header, such as Accept-Language in the Accept-Language:zh-CN,zh;q=0.9 header.&amp;lt;/li&amp;gt;.",
																		},
																		"ignore_name_case": {
																			Type:        schema.TypeBool,
																			Optional:    true,
																			Description: "Whether the parameter name is case insensitive. Default value: `false`.Note: This field may return null, indicating that no valid values can be obtained.",
																		},
																	},
																},
															},
														},
													},
												},
												"actions": {
													Type:        schema.TypeList,
													Required:    true,
													Description: "The feature to be executed.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"normal_action": {
																Type:        schema.TypeList,
																MaxItems:    1,
																Optional:    true,
																Description: "Common operation. Values:&amp;lt;li&amp;gt;`AccessUrlRedirect`: Access URL rewrite&amp;lt;/li&amp;gt;&amp;lt;li&amp;gt;`UpstreamUrlRedirect`: Origin-pull URL rewrite&amp;lt;/li&amp;gt;&amp;lt;li&amp;gt;`QUIC`: QUIC&amp;lt;/li&amp;gt;&amp;lt;li&amp;gt;`WebSocket`: WebSocket&amp;lt;/li&amp;gt;&amp;lt;li&amp;gt;`VideoSeek`: Video dragging&amp;lt;/li&amp;gt;&amp;lt;li&amp;gt;`Authentication`: Token authentication&amp;lt;/li&amp;gt;&amp;lt;li&amp;gt;`CacheKey`: Custom cache key&amp;lt;/li&amp;gt;&amp;lt;li&amp;gt;`Cache`: Node cache TTL&amp;lt;/li&amp;gt;&amp;lt;li&amp;gt;`MaxAge`: Browser cache TTL&amp;lt;/li&amp;gt;&amp;lt;li&amp;gt;`OfflineCache`: Offline cache&amp;lt;/li&amp;gt;&amp;lt;li&amp;gt;`SmartRouting`: Smart acceleration&amp;lt;/li&amp;gt;&amp;lt;li&amp;gt;`RangeOriginPull`: Range GETs&amp;lt;/li&amp;gt;&amp;lt;li&amp;gt;`UpstreamHttp2`: HTTP/2 forwarding&amp;lt;/li&amp;gt;&amp;lt;li&amp;gt;`HostHeader`: Host header rewrite&amp;lt;/li&amp;gt;&amp;lt;li&amp;gt;`ForceRedirect`: Force HTTPS&amp;lt;/li&amp;gt;&amp;lt;li&amp;gt;`OriginPullProtocol`: Origin-pull HTTPS&amp;lt;/li&amp;gt;&amp;lt;li&amp;gt;`CachePrefresh`: Cache prefresh&amp;lt;/li&amp;gt;&amp;lt;li&amp;gt;`Compression`: Smart compression&amp;lt;/li&amp;gt;&amp;lt;li&amp;gt;`Hsts`&amp;lt;/li&amp;gt;&amp;lt;li&amp;gt;`ClientIpHeader`&amp;lt;/li&amp;gt;&amp;lt;li&amp;gt;`SslTlsSecureConf`&amp;lt;/li&amp;gt;&amp;lt;li&amp;gt;`OcspStapling`&amp;lt;/li&amp;gt;&amp;lt;li&amp;gt;`Http2`: HTTP/2 access&amp;lt;/li&amp;gt;&amp;lt;li&amp;gt;`UpstreamFollowRedirect`: Follow origin redirect&amp;lt;/li&amp;gt;&amp;lt;li&amp;gt;`Origin`: Origin&amp;lt;/li&amp;gt;Note: This field may return `null`, indicating that no valid value can be obtained.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"action": {
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "Feature name. You can call the [DescribeRulesSetting](https://tcloud4api.woa.com/document/product/1657/79433?!preview&amp;amp;!document=1) API to view the requirements for entering the feature name.",
																		},
																		"parameters": {
																			Type:        schema.TypeList,
																			Required:    true,
																			Description: "Parameter.",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"name": {
																						Type:        schema.TypeString,
																						Required:    true,
																						Description: "Parameter name. You can call the [DescribeRulesSetting](https://tcloud4api.woa.com/document/product/1657/79433?!preview&amp;amp;!document=1) API to view the requirements for entering the parameter name.",
																					},
																					"values": {
																						Type: schema.TypeSet,
																						Elem: &schema.Schema{
																							Type: schema.TypeString,
																						},
																						Required:    true,
																						Description: "The parameter value.",
																					},
																				},
																			},
																		},
																	},
																},
															},
															"rewrite_action": {
																Type:        schema.TypeList,
																MaxItems:    1,
																Optional:    true,
																Description: "Feature operation with a request/response header. Features of this type include:&amp;lt;li&amp;gt;`RequestHeader`: HTTP request header modification.&amp;lt;/li&amp;gt;&amp;lt;li&amp;gt;`ResponseHeader`: HTTP response header modification.&amp;lt;/li&amp;gt;Note: This field may return null, indicating that no valid values can be obtained.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"action": {
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "Feature name. You can call the [DescribeRulesSetting](https://tcloud4api.woa.com/document/product/1657/79433?!preview&amp;amp;!document=1) API to view the requirements for entering the feature name.",
																		},
																		"parameters": {
																			Type:        schema.TypeList,
																			Required:    true,
																			Description: "Parameter.",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"action": {
																						Type:        schema.TypeString,
																						Required:    true,
																						Description: "Feature parameter name. You can call the [DescribeRulesSetting](https://tcloud4api.woa.com/document/product/1657/79433?!preview&amp;amp;!document=1) API to view the requirements for entering the parameter name, which has three values:&amp;lt;li&amp;gt;add: Add the HTTP header.&amp;lt;/li&amp;gt;&amp;lt;li&amp;gt;set: Rewrite the HTTP header.&amp;lt;/li&amp;gt;&amp;lt;li&amp;gt;del: Delete the HTTP header.&amp;lt;/li&amp;gt;.",
																					},
																					"name": {
																						Type:        schema.TypeString,
																						Required:    true,
																						Description: "Parameter name.",
																					},
																					"values": {
																						Type: schema.TypeSet,
																						Elem: &schema.Schema{
																							Type: schema.TypeString,
																						},
																						Required:    true,
																						Description: "Parameter value.",
																					},
																				},
																			},
																		},
																	},
																},
															},
															"code_action": {
																Type:        schema.TypeList,
																MaxItems:    1,
																Optional:    true,
																Description: "Feature operation with a status code. Features of this type include:&amp;lt;li&amp;gt;`ErrorPage`: Custom error page.&amp;lt;/li&amp;gt;&amp;lt;li&amp;gt;`StatusCodeCache`: Status code cache TTL.&amp;lt;/li&amp;gt;Note: This field may return null, indicating that no valid values can be obtained.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"action": {
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "Feature name. You can call the [DescribeRulesSetting](https://tcloud4api.woa.com/document/product/1657/79433?!preview&amp;amp;!document=1) API to view the requirements for entering the feature name.",
																		},
																		"parameters": {
																			Type:        schema.TypeList,
																			Required:    true,
																			Description: "Operation parameter.",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"status_code": {
																						Type:        schema.TypeInt,
																						Required:    true,
																						Description: "The status code.",
																					},
																					"name": {
																						Type:        schema.TypeString,
																						Required:    true,
																						Description: "The parameter name. You can call the [DescribeRulesSetting](https://tcloud4api.woa.com/document/product/1657/79433?!preview&amp;amp;!document=1) API to view the requirements for entering the parameter name.",
																					},
																					"values": {
																						Type: schema.TypeSet,
																						Elem: &schema.Schema{
																							Type: schema.TypeString,
																						},
																						Required:    true,
																						Description: "The parameter value.",
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
										},
									},
									"tags": {
										Type: schema.TypeSet,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
										Optional:    true,
										Description: "Tag of the rule.",
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudTeoRule_engineCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_teo_rule_engine.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request  = teo.NewCreateRuleRequest()
		response = teo.NewCreateRuleResponse()
		zoneId   string
		ruleId   string
	)
	if v, ok := d.GetOk("zone_id"); ok {
		zoneId = v.(string)
		request.ZoneId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("rule_name"); ok {
		request.RuleName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("status"); ok {
		request.Status = helper.String(v.(string))
	}

	if v, ok := d.GetOk("rules"); ok {
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			rule := teo.Rule{}
			if v, ok := dMap["actions"]; ok {
				for _, item := range v.([]interface{}) {
					actionsMap := item.(map[string]interface{})
					action := teo.Action{}
					if normalActionMap, ok := helper.InterfaceToMap(actionsMap, "normal_action"); ok {
						normalAction := teo.NormalAction{}
						if v, ok := normalActionMap["action"]; ok {
							normalAction.Action = helper.String(v.(string))
						}
						if v, ok := normalActionMap["parameters"]; ok {
							for _, item := range v.([]interface{}) {
								parametersMap := item.(map[string]interface{})
								ruleNormalActionParams := teo.RuleNormalActionParams{}
								if v, ok := parametersMap["name"]; ok {
									ruleNormalActionParams.Name = helper.String(v.(string))
								}
								if v, ok := parametersMap["values"]; ok {
									valuesSet := v.(*schema.Set).List()
									for i := range valuesSet {
										values := valuesSet[i].(string)
										ruleNormalActionParams.Values = append(ruleNormalActionParams.Values, &values)
									}
								}
								normalAction.Parameters = append(normalAction.Parameters, &ruleNormalActionParams)
							}
						}
						action.NormalAction = &normalAction
					}
					if rewriteActionMap, ok := helper.InterfaceToMap(actionsMap, "rewrite_action"); ok {
						rewriteAction := teo.RewriteAction{}
						if v, ok := rewriteActionMap["action"]; ok {
							rewriteAction.Action = helper.String(v.(string))
						}
						if v, ok := rewriteActionMap["parameters"]; ok {
							for _, item := range v.([]interface{}) {
								parametersMap := item.(map[string]interface{})
								ruleRewriteActionParams := teo.RuleRewriteActionParams{}
								if v, ok := parametersMap["action"]; ok {
									ruleRewriteActionParams.Action = helper.String(v.(string))
								}
								if v, ok := parametersMap["name"]; ok {
									ruleRewriteActionParams.Name = helper.String(v.(string))
								}
								if v, ok := parametersMap["values"]; ok {
									valuesSet := v.(*schema.Set).List()
									for i := range valuesSet {
										values := valuesSet[i].(string)
										ruleRewriteActionParams.Values = append(ruleRewriteActionParams.Values, &values)
									}
								}
								rewriteAction.Parameters = append(rewriteAction.Parameters, &ruleRewriteActionParams)
							}
						}
						action.RewriteAction = &rewriteAction
					}
					if codeActionMap, ok := helper.InterfaceToMap(actionsMap, "code_action"); ok {
						codeAction := teo.CodeAction{}
						if v, ok := codeActionMap["action"]; ok {
							codeAction.Action = helper.String(v.(string))
						}
						if v, ok := codeActionMap["parameters"]; ok {
							for _, item := range v.([]interface{}) {
								parametersMap := item.(map[string]interface{})
								ruleCodeActionParams := teo.RuleCodeActionParams{}
								if v, ok := parametersMap["status_code"]; ok {
									ruleCodeActionParams.StatusCode = helper.IntInt64(v.(int))
								}
								if v, ok := parametersMap["name"]; ok {
									ruleCodeActionParams.Name = helper.String(v.(string))
								}
								if v, ok := parametersMap["values"]; ok {
									valuesSet := v.(*schema.Set).List()
									for i := range valuesSet {
										values := valuesSet[i].(string)
										ruleCodeActionParams.Values = append(ruleCodeActionParams.Values, &values)
									}
								}
								codeAction.Parameters = append(codeAction.Parameters, &ruleCodeActionParams)
							}
						}
						action.CodeAction = &codeAction
					}
					rule.Actions = append(rule.Actions, &action)
				}
			}
			if v, ok := dMap["conditions"]; ok {
				for _, item := range v.([]interface{}) {
					conditionsMap := item.(map[string]interface{})
					ruleAndConditions := teo.RuleAndConditions{}
					if v, ok := conditionsMap["conditions"]; ok {
						for _, item := range v.([]interface{}) {
							conditionsMap := item.(map[string]interface{})
							ruleCondition := teo.RuleCondition{}
							if v, ok := conditionsMap["operator"]; ok {
								ruleCondition.Operator = helper.String(v.(string))
							}
							if v, ok := conditionsMap["target"]; ok {
								ruleCondition.Target = helper.String(v.(string))
							}
							if v, ok := conditionsMap["values"]; ok {
								valuesSet := v.(*schema.Set).List()
								for i := range valuesSet {
									values := valuesSet[i].(string)
									ruleCondition.Values = append(ruleCondition.Values, &values)
								}
							}
							if v, ok := conditionsMap["ignore_case"]; ok {
								ruleCondition.IgnoreCase = helper.Bool(v.(bool))
							}
							if v, ok := conditionsMap["name"]; ok {
								ruleCondition.Name = helper.String(v.(string))
							}
							if v, ok := conditionsMap["ignore_name_case"]; ok {
								ruleCondition.IgnoreNameCase = helper.Bool(v.(bool))
							}
							ruleAndConditions.Conditions = append(ruleAndConditions.Conditions, &ruleCondition)
						}
					}
					rule.Conditions = append(rule.Conditions, &ruleAndConditions)
				}
			}
			if v, ok := dMap["sub_rules"]; ok {
				for _, item := range v.([]interface{}) {
					subRulesMap := item.(map[string]interface{})
					subRuleItem := teo.SubRuleItem{}
					if v, ok := subRulesMap["rules"]; ok {
						for _, item := range v.([]interface{}) {
							rulesMap := item.(map[string]interface{})
							subRule := teo.SubRule{}
							if v, ok := rulesMap["conditions"]; ok {
								for _, item := range v.([]interface{}) {
									conditionsMap := item.(map[string]interface{})
									ruleAndConditions := teo.RuleAndConditions{}
									if v, ok := conditionsMap["conditions"]; ok {
										for _, item := range v.([]interface{}) {
											conditionsMap := item.(map[string]interface{})
											ruleCondition := teo.RuleCondition{}
											if v, ok := conditionsMap["operator"]; ok {
												ruleCondition.Operator = helper.String(v.(string))
											}
											if v, ok := conditionsMap["target"]; ok {
												ruleCondition.Target = helper.String(v.(string))
											}
											if v, ok := conditionsMap["values"]; ok {
												valuesSet := v.(*schema.Set).List()
												for i := range valuesSet {
													values := valuesSet[i].(string)
													ruleCondition.Values = append(ruleCondition.Values, &values)
												}
											}
											if v, ok := conditionsMap["ignore_case"]; ok {
												ruleCondition.IgnoreCase = helper.Bool(v.(bool))
											}
											if v, ok := conditionsMap["name"]; ok {
												ruleCondition.Name = helper.String(v.(string))
											}
											if v, ok := conditionsMap["ignore_name_case"]; ok {
												ruleCondition.IgnoreNameCase = helper.Bool(v.(bool))
											}
											ruleAndConditions.Conditions = append(ruleAndConditions.Conditions, &ruleCondition)
										}
									}
									subRule.Conditions = append(subRule.Conditions, &ruleAndConditions)
								}
							}
							if v, ok := rulesMap["actions"]; ok {
								for _, item := range v.([]interface{}) {
									actionsMap := item.(map[string]interface{})
									action := teo.Action{}
									if normalActionMap, ok := helper.InterfaceToMap(actionsMap, "normal_action"); ok {
										normalAction := teo.NormalAction{}
										if v, ok := normalActionMap["action"]; ok {
											normalAction.Action = helper.String(v.(string))
										}
										if v, ok := normalActionMap["parameters"]; ok {
											for _, item := range v.([]interface{}) {
												parametersMap := item.(map[string]interface{})
												ruleNormalActionParams := teo.RuleNormalActionParams{}
												if v, ok := parametersMap["name"]; ok {
													ruleNormalActionParams.Name = helper.String(v.(string))
												}
												if v, ok := parametersMap["values"]; ok {
													valuesSet := v.(*schema.Set).List()
													for i := range valuesSet {
														values := valuesSet[i].(string)
														ruleNormalActionParams.Values = append(ruleNormalActionParams.Values, &values)
													}
												}
												normalAction.Parameters = append(normalAction.Parameters, &ruleNormalActionParams)
											}
										}
										action.NormalAction = &normalAction
									}
									if rewriteActionMap, ok := helper.InterfaceToMap(actionsMap, "rewrite_action"); ok {
										rewriteAction := teo.RewriteAction{}
										if v, ok := rewriteActionMap["action"]; ok {
											rewriteAction.Action = helper.String(v.(string))
										}
										if v, ok := rewriteActionMap["parameters"]; ok {
											for _, item := range v.([]interface{}) {
												parametersMap := item.(map[string]interface{})
												ruleRewriteActionParams := teo.RuleRewriteActionParams{}
												if v, ok := parametersMap["action"]; ok {
													ruleRewriteActionParams.Action = helper.String(v.(string))
												}
												if v, ok := parametersMap["name"]; ok {
													ruleRewriteActionParams.Name = helper.String(v.(string))
												}
												if v, ok := parametersMap["values"]; ok {
													valuesSet := v.(*schema.Set).List()
													for i := range valuesSet {
														values := valuesSet[i].(string)
														ruleRewriteActionParams.Values = append(ruleRewriteActionParams.Values, &values)
													}
												}
												rewriteAction.Parameters = append(rewriteAction.Parameters, &ruleRewriteActionParams)
											}
										}
										action.RewriteAction = &rewriteAction
									}
									if codeActionMap, ok := helper.InterfaceToMap(actionsMap, "code_action"); ok {
										codeAction := teo.CodeAction{}
										if v, ok := codeActionMap["action"]; ok {
											codeAction.Action = helper.String(v.(string))
										}
										if v, ok := codeActionMap["parameters"]; ok {
											for _, item := range v.([]interface{}) {
												parametersMap := item.(map[string]interface{})
												ruleCodeActionParams := teo.RuleCodeActionParams{}
												if v, ok := parametersMap["status_code"]; ok {
													ruleCodeActionParams.StatusCode = helper.IntInt64(v.(int))
												}
												if v, ok := parametersMap["name"]; ok {
													ruleCodeActionParams.Name = helper.String(v.(string))
												}
												if v, ok := parametersMap["values"]; ok {
													valuesSet := v.(*schema.Set).List()
													for i := range valuesSet {
														values := valuesSet[i].(string)
														ruleCodeActionParams.Values = append(ruleCodeActionParams.Values, &values)
													}
												}
												codeAction.Parameters = append(codeAction.Parameters, &ruleCodeActionParams)
											}
										}
										action.CodeAction = &codeAction
									}
									subRule.Actions = append(subRule.Actions, &action)
								}
							}
							subRuleItem.Rules = append(subRuleItem.Rules, &subRule)
						}
					}
					if v, ok := subRulesMap["tags"]; ok {
						tagsSet := v.(*schema.Set).List()
						for i := range tagsSet {
							tags := tagsSet[i].(string)
							subRuleItem.Tags = append(subRuleItem.Tags, &tags)
						}
					}
					rule.SubRules = append(rule.SubRules, &subRuleItem)
				}
			}
			request.Rules = append(request.Rules, &rule)
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTeoClient().CreateRule(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create teo rule_engine failed, reason:%+v", logId, err)
		return err
	}

	zoneId = *response.Response.ZoneId
	d.SetId(strings.Join([]string{zoneId, ruleId}, FILED_SP))

	return resourceTencentCloudTeoRule_engineRead(d, meta)
}

func resourceTencentCloudTeoRule_engineRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_teo_rule_engine.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TeoService{client: meta.(*TencentCloudClient).apiV3Conn}

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	zoneId := idSplit[0]
	ruleId := idSplit[1]

	rule_engine, err := service.DescribeTeoRule_engineById(ctx, zoneId, ruleId)
	if err != nil {
		return err
	}

	if rule_engine == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `TeoRule_engine` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if rule_engine.ZoneId != nil {
		_ = d.Set("zone_id", rule_engine.ZoneId)
	}

	if rule_engine.RuleName != nil {
		_ = d.Set("rule_name", rule_engine.RuleName)
	}

	if rule_engine.Status != nil {
		_ = d.Set("status", rule_engine.Status)
	}

	if rule_engine.Rules != nil {
		rulesList := []interface{}{}
		for _, rules := range rule_engine.Rules {
			rulesMap := map[string]interface{}{}

			if rule_engine.Rules.Actions != nil {
				actionsList := []interface{}{}
				for _, actions := range rule_engine.Rules.Actions {
					actionsMap := map[string]interface{}{}

					if actions.NormalAction != nil {
						normalActionMap := map[string]interface{}{}

						if actions.NormalAction.Action != nil {
							normalActionMap["action"] = actions.NormalAction.Action
						}

						if actions.NormalAction.Parameters != nil {
							parametersList := []interface{}{}
							for _, parameters := range actions.NormalAction.Parameters {
								parametersMap := map[string]interface{}{}

								if parameters.Name != nil {
									parametersMap["name"] = parameters.Name
								}

								if parameters.Values != nil {
									parametersMap["values"] = parameters.Values
								}

								parametersList = append(parametersList, parametersMap)
							}

							normalActionMap["parameters"] = []interface{}{parametersList}
						}

						actionsMap["normal_action"] = []interface{}{normalActionMap}
					}

					if actions.RewriteAction != nil {
						rewriteActionMap := map[string]interface{}{}

						if actions.RewriteAction.Action != nil {
							rewriteActionMap["action"] = actions.RewriteAction.Action
						}

						if actions.RewriteAction.Parameters != nil {
							parametersList := []interface{}{}
							for _, parameters := range actions.RewriteAction.Parameters {
								parametersMap := map[string]interface{}{}

								if parameters.Action != nil {
									parametersMap["action"] = parameters.Action
								}

								if parameters.Name != nil {
									parametersMap["name"] = parameters.Name
								}

								if parameters.Values != nil {
									parametersMap["values"] = parameters.Values
								}

								parametersList = append(parametersList, parametersMap)
							}

							rewriteActionMap["parameters"] = []interface{}{parametersList}
						}

						actionsMap["rewrite_action"] = []interface{}{rewriteActionMap}
					}

					if actions.CodeAction != nil {
						codeActionMap := map[string]interface{}{}

						if actions.CodeAction.Action != nil {
							codeActionMap["action"] = actions.CodeAction.Action
						}

						if actions.CodeAction.Parameters != nil {
							parametersList := []interface{}{}
							for _, parameters := range actions.CodeAction.Parameters {
								parametersMap := map[string]interface{}{}

								if parameters.StatusCode != nil {
									parametersMap["status_code"] = parameters.StatusCode
								}

								if parameters.Name != nil {
									parametersMap["name"] = parameters.Name
								}

								if parameters.Values != nil {
									parametersMap["values"] = parameters.Values
								}

								parametersList = append(parametersList, parametersMap)
							}

							codeActionMap["parameters"] = []interface{}{parametersList}
						}

						actionsMap["code_action"] = []interface{}{codeActionMap}
					}

					actionsList = append(actionsList, actionsMap)
				}

				rulesMap["actions"] = []interface{}{actionsList}
			}

			if rule_engine.Rules.Conditions != nil {
				conditionsList := []interface{}{}
				for _, conditions := range rule_engine.Rules.Conditions {
					conditionsMap := map[string]interface{}{}

					if conditions.Conditions != nil {
						conditionsList := []interface{}{}
						for _, conditions := range conditions.Conditions {
							conditionsMap := map[string]interface{}{}

							if conditions.Operator != nil {
								conditionsMap["operator"] = conditions.Operator
							}

							if conditions.Target != nil {
								conditionsMap["target"] = conditions.Target
							}

							if conditions.Values != nil {
								conditionsMap["values"] = conditions.Values
							}

							if conditions.IgnoreCase != nil {
								conditionsMap["ignore_case"] = conditions.IgnoreCase
							}

							if conditions.Name != nil {
								conditionsMap["name"] = conditions.Name
							}

							if conditions.IgnoreNameCase != nil {
								conditionsMap["ignore_name_case"] = conditions.IgnoreNameCase
							}

							conditionsList = append(conditionsList, conditionsMap)
						}

						conditionsMap["conditions"] = []interface{}{conditionsList}
					}

					conditionsList = append(conditionsList, conditionsMap)
				}

				rulesMap["conditions"] = []interface{}{conditionsList}
			}

			if rule_engine.Rules.SubRules != nil {
				subRulesList := []interface{}{}
				for _, subRules := range rule_engine.Rules.SubRules {
					subRulesMap := map[string]interface{}{}

					if subRules.Rules != nil {
						rulesList := []interface{}{}
						for _, rules := range subRules.Rules {
							rulesMap := map[string]interface{}{}

							if rules.Conditions != nil {
								conditionsList := []interface{}{}
								for _, conditions := range rules.Conditions {
									conditionsMap := map[string]interface{}{}

									if conditions.Conditions != nil {
										conditionsList := []interface{}{}
										for _, conditions := range conditions.Conditions {
											conditionsMap := map[string]interface{}{}

											if conditions.Operator != nil {
												conditionsMap["operator"] = conditions.Operator
											}

											if conditions.Target != nil {
												conditionsMap["target"] = conditions.Target
											}

											if conditions.Values != nil {
												conditionsMap["values"] = conditions.Values
											}

											if conditions.IgnoreCase != nil {
												conditionsMap["ignore_case"] = conditions.IgnoreCase
											}

											if conditions.Name != nil {
												conditionsMap["name"] = conditions.Name
											}

											if conditions.IgnoreNameCase != nil {
												conditionsMap["ignore_name_case"] = conditions.IgnoreNameCase
											}

											conditionsList = append(conditionsList, conditionsMap)
										}

										conditionsMap["conditions"] = []interface{}{conditionsList}
									}

									conditionsList = append(conditionsList, conditionsMap)
								}

								rulesMap["conditions"] = []interface{}{conditionsList}
							}

							if rules.Actions != nil {
								actionsList := []interface{}{}
								for _, actions := range rules.Actions {
									actionsMap := map[string]interface{}{}

									if actions.NormalAction != nil {
										normalActionMap := map[string]interface{}{}

										if actions.NormalAction.Action != nil {
											normalActionMap["action"] = actions.NormalAction.Action
										}

										if actions.NormalAction.Parameters != nil {
											parametersList := []interface{}{}
											for _, parameters := range actions.NormalAction.Parameters {
												parametersMap := map[string]interface{}{}

												if parameters.Name != nil {
													parametersMap["name"] = parameters.Name
												}

												if parameters.Values != nil {
													parametersMap["values"] = parameters.Values
												}

												parametersList = append(parametersList, parametersMap)
											}

											normalActionMap["parameters"] = []interface{}{parametersList}
										}

										actionsMap["normal_action"] = []interface{}{normalActionMap}
									}

									if actions.RewriteAction != nil {
										rewriteActionMap := map[string]interface{}{}

										if actions.RewriteAction.Action != nil {
											rewriteActionMap["action"] = actions.RewriteAction.Action
										}

										if actions.RewriteAction.Parameters != nil {
											parametersList := []interface{}{}
											for _, parameters := range actions.RewriteAction.Parameters {
												parametersMap := map[string]interface{}{}

												if parameters.Action != nil {
													parametersMap["action"] = parameters.Action
												}

												if parameters.Name != nil {
													parametersMap["name"] = parameters.Name
												}

												if parameters.Values != nil {
													parametersMap["values"] = parameters.Values
												}

												parametersList = append(parametersList, parametersMap)
											}

											rewriteActionMap["parameters"] = []interface{}{parametersList}
										}

										actionsMap["rewrite_action"] = []interface{}{rewriteActionMap}
									}

									if actions.CodeAction != nil {
										codeActionMap := map[string]interface{}{}

										if actions.CodeAction.Action != nil {
											codeActionMap["action"] = actions.CodeAction.Action
										}

										if actions.CodeAction.Parameters != nil {
											parametersList := []interface{}{}
											for _, parameters := range actions.CodeAction.Parameters {
												parametersMap := map[string]interface{}{}

												if parameters.StatusCode != nil {
													parametersMap["status_code"] = parameters.StatusCode
												}

												if parameters.Name != nil {
													parametersMap["name"] = parameters.Name
												}

												if parameters.Values != nil {
													parametersMap["values"] = parameters.Values
												}

												parametersList = append(parametersList, parametersMap)
											}

											codeActionMap["parameters"] = []interface{}{parametersList}
										}

										actionsMap["code_action"] = []interface{}{codeActionMap}
									}

									actionsList = append(actionsList, actionsMap)
								}

								rulesMap["actions"] = []interface{}{actionsList}
							}

							rulesList = append(rulesList, rulesMap)
						}

						subRulesMap["rules"] = []interface{}{rulesList}
					}

					if subRules.Tags != nil {
						subRulesMap["tags"] = subRules.Tags
					}

					subRulesList = append(subRulesList, subRulesMap)
				}

				rulesMap["sub_rules"] = []interface{}{subRulesList}
			}

			rulesList = append(rulesList, rulesMap)
		}

		_ = d.Set("rules", rulesList)

	}

	return nil
}

func resourceTencentCloudTeoRule_engineUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_teo_rule_engine.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := teo.NewModifyRuleRequest()

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	zoneId := idSplit[0]
	ruleId := idSplit[1]

	request.ZoneId = &zoneId
	request.RuleId = &ruleId

	immutableArgs := []string{"zone_id", "rule_name", "status", "rules"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	if d.HasChange("zone_id") {
		if v, ok := d.GetOk("zone_id"); ok {
			request.ZoneId = helper.String(v.(string))
		}
	}

	if d.HasChange("rule_name") {
		if v, ok := d.GetOk("rule_name"); ok {
			request.RuleName = helper.String(v.(string))
		}
	}

	if d.HasChange("status") {
		if v, ok := d.GetOk("status"); ok {
			request.Status = helper.String(v.(string))
		}
	}

	if d.HasChange("rules") {
		if v, ok := d.GetOk("rules"); ok {
			for _, item := range v.([]interface{}) {
				rule := teo.Rule{}
				if v, ok := dMap["actions"]; ok {
					for _, item := range v.([]interface{}) {
						actionsMap := item.(map[string]interface{})
						action := teo.Action{}
						if normalActionMap, ok := helper.InterfaceToMap(actionsMap, "normal_action"); ok {
							normalAction := teo.NormalAction{}
							if v, ok := normalActionMap["action"]; ok {
								normalAction.Action = helper.String(v.(string))
							}
							if v, ok := normalActionMap["parameters"]; ok {
								for _, item := range v.([]interface{}) {
									parametersMap := item.(map[string]interface{})
									ruleNormalActionParams := teo.RuleNormalActionParams{}
									if v, ok := parametersMap["name"]; ok {
										ruleNormalActionParams.Name = helper.String(v.(string))
									}
									if v, ok := parametersMap["values"]; ok {
										valuesSet := v.(*schema.Set).List()
										for i := range valuesSet {
											values := valuesSet[i].(string)
											ruleNormalActionParams.Values = append(ruleNormalActionParams.Values, &values)
										}
									}
									normalAction.Parameters = append(normalAction.Parameters, &ruleNormalActionParams)
								}
							}
							action.NormalAction = &normalAction
						}
						if rewriteActionMap, ok := helper.InterfaceToMap(actionsMap, "rewrite_action"); ok {
							rewriteAction := teo.RewriteAction{}
							if v, ok := rewriteActionMap["action"]; ok {
								rewriteAction.Action = helper.String(v.(string))
							}
							if v, ok := rewriteActionMap["parameters"]; ok {
								for _, item := range v.([]interface{}) {
									parametersMap := item.(map[string]interface{})
									ruleRewriteActionParams := teo.RuleRewriteActionParams{}
									if v, ok := parametersMap["action"]; ok {
										ruleRewriteActionParams.Action = helper.String(v.(string))
									}
									if v, ok := parametersMap["name"]; ok {
										ruleRewriteActionParams.Name = helper.String(v.(string))
									}
									if v, ok := parametersMap["values"]; ok {
										valuesSet := v.(*schema.Set).List()
										for i := range valuesSet {
											values := valuesSet[i].(string)
											ruleRewriteActionParams.Values = append(ruleRewriteActionParams.Values, &values)
										}
									}
									rewriteAction.Parameters = append(rewriteAction.Parameters, &ruleRewriteActionParams)
								}
							}
							action.RewriteAction = &rewriteAction
						}
						if codeActionMap, ok := helper.InterfaceToMap(actionsMap, "code_action"); ok {
							codeAction := teo.CodeAction{}
							if v, ok := codeActionMap["action"]; ok {
								codeAction.Action = helper.String(v.(string))
							}
							if v, ok := codeActionMap["parameters"]; ok {
								for _, item := range v.([]interface{}) {
									parametersMap := item.(map[string]interface{})
									ruleCodeActionParams := teo.RuleCodeActionParams{}
									if v, ok := parametersMap["status_code"]; ok {
										ruleCodeActionParams.StatusCode = helper.IntInt64(v.(int))
									}
									if v, ok := parametersMap["name"]; ok {
										ruleCodeActionParams.Name = helper.String(v.(string))
									}
									if v, ok := parametersMap["values"]; ok {
										valuesSet := v.(*schema.Set).List()
										for i := range valuesSet {
											values := valuesSet[i].(string)
											ruleCodeActionParams.Values = append(ruleCodeActionParams.Values, &values)
										}
									}
									codeAction.Parameters = append(codeAction.Parameters, &ruleCodeActionParams)
								}
							}
							action.CodeAction = &codeAction
						}
						rule.Actions = append(rule.Actions, &action)
					}
				}
				if v, ok := dMap["conditions"]; ok {
					for _, item := range v.([]interface{}) {
						conditionsMap := item.(map[string]interface{})
						ruleAndConditions := teo.RuleAndConditions{}
						if v, ok := conditionsMap["conditions"]; ok {
							for _, item := range v.([]interface{}) {
								conditionsMap := item.(map[string]interface{})
								ruleCondition := teo.RuleCondition{}
								if v, ok := conditionsMap["operator"]; ok {
									ruleCondition.Operator = helper.String(v.(string))
								}
								if v, ok := conditionsMap["target"]; ok {
									ruleCondition.Target = helper.String(v.(string))
								}
								if v, ok := conditionsMap["values"]; ok {
									valuesSet := v.(*schema.Set).List()
									for i := range valuesSet {
										values := valuesSet[i].(string)
										ruleCondition.Values = append(ruleCondition.Values, &values)
									}
								}
								if v, ok := conditionsMap["ignore_case"]; ok {
									ruleCondition.IgnoreCase = helper.Bool(v.(bool))
								}
								if v, ok := conditionsMap["name"]; ok {
									ruleCondition.Name = helper.String(v.(string))
								}
								if v, ok := conditionsMap["ignore_name_case"]; ok {
									ruleCondition.IgnoreNameCase = helper.Bool(v.(bool))
								}
								ruleAndConditions.Conditions = append(ruleAndConditions.Conditions, &ruleCondition)
							}
						}
						rule.Conditions = append(rule.Conditions, &ruleAndConditions)
					}
				}
				if v, ok := dMap["sub_rules"]; ok {
					for _, item := range v.([]interface{}) {
						subRulesMap := item.(map[string]interface{})
						subRuleItem := teo.SubRuleItem{}
						if v, ok := subRulesMap["rules"]; ok {
							for _, item := range v.([]interface{}) {
								rulesMap := item.(map[string]interface{})
								subRule := teo.SubRule{}
								if v, ok := rulesMap["conditions"]; ok {
									for _, item := range v.([]interface{}) {
										conditionsMap := item.(map[string]interface{})
										ruleAndConditions := teo.RuleAndConditions{}
										if v, ok := conditionsMap["conditions"]; ok {
											for _, item := range v.([]interface{}) {
												conditionsMap := item.(map[string]interface{})
												ruleCondition := teo.RuleCondition{}
												if v, ok := conditionsMap["operator"]; ok {
													ruleCondition.Operator = helper.String(v.(string))
												}
												if v, ok := conditionsMap["target"]; ok {
													ruleCondition.Target = helper.String(v.(string))
												}
												if v, ok := conditionsMap["values"]; ok {
													valuesSet := v.(*schema.Set).List()
													for i := range valuesSet {
														values := valuesSet[i].(string)
														ruleCondition.Values = append(ruleCondition.Values, &values)
													}
												}
												if v, ok := conditionsMap["ignore_case"]; ok {
													ruleCondition.IgnoreCase = helper.Bool(v.(bool))
												}
												if v, ok := conditionsMap["name"]; ok {
													ruleCondition.Name = helper.String(v.(string))
												}
												if v, ok := conditionsMap["ignore_name_case"]; ok {
													ruleCondition.IgnoreNameCase = helper.Bool(v.(bool))
												}
												ruleAndConditions.Conditions = append(ruleAndConditions.Conditions, &ruleCondition)
											}
										}
										subRule.Conditions = append(subRule.Conditions, &ruleAndConditions)
									}
								}
								if v, ok := rulesMap["actions"]; ok {
									for _, item := range v.([]interface{}) {
										actionsMap := item.(map[string]interface{})
										action := teo.Action{}
										if normalActionMap, ok := helper.InterfaceToMap(actionsMap, "normal_action"); ok {
											normalAction := teo.NormalAction{}
											if v, ok := normalActionMap["action"]; ok {
												normalAction.Action = helper.String(v.(string))
											}
											if v, ok := normalActionMap["parameters"]; ok {
												for _, item := range v.([]interface{}) {
													parametersMap := item.(map[string]interface{})
													ruleNormalActionParams := teo.RuleNormalActionParams{}
													if v, ok := parametersMap["name"]; ok {
														ruleNormalActionParams.Name = helper.String(v.(string))
													}
													if v, ok := parametersMap["values"]; ok {
														valuesSet := v.(*schema.Set).List()
														for i := range valuesSet {
															values := valuesSet[i].(string)
															ruleNormalActionParams.Values = append(ruleNormalActionParams.Values, &values)
														}
													}
													normalAction.Parameters = append(normalAction.Parameters, &ruleNormalActionParams)
												}
											}
											action.NormalAction = &normalAction
										}
										if rewriteActionMap, ok := helper.InterfaceToMap(actionsMap, "rewrite_action"); ok {
											rewriteAction := teo.RewriteAction{}
											if v, ok := rewriteActionMap["action"]; ok {
												rewriteAction.Action = helper.String(v.(string))
											}
											if v, ok := rewriteActionMap["parameters"]; ok {
												for _, item := range v.([]interface{}) {
													parametersMap := item.(map[string]interface{})
													ruleRewriteActionParams := teo.RuleRewriteActionParams{}
													if v, ok := parametersMap["action"]; ok {
														ruleRewriteActionParams.Action = helper.String(v.(string))
													}
													if v, ok := parametersMap["name"]; ok {
														ruleRewriteActionParams.Name = helper.String(v.(string))
													}
													if v, ok := parametersMap["values"]; ok {
														valuesSet := v.(*schema.Set).List()
														for i := range valuesSet {
															values := valuesSet[i].(string)
															ruleRewriteActionParams.Values = append(ruleRewriteActionParams.Values, &values)
														}
													}
													rewriteAction.Parameters = append(rewriteAction.Parameters, &ruleRewriteActionParams)
												}
											}
											action.RewriteAction = &rewriteAction
										}
										if codeActionMap, ok := helper.InterfaceToMap(actionsMap, "code_action"); ok {
											codeAction := teo.CodeAction{}
											if v, ok := codeActionMap["action"]; ok {
												codeAction.Action = helper.String(v.(string))
											}
											if v, ok := codeActionMap["parameters"]; ok {
												for _, item := range v.([]interface{}) {
													parametersMap := item.(map[string]interface{})
													ruleCodeActionParams := teo.RuleCodeActionParams{}
													if v, ok := parametersMap["status_code"]; ok {
														ruleCodeActionParams.StatusCode = helper.IntInt64(v.(int))
													}
													if v, ok := parametersMap["name"]; ok {
														ruleCodeActionParams.Name = helper.String(v.(string))
													}
													if v, ok := parametersMap["values"]; ok {
														valuesSet := v.(*schema.Set).List()
														for i := range valuesSet {
															values := valuesSet[i].(string)
															ruleCodeActionParams.Values = append(ruleCodeActionParams.Values, &values)
														}
													}
													codeAction.Parameters = append(codeAction.Parameters, &ruleCodeActionParams)
												}
											}
											action.CodeAction = &codeAction
										}
										subRule.Actions = append(subRule.Actions, &action)
									}
								}
								subRuleItem.Rules = append(subRuleItem.Rules, &subRule)
							}
						}
						if v, ok := subRulesMap["tags"]; ok {
							tagsSet := v.(*schema.Set).List()
							for i := range tagsSet {
								tags := tagsSet[i].(string)
								subRuleItem.Tags = append(subRuleItem.Tags, &tags)
							}
						}
						rule.SubRules = append(rule.SubRules, &subRuleItem)
					}
				}
				request.Rules = append(request.Rules, &rule)
			}
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTeoClient().ModifyRule(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update teo rule_engine failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudTeoRule_engineRead(d, meta)
}

func resourceTencentCloudTeoRule_engineDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_teo_rule_engine.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TeoService{client: meta.(*TencentCloudClient).apiV3Conn}
	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	zoneId := idSplit[0]
	ruleId := idSplit[1]

	if err := service.DeleteTeoRule_engineById(ctx, zoneId, ruleId); err != nil {
		return err
	}

	return nil
}
