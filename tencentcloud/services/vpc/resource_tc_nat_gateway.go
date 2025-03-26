package vpc

import (
	"context"
	"fmt"
	"log"
	"time"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	svctag "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/tag"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudNatGateway() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudNatGatewayCreate,
		Read:   resourceTencentCloudNatGatewayRead,
		Update: resourceTencentCloudNatGatewayUpdate,
		Delete: resourceTencentCloudNatGatewayDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"vpc_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ID of the vpc.",
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: tccommon.ValidateStringLengthInRange(1, 60),
				Description:  "Name of the NAT gateway.",
			},
			"max_concurrent": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: tccommon.ValidateAllowedIntValue([]int{1000000, 3000000, 10000000, 2000000}),
				Description:  "The upper limit of concurrent connection of NAT gateway. Valid values: `1000000`, `3000000`, `10000000`. Default is `1000000`. When the value of parameter `nat_product_version` is 2, which is the standard NAT type, this parameter does not need to be filled in and defaults to `2000000`.",
			},
			"bandwidth": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: tccommon.ValidateAllowedIntValue([]int{20, 50, 100, 200, 500, 1000, 2000, 5000}),
				Description:  "The maximum public network output bandwidth of NAT gateway (unit: Mbps). Valid values: `20`, `50`, `100`, `200`, `500`, `1000`, `2000`, `5000`. Default is `100`. When the value of parameter `nat_product_version` is 2, which is the standard NAT type, this parameter does not need to be filled in and defaults to `5000`.",
			},
			"assigned_eip_set": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: tccommon.ValidateIp,
				},
				MinItems:    1,
				Description: "EIP IP address set bound to the gateway. The value of at least 1 and at most 10 if do not apply for a whitelist.",
			},
			"zone": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The availability zone, such as `ap-guangzhou-3`.",
			},
			"subnet_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "Subnet of NAT.",
			},
			"nat_product_version": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: tccommon.ValidateAllowedIntValue([]int{1, 2}),
				Description:  "1: traditional NAT, 2: standard NAT, default value is 1.",
			},
			"stock_public_ip_addresses_bandwidth_out": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "The elastic public IP bandwidth value (unit: Mbps) for binding NAT gateway. When this parameter is not filled in, it defaults to the bandwidth value of the elastic public IP, and for some users, it defaults to the bandwidth limit of the elastic public IP of that user type.",
			},
			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "The available tags within this NAT gateway.",
			},
			//computed
			"created_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Create time of the NAT gateway.",
			},
		},
	}
}

func resourceTencentCloudNatGatewayCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_nat_gateway.create")()

	var (
		logId                      = tccommon.GetLogId(tccommon.ContextNil)
		ctx                        = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		request                    = vpc.NewCreateNatGatewayRequest()
		maxConcurrentConnection    int
		internetMaxBandwidthOut    int
		hasMaxConcurrentConnection bool
		hasInternetMaxBandwidthOut bool
	)

	if v, ok := d.GetOkExists("max_concurrent"); ok {
		maxConcurrentConnection = v.(int)
		hasMaxConcurrentConnection = true
		request.MaxConcurrentConnection = helper.IntUint64(maxConcurrentConnection)
	}

	if v, ok := d.GetOkExists("bandwidth"); ok {
		internetMaxBandwidthOut = v.(int)
		hasInternetMaxBandwidthOut = true
		request.InternetMaxBandwidthOut = helper.IntUint64(internetMaxBandwidthOut)
	}

	if v, ok := d.GetOkExists("nat_product_version"); ok {
		request.NatProductVersion = helper.IntUint64(v.(int))
		if v.(int) == 2 {
			if hasMaxConcurrentConnection && maxConcurrentConnection != 2000000 {
				return fmt.Errorf("If `nat_product_version` is 2, `max_concurrent` can only be set to `2000000` or not set at all.")
			}

			if hasInternetMaxBandwidthOut && internetMaxBandwidthOut != 5000 {
				return fmt.Errorf("If `nat_product_version` is 2, `bandwidth` can only be set to `5000` or not set at all.")
			}
		}
	}

	if v, ok := d.GetOk("vpc_id"); ok {
		request.VpcId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("name"); ok {
		request.NatGatewayName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("assigned_eip_set"); ok {
		eipSet := v.(*schema.Set).List()
		for i := range eipSet {
			publicIp := eipSet[i].(string)
			request.PublicIpAddresses = append(request.PublicIpAddresses, &publicIp)
		}
	}

	if v, ok := d.GetOk("zone"); ok {
		request.Zone = helper.String(v.(string))
	}

	if v, ok := d.GetOk("subnet_id"); ok {
		request.SubnetId = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("stock_public_ip_addresses_bandwidth_out"); ok {
		request.StockPublicIpAddressesBandwidthOut = helper.IntUint64(v.(int))
	}

	if v := helper.GetTags(d, "tags"); len(v) > 0 {
		for tagKey, tagValue := range v {
			tag := vpc.Tag{
				Key:   helper.String(tagKey),
				Value: helper.String(tagValue),
			}

			request.Tags = append(request.Tags, &tag)
		}
	}

	var response *vpc.CreateNatGatewayResponse
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseVpcClient().CreateNatGateway(request)
		if e != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), e.Error())
			return tccommon.RetryError(e)
		}

		if result == nil || result.Response == nil || len(result.Response.NatGatewaySet) < 1 {
			e = fmt.Errorf("create NAT gateway failed.")
			return resource.NonRetryableError(e)
		}

		response = result
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create NAT gateway failed, reason:%s\n", logId, err.Error())
		return err
	}

	natGatewayId := *response.Response.NatGatewaySet[0].NatGatewayId
	d.SetId(natGatewayId)

	// must wait for finishing creating NAT
	statRequest := vpc.NewDescribeNatGatewaysRequest()
	statRequest.NatGatewayIds = []*string{&natGatewayId}
	err = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseVpcClient().DescribeNatGateways(statRequest)
		if e != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, statRequest.GetAction(), statRequest.ToJsonString(), e.Error())
			return tccommon.RetryError(e)
		} else {
			//if not, quit
			if len(result.Response.NatGatewaySet) != 1 {
				return resource.NonRetryableError(fmt.Errorf("creating error"))
			}

			//else get stat
			nat := result.Response.NatGatewaySet[0]
			stat := *nat.State
			if stat == "AVAILABLE" {
				return nil
			}

			return resource.RetryableError(fmt.Errorf("creating not ready retry"))
		}
	})

	if err != nil {
		log.Printf("[CRITAL]%s create NAT gateway failed, reason:%s\n", logId, err.Error())
		return err
	}

	//cs::vpc:ap-guangzhou:uin/12345:nat/nat-nxxx
	if tags := helper.GetTags(d, "tags"); len(tags) > 0 {
		tcClient := meta.(tccommon.ProviderMeta).GetAPIV3Conn()
		tagService := svctag.NewTagService(tcClient)
		resourceName := tccommon.BuildTagResourceName("vpc", "nat", tcClient.Region, d.Id())
		if err := tagService.ModifyTags(ctx, resourceName, tags, nil); err != nil {
			return err
		}
	}

	return resourceTencentCloudNatGatewayRead(d, meta)
}

func resourceTencentCloudNatGatewayRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_nat_gateway.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	natGatewayId := d.Id()
	request := vpc.NewDescribeNatGatewaysRequest()
	request.NatGatewayIds = []*string{&natGatewayId}
	var response *vpc.DescribeNatGatewaysResponse
	var iacExtInfo connectivity.IacExtInfo
	iacExtInfo.InstanceId = natGatewayId
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseVpcClient(iacExtInfo).DescribeNatGateways(request)
		if e != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), e.Error())
			return tccommon.RetryError(e)
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s read NAT gateway failed, reason:%s\n", logId, err.Error())
		return err
	}
	if len(response.Response.NatGatewaySet) < 1 {
		d.SetId("")
		return nil
	}

	nat := response.Response.NatGatewaySet[0]

	_ = d.Set("vpc_id", *nat.VpcId)
	_ = d.Set("name", *nat.NatGatewayName)
	_ = d.Set("max_concurrent", *nat.MaxConcurrentConnection)
	_ = d.Set("bandwidth", *nat.InternetMaxBandwidthOut)
	_ = d.Set("created_time", *nat.CreatedTime)
	_ = d.Set("assigned_eip_set", flattenAddressList((*nat).PublicIpAddressSet))
	_ = d.Set("zone", *nat.Zone)
	if nat.SubnetId != nil {
		_ = d.Set("subnet_id", *nat.SubnetId)
	}
	if nat.NatProductVersion != nil {
		_ = d.Set("nat_product_version", *nat.NatProductVersion)
	}

	// set `stock_public_ip_addresses_bandwidth_out`
	bandwidthRequest := vpc.NewDescribeAddressesRequest()
	bandwidthResponse := vpc.NewDescribeAddressesResponse()
	bandwidthRequest.Filters = []*vpc.Filter{
		{
			Name:   common.StringPtr("address-ip"),
			Values: flattenAddressList((*nat).PublicIpAddressSet),
		},
	}

	err = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseVpcClient(iacExtInfo).DescribeAddresses(bandwidthRequest)
		if e != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, bandwidthRequest.GetAction(), bandwidthRequest.ToJsonString(), e.Error())
			return tccommon.RetryError(e)
		}
		bandwidthResponse = result
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s read NAT gateway addresses failed, reason:%s\n", logId, err.Error())
		return err
	}

	if bandwidthResponse != nil && len(bandwidthResponse.Response.AddressSet) > 0 {
		address := bandwidthResponse.Response.AddressSet[0]
		if address.Bandwidth != nil {
			_ = d.Set("stock_public_ip_addresses_bandwidth_out", address.Bandwidth)
		}
	}

	tcClient := meta.(tccommon.ProviderMeta).GetAPIV3Conn()
	tagService := svctag.NewTagService(tcClient)
	tags, err := tagService.DescribeResourceTags(ctx, "vpc", "nat", tcClient.Region, d.Id())
	if err != nil {
		return err
	}
	_ = d.Set("tags", tags)

	return nil
}

func resourceTencentCloudNatGatewayUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_nat_gateway.update")()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		vpcService = VpcService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	immutableArgs := []string{"zone"}
	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	d.Partial(true)
	natGatewayId := d.Id()
	request := vpc.NewModifyNatGatewayAttributeRequest()
	request.NatGatewayId = &natGatewayId
	changed := false
	if d.HasChange("name") {
		request.NatGatewayName = helper.String(d.Get("name").(string))
		changed = true
	}

	if d.HasChange("bandwidth") {
		bandwidth := d.Get("bandwidth").(int)
		bandwidth64 := uint64(bandwidth)
		if v, ok := d.GetOkExists("nat_product_version"); ok {
			if v.(int) == 2 && bandwidth64 != 5000 {
				return fmt.Errorf("If `nat_product_version` is 2, `bandwidth` can only be set to `5000` or not set at all.")
			}
		}
		request.InternetMaxBandwidthOut = &bandwidth64
		changed = true
	}

	if changed {
		err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			_, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseVpcClient().ModifyNatGatewayAttribute(request)
			if e != nil {
				log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
					logId, request.GetAction(), request.ToJsonString(), e.Error())
				return tccommon.RetryError(e)
			}

			return nil
		})

		if err != nil {
			log.Printf("[CRITAL]%s modify NAT gateway failed, reason:%s\n", logId, err.Error())
			return err
		}
	}

	//max concurrent
	if d.HasChange("max_concurrent") {
		concurrentReq := vpc.NewResetNatGatewayConnectionRequest()
		concurrentReq.NatGatewayId = &natGatewayId
		concurrent := d.Get("max_concurrent").(int)
		concurrent64 := uint64(concurrent)
		if v, ok := d.GetOkExists("nat_product_version"); ok {
			if v.(int) == 2 && concurrent64 != 2000000 {
				return fmt.Errorf("If `nat_product_version` is 2, `max_concurrent` can only be set to `2000000` or not set at all.")
			}
		}
		concurrentReq.MaxConcurrentConnection = &concurrent64
		err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			_, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseVpcClient().ResetNatGatewayConnection(concurrentReq)
			if e != nil {
				log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
					logId, concurrentReq.GetAction(), concurrentReq.ToJsonString(), e.Error())
				return tccommon.RetryError(e, tccommon.InternalError)
			}
			return nil
		})

		if err != nil {
			log.Printf("[CRITAL]%s modify NAT gateway concurrent failed, reason:%s\n", logId, err.Error())
			return err
		}
	}

	//eip
	if d.HasChange("assigned_eip_set") {
		eipSetLength := 0
		if v, ok := d.GetOk("assigned_eip_set"); ok {
			eipSet := v.(*schema.Set).List()
			eipSetLength = len(eipSet)
		}

		if d.HasChange("assigned_eip_set") {
			o, n := d.GetChange("assigned_eip_set")
			os := o.(*schema.Set)
			ns := n.(*schema.Set)
			oldEipSet := os.List()
			newEipSet := ns.List()

			//in case of no union set
			backUpOldIp := ""
			backUpNewIp := ""
			//Unassign eips
			if len(oldEipSet) > 0 {
				unassignedRequest := vpc.NewDisassociateNatGatewayAddressRequest()
				unassignedRequest.PublicIpAddresses = make([]*string, 0, len(oldEipSet))
				unassignedRequest.NatGatewayId = &natGatewayId
				//set request public ips
				for i := range oldEipSet {
					publicIp := oldEipSet[i].(string)
					isIn := false
					for j := range newEipSet {
						if publicIp == newEipSet[j] {
							isIn = true
						}
					}

					if !isIn {
						if len(unassignedRequest.PublicIpAddresses)+1 == len(oldEipSet) {
							backUpOldIp = publicIp
						} else {
							unassignedRequest.PublicIpAddresses = append(unassignedRequest.PublicIpAddresses, &publicIp)
						}
					}
				}

				if len(unassignedRequest.PublicIpAddresses) > 0 {
					var response *vpc.DisassociateNatGatewayAddressResponseParams
					err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
						result, e := vpcService.DisassociateNatGatewayAddress(ctx, unassignedRequest)
						if e != nil {
							return tccommon.RetryError(e)
						}

						if result != nil && result.Response != nil {
							response = result.Response
						}

						return nil
					})

					if err != nil {
						log.Printf("[CRITAL]%s modify NAT gateway EIP failed, reason:%s\n", logId, err.Error())
						return err
					}

					if response != nil && response.RequestId != nil {
						err = vpcService.DescribeVpcTaskResult(ctx, response.RequestId)
						if err != nil {
							return err
						}
					}
				}
			}

			//Assign new EIP
			if len(newEipSet) > 0 {
				assignedRequest := vpc.NewAssociateNatGatewayAddressRequest()
				assignedRequest.PublicIpAddresses = make([]*string, 0, len(newEipSet))
				assignedRequest.NatGatewayId = &natGatewayId
				//set request public ips
				for i := range newEipSet {
					publicIp := newEipSet[i].(string)
					isIn := false
					for j := range oldEipSet {
						if publicIp == oldEipSet[j] {
							isIn = true
						}
					}

					if !isIn {
						if len(assignedRequest.PublicIpAddresses)+eipSetLength+1 == NAT_EIP_MAX_LIMIT {
							backUpNewIp = publicIp
						} else {
							assignedRequest.PublicIpAddresses = append(assignedRequest.PublicIpAddresses, &publicIp)
						}
					}
				}

				if len(assignedRequest.PublicIpAddresses) > 0 {
					var response *vpc.AssociateNatGatewayAddressResponseParams
					err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
						result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseVpcClient().AssociateNatGatewayAddress(assignedRequest)
						if e != nil {
							log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
								logId, assignedRequest.GetAction(), assignedRequest.ToJsonString(), e.Error())
							return tccommon.RetryError(e)
						}

						if result != nil && result.Response != nil {
							response = result.Response
						}

						return nil
					})

					if err != nil {
						log.Printf("[CRITAL]%s modify NAT gateway EIP failed, reason:%s\n", logId, err.Error())
						return err
					}

					if response != nil && response.RequestId != nil {
						err = vpcService.DescribeVpcTaskResult(ctx, response.RequestId)
						if err != nil {
							return err
						}
					}
				}
			}

			if backUpOldIp != "" {
				//disassociate one old ip
				var response *vpc.DisassociateNatGatewayAddressResponseParams
				unassignedRequest := vpc.NewDisassociateNatGatewayAddressRequest()
				unassignedRequest.NatGatewayId = &natGatewayId
				unassignedRequest.PublicIpAddresses = []*string{&backUpOldIp}
				err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
					result, e := vpcService.DisassociateNatGatewayAddress(ctx, unassignedRequest)
					if e != nil {
						return tccommon.RetryError(e)
					}

					if result != nil && result.Response != nil {
						response = result.Response
					}

					return nil
				})

				if err != nil {
					log.Printf("[CRITAL]%s modify NAT gateway EIP failed, reason:%s\n", logId, err.Error())
					return err
				}

				if response != nil && response.RequestId != nil {
					err = vpcService.DescribeVpcTaskResult(ctx, response.RequestId)
					if err != nil {
						return err
					}
				}
			}

			if backUpNewIp != "" {
				//associate one new ip
				var response *vpc.AssociateNatGatewayAddressResponseParams
				assignedRequest := vpc.NewAssociateNatGatewayAddressRequest()
				assignedRequest.NatGatewayId = &natGatewayId
				assignedRequest.PublicIpAddresses = []*string{&backUpNewIp}
				err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
					result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseVpcClient().AssociateNatGatewayAddress(assignedRequest)
					if e != nil {
						log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
							logId, assignedRequest.GetAction(), assignedRequest.ToJsonString(), e.Error())
						return tccommon.RetryError(e)
					}

					if result != nil && result.Response != nil {
						response = result.Response
					}

					return nil
				})

				if err != nil {
					log.Printf("[CRITAL]%s modify NAT gateway EIP failed, reason:%s\n", logId, err.Error())
					return err
				}

				if response != nil && response.RequestId != nil {
					err = vpcService.DescribeVpcTaskResult(ctx, response.RequestId)
					if err != nil {
						return err
					}
				}
			}
		}
	}

	if d.HasChange("tags") {
		oldValue, newValue := d.GetChange("tags")
		replaceTags, deleteTags := svctag.DiffTags(oldValue.(map[string]interface{}), newValue.(map[string]interface{}))
		tcClient := meta.(tccommon.ProviderMeta).GetAPIV3Conn()
		tagService := svctag.NewTagService(tcClient)
		resourceName := tccommon.BuildTagResourceName("vpc", "nat", tcClient.Region, d.Id())
		err := tagService.ModifyTags(ctx, resourceName, replaceTags, deleteTags)
		if err != nil {
			return err
		}
	}

	d.Partial(false)

	return resourceTencentCloudNatGatewayRead(d, meta)
}

func resourceTencentCloudNatGatewayDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_nat_gateway.delete")()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	natGatewayId := d.Id()
	request := vpc.NewDeleteNatGatewayRequest()
	request.NatGatewayId = &natGatewayId
	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		_, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseVpcClient().DeleteNatGateway(request)
		if e != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), e.Error())
			return tccommon.RetryError(e)
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s delete NAT gateway failed, reason:%s\n", logId, err.Error())
		return err
	}
	// must wait for finishing deleting NAT
	time.Sleep(10 * time.Second)
	//to get the status of NAT

	statRequest := vpc.NewDescribeNatGatewaysRequest()
	statRequest.NatGatewayIds = []*string{&natGatewayId}
	err = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseVpcClient().DescribeNatGateways(statRequest)
		if e != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), e.Error())
			return tccommon.RetryError(e)
		} else {
			//if not, quit
			if len(result.Response.NatGatewaySet) == 0 {
				log.Printf("deleting done")
				return nil
			}
			//else get stat
			nat := result.Response.NatGatewaySet[0]
			stat := *nat.State
			if stat == NAT_FAILED_STATE {
				return resource.NonRetryableError(fmt.Errorf("delete NAT failed"))
			}
			time.Sleep(3 * time.Second)

			return resource.RetryableError(fmt.Errorf("deleting retry"))
		}
	})
	if err != nil {
		log.Printf("[CRITAL]%s delete NAT gateway failed, reason:%s\n", logId, err.Error())
		return err
	}
	return nil
}

func flattenAddressList(addresses []*vpc.NatGatewayAddress) (eips []*string) {
	for _, address := range addresses {
		eips = append(eips, address.PublicIpAddress)
	}
	return
}
