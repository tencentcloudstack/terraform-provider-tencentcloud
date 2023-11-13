/*
Provides a resource to create a tem gateway

Example Usage

```hcl
resource "tencentcloud_tem_gateway" "gateway" {
  ingress {
		ingress_name = "en-xxx"
		environment_id = "en-xxx"
		cluster_namespace = "default"
		address_ip_version = "IPV4"
		rewrite_type = "AUTO"
		mixed = false
		tls {
			hosts =
			secret_name = &lt;nil&gt;
			certificate_id = &lt;nil&gt;
		}
		rules {
			host = &lt;nil&gt;
			protocol = "http"
			http {
				paths {
					path = &lt;nil&gt;
					backend {
						service_name = &lt;nil&gt;
						service_port = &lt;nil&gt;
					}
				}
			}
		}
		clb_id = "xxx"

  }
}
```

Import

tem gateway can be imported using the id, e.g.

```
terraform import tencentcloud_tem_gateway.gateway gateway_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tem "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tem/v20210701"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
	"strings"
)

func resourceTencentCloudTemGateway() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTemGatewayCreate,
		Read:   resourceTencentCloudTemGatewayRead,
		Update: resourceTencentCloudTemGatewayUpdate,
		Delete: resourceTencentCloudTemGatewayDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"ingress": {
				Optional:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Gateway properties.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ingress_name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Gateway name.",
						},
						"environment_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Environment ID.",
						},
						"cluster_namespace": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Inner namespace, default only.",
						},
						"address_ip_version": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Ip version, support IPV4.",
						},
						"rewrite_type": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Redirect mode, support AUTO and NONE.",
						},
						"mixed": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "Mixing HTTP and HTTPS.",
						},
						"tls": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Ingress TLS configurations.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"hosts": {
										Type: schema.TypeSet,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
										Required:    true,
										Description: "Host names.",
									},
									"secret_name": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Secret name.",
									},
									"certificate_id": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Certificate ID.",
									},
								},
							},
						},
						"rules": {
							Type:        schema.TypeList,
							Required:    true,
							Description: "Proxy rules.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"host": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Host name.",
									},
									"protocol": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Protocol.",
									},
									"http": {
										Type:        schema.TypeList,
										MaxItems:    1,
										Required:    true,
										Description: "Rule payload.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"paths": {
													Type:        schema.TypeList,
													Required:    true,
													Description: "Path payload.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"path": {
																Type:        schema.TypeString,
																Required:    true,
																Description: "Path.",
															},
															"backend": {
																Type:        schema.TypeList,
																MaxItems:    1,
																Required:    true,
																Description: "Backend payload.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"service_name": {
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "Backend name.",
																		},
																		"service_port": {
																			Type:        schema.TypeInt,
																			Required:    true,
																			Description: "Backend port.",
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
						"vip": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Gateway vip.",
						},
						"clb_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Related CLB ID.",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Creation time.",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudTemGatewayCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tem_gateway.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request       = tem.NewModifyIngressRequest()
		response      = tem.NewModifyIngressResponse()
		environmentId string
		ingressName   string
	)
	if dMap, ok := helper.InterfacesHeadMap(d, "ingress"); ok {
		ingressInfo := tem.IngressInfo{}
		if v, ok := dMap["ingress_name"]; ok {
			ingressInfo.IngressName = helper.String(v.(string))
		}
		if v, ok := dMap["environment_id"]; ok {
			ingressInfo.EnvironmentId = helper.String(v.(string))
		}
		if v, ok := dMap["cluster_namespace"]; ok {
			ingressInfo.ClusterNamespace = helper.String(v.(string))
		}
		if v, ok := dMap["address_ip_version"]; ok {
			ingressInfo.AddressIpVersion = helper.String(v.(string))
		}
		if v, ok := dMap["rewrite_type"]; ok {
			ingressInfo.RewriteType = helper.String(v.(string))
		}
		if v, ok := dMap["mixed"]; ok {
			ingressInfo.Mixed = helper.Bool(v.(bool))
		}
		if v, ok := dMap["tls"]; ok {
			for _, item := range v.([]interface{}) {
				tlsMap := item.(map[string]interface{})
				ingressTls := tem.IngressTls{}
				if v, ok := tlsMap["hosts"]; ok {
					hostsSet := v.(*schema.Set).List()
					for i := range hostsSet {
						hosts := hostsSet[i].(string)
						ingressTls.Hosts = append(ingressTls.Hosts, &hosts)
					}
				}
				if v, ok := tlsMap["secret_name"]; ok {
					ingressTls.SecretName = helper.String(v.(string))
				}
				if v, ok := tlsMap["certificate_id"]; ok {
					ingressTls.CertificateId = helper.String(v.(string))
				}
				ingressInfo.Tls = append(ingressInfo.Tls, &ingressTls)
			}
		}
		if v, ok := dMap["rules"]; ok {
			for _, item := range v.([]interface{}) {
				rulesMap := item.(map[string]interface{})
				ingressRule := tem.IngressRule{}
				if v, ok := rulesMap["host"]; ok {
					ingressRule.Host = helper.String(v.(string))
				}
				if v, ok := rulesMap["protocol"]; ok {
					ingressRule.Protocol = helper.String(v.(string))
				}
				if httpMap, ok := helper.InterfaceToMap(rulesMap, "http"); ok {
					ingressRuleValue := tem.IngressRuleValue{}
					if v, ok := httpMap["paths"]; ok {
						for _, item := range v.([]interface{}) {
							pathsMap := item.(map[string]interface{})
							ingressRulePath := tem.IngressRulePath{}
							if v, ok := pathsMap["path"]; ok {
								ingressRulePath.Path = helper.String(v.(string))
							}
							if backendMap, ok := helper.InterfaceToMap(pathsMap, "backend"); ok {
								ingressRuleBackend := tem.IngressRuleBackend{}
								if v, ok := backendMap["service_name"]; ok {
									ingressRuleBackend.ServiceName = helper.String(v.(string))
								}
								if v, ok := backendMap["service_port"]; ok {
									ingressRuleBackend.ServicePort = helper.IntInt64(v.(int))
								}
								ingressRulePath.Backend = &ingressRuleBackend
							}
							ingressRuleValue.Paths = append(ingressRuleValue.Paths, &ingressRulePath)
						}
					}
					ingressRule.Http = &ingressRuleValue
				}
				ingressInfo.Rules = append(ingressInfo.Rules, &ingressRule)
			}
		}
		if v, ok := dMap["clb_id"]; ok {
			ingressInfo.ClbId = helper.String(v.(string))
		}
		request.Ingress = &ingressInfo
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTemClient().ModifyIngress(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create tem gateway failed, reason:%+v", logId, err)
		return err
	}

	environmentId = *response.Response.EnvironmentId
	d.SetId(strings.Join([]string{environmentId, ingressName}, FILED_SP))

	return resourceTencentCloudTemGatewayRead(d, meta)
}

func resourceTencentCloudTemGatewayRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tem_gateway.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TemService{client: meta.(*TencentCloudClient).apiV3Conn}

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	environmentId := idSplit[0]
	ingressName := idSplit[1]

	gateway, err := service.DescribeTemGatewayById(ctx, environmentId, ingressName)
	if err != nil {
		return err
	}

	if gateway == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `TemGateway` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if gateway.Ingress != nil {
		ingressMap := map[string]interface{}{}

		if gateway.Ingress.IngressName != nil {
			ingressMap["ingress_name"] = gateway.Ingress.IngressName
		}

		if gateway.Ingress.EnvironmentId != nil {
			ingressMap["environment_id"] = gateway.Ingress.EnvironmentId
		}

		if gateway.Ingress.ClusterNamespace != nil {
			ingressMap["cluster_namespace"] = gateway.Ingress.ClusterNamespace
		}

		if gateway.Ingress.AddressIpVersion != nil {
			ingressMap["address_ip_version"] = gateway.Ingress.AddressIpVersion
		}

		if gateway.Ingress.RewriteType != nil {
			ingressMap["rewrite_type"] = gateway.Ingress.RewriteType
		}

		if gateway.Ingress.Mixed != nil {
			ingressMap["mixed"] = gateway.Ingress.Mixed
		}

		if gateway.Ingress.Tls != nil {
			tlsList := []interface{}{}
			for _, tls := range gateway.Ingress.Tls {
				tlsMap := map[string]interface{}{}

				if tls.Hosts != nil {
					tlsMap["hosts"] = tls.Hosts
				}

				if tls.SecretName != nil {
					tlsMap["secret_name"] = tls.SecretName
				}

				if tls.CertificateId != nil {
					tlsMap["certificate_id"] = tls.CertificateId
				}

				tlsList = append(tlsList, tlsMap)
			}

			ingressMap["tls"] = []interface{}{tlsList}
		}

		if gateway.Ingress.Rules != nil {
			rulesList := []interface{}{}
			for _, rules := range gateway.Ingress.Rules {
				rulesMap := map[string]interface{}{}

				if rules.Host != nil {
					rulesMap["host"] = rules.Host
				}

				if rules.Protocol != nil {
					rulesMap["protocol"] = rules.Protocol
				}

				if rules.Http != nil {
					httpMap := map[string]interface{}{}

					if rules.Http.Paths != nil {
						pathsList := []interface{}{}
						for _, paths := range rules.Http.Paths {
							pathsMap := map[string]interface{}{}

							if paths.Path != nil {
								pathsMap["path"] = paths.Path
							}

							if paths.Backend != nil {
								backendMap := map[string]interface{}{}

								if paths.Backend.ServiceName != nil {
									backendMap["service_name"] = paths.Backend.ServiceName
								}

								if paths.Backend.ServicePort != nil {
									backendMap["service_port"] = paths.Backend.ServicePort
								}

								pathsMap["backend"] = []interface{}{backendMap}
							}

							pathsList = append(pathsList, pathsMap)
						}

						httpMap["paths"] = []interface{}{pathsList}
					}

					rulesMap["http"] = []interface{}{httpMap}
				}

				rulesList = append(rulesList, rulesMap)
			}

			ingressMap["rules"] = []interface{}{rulesList}
		}

		if gateway.Ingress.Vip != nil {
			ingressMap["vip"] = gateway.Ingress.Vip
		}

		if gateway.Ingress.ClbId != nil {
			ingressMap["clb_id"] = gateway.Ingress.ClbId
		}

		if gateway.Ingress.CreateTime != nil {
			ingressMap["create_time"] = gateway.Ingress.CreateTime
		}

		_ = d.Set("ingress", []interface{}{ingressMap})
	}

	return nil
}

func resourceTencentCloudTemGatewayUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tem_gateway.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := tem.NewModifyIngressRequest()

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	environmentId := idSplit[0]
	ingressName := idSplit[1]

	request.EnvironmentId = &environmentId
	request.IngressName = &ingressName

	immutableArgs := []string{"ingress"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	if d.HasChange("ingress") {
		if dMap, ok := helper.InterfacesHeadMap(d, "ingress"); ok {
			ingressInfo := tem.IngressInfo{}
			if v, ok := dMap["ingress_name"]; ok {
				ingressInfo.IngressName = helper.String(v.(string))
			}
			if v, ok := dMap["environment_id"]; ok {
				ingressInfo.EnvironmentId = helper.String(v.(string))
			}
			if v, ok := dMap["cluster_namespace"]; ok {
				ingressInfo.ClusterNamespace = helper.String(v.(string))
			}
			if v, ok := dMap["address_ip_version"]; ok {
				ingressInfo.AddressIpVersion = helper.String(v.(string))
			}
			if v, ok := dMap["rewrite_type"]; ok {
				ingressInfo.RewriteType = helper.String(v.(string))
			}
			if v, ok := dMap["mixed"]; ok {
				ingressInfo.Mixed = helper.Bool(v.(bool))
			}
			if v, ok := dMap["tls"]; ok {
				for _, item := range v.([]interface{}) {
					tlsMap := item.(map[string]interface{})
					ingressTls := tem.IngressTls{}
					if v, ok := tlsMap["hosts"]; ok {
						hostsSet := v.(*schema.Set).List()
						for i := range hostsSet {
							hosts := hostsSet[i].(string)
							ingressTls.Hosts = append(ingressTls.Hosts, &hosts)
						}
					}
					if v, ok := tlsMap["secret_name"]; ok {
						ingressTls.SecretName = helper.String(v.(string))
					}
					if v, ok := tlsMap["certificate_id"]; ok {
						ingressTls.CertificateId = helper.String(v.(string))
					}
					ingressInfo.Tls = append(ingressInfo.Tls, &ingressTls)
				}
			}
			if v, ok := dMap["rules"]; ok {
				for _, item := range v.([]interface{}) {
					rulesMap := item.(map[string]interface{})
					ingressRule := tem.IngressRule{}
					if v, ok := rulesMap["host"]; ok {
						ingressRule.Host = helper.String(v.(string))
					}
					if v, ok := rulesMap["protocol"]; ok {
						ingressRule.Protocol = helper.String(v.(string))
					}
					if httpMap, ok := helper.InterfaceToMap(rulesMap, "http"); ok {
						ingressRuleValue := tem.IngressRuleValue{}
						if v, ok := httpMap["paths"]; ok {
							for _, item := range v.([]interface{}) {
								pathsMap := item.(map[string]interface{})
								ingressRulePath := tem.IngressRulePath{}
								if v, ok := pathsMap["path"]; ok {
									ingressRulePath.Path = helper.String(v.(string))
								}
								if backendMap, ok := helper.InterfaceToMap(pathsMap, "backend"); ok {
									ingressRuleBackend := tem.IngressRuleBackend{}
									if v, ok := backendMap["service_name"]; ok {
										ingressRuleBackend.ServiceName = helper.String(v.(string))
									}
									if v, ok := backendMap["service_port"]; ok {
										ingressRuleBackend.ServicePort = helper.IntInt64(v.(int))
									}
									ingressRulePath.Backend = &ingressRuleBackend
								}
								ingressRuleValue.Paths = append(ingressRuleValue.Paths, &ingressRulePath)
							}
						}
						ingressRule.Http = &ingressRuleValue
					}
					ingressInfo.Rules = append(ingressInfo.Rules, &ingressRule)
				}
			}
			if v, ok := dMap["clb_id"]; ok {
				ingressInfo.ClbId = helper.String(v.(string))
			}
			request.Ingress = &ingressInfo
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTemClient().ModifyIngress(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update tem gateway failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudTemGatewayRead(d, meta)
}

func resourceTencentCloudTemGatewayDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tem_gateway.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TemService{client: meta.(*TencentCloudClient).apiV3Conn}
	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	environmentId := idSplit[0]
	ingressName := idSplit[1]

	if err := service.DeleteTemGatewayById(ctx, environmentId, ingressName); err != nil {
		return err
	}

	return nil
}
