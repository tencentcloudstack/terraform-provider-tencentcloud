package emr

import (
	"context"
	"fmt"
	"log"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	emr "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/emr/v20190103"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudEmrClusterV2() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudEmrClusterV2Create,
		Read:   resourceTencentCloudEmrClusterV2Read,
		Update: resourceTencentCloudEmrClusterV2Update,
		Delete: resourceTencentCloudEmrClusterV2Delete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"product_version": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "EMR product version name, e.g., `EMR-V3.5.0`.",
			},
			"enable_support_ha_flag": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Whether to enable node high availability. `true` means enabled, `false` means disabled.",
			},
			"instance_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Instance name. Length 6-36 characters. Only Chinese, letters, numbers, `-`, `_` are allowed.",
			},
			"instance_charge_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Instance billing mode. Valid values: `PREPAID` (monthly/yearly subscription), `POSTPAID_BY_HOUR` (pay-as-you-go).",
			},
			"login_settings": {
				Type:        schema.TypeList,
				Required:    true,
				MaxItems:    1,
				Description: "Login settings for purchased nodes.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"password": {
							Type:        schema.TypeString,
							Optional:    true,
							Sensitive:   true,
							Description: "Login password, 8-16 characters, must contain uppercase, lowercase, digit and special character (supported special chars: `!@%^*`). First char cannot be special.",
						},
						"public_key_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Public key ID for key-based login.",
						},
					},
				},
			},
			"scene_software_config": {
				Type:        schema.TypeList,
				Required:    true,
				MaxItems:    1,
				Description: "Cluster scenario and components to deploy.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"software": {
							Type:        schema.TypeList,
							Required:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "List of components with versions, e.g., `[\"hdfs-3.2.2\", \"yarn-3.2.2\"]`.",
						},
						"scene_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Scenario name, e.g., `Hadoop-Default`, `Hadoop-Kudu`, `Hadoop-Zookeeper`, `Hadoop-Presto`, `Hadoop-Hbase`.",
						},
					},
				},
			},
			"instance_charge_prepaid": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "Prepaid (monthly/yearly) billing parameters. Required when `instance_charge_type` is `PREPAID`.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"period": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Purchase duration in months. Valid values: 1-12, 24, 36, 48, 60.",
						},
						"renew_flag": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Auto-renew flag. Default is false.",
						},
					},
				},
			},
			"security_group_ids": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Security group IDs bound to the instance, e.g., `[\"sg-xxxxxxxx\"]`.",
			},
			"script_bootstrap_action_config": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Bootstrap script configurations.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cos_file_uri": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "COS URI of the script.",
						},
						"execution_moment": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Execution timing. Valid values: `resourceAfter`, `clusterAfter`, `clusterBefore`.",
						},
						"args": {
							Type:        schema.TypeList,
							Optional:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "Script arguments, following standard Shell convention.",
						},
						"cos_file_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Script file name.",
						},
						"remark": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Remark.",
						},
					},
				},
			},
			"need_master_wan": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Whether to enable master public network. Valid values: `NEED_MASTER_WAN` (default), `NOT_NEED_MASTER_WAN`.",
			},
			"enable_remote_login_flag": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether to enable external remote login. Invalid when `security_group_ids` is set. Default is false.",
			},
			"enable_kerberos_flag": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether to enable Kerberos authentication. Default is false.",
			},
			"custom_conf": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Custom software configuration in JSON format.",
			},
			"tags": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Tags to bind to the cluster instance.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"tag_key": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Tag key.",
						},
						"tag_value": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Tag value.",
						},
					},
				},
			},
			"disaster_recover_group_ids": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Spread placement group IDs. Currently supports only one ID.",
			},
			"enable_cbs_encrypt_flag": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether to enable cluster-level CBS encryption. Default is false.",
			},
			"meta_db_info": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "Metadata database information. When `meta_type` is `EMR_NEW_META`/`EMR_DEFAULT_META`, no extra fields are required; when `EMR_EXIT_META`, `unify_meta_instance_id` must be set; when `USER_CUSTOM_META`, `meta_data_jdbc_url`/`meta_data_user`/`meta_data_pass` must be set.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"meta_data_jdbc_url": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Custom MetaDB JDBC URL, e.g., `jdbc:mysql://10.10.10.10:3306/dbname`.",
						},
						"meta_data_user": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Custom MetaDB username.",
						},
						"meta_data_pass": {
							Type:        schema.TypeString,
							Optional:    true,
							Sensitive:   true,
							Description: "Custom MetaDB password.",
						},
						"meta_type": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Hive shared metadata DB type. Valid values: `EMR_DEFAULT_META`, `EMR_EXIT_META`, `USER_CUSTOM_META`.",
						},
						"unify_meta_instance_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "EMR-MetaDB instance ID.",
						},
					},
				},
			},
			"depend_service": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Shared component dependency information.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"service_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Shared component name.",
						},
						"instance_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Shared component cluster instance ID.",
						},
					},
				},
			},
			"zone_resource_configuration": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    3,
				Description: "Per-zone resource configuration. Supports 1 (single-AZ) or up to 3 entries (multi-AZ: primary, backup, arbitration). ZoneTag is derived automatically from the list index.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"virtual_private_cloud": {
							Type:        schema.TypeList,
							Optional:    true,
							MaxItems:    1,
							Description: "VPC/Subnet information.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"vpc_id": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "VPC ID.",
									},
									"subnet_id": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Subnet ID.",
									},
								},
							},
						},
						"placement": {
							Type:        schema.TypeList,
							Optional:    true,
							MaxItems:    1,
							Description: "Zone and project placement.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"zone": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Availability zone, e.g., `ap-guangzhou-7`.",
									},
									"project_id": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "Project ID. Defaults to default project if omitted.",
									},
								},
							},
						},
						"all_node_resource_spec": {
							Type:        schema.TypeList,
							Optional:    true,
							MaxItems:    1,
							Description: "Resource specifications for all node roles.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"master_resource_spec": {
										Type:     schema.TypeList,
										Optional: true,
										Description: "Master node resource specifications. The number of blocks determines `MasterCount` passed to `CreateCluster`. " +
											"**All blocks must have identical configuration**; the first block is used as the single resource template by the API. " +
											"A validation error is returned at create time if any block differs from the first.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"instance_type": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "CVM instance type, e.g., `S6.2XLARGE32`, `SA4.8XLARGE64`.",
												},
												"system_disk": {
													Type:        schema.TypeList,
													Optional:    true,
													MaxItems:    1,
													Description: "System disk specifications.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"disk_size": {
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "Disk size in GB.",
															},
															"disk_type": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Disk type. Valid values: `CLOUD_SSD`, `CLOUD_PREMIUM`, `CLOUD_BASIC`, `LOCAL_BASIC`, `LOCAL_SSD`.",
															},
														},
													},
												},
												"data_disk": {
													Type:        schema.TypeSet,
													Optional:    true,
													Description: "Cloud data disk specifications.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"count": {
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "Number of disks.",
															},
															"disk_size": {
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "Disk size in GB.",
															},
															"disk_type": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Disk type. Valid values: `CLOUD_SSD`, `CLOUD_PREMIUM`, `CLOUD_BASIC`, `LOCAL_BASIC`, `LOCAL_SSD`, `CLOUD_HSSD`, `CLOUD_THROUGHPUT`, `CLOUD_TSSD`, `CLOUD_BIGDATA`, `CLOUD_HIGHIO`, `CLOUD_BSSD`, `REMOTE_SSD`.",
															},
														},
													},
												},
												// "tags": {
												// 	Type:        schema.TypeList,
												// 	Optional:    true,
												// 	Description: "Tags to bind to the node.",
												// 	Elem: &schema.Resource{
												// 		Schema: map[string]*schema.Schema{
												// 			"tag_key": {
												// 				Type:        schema.TypeString,
												// 				Optional:    true,
												// 				Description: "Tag key.",
												// 			},
												// 			"tag_value": {
												// 				Type:        schema.TypeString,
												// 				Optional:    true,
												// 				Description: "Tag value.",
												// 			},
												// 		},
												// 	},
												// },
												"emr_resource_id": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "EMR node resource ID.",
												},
											},
										},
									},
									"core_resource_spec": {
										Type:     schema.TypeList,
										Optional: true,
										Description: "Core node resource specifications. The number of blocks determines `CoreCount` passed to `CreateCluster`. " +
											"**All blocks must have identical configuration**; the first block is used as the single resource template by the API. " +
											"A validation error is returned at create time if any block differs from the first.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"instance_type": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "CVM instance type, e.g., `S6.2XLARGE32`, `SA4.8XLARGE64`.",
												},
												"system_disk": {
													Type:        schema.TypeList,
													Optional:    true,
													MaxItems:    1,
													Description: "System disk specifications.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"disk_size": {
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "Disk size in GB.",
															},
															"disk_type": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Disk type. Valid values: `CLOUD_SSD`, `CLOUD_PREMIUM`, `CLOUD_BASIC`, `LOCAL_BASIC`, `LOCAL_SSD`.",
															},
														},
													},
												},
												"data_disk": {
													Type:        schema.TypeSet,
													Optional:    true,
													Description: "Cloud data disk specifications.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"count": {
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "Number of disks.",
															},
															"disk_size": {
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "Disk size in GB.",
															},
															"disk_type": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Disk type. Valid values: `CLOUD_SSD`, `CLOUD_PREMIUM`, `CLOUD_BASIC`, `LOCAL_BASIC`, `LOCAL_SSD`, `CLOUD_HSSD`, `CLOUD_THROUGHPUT`, `CLOUD_TSSD`, `CLOUD_BIGDATA`, `CLOUD_HIGHIO`, `CLOUD_BSSD`, `REMOTE_SSD`.",
															},
														},
													},
												},
												// "tags": {
												// 	Type:        schema.TypeList,
												// 	Optional:    true,
												// 	Description: "Tags to bind to the node.",
												// 	Elem: &schema.Resource{
												// 		Schema: map[string]*schema.Schema{
												// 			"tag_key": {
												// 				Type:        schema.TypeString,
												// 				Optional:    true,
												// 				Description: "Tag key.",
												// 			},
												// 			"tag_value": {
												// 				Type:        schema.TypeString,
												// 				Optional:    true,
												// 				Description: "Tag value.",
												// 			},
												// 		},
												// 	},
												// },
												"emr_resource_id": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "EMR node resource ID.",
												},
											},
										},
									},
									"task_resource_spec": {
										Type:     schema.TypeList,
										Optional: true,
										Description: "Task node resource specifications. The number of blocks determines `TaskCount` passed to `CreateCluster`. " +
											"**All blocks must have identical configuration**; the first block is used as the single resource template by the API. " +
											"A validation error is returned at create time if any block differs from the first.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"instance_type": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "CVM instance type, e.g., `S6.2XLARGE32`, `SA4.8XLARGE64`.",
												},
												"system_disk": {
													Type:        schema.TypeList,
													Optional:    true,
													MaxItems:    1,
													Description: "System disk specifications.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"disk_size": {
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "Disk size in GB.",
															},
															"disk_type": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Disk type. Valid values: `CLOUD_SSD`, `CLOUD_PREMIUM`, `CLOUD_BASIC`, `LOCAL_BASIC`, `LOCAL_SSD`.",
															},
														},
													},
												},
												"data_disk": {
													Type:        schema.TypeSet,
													Optional:    true,
													Description: "Cloud data disk specifications.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"count": {
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "Number of disks.",
															},
															"disk_size": {
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "Disk size in GB.",
															},
															"disk_type": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Disk type. Valid values: `CLOUD_SSD`, `CLOUD_PREMIUM`, `CLOUD_BASIC`, `LOCAL_BASIC`, `LOCAL_SSD`, `CLOUD_HSSD`, `CLOUD_THROUGHPUT`, `CLOUD_TSSD`, `CLOUD_BIGDATA`, `CLOUD_HIGHIO`, `CLOUD_BSSD`, `REMOTE_SSD`.",
															},
														},
													},
												},
												// "tags": {
												// 	Type:        schema.TypeList,
												// 	Optional:    true,
												// 	Description: "Tags to bind to the node.",
												// 	Elem: &schema.Resource{
												// 		Schema: map[string]*schema.Schema{
												// 			"tag_key": {
												// 				Type:        schema.TypeString,
												// 				Optional:    true,
												// 				Description: "Tag key.",
												// 			},
												// 			"tag_value": {
												// 				Type:        schema.TypeString,
												// 				Optional:    true,
												// 				Description: "Tag value.",
												// 			},
												// 		},
												// 	},
												// },
												"emr_resource_id": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "EMR node resource ID.",
												},
											},
										},
									},
									"common_resource_spec": {
										Type:     schema.TypeList,
										Optional: true,
										Description: "Common node resource specifications. The number of blocks determines `CommonCount` passed to `CreateCluster`. " +
											"**All blocks must have identical configuration**; the first block is used as the single resource template by the API. " +
											"A validation error is returned at create time if any block differs from the first.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"instance_type": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "CVM instance type, e.g., `S6.2XLARGE32`, `SA4.8XLARGE64`.",
												},
												"system_disk": {
													Type:        schema.TypeList,
													Optional:    true,
													MaxItems:    1,
													Description: "System disk specifications.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"disk_size": {
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "Disk size in GB.",
															},
															"disk_type": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Disk type. Valid values: `CLOUD_SSD`, `CLOUD_PREMIUM`, `CLOUD_BASIC`, `LOCAL_BASIC`, `LOCAL_SSD`.",
															},
														},
													},
												},
												"data_disk": {
													Type:        schema.TypeSet,
													Optional:    true,
													Description: "Cloud data disk specifications.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"count": {
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "Number of disks.",
															},
															"disk_size": {
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "Disk size in GB.",
															},
															"disk_type": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Disk type. Valid values: `CLOUD_SSD`, `CLOUD_PREMIUM`, `CLOUD_BASIC`, `LOCAL_BASIC`, `LOCAL_SSD`, `CLOUD_HSSD`, `CLOUD_THROUGHPUT`, `CLOUD_TSSD`, `CLOUD_BIGDATA`, `CLOUD_HIGHIO`, `CLOUD_BSSD`, `REMOTE_SSD`.",
															},
														},
													},
												},
												// "tags": {
												// 	Type:        schema.TypeList,
												// 	Optional:    true,
												// 	Description: "Tags to bind to the node.",
												// 	Elem: &schema.Resource{
												// 		Schema: map[string]*schema.Schema{
												// 			"tag_key": {
												// 				Type:        schema.TypeString,
												// 				Optional:    true,
												// 				Description: "Tag key.",
												// 			},
												// 			"tag_value": {
												// 				Type:        schema.TypeString,
												// 				Optional:    true,
												// 				Description: "Tag value.",
												// 			},
												// 		},
												// 	},
												// },
												"emr_resource_id": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "EMR node resource ID.",
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			"cos_bucket": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "COS bucket path, used when creating StarRocks storage-compute separation clusters.",
			},
			"load_balancer_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "CLB instance ID, e.g., `lb-xxxxxxxx`.",
			},
			"default_meta_version": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Default metadata DB version. Valid values: `mysql8`, `tdsql8`, `mysql5`.",
			},
			"need_cdb_audit": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Whether to enable database auditing.",
			},
			"sg_ip": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Security source IP, e.g., `10.0.0.0/8`.",
			},
			"partition_number": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Partition placement group partition number.",
			},
			"web_ui_version": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Service UI address version. `0`: single URL (default); `1`: all URLs.",
			},

			// Computed read-back fields (populated by Read from DescribeInstances).
			"status": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Cluster status code. `2` indicates the cluster is running.",
			},
			"cluster_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Cluster ID (same as the resource ID).",
			},
		},
	}
}

// -----------------------------------------------------------------------------
// Create
// -----------------------------------------------------------------------------

func resourceTencentCloudEmrClusterV2Create(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_emr_cluster_v2.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId    = tccommon.GetLogId(tccommon.ContextNil)
		ctx      = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service  = EMRService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		request  = emr.NewCreateClusterRequest()
		response = emr.NewCreateClusterResponse()
	)

	if v, ok := d.GetOk("product_version"); ok {
		request.ProductVersion = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("enable_support_ha_flag"); ok {
		request.EnableSupportHAFlag = helper.Bool(v.(bool))
	}

	if v, ok := d.GetOk("instance_name"); ok {
		request.InstanceName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("instance_charge_type"); ok {
		request.InstanceChargeType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("login_settings"); ok {
		for _, item := range v.([]interface{}) {
			loginMap := item.(map[string]interface{})
			loginSettings := emr.LoginSettings{}
			if val, ok := loginMap["password"].(string); ok && val != "" {
				loginSettings.Password = helper.String(val)
			}
			if val, ok := loginMap["public_key_id"].(string); ok && val != "" {
				loginSettings.PublicKeyId = helper.String(val)
			}
			request.LoginSettings = &loginSettings
		}
	}

	if v, ok := d.GetOk("scene_software_config"); ok {
		for _, item := range v.([]interface{}) {
			sceneMap := item.(map[string]interface{})
			sceneSoftwareConfig := emr.SceneSoftwareConfig{}
			if val, ok := sceneMap["software"].([]interface{}); ok {
				for _, s := range val {
					sceneSoftwareConfig.Software = append(sceneSoftwareConfig.Software, helper.String(s.(string)))
				}
			}
			if val, ok := sceneMap["scene_name"].(string); ok && val != "" {
				sceneSoftwareConfig.SceneName = helper.String(val)
			}
			request.SceneSoftwareConfig = &sceneSoftwareConfig
		}
	}

	if v, ok := d.GetOk("instance_charge_prepaid"); ok {
		for _, item := range v.([]interface{}) {
			prepaidMap := item.(map[string]interface{})
			prepaid := emr.InstanceChargePrepaid{}
			if val, ok := prepaidMap["period"].(int); ok && val != 0 {
				prepaid.Period = helper.IntInt64(val)
			}
			if val, ok := prepaidMap["renew_flag"].(bool); ok {
				prepaid.RenewFlag = helper.Bool(val)
			}
			request.InstanceChargePrepaid = &prepaid
		}
	}

	if v, ok := d.GetOk("security_group_ids"); ok {
		for _, id := range v.([]interface{}) {
			request.SecurityGroupIds = append(request.SecurityGroupIds, helper.String(id.(string)))
		}
	}

	if v, ok := d.GetOk("script_bootstrap_action_config"); ok {
		for _, item := range v.([]interface{}) {
			scriptMap := item.(map[string]interface{})
			script := emr.ScriptBootstrapActionConfig{}
			if val, ok := scriptMap["cos_file_uri"].(string); ok && val != "" {
				script.CosFileURI = helper.String(val)
			}
			if val, ok := scriptMap["execution_moment"].(string); ok && val != "" {
				script.ExecutionMoment = helper.String(val)
			}
			if val, ok := scriptMap["args"].([]interface{}); ok {
				for _, a := range val {
					script.Args = append(script.Args, helper.String(a.(string)))
				}
			}
			if val, ok := scriptMap["cos_file_name"].(string); ok && val != "" {
				script.CosFileName = helper.String(val)
			}
			if val, ok := scriptMap["remark"].(string); ok && val != "" {
				script.Remark = helper.String(val)
			}
			request.ScriptBootstrapActionConfig = append(request.ScriptBootstrapActionConfig, &script)
		}
	}

	if v, ok := d.GetOk("need_master_wan"); ok {
		request.NeedMasterWan = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("enable_remote_login_flag"); ok {
		request.EnableRemoteLoginFlag = helper.Bool(v.(bool))
	}

	if v, ok := d.GetOkExists("enable_kerberos_flag"); ok {
		request.EnableKerberosFlag = helper.Bool(v.(bool))
	}

	if v, ok := d.GetOk("custom_conf"); ok {
		request.CustomConf = helper.String(v.(string))
	}

	if v, ok := d.GetOk("tags"); ok {
		for _, item := range v.([]interface{}) {
			tagMap := item.(map[string]interface{})
			tag := emr.Tag{}
			if val, ok := tagMap["tag_key"].(string); ok && val != "" {
				tag.TagKey = helper.String(val)
			}
			if val, ok := tagMap["tag_value"].(string); ok && val != "" {
				tag.TagValue = helper.String(val)
			}
			request.Tags = append(request.Tags, &tag)
		}
	}

	if v, ok := d.GetOk("disaster_recover_group_ids"); ok {
		for _, id := range v.([]interface{}) {
			request.DisasterRecoverGroupIds = append(request.DisasterRecoverGroupIds, helper.String(id.(string)))
		}
	}

	if v, ok := d.GetOkExists("enable_cbs_encrypt_flag"); ok {
		request.EnableCbsEncryptFlag = helper.Bool(v.(bool))
	}

	if v, ok := d.GetOk("meta_db_info"); ok {
		for _, item := range v.([]interface{}) {
			metaMap := item.(map[string]interface{})
			metaDB := emr.CustomMetaDBInfo{}
			if val, ok := metaMap["meta_data_jdbc_url"].(string); ok && val != "" {
				metaDB.MetaDataJdbcUrl = helper.String(val)
			}
			if val, ok := metaMap["meta_data_user"].(string); ok && val != "" {
				metaDB.MetaDataUser = helper.String(val)
			}
			if val, ok := metaMap["meta_data_pass"].(string); ok && val != "" {
				metaDB.MetaDataPass = helper.String(val)
			}
			if val, ok := metaMap["meta_type"].(string); ok && val != "" {
				metaDB.MetaType = helper.String(val)
			}
			if val, ok := metaMap["unify_meta_instance_id"].(string); ok && val != "" {
				metaDB.UnifyMetaInstanceId = helper.String(val)
			}
			request.MetaDBInfo = &metaDB
		}
	}

	if v, ok := d.GetOk("depend_service"); ok {
		for _, item := range v.([]interface{}) {
			dsMap := item.(map[string]interface{})
			ds := emr.DependService{}
			if val, ok := dsMap["service_name"].(string); ok && val != "" {
				ds.ServiceName = helper.String(val)
			}
			if val, ok := dsMap["instance_id"].(string); ok && val != "" {
				ds.InstanceId = helper.String(val)
			}
			request.DependService = append(request.DependService, &ds)
		}
	}

	if v, ok := d.GetOk("zone_resource_configuration"); ok {
		// Validate that all blocks of the same role within each zone carry
		// identical configuration. The SDK only accepts a single resource
		// template per role; users express count by repeating identical blocks.
		for zoneIdx, item := range v.([]interface{}) {
			zrcMap := item.(map[string]interface{})
			allList, _ := zrcMap["all_node_resource_spec"].([]interface{})
			if len(allList) == 0 {
				continue
			}
			allMap := allList[0].(map[string]interface{})
			for _, role := range []string{"master_resource_spec", "core_resource_spec", "task_resource_spec", "common_resource_spec"} {
				rawList, _ := allMap[role].([]interface{})
				if err := validateEmrNodeResourceSpecUniformity(role, zoneIdx, rawList); err != nil {
					return err
				}
			}
		}

		// zone_tag is derived from list index:
		//   single-AZ (len==1): ""
		//   multi-AZ:  0→"master", 1→"standby", 2→"third-party"
		zrcList := v.([]interface{})
		zoneTags := []string{"master", "standby", "third-party"}
		isSingleAZ := len(zrcList) == 1

		for zoneIdx, item := range zrcList {
			zrcMap := item.(map[string]interface{})
			zrc := emr.ZoneResourceConfiguration{}

			if vpcList, ok := zrcMap["virtual_private_cloud"].([]interface{}); ok && len(vpcList) > 0 {
				vpcMap := vpcList[0].(map[string]interface{})
				vpc := emr.VirtualPrivateCloud{}
				if val, ok := vpcMap["vpc_id"].(string); ok && val != "" {
					vpc.VpcId = helper.String(val)
				}
				if val, ok := vpcMap["subnet_id"].(string); ok && val != "" {
					vpc.SubnetId = helper.String(val)
				}
				zrc.VirtualPrivateCloud = &vpc
			}

			if plList, ok := zrcMap["placement"].([]interface{}); ok && len(plList) > 0 {
				plMap := plList[0].(map[string]interface{})
				pl := emr.Placement{}
				if val, ok := plMap["zone"].(string); ok && val != "" {
					pl.Zone = helper.String(val)
				}
				if val, ok := plMap["project_id"].(int); ok {
					pl.ProjectId = helper.IntInt64(val)
				}
				zrc.Placement = &pl
			}

			// Automatically derive zone_tag from position; omit for single-AZ.
			if !isSingleAZ && zoneIdx < len(zoneTags) {
				zrc.ZoneTag = helper.String(zoneTags[zoneIdx])
			}

			if allList, ok := zrcMap["all_node_resource_spec"].([]interface{}); ok && len(allList) > 0 {
				allMap := allList[0].(map[string]interface{})
				allSpec := emr.AllNodeResourceSpec{}

				// The node count for each role is derived from the length of
				// its `*_resource_spec` list. The first block is taken as the
				// single resource template accepted by the SDK; all blocks must
				// have identical configuration (validated above).
				if rawList, ok := allMap["master_resource_spec"].([]interface{}); ok && len(rawList) > 0 {
					allSpec.MasterCount = helper.IntInt64(len(rawList))
					if firstMap, ok := rawList[0].(map[string]interface{}); ok {
						allSpec.MasterResourceSpec = buildEmrClusterV2NodeResourceSpec(firstMap)
					}
				}
				if rawList, ok := allMap["core_resource_spec"].([]interface{}); ok && len(rawList) > 0 {
					allSpec.CoreCount = helper.IntInt64(len(rawList))
					if firstMap, ok := rawList[0].(map[string]interface{}); ok {
						allSpec.CoreResourceSpec = buildEmrClusterV2NodeResourceSpec(firstMap)
					}
				}
				if rawList, ok := allMap["task_resource_spec"].([]interface{}); ok && len(rawList) > 0 {
					allSpec.TaskCount = helper.IntInt64(len(rawList))
					if firstMap, ok := rawList[0].(map[string]interface{}); ok {
						allSpec.TaskResourceSpec = buildEmrClusterV2NodeResourceSpec(firstMap)
					}
				}
				if rawList, ok := allMap["common_resource_spec"].([]interface{}); ok && len(rawList) > 0 {
					allSpec.CommonCount = helper.IntInt64(len(rawList))
					if firstMap, ok := rawList[0].(map[string]interface{}); ok {
						allSpec.CommonResourceSpec = buildEmrClusterV2NodeResourceSpec(firstMap)
					}
				}
				zrc.AllNodeResourceSpec = &allSpec
			}

			request.ZoneResourceConfiguration = append(request.ZoneResourceConfiguration, &zrc)
		}
	}

	if v, ok := d.GetOk("cos_bucket"); ok {
		request.CosBucket = helper.String(v.(string))
	}

	if v, ok := d.GetOk("load_balancer_id"); ok {
		request.LoadBalancerId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("default_meta_version"); ok {
		request.DefaultMetaVersion = helper.String(v.(string))
	}

	if v, ok := d.GetOk("need_cdb_audit"); ok {
		request.NeedCdbAudit = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("sg_ip"); ok {
		request.SgIP = helper.String(v.(string))
	}

	if v, ok := d.GetOk("partition_number"); ok {
		request.PartitionNumber = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("web_ui_version"); ok {
		request.WebUiVersion = helper.IntInt64(v.(int))
	}

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseEmrClient().CreateClusterWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Create emr cluster v2 failed, Response is nil."))
		}
		response = result
		return nil
	})
	if reqErr != nil {
		log.Printf("[CRITAL]%s create emr cluster v2 failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	if response.Response.InstanceId == nil {
		return fmt.Errorf("InstanceId is nil.")
	}

	instanceId := *response.Response.InstanceId
	d.SetId(instanceId)

	// Wait for async provisioning to finish: Clusters[0].Status == 2 (running).
	waitErr := resource.Retry(d.Timeout(schema.TimeoutCreate)-time.Minute, func() *resource.RetryError {
		cluster, e := service.DescribeEmrClusterV2ById(ctx, instanceId)
		if e != nil {
			return tccommon.RetryError(e)
		}
		if cluster == nil || cluster.Status == nil {
			return resource.RetryableError(fmt.Errorf("emr cluster v2 %s not yet visible, retrying", instanceId))
		}
		status := *cluster.Status
		switch status {
		case EmrInternetStatusCreated:
			return nil
		case 301, 302, EmrInternetStatusDeleted:
			return resource.NonRetryableError(fmt.Errorf("emr cluster v2 %s entered terminal status %d", instanceId, status))
		default:
			return resource.RetryableError(fmt.Errorf("emr cluster v2 %s status=%d, waiting for running (2)", instanceId, status))
		}
	})
	if waitErr != nil {
		log.Printf("[CRITAL]%s wait emr cluster v2 creation failed, reason:%+v", logId, waitErr)
		return waitErr
	}

	return resourceTencentCloudEmrClusterV2Read(d, meta)
}

// -----------------------------------------------------------------------------
// Read
// -----------------------------------------------------------------------------

func resourceTencentCloudEmrClusterV2Read(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_emr_cluster_v2.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = EMRService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	instanceId := d.Id()
	cluster, err := service.DescribeEmrClusterV2ById(ctx, instanceId)
	if err != nil {
		return err
	}

	if cluster == nil {
		log.Printf("[WARN]%s resource `tencentcloud_emr_cluster_v2` [%s] not found, please check if it has been deleted.\n", logId, instanceId)
		d.SetId("")
		return nil
	}

	if cluster.EmrVersion != nil {
		_ = d.Set("product_version", cluster.EmrVersion)
	}

	if cluster.Config != nil && cluster.Config.SupportHA != nil {
		_ = d.Set("enable_support_ha_flag", cluster.Config.SupportHA)
	}

	if cluster.ClusterName != nil {
		_ = d.Set("instance_name", cluster.ClusterName)
	}

	if cluster.ChargeType != nil {
		// 0/1 POSTPAID_BY_HOUR/PREPAID; the API returns the numeric code, keep the
		// user's configured string value and only round-trip when both map cleanly.
		switch *cluster.ChargeType {
		case 1:
			_ = d.Set("instance_charge_type", "PREPAID")
		case 0:
			_ = d.Set("instance_charge_type", "POSTPAID_BY_HOUR")
		}
	}

	if cluster.SceneName != nil {
		// Preserve the scene_name inside scene_software_config if the API returns it.
		if existing, ok := d.GetOk("scene_software_config"); ok {
			list := existing.([]interface{})
			if len(list) > 0 {
				item := list[0].(map[string]interface{})
				if cluster.Config != nil && cluster.Config.SoftInfo != nil {
					item["software"] = cluster.Config.SoftInfo
				}

				item["scene_name"] = *cluster.SceneName
				list[0] = item
				_ = d.Set("scene_software_config", list)
			}
		}
	}

	if cluster.Config.SecurityGroups != nil {
		_ = d.Set("security_group_ids", cluster.Config.SecurityGroups)
	}

	if cluster.Config.CbsEncrypt != nil {
		if *cluster.Config.CbsEncrypt == 1 {
			_ = d.Set("enable_cbs_encrypt_flag", true)
		} else {
			_ = d.Set("enable_cbs_encrypt_flag", false)
		}
	}

	if len(cluster.Tags) > 0 {
		tagList := make([]map[string]interface{}, 0, len(cluster.Tags))
		for _, t := range cluster.Tags {
			m := map[string]interface{}{}
			if t.TagKey != nil {
				m["tag_key"] = t.TagKey
			}
			if t.TagValue != nil {
				m["tag_value"] = t.TagValue
			}
			tagList = append(tagList, m)
		}
		_ = d.Set("tags", tagList)
	}

	if _, existing := d.GetOk("zone_resource_configuration"); existing {
		nodes, nerr := service.DescribeEmrClusterV2Nodes(ctx, instanceId, "all")
		if nerr != nil {
			log.Printf("[WARN]%s DescribeEmrClusterV2Nodes failed, skip zone_resource_configuration read: %v\n", logId, nerr)
		} else if len(nodes) > 0 {
			log.Printf("[DEBUG]%s emr cluster v2 [%s] has %d nodes\n", logId, instanceId, len(nodes))

			// Step 1: group nodes by Zone.
			zoneMap := make(map[string][]*emr.NodeHardwareInfo)
			for _, item := range nodes {
				if item.Zone == nil {
					continue
				}
				z := *item.Zone
				zoneMap[z] = append(zoneMap[z], item)
			}

			// Determine zone order:
			//   - If state already has zone_resource_configuration (normal apply/refresh),
			//     preserve the user-defined order from state (placement.zone sequence).
			//     Any API zone not present in state is appended at the end.
			//   - If state is absent (terraform import), fall back to API insertion order.
			zoneOrder := make([]string, 0, len(zoneMap))
			if raw, ok := d.GetOk("zone_resource_configuration"); ok {
				for _, item := range raw.([]interface{}) {
					zrcMap, _ := item.(map[string]interface{})
					if zrcMap == nil {
						continue
					}
					plList, _ := zrcMap["placement"].([]interface{})
					if len(plList) == 0 {
						continue
					}
					plMap, _ := plList[0].(map[string]interface{})
					if plMap == nil {
						continue
					}
					z, _ := plMap["zone"].(string)
					if z != "" && len(zoneMap[z]) > 0 {
						zoneOrder = append(zoneOrder, z)
					}
				}
				// Append any API zones not covered by state (e.g. newly added zones).
				inOrder := make(map[string]bool, len(zoneOrder))
				for _, z := range zoneOrder {
					inOrder[z] = true
				}
				for _, item := range nodes {
					if item.Zone == nil {
						continue
					}
					z := *item.Zone
					if !inOrder[z] {
						zoneOrder = append(zoneOrder, z)
						inOrder[z] = true
					}
				}
			} else {
				// Import path: use API insertion order.
				seen := make(map[string]bool)
				for _, item := range nodes {
					if item.Zone == nil {
						continue
					}
					z := *item.Zone
					if !seen[z] {
						zoneOrder = append(zoneOrder, z)
						seen[z] = true
					}
				}
			}

			// Step 2: build one configuration block per Zone.
			configuration := make([]map[string]interface{}, 0, len(zoneOrder))
			for _, zone := range zoneOrder {
				zoneNodes := zoneMap[zone]
				// Use the first node as the representative for zone-level attributes.
				representative := zoneNodes[0]
				node := make(map[string]interface{})

				// virtual_private_cloud
				vpcItem := make(map[string]interface{})
				if cluster.UniqVpcId != nil {
					vpcItem["vpc_id"] = cluster.UniqVpcId
				}
				if representative.SubnetInfo != nil && representative.SubnetInfo.SubnetId != nil {
					vpcItem["subnet_id"] = representative.SubnetInfo.SubnetId
				}
				node["virtual_private_cloud"] = []interface{}{vpcItem}

				// placement
				placementItem := make(map[string]interface{})
				placementItem["zone"] = zone
				if cluster.ProductId != nil {
					placementItem["project_id"] = cluster.ProductId
				}
				node["placement"] = []interface{}{placementItem}

				// all_node_resource_spec: group nodes in this Zone by Flag.
				// Flag: 1=master, 2=core, 3=task, 0=common
				roleKeyMap := map[int64]string{
					1: "master_resource_spec",
					2: "core_resource_spec",
					3: "task_resource_spec",
					0: "common_resource_spec",
				}
				grouped := make(map[int64][]*emr.NodeHardwareInfo)
				for _, n := range zoneNodes {
					if n.Flag == nil {
						continue
					}
					grouped[*n.Flag] = append(grouped[*n.Flag], n)
				}

				// Retrieve the state spec list for this zone (for ordering and
				// emr_resource_id-based matching). stateAllSpec[roleKey] is a
				// []interface{} of spec blocks from the previous state.
				stateAllSpec := map[string][]interface{}{}
				if raw, ok := d.GetOk("zone_resource_configuration"); ok {
					for _, zrcRaw := range raw.([]interface{}) {
						zrcMap, _ := zrcRaw.(map[string]interface{})
						if zrcMap == nil {
							continue
						}
						plList, _ := zrcMap["placement"].([]interface{})
						if len(plList) == 0 {
							continue
						}
						plMap, _ := plList[0].(map[string]interface{})
						if plMap == nil {
							continue
						}
						if sz, _ := plMap["zone"].(string); sz != zone {
							continue
						}
						allList, _ := zrcMap["all_node_resource_spec"].([]interface{})
						if len(allList) == 0 {
							break
						}
						allMap, _ := allList[0].(map[string]interface{})
						if allMap == nil {
							break
						}
						for _, rk := range []string{"master_resource_spec", "core_resource_spec", "task_resource_spec", "common_resource_spec"} {
							if sl, ok := allMap[rk].([]interface{}); ok {
								stateAllSpec[rk] = sl
							}
						}
						break
					}
				}

				allSpec := make(map[string]interface{})
				for flag, roleKey := range roleKeyMap {
					roleNodes := grouped[flag]

					// Sort role nodes by the numeric index in NameTag
					// (e.g. "master.0" < "master.1" < "master.10").
					// NameTag format: "<role>.<index>", e.g. "master.0", "core.2".
					// We extract the suffix after the last '.' and sort numerically.
					nameTagIndex := func(n *emr.NodeHardwareInfo) int {
						if n.NameTag == nil {
							return 0
						}
						tag := *n.NameTag
						if dot := strings.LastIndex(tag, "."); dot >= 0 {
							if idx, err := strconv.Atoi(tag[dot+1:]); err == nil {
								return idx
							}
						}
						return 0
					}
					sort.Slice(roleNodes, func(i, j int) bool {
						return nameTagIndex(roleNodes[i]) < nameTagIndex(roleNodes[j])
					})

					stateSpecs := stateAllSpec[roleKey] // may be nil

					// Determine whether state has emr_resource_id populated.
					// If any state spec carries a non-empty emr_resource_id, use
					// id-based matching; otherwise use positional (index) matching.
					useIdMatch := false
					for _, ss := range stateSpecs {
						ssMap, _ := ss.(map[string]interface{})
						if ssMap == nil {
							continue
						}
						if rid, _ := ssMap["emr_resource_id"].(string); rid != "" {
							useIdMatch = true
							break
						}
					}

					// Build emr_resource_id → state spec index map for id-based matching.
					stateByResourceId := map[string]map[string]interface{}{}
					if useIdMatch {
						for _, ss := range stateSpecs {
							ssMap, _ := ss.(map[string]interface{})
							if ssMap == nil {
								continue
							}
							if rid, _ := ssMap["emr_resource_id"].(string); rid != "" {
								stateByResourceId[rid] = ssMap
							}
						}
					}

					specs := make([]interface{}, 0, len(roleNodes))
					for idx, n := range roleNodes {
						specItem := make(map[string]interface{})
						specItem["instance_type"] = emrInstanceTypeFromNode(n)
						if n.EmrResourceId != nil {
							specItem["emr_resource_id"] = *n.EmrResourceId
						} else {
							specItem["emr_resource_id"] = ""
						}

						// system_disk
						if n.RootSize != nil && n.RootStorageType != nil {
							specItem["system_disk"] = []interface{}{
								map[string]interface{}{
									"disk_size": int(*n.RootSize),
									"disk_type": emrDiskTypeIntToString(*n.RootStorageType),
								},
							}
						} else {
							specItem["system_disk"] = []interface{}{}
						}

						// data_disk (TypeSet → []interface{})
						dataDisks := make([]interface{}, 0, len(n.MCMultiDisk))
						for _, d := range n.MCMultiDisk {
							if d == nil || d.Count == nil || d.Type == nil || d.Size == nil {
								continue
							}
							sizeStr := strings.TrimRightFunc(*d.Size, func(r rune) bool {
								return r < '0' || r > '9'
							})
							sizeGB := 0
							if v, err := strconv.Atoi(sizeStr); err == nil {
								sizeGB = v
							}
							dataDisks = append(dataDisks, map[string]interface{}{
								"count":     int(*d.Count),
								"disk_size": sizeGB,
								"disk_type": emrDiskTypeIntToString(*d.Type),
							})
						}
						specItem["data_disk"] = dataDisks

						// tags
						// nodeTags := make([]interface{}, 0, len(n.Tags))
						// for _, t := range n.Tags {
						// 	if t == nil {
						// 		continue
						// 	}
						// 	tagItem := map[string]interface{}{
						// 		"tag_key":   "",
						// 		"tag_value": "",
						// 	}
						// 	if t.TagKey != nil {
						// 		tagItem["tag_key"] = *t.TagKey
						// 	}
						// 	if t.TagValue != nil {
						// 		tagItem["tag_value"] = *t.TagValue
						// 	}
						// 	nodeTags = append(nodeTags, tagItem)
						// }
						// specItem["tags"] = nodeTags

						// Align with state to avoid spurious diffs:
						//   - id-based: match by emr_resource_id (post-update scenario)
						//   - index-based: match by position (post-create, all nodes identical)
						var stateSpec map[string]interface{}
						if useIdMatch && n.EmrResourceId != nil {
							stateSpec = stateByResourceId[*n.EmrResourceId]
						} else if !useIdMatch && idx < len(stateSpecs) {
							stateSpec, _ = stateSpecs[idx].(map[string]interface{})
						}

						// If a matching state spec exists, carry over emr_resource_id
						// (already set above from API) – primary value is always from API.
						// For fields the API returns authoritatively, API wins.
						// stateSpec is kept here as a hook for future non-API-returnable fields.
						_ = stateSpec

						specs = append(specs, specItem)
					}
					allSpec[roleKey] = specs
				}
				node["all_node_resource_spec"] = []interface{}{allSpec}

				configuration = append(configuration, node)
			}

			if err := d.Set("zone_resource_configuration", configuration); err != nil {
				log.Printf("[WARN]%s set zone_resource_configuration failed: %v\n", logId, err)
			}
		}
	}

	if cluster.ClusterId != nil {
		_ = d.Set("cluster_id", cluster.ClusterId)
	}

	if cluster.Status != nil {
		_ = d.Set("status", cluster.Status)
	}

	return nil
}

// -----------------------------------------------------------------------------
// Update (no-op)
// -----------------------------------------------------------------------------

func resourceTencentCloudEmrClusterV2Update(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_emr_cluster_v2.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	// NOTE: Modify APIs (scale-out/scale-in, ModifyResourcesTags, rename, ...)
	// will be added in a follow-up change. The Update handler is intentionally a
	// no-op for now so that schema fields can stay non-ForceNew without forcing
	// recreation on accidental re-apply.
	return resourceTencentCloudEmrClusterV2Read(d, meta)
}

// -----------------------------------------------------------------------------
// Delete
// -----------------------------------------------------------------------------

func resourceTencentCloudEmrClusterV2Delete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_emr_cluster_v2.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = EMRService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		request = emr.NewTerminateInstanceRequest()
	)

	instanceId := d.Id()
	request.InstanceId = helper.String(instanceId)

	// First termination call.
	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseEmrClient().TerminateInstanceWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		}
		if result != nil {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if reqErr != nil {
		log.Printf("[CRITAL]%s terminate emr cluster v2 failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	// Wait until the cluster enters a terminated status or disappears.
	waitErr := resource.Retry(d.Timeout(schema.TimeoutDelete)-time.Minute, func() *resource.RetryError {
		cluster, e := service.DescribeEmrClusterV2ById(ctx, instanceId)
		if e != nil {
			return tccommon.RetryError(e)
		}
		if cluster == nil {
			return nil
		}
		if cluster.Status != nil && *cluster.Status == EmrInternetStatusDeleted {
			return nil
		}
		return resource.RetryableError(fmt.Errorf("emr cluster v2 %s being destroyed", instanceId))
	})
	if waitErr != nil {
		log.Printf("[CRITAL]%s wait emr cluster v2 deletion failed, reason:%+v", logId, waitErr)
		return waitErr
	}

	// For PREPAID clusters a second termination is required to actually release
	// the underlying resources. The polling condition is stricter: done only when
	// the instance is no longer found (nil result) OR the SDK returns
	// ResourceNotFound.InstanceNotFound.
	if v, ok := d.GetOk("instance_charge_type"); ok && v.(string) == "PREPAID" {
		reqErr2 := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseEmrClient().TerminateInstanceWithContext(ctx, request)
			if e != nil {
				if sdkErr, ok := e.(*errors.TencentCloudSDKError); ok {
					if sdkErr.GetCode() == "ResourceNotFound.InstanceNotFound" {
						return nil
					}
				}
				return tccommon.RetryError(e)
			}
			if result != nil {
				log.Printf("[DEBUG]%s api[%s] (2nd) success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}
			return nil
		})
		if reqErr2 != nil {
			log.Printf("[CRITAL]%s terminate emr cluster v2 (2nd) failed, reason:%+v", logId, reqErr2)
			return reqErr2
		}

		waitErr2 := resource.Retry(d.Timeout(schema.TimeoutDelete)-time.Minute, func() *resource.RetryError {
			cluster, e := service.DescribeEmrClusterV2ById(ctx, instanceId)
			if e != nil {
				if sdkErr, ok := e.(*errors.TencentCloudSDKError); ok {
					if sdkErr.GetCode() == "ResourceNotFound.InstanceNotFound" {
						return nil
					}
				}
				return tccommon.RetryError(e)
			}
			if cluster == nil {
				return nil
			}
			return resource.RetryableError(fmt.Errorf("emr cluster v2 %s (PREPAID) still being destroyed", instanceId))
		})
		if waitErr2 != nil {
			log.Printf("[CRITAL]%s wait emr cluster v2 (PREPAID 2nd) deletion failed, reason:%+v", logId, waitErr2)
			return waitErr2
		}
	}

	return nil
}

// -----------------------------------------------------------------------------
// Helper functions
// -----------------------------------------------------------------------------

// buildEmrClusterV2NodeResourceSpec converts a single Terraform node resource
// spec block (the first element of a *_resource_spec list) into an
// *emr.NodeResourceSpec. The caller is responsible for passing specList[0].
func buildEmrClusterV2NodeResourceSpec(specMap map[string]interface{}) *emr.NodeResourceSpec {
	if specMap == nil {
		return nil
	}
	spec := &emr.NodeResourceSpec{}
	if val, ok := specMap["instance_type"].(string); ok && val != "" {
		spec.InstanceType = helper.String(val)
	}
	if val, ok := specMap["system_disk"].([]interface{}); ok && len(val) > 0 {
		// system_disk is MaxItems:1; count is always 1 (schema field removed).
		disk := buildEmrClusterV2DiskSpecInfo(val[0])
		if disk != nil {
			disk.Count = helper.IntInt64(1)
			spec.SystemDisk = append(spec.SystemDisk, disk)
		}
	}
	if val, ok := specMap["data_disk"].(*schema.Set); ok && val != nil {
		for _, d := range val.List() {
			spec.DataDisk = append(spec.DataDisk, buildEmrClusterV2DiskSpecInfo(d))
		}
	}
	// if val, ok := specMap["tags"].([]interface{}); ok {
	// 	for _, t := range val {
	// 		tagMap, ok := t.(map[string]interface{})
	// 		if !ok {
	// 			continue
	// 		}
	// 		tag := &emr.Tag{}
	// 		if tk, ok := tagMap["tag_key"].(string); ok && tk != "" {
	// 			tag.TagKey = helper.String(tk)
	// 		}
	// 		if tv, ok := tagMap["tag_value"].(string); ok && tv != "" {
	// 			tag.TagValue = helper.String(tv)
	// 		}
	// 		spec.Tags = append(spec.Tags, tag)
	// 	}
	// }
	return spec
}

func buildEmrClusterV2DiskSpecInfo(raw interface{}) *emr.DiskSpecInfo {
	diskMap, ok := raw.(map[string]interface{})
	if !ok {
		return nil
	}
	disk := &emr.DiskSpecInfo{}
	if val, ok := diskMap["count"].(int); ok && val != 0 {
		disk.Count = helper.IntInt64(val)
	}
	if val, ok := diskMap["disk_size"].(int); ok && val != 0 {
		disk.DiskSize = helper.IntInt64(val)
	}
	if val, ok := diskMap["disk_type"].(string); ok && val != "" {
		disk.DiskType = helper.String(val)
	}
	return disk
}

// validateEmrNodeResourceSpecUniformity checks that every block in specList is
// identical to the first block. This is required because the SDK only accepts a
// single NodeResourceSpec template per role; users express the desired node count
// by repeating blocks with the same content.
//
// zoneIdx is used only for human-readable error messages.
func validateEmrNodeResourceSpecUniformity(role string, zoneIdx int, specList []interface{}) error {
	if len(specList) <= 1 {
		return nil
	}
	first, ok := specList[0].(map[string]interface{})
	if !ok {
		return nil
	}
	for i := 1; i < len(specList); i++ {
		cur, ok := specList[i].(map[string]interface{})
		if !ok {
			continue
		}
		if err := emrNodeResourceSpecEqual(first, cur); err != nil {
			return fmt.Errorf(
				"When creating a cluster, the attribute values must be consistent for all `master_desource_spec`, `core_desource_spec`, `task_desource_spec`, or `common_desource_spec` parameters in the `zone_resource_configuration[n].all_node_resource_spec` parameter when multiple parameters are entered.\n"+
					"zone_resource_configuration[%d].all_node_resource_spec.%s[%d] differs from [0]: %s",
				zoneIdx, role, i, err.Error(),
			)
		}
	}
	return nil
}

// emrNodeResourceSpecEqual returns a non-nil error describing the first
// field difference found between two node resource spec maps.
func emrNodeResourceSpecEqual(a, b map[string]interface{}) error {
	if a["instance_type"] != b["instance_type"] {
		return fmt.Errorf("instance_type mismatch: %q vs %q", a["instance_type"], b["instance_type"])
	}
	// system_disk is TypeList (MaxItems:1); compare by index.
	aSD, _ := a["system_disk"].([]interface{})
	bSD, _ := b["system_disk"].([]interface{})
	if len(aSD) != len(bSD) {
		return fmt.Errorf("system_disk: block count mismatch (%d vs %d)", len(aSD), len(bSD))
	}
	for i := range aSD {
		aD, _ := aSD[i].(map[string]interface{})
		bD, _ := bSD[i].(map[string]interface{})
		for _, k := range []string{"disk_size", "disk_type"} {
			if aD[k] != bD[k] {
				return fmt.Errorf("system_disk[%d].%s mismatch: %v vs %v", i, k, aD[k], bD[k])
			}
		}
	}
	// data_disk is TypeSet; compare as unordered bags of {count,disk_size,disk_type}.
	diskKey := func(m map[string]interface{}) string {
		return fmt.Sprintf("%v|%v|%v", m["count"], m["disk_size"], m["disk_type"])
	}
	aDS, aIsSet := a["data_disk"].(*schema.Set)
	bDS, bIsSet := b["data_disk"].(*schema.Set)
	var aDisks, bDisks []map[string]interface{}
	if aIsSet && aDS != nil {
		for _, item := range aDS.List() {
			if m, ok := item.(map[string]interface{}); ok {
				aDisks = append(aDisks, m)
			}
		}
	}
	if bIsSet && bDS != nil {
		for _, item := range bDS.List() {
			if m, ok := item.(map[string]interface{}); ok {
				bDisks = append(bDisks, m)
			}
		}
	}
	if len(aDisks) != len(bDisks) {
		return fmt.Errorf("data_disk: block count mismatch (%d vs %d)", len(aDisks), len(bDisks))
	}
	aKeys := make(map[string]int, len(aDisks))
	for _, m := range aDisks {
		aKeys[diskKey(m)]++
	}
	for _, m := range bDisks {
		k := diskKey(m)
		if aKeys[k] <= 0 {
			return fmt.Errorf("data_disk: entry {%s} in second spec not found in first", k)
		}
		aKeys[k]--
	}
	aTags, _ := a["tags"].([]interface{})
	bTags, _ := b["tags"].([]interface{})
	if len(aTags) != len(bTags) {
		return fmt.Errorf("tags: block count mismatch (%d vs %d)", len(aTags), len(bTags))
	}
	for i := range aTags {
		aT, _ := aTags[i].(map[string]interface{})
		bT, _ := bTags[i].(map[string]interface{})
		if aT["tag_key"] != bT["tag_key"] || aT["tag_value"] != bT["tag_value"] {
			return fmt.Errorf("tags[%d] mismatch: {%v:%v} vs {%v:%v}",
				i, aT["tag_key"], aT["tag_value"], bT["tag_key"], bT["tag_value"])
		}
	}
	return nil
}

// emrInstanceTypeFromNode reconstructs the instance_type string from the three
// node fields Spec / CpuNum / MemDesc.
//
//	Example: Spec="CVM.S6", CpuNum=8, MemDesc="16GB"  →  "S6.2XLARGE16"
func emrInstanceTypeFromNode(n *emr.NodeHardwareInfo) string {
	if n.Spec == nil || n.CpuNum == nil || n.MemDesc == nil {
		return ""
	}
	// Extract the suffix after '.', e.g. "CVM.S6" → "S6"
	specParts := strings.SplitN(*n.Spec, ".", 2)
	specSuffix := specParts[len(specParts)-1]

	cpuSizeLabel := emrCpuSizeLabel(*n.CpuNum)

	// Strip non-digit suffix from MemDesc, e.g. "16GB" → "16"
	memGB := strings.TrimRightFunc(*n.MemDesc, func(r rune) bool {
		return r < '0' || r > '9'
	})

	return fmt.Sprintf("%s.%s%s", specSuffix, cpuSizeLabel, memGB)
}

// emrCpuSizeLabel maps CPU count to the EMR instance size label.
//
//	1→SMALL, 2→MEDIUM, 4→LARGE, 8→2XLARGE, 16→4XLARGE, 32→8XLARGE, ...
func emrCpuSizeLabel(cpu int64) string {
	switch cpu {
	case 1:
		return "SMALL"
	case 2:
		return "MEDIUM"
	case 4:
		return "LARGE"
	default:
		if cpu > 4 {
			mult := cpu / 4
			return fmt.Sprintf("%dXLARGE", mult)
		}
		return fmt.Sprintf("%dXLARGE", cpu)
	}
}

// emrDiskTypeIntToString maps the integer StorageType / RootStorageType
// returned by DescribeClusterNodes to the string disk type accepted by CreateCluster.
func emrDiskTypeIntToString(t int64) string {
	switch t {
	case 1:
		return "LOCAL_BASIC"
	case 2:
		return "CLOUD_BASIC"
	case 3:
		return "LOCAL_SSD"
	case 4:
		return "CLOUD_SSD"
	case 5:
		return "CLOUD_PREMIUM"
	case 6:
		return "CLOUD_HSSD"
	case 11:
		return "CLOUD_THROUGHPUT"
	case 12:
		return "CLOUD_TSSD"
	case 13:
		return "CLOUD_BSSD"
	case 14:
		return "CLOUD_BIGDATA"
	case 15:
		return "CLOUD_HIGHIO"
	case 16:
		return "REMOTE_SSD"
	default:
		return ""
	}
}
