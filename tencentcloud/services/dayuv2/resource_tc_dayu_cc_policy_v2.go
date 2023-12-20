package dayuv2

import (
	"context"
	"fmt"
	"strings"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	antiddos "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/antiddos/v20200309"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudDayuCCPolicyV2() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudDayuCCPolicyV2Create,
		Read:   resourceTencentCloudDayuCCPolicyV2Read,
		Update: resourceTencentCloudDayuCCPolicyV2Update,
		Delete: resourceTencentCloudDayuCCPolicyV2Delete,

		Schema: map[string]*schema.Schema{
			"resource_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of the resource instance.",
			},
			"business": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Business of resource instance. bgpip indicates anti-anti-ip ip; bgp means exclusive package; bgp-multip means shared packet; net indicates anti-anti-ip pro version.",
			},
			"thresholds": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"threshold": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Cleaning threshold, -1 indicates that the `default` mode is turned on.",
						},
						"domain": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "domain.",
						},
					},
				},
				Description: "List of protection threshold configurations.",
			},
			"cc_geo_ip_policys": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"action": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "User action, drop or arg.",
						},
						"area_list": {
							Type:        schema.TypeList,
							Elem:        &schema.Schema{Type: schema.TypeInt},
							Optional:    true,
							Computed:    true,
							Description: "The list of region IDs that the user selects to block.",
						},
						"region_type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Regional types, divided into china, oversea and customized.",
						},
						"create_time": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Create time.",
						},
						"modify_time": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Modify time.",
						},
						"domain": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "domain.",
						},
						"protocol": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Protocol, preferably HTTP, HTTPS.",
						},
					},
				},
				Description: "Details of the CC region blocking policy list.",
			},
			"cc_black_white_ips": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"black_white_ip": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Blacklist and whitelist IP addresses.",
						},
						"domain": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Domain.",
						},
						"protocol": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Protocol.",
						},
						"type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "IP type, value [black(blacklist IP), white (whitelist IP)].",
						},
						"create_time": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Create time.",
						},
						"modify_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Optional:    true,
							Description: "Modify time.",
						},
					},
				},
				Description: "Blacklist and whitelist.",
			},
			"cc_precision_policys": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"policy_action": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Policy mode (discard or captcha).",
						},
						"domain": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Domain.",
						},
						"protocol": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Protocol.",
						},
						"ip": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Ip address.",
						},
						"policy_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Policy Id.",
						},
						"policys": {
							Type: schema.TypeList,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"field_name": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Configuration item types, currently only support value.",
									},
									"field_type": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Configuration fields with the desirable values cgi, ua, cookie, referer, accept, srcip.",
									},
									"value": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Configure the value.",
									},
									"value_operator": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Configure the item-value comparison mode, which can be taken as the value of evaluate, not_equal, include.",
									},
								},
							},
							Required:    true,
							Description: "A list of policies.",
						},
					},
				},
				Description: "CC Precision Protection List.",
			},
			"cc_precision_req_limits": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"domain": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Domain.",
						},
						"protocol": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Protocol, preferably HTTP, HTTPS.",
						},
						"ip": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "IP address.",
						},
						"level": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Protection rating, the optional value of default means default policy, loose means loose, and strict means strict.",
						},
						"instance_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Instance id.",
						},
						"policys": {
							Type: schema.TypeList,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"action": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "The frequency limit policy mode, the optional value of arg indicates the verification code, and drop indicates the discard.",
									},
									"cookie": {
										Type:        schema.TypeString,
										Optional:    true,
										Default:     "",
										Description: "Cookies, one of the three policy entries can only be filled in.",
									},
									"execute_duration": {
										Type:        schema.TypeInt,
										Required:    true,
										Description: "The duration of the frequency limit policy can be taken from 1 to 86400 per second.",
									},
									"mode": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "The policy item is compared, and the optional value include indicates inclusion, and equal means equal.",
									},
									"period": {
										Type:        schema.TypeInt,
										Required:    true,
										Description: "Statistical period, take values 1, 10, 30, 60, in seconds.",
									},
									"request_num": {
										Type:        schema.TypeInt,
										Required:    true,
										Description: "The number of requests, the value is 1 to 20000.",
									},
									"uri": {
										Type:        schema.TypeString,
										Optional:    true,
										Default:     "",
										Description: "Uri, one of the three policy entries can only be filled in.",
									},
									"user_agent": {
										Type:        schema.TypeString,
										Optional:    true,
										Default:     "",
										Description: "User-Agent, only one of the three policy entries can be filled in.",
									},
								},
							},
							Required:    true,
							Description: "The CC Frequency Limit Policy Item field.",
						},
					},
				},
				Description: "CC frequency throttling policy.",
			},
		},
	}
}

func resourceTencentCloudDayuCCPolicyV2Create(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dayu_cc_policy_v2.create")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	antiddosService := AntiddosService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	resourceId := d.Get("resource_id").(string)
	business := d.Get("business").(string)
	protectThresholdConfig, err := antiddosService.DescribeListProtectThresholdConfig(ctx, resourceId)
	if err != nil {
		return err
	}
	instanceDetailList := protectThresholdConfig.InstanceDetailList
	if len(instanceDetailList) != 1 {
		return fmt.Errorf("Can not fetch Eip.")
	}
	ip := *instanceDetailList[0].EipList[0]
	if v, ok := d.GetOk("thresholds"); ok {
		thresholds := v.([]interface{})
		for _, threshold := range thresholds {
			thresholdMap := threshold.(map[string]interface{})
			domain := thresholdMap["domain"].(string)
			threshold := thresholdMap["threshold"].(int)
			err := antiddosService.ModifyCCThresholdPolicy(ctx, resourceId, "http", ip, domain, threshold)
			if err != nil {
				return err
			}
		}
	}

	if v, ok := d.GetOk("cc_black_white_ips"); ok {
		ccBlackWhiteIps := v.([]interface{})
		for _, ccBlackWhiteIp := range ccBlackWhiteIps {
			ccBlackWhiteIpMap := ccBlackWhiteIp.(map[string]interface{})
			protocol := ccBlackWhiteIpMap["protocol"].(string)
			domain := ccBlackWhiteIpMap["domain"].(string)
			blackWhiteIp := ccBlackWhiteIpMap["black_white_ip"].(string)
			blackWhiteType := ccBlackWhiteIpMap["type"].(string)
			err := antiddosService.CreateCcBlackWhiteIpList(ctx, resourceId, protocol, ip, domain, blackWhiteType, []string{blackWhiteIp})
			if err != nil {
				return err
			}
		}
	}

	if v, ok := d.GetOk("cc_geo_ip_policys"); ok {
		ccGeoIpPolicys := v.([]interface{})
		for _, ccGeoIpPolicy := range ccGeoIpPolicys {
			ccGeoIpPolicyMap := ccGeoIpPolicy.(map[string]interface{})
			action := ccGeoIpPolicyMap["action"].(string)
			regionType := ccGeoIpPolicyMap["region_type"].(string)
			domain := ccGeoIpPolicyMap["domain"].(string)
			protocol := ccGeoIpPolicyMap["protocol"].(string)
			ccGeoIPBlockConfig := antiddos.CcGeoIPBlockConfig{
				Action:     &action,
				RegionType: &regionType,
			}
			err := antiddosService.CreateCcGeoIPBlockConfig(ctx, resourceId, protocol, ip, domain, ccGeoIPBlockConfig)
			if err != nil {
				return err
			}
		}
	}

	if v, ok := d.GetOk("cc_precision_policys"); ok {
		ccPrecisionPolicys := v.([]interface{})
		for _, ccPrecisionPolicy := range ccPrecisionPolicys {
			ccPrecisionPolicyMap := ccPrecisionPolicy.(map[string]interface{})
			policyAction := ccPrecisionPolicyMap["policy_action"].(string)
			domain := ccPrecisionPolicyMap["domain"].(string)
			protocol := ccPrecisionPolicyMap["protocol"].(string)
			ip := ccPrecisionPolicyMap["ip"].(string)
			ccPrecisionPlyRecords := make([]*antiddos.CCPrecisionPlyRecord, 0)
			policys := ccPrecisionPolicyMap["policys"].([]interface{})
			for _, policy := range policys {
				policyMap := policy.(map[string]interface{})
				fieldName := policyMap["field_name"].(string)
				fieldType := policyMap["field_type"].(string)
				value := policyMap["value"].(string)
				valueOperator := policyMap["value_operator"].(string)
				tmpCCPrecisionPlyRecord := antiddos.CCPrecisionPlyRecord{}
				tmpCCPrecisionPlyRecord.FieldName = &fieldName
				tmpCCPrecisionPlyRecord.FieldType = &fieldType
				tmpCCPrecisionPlyRecord.Value = &value
				tmpCCPrecisionPlyRecord.ValueOperator = &valueOperator
				ccPrecisionPlyRecords = append(ccPrecisionPlyRecords, &tmpCCPrecisionPlyRecord)
			}
			err := antiddosService.CreateCCPrecisionPolicy(ctx, resourceId, protocol, ip, domain, policyAction, ccPrecisionPlyRecords)
			if err != nil {
				return err
			}
		}
	}

	if v, ok := d.GetOk("cc_precision_req_limits"); ok {
		ccPrecisionReqLimits := v.([]interface{})
		for _, ccPrecisionReqLimit := range ccPrecisionReqLimits {
			ccPrecisionReqLimitMap := ccPrecisionReqLimit.(map[string]interface{})
			domain := ccPrecisionReqLimitMap["domain"].(string)
			protocol := ccPrecisionReqLimitMap["protocol"].(string)
			level := ccPrecisionReqLimitMap["level"].(string)
			err := antiddosService.ModifyCCLevelPolicy(ctx, resourceId, ip, domain, protocol, level)
			if err != nil {
				return err
			}
			policys := ccPrecisionReqLimitMap["policys"].([]interface{})
			for _, policy := range policys {
				policyMap := policy.(map[string]interface{})
				action := policyMap["action"].(string)
				executeDuration := policyMap["execute_duration"].(int)
				mode := policyMap["mode"].(string)
				period := policyMap["period"].(int)
				requestNum := policyMap["request_num"].(int)
				uri := policyMap["uri"].(string)
				cookie := policyMap["cookie"].(string)
				userAgent := policyMap["user_agent"].(string)
				ccPolicyRecord := antiddos.CCReqLimitPolicyRecord{
					Action:          &action,
					ExecuteDuration: helper.IntUint64(executeDuration),
					Mode:            &mode,
					Period:          helper.IntUint64(period),
					RequestNum:      helper.IntUint64(requestNum),
				}
				if uri != "" {
					ccPolicyRecord.Uri = &uri
				} else if cookie != "" {
					ccPolicyRecord.Cookie = &cookie
				} else if userAgent != "" {
					ccPolicyRecord.UserAgent = &userAgent
				}
				err := antiddosService.CreateCCReqLimitPolicy(ctx, resourceId, protocol, ip, domain, ccPolicyRecord)
				if err != nil {
					return err
				}
			}

		}
	}
	d.SetId(resourceId + tccommon.FILED_SP + business)
	return resourceTencentCloudDayuCCPolicyV2Read(d, meta)
}

func resourceTencentCloudDayuCCPolicyV2Read(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dayu_cc_policy_v2.read")()
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	items := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(items) < 2 {
		return fmt.Errorf("broken ID of DDoS policy")
	}
	instanceId := items[0]
	business := items[1]
	antiddosService := AntiddosService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	thresholdList, err := antiddosService.DescribeCCThresholdList(ctx, business, instanceId)
	if err != nil {
		return err
	}
	resultThresholds := make([]map[string]interface{}, 0)
	for _, threshold := range thresholdList {
		tmpThreshold := make(map[string]interface{})
		tmpThreshold["threshold"] = threshold.Threshold
		tmpThreshold["domain"] = threshold.Domain
		resultThresholds = append(resultThresholds, tmpThreshold)
	}
	_ = d.Set("thresholds", resultThresholds)

	ccGeoIpPolicys, err := antiddosService.DescribeCcGeoIPBlockConfigList(ctx, business, instanceId)
	if err != nil {
		return err
	}
	resultCCGeoIpPolicys := make([]map[string]interface{}, 0)
	for _, ccGeoIpPolicy := range ccGeoIpPolicys {
		tmpResultCCGeoIpPolicy := make(map[string]interface{})
		tmpResultCCGeoIpPolicy["action"] = *ccGeoIpPolicy.Action
		tmpResultCCGeoIpPolicy["area_list"] = ccGeoIpPolicy.AreaList
		tmpResultCCGeoIpPolicy["create_time"] = *ccGeoIpPolicy.CreateTime
		tmpResultCCGeoIpPolicy["modify_time"] = *ccGeoIpPolicy.ModifyTime
		tmpResultCCGeoIpPolicy["domain"] = *ccGeoIpPolicy.Domain
		tmpResultCCGeoIpPolicy["protocol"] = *ccGeoIpPolicy.Protocol
		tmpResultCCGeoIpPolicy["region_type"] = *ccGeoIpPolicy.RegionType
		resultCCGeoIpPolicys = append(resultCCGeoIpPolicys, tmpResultCCGeoIpPolicy)
	}
	_ = d.Set("cc_geo_ip_policys", resultCCGeoIpPolicys)

	ccBlackWhiteIpList, err := antiddosService.DescribeCcBlackWhiteIpList(ctx, business, instanceId)
	if err != nil {
		return err
	}
	resultCCBlackWhiteIpList := make([]map[string]interface{}, 0)
	for _, ccBlackWhiteIp := range ccBlackWhiteIpList {
		tmpResultCCBlackWhiteIp := make(map[string]interface{})
		tmpResultCCBlackWhiteIp["black_white_ip"] = *ccBlackWhiteIp.BlackWhiteIp
		tmpResultCCBlackWhiteIp["create_time"] = *ccBlackWhiteIp.CreateTime
		tmpResultCCBlackWhiteIp["modify_time"] = *ccBlackWhiteIp.ModifyTime
		tmpResultCCBlackWhiteIp["domain"] = *ccBlackWhiteIp.Domain
		tmpResultCCBlackWhiteIp["protocol"] = *ccBlackWhiteIp.Protocol
		tmpResultCCBlackWhiteIp["type"] = *ccBlackWhiteIp.Type
		resultCCBlackWhiteIpList = append(resultCCBlackWhiteIpList, tmpResultCCBlackWhiteIp)
	}
	_ = d.Set("cc_black_white_ips", resultCCBlackWhiteIpList)

	ccPrecisionPlyList, err := antiddosService.DescribeCCPrecisionPlyList(ctx, business, instanceId)
	if err != nil {
		return err
	}
	resultCCPrecisionPlyList := make([]map[string]interface{}, 0)
	for _, ccPrecisionPly := range ccPrecisionPlyList {
		tmpResultCCPrecisionPly := make(map[string]interface{})
		tmpResultCCPrecisionPly["policy_action"] = *ccPrecisionPly.PolicyAction
		tmpResultCCPrecisionPly["policy_id"] = *ccPrecisionPly.PolicyId
		tmpResultCCPrecisionPly["domain"] = *ccPrecisionPly.Domain
		tmpResultCCPrecisionPly["ip"] = *ccPrecisionPly.Ip
		tmpResultCCPrecisionPly["protocol"] = *ccPrecisionPly.Protocol
		policys := make([]map[string]interface{}, 0)
		for _, policy := range ccPrecisionPly.PolicyList {
			tmpPolicy := make(map[string]interface{})
			tmpPolicy["field_name"] = *policy.FieldName
			tmpPolicy["field_type"] = *policy.FieldType
			tmpPolicy["value"] = *policy.Value
			tmpPolicy["value_operator"] = *policy.ValueOperator
			policys = append(policys, tmpPolicy)
		}
		tmpResultCCPrecisionPly["policys"] = policys
		resultCCPrecisionPlyList = append(resultCCPrecisionPlyList, tmpResultCCPrecisionPly)
	}
	_ = d.Set("cc_precision_policys", resultCCPrecisionPlyList)

	ccLevelList, err := antiddosService.DescribeCCLevelList(ctx, business, instanceId)
	if err != nil {
		return err
	}

	resultCCReqLimitList := make([]map[string]interface{}, 0)
	for _, ccLevelItem := range ccLevelList {
		ccLevelIp := *ccLevelItem.Ip
		ccLevelProtocol := *ccLevelItem.Protocol
		ccLevel := *ccLevelItem.Level
		ccLevelDomain := *ccLevelItem.Domain
		ccLevelInstanceId := *ccLevelItem.InstanceId
		tmpResultCCReqLimitList := make(map[string]interface{})
		tmpResultCCReqLimitList["ip"] = ccLevelIp
		tmpResultCCReqLimitList["protocol"] = ccLevelProtocol
		tmpResultCCReqLimitList["level"] = ccLevel
		tmpResultCCReqLimitList["domain"] = ccLevelDomain
		tmpResultCCReqLimitList["instance_id"] = ccLevelInstanceId
		ccReqLimitPolicyList, err := antiddosService.DescribeCCReqLimitPolicyList(ctx, business, instanceId)
		if err != nil {
			return err
		}
		resultCCReqLimitPolicyList := make([]map[string]interface{}, 0)
		for _, ccReqLimitPolicy := range ccReqLimitPolicyList {
			if ccLevelDomain != *ccReqLimitPolicy.Domain {
				continue
			}
			policyRecord := make(map[string]interface{})
			policyRecord["action"] = *ccReqLimitPolicy.PolicyRecord.Action
			policyRecord["cookie"] = *ccReqLimitPolicy.PolicyRecord.Cookie
			policyRecord["execute_duration"] = *ccReqLimitPolicy.PolicyRecord.ExecuteDuration
			policyRecord["mode"] = *ccReqLimitPolicy.PolicyRecord.Mode
			policyRecord["period"] = *ccReqLimitPolicy.PolicyRecord.Period
			policyRecord["request_num"] = *ccReqLimitPolicy.PolicyRecord.RequestNum
			policyRecord["uri"] = *ccReqLimitPolicy.PolicyRecord.Uri
			policyRecord["user_agent"] = *ccReqLimitPolicy.PolicyRecord.UserAgent
			resultCCReqLimitPolicyList = append(resultCCReqLimitPolicyList, policyRecord)
		}
		tmpResultCCReqLimitList["policys"] = resultCCReqLimitPolicyList
		resultCCReqLimitList = append(resultCCReqLimitList, tmpResultCCReqLimitList)
	}

	_ = d.Set("cc_precision_req_limits", resultCCReqLimitList)
	return nil
}

func resourceTencentCloudDayuCCPolicyV2Update(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dayu_cc_policy_v2.update")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	items := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(items) < 2 {
		return fmt.Errorf("broken ID of DDoS policy")
	}
	instanceId := items[0]
	business := items[1]
	antiddosService := AntiddosService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	protectThresholdConfig, err := antiddosService.DescribeListProtectThresholdConfig(ctx, instanceId)
	if err != nil {
		return err
	}
	instanceDetailList := protectThresholdConfig.InstanceDetailList
	if len(instanceDetailList) != 1 {
		return fmt.Errorf("Can not fetch Eip.")
	}
	ip := *instanceDetailList[0].EipList[0]
	if d.HasChange("thresholds") {
		oldthresholdMap := make(map[string]int)
		newthresholdMap := make(map[string]int)
		thresholdList, err := antiddosService.DescribeCCThresholdList(ctx, business, instanceId)
		if err != nil {
			return err
		}
		for _, threshold := range thresholdList {
			key := *threshold.Domain + "_" + fmt.Sprint(*threshold.Threshold)
			oldthresholdMap[key] = 1
		}
		thresholds := d.Get("thresholds").([]interface{})
		for _, threshold := range thresholds {
			thresholdMap := threshold.(map[string]interface{})
			subDomain := thresholdMap["domain"].(string)
			subThreshold := thresholdMap["threshold"].(int)
			key := subDomain + "_" + fmt.Sprint(subThreshold)
			newthresholdMap[key] = 1
			if oldthresholdMap[key] == 0 {
				err := antiddosService.ModifyCCThresholdPolicy(ctx, instanceId, "http", ip, subDomain, subThreshold)
				if err != nil {
					return err
				}
			}
		}
		for _, threshold := range thresholdList {
			key := *threshold.Domain + "_" + fmt.Sprint(*threshold.Threshold)
			if newthresholdMap[key] == 0 {
				err := antiddosService.DeleteCCLevelPolicy(ctx, instanceId, ip, *threshold.Domain)
				if err != nil {
					return err
				}
				err = antiddosService.DeleteCCThresholdPolicy(ctx, instanceId, ip, *threshold.Domain)
				if err != nil {
					return err
				}
			}
		}
	}

	if d.HasChange("cc_geo_ip_policys") {
		newCCGeoIpPolicys := d.Get("cc_geo_ip_policys").([]interface{})
		oldCCGeoIpPolicys, err := antiddosService.DescribeCcGeoIPBlockConfigList(ctx, business, instanceId)
		if err != nil {
			return err
		}
		oldCCGeoIpPolicyMap := make(map[string]int)
		newCCGeoIpPolicyMap := make(map[string]int)

		for _, oldCCGeoIpPolicy := range oldCCGeoIpPolicys {
			action := *oldCCGeoIpPolicy.Action
			regionType := *oldCCGeoIpPolicy.RegionType
			domain := *oldCCGeoIpPolicy.Domain
			protocol := *oldCCGeoIpPolicy.Protocol
			key := action + "_" + regionType + "_" + domain + "_" + protocol
			oldCCGeoIpPolicyMap[key] = 1

		}
		for _, newCCGeoIpPolicyItem := range newCCGeoIpPolicys {
			newCCGeoIpPolicyItemMap := newCCGeoIpPolicyItem.(map[string]interface{})
			action := newCCGeoIpPolicyItemMap["action"].(string)
			regionType := newCCGeoIpPolicyItemMap["region_type"].(string)
			domain := newCCGeoIpPolicyItemMap["domain"].(string)
			protocol := newCCGeoIpPolicyItemMap["protocol"].(string)
			key := action + "_" + regionType + "_" + domain + "_" + protocol
			newCCGeoIpPolicyMap[key] = 1
			if oldCCGeoIpPolicyMap[key] == 0 {
				ccGeoIPBlockConfig := antiddos.CcGeoIPBlockConfig{
					Action:     &action,
					RegionType: &regionType,
				}
				err := antiddosService.CreateCcGeoIPBlockConfig(ctx, instanceId, protocol, ip, domain, ccGeoIPBlockConfig)
				if err != nil {
					return err
				}
			}
		}
		for _, oldCCGeoIpPolicy := range oldCCGeoIpPolicys {
			action := *oldCCGeoIpPolicy.Action
			regionType := *oldCCGeoIpPolicy.RegionType
			domain := *oldCCGeoIpPolicy.Domain
			protocol := *oldCCGeoIpPolicy.Protocol
			key := action + "_" + regionType + "_" + domain + "_" + protocol
			if newCCGeoIpPolicyMap[key] == 0 {
				areaInt64List := make([]*int64, 0)
				for _, area := range oldCCGeoIpPolicy.AreaList {
					ateaInt64 := int64(*area)
					areaInt64List = append(areaInt64List, &ateaInt64)
				}
				ccGeoIPBlockConfig := antiddos.CcGeoIPBlockConfig{
					Action:     oldCCGeoIpPolicy.Action,
					AreaList:   areaInt64List,
					Id:         oldCCGeoIpPolicy.PolicyId,
					RegionType: oldCCGeoIpPolicy.RegionType,
				}
				_ = antiddosService.DeleteCcGeoIPBlockConfig(ctx, instanceId, ccGeoIPBlockConfig)
			}
		}

	}

	if d.HasChange("cc_black_white_ips") {
		newCCBlackWhiteIps := d.Get("cc_black_white_ips").([]interface{})
		oldCCBlackWhiteIps, err := antiddosService.DescribeCcBlackWhiteIpList(ctx, business, instanceId)
		if err != nil {
			return err
		}
		oldCCBlackWhiteIpMap := make(map[string]int)
		newCCBlackWhiteIpMap := make(map[string]int)
		oldCCBlackWhiteIpPolicyIdMap := make(map[string]string)

		for _, oldCCBlackWhiteIp := range oldCCBlackWhiteIps {
			blackWhiteIp := *oldCCBlackWhiteIp.BlackWhiteIp
			blackWhiteIpType := *oldCCBlackWhiteIp.Type
			domain := *oldCCBlackWhiteIp.Domain
			protocol := *oldCCBlackWhiteIp.Protocol
			key := blackWhiteIpType + "_" + blackWhiteIp + "_" + domain + "_" + protocol
			oldCCBlackWhiteIpMap[key] = 1
			oldCCBlackWhiteIpPolicyIdMap[key] = *oldCCBlackWhiteIp.PolicyId

		}
		for _, newCCBlackWhiteIpItem := range newCCBlackWhiteIps {
			newCCBlackWhiteIpItemMap := newCCBlackWhiteIpItem.(map[string]interface{})
			blackWhiteIp := newCCBlackWhiteIpItemMap["black_white_ip"].(string)
			blackWhiteIpType := newCCBlackWhiteIpItemMap["type"].(string)
			domain := newCCBlackWhiteIpItemMap["domain"].(string)
			protocol := newCCBlackWhiteIpItemMap["protocol"].(string)
			key := blackWhiteIpType + "_" + blackWhiteIp + "_" + domain + "_" + protocol
			newCCBlackWhiteIpMap[key] = 1
			if oldCCBlackWhiteIpMap[key] == 0 {
				err := antiddosService.CreateCcBlackWhiteIpList(ctx, instanceId, protocol, ip, domain, blackWhiteIpType, []string{blackWhiteIp})
				if err != nil {
					return err
				}
			}
		}
		for _, oldCCBlackWhiteIp := range oldCCBlackWhiteIps {
			blackWhiteIp := *oldCCBlackWhiteIp.BlackWhiteIp
			blackWhiteIpType := *oldCCBlackWhiteIp.Type
			domain := *oldCCBlackWhiteIp.Domain
			protocol := *oldCCBlackWhiteIp.Protocol
			key := blackWhiteIpType + "_" + blackWhiteIp + "_" + domain + "_" + protocol
			if newCCBlackWhiteIpMap[key] == 0 {
				err := antiddosService.DeleteCcBlackWhiteIpList(ctx, instanceId, oldCCBlackWhiteIpPolicyIdMap[key])
				if err != nil {
					return err
				}
			}
		}
	}

	if d.HasChange("cc_precision_policys") {
		newccPrecisionPolicys := d.Get("cc_precision_policys").([]interface{})
		oldCCPrecisionPolicys, err := antiddosService.DescribeCCPrecisionPlyList(ctx, business, instanceId)
		if err != nil {
			return err
		}
		oldCCPrecisionPolicyMap := make(map[string]int)
		newCCPrecisionPolicyMap := make(map[string]int)
		oldCCPrecisionPolicyPolicyIdMap := make(map[string]string)

		for _, oldCCPrecisionPolicy := range oldCCPrecisionPolicys {
			domain := *oldCCPrecisionPolicy.Domain
			protocol := *oldCCPrecisionPolicy.Protocol
			policyAction := *oldCCPrecisionPolicy.PolicyAction
			policyList := oldCCPrecisionPolicy.PolicyList
			policyId := oldCCPrecisionPolicy.PolicyId
			for _, policy := range policyList {
				fieldName := *policy.FieldName
				fieldType := *policy.FieldType
				value := *policy.Value
				valueOperator := *policy.ValueOperator
				key := domain + "_" + protocol + "+" + policyAction + "_" + fieldName + "_" + fieldType + "_" + value + "_" + valueOperator
				oldCCPrecisionPolicyMap[key] = 1
				oldCCPrecisionPolicyPolicyIdMap[key] = *policyId
			}

		}
		for _, newccPrecisionPolicyItem := range newccPrecisionPolicys {
			newccPrecisionPolicyItemMap := newccPrecisionPolicyItem.(map[string]interface{})
			policyAction := newccPrecisionPolicyItemMap["policy_action"].(string)
			protocol := newccPrecisionPolicyItemMap["protocol"].(string)
			domain := newccPrecisionPolicyItemMap["domain"].(string)

			policys := newccPrecisionPolicyItemMap["policys"].([]interface{})
			for _, policy := range policys {
				policyItemMap := policy.(map[string]interface{})
				fieldName := policyItemMap["field_name"].(string)
				fieldType := policyItemMap["field_type"].(string)
				value := policyItemMap["value"].(string)
				valueOperator := policyItemMap["value_operator"].(string)
				key := domain + "_" + protocol + "+" + policyAction + "_" + fieldName + "_" + fieldType + "_" + value + "_" + valueOperator
				newCCPrecisionPolicyMap[key] = 1
				if oldCCPrecisionPolicyMap[key] == 0 {
					tmpCCPrecisionPlyRecord := antiddos.CCPrecisionPlyRecord{}
					tmpCCPrecisionPlyRecord.FieldName = &fieldName
					tmpCCPrecisionPlyRecord.FieldType = &fieldType
					tmpCCPrecisionPlyRecord.Value = &value
					tmpCCPrecisionPlyRecord.ValueOperator = &valueOperator
					err := antiddosService.CreateCCPrecisionPolicy(ctx, instanceId, protocol, ip, domain, policyAction, []*antiddos.CCPrecisionPlyRecord{&tmpCCPrecisionPlyRecord})
					if err != nil {
						return err
					}
				}
			}

		}
		for _, oldCCPrecisionPolicy := range oldCCPrecisionPolicys {
			policyAction := *oldCCPrecisionPolicy.PolicyAction
			policyId := *oldCCPrecisionPolicy.PolicyId
			policyList := oldCCPrecisionPolicy.PolicyList
			domain := *oldCCPrecisionPolicy.Domain
			protocol := *oldCCPrecisionPolicy.Protocol
			for _, policy := range policyList {
				fieldName := *policy.FieldName
				fieldType := *policy.FieldType
				value := *policy.Value
				valueOperator := *policy.ValueOperator
				key := domain + "_" + protocol + "+" + policyAction + "_" + fieldName + "_" + fieldType + "_" + value + "_" + valueOperator
				if newCCPrecisionPolicyMap[key] == 0 {
					err := antiddosService.DeleteCCPrecisionPolicy(ctx, instanceId, policyId)
					if err != nil {
						return err
					}
				}
			}

		}
	}

	if d.HasChange("cc_precision_req_limits") {
		newCCPrecisionReqLimits := d.Get("cc_precision_req_limits").([]interface{})
		oldCCPrecisionReqLimitPolicyMap := make(map[string]int)
		newCCPrecisionReqLimitPolicyMap := make(map[string]int)
		newDomainSetMap := make(map[string]int)

		ccReqLimitPolicyList, err := antiddosService.DescribeCCReqLimitPolicyList(ctx, business, instanceId)
		if err != nil {
			return err
		}
		for _, ccReqLimitPolicy := range ccReqLimitPolicyList {
			action := *ccReqLimitPolicy.PolicyRecord.Action
			executeDuration := *ccReqLimitPolicy.PolicyRecord.ExecuteDuration
			mode := *ccReqLimitPolicy.PolicyRecord.Mode
			period := *ccReqLimitPolicy.PolicyRecord.Period
			requestNum := *ccReqLimitPolicy.PolicyRecord.RequestNum
			uri := *ccReqLimitPolicy.PolicyRecord.Uri
			cookie := *ccReqLimitPolicy.PolicyRecord.Cookie
			userAgent := *ccReqLimitPolicy.PolicyRecord.UserAgent
			domain := *ccReqLimitPolicy.Domain
			protocol := *ccReqLimitPolicy.Protocol
			key := domain + "_" + protocol + "_" + action + "_" + fmt.Sprint(executeDuration) + "_" + fmt.Sprint(mode) + "_" + fmt.Sprint(period) + "_" + fmt.Sprint(requestNum) + "_" + uri + "_" + cookie + "_" + userAgent
			oldCCPrecisionReqLimitPolicyMap[key] = 1
		}

		for _, newCCPrecisionReqLimit := range newCCPrecisionReqLimits {
			newCCPrecisionReqLimitMap := newCCPrecisionReqLimit.(map[string]interface{})
			domain := newCCPrecisionReqLimitMap["domain"].(string)
			protocol := newCCPrecisionReqLimitMap["protocol"].(string)
			newCCLevel := newCCPrecisionReqLimitMap["level"].(string)
			newDomainSetMap[domain] = 1
			policys := newCCPrecisionReqLimitMap["policys"].([]interface{})
			for _, policy := range policys {
				policyItemMap := policy.(map[string]interface{})
				action := policyItemMap["action"].(string)
				executeDuration := policyItemMap["execute_duration"].(int)
				mode := policyItemMap["mode"].(string)
				period := policyItemMap["period"].(int)
				requestNum := policyItemMap["request_num"].(int)
				uri := policyItemMap["uri"].(string)
				cookie := policyItemMap["cookie"].(string)
				userAgent := policyItemMap["user_agent"].(string)
				key := domain + "_" + protocol + "_" + action + "_" + fmt.Sprint(executeDuration) + "_" + fmt.Sprint(mode) + "_" + fmt.Sprint(period) + "_" + fmt.Sprint(requestNum) + "_" + uri + "_" + cookie + "_" + userAgent
				newCCPrecisionReqLimitPolicyMap[key] = 1
				if oldCCPrecisionReqLimitPolicyMap[key] == 0 {
					ccPolicyRecord := antiddos.CCReqLimitPolicyRecord{
						Action:          &action,
						ExecuteDuration: helper.IntUint64(executeDuration),
						Mode:            &mode,
						Period:          helper.IntUint64(period),
						RequestNum:      helper.IntUint64(requestNum),
					}
					if uri != "" {
						ccPolicyRecord.Uri = &uri
					} else if cookie != "" {
						ccPolicyRecord.Cookie = &cookie
					} else if userAgent != "" {
						ccPolicyRecord.UserAgent = &userAgent
					}
					err := antiddosService.CreateCCReqLimitPolicy(ctx, instanceId, protocol, ip, domain, ccPolicyRecord)
					if err != nil {
						return err
					}
				}
			}
			err = antiddosService.ModifyCCLevelPolicy(ctx, instanceId, ip, domain, protocol, newCCLevel)
			if err != nil {
				return err
			}

		}

		for _, ccReqLimitPolicy := range ccReqLimitPolicyList {
			protocol := *ccReqLimitPolicy.Protocol
			domain := *ccReqLimitPolicy.Domain
			action := *ccReqLimitPolicy.PolicyRecord.Action
			executeDuration := *ccReqLimitPolicy.PolicyRecord.ExecuteDuration
			mode := *ccReqLimitPolicy.PolicyRecord.Mode
			period := *ccReqLimitPolicy.PolicyRecord.Period
			requestNum := *ccReqLimitPolicy.PolicyRecord.RequestNum
			uri := *ccReqLimitPolicy.PolicyRecord.Uri
			cookie := *ccReqLimitPolicy.PolicyRecord.Cookie
			userAgent := *ccReqLimitPolicy.PolicyRecord.UserAgent
			key := domain + "_" + protocol + "_" + action + "_" + fmt.Sprint(executeDuration) + "_" + fmt.Sprint(mode) + "_" + fmt.Sprint(period) + "_" + fmt.Sprint(requestNum) + "_" + uri + "_" + cookie + "_" + userAgent
			if newCCPrecisionReqLimitPolicyMap[key] == 0 {
				err := antiddosService.DeleteCCRequestLimitPolicy(ctx, instanceId, *ccReqLimitPolicy.PolicyId)
				if err != nil {
					return err
				}
			}
			if newDomainSetMap[domain] == 0 {
				err := antiddosService.DeleteCCLevelPolicy(ctx, instanceId, ip, domain)
				if err != nil {
					return err
				}
			}
		}
	}

	return resourceTencentCloudDayuCCPolicyV2Read(d, meta)
}

func resourceTencentCloudDayuCCPolicyV2Delete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dayu_cc_policy_v2.delete")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	items := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(items) < 2 {
		return fmt.Errorf("broken ID of DDoS policy")
	}
	instanceId := items[0]
	business := items[1]
	antiddosService := AntiddosService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	thresholdList, err := antiddosService.DescribeCCThresholdList(ctx, business, instanceId)
	if err != nil {
		return err
	}
	for _, threshold := range thresholdList {
		err := antiddosService.DeleteCCLevelPolicy(ctx, instanceId, *threshold.Ip, *threshold.Domain)
		if err != nil {
			return err
		}
		err = antiddosService.DeleteCCThresholdPolicy(ctx, instanceId, *threshold.Ip, *threshold.Domain)
		if err != nil {
			return err
		}
	}

	ccGeoIpPolicys, err := antiddosService.DescribeCcGeoIPBlockConfigList(ctx, business, instanceId)
	if err != nil {
		return err
	}
	for _, ccGeoIpPolicy := range ccGeoIpPolicys {
		areaInt64List := make([]*int64, 0)
		for _, area := range ccGeoIpPolicy.AreaList {
			ateaInt64 := int64(*area)
			areaInt64List = append(areaInt64List, &ateaInt64)
		}
		ccGeoIPBlockConfig := antiddos.CcGeoIPBlockConfig{
			Action:     ccGeoIpPolicy.Action,
			AreaList:   areaInt64List,
			Id:         ccGeoIpPolicy.PolicyId,
			RegionType: ccGeoIpPolicy.RegionType,
		}
		_ = antiddosService.DeleteCcGeoIPBlockConfig(ctx, instanceId, ccGeoIPBlockConfig)
	}

	ccBlackWhiteIpList, err := antiddosService.DescribeCcBlackWhiteIpList(ctx, business, instanceId)
	if err != nil {
		return err
	}
	for _, ccBlackWhiteIp := range ccBlackWhiteIpList {
		err := antiddosService.DeleteCcBlackWhiteIpList(ctx, instanceId, *ccBlackWhiteIp.PolicyId)
		if err != nil {
			return err
		}
	}

	ccPrecisionPlyList, err := antiddosService.DescribeCCPrecisionPlyList(ctx, business, instanceId)
	if err != nil {
		return err
	}
	for _, ccPrecisionPly := range ccPrecisionPlyList {
		err := antiddosService.DeleteCCPrecisionPolicy(ctx, instanceId, *ccPrecisionPly.PolicyId)
		if err != nil {
			return err
		}
	}
	ccLevelPolicyDomainMap := make(map[string]interface{})
	ccLevelList, err := antiddosService.DescribeCCLevelList(ctx, business, instanceId)
	if err != nil {
		return err
	}
	for _, ccLevel := range ccLevelList {
		key := *ccLevel.Ip + "_" + *ccLevel.Domain
		ccLevelPolicyDomainMap[key] = 1

	}
	ccReqLimitPolicyList, err := antiddosService.DescribeCCReqLimitPolicyList(ctx, business, instanceId)
	if err != nil {
		return err
	}
	for _, ccReqLimitPolicy := range ccReqLimitPolicyList {
		key := *ccReqLimitPolicy.Ip + "_" + *ccReqLimitPolicy.Domain
		ccLevelPolicyDomainMap[key] = 1
	}
	for item := range ccLevelPolicyDomainMap {
		itemList := strings.Split(item, "_")
		err := antiddosService.DeleteCCLevelPolicy(ctx, instanceId, itemList[0], itemList[1])
		if err != nil {
			return err
		}
	}

	return nil
}
