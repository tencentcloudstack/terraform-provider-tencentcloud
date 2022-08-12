/*
Provides a resource to create a tem gateway

Example Usage

```hcl
resource "tencentcloud_tem_gateway" "gateway" {
  ingress {
    ingress_name = "demo"
    environment_id = "en-853mggjm"
    address_ip_version = "IPV4"
    rewrite_type = "NONE"
    mixed = false
    rules {
      host = "test.com"
      protocol = "http"
      http {
        paths {
          path = "/"
          backend {
            service_name = "demo"
            service_port = 80
          }
        }
      }
    }
    rules {
      host = "hello.com"
      protocol = "http"
      http {
        paths {
          path = "/"
          backend {
            service_name = "hello"
            service_port = 36000
          }
        }
      }
    }
  }
}

```
Import

tem gateway can be imported using the id, e.g.
```
$ terraform import tencentcloud_tem_gateway.gateway environmentId#gatewayName
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	tem "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tem/v20210701"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudTemGateway() *schema.Resource {
	return &schema.Resource{
		Read:   resourceTencentCloudTemGatewayRead,
		Create: resourceTencentCloudTemGatewayCreate,
		Update: resourceTencentCloudTemGatewayUpdate,
		Delete: resourceTencentCloudTemGatewayDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"ingress": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Description: "gateway properties.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ingress_name": {
							Type:        schema.TypeString,
							Required:    true,
							ForceNew:    true,
							Description: "gateway name.",
						},
						"environment_id": {
							Type:        schema.TypeString,
							Required:    true,
							ForceNew:    true,
							Description: "environment ID.",
						},
						"address_ip_version": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "ip version, support IPV4.",
						},
						"rewrite_type": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "redirect mode, support AUTO and NONE.",
						},
						"mixed": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "mixing HTTP and HTTPS.",
						},
						"tls": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "ingress TLS configurations.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"hosts": {
										Type: schema.TypeSet,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
										Required:    true,
										Description: "host names.",
									},
									"secret_name": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "secret name.",
									},
									"certificate_id": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "certificate ID.",
									},
								},
							},
						},
						"rules": {
							Type:        schema.TypeList,
							Required:    true,
							Description: "proxy rules.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"host": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "host name.",
									},
									"protocol": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "protocol.",
									},
									"http": {
										Type:        schema.TypeList,
										MaxItems:    1,
										Required:    true,
										Description: "rule payload.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"paths": {
													Type:        schema.TypeList,
													Required:    true,
													Description: "path payload.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"path": {
																Type:        schema.TypeString,
																Required:    true,
																Description: "path.",
															},
															"backend": {
																Type:        schema.TypeList,
																MaxItems:    1,
																Required:    true,
																Description: "backend payload.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"service_name": {
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "backend name.",
																		},
																		"service_port": {
																			Type:        schema.TypeInt,
																			Required:    true,
																			Description: "backend port.",
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
							Description: "gateway vip.",
						},
						"clb_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "related CLB ID.",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "creation time.",
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
		environmentId string
		ingressName   string
	)

	if dMap, ok := helper.InterfacesHeadMap(d, "ingress"); ok {
		ingressInfo := tem.IngressInfo{}
		if v, ok := dMap["ingress_name"]; ok {
			ingressName = v.(string)
			ingressInfo.IngressName = helper.String(v.(string))
		}
		if v, ok := dMap["environment_id"]; ok {
			environmentId = v.(string)
			ingressInfo.EnvironmentId = helper.String(v.(string))
		}

		ingressInfo.ClusterNamespace = helper.String("default")

		if v, ok := dMap["address_ip_version"]; ok {
			ingressInfo.AddressIPVersion = helper.String(v.(string))
		}
		if v, ok := dMap["rewrite_type"]; ok {
			ingressInfo.RewriteType = helper.String(v.(string))
		}
		if v, ok := dMap["mixed"]; ok {
			ingressInfo.Mixed = helper.Bool(v.(bool))
		}
		if v, ok := dMap["tls"]; ok {
			for _, item := range v.([]interface{}) {
				TlsMap := item.(map[string]interface{})
				ingressTls := tem.IngressTls{}
				if v, ok := TlsMap["hosts"]; ok {
					hostsSet := v.(*schema.Set).List()
					for i := range hostsSet {
						hosts := hostsSet[i].(string)
						ingressTls.Hosts = append(ingressTls.Hosts, &hosts)
					}
				}
				if v, ok := TlsMap["secret_name"]; ok {
					ingressTls.SecretName = helper.String(v.(string))
				}
				if v, ok := TlsMap["certificate_id"]; ok {
					ingressTls.CertificateId = helper.String(v.(string))
				}
				ingressInfo.Tls = append(ingressInfo.Tls, &ingressTls)
			}
		}
		if v, ok := dMap["rules"]; ok {
			for _, item := range v.([]interface{}) {
				RulesMap := item.(map[string]interface{})
				ingressRule := tem.IngressRule{}
				if v, ok := RulesMap["host"]; ok {
					ingressRule.Host = helper.String(v.(string))
				}
				if v, ok := RulesMap["protocol"]; ok {
					ingressRule.Protocol = helper.String(v.(string))
				}
				if HttpMap, ok := helper.InterfaceToMap(RulesMap, "http"); ok {
					ingressRuleValue := tem.IngressRuleValue{}
					if v, ok := HttpMap["paths"]; ok {
						for _, item := range v.([]interface{}) {
							PathsMap := item.(map[string]interface{})
							ingressRulePath := tem.IngressRulePath{}
							if v, ok := PathsMap["path"]; ok {
								ingressRulePath.Path = helper.String(v.(string))
							}
							if BackendMap, ok := helper.InterfaceToMap(PathsMap, "backend"); ok {
								ingressRuleBackend := tem.IngressRuleBackend{}
								if v, ok := BackendMap["service_name"]; ok {
									ingressRuleBackend.ServiceName = helper.String(v.(string))
								}
								if v, ok := BackendMap["service_port"]; ok {
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
		request.Ingress = &ingressInfo
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTemClient().ModifyIngress(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create tem gateway failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(environmentId + FILED_SP + ingressName)

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

	gateway, err := service.DescribeTemGateway(ctx, environmentId, ingressName)

	if err != nil {
		return err
	}

	if gateway == nil {
		d.SetId("")
		return fmt.Errorf("resource `gateway` %s does not exist", ingressName)
	}

	ingressMap := map[string]interface{}{}
	if gateway.IngressName != nil {
		ingressMap["ingress_name"] = gateway.IngressName
	}
	if gateway.EnvironmentId != nil {
		ingressMap["environment_id"] = gateway.EnvironmentId
	}
	if gateway.AddressIPVersion != nil {
		ingressMap["address_ip_version"] = gateway.AddressIPVersion
	}
	if gateway.RewriteType != nil {
		ingressMap["rewrite_type"] = gateway.RewriteType
	}
	if gateway.Mixed != nil {
		ingressMap["mixed"] = gateway.Mixed
	}
	if gateway.Tls != nil {
		tlsList := []interface{}{}
		for _, tls := range gateway.Tls {
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
		ingressMap["tls"] = tlsList
	}
	if gateway.Rules != nil {
		rulesList := []interface{}{}
		for _, rules := range gateway.Rules {
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
					httpMap["paths"] = pathsList
				}

				rulesMap["http"] = []interface{}{httpMap}
			}

			rulesList = append(rulesList, rulesMap)
		}
		ingressMap["rules"] = rulesList
	}
	if gateway.Vip != nil {
		ingressMap["vip"] = gateway.Vip
	}
	if gateway.ClbId != nil {
		ingressMap["clb_id"] = gateway.ClbId
	}
	if gateway.CreateTime != nil {
		ingressMap["create_time"] = gateway.CreateTime
	}

	_ = d.Set("ingress", []interface{}{ingressMap})

	return nil
}

func resourceTencentCloudTemGatewayUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tem_gateway.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := tem.NewModifyIngressRequest()

	if d.HasChange("ingress") {
		if dMap, ok := helper.InterfacesHeadMap(d, "ingress"); ok {
			ingressInfo := tem.IngressInfo{}

			idSplit := strings.Split(d.Id(), FILED_SP)
			if len(idSplit) != 2 {
				return fmt.Errorf("id is broken,%s", d.Id())
			}
			environmentId := idSplit[0]
			ingressName := idSplit[1]
			ingressInfo.IngressName = helper.String(ingressName)
			ingressInfo.EnvironmentId = helper.String(environmentId)
			ingressInfo.ClusterNamespace = helper.String("default")

			if v, ok := dMap["address_ip_version"]; ok {
				ingressInfo.AddressIPVersion = helper.String(v.(string))
			}
			if v, ok := dMap["rewrite_type"]; ok {
				ingressInfo.RewriteType = helper.String(v.(string))
			}
			if v, ok := dMap["mixed"]; ok {
				ingressInfo.Mixed = helper.Bool(v.(bool))
			}
			if v, ok := dMap["tls"]; ok {
				for _, item := range v.([]interface{}) {
					TlsMap := item.(map[string]interface{})
					ingressTls := tem.IngressTls{}
					if v, ok := TlsMap["hosts"]; ok {
						hostsSet := v.(*schema.Set).List()
						for i := range hostsSet {
							hosts := hostsSet[i].(string)
							ingressTls.Hosts = append(ingressTls.Hosts, &hosts)
						}
					}
					if v, ok := TlsMap["secret_name"]; ok {
						ingressTls.SecretName = helper.String(v.(string))
					}
					if v, ok := TlsMap["certificate_id"]; ok {
						ingressTls.CertificateId = helper.String(v.(string))
					}
					ingressInfo.Tls = append(ingressInfo.Tls, &ingressTls)
				}
			}
			if v, ok := dMap["rules"]; ok {
				for _, item := range v.([]interface{}) {
					RulesMap := item.(map[string]interface{})
					ingressRule := tem.IngressRule{}
					if v, ok := RulesMap["host"]; ok {
						ingressRule.Host = helper.String(v.(string))
					}
					if v, ok := RulesMap["protocol"]; ok {
						ingressRule.Protocol = helper.String(v.(string))
					}
					if HttpMap, ok := helper.InterfaceToMap(RulesMap, "http"); ok {
						ingressRuleValue := tem.IngressRuleValue{}
						if v, ok := HttpMap["paths"]; ok {
							for _, item := range v.([]interface{}) {
								PathsMap := item.(map[string]interface{})
								ingressRulePath := tem.IngressRulePath{}
								if v, ok := PathsMap["path"]; ok {
									ingressRulePath.Path = helper.String(v.(string))
								}
								if BackendMap, ok := helper.InterfaceToMap(PathsMap, "backend"); ok {
									ingressRuleBackend := tem.IngressRuleBackend{}
									if v, ok := BackendMap["service_name"]; ok {
										ingressRuleBackend.ServiceName = helper.String(v.(string))
									}
									if v, ok := BackendMap["service_port"]; ok {
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
			request.Ingress = &ingressInfo
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTemClient().ModifyIngress(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})

	if err != nil {
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
