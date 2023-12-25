package dayuv2

import (
	"context"
	"fmt"
	"log"
	"reflect"
	"sort"
	"strings"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svcantiddos "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/antiddos"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	antiddos "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/antiddos/v20200309"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudAntiddosDdosGeoIpBlockConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudAntiddosDdosGeoIpBlockConfigCreate,
		Read:   resourceTencentCloudAntiddosDdosGeoIpBlockConfigRead,
		Update: resourceTencentCloudAntiddosDdosGeoIpBlockConfigUpdate,
		Delete: resourceTencentCloudAntiddosDdosGeoIpBlockConfigDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "InstanceId.",
			},

			"ddos_geo_ip_block_config": {
				Required:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "DDoS region blocking configuration, configuration ID cannot be empty when filling in parameters.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"region_type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Region type, value [oversea (overseas) China (domestic) customized (custom region)].",
						},
						"action": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Blocking action, value [drop (intercept) trans (release)].",
						},
						"area_list": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeInt,
							},
							Optional:    true,
							Description: "When RegionType is customized, an AreaList must be filled in, with a maximum of 128 entries;.",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudAntiddosDdosGeoIpBlockConfigCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_antiddos_ddos_geo_ip_block_config.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	var (
		request    = antiddos.NewCreateDDoSGeoIPBlockConfigRequest()
		instanceId string
		regionType string
		action     string
	)
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
		request.InstanceId = helper.String(instanceId)
	}

	areaList := make([]int64, 0)
	if dMap, ok := helper.InterfacesHeadMap(d, "ddos_geo_ip_block_config"); ok {
		dDoSGeoIPBlockConfig := antiddos.DDoSGeoIPBlockConfig{}
		if v, ok := dMap["region_type"]; ok {
			regionType = v.(string)
			dDoSGeoIPBlockConfig.RegionType = helper.String(regionType)
		}
		if v, ok := dMap["action"]; ok {
			action = v.(string)
			dDoSGeoIPBlockConfig.Action = helper.String(action)
		}

		if v, ok := dMap["area_list"]; ok {
			areaListSet := v.(*schema.Set).List()
			for i := range areaListSet {
				area := areaListSet[i].(int)
				areaList = append(areaList, int64(area))
				dDoSGeoIPBlockConfig.AreaList = append(dDoSGeoIPBlockConfig.AreaList, helper.IntInt64(area))
			}
		}
		request.DDoSGeoIPBlockConfig = &dDoSGeoIPBlockConfig
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseAntiddosClient().CreateDDoSGeoIPBlockConfig(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create antiddos ddosGeoIpBlockConfig failed, reason:%+v", logId, err)
		return err
	}

	service := svcantiddos.NewAntiddosService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())

	configList, err := service.DescribeAntiddosDdosGeoIpBlockConfigById(ctx, instanceId)
	if err != nil {
		return err
	}

	var targetConfig *antiddos.DDoSGeoIPBlockConfig
	for _, item := range configList {
		if *item.GeoIPBlockConfig.Action != action {
			continue
		}
		if *item.GeoIPBlockConfig.RegionType != regionType {
			continue
		}
		sort.Slice(areaList, func(i, j int) bool {
			return areaList[i] < areaList[j]
		})
		tmpAreaList := make([]int64, 0)
		for _, v := range item.GeoIPBlockConfig.AreaList {
			area := *v
			tmpAreaList = append(tmpAreaList, area)
		}
		sort.Slice(tmpAreaList, func(i, j int) bool {
			return tmpAreaList[i] < tmpAreaList[j]
		})
		if !reflect.DeepEqual(areaList, tmpAreaList) {
			continue
		}
		targetConfig = item.GeoIPBlockConfig
	}

	if targetConfig == nil {
		return fmt.Errorf("can not find geo ip block config")
	}
	d.SetId(instanceId + tccommon.FILED_SP + *targetConfig.Id)

	return resourceTencentCloudAntiddosDdosGeoIpBlockConfigRead(d, meta)
}

func resourceTencentCloudAntiddosDdosGeoIpBlockConfigRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_antiddos_ddos_geo_ip_block_config.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := svcantiddos.NewAntiddosService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", idSplit)
	}
	configList, err := service.DescribeAntiddosDdosGeoIpBlockConfigById(ctx, idSplit[0])
	if err != nil {
		return err
	}

	if len(configList) == 0 {
		d.SetId("")
		log.Printf("[WARN]%s resource `AntiddosDdosGeoIpBlockConfig` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	var targetConfig *antiddos.DDoSGeoIPBlockConfig
	for _, item := range configList {
		if *item.GeoIPBlockConfig.Id == idSplit[1] {
			targetConfig = item.GeoIPBlockConfig
			break
		}
	}

	if targetConfig != nil {
		_ = d.Set("instance_id", idSplit[0])

		if targetConfig != nil {
			ddoSGeoIPBlockConfigMap := map[string]interface{}{}

			if targetConfig.RegionType != nil {
				ddoSGeoIPBlockConfigMap["region_type"] = targetConfig.RegionType
			}

			if targetConfig.Action != nil {
				ddoSGeoIPBlockConfigMap["action"] = targetConfig.Action
			}

			if targetConfig.AreaList != nil {
				ddoSGeoIPBlockConfigMap["area_list"] = targetConfig.AreaList
			}

			_ = d.Set("ddos_geo_ip_block_config", []interface{}{ddoSGeoIPBlockConfigMap})
		}
	}

	return nil
}

func resourceTencentCloudAntiddosDdosGeoIpBlockConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_antiddos_ddos_geo_ip_block_config.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	request := antiddos.NewModifyDDoSGeoIPBlockConfigRequest()

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", idSplit)
	}

	request.InstanceId = helper.String(idSplit[0])

	immutableArgs := []string{"instance_id"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	if d.HasChange("ddos_geo_ip_block_config") {
		if dMap, ok := helper.InterfacesHeadMap(d, "ddos_geo_ip_block_config"); ok {
			dDoSGeoIPBlockConfig := antiddos.DDoSGeoIPBlockConfig{}
			if v, ok := dMap["region_type"]; ok {
				dDoSGeoIPBlockConfig.RegionType = helper.String(v.(string))
			}
			if v, ok := dMap["action"]; ok {
				dDoSGeoIPBlockConfig.Action = helper.String(v.(string))
			}
			dDoSGeoIPBlockConfig.Id = helper.String(idSplit[1])
			if v, ok := dMap["area_list"]; ok {
				areaListSet := v.(*schema.Set).List()
				for i := range areaListSet {
					areaList := areaListSet[i].(int)
					dDoSGeoIPBlockConfig.AreaList = append(dDoSGeoIPBlockConfig.AreaList, helper.IntInt64(areaList))
				}
			}
			request.DDoSGeoIPBlockConfig = &dDoSGeoIPBlockConfig
		}

		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseAntiddosClient().ModifyDDoSGeoIPBlockConfig(request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s update antiddos ddosGeoIpBlockConfig failed, reason:%+v", logId, err)
			return err
		}
	}

	return resourceTencentCloudAntiddosDdosGeoIpBlockConfigRead(d, meta)
}

func resourceTencentCloudAntiddosDdosGeoIpBlockConfigDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_antiddos_ddos_geo_ip_block_config.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := svcantiddos.NewAntiddosService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())
	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", idSplit)
	}

	configList, err := service.DescribeAntiddosDdosGeoIpBlockConfigById(ctx, idSplit[0])
	if err != nil {
		return err
	}

	if len(configList) == 0 {
		d.SetId("")
		log.Printf("[WARN]%s resource `AntiddosDdosGeoIpBlockConfig` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	var targetConfig *antiddos.DDoSGeoIPBlockConfig
	for _, item := range configList {
		if *item.GeoIPBlockConfig.Id == idSplit[1] {
			targetConfig = item.GeoIPBlockConfig
			break
		}
	}

	if err := service.DeleteAntiddosDdosGeoIpBlockConfigById(ctx, idSplit[0], targetConfig); err != nil {
		return err
	}

	return nil
}
