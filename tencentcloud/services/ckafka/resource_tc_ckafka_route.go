package ckafka

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ckafka "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ckafka/v20190819"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudCkafkaRoute() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCkafkaRouteCreate,
		Read:   resourceTencentCloudCkafkaRouteRead,
		Update: resourceTencentCloudCkafkaRouteUpdate,
		Delete: resourceTencentCloudCkafkaRouteDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Instance id.",
			},

			"vip_type": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "Routing network type (3:vpc routing; 4: standard support routing; 7: professional support routing).",
			},

			"vpc_id": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Vpc id.",
			},

			"subnet_id": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Subnet id.",
			},

			"access_type": {
				Optional: true,
				Computed: true,
				Type:     schema.TypeInt,
				Description: "Access type. Valid values:\n" +
					"- 0: PLAINTEXT (in clear text, supported by both the old version and the community version without user information)\n" +
					"- 1: SASL_PLAINTEXT (in clear text, but at the beginning of the data, authentication will be logged in through SASL, which is only supported by the community version)\n" +
					"- 2: SSL (SSL encrypted communication without user information, supported by both older and community versions)\n" +
					"- 3: SASL_SSL (SSL encrypted communication. When the data starts, authentication will be logged in through SASL. Only the community version supports it).",
			},

			"auth_flag": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Auth flag.",
			},

			"caller_appid": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Caller appid.",
			},

			"public_network": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Public network.",
			},

			"ip": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Ip.",
			},
			"vip_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Virtual IP list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"vip": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Virtual IP.",
						},
						"vport": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Virtual port.",
						},
					},
				},
			},
			"broker_vip_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Virtual IP list (1 to 1 broker nodes).",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"vip": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Virtual IP.",
						},
						"vport": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Virtual port.",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudCkafkaRouteCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ckafka_route.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		request    = ckafka.NewCreateRouteRequest()
		response   = ckafka.NewCreateRouteResponse()
		instanceId string
	)
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
		request.InstanceId = helper.String(instanceId)
	}

	if v, ok := d.GetOkExists("vip_type"); ok {
		request.VipType = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("vpc_id"); ok {
		request.VpcId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("subnet_id"); ok {
		request.SubnetId = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("access_type"); ok {
		request.AccessType = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("auth_flag"); ok {
		request.AuthFlag = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("caller_appid"); ok {
		request.CallerAppid = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("public_network"); ok {
		request.PublicNetwork = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("ip"); ok {
		request.Ip = helper.String(v.(string))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCkafkaClient().CreateRoute(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create ckafka route failed, reason:%+v", logId, err)
		return err
	}
	log.Printf("Result.Data: %+v", response.Response.Result.Data)
	routeIdInt64 := *response.Response.Result.Data.RouteDTO.RouteId
	flowIdInt64 := *response.Response.Result.Data.FlowId
	reouteIdString := strconv.FormatInt(routeIdInt64, 10)
	d.SetId(instanceId + tccommon.FILED_SP + reouteIdString)

	service := CkafkaService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	conf := tccommon.BuildStateChangeConf([]string{}, []string{"0"}, 1*tccommon.ReadRetryTimeout, time.Second, service.CkafkaRouteStateRefreshFunc(flowIdInt64, []string{}))

	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	return resourceTencentCloudCkafkaRouteRead(d, meta)
}

func resourceTencentCloudCkafkaRouteRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ckafka_route.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := CkafkaService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	items := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(items) < 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	instanceId := items[0]
	routeId := items[1]
	routeIdInt64, err := strconv.ParseInt(routeId, 10, 64)
	if err != nil {
		return err
	}
	route, err := service.DescribeCkafkaRouteById(ctx, instanceId, routeIdInt64)
	if err != nil {
		return err
	}
	_ = d.Set("instance_id", instanceId)
	if route == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `CkafkaRoute` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if route.VipType != nil {
		_ = d.Set("vip_type", route.VipType)
	}

	if route.VpcId != nil {
		_ = d.Set("vpc_id", route.VpcId)
	}

	if route.Subnet != nil {
		_ = d.Set("subnet_id", route.Subnet)
	}

	if route.AccessType != nil {
		_ = d.Set("access_type", route.AccessType)
	}

	if len(route.VipList) > 0 {
		_ = d.Set("ip", route.VipList[0].Vip)
	}
	vipList := make([]map[string]interface{}, 0)
	for _, vip := range route.VipList {
		vipList = append(vipList, map[string]interface{}{
			"vip":   vip.Vip,
			"vport": vip.Vport,
		})
	}
	_ = d.Set("vip_list", vipList)
	brokerVipList := make([]map[string]interface{}, 0)
	for _, brokerVip := range route.BrokerVipList {
		brokerVipList = append(brokerVipList, map[string]interface{}{
			"vip":   brokerVip.Vip,
			"vport": brokerVip.Vport,
		})
	}
	_ = d.Set("broker_vip_list", brokerVipList)

	return nil
}

func resourceTencentCloudCkafkaRouteUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ckafka_route.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	immutableArgs := []string{"instance_id", "vip_type", "vpc_id", "subnet_id", "access_type", "auth_flag", "caller_appid", "public_network", "ip"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}
	return resourceTencentCloudCkafkaRouteRead(d, meta)
}

func resourceTencentCloudCkafkaRouteDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ckafka_route.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := CkafkaService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	items := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(items) < 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	instanceId := items[0]
	routeId := items[1]
	routeIdInt64, err := strconv.ParseInt(routeId, 10, 64)
	if err != nil {
		return err
	}
	if err := service.DeleteCkafkaRouteById(ctx, instanceId, routeIdInt64); err != nil {
		return err
	}

	return nil
}
