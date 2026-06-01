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
		// CustomizeDiff enforces business immutability rules at plan time;
		// see customizeDiffEmrClusterV2 for the full rule set.
		CustomizeDiff: customizeDiffEmrClusterV2,
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
				Description: "Cluster scenario configuration. Components to deploy are now declared per node role via the `soft_ware` field inside each `*_resource_spec` block; the cluster-level `software` list below is computed (read-back from API) for observability.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"scene_name": {
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
							Description: "Scenario name, e.g., `Hadoop-Default`, `Hadoop-Kudu`, `Hadoop-Zookeeper`, `Hadoop-Presto`, `Hadoop-Hbase`.",
						},
						"software": {
							Type:        schema.TypeList,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "Deduped list of components actually deployed on the cluster (read-back from API).",
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
										Type:     schema.TypeSet,
										Optional: true,
										Set:      hashEmrNodeResourceSpec,
										Description: "Master node resource specifications. Number of blocks = `MasterCount`. " +
											"All blocks must have identical configuration; the first block is the single resource template sent to the API. " +
											"This field is a `TypeSet` keyed by `_node_index` only  block order in HCL is irrelevant.",
										Elem: emrNodeSpecElem(),
									},
									"core_resource_spec": {
										Type:     schema.TypeSet,
										Optional: true,
										Set:      hashEmrNodeResourceSpec,
										Description: "Core node resource specifications. Number of blocks = `CoreCount`. " +
											"All blocks must have identical configuration; the first block is the single resource template sent to the API. " +
											"This field is a `TypeSet` keyed by `_node_index` only  block order in HCL is irrelevant.",
										Elem: emrNodeSpecElem(),
									},
									"task_resource_spec": {
										Type:     schema.TypeSet,
										Optional: true,
										Set:      hashEmrNodeResourceSpec,
										Description: "Task node resource specifications. Number of blocks = `TaskCount`. " +
											"All blocks must have identical configuration; the first block is the single resource template sent to the API. " +
											"This field is a `TypeSet` keyed by `_node_index` only  block order in HCL is irrelevant.",
										Elem: emrNodeSpecElem(),
									},
									"common_resource_spec": {
										Type:     schema.TypeSet,
										Optional: true,
										Set:      hashEmrNodeResourceSpec,
										Description: "Common node resource specifications. Number of blocks = `CommonCount`. " +
											"All blocks must have identical configuration; the first block is the single resource template sent to the API. " +
											"This field is a `TypeSet` keyed by `_node_index` only  block order in HCL is irrelevant.",
										Elem: emrNodeSpecElem(),
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
				rawList := emrNodeSetToList(allMap[role])
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
				// its `*_resource_spec` set. The first block is taken as the
				// single resource template accepted by the SDK; all blocks must
				// have identical configuration (validated above).
				if rawList := emrNodeSetToList(allMap["master_resource_spec"]); len(rawList) > 0 {
					allSpec.MasterCount = helper.IntInt64(len(rawList))
					if firstMap, ok := rawList[0].(map[string]interface{}); ok {
						allSpec.MasterResourceSpec = buildEmrClusterV2NodeResourceSpec(firstMap)
					}
				}
				if rawList := emrNodeSetToList(allMap["core_resource_spec"]); len(rawList) > 0 {
					allSpec.CoreCount = helper.IntInt64(len(rawList))
					if firstMap, ok := rawList[0].(map[string]interface{}); ok {
						allSpec.CoreResourceSpec = buildEmrClusterV2NodeResourceSpec(firstMap)
					}
				}
				if rawList := emrNodeSetToList(allMap["task_resource_spec"]); len(rawList) > 0 {
					allSpec.TaskCount = helper.IntInt64(len(rawList))
					if firstMap, ok := rawList[0].(map[string]interface{}); ok {
						allSpec.TaskResourceSpec = buildEmrClusterV2NodeResourceSpec(firstMap)
					}
				}
				if rawList := emrNodeSetToList(allMap["common_resource_spec"]); len(rawList) > 0 {
					allSpec.CommonCount = helper.IntInt64(len(rawList))
					if firstMap, ok := rawList[0].(map[string]interface{}); ok {
						allSpec.CommonResourceSpec = buildEmrClusterV2NodeResourceSpec(firstMap)
					}
				}
				zrc.AllNodeResourceSpec = &allSpec
			}

			request.ZoneResourceConfiguration = append(request.ZoneResourceConfiguration, &zrc)
		}

		// Aggregate `soft_ware[*].services` across all roles of all zones
		// (deduped, order-preserving) and feed it to SceneSoftwareConfig.Software.
		// This is the only path that supplies the cluster-wide component list
		// to CreateCluster — the user-facing `software` field is Computed.
		// At least one component must be declared somewhere; otherwise we
		// would send an empty Software list and EMR would reject the create
		// with an unhelpful error.
		serviceList := emrCollectServiceNames(zrcList)
		if len(serviceList) == 0 {
			return fmt.Errorf("at least one `soft_ware` block must be declared on a node role: aggregated component list (SceneSoftwareConfig.Software) is empty, EMR cluster cannot be created without components")
		}
		if request.SceneSoftwareConfig == nil {
			request.SceneSoftwareConfig = &emr.SceneSoftwareConfig{}
		}

		tmpList := []*string{}
		for _, item := range serviceList {
			if item != nil {
				if strings.HasPrefix(*item, "RUNTIME") || strings.HasPrefix(*item, "FILEBEAT") {
					continue
				}
			}

			tmpList = append(tmpList, item)
		}

		request.SceneSoftwareConfig.Software = tmpList
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
		// Round-trip scene_name and the cluster-wide deployed software list
		// (Computed). soft_ware (per-role) is intentionally NOT touched here —
		// the API does not return the role/process breakdown.
		if existing, ok := d.GetOk("scene_software_config"); ok {
			list, _ := existing.([]interface{})
			if len(list) > 0 {
				if item, ok := list[0].(map[string]interface{}); ok {
					if cluster.Config != nil && cluster.Config.SoftInfo != nil {
						item["software"] = cluster.Config.SoftInfo
					}
					item["scene_name"] = *cluster.SceneName
					list[0] = item
					_ = d.Set("scene_software_config", list)
				}
			}
		}
	}

	if cluster.Config != nil {
		if cluster.Config.SecurityGroups != nil {
			_ = d.Set("security_group_ids", cluster.Config.SecurityGroups)
		}

		if cluster.Config.CbsEncrypt != nil {
			_ = d.Set("enable_cbs_encrypt_flag", *cluster.Config.CbsEncrypt == 1)
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

				// Build state-side maps so plan-refresh Read keeps the user's
				// `_node_index` / `_disk_index` values stable across reads:
				//   - stateNodeByRID[roleKey][emr_resource_id] = _node_index
				//   - stateDiskByID[roleKey][emr_resource_id][disk_id] = _disk_index
				//   - stateSoftWareByRID[roleKey][emr_resource_id] = soft_ware
				//     (preserved verbatim because the EMR API does not return
				//     the per-role services/roles breakdown).
				stateNodeByRID := map[string]map[string]string{}
				stateDiskByID := map[string]map[string]map[string]string{}
				stateSoftWareByRID := map[string]map[string][]interface{}{}
				if oldZrcRaw, _ := d.GetChange("zone_resource_configuration"); oldZrcRaw != nil {
					list, _ := oldZrcRaw.([]interface{})
					for _, zrc := range list {
						zrcMap, _ := zrc.(map[string]interface{})
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
							rawList := emrNodeSetToList(allMap[rk])
							nodeMap := map[string]string{}
							diskMap := map[string]map[string]string{}
							softWareMap := map[string][]interface{}{}
							for _, raw := range rawList {
								m, ok := raw.(map[string]interface{})
								if !ok {
									continue
								}
								nodeIdxStr, _ := m["_node_index"].(string)
								rid, _ := m["emr_resource_id"].(string)
								if rid != "" && nodeIdxStr != "" {
									nodeMap[rid] = nodeIdxStr
								}
								perDisk := map[string]string{}
								for _, dRaw := range emrDiskSetToList(m["data_disk"]) {
									dm, ok := dRaw.(map[string]interface{})
									if !ok {
										continue
									}
									dID, _ := dm["disk_id"].(string)
									diskIdxStr, _ := dm["_disk_index"].(string)
									if dID != "" && diskIdxStr != "" {
										perDisk[dID] = diskIdxStr
									}
								}
								if rid != "" {
									diskMap[rid] = perDisk
								}
								if rid != "" {
									softWareMap[rid] = emrNodeSetToList(m["software"])
								}
							}
							stateNodeByRID[rk] = nodeMap
							stateDiskByID[rk] = diskMap
							stateSoftWareByRID[rk] = softWareMap
						}
						break
					}
				}

				// Build a fallback pool of `_node_index` / `_disk_index` values
				// from the current config so the FIRST Read after Create can
				// stamp newly-Read nodes with user-supplied identity strings
				// (state has no rid→_node_index mapping yet). For subsequent
				// plan-refresh Reads the pool simply mirrors state, so this
				// fallback is a no-op there.
				type cfgDiskEntry struct {
					diskIndex string
					diskSize  int
					diskType  string
				}
				configNodeIndexes := map[string][]string{}
				configDisksByNodeIndex := map[string]map[string][]cfgDiskEntry{}
				configSoftWareByNodeIndex := map[string]map[string][]interface{}{}
				if cfgRaw, ok := d.GetOk("zone_resource_configuration"); ok {
					list, _ := cfgRaw.([]interface{})
					for _, zrc := range list {
						zrcMap, _ := zrc.(map[string]interface{})
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
							rawList := emrNodeSetToList(allMap[rk])
							nodeOrd := make([]string, 0, len(rawList))
							perNodeMap := make(map[string][]cfgDiskEntry, len(rawList))
							perNodeSoftWare := make(map[string][]interface{}, len(rawList))
							for _, raw := range rawList {
								m, ok := raw.(map[string]interface{})
								if !ok {
									continue
								}
								nodeIdxStr, _ := m["_node_index"].(string)
								if nodeIdxStr != "" {
									nodeOrd = append(nodeOrd, nodeIdxStr)
								}
								perDisk := []cfgDiskEntry{}
								for _, diskRaw := range emrDiskSetToList(m["data_disk"]) {
									dm, ok := diskRaw.(map[string]interface{})
									if !ok {
										continue
									}
									diskIdxStr, _ := dm["_disk_index"].(string)
									dSz, _ := dm["disk_size"].(int)
									dDt, _ := dm["disk_type"].(string)
									perDisk = append(perDisk, cfgDiskEntry{diskIndex: diskIdxStr, diskSize: dSz, diskType: dDt})
								}
								if nodeIdxStr != "" {
									perNodeMap[nodeIdxStr] = perDisk
									perNodeSoftWare[nodeIdxStr] = emrNodeSetToList(m["software"])
								}
							}
							configNodeIndexes[rk] = nodeOrd
							configDisksByNodeIndex[rk] = perNodeMap
							configSoftWareByNodeIndex[rk] = perNodeSoftWare
						}
						break
					}
				}

				allSpec := make(map[string]interface{})
				for flag, roleKey := range roleKeyMap {
					roleNodes := grouped[flag]

					// Sort API nodes by emr_resource_id so first-Read fallback
					// (when state has no rid mapping yet) consumes the user's
					// `_node_index` pool deterministically.
					sort.Slice(roleNodes, func(i, j int) bool {
						ri, rj := "", ""
						if roleNodes[i].EmrResourceId != nil {
							ri = *roleNodes[i].EmrResourceId
						}
						if roleNodes[j].EmrResourceId != nil {
							rj = *roleNodes[j].EmrResourceId
						}
						return ri < rj
					})

					// Resolve `_node_index` per API node, two-pass:
					//   1. Existing state mapping (rid → _node_index) wins.
					//   2. Unmapped nodes consume unused names from the
					//      user's config pool (first-Read-after-Create path).
					userNodeIdxByPos := make([]string, len(roleNodes))
					{
						cfgList := configNodeIndexes[roleKey]
						cfgConsumed := make(map[string]bool, len(cfgList))
						unmapped := make([]int, 0, len(roleNodes))
						for i, n := range roleNodes {
							rid := ""
							if n.EmrResourceId != nil {
								rid = *n.EmrResourceId
							}
							if rid != "" {
								if t, ok := stateNodeByRID[roleKey][rid]; ok && t != "" {
									userNodeIdxByPos[i] = t
									cfgConsumed[t] = true
									continue
								}
							}
							unmapped = append(unmapped, i)
						}
						cfgFree := make([]string, 0, len(cfgList))
						for _, name := range cfgList {
							if !cfgConsumed[name] {
								cfgFree = append(cfgFree, name)
							}
						}
						for k, i := range unmapped {
							if k < len(cfgFree) {
								userNodeIdxByPos[i] = cfgFree[k]
							}
							// else: API has more nodes than config declares
							// (drift). Leave `_node_index` blank — TypeSet
							// hash collision will surface the drift on next plan.
						}
					}

					specs := make([]interface{}, 0, len(roleNodes))
					for nodeIdx, n := range roleNodes {
						if n.OrderNo == nil {
							continue
						}
						specItem := make(map[string]interface{})
						specItem["instance_type"] = emrInstanceTypeFromNode(n)
						rid := ""
						if n.EmrResourceId != nil {
							rid = *n.EmrResourceId
							specItem["emr_resource_id"] = rid
						} else {
							specItem["emr_resource_id"] = ""
						}
						specItem["order_no"] = *n.OrderNo

						userNodeIdx := userNodeIdxByPos[nodeIdx]
						specItem["_node_index"] = userNodeIdx

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

						dataResp, derr := service.DescribeEmrNodeDataDisks(ctx, instanceId, *n.OrderNo)
						if derr != nil {
							return derr
						}

						apiDisks := make([]map[string]interface{}, 0, len(dataResp))
						for _, item := range dataResp {
							tmpObj := map[string]interface{}{"disk_size": 0, "disk_type": "", "disk_id": ""}
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
						// Sort API disks by disk_id so the (size, type)
						// fallback below is deterministic.
						sort.Slice(apiDisks, func(i, j int) bool {
							di, _ := apiDisks[i]["disk_id"].(string)
							dj, _ := apiDisks[j]["disk_id"].(string)
							return di < dj
						})

						// Resolve `_disk_index` per API disk:
						//   1. Existing state mapping by disk_id wins.
						//   2. Fallback: match unconsumed config disk by
						//      (disk_size, disk_type) — covers brand-new
						//      disks attached during this Apply.
						// Look up per-node config by `_node_index` (already
						// resolved as userNodeIdx), NOT by positional nodeIdx.
						var cfgDiskList []cfgDiskEntry
						if userNodeIdx != "" {
							if perNodeMap := configDisksByNodeIndex[roleKey]; perNodeMap != nil {
								cfgDiskList = perNodeMap[userNodeIdx]
							}
						}
						cfgUsed := make([]bool, len(cfgDiskList))
						cfgIndexByName := make(map[string][]int, len(cfgDiskList))
						for ci, ce := range cfgDiskList {
							cfgIndexByName[ce.diskIndex] = append(cfgIndexByName[ce.diskIndex], ci)
						}
						resolvedIdx := make([]string, len(apiDisks))
						// Pass 1: state-mapped disks claim their `_disk_index`
						// slot first so Pass 2 cannot steal those names.
						for i, dm := range apiDisks {
							diskID, _ := dm["disk_id"].(string)
							if rid == "" || diskID == "" {
								continue
							}
							t, ok := stateDiskByID[roleKey][rid][diskID]
							if !ok || t == "" {
								continue
							}
							resolvedIdx[i] = t
							for _, ci := range cfgIndexByName[t] {
								if !cfgUsed[ci] {
									cfgUsed[ci] = true
									break
								}
							}
						}
						// Pass 2: brand-new disks (no state mapping) match
						// unconsumed cfgDiskList entries by (disk_size, disk_type).
						for i, dm := range apiDisks {
							if resolvedIdx[i] != "" {
								continue
							}
							apiSz, _ := dm["disk_size"].(int)
							apiDt, _ := dm["disk_type"].(string)
							for ci, ce := range cfgDiskList {
								if cfgUsed[ci] {
									continue
								}
								if ce.diskSize == apiSz && ce.diskType == apiDt {
									resolvedIdx[i] = ce.diskIndex
									cfgUsed[ci] = true
									break
								}
							}
						}
						for i, dm := range apiDisks {
							dm["_disk_index"] = resolvedIdx[i]
						}

						diskList := make([]interface{}, len(apiDisks))
						for i, d := range apiDisks {
							diskList[i] = d
						}
						specItem["data_disk"] = diskList

						// soft_ware: the EMR API does not return per-role
						// services/roles, so we restore it from previous
						// state (rid → soft_ware). For the FIRST Read after
						// Create the state mapping is empty, so we fall back
						// to the user's config pool indexed by `_node_index`.
						// Without this, the Set diff would treat every
						// `soft_ware` block as drift and force-replace the
						// whole cluster on the next plan.
						var softWare []interface{}
						if rid != "" {
							if m := stateSoftWareByRID[roleKey]; m != nil {
								softWare = m[rid]
							}
						}
						if len(softWare) == 0 && userNodeIdx != "" {
							if m := configSoftWareByNodeIndex[roleKey]; m != nil {
								softWare = m[userNodeIdx]
							}
						}
						specItem["software"] = softWare

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
				pairedOldMaster, pairedNewMaster, addedMaster, removedMaster := alignNodeListByNodeIndex(oldMasterList, newMasterList)
				if len(addedMaster) > 0 || len(removedMaster) > 0 {
					return fmt.Errorf("zone_resource_configuration[%d].all_node_resource_spec.master_resource_spec does not support scaling (add/remove nodes)", zoneIdx)
				}

				for nodeIdx := 0; nodeIdx < len(pairedNewMaster); nodeIdx++ {
					oldSpec, _ := pairedOldMaster[nodeIdx].(map[string]interface{})
					newSpec, _ := pairedNewMaster[nodeIdx].(map[string]interface{})

					// Skip nodes that were never actually provisioned (order_no is empty).
					if orderNo, _ := oldSpec["order_no"].(string); orderNo == "" {
						continue
					}

					// Prevent modifying system_disk of existing nodes (not supported by API).
					if rid, _ := oldSpec["emr_resource_id"].(string); rid != "" {
						if sysDiskChanged(oldSpec, newSpec) {
							return fmt.Errorf("modifying system_disk of an existing node (emr_resource_id=%s) is not supported", rid)
						}
					}
					if instanceTypeChanged(oldSpec, newSpec) {
						if err := emrModifyNodeInstanceType(ctx, meta, d, logId, instanceId,
							"master_resource_spec", zoneIdx, nodeIdx, oldSpec, newSpec); err != nil {
							return err
						}
					}

					if dataDiskChanged(oldSpec, newSpec) {
						orderNo, _ := oldSpec["order_no"].(string)
						oldDisks := emrDiskSetToList(oldSpec["data_disk"])
						newDisks := emrDiskSetToList(newSpec["data_disk"])

						if err := handleNodeDataDiskChange(ctx, meta, d, logId, instanceId, orderNo,
							oldDisks, newDisks, "master_resource_spec", zoneIdx, nodeIdx); err != nil {
							return err
						}
					}
				}
			}

			// ---------- core_resource_spec ----------
			if oldCoreList, newCoreList, changed := nodeRoleChanged(oldAll, newAll, "core_resource_spec"); changed {
				pairedOldCore, pairedNewCore, addedCore, removedOldCore := alignNodeListByNodeIndex(oldCoreList, newCoreList)

				// Content changes on existing nodes.
				for nodeIdx := 0; nodeIdx < len(pairedNewCore); nodeIdx++ {
					oldSpec, _ := pairedOldCore[nodeIdx].(map[string]interface{})
					newSpec, _ := pairedNewCore[nodeIdx].(map[string]interface{})

					// Skip nodes that were never actually provisioned.
					if orderNo, _ := oldSpec["order_no"].(string); orderNo == "" {
						continue
					}

					if rid, _ := oldSpec["emr_resource_id"].(string); rid != "" {
						if sysDiskChanged(oldSpec, newSpec) {
							return fmt.Errorf("modifying system_disk of an existing node (emr_resource_id=%s) is not supported", rid)
						}
					}
					if instanceTypeChanged(oldSpec, newSpec) {
						if err := emrModifyNodeInstanceType(ctx, meta, d, logId, instanceId,
							"core_resource_spec", zoneIdx, nodeIdx, oldSpec, newSpec); err != nil {
							return err
						}
					}
					if dataDiskChanged(oldSpec, newSpec) {
						orderNo, _ := oldSpec["order_no"].(string)
						oldDisks := emrDiskSetToList(oldSpec["data_disk"])
						newDisks := emrDiskSetToList(newSpec["data_disk"])
						if err := handleNodeDataDiskChange(ctx, meta, d, logId, instanceId, orderNo, oldDisks, newDisks, "core_resource_spec", zoneIdx, nodeIdx); err != nil {
							return err
						}
					}
				}

				// Scale-in first: terminate old nodes that have no matching new
				// entry. Doing scale-in before scale-out frees capacity before
				// new nodes are added.
				var removedCoreOrderNos []*string
				for _, raw := range removedOldCore {
					m, _ := raw.(map[string]interface{})
					orderNo, _ := m["order_no"].(string)
					if orderNo != "" {
						removedCoreOrderNos = append(removedCoreOrderNos, helper.String(orderNo))
					}
				}
				if err := emrTerminateNodes(ctx, meta, d, logId, instanceId, "CORE", zoneIdx, removedCoreOrderNos); err != nil {
					return err
				}

				// Scale-out: new nodes.
				zone := emrZoneOf(newZrcMap)
				for addedIdx, addedRaw := range addedCore {
					addedSpec, _ := addedRaw.(map[string]interface{})
					if addedSpec == nil {
						continue
					}
					if err := emrScaleOutSingleNode(ctx, meta, d, logId, instanceId,
						"CORE", zone, zoneIdx, addedIdx, addedSpec); err != nil {
						return err
					}
				}
			}

			// ---------- task_resource_spec ----------
			if oldTaskList, newTaskList, changed := nodeRoleChanged(oldAll, newAll, "task_resource_spec"); changed {
				pairedOldTask, pairedNewTask, addedTask, removedOldTask := alignNodeListByNodeIndex(oldTaskList, newTaskList)

				// Content changes on existing nodes.
				for nodeIdx := 0; nodeIdx < len(pairedNewTask); nodeIdx++ {
					oldSpec, _ := pairedOldTask[nodeIdx].(map[string]interface{})
					newSpec, _ := pairedNewTask[nodeIdx].(map[string]interface{})

					// Skip nodes that were never actually provisioned.
					if orderNo, _ := oldSpec["order_no"].(string); orderNo == "" {
						continue
					}

					if rid, _ := oldSpec["emr_resource_id"].(string); rid != "" {
						if sysDiskChanged(oldSpec, newSpec) {
							return fmt.Errorf("modifying system_disk of an existing node (emr_resource_id=%s) is not supported", rid)
						}
					}
					if instanceTypeChanged(oldSpec, newSpec) {
						if err := emrModifyNodeInstanceType(ctx, meta, d, logId, instanceId,
							"task_resource_spec", zoneIdx, nodeIdx, oldSpec, newSpec); err != nil {
							return err
						}
					}
					if dataDiskChanged(oldSpec, newSpec) {
						orderNo, _ := oldSpec["order_no"].(string)
						oldDisks := emrDiskSetToList(oldSpec["data_disk"])
						newDisks := emrDiskSetToList(newSpec["data_disk"])
						if err := handleNodeDataDiskChange(ctx, meta, d, logId, instanceId, orderNo, oldDisks, newDisks, "task_resource_spec", zoneIdx, nodeIdx); err != nil {
							return err
						}
					}
				}

				// Scale-in first: terminate old task nodes that have no matching
				// new entry. Doing scale-in before scale-out frees capacity
				// before new nodes are added.
				var removedTaskOrderNos []*string
				for _, raw := range removedOldTask {
					m, _ := raw.(map[string]interface{})
					orderNo, _ := m["order_no"].(string)
					if orderNo != "" {
						removedTaskOrderNos = append(removedTaskOrderNos, helper.String(orderNo))
					}
				}
				if err := emrTerminateNodes(ctx, meta, d, logId, instanceId, "TASK", zoneIdx, removedTaskOrderNos); err != nil {
					return err
				}

				// Scale-out: new nodes.
				zone := emrZoneOf(newZrcMap)
				for addedIdx, addedRaw := range addedTask {
					addedSpec, _ := addedRaw.(map[string]interface{})
					if addedSpec == nil {
						continue
					}
					if err := emrScaleOutSingleNode(ctx, meta, d, logId, instanceId,
						"TASK", zone, zoneIdx, addedIdx, addedSpec); err != nil {
						return err
					}
				}
			}

			// ---------- common_resource_spec ----------
			if oldCommonList, newCommonList, changed := nodeRoleChanged(oldAll, newAll, "common_resource_spec"); changed {
				// common_resource_spec does not support scaling (add or remove nodes).
				pairedOldCommon, pairedNewCommon, addedCommon, removedCommon := alignNodeListByNodeIndex(oldCommonList, newCommonList)
				if len(addedCommon) > 0 || len(removedCommon) > 0 {
					return fmt.Errorf("zone_resource_configuration[%d].all_node_resource_spec.common_resource_spec does not support scaling (add/remove nodes)", zoneIdx)
				}

				for nodeIdx := 0; nodeIdx < len(pairedNewCommon); nodeIdx++ {
					oldSpec, _ := pairedOldCommon[nodeIdx].(map[string]interface{})
					newSpec, _ := pairedNewCommon[nodeIdx].(map[string]interface{})

					// Skip nodes that were never actually provisioned.
					if orderNo, _ := oldSpec["order_no"].(string); orderNo == "" {
						continue
					}

					if rid, _ := oldSpec["emr_resource_id"].(string); rid != "" {
						if sysDiskChanged(oldSpec, newSpec) {
							return fmt.Errorf("modifying system_disk of an existing node (emr_resource_id=%s) is not supported", rid)
						}
					}
					if instanceTypeChanged(oldSpec, newSpec) {
						if err := emrModifyNodeInstanceType(ctx, meta, d, logId, instanceId,
							"common_resource_spec", zoneIdx, nodeIdx, oldSpec, newSpec); err != nil {
							return err
						}
					}
					if dataDiskChanged(oldSpec, newSpec) {
						orderNo, _ := oldSpec["order_no"].(string)
						oldDisks := emrDiskSetToList(oldSpec["data_disk"])
						newDisks := emrDiskSetToList(newSpec["data_disk"])
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

// customizeDiffEmrClusterV2 enforces business immutability rules at plan time
// for the four `*_resource_spec` TypeSet fields and the inner `data_disk`
// TypeSet. Both layers are Sets keyed via their full content hash, so the
// SDK's set-diff handles add/remove/in-place-modify automatically; this hook
// performs validation only — it never mutates `new` and never overwrites
// Computed fields.
//
// Rules enforced:
//   - master/common: list length is fixed; rename of `_node_index` is
//     rejected (length-immutable roles cannot tolerate add+remove).
//   - core/task: scale-out and scale-in are both allowed, including
//     simultaneous add+remove in a single apply (legitimate node
//     replacement). Update processes scale-in before scale-out.
//   - For nodes paired by `_node_index`:
//   - `system_disk` is immutable in every dimension.
//   - `data_disk` allows attach (new `_disk_index`) but never detach.
//   - For paired disks: `disk_size` only grows; `disk_type` is immutable.
//   - `_node_index` must be non-empty and unique within each role.
//   - `_disk_index` must be non-empty and unique within each node.
//
// Skipped during create (`d.Id() == ""`): no prior state to compare against,
// and Create has its own validation path.
func customizeDiffEmrClusterV2(ctx context.Context, d *schema.ResourceDiff, meta interface{}) error {
	if d.Id() == "" {
		return nil
	}
	oldRaw, newRaw := d.GetChange("zone_resource_configuration")
	oldZones, _ := oldRaw.([]interface{})
	newZones, _ := newRaw.([]interface{})
	if len(newZones) == 0 {
		return nil
	}

	// Pair zones by user-supplied placement.0.zone (zones are user-keyed).
	oldZoneByName := make(map[string]map[string]interface{}, len(oldZones))
	for _, z := range oldZones {
		zm, _ := z.(map[string]interface{})
		if zm == nil {
			continue
		}
		oldZoneByName[zoneNameOfBlock(zm)] = zm
	}

	immutableLengthRoles := map[string]bool{
		"master_resource_spec": true,
		"common_resource_spec": true,
	}
	roleKeys := []string{
		"master_resource_spec",
		"core_resource_spec",
		"task_resource_spec",
		"common_resource_spec",
	}
	roleLabel := map[string]string{
		"master_resource_spec": "master",
		"core_resource_spec":   "core",
		"task_resource_spec":   "task",
		"common_resource_spec": "common",
	}

	for zoneIdx, nz := range newZones {
		newZone, _ := nz.(map[string]interface{})
		if newZone == nil {
			continue
		}
		zoneName := zoneNameOfBlock(newZone)
		oldZone := oldZoneByName[zoneName]

		newAll := emrFirstAllNodeResourceSpec(newZone)
		if newAll == nil {
			continue
		}
		var oldAll map[string]interface{}
		if oldZone != nil {
			oldAll = emrFirstAllNodeResourceSpec(oldZone)
		}

		for _, rk := range roleKeys {
			newList := emrNodeSetToList(newAll[rk])

			// Validate _node_index: non-empty + unique within this role.
			seenNodeIdx := make(map[string]bool, len(newList))
			for _, raw := range newList {
				nm, _ := raw.(map[string]interface{})
				if nm == nil {
					continue
				}
				nodeIdxStr, _ := nm["_node_index"].(string)
				if nodeIdxStr == "" {
					return fmt.Errorf(
						"%s.zone_resource_configuration[%d].%s: every block must set a non-empty `_node_index`",
						zoneName, zoneIdx, rk)
				}
				if seenNodeIdx[nodeIdxStr] {
					return fmt.Errorf(
						"%s.zone_resource_configuration[%d].%s: duplicate `_node_index` %q within the same role",
						zoneName, zoneIdx, rk, nodeIdxStr)
				}
				seenNodeIdx[nodeIdxStr] = true
			}

			// Validate _disk_index: non-empty + unique within each node.
			for _, raw := range newList {
				nm, _ := raw.(map[string]interface{})
				if nm == nil {
					continue
				}
				nodeIdxStr, _ := nm["_node_index"].(string)
				newDisks := emrDiskSetToList(nm["data_disk"])
				seenDiskIdx := make(map[string]bool, len(newDisks))
				for _, dRaw := range newDisks {
					dm, _ := dRaw.(map[string]interface{})
					if dm == nil {
						continue
					}
					diskIdxStr, _ := dm["_disk_index"].(string)
					if diskIdxStr == "" {
						return fmt.Errorf(
							"%s.zone_resource_configuration[%d].%s[_node_index=%s].data_disk: every block must set a non-empty `_disk_index`",
							zoneName, zoneIdx, rk, nodeIdxStr)
					}
					if seenDiskIdx[diskIdxStr] {
						return fmt.Errorf(
							"%s.zone_resource_configuration[%d].%s[_node_index=%s].data_disk: duplicate `_disk_index` %q within the same node",
							zoneName, zoneIdx, rk, nodeIdxStr, diskIdxStr)
					}
					seenDiskIdx[diskIdxStr] = true
				}
			}

			// If the corresponding old zone/role does not exist (brand-new
			// zone), skip cross-old/new comparisons — there is nothing to
			// compare against, and Create-time validation has already run.
			if oldAll == nil {
				continue
			}
			oldList := emrNodeSetToList(oldAll[rk])

			// Build _node_index → block maps for both sides.
			oldByIdx := indexListBy(oldList, "_node_index")
			newByIdx := indexListBy(newList, "_node_index")

			// Compute set diff at the `_node_index` granularity.
			var addedNodes, removedNodes []string
			for idx := range newByIdx {
				if _, ok := oldByIdx[idx]; !ok {
					addedNodes = append(addedNodes, idx)
				}
			}
			for idx := range oldByIdx {
				if _, ok := newByIdx[idx]; !ok {
					removedNodes = append(removedNodes, idx)
				}
			}

			// Length-immutable roles: any add or remove is forbidden.
			if immutableLengthRoles[rk] {
				if len(addedNodes) > 0 || len(removedNodes) > 0 {
					return fmt.Errorf(
						"%s.zone_resource_configuration[%d].%s: %s list length is immutable and `_node_index` cannot be renamed (added=%v, removed=%v, old_count=%d, new_count=%d)",
						zoneName, zoneIdx, rk, roleLabel[rk],
						addedNodes, removedNodes, len(oldList), len(newList))
				}
			}
			// core/task: simultaneous add+remove is allowed (legitimate node
			// replacement). Update processes scale-in before scale-out.

			// For nodes present in both old and new (paired by _node_index):
			// validate system_disk immutability and data_disk constraints.
			for nodeIdxStr, newRawNode := range newByIdx {
				oldRawNode, paired := oldByIdx[nodeIdxStr]
				if !paired {
					continue
				}
				newNode, _ := newRawNode.(map[string]interface{})
				oldNode, _ := oldRawNode.(map[string]interface{})
				if newNode == nil || oldNode == nil {
					continue
				}

				// system_disk immutable.
				if err := emrAssertSystemDiskImmutable(
					oldNode["system_disk"], newNode["system_disk"],
					zoneName, zoneIdx, rk, nodeIdxStr,
				); err != nil {
					return err
				}

				// software immutable: any add/remove/content change is rejected.
				if err := emrSoftWareSetsEqual(oldNode["software"], newNode["software"]); err != nil {
					return fmt.Errorf(
						"%s.zone_resource_configuration[%d].%s[_node_index=%s].software: software is immutable after create: %v",
						zoneName, zoneIdx, rk, nodeIdxStr, err)
				}

				// data_disk validation: no detach, no shrink, no type change,
				// no rename (rename surfaces as remove+add at disk level).
				oldDisks := emrDiskSetToList(oldNode["data_disk"])
				newDisks := emrDiskSetToList(newNode["data_disk"])
				oldDiskByIdx := indexListBy(oldDisks, "_disk_index")
				newDiskByIdx := indexListBy(newDisks, "_disk_index")

				var addedDisks, removedDisks []string
				for idx := range newDiskByIdx {
					if _, ok := oldDiskByIdx[idx]; !ok {
						addedDisks = append(addedDisks, idx)
					}
				}
				for idx := range oldDiskByIdx {
					if _, ok := newDiskByIdx[idx]; !ok {
						removedDisks = append(removedDisks, idx)
					}
				}
				if len(removedDisks) > 0 {
					return fmt.Errorf(
						"%s.zone_resource_configuration[%d].%s[_node_index=%s].data_disk: removing data_disk is not allowed (_disk_index=%v removed from config). To grow capacity, add a new `_disk_index` instead of removing existing ones",
						zoneName, zoneIdx, rk, nodeIdxStr, removedDisks)
				}
				_ = addedDisks // attaching new disks is allowed; nothing to validate beyond uniqueness.

				// Paired disks: disk_size only grows; disk_type immutable.
				for diskIdxStr, oldRawDisk := range oldDiskByIdx {
					newRawDisk, ok := newDiskByIdx[diskIdxStr]
					if !ok {
						continue
					}
					oldDisk, _ := oldRawDisk.(map[string]interface{})
					newDisk, _ := newRawDisk.(map[string]interface{})
					if oldDisk == nil || newDisk == nil {
						continue
					}
					oldSz, _ := oldDisk["disk_size"].(int)
					newSz, _ := newDisk["disk_size"].(int)
					if newSz < oldSz {
						return fmt.Errorf(
							"%s.zone_resource_configuration[%d].%s[_node_index=%s].data_disk[_disk_index=%s]: data_disk.disk_size can only grow (old=%d, new=%d)",
							zoneName, zoneIdx, rk, nodeIdxStr, diskIdxStr, oldSz, newSz)
					}
					oldDt, _ := oldDisk["disk_type"].(string)
					newDt, _ := newDisk["disk_type"].(string)
					if oldDt != "" && newDt != "" && oldDt != newDt {
						return fmt.Errorf(
							"%s.zone_resource_configuration[%d].%s[_node_index=%s].data_disk[_disk_index=%s]: data_disk.disk_type is immutable after create (old=%s, new=%s)",
							zoneName, zoneIdx, rk, nodeIdxStr, diskIdxStr, oldDt, newDt)
					}
				}
			}
		}
	}
	return nil
}

// emrAssertSystemDiskImmutable returns a non-nil error if any field of the
// system_disk block (size, type) differs between old and new. Empty old =
// node had no system_disk in state (shouldn't happen for an existing
// cluster); in that case we skip the check and let Read repair state on
// next refresh.
func emrAssertSystemDiskImmutable(oldRaw, newRaw interface{}, zoneName string, zoneIdx int, rk, nodeIdxStr string) error {
	oldL, _ := oldRaw.([]interface{})
	newL, _ := newRaw.([]interface{})
	if len(oldL) == 0 || len(newL) == 0 {
		return nil
	}
	om, _ := oldL[0].(map[string]interface{})
	nm, _ := newL[0].(map[string]interface{})
	if om == nil || nm == nil {
		return nil
	}
	for _, f := range []string{"disk_size", "disk_type"} {
		ov := fmt.Sprintf("%v", om[f])
		nv := fmt.Sprintf("%v", nm[f])
		if ov != nv {
			return fmt.Errorf(
				"%s.zone_resource_configuration[%d].%s[_node_index=%s].system_disk.%s: system_disk is immutable after create (old=%s, new=%s)",
				zoneName, zoneIdx, rk, nodeIdxStr, f, ov, nv)
		}
	}
	return nil
}

// zoneNameOfBlock returns the zone string of a zone_resource_configuration
// block, or "" if missing.
func zoneNameOfBlock(zone map[string]interface{}) string {
	pl, _ := zone["placement"].([]interface{})
	if len(pl) == 0 {
		return ""
	}
	pm, _ := pl[0].(map[string]interface{})
	if pm == nil {
		return ""
	}
	z, _ := pm["zone"].(string)
	return z
}

// emrFirstAllNodeResourceSpec returns the first all_node_resource_spec block
// (it is MaxItems:1) of the given zone, as a writable map; nil if missing.
func emrFirstAllNodeResourceSpec(zone map[string]interface{}) map[string]interface{} {
	all, _ := zone["all_node_resource_spec"].([]interface{})
	if len(all) == 0 {
		return nil
	}
	m, _ := all[0].(map[string]interface{})
	return m
}

// indexListBy builds a lookup map keyed by the string value of `key` field
// (typically "_node_index" or "_disk_index"). Empty-keyed and non-map entries
// are silently skipped.
func indexListBy(list []interface{}, key string) map[string]interface{} {
	out := make(map[string]interface{}, len(list))
	for _, raw := range list {
		m, _ := raw.(map[string]interface{})
		if m == nil {
			continue
		}
		k, _ := m[key].(string)
		if k == "" {
			continue
		}
		out[k] = m
	}
	return out
}

// hashEmrNodeResourceSpec is the Set hash function used by all four
// *_resource_spec TypeSet fields. It hashes the entire user-supplied content
// (skipping Computed-only fields like emr_resource_id / order_no / disk_id
// whose plan-time values are unknown).
//
// We hash full content rather than only `_node_index` because the SDK's
// diffSet short-circuits when the hash multiset is identical: a `_node_index`-
// only hash would mask in-place edits (resize disk, change instance_type, ...)
// from the plan. Full-content hash makes any user-visible change shift the
// element's hash; the Update handler then re-pairs old/new by `_node_index`
// string match so the cloud operation is still in-place — no destroy/recreate.
//
// Trade-off: plan output for an in-place modification visually appears as
// "remove old block + add new block". This is purely cosmetic; the actual
// Update is in-place.
var hashEmrNodeResourceSpec = schema.HashResource(emrNodeSpecElem())

// hashEmrDataDisk is the Set hash function for the inner `data_disk` TypeSet.
// Same rationale as hashEmrNodeResourceSpec: hash the full content so the SDK
// diffSet does not short-circuit on inner field edits (e.g., disk_size).
var hashEmrDataDisk = schema.HashResource(emrDataDiskElem())

// emrDataDiskElem returns the inner `data_disk` Resource schema, used both
// as the TypeSet's Elem and as input to schema.HashResource above.
func emrDataDiskElem() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"disk_size": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Disk size in GB. Can only be increased after creation; shrinking is rejected at plan time.",
			},
			"disk_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Disk type. Valid values: `CLOUD_SSD`, `CLOUD_PREMIUM`, `CLOUD_BASIC`, `LOCAL_BASIC`, `LOCAL_SSD`, `CLOUD_HSSD`, `CLOUD_THROUGHPUT`, `CLOUD_TSSD`, `CLOUD_BIGDATA`, `CLOUD_HIGHIO`, `CLOUD_BSSD`, `REMOTE_SSD`. Immutable after creation.",
			},
			"disk_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Disk ID (read-only, populated from API).",
			},
			"_disk_index": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "**Required** stable identity key for this `data_disk` block. Must be unique within the same node's `data_disk` set, and must remain stable across plan/apply once written to state. Renaming an existing `_disk_index` is rejected (treated as remove+add, which violates the no-shrink rule).",
			},
		},
	}
}

// emrNodeSpecElem returns the schema.Resource Elem shared by all four
// *_resource_spec fields. Those fields are TypeSet keyed by `_node_index`;
// the inner `data_disk` is a TypeSet keyed by `_disk_index`. CustomizeDiff
// enforces all business immutability rules.
func emrNodeSpecElem() *schema.Resource {
	return &schema.Resource{
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
				Type:        schema.TypeSet,
				Optional:    true,
				Set:         hashEmrDataDisk,
				Description: "Cloud data disk specifications. `TypeSet` keyed by full content (including `_disk_index`); block order in HCL is irrelevant.",
				Elem:        emrDataDiskElem(),
			},
			"software": {
				Type:        schema.TypeSet,
				Required:    true,
				MinItems:    1,
				Description: "Per-role software components (with their role/process lists) deployed on this node role. Must be identical across every block of the same role at create time. Aggregated across all four roles (deduped by `services`) and passed to `CreateCluster` as `SceneSoftwareConfig.Software`. Immutable after create modification is rejected at plan time by CustomizeDiff.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"services": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Component name with version, e.g., `hdfs-3.2.2`.",
						},
						"roles": {
							Type:        schema.TypeSet,
							Required:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "Process list for this component on this role, e.g., `[\"NameNode\", \"ZKFailoverController\"]` for hdfs.",
						},
					},
				},
			},
			"emr_resource_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "EMR node resource ID (read-only).",
			},
			"order_no": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Machine instance ID (read-only).",
			},
			"_node_index": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "**Required** stable identity key for this node spec block. Must be unique within the same role's set, and must remain stable across plan/apply. Used by the Update handler to pair old/new blocks for in-place modification. Renaming an existing `_node_index` is rejected by CustomizeDiff for master/common (length immutable); for core/task it is interpreted as scale-in old + scale-out new.",
			},
		},
	}
}

// buildEmrClusterV2NodeResourceSpec converts a single Terraform node resource
// spec block (the first element of a *_resource_spec set) into an
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
	for _, d := range emrDiskSetToList(specMap["data_disk"]) {
		spec.DataDisk = append(spec.DataDisk, buildEmrClusterV2DiskSpecInfo(d))
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

// emrNodeSetToList extracts []interface{} from either a *schema.Set (the
// current TypeSet shape — `(*schema.Set).List()` returns entries in hash
// order) or a []interface{} (kept for forward compatibility), returning nil
// for any other type.
func emrNodeSetToList(raw interface{}) []interface{} {
	if raw == nil {
		return nil
	}
	if s, ok := raw.(*schema.Set); ok {
		return s.List()
	}
	if l, ok := raw.([]interface{}); ok {
		return l
	}
	return nil
}

// emrDiskSetToList extracts []interface{} from a data_disk value, accepting
// both *schema.Set (TypeSet) and []interface{} shapes. Identical in logic to
// emrNodeSetToList; kept separate for naming clarity at call sites.
func emrDiskSetToList(raw interface{}) []interface{} {
	return emrNodeSetToList(raw)
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
	// data_disk is a TypeSet; compare as multisets by (disk_type, disk_size) — order-independent.
	aDataDisks := emrDiskSetToList(a["data_disk"])
	bDataDisks := emrDiskSetToList(b["data_disk"])
	if len(aDataDisks) != len(bDataDisks) {
		return fmt.Errorf("data_disk: block count mismatch (%d vs %d)", len(aDataDisks), len(bDataDisks))
	}
	diskFP := func(raw interface{}) string {
		m, _ := raw.(map[string]interface{})
		return fmt.Sprintf("%v:%v", m["disk_type"], m["disk_size"])
	}
	aCounts := make(map[string]int, len(aDataDisks))
	for _, d := range aDataDisks {
		aCounts[diskFP(d)]++
	}
	for _, d := range bDataDisks {
		fp := diskFP(d)
		aCounts[fp]--
		if aCounts[fp] < 0 {
			return fmt.Errorf("data_disk: disk %q present in b but not matched in a", fp)
		}
	}
	// soft_ware is a TypeSet of {services, roles(TypeSet of strings)}; compare
	// as multisets by full content (services + sorted roles list).
	if err := emrSoftWareSetsEqual(a["software"], b["software"]); err != nil {
		return err
	}
	return nil
}

// emrSoftWareSetsEqual returns a non-nil error if the two `soft_ware` set
// values differ in content. Equality is full-content (`services` + `roles`
// multiset), order-independent.
func emrSoftWareSetsEqual(aRaw, bRaw interface{}) error {
	aList := emrNodeSetToList(aRaw)
	bList := emrNodeSetToList(bRaw)
	if len(aList) != len(bList) {
		return fmt.Errorf("soft_ware: block count mismatch (%d vs %d)", len(aList), len(bList))
	}
	fp := func(raw interface{}) string {
		m, _ := raw.(map[string]interface{})
		if m == nil {
			return ""
		}
		services, _ := m["services"].(string)
		rolesList := emrNodeSetToList(m["roles"])
		roles := make([]string, 0, len(rolesList))
		for _, r := range rolesList {
			if s, ok := r.(string); ok {
				roles = append(roles, s)
			}
		}
		sort.Strings(roles)
		return services + "|" + strings.Join(roles, ",")
	}
	counts := make(map[string]int, len(aList))
	for _, raw := range aList {
		counts[fp(raw)]++
	}
	for _, raw := range bList {
		key := fp(raw)
		counts[key]--
		if counts[key] < 0 {
			return fmt.Errorf("soft_ware: entry %q present in b but not matched in a", key)
		}
	}
	return nil
}

// emrCollectServiceNames walks every zone in the user-provided
// `zone_resource_configuration` list and collects the union of
// `soft_ware[*].services` values across all four roles, deduping while
// preserving first-seen order. Returns the result as a slice of *string ready
// to be assigned to `request.SceneSoftwareConfig.Software`.
func emrCollectServiceNames(zrcList []interface{}) []*string {
	seen := make(map[string]bool)
	var out []*string
	for _, zrc := range zrcList {
		zrcMap, _ := zrc.(map[string]interface{})
		if zrcMap == nil {
			continue
		}
		allList, _ := zrcMap["all_node_resource_spec"].([]interface{})
		if len(allList) == 0 {
			continue
		}
		allMap, _ := allList[0].(map[string]interface{})
		if allMap == nil {
			continue
		}
		for _, role := range []string{"master_resource_spec", "core_resource_spec", "task_resource_spec", "common_resource_spec"} {
			for _, raw := range emrNodeSetToList(allMap[role]) {
				m, _ := raw.(map[string]interface{})
				if m == nil {
					continue
				}
				for _, sw := range emrNodeSetToList(m["software"]) {
					sm, _ := sw.(map[string]interface{})
					if sm == nil {
						continue
					}
					services, _ := sm["services"].(string)
					if services == "" || seen[services] {
						continue
					}
					seen[services] = true
					out = append(out, helper.String(services))
				}
			}
		}
	}
	return out
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

// nodeRoleChanged returns the old/new spec lists for a given node role key
// and whether any difference exists between them.
func nodeRoleChanged(oldAll, newAll map[string]interface{}, roleKey string) (oldList, newList []interface{}, changed bool) {
	oldList = emrNodeSetToList(oldAll[roleKey])
	newList = emrNodeSetToList(newAll[roleKey])
	if len(oldList) == 0 && len(newList) == 0 {
		return oldList, newList, false
	}
	return oldList, newList, fmt.Sprintf("%v", oldList) != fmt.Sprintf("%v", newList)
}

// instanceTypeChanged reports whether the instance_type field differs between
// the old and new node spec maps.
func instanceTypeChanged(oldSpec, newSpec map[string]interface{}) bool {
	oldType, _ := oldSpec["instance_type"].(string)
	newType, _ := newSpec["instance_type"].(string)
	return oldType != newType
}

// dataDiskChanged reports whether data_disk fields differ between old and new node
// spec maps. Comparison is multiset-based on (disk_type, disk_size) so that
// ordering differences in the TypeList do not produce false positives.
func dataDiskChanged(oldSpec, newSpec map[string]interface{}) bool {
	diskFP := func(raw interface{}) string {
		m, _ := raw.(map[string]interface{})
		return fmt.Sprintf("%v:%v", m["disk_type"], m["disk_size"])
	}
	oldDisks := emrDiskSetToList(oldSpec["data_disk"])
	newDisks := emrDiskSetToList(newSpec["data_disk"])
	if len(oldDisks) != len(newDisks) {
		return true
	}
	counts := make(map[string]int, len(oldDisks))
	for _, d := range oldDisks {
		counts[diskFP(d)]++
	}
	for _, d := range newDisks {
		fp := diskFP(d)
		counts[fp]--
		if counts[fp] < 0 {
			return true
		}
	}
	return false
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

// handleNodeDataDiskChange handles data_disk changes for a single node.
//
// oldDisks comes from the paired old node spec (state layer, authoritative).
// newDisks comes from the user config (diff layer).
//
// Matching is done strictly by _disk_index (the user-supplied stable index):
//
//   - Same _disk_index in both old and new:
//
//   - same disk_size  → no-op
//
//   - new disk_size > old disk_size → resize old disk's disk_id to new size
//
//   - new disk_size < old disk_size → error (shrinking not supported)
//
//   - _disk_index only in new → attach as a brand-new disk.
//
//   - _disk_index only in old → error (removing data_disk is not supported).
//
// Because _disk_index is required and remains stable across user reorderings,
// we never need to guess which old disk corresponds to which new entry by
// content fingerprint. The diff-layer disk_id is ignored entirely.
func handleNodeDataDiskChange(
	ctx context.Context,
	meta interface{},
	d *schema.ResourceData,
	logId, instanceId, orderNo string,
	oldDisks, newDisks []interface{},
	roleKey string, zoneIdx, nodeIdx int,
) error {
	type oldDisk struct {
		size     int
		diskID   string
		diskType string
	}
	oldByDiskIndex := make(map[string]oldDisk, len(oldDisks))
	for _, item := range oldDisks {
		m, _ := item.(map[string]interface{})
		if m == nil {
			continue
		}
		diskIdxStr, _ := m["_disk_index"].(string)
		if diskIdxStr == "" {
			continue
		}
		sz, _ := m["disk_size"].(int)
		id, _ := m["disk_id"].(string)
		dt, _ := m["disk_type"].(string)
		oldByDiskIndex[diskIdxStr] = oldDisk{size: sz, diskID: id, diskType: dt}
	}

	type resizeOp struct {
		diskID  string
		newSize int
	}
	type attachOp struct {
		diskType string
		size     int
	}
	var resizes []resizeOp
	var attaches []attachOp

	matchedDiskIndexes := make(map[string]bool, len(oldByDiskIndex))

	for _, item := range newDisks {
		m, _ := item.(map[string]interface{})
		if m == nil {
			continue
		}
		diskIdxStr, _ := m["_disk_index"].(string)
		if diskIdxStr == "" {
			return fmt.Errorf("%s[%d][%d]: data_disk._disk_index must not be empty", roleKey, zoneIdx, nodeIdx)
		}
		newSz, _ := m["disk_size"].(int)
		newDt, _ := m["disk_type"].(string)

		if old, ok := oldByDiskIndex[diskIdxStr]; ok {
			matchedDiskIndexes[diskIdxStr] = true
			if newSz < old.size {
				return fmt.Errorf(
					"%s[%d][%d]: shrinking data_disk is not supported (_disk_index=%s, old=%d, new=%d)",
					roleKey, zoneIdx, nodeIdx, diskIdxStr, old.size, newSz)
			}
			if newSz > old.size {
				resizes = append(resizes, resizeOp{old.diskID, newSz})
			}
			// disk_type change on existing disk is not supported by the API; warn via error.
			if newDt != "" && old.diskType != "" && newDt != old.diskType {
				return fmt.Errorf(
					"%s[%d][%d]: changing data_disk.disk_type of an existing disk is not supported (_disk_index=%s, old=%s, new=%s)",
					roleKey, zoneIdx, nodeIdx, diskIdxStr, old.diskType, newDt)
			}
			continue
		}
		// New _disk_index → attach a new disk.
		attaches = append(attaches, attachOp{newDt, newSz})
	}

	for diskIdxStr, old := range oldByDiskIndex {
		if !matchedDiskIndexes[diskIdxStr] {
			return fmt.Errorf(
				"%s[%d][%d]: removing data_disk is not supported (_disk_index=%s, disk_id=%s)",
				roleKey, zoneIdx, nodeIdx, diskIdxStr, old.diskID)
		}
	}

	// Execute resize operations.
	for _, op := range resizes {
		resizeReq := emr.NewResizeDataDisksRequest()
		resizeResp := emr.NewResizeDataDisksResponse()
		resizeReq.InstanceId = helper.String(instanceId)
		resizeReq.DiskSize = helper.Int64(int64(op.newSize))
		resizeReq.CvmInstanceIds = []*string{helper.String(orderNo)}
		resizeReq.ResizeAll = helper.Bool(false)
		resizeReq.DiskIds = []*string{helper.String(op.diskID)}

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
			log.Printf("[CRITAL]%s %s[%d][%d] data_disk (disk_id=%s) resize to %d failed: %+v",
				logId, roleKey, zoneIdx, nodeIdx, op.diskID, op.newSize, resizeErr)
			return resizeErr
		}
		flowId := *resizeResp.Response.FlowId
		conf := tccommon.BuildStateChangeConf([]string{"0", "1"}, []string{"2", "-1"}, d.Timeout(schema.TimeoutUpdate)-time.Minute, time.Second,
			(&EMRService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}).FlowStatusRefreshFunc(instanceId, strconv.FormatInt(flowId, 10), F_KEY_FLOW_ID, []string{}))
		if object, e := conf.WaitForState(); e != nil {
			return e
		} else if status, ok := object.(*int64); ok && status != nil && *status == -1 {
			return fmt.Errorf("EMR cluster flow failed (FlowId=%s), flow total status is -1.", strconv.FormatInt(flowId, 10))
		}
	}

	// Execute attach (new disk) operations.
	for diskIdx, op := range attaches {
		attachReq := emr.NewAttachDisksRequest()
		attachResp := emr.NewAttachDisksResponse()
		attachReq.InstanceId = helper.String(instanceId)
		attachReq.CvmInstanceIds = []*string{helper.String(orderNo)}
		attachReq.CreateDisk = helper.Bool(true)
		attachReq.DiskSpec = &emr.NodeSpecDiskV2{
			Count:           helper.Int64(1),
			DiskType:        helper.String(op.diskType),
			DefaultDiskSize: helper.Int64(int64(op.size)),
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
			log.Printf("[CRITAL]%s %s[%d][%d] attach new disk[%d] (type=%s size=%d) failed: %+v",
				logId, roleKey, zoneIdx, nodeIdx, diskIdx, op.diskType, op.size, attachErr)
			return attachErr
		}
		flowId := int64(*attachResp.Response.FlowId)
		conf := tccommon.BuildStateChangeConf([]string{"0", "1"}, []string{"2", "-1"}, d.Timeout(schema.TimeoutUpdate)-time.Minute, time.Second,
			(&EMRService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}).FlowStatusRefreshFunc(instanceId, strconv.FormatInt(flowId, 10), F_KEY_FLOW_ID, []string{}))
		if object, e := conf.WaitForState(); e != nil {
			return e
		} else if status, ok := object.(*int64); ok && status != nil && *status == -1 {
			return fmt.Errorf("EMR cluster flow failed (FlowId=%s), flow total status is -1.", strconv.FormatInt(flowId, 10))
		}
	}

	return nil
}

func sysDiskChanged(oldSpec, newSpec map[string]interface{}) bool {
	return fmt.Sprintf("%v", oldSpec["system_disk"]) != fmt.Sprintf("%v", newSpec["system_disk"])
}

// alignNodeListByNodeIndex pairs old (state) and new (config) node lists for
// Update processing using the user-supplied `_node_index` as the stable
// identity key (Required + user-controlled, unlike the Computed
// `emr_resource_id` which can be unreliable across reorderings).
//
// Returns:
//   - pairedOld:  old node specs for existing nodes (in pairedNew order)
//   - pairedNew:  new node specs matched to existing old nodes
//   - addedNew:   new node specs whose `_node_index` is not in oldList
//   - removedOld: old node specs whose `_node_index` is not in newList
func alignNodeListByNodeIndex(
	oldList, newList []interface{},
) (pairedOld, pairedNew []interface{}, addedNew []interface{}, removedOld []interface{}) {
	oldByNodeIndex := make(map[string]interface{}, len(oldList))
	for _, raw := range oldList {
		m, _ := raw.(map[string]interface{})
		if m == nil {
			continue
		}
		if nodeIdxStr, _ := m["_node_index"].(string); nodeIdxStr != "" {
			oldByNodeIndex[nodeIdxStr] = raw
		}
	}

	matchedNodeIndexes := make(map[string]bool, len(oldList))

	for _, raw := range newList {
		m, _ := raw.(map[string]interface{})
		if m == nil {
			continue
		}
		nodeIdxStr, _ := m["_node_index"].(string)
		if nodeIdxStr != "" {
			if oldSpec, ok := oldByNodeIndex[nodeIdxStr]; ok && !matchedNodeIndexes[nodeIdxStr] {
				pairedOld = append(pairedOld, oldSpec)
				pairedNew = append(pairedNew, raw)
				matchedNodeIndexes[nodeIdxStr] = true
				continue
			}
		}
		// _node_index not found in old → new node (scale-out).
		addedNew = append(addedNew, raw)
	}

	for nodeIdxStr, oldSpec := range oldByNodeIndex {
		if !matchedNodeIndexes[nodeIdxStr] {
			removedOld = append(removedOld, oldSpec)
		}
	}
	return
}

// emrModifyNodeInstanceType issues a single ModifyResource call to change a
// node's instance_type and waits for the resulting flow to finish. It is
// shared by master/core/task/common Update branches.
func emrModifyNodeInstanceType(
	ctx context.Context,
	meta interface{},
	d *schema.ResourceData,
	logId, instanceId, roleKey string,
	zoneIdx, nodeIdx int,
	oldSpec, newSpec map[string]interface{},
) error {
	newInstanceType, _ := newSpec["instance_type"].(string)
	resourceId, _ := oldSpec["emr_resource_id"].(string)
	newCpu, newMem, err := emrParseInstanceTypeCpuMem(newInstanceType)
	if err != nil {
		return fmt.Errorf("zone[%d] %s[%d]: failed to parse instance_type %q: %v",
			zoneIdx, roleKey, nodeIdx, newInstanceType, err)
	}

	req := emr.NewModifyResourceRequest()
	resp := emr.NewModifyResourceResponse()
	req.InstanceId = helper.String(instanceId)
	req.PayMode = emrPayModeFromChargeType(d.Get("instance_charge_type").(string))
	req.InstanceType = helper.String(newInstanceType)
	req.NewCpu = helper.Int64(newCpu)
	req.NewMem = helper.Int64(newMem)
	req.ResourceIdList = []*string{helper.String(resourceId)}

	conn := meta.(tccommon.ProviderMeta).GetAPIV3Conn()
	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := conn.UseEmrClient().ModifyResourceWithContext(ctx, req)
		if e != nil {
			return tccommon.RetryError(e)
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
			logId, req.GetAction(), req.ToJsonString(), result.ToJsonString())
		if result == nil || result.Response == nil || result.Response.TraceId == nil {
			return resource.NonRetryableError(fmt.Errorf("ModifyResource: response is nil"))
		}
		resp = result
		return nil
	})
	if reqErr != nil {
		log.Printf("[CRITAL]%s zone[%d] %s[%d] modify instance type failed: %+v",
			logId, zoneIdx, roleKey, nodeIdx, reqErr)
		return reqErr
	}

	traceId := *resp.Response.TraceId
	service := EMRService{client: conn}
	conf := tccommon.BuildStateChangeConf(
		[]string{"0", "1"}, []string{"2", "-1"},
		d.Timeout(schema.TimeoutUpdate)-time.Minute, time.Second,
		service.FlowStatusRefreshFunc(instanceId, traceId, F_KEY_TRACE_ID, []string{}),
	)
	if object, e := conf.WaitForState(); e != nil {
		return e
	} else if status, ok := object.(*int64); ok && status != nil && *status == -1 {
		return fmt.Errorf("EMR cluster flow failed (TraceId=%s), flow total status is -1", traceId)
	}
	return nil
}

// emrTerminateNodes issues a TerminateClusterNodes call to scale-in the given
// CVM instance IDs (orderNos) for the specified node role (CORE / TASK) and
// waits for completion.
func emrTerminateNodes(
	ctx context.Context,
	meta interface{},
	d *schema.ResourceData,
	logId, instanceId, nodeFlag string,
	zoneIdx int,
	orderNos []*string,
) error {
	if len(orderNos) == 0 {
		return nil
	}

	req := emr.NewTerminateClusterNodesRequest()
	resp := emr.NewTerminateClusterNodesResponse()
	req.InstanceId = helper.String(instanceId)
	req.CvmInstanceIds = orderNos
	req.NodeFlag = helper.String(nodeFlag)

	conn := meta.(tccommon.ProviderMeta).GetAPIV3Conn()
	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := conn.UseEmrClient().TerminateClusterNodesWithContext(ctx, req)
		if e != nil {
			return tccommon.RetryError(e)
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
			logId, req.GetAction(), req.ToJsonString(), result.ToJsonString())
		if result == nil || result.Response == nil || result.Response.FlowId == nil {
			return resource.NonRetryableError(fmt.Errorf("TerminateClusterNodes: response is nil"))
		}
		resp = result
		return nil
	})
	if reqErr != nil {
		log.Printf("[CRITAL]%s zone[%d] %s scale-in failed: %+v", logId, zoneIdx, nodeFlag, reqErr)
		return reqErr
	}

	flowId := strconv.FormatInt(*resp.Response.FlowId, 10)
	service := EMRService{client: conn}
	conf := tccommon.BuildStateChangeConf(
		[]string{"0", "1"}, []string{"2", "-1"},
		d.Timeout(schema.TimeoutUpdate)-time.Minute, time.Second,
		service.FlowStatusRefreshFunc(instanceId, flowId, F_KEY_FLOW_ID, []string{}),
	)
	if object, e := conf.WaitForState(); e != nil {
		return e
	} else if status, ok := object.(*int64); ok && status != nil && *status == -1 {
		return fmt.Errorf("EMR cluster flow failed (FlowId=%s), flow total status is -1", flowId)
	}
	return nil
}

// emrScaleOutSingleNode issues a ScaleOutCluster call to add ONE node of the
// specified role (CORE / TASK) using the spec block `addedSpec` and waits for
// the resulting flow to finish. zone is the placement zone string for the
// new node.
func emrScaleOutSingleNode(
	ctx context.Context,
	meta interface{},
	d *schema.ResourceData,
	logId, instanceId, nodeFlag, zone string,
	zoneIdx, addedIdx int,
	addedSpec map[string]interface{},
) error {
	chargeType := d.Get("instance_charge_type").(string)

	req := emr.NewScaleOutClusterRequest()
	resp := emr.NewScaleOutClusterResponse()
	req.InstanceId = helper.String(instanceId)
	req.InstanceChargeType = helper.String(chargeType)
	req.ScaleOutNodeConfig = &emr.ScaleOutNodeConfig{
		NodeFlag:  common.StringPtr(nodeFlag),
		NodeCount: common.Uint64Ptr(1),
	}
	if zone != "" {
		req.Zone = helper.String(zone)
	}

	if chargeType == "PREPAID" {
		if v, ok := d.GetOk("instance_charge_prepaid"); ok {
			prepaidList, _ := v.([]interface{})
			if len(prepaidList) > 0 {
				if prepaidMap, ok := prepaidList[0].(map[string]interface{}); ok {
					prepaid := &emr.InstanceChargePrepaid{}
					if val, ok := prepaidMap["period"].(int); ok && val != 0 {
						prepaid.Period = helper.IntInt64(val)
					}
					if val, ok := prepaidMap["renew_flag"].(bool); ok {
						prepaid.RenewFlag = helper.Bool(val)
					}
					req.InstanceChargePrepaid = prepaid
				}
			}
		}
	}

	resourceSpec := &emr.NodeResourceSpec{}
	if instanceType, ok := addedSpec["instance_type"].(string); ok && instanceType != "" {
		resourceSpec.InstanceType = helper.String(instanceType)
	}
	if sdList, ok := addedSpec["system_disk"].([]interface{}); ok && len(sdList) > 0 {
		if sdMap, ok := sdList[0].(map[string]interface{}); ok {
			resourceSpec.SystemDisk = []*emr.DiskSpecInfo{{
				DiskType: helper.String(sdMap["disk_type"].(string)),
				Count:    helper.Int64(1),
				DiskSize: helper.Int64(int64(sdMap["disk_size"].(int))),
			}}
		}
	}
	for _, ddItem := range emrDiskSetToList(addedSpec["data_disk"]) {
		ddMap, ok := ddItem.(map[string]interface{})
		if !ok {
			continue
		}
		resourceSpec.DataDisk = append(resourceSpec.DataDisk, &emr.DiskSpecInfo{
			DiskType: helper.String(ddMap["disk_type"].(string)),
			Count:    helper.Int64(1),
			DiskSize: helper.Int64(int64(ddMap["disk_size"].(int))),
		})
	}
	req.ResourceSpec = resourceSpec

	conn := meta.(tccommon.ProviderMeta).GetAPIV3Conn()
	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := conn.UseEmrClient().ScaleOutClusterWithContext(ctx, req)
		if e != nil {
			return tccommon.RetryError(e)
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
			logId, req.GetAction(), req.ToJsonString(), result.ToJsonString())
		if result == nil || result.Response == nil || result.Response.TraceId == nil {
			return resource.NonRetryableError(fmt.Errorf("ScaleOutCluster: response is nil"))
		}
		resp = result
		return nil
	})
	if reqErr != nil {
		log.Printf("[CRITAL]%s zone[%d] %s[added:%d] scale-out failed: %+v",
			logId, zoneIdx, nodeFlag, addedIdx, reqErr)
		return reqErr
	}

	traceId := *resp.Response.TraceId
	service := EMRService{client: conn}
	conf := tccommon.BuildStateChangeConf(
		[]string{"0", "1"}, []string{"2", "-1"},
		d.Timeout(schema.TimeoutUpdate)-time.Minute, time.Second,
		service.FlowStatusRefreshFunc(instanceId, traceId, F_KEY_TRACE_ID, []string{}),
	)
	if object, e := conf.WaitForState(); e != nil {
		return e
	} else if status, ok := object.(*int64); ok && status != nil && *status == -1 {
		return fmt.Errorf("EMR cluster flow failed (TraceId=%s), flow total status is -1", traceId)
	}
	return nil
}

// emrZoneOf returns placement.0.zone of a zone_resource_configuration block,
// or "" if missing.
func emrZoneOf(zrcMap map[string]interface{}) string {
	plList, _ := zrcMap["placement"].([]interface{})
	if len(plList) == 0 {
		return ""
	}
	plMap, _ := plList[0].(map[string]interface{})
	if plMap == nil {
		return ""
	}
	z, _ := plMap["zone"].(string)
	return z
}
