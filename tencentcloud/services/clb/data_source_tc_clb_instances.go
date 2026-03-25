package clb

import (
	"context"
	"encoding/json"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	clb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/clb/v20180317"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudClbInstances() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudClbInstancesRead,

		Schema: map[string]*schema.Schema{
			"clb_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "ID of the CLB to be queried.",
			},
			"network_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: tccommon.ValidateAllowedStringValue(CLB_NETWORK_TYPE),
				Description:  "Type of CLB instance, and available values include `OPEN` and `INTERNAL`.",
			},
			"clb_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Name of the CLB to be queried.",
			},
			"project_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Project ID of the CLB.",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
			"master_zone": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Master available zone id.",
			},
			"clb_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A list of cloud load balancers. Each element contains the following attributes:",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"clb_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of CLB.",
						},
						"clb_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of CLB.",
						},
						"network_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Types of CLB.",
						},
						"project_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "ID of the project.",
						},
						"cluster_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of the cluster.",
						},
						"clb_vips": {
							Type:        schema.TypeList,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "The virtual service address table of the CLB.",
						},
						"status": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The status of CLB.",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Create time of the CLB.",
						},
						"status_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Latest state transition time of CLB.",
						},
						"vpc_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of the VPC.",
						},
						"subnet_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of the subnet.",
						},
						"security_groups": {
							Type:        schema.TypeList,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "ID set of the security groups.",
						},
						"target_region_info_region": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Region information of backend service are attached the CLB.",
						},
						"target_region_info_vpc_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "VpcId information of backend service are attached the CLB.",
						},
						"tags": {
							Type:        schema.TypeMap,
							Computed:    true,
							Description: "The available tags within this CLB.",
						},
						"address_ip_version": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "IP version, only applicable to open CLB. Valid values are `IPV4`, `IPV6` and `IPv6FullChain`.",
						},
						"vip_isp": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Network operator, only applicable to open CLB. Valid values are `CMCC`(China Mobile), `CTCC`(Telecom), `CUCC`(China Unicom) and `BGP`. If this ISP is specified, network billing method can only use the bandwidth package billing (BANDWIDTH_PACKAGE).",
						},
						"internet_charge_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Internet charge type, only applicable to open CLB. Valid values are `TRAFFIC_POSTPAID_BY_HOUR`, `BANDWIDTH_POSTPAID_BY_HOUR` and `BANDWIDTH_PACKAGE`.",
						},
						"internet_bandwidth_max_out": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Max bandwidth out, only applicable to open CLB. Valid value ranges is [1, 2048]. Unit is MB.",
						},
						"zone_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Available zone unique id(numerical representation), This field maybe null, means cannot get a valid value.",
						},
						"zone": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Available zone unique id(string representation), This field maybe null, means cannot get a valid value.",
						},
						"zone_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Available zone name, This field maybe null, means cannot get a valid value.",
						},
						"zone_region": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Region that this available zone belong to, This field maybe null, means cannot get a valid value.",
						},
						"local_zone": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether this available zone is local zone, This field maybe null, means cannot get a valid value.",
						},
						"numerical_vpc_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "VPC ID in a numeric form. Note: This field may return null, indicating that no valid values can be obtained.",
						},
						"zones": {
							Type:        schema.TypeList,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "Zones where rules are deployed for VPC internal load balancers with nearby access mode. Note: This field may return null, indicating no valid values can be obtained.",
						},
						// Basic info fields
						"forward": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "CLB type identifier, 1: CLB, 0: Classic CLB.",
						},
						"domain": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "CLB domain (only for public network Classic CLB), gradually deprecated.",
						},
						"load_balancer_domain": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Domain of the CLB instance.",
						},
						// Network config fields
						"address_ipv6": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "IPv6 address of the CLB instance.",
						},
						"ipv6_mode": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "IPv6 mode when IP version is ipv6, IPv6Nat64 or IPv6FullChain.",
						},
						"mix_ip_target": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "IPv6FullChain CLB layer-7 listener supports mixed binding of IPv4/IPv6 targets.",
						},
						"anycast_zone": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Anycast CLB publishing region, returns empty string for non-anycast CLB.",
						},
						"egress": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Network egress.",
						},
						"local_bgp": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether the IP type is local BGP.",
						},
						// Billing and lifecycle fields
						"charge_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Billing type, PREPAID: Prepaid, POSTPAID_BY_HOUR: Pay-as-you-go.",
						},
						"expire_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Expiration time of the CLB instance, only for prepaid CLB, format: YYYY-MM-DD HH:mm:ss.",
						},
						"prepaid_period": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Prepaid purchase period, unit: month.",
						},
						"prepaid_renew_flag": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Prepaid renewal flag, NOTIFY_AND_AUTO_RENEW: Notify and auto-renew, NOTIFY_AND_MANUAL_RENEW: Notify but not auto-renew, DISABLE_NOTIFY_AND_MANUAL_RENEW: No notification and not auto-renew.",
						},
						// Log config fields
						"log_set_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Log service (CLS) log set ID.",
						},
						"log_topic_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Log service (CLS) log topic ID.",
						},
						"health_log_set_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Log service (CLS) health check log set ID.",
						},
						"health_log_topic_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Log service (CLS) health check log topic ID.",
						},
						// Security and isolation fields
						"open_bgp": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Anti-DDoS Pro LB identifier, 1: Anti-DDoS Pro, 0: Not Anti-DDoS Pro.",
						},
						"snat": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether SNAT is enabled.",
						},
						"snat_pro": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether SnatPro is enabled.",
						},
						"snat_ips": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "SnatIp list after enabling SnatPro (JSON format).",
						},
						"isolation": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Whether isolated, 0: Not isolated, 1: Isolated.",
						},
						"isolated_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Time when the CLB instance was isolated, format: YYYY-MM-DD HH:mm:ss.",
						},
						"is_block": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether the VIP is blocked.",
						},
						"is_block_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Time of blocking or unblocking, format: YYYY-MM-DD HH:mm:ss.",
						},
						"is_ddos": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether Anti-DDoS Pro can be bound.",
						},
						// Performance and capacity fields
						"sla_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Performance capacity type specification (clb.c1.small/clb.c2.medium/clb.c3.small/clb.c3.medium/clb.c4.small/clb.c4.medium/clb.c4.large/clb.c4.xlarge or empty string).",
						},
						"exclusive": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Whether the instance type is exclusive, 1: Exclusive, 0: Not exclusive.",
						},
						"target_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Number of bound backend services.",
						},
						// Cluster and deployment fields
						"cluster_ids": {
							Type:        schema.TypeList,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "Cluster ID list.",
						},
						"cluster_tag": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Layer-7 exclusive tag.",
						},
						"nfv_info": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Whether CLB is NFV, empty: No, l7nfv: Layer-7 is NFV.",
						},
						"backup_zone_set": {
							Type:        schema.TypeList,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeMap},
							Description: "Backup zone list, each element contains zone_id/zone/zone_name/zone_region/local_zone.",
						},
						"available_zone_affinity_info": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Available zone forwarding affinity information (JSON format).",
						},
						// Advanced config fields
						"config_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "CLB dimension personalized configuration ID.",
						},
						"load_balancer_pass_to_target": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether backend services allow traffic from CLB.",
						},
						"attribute_flags": {
							Type:        schema.TypeList,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "CLB attribute flags array.",
						},
						"exclusive_cluster": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Internal exclusive cluster information (JSON format).",
						},
						"extra_info": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Reserved field, generally no need to pay attention (JSON format).",
						},
						"associate_endpoint": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Endpoint ID associated with the CLB instance.",
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentCloudClbInstancesRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_clb_instances.read")()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		clbService = ClbService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		clbs       []*clb.LoadBalancer
	)

	params := make(map[string]interface{})
	if v, ok := d.GetOk("clb_id"); ok {
		params["clb_id"] = v.(string)
	}

	if v, ok := d.GetOk("clb_name"); ok {
		params["clb_name"] = v.(string)
	}

	if v, ok := d.GetOkExists("project_id"); ok {
		params["project_id"] = v.(int)
	}

	if v, ok := d.GetOk("network_type"); ok {
		params["network_type"] = v.(string)
	}

	if v, ok := d.GetOk("master_zone"); ok {
		params["master_zone"] = v.(string)
	}

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		results, e := clbService.DescribeLoadBalancerByFilter(ctx, params)
		if e != nil {
			return tccommon.RetryError(e)
		}

		clbs = results
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s read CLB instances failed, reason:%+v", logId, err)
		return err
	}

	clbList := make([]map[string]interface{}, 0, len(clbs))
	ids := make([]string, 0, len(clbs))
	for _, clbInstance := range clbs {
		mapping := map[string]interface{}{
			"clb_id":                    clbInstance.LoadBalancerId,
			"clb_name":                  clbInstance.LoadBalancerName,
			"network_type":              clbInstance.LoadBalancerType,
			"status":                    clbInstance.Status,
			"create_time":               clbInstance.CreateTime,
			"status_time":               clbInstance.StatusTime,
			"project_id":                clbInstance.ProjectId,
			"vpc_id":                    clbInstance.VpcId,
			"subnet_id":                 clbInstance.SubnetId,
			"clb_vips":                  helper.StringsInterfaces(clbInstance.LoadBalancerVips),
			"target_region_info_region": clbInstance.TargetRegionInfo.Region,
			"target_region_info_vpc_id": clbInstance.TargetRegionInfo.VpcId,
			"address_ip_version":        clbInstance.AddressIPVersion,
			"vip_isp":                   clbInstance.VipIsp,
			"security_groups":           helper.StringsInterfaces(clbInstance.SecureGroups),
		}

		if clbInstance.ClusterIds != nil && len(clbInstance.ClusterIds) > 0 {
			mapping["cluster_id"] = *clbInstance.ClusterIds[0]
		}

		if clbInstance.NetworkAttributes != nil {
			mapping["internet_charge_type"] = *clbInstance.NetworkAttributes.InternetChargeType
			mapping["internet_bandwidth_max_out"] = *clbInstance.NetworkAttributes.InternetMaxBandwidthOut
		}

		if clbInstance.MasterZone != nil {
			mapping["zone_id"] = *clbInstance.MasterZone.ZoneId
			mapping["zone"] = *clbInstance.MasterZone.Zone
			mapping["zone_name"] = *clbInstance.MasterZone.ZoneName
			mapping["zone_region"] = *clbInstance.MasterZone.ZoneRegion
			mapping["local_zone"] = *clbInstance.MasterZone.LocalZone
		}

		if clbInstance.Tags != nil {
			tags := make(map[string]interface{}, len(clbInstance.Tags))
			for _, t := range clbInstance.Tags {
				tags[*t.TagKey] = *t.TagValue
			}

			mapping["tags"] = tags
		}

		if clbInstance.NumericalVpcId != nil {
			mapping["numerical_vpc_id"] = clbInstance.NumericalVpcId
		}

		if clbInstance.Zones != nil {
			mapping["zones"] = helper.StringsInterfaces(clbInstance.Zones)
		}

		// Basic info fields
		if clbInstance.Forward != nil {
			mapping["forward"] = clbInstance.Forward
		}
		if clbInstance.Domain != nil {
			mapping["domain"] = clbInstance.Domain
		}
		if clbInstance.LoadBalancerDomain != nil {
			mapping["load_balancer_domain"] = clbInstance.LoadBalancerDomain
		}

		// Network config fields
		if clbInstance.AddressIPv6 != nil {
			mapping["address_ipv6"] = clbInstance.AddressIPv6
		}
		if clbInstance.IPv6Mode != nil {
			mapping["ipv6_mode"] = clbInstance.IPv6Mode
		}
		if clbInstance.MixIpTarget != nil {
			mapping["mix_ip_target"] = clbInstance.MixIpTarget
		}
		if clbInstance.AnycastZone != nil {
			mapping["anycast_zone"] = clbInstance.AnycastZone
		}
		if clbInstance.Egress != nil {
			mapping["egress"] = clbInstance.Egress
		}
		if clbInstance.LocalBgp != nil {
			mapping["local_bgp"] = clbInstance.LocalBgp
		}

		// Billing and lifecycle fields
		if clbInstance.ChargeType != nil {
			mapping["charge_type"] = clbInstance.ChargeType
		}
		if clbInstance.ExpireTime != nil {
			mapping["expire_time"] = clbInstance.ExpireTime
		}
		if clbInstance.PrepaidAttributes != nil {
			if clbInstance.PrepaidAttributes.Period != nil {
				mapping["prepaid_period"] = clbInstance.PrepaidAttributes.Period
			}
			if clbInstance.PrepaidAttributes.RenewFlag != nil {
				mapping["prepaid_renew_flag"] = clbInstance.PrepaidAttributes.RenewFlag
			}
		}

		// Log config fields
		if clbInstance.LogSetId != nil {
			mapping["log_set_id"] = clbInstance.LogSetId
		}
		if clbInstance.LogTopicId != nil {
			mapping["log_topic_id"] = clbInstance.LogTopicId
		}
		if clbInstance.HealthLogSetId != nil {
			mapping["health_log_set_id"] = clbInstance.HealthLogSetId
		}
		if clbInstance.HealthLogTopicId != nil {
			mapping["health_log_topic_id"] = clbInstance.HealthLogTopicId
		}

		// Security and isolation fields
		if clbInstance.OpenBgp != nil {
			mapping["open_bgp"] = clbInstance.OpenBgp
		}
		if clbInstance.Snat != nil {
			mapping["snat"] = clbInstance.Snat
		}
		if clbInstance.SnatPro != nil {
			mapping["snat_pro"] = clbInstance.SnatPro
		}
		if clbInstance.SnatIps != nil {
			snatIpsJSON, _ := json.Marshal(clbInstance.SnatIps)
			mapping["snat_ips"] = string(snatIpsJSON)
		}
		if clbInstance.Isolation != nil {
			mapping["isolation"] = clbInstance.Isolation
		}
		if clbInstance.IsolatedTime != nil {
			mapping["isolated_time"] = clbInstance.IsolatedTime
		}
		if clbInstance.IsBlock != nil {
			mapping["is_block"] = clbInstance.IsBlock
		}
		if clbInstance.IsBlockTime != nil {
			mapping["is_block_time"] = clbInstance.IsBlockTime
		}
		if clbInstance.IsDDos != nil {
			mapping["is_ddos"] = clbInstance.IsDDos
		}

		// Performance and capacity fields
		if clbInstance.SlaType != nil {
			mapping["sla_type"] = clbInstance.SlaType
		}
		if clbInstance.Exclusive != nil {
			mapping["exclusive"] = clbInstance.Exclusive
		}
		if clbInstance.TargetCount != nil {
			mapping["target_count"] = clbInstance.TargetCount
		}

		// Cluster and deployment fields
		if clbInstance.ClusterIds != nil {
			mapping["cluster_ids"] = helper.StringsInterfaces(clbInstance.ClusterIds)
		}
		if clbInstance.ClusterTag != nil {
			mapping["cluster_tag"] = clbInstance.ClusterTag
		}
		if clbInstance.NfvInfo != nil {
			mapping["nfv_info"] = clbInstance.NfvInfo
		}
		if clbInstance.BackupZoneSet != nil {
			backupZones := make([]map[string]interface{}, 0, len(clbInstance.BackupZoneSet))
			for _, zone := range clbInstance.BackupZoneSet {
				backupZone := make(map[string]interface{})
				if zone.ZoneId != nil {
					backupZone["zone_id"] = *zone.ZoneId
				}
				if zone.Zone != nil {
					backupZone["zone"] = *zone.Zone
				}
				if zone.ZoneName != nil {
					backupZone["zone_name"] = *zone.ZoneName
				}
				if zone.ZoneRegion != nil {
					backupZone["zone_region"] = *zone.ZoneRegion
				}
				if zone.LocalZone != nil {
					backupZone["local_zone"] = *zone.LocalZone
				}
				backupZones = append(backupZones, backupZone)
			}
			mapping["backup_zone_set"] = backupZones
		}
		if clbInstance.AvailableZoneAffinityInfo != nil {
			availableZoneAffinityJSON, _ := json.Marshal(clbInstance.AvailableZoneAffinityInfo)
			mapping["available_zone_affinity_info"] = string(availableZoneAffinityJSON)
		}

		// Advanced config fields
		if clbInstance.ConfigId != nil {
			mapping["config_id"] = clbInstance.ConfigId
		}
		if clbInstance.LoadBalancerPassToTarget != nil {
			mapping["load_balancer_pass_to_target"] = clbInstance.LoadBalancerPassToTarget
		}
		if clbInstance.AttributeFlags != nil {
			mapping["attribute_flags"] = helper.StringsInterfaces(clbInstance.AttributeFlags)
		}
		if clbInstance.ExclusiveCluster != nil {
			exclusiveClusterJSON, _ := json.Marshal(clbInstance.ExclusiveCluster)
			mapping["exclusive_cluster"] = string(exclusiveClusterJSON)
		}
		if clbInstance.ExtraInfo != nil {
			extraInfoJSON, _ := json.Marshal(clbInstance.ExtraInfo)
			mapping["extra_info"] = string(extraInfoJSON)
		}
		if clbInstance.AssociateEndpoint != nil {
			mapping["associate_endpoint"] = clbInstance.AssociateEndpoint
		}

		clbList = append(clbList, mapping)
		ids = append(ids, *clbInstance.LoadBalancerId)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	if e := d.Set("clb_list", clbList); e != nil {
		log.Printf("[CRITAL]%s provider set CLB list fail, reason:%+v", logId, e)
		return e
	}

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), clbList); e != nil {
			return e
		}
	}

	return nil
}
