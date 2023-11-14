/*
Provides a resource to create a tse instance

Example Usage

```hcl
resource "tencentcloud_tse_instance" "instance" {
  engine_type = "nacos"
  engine_version = "2.0.3"
  engine_product_version = "STANDARD"
  engine_region = "ap-guangzhou"
  engine_name = "nacos-test"
  trade_type = 0
  engine_resource_spec = "STANDARD"
  engine_node_num = 3
  vpc_id = "vpc-xxxxxx"
  subnet_id = "subnet-xxxxxx"
  apollo_env_params {
		name = "dev"
		engine_resource_spec = "1C2G"
		engine_node_num = 3
		storage_capacity = 20
		vpc_id = "vpc-xxxxxx"
		subnet_id = "subnet-xxxxxx"
		env_desc = "dev env"

  }
  engine_tags {
		tag_key = ""
		tag_value = ""

  }
  engine_admin {
		name = "admin"
		password = "admin"
		token = "xxxxxx"

  }
  prepaid_period = 0
  prepaid_renew_flag = 1
  engine_region_infos {
		engine_region = "ap-guangzhou"
		replica = 3
		vpc_infos {
			vpc_id = "vpc-xxxxxx"
			subnet_id = "subnet-xxxxxx"
			intranet_address = ""
		}

  }
  tags = {
    "createdBy" = "terraform"
  }
}
```

Import

tse instance can be imported using the id, e.g.

```
terraform import tencentcloud_tse_instance.instance instance_id
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
				Description: "Engine type. Reference value：- zookeeper- nacos- consul- apollo- eureka- polaris.",
			},

			"engine_version": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "An open source version of the engine. Each engine supports different open source versions, refer to the product documentation or console purchase page.",
			},

			"engine_product_version": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Engine product version. Reference value：- STANDARD： Standard editionEngine versions and optional specifications and number of nodes:：apollo - STANDARD versionspec list：1C2G、2C4G、4C8G、8C16G、16C32Gnode num：1，2，3，4，5eureka - STANDARD versionspec list：1C2G、2C4G、4C8G、8C16G、16C32Gnode num：3，4，5polarismesh - STANDARD versionspec list：NUM50、NUM100、NUM200、NUM500、NUM1000、NUM5000、NUM10000、NUM50000.",
			},

			"engine_region": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Engine deploy region. Reference value： China area Reference value：- ap-guangzhou：guangzhou- ap-beijing：beijing- ap-chengdu：chengdu- ap-chongqing：chongqing- ap-nanjing：nanjing- ap-shanghai：shanghai- ap-hongkong：hongkong- ap-taipei：taipeiAsia Pacific area Reference value：- ap-jakarta：jakarta- ap-singapore：singaporeNorth America area Reference value：- na-toronto：torontoFinancial area Reference value- ap-beijing-fsi：beijing-fsi- ap-shanghai-fsi：shanghai-fsi- ap-shenzhen-fsi：shenzhen-fsi.",
			},

			"engine_name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Engien name. Reference value：- nacos-test.",
			},

			"trade_type": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "Trade type. Reference value：- 0：postpaid- 1：Prepaid (Interface does not support the creation of prepaid instances yet).",
			},

			"engine_resource_spec": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Engine spec ID. see EngineProductVersion.",
			},

			"engine_node_num": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Engine node num. see EngineProductVersion.",
			},

			"vpc_id": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "VPC ID. Assign an IP address to the engine in the VPC subnet. Reference value：- vpc-conz6aix.",
			},

			"subnet_id": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Subnet ID. Assign an IP address to the engine in the VPC subnet. Reference value：- subnet-ahde9me9.",
			},

			"apollo_env_params": {
				Optional:    true,
				Type:        schema.TypeList,
				Description: "Apollo env config param list. Reference value：If the Apollo type is created, this parameter is a mandatory environment list. You can select a maximum of four environments. Description of environmental information parameters：- Name：env name. Reference value：prod, dev, fat, uat- EngineResourceSpec：env engine spec id. see EngineProductVersion- EngineNodeNum：env engien node num. see EngineProductVersion，prod  env 2，3，4，5- StorageCapacity：Configure the storage space size，The unit is GB， The step is 5， Reference value： 35- VpcId：VPC ID. Reference value：vpc-conz6aix- SubnetId：subnet ID. Reference value：subnet-ahde9me9.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Env name.",
						},
						"engine_resource_spec": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Node spec ID of the engine in the environment -1C2G-2C4G.",
						},
						"engine_node_num": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Number of engine nodes in the environment.",
						},
						"storage_capacity": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Configure the storage space size， in GB.",
						},
						"vpc_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "VPC ID. Assign an IP address to the VPC subnet as the ConfigServer access address.",
						},
						"subnet_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Subnet ID. Assign an IP address to the VPC subnet as the ConfigServer access address.",
						},
						"env_desc": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Env desc.",
						},
					},
				},
			},

			"engine_tags": {
				Optional:    true,
				Type:        schema.TypeList,
				Description: "List of tags for the engine. The value is a user-defined key/value without reference value.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"tag_key": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Tag key.",
						},
						"tag_value": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Tag val.",
						},
					},
				},
			},

			"engine_admin": {
				Optional:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Initial account information of the engine. Configurable parameters：- Name：Initial user name of the console- Password：Console initial password- Token：Engine interface admin Token.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Initial user name of the console.",
						},
						"password": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Console initial password.",
						},
						"token": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Engine interface admin Token.",
						},
					},
				},
			},

			"prepaid_period": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Prepaid time， in monthly units.",
			},

			"prepaid_renew_flag": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Automatic renewal mark， prepaid only.  Reference value：- 0：No automatic renewal- 1：Automatic renewal.",
			},

			"engine_region_infos": {
				Optional:    true,
				Type:        schema.TypeList,
				Description: "Details about the regional configuration of the engine in cross-region deployment.",
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
										Description: "Intranet access addressNote: This field may return null, indicating that a valid value is not available.. .",
									},
								},
							},
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

func resourceTencentCloudTseInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tse_instance.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

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

	if v, ok := d.GetOk("apollo_env_params"); ok {
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			apolloEnvParam := tse.ApolloEnvParam{}
			if v, ok := dMap["name"]; ok {
				apolloEnvParam.Name = helper.String(v.(string))
			}
			if v, ok := dMap["engine_resource_spec"]; ok {
				apolloEnvParam.EngineResourceSpec = helper.String(v.(string))
			}
			if v, ok := dMap["engine_node_num"]; ok {
				apolloEnvParam.EngineNodeNum = helper.IntInt64(v.(int))
			}
			if v, ok := dMap["storage_capacity"]; ok {
				apolloEnvParam.StorageCapacity = helper.IntInt64(v.(int))
			}
			if v, ok := dMap["vpc_id"]; ok {
				apolloEnvParam.VpcId = helper.String(v.(string))
			}
			if v, ok := dMap["subnet_id"]; ok {
				apolloEnvParam.SubnetId = helper.String(v.(string))
			}
			if v, ok := dMap["env_desc"]; ok {
				apolloEnvParam.EnvDesc = helper.String(v.(string))
			}
			request.ApolloEnvParams = append(request.ApolloEnvParams, &apolloEnvParam)
		}
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

	if dMap, ok := helper.InterfacesHeadMap(d, "engine_admin"); ok {
		engineAdmin := tse.EngineAdmin{}
		if v, ok := dMap["name"]; ok {
			engineAdmin.Name = helper.String(v.(string))
		}
		if v, ok := dMap["password"]; ok {
			engineAdmin.Password = helper.String(v.(string))
		}
		if v, ok := dMap["token"]; ok {
			engineAdmin.Token = helper.String(v.(string))
		}
		request.EngineAdmin = &engineAdmin
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

	ctx := context.WithValue(context.TODO(), logIdKey, logId)
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

	if instance.EngineType != nil {
		_ = d.Set("engine_type", instance.EngineType)
	}

	if instance.EngineVersion != nil {
		_ = d.Set("engine_version", instance.EngineVersion)
	}

	if instance.EngineProductVersion != nil {
		_ = d.Set("engine_product_version", instance.EngineProductVersion)
	}

	if instance.EngineRegion != nil {
		_ = d.Set("engine_region", instance.EngineRegion)
	}

	if instance.EngineName != nil {
		_ = d.Set("engine_name", instance.EngineName)
	}

	if instance.TradeType != nil {
		_ = d.Set("trade_type", instance.TradeType)
	}

	if instance.EngineResourceSpec != nil {
		_ = d.Set("engine_resource_spec", instance.EngineResourceSpec)
	}

	if instance.EngineNodeNum != nil {
		_ = d.Set("engine_node_num", instance.EngineNodeNum)
	}

	if instance.VpcId != nil {
		_ = d.Set("vpc_id", instance.VpcId)
	}

	if instance.SubnetId != nil {
		_ = d.Set("subnet_id", instance.SubnetId)
	}

	if instance.ApolloEnvParams != nil {
		apolloEnvParamsList := []interface{}{}
		for _, apolloEnvParams := range instance.ApolloEnvParams {
			apolloEnvParamsMap := map[string]interface{}{}

			if instance.ApolloEnvParams.Name != nil {
				apolloEnvParamsMap["name"] = instance.ApolloEnvParams.Name
			}

			if instance.ApolloEnvParams.EngineResourceSpec != nil {
				apolloEnvParamsMap["engine_resource_spec"] = instance.ApolloEnvParams.EngineResourceSpec
			}

			if instance.ApolloEnvParams.EngineNodeNum != nil {
				apolloEnvParamsMap["engine_node_num"] = instance.ApolloEnvParams.EngineNodeNum
			}

			if instance.ApolloEnvParams.StorageCapacity != nil {
				apolloEnvParamsMap["storage_capacity"] = instance.ApolloEnvParams.StorageCapacity
			}

			if instance.ApolloEnvParams.VpcId != nil {
				apolloEnvParamsMap["vpc_id"] = instance.ApolloEnvParams.VpcId
			}

			if instance.ApolloEnvParams.SubnetId != nil {
				apolloEnvParamsMap["subnet_id"] = instance.ApolloEnvParams.SubnetId
			}

			if instance.ApolloEnvParams.EnvDesc != nil {
				apolloEnvParamsMap["env_desc"] = instance.ApolloEnvParams.EnvDesc
			}

			apolloEnvParamsList = append(apolloEnvParamsList, apolloEnvParamsMap)
		}

		_ = d.Set("apollo_env_params", apolloEnvParamsList)

	}

	if instance.EngineTags != nil {
		engineTagsList := []interface{}{}
		for _, engineTags := range instance.EngineTags {
			engineTagsMap := map[string]interface{}{}

			if instance.EngineTags.TagKey != nil {
				engineTagsMap["tag_key"] = instance.EngineTags.TagKey
			}

			if instance.EngineTags.TagValue != nil {
				engineTagsMap["tag_value"] = instance.EngineTags.TagValue
			}

			engineTagsList = append(engineTagsList, engineTagsMap)
		}

		_ = d.Set("engine_tags", engineTagsList)

	}

	if instance.EngineAdmin != nil {
		engineAdminMap := map[string]interface{}{}

		if instance.EngineAdmin.Name != nil {
			engineAdminMap["name"] = instance.EngineAdmin.Name
		}

		if instance.EngineAdmin.Password != nil {
			engineAdminMap["password"] = instance.EngineAdmin.Password
		}

		if instance.EngineAdmin.Token != nil {
			engineAdminMap["token"] = instance.EngineAdmin.Token
		}

		_ = d.Set("engine_admin", []interface{}{engineAdminMap})
	}

	if instance.PrepaidPeriod != nil {
		_ = d.Set("prepaid_period", instance.PrepaidPeriod)
	}

	if instance.PrepaidRenewFlag != nil {
		_ = d.Set("prepaid_renew_flag", instance.PrepaidRenewFlag)
	}

	if instance.EngineRegionInfos != nil {
		engineRegionInfosList := []interface{}{}
		for _, engineRegionInfos := range instance.EngineRegionInfos {
			engineRegionInfosMap := map[string]interface{}{}

			if instance.EngineRegionInfos.EngineRegion != nil {
				engineRegionInfosMap["engine_region"] = instance.EngineRegionInfos.EngineRegion
			}

			if instance.EngineRegionInfos.Replica != nil {
				engineRegionInfosMap["replica"] = instance.EngineRegionInfos.Replica
			}

			if instance.EngineRegionInfos.VpcInfos != nil {
				vpcInfosList := []interface{}{}
				for _, vpcInfos := range instance.EngineRegionInfos.VpcInfos {
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

				engineRegionInfosMap["vpc_infos"] = []interface{}{vpcInfosList}
			}

			engineRegionInfosList = append(engineRegionInfosList, engineRegionInfosMap)
		}

		_ = d.Set("engine_region_infos", engineRegionInfosList)

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

	request := tse.NewUpdateEngineInternetAccessRequest()

	instanceId := d.Id()

	request.InstanceId = &instanceId

	immutableArgs := []string{"engine_type", "engine_version", "engine_product_version", "engine_region", "engine_name", "trade_type", "engine_resource_spec", "engine_node_num", "vpc_id", "subnet_id", "apollo_env_params", "engine_tags", "engine_admin", "prepaid_period", "prepaid_renew_flag", "engine_region_infos"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	if d.HasChange("engine_type") {
		if v, ok := d.GetOk("engine_type"); ok {
			request.EngineType = helper.String(v.(string))
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

	return nil
}
