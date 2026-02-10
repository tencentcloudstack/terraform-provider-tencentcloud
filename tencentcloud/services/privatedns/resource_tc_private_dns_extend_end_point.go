package privatedns

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	privatednsIntlv20201028 "github.com/tencentcloud/tencentcloud-sdk-go-intl-en/tencentcloud/privatedns/v20201028"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudPrivateDnsExtendEndPoint() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudPrivateDnsExtendEndPointCreate,
		Read:   resourceTencentCloudPrivateDnsExtendEndPointRead,
		Delete: resourceTencentCloudPrivateDnsExtendEndPointDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"end_point_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Outbound endpoint name.",
			},
			"end_point_region": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The region of the outbound endpoint must be consistent with the region of the forwarding target VIP.",
			},
			"forward_ip": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				MaxItems:    1,
				Description: "Forwarding target.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"access_type": {
							Type:        schema.TypeString,
							Required:    true,
							ForceNew:    true,
							Description: "Forwarding target IP network access type. CLB: The forwarding IP is the internal CLB VIP. CCN: Forwarding IP through CCN routing.",
						},
						"host": {
							Type:          schema.TypeString,
							Optional:      true,
							ForceNew:      true,
							ConflictsWith: []string{"forward_ip.0.hosts"},
							Description:   "Forwarding target IP address.",
						},
						"hosts": {
							Type:          schema.TypeSet,
							Optional:      true,
							ForceNew:      true,
							MinItems:      2,
							ConflictsWith: []string{"forward_ip.0.host"},
							Description:   "Forwarding target IPs address.",
							Elem:          &schema.Schema{Type: schema.TypeString},
						},
						"port": {
							Type:     schema.TypeInt,
							Required: true,
							ForceNew: true,
							DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
								forwardIpList := d.Get("forward_ip").([]interface{})
								if len(forwardIpList) == 0 {
									return false
								}

								_, hostsExplicitlySet := d.GetOkExists("forward_ip.0.hosts")
								if hostsExplicitlySet {
									return true
								}

								return false
							},
							Description: "Specifies the forwarding IP port number. This only applies when configuring `host`; if using `hosts`, you can choose any port value from `hosts`.",
						},
						"vpc_id": {
							Type:        schema.TypeString,
							Required:    true,
							ForceNew:    true,
							Description: "Unique VPC ID.",
						},
						// "subnet_id": {
						// 	Type:        schema.TypeString,
						// 	Optional:    true,
						// 	Description: "Unique subnet ID. Required when the access type is CCN.",
						// },
						"access_gateway_id": {
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
							Description: "CCN id. Required when the access type is CCN.",
						},
						// computed
						"vip": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Specifies the forwarding target IP proxy IP.",
						},
						"vport": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Specifies the forwarding target IP proxy port.",
						},
						"proto": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Specifies the forwarding target IP protocol.",
						},
						"snat_vip_cidr": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The SNAT CIDR block of the outbound endpoint.",
						},
						"snat_vip_set": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The SNAT IP list of the outbound endpoint.",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudPrivateDnsExtendEndPointCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_private_dns_extend_end_point.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request    = privatednsIntlv20201028.NewCreateExtendEndpointRequest()
		response   = privatednsIntlv20201028.NewCreateExtendEndpointResponse()
		endPointId string
	)

	if v, ok := d.GetOk("end_point_name"); ok {
		request.EndpointName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("end_point_region"); ok {
		request.EndpointRegion = helper.String(v.(string))
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "forward_ip"); ok {
		forwardIp := &privatednsIntlv20201028.ForwardIp{}
		if v, ok := dMap["access_type"]; ok {
			forwardIp.AccessType = helper.String(v.(string))
		}

		if v, ok := dMap["host"].(string); ok && v != "" {
			forwardIp.Host = helper.String(v)
			forwardIp.IpNum = helper.IntInt64(1)
		}

		if v, ok := dMap["hosts"].(*schema.Set); ok && len(v.List()) > 0 {
			var tmpStr string
			for _, item := range v.List() {
				tmpStr += item.(string) + ";"
			}

			forwardIp.Host = helper.String(tmpStr[:len(tmpStr)-1])
			forwardIp.IpNum = helper.IntInt64(len(v.List()))
		}

		if v, ok := dMap["port"]; ok {
			forwardIp.Port = helper.IntInt64(v.(int))
		}

		if v, ok := dMap["vpc_id"]; ok {
			forwardIp.VpcId = helper.String(v.(string))
		}

		// if v, ok := dMap["subnet_id"]; ok {
		// 	forwardIp.SubnetId = helper.String(v.(string))
		// }

		if v, ok := dMap["access_gateway_id"]; ok {
			forwardIp.AccessGatewayId = helper.String(v.(string))
		}

		request.ForwardIp = forwardIp
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UsePrivatednsIntlV20201028Client().CreateExtendEndpointWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e, PRIVATEDNS_CUSTOM_RETRY_SDK_ERROR...)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Create private dns extend end point failed, Response is nil."))
		}

		response = result
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create private dns extend end point failed, reason:%+v", logId, err)
		return err
	}

	if response.Response.EndpointId == nil {
		return fmt.Errorf("EndpointId is nil.")
	}

	endPointId = *response.Response.EndpointId
	d.SetId(endPointId)
	return resourceTencentCloudPrivateDnsExtendEndPointRead(d, meta)
}

func resourceTencentCloudPrivateDnsExtendEndPointRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_private_dns_extend_end_point.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service    = PrivatednsService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		endPointId = d.Id()
	)

	respData, err := service.DescribePrivateDnsExtendEndPointById(ctx, endPointId)
	if err != nil {
		return err
	}

	if respData == nil {
		log.Printf("[WARN]%s resource `tencentcloud_private_dns_extend_end_point` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		d.SetId("")
		return nil
	}

	if err := resourceTencentCloudPrivateDnsExtendEndPointReadPreHandleResponse0(ctx, respData); err != nil {
		return err
	}

	return nil
}

func resourceTencentCloudPrivateDnsExtendEndPointDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_private_dns_extend_end_point.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request    = privatednsIntlv20201028.NewDeleteEndPointRequest()
		endPointId = d.Id()
	)

	request.EndPointId = helper.String(endPointId)
	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UsePrivatednsIntlV20201028Client().DeleteEndPointWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e, PRIVATEDNS_CUSTOM_RETRY_SDK_ERROR...)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s delete private dns extend end point failed, reason:%+v", logId, err)
		return err
	}

	return nil
}
