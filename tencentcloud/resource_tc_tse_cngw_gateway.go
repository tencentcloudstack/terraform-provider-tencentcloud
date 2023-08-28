/*
Provides a resource to create a tse cngw_gateway

Example Usage

```hcl
resource "tencentcloud_tse_cngw_gateway" "cngw_gateway" {
  name = "test"
  type = "kong"
  gateway_version = "2.4.1"
  node_config {
		specification = ""
		number =

  }
  vpc_config {
		vpc_id = "vpc-xxxxxx"
		subnet_id = "subnet-xxxxxx"

  }
  description = "for test"
  enable_cls = false
  feature_version = ""
  internet_max_bandwidth_out =
  engine_region = "ap-guangzhou"
  ingress_class_name = ""
  trade_type =
  internet_config {
		internet_address_version = ""
		internet_pay_mode = ""
		internet_max_bandwidth_out =
		description = ""
		sla_type = ""
		multi_zone_flag =
		master_zone_id = ""
		slave_zone_id = ""

  }
  tags = {
    "createdBy" = "terraform"
  }
}
```

Import

tse cngw_gateway can be imported using the id, e.g.

```
terraform import tencentcloud_tse_cngw_gateway.cngw_gateway cngw_gateway_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tse "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tse/v20201207"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
)

func resourceTencentCloudTseCngwGateway() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTseCngwGatewayCreate,
		Read:   resourceTencentCloudTseCngwGatewayRead,
		Update: resourceTencentCloudTseCngwGatewayUpdate,
		Delete: resourceTencentCloudTseCngwGatewayDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "gateway name, supports up to 60 characters.",
			},

			"type": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "gateway type,currently only supports kong.",
			},

			"gateway_version": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "gateway vwersion. Reference value：- 2.4.1- 2.5.1.",
			},

			"node_config": {
				Required:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "gateway node configration.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"specification": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "specification, 1c2g|2c4g|4c8g|8c16g.",
						},
						"number": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "node number, 2-50.",
						},
					},
				},
			},

			"vpc_config": {
				Required:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "vpc infomation.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"vpc_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "VPC ID. Assign an IP address to the engine in the VPC subnet. Reference value: vpc-conz6aix.",
						},
						"subnet_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "subnet ID. Assign an IP address to the engine in the VPC subnet. Reference value: subnet-ahde9me9.",
						},
					},
				},
			},

			"description": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "description information, up to 120 characters.",
			},

			"enable_cls": {
				Optional:    true,
				Type:        schema.TypeBool,
				Description: "whether to enable CLS log. Default value: fasle.",
			},

			"feature_version": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "product version. Reference value: `TRIAL`, `STANDARD`(default value), `PROFESSIONAL`.",
			},

			"internet_max_bandwidth_out": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "public network outbound traffic bandwidth,[1,2048]Mbps.",
			},

			"engine_region": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "engine region of gateway.",
			},

			"ingress_class_name": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "ingress class name.",
			},

			"trade_type": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "trade type. Reference value: `0`: postpaid, `1`:Prepaid (Interface does not support the creation of prepaid instances yet).",
			},

			"internet_config": {
				Optional:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "internet configration.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"internet_address_version": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "internet type. Reference value: `IPV4`(default value), `IPV6`.",
						},
						"internet_pay_mode": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "trade type of internet. Reference value: `BANDWIDTH`, `TRAFFIC`(default value).",
						},
						"internet_max_bandwidth_out": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "public network bandwidth.",
						},
						"description": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "description of clb.",
						},
						"sla_type": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "specification type of clb. Default shared type when this parameter is empty. Reference value:- SLA LCU-supported.",
						},
						"multi_zone_flag": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Whether load balancing has multiple availability zones.",
						},
						"master_zone_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "primary availability zone.",
						},
						"slave_zone_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "alternate availability zone.",
						},
					},
				},
			},

			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Tag description list.",
			},
		},
	}
}

func resourceTencentCloudTseCngwGatewayCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tse_cngw_gateway.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	var (
		request   = tse.NewCreateCloudNativeAPIGatewayRequest()
		response  = tse.NewCreateCloudNativeAPIGatewayResponse()
		gatewayId string
	)
	if v, ok := d.GetOk("name"); ok {
		request.Name = helper.String(v.(string))
	}

	if v, ok := d.GetOk("type"); ok {
		request.Type = helper.String(v.(string))
	}

	if v, ok := d.GetOk("gateway_version"); ok {
		request.GatewayVersion = helper.String(v.(string))
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "node_config"); ok {
		cloudNativeAPIGatewayNodeConfig := tse.CloudNativeAPIGatewayNodeConfig{}
		if v, ok := dMap["specification"]; ok {
			cloudNativeAPIGatewayNodeConfig.Specification = helper.String(v.(string))
		}
		if v, ok := dMap["number"]; ok {
			cloudNativeAPIGatewayNodeConfig.Number = helper.IntInt64(v.(int))
		}
		request.NodeConfig = &cloudNativeAPIGatewayNodeConfig
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "vpc_config"); ok {
		cloudNativeAPIGatewayVpcConfig := tse.CloudNativeAPIGatewayVpcConfig{}
		if v, ok := dMap["vpc_id"]; ok {
			cloudNativeAPIGatewayVpcConfig.VpcId = helper.String(v.(string))
		}
		if v, ok := dMap["subnet_id"]; ok {
			cloudNativeAPIGatewayVpcConfig.SubnetId = helper.String(v.(string))
		}
		request.VpcConfig = &cloudNativeAPIGatewayVpcConfig
	}

	if v, ok := d.GetOk("description"); ok {
		request.Description = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("enable_cls"); ok {
		request.EnableCls = helper.Bool(v.(bool))
	}

	if v, ok := d.GetOk("feature_version"); ok {
		request.FeatureVersion = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("internet_max_bandwidth_out"); ok {
		request.InternetMaxBandwidthOut = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("engine_region"); ok {
		request.EngineRegion = helper.String(v.(string))
	}

	if v, ok := d.GetOk("ingress_class_name"); ok {
		request.IngressClassName = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("trade_type"); ok {
		request.TradeType = helper.IntInt64(v.(int))
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "internet_config"); ok {
		internetConfig := tse.InternetConfig{}
		if v, ok := dMap["internet_address_version"]; ok {
			internetConfig.InternetAddressVersion = helper.String(v.(string))
		}
		if v, ok := dMap["internet_pay_mode"]; ok {
			internetConfig.InternetPayMode = helper.String(v.(string))
		}
		if v, ok := dMap["internet_max_bandwidth_out"]; ok {
			internetConfig.InternetMaxBandwidthOut = helper.IntUint64(v.(int))
		}
		if v, ok := dMap["description"]; ok {
			internetConfig.Description = helper.String(v.(string))
		}
		if v, ok := dMap["sla_type"]; ok {
			internetConfig.SlaType = helper.String(v.(string))
		}
		if v, ok := dMap["multi_zone_flag"]; ok {
			internetConfig.MultiZoneFlag = helper.Bool(v.(bool))
		}
		if v, ok := dMap["master_zone_id"]; ok {
			internetConfig.MasterZoneId = helper.String(v.(string))
		}
		if v, ok := dMap["slave_zone_id"]; ok {
			internetConfig.SlaveZoneId = helper.String(v.(string))
		}
		request.InternetConfig = &internetConfig
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTseClient().CreateCloudNativeAPIGateway(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create tse cngwGateway failed, reason:%+v", logId, err)
		return err
	}

	gatewayId = *response.Response.Result.GatewayId
	d.SetId(gatewayId)

	service := TseService{client: meta.(*TencentCloudClient).apiV3Conn}
	if err := service.CheckTseNativeAPIGatewayStatusById(ctx, gatewayId, "create"); err != nil {
		return err
	}

	if tags := helper.GetTags(d, "tags"); len(tags) > 0 {
		tagService := TagService{client: meta.(*TencentCloudClient).apiV3Conn}
		region := meta.(*TencentCloudClient).apiV3Conn.Region
		resourceName := fmt.Sprintf("qcs::tse:%s:uin/:gateway/%s", region, d.Id())
		if err := tagService.ModifyTags(ctx, resourceName, tags, nil); err != nil {
			return err
		}
	}

	return resourceTencentCloudTseCngwGatewayRead(d, meta)
}

func resourceTencentCloudTseCngwGatewayRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tse_cngw_gateway.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TseService{client: meta.(*TencentCloudClient).apiV3Conn}

	gatewayId := d.Id()

	cngwGateway, err := service.DescribeTseCngwGatewayById(ctx, gatewayId)
	if err != nil {
		return err
	}

	if cngwGateway == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `TseCngwGateway` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if cngwGateway.Name != nil {
		_ = d.Set("name", cngwGateway.Name)
	}

	if cngwGateway.Type != nil {
		_ = d.Set("type", cngwGateway.Type)
	}

	if cngwGateway.GatewayVersion != nil {
		_ = d.Set("gateway_version", cngwGateway.GatewayVersion)
	}

	if cngwGateway.NodeConfig != nil {
		nodeConfigMap := map[string]interface{}{}

		if cngwGateway.NodeConfig.Specification != nil {
			nodeConfigMap["specification"] = cngwGateway.NodeConfig.Specification
		}

		if cngwGateway.NodeConfig.Number != nil {
			nodeConfigMap["number"] = cngwGateway.NodeConfig.Number
		}

		_ = d.Set("node_config", []interface{}{nodeConfigMap})
	}

	if cngwGateway.VpcConfig != nil {
		vpcConfigMap := map[string]interface{}{}

		if cngwGateway.VpcConfig.VpcId != nil {
			vpcConfigMap["vpc_id"] = cngwGateway.VpcConfig.VpcId
		}

		if cngwGateway.VpcConfig.SubnetId != nil {
			vpcConfigMap["subnet_id"] = cngwGateway.VpcConfig.SubnetId
		}

		_ = d.Set("vpc_config", []interface{}{vpcConfigMap})
	}

	if cngwGateway.Description != nil {
		_ = d.Set("description", cngwGateway.Description)
	}

	if cngwGateway.EnableCls != nil {
		_ = d.Set("enable_cls", cngwGateway.EnableCls)
	}

	if cngwGateway.FeatureVersion != nil {
		_ = d.Set("feature_version", cngwGateway.FeatureVersion)
	}

	if cngwGateway.InternetMaxBandwidthOut != nil {
		_ = d.Set("internet_max_bandwidth_out", cngwGateway.InternetMaxBandwidthOut)
	}

	if cngwGateway.EngineRegion != nil {
		_ = d.Set("engine_region", cngwGateway.EngineRegion)
	}

	if cngwGateway.IngressClassName != nil {
		_ = d.Set("ingress_class_name", cngwGateway.IngressClassName)
	}

	if cngwGateway.TradeType != nil {
		_ = d.Set("trade_type", cngwGateway.TradeType)
	}

	// if cngwGateway.InternetConfig != nil {
	// 	internetConfigMap := map[string]interface{}{}

	// 	if cngwGateway.InternetConfig.InternetAddressVersion != nil {
	// 		internetConfigMap["internet_address_version"] = cngwGateway.InternetConfig.InternetAddressVersion
	// 	}

	// 	if cngwGateway.InternetConfig.InternetPayMode != nil {
	// 		internetConfigMap["internet_pay_mode"] = cngwGateway.InternetConfig.InternetPayMode
	// 	}

	// 	if cngwGateway.InternetConfig.InternetMaxBandwidthOut != nil {
	// 		internetConfigMap["internet_max_bandwidth_out"] = cngwGateway.InternetConfig.InternetMaxBandwidthOut
	// 	}

	// 	if cngwGateway.InternetConfig.Description != nil {
	// 		internetConfigMap["description"] = cngwGateway.InternetConfig.Description
	// 	}

	// 	if cngwGateway.InternetConfig.SlaType != nil {
	// 		internetConfigMap["sla_type"] = cngwGateway.InternetConfig.SlaType
	// 	}

	// 	if cngwGateway.InternetConfig.MultiZoneFlag != nil {
	// 		internetConfigMap["multi_zone_flag"] = cngwGateway.InternetConfig.MultiZoneFlag
	// 	}

	// 	if cngwGateway.InternetConfig.MasterZoneId != nil {
	// 		internetConfigMap["master_zone_id"] = cngwGateway.InternetConfig.MasterZoneId
	// 	}

	// 	if cngwGateway.InternetConfig.SlaveZoneId != nil {
	// 		internetConfigMap["slave_zone_id"] = cngwGateway.InternetConfig.SlaveZoneId
	// 	}

	// 	_ = d.Set("internet_config", []interface{}{internetConfigMap})
	// }

	tcClient := meta.(*TencentCloudClient).apiV3Conn
	tagService := &TagService{client: tcClient}
	tags, err := tagService.DescribeResourceTags(ctx, "tse", "gateway", tcClient.Region, d.Id())
	if err != nil {
		return err
	}
	_ = d.Set("tags", tags)

	return nil
}

func resourceTencentCloudTseCngwGatewayUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tse_cngw_gateway.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	request := tse.NewModifyCloudNativeAPIGatewayRequest()

	gatewayId := d.Id()

	request.GatewayId = &gatewayId

	immutableArgs := []string{"type", "gateway_version", "node_config", "vpc_config", "feature_version", "internet_max_bandwidth_out", "engine_region", "ingress_class_name", "trade_type", "internet_config"}
	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	changeFlag := false
	if d.HasChange("name") {
		if v, ok := d.GetOk("name"); ok {
			changeFlag = true
			request.Name = helper.String(v.(string))
		}
	}

	if d.HasChange("description") {
		if v, ok := d.GetOk("description"); ok {
			changeFlag = true
			request.Description = helper.String(v.(string))
		}
	}

	if d.HasChange("enable_cls") {
		if v, ok := d.GetOkExists("enable_cls"); ok {
			changeFlag = true
			request.EnableCls = helper.Bool(v.(bool))
		}
	}

	if changeFlag {
		err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			result, e := meta.(*TencentCloudClient).apiV3Conn.UseTseClient().ModifyCloudNativeAPIGateway(request)
			if e != nil {
				return retryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s update tse cngwGateway failed, reason:%+v", logId, err)
			return err
		}

		service := TseService{client: meta.(*TencentCloudClient).apiV3Conn}
		if err := service.CheckTseNativeAPIGatewayStatusById(ctx, gatewayId, "update"); err != nil {
			return err
		}
	}

	if d.HasChange("tags") {
		ctx := context.WithValue(context.TODO(), logIdKey, logId)
		tcClient := meta.(*TencentCloudClient).apiV3Conn
		tagService := &TagService{client: tcClient}
		oldTags, newTags := d.GetChange("tags")
		replaceTags, deleteTags := diffTags(oldTags.(map[string]interface{}), newTags.(map[string]interface{}))
		resourceName := BuildTagResourceName("tse", "gateway", tcClient.Region, d.Id())
		if err := tagService.ModifyTags(ctx, resourceName, replaceTags, deleteTags); err != nil {
			return err
		}
	}

	return resourceTencentCloudTseCngwGatewayRead(d, meta)
}

func resourceTencentCloudTseCngwGatewayDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tse_cngw_gateway.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TseService{client: meta.(*TencentCloudClient).apiV3Conn}
	gatewayId := d.Id()

	if err := service.DeleteTseCngwGatewayById(ctx, gatewayId); err != nil {
		return err
	}
	if err := service.CheckTseNativeAPIGatewayStatusById(ctx, gatewayId, "delete"); err != nil {
		return err
	}

	return nil
}
