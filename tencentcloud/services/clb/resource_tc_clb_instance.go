package clb

import (
	"context"
	"fmt"
	"log"
	"sync"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svcas "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/as"
	svctag "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/tag"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/pkg/errors"
	clb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/clb/v20180317"
	sdkErrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

var clbActionMu = &sync.Mutex{}

func ResourceTencentCloudClbInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudClbInstanceCreate,
		Read:   resourceTencentCloudClbInstanceRead,
		Update: resourceTencentCloudClbInstanceUpdate,
		Delete: resourceTencentCloudClbInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"network_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: tccommon.ValidateAllowedStringValue(CLB_NETWORK_TYPE),
				Description:  "Type of CLB instance. Valid values: `OPEN` and `INTERNAL`.",
			},
			"clb_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: tccommon.ValidateStringLengthInRange(1, 60),
				Description:  "Name of the CLB. The name can only contain Chinese characters, English letters, numbers, underscore and hyphen '-'.",
			},
			"clb_vips": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "The virtual service address table of the CLB.",
			},
			"project_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     0,
				Description: "ID of the project within the CLB instance, `0` - Default Project.",
			},
			"vpc_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Computed:    true,
				Description: "VPC ID of the CLB.",
			},
			"subnet_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: tccommon.ValidateStringLengthInRange(2, 60),
				Description:  "In the case of purchasing a `INTERNAL` clb instance, the subnet id must be specified. The VIP of the `INTERNAL` clb instance will be generated from this subnet.",
			},
			"cluster_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Cluster ID.",
			},
			"address_ip_version": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "It's only applicable to public network CLB instances. IP version. Values: `IPV4`, `IPV6` and `IPv6FullChain` (case-insensitive). Default: `IPV4`. Note: IPV6 indicates IPv6 NAT64, while IPv6FullChain indicates IPv6.",
			},
			"internet_charge_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Internet charge type, only applicable to open CLB. Valid values are `TRAFFIC_POSTPAID_BY_HOUR`, `BANDWIDTH_POSTPAID_BY_HOUR` and `BANDWIDTH_PACKAGE`.",
			},
			"delete_protect": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether to enable delete protection.",
			},
			"bandwidth_package_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Bandwidth package id. If set, the `internet_charge_type` must be `BANDWIDTH_PACKAGE`.",
			},
			"internet_bandwidth_max_out": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Max bandwidth out, only applicable to open CLB. Valid value ranges is [1, 2048]. Unit is Mbps.",
			},
			"security_groups": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Security groups of the CLB instance. Supports both `OPEN` and `INTERNAL` CLBs.",
			},
			"target_region_info_region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Region information of backend services are attached the CLB instance. Only supports `OPEN` CLBs.",
			},
			"target_region_info_vpc_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Vpc information of backend services are attached the CLB instance. Only supports `OPEN` CLBs.",
			},
			"snat_pro": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Indicates whether Binding IPs of other VPCs feature switch.",
			},
			"snat_ips": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Snat Ip List, required with `snat_pro=true`. NOTE: This argument cannot be read and modified here because dynamic ip is untraceable, please import resource `tencentcloud_clb_snat_ip` to handle fixed ips.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ip": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Snat IP address, If set to empty will auto allocated.",
						},
						"subnet_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Snat subnet ID.",
						},
					},
				},
			},
			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "The available tags within this CLB.",
			},
			"sla_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				Description: "This parameter is required to create LCU-supported instances. Values:" +
					"`SLA`: Super Large 4. When you have activated Super Large models, `SLA` refers to Super Large 4; " +
					"`clb.c2.medium`: Standard; " +
					"`clb.c3.small`: Advanced 1; " +
					"`clb.c3.medium`: Advanced 1; " +
					"`clb.c4.small`: Super Large 1; " +
					"`clb.c4.medium`: Super Large 2; " +
					"`clb.c4.large`: Super Large 3; " +
					"`clb.c4.xlarge`: Super Large 4. " +
					"For more details, see [Instance Specifications](https://intl.cloud.tencent.com/document/product/214/84689?from_cn_redirect=1).",
			},
			"vip_isp": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Computed:    true,
				Description: "Network operator, only applicable to open CLB. Valid values are `CMCC`(China Mobile), `CTCC`(Telecom), `CUCC`(China Unicom) and `BGP`. If this ISP is specified, network billing method can only use the bandwidth package billing (BANDWIDTH_PACKAGE).",
			},
			"vip": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Computed:    true,
				Description: "Specifies the VIP for the application of a CLB instance. This parameter is optional. If you do not specify this parameter, the system automatically assigns a value for the parameter. IPv4 and IPv6 CLB instances support this parameter, but IPv6 NAT64 CLB instances do not.",
			},
			"load_balancer_pass_to_target": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Whether the target allow flow come from clb. If value is true, only check security group of clb, or check both clb and backend instance security group.",
			},
			"master_zone_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Setting master zone id of cross available zone disaster recovery, only applicable to open CLB.",
			},
			"zone_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Available zone id, only applicable to open CLB.",
			},
			"slave_zone_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Setting slave zone id of cross available zone disaster recovery, only applicable to open CLB. this zone will undertake traffic when the master is down.",
			},
			"log_set_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The id of log set.",
			},
			"log_topic_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The id of log topic.",
			},
			"dynamic_vip": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "If create dynamic vip CLB instance, `true` or `false`.",
			},
			"eip_address_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The unique ID of the EIP, such as eip-1v2rmbwk, is only applicable to the intranet load balancing binding EIP. During the EIP change, there may be a brief network interruption.",
			},
			"associate_endpoint": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The associated terminal node ID; passing an empty string indicates unassociating the node.",
			},
			"domain": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Domain name of the CLB instance.",
			},
			"ipv6_mode": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "This field is meaningful when the IP address version is ipv6, `IPv6Nat64` | `IPv6FullChain`.",
			},
			"address_ipv6": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The IPv6 address of the load balancing instance.",
			},
		},
	}
}

func resourceTencentCloudClbInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_clb_instance.create")()

	clbActionMu.Lock()
	defer clbActionMu.Unlock()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	networkType := d.Get("network_type").(string)
	clbName := d.Get("clb_name").(string)
	flag, e := checkSameName(clbName, meta)
	if e != nil {
		return e
	}

	if flag {
		return fmt.Errorf("[CHECK][CLB instance][Create] check: Same CLB name %s exists!", clbName)
	}

	targetRegionInfoRegion := ""
	targetRegionInfoVpcId := ""
	if v, ok := d.GetOk("target_region_info_region"); ok {
		targetRegionInfoRegion = v.(string)
		if networkType == CLB_NETWORK_TYPE_INTERNAL {
			return fmt.Errorf("[CHECK][CLB instance][Create] check: INTERNAL network_type do not support this operation with target_region_info")
		}
	}

	if v, ok := d.GetOk("target_region_info_vpc_id"); ok {
		targetRegionInfoVpcId = v.(string)
		if networkType == CLB_NETWORK_TYPE_INTERNAL {
			return fmt.Errorf("[CHECK][CLB instance][Create] check: INTERNAL network_type do not support this operation with target_region_info")
		}
	}

	if (targetRegionInfoRegion != "" && targetRegionInfoVpcId == "") || (targetRegionInfoRegion == "" && targetRegionInfoVpcId != "") {
		return fmt.Errorf("[CHECK][CLB instance][Create] check: region and vpc_id must be set at same time")
	}

	request := clb.NewCreateLoadBalancerRequest()
	request.LoadBalancerType = helper.String(networkType)
	request.LoadBalancerName = helper.String(clbName)
	if v, ok := d.GetOk("vpc_id"); ok {
		request.VpcId = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("project_id"); ok {
		projectId := int64(v.(int))
		request.ProjectId = &projectId
	}

	if v, ok := d.GetOk("subnet_id"); ok {
		request.SubnetId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("cluster_id"); ok {
		request.ClusterIds = []*string{helper.String(v.(string))}
	}

	//vip_isp
	if v, ok := d.GetOk("vip_isp"); ok {
		if networkType == CLB_NETWORK_TYPE_INTERNAL {
			return fmt.Errorf("[CHECK][CLB instance][Create] check: INTERNAL network_type do not support vip ISP setting")
		}

		request.VipIsp = helper.String(v.(string))
	}

	//vip
	if v, ok := d.GetOk("vip"); ok {
		request.Vip = helper.String(v.(string))
	}

	//SlaType
	if v, ok := d.GetOk("sla_type"); ok {
		request.SlaType = helper.String(v.(string))
	}

	//ip version
	if v, ok := d.GetOk("address_ip_version"); ok {
		if networkType == CLB_NETWORK_TYPE_INTERNAL {
			return fmt.Errorf("[CHECK][CLB instance][Create] check: INTERNAL network_type do not support IP version setting")
		}

		request.AddressIPVersion = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("snat_pro"); ok {
		request.SnatPro = helper.Bool(v.(bool))
	}

	if v, ok := d.Get("snat_ips").([]interface{}); ok && len(v) > 0 {
		for i := range v {
			item := v[i].(map[string]interface{})
			snatIp := &clb.SnatIp{}
			if v, ok := item["subnet_id"].(string); ok && v != "" {
				snatIp.SubnetId = &v
			}

			if v, ok := item["ip"].(string); ok && v != "" {
				snatIp.Ip = &v
			}

			request.SnatIps = append(request.SnatIps, snatIp)
		}
	}

	v, ok := d.GetOk("internet_charge_type")
	bv, bok := d.GetOk("internet_bandwidth_max_out")
	pv, pok := d.GetOk("bandwidth_package_id")

	chargeType := v.(string)

	//internet charge type
	if ok {
		if networkType == CLB_NETWORK_TYPE_INTERNAL {
			return fmt.Errorf("[CHECK][CLB instance][Create] check: INTERNAL network_type do not support internet charge type setting")
		}

		request.InternetAccessible = &clb.InternetAccessible{}
		if ok {
			request.InternetAccessible.InternetChargeType = helper.String(chargeType)
		}

		if pok {
			if chargeType != svcas.INTERNET_CHARGE_TYPE_BANDWIDTH_PACKAGE {
				return fmt.Errorf("[CHECK][CLB instance][Create] check: internet_charge_type must `BANDWIDTH_PACKAGE` when bandwidth_package_id was set")
			}

			request.BandwidthPackageId = helper.String(pv.(string))
		}
	}

	// open or internal
	if bok {
		request.InternetAccessible.InternetMaxBandwidthOut = helper.IntInt64(bv.(int))
	}

	if v, ok := d.GetOk("master_zone_id"); ok {
		if networkType == CLB_NETWORK_TYPE_INTERNAL {
			return fmt.Errorf("[CHECK][CLB instance][Create] check: INTERNAL network_type do not support master zone id setting")
		}

		request.MasterZoneId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("zone_id"); ok {
		if networkType == CLB_NETWORK_TYPE_INTERNAL {
			return fmt.Errorf("[CHECK][CLB instance][Create] check: INTERNAL network_type do not support zone id setting")
		}

		request.ZoneId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("slave_zone_id"); ok {
		if networkType == CLB_NETWORK_TYPE_INTERNAL {
			return fmt.Errorf("[CHECK][CLB instance][Create] check: INTERNAL network_type do not support slave zone id setting")
		}

		request.SlaveZoneId = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("load_balancer_pass_to_target"); ok {
		request.LoadBalancerPassToTarget = helper.Bool(v.(bool))
	}

	if v, ok := d.GetOkExists("dynamic_vip"); ok {
		request.DynamicVip = helper.Bool(v.(bool))
	}

	if v, ok := d.GetOk("eip_address_id"); ok {
		request.EipAddressId = helper.String(v.(string))
	}

	if tags := helper.GetTags(d, "tags"); len(tags) > 0 {
		for k, v := range tags {
			tmpKey := k
			tmpValue := v
			request.Tags = append(request.Tags, &clb.TagInfo{
				TagKey:   &tmpKey,
				TagValue: &tmpValue,
			})
		}
	}

	var response *clb.CreateLoadBalancerResponse
	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseClbClient().CreateLoadBalancer(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil || result.Response.RequestId == nil {
			return resource.NonRetryableError(fmt.Errorf("Create CLB instance failed, Response is nil."))
		}

		response = result
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create CLB instance failed, reason:%+v", logId, err)
		return err
	}

	// wait
	requestId := *response.Response.RequestId
	clbId, err := waitForTaskFinishGetID(requestId, meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseClbClient())
	if err != nil {
		return err
	}

	if clbId == "" {
		return fmt.Errorf("[CHECK][CLB instance][Create] check: response error, load balancer id is nil")
	}

	d.SetId(clbId)

	if v, ok := d.GetOk("security_groups"); ok {
		sgRequest := clb.NewSetLoadBalancerSecurityGroupsRequest()
		sgRequest.LoadBalancerId = helper.String(clbId)
		securityGroups := v.([]interface{})
		sgRequest.SecurityGroups = make([]*string, 0, len(securityGroups))
		for i := range securityGroups {
			if securityGroups[i] != nil {
				if securityGroup, ok := securityGroups[i].(string); ok && securityGroup != "" {
					sgRequest.SecurityGroups = append(sgRequest.SecurityGroups, &securityGroup)
				}
			}
		}

		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			sgResponse, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseClbClient().SetLoadBalancerSecurityGroups(sgRequest)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
					logId, sgRequest.GetAction(), sgRequest.ToJsonString(), sgResponse.ToJsonString())
				requestId := *sgResponse.Response.RequestId
				retryErr := waitForTaskFinish(requestId, meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseClbClient())
				if retryErr != nil {
					return tccommon.RetryError(errors.WithStack(retryErr))
				}
			}

			return nil
		})

		if err != nil {
			log.Printf("[CRITAL]%s create CLB instance security_groups failed, reason:%+v", logId, err)
			return err
		}
	}

	if v, ok := d.GetOk("log_set_id"); ok {
		if u, ok := d.GetOk("log_topic_id"); ok {
			logRequest := clb.NewSetLoadBalancerClsLogRequest()
			logRequest.LoadBalancerId = helper.String(clbId)
			logRequest.LogSetId = helper.String(v.(string))
			logRequest.LogTopicId = helper.String(u.(string))
			err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
				logResponse, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseClbClient().SetLoadBalancerClsLog(logRequest)
				if e != nil {
					return tccommon.RetryError(e)
				} else {
					log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
						logId, logRequest.GetAction(), logRequest.ToJsonString(), logResponse.ToJsonString())
					requestId := *logResponse.Response.RequestId
					retryErr := waitForTaskFinish(requestId, meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseClbClient())
					if retryErr != nil {
						return tccommon.RetryError(errors.WithStack(retryErr))
					}
				}

				return nil
			})

			if err != nil {
				log.Printf("[CRITAL]%s set CLB instance log failed, reason:%+v", logId, err)
				return err
			}

		} else {
			return fmt.Errorf("log_topic_id and log_set_id must be set together.")
		}
	}

	if targetRegionInfoRegion != "" {
		isLoadBalancePassToTgt := d.Get("load_balancer_pass_to_target").(bool)
		targetRegionInfo := clb.TargetRegionInfo{
			Region: &targetRegionInfoRegion,
			VpcId:  &targetRegionInfoVpcId,
		}

		mRequest := clb.NewModifyLoadBalancerAttributesRequest()
		mRequest.LoadBalancerId = helper.String(clbId)
		mRequest.TargetRegionInfo = &targetRegionInfo
		mRequest.LoadBalancerPassToTarget = &isLoadBalancePassToTgt
		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			mResponse, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseClbClient().ModifyLoadBalancerAttributes(mRequest)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
					logId, mRequest.GetAction(), mRequest.ToJsonString(), mResponse.ToJsonString())
				requestId := *mResponse.Response.RequestId
				retryErr := waitForTaskFinish(requestId, meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseClbClient())
				if retryErr != nil {
					return tccommon.RetryError(errors.WithStack(retryErr))
				}
			}

			return nil
		})

		if err != nil {
			log.Printf("[CRITAL]%s create CLB instance failed, reason:%+v", logId, err)
			return err
		}
	}

	if v, ok := d.GetOkExists("delete_protect"); ok {
		isDeleteProect := v.(bool)
		if isDeleteProect {
			mRequest := clb.NewModifyLoadBalancerAttributesRequest()
			mRequest.LoadBalancerId = helper.String(clbId)
			mRequest.DeleteProtect = &isDeleteProect
			err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
				mResponse, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseClbClient().ModifyLoadBalancerAttributes(mRequest)
				if e != nil {
					return tccommon.RetryError(e)
				} else {
					log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
						logId, mRequest.GetAction(), mRequest.ToJsonString(), mResponse.ToJsonString())
					requestId := *mResponse.Response.RequestId
					retryErr := waitForTaskFinish(requestId, meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseClbClient())
					if retryErr != nil {
						return tccommon.RetryError(errors.WithStack(retryErr))
					}
				}

				return nil
			})

			if err != nil {
				log.Printf("[CRITAL]%s create CLB instance failed, reason:%+v", logId, err)
				return err
			}
		}
	}

	if v, ok := d.GetOkExists("associate_endpoint"); ok {
		endpointId := v.(string)
		if endpointId != "" {
			mRequest := clb.NewModifyLoadBalancerAttributesRequest()
			mRequest.LoadBalancerId = helper.String(clbId)
			mRequest.AssociateEndpoint = &endpointId
			err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
				mResponse, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseClbClient().ModifyLoadBalancerAttributes(mRequest)
				if e != nil {
					return tccommon.RetryError(e)
				} else {
					log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, mRequest.GetAction(), mRequest.ToJsonString(), mResponse.ToJsonString())
					if mResponse == nil || mResponse.Response == nil || mResponse.Response.RequestId == nil {
						return resource.NonRetryableError(fmt.Errorf("Modify load balancer attributes failed, Response is nil."))
					}

					requestId := *mResponse.Response.RequestId
					retryErr := waitForTaskFinish(requestId, meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseClbClient())
					if retryErr != nil {
						return tccommon.RetryError(errors.WithStack(retryErr))
					}
				}

				return nil
			})

			if err != nil {
				log.Printf("[CRITAL]%s create CLB instance failed, reason:%+v", logId, err)
				return err
			}
		}
	}

	return resourceTencentCloudClbInstanceRead(d, meta)
}

func resourceTencentCloudClbInstanceRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_clb_instance.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		clbService = ClbService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		instance   *clb.LoadBalancer
		clbId      = d.Id()
	)

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := clbService.DescribeLoadBalancerById(ctx, clbId)
		if e != nil {
			return tccommon.RetryError(e)
		}

		instance = result
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s read CLB instance failed, reason:%+v", logId, err)
		return err
	}

	if instance == nil {
		d.SetId("")
		return nil
	}

	_ = d.Set("network_type", instance.LoadBalancerType)
	_ = d.Set("clb_name", instance.LoadBalancerName)
	_ = d.Set("clb_vips", helper.StringsInterfaces(instance.LoadBalancerVips))
	_ = d.Set("subnet_id", instance.SubnetId)
	_ = d.Set("vpc_id", instance.VpcId)
	_ = d.Set("target_region_info_region", instance.TargetRegionInfo.Region)
	_ = d.Set("target_region_info_vpc_id", instance.TargetRegionInfo.VpcId)
	_ = d.Set("project_id", instance.ProjectId)
	_ = d.Set("security_groups", helper.StringsInterfaces(instance.SecureGroups))
	_ = d.Set("domain", instance.LoadBalancerDomain)
	_ = d.Set("ipv6_mode", instance.IPv6Mode)
	_ = d.Set("address_ipv6", instance.AddressIPv6)

	if instance.ClusterIds != nil && len(instance.ClusterIds) > 0 {
		_ = d.Set("cluster_id", instance.ClusterIds[0])
	}

	if instance.SlaType != nil {
		_ = d.Set("sla_type", instance.SlaType)
	}

	if instance.VipIsp != nil {
		_ = d.Set("vip_isp", instance.VipIsp)
	}

	if instance.LoadBalancerVips != nil && len(instance.LoadBalancerVips) > 0 {
		_ = d.Set("vip", instance.LoadBalancerVips[0])
	}

	if instance.AddressIPVersion != nil {
		if *instance.AddressIPVersion == "ipv6" && instance.IPv6Mode != nil && *instance.IPv6Mode == "IPv6FullChain" {
			_ = d.Set("address_ip_version", instance.IPv6Mode)
		} else {
			_ = d.Set("address_ip_version", instance.AddressIPVersion)
		}
	}

	if instance.NetworkAttributes != nil {
		_ = d.Set("internet_bandwidth_max_out", instance.NetworkAttributes.InternetMaxBandwidthOut)
		_ = d.Set("internet_charge_type", instance.NetworkAttributes.InternetChargeType)
	}

	if instance.MasterZone != nil {
		_ = d.Set("master_zone_id", instance.MasterZone.Zone)
		_ = d.Set("zone_id", instance.MasterZone.Zone)
	}

	if instance.BackupZoneSet != nil && len(instance.BackupZoneSet) > 0 {
		_ = d.Set("slave_zone_id", instance.BackupZoneSet[0].Zone)
	}

	_ = d.Set("load_balancer_pass_to_target", instance.LoadBalancerPassToTarget)
	_ = d.Set("log_set_id", instance.LogSetId)
	_ = d.Set("log_topic_id", instance.LogTopicId)

	if _, ok := d.GetOk("snat_pro"); ok {
		_ = d.Set("snat_pro", instance.SnatPro)
	}

	if *instance.LoadBalancerType == "INTERNAL" {
		request := vpc.NewDescribeAddressesRequest()
		request.Filters = []*vpc.Filter{
			{
				Name:   helper.String("instance-id"),
				Values: helper.Strings([]string{clbId}),
			},
		}
		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseVpcClient().DescribeAddresses(request)
			if e != nil {
				return tccommon.RetryError(e)
			}

			if result == nil || result.Response == nil || result.Response.AddressSet == nil {
				e = fmt.Errorf("Describe CLB instance EIP failed")
				return resource.NonRetryableError(e)
			}

			if len(result.Response.AddressSet) == 1 {
				if result.Response.AddressSet[0].AddressId != nil {
					_ = d.Set("eip_address_id", result.Response.AddressSet[0].AddressId)
				}
			}

			return nil
		})

		if err != nil {
			log.Printf("[CRITAL]%s Describe CLB instance EIP failed, reason:%+v", logId, err)
			return err
		}
	}

	if instance.AssociateEndpoint != nil {
		_ = d.Set("associate_endpoint", instance.AssociateEndpoint)
	}

	tcClient := meta.(tccommon.ProviderMeta).GetAPIV3Conn()
	tagService := svctag.NewTagService(tcClient)
	tags, err := tagService.DescribeResourceTags(ctx, "clb", "clb", tcClient.Region, d.Id())
	if err != nil {
		return err
	}

	_ = d.Set("tags", tags)
	return nil
}

func resourceTencentCloudClbInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_clb_instance.update")()

	clbActionMu.Lock()
	defer clbActionMu.Unlock()

	var (
		logId = tccommon.GetLogId(tccommon.ContextNil)
		ctx   = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		clbId = d.Id()
	)

	immutableArgs := []string{"snat_ips", "dynamic_vip", "master_zone_id", "slave_zone_id", "vpc_id", "subnet_id", "address_ip_version", "bandwidth_package_id", "zone_id"}
	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	d.Partial(true)

	request := clb.NewModifyLoadBalancerAttributesRequest()
	request.LoadBalancerId = helper.String(clbId)
	clbName := ""
	targetRegionInfo := clb.TargetRegionInfo{}
	internet := clb.InternetAccessible{}
	changed := false
	isLoadBalancerPassToTgt := false
	isDeleteProtect := false
	snatPro := d.Get("snat_pro").(bool)

	if d.HasChange("clb_name") {
		changed = true
		clbName = d.Get("clb_name").(string)
		flag, err := checkSameName(clbName, meta)
		if err != nil {
			return err
		}

		if flag {
			return fmt.Errorf("[CHECK][CLB instance][Update] check: Same CLB name %s exists!", clbName)
		}

		request.LoadBalancerName = helper.String(clbName)
	}

	if d.HasChange("target_region_info_region") || d.HasChange("target_region_info_vpc_id") {
		if d.Get("network_type") == CLB_NETWORK_TYPE_INTERNAL {
			return fmt.Errorf("[CHECK][CLB instance %s][Update] check: INTERNAL network_type do not support this operation with target_region_info", clbId)
		}

		changed = true
		region := d.Get("target_region_info_region").(string)
		vpcId := d.Get("target_region_info_vpc_id").(string)
		targetRegionInfo = clb.TargetRegionInfo{
			Region: &region,
			VpcId:  &vpcId,
		}
		request.TargetRegionInfo = &targetRegionInfo
	}

	if d.HasChange("sla_type") {
		slaRequest := clb.NewModifyLoadBalancerSlaRequest()
		param := clb.SlaUpdateParam{}
		param.LoadBalancerId = &clbId
		param.SlaType = helper.String(d.Get("sla_type").(string))
		slaRequest.LoadBalancerSla = []*clb.SlaUpdateParam{&param}
		var taskId string
		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseClbClient().ModifyLoadBalancerSla(slaRequest)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}

			taskId = *result.Response.RequestId
			return nil
		})

		if err != nil {
			log.Printf("[CRITAL]%s update clb instanceSlaConfig failed, reason:%+v", logId, err)
			return err
		}

		retryErr := waitForTaskFinish(taskId, meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseClbClient())
		if retryErr != nil {
			return retryErr
		}
	}

	if d.HasChange("internet_charge_type") || d.HasChange("internet_bandwidth_max_out") {
		changed = true
		chargeType := d.Get("internet_charge_type").(string)
		bandwidth := d.Get("internet_bandwidth_max_out").(int)
		if chargeType != "" {
			internet.InternetChargeType = &chargeType
		}

		if bandwidth > 0 {
			internet.InternetMaxBandwidthOut = helper.IntInt64(bandwidth)
		}

		request.InternetChargeInfo = &internet
	}

	if d.HasChange("load_balancer_pass_to_target") {
		changed = true
		isLoadBalancerPassToTgt = d.Get("load_balancer_pass_to_target").(bool)
		request.LoadBalancerPassToTarget = &isLoadBalancerPassToTgt
	}

	if d.HasChange("snat_pro") {
		changed = true
		request.SnatPro = &snatPro
	}

	if d.HasChange("delete_protect") {
		changed = true
		isDeleteProtect = d.Get("delete_protect").(bool)
		request.DeleteProtect = &isDeleteProtect
	}

	if d.HasChange("associate_endpoint") {
		changed = true
		associateEndpoint := d.Get("associate_endpoint").(string)
		request.AssociateEndpoint = &associateEndpoint
	}

	if changed {
		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			response, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseClbClient().ModifyLoadBalancerAttributes(request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
					logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
				requestId := *response.Response.RequestId
				retryErr := waitForTaskFinish(requestId, meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseClbClient())
				if retryErr != nil {
					return tccommon.RetryError(retryErr)
				}
			}

			return nil
		})

		if err != nil {
			log.Printf("[CRITAL]%s update CLB instance failed, reason:%+v", logId, err)
			return err
		}
	}

	if d.HasChange("security_groups") {
		sgRequest := clb.NewSetLoadBalancerSecurityGroupsRequest()
		sgRequest.LoadBalancerId = helper.String(clbId)
		securityGroups := d.Get("security_groups").([]interface{})
		sgRequest.SecurityGroups = make([]*string, 0, len(securityGroups))
		for i := range securityGroups {
			if securityGroup, ok := securityGroups[i].(string); ok && securityGroup != "" {
				sgRequest.SecurityGroups = append(sgRequest.SecurityGroups, &securityGroup)
			}
		}

		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			sgResponse, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseClbClient().SetLoadBalancerSecurityGroups(sgRequest)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
					logId, sgRequest.GetAction(), sgRequest.ToJsonString(), sgResponse.ToJsonString())
				requestId := *sgResponse.Response.RequestId
				retryErr := waitForTaskFinish(requestId, meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseClbClient())
				if retryErr != nil {
					return tccommon.RetryError(errors.WithStack(retryErr))
				}
			}

			return nil
		})

		if err != nil {
			log.Printf("[CRITAL]%s update CLB instance security_group failed, reason:%+v", logId, err)
			return err
		}
	}

	if d.HasChange("log_set_id") || d.HasChange("log_topic_id") {
		logSetId := d.Get("log_set_id")
		logTopicId := d.Get("log_topic_id")
		logRequest := clb.NewSetLoadBalancerClsLogRequest()
		logRequest.LoadBalancerId = helper.String(clbId)
		logRequest.LogSetId = helper.String(logSetId.(string))
		logRequest.LogTopicId = helper.String(logTopicId.(string))
		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			logResponse, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseClbClient().SetLoadBalancerClsLog(logRequest)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
					logId, logRequest.GetAction(), logRequest.ToJsonString(), logResponse.ToJsonString())
				requestId := *logResponse.Response.RequestId
				retryErr := waitForTaskFinish(requestId, meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseClbClient())
				if retryErr != nil {
					return tccommon.RetryError(errors.WithStack(retryErr))
				}
			}

			return nil
		})

		if err != nil {
			log.Printf("[CRITAL]%s set CLB instance log failed, reason:%+v", logId, err)
			return err
		}
	}

	if d.HasChange("project_id") {
		var projectId int
		if v, ok := d.GetOkExists("project_id"); ok {
			projectId = v.(int)
		}

		pRequest := clb.NewModifyLoadBalancersProjectRequest()
		pRequest.LoadBalancerIds = []*string{&clbId}
		pRequest.ProjectId = helper.IntUint64(projectId)
		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			pResponse, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseClbClient().ModifyLoadBalancersProject(pRequest)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
					logId, pRequest.GetAction(), pRequest.ToJsonString(), pResponse.ToJsonString())
				requestId := *pResponse.Response.RequestId
				retryErr := waitForTaskFinish(requestId, meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseClbClient())
				if retryErr != nil {
					return tccommon.RetryError(errors.WithStack(retryErr))
				}
			}

			return nil
		})

		if err != nil {
			log.Printf("[CRITAL]%s update CLB instance project_id failed, reason:%+v", logId, err)
			return err
		}
	}

	if d.HasChange("eip_address_id") {
		oldEip, newEip := d.GetChange("eip_address_id")
		oldEipStr := oldEip.(string)
		newEipStr := newEip.(string)
		// delete old first
		if oldEipStr != "" {
			request := vpc.NewDisassociateAddressRequest()
			request.AddressId = helper.String(oldEipStr)
			err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
				_, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseVpcClient().DisassociateAddress(request)
				if e != nil {
					return tccommon.RetryError(e)
				}

				return nil
			})

			if err != nil {
				log.Printf("[CRITAL]%s Disassociate EIP failed, reason:%+v", logId, err)
				return err
			}

			// wait
			eipRequest := vpc.NewDescribeAddressesRequest()
			eipRequest.AddressIds = helper.Strings([]string{oldEipStr})
			err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
				result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseVpcClient().DescribeAddresses(eipRequest)
				if e != nil {
					return tccommon.RetryError(e)
				}

				if result == nil || result.Response == nil || result.Response.AddressSet == nil || len(result.Response.AddressSet) != 1 {
					e = fmt.Errorf("Describe CLB instance EIP failed")
					return resource.NonRetryableError(e)
				}

				if *result.Response.AddressSet[0].AddressStatus != "UNBIND" {
					return resource.RetryableError(fmt.Errorf("EIP status is still %s", *result.Response.AddressSet[0].AddressStatus))
				}

				return nil
			})

			if err != nil {
				log.Printf("[CRITAL]%s Describe CLB instance EIP failed, reason:%+v", logId, err)
				return err
			}
		}

		// attach new
		if newEipStr != "" {
			request := vpc.NewAssociateAddressRequest()
			request.AddressId = helper.String(newEipStr)
			request.InstanceId = helper.String(clbId)
			err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
				_, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseVpcClient().AssociateAddress(request)
				if e != nil {
					return tccommon.RetryError(e)
				}

				return nil
			})

			if err != nil {
				log.Printf("[CRITAL]%s Associate EIP failed, reason:%+v", logId, err)
				return err
			}

			// wait
			eipRequest := vpc.NewDescribeAddressesRequest()
			eipRequest.AddressIds = helper.Strings([]string{newEipStr})
			err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
				result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseVpcClient().DescribeAddresses(eipRequest)
				if e != nil {
					return tccommon.RetryError(e)
				}

				if result == nil || result.Response == nil || result.Response.AddressSet == nil || len(result.Response.AddressSet) != 1 {
					e = fmt.Errorf("Describe CLB instance EIP failed")
					return resource.NonRetryableError(e)
				}

				if *result.Response.AddressSet[0].AddressStatus != "BIND" {
					return resource.RetryableError(fmt.Errorf("EIP status is still %s", *result.Response.AddressSet[0].AddressStatus))
				}

				return nil
			})

			if err != nil {
				log.Printf("[CRITAL]%s Describe CLB instance EIP failed, reason:%+v", logId, err)
				return err
			}
		}
	}

	if d.HasChange("tags") {
		oldValue, newValue := d.GetChange("tags")
		replaceTags, deleteTags := svctag.DiffTags(oldValue.(map[string]interface{}), newValue.(map[string]interface{}))
		tcClient := meta.(tccommon.ProviderMeta).GetAPIV3Conn()
		tagService := svctag.NewTagService(tcClient)
		resourceName := tccommon.BuildTagResourceName("clb", "clb", tcClient.Region, d.Id())
		err := tagService.ModifyTags(ctx, resourceName, replaceTags, deleteTags)
		if err != nil {
			return err
		}
	}

	d.Partial(false)
	return nil
}

func resourceTencentCloudClbInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_clb_instance.delete")()

	clbActionMu.Lock()
	defer clbActionMu.Unlock()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		clbService = ClbService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		clbId      = d.Id()
	)

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		e := clbService.DeleteLoadBalancerById(ctx, clbId)
		if e != nil {
			if ve, ok := e.(*sdkErrors.TencentCloudSDKError); ok {
				if ve.GetCode() == "FailedOperation.ResourceInOperating" {
					return tccommon.RetryError(e, "FailedOperation.ResourceInOperating")
				}
			}

			return tccommon.RetryError(e)
		}

		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s delete CLB instance failed, reason:%+v", logId, err)
		return err
	}

	return nil
}

func checkSameName(name string, meta interface{}) (flag bool, errRet error) {
	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		clbService = ClbService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	flag = false
	params := make(map[string]interface{})
	params["clb_name"] = name
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		clbs, e := clbService.DescribeLoadBalancerByFilter(ctx, params)
		if e != nil {
			return tccommon.RetryError(e)
		}

		if len(clbs) > 0 {
			//this function is a fuzzy query
			// so take a further check
			for _, clbInfo := range clbs {
				if *clbInfo.LoadBalancerName == name {
					flag = true
					return nil
				}
			}
		}

		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s read CLB instance failed, reason:%+v", logId, err)
	}

	errRet = err
	return
}
