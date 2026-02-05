package postgresql

import (
	"context"
	"log"
	"strings"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svctag "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/tag"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	postgresql "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/postgres/v20170312"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudPostgresqlInstances() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudPostgresqlInstanceRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Name of the postgresql instance to be query.",
			},
			"id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "ID of the postgresql instance to be query.",
			},
			"project_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Project ID of the postgresql instance to be query.",
			},
			"instance_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Deprecated:  "It has been deprecated from version 1.82.64. Please use `db_instance_set` instead.",
				Description: "A list of postgresql instances. Each element contains the following attributes.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of the postgresql instance.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of the postgresql instance.",
						},
						"charge_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Pay type of the postgresql instance.",
						},
						"auto_renew_flag": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Auto renew flag.",
						},
						"engine_version": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Version of the postgresql database engine.",
						},
						"db_kernel_version": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "PostgreSQL kernel version number.",
						},
						"db_major_version": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "PostgreSQL major version number.",
						},
						"vpc_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of VPC.",
						},
						"subnet_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of subnet.",
						},
						"storage": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Volume size(in GB).",
						},
						"memory": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Memory size(in GB).",
						},
						"project_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Project id, default value is 0.",
						},
						"availability_zone": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Availability zone.",
						},
						"root_user": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Instance root account name, default value is `root`.",
						},
						"public_access_switch": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Indicates whether to enable the access to an instance from public network or not.",
						},
						"public_access_host": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Host for public access.",
						},
						"public_access_port": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Port for public access.",
						},
						"private_access_ip": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "IP address for private access.",
						},
						"private_access_port": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Port for private access.",
						},
						"charset": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Charset of the postgresql instance.",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Create time of the postgresql instance.",
						},
						"tags": {
							Type:        schema.TypeMap,
							Computed:    true,
							Description: "The available tags within this postgresql.",
						},
					},
				},
			},
			"db_instance_set": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Instance details set.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"region": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Instance region such as ap-guangzhou, which corresponds to the`Region` field in `RegionSet`.",
						},
						"zone": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Instance AZ such as ap-guangzhou-3, which corresponds to the `Zone` field of `ZoneSet`.",
						},
						"vpc_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "vpc ID, such as vpc-e6w23k31. a valid vpc ID can be obtained by logging in to the console to query or by calling the api [DescribeVpcs](https://www.tencentcloud.comom/document/api/215/15778?from_cn_redirect=1) and acquiring the unVpcId field in the api return.",
						},
						"subnet_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "VPC subnet ID, such as subnet-51lcif9y. effective VPC subnet ids can be obtained by logging in to the console or calling the api [DescribeSubnets](https://www.tencentcloud.comom/document/api/215/15784?from_cn_redirect=1) to acquire the unSubnetId field in the api return.",
						},
						"db_instance_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Instance ID.",
						},
						"db_instance_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Instance name.",
						},
						"db_instance_status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Instance status, including: `applying` (applying), `init` (to be initialized), `initing` (initializing), `running` (running), `limited run` (restricted operation), `isolating` (isolating), `isolated` (isolated), `disisolating` (de-isolating), `recycling` (recycling), `recycled` (recycled), `job running` (task executing), `offline` (offline), `migrating` (migrating), `expanding` (scaling out), `waitSwitch` (waiting to switch), `switching` (switching), `readonly` (readonly), `restarting` (restarting), `network changing` (network modification in progress), `upgrading` (kernel version upgrading), `audit-switching` (audit status changing), `primary-switching` (primary-secondary switching), `offlining` (offline), `deployment changing` (modify az), `cloning` (restoring data), `parameter modifying` (parameter modification in progress), `log-switching` (log status change), `restoring` (recovering), and `expanding` (scaling out).",
						},
						"db_instance_memory": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Assigned instance memory size in GB.",
						},
						"db_instance_storage": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Assigned instance storage capacity in GB.",
						},
						"db_instance_cpu": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Number of assigned CPUs.",
						},
						"db_instance_class": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Purchasable specification ID.",
						},
						"db_major_version": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "PostgreSQL major version number. specifies the version information that can be obtained from the [DescribeDBVersions](https://www.tencentcloud.comom/document/api/409/89018?from_cn_redirect=1) api. currently supports major versions 10, 11, 12, 13, 14, 15.",
						},
						"db_version": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Number of the major PostgreSQL community version and minor version, such as 12.4, which can be queried by the [DescribeDBVersions](https://intl.cloud.tencent.com/document/api/409/89018?from_cn_redirect=1) API.",
						},
						"db_kernel_version": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "PostgreSQL kernel version, such as v12.7_r1.8. version information can be obtained from [DescribeDBVersions](https://www.tencentcloud.comom/document/api/409/89018?from_cn_redirect=1).",
						},
						"db_instance_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Instance type, which includes:\n<li>primary: primary instance </li>\n<li>readonly: read-only instance</li>\n<li>guard: disaster recovery instance</li>\n<li>temp: temporary instance</li>.",
						},
						"db_instance_version": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Instance version. Valid value: `standard` (dual-server high-availability; one-primary-one-standby).",
						},
						"db_charset": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Instance character set, which currently supports only:\n<li>UTF8</li>\n<li>LATIN1</li>.",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Instance creation time.",
						},
						"update_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Last updated time of the instance attribute.",
						},
						"expire_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Instance expiration time.",
						},
						"isolated_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Instance isolation time.",
						},
						"pay_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Billing mode:\n<li>prepaid: monthly subscription, prepaid</li>\n<li>postpaid: pay-as-you-go, postpaid</li>.",
						},
						"auto_renew": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Auto-renewal or not:\n<li>`0`: manual renewal</li>\n<li>`1`: auto-renewal</li>\nDefault value: 0.",
						},
						"db_instance_net_info": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Instance network connection information.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"address": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "DNS domain name.",
									},
									"ip": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Ip.",
									},
									"port": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Connection port address.",
									},
									"net_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Network type. 1: inner (private network address), 2: public (public network address).",
									},
									"status": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Network connection status. Valid values: `initing` (never enabled before), `opened` (enabled), `closed` (disabled), `opening` (enabling), `closing` (disabling).",
									},
									"vpc_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "VPC ID. specifies the ID of the virtual private cloud.",
									},
									"subnet_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Subnet ID.",
									},
									"protocol_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Specifies the protocol type to connect to the database. currently supported: postgresql, mssql (mssql compatible syntax).",
									},
								},
							},
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Machine type.",
						},
						"app_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "User `AppId`.",
						},
						"uid": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Instance `Uid`.",
						},
						"project_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Project ID.",
						},
						"tag_list": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Describes the Tag information associated with the instance.",
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
						"master_db_instance_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Primary instance information. returned only when the instance is a read-only instance.",
						},
						"read_only_instance_num": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Specifies the number of read-only instances.",
						},
						"status_in_readonly_group": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Describes the state of the read-only instance in the read-only group.",
						},
						"offline_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Decommissioning time.",
						},
						"db_node_set": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Instance node information\nNote: This field may return null, indicating that no valid values can be obtained.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"role": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Node type. Valid values:\n`Primary`;\n`Standby`.",
									},
									"zone": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "AZ where the node resides, such as ap-guangzhou-1.",
									},
									"dedicated_cluster_id": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "CDC ID.",
									},
								},
							},
						},
						"is_support_tde": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Whether the instance supports TDE data encryption.\n<Li>0: not supported</li>.\n<Li>1: supported.</li>.\nDefault value: 0\n\nFor TDE data encryption, see [overview of data transparent encryption](https://www.tencentcloud.comom/document/product/409/71748?from_cn_redirect=1).",
						},
						"db_engine": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Database engine, which supports:\n<li>`postgresql`: tencentdb for postgresql</li>.\n<li>`mssql_compatible`: specifies mssql compatible - tencentdb for PostgreSQL.</li>.\nDefault value: `postgresql`.",
						},
						"db_engine_config": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Configuration information for the database engine, and the configuration format is as follows:.\n{`$key1`:`$value1`, `$key2`:`$value2`}\nSupported engines include:.\nmssql_compatible engine:.\n<li>migrationMode: specifies the database mode. optional parameter. valid values: single-db (single-database schema) and multi-db (multiple database schemas). defaults to single-db.</li>.\n<li>defaultLocale: specifies the sorting area rule, an optional parameter that cannot be modified after initialization. default value is en_US. valid values include:.\n`af_ZA`, `sq_AL`, `ar_DZ`, `ar_BH`, `ar_EG`, `ar_IQ`, `ar_JO`, `ar_KW`, `ar_LB`, `ar_LY`, `ar_MA`, `ar_OM`, `ar_QA`, `ar_SA`, `ar_SY`, `ar_TN`, `ar_AE`, `ar_YE`, `hy_AM`, `az_Cyrl_AZ`, `az_Latn_AZ`, `eu_ES`, `be_BY`, `bg_BG`, `ca_ES`, `zh_HK`, `zh_MO`, `zh_CN`, `zh_SG`, `zh_TW`, `hr_HR`, `cs_CZ`, `da_DK`, `nl_BE`, `nl_NL`, `en_AU`, `en_BZ`, `en_CA`, `en_IE`, `en_JM`, `en_NZ`, `en_PH`, `en_ZA`, `en_TT`, `en_GB`, `en_US`, `en_ZW`, `et_EE`, `fo_FO`, `fa_IR`, `fi_FI`, `fr_BE`, `fr_CA`, `fr_FR`, `fr_LU`, `fr_MC`, `fr_CH`, `mk_MK`, `ka_GE`, `de_AT`, `de_DE`, `de_LI`, `de_LU`, `de_CH`, `el_GR`, `gu_IN`, `he_IL`, `hi_IN`, `hu_HU`, `is_IS`, `id_ID`, `it_IT`, `it_CH`, `ja_JP`, `kn_IN`, `kok_IN`, `ko_KR`, `ky_KG`, `lv_LV`, `lt_LT`, `ms_BN`, `ms_MY`, `mr_IN`, `mn_MN`, `nb_NO`, `nn_NO`, `pl_PL`, `pt_BR`, `pt_PT`, `pa_IN`, `ro_RO`, `ru_RU`, `sa_IN`, `sr_Cyrl_RS`, `sr_Latn_RS`, `sk_SK`, `sl_SI`, `es_AR`, `es_BO`, `es_CL`, `es_CO`, `es_CR`, `es_DO`, `es_EC`, `es_SV`, `es_GT`, `es_HN`, `es_MX`, `es_NI`, `es_PA`, `es_PY`,`es_PE`, `es_PR`, `es_ES`, `es_TRADITIONAL`, `es_UY`, `es_VE`, `sw_KE`, `sv_FI`, `sv_SE`, `tt_RU`, `te_IN`, `th_TH`, `tr_TR`, `uk_UA`, `ur_IN`, `ur_PK`, `uz_Cyrl_UZ`, `uz_Latn_UZ`, `vi_VN`.</li>\n<li>serverCollationName: Sorting rule name, an optional parameter, which cannot be modified after initialization, its default value is sql_latin1_general_cp1_ci_as, and its valid values include: `bbf_unicode_general_ci_as`, `bbf_unicode_cp1_ci_as`, `bbf_unicode_CP1250_ci_as`, `bbf_unicode_CP1251_ci_as`, `bbf_unicode_cp1253_ci_as`, `bbf_unicode_cp1254_ci_as`, `bbf_unicode_cp1255_ci_as`, `bbf_unicode_cp1256_ci_as`, `bbf_unicode_cp1257_ci_as`, `bbf_unicode_cp1258_ci_as`, `bbf_unicode_cp874_ci_as`, `sql_latin1_general_cp1250_ci_as`, `sql_latin1_general_cp1251_ci_as`, `sql_latin1_general_cp1_ci_as`, `sql_latin1_general_cp1253_ci_as`, `sql_latin1_general_cp1254_ci_as`, `sql_latin1_general_cp1255_ci_as`, `sql_latin1_general_cp1256_ci_as`, `sql_latin1_general_cp1257_ci_as`, `sql_latin1_general_cp1258_ci_as`, `chinese_prc_ci_as`, `cyrillic_general_ci_as`, `finnish_swedish_ci_as`, `french_ci_as`, `japanese_ci_as`, `korean_wansung_ci_as`, `latin1_general_ci_as`, `modern_spanish_ci_as`, `polish_ci_as`, `thai_ci_as`, `traditional_spanish_ci_as`, `turkish_ci_as`, `ukrainian_ci_as`, and `vietnamese_ci_as`.</li>.",
						},
						"network_access_list": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Network access list of the instance (this field has been deprecated)\nNote: this field may return `null`, indicating that no valid values can be obtained.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"resource_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Network resource id, instance id, or RO group id.",
									},
									"resource_type": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Resource type. valid values: 1 (instance), 2 (RO group).",
									},
									"vpc_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "VPC ID. specifies the ID of the virtual private cloud.",
									},
									"vip": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "IPv4 Address.",
									},
									"vip6": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "IPv6 Address.",
									},
									"vport": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Specifies the access port.",
									},
									"subnet_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Subnet ID.",
									},
									"vpc_status": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Network status. valid values: 1-applying, 2-active, 3-deleting, 4-deleted.",
									},
								},
							},
						},
						"support_ipv6": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Whether the instance supports IPv6:\n<li>`0`: no</li>\n<li>`1`: yes</li>\nDefault value: 0.",
						},
						"expanded_cpu": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Number of cpu cores that have been elastically scaled out.",
						},
						"deletion_protection": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Specifies whether to enable deletion protection for the instance. valid values as follows:.\n-Specifies whether to enable deletion protection. valid values: true (enable deletion protection).\n-Specifies whether to disable deletion protection. valid values: false (disable deletion protection).",
						},
						"root_user": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Instance root account name, default value is `root`.",
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

func dataSourceTencentCloudPostgresqlInstanceRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_postgresql_instances.read")()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service = PostgresqlService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	filter := make([]*postgresql.Filter, 0)
	if v, ok := d.GetOk("name"); ok {
		filter = append(filter, &postgresql.Filter{Name: helper.String("db-instance-name"), Values: []*string{helper.String(v.(string))}})
	}

	if v, ok := d.GetOk("id"); ok {
		filter = append(filter, &postgresql.Filter{Name: helper.String("db-instance-id"), Values: []*string{helper.String(v.(string))}})
	}

	if v, ok := d.GetOk("project_id"); ok {
		filter = append(filter, &postgresql.Filter{Name: helper.String("db-project-id"), Values: []*string{helper.String(v.(string))}})
	}

	instanceList, err := service.DescribePostgresqlInstances(ctx, filter)
	if err != nil {
		instanceList, err = service.DescribePostgresqlInstances(ctx, filter)
	}

	if err != nil {
		return err
	}

	ids := make([]string, 0, len(instanceList))
	list := make([]map[string]interface{}, 0, len(instanceList))
	tcClient := meta.(tccommon.ProviderMeta).GetAPIV3Conn()
	tagService := svctag.NewTagService(tcClient)

	// old
	for _, v := range instanceList {
		listItem := make(map[string]interface{})
		listItem["id"] = v.DBInstanceId
		listItem["name"] = v.DBInstanceName
		listItem["auto_renew_flag"] = v.AutoRenew
		listItem["project_id"] = v.ProjectId
		listItem["storage"] = v.DBInstanceStorage
		listItem["memory"] = v.DBInstanceMemory
		listItem["availability_zone"] = v.Zone
		listItem["create_time"] = v.CreateTime
		listItem["vpc_id"] = v.VpcId
		listItem["subnet_id"] = v.SubnetId
		listItem["engine_version"] = v.DBVersion
		if v.DBKernelVersion != nil {
			listItem["db_kernel_version"] = v.DBKernelVersion
		}

		if v.DBMajorVersion != nil {
			listItem["db_major_version"] = v.DBMajorVersion
		}

		listItem["public_access_switch"] = false
		listItem["charset"] = v.DBCharset
		listItem["public_access_host"] = ""

		// rootUser
		if v.DBInstanceId != nil && strings.HasPrefix(*v.DBInstanceId, "postgres-") {
			accounts, outErr := service.DescribeRootUser(ctx, *v.DBInstanceId)
			if outErr != nil {
				continue
			}

			if len(accounts) > 0 {
				listItem["root_user"] = accounts[0].UserName
			}
		}

		for _, netInfo := range v.DBInstanceNetInfo {
			if *netInfo.NetType == "public" {
				if *netInfo.Status == "opened" || *netInfo.Status == "1" {
					listItem["public_access_switch"] = true
				}
				listItem["public_access_host"] = netInfo.Address
				listItem["public_access_port"] = netInfo.Port
			}
			if (*netInfo.NetType == "private" || *netInfo.NetType == "inner") && *netInfo.Ip != "" {
				listItem["private_access_ip"] = netInfo.Ip
				listItem["private_access_port"] = netInfo.Port
			}
		}

		if *v.PayType == POSTGRESQL_PAYTYPE_PREPAID || *v.PayType == COMMON_PAYTYPE_PREPAID {
			listItem["charge_type"] = COMMON_PAYTYPE_PREPAID
		} else {
			listItem["charge_type"] = COMMON_PAYTYPE_POSTPAID
		}

		//the describe list API is delayed with argument `tag`
		tagList, err := tagService.DescribeResourceTags(ctx, "postgres", "DBInstanceId", tcClient.Region, *v.DBInstanceId)
		if err != nil {
			return err
		}

		listItem["tags"] = tagList

		list = append(list, listItem)
		ids = append(ids, *v.DBInstanceId)
	}

	// new
	dBInstanceSetList := make([]map[string]interface{}, 0, len(instanceList))
	for _, dBInstanceSet := range instanceList {
		dBInstanceSetMap := map[string]interface{}{}
		if dBInstanceSet.Region != nil {
			dBInstanceSetMap["region"] = dBInstanceSet.Region
		}

		if dBInstanceSet.Zone != nil {
			dBInstanceSetMap["zone"] = dBInstanceSet.Zone
		}

		if dBInstanceSet.VpcId != nil {
			dBInstanceSetMap["vpc_id"] = dBInstanceSet.VpcId
		}

		if dBInstanceSet.SubnetId != nil {
			dBInstanceSetMap["subnet_id"] = dBInstanceSet.SubnetId
		}

		if dBInstanceSet.DBInstanceId != nil {
			dBInstanceSetMap["db_instance_id"] = dBInstanceSet.DBInstanceId
		}

		if dBInstanceSet.DBInstanceName != nil {
			dBInstanceSetMap["db_instance_name"] = dBInstanceSet.DBInstanceName
		}

		if dBInstanceSet.DBInstanceStatus != nil {
			dBInstanceSetMap["db_instance_status"] = dBInstanceSet.DBInstanceStatus
		}

		if dBInstanceSet.DBInstanceMemory != nil {
			dBInstanceSetMap["db_instance_memory"] = dBInstanceSet.DBInstanceMemory
		}

		if dBInstanceSet.DBInstanceStorage != nil {
			dBInstanceSetMap["db_instance_storage"] = dBInstanceSet.DBInstanceStorage
		}

		if dBInstanceSet.DBInstanceCpu != nil {
			dBInstanceSetMap["db_instance_cpu"] = dBInstanceSet.DBInstanceCpu
		}

		if dBInstanceSet.DBInstanceClass != nil {
			dBInstanceSetMap["db_instance_class"] = dBInstanceSet.DBInstanceClass
		}

		if dBInstanceSet.DBMajorVersion != nil {
			dBInstanceSetMap["db_major_version"] = dBInstanceSet.DBMajorVersion
		}

		if dBInstanceSet.DBVersion != nil {
			dBInstanceSetMap["db_version"] = dBInstanceSet.DBVersion
		}

		if dBInstanceSet.DBKernelVersion != nil {
			dBInstanceSetMap["db_kernel_version"] = dBInstanceSet.DBKernelVersion
		}

		if dBInstanceSet.DBInstanceType != nil {
			dBInstanceSetMap["db_instance_type"] = dBInstanceSet.DBInstanceType
		}

		if dBInstanceSet.DBInstanceVersion != nil {
			dBInstanceSetMap["db_instance_version"] = dBInstanceSet.DBInstanceVersion
		}

		if dBInstanceSet.DBCharset != nil {
			dBInstanceSetMap["db_charset"] = dBInstanceSet.DBCharset
		}

		if dBInstanceSet.CreateTime != nil {
			dBInstanceSetMap["create_time"] = dBInstanceSet.CreateTime
		}

		if dBInstanceSet.UpdateTime != nil {
			dBInstanceSetMap["update_time"] = dBInstanceSet.UpdateTime
		}

		if dBInstanceSet.ExpireTime != nil {
			dBInstanceSetMap["expire_time"] = dBInstanceSet.ExpireTime
		}

		if dBInstanceSet.IsolatedTime != nil {
			dBInstanceSetMap["isolated_time"] = dBInstanceSet.IsolatedTime
		}

		if dBInstanceSet.PayType != nil {
			dBInstanceSetMap["pay_type"] = dBInstanceSet.PayType
		}

		if dBInstanceSet.AutoRenew != nil {
			dBInstanceSetMap["auto_renew"] = dBInstanceSet.AutoRenew
		}

		dBInstanceNetInfoList := make([]map[string]interface{}, 0, len(dBInstanceSet.DBInstanceNetInfo))
		if dBInstanceSet.DBInstanceNetInfo != nil {
			for _, dBInstanceNetInfo := range dBInstanceSet.DBInstanceNetInfo {
				dBInstanceNetInfoMap := map[string]interface{}{}
				if dBInstanceNetInfo.Address != nil {
					dBInstanceNetInfoMap["address"] = dBInstanceNetInfo.Address
				}

				if dBInstanceNetInfo.Ip != nil {
					dBInstanceNetInfoMap["ip"] = dBInstanceNetInfo.Ip
				}

				if dBInstanceNetInfo.Port != nil {
					dBInstanceNetInfoMap["port"] = dBInstanceNetInfo.Port
				}

				if dBInstanceNetInfo.NetType != nil {
					dBInstanceNetInfoMap["net_type"] = dBInstanceNetInfo.NetType
				}

				if dBInstanceNetInfo.Status != nil {
					dBInstanceNetInfoMap["status"] = dBInstanceNetInfo.Status
				}

				if dBInstanceNetInfo.VpcId != nil {
					dBInstanceNetInfoMap["vpc_id"] = dBInstanceNetInfo.VpcId
				}

				if dBInstanceNetInfo.SubnetId != nil {
					dBInstanceNetInfoMap["subnet_id"] = dBInstanceNetInfo.SubnetId
				}

				if dBInstanceNetInfo.ProtocolType != nil {
					dBInstanceNetInfoMap["protocol_type"] = dBInstanceNetInfo.ProtocolType
				}

				dBInstanceNetInfoList = append(dBInstanceNetInfoList, dBInstanceNetInfoMap)
			}

			dBInstanceSetMap["db_instance_net_info"] = dBInstanceNetInfoList
		}

		if dBInstanceSet.Type != nil {
			dBInstanceSetMap["type"] = dBInstanceSet.Type
		}

		if dBInstanceSet.AppId != nil {
			dBInstanceSetMap["app_id"] = dBInstanceSet.AppId
		}

		if dBInstanceSet.Uid != nil {
			dBInstanceSetMap["uid"] = dBInstanceSet.Uid
		}

		if dBInstanceSet.ProjectId != nil {
			dBInstanceSetMap["project_id"] = dBInstanceSet.ProjectId
		}

		tagListList := make([]map[string]interface{}, 0, len(dBInstanceSet.TagList))
		if dBInstanceSet.TagList != nil {
			for _, tagList := range dBInstanceSet.TagList {
				tagListMap := map[string]interface{}{}
				if tagList.TagKey != nil {
					tagListMap["tag_key"] = tagList.TagKey
				}

				if tagList.TagValue != nil {
					tagListMap["tag_value"] = tagList.TagValue
				}

				tagListList = append(tagListList, tagListMap)
			}

			dBInstanceSetMap["tag_list"] = tagListList
		}

		if dBInstanceSet.MasterDBInstanceId != nil {
			dBInstanceSetMap["master_db_instance_id"] = dBInstanceSet.MasterDBInstanceId
		}

		if dBInstanceSet.ReadOnlyInstanceNum != nil {
			dBInstanceSetMap["read_only_instance_num"] = dBInstanceSet.ReadOnlyInstanceNum
		}

		if dBInstanceSet.StatusInReadonlyGroup != nil {
			dBInstanceSetMap["status_in_readonly_group"] = dBInstanceSet.StatusInReadonlyGroup
		}

		if dBInstanceSet.OfflineTime != nil {
			dBInstanceSetMap["offline_time"] = dBInstanceSet.OfflineTime
		}

		dBNodeSetList := make([]map[string]interface{}, 0, len(dBInstanceSet.DBNodeSet))
		if dBInstanceSet.DBNodeSet != nil {
			for _, dBNodeSet := range dBInstanceSet.DBNodeSet {
				dBNodeSetMap := map[string]interface{}{}
				if dBNodeSet.Role != nil {
					dBNodeSetMap["role"] = dBNodeSet.Role
				}

				if dBNodeSet.Zone != nil {
					dBNodeSetMap["zone"] = dBNodeSet.Zone
				}

				if dBNodeSet.DedicatedClusterId != nil {
					dBNodeSetMap["dedicated_cluster_id"] = dBNodeSet.DedicatedClusterId
				}

				dBNodeSetList = append(dBNodeSetList, dBNodeSetMap)
			}

			dBInstanceSetMap["db_node_set"] = dBNodeSetList
		}

		if dBInstanceSet.IsSupportTDE != nil {
			dBInstanceSetMap["is_support_tde"] = dBInstanceSet.IsSupportTDE
		}

		if dBInstanceSet.DBEngine != nil {
			dBInstanceSetMap["db_engine"] = dBInstanceSet.DBEngine
		}

		if dBInstanceSet.DBEngineConfig != nil {
			dBInstanceSetMap["db_engine_config"] = dBInstanceSet.DBEngineConfig
		}

		networkAccessListList := make([]map[string]interface{}, 0, len(dBInstanceSet.NetworkAccessList))
		if dBInstanceSet.NetworkAccessList != nil {
			for _, networkAccessList := range dBInstanceSet.NetworkAccessList {
				networkAccessListMap := map[string]interface{}{}
				if networkAccessList.ResourceId != nil {
					networkAccessListMap["resource_id"] = networkAccessList.ResourceId
				}

				if networkAccessList.ResourceType != nil {
					networkAccessListMap["resource_type"] = networkAccessList.ResourceType
				}

				if networkAccessList.VpcId != nil {
					networkAccessListMap["vpc_id"] = networkAccessList.VpcId
				}

				if networkAccessList.Vip != nil {
					networkAccessListMap["vip"] = networkAccessList.Vip
				}

				if networkAccessList.Vip6 != nil {
					networkAccessListMap["vip6"] = networkAccessList.Vip6
				}

				if networkAccessList.Vport != nil {
					networkAccessListMap["vport"] = networkAccessList.Vport
				}

				if networkAccessList.SubnetId != nil {
					networkAccessListMap["subnet_id"] = networkAccessList.SubnetId
				}

				if networkAccessList.VpcStatus != nil {
					networkAccessListMap["vpc_status"] = networkAccessList.VpcStatus
				}

				networkAccessListList = append(networkAccessListList, networkAccessListMap)
			}

			dBInstanceSetMap["network_access_list"] = networkAccessListList
		}

		if dBInstanceSet.SupportIpv6 != nil {
			dBInstanceSetMap["support_ipv6"] = dBInstanceSet.SupportIpv6
		}

		if dBInstanceSet.ExpandedCpu != nil {
			dBInstanceSetMap["expanded_cpu"] = dBInstanceSet.ExpandedCpu
		}

		if dBInstanceSet.DeletionProtection != nil {
			dBInstanceSetMap["deletion_protection"] = dBInstanceSet.DeletionProtection
		}

		// rootUser
		if dBInstanceSet.DBInstanceId != nil && strings.HasPrefix(*dBInstanceSet.DBInstanceId, "postgres-") {
			accounts, outErr := service.DescribeRootUser(ctx, *dBInstanceSet.DBInstanceId)
			if outErr != nil {
				continue
			}

			if len(accounts) > 0 {
				dBInstanceSetMap["root_user"] = accounts[0].UserName
			}
		}

		dBInstanceSetList = append(dBInstanceSetList, dBInstanceSetMap)
	}

	_ = d.Set("db_instance_set", dBInstanceSetList)

	d.SetId(helper.DataResourceIdsHash(ids))
	if e := d.Set("instance_list", list); e != nil {
		log.Printf("[CRITAL]%s provider set list fail, reason:%s\n", logId, e.Error())
		return e
	}

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		return tccommon.WriteToFile(output.(string), list)
	}

	return nil
}
