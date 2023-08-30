/*
Provides a resource to create a tse cngw_service

Example Usage

```hcl
variable "availability_zone" {
  default = "ap-guangzhou-4"
}

resource "tencentcloud_vpc" "vpc" {
  cidr_block = "10.0.0.0/16"
  name       = "tf_tse_vpc"
}

resource "tencentcloud_subnet" "subnet" {
  vpc_id            = tencentcloud_vpc.vpc.id
  availability_zone = var.availability_zone
  name              = "tf_tse_subnet"
  cidr_block        = "10.0.1.0/24"
}

resource "tencentcloud_tse_cngw_gateway" "cngw_gateway" {
  description                = "terraform test1"
  enable_cls                 = true
  engine_region              = "ap-guangzhou"
  feature_version            = "STANDARD"
  gateway_version            = "2.5.1"
  ingress_class_name         = "tse-nginx-ingress"
  internet_max_bandwidth_out = 0
  name                       = "terraform-gateway1"
  trade_type                 = 0
  type                       = "kong"

  node_config {
    number        = 2
    specification = "1c2g"
  }

  vpc_config {
    subnet_id = tencentcloud_subnet.subnet.id
    vpc_id    = tencentcloud_vpc.vpc.id
  }

  tags = {
    "createdBy" = "terraform"
  }
}

resource "tencentcloud_tse_cngw_service" "cngw_service" {
  gateway_id    = tencentcloud_tse_cngw_gateway.cngw_gateway.id
  name          = "terraform-test"
  path          = "/test"
  protocol      = "http"
  retries       = 5
  timeout       = 60000
  upstream_type = "HostIP"

  upstream_info {
    algorithm             = "round-robin"
    auto_scaling_cvm_port = 0
    host                  = "arunma.cn"
    port                  = 8012
    slow_start            = 0
  }
}
```

Import

tse cngw_service can be imported using the id, e.g.

```
terraform import tencentcloud_tse_cngw_service.cngw_service gatewayId#name
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tse "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tse/v20201207"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudTseCngwService() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTseCngwServiceCreate,
		Read:   resourceTencentCloudTseCngwServiceRead,
		Update: resourceTencentCloudTseCngwServiceUpdate,
		Delete: resourceTencentCloudTseCngwServiceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"gateway_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "gateway ID.",
			},

			"name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "service name.",
			},

			"protocol": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "protocol. Reference value:`https`, `http`, `tcp`, `udp`.",
			},

			"path": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "path.",
			},

			"timeout": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "time out, unit:ms.",
			},

			"retries": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "retry times.",
			},

			"upstream_type": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "service type. Reference value:`Kubernetes`, `Registry`, `IPList`, `HostIP`, `Scf`.",
			},

			"upstream_info": {
				Required:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "service config information.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"host": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "an IP address or domain name.",
						},
						"port": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "backend service port.valid values: `1` to `65535`.",
						},
						"source_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "service source ID.",
						},
						"namespace": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "namespace.",
						},
						"service_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "the name of the service in registry or kubernetes.",
						},
						"targets": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "provided when service type is IPList.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"host": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "host.",
									},
									"port": {
										Type:        schema.TypeInt,
										Required:    true,
										Description: "port.",
									},
									"weight": {
										Type:        schema.TypeInt,
										Required:    true,
										Description: "weight.",
									},
									"health": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "health.",
									},
									"created_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "created time.",
									},
									"source": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "source of target.",
									},
								},
							},
						},
						"source_type": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "source service type.",
						},
						"scf_type": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "scf lambda type.",
						},
						"scf_namespace": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "scf lambda namespace.",
						},
						"scf_lambda_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "scf lambda name.",
						},
						"scf_lambda_qualifier": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "scf lambda version.",
						},
						"slow_start": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "slow start time, unit: `second`, when it is enabled, weight of the node is increased from 1 to the target value gradually.",
						},
						"algorithm": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "load balance algorithm,default: `round-robin`, `least-connections` and `consisten_hashing` also support.",
						},
						"auto_scaling_group_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "auto scaling group ID of cvm.",
						},
						"auto_scaling_cvm_port": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "auto scaling group port of cvm.",
						},
						"auto_scaling_tat_cmd_status": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "tat cmd status in auto scaling group of cvm.",
						},
						"auto_scaling_hook_status": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "hook status in auto scaling group of cvm.",
						},
						"source_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "the name of source service.",
						},
						"real_source_type": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "exact source service type.",
						},
					},
				},
			},

			"service_id": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "service id.",
			},
		},
	}
}

func resourceTencentCloudTseCngwServiceCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tse_cngw_service.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request   = tse.NewCreateCloudNativeAPIGatewayServiceRequest()
		gatewayId string
		name      string
	)
	if v, ok := d.GetOk("gateway_id"); ok {
		gatewayId = v.(string)
		request.GatewayId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("name"); ok {
		name = v.(string)
		request.Name = helper.String(v.(string))
	}

	if v, ok := d.GetOk("protocol"); ok {
		request.Protocol = helper.String(v.(string))
	}

	if v, ok := d.GetOk("path"); ok {
		request.Path = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("timeout"); ok {
		request.Timeout = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("retries"); ok {
		request.Retries = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("upstream_type"); ok {
		request.UpstreamType = helper.String(v.(string))
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "upstream_info"); ok {
		kongUpstreamInfo := tse.KongUpstreamInfo{}
		if v, ok := dMap["host"]; ok {
			kongUpstreamInfo.Host = helper.String(v.(string))
		}
		if v, ok := dMap["port"]; ok {
			kongUpstreamInfo.Port = helper.IntInt64(v.(int))
		}
		if v, ok := dMap["source_id"]; ok {
			kongUpstreamInfo.SourceID = helper.String(v.(string))
		}
		if v, ok := dMap["namespace"]; ok {
			kongUpstreamInfo.Namespace = helper.String(v.(string))
		}
		if v, ok := dMap["service_name"]; ok {
			kongUpstreamInfo.ServiceName = helper.String(v.(string))
		}
		if v, ok := dMap["targets"]; ok {
			for _, item := range v.([]interface{}) {
				targetsMap := item.(map[string]interface{})
				kongTarget := tse.KongTarget{}
				if v, ok := targetsMap["host"]; ok {
					kongTarget.Host = helper.String(v.(string))
				}
				if v, ok := targetsMap["port"]; ok {
					kongTarget.Port = helper.IntInt64(v.(int))
				}
				if v, ok := targetsMap["weight"]; ok {
					kongTarget.Weight = helper.IntInt64(v.(int))
				}
				if v, ok := targetsMap["source"]; ok {
					kongTarget.Source = helper.String(v.(string))
				}
				kongUpstreamInfo.Targets = append(kongUpstreamInfo.Targets, &kongTarget)
			}
		}
		if v, ok := dMap["source_type"]; ok {
			kongUpstreamInfo.SourceType = helper.String(v.(string))
		}
		if v, ok := dMap["scf_type"]; ok {
			kongUpstreamInfo.ScfType = helper.String(v.(string))
		}
		if v, ok := dMap["scf_namespace"]; ok {
			kongUpstreamInfo.ScfNamespace = helper.String(v.(string))
		}
		if v, ok := dMap["scf_lambda_name"]; ok {
			kongUpstreamInfo.ScfLambdaName = helper.String(v.(string))
		}
		if v, ok := dMap["scf_lambda_qualifier"]; ok {
			kongUpstreamInfo.ScfLambdaQualifier = helper.String(v.(string))
		}
		if v, ok := dMap["slow_start"]; ok {
			kongUpstreamInfo.SlowStart = helper.IntInt64(v.(int))
		}
		if v, ok := dMap["algorithm"]; ok {
			kongUpstreamInfo.Algorithm = helper.String(v.(string))
		}
		if v, ok := dMap["auto_scaling_group_id"]; ok {
			kongUpstreamInfo.AutoScalingGroupID = helper.String(v.(string))
		}
		if v, ok := dMap["auto_scaling_cvm_port"]; ok {
			kongUpstreamInfo.AutoScalingCvmPort = helper.IntUint64(v.(int))
		}
		if v, ok := dMap["auto_scaling_tat_cmd_status"]; ok {
			kongUpstreamInfo.AutoScalingTatCmdStatus = helper.String(v.(string))
		}
		if v, ok := dMap["auto_scaling_hook_status"]; ok {
			kongUpstreamInfo.AutoScalingHookStatus = helper.String(v.(string))
		}
		if v, ok := dMap["source_name"]; ok {
			kongUpstreamInfo.SourceName = helper.String(v.(string))
		}
		if v, ok := dMap["real_source_type"]; ok {
			kongUpstreamInfo.RealSourceType = helper.String(v.(string))
		}
		request.UpstreamInfo = &kongUpstreamInfo
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTseClient().CreateCloudNativeAPIGatewayService(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create tse cngwService failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(gatewayId + FILED_SP + name)

	return resourceTencentCloudTseCngwServiceRead(d, meta)
}

func resourceTencentCloudTseCngwServiceRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tse_cngw_service.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TseService{client: meta.(*TencentCloudClient).apiV3Conn}

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	gatewayId := idSplit[0]
	name := idSplit[1]
	cngwService, err := service.DescribeTseCngwServiceById(ctx, gatewayId, name)
	if err != nil {
		return err
	}

	if cngwService == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `TseCngwService` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	_ = d.Set("gateway_id", gatewayId)

	if cngwService.Name != nil {
		_ = d.Set("name", cngwService.Name)
	}

	if cngwService.Protocol != nil {
		_ = d.Set("protocol", cngwService.Protocol)
	}

	if cngwService.Path != nil {
		_ = d.Set("path", cngwService.Path)
	}

	if cngwService.Timeout != nil {
		_ = d.Set("timeout", cngwService.Timeout)
	}

	if cngwService.Retries != nil {
		_ = d.Set("retries", cngwService.Retries)
	}

	if cngwService.UpstreamType != nil {
		_ = d.Set("upstream_type", cngwService.UpstreamType)
	}

	if cngwService.UpstreamInfo != nil {
		upstreamInfoMap := map[string]interface{}{}

		if cngwService.UpstreamInfo.Host != nil {
			upstreamInfoMap["host"] = cngwService.UpstreamInfo.Host
		}

		if cngwService.UpstreamInfo.Port != nil {
			upstreamInfoMap["port"] = cngwService.UpstreamInfo.Port
		}

		if cngwService.UpstreamInfo.SourceID != nil {
			upstreamInfoMap["source_id"] = cngwService.UpstreamInfo.SourceID
		}

		if cngwService.UpstreamInfo.Namespace != nil {
			upstreamInfoMap["namespace"] = cngwService.UpstreamInfo.Namespace
		}

		if cngwService.UpstreamInfo.ServiceName != nil {
			upstreamInfoMap["service_name"] = cngwService.UpstreamInfo.ServiceName
		}

		if cngwService.UpstreamInfo.Targets != nil {
			targetsList := []interface{}{}
			for _, targets := range cngwService.UpstreamInfo.Targets {
				targetsMap := map[string]interface{}{}

				if targets.Host != nil {
					targetsMap["host"] = targets.Host
				}

				if targets.Port != nil {
					targetsMap["port"] = targets.Port
				}

				if targets.Weight != nil {
					targetsMap["weight"] = targets.Weight
				}

				if targets.Health != nil {
					targetsMap["health"] = targets.Health
				}

				if targets.CreatedTime != nil {
					targetsMap["created_time"] = targets.CreatedTime
				}

				if targets.Source != nil {
					targetsMap["source"] = targets.Source
				}

				targetsList = append(targetsList, targetsMap)
			}

			upstreamInfoMap["targets"] = targetsList
		}

		if cngwService.UpstreamInfo.SourceType != nil {
			upstreamInfoMap["source_type"] = cngwService.UpstreamInfo.SourceType
		}

		if cngwService.UpstreamInfo.ScfType != nil {
			upstreamInfoMap["scf_type"] = cngwService.UpstreamInfo.ScfType
		}

		if cngwService.UpstreamInfo.ScfNamespace != nil {
			upstreamInfoMap["scf_namespace"] = cngwService.UpstreamInfo.ScfNamespace
		}

		if cngwService.UpstreamInfo.ScfLambdaName != nil {
			upstreamInfoMap["scf_lambda_name"] = cngwService.UpstreamInfo.ScfLambdaName
		}

		if cngwService.UpstreamInfo.ScfLambdaQualifier != nil {
			upstreamInfoMap["scf_lambda_qualifier"] = cngwService.UpstreamInfo.ScfLambdaQualifier
		}

		if cngwService.UpstreamInfo.SlowStart != nil {
			upstreamInfoMap["slow_start"] = cngwService.UpstreamInfo.SlowStart
		}

		if cngwService.UpstreamInfo.Algorithm != nil {
			upstreamInfoMap["algorithm"] = cngwService.UpstreamInfo.Algorithm
		}

		if cngwService.UpstreamInfo.AutoScalingGroupID != nil {
			upstreamInfoMap["auto_scaling_group_id"] = cngwService.UpstreamInfo.AutoScalingGroupID
		}

		if cngwService.UpstreamInfo.AutoScalingCvmPort != nil {
			upstreamInfoMap["auto_scaling_cvm_port"] = cngwService.UpstreamInfo.AutoScalingCvmPort
		}

		if cngwService.UpstreamInfo.AutoScalingTatCmdStatus != nil {
			upstreamInfoMap["auto_scaling_tat_cmd_status"] = cngwService.UpstreamInfo.AutoScalingTatCmdStatus
		}

		if cngwService.UpstreamInfo.AutoScalingHookStatus != nil {
			upstreamInfoMap["auto_scaling_hook_status"] = cngwService.UpstreamInfo.AutoScalingHookStatus
		}

		if cngwService.UpstreamInfo.SourceName != nil {
			upstreamInfoMap["source_name"] = cngwService.UpstreamInfo.SourceName
		}

		if cngwService.UpstreamInfo.RealSourceType != nil {
			upstreamInfoMap["real_source_type"] = cngwService.UpstreamInfo.RealSourceType
		}

		_ = d.Set("upstream_info", []interface{}{upstreamInfoMap})
	}

	if cngwService.ID != nil {
		_ = d.Set("service_id", cngwService.ID)
	}

	return nil
}

func resourceTencentCloudTseCngwServiceUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tse_cngw_service.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	request := tse.NewModifyCloudNativeAPIGatewayServiceRequest()

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	gatewayId := idSplit[0]
	name := idSplit[1]
	request.GatewayId = &gatewayId
	request.Name = &name

	service := TseService{client: meta.(*TencentCloudClient).apiV3Conn}
	cngwService, err := service.DescribeTseCngwServiceById(ctx, gatewayId, name)
	if err != nil {
		return err
	}

	serviceId := ""
	if cngwService.ID != nil {
		serviceId = *cngwService.ID
	}
	request.ID = &serviceId

	if cngwService == nil {
		return fmt.Errorf("[WARN]%s resource `TseCngwService` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
	}

	immutableArgs := []string{"gateway_id", "name"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	if v, ok := d.GetOk("protocol"); ok {
		request.Protocol = helper.String(v.(string))
	}

	if v, ok := d.GetOk("path"); ok {
		request.Path = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("timeout"); ok {
		request.Timeout = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("retries"); ok {
		request.Retries = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("upstream_type"); ok {
		request.UpstreamType = helper.String(v.(string))
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "upstream_info"); ok {
		kongUpstreamInfo := tse.KongUpstreamInfo{}
		if v, ok := dMap["host"]; ok {
			kongUpstreamInfo.Host = helper.String(v.(string))
		}
		if v, ok := dMap["port"]; ok {
			kongUpstreamInfo.Port = helper.IntInt64(v.(int))
		}
		if v, ok := dMap["source_id"]; ok {
			kongUpstreamInfo.SourceID = helper.String(v.(string))
		}
		if v, ok := dMap["namespace"]; ok {
			kongUpstreamInfo.Namespace = helper.String(v.(string))
		}
		if v, ok := dMap["service_name"]; ok {
			kongUpstreamInfo.ServiceName = helper.String(v.(string))
		}
		if v, ok := dMap["targets"]; ok {
			for _, item := range v.([]interface{}) {
				targetsMap := item.(map[string]interface{})
				kongTarget := tse.KongTarget{}
				if v, ok := targetsMap["host"]; ok {
					kongTarget.Host = helper.String(v.(string))
				}
				if v, ok := targetsMap["port"]; ok {
					kongTarget.Port = helper.IntInt64(v.(int))
				}
				if v, ok := targetsMap["weight"]; ok {
					kongTarget.Weight = helper.IntInt64(v.(int))
				}
				if v, ok := targetsMap["source"]; ok {
					kongTarget.Source = helper.String(v.(string))
				}
				kongUpstreamInfo.Targets = append(kongUpstreamInfo.Targets, &kongTarget)
			}
		}
		if v, ok := dMap["source_type"]; ok {
			kongUpstreamInfo.SourceType = helper.String(v.(string))
		}
		if v, ok := dMap["scf_type"]; ok {
			kongUpstreamInfo.ScfType = helper.String(v.(string))
		}
		if v, ok := dMap["scf_namespace"]; ok {
			kongUpstreamInfo.ScfNamespace = helper.String(v.(string))
		}
		if v, ok := dMap["scf_lambda_name"]; ok {
			kongUpstreamInfo.ScfLambdaName = helper.String(v.(string))
		}
		if v, ok := dMap["scf_lambda_qualifier"]; ok {
			kongUpstreamInfo.ScfLambdaQualifier = helper.String(v.(string))
		}
		if v, ok := dMap["slow_start"]; ok {
			kongUpstreamInfo.SlowStart = helper.IntInt64(v.(int))
		}
		if v, ok := dMap["algorithm"]; ok {
			kongUpstreamInfo.Algorithm = helper.String(v.(string))
		}
		if v, ok := dMap["auto_scaling_group_id"]; ok {
			kongUpstreamInfo.AutoScalingGroupID = helper.String(v.(string))
		}
		if v, ok := dMap["auto_scaling_cvm_port"]; ok {
			kongUpstreamInfo.AutoScalingCvmPort = helper.IntUint64(v.(int))
		}
		if v, ok := dMap["auto_scaling_tat_cmd_status"]; ok {
			kongUpstreamInfo.AutoScalingTatCmdStatus = helper.String(v.(string))
		}
		if v, ok := dMap["auto_scaling_hook_status"]; ok {
			kongUpstreamInfo.AutoScalingHookStatus = helper.String(v.(string))
		}
		if v, ok := dMap["source_name"]; ok {
			kongUpstreamInfo.SourceName = helper.String(v.(string))
		}
		if v, ok := dMap["real_source_type"]; ok {
			kongUpstreamInfo.RealSourceType = helper.String(v.(string))
		}
		request.UpstreamInfo = &kongUpstreamInfo
	}

	err = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTseClient().ModifyCloudNativeAPIGatewayService(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update tse cngwService failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudTseCngwServiceRead(d, meta)
}

func resourceTencentCloudTseCngwServiceDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tse_cngw_service.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TseService{client: meta.(*TencentCloudClient).apiV3Conn}
	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	gatewayId := idSplit[0]
	name := idSplit[1]
	if err := service.DeleteTseCngwServiceById(ctx, gatewayId, name); err != nil {
		return err
	}

	return nil
}
