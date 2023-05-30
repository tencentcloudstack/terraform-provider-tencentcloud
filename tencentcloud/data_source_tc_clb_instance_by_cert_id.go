/*
Use this data source to query detailed information of clb instance_by_cert_id

Example Usage

```hcl
data "tencentcloud_clb_instance_by_cert_id" "instance_by_cert_id" {
  cert_ids = ["3a6B5y8v"]
}
```
*/
package tencentcloud

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	clb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/clb/v20180317"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudClbInstanceByCertId() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudClbInstanceByCertIdRead,
		Schema: map[string]*schema.Schema{
			"cert_ids": {
				Required: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Server or client certificate ID.",
			},

			"cert_set": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Certificate ID and list of CLB instances associated with it.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cert_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Certificate ID.",
						},
						"load_balancers": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "List of CLB instances associated with certificate. Note: this field may return null, indicating that no valid values can be obtained.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"load_balancer_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "CLB instance ID.",
									},
									"load_balancer_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "CLB instance name.",
									},
									"load_balancer_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "CLB instance network type:OPEN: public network; INTERNAL: private network.",
									},
									"forward": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "CLB type identifier. Value range: 1 (CLB); 0 (classic CLB).",
									},
									"domain": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Domain name of the CLB instance. It is only available for public classic CLBs. This parameter will be discontinued soon. Please use LoadBalancerDomain instead. Note: This field may return null, indicating that no valid values can be obtained.",
									},
									"load_balancer_vips": {
										Type: schema.TypeSet,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
										Computed:    true,
										Description: "List of VIPs of a CLB instance. Note: This field may return null, indicating that no valid values can be obtained.",
									},
									"status": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "CLB instance status, including:0: creating; 1: running. Note: This field may return null, indicating that no valid values can be obtained.",
									},
									"create_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "CLB instance creation time. Note: This field may return null, indicating that no valid values can be obtained.",
									},
									"status_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Last status change time of a CLB instance. Note: This field may return null, indicating that no valid values can be obtained.",
									},
									"project_id": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "ID of the project to which a CLB instance belongs. 0: default project.",
									},
									"vpc_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "VPC ID Note: This field may return null, indicating that no valid values can be obtained.",
									},
									"open_bgp": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Protective CLB identifier. Value range: 1 (protective), 0 (non-protective). Note: This field may return null, indicating that no valid values can be obtained.",
									},
									"snat": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "SNAT is enabled for all private network classic CLB created before December 2016. Note: This field may return null, indicating that no valid values can be obtained.",
									},
									"isolation": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "0: not isolated; 1: isolated. Note: This field may return null, indicating that no valid values can be obtained.",
									},
									"log": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Log information. Only the public network CLB that have HTTP or HTTPS listeners can generate logs. Note: This field may return null, indicating that no valid values can be obtained.",
									},
									"subnet_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Subnet where a CLB instance resides (meaningful only for private network VPC CLB). Note: This field may return null, indicating that no valid values can be obtained.",
									},
									"tags": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "CLB instance tag information. Note: This field may return null, indicating that no valid values can be obtained.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"tag_key": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Tag key.",
												},
												"tag_value": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Tag value.",
												},
											},
										},
									},
									"secure_groups": {
										Type: schema.TypeSet,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
										Computed:    true,
										Description: "Security group of a CLB instance. Note: This field may return null, indicating that no valid values can be obtained.",
									},
									"target_region_info": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Basic information of a backend server bound to a CLB instance. Note: This field may return null, indicating that no valid values can be obtained.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"region": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Region of the target, such as ap-guangzhou.",
												},
												"vpc_id": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Network of the target, which is in the format of vpc-abcd1234 for VPC or 0 for basic network.",
												},
											},
										},
									},
									"anycast_zone": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Anycast CLB publishing region. For non-anycast CLB, this field returns an empty string. Note: This field may return null, indicating that no valid values can be obtained.",
									},
									"address_i_p_version": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "IP version. Valid values: ipv4, ipv6. Note: this field may return null, indicating that no valid values can be obtained.",
									},
									"numerical_vpc_id": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "VPC ID in a numeric form. Note: This field may return null, indicating that no valid values can be obtained.",
									},
									"vip_isp": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "ISP to which a CLB IP address belongs. Note: This field may return null, indicating that no valid values can be obtained.",
									},
									"master_zone": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Primary AZ. Note: This field may return null, indicating that no valid values can be obtained.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"zone_id": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: ".",
												},
												"zone": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Unique AZ ID in a numeric form, such as 100001. Note: This field may return null, indicating that no valid values can be obtained.",
												},
												"zone_name": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "AZ name, such as Guangzhou Zone 1. Note: This field may return null, indicating that no valid values can be obtained.",
												},
												"zone_region": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "AZ region, e.g., ap-guangzhou. Note: This field may return null, indicating that no valid values can be obtained.",
												},
												"local_zone": {
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "Whether the AZ is the LocalZone, e.g., false. Note: This field may return null, indicating that no valid values can be obtained.",
												},
												"edge_zone": {
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "Whether the AZ is an edge zone. Values: true, false. Note: This field may return null, indicating that no valid values can be obtained.",
												},
											},
										},
									},
									"backup_zone_set": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "backup zone.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"zone_id": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Unique AZ ID in a numeric form, such as 100001. Note: This field may return null, indicating that no valid values can be obtained.",
												},
												"zone": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Unique AZ ID in a numeric form, such as 100001. Note: This field may return null, indicating that no valid values can be obtained.",
												},
												"zone_name": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "AZ name, such as Guangzhou Zone 1. Note: This field may return null, indicating that no valid values can be obtained.",
												},
												"zone_region": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "AZ region, e.g., ap-guangzhou. Note: This field may return null, indicating that no valid values can be obtained.",
												},
												"local_zone": {
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "Whether the AZ is the LocalZone, e.g., false. Note: This field may return null, indicating that no valid values can be obtained.",
												},
												"edge_zone": {
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "Whether the AZ is an edge zone. Values: true, false. Note: This field may return null, indicating that no valid values can be obtained.",
												},
											},
										},
									},
									"isolated_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "CLB instance isolation time. Note: This field may return null, indicating that no valid values can be obtained.",
									},
									"expire_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "CLB instance expiration time, which takes effect only for prepaid instances. Note: This field may return null, indicating that no valid values can be obtained.",
									},
									"charge_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Billing mode of CLB instance. Valid values: PREPAID (monthly subscription), POSTPAID_BY_HOUR (pay as you go). Note: this field may return null, indicating that no valid values can be obtained.",
									},
									"network_attributes": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "CLB instance network attributes. Note: This field may return null, indicating that no valid values can be obtained.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"internet_charge_type": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "TRAFFIC_POSTPAID_BY_HOUR: hourly pay-as-you-go by traffic; BANDWIDTH_POSTPAID_BY_HOUR: hourly pay-as-you-go by bandwidth; BANDWIDTH_PACKAGE: billed by bandwidth package (currently, this method is supported only if the ISP is specified).",
												},
												"internet_max_bandwidth_out": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Maximum outbound bandwidth in Mbps, which applies only to public network CLB. Value range: 0-65,535. Default value: 10.",
												},
												"bandwidthpkg_sub_type": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Bandwidth package type, such as SINGLEISP. Note: This field may return null, indicating that no valid values can be obtained.",
												},
											},
										},
									},
									"prepaid_attributes": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Prepaid billing attributes of a CLB instance. Note: This field may return null, indicating that no valid values can be obtained.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"renew_flag": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Renewal type. AUTO_RENEW: automatic renewal; MANUAL_RENEW: manual renewal. Note: This field may return null, indicating that no valid values can be obtained.",
												},
												"period": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Cycle, indicating the number of months (reserved field). Note: This field may return null, indicating that no valid values can be obtained.",
												},
											},
										},
									},
									"log_set_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Logset ID of CLB Log Service (CLS). Note: This field may return null, indicating that no valid values can be obtained.",
									},
									"log_topic_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Log topic ID of CLB Log Service (CLS). Note: This field may return null, indicating that no valid values can be obtained.",
									},
									"address_i_pv6": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "IPv6 address of a CLB instance. Note: This field may return null, indicating that no valid values can be obtained.",
									},
									"extra_info": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Reserved field which can be ignored generally.Note: This field may return null, indicating that no valid values can be obtained.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"zhi_tong": {
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "Whether to enable VIP direct connection. Note: This field may return null, indicating that no valid values can be obtained.",
												},
												"tgw_group_name": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "TgwGroup name. Note: This field may return null, indicating that no valid values can be obtained.",
												},
											},
										},
									},
									"is_ddos": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "Whether an Anti-DDoS Pro instance can be bound. Note: This field may return null, indicating that no valid values can be obtained.",
									},
									"config_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Custom configuration ID at the CLB instance level. Note: This field may return null, indicating that no valid values can be obtained.",
									},
									"load_balancer_pass_to_target": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "Whether a real server opens the traffic from a CLB instance to the internet. Note: this field may return null, indicating that no valid values can be obtained.",
									},
									"exclusive_cluster": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Private network dedicated cluster. Note: this field may return null, indicating that no valid values can be obtained.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"l4_clusters": {
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Layer-4 dedicated cluster list. Note: this field may return null, indicating that no valid values can be obtained.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"cluster_id": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Unique cluster ID.",
															},
															"cluster_name": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Cluster name. Note: this field may return null, indicating that no valid values can be obtained.",
															},
															"zone": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Cluster AZ, such as ap-guangzhou-1. Note: this field may return null, indicating that no valid values can be obtained.",
															},
														},
													},
												},
												"l7_clusters": {
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Layer-7 dedicated cluster list. Note: this field may return null, indicating that no valid values can be obtained.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"cluster_id": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Unique cluster ID.",
															},
															"cluster_name": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Cluster name. Note: this field may return null, indicating that no valid values can be obtained.",
															},
															"zone": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Cluster AZ, such as ap-guangzhou-1. Note: this field may return null, indicating that no valid values can be obtained.",
															},
														},
													},
												},
												"classical_cluster": {
													Type:        schema.TypeList,
													Computed:    true,
													Description: "vpcgw cluster. Note: this field may return null, indicating that no valid values can be obtained.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"cluster_id": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Unique cluster ID.",
															},
															"cluster_name": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Cluster name. Note: this field may return null, indicating that no valid values can be obtained.",
															},
															"zone": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Cluster AZ, such as ap-guangzhou-1. Note: this field may return null, indicating that no valid values can be obtained.",
															},
														},
													},
												},
											},
										},
									},
									"ipv6_mode": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "This field is meaningful only when the IP address version is ipv6. Valid values: IPv6Nat64, IPv6FullChain. Note: this field may return null, indicating that no valid values can be obtained.",
									},
									"snat_pro": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "Whether to enable SnatPro. Note: this field may return null, indicating that no valid values can be obtained.",
									},
									"snat_ips": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "SnatIp list after SnatPro load balancing is enabled. Note: this field may return null, indicating that no valid values can be obtained.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"subnet_id": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Unique VPC subnet ID, such as subnet-12345678.",
												},
												"ip": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "IP address, such as 192.168.0.1.",
												},
											},
										},
									},
									"sla_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Specification of the LCU-supported instance. Note: This field may return null, indicating that no valid values can be obtained.",
									},
									"is_block": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "Whether VIP is blocked. Note: this field may return null, indicating that no valid values can be obtained.",
									},
									"is_block_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Time blocked or unblocked. Note: this field may return null, indicating that no valid values can be obtained.",
									},
									"local_bgp": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "Whether the IP type is the local BGP.",
									},
									"cluster_tag": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Dedicated layer-7 tag. Note: this field may return null, indicating that no valid values can be obtained.",
									},
									"mix_ip_target": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "If the layer-7 listener of an IPv6FullChain CLB instance is enabled, the CLB instance can be bound with an IPv4 and an IPv6 CVM instance simultaneously. Note: this field may return null, indicating that no valid values can be obtained.",
									},
									"zones": {
										Type: schema.TypeSet,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
										Computed:    true,
										Description: "Availability zone of a VPC-based private network CLB instance. Note: this field may return null, indicating that no valid values can be obtained.",
									},
									"nfv_info": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Whether it is an NFV CLB instance. No returned information: no; l7nfv: yes. Note: this field may return null, indicating that no valid values can be obtained.",
									},
									"health_log_set_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Health check logset ID of CLB CLS. Note: this field may return null, indicating that no valid values can be obtained.",
									},
									"health_log_topic_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Health check log topic ID of CLB CLS. Note: this field may return null, indicating that no valid values can be obtained.",
									},
									"cluster_ids": {
										Type: schema.TypeSet,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
										Computed:    true,
										Description: "Cluster ID. Note: This field may return null, indicating that no valid values can be obtained.",
									},
									"attribute_flags": {
										Type: schema.TypeSet,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
										Computed:    true,
										Description: "Cluster ID.Note: This field may return null, indicating that no valid values can be obtained.",
									},
									"load_balancer_domain": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Domain name of the CLB instance. Note: This field may return null, indicating that no valid values can be obtained.",
									},
								},
							},
						},
					},
				},
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
		},
	}
}

func dataSourceTencentCloudClbInstanceByCertIdRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_clb_instance_by_cert_id.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("cert_ids"); ok {
		certIdsSet := v.(*schema.Set).List()
		paramMap["CertIds"] = helper.InterfacesStringsPoint(certIdsSet)
	}

	service := ClbService{client: meta.(*TencentCloudClient).apiV3Conn}

	var certSet []*clb.CertIdRelatedWithLoadBalancers

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeClbInstanceByCertId(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		certSet = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(certSet))
	tmpList := make([]map[string]interface{}, 0, len(certSet))

	if certSet != nil {
		for _, certIdRelatedWithLoadBalancers := range certSet {
			certIdRelatedWithLoadBalancersMap := map[string]interface{}{}

			if certIdRelatedWithLoadBalancers.CertId != nil {
				certIdRelatedWithLoadBalancersMap["cert_id"] = certIdRelatedWithLoadBalancers.CertId
			}

			if certIdRelatedWithLoadBalancers.LoadBalancers != nil {
				loadBalancersList := []interface{}{}
				for _, loadBalancers := range certIdRelatedWithLoadBalancers.LoadBalancers {
					loadBalancersMap := map[string]interface{}{}

					if loadBalancers.LoadBalancerId != nil {
						loadBalancersMap["load_balancer_id"] = loadBalancers.LoadBalancerId
					}

					if loadBalancers.LoadBalancerName != nil {
						loadBalancersMap["load_balancer_name"] = loadBalancers.LoadBalancerName
					}

					if loadBalancers.LoadBalancerType != nil {
						loadBalancersMap["load_balancer_type"] = loadBalancers.LoadBalancerType
					}

					if loadBalancers.Forward != nil {
						loadBalancersMap["forward"] = loadBalancers.Forward
					}

					if loadBalancers.Domain != nil {
						loadBalancersMap["domain"] = loadBalancers.Domain
					}

					if loadBalancers.LoadBalancerVips != nil {
						loadBalancersMap["load_balancer_vips"] = loadBalancers.LoadBalancerVips
					}

					if loadBalancers.Status != nil {
						loadBalancersMap["status"] = loadBalancers.Status
					}

					if loadBalancers.CreateTime != nil {
						loadBalancersMap["create_time"] = loadBalancers.CreateTime
					}

					if loadBalancers.StatusTime != nil {
						loadBalancersMap["status_time"] = loadBalancers.StatusTime
					}

					if loadBalancers.ProjectId != nil {
						loadBalancersMap["project_id"] = loadBalancers.ProjectId
					}

					if loadBalancers.VpcId != nil {
						loadBalancersMap["vpc_id"] = loadBalancers.VpcId
					}

					if loadBalancers.OpenBgp != nil {
						loadBalancersMap["open_bgp"] = loadBalancers.OpenBgp
					}

					if loadBalancers.Snat != nil {
						loadBalancersMap["snat"] = loadBalancers.Snat
					}

					if loadBalancers.Isolation != nil {
						loadBalancersMap["isolation"] = loadBalancers.Isolation
					}

					if loadBalancers.Log != nil {
						loadBalancersMap["log"] = loadBalancers.Log
					}

					if loadBalancers.SubnetId != nil {
						loadBalancersMap["subnet_id"] = loadBalancers.SubnetId
					}

					if loadBalancers.Tags != nil {
						tagsList := []interface{}{}
						for _, tags := range loadBalancers.Tags {
							tagsMap := map[string]interface{}{}

							if tags.TagKey != nil {
								tagsMap["tag_key"] = tags.TagKey
							}

							if tags.TagValue != nil {
								tagsMap["tag_value"] = tags.TagValue
							}

							tagsList = append(tagsList, tagsMap)
						}

						loadBalancersMap["tags"] = tagsList
					}

					if loadBalancers.SecureGroups != nil {
						loadBalancersMap["secure_groups"] = loadBalancers.SecureGroups
					}

					if loadBalancers.TargetRegionInfo != nil {
						targetRegionInfoMap := map[string]interface{}{}

						if loadBalancers.TargetRegionInfo.Region != nil {
							targetRegionInfoMap["region"] = loadBalancers.TargetRegionInfo.Region
						}

						if loadBalancers.TargetRegionInfo.VpcId != nil {
							targetRegionInfoMap["vpc_id"] = loadBalancers.TargetRegionInfo.VpcId
						}

						loadBalancersMap["target_region_info"] = []interface{}{targetRegionInfoMap}
					}

					if loadBalancers.AnycastZone != nil {
						loadBalancersMap["anycast_zone"] = loadBalancers.AnycastZone
					}

					if loadBalancers.AddressIPVersion != nil {
						loadBalancersMap["address_i_p_version"] = loadBalancers.AddressIPVersion
					}

					if loadBalancers.NumericalVpcId != nil {
						loadBalancersMap["numerical_vpc_id"] = loadBalancers.NumericalVpcId
					}

					if loadBalancers.VipIsp != nil {
						loadBalancersMap["vip_isp"] = loadBalancers.VipIsp
					}

					if loadBalancers.MasterZone != nil {
						masterZoneMap := map[string]interface{}{}

						if loadBalancers.MasterZone.ZoneId != nil {
							masterZoneMap["zone_id"] = loadBalancers.MasterZone.ZoneId
						}

						if loadBalancers.MasterZone.Zone != nil {
							masterZoneMap["zone"] = loadBalancers.MasterZone.Zone
						}

						if loadBalancers.MasterZone.ZoneName != nil {
							masterZoneMap["zone_name"] = loadBalancers.MasterZone.ZoneName
						}

						if loadBalancers.MasterZone.ZoneRegion != nil {
							masterZoneMap["zone_region"] = loadBalancers.MasterZone.ZoneRegion
						}

						if loadBalancers.MasterZone.LocalZone != nil {
							masterZoneMap["local_zone"] = loadBalancers.MasterZone.LocalZone
						}

						if loadBalancers.MasterZone.EdgeZone != nil {
							masterZoneMap["edge_zone"] = loadBalancers.MasterZone.EdgeZone
						}

						loadBalancersMap["master_zone"] = []interface{}{masterZoneMap}
					}

					if loadBalancers.BackupZoneSet != nil {
						backupZoneSetList := []interface{}{}
						for _, backupZoneSet := range loadBalancers.BackupZoneSet {
							backupZoneSetMap := map[string]interface{}{}

							if backupZoneSet.ZoneId != nil {
								backupZoneSetMap["zone_id"] = backupZoneSet.ZoneId
							}

							if backupZoneSet.Zone != nil {
								backupZoneSetMap["zone"] = backupZoneSet.Zone
							}

							if backupZoneSet.ZoneName != nil {
								backupZoneSetMap["zone_name"] = backupZoneSet.ZoneName
							}

							if backupZoneSet.ZoneRegion != nil {
								backupZoneSetMap["zone_region"] = backupZoneSet.ZoneRegion
							}

							if backupZoneSet.LocalZone != nil {
								backupZoneSetMap["local_zone"] = backupZoneSet.LocalZone
							}

							if backupZoneSet.EdgeZone != nil {
								backupZoneSetMap["edge_zone"] = backupZoneSet.EdgeZone
							}

							backupZoneSetList = append(backupZoneSetList, backupZoneSetMap)
						}

						loadBalancersMap["backup_zone_set"] = backupZoneSetList
					}

					if loadBalancers.IsolatedTime != nil {
						loadBalancersMap["isolated_time"] = loadBalancers.IsolatedTime
					}

					if loadBalancers.ExpireTime != nil {
						loadBalancersMap["expire_time"] = loadBalancers.ExpireTime
					}

					if loadBalancers.ChargeType != nil {
						loadBalancersMap["charge_type"] = loadBalancers.ChargeType
					}

					if loadBalancers.NetworkAttributes != nil {
						networkAttributesMap := map[string]interface{}{}

						if loadBalancers.NetworkAttributes.InternetChargeType != nil {
							networkAttributesMap["internet_charge_type"] = loadBalancers.NetworkAttributes.InternetChargeType
						}

						if loadBalancers.NetworkAttributes.InternetMaxBandwidthOut != nil {
							networkAttributesMap["internet_max_bandwidth_out"] = loadBalancers.NetworkAttributes.InternetMaxBandwidthOut
						}

						if loadBalancers.NetworkAttributes.BandwidthpkgSubType != nil {
							networkAttributesMap["bandwidthpkg_sub_type"] = loadBalancers.NetworkAttributes.BandwidthpkgSubType
						}

						loadBalancersMap["network_attributes"] = []interface{}{networkAttributesMap}
					}

					if loadBalancers.PrepaidAttributes != nil {
						prepaidAttributesMap := map[string]interface{}{}

						if loadBalancers.PrepaidAttributes.RenewFlag != nil {
							prepaidAttributesMap["renew_flag"] = loadBalancers.PrepaidAttributes.RenewFlag
						}

						if loadBalancers.PrepaidAttributes.Period != nil {
							prepaidAttributesMap["period"] = loadBalancers.PrepaidAttributes.Period
						}

						loadBalancersMap["prepaid_attributes"] = []interface{}{prepaidAttributesMap}
					}

					if loadBalancers.LogSetId != nil {
						loadBalancersMap["log_set_id"] = loadBalancers.LogSetId
					}

					if loadBalancers.LogTopicId != nil {
						loadBalancersMap["log_topic_id"] = loadBalancers.LogTopicId
					}

					if loadBalancers.AddressIPv6 != nil {
						loadBalancersMap["address_i_pv6"] = loadBalancers.AddressIPv6
					}

					if loadBalancers.ExtraInfo != nil {
						extraInfoMap := map[string]interface{}{}

						if loadBalancers.ExtraInfo.ZhiTong != nil {
							extraInfoMap["zhi_tong"] = loadBalancers.ExtraInfo.ZhiTong
						}

						if loadBalancers.ExtraInfo.TgwGroupName != nil {
							extraInfoMap["tgw_group_name"] = loadBalancers.ExtraInfo.TgwGroupName
						}

						loadBalancersMap["extra_info"] = []interface{}{extraInfoMap}
					}

					if loadBalancers.IsDDos != nil {
						loadBalancersMap["is_ddos"] = loadBalancers.IsDDos
					}

					if loadBalancers.ConfigId != nil {
						loadBalancersMap["config_id"] = loadBalancers.ConfigId
					}

					if loadBalancers.LoadBalancerPassToTarget != nil {
						loadBalancersMap["load_balancer_pass_to_target"] = loadBalancers.LoadBalancerPassToTarget
					}

					if loadBalancers.ExclusiveCluster != nil {
						exclusiveClusterMap := map[string]interface{}{}

						if loadBalancers.ExclusiveCluster.L4Clusters != nil {
							l4ClustersList := []interface{}{}
							for _, l4Clusters := range loadBalancers.ExclusiveCluster.L4Clusters {
								l4ClustersMap := map[string]interface{}{}

								if l4Clusters.ClusterId != nil {
									l4ClustersMap["cluster_id"] = l4Clusters.ClusterId
								}

								if l4Clusters.ClusterName != nil {
									l4ClustersMap["cluster_name"] = l4Clusters.ClusterName
								}

								if l4Clusters.Zone != nil {
									l4ClustersMap["zone"] = l4Clusters.Zone
								}

								l4ClustersList = append(l4ClustersList, l4ClustersMap)
							}

							exclusiveClusterMap["l4_clusters"] = l4ClustersList
						}

						if loadBalancers.ExclusiveCluster.L7Clusters != nil {
							l7ClustersList := []interface{}{}
							for _, l7Clusters := range loadBalancers.ExclusiveCluster.L7Clusters {
								l7ClustersMap := map[string]interface{}{}

								if l7Clusters.ClusterId != nil {
									l7ClustersMap["cluster_id"] = l7Clusters.ClusterId
								}

								if l7Clusters.ClusterName != nil {
									l7ClustersMap["cluster_name"] = l7Clusters.ClusterName
								}

								if l7Clusters.Zone != nil {
									l7ClustersMap["zone"] = l7Clusters.Zone
								}

								l7ClustersList = append(l7ClustersList, l7ClustersMap)
							}

							exclusiveClusterMap["l7_clusters"] = l7ClustersList
						}

						if loadBalancers.ExclusiveCluster.ClassicalCluster != nil {
							classicalClusterMap := map[string]interface{}{}

							if loadBalancers.ExclusiveCluster.ClassicalCluster.ClusterId != nil {
								classicalClusterMap["cluster_id"] = loadBalancers.ExclusiveCluster.ClassicalCluster.ClusterId
							}

							if loadBalancers.ExclusiveCluster.ClassicalCluster.ClusterName != nil {
								classicalClusterMap["cluster_name"] = loadBalancers.ExclusiveCluster.ClassicalCluster.ClusterName
							}

							if loadBalancers.ExclusiveCluster.ClassicalCluster.Zone != nil {
								classicalClusterMap["zone"] = loadBalancers.ExclusiveCluster.ClassicalCluster.Zone
							}

							exclusiveClusterMap["classical_cluster"] = []interface{}{classicalClusterMap}
						}

						loadBalancersMap["exclusive_cluster"] = []interface{}{exclusiveClusterMap}
					}

					if loadBalancers.IPv6Mode != nil {
						loadBalancersMap["ipv6_mode"] = loadBalancers.IPv6Mode
					}

					if loadBalancers.SnatPro != nil {
						loadBalancersMap["snat_pro"] = loadBalancers.SnatPro
					}

					if loadBalancers.SnatIps != nil {
						snatIpsList := []interface{}{}
						for _, snatIps := range loadBalancers.SnatIps {
							snatIpsMap := map[string]interface{}{}

							if snatIps.SubnetId != nil {
								snatIpsMap["subnet_id"] = snatIps.SubnetId
							}

							if snatIps.Ip != nil {
								snatIpsMap["ip"] = snatIps.Ip
							}

							snatIpsList = append(snatIpsList, snatIpsMap)
						}

						loadBalancersMap["snat_ips"] = snatIpsList
					}

					if loadBalancers.SlaType != nil {
						loadBalancersMap["sla_type"] = loadBalancers.SlaType
					}

					if loadBalancers.IsBlock != nil {
						loadBalancersMap["is_block"] = loadBalancers.IsBlock
					}

					if loadBalancers.IsBlockTime != nil {
						loadBalancersMap["is_block_time"] = loadBalancers.IsBlockTime
					}

					if loadBalancers.LocalBgp != nil {
						loadBalancersMap["local_bgp"] = loadBalancers.LocalBgp
					}

					if loadBalancers.ClusterTag != nil {
						loadBalancersMap["cluster_tag"] = loadBalancers.ClusterTag
					}

					if loadBalancers.MixIpTarget != nil {
						loadBalancersMap["mix_ip_target"] = loadBalancers.MixIpTarget
					}

					if loadBalancers.Zones != nil {
						loadBalancersMap["zones"] = loadBalancers.Zones
					}

					if loadBalancers.NfvInfo != nil {
						loadBalancersMap["nfv_info"] = loadBalancers.NfvInfo
					}

					if loadBalancers.HealthLogSetId != nil {
						loadBalancersMap["health_log_set_id"] = loadBalancers.HealthLogSetId
					}

					if loadBalancers.HealthLogTopicId != nil {
						loadBalancersMap["health_log_topic_id"] = loadBalancers.HealthLogTopicId
					}

					if loadBalancers.ClusterIds != nil {
						loadBalancersMap["cluster_ids"] = loadBalancers.ClusterIds
					}

					if loadBalancers.AttributeFlags != nil {
						loadBalancersMap["attribute_flags"] = loadBalancers.AttributeFlags
					}

					if loadBalancers.LoadBalancerDomain != nil {
						loadBalancersMap["load_balancer_domain"] = loadBalancers.LoadBalancerDomain
					}

					loadBalancersList = append(loadBalancersList, loadBalancersMap)
				}

				certIdRelatedWithLoadBalancersMap["load_balancers"] = loadBalancersList
			}

			ids = append(ids, *certIdRelatedWithLoadBalancers.CertId)
			tmpList = append(tmpList, certIdRelatedWithLoadBalancersMap)
		}

		_ = d.Set("cert_set", tmpList)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), tmpList); e != nil {
			return e
		}
	}
	return nil
}
