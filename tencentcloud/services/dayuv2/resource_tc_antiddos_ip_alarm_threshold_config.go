package dayuv2

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svcantiddos "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/antiddos"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	antiddos "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/antiddos/v20200309"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudAntiddosIpAlarmThresholdConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudAntiddosIpAlarmThresholdConfigCreate,
		Read:   resourceTencentCloudAntiddosIpAlarmThresholdConfigRead,
		Update: resourceTencentCloudAntiddosIpAlarmThresholdConfigUpdate,
		Delete: resourceTencentCloudAntiddosIpAlarmThresholdConfigDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"alarm_type": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Alarm threshold type, value [1 (incoming traffic alarm threshold) 2 (attack cleaning traffic alarm threshold)].",
			},
			"alarm_threshold": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Alarm threshold, in Mbps, with a value of&gt;=0; When used as an input parameter, setting 0 will delete the alarm threshold configuration;.",
			},

			"instance_ip": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Instance ip.",
			},
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Instance id.",
			},
		},
	}
}

func resourceTencentCloudAntiddosIpAlarmThresholdConfigCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_antiddos_ip_alarm_threshold_config.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	instanceId := d.Get("instance_id").(string)
	instanceIp := d.Get("instance_ip").(string)
	alarmType := d.Get("alarm_type").(int)
	alarmTypeString := strconv.Itoa(alarmType)
	d.SetId(instanceId + tccommon.FILED_SP + instanceIp + tccommon.FILED_SP + alarmTypeString)

	return resourceTencentCloudAntiddosIpAlarmThresholdConfigUpdate(d, meta)
}

func resourceTencentCloudAntiddosIpAlarmThresholdConfigRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_antiddos_ip_alarm_threshold_config.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := svcantiddos.NewAntiddosService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", idSplit)
	}
	alarmType, err := strconv.Atoi(idSplit[2])
	if err != nil {
		return err
	}
	ipAlarmThresholdConfig, err := service.DescribeAntiddosIpAlarmThresholdConfigById(ctx, idSplit[0], idSplit[1], alarmType)
	if err != nil {
		return err
	}

	if ipAlarmThresholdConfig == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `AntiddosIpAlarmThresholdConfig` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if ipAlarmThresholdConfig.AlarmThreshold != nil {
		_ = d.Set("alarm_threshold", ipAlarmThresholdConfig.AlarmThreshold)
	}
	_ = d.Set("instance_id", idSplit[0])
	_ = d.Set("instance_ip", idSplit[1])
	_ = d.Set("alarm_type", alarmType)

	return nil
}

func resourceTencentCloudAntiddosIpAlarmThresholdConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_antiddos_ip_alarm_threshold_config.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	request := antiddos.NewCreateIPAlarmThresholdConfigRequest()

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", idSplit)
	}

	instanceId := idSplit[0]
	instanceIp := idSplit[1]
	alarmType, err := strconv.Atoi(idSplit[2])
	if err != nil {
		return err
	}

	if d.HasChange("alarm_threshold") {

		ipAlarmThresholdRelation := antiddos.IPAlarmThresholdRelation{}
		ipAlarmThresholdRelation.AlarmType = helper.IntUint64(alarmType)

		if v, ok := d.GetOkExists("alarm_threshold"); ok {
			ipAlarmThresholdRelation.AlarmThreshold = helper.IntUint64(v.(int))
		}
		ipAlarmThresholdRelation.InstanceDetailList = []*antiddos.InstanceRelation{
			{
				EipList:    []*string{&instanceIp},
				InstanceId: &instanceId,
			},
		}

		request.IpAlarmThresholdConfigList = []*antiddos.IPAlarmThresholdRelation{&ipAlarmThresholdRelation}

		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseAntiddosClient().CreateIPAlarmThresholdConfig(request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s update antiddos ipAlarmThresholdConfig failed, reason:%+v", logId, err)
			return err
		}
	}

	return resourceTencentCloudAntiddosIpAlarmThresholdConfigRead(d, meta)
}

func resourceTencentCloudAntiddosIpAlarmThresholdConfigDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_antiddos_ip_alarm_threshold_config.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
