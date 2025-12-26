package bh

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	bhv20230418 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/bh/v20230418"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudBhDevices() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudBhDevicesRead,
		Schema: map[string]*schema.Schema{
			"id_set": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Asset ID collection.",
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},

			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Asset name or asset IP, fuzzy search.",
			},

			"ip": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Not currently used.",
			},

			"ap_code_set": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Region code collection.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"kind": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Operating system type, 1 - Linux, 2 - Windows, 3 - MySQL, 4 - SQLServer.",
			},

			"authorized_user_id_set": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "User ID collection with access to this asset.",
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},

			"resource_id_set": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Filter condition, asset-bound bastion host service ID collection.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"kind_set": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Can filter by multiple types, 1 - Linux, 2 - Windows, 3 - MySQL, 4 - SQLServer.",
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},

			"managed_account": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Whether the asset contains managed accounts. 1, contains; 0, does not contain.",
			},

			"department_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Filter condition, can filter by department ID.",
			},

			"account_id_set": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Cloud account ID to which the asset belongs.",
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},

			"provider_type_set": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Cloud provider type, 1 - Tencent Cloud, 2 - Alibaba Cloud.",
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},

			"cloud_device_status_set": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Synchronized cloud asset status, marking the status of synchronized assets, 0 - deleted, 1 - normal, 2 - isolated, 3 - expired.",
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},

			"tag_filters": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Filter condition, can filter by tag key and tag value. If both tag key and tag value filter conditions are specified, they have an \"AND\" relationship.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"tag_key": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Tag key.",
						},
						"tag_value": {
							Type:        schema.TypeSet,
							Optional:    true,
							Description: "Tag value.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},

			"filters": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Filter array.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Field to filter. Support: BindingStatus, InstanceId, DeviceAccount, VpcId, DomainId, ResourceId, Name, Ip, ManageDimension.",
						},
						"values": {
							Type:        schema.TypeSet,
							Required:    true,
							Description: "Filter values for the field. \nIf multiple Filters exist, the relationship between Filters is logical AND. \nIf multiple Values exist for the same Filter, the relationship between Values under the same Filter is logical OR.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},

			"device_set": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Asset information list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Asset ID.",
						},
						"instance_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Instance ID, corresponding to CVM, CDB and other instance IDs.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Asset name.",
						},
						"public_ip": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Public IP.",
						},
						"private_ip": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Private IP.",
						},
						"ap_code": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Region code.",
						},
						"ap_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Region name.",
						},
						"os_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Operating system name.",
						},
						"kind": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Asset type 1 - Linux, 2 - Windows, 3 - MySQL, 4 - SQLServer.",
						},
						"port": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Management port.",
						},
						"group_set": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Asset group list to which it belongs.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Group ID.",
									},
									"name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Group name.",
									},
									"department": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Department information to which it belongs.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"id": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Department ID.",
												},
												"name": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Department name, 1 - 256 characters.",
												},
												"managers": {
													Type:        schema.TypeSet,
													Computed:    true,
													Description: "Department administrator account ID.",
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
												},
												"manager_users": {
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Administrator users.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"manager_id": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Administrator ID.",
															},
															"manager_name": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Administrator name.",
															},
														},
													},
												},
											},
										},
									},
									"count": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Count.",
									},
								},
							},
						},
						"account_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Number of accounts bound to the asset.",
						},
						"vpc_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "VPC ID.",
						},
						"subnet_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Subnet ID.",
						},
						"resource": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Bastion host service information, note that it is null when no service is bound.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"resource_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Service instance ID, such as bh-saas-s3ed4r5e.",
									},
									"ap_code": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Region code.",
									},
									"sv_args": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Service instance specification information.",
									},
									"vpc_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "VPC ID.",
									},
									"nodes": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Number of assets corresponding to the service specification.",
									},
									"renew_flag": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Auto-renewal flag, 0 - default state, 1 - auto-renewal, 2 - explicitly not auto-renewal.",
									},
									"expire_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Expiration time.",
									},
									"status": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Resource status, 0 - not initialized, 1 - normal, 2 - isolated, 3 - destroyed, 4 - initialization failed, 5 - initializing.",
									},
									"resource_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Service instance name, such as T-Sec-Bastion Host (SaaS type).",
									},
									"pid": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Pricing model ID.",
									},
									"create_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Resource creation time.",
									},
									"product_code": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Product code, p_cds_dasb.",
									},
									"sub_product_code": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Sub-product code, sp_cds_dasb_bh_saas.",
									},
									"zone": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Availability zone.",
									},
									"expired": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "Whether expired, true - expired, false - not expired.",
									},
									"deployed": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "Whether deployed, true - deployed, false - not deployed.",
									},
									"vpc_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "VPC name where the service is deployed.",
									},
									"vpc_cidr_block": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "CIDR block of the VPC where the service is deployed.",
									},
									"subnet_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Subnet ID where the service is deployed.",
									},
									"subnet_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Subnet name where the service is deployed.",
									},
									"cidr_block": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "CIDR block of the subnet where the service is deployed.",
									},
									"public_ip_set": {
										Type:        schema.TypeSet,
										Computed:    true,
										Description: "External IP.",
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									"private_ip_set": {
										Type:        schema.TypeSet,
										Computed:    true,
										Description: "Internal IP.",
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									"module_set": {
										Type:        schema.TypeSet,
										Computed:    true,
										Description: "Advanced feature list enabled for the service, such as: [DB].",
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									"used_nodes": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Number of used authorization points.",
									},
									"extend_points": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Extension points.",
									},
									"package_bandwidth": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Number of bandwidth extension packages (4M).",
									},
									"package_node": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Number of authorization point extension packages (50 points).",
									},
									"log_delivery_args": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Log delivery specification information.",
									},
									"clb_set": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Bastion host resource load balancer.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"clb_ip": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Load balancer IP.",
												},
											},
										},
									},
									"domain_count": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Number of network domains.",
									},
									"used_domain_count": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Number of network domains already used.",
									},
									"trial": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "0 non-trial version, 1 trial version.",
									},
									"log_delivery": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Log delivery specification information.",
									},
									"cdc_cluster_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "CDC cluster ID.",
									},
									"deploy_model": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Deployment mode, default 0, 0-cvm 1-tke.",
									},
									"intranet_access": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "0 default value, non-intranet access, 1 intranet access, 2 intranet access opening, 3 intranet access closing.",
									},
									"intranet_private_ip_set": {
										Type:        schema.TypeSet,
										Computed:    true,
										Description: "IP addresses for intranet access.",
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									"intranet_vpc_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "VPC for enabling intranet access.",
									},
									"intranet_subnet_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Subnet ID for enabling intranet access.",
									},
									"intranet_vpc_cidr": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "CIDR block of the VPC for enabling intranet access.",
									},
									"domain_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Custom domain name for bastion host intranet IP.",
									},
									"share_clb": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "Whether to share CLB, true-shared CLB, false-dedicated CLB.",
									},
									"open_clb_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Shared CLB ID.",
									},
									"lb_vip_isp": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "ISP information.",
									},
									"tui_cmd_port": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Linux asset command line operation port.",
									},
									"tui_direct_port": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Linux asset direct connection port.",
									},
									"web_access": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "1 default value, web access enabled, 0 web access disabled, 2 web access opening, 3 web access closing.",
									},
									"client_access": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "1 default value, client access enabled, 0 client access disabled, 2 client access opening, 3 client access closing.",
									},
									"external_access": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "1 default value, external access enabled, 0 external access disabled, 2 external access opening, 3 external access closing.",
									},
									"ioa_resource": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "0 default value, 0-free version (trial version) IOA, 1-paid version IOA.",
									},
									"package_ioa_user_count": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Number of zero trust bastion host user extension packages, 1 extension package corresponds to 20 users.",
									},
									"package_ioa_bandwidth": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Number of zero trust bastion host bandwidth extension packages, one extension package represents 4M bandwidth.",
									},
									"ioa_resource_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Zero trust instance ID corresponding to the bastion host instance.",
									},
								},
							},
						},
						"department": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Department to which the asset belongs.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Department ID.",
									},
									"name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Department name, 1 - 256 characters.",
									},
									"managers": {
										Type:        schema.TypeSet,
										Computed:    true,
										Description: "Department administrator account ID.",
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									"manager_users": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Administrator users.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"manager_id": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Administrator ID.",
												},
												"manager_name": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Administrator name.",
												},
											},
										},
									},
								},
							},
						},
						"ip_port_set": {
							Type:        schema.TypeSet,
							Computed:    true,
							Description: "Multi-node information for database assets.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"domain_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Network domain ID.",
						},
						"domain_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Network domain name.",
						},
						"enable_ssl": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Whether SSL is enabled, only supports Redis assets, 0: disabled 1: enabled.",
						},
						"ssl_cert_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of the uploaded SSL certificate.",
						},
						"ioa_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Resource ID on the IOA side.",
						},
						"manage_dimension": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "K8S cluster management dimension, 1-cluster, 2-namespace, 3-workload.",
						},
						"manage_account_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "K8S cluster management account ID.",
						},
						"namespace": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "K8S cluster namespace.",
						},
						"workload": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "K8S cluster workload.",
						},
						"sync_pod_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Number of synchronized pods in K8S cluster.",
						},
						"total_pod_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Total number of pods in K8S cluster.",
						},
						"cloud_account_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Cloud account ID.",
						},
						"cloud_account_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Cloud account name.",
						},
						"provider_type": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Cloud provider type, 1-Tencent Cloud, 2-Alibaba Cloud.",
						},
						"provider_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Cloud provider name.",
						},
						"sync_cloud_device_status": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Synchronized cloud asset status, marking the status of synchronized assets, 0-deleted, 1-normal, 2-isolated, 3-expired.",
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

func dataSourceTencentCloudBhDevicesRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_bh_devices.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(nil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = BhService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("id_set"); ok {
		idSetList := []*uint64{}
		idSetSet := v.(*schema.Set).List()
		for i := range idSetSet {
			idSet := idSetSet[i].(int)
			idSetList = append(idSetList, helper.IntUint64(idSet))
		}

		paramMap["IdSet"] = idSetList
	}

	if v, ok := d.GetOk("name"); ok {
		paramMap["Name"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("ip"); ok {
		paramMap["Ip"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("ap_code_set"); ok {
		apCodeSetList := []*string{}
		apCodeSetSet := v.(*schema.Set).List()
		for i := range apCodeSetSet {
			apCodeSet := apCodeSetSet[i].(string)
			apCodeSetList = append(apCodeSetList, helper.String(apCodeSet))
		}

		paramMap["ApCodeSet"] = apCodeSetList
	}

	if v, ok := d.GetOkExists("kind"); ok {
		paramMap["Kind"] = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("authorized_user_id_set"); ok {
		authorizedUserIdSetList := []*uint64{}
		authorizedUserIdSetSet := v.(*schema.Set).List()
		for i := range authorizedUserIdSetSet {
			authorizedUserIdSet := authorizedUserIdSetSet[i].(int)
			authorizedUserIdSetList = append(authorizedUserIdSetList, helper.IntUint64(authorizedUserIdSet))
		}

		paramMap["AuthorizedUserIdSet"] = authorizedUserIdSetList
	}

	if v, ok := d.GetOk("resource_id_set"); ok {
		resourceIdSetList := []*string{}
		resourceIdSetSet := v.(*schema.Set).List()
		for i := range resourceIdSetSet {
			resourceIdSet := resourceIdSetSet[i].(string)
			resourceIdSetList = append(resourceIdSetList, helper.String(resourceIdSet))
		}

		paramMap["ResourceIdSet"] = resourceIdSetList
	}

	if v, ok := d.GetOk("kind_set"); ok {
		kindSetList := []*uint64{}
		kindSetSet := v.(*schema.Set).List()
		for i := range kindSetSet {
			kindSet := kindSetSet[i].(int)
			kindSetList = append(kindSetList, helper.IntUint64(kindSet))
		}

		paramMap["KindSet"] = kindSetList
	}

	if v, ok := d.GetOk("managed_account"); ok {
		paramMap["ManagedAccount"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("department_id"); ok {
		paramMap["DepartmentId"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("account_id_set"); ok {
		accountIdSetList := []*uint64{}
		accountIdSetSet := v.(*schema.Set).List()
		for i := range accountIdSetSet {
			accountIdSet := accountIdSetSet[i].(int)
			accountIdSetList = append(accountIdSetList, helper.IntUint64(accountIdSet))
		}

		paramMap["AccountIdSet"] = accountIdSetList
	}

	if v, ok := d.GetOk("provider_type_set"); ok {
		providerTypeSetList := []*uint64{}
		providerTypeSetSet := v.(*schema.Set).List()
		for i := range providerTypeSetSet {
			providerTypeSet := providerTypeSetSet[i].(int)
			providerTypeSetList = append(providerTypeSetList, helper.IntUint64(providerTypeSet))
		}

		paramMap["ProviderTypeSet"] = providerTypeSetList
	}

	if v, ok := d.GetOk("cloud_device_status_set"); ok {
		cloudDeviceStatusSetList := []*uint64{}
		cloudDeviceStatusSetSet := v.(*schema.Set).List()
		for i := range cloudDeviceStatusSetSet {
			cloudDeviceStatusSet := cloudDeviceStatusSetSet[i].(int)
			cloudDeviceStatusSetList = append(cloudDeviceStatusSetList, helper.IntUint64(cloudDeviceStatusSet))
		}

		paramMap["CloudDeviceStatusSet"] = cloudDeviceStatusSetList
	}

	if v, ok := d.GetOk("tag_filters"); ok {
		tagFiltersSet := v.([]interface{})
		tmpSet := make([]*bhv20230418.TagFilter, 0, len(tagFiltersSet))
		for _, item := range tagFiltersSet {
			tagFiltersMap := item.(map[string]interface{})
			tagFilter := bhv20230418.TagFilter{}
			if v, ok := tagFiltersMap["tag_key"].(string); ok && v != "" {
				tagFilter.TagKey = helper.String(v)
			}

			if v, ok := tagFiltersMap["tag_value"]; ok {
				tagValueSet := v.(*schema.Set).List()
				for i := range tagValueSet {
					tagValue := tagValueSet[i].(string)
					tagFilter.TagValue = append(tagFilter.TagValue, helper.String(tagValue))
				}
			}

			tmpSet = append(tmpSet, &tagFilter)
		}

		paramMap["TagFilters"] = tmpSet
	}

	if v, ok := d.GetOk("filters"); ok {
		filtersSet := v.([]interface{})
		tmpSet := make([]*bhv20230418.Filter, 0, len(filtersSet))
		for _, item := range filtersSet {
			filtersMap := item.(map[string]interface{})
			filter := bhv20230418.Filter{}
			if v, ok := filtersMap["name"].(string); ok && v != "" {
				filter.Name = helper.String(v)
			}

			if v, ok := filtersMap["values"]; ok {
				valuesSet := v.(*schema.Set).List()
				for i := range valuesSet {
					values := valuesSet[i].(string)
					filter.Values = append(filter.Values, helper.String(values))
				}
			}

			tmpSet = append(tmpSet, &filter)
		}

		paramMap["Filters"] = tmpSet
	}

	var respData []*bhv20230418.Device
	reqErr := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeBhDevicesByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}

		respData = result
		return nil
	})

	if reqErr != nil {
		return reqErr
	}

	deviceSetList := make([]map[string]interface{}, 0, len(respData))
	if respData != nil {
		for _, deviceSet := range respData {
			deviceSetMap := map[string]interface{}{}
			if deviceSet.Id != nil {
				deviceSetMap["id"] = deviceSet.Id
			}

			if deviceSet.InstanceId != nil {
				deviceSetMap["instance_id"] = deviceSet.InstanceId
			}

			if deviceSet.Name != nil {
				deviceSetMap["name"] = deviceSet.Name
			}

			if deviceSet.PublicIp != nil {
				deviceSetMap["public_ip"] = deviceSet.PublicIp
			}

			if deviceSet.PrivateIp != nil {
				deviceSetMap["private_ip"] = deviceSet.PrivateIp
			}

			if deviceSet.ApCode != nil {
				deviceSetMap["ap_code"] = deviceSet.ApCode
			}

			if deviceSet.ApName != nil {
				deviceSetMap["ap_name"] = deviceSet.ApName
			}

			if deviceSet.OsName != nil {
				deviceSetMap["os_name"] = deviceSet.OsName
			}

			if deviceSet.Kind != nil {
				deviceSetMap["kind"] = deviceSet.Kind
			}

			if deviceSet.Port != nil {
				deviceSetMap["port"] = deviceSet.Port
			}

			groupSetList := make([]map[string]interface{}, 0, len(deviceSet.GroupSet))
			if deviceSet.GroupSet != nil {
				for _, groupSet := range deviceSet.GroupSet {
					groupSetMap := map[string]interface{}{}
					if groupSet.Id != nil {
						groupSetMap["id"] = groupSet.Id
					}

					if groupSet.Name != nil {
						groupSetMap["name"] = groupSet.Name
					}

					departmentMap := map[string]interface{}{}
					if groupSet.Department != nil {
						if groupSet.Department.Id != nil {
							departmentMap["id"] = groupSet.Department.Id
						}

						if groupSet.Department.Name != nil {
							departmentMap["name"] = groupSet.Department.Name
						}

						if groupSet.Department.Managers != nil {
							departmentMap["managers"] = groupSet.Department.Managers
						}

						managerUsersList := make([]map[string]interface{}, 0, len(groupSet.Department.ManagerUsers))
						if groupSet.Department.ManagerUsers != nil {
							for _, managerUsers := range groupSet.Department.ManagerUsers {
								managerUsersMap := map[string]interface{}{}
								if managerUsers.ManagerId != nil {
									managerUsersMap["manager_id"] = managerUsers.ManagerId
								}

								if managerUsers.ManagerName != nil {
									managerUsersMap["manager_name"] = managerUsers.ManagerName
								}

								managerUsersList = append(managerUsersList, managerUsersMap)
							}

							departmentMap["manager_users"] = managerUsersList
						}

						groupSetMap["department"] = []interface{}{departmentMap}
					}

					if groupSet.Count != nil {
						groupSetMap["count"] = groupSet.Count
					}

					groupSetList = append(groupSetList, groupSetMap)
				}

				deviceSetMap["group_set"] = groupSetList
			}

			if deviceSet.AccountCount != nil {
				deviceSetMap["account_count"] = deviceSet.AccountCount
			}

			if deviceSet.VpcId != nil {
				deviceSetMap["vpc_id"] = deviceSet.VpcId
			}

			if deviceSet.SubnetId != nil {
				deviceSetMap["subnet_id"] = deviceSet.SubnetId
			}

			resourceMap := map[string]interface{}{}
			if deviceSet.Resource != nil {
				if deviceSet.Resource.ResourceId != nil {
					resourceMap["resource_id"] = deviceSet.Resource.ResourceId
				}

				if deviceSet.Resource.ApCode != nil {
					resourceMap["ap_code"] = deviceSet.Resource.ApCode
				}

				if deviceSet.Resource.SvArgs != nil {
					resourceMap["sv_args"] = deviceSet.Resource.SvArgs
				}

				if deviceSet.Resource.VpcId != nil {
					resourceMap["vpc_id"] = deviceSet.Resource.VpcId
				}

				if deviceSet.Resource.Nodes != nil {
					resourceMap["nodes"] = deviceSet.Resource.Nodes
				}

				if deviceSet.Resource.RenewFlag != nil {
					resourceMap["renew_flag"] = deviceSet.Resource.RenewFlag
				}

				if deviceSet.Resource.ExpireTime != nil {
					resourceMap["expire_time"] = deviceSet.Resource.ExpireTime
				}

				if deviceSet.Resource.Status != nil {
					resourceMap["status"] = deviceSet.Resource.Status
				}

				if deviceSet.Resource.ResourceName != nil {
					resourceMap["resource_name"] = deviceSet.Resource.ResourceName
				}

				if deviceSet.Resource.Pid != nil {
					resourceMap["pid"] = deviceSet.Resource.Pid
				}

				if deviceSet.Resource.CreateTime != nil {
					resourceMap["create_time"] = deviceSet.Resource.CreateTime
				}

				if deviceSet.Resource.ProductCode != nil {
					resourceMap["product_code"] = deviceSet.Resource.ProductCode
				}

				if deviceSet.Resource.SubProductCode != nil {
					resourceMap["sub_product_code"] = deviceSet.Resource.SubProductCode
				}

				if deviceSet.Resource.Zone != nil {
					resourceMap["zone"] = deviceSet.Resource.Zone
				}

				if deviceSet.Resource.Expired != nil {
					resourceMap["expired"] = deviceSet.Resource.Expired
				}

				if deviceSet.Resource.Deployed != nil {
					resourceMap["deployed"] = deviceSet.Resource.Deployed
				}

				if deviceSet.Resource.VpcName != nil {
					resourceMap["vpc_name"] = deviceSet.Resource.VpcName
				}

				if deviceSet.Resource.VpcCidrBlock != nil {
					resourceMap["vpc_cidr_block"] = deviceSet.Resource.VpcCidrBlock
				}

				if deviceSet.Resource.SubnetId != nil {
					resourceMap["subnet_id"] = deviceSet.Resource.SubnetId
				}

				if deviceSet.Resource.SubnetName != nil {
					resourceMap["subnet_name"] = deviceSet.Resource.SubnetName
				}

				if deviceSet.Resource.CidrBlock != nil {
					resourceMap["cidr_block"] = deviceSet.Resource.CidrBlock
				}

				if deviceSet.Resource.PublicIpSet != nil {
					resourceMap["public_ip_set"] = deviceSet.Resource.PublicIpSet
				}

				if deviceSet.Resource.PrivateIpSet != nil {
					resourceMap["private_ip_set"] = deviceSet.Resource.PrivateIpSet
				}

				if deviceSet.Resource.ModuleSet != nil {
					resourceMap["module_set"] = deviceSet.Resource.ModuleSet
				}

				if deviceSet.Resource.UsedNodes != nil {
					resourceMap["used_nodes"] = deviceSet.Resource.UsedNodes
				}

				if deviceSet.Resource.ExtendPoints != nil {
					resourceMap["extend_points"] = deviceSet.Resource.ExtendPoints
				}

				if deviceSet.Resource.PackageBandwidth != nil {
					resourceMap["package_bandwidth"] = deviceSet.Resource.PackageBandwidth
				}

				if deviceSet.Resource.PackageNode != nil {
					resourceMap["package_node"] = deviceSet.Resource.PackageNode
				}

				if deviceSet.Resource.LogDeliveryArgs != nil {
					resourceMap["log_delivery_args"] = deviceSet.Resource.LogDeliveryArgs
				}

				clbSetList := make([]map[string]interface{}, 0, len(deviceSet.Resource.ClbSet))
				if deviceSet.Resource.ClbSet != nil {
					for _, clbSet := range deviceSet.Resource.ClbSet {
						clbSetMap := map[string]interface{}{}
						if clbSet.ClbIp != nil {
							clbSetMap["clb_ip"] = clbSet.ClbIp
						}

						clbSetList = append(clbSetList, clbSetMap)
					}

					resourceMap["clb_set"] = clbSetList
				}

				if deviceSet.Resource.DomainCount != nil {
					resourceMap["domain_count"] = deviceSet.Resource.DomainCount
				}

				if deviceSet.Resource.UsedDomainCount != nil {
					resourceMap["used_domain_count"] = deviceSet.Resource.UsedDomainCount
				}

				if deviceSet.Resource.Trial != nil {
					resourceMap["trial"] = deviceSet.Resource.Trial
				}

				if deviceSet.Resource.LogDelivery != nil {
					resourceMap["log_delivery"] = deviceSet.Resource.LogDelivery
				}

				if deviceSet.Resource.CdcClusterId != nil {
					resourceMap["cdc_cluster_id"] = deviceSet.Resource.CdcClusterId
				}

				if deviceSet.Resource.DeployModel != nil {
					resourceMap["deploy_model"] = deviceSet.Resource.DeployModel
				}

				if deviceSet.Resource.IntranetAccess != nil {
					resourceMap["intranet_access"] = deviceSet.Resource.IntranetAccess
				}

				if deviceSet.Resource.IntranetPrivateIpSet != nil {
					resourceMap["intranet_private_ip_set"] = deviceSet.Resource.IntranetPrivateIpSet
				}

				if deviceSet.Resource.IntranetVpcId != nil {
					resourceMap["intranet_vpc_id"] = deviceSet.Resource.IntranetVpcId
				}

				if deviceSet.Resource.IntranetSubnetId != nil {
					resourceMap["intranet_subnet_id"] = deviceSet.Resource.IntranetSubnetId
				}

				if deviceSet.Resource.IntranetVpcCidr != nil {
					resourceMap["intranet_vpc_cidr"] = deviceSet.Resource.IntranetVpcCidr
				}

				if deviceSet.Resource.DomainName != nil {
					resourceMap["domain_name"] = deviceSet.Resource.DomainName
				}

				if deviceSet.Resource.ShareClb != nil {
					resourceMap["share_clb"] = deviceSet.Resource.ShareClb
				}

				if deviceSet.Resource.OpenClbId != nil {
					resourceMap["open_clb_id"] = deviceSet.Resource.OpenClbId
				}

				if deviceSet.Resource.LbVipIsp != nil {
					resourceMap["lb_vip_isp"] = deviceSet.Resource.LbVipIsp
				}

				if deviceSet.Resource.TUICmdPort != nil {
					resourceMap["tui_cmd_port"] = deviceSet.Resource.TUICmdPort
				}

				if deviceSet.Resource.TUIDirectPort != nil {
					resourceMap["tui_direct_port"] = deviceSet.Resource.TUIDirectPort
				}

				if deviceSet.Resource.WebAccess != nil {
					resourceMap["web_access"] = deviceSet.Resource.WebAccess
				}

				if deviceSet.Resource.ClientAccess != nil {
					resourceMap["client_access"] = deviceSet.Resource.ClientAccess
				}

				if deviceSet.Resource.ExternalAccess != nil {
					resourceMap["external_access"] = deviceSet.Resource.ExternalAccess
				}

				if deviceSet.Resource.IOAResource != nil {
					resourceMap["ioa_resource"] = deviceSet.Resource.IOAResource
				}

				if deviceSet.Resource.PackageIOAUserCount != nil {
					resourceMap["package_ioa_user_count"] = deviceSet.Resource.PackageIOAUserCount
				}

				if deviceSet.Resource.PackageIOABandwidth != nil {
					resourceMap["package_ioa_bandwidth"] = deviceSet.Resource.PackageIOABandwidth
				}

				if deviceSet.Resource.IOAResourceId != nil {
					resourceMap["ioa_resource_id"] = deviceSet.Resource.IOAResourceId
				}

				deviceSetMap["resource"] = []interface{}{resourceMap}
			}

			departmentMap := map[string]interface{}{}
			if deviceSet.Department != nil {
				if deviceSet.Department.Id != nil {
					departmentMap["id"] = deviceSet.Department.Id
				}

				if deviceSet.Department.Name != nil {
					departmentMap["name"] = deviceSet.Department.Name
				}

				if deviceSet.Department.Managers != nil {
					departmentMap["managers"] = deviceSet.Department.Managers
				}

				managerUsersList := make([]map[string]interface{}, 0, len(deviceSet.Department.ManagerUsers))
				if deviceSet.Department.ManagerUsers != nil {
					for _, managerUsers := range deviceSet.Department.ManagerUsers {
						managerUsersMap := map[string]interface{}{}
						if managerUsers.ManagerId != nil {
							managerUsersMap["manager_id"] = managerUsers.ManagerId
						}

						if managerUsers.ManagerName != nil {
							managerUsersMap["manager_name"] = managerUsers.ManagerName
						}

						managerUsersList = append(managerUsersList, managerUsersMap)
					}

					departmentMap["manager_users"] = managerUsersList
				}

				deviceSetMap["department"] = []interface{}{departmentMap}
			}

			if deviceSet.IpPortSet != nil {
				deviceSetMap["ip_port_set"] = deviceSet.IpPortSet
			}

			if deviceSet.DomainId != nil {
				deviceSetMap["domain_id"] = deviceSet.DomainId
			}

			if deviceSet.DomainName != nil {
				deviceSetMap["domain_name"] = deviceSet.DomainName
			}

			if deviceSet.EnableSSL != nil {
				deviceSetMap["enable_ssl"] = deviceSet.EnableSSL
			}

			if deviceSet.SSLCertName != nil {
				deviceSetMap["ssl_cert_name"] = deviceSet.SSLCertName
			}

			if deviceSet.IOAId != nil {
				deviceSetMap["ioa_id"] = deviceSet.IOAId
			}

			if deviceSet.ManageDimension != nil {
				deviceSetMap["manage_dimension"] = deviceSet.ManageDimension
			}

			if deviceSet.ManageAccountId != nil {
				deviceSetMap["manage_account_id"] = deviceSet.ManageAccountId
			}

			if deviceSet.Namespace != nil {
				deviceSetMap["namespace"] = deviceSet.Namespace
			}

			if deviceSet.Workload != nil {
				deviceSetMap["workload"] = deviceSet.Workload
			}

			if deviceSet.SyncPodCount != nil {
				deviceSetMap["sync_pod_count"] = deviceSet.SyncPodCount
			}

			if deviceSet.TotalPodCount != nil {
				deviceSetMap["total_pod_count"] = deviceSet.TotalPodCount
			}

			if deviceSet.CloudAccountId != nil {
				deviceSetMap["cloud_account_id"] = deviceSet.CloudAccountId
			}

			if deviceSet.CloudAccountName != nil {
				deviceSetMap["cloud_account_name"] = deviceSet.CloudAccountName
			}

			if deviceSet.ProviderType != nil {
				deviceSetMap["provider_type"] = deviceSet.ProviderType
			}

			if deviceSet.ProviderName != nil {
				deviceSetMap["provider_name"] = deviceSet.ProviderName
			}

			if deviceSet.SyncCloudDeviceStatus != nil {
				deviceSetMap["sync_cloud_device_status"] = deviceSet.SyncCloudDeviceStatus
			}

			deviceSetList = append(deviceSetList, deviceSetMap)
		}

		_ = d.Set("device_set", deviceSetList)
	}

	d.SetId(helper.BuildToken())
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), deviceSetList); e != nil {
			return e
		}
	}

	return nil
}
