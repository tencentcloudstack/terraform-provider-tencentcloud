/*
Provides a resource to create a waf clb_domain

Example Usage

```hcl
resource "tencentcloud_waf_clb_domain" "clb_domain" {
  host {
		domain = ""
		domain_id = ""
		main_domain = ""
		mode =
		status =
		state =
		engine =
		is_cdn =
		load_balancer_set {
			load_balancer_id = ""
			load_balancer_name = ""
			listener_id = ""
			listener_name = ""
			vip = ""
			vport =
			region = ""
			protocol = ""
			zone = ""
			numerical_vpc_id =
			load_balancer_type = ""
		}
		region = ""
		edition = ""
		flow_mode =
		cls_status =
		level =
		cdc_clusters =
		alb_type = ""
		ip_headers =
		engine_type =

  }
  instance_i_d = ""
}
```

Import

waf clb_domain can be imported using the id, e.g.

```
terraform import tencentcloud_waf_clb_domain.clb_domain clb_domain_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	waf "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/waf/v20180125"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
)

func resourceTencentCloudWafClbDomain() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudWafClbDomainCreate,
		Read:   resourceTencentCloudWafClbDomainRead,
		Update: resourceTencentCloudWafClbDomainUpdate,
		Delete: resourceTencentCloudWafClbDomainDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"host": {
				Required:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Configuration of domain.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"domain": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Domain name.",
						},
						"domain_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Domain ID.",
						},
						"main_domain": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Primary domain name, empty when used as input parameters.",
						},
						"mode": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Rule defense mode, 0 observation mode, 1 interception mode.",
						},
						"status": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "WAF switch，0 off，1 on.",
						},
						"state": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Domain name status, 0 normal, 1 no traffic detected, 2 about to expire, 3 expired, use as output parameter.",
						},
						"engine": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Rule and AI Defense Mode, 10 Rule Engine Observation&amp;amp;&amp;amp;AI Engine Shutdown Mode 11 Rule Engine Observation&amp;amp;&amp;amp;AI Engine Observation Mode 12 Rule Engine Observation&amp;amp;&amp;amp;AI Engine Interception Mode 20 Rule Engine Interception&amp;amp;&amp;amp;AI Engine Shutdown Mode 21 Rule Engine Interception&amp;amp;&amp;amp;AI Engine Observation Mode 22 Rule Engine Interception&amp;amp;&amp;amp;AI Engine Interception Mode.",
						},
						"is_cdn": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Whether a proxy has been enabled before WAF, 0 no deployment, 1 deployment and use first IP in X-Forwarded-For as client IP, 2 deployment and use remote_addr as client IP, 3 deployment and use values of custom headers as client IP.",
						},
						"load_balancer_set": {
							Type:        schema.TypeList,
							Required:    true,
							Description: "List of bound LB.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"load_balancer_id": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "LoadBalancer ID.",
									},
									"load_balancer_name": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "LoadBalancer name.",
									},
									"listener_id": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Unique ID of listener in LB.",
									},
									"listener_name": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Listener name.",
									},
									"vip": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "LoadBalancer IP.",
									},
									"vport": {
										Type:        schema.TypeInt,
										Required:    true,
										Description: "LoadBalancer port.",
									},
									"region": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "LoadBalancer region.",
									},
									"protocol": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Protocol of listener，http or https.",
									},
									"zone": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "LoadBalancer zone.",
									},
									"numerical_vpc_id": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "VPC ID for load balancer, public network is -1, and internal network is filled in according to actual conditionsNote: This field may return null, indicating that a valid value cannot be obtained.",
									},
									"load_balancer_type": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Network type for load balancerNote: This field may return null, indicating that a valid value cannot be obtained.",
									},
								},
							},
						},
						"region": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Regions of LB bound by domain, and multiple regions are divided by ,.",
						},
						"edition": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "WAF edition, sparta-waf or clb-waf.",
						},
						"flow_mode": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "WAF traffic mode, 1 cleaning mode, 0 mirroring mode.",
						},
						"cls_status": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Whether to enable access logs, 1 enable, 0 disable, use as output parameter.",
						},
						"level": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Instance level, use as output parameterNote: This field may return null, indicating that a valid value cannot be obtained.",
						},
						"cdc_clusters": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Optional:    true,
							Description: "List of CDC clusters to which domain names need to be distributed, use as output parameterNote: This field may return null, indicating that a valid value cannot be obtained.。.",
						},
						"alb_type": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Traffic Source: clb represents Tencent Cloud clb, apisix represents apisix gateway, tsegw represents Tencent Cloud API gateway, default clbNote: This field may return null, indicating that a valid value cannot be obtained.",
						},
						"ip_headers": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Optional:    true,
							Description: "When IsCdn=3, this parameter needs to be filled in to indicate custom headersNote: This field may return null, indicating that a valid value cannot be obtained.",
						},
						"engine_type": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Rule engine type, 1: menshen, 2: tigaNote: This field may return null, indicating that a valid value cannot be obtained.",
						},
					},
				},
			},

			"instance_i_d": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Instance unique ID.",
			},
		},
	}
}

func resourceTencentCloudWafClbDomainCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_waf_clb_domain.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request  = waf.NewCreateHostRequest()
		response = waf.NewCreateHostResponse()
		domainId string
	)
	if dMap, ok := helper.InterfacesHeadMap(d, "host"); ok {
		hostRecord := waf.HostRecord{}
		if v, ok := dMap["domain"]; ok {
			hostRecord.Domain = helper.String(v.(string))
		}
		if v, ok := dMap["domain_id"]; ok {
			hostRecord.DomainId = helper.String(v.(string))
		}
		if v, ok := dMap["main_domain"]; ok {
			hostRecord.MainDomain = helper.String(v.(string))
		}
		if v, ok := dMap["mode"]; ok {
			hostRecord.Mode = helper.IntUint64(v.(int))
		}
		if v, ok := dMap["status"]; ok {
			hostRecord.Status = helper.IntUint64(v.(int))
		}
		if v, ok := dMap["state"]; ok {
			hostRecord.State = helper.IntUint64(v.(int))
		}
		if v, ok := dMap["engine"]; ok {
			hostRecord.Engine = helper.IntUint64(v.(int))
		}
		if v, ok := dMap["is_cdn"]; ok {
			hostRecord.IsCdn = helper.IntUint64(v.(int))
		}
		if v, ok := dMap["load_balancer_set"]; ok {
			for _, item := range v.([]interface{}) {
				loadBalancerSetMap := item.(map[string]interface{})
				loadBalancer := waf.LoadBalancer{}
				if v, ok := loadBalancerSetMap["load_balancer_id"]; ok {
					loadBalancer.LoadBalancerId = helper.String(v.(string))
				}
				if v, ok := loadBalancerSetMap["load_balancer_name"]; ok {
					loadBalancer.LoadBalancerName = helper.String(v.(string))
				}
				if v, ok := loadBalancerSetMap["listener_id"]; ok {
					loadBalancer.ListenerId = helper.String(v.(string))
				}
				if v, ok := loadBalancerSetMap["listener_name"]; ok {
					loadBalancer.ListenerName = helper.String(v.(string))
				}
				if v, ok := loadBalancerSetMap["vip"]; ok {
					loadBalancer.Vip = helper.String(v.(string))
				}
				if v, ok := loadBalancerSetMap["vport"]; ok {
					loadBalancer.Vport = helper.IntUint64(v.(int))
				}
				if v, ok := loadBalancerSetMap["region"]; ok {
					loadBalancer.Region = helper.String(v.(string))
				}
				if v, ok := loadBalancerSetMap["protocol"]; ok {
					loadBalancer.Protocol = helper.String(v.(string))
				}
				if v, ok := loadBalancerSetMap["zone"]; ok {
					loadBalancer.Zone = helper.String(v.(string))
				}
				if v, ok := loadBalancerSetMap["numerical_vpc_id"]; ok {
					loadBalancer.NumericalVpcId = helper.IntInt64(v.(int))
				}
				if v, ok := loadBalancerSetMap["load_balancer_type"]; ok {
					loadBalancer.LoadBalancerType = helper.String(v.(string))
				}
				hostRecord.LoadBalancerSet = append(hostRecord.LoadBalancerSet, &loadBalancer)
			}
		}
		if v, ok := dMap["region"]; ok {
			hostRecord.Region = helper.String(v.(string))
		}
		if v, ok := dMap["edition"]; ok {
			hostRecord.Edition = helper.String(v.(string))
		}
		if v, ok := dMap["flow_mode"]; ok {
			hostRecord.FlowMode = helper.IntUint64(v.(int))
		}
		if v, ok := dMap["cls_status"]; ok {
			hostRecord.ClsStatus = helper.IntUint64(v.(int))
		}
		if v, ok := dMap["level"]; ok {
			hostRecord.Level = helper.IntUint64(v.(int))
		}
		if v, ok := dMap["cdc_clusters"]; ok {
			cdcClustersSet := v.(*schema.Set).List()
			for i := range cdcClustersSet {
				cdcClusters := cdcClustersSet[i].(string)
				hostRecord.CdcClusters = append(hostRecord.CdcClusters, &cdcClusters)
			}
		}
		if v, ok := dMap["alb_type"]; ok {
			hostRecord.AlbType = helper.String(v.(string))
		}
		if v, ok := dMap["ip_headers"]; ok {
			ipHeadersSet := v.(*schema.Set).List()
			for i := range ipHeadersSet {
				ipHeaders := ipHeadersSet[i].(string)
				hostRecord.IpHeaders = append(hostRecord.IpHeaders, &ipHeaders)
			}
		}
		if v, ok := dMap["engine_type"]; ok {
			hostRecord.EngineType = helper.IntInt64(v.(int))
		}
		request.Host = &hostRecord
	}

	if v, ok := d.GetOk("instance_i_d"); ok {
		request.InstanceID = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseWafClient().CreateHost(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create waf clbDomain failed, reason:%+v", logId, err)
		return err
	}

	domainId = *response.Response.DomainId
	d.SetId(domainId)

	return resourceTencentCloudWafClbDomainRead(d, meta)
}

func resourceTencentCloudWafClbDomainRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_waf_clb_domain.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := WafService{client: meta.(*TencentCloudClient).apiV3Conn}

	clbDomainId := d.Id()

	clbDomain, err := service.DescribeWafClbDomainById(ctx, domainId)
	if err != nil {
		return err
	}

	if clbDomain == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `WafClbDomain` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if clbDomain.Host != nil {
		hostMap := map[string]interface{}{}

		if clbDomain.Host.Domain != nil {
			hostMap["domain"] = clbDomain.Host.Domain
		}

		if clbDomain.Host.DomainId != nil {
			hostMap["domain_id"] = clbDomain.Host.DomainId
		}

		if clbDomain.Host.MainDomain != nil {
			hostMap["main_domain"] = clbDomain.Host.MainDomain
		}

		if clbDomain.Host.Mode != nil {
			hostMap["mode"] = clbDomain.Host.Mode
		}

		if clbDomain.Host.Status != nil {
			hostMap["status"] = clbDomain.Host.Status
		}

		if clbDomain.Host.State != nil {
			hostMap["state"] = clbDomain.Host.State
		}

		if clbDomain.Host.Engine != nil {
			hostMap["engine"] = clbDomain.Host.Engine
		}

		if clbDomain.Host.IsCdn != nil {
			hostMap["is_cdn"] = clbDomain.Host.IsCdn
		}

		if clbDomain.Host.LoadBalancerSet != nil {
			loadBalancerSetList := []interface{}{}
			for _, loadBalancerSet := range clbDomain.Host.LoadBalancerSet {
				loadBalancerSetMap := map[string]interface{}{}

				if loadBalancerSet.LoadBalancerId != nil {
					loadBalancerSetMap["load_balancer_id"] = loadBalancerSet.LoadBalancerId
				}

				if loadBalancerSet.LoadBalancerName != nil {
					loadBalancerSetMap["load_balancer_name"] = loadBalancerSet.LoadBalancerName
				}

				if loadBalancerSet.ListenerId != nil {
					loadBalancerSetMap["listener_id"] = loadBalancerSet.ListenerId
				}

				if loadBalancerSet.ListenerName != nil {
					loadBalancerSetMap["listener_name"] = loadBalancerSet.ListenerName
				}

				if loadBalancerSet.Vip != nil {
					loadBalancerSetMap["vip"] = loadBalancerSet.Vip
				}

				if loadBalancerSet.Vport != nil {
					loadBalancerSetMap["vport"] = loadBalancerSet.Vport
				}

				if loadBalancerSet.Region != nil {
					loadBalancerSetMap["region"] = loadBalancerSet.Region
				}

				if loadBalancerSet.Protocol != nil {
					loadBalancerSetMap["protocol"] = loadBalancerSet.Protocol
				}

				if loadBalancerSet.Zone != nil {
					loadBalancerSetMap["zone"] = loadBalancerSet.Zone
				}

				if loadBalancerSet.NumericalVpcId != nil {
					loadBalancerSetMap["numerical_vpc_id"] = loadBalancerSet.NumericalVpcId
				}

				if loadBalancerSet.LoadBalancerType != nil {
					loadBalancerSetMap["load_balancer_type"] = loadBalancerSet.LoadBalancerType
				}

				loadBalancerSetList = append(loadBalancerSetList, loadBalancerSetMap)
			}

			hostMap["load_balancer_set"] = []interface{}{loadBalancerSetList}
		}

		if clbDomain.Host.Region != nil {
			hostMap["region"] = clbDomain.Host.Region
		}

		if clbDomain.Host.Edition != nil {
			hostMap["edition"] = clbDomain.Host.Edition
		}

		if clbDomain.Host.FlowMode != nil {
			hostMap["flow_mode"] = clbDomain.Host.FlowMode
		}

		if clbDomain.Host.ClsStatus != nil {
			hostMap["cls_status"] = clbDomain.Host.ClsStatus
		}

		if clbDomain.Host.Level != nil {
			hostMap["level"] = clbDomain.Host.Level
		}

		if clbDomain.Host.CdcClusters != nil {
			hostMap["cdc_clusters"] = clbDomain.Host.CdcClusters
		}

		if clbDomain.Host.AlbType != nil {
			hostMap["alb_type"] = clbDomain.Host.AlbType
		}

		if clbDomain.Host.IpHeaders != nil {
			hostMap["ip_headers"] = clbDomain.Host.IpHeaders
		}

		if clbDomain.Host.EngineType != nil {
			hostMap["engine_type"] = clbDomain.Host.EngineType
		}

		_ = d.Set("host", []interface{}{hostMap})
	}

	if clbDomain.InstanceID != nil {
		_ = d.Set("instance_i_d", clbDomain.InstanceID)
	}

	return nil
}

func resourceTencentCloudWafClbDomainUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_waf_clb_domain.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		modifyHostRequest  = waf.NewModifyHostRequest()
		modifyHostResponse = waf.NewModifyHostResponse()
	)

	clbDomainId := d.Id()

	request.DomainId = &domainId

	immutableArgs := []string{"host", "instance_i_d"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	if d.HasChange("host") {
		if dMap, ok := helper.InterfacesHeadMap(d, "host"); ok {
			hostRecord := waf.HostRecord{}
			if v, ok := dMap["domain"]; ok {
				hostRecord.Domain = helper.String(v.(string))
			}
			if v, ok := dMap["domain_id"]; ok {
				hostRecord.DomainId = helper.String(v.(string))
			}
			if v, ok := dMap["main_domain"]; ok {
				hostRecord.MainDomain = helper.String(v.(string))
			}
			if v, ok := dMap["mode"]; ok {
				hostRecord.Mode = helper.IntUint64(v.(int))
			}
			if v, ok := dMap["status"]; ok {
				hostRecord.Status = helper.IntUint64(v.(int))
			}
			if v, ok := dMap["state"]; ok {
				hostRecord.State = helper.IntUint64(v.(int))
			}
			if v, ok := dMap["engine"]; ok {
				hostRecord.Engine = helper.IntUint64(v.(int))
			}
			if v, ok := dMap["is_cdn"]; ok {
				hostRecord.IsCdn = helper.IntUint64(v.(int))
			}
			if v, ok := dMap["load_balancer_set"]; ok {
				for _, item := range v.([]interface{}) {
					loadBalancerSetMap := item.(map[string]interface{})
					loadBalancer := waf.LoadBalancer{}
					if v, ok := loadBalancerSetMap["load_balancer_id"]; ok {
						loadBalancer.LoadBalancerId = helper.String(v.(string))
					}
					if v, ok := loadBalancerSetMap["load_balancer_name"]; ok {
						loadBalancer.LoadBalancerName = helper.String(v.(string))
					}
					if v, ok := loadBalancerSetMap["listener_id"]; ok {
						loadBalancer.ListenerId = helper.String(v.(string))
					}
					if v, ok := loadBalancerSetMap["listener_name"]; ok {
						loadBalancer.ListenerName = helper.String(v.(string))
					}
					if v, ok := loadBalancerSetMap["vip"]; ok {
						loadBalancer.Vip = helper.String(v.(string))
					}
					if v, ok := loadBalancerSetMap["vport"]; ok {
						loadBalancer.Vport = helper.IntUint64(v.(int))
					}
					if v, ok := loadBalancerSetMap["region"]; ok {
						loadBalancer.Region = helper.String(v.(string))
					}
					if v, ok := loadBalancerSetMap["protocol"]; ok {
						loadBalancer.Protocol = helper.String(v.(string))
					}
					if v, ok := loadBalancerSetMap["zone"]; ok {
						loadBalancer.Zone = helper.String(v.(string))
					}
					if v, ok := loadBalancerSetMap["numerical_vpc_id"]; ok {
						loadBalancer.NumericalVpcId = helper.IntInt64(v.(int))
					}
					if v, ok := loadBalancerSetMap["load_balancer_type"]; ok {
						loadBalancer.LoadBalancerType = helper.String(v.(string))
					}
					hostRecord.LoadBalancerSet = append(hostRecord.LoadBalancerSet, &loadBalancer)
				}
			}
			if v, ok := dMap["region"]; ok {
				hostRecord.Region = helper.String(v.(string))
			}
			if v, ok := dMap["edition"]; ok {
				hostRecord.Edition = helper.String(v.(string))
			}
			if v, ok := dMap["flow_mode"]; ok {
				hostRecord.FlowMode = helper.IntUint64(v.(int))
			}
			if v, ok := dMap["cls_status"]; ok {
				hostRecord.ClsStatus = helper.IntUint64(v.(int))
			}
			if v, ok := dMap["level"]; ok {
				hostRecord.Level = helper.IntUint64(v.(int))
			}
			if v, ok := dMap["cdc_clusters"]; ok {
				cdcClustersSet := v.(*schema.Set).List()
				for i := range cdcClustersSet {
					cdcClusters := cdcClustersSet[i].(string)
					hostRecord.CdcClusters = append(hostRecord.CdcClusters, &cdcClusters)
				}
			}
			if v, ok := dMap["alb_type"]; ok {
				hostRecord.AlbType = helper.String(v.(string))
			}
			if v, ok := dMap["ip_headers"]; ok {
				ipHeadersSet := v.(*schema.Set).List()
				for i := range ipHeadersSet {
					ipHeaders := ipHeadersSet[i].(string)
					hostRecord.IpHeaders = append(hostRecord.IpHeaders, &ipHeaders)
				}
			}
			if v, ok := dMap["engine_type"]; ok {
				hostRecord.EngineType = helper.IntInt64(v.(int))
			}
			request.Host = &hostRecord
		}
	}

	if d.HasChange("instance_i_d") {
		if v, ok := d.GetOk("instance_i_d"); ok {
			request.InstanceID = helper.String(v.(string))
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseWafClient().ModifyHost(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update waf clbDomain failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudWafClbDomainRead(d, meta)
}

func resourceTencentCloudWafClbDomainDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_waf_clb_domain.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := WafService{client: meta.(*TencentCloudClient).apiV3Conn}
	clbDomainId := d.Id()

	if err := service.DeleteWafClbDomainById(ctx, domainId); err != nil {
		return err
	}

	return nil
}
