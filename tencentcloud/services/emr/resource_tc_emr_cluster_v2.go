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
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
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
				ForceNew:    true,
				Description: "EMR product version name, e.g., `EMR-V3.5.0`.",
			},
			"enable_support_ha_flag": {
				Type:        schema.TypeBool,
				Required:    true,
				ForceNew:    true,
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
				ForceNew:    true,
				Description: "Instance billing mode. Valid values: `PREPAID` (monthly/yearly subscription), `POSTPAID_BY_HOUR` (pay-as-you-go).",
			},
			"login_settings": {
				Type:        schema.TypeList,
				Required:    true,
				ForceNew:    true,
				MaxItems:    1,
				Description: "Login settings for purchased nodes.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"password": {
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
							Sensitive:   true,
							Description: "Login password, 8-16 characters, must contain uppercase, lowercase, digit and special character (supported special chars: `!@%^*`). First char cannot be special.",
						},
						"public_key_id": {
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
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
							ForceNew:    true,
							Description: "Scenario name, e.g., `Hadoop-Default`, `Hadoop-Kudu`, `Hadoop-Zookeeper`, `Hadoop-Presto`, `Hadoop-Hbase`.",
						},
					},
				},
			},
			"instance_charge_prepaid": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				MaxItems:    1,
				Description: "Prepaid (monthly/yearly) billing parameters. Required when `instance_charge_type` is `PREPAID`.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"period": {
							Type:        schema.TypeInt,
							Optional:    true,
							ForceNew:    true,
							Description: "Purchase duration in months. Valid values: 1-12, 24, 36, 48, 60.",
						},
						"renew_flag": {
							Type:        schema.TypeBool,
							Optional:    true,
							ForceNew:    true,
							Description: "Auto-renew flag. Default is false.",
						},
					},
				},
			},
			"security_group_ids": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				MaxItems:    1,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Security group IDs bound to the instance, e.g., `[\"sg-xxxxxxxx\"]`.",
			},
			"script_bootstrap_action_config": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				Description: "Bootstrap script configurations.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cos_file_uri": {
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
							Description: "COS URI of the script.",
						},
						"execution_moment": {
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
							Description: "Execution timing. Valid values: `resourceAfter`, `clusterAfter`, `clusterBefore`.",
						},
						"args": {
							Type:        schema.TypeList,
							Optional:    true,
							ForceNew:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "Script arguments, following standard Shell convention.",
						},
						"cos_file_name": {
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
							Description: "Script file name.",
						},
						"remark": {
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
							Description: "Remark.",
						},
					},
				},
			},
			"need_master_wan": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Whether to enable master public network. Valid values: `NEED_MASTER_WAN` (default), `NOT_NEED_MASTER_WAN`.",
			},
			"enable_remote_login_flag": {
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
				Description: "Whether to enable external remote login. Invalid when `security_group_ids` is set. Default is false.",
			},
			"enable_kerberos_flag": {
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
				Description: "Whether to enable Kerberos authentication. Default is false.",
			},
			"custom_conf": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
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
				ForceNew:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Spread placement group IDs. Currently supports only one ID.",
			},
			"enable_cbs_encrypt_flag": {
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
				Description: "Whether to enable cluster-level CBS encryption. Default is false.",
			},
			"meta_db_info": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				MaxItems:    1,
				Description: "Metadata database information. When `meta_type` is `EMR_NEW_META`/`EMR_DEFAULT_META`, no extra fields are required; when `EMR_EXIT_META`, `unify_meta_instance_id` must be set; when `USER_CUSTOM_META`, `meta_data_jdbc_url`/`meta_data_user`/`meta_data_pass` must be set.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"meta_data_jdbc_url": {
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
							Description: "Custom MetaDB JDBC URL, e.g., `jdbc:mysql://10.10.10.10:3306/dbname`.",
						},
						"meta_data_user": {
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
							Description: "Custom MetaDB username.",
						},
						"meta_data_pass": {
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
							Sensitive:   true,
							Description: "Custom MetaDB password.",
						},
						"meta_type": {
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
							Description: "Hive shared metadata DB type. Valid values: `EMR_DEFAULT_META`, `EMR_EXIT_META`, `USER_CUSTOM_META`.",
						},
						"unify_meta_instance_id": {
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
							Description: "EMR-MetaDB instance ID.",
						},
					},
				},
			},
			"depend_service": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				Description: "Shared component dependency information.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"service_name": {
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
							Description: "Shared component name.",
						},
						"instance_id": {
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
							Description: "Shared component cluster instance ID.",
						},
					},
				},
			},
			"zone_resource_configuration": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				MaxItems:    3,
				Description: "Per-zone resource configuration. Supports 1 (single-AZ) or up to 3 entries (multi-AZ: primary, backup, arbitration). ZoneTag is derived automatically from the list index.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"virtual_private_cloud": {
							Type:        schema.TypeList,
							Optional:    true,
							ForceNew:    true,
							MaxItems:    1,
							Description: "VPC/Subnet information.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"vpc_id": {
										Type:        schema.TypeString,
										Optional:    true,
										ForceNew:    true,
										Description: "VPC ID.",
									},
									"subnet_id": {
										Type:        schema.TypeString,
										Optional:    true,
										ForceNew:    true,
										Description: "Subnet ID.",
									},
								},
							},
						},
						"placement": {
							Type:        schema.TypeList,
							Optional:    true,
							ForceNew:    true,
							MaxItems:    1,
							Description: "Zone and project placement.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"zone": {
										Type:        schema.TypeString,
										Optional:    true,
										ForceNew:    true,
										Description: "Availability zone, e.g., `ap-guangzhou-7`.",
									},
									"project_id": {
										Type:        schema.TypeInt,
										Optional:    true,
										Computed:    true,
										ForceNew:    true,
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
															"disk_id": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Disk ID.",
															},
														},
													},
												},
												"data_disk": {
													Type:             schema.TypeList,
													Optional:         true,
													DiffSuppressFunc: emrDataDiskOrderSuppressFunc,
													Description:      "Cloud data disk specifications.",
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
																Description: "Disk type. Valid values: `CLOUD_SSD`, `CLOUD_PREMIUM`, `CLOUD_BASIC`, `LOCAL_BASIC`, `LOCAL_SSD`, `CLOUD_HSSD`, `CLOUD_THROUGHPUT`, `CLOUD_TSSD`, `CLOUD_BIGDATA`, `CLOUD_HIGHIO`, `CLOUD_BSSD`, `REMOTE_SSD`.",
															},
															"disk_id": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Disk ID.",
															},
														},
													},
												},
												"emr_resource_id": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "EMR node resource ID.",
												},
												"order_no": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Machine instance ID.",
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
															"disk_id": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Disk ID.",
															},
														},
													},
												},
												"data_disk": {
													Type:             schema.TypeList,
													Optional:         true,
													DiffSuppressFunc: emrDataDiskOrderSuppressFunc,
													Description:      "Cloud data disk specifications.",
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
																Description: "Disk type. Valid values: `CLOUD_SSD`, `CLOUD_PREMIUM`, `CLOUD_BASIC`, `LOCAL_BASIC`, `LOCAL_SSD`, `CLOUD_HSSD`, `CLOUD_THROUGHPUT`, `CLOUD_TSSD`, `CLOUD_BIGDATA`, `CLOUD_HIGHIO`, `CLOUD_BSSD`, `REMOTE_SSD`.",
															},
															"disk_id": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Disk ID.",
															},
														},
													},
												},
												"emr_resource_id": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "EMR node resource ID.",
												},
												"order_no": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Machine instance ID.",
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
															"disk_id": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Disk ID.",
															},
														},
													},
												},
												"data_disk": {
													Type:             schema.TypeList,
													Optional:         true,
													DiffSuppressFunc: emrDataDiskOrderSuppressFunc,
													Description:      "Cloud data disk specifications.",
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
																Description: "Disk type. Valid values: `CLOUD_SSD`, `CLOUD_PREMIUM`, `CLOUD_BASIC`, `LOCAL_BASIC`, `LOCAL_SSD`, `CLOUD_HSSD`, `CLOUD_THROUGHPUT`, `CLOUD_TSSD`, `CLOUD_BIGDATA`, `CLOUD_HIGHIO`, `CLOUD_BSSD`, `REMOTE_SSD`.",
															},
															"disk_id": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Disk ID.",
															},
														},
													},
												},
												"emr_resource_id": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "EMR node resource ID.",
												},
												"order_no": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Machine instance ID.",
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
															"disk_id": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Disk ID.",
															},
														},
													},
												},
												"data_disk": {
													Type:             schema.TypeList,
													Optional:         true,
													DiffSuppressFunc: emrDataDiskOrderSuppressFunc,
													Description:      "Cloud data disk specifications.",
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
																Description: "Disk type. Valid values: `CLOUD_SSD`, `CLOUD_PREMIUM`, `CLOUD_BASIC`, `LOCAL_BASIC`, `LOCAL_SSD`, `CLOUD_HSSD`, `CLOUD_THROUGHPUT`, `CLOUD_TSSD`, `CLOUD_BIGDATA`, `CLOUD_HIGHIO`, `CLOUD_BSSD`, `REMOTE_SSD`.",
															},
															"disk_id": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Disk ID.",
															},
														},
													},
												},
												"emr_resource_id": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "EMR node resource ID.",
												},
												"order_no": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Machine instance ID.",
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
				ForceNew:    true,
				Description: "COS bucket path, used when creating StarRocks storage-compute separation clusters.",
			},
			"load_balancer_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "CLB instance ID, e.g., `lb-xxxxxxxx`.",
			},
			"default_meta_version": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Default metadata DB version. Valid values: `mysql8`, `tdsql8`, `mysql5`.",
			},
			"need_cdb_audit": {
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
				Description: "Whether to enable database auditing.",
			},
			"sg_ip": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Security source IP, e.g., `10.0.0.0/8`.",
			},
			"partition_number": {
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
				Description: "Partition placement group partition number.",
			},
			"web_ui_version": {
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
				Description: "Service UI address version. `0`: single URL (default); `1`: all URLs.",
			},
			"enable_cbs_sys_encrypt_flag": {
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
				Description: "Whether to enable CBS system encryption.",
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

	if v, ok := d.GetOkExists("enable_cbs_sys_encrypt_flag"); ok {
		request.EnableCbsSysEncryptFlag = helper.Bool(v.(bool))
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
				if cluster.ProjectId != nil {
					placementItem["project_id"] = cluster.ProjectId
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

						if n.OrderNo != nil {
							specItem["order_no"] = *n.OrderNo
						} else {
							specItem["order_no"] = ""
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

						// data_disk
						if n.OrderNo == nil {
							continue
						}

						dataResp, derr := service.DescribeEmrNodeDataDisks(ctx, instanceId, specItem["order_no"].(string))
						if derr != nil {
							return derr
						}

						// Convert API response to a list of disk maps (each has disk_id).
						apiDisks := make([]map[string]interface{}, 0, len(dataResp))
						for _, item := range dataResp {
							tmpObj := map[string]interface{}{
								"disk_size": 0,
								"disk_type": "",
								"disk_id":   "",
							}
							if item.DiskSize != nil {
								tmpObj["disk_size"] = int(*item.DiskSize)
							}
							if item.DiskType != nil {
								tmpObj["disk_type"] = *item.DiskType
							}
							if item.DiskId != nil {
								tmpObj["disk_id"] = *item.DiskId
							}
							apiDisks = append(apiDisks, tmpObj)
						}

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

						// Align with previous state to keep disk_id→index bindings stable.
						// stateSpec is fetched from the raw state layer (not the diff layer),
						// so its disk_id values are authoritative and not index-shifted.
						var stateSpec map[string]interface{}
						if useIdMatch && n.EmrResourceId != nil {
							stateSpec = stateByResourceId[*n.EmrResourceId]
						} else if !useIdMatch && idx < len(stateSpecs) {
							stateSpec, _ = stateSpecs[idx].(map[string]interface{})
						}
						var stateDataDisks []interface{}
						if stateSpec != nil {
							stateDataDisks, _ = stateSpec["data_disk"].([]interface{})
						}

						// Re-order API disks to match the previous state's disk_id order.
						// This ensures that each position in state always refers to the same
						// physical disk, so Terraform's per-index diff shows the correct change.
						// Newly attached disks (not in previous state) are appended at the end.
						specItem["data_disk"] = alignDataDisksToConfig(stateDataDisks, apiDisks)

						specs = append(specs, specItem)
					}

					// Re-order the freshly-built specs to match the node order that was
					// previously recorded in state (by emr_resource_id).  This keeps the
					// TypeList index of every existing node stable across refreshes so
					// that Terraform never sees a "node swap" diff.  Newly attached nodes
					// (emr_resource_id not found in state) are appended after all existing ones.
					specs = alignNodeSpecsToState(specs, stateSpecs)

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

func resourceTencentCloudEmrClusterV2Update(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_emr_cluster_v2.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service    = EMRService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		instanceId = d.Id()
	)

	if d.HasChange("instance_name") {
		request := emr.NewModifyInstanceBasicRequest()
		request.InstanceId = helper.String(instanceId)
		if v, ok := d.GetOk("instance_name"); ok {
			request.ClusterName = helper.String(v.(string))
		}

		reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseEmrClient().ModifyInstanceBasicWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			}
			if result != nil {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}
			return nil
		})

		if reqErr != nil {
			log.Printf("[CRITAL]%s update emr cluster name failed, reason:%+v", logId, reqErr)
			return reqErr
		}
	}

	if d.HasChange("tags") {
		oldInterface, newInterface := d.GetChange("tags")
		oldTagList := oldInterface.([]interface{})
		newTagsList := newInterface.([]interface{})
		region := meta.(tccommon.ProviderMeta).GetAPIV3Conn().Region
		resourceName := tccommon.BuildTagResourceName("emr", "emr-instance", region, instanceId)
		modifyResourceTags := &emr.ModifyResourceTags{
			ResourceId:     helper.String(instanceId),
			Resource:       helper.String(resourceName),
			ResourcePrefix: helper.String("emr"),
			ResourceRegion: helper.String(region),
			ServiceType:    helper.String("emr"),
		}

		deleteTags := make([]*emr.Tag, 0, len(oldTagList))
		addTags := make([]*emr.Tag, 0, len(newTagsList))
		modifyTags := make([]*emr.Tag, 0, len(newTagsList))

		// Build a map of old tags: key -> value
		oldTagMap := make(map[string]string, len(oldTagList))
		for _, oItem := range oldTagList {
			odMap := oItem.(map[string]interface{})
			k, _ := odMap["tag_key"].(string)
			v, _ := odMap["tag_value"].(string)
			oldTagMap[k] = v
		}

		// Build a map of new tags: key -> value
		newTagMap := make(map[string]string, len(newTagsList))
		for _, nItem := range newTagsList {
			ndMap := nItem.(map[string]interface{})
			k, _ := ndMap["tag_key"].(string)
			v, _ := ndMap["tag_value"].(string)
			newTagMap[k] = v
		}

		// key only in old -> deleteTags
		for k, v := range oldTagMap {
			if _, exists := newTagMap[k]; !exists {
				deleteTags = append(deleteTags, &emr.Tag{
					TagKey:   helper.String(k),
					TagValue: helper.String(v),
				})
			}
		}

		// key only in new -> addTags; key in both -> modifyTags
		for k, v := range newTagMap {
			if _, exists := oldTagMap[k]; !exists {
				addTags = append(addTags, &emr.Tag{
					TagKey:   helper.String(k),
					TagValue: helper.String(v),
				})
			} else {
				modifyTags = append(modifyTags, &emr.Tag{
					TagKey:   helper.String(k),
					TagValue: helper.String(v),
				})
			}
		}

		if len(deleteTags) > 0 {
			modifyResourceTags.DeleteTags = deleteTags
		}

		if len(addTags) > 0 {
			modifyResourceTags.AddTags = addTags
		}

		if len(modifyTags) > 0 {
			modifyResourceTags.ModifyTags = modifyTags
		}

		request := emr.NewModifyResourcesTagsRequest()
		response := emr.NewModifyResourcesTagsResponse()
		request.ModifyType = helper.String("Cluster")
		request.ModifyResourceTagsInfoList = []*emr.ModifyResourceTags{
			modifyResourceTags,
		}
		reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseEmrClient().ModifyResourcesTagsWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}

			if result == nil || result.Response == nil || result.Response.ClusterToFlowIdList == nil || len(result.Response.ClusterToFlowIdList) == 0 {
				return resource.NonRetryableError(fmt.Errorf("Update emr cluster modify tags failed, Response is nil."))
			}

			response = result
			return nil
		})

		if reqErr != nil {
			log.Printf("[CRITAL]%s update emr cluster tags failed, reason:%+v", logId, reqErr)
			return reqErr
		}

		if response.Response.ClusterToFlowIdList[0].FlowId == nil {
			return fmt.Errorf("Update emr cluster modify tags failed, FlowId is nil.")
		}

		// wait
		flowId := int64(*response.Response.ClusterToFlowIdList[0].FlowId)
		conf := tccommon.BuildStateChangeConf([]string{"0", "1"}, []string{"2", "-1"}, d.Timeout(schema.TimeoutUpdate)-time.Minute, time.Second, service.FlowStatusRefreshFunc(instanceId, strconv.FormatInt(flowId, 10), F_KEY_FLOW_ID, []string{}))
		if object, e := conf.WaitForState(); e != nil {
			return e
		} else if status, ok := object.(*int64); ok && status != nil && *status == -1 {
			return fmt.Errorf("Update emr cluster modify tags failed, flow total status is -1.")
		}
	}

	if d.HasChange("zone_resource_configuration") {
		oldZrcRaw, newZrcRaw := d.GetChange("zone_resource_configuration")
		oldZrcList := oldZrcRaw.([]interface{})
		newZrcList := newZrcRaw.([]interface{})

		if len(oldZrcList) != len(newZrcList) {
			return fmt.Errorf("the number of zone_resource_configuration blocks cannot be changed after creation (old: %d, new: %d)",
				len(oldZrcList), len(newZrcList))
		}

		for zoneIdx := 0; zoneIdx < len(newZrcList); zoneIdx++ {
			oldZrcMap, ok1 := oldZrcList[zoneIdx].(map[string]interface{})
			newZrcMap, ok2 := newZrcList[zoneIdx].(map[string]interface{})
			if !ok1 || !ok2 {
				continue
			}

			oldAllList, _ := oldZrcMap["all_node_resource_spec"].([]interface{})
			newAllList, _ := newZrcMap["all_node_resource_spec"].([]interface{})
			if len(oldAllList) == 0 || len(newAllList) == 0 {
				continue
			}

			oldAll := oldAllList[0].(map[string]interface{})
			newAll := newAllList[0].(map[string]interface{})

			// ---------- master_resource_spec ----------
			if oldMasterList, newMasterList, changed := nodeRoleChanged(oldAll, newAll, "master_resource_spec"); changed {
				// master_resource_spec does not support scaling (add or remove nodes).
				// Re-align by emr_resource_id first so that mid-list insertions in config
				// do not look like scaling.
				alignedOldMaster, alignedNewMaster, addedMaster := alignNodeListByResourceId(oldMasterList, newMasterList)
				// Any node present in old but absent from new is a removal attempt.
				removedMaster := len(oldMasterList) - len(alignedOldMaster)
				if len(addedMaster) > 0 || removedMaster > 0 {
					return fmt.Errorf("zone_resource_configuration[%d].all_node_resource_spec.master_resource_spec does not support scaling (add/remove nodes)", zoneIdx)
				}

				for nodeIdx := 0; nodeIdx < len(alignedNewMaster); nodeIdx++ {
					oldSpec, _ := alignedOldMaster[nodeIdx].(map[string]interface{})
					newSpec, _ := alignedNewMaster[nodeIdx].(map[string]interface{})

					// Prevent modifying system_disk of existing nodes (not supported by API).
					if rid, _ := oldSpec["emr_resource_id"].(string); rid != "" {
						if sysDiskChanged(oldSpec, newSpec) {
							return fmt.Errorf("modifying system_disk of an existing node (emr_resource_id=%s) is not supported", rid)
						}
					}
					if instanceTypeChanged(oldSpec, newSpec) {
						newInstanceType, _ := newSpec["instance_type"].(string)
						resourceId, _ := oldSpec["emr_resource_id"].(string)

						newCpu, newMem, err := emrParseInstanceTypeCpuMem(newInstanceType)
						if err != nil {
							return fmt.Errorf("zone[%d] master_resource_spec[%d]: failed to parse instance_type %q: %v",
								zoneIdx, nodeIdx, newInstanceType, err)
						}

						modifyReq := emr.NewModifyResourceRequest()
						modifyResp := emr.NewModifyResourceResponse()
						modifyReq.InstanceId = helper.String(instanceId)
						modifyReq.PayMode = emrPayModeFromChargeType(d.Get("instance_charge_type").(string))
						modifyReq.InstanceType = helper.String(newInstanceType)
						modifyReq.NewCpu = helper.Int64(newCpu)
						modifyReq.NewMem = helper.Int64(newMem)
						modifyReq.ResourceIdList = []*string{helper.String(resourceId)}

						reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
							result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseEmrClient().ModifyResourceWithContext(ctx, modifyReq)
							if e != nil {
								return tccommon.RetryError(e)
							}
							log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
								logId, modifyReq.GetAction(), modifyReq.ToJsonString(), result.ToJsonString())

							if result == nil || result.Response == nil || result.Response.TraceId == nil {
								return resource.NonRetryableError(fmt.Errorf("Update emr cluster modify resource failed, Response is nil."))
							}

							modifyResp = result
							return nil
						})
						if reqErr != nil {
							log.Printf("[CRITAL]%s zone[%d] master_resource_spec[%d] modify instance type failed: %+v",
								logId, zoneIdx, nodeIdx, reqErr)
							return reqErr
						}

						// wait
						traceId := *modifyResp.Response.TraceId
						conf := tccommon.BuildStateChangeConf([]string{"0", "1"}, []string{"2", "-1"}, d.Timeout(schema.TimeoutUpdate)-time.Minute, time.Second, service.FlowStatusRefreshFunc(instanceId, traceId, F_KEY_TRACE_ID, []string{}))
						if object, e := conf.WaitForState(); e != nil {
							return e
						} else if status, ok := object.(*int64); ok && status != nil && *status == -1 {
							return fmt.Errorf("EMR cluster flow failed (TraceId=%s), flow total status is -1.", traceId)
						}
					}

					if dataDiskChanged(oldSpec, newSpec) {
						orderNo, _ := oldSpec["order_no"].(string)
						oldDisks, _ := oldSpec["data_disk"].([]interface{})
						newDisks, _ := newSpec["data_disk"].([]interface{})

						if err := handleNodeDataDiskChange(ctx, meta, d, logId, instanceId, orderNo,
							oldDisks, newDisks, "master_resource_spec", zoneIdx, nodeIdx); err != nil {
							return err
						}
					}
				}
			}

			// ---------- core_resource_spec ----------
			if oldCoreList, newCoreList, changed := nodeRoleChanged(oldAll, newAll, "core_resource_spec"); changed {
				alignedOldCore, alignedNewCore, addedCore := alignNodeListByResourceId(oldCoreList, newCoreList)

				// Content changes on existing nodes.
				for nodeIdx := 0; nodeIdx < len(alignedNewCore); nodeIdx++ {
					oldSpec, _ := alignedOldCore[nodeIdx].(map[string]interface{})
					newSpec, _ := alignedNewCore[nodeIdx].(map[string]interface{})

					if rid, _ := oldSpec["emr_resource_id"].(string); rid != "" {
						if sysDiskChanged(oldSpec, newSpec) {
							return fmt.Errorf("modifying system_disk of an existing node (emr_resource_id=%s) is not supported", rid)
						}
					}
					if instanceTypeChanged(oldSpec, newSpec) {
						newInstanceType, _ := newSpec["instance_type"].(string)
						resourceId, _ := oldSpec["emr_resource_id"].(string)
						newCpu, newMem, err := emrParseInstanceTypeCpuMem(newInstanceType)
						if err != nil {
							return fmt.Errorf("zone[%d] core_resource_spec[%d]: failed to parse instance_type %q: %v", zoneIdx, nodeIdx, newInstanceType, err)
						}
						modifyReq := emr.NewModifyResourceRequest()
						modifyResp := emr.NewModifyResourceResponse()
						modifyReq.InstanceId = helper.String(instanceId)
						modifyReq.PayMode = emrPayModeFromChargeType(d.Get("instance_charge_type").(string))
						modifyReq.InstanceType = helper.String(newInstanceType)
						modifyReq.NewCpu = helper.Int64(newCpu)
						modifyReq.NewMem = helper.Int64(newMem)
						modifyReq.ResourceIdList = []*string{helper.String(resourceId)}
						reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
							result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseEmrClient().ModifyResourceWithContext(ctx, modifyReq)
							if e != nil {
								return tccommon.RetryError(e)
							}
							log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, modifyReq.GetAction(), modifyReq.ToJsonString(), result.ToJsonString())
							if result == nil || result.Response == nil || result.Response.TraceId == nil {
								return resource.NonRetryableError(fmt.Errorf("Update emr cluster modify resource failed, Response is nil."))
							}
							modifyResp = result
							return nil
						})
						if reqErr != nil {
							log.Printf("[CRITAL]%s zone[%d] core_resource_spec[%d] modify instance type failed: %+v", logId, zoneIdx, nodeIdx, reqErr)
							return reqErr
						}
						traceId := *modifyResp.Response.TraceId
						conf := tccommon.BuildStateChangeConf([]string{"0", "1"}, []string{"2", "-1"}, d.Timeout(schema.TimeoutUpdate)-time.Minute, time.Second, service.FlowStatusRefreshFunc(instanceId, traceId, F_KEY_TRACE_ID, []string{}))
						if object, e := conf.WaitForState(); e != nil {
							return e
						} else if status, ok := object.(*int64); ok && status != nil && *status == -1 {
							return fmt.Errorf("EMR cluster flow failed (TraceId=%s), flow total status is -1.", traceId)
						}
					}
					if dataDiskChanged(oldSpec, newSpec) {
						orderNo, _ := oldSpec["order_no"].(string)
						oldDisks, _ := oldSpec["data_disk"].([]interface{})
						newDisks, _ := newSpec["data_disk"].([]interface{})
						if err := handleNodeDataDiskChange(ctx, meta, d, logId, instanceId, orderNo, oldDisks, newDisks, "core_resource_spec", zoneIdx, nodeIdx); err != nil {
							return err
						}
					}
				}

				// Scale-out: new nodes.
				for addedIdx, addedRaw := range addedCore {
					addedSpec, _ := addedRaw.(map[string]interface{})
					scaleOutReq := emr.NewScaleOutClusterRequest()
					scaleOutResp := emr.NewScaleOutClusterResponse()
					scaleOutReq.InstanceId = helper.String(instanceId)
					scaleOutReq.InstanceChargeType = helper.String(d.Get("instance_charge_type").(string))
					scaleOutReq.ScaleOutNodeConfig = &emr.ScaleOutNodeConfig{NodeFlag: common.StringPtr("CORE"), NodeCount: common.Uint64Ptr(1)}
					if d.Get("instance_charge_type").(string) == "PREPAID" {
						if v, ok := d.GetOk("instance_charge_prepaid"); ok {
							prepaidList := v.([]interface{})
							if len(prepaidList) > 0 {
								prepaidMap := prepaidList[0].(map[string]interface{})
								prepaid := &emr.InstanceChargePrepaid{}
								if val, ok := prepaidMap["period"].(int); ok && val != 0 {
									prepaid.Period = helper.IntInt64(val)
								}
								if val, ok := prepaidMap["renew_flag"].(bool); ok {
									prepaid.RenewFlag = helper.Bool(val)
								}
								scaleOutReq.InstanceChargePrepaid = prepaid
							}
						}
					}
					coreResourceSpec := &emr.NodeResourceSpec{}
					if instanceType, ok := addedSpec["instance_type"].(string); ok && instanceType != "" {
						coreResourceSpec.InstanceType = helper.String(instanceType)
					}
					if sdList, ok := addedSpec["system_disk"].([]interface{}); ok && len(sdList) > 0 {
						sdMap, _ := sdList[0].(map[string]interface{})
						coreResourceSpec.SystemDisk = []*emr.DiskSpecInfo{{DiskType: helper.String(sdMap["disk_type"].(string)), Count: helper.Int64(1), DiskSize: helper.Int64(int64(sdMap["disk_size"].(int)))}}
					}
					if ddList, ok := addedSpec["data_disk"].([]interface{}); ok {
						for _, ddItem := range ddList {
							ddMap, _ := ddItem.(map[string]interface{})
							coreResourceSpec.DataDisk = append(coreResourceSpec.DataDisk, &emr.DiskSpecInfo{DiskType: helper.String(ddMap["disk_type"].(string)), Count: helper.Int64(1), DiskSize: helper.Int64(int64(ddMap["disk_size"].(int)))})
						}
					}
					scaleOutReq.ResourceSpec = coreResourceSpec
					if plList, ok := newZrcMap["placement"].([]interface{}); ok && len(plList) > 0 {
						plMap, _ := plList[0].(map[string]interface{})
						if zone, ok := plMap["zone"].(string); ok && zone != "" {
							scaleOutReq.Zone = helper.String(zone)
						}
					}
					scaleOutErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
						result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseEmrClient().ScaleOutClusterWithContext(ctx, scaleOutReq)
						if e != nil {
							return tccommon.RetryError(e)
						}
						log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, scaleOutReq.GetAction(), scaleOutReq.ToJsonString(), result.ToJsonString())
						if result == nil || result.Response == nil || result.Response.TraceId == nil {
							return resource.NonRetryableError(fmt.Errorf("scale out core nodes failed, Response is nil"))
						}
						scaleOutResp = result
						return nil
					})
					if scaleOutErr != nil {
						log.Printf("[CRITAL]%s zone[%d] core_resource_spec[added:%d] scale-out failed: %+v", logId, zoneIdx, addedIdx, scaleOutErr)
						return scaleOutErr
					}
					traceId := *scaleOutResp.Response.TraceId
					conf := tccommon.BuildStateChangeConf([]string{"0", "1"}, []string{"2", "-1"}, d.Timeout(schema.TimeoutUpdate)-time.Minute, time.Second, service.FlowStatusRefreshFunc(instanceId, traceId, F_KEY_TRACE_ID, []string{}))
					if object, e := conf.WaitForState(); e != nil {
						return e
					} else if status, ok := object.(*int64); ok && status != nil && *status == -1 {
						return fmt.Errorf("EMR cluster flow failed (TraceId=%s), flow total status is -1.", traceId)
					}
				}

				// Scale-in: nodes removed — find old nodes whose emr_resource_id is not referenced by new list.
				newResourceIDs := make(map[string]bool, len(newCoreList))
				for _, raw := range newCoreList {
					m, _ := raw.(map[string]interface{})
					if rid, _ := m["emr_resource_id"].(string); rid != "" {
						newResourceIDs[rid] = true
					}
				}
				var removedCoreOrderNos []*string
				for _, raw := range oldCoreList {
					m, _ := raw.(map[string]interface{})
					rid, _ := m["emr_resource_id"].(string)
					orderNo, _ := m["order_no"].(string)
					if rid != "" && !newResourceIDs[rid] && orderNo != "" {
						removedCoreOrderNos = append(removedCoreOrderNos, helper.String(orderNo))
					}
				}
				if len(removedCoreOrderNos) > 0 {
					terminateReq := emr.NewTerminateClusterNodesRequest()
					terminateResp := emr.NewTerminateClusterNodesResponse()
					terminateReq.InstanceId = helper.String(instanceId)
					terminateReq.CvmInstanceIds = removedCoreOrderNos
					terminateReq.NodeFlag = helper.String("CORE")
					terminateErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
						result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseEmrClient().TerminateClusterNodesWithContext(ctx, terminateReq)
						if e != nil {
							return tccommon.RetryError(e)
						}
						log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, terminateReq.GetAction(), terminateReq.ToJsonString(), result.ToJsonString())
						if result == nil || result.Response == nil || result.Response.FlowId == nil {
							return resource.NonRetryableError(fmt.Errorf("terminate core nodes failed, Response is nil"))
						}
						terminateResp = result
						return nil
					})
					if terminateErr != nil {
						log.Printf("[CRITAL]%s zone[%d] core_resource_spec scale-in failed: %+v", logId, zoneIdx, terminateErr)
						return terminateErr
					}
					flowId := int64(*terminateResp.Response.FlowId)
					conf := tccommon.BuildStateChangeConf([]string{"0", "1"}, []string{"2", "-1"}, d.Timeout(schema.TimeoutUpdate)-time.Minute, time.Second, service.FlowStatusRefreshFunc(instanceId, strconv.FormatInt(flowId, 10), F_KEY_FLOW_ID, []string{}))
					if object, e := conf.WaitForState(); e != nil {
						return e
					} else if status, ok := object.(*int64); ok && status != nil && *status == -1 {
						return fmt.Errorf("EMR cluster flow failed (FlowId=%s), flow total status is -1.", strconv.FormatInt(flowId, 10))
					}
				}
			}

			// ---------- task_resource_spec ----------
			if oldTaskList, newTaskList, changed := nodeRoleChanged(oldAll, newAll, "task_resource_spec"); changed {
				alignedOldTask, alignedNewTask, addedTask := alignNodeListByResourceId(oldTaskList, newTaskList)

				// Content changes on existing nodes.
				for nodeIdx := 0; nodeIdx < len(alignedNewTask); nodeIdx++ {
					oldSpec, _ := alignedOldTask[nodeIdx].(map[string]interface{})
					newSpec, _ := alignedNewTask[nodeIdx].(map[string]interface{})

					if rid, _ := oldSpec["emr_resource_id"].(string); rid != "" {
						if sysDiskChanged(oldSpec, newSpec) {
							return fmt.Errorf("modifying system_disk of an existing node (emr_resource_id=%s) is not supported", rid)
						}
					}
					if instanceTypeChanged(oldSpec, newSpec) {
						newInstanceType, _ := newSpec["instance_type"].(string)
						resourceId, _ := oldSpec["emr_resource_id"].(string)
						newCpu, newMem, err := emrParseInstanceTypeCpuMem(newInstanceType)
						if err != nil {
							return fmt.Errorf("zone[%d] task_resource_spec[%d]: failed to parse instance_type %q: %v", zoneIdx, nodeIdx, newInstanceType, err)
						}
						modifyReq := emr.NewModifyResourceRequest()
						modifyResp := emr.NewModifyResourceResponse()
						modifyReq.InstanceId = helper.String(instanceId)
						modifyReq.PayMode = emrPayModeFromChargeType(d.Get("instance_charge_type").(string))
						modifyReq.InstanceType = helper.String(newInstanceType)
						modifyReq.NewCpu = helper.Int64(newCpu)
						modifyReq.NewMem = helper.Int64(newMem)
						modifyReq.ResourceIdList = []*string{helper.String(resourceId)}
						reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
							result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseEmrClient().ModifyResourceWithContext(ctx, modifyReq)
							if e != nil {
								return tccommon.RetryError(e)
							}
							log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, modifyReq.GetAction(), modifyReq.ToJsonString(), result.ToJsonString())
							if result == nil || result.Response == nil || result.Response.TraceId == nil {
								return resource.NonRetryableError(fmt.Errorf("Update emr cluster modify resource failed, Response is nil."))
							}
							modifyResp = result
							return nil
						})
						if reqErr != nil {
							log.Printf("[CRITAL]%s zone[%d] task_resource_spec[%d] modify instance type failed: %+v", logId, zoneIdx, nodeIdx, reqErr)
							return reqErr
						}
						traceId := *modifyResp.Response.TraceId
						conf := tccommon.BuildStateChangeConf([]string{"0", "1"}, []string{"2", "-1"}, d.Timeout(schema.TimeoutUpdate)-time.Minute, time.Second, service.FlowStatusRefreshFunc(instanceId, traceId, F_KEY_TRACE_ID, []string{}))
						if object, e := conf.WaitForState(); e != nil {
							return e
						} else if status, ok := object.(*int64); ok && status != nil && *status == -1 {
							return fmt.Errorf("EMR cluster flow failed (TraceId=%s), flow total status is -1.", traceId)
						}
					}
					if dataDiskChanged(oldSpec, newSpec) {
						orderNo, _ := oldSpec["order_no"].(string)
						oldDisks, _ := oldSpec["data_disk"].([]interface{})
						newDisks, _ := newSpec["data_disk"].([]interface{})
						if err := handleNodeDataDiskChange(ctx, meta, d, logId, instanceId, orderNo, oldDisks, newDisks, "task_resource_spec", zoneIdx, nodeIdx); err != nil {
							return err
						}
					}
				}

				// Scale-out: new nodes.
				for addedIdx, addedRaw := range addedTask {
					addedSpec, _ := addedRaw.(map[string]interface{})
					scaleOutReq := emr.NewScaleOutClusterRequest()
					scaleOutResp := emr.NewScaleOutClusterResponse()
					scaleOutReq.InstanceId = helper.String(instanceId)
					scaleOutReq.InstanceChargeType = helper.String(d.Get("instance_charge_type").(string))
					scaleOutReq.ScaleOutNodeConfig = &emr.ScaleOutNodeConfig{NodeFlag: common.StringPtr("TASK"), NodeCount: common.Uint64Ptr(1)}
					if d.Get("instance_charge_type").(string) == "PREPAID" {
						if v, ok := d.GetOk("instance_charge_prepaid"); ok {
							prepaidList := v.([]interface{})
							if len(prepaidList) > 0 {
								prepaidMap := prepaidList[0].(map[string]interface{})
								prepaid := &emr.InstanceChargePrepaid{}
								if val, ok := prepaidMap["period"].(int); ok && val != 0 {
									prepaid.Period = helper.IntInt64(val)
								}
								if val, ok := prepaidMap["renew_flag"].(bool); ok {
									prepaid.RenewFlag = helper.Bool(val)
								}
								scaleOutReq.InstanceChargePrepaid = prepaid
							}
						}
					}
					taskResourceSpec := &emr.NodeResourceSpec{}
					if instanceType, ok := addedSpec["instance_type"].(string); ok && instanceType != "" {
						taskResourceSpec.InstanceType = helper.String(instanceType)
					}
					if sdList, ok := addedSpec["system_disk"].([]interface{}); ok && len(sdList) > 0 {
						sdMap, _ := sdList[0].(map[string]interface{})
						taskResourceSpec.SystemDisk = []*emr.DiskSpecInfo{{DiskType: helper.String(sdMap["disk_type"].(string)), Count: helper.Int64(1), DiskSize: helper.Int64(int64(sdMap["disk_size"].(int)))}}
					}
					if ddList, ok := addedSpec["data_disk"].([]interface{}); ok {
						for _, ddItem := range ddList {
							ddMap, _ := ddItem.(map[string]interface{})
							taskResourceSpec.DataDisk = append(taskResourceSpec.DataDisk, &emr.DiskSpecInfo{DiskType: helper.String(ddMap["disk_type"].(string)), Count: helper.Int64(1), DiskSize: helper.Int64(int64(ddMap["disk_size"].(int)))})
						}
					}
					scaleOutReq.ResourceSpec = taskResourceSpec
					if plList, ok := newZrcMap["placement"].([]interface{}); ok && len(plList) > 0 {
						plMap, _ := plList[0].(map[string]interface{})
						if zone, ok := plMap["zone"].(string); ok && zone != "" {
							scaleOutReq.Zone = helper.String(zone)
						}
					}
					scaleOutErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
						result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseEmrClient().ScaleOutClusterWithContext(ctx, scaleOutReq)
						if e != nil {
							return tccommon.RetryError(e)
						}
						log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, scaleOutReq.GetAction(), scaleOutReq.ToJsonString(), result.ToJsonString())
						if result == nil || result.Response == nil || result.Response.TraceId == nil {
							return resource.NonRetryableError(fmt.Errorf("scale out task nodes failed, Response is nil"))
						}
						scaleOutResp = result
						return nil
					})
					if scaleOutErr != nil {
						log.Printf("[CRITAL]%s zone[%d] task_resource_spec[added:%d] scale-out failed: %+v", logId, zoneIdx, addedIdx, scaleOutErr)
						return scaleOutErr
					}
					traceId := *scaleOutResp.Response.TraceId
					conf := tccommon.BuildStateChangeConf([]string{"0", "1"}, []string{"2", "-1"}, d.Timeout(schema.TimeoutUpdate)-time.Minute, time.Second, service.FlowStatusRefreshFunc(instanceId, traceId, F_KEY_TRACE_ID, []string{}))
					if object, e := conf.WaitForState(); e != nil {
						return e
					} else if status, ok := object.(*int64); ok && status != nil && *status == -1 {
						return fmt.Errorf("EMR cluster flow failed (TraceId=%s), flow total status is -1.", traceId)
					}
				}

				// Scale-in: nodes removed.
				newTaskResourceIDs := make(map[string]bool, len(newTaskList))
				for _, raw := range newTaskList {
					m, _ := raw.(map[string]interface{})
					if rid, _ := m["emr_resource_id"].(string); rid != "" {
						newTaskResourceIDs[rid] = true
					}
				}
				var removedTaskOrderNos []*string
				for _, raw := range oldTaskList {
					m, _ := raw.(map[string]interface{})
					rid, _ := m["emr_resource_id"].(string)
					orderNo, _ := m["order_no"].(string)
					if rid != "" && !newTaskResourceIDs[rid] && orderNo != "" {
						removedTaskOrderNos = append(removedTaskOrderNos, helper.String(orderNo))
					}
				}
				if len(removedTaskOrderNos) > 0 {
					terminateReq := emr.NewTerminateClusterNodesRequest()
					terminateResp := emr.NewTerminateClusterNodesResponse()
					terminateReq.InstanceId = helper.String(instanceId)
					terminateReq.CvmInstanceIds = removedTaskOrderNos
					terminateReq.NodeFlag = helper.String("TASK")
					terminateErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
						result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseEmrClient().TerminateClusterNodesWithContext(ctx, terminateReq)
						if e != nil {
							return tccommon.RetryError(e)
						}
						log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, terminateReq.GetAction(), terminateReq.ToJsonString(), result.ToJsonString())
						if result == nil || result.Response == nil || result.Response.FlowId == nil {
							return resource.NonRetryableError(fmt.Errorf("terminate task nodes failed, Response is nil"))
						}
						terminateResp = result
						return nil
					})
					if terminateErr != nil {
						log.Printf("[CRITAL]%s zone[%d] task_resource_spec scale-in failed: %+v", logId, zoneIdx, terminateErr)
						return terminateErr
					}
					flowId := int64(*terminateResp.Response.FlowId)
					conf := tccommon.BuildStateChangeConf([]string{"0", "1"}, []string{"2", "-1"}, d.Timeout(schema.TimeoutUpdate)-time.Minute, time.Second, service.FlowStatusRefreshFunc(instanceId, strconv.FormatInt(flowId, 10), F_KEY_FLOW_ID, []string{}))
					if object, e := conf.WaitForState(); e != nil {
						return e
					} else if status, ok := object.(*int64); ok && status != nil && *status == -1 {
						return fmt.Errorf("EMR cluster flow failed (FlowId=%s), flow total status is -1.", strconv.FormatInt(flowId, 10))
					}
				}
			}

			// ---------- common_resource_spec ----------
			if oldCommonList, newCommonList, changed := nodeRoleChanged(oldAll, newAll, "common_resource_spec"); changed {
				// common_resource_spec does not support scaling (add or remove nodes).
				alignedOldCommon, alignedNewCommon, addedCommon := alignNodeListByResourceId(oldCommonList, newCommonList)
				removedCommon := len(oldCommonList) - len(alignedOldCommon)
				if len(addedCommon) > 0 || removedCommon > 0 {
					return fmt.Errorf("zone_resource_configuration[%d].all_node_resource_spec.common_resource_spec does not support scaling (add/remove nodes)", zoneIdx)
				}

				for nodeIdx := 0; nodeIdx < len(alignedNewCommon); nodeIdx++ {
					oldSpec, _ := alignedOldCommon[nodeIdx].(map[string]interface{})
					newSpec, _ := alignedNewCommon[nodeIdx].(map[string]interface{})

					if rid, _ := oldSpec["emr_resource_id"].(string); rid != "" {
						if sysDiskChanged(oldSpec, newSpec) {
							return fmt.Errorf("modifying system_disk of an existing node (emr_resource_id=%s) is not supported", rid)
						}
					}
					if instanceTypeChanged(oldSpec, newSpec) {
						newInstanceType, _ := newSpec["instance_type"].(string)
						resourceId, _ := oldSpec["emr_resource_id"].(string)
						newCpu, newMem, err := emrParseInstanceTypeCpuMem(newInstanceType)
						if err != nil {
							return fmt.Errorf("zone[%d] common_resource_spec[%d]: failed to parse instance_type %q: %v", zoneIdx, nodeIdx, newInstanceType, err)
						}
						modifyReq := emr.NewModifyResourceRequest()
						modifyResp := emr.NewModifyResourceResponse()
						modifyReq.InstanceId = helper.String(instanceId)
						modifyReq.PayMode = emrPayModeFromChargeType(d.Get("instance_charge_type").(string))
						modifyReq.InstanceType = helper.String(newInstanceType)
						modifyReq.NewCpu = helper.Int64(newCpu)
						modifyReq.NewMem = helper.Int64(newMem)
						modifyReq.ResourceIdList = []*string{helper.String(resourceId)}
						reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
							result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseEmrClient().ModifyResourceWithContext(ctx, modifyReq)
							if e != nil {
								return tccommon.RetryError(e)
							}
							log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, modifyReq.GetAction(), modifyReq.ToJsonString(), result.ToJsonString())
							if result == nil || result.Response == nil || result.Response.TraceId == nil {
								return resource.NonRetryableError(fmt.Errorf("Update emr cluster modify resource failed, Response is nil."))
							}
							modifyResp = result
							return nil
						})
						if reqErr != nil {
							log.Printf("[CRITAL]%s zone[%d] common_resource_spec[%d] modify instance type failed: %+v", logId, zoneIdx, nodeIdx, reqErr)
							return reqErr
						}
						traceId := *modifyResp.Response.TraceId
						conf := tccommon.BuildStateChangeConf([]string{"0", "1"}, []string{"2", "-1"}, d.Timeout(schema.TimeoutUpdate)-time.Minute, time.Second, service.FlowStatusRefreshFunc(instanceId, traceId, F_KEY_TRACE_ID, []string{}))
						if object, e := conf.WaitForState(); e != nil {
							return e
						} else if status, ok := object.(*int64); ok && status != nil && *status == -1 {
							return fmt.Errorf("EMR cluster flow failed (TraceId=%s), flow total status is -1.", traceId)
						}
					}
					if dataDiskChanged(oldSpec, newSpec) {
						orderNo, _ := oldSpec["order_no"].(string)
						oldDisks, _ := oldSpec["data_disk"].([]interface{})
						newDisks, _ := newSpec["data_disk"].([]interface{})
						if err := handleNodeDataDiskChange(ctx, meta, d, logId, instanceId, orderNo, oldDisks, newDisks, "common_resource_spec", zoneIdx, nodeIdx); err != nil {
							return err
						}
					}
				}
			}

		}
	}

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
	if val, ok := specMap["data_disk"].([]interface{}); ok && val != nil {
		for _, d := range val {
			spec.DataDisk = append(spec.DataDisk, buildEmrClusterV2DiskSpecInfo(d))
		}
	}
	return spec
}

func buildEmrClusterV2DiskSpecInfo(raw interface{}) *emr.DiskSpecInfo {
	diskMap, ok := raw.(map[string]interface{})
	if !ok {
		return nil
	}
	disk := &emr.DiskSpecInfo{}
	disk.Count = helper.IntInt64(1)
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
	// data_disk is TypeList; compare by index in declaration order.
	aDataDisks, _ := a["data_disk"].([]interface{})
	bDataDisks, _ := b["data_disk"].([]interface{})
	if len(aDataDisks) != len(bDataDisks) {
		return fmt.Errorf("data_disk: block count mismatch (%d vs %d)", len(aDataDisks), len(bDataDisks))
	}
	for i := range aDataDisks {
		aD, _ := aDataDisks[i].(map[string]interface{})
		bD, _ := bDataDisks[i].(map[string]interface{})
		for _, k := range []string{"count", "disk_size", "disk_type"} {
			if aD[k] != bD[k] {
				return fmt.Errorf("data_disk[%d].%s mismatch: %v vs %v", i, k, aD[k], bD[k])
			}
		}
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

// nodeRoleChanged returns the old/new spec lists for a given node role key and
// whether any difference exists between them.
func nodeRoleChanged(oldAll, newAll map[string]interface{}, roleKey string) (oldList, newList []interface{}, changed bool) {
	oldList, _ = oldAll[roleKey].([]interface{})
	newList, _ = newAll[roleKey].([]interface{})
	if len(oldList) == 0 && len(newList) == 0 {
		return oldList, newList, false
	}
	return oldList, newList, fmt.Sprintf("%v", oldList) != fmt.Sprintf("%v", newList)
}

// firstNodeSpec returns the first element of a node spec list as a map, or an
// empty map if the list is empty.
// instanceTypeChanged reports whether the instance_type field differs between
// the old and new node spec maps.
func instanceTypeChanged(oldSpec, newSpec map[string]interface{}) bool {
	oldType, _ := oldSpec["instance_type"].(string)
	newType, _ := newSpec["instance_type"].(string)
	return oldType != newType
}

// diskChanged reports whether system_disk or data_disk fields differ between
// the old and new node spec maps.
func dataDiskChanged(oldSpec, newSpec map[string]interface{}) bool {
	return fmt.Sprintf("%v", oldSpec["data_disk"]) != fmt.Sprintf("%v", newSpec["data_disk"])
}

// emrParseInstanceTypeCpuMem parses an EMR instance type string (e.g. "S6.2XLARGE32")
// and returns the corresponding CPU count and memory (GB).
//
// Format: {spec}.{sizeLabel}{memGB}
// sizeLabel → cpu:  SMALL=1, MEDIUM=2, LARGE=4, {n}XLARGE=n*4
func emrParseInstanceTypeCpuMem(instanceType string) (cpu, mem int64, err error) {
	// Split on the first '.': "S6.2XLARGE32" → ["S6", "2XLARGE32"]
	parts := strings.SplitN(instanceType, ".", 2)
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("invalid instance_type format %q: expected {spec}.{sizeLabel}{mem}", instanceType)
	}
	suffix := parts[1] // e.g. "2XLARGE32"

	// Extract trailing digits as memory (GB).
	memStr := ""
	labelStr := suffix
	for i := len(suffix) - 1; i >= 0; i-- {
		if suffix[i] >= '0' && suffix[i] <= '9' {
			memStr = string(suffix[i]) + memStr
		} else {
			labelStr = suffix[:i+1]
			break
		}
	}
	if memStr == "" {
		return 0, 0, fmt.Errorf("invalid instance_type format %q: cannot parse memory", instanceType)
	}
	memVal, e := strconv.ParseInt(memStr, 10, 64)
	if e != nil {
		return 0, 0, fmt.Errorf("invalid instance_type format %q: %v", instanceType, e)
	}

	// Parse sizeLabel → cpu count.
	var cpuVal int64
	switch labelStr {
	case "SMALL":
		cpuVal = 1
	case "MEDIUM":
		cpuVal = 2
	case "LARGE":
		cpuVal = 4
	default:
		// Pattern: {n}XLARGE  e.g. "2XLARGE" → n=2, cpu=n*4
		if strings.HasSuffix(labelStr, "XLARGE") {
			nStr := strings.TrimSuffix(labelStr, "XLARGE")
			n, e2 := strconv.ParseInt(nStr, 10, 64)
			if e2 != nil {
				return 0, 0, fmt.Errorf("invalid instance_type size label %q in %q", labelStr, instanceType)
			}
			cpuVal = n * 4
		} else {
			return 0, 0, fmt.Errorf("unknown instance_type size label %q in %q", labelStr, instanceType)
		}
	}

	return cpuVal, memVal, nil
}

// emrPayModeFromChargeType maps instance_charge_type to the PayMode integer
// expected by ModifyResource: PREPAID → 1, anything else (POSTPAID_BY_HOUR) → 0.
func emrPayModeFromChargeType(chargeType string) *uint64 {
	if chargeType == "PREPAID" {
		return helper.Uint64(1)
	}
	return helper.Uint64(0)
}

// handleNodeDataDiskChange handles data_disk changes for a single node by matching
// disks via disk_id rather than by index, so disk reordering is handled correctly.
//
//   - Existing disks (disk_id found in oldDisks): validate and resize if disk_size grew.
//   - New disks (disk_id empty or not in oldDisks): attach as new disk.
//   - Removed disks (disk_id in oldDisks but not in newDisks): return error (not supported).
func handleNodeDataDiskChange(
	ctx context.Context,
	meta interface{},
	d *schema.ResourceData,
	logId, instanceId, orderNo string,
	oldDisks, newDisks []interface{},
	roleKey string, zoneIdx, nodeIdx int,
) error {
	// Build old disk map: disk_id -> disk map
	oldDiskByID := make(map[string]map[string]interface{}, len(oldDisks))
	for _, item := range oldDisks {
		m, _ := item.(map[string]interface{})
		id, _ := m["disk_id"].(string)
		if id != "" {
			oldDiskByID[id] = m
		}
	}

	// Track which old disk_ids are referenced by new list (to detect removals)
	referencedOldIDs := make(map[string]bool, len(oldDisks))

	for diskIdx, item := range newDisks {
		newDisk, _ := item.(map[string]interface{})
		diskId, _ := newDisk["disk_id"].(string)
		newSize, _ := newDisk["disk_size"].(int)
		diskType, _ := newDisk["disk_type"].(string)

		if diskId != "" {
			// Existing disk — match by disk_id
			referencedOldIDs[diskId] = true
			oldDisk, exists := oldDiskByID[diskId]
			if !exists {
				continue
			}
			oldSize, _ := oldDisk["disk_size"].(int)

			if newSize < oldSize {
				return fmt.Errorf("%s[%d][%d] data_disk (disk_id=%s): shrinking disk_size is not supported (old: %d, new: %d)",
					roleKey, zoneIdx, nodeIdx, diskId, oldSize, newSize)
			}

			if newSize > oldSize {
				resizeReq := emr.NewResizeDataDisksRequest()
				resizeResp := emr.NewResizeDataDisksResponse()
				resizeReq.InstanceId = helper.String(instanceId)
				resizeReq.DiskSize = helper.Int64(int64(newSize))
				resizeReq.CvmInstanceIds = []*string{helper.String(orderNo)}
				resizeReq.ResizeAll = helper.Bool(false)
				resizeReq.DiskIds = []*string{helper.String(diskId)}

				resizeErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
					result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseEmrClient().ResizeDataDisksWithContext(ctx, resizeReq)
					if e != nil {
						return tccommon.RetryError(e)
					}
					log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
						logId, resizeReq.GetAction(), resizeReq.ToJsonString(), result.ToJsonString())

					if result == nil || result.Response == nil || result.Response.FlowId == nil {
						return resource.NonRetryableError(fmt.Errorf("resize data disks failed, Response is nil"))
					}

					resizeResp = result
					return nil
				})
				if resizeErr != nil {
					log.Printf("[CRITAL]%s %s[%d][%d] data_disk (disk_id=%s) resize failed: %+v",
						logId, roleKey, zoneIdx, nodeIdx, diskId, resizeErr)
					return resizeErr
				}

				// wait
				flowId := int64(*resizeResp.Response.FlowId)
				conf := tccommon.BuildStateChangeConf([]string{"0", "1"}, []string{"2", "-1"}, d.Timeout(schema.TimeoutUpdate)-time.Minute, time.Second,
					(&EMRService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}).FlowStatusRefreshFunc(instanceId, strconv.FormatInt(flowId, 10), F_KEY_FLOW_ID, []string{}))
				if object, e := conf.WaitForState(); e != nil {
					return e
				} else if status, ok := object.(*int64); ok && status != nil && *status == -1 {
					return fmt.Errorf("EMR cluster flow failed (FlowId=%s), flow total status is -1.", strconv.FormatInt(flowId, 10))
				}
			}
		} else {
			// New disk — disk_id is empty, attach as new disk
			attachReq := emr.NewAttachDisksRequest()
			attachResp := emr.NewAttachDisksResponse()
			attachReq.InstanceId = helper.String(instanceId)
			attachReq.CvmInstanceIds = []*string{helper.String(orderNo)}
			attachReq.CreateDisk = helper.Bool(true)
			attachReq.DiskSpec = &emr.NodeSpecDiskV2{
				Count:           helper.Int64(1),
				DiskType:        helper.String(diskType),
				DefaultDiskSize: helper.Int64(int64(newSize)),
			}

			attachErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
				result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseEmrClient().AttachDisksWithContext(ctx, attachReq)
				if e != nil {
					return tccommon.RetryError(e)
				}
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
					logId, attachReq.GetAction(), attachReq.ToJsonString(), result.ToJsonString())

				if result == nil || result.Response == nil || result.Response.FlowId == nil {
					return resource.NonRetryableError(fmt.Errorf("attach data disks failed, Response is nil"))
				}

				attachResp = result
				return nil
			})
			if attachErr != nil {
				log.Printf("[CRITAL]%s %s[%d][%d] attach disk[%d] failed: %+v",
					logId, roleKey, zoneIdx, nodeIdx, diskIdx, attachErr)
				return attachErr
			}

			// wait
			flowId := int64(*attachResp.Response.FlowId)
			conf := tccommon.BuildStateChangeConf([]string{"0", "1"}, []string{"2", "-1"}, d.Timeout(schema.TimeoutUpdate)-time.Minute, time.Second,
				(&EMRService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}).FlowStatusRefreshFunc(instanceId, strconv.FormatInt(flowId, 10), F_KEY_FLOW_ID, []string{}))
			if object, e := conf.WaitForState(); e != nil {
				return e
			} else if status, ok := object.(*int64); ok && status != nil && *status == -1 {
				return fmt.Errorf("EMR cluster flow failed (FlowId=%s), flow total status is -1.", strconv.FormatInt(flowId, 10))
			}
		}
	}

	// Detect removed disks (in old but not referenced by new)
	for id := range oldDiskByID {
		if !referencedOldIDs[id] {
			return fmt.Errorf("%s[%d][%d]: removing data_disk (disk_id=%s) is not supported",
				roleKey, zoneIdx, nodeIdx, id)
		}
	}

	return nil
}

// emrDataDiskOrderSuppressFunc suppresses diffs on data_disk list elements that
// are caused purely by the user's config having a different ordering than the
// state, rather than genuine value changes.
//
// Because Read writes state with disk_id-stable ordering (each index always
// refers to the same physical disk), a diff where old==new means the field
// value did not actually change — the apparent diff is only due to the user
// declaring disks in a different order than the state.  Such diffs are
// suppressed.
//
// Diffs where old!=new are always surfaced (genuine resize/retype).
// Purely new elements (index ≥ len(oldList)) are never suppressed.
// Removal (len(new) < len(old)) is never suppressed.
// emrDataDiskOrderSuppressFunc suppresses diffs on data_disk list elements that
// are caused purely by the user's config having a different ordering than the
// state, rather than genuine value changes.
//
// data_disk is treated as an unordered repeatable list.  Read writes state in
// disk_id-stable order (each index always refers to the same physical disk).
// The suppress function handles two cases where a diff is not genuine:
//
//  1. old == new at this position: the value did not change; the apparent diff
//     exists only because the user declared disks in a different order than the
//     state's disk_id-stable order.  Always suppress.
//
//  2. old != new, but old multiset ⊆ new multiset (only reordering / appending):
//     no disk was genuinely modified; the apparent change at this position is
//     caused by a neighbouring disk being shifted.  Suppress only for existing
//     positions (elemIdx < len(old)).  Newly added positions are never suppressed.
//
// Genuine changes (resize, retype, removal) are never suppressed.
func emrDataDiskOrderSuppressFunc(k, old, new string, d *schema.ResourceData) bool {
	// k example: "zone_resource_configuration.0.all_node_resource_spec.0.core_resource_spec.0.data_disk.1.disk_size"
	parts := strings.Split(k, ".")
	if len(parts) < 3 {
		return false
	}

	// Find the "data_disk" segment position.
	diskListIdx := -1
	for i := len(parts) - 1; i >= 0; i-- {
		if parts[i] == "data_disk" {
			diskListIdx = i
			break
		}
	}
	if diskListIdx < 0 || diskListIdx+1 >= len(parts) {
		return false
	}

	// Parse the element index within data_disk.
	elemIdx, err := strconv.Atoi(parts[diskListIdx+1])
	if err != nil {
		return false
	}

	listPath := strings.Join(parts[:diskListIdx+1], ".")
	oldRaw, newRaw := d.GetChange(listPath)
	oldList, _ := oldRaw.([]interface{})
	newList, _ := newRaw.([]interface{})

	// Removal: never suppress.
	if len(newList) < len(oldList) {
		return false
	}

	// No existing disks (new node being added): never suppress.
	if len(oldList) == 0 {
		return false
	}

	// Purely new element (beyond existing list bounds): never suppress.
	if elemIdx >= len(oldList) {
		return false
	}

	// Case 1: value unchanged at this position — ordering diff only.
	if old == new {
		return true
	}

	// Case 2: old != new — suppress only if old multiset ⊆ new multiset,
	// meaning no disk was removed or genuinely modified (just reordered).
	diskFP := func(item interface{}) string {
		m, _ := item.(map[string]interface{})
		return fmt.Sprintf("%v:%v", m["disk_size"], m["disk_type"])
	}
	oldCounts := make(map[string]int, len(oldList))
	for _, item := range oldList {
		oldCounts[diskFP(item)]++
	}
	newCounts := make(map[string]int, len(newList))
	for _, item := range newList {
		newCounts[diskFP(item)]++
	}
	for fp, cnt := range oldCounts {
		if newCounts[fp] < cnt {
			// A disk was removed or its size/type changed — genuine change.
			return false
		}
	}
	// Old multiset is a subset of the new multiset: only reordering or appending.
	return true
}

// sysDiskChanged reports whether system_disk fields differ between the old and new node spec maps.
func sysDiskChanged(oldSpec, newSpec map[string]interface{}) bool {
	return fmt.Sprintf("%v", oldSpec["system_disk"]) != fmt.Sprintf("%v", newSpec["system_disk"])
}

// alignDataDisksToConfig re-orders apiDisks to match the reference order given
// by refDisks (previous state or current config/diff), so that the state written
// by Read keeps the same physical disk at a stable TypeList index.
//
// Matching uses two passes:
//
//  1. disk_id exact match + disk_size agreement.  When refDisks comes from the
//     diff layer (Update→Read), Terraform inherits disk_id by index (not by
//     identity), so we also require size agreement to avoid false matches.
//
//  2. Greedy match by {disk_type, disk_size}.  Handles entries where disk_id is
//     absent or was not matched with correct size.
//
// Unmatched API disks (newly attached) are appended in API order at the end.
func alignDataDisksToConfig(refDisks []interface{}, apiDisks []map[string]interface{}) []interface{} {
	usedAPIIdx := make([]bool, len(apiDisks))

	// Build lookup: disk_id → api index.
	idxByDiskID := make(map[string]int, len(apiDisks))
	for i, d := range apiDisks {
		if id, _ := d["disk_id"].(string); id != "" {
			idxByDiskID[id] = i
		}
	}

	result := make([]interface{}, 0, len(apiDisks))

	for _, refRaw := range refDisks {
		refDisk, _ := refRaw.(map[string]interface{})
		if refDisk == nil {
			continue
		}

		refID, _ := refDisk["disk_id"].(string)
		refSize, _ := refDisk["disk_size"].(int)
		refType, _ := refDisk["disk_type"].(string)

		matched := false

		// Pass 1: disk_id + disk_size (guards against index-shifted Computed inheritance).
		if refID != "" {
			if apiIdx, ok := idxByDiskID[refID]; ok && !usedAPIIdx[apiIdx] {
				apiSize, _ := apiDisks[apiIdx]["disk_size"].(int)
				if apiSize == refSize {
					result = append(result, apiDisks[apiIdx])
					usedAPIIdx[apiIdx] = true
					matched = true
				}
			}
		}

		// Pass 2: greedy {disk_type, disk_size} match.
		if !matched {
			for i, apiDisk := range apiDisks {
				if usedAPIIdx[i] {
					continue
				}
				if apiDisk["disk_type"] == refType {
					apiSize, _ := apiDisk["disk_size"].(int)
					if apiSize == refSize {
						result = append(result, apiDisk)
						usedAPIIdx[i] = true
						matched = true
						break
					}
				}
			}
		}
		// No match: disk removed externally — omit.
	}

	// Append newly attached disks (not consumed by either pass) in API order.
	for i, apiDisk := range apiDisks {
		if !usedAPIIdx[i] {
			result = append(result, apiDisk)
		}
	}

	return result
}

// alignNodeSpecsToState re-orders freshly-built node specs (from the API) to
// match the order that was previously recorded in state, identified by
// emr_resource_id.
//
// This keeps the TypeList index of every existing node stable across refreshes:
// a node that was at position i in the last state will still be at position i
// after the next refresh, so Terraform never sees a spurious "node swap" diff.
// Newly added nodes (whose emr_resource_id is not present in state) are
// appended after all existing nodes in the order they appear in apiSpecs.
func alignNodeSpecsToState(apiSpecs []interface{}, stateSpecs []interface{}) []interface{} {
	if len(stateSpecs) == 0 {
		// No previous state — keep API order.
		return apiSpecs
	}

	// Build a map from emr_resource_id to the corresponding api spec.
	apiByID := make(map[string]interface{}, len(apiSpecs))
	for _, raw := range apiSpecs {
		m, _ := raw.(map[string]interface{})
		if m == nil {
			continue
		}
		if rid, _ := m["emr_resource_id"].(string); rid != "" {
			apiByID[rid] = raw
		}
	}

	result := make([]interface{}, 0, len(apiSpecs))
	usedIDs := make(map[string]bool, len(apiSpecs))

	// First pass: place existing nodes in state order.
	for _, stRaw := range stateSpecs {
		stMap, _ := stRaw.(map[string]interface{})
		if stMap == nil {
			continue
		}
		rid, _ := stMap["emr_resource_id"].(string)
		if rid == "" {
			continue
		}
		if apiSpec, ok := apiByID[rid]; ok {
			result = append(result, apiSpec)
			usedIDs[rid] = true
		}
		// If a state node is no longer returned by API, skip it (decommissioned).
	}

	// Second pass: append newly added nodes (not seen in previous state).
	for _, raw := range apiSpecs {
		m, _ := raw.(map[string]interface{})
		if m == nil {
			continue
		}
		rid, _ := m["emr_resource_id"].(string)
		if !usedIDs[rid] {
			result = append(result, raw)
		}
	}

	return result
}

// alignNodeListByResourceId re-aligns a pair of old/new TypeList node slices
// (from d.GetChange) so that each position holds the same physical node
// (identified by emr_resource_id) in both lists.
//
// Background: Terraform's TypeList SDK copies Computed fields (emr_resource_id,
// order_no) from state to diff by index, not by identity.  When the user inserts
// a new node in the middle of the config list, the SDK shifts all subsequent
// nodes, causing their emr_resource_id slots to be filled with wrong values (or
// empty).  This function recovers the correct pairing so that Update can
// accurately determine which nodes need resizing and which are brand-new.
//
// Returns:
//   - alignedOld: existing nodes in the same order as alignedNew
//   - alignedNew: config entries paired with their matching old node, with
//     genuinely new entries (emr_resource_id=="") appended at the end
//   - addedNew:   config entries that have no matching old node (scale-out)
func alignNodeListByResourceId(
	oldList, newList []interface{},
) (alignedOld, alignedNew []interface{}, addedNew []interface{}) {
	// Build old map: emr_resource_id → spec
	oldByID := make(map[string]interface{}, len(oldList))
	for _, raw := range oldList {
		m, _ := raw.(map[string]interface{})
		if m == nil {
			continue
		}
		if rid, _ := m["emr_resource_id"].(string); rid != "" {
			oldByID[rid] = raw
		}
	}

	// Walk new list.  For entries that carry a valid emr_resource_id that
	// matches an old node, pair them.  Entries with no id are new nodes.
	for _, raw := range newList {
		m, _ := raw.(map[string]interface{})
		if m == nil {
			continue
		}
		rid, _ := m["emr_resource_id"].(string)
		if rid != "" {
			if oldSpec, ok := oldByID[rid]; ok {
				alignedOld = append(alignedOld, oldSpec)
				alignedNew = append(alignedNew, raw)
				delete(oldByID, rid)
				continue
			}
		}
		// No matching old node — this is a new node.
		addedNew = append(addedNew, raw)
	}
	return
}
