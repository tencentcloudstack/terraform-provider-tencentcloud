package es

import (
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svccvm "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/cvm"
	svctag "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/tag"

	"context"
	"errors"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sdkErrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	es "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/es/v20180416"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

//internal version: replace import begin, please do not modify this annotation and refrain from inserting any code between the beginning and end lines of the annotation.
//internal version: replace import end, please do not modify this annotation and refrain from inserting any code between the beginning and end lines of the annotation.

func ResourceTencentCloudElasticsearchInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudElasticsearchInstanceCreate,
		Read:   resourceTencentCloudElasticsearchInstanceRead,
		Update: resourceTencentCloudElasticsearchInstanceUpdate,
		Delete: resourceTencentCloudElasticsearchInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: helper.ImportWithDefaultValue(map[string]interface{}{
				"basic_security_type": ES_BASIC_SECURITY_TYPE_OFF,
			}),
		},

		Schema: map[string]*schema.Schema{
			"instance_name": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: tccommon.ValidateStringLengthInRange(1, 50),
				Description:  "Name of the instance, which can contain 1 to 50 English letters, Chinese characters, digits, dashes(-), or underscores(_).",
			},
			"availability_zone": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "-",
				ForceNew:    true,
				Description: "Availability zone. When create multi-az es, this parameter must be omitted or `-`.",
			},
			"version": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Version of the instance. Valid values are `5.6.4`, `6.4.3`, `6.8.2`, `7.5.1` and `7.10.1`.",
			},
			"vpc_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The ID of a VPC network.",
			},
			"subnet_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "-",
				ForceNew:    true,
				Description: "The ID of a VPC subnetwork. When create multi-az es, this parameter must be omitted or `-`.",
			},
			"password": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				Description: "Password to an instance, the password needs to be 8 to 16 characters, including at least two items ([a-z,A-Z], [0-9] and [-!@#$%&^*+=_:;,.?] special symbols.",
			},
			"charge_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      ES_CHARGE_TYPE_POSTPAID_BY_HOUR,
				ValidateFunc: tccommon.ValidateAllowedStringValue(ES_CHARGE_TYPE),
				Description:  "The charge type of instance. Valid values are `PREPAID` and `POSTPAID_BY_HOUR`.",
			},
			"charge_period": {
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
				Default:     0,
				Description: "The tenancy of the prepaid instance, and uint is month. NOTE: it only works when charge_type is set to `PREPAID`.",
			},
			"renew_flag": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: tccommon.ValidateAllowedStringValue(ES_RENEW_FLAG),
				Description:  "When enabled, the instance will be renew automatically when it reach the end of the prepaid tenancy. Valid values are `RENEW_FLAG_AUTO` and `RENEW_FLAG_MANUAL`. NOTE: it only works when charge_type is set to `PREPAID`.",
			},
			"deploy_mode": {
				Type:         schema.TypeInt,
				Optional:     true,
				ForceNew:     true,
				Default:      ES_DEPLOY_MODE_SINGLE_REGION,
				ValidateFunc: tccommon.ValidateAllowedIntValue(ES_DEPLOY_MODE),
				Description:  "Cluster deployment mode. Valid values are `0` and `1`. `0` is single-AZ deployment, and `1` is multi-AZ deployment. Default value is `0`.",
			},
			"multi_zone_infos": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				Computed:    true,
				Description: "Details of AZs in multi-AZ deployment mode (which is required when deploy_mode is `1`).",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"availability_zone": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Availability zone.",
						},
						"subnet_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The ID of a VPC subnetwork.",
						},
					},
				},
			},
			"web_node_type_info": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Visual node configuration.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"node_num": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Visual node number.",
						},
						"node_type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Visual node specifications.",
						},
					},
				},
			},
			"es_acl": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				MaxItems:    1,
				Description: "Kibana Access Control Configuration.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"black_list": {
							Type:        schema.TypeSet,
							Optional:    true,
							Computed:    true,
							Description: "Blacklist of kibana access.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"white_list": {
							Type:        schema.TypeSet,
							Optional:    true,
							Computed:    true,
							Description: "Whitelist of kibana access.",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"license_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      ES_LICENSE_TYPE_PLATINUM,
				ValidateFunc: tccommon.ValidateAllowedStringValue(ES_LICENSE_TYPE),
				Description:  "License type. Valid values are `oss`, `basic` and `platinum`. The default value is `platinum`.",
			},
			"node_info_list": {
				Type:        schema.TypeList,
				Required:    true,
				MinItems:    1,
				Description: "Node information list, which is used to describe the specification information of various types of nodes in the cluster, such as node type, node quantity, node specification, disk type, and disk size.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"node_num": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Number of nodes.",
						},
						"node_type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Node specification, and valid values refer to [document of tencentcloud](https://intl.cloud.tencent.com/document/product/845/18376).",
						},
						"type": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      ES_NODE_TYPE_HOT_DATA,
							ValidateFunc: tccommon.ValidateAllowedStringValue(ES_NODE_TYPE),
							Description:  "Node type. Valid values are `hotData`, `warmData` and `dedicatedMaster`. The default value is 'hotData`.",
						},
						"disk_type": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      svccvm.CVM_DISK_TYPE_CLOUD_SSD,
							ValidateFunc: tccommon.ValidateAllowedStringValue(ES_NODE_DISK_TYPE),
							Description:  "Node disk type. Valid values are `CLOUD_SSD` and `CLOUD_PREMIUM`, `CLOUD_HSSD`. The default value is `CLOUD_SSD`.",
						},
						"disk_size": {
							Type:        schema.TypeInt,
							Optional:    true,
							Default:     100,
							Description: "Node disk size. Unit is GB, and default value is `100`.",
						},
						"encrypt": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "Decides to encrypt this disk or not.",
						},
					},
				},
			},
			"basic_security_type": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      ES_BASIC_SECURITY_TYPE_OFF,
				ValidateFunc: tccommon.ValidateAllowedIntValue(ES_BASIC_SECURITY_TYPE),
				Description:  "Whether to enable X-Pack security authentication in Basic Edition 6.8 and above. Valid values are `1` and `2`. `1` is disabled, `2` is enabled, and default value is `1`. Notice: this parameter is only take effect on `basic` license.",
			},
			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "A mapping of tags to assign to the instance. For tag limits, please refer to [Use Limits](https://intl.cloud.tencent.com/document/product/651/13354).",
			},

			// computed
			"elasticsearch_domain": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Elasticsearch domain name.",
			},
			"elasticsearch_vip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Elasticsearch VIP.",
			},
			"elasticsearch_port": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Elasticsearch port.",
			},
			"kibana_url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Kibana access URL.",
			},
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Instance creation time.",
			},
		},
	}
}

func resourceTencentCloudElasticsearchInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_elasticsearch_instance.create")()
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	elasticsearchService := ElasticsearchService{
		client: meta.(tccommon.ProviderMeta).GetAPIV3Conn(),
	}
	request := es.NewCreateInstanceRequest()

	//internal version: replace var begin, please do not modify this annotation and refrain from inserting any code between the beginning and end lines of the annotation.
	//internal version: replace var end, please do not modify this annotation and refrain from inserting any code between the beginning and end lines of the annotation.
	request.Zone = helper.String(d.Get("availability_zone").(string))
	request.EsVersion = helper.String(d.Get("version").(string))
	request.VpcId = helper.String(d.Get("vpc_id").(string))
	request.SubnetId = helper.String(d.Get("subnet_id").(string))
	request.Password = helper.String(d.Get("password").(string))
	if v, ok := d.GetOk("instance_name"); ok {
		request.InstanceName = helper.String(v.(string))
	}
	if v, ok := d.GetOk("charge_type"); ok {
		//internal version: replace strCharge begin, please do not modify this annotation and refrain from inserting any code between the beginning and end lines of the annotation.
		chargeType := v.(string)
		//internal version: replace strCharge end, please do not modify this annotation and refrain from inserting any code between the beginning and end lines of the annotation.
		request.ChargeType = &chargeType
		if chargeType == ES_CHARGE_TYPE_PREPAID {
			if v, ok := d.GetOk("charge_period"); ok {
				request.ChargePeriod = helper.IntUint64(v.(int))
			} else {
				return fmt.Errorf("elasticsearch charge period can not be empty when charge type is %s", chargeType)
			}
			if v, ok := d.GetOk("renew_flag"); ok {
				request.RenewFlag = helper.String(v.(string))
			}
		}
	}
	if v, ok := d.GetOk("deploy_mode"); ok {
		deployMode := v.(int)
		request.DeployMode = helper.IntUint64(deployMode)
		if deployMode == ES_DEPLOY_MODE_MULTI_REGION {
			if v, ok := d.GetOk("multi_zone_infos"); ok {
				infos := v.([]interface{})
				request.MultiZoneInfo = make([]*es.ZoneDetail, 0, len(infos))
				for _, item := range infos {
					value := item.(map[string]interface{})
					info := es.ZoneDetail{
						Zone:     helper.String(value["availability_zone"].(string)),
						SubnetId: helper.String(value["subnet_id"].(string)),
					}
					request.MultiZoneInfo = append(request.MultiZoneInfo, &info)
				}
			} else {
				return fmt.Errorf("elasticsearch multi_zone_infos can not be empty when deploy mode is %d", deployMode)
			}
		}
	}
	var licenseType string
	if v, ok := d.GetOk("license_type"); ok {
		licenseType = v.(string)
		request.LicenseType = helper.String(licenseType)
	}
	if v, ok := d.GetOk("basic_security_type"); ok {
		if licenseType == ES_LICENSE_TYPE_BASIC { // this field is only valid for the basic edition
			request.BasicSecurityType = helper.IntUint64(v.(int))
		}
	}
	if v, ok := d.GetOk("web_node_type_info"); ok {
		infos := v.([]interface{})
		for _, item := range infos {
			value := item.(map[string]interface{})
			info := &es.WebNodeTypeInfo{
				NodeNum:  helper.IntUint64(value["node_num"].(int)),
				NodeType: helper.String(value["node_type"].(string)),
			}
			request.WebNodeTypeInfo = info
			break
		}
	}

	if v, ok := d.GetOk("node_info_list"); ok {
		infos := v.([]interface{})
		request.NodeInfoList = make([]*es.NodeInfo, 0, len(infos))
		for _, item := range infos {
			value := item.(map[string]interface{})
			info := es.NodeInfo{
				NodeNum:  helper.IntUint64(value["node_num"].(int)),
				NodeType: helper.String(value["node_type"].(string)),
			}
			if v := value["type"].(string); v != "" {
				info.Type = &v
			}
			if v := value["disk_type"].(string); v != "" {
				info.DiskType = &v
			}
			if v := value["disk_size"].(int); v > 0 {
				info.DiskSize = helper.IntUint64(v)
			}
			if v := value["encrypt"].(bool); v {
				info.DiskEncrypt = helper.BoolToInt64Pointer(v)
			}
			request.NodeInfoList = append(request.NodeInfoList, &info)
		}
	}
	//internal version: replace reqTag begin, please do not modify this annotation and refrain from inserting any code between the beginning and end lines of the annotation.
	//internal version: replace reqTag end, please do not modify this annotation and refrain from inserting any code between the beginning and end lines of the annotation.
	instanceId := ""
	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		response, err := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseEsClient().CreateInstance(request)
		if err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), err.Error())
			//internal version: replace bpass begin, please do not modify this annotation and refrain from inserting any code between the beginning and end lines of the annotation.
			//internal version: replace bpass end, please do not modify this annotation and refrain from inserting any code between the beginning and end lines of the annotation.
			return tccommon.RetryError(err)
		}
		instanceId = *response.Response.InstanceId
		return nil
	})
	if err != nil {
		return err
	}
	d.SetId(instanceId)

	//internal version: replace setTag begin, please do not modify this annotation and refrain from inserting any code between the beginning and end lines of the annotation.
	//internal version: replace setTag end, please do not modify this annotation and refrain from inserting any code between the beginning and end lines of the annotation.
	instanceEmptyRetries := 5
	err = resource.Retry(15*tccommon.ReadRetryTimeout, func() *resource.RetryError {
		instance, errRet := elasticsearchService.DescribeInstanceById(ctx, instanceId)
		if errRet != nil {
			return tccommon.RetryError(errRet, tccommon.InternalError)
		}
		if instance == nil {
			if instanceEmptyRetries > 0 {
				instanceEmptyRetries -= 1
				return resource.RetryableError(fmt.Errorf("cannot find instance %s, retrying", instanceId))
			}
			return resource.NonRetryableError(fmt.Errorf("instance %s not exists", instanceId))
		}
		if *instance.Status != ES_INSTANCE_STATUS_NORMAL {
			return resource.RetryableError(fmt.Errorf("elasticsearch instance status is %v, retrying", *instance.Status))
		}
		return nil
	})
	if err != nil {
		return err
	}

	// es acl
	esAcl := es.EsAcl{}
	if aclMap, ok := helper.InterfacesHeadMap(d, "es_acl"); ok {
		if v, ok := aclMap["black_list"]; ok {
			bList := v.(*schema.Set).List()
			for _, d := range bList {
				esAcl.BlackIpList = append(esAcl.BlackIpList, helper.String(d.(string)))
			}
		}
		if v, ok := aclMap["white_list"]; ok {
			wList := v.(*schema.Set).List()
			for _, d := range wList {
				esAcl.WhiteIpList = append(esAcl.WhiteIpList, helper.String(d.(string)))
			}
		}
	}

	err = resource.Retry(tccommon.WriteRetryTimeout*2, func() *resource.RetryError {
		errRet := elasticsearchService.UpdateInstance(ctx, instanceId, "", "", 0, nil, nil, &esAcl)
		if errRet != nil {
			return tccommon.RetryError(errRet)
		}
		return nil
	})
	if err != nil {
		return err
	}

	// tags
	if tags := helper.GetTags(d, "tags"); len(tags) > 0 {
		client := meta.(tccommon.ProviderMeta).GetAPIV3Conn()
		tagService := svctag.NewTagService(client)
		region := client.Region
		resourceName := fmt.Sprintf("qcs::es:%s:uin/:instance/%s", region, instanceId)
		err := tagService.ModifyTags(ctx, resourceName, tags, nil)
		if err != nil {
			return err
		}
	}

	return resourceTencentCloudElasticsearchInstanceRead(d, meta)
}

func resourceTencentCloudElasticsearchInstanceRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_elasticsearch_instance.read")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	instanceId := d.Id()
	elasticsearchService := ElasticsearchService{
		client: meta.(tccommon.ProviderMeta).GetAPIV3Conn(),
	}

	var instance *es.InstanceInfo
	var errRet error
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		instance, errRet = elasticsearchService.DescribeInstanceById(ctx, instanceId)
		if errRet != nil {
			return tccommon.RetryError(errRet, tccommon.InternalError)
		}
		return nil
	})
	if err != nil {
		return err
	}
	if instance == nil {
		d.SetId("")
		return nil
	}

	_ = d.Set("instance_name", instance.InstanceName)
	_ = d.Set("availability_zone", instance.Zone)
	_ = d.Set("version", instance.EsVersion)
	_ = d.Set("vpc_id", instance.VpcUid)
	_ = d.Set("subnet_id", instance.SubnetUid)
	_ = d.Set("charge_type", instance.ChargeType)
	_ = d.Set("charge_period", instance.ChargePeriod)
	_ = d.Set("renew_flag", instance.RenewFlag)
	_ = d.Set("deploy_mode", instance.DeployMode)
	_ = d.Set("license_type", instance.LicenseType)
	licenseType := instance.LicenseType
	if licenseType != nil && *licenseType == ES_LICENSE_TYPE_BASIC { // this field is only valid for the basic edition
		_ = d.Set("basic_security_type", instance.SecurityType)
	}
	_ = d.Set("elasticsearch_domain", instance.EsDomain)
	_ = d.Set("elasticsearch_vip", instance.EsVip)
	_ = d.Set("elasticsearch_port", instance.EsPort)
	_ = d.Set("kibana_url", instance.KibanaUrl)
	_ = d.Set("create_time", instance.CreateTime)

	multiZoneInfos := make([]map[string]interface{}, 0, len(instance.MultiZoneInfo))
	for _, item := range instance.MultiZoneInfo {
		info := make(map[string]interface{}, 2)
		info["availability_zone"] = item.Zone
		info["subnet_id"] = item.SubnetId
		multiZoneInfos = append(multiZoneInfos, info)
	}
	_ = d.Set("multi_zone_infos", multiZoneInfos)

	nodeInfoList := make([]map[string]interface{}, 0, len(instance.NodeInfoList))
	for _, item := range instance.NodeInfoList {
		// this will not keep longer as long as cloud api response update
		if *item.Type == "kibana" {
			continue
		}
		info := make(map[string]interface{}, 5)
		info["node_num"] = item.NodeNum
		info["node_type"] = item.NodeType
		info["type"] = item.Type
		info["disk_type"] = item.DiskType
		info["disk_size"] = item.DiskSize
		info["encrypt"] = *item.DiskEncrypt > 0
		nodeInfoList = append(nodeInfoList, info)
	}
	_ = d.Set("node_info_list", nodeInfoList)

	if webInfo := instance.WebNodeTypeInfo; webInfo != nil {
		_ = helper.SetMapInterfaces(d, "web_node_type_info", map[string]interface{}{
			"node_type": webInfo.NodeType,
			"node_num":  webInfo.NodeNum,
		})
	}

	if instance.EsAcl != nil {
		esAcls := make([]map[string]interface{}, 0, 1)
		esAcl := map[string]interface{}{
			"black_list": instance.EsAcl.BlackIpList,
			"white_list": instance.EsAcl.WhiteIpList,
		}
		esAcls = append(esAcls, esAcl)
		_ = d.Set("es_acl", esAcls)
	}

	if len(instance.TagList) > 0 {
		tags := make(map[string]string)
		for _, tag := range instance.TagList {
			tags[*tag.TagKey] = *tag.TagValue
		}
		_ = d.Set("tags", tags)
	}

	return nil
}

func resourceTencentCloudElasticsearchInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_elasticsearch.update")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	instanceId := d.Id()
	elasticsearchService := ElasticsearchService{
		client: meta.(tccommon.ProviderMeta).GetAPIV3Conn(),
	}

	d.Partial(true)

	if d.HasChange("instance_name") {
		instanceName := d.Get("instance_name").(string)
		// Update operation support at most one item at the same time
		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			errRet := elasticsearchService.UpdateInstance(ctx, instanceId, instanceName, "", 0, nil, nil, nil)
			if errRet != nil {
				return tccommon.RetryError(errRet)
			}
			return nil
		})
		if err != nil {
			return err
		}
	}
	if d.HasChange("password") {
		password := d.Get("password").(string)
		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			errRet := elasticsearchService.UpdateInstance(ctx, instanceId, "", password, 0, nil, nil, nil)
			if errRet != nil {
				return tccommon.RetryError(errRet)
			}
			return nil
		})
		if err != nil {
			return err
		}
	}

	if d.HasChange("version") {
		version := d.Get("version").(string)
		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			errRet := elasticsearchService.UpdateInstanceVersion(ctx, instanceId, version)
			if errRet != nil {
				return tccommon.RetryError(errRet)
			}
			return nil
		})
		if err != nil {
			return err
		}

		err = tencentCloudElasticsearchInstanceUpgradeWaiting(ctx, &elasticsearchService, instanceId)
		if err != nil {
			return err
		}
	}

	if d.HasChange("license_type") {
		licenseType := d.Get("license_type").(string)
		err := resource.Retry(tccommon.WriteRetryTimeout*2, func() *resource.RetryError {
			errRet := elasticsearchService.UpdateInstanceLicense(ctx, instanceId, licenseType)
			if errRet != nil {
				return tccommon.RetryError(errRet)
			}
			return nil
		})
		if err != nil {
			return err
		}

		err = tencentCloudElasticsearchInstanceUpgradeWaiting(ctx, &elasticsearchService, instanceId)
		if err != nil {
			return err
		}
	}

	if d.HasChange("basic_security_type") {
		basicSecurityType := d.Get("basic_security_type").(int)
		licenseType := d.Get("license_type").(string)
		licenseTypeUpgrading := licenseType != "oss"
		err := resource.Retry(tccommon.WriteRetryTimeout*2, func() *resource.RetryError {
			errRet := elasticsearchService.UpdateInstance(ctx, instanceId, "", "", int64(basicSecurityType), nil, nil, nil)
			if errRet != nil {
				err := errRet.(*sdkErrors.TencentCloudSDKError)
				if err.Code == es.INVALIDPARAMETER && licenseTypeUpgrading {
					return resource.RetryableError(fmt.Errorf("waiting for licenseType update"))
				}
				return tccommon.RetryError(errRet)
			}

			return nil
		})
		if err != nil {
			return err
		}
		err = tencentCloudElasticsearchInstanceUpgradeWaiting(ctx, &elasticsearchService, instanceId)

		if err != nil {
			return err
		}
	}

	if d.HasChange("web_node_type_info") {
		var err error
		infos := d.Get("web_node_type_info").([]interface{})
		for _, item := range infos {
			value := item.(map[string]interface{})
			info := &es.WebNodeTypeInfo{
				NodeNum:  helper.IntUint64(value["node_num"].(int)),
				NodeType: helper.String(value["node_type"].(string)),
			}
			err = resource.Retry(tccommon.WriteRetryTimeout*2, func() *resource.RetryError {
				errRet := elasticsearchService.UpdateInstance(ctx, instanceId, "", "", 0, nil, info, nil)
				if errRet != nil {
					return tccommon.RetryError(errRet)
				}
				return nil
			})
			break
		}
		if err != nil {
			return err
		}
		err = tencentCloudElasticsearchInstanceUpgradeWaiting(ctx, &elasticsearchService, instanceId)
		if err != nil {
			return err
		}
	}

	if d.HasChange("node_info_list") {
		nodeInfos := d.Get("node_info_list").([]interface{})
		nodeInfoList := make([]*es.NodeInfo, 0, len(nodeInfos))
		for _, d := range nodeInfos {
			value := d.(map[string]interface{})
			nodeType := value["node_type"].(string)
			diskSize := uint64(value["disk_size"].(int))
			nodeNum := uint64(value["node_num"].(int))
			types := value["type"].(string)
			diskType := value["disk_type"].(string)
			encrypt := value["encrypt"].(bool)
			dataDisk := es.NodeInfo{
				NodeType:    &nodeType,
				DiskSize:    &diskSize,
				NodeNum:     &nodeNum,
				Type:        &types,
				DiskType:    &diskType,
				DiskEncrypt: helper.BoolToInt64Pointer(encrypt),
			}
			nodeInfoList = append(nodeInfoList, &dataDisk)
		}
		err := resource.Retry(tccommon.WriteRetryTimeout*2, func() *resource.RetryError {
			errRet := elasticsearchService.UpdateInstance(ctx, instanceId, "", "", 0, nodeInfoList, nil, nil)
			if errRet != nil {
				return tccommon.RetryError(errRet)
			}
			return nil
		})
		if err != nil {
			return err
		}
		err = tencentCloudElasticsearchInstanceUpgradeWaiting(ctx, &elasticsearchService, instanceId)
		if err != nil {
			return err
		}
	}

	if d.HasChange("tags") {
		oldInterface, newInterface := d.GetChange("tags")
		replaceTags, deleteTags := svctag.DiffTags(oldInterface.(map[string]interface{}), newInterface.(map[string]interface{}))
		tagService := svctag.NewTagService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())
		region := meta.(tccommon.ProviderMeta).GetAPIV3Conn().Region

		//internal version: replace null begin, please do not modify this annotation and refrain from inserting any code between the beginning and end lines of the annotation.
		resourceName := fmt.Sprintf("qcs::es:%s:uin/:instance/%s", region, instanceId)
		err := tagService.ModifyTags(ctx, resourceName, replaceTags, deleteTags)
		//internal version: replace null end, please do not modify this annotation and refrain from inserting any code between the beginning and end lines of the annotation.

		//internal version: replace waitTag begin, please do not modify this annotation and refrain from inserting any code between the beginning and end lines of the annotation.
		//internal version: replace waitTag end, please do not modify this annotation and refrain from inserting any code between the beginning and end lines of the annotation.

		if err != nil {
			return err
		}
	}
	if d.HasChange("es_acl") {
		esAcl := es.EsAcl{}
		if aclMap, ok := helper.InterfacesHeadMap(d, "es_acl"); ok {
			if v, ok := aclMap["black_list"]; ok {
				bList := v.(*schema.Set).List()
				for _, d := range bList {
					esAcl.BlackIpList = append(esAcl.BlackIpList, helper.String(d.(string)))
				}
			}
			if v, ok := aclMap["white_list"]; ok {
				wList := v.(*schema.Set).List()
				for _, d := range wList {
					esAcl.WhiteIpList = append(esAcl.WhiteIpList, helper.String(d.(string)))
				}
			}
		}

		err := resource.Retry(tccommon.WriteRetryTimeout*2, func() *resource.RetryError {
			errRet := elasticsearchService.UpdateInstance(ctx, instanceId, "", "", 0, nil, nil, &esAcl)
			if errRet != nil {
				return tccommon.RetryError(errRet)
			}
			return nil
		})
		if err != nil {
			return err
		}
	}

	d.Partial(false)

	return resourceTencentCloudElasticsearchInstanceRead(d, meta)
}

func resourceTencentCloudElasticsearchInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_elasticsearch.delete")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	instanceId := d.Id()
	elasticsearchService := ElasticsearchService{
		client: meta.(tccommon.ProviderMeta).GetAPIV3Conn(),
	}
	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		errRet := elasticsearchService.DeleteInstance(ctx, instanceId)
		if errRet != nil {
			return tccommon.RetryError(errRet)
		}
		return nil
	})
	if err != nil {
		return err
	}

	err = resource.Retry(2*tccommon.ReadRetryTimeout, func() *resource.RetryError {
		instance, errRet := elasticsearchService.DescribeInstanceById(ctx, instanceId)
		if errRet != nil {
			return tccommon.RetryError(errRet, tccommon.InternalError)
		}
		if instance == nil {
			return nil
		}
		return resource.RetryableError(fmt.Errorf("elasticsearch instance status is %d, retry...", *instance.Status))
	})
	if err != nil {
		return err
	}

	return nil
}

func tencentCloudElasticsearchInstanceUpgradeWaiting(ctx context.Context, service *ElasticsearchService, instanceId string) error {
	statusChangeRetries := 5
	return resource.Retry(10*tccommon.ReadRetryTimeout, func() *resource.RetryError {
		instance, errRet := service.DescribeInstanceById(ctx, instanceId)
		if errRet != nil {
			return tccommon.RetryError(errRet, tccommon.InternalError)
		}
		if instance == nil {
			return resource.NonRetryableError(fmt.Errorf("instance %s not exist", instanceId))
		}
		if *instance.Status == ES_INSTANCE_STATUS_NORMAL && statusChangeRetries > 0 {
			statusChangeRetries -= 1
			err := fmt.Errorf("instance %s waiting for upgrade status change, %d times remaining", instanceId, statusChangeRetries)
			return resource.RetryableError(err)
		}
		if *instance.Status == ES_INSTANCE_STATUS_PROCESSING {
			return resource.RetryableError(errors.New("elasticsearch instance status is processing, retry..."))
		}
		return nil
	})
}
