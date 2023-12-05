package tencentcloud

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tse "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tse/v20201207"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudTseInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTseInstanceCreate,
		Read:   resourceTencentCloudTseInstanceRead,
		Update: resourceTencentCloudTseInstanceUpdate,
		Delete: resourceTencentCloudTseInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"engine_type": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "engine type. Reference value: `zookeeper`, `nacos`, `polaris`.",
			},

			"engine_version": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "An open source version of the engine. Each engine supports different open source versions, refer to the product documentation or console purchase page.",
			},

			"engine_product_version": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Engine product version. Reference value: `Nacos`: `TRIAL`: Development version, optional node num: `1`, optional spec list: `1C1G`; `STANDARD`: Standard versions, optional node num: `3`, `5`, `7`, optional spec list: `1C2G`, `2C4G`, `4C8G`, `8C16G`, `16C32G`. `Zookeeper`: `TRIAL`: Development version, optional node num: `1`, optional spec list: `1C1G`; `STANDARD`: Standard versions, optional node num: `3`, `5`, `7`, optional spec list: `1C2G`, `2C4G`, `4C8G`, `8C16G`, `16C32G`; `PROFESSIONAL`: professional versions, optional node num: `3`, `5`, `7`, optional spec list: `1C2G`, `2C4G`, `4C8G`, `8C16G`, `16C32G`. `Polarismesh`: `BASE`: Base version, optional node num: `1`, optional spec list: `NUM50`; `PROFESSIONAL`: Enterprise versions, optional node num: `2`, `3`, optional spec list: `NUM50`, `NUM100`, `NUM200`, `NUM500`, `NUM1000`, `NUM5000`, `NUM10000`, `NUM50000`.",
			},

			"engine_region": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "engine deploy region. Reference value: `China area` Reference value: `ap-guangzhou`, `ap-beijing`, `ap-chengdu`, `ap-chongqing`, `ap-nanjing`, `ap-shanghai` `ap-beijing-fsi`, `ap-shanghai-fsi`, `ap-shenzhen-fsi`. `Asia Pacific` area Reference value: `ap-hongkong`, `ap-taipei`, `ap-jakarta`, `ap-singapore`, `ap-bangkok`, `ap-seoul`, `ap-tokyo`. `North America area` Reference value: `na-toronto`, `sa-saopaulo`, `na-siliconvalley`, `na-ashburn`.",
			},

			"engine_name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "engien name. Reference value: nacos-test.",
			},

			"trade_type": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "trade type. Reference value:- 0:postpaid- 1:Prepaid (Interface does not support the creation of prepaid instances yet).",
			},

			"engine_resource_spec": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "engine spec ID. see EngineProductVersion.",
			},

			"engine_node_num": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "engine node num. see EngineProductVersion.",
			},

			"vpc_id": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "VPC ID. Assign an IP address to the engine in the VPC subnet. Reference value: vpc-conz6aix.",
			},

			"subnet_id": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "subnet ID. Assign an IP address to the engine in the VPC subnet. Reference value: subnet-ahde9me9.",
			},

			"prepaid_period": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Prepaid time, in monthly units.",
			},

			"prepaid_renew_flag": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Automatic renewal mark, prepaid only.  Reference value: `0`: No automatic renewal, `1`: Automatic renewal.",
			},

			"engine_region_infos": {
				Optional:    true,
				Type:        schema.TypeList,
				Description: "Details about the regional configuration of the engine in cross-region deployment, only zookeeper professional requires the use of the EngineRegionInfos parameter.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"engine_region": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Engine node region.",
						},
						"replica": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "The number of nodes allocated in this region.",
						},
						"vpc_infos": {
							Type:        schema.TypeList,
							Required:    true,
							Description: "Cluster network information.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"vpc_id": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Vpc Id.",
									},
									"subnet_id": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Subnet ID.",
									},
									"intranet_address": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Intranet access addressNote: This field may return null, indicating that a valid value is not available..",
									},
								},
							},
						},
					},
				},
			},

			"enable_client_internet_access": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Client public network access, `true`: on, `false`: off, default: false.",
			},

			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Tag description list.",
			},
		},
	}
}

func resourceTencentCloudTseInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tse_instance.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	var (
		request    = tse.NewCreateEngineRequest()
		response   = tse.NewCreateEngineResponse()
		instanceId string
	)
	if v, ok := d.GetOk("engine_type"); ok {
		request.EngineType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("engine_version"); ok {
		request.EngineVersion = helper.String(v.(string))
	}

	if v, ok := d.GetOk("engine_product_version"); ok {
		request.EngineProductVersion = helper.String(v.(string))
	}

	if v, ok := d.GetOk("engine_region"); ok {
		request.EngineRegion = helper.String(v.(string))
	}

	if v, ok := d.GetOk("engine_name"); ok {
		request.EngineName = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("trade_type"); ok {
		request.TradeType = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("engine_resource_spec"); ok {
		request.EngineResourceSpec = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("engine_node_num"); ok {
		request.EngineNodeNum = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("vpc_id"); ok {
		request.VpcId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("subnet_id"); ok {
		request.SubnetId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("engine_tags"); ok {
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			instanceTagInfo := tse.InstanceTagInfo{}
			if v, ok := dMap["tag_key"]; ok {
				instanceTagInfo.TagKey = helper.String(v.(string))
			}
			if v, ok := dMap["tag_value"]; ok {
				instanceTagInfo.TagValue = helper.String(v.(string))
			}
			request.EngineTags = append(request.EngineTags, &instanceTagInfo)
		}
	}

	if v, ok := d.GetOkExists("prepaid_period"); ok {
		request.PrepaidPeriod = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("prepaid_renew_flag"); ok {
		request.PrepaidRenewFlag = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("engine_region_infos"); ok {
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			engineRegionInfo := tse.EngineRegionInfo{}
			if v, ok := dMap["engine_region"]; ok {
				engineRegionInfo.EngineRegion = helper.String(v.(string))
			}
			if v, ok := dMap["replica"]; ok {
				engineRegionInfo.Replica = helper.IntInt64(v.(int))
			}
			if v, ok := dMap["vpc_infos"]; ok {
				for _, item := range v.([]interface{}) {
					vpcInfosMap := item.(map[string]interface{})
					vpcInfo := tse.VpcInfo{}
					if v, ok := vpcInfosMap["vpc_id"]; ok {
						vpcInfo.VpcId = helper.String(v.(string))
					}
					if v, ok := vpcInfosMap["subnet_id"]; ok {
						vpcInfo.SubnetId = helper.String(v.(string))
					}
					if v, ok := vpcInfosMap["intranet_address"]; ok {
						vpcInfo.IntranetAddress = helper.String(v.(string))
					}
					engineRegionInfo.VpcInfos = append(engineRegionInfo.VpcInfos, &vpcInfo)
				}
			}
			request.EngineRegionInfos = append(request.EngineRegionInfos, &engineRegionInfo)
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTseClient().CreateEngine(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create tse instance failed, reason:%+v", logId, err)
		return err
	}

	instanceId = *response.Response.InstanceId
	d.SetId(instanceId)

	service := TseService{client: meta.(*TencentCloudClient).apiV3Conn}
	if err := service.CheckTseInstanceStatusById(ctx, instanceId, "create"); err != nil {
		return err
	}

	if tags := helper.GetTags(d, "tags"); len(tags) > 0 {
		tagService := TagService{client: meta.(*TencentCloudClient).apiV3Conn}
		region := meta.(*TencentCloudClient).apiV3Conn.Region
		resourceName := fmt.Sprintf("qcs::tse:%s:uin/:instance/%s", region, d.Id())
		if err := tagService.ModifyTags(ctx, resourceName, tags, nil); err != nil {
			return err
		}
	}

	return resourceTencentCloudTseInstanceRead(d, meta)
}

func resourceTencentCloudTseInstanceRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tse_instance.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TseService{client: meta.(*TencentCloudClient).apiV3Conn}

	instanceId := d.Id()

	instance, err := service.DescribeTseInstanceById(ctx, instanceId)
	if err != nil {
		return err
	}

	if instance == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `TseInstance` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if instance.Type != nil {
		_ = d.Set("engine_type", instance.Type)
	}

	if instance.Edition != nil {
		_ = d.Set("engine_version", instance.Edition)
	}

	if instance.FeatureVersion != nil {
		_ = d.Set("engine_product_version", instance.FeatureVersion)
	}

	if instance.EngineRegion != nil {
		_ = d.Set("engine_region", instance.EngineRegion)
	}

	if instance.Name != nil {
		_ = d.Set("engine_name", instance.Name)
	}

	if instance.TradeType != nil {
		_ = d.Set("trade_type", instance.TradeType)
	}

	if instance.SpecId != nil {
		_ = d.Set("engine_resource_spec", instance.SpecId)
	}

	if instance.Replica != nil {
		_ = d.Set("engine_node_num", instance.Replica)
	}

	if instance.VpcId != nil {
		_ = d.Set("vpc_id", instance.VpcId)
	}

	if instance.SubnetIds != nil {
		_ = d.Set("subnet_id", instance.SubnetIds[0])
	}

	if instance.Tags != nil {
		engineTagsList := []interface{}{}
		for _, engineTags := range instance.Tags {
			engineTagsMap := map[string]interface{}{}

			if engineTags.Key != nil {
				engineTagsMap["tag_key"] = engineTags.Key
			}

			if engineTags.Value != nil {
				engineTagsMap["tag_value"] = engineTags.Value
			}

			engineTagsList = append(engineTagsList, engineTagsMap)
		}

		_ = d.Set("engine_tags", engineTagsList)

	}

	if instance.RegionInfos != nil && *instance.Type == "zookeeper" && *instance.FeatureVersion == "PROFESSIONAL" {
		engineRegionInfosList := []interface{}{}
		for _, engineRegionInfos := range instance.RegionInfos {
			engineRegionInfosMap := map[string]interface{}{}

			if engineRegionInfos.EngineRegion != nil {
				engineRegionInfosMap["engine_region"] = engineRegionInfos.EngineRegion
			}

			if engineRegionInfos.Replica != nil {
				engineRegionInfosMap["replica"] = engineRegionInfos.Replica
			}

			if engineRegionInfos.IntranetVpcInfos != nil {
				vpcInfosList := []interface{}{}
				for _, vpcInfos := range engineRegionInfos.IntranetVpcInfos {
					vpcInfosMap := map[string]interface{}{}

					if vpcInfos.VpcId != nil {
						vpcInfosMap["vpc_id"] = vpcInfos.VpcId
					}

					if vpcInfos.SubnetId != nil {
						vpcInfosMap["subnet_id"] = vpcInfos.SubnetId
					}

					if vpcInfos.IntranetAddress != nil {
						vpcInfosMap["intranet_address"] = vpcInfos.IntranetAddress
					}

					vpcInfosList = append(vpcInfosList, vpcInfosMap)
				}

				engineRegionInfosMap["vpc_infos"] = vpcInfosList
			}

			engineRegionInfosList = append(engineRegionInfosList, engineRegionInfosMap)
		}

		_ = d.Set("engine_region_infos", engineRegionInfosList)

	}

	if instance.EnableInternet != nil {
		_ = d.Set("enable_client_internet_access", instance.EnableInternet)
	}

	tcClient := meta.(*TencentCloudClient).apiV3Conn
	tagService := &TagService{client: tcClient}
	tags, err := tagService.DescribeResourceTags(ctx, "tse", "instance", tcClient.Region, d.Id())
	if err != nil {
		return err
	}
	_ = d.Set("tags", tags)

	return nil
}

func resourceTencentCloudTseInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tse_instance.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	request := tse.NewUpdateEngineInternetAccessRequest()
	instanceId := d.Id()
	request.InstanceId = &instanceId

	immutableArgs := []string{"engine_type", "engine_version", "engine_product_version", "engine_region", "engine_name", "trade_type", "engine_resource_spec", "engine_node_num", "vpc_id", "subnet_id", "engine_tags", "prepaid_period", "prepaid_renew_flag", "engine_region_infos"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	if d.HasChange("enable_client_internet_access") {
		if v, ok := d.GetOk("engine_type"); ok {
			request.EngineType = helper.String(v.(string))
		}

		if v, _ := d.GetOk("enable_client_internet_access"); v != nil {
			request.EnableClientInternetAccess = helper.Bool(v.(bool))
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTseClient().UpdateEngineInternetAccess(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update tse instance failed, reason:%+v", logId, err)
		return err
	}

	service := TseService{client: meta.(*TencentCloudClient).apiV3Conn}
	if err := service.CheckTseInstanceStatusById(ctx, instanceId, "update"); err != nil {
		return err
	}

	if d.HasChange("tags") {
		ctx := context.WithValue(context.TODO(), logIdKey, logId)
		tcClient := meta.(*TencentCloudClient).apiV3Conn
		tagService := &TagService{client: tcClient}
		oldTags, newTags := d.GetChange("tags")
		replaceTags, deleteTags := diffTags(oldTags.(map[string]interface{}), newTags.(map[string]interface{}))
		resourceName := BuildTagResourceName("tse", "instance", tcClient.Region, d.Id())
		if err := tagService.ModifyTags(ctx, resourceName, replaceTags, deleteTags); err != nil {
			return err
		}
	}

	return resourceTencentCloudTseInstanceRead(d, meta)
}

func resourceTencentCloudTseInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tse_instance.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TseService{client: meta.(*TencentCloudClient).apiV3Conn}
	instanceId := d.Id()

	if err := service.DeleteTseInstanceById(ctx, instanceId); err != nil {
		return err
	}
	if err := service.CheckTseInstanceStatusById(ctx, instanceId, "delete"); err != nil {
		return err
	}

	return nil
}
