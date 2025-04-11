package emr

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	emr "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/emr/v20190103"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
	svctag "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/tag"
)

func ResourceTencentCloudServerlessHbaseInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudServerlessHbaseInstanceCreate,
		Read:   resourceTencentCloudServerlessHbaseInstanceRead,
		Update: resourceTencentCloudServerlessHbaseInstanceUpdate,
		Delete: resourceTencentCloudServerlessHbaseInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Instance name. Length limit is 6-36 characters. Only Chinese characters, letters, numbers, -, and _ are allowed.",
			},

			"pay_mode": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Instance pay mode. Value range: 0: indicates post-pay mode, that is, pay-as-you-go. 1: indicates pre-pay mode, that is, monthly subscription.",
			},

			"disk_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Instance disk type, Value range: CLOUD_HSSD: indicate performance cloud storage(ESSD). CLOUD_BSSD: indicate standard cloud storage(SSD).",
			},

			"disk_size": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Instance single-node disk capacity, in GB. The single-node disk capacity must be greater than or equal to 100 and less than or equal to 250 times the number of CPU cores. The capacity adjustment step is 100.",
			},

			"node_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Instance node type, can be filled in as 4C16G, 8C32G, 16C64G, 32C128G, case insensitive.",
			},

			"zone_settings": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "Detailed configuration of the instance availability zone, currently supports multiple availability zones, the number of availability zones can only be 1 or 3, including zone name, VPC information, and number of nodes. The total number of nodes across all zones must be greater than or equal to 3 and less than or equal to 50.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"zone": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The availability zone to which the instance belongs, such as ap-guangzhou-1.",
						},
						"vpc_settings": {
							Type:        schema.TypeList,
							Required:    true,
							MaxItems:    1,
							Description: "Private network related information configuration. This parameter can be used to specify the ID of the private network, subnet ID, and other information.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"vpc_id": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "VPC ID.",
									},
									"subnet_id": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Subnet ID.",
									},
								},
							},
						},
						"node_num": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Number of nodes.",
						},
					},
				},
			},

			"tags": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "List of tags to bind to the instance.",
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
							Description: "Tag value.",
						},
					},
				},
			},

			"time_span": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Time span.",
			},
			"time_unit": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Time unit, fill in m which means month.",
			},
			"auto_renew_flag": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "AutoRenewFlag, Value range: 0: indicates NOTIFY_AND_MANUAL_RENEW; 1: indicates NOTIFY_AND_AUTO_RENEW; 2: indicates DISABLE_NOTIFY_AND_MANUAL_RENEW.",
			},
		},
	}
}

func resourceTencentCloudServerlessHbaseInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_serverless_hbase_instance.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	var (
		request  = emr.NewCreateSLInstanceRequest()
		response = emr.NewCreateSLInstanceResponse()
	)

	if v, ok := d.GetOk("instance_name"); ok {
		request.InstanceName = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("pay_mode"); ok {
		request.PayMode = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("disk_type"); ok {
		request.DiskType = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("disk_size"); ok {
		request.DiskSize = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("node_type"); ok {
		request.NodeType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("zone_settings"); ok {
		for _, item := range v.([]interface{}) {
			zoneSettingsMap := item.(map[string]interface{})
			zoneSetting := emr.ZoneSetting{}
			if v, ok := zoneSettingsMap["zone"]; ok {
				zoneSetting.Zone = helper.String(v.(string))
			}
			if vPCSettingsMap, ok := helper.ConvertInterfacesHeadToMap(zoneSettingsMap["vpc_settings"]); ok {
				vPCSettings := emr.VPCSettings{}
				if v, ok := vPCSettingsMap["vpc_id"]; ok {
					vPCSettings.VpcId = helper.String(v.(string))
				}
				if v, ok := vPCSettingsMap["subnet_id"]; ok {
					vPCSettings.SubnetId = helper.String(v.(string))
				}
				zoneSetting.VPCSettings = &vPCSettings
			}
			if v, ok := zoneSettingsMap["node_num"]; ok {
				zoneSetting.NodeNum = helper.IntInt64(v.(int))
			}
			request.ZoneSettings = append(request.ZoneSettings, &zoneSetting)
		}
	}

	if v, ok := d.GetOk("tags"); ok {
		for idx, item := range v.(*schema.Set).List() {
			if item == nil {
				return fmt.Errorf("tags element with index %d is nil", idx+1)
			}
			tagsMap := item.(map[string]interface{})
			tag := emr.Tag{}
			if v, ok := tagsMap["tag_key"]; ok {
				tag.TagKey = helper.String(v.(string))
			}
			if v, ok := tagsMap["tag_value"]; ok {
				tag.TagValue = helper.String(v.(string))
			}
			request.Tags = append(request.Tags, &tag)
		}
	}

	var prePaySetting *emr.PrePaySetting
	if v, ok := d.GetOk("time_span"); ok {
		prePaySetting = &emr.PrePaySetting{}
		prePaySetting.Period = &emr.Period{}
		prePaySetting.Period.TimeSpan = helper.IntInt64(v.(int))
	}
	if v, ok := d.GetOk("time_unit"); ok {
		if prePaySetting == nil {
			prePaySetting = &emr.PrePaySetting{}
		}
		if prePaySetting.Period == nil {
			prePaySetting.Period = &emr.Period{}
		}
		prePaySetting.Period.TimeUnit = helper.String(v.(string))
	}
	if v, ok := d.GetOk("auto_renew_flag"); ok {
		if prePaySetting == nil {
			prePaySetting = &emr.PrePaySetting{}
		}
		prePaySetting.AutoRenewFlag = helper.IntInt64(v.(int))
	}
	if prePaySetting != nil {
		request.PrePaySetting = prePaySetting
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseEmrClient().CreateSLInstanceWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create serverless hbase instance failed, reason:%+v", logId, err)
		return err
	}

	instanceId := *response.Response.InstanceId
	d.SetId(instanceId)

	emrService := EMRService{
		client: meta.(tccommon.ProviderMeta).GetAPIV3Conn(),
	}
	conf := tccommon.BuildStateChangeConf([]string{}, []string{"2"}, 10*tccommon.ReadRetryTimeout, time.Second, emrService.SLInstanceStateRefreshFunc(instanceId, []string{}))
	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	return resourceTencentCloudServerlessHbaseInstanceRead(d, meta)
}

func resourceTencentCloudServerlessHbaseInstanceRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_serverless_hbase_instance.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	instanceId := d.Id()

	request := emr.NewDescribeSLInstanceRequest()
	response := emr.NewDescribeSLInstanceResponse()
	request.InstanceId = helper.String(instanceId)
	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseEmrClient().DescribeSLInstance(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update serverless hbase instance failed, reason:%+v", logId, err)
		return err
	}

	if response.Response == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `serverless_hbase_instance` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}
	respData := response.Response
	if respData.InstanceName != nil {
		_ = d.Set("instance_name", respData.InstanceName)
	}

	if respData.PayMode != nil {
		_ = d.Set("pay_mode", respData.PayMode)
	}

	if respData.DiskType != nil {
		_ = d.Set("disk_type", respData.DiskType)
	}

	if respData.DiskSize != nil {
		_ = d.Set("disk_size", respData.DiskSize)
	}

	zoneSettingsList := make([]map[string]interface{}, 0, len(respData.ZoneSettings))
	if respData.ZoneSettings != nil {
		for _, zoneSettings := range respData.ZoneSettings {
			zoneSettingsMap := map[string]interface{}{}

			if zoneSettings.Zone != nil {
				zoneSettingsMap["zone"] = zoneSettings.Zone
			}

			vPCSettingsMap := map[string]interface{}{}

			if zoneSettings.VPCSettings != nil {
				if zoneSettings.VPCSettings.VpcId != nil {
					vPCSettingsMap["vpc_id"] = zoneSettings.VPCSettings.VpcId
				}

				if zoneSettings.VPCSettings.SubnetId != nil {
					vPCSettingsMap["subnet_id"] = zoneSettings.VPCSettings.SubnetId
				}

				zoneSettingsMap["vpc_settings"] = []interface{}{vPCSettingsMap}
			}

			if zoneSettings.NodeNum != nil {
				zoneSettingsMap["node_num"] = zoneSettings.NodeNum
			}

			zoneSettingsList = append(zoneSettingsList, zoneSettingsMap)
		}

		_ = d.Set("zone_settings", zoneSettingsList)
	}

	tagsList := make([]map[string]interface{}, 0, len(respData.Tags))
	if respData.Tags != nil {
		for _, tags := range respData.Tags {
			tagsMap := map[string]interface{}{}

			if tags.TagKey != nil {
				tagsMap["tag_key"] = tags.TagKey
			}

			if tags.TagValue != nil {
				tagsMap["tag_value"] = tags.TagValue
			}

			tagsList = append(tagsList, tagsMap)
		}

		_ = d.Set("tags", tagsList)
	}
	if respData.AutoRenewFlag != nil {
		_ = d.Set("auto_renew_flag", respData.AutoRenewFlag)
	}

	if respData.NodeType != nil {
		_ = d.Set("node_type", respData.NodeType)
	}
	_ = instanceId
	return nil
}

func resourceTencentCloudServerlessHbaseInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_serverless_hbase_instance.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	immutableArgs := []string{"pay_mode", "disk_type", "disk_size", "node_type", "time_span", "time_unit", "auto_renew_flag"}
	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}
	instanceId := d.Id()

	if d.HasChange("instance_name") {
		if v, ok := d.GetOk("instance_name"); ok {
			request := emr.NewModifySLInstanceBasicRequest()
			request.InstanceId = helper.String(instanceId)
			request.ClusterName = helper.String(v.(string))
			err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
				result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseEmrClient().ModifySLInstanceBasic(request)
				if e != nil {
					return tccommon.RetryError(e)
				} else {
					log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
				}
				return nil
			})
			if err != nil {
				log.Printf("[CRITAL]%s update serverless hbase instance_name failed, reason:%+v", logId, err)
				return err
			}
		}
	}
	if d.HasChange("tags") {

		oldValue, newValue := d.GetChange("tags")
		oldMap := make(map[string]interface{})
		newMap := make(map[string]interface{})

		for _, o := range oldValue.(*schema.Set).List() {
			oMap := o.(map[string]interface{})
			oldMap[oMap["tag_key"].(string)] = oMap["tag_value"].(string)
		}
		for _, n := range newValue.(*schema.Set).List() {
			nMap := n.(map[string]interface{})
			newMap[nMap["tag_key"].(string)] = nMap["tag_value"].(string)
		}

		replaceTags, deleteTags := svctag.DiffTags(oldMap, newMap)

		tcClient := meta.(tccommon.ProviderMeta).GetAPIV3Conn()
		tagService := svctag.NewTagService(tcClient)
		resourceName := tccommon.BuildTagResourceName("emr", "emr-serverless-instance", tcClient.Region, d.Id())
		err := tagService.ModifyTags(ctx, resourceName, replaceTags, deleteTags)
		if err != nil {
			return err
		}
	}
	needChange := false
	if d.HasChange("zone_settings") {
		for idx, zoneSetting := range d.Get("zone_settings").([]interface{}) {
			for k := range zoneSetting.(map[string]interface{}) {
				param := fmt.Sprintf("zone_settings.%d.%s", idx, k)
				if d.HasChange(param) {
					if k == "node_num" {
						needChange = true
					} else {
						return fmt.Errorf("argument `%s` cannot be changed", param)
					}
				}
			}
		}
	}
	if needChange {
		for idx, zoneSetting := range d.Get("zone_settings").([]interface{}) {
			param := fmt.Sprintf("zone_settings.%d.node_num", idx)
			if !d.HasChange(param) {
				continue
			}
			zoneSettingMap := zoneSetting.(map[string]interface{})
			request := emr.NewModifySLInstanceRequest()

			request.InstanceId = helper.String(instanceId)

			if v, ok := zoneSettingMap["zone"]; ok {
				request.Zone = helper.String(v.(string))
			}

			if v, ok := zoneSettingMap["node_num"]; ok {
				request.NodeNum = helper.IntInt64(v.(int))
			}

			err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
				result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseEmrClient().ModifySLInstanceWithContext(ctx, request)
				if e != nil {
					return tccommon.RetryError(e)
				} else {
					log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
				}
				return nil
			})
			if err != nil {
				log.Printf("[CRITAL]%s update serverless hbase instance failed, reason:%+v", logId, err)
				return err
			}

			emrService := EMRService{
				client: meta.(tccommon.ProviderMeta).GetAPIV3Conn(),
			}
			conf := tccommon.BuildStateChangeConf([]string{}, []string{"2"}, 10*tccommon.ReadRetryTimeout, time.Second, emrService.SLInstanceStateRefreshFunc(instanceId, []string{}))
			if _, e := conf.WaitForState(); e != nil {
				return e
			}
		}
	}

	_ = instanceId
	return resourceTencentCloudServerlessHbaseInstanceRead(d, meta)
}

func resourceTencentCloudServerlessHbaseInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_serverless_hbase_instance.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	instanceId := d.Id()

	var (
		request  = emr.NewTerminateSLInstanceRequest()
		response = emr.NewTerminateSLInstanceResponse()
	)

	request.InstanceId = helper.String(instanceId)

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseEmrClient().TerminateSLInstanceWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s delete serverless hbase instance failed, reason:%+v", logId, err)
		return err
	}

	emrService := EMRService{
		client: meta.(tccommon.ProviderMeta).GetAPIV3Conn(),
	}

	if d.Get("pay_mode").(int) == 1 {
		conf := tccommon.BuildStateChangeConf([]string{}, []string{"201"}, 10*tccommon.ReadRetryTimeout, time.Second, emrService.SLInstanceStateRefreshFunc(instanceId, []string{}))
		if _, e := conf.WaitForState(); e != nil {
			return e
		}

		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseEmrClient().TerminateSLInstanceWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}
			response = result
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s delete serverless hbase instance failed, reason:%+v", logId, err)
			return err
		}
	}

	conf := tccommon.BuildStateChangeConf([]string{}, []string{"-2"}, 10*tccommon.ReadRetryTimeout, time.Second, emrService.SLInstanceStateRefreshFunc(instanceId, []string{}))
	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	_ = response
	_ = instanceId
	return nil
}
