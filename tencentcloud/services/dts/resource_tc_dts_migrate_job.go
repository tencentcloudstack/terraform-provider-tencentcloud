package dts

import (
	"context"
	"fmt"
	"log"
	"time"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dts "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dts/v20211206"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudDtsMigrateJob() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudDtsMigrateJobCreate,
		Read:   resourceTencentCloudDtsMigrateJobRead,
		Update: resourceTencentCloudDtsMigrateJobUpdate,
		Delete: resourceTencentCloudDtsMigrateJobDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"service_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Migrate service Id from `tencentcloud_dts_migrate_service`.",
			},

			"status": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Task status. Valid values: created(Created), checking (Checking), checkPass (Check passed), checkNotPass (Check not passed), readyRun (Ready for running), running (Running), readyComplete (Preparation completed), success (Successful), failed (Failed), stopping (Stopping), completing (Completing), pausing (Pausing), manualPaused (Paused).",
			},

			// for modify operation
			"run_mode": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Running mode. Valid values: immediate, timed.",
			},

			"migrate_option": {
				Required:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Migration job configuration options, used to describe how the task performs migration.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"database_table": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Required:    true,
							Description: "Migration object option, you need to tell the migration service which library table objects to migrate.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"object_mode": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Migration object type. Valid values: all, partial.",
									},
									"databases": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "Migration object, which is required if ObjectMode is partial.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"db_name": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Name of the database to be migrated or synced, which is required if ObjectMode is partial.",
												},
												"new_db_name": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Name of the database after migration or sync, which is the same as the source database name by default.",
												},
												"schema_name": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "The schema to be migrated or synced.",
												},
												"new_schema_name": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Name of the schema after migration or sync.",
												},
												"db_mode": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Database selection mode, which is required if ObjectMode is partial. Valid values: all, partial.",
												},
												"schema_mode": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Schema selection mode. Valid values: all, partial.",
												},
												"table_mode": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Table selection mode, which is required if DBMode is partial. Valid values: all, partial.",
												},
												"tables": {
													Type:        schema.TypeList,
													Optional:    true,
													Description: "The set of table objects, which is required if TableMode is partial.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"table_name": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Name of the migrated table, which is case-sensitive.",
															},
															"new_table_name": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "New name of the migrated table. This parameter is required when TableEditMode is rename. It is mutually exclusive with TmpTables..",
															},
															"tmp_tables": {
																Type: schema.TypeSet,
																Elem: &schema.Schema{
																	Type: schema.TypeString,
																},
																Optional:    true,
																Computed:    true,
																Description: "The temp tables to be migrated. This parameter is mutually exclusive with NewTableName. It is valid only when the configured migration objects are table-level ones and TableEditMode is pt. To migrate temp tables generated when pt-osc or other tools are used during the migration process, you must configure this parameter first. For example, if you want to perform the pt-osc operation on a table named 't1', configure this parameter as ['_t1_new','_t1_old']; to perform the gh-ost operation on t1, configure it as ['_t1_ghc','_t1_gho','_t1_del']. Temp tables generated by pt-osc and gh-ost operations can be configured at the same time.",
															},
															"table_edit_mode": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Table editing type. Valid values: rename (table mapping); pt (additional table sync).",
															},
														},
													},
												},
												"view_mode": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "View selection mode. Valid values: all, partial.",
												},
												"views": {
													Type:        schema.TypeList,
													Optional:    true,
													Description: "The set of view objects, which is required if ViewMode is partial.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"view_name": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "View name.",
															},
															"new_view_name": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "View name after migration.",
															},
														},
													},
												},
												"role_mode": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Role selection mode, which is exclusive to PostgreSQL. Valid values: all, partial.",
												},
												"roles": {
													Type:        schema.TypeList,
													Optional:    true,
													Description: "Role, which is exclusive to PostgreSQL and required if RoleMode is partial.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"role_name": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Role name.",
															},
															"new_role_name": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Role name after migration.",
															},
														},
													},
												},
												"function_mode": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Sync mode. Valid values: partial, all.",
												},
												"trigger_mode": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Sync mode. Valid values: partial, all.",
												},
												"event_mode": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Sync mode. Valid values: partial, all.",
												},
												"procedure_mode": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Sync mode. Valid values: partial, all.",
												},
												"functions": {
													Type: schema.TypeSet,
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
													Optional:    true,
													Computed:    true,
													Description: "This parameter is required if FunctionMode is partial.",
												},
												"procedures": {
													Type: schema.TypeSet,
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
													Optional:    true,
													Computed:    true,
													Description: "This parameter is required if ProcedureMode is partial.",
												},
												"events": {
													Type: schema.TypeSet,
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
													Optional:    true,
													Computed:    true,
													Description: "This parameter is required if EventMode is partial.",
												},
												"triggers": {
													Type: schema.TypeSet,
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
													Optional:    true,
													Computed:    true,
													Description: "This parameter is required if TriggerMode is partial.",
												},
											},
										},
									},
									"advanced_objects": {
										Type: schema.TypeSet,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
										Optional:    true,
										Computed:    true,
										Description: "Advanced object types, such as trigger, function, procedure, event. Note: If you want to migrate and synchronize advanced objects, the corresponding advanced object type should be included in this configuration.",
									},
								},
							},
						},
						"migrate_type": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Migration type. Valid values: full, structure, fullAndIncrement. Default value: fullAndIncrement.",
						},
						"consistency": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Computed:    true,
							Description: "Data consistency check option. Data consistency check is disabled by default.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"mode": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Data consistency check type. Valid values: full, noCheck, notConfigured.",
									},
								},
							},
						},
						"is_migrate_account": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Whether to migrate accounts.",
						},
						"is_override_root": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Whether to use the Root account in the source database to overwrite that in the target database. Valid values: false, true. For database/table or structural migration, you should specify false. Note that this parameter takes effect only for OldDTS.",
						},
						"is_dst_read_only": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Whether to set the target database to read-only during migration, which takes effect only for MySQL databases. Valid values: true, false. Default value: false.",
						},
						"extra_attr": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Additional information. You can set additional parameters for certain database types.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"key": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Option key.",
									},
									"value": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Option value.",
									},
								},
							},
						},
					},
				},
			},

			"src_info": {
				Required:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Source instance information.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"region": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Instance region.",
						},
						"access_type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Instances network access type. Valid values: extranet (public network); ipv6 (public IPv6); cvm (self-build on CVM); dcg (Direct Connect); vpncloud (VPN access); cdb (database); ccn (CCN); intranet (intranet); vpc (VPC). Note that the valid values are subject to the current link.",
						},
						"database_type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Database type, such as mysql, redis, mongodb, postgresql, mariadb, and percona.",
						},
						"node_type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Node type, empty or simple indicates a general node, cluster indicates a cluster node; for mongo services, valid values: replicaset (mongodb replica set), standalone (mongodb single node), cluster (mongodb cluster); for redis instances, valid values: empty or simple (single node), cluster (cluster), cluster-cache (cache cluster), cluster-proxy (proxy cluster).",
						},
						"info": {
							Type:        schema.TypeList,
							Required:    true,
							Description: "Database information.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"role": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Node role in a distributed database, such as the mongos node in MongoDB.",
									},
									"db_kernel": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Kernel version, such as the different kernel versions of MariaDB.",
									},
									"host": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Instance IP address, which is required for the following access types: public network, Direct Connect, VPN, CCN, intranet, and VPC.",
									},
									"port": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "Instance port, which is required for the following access types: public network, self-build on CVM, Direct Connect, VPN, CCN, intranet, and VPC.",
									},
									"user": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Instance username.",
									},
									"password": {
										Type:        schema.TypeString,
										Optional:    true,
										Sensitive:   true,
										Description: "Instance password.",
									},
									"cvm_instance_id": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Short CVM instance ID in the format of ins-olgl39y8, which is required if the access type is cvm. It is the same as the instance ID displayed in the CVM console.",
									},
									"uniq_vpn_gw_id": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "VPN gateway ID in the format of vpngw-9ghexg7q, which is required if the access type is vpncloud.",
									},
									"uniq_dcg_id": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Direct Connect gateway ID in the format of dcg-0rxtqqxb, which is required if the access type is dcg.",
									},
									"instance_id": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Database instance ID in the format of cdb-powiqx8q, which is required if the access type is cdb.",
									},
									"ccn_gw_id": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "CCN instance ID such as ccn-afp6kltc.",
									},
									"vpc_id": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "VPC ID in the format of vpc-92jblxto, which is required if the access type is vpc, vpncloud, ccn, or dcg.",
									},
									"subnet_id": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "ID of the subnet in the VPC in the format of subnet-3paxmkdz, which is required if the access type is vpc, vpncloud, ccn, or dcg.",
									},
									"engine_version": {
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "Database version in the format of 5.6 or 5.7, which takes effect only if the instance is an RDS instance. Default value: 5.6.",
									},
									"account": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Instance account.",
									},
									"account_role": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "The role used for cross-account migration, which can contain [a-zA-Z0-9-_]+.",
									},
									"account_mode": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "The account to which the resource belongs. Valid values: empty or self (the current account); other (another account).",
									},
									"tmp_secret_id": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Temporary SecretId, you can obtain the temporary key by GetFederationToken.",
									},
									"tmp_secret_key": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Temporary SecretKey, you can obtain the temporary key by GetFederationToken.",
									},
									"tmp_token": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Temporary token, you can obtain the temporary key by GetFederationToken.",
									},
								},
							},
						},
						"supplier": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Instance service provider, such as `aliyun` and `others`.",
						},
						"extra_attr": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "For MongoDB, you can define the following parameters: ['AuthDatabase':'admin', 'AuthFlag': '1', 'AuthMechanism':'SCRAM-SHA-1'].",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"key": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Option key.",
									},
									"value": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Option value.",
									},
								},
							},
						},
					},
				},
			},

			"dst_info": {
				Required:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Target database information.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"region": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Instance region.",
						},
						"access_type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Instances network access type. Valid values: extranet (public network); ipv6 (public IPv6); cvm (self-build on CVM); dcg (Direct Connect); vpncloud (VPN access); cdb (database); ccn (CCN); intranet (intranet); vpc (VPC). Note that the valid values are subject to the current link.",
						},
						"database_type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Database type, such as mysql, redis, mongodb, postgresql, mariadb, and percona.",
						},
						"node_type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Node type, empty or simple indicates a general node, cluster indicates a cluster node; for mongo services, valid values: replicaset (mongodb replica set), standalone (mongodb single node), cluster (mongodb cluster); for redis instances, valid values: empty or simple (single node), cluster (cluster), cluster-cache (cache cluster), cluster-proxy (proxy cluster).",
						},
						"info": {
							Type:        schema.TypeList,
							Required:    true,
							Description: "Database information.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"role": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Node role in a distributed database, such as the mongos node in MongoDB.",
									},
									"db_kernel": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Kernel version, such as the different kernel versions of MariaDB.",
									},
									"host": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Instance IP address, which is required for the following access types: public network, Direct Connect, VPN, CCN, intranet, and VPC.",
									},
									"port": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "Instance port, which is required for the following access types: public network, self-build on CVM, Direct Connect, VPN, CCN, intranet, and VPC.",
									},
									"user": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Instance username.",
									},
									"password": {
										Type:        schema.TypeString,
										Optional:    true,
										Sensitive:   true,
										Description: "Instance password.",
									},
									"cvm_instance_id": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Short CVM instance ID in the format of ins-olgl39y8, which is required if the access type is cvm. It is the same as the instance ID displayed in the CVM console.",
									},
									"uniq_vpn_gw_id": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "VPN gateway ID in the format of vpngw-9ghexg7q, which is required if the access type is vpncloud.",
									},
									"uniq_dcg_id": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Direct Connect gateway ID in the format of dcg-0rxtqqxb, which is required if the access type is dcg.",
									},
									"instance_id": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Database instance ID in the format of cdb-powiqx8q, which is required if the access type is cdb.",
									},
									"ccn_gw_id": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "CCN instance ID such as ccn-afp6kltc.",
									},
									"vpc_id": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "VPC ID in the format of vpc-92jblxto, which is required if the access type is vpc, vpncloud, ccn, or dcg.",
									},
									"subnet_id": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "ID of the subnet in the VPC in the format of subnet-3paxmkdz, which is required if the access type is vpc, vpncloud, ccn, or dcg.",
									},
									"engine_version": {
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "Database version in the format of 5.6 or 5.7, which takes effect only if the instance is an RDS instance. Default value: 5.6.",
									},
									"account": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Instance account.",
									},
									"account_role": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "The role used for cross-account migration, which can contain [a-zA-Z0-9-_]+.",
									},
									"account_mode": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "The account to which the resource belongs. Valid values: empty or self (the current account); other (another account).",
									},
									"tmp_secret_id": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Temporary SecretId, you can obtain the temporary key by GetFederationToken.",
									},
									"tmp_secret_key": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Temporary SecretKey, you can obtain the temporary key by GetFederationToken.",
									},
									"tmp_token": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Temporary token, you can obtain the temporary key by GetFederationToken.",
									},
								},
							},
						},
						"supplier": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Instance service provider, such as `aliyun` and `others`.",
						},
						"extra_attr": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "For MongoDB, you can define the following parameters: ['AuthDatabase':'admin','AuthFlag': '1', 'AuthMechanism':'SCRAM-SHA-1'].",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"key": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Option key.",
									},
									"value": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Option value.",
									},
								},
							},
						},
					},
				},
			},

			"expect_run_time": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Expected start time in the format of `2006-01-02 15:04:05`, which is required if RunMode is timed.",
			},

			"auto_retry_time_range_minutes": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "The automatic retry time period can be set from 5 to 720 minutes, with 0 indicating no retry.",
			},
		},
	}
}

func resourceTencentCloudDtsMigrateJobCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dts_migrate_job.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		tcClient  = meta.(tccommon.ProviderMeta).GetAPIV3Conn()
		service   = DtsService{client: tcClient}
		conf      *resource.StateChangeConf
		serviceId string
	)

	if v, ok := d.GetOk("service_id"); ok {
		serviceId = v.(string)
	}

	// case "modify":
	err := handleModifyMigrate(d, tcClient, logId, serviceId)
	if err != nil {
		return err
	}

	conf = tccommon.BuildStateChangeConf([]string{}, []string{"created"}, 3*tccommon.ReadRetryTimeout, time.Second, service.DtsMigrateJobStateRefreshFunc(serviceId, []string{}))
	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	// case "check":
	err = handleCheckMigrate(d, tcClient, logId, serviceId)
	if err != nil {
		return err
	}

	conf = tccommon.BuildStateChangeConf([]string{}, []string{"checkPass", "checkNotPass"}, 3*tccommon.ReadRetryTimeout, time.Second, service.DtsMigrateCheckConfigStateRefreshFunc(serviceId, []string{}))
	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	d.SetId(serviceId)
	return resourceTencentCloudDtsMigrateJobRead(d, meta)
}

func resourceTencentCloudDtsMigrateJobUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dts_migrate_job.update")()
	defer tccommon.InconsistentCheck(d, meta)()
	logId := tccommon.GetLogId(tccommon.ContextNil)

	log.Printf("[DEBUG]%s tencentcloud_dts_migrate_job.update in. id:[%s]\n", logId, d.Id())

	return resourceTencentCloudDtsMigrateJobCreate(d, meta)
}

func resourceTencentCloudDtsMigrateJobRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dts_migrate_job.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := DtsService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	jobId := d.Id()
	log.Printf("[DEBUG]%s tencentcloud_dts_migrate_job.read trying to call DescribeDtsMigrateJobById. jobId:[%s]\n", logId, jobId)
	migrateJob, err := service.DescribeDtsMigrateJobById(ctx, jobId)
	if err != nil {
		return err
	}

	if migrateJob == nil {
		d.SetId("")
		return fmt.Errorf("resource `track` %s does not exist", d.Id())
	}

	if migrateJob.JobId != nil {
		_ = d.Set("service_id", migrateJob.JobId)
	}

	if migrateJob.Status != nil {
		_ = d.Set("status", migrateJob.Status)
	}

	// for modify operation
	if migrateJob.RunMode != nil {
		_ = d.Set("run_mode", migrateJob.RunMode)
	}

	if migrateJob.MigrateOption != nil {
		migrateOptionMap := make(map[string]interface{})

		if migrateJob.MigrateOption.DatabaseTable != nil {
			databaseTableMap := make(map[string]interface{})

			if migrateJob.MigrateOption.DatabaseTable.ObjectMode != nil {
				databaseTableMap["object_mode"] = migrateJob.MigrateOption.DatabaseTable.ObjectMode
			}

			if migrateJob.MigrateOption.DatabaseTable.Databases != nil {
				databasesList := make([]interface{}, 0, len(migrateJob.MigrateOption.DatabaseTable.Databases))
				for _, databases := range migrateJob.MigrateOption.DatabaseTable.Databases {
					databasesMap := make(map[string]interface{})

					if databases.DbName != nil {
						databasesMap["db_name"] = databases.DbName
					}

					if databases.NewDbName != nil {
						databasesMap["new_db_name"] = databases.NewDbName
					}

					if databases.SchemaName != nil {
						databasesMap["schema_name"] = databases.SchemaName
					}

					if databases.NewSchemaName != nil {
						databasesMap["new_schema_name"] = databases.NewSchemaName
					}

					if databases.DBMode != nil {
						databasesMap["db_mode"] = databases.DBMode
					}

					if databases.SchemaMode != nil {
						databasesMap["schema_mode"] = databases.SchemaMode
					}

					if databases.TableMode != nil {
						databasesMap["table_mode"] = databases.TableMode
					}

					if databases.Tables != nil {
						tablesList := make([]interface{}, 0, len(databases.Tables))
						for _, tables := range databases.Tables {
							tablesMap := make(map[string]interface{})

							if tables.TableName != nil {
								tablesMap["table_name"] = tables.TableName
							}

							if tables.NewTableName != nil {
								tablesMap["new_table_name"] = tables.NewTableName
							}

							if tables.TmpTables != nil {
								tablesMap["tmp_tables"] = tables.TmpTables
							}

							if tables.TableEditMode != nil {
								tablesMap["table_edit_mode"] = tables.TableEditMode
							}

							tablesList = append(tablesList, tablesMap)
						}

						databasesMap["tables"] = tablesList
					}

					if databases.ViewMode != nil {
						databasesMap["view_mode"] = databases.ViewMode
					}

					if databases.Views != nil {
						viewsList := make([]interface{}, 0, len(databases.Views))
						for _, views := range databases.Views {
							viewsMap := make(map[string]interface{})

							if views.ViewName != nil {
								viewsMap["view_name"] = views.ViewName
							}

							if views.NewViewName != nil {
								viewsMap["new_view_name"] = views.NewViewName
							}

							viewsList = append(viewsList, viewsMap)
						}

						databasesMap["views"] = viewsList
					}

					if databases.RoleMode != nil {
						databasesMap["role_mode"] = databases.RoleMode
					}

					if databases.Roles != nil {
						rolesList := make([]interface{}, 0, len(databases.Roles))
						for _, roles := range databases.Roles {
							rolesMap := make(map[string]interface{})

							if roles.RoleName != nil {
								rolesMap["role_name"] = roles.RoleName
							}

							if roles.NewRoleName != nil {
								rolesMap["new_role_name"] = roles.NewRoleName
							}

							rolesList = append(rolesList, rolesMap)
						}

						databasesMap["roles"] = rolesList
					}

					if databases.FunctionMode != nil {
						databasesMap["function_mode"] = databases.FunctionMode
					}

					if databases.TriggerMode != nil {
						databasesMap["trigger_mode"] = databases.TriggerMode
					}

					if databases.EventMode != nil {
						databasesMap["event_mode"] = databases.EventMode
					}

					if databases.ProcedureMode != nil {
						databasesMap["procedure_mode"] = databases.ProcedureMode
					}

					log.Printf("[DEBUG]%s read databases.Functions:[%v],len[%d]", logId, databases.Functions, len(databases.Functions))
					for _, fun := range databases.Functions {
						log.Printf("[DEBUG]%s read databases.Functions: iterate fun:[%s]", logId, *fun)
					}

					if databases.Functions != nil {
						databasesMap["functions"] = databases.Functions
						log.Printf("[DEBUG]%s read databases.Functions: i'm in. databasesMap:[%v]", logId, databasesMap["functions"])
					}

					if databases.Procedures != nil {
						databasesMap["procedures"] = databases.Procedures
					}

					if databases.Events != nil {
						databasesMap["events"] = databases.Events
					}

					if databases.Triggers != nil {
						databasesMap["triggers"] = databases.Triggers
					}

					databasesList = append(databasesList, databasesMap)
				}

				// databaseTableMap["databases"] = []interface{}{databasesList}
				databaseTableMap["databases"] = databasesList
			}

			if migrateJob.MigrateOption.DatabaseTable.AdvancedObjects != nil {
				databaseTableMap["advanced_objects"] = migrateJob.MigrateOption.DatabaseTable.AdvancedObjects
			}

			migrateOptionMap["database_table"] = []interface{}{databaseTableMap}
		}

		if migrateJob.MigrateOption.MigrateType != nil {
			migrateOptionMap["migrate_type"] = migrateJob.MigrateOption.MigrateType
		}

		log.Printf("[DEBUG]%s read  migrateJob.MigrateOption.Consistency:[%v]", logId, migrateJob.MigrateOption.Consistency)
		if migrateJob.MigrateOption.Consistency != nil {
			consistencyMap := make(map[string]interface{})

			mode := migrateJob.MigrateOption.Consistency.Mode
			if mode != nil && *mode != "" {
				consistencyMap["mode"] = migrateJob.MigrateOption.Consistency.Mode
			}
			migrateOptionMap["consistency"] = []interface{}{consistencyMap}
		}

		if migrateJob.MigrateOption.IsMigrateAccount != nil {
			migrateOptionMap["is_migrate_account"] = migrateJob.MigrateOption.IsMigrateAccount
		}

		if migrateJob.MigrateOption.IsOverrideRoot != nil {
			migrateOptionMap["is_override_root"] = migrateJob.MigrateOption.IsOverrideRoot
		}

		if migrateJob.MigrateOption.IsDstReadOnly != nil {
			migrateOptionMap["is_dst_read_only"] = migrateJob.MigrateOption.IsDstReadOnly
		}

		if migrateJob.MigrateOption.ExtraAttr != nil {
			extraAttrList := make([]interface{}, 0, len(migrateJob.MigrateOption.ExtraAttr))
			for _, extraAttr := range migrateJob.MigrateOption.ExtraAttr {
				extraAttrMap := make(map[string]interface{})

				if extraAttr.Key != nil {
					extraAttrMap["key"] = extraAttr.Key
				}

				if extraAttr.Value != nil {
					extraAttrMap["value"] = extraAttr.Value
				}

				extraAttrList = append(extraAttrList, extraAttrMap)
			}

			migrateOptionMap["extra_attr"] = extraAttrList
		}

		_ = d.Set("migrate_option", []interface{}{migrateOptionMap})
	}

	if migrateJob.SrcInfo != nil {
		srcInfoMap := make(map[string]interface{})

		if migrateJob.SrcInfo.Region != nil {
			srcInfoMap["region"] = migrateJob.SrcInfo.Region
		}

		if migrateJob.SrcInfo.AccessType != nil {
			srcInfoMap["access_type"] = migrateJob.SrcInfo.AccessType
		}

		if migrateJob.SrcInfo.DatabaseType != nil {
			srcInfoMap["database_type"] = migrateJob.SrcInfo.DatabaseType
		}

		if migrateJob.SrcInfo.NodeType != nil {
			srcInfoMap["node_type"] = migrateJob.SrcInfo.NodeType
		}

		if migrateJob.SrcInfo.Info != nil {
			infoList := make([]interface{}, 0, len(migrateJob.SrcInfo.Info))
			for i, info := range migrateJob.SrcInfo.Info {
				infoMap := make(map[string]interface{})

				if info.Password == nil || *info.Password == "" {
					//reset password
					key := fmt.Sprintf("src_info.0.info.%v.password", i)
					if v, ok := d.GetOk(key); ok {
						infoMap["password"] = helper.String(v.(string))
						log.Printf("[DEBUG]%s set src_info.0.info.%v.password:[key:%s]", logId, i, key)
					}
				} else {
					infoMap["password"] = info.Password
				}

				if info.Role != nil {
					infoMap["role"] = info.Role
				}

				if info.DbKernel != nil {
					infoMap["db_kernel"] = info.DbKernel
				}

				if info.Host != nil {
					infoMap["host"] = info.Host
				}

				if info.Port != nil {
					infoMap["port"] = info.Port
				}

				if info.User != nil {
					infoMap["user"] = info.User
				}

				if info.CvmInstanceId != nil {
					infoMap["cvm_instance_id"] = info.CvmInstanceId
				}

				if info.UniqVpnGwId != nil {
					infoMap["uniq_vpn_gw_id"] = info.UniqVpnGwId
				}

				if info.UniqDcgId != nil {
					infoMap["uniq_dcg_id"] = info.UniqDcgId
				}

				if info.InstanceId != nil {
					infoMap["instance_id"] = info.InstanceId
				}

				if info.CcnGwId != nil {
					infoMap["ccn_gw_id"] = info.CcnGwId
				}

				if info.VpcId != nil {
					infoMap["vpc_id"] = info.VpcId
				}

				if info.SubnetId != nil {
					infoMap["subnet_id"] = info.SubnetId
				}

				log.Printf("[DEBUG]%s read  migrateJob.SrcInfo.Info.EngineVersion:[%v,%s]", logId, info.EngineVersion, *info.EngineVersion)
				if info.EngineVersion != nil {
					infoMap["engine_version"] = info.EngineVersion
				}

				if info.Account != nil {
					infoMap["account"] = info.Account
				}

				if info.AccountRole != nil {
					infoMap["account_role"] = info.AccountRole
				}

				if info.AccountMode != nil {
					infoMap["account_mode"] = info.AccountMode
				}

				if info.TmpSecretId != nil {
					infoMap["tmp_secret_id"] = info.TmpSecretId
				}

				if info.TmpSecretKey != nil {
					infoMap["tmp_secret_key"] = info.TmpSecretKey
				}

				if info.TmpToken != nil {
					infoMap["tmp_token"] = info.TmpToken
				}

				infoList = append(infoList, infoMap)
			}

			srcInfoMap["info"] = infoList
		}

		if migrateJob.SrcInfo.Supplier != nil {
			srcInfoMap["supplier"] = migrateJob.SrcInfo.Supplier
		}

		if migrateJob.SrcInfo.ExtraAttr != nil {
			extraAttrList := make([]interface{}, 0, len(migrateJob.SrcInfo.ExtraAttr))
			for _, extraAttr := range migrateJob.SrcInfo.ExtraAttr {
				extraAttrMap := make(map[string]interface{})

				if extraAttr.Key != nil {
					extraAttrMap["key"] = extraAttr.Key
				}

				if extraAttr.Value != nil {
					extraAttrMap["value"] = extraAttr.Value
				}

				extraAttrList = append(extraAttrList, extraAttrMap)
			}

			srcInfoMap["extra_attr"] = extraAttrList
		}

		_ = d.Set("src_info", []interface{}{srcInfoMap})
	}

	if migrateJob.DstInfo != nil {
		dstInfoMap := make(map[string]interface{})
		if migrateJob.DstInfo.Region != nil {
			dstInfoMap["region"] = migrateJob.DstInfo.Region
		}

		if migrateJob.DstInfo.AccessType != nil {
			dstInfoMap["access_type"] = migrateJob.DstInfo.AccessType
		}

		if migrateJob.DstInfo.DatabaseType != nil {
			dstInfoMap["database_type"] = migrateJob.DstInfo.DatabaseType
		}

		if migrateJob.DstInfo.NodeType != nil {
			dstInfoMap["node_type"] = migrateJob.DstInfo.NodeType
		}

		log.Printf("[DEBUG]%s read migrateJob.DstInfo.Info :[%v], len:[%v]", logId, migrateJob.DstInfo.Info, len(migrateJob.DstInfo.Info))
		if migrateJob.DstInfo.Info != nil {
			infoList := make([]interface{}, 0, len(migrateJob.DstInfo.Info))
			for i, info := range migrateJob.DstInfo.Info {
				infoMap := make(map[string]interface{})

				if info.Password == nil || *info.Password == "" {
					//reset password
					key := fmt.Sprintf("dst_info.0.info.%v.password", i)
					if v, ok := d.GetOk(key); ok {
						infoMap["password"] = helper.String(v.(string))
						log.Printf("[DEBUG]%s set dst_info.0.info.%v.password:[key:%s]", logId, i, key)
					}
				} else {
					infoMap["password"] = info.Password
				}

				if info.Role != nil {
					infoMap["role"] = info.Role
				}

				if info.DbKernel != nil {
					infoMap["db_kernel"] = info.DbKernel
				}

				if info.Host != nil {
					infoMap["host"] = info.Host
				}

				if info.Port != nil {
					infoMap["port"] = info.Port
				}

				if info.User != nil {
					infoMap["user"] = info.User
				}

				if info.CvmInstanceId != nil {
					infoMap["cvm_instance_id"] = info.CvmInstanceId
				}

				if info.UniqVpnGwId != nil {
					infoMap["uniq_vpn_gw_id"] = info.UniqVpnGwId
				}

				if info.UniqDcgId != nil {
					infoMap["uniq_dcg_id"] = info.UniqDcgId
				}

				if info.InstanceId != nil {
					infoMap["instance_id"] = info.InstanceId
				}

				if info.CcnGwId != nil {
					infoMap["ccn_gw_id"] = info.CcnGwId
				}

				if info.VpcId != nil {
					infoMap["vpc_id"] = info.VpcId
				}

				if info.SubnetId != nil {
					infoMap["subnet_id"] = info.SubnetId
				}

				log.Printf("[DEBUG]%s read  migrateJob.DstInfo.Info.EngineVersion:[%v,%s]", logId, info.EngineVersion, *info.EngineVersion)
				if d.HasChange("engine_version") {
					if info.EngineVersion != nil {
						infoMap["engine_version"] = info.EngineVersion
					}
				}

				if info.Account != nil {
					infoMap["account"] = info.Account
				}

				if info.AccountRole != nil {
					infoMap["account_role"] = info.AccountRole
				}

				if info.AccountMode != nil {
					infoMap["account_mode"] = info.AccountMode
				}

				if info.TmpSecretId != nil {
					infoMap["tmp_secret_id"] = info.TmpSecretId
				}

				if info.TmpSecretKey != nil {
					infoMap["tmp_secret_key"] = info.TmpSecretKey
				}

				if info.TmpToken != nil {
					infoMap["tmp_token"] = info.TmpToken
				}

				infoList = append(infoList, infoMap)
			}

			dstInfoMap["info"] = infoList
		}

		log.Printf("[DEBUG]%s read migrateJob.DstInfo.Supplier :[%s]", logId, *migrateJob.DstInfo.Supplier)
		if migrateJob.DstInfo.Supplier != nil {
			dstInfoMap["supplier"] = migrateJob.DstInfo.Supplier
		}

		if migrateJob.DstInfo.ExtraAttr != nil {
			extraAttrList := make([]interface{}, 0, len(migrateJob.DstInfo.ExtraAttr))
			for _, extraAttr := range migrateJob.DstInfo.ExtraAttr {
				extraAttrMap := make(map[string]interface{})

				if extraAttr.Key != nil {
					extraAttrMap["key"] = extraAttr.Key
				}

				if extraAttr.Value != nil {
					extraAttrMap["value"] = extraAttr.Value
				}

				extraAttrList = append(extraAttrList, extraAttrMap)
			}

			dstInfoMap["extra_attr"] = extraAttrList
		}

		_ = d.Set("dst_info", []interface{}{dstInfoMap})
	}

	if migrateJob.ExpectRunTime != nil {
		_ = d.Set("expect_run_time", migrateJob.ExpectRunTime)
	}

	return nil
}

func handleModifyMigrate(d *schema.ResourceData, tcClient *connectivity.TencentCloudClient, logId, jobId string) error {
	configMigrationJobRequest := dts.NewModifyMigrationJobRequest()
	configMigrationJobRequest.JobId = helper.String(jobId)

	if v, ok := d.GetOk("run_mode"); ok {
		configMigrationJobRequest.RunMode = helper.String(v.(string))
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "migrate_option"); ok {
		migrateOption := dts.MigrateOption{}
		if databaseTableMap, ok := helper.InterfaceToMap(dMap, "database_table"); ok {
			databaseTableObject := dts.DatabaseTableObject{}
			if v, ok := databaseTableMap["object_mode"]; ok && v.(string) != "" {
				databaseTableObject.ObjectMode = helper.String(v.(string))
			}
			if v, ok := databaseTableMap["databases"]; ok {
				for _, item := range v.([]interface{}) {
					databasesMap := item.(map[string]interface{})
					dBItem := dts.DBItem{}
					if v, ok := databasesMap["db_name"]; ok && v.(string) != "" {
						dBItem.DbName = helper.String(v.(string))
					}
					if v, ok := databasesMap["new_db_name"]; ok && v.(string) != "" {
						dBItem.NewDbName = helper.String(v.(string))
					}
					if v, ok := databasesMap["schema_name"]; ok && v.(string) != "" {
						dBItem.SchemaName = helper.String(v.(string))
					}
					if v, ok := databasesMap["new_schema_name"]; ok && v.(string) != "" {
						dBItem.NewSchemaName = helper.String(v.(string))
					}
					if v, ok := databasesMap["db_mode"]; ok && v.(string) != "" {
						dBItem.DBMode = helper.String(v.(string))
					}
					if v, ok := databasesMap["schema_mode"]; ok && v.(string) != "" {
						dBItem.SchemaMode = helper.String(v.(string))
					}
					if v, ok := databasesMap["table_mode"]; ok && v.(string) != "" {
						dBItem.TableMode = helper.String(v.(string))
					}
					if v, ok := databasesMap["tables"]; ok {
						for _, item := range v.([]interface{}) {
							tablesMap := item.(map[string]interface{})
							tableItem := dts.TableItem{}
							if v, ok := tablesMap["table_name"]; ok && v.(string) != "" {
								tableItem.TableName = helper.String(v.(string))
							}
							if v, ok := tablesMap["new_table_name"]; ok && v.(string) != "" {
								tableItem.NewTableName = helper.String(v.(string))
							}
							if v, ok := tablesMap["tmp_tables"]; ok {
								tmpTablesSet := v.(*schema.Set).List()
								for i := range tmpTablesSet {
									tmpTables := tmpTablesSet[i].(string)
									tableItem.TmpTables = append(tableItem.TmpTables, &tmpTables)
								}
							}
							if v, ok := tablesMap["table_edit_mode"]; ok && v.(string) != "" {
								tableItem.TableEditMode = helper.String(v.(string))
							}
							dBItem.Tables = append(dBItem.Tables, &tableItem)
						}
					}
					if v, ok := databasesMap["view_mode"]; ok && v.(string) != "" {
						dBItem.ViewMode = helper.String(v.(string))
					}
					if v, ok := databasesMap["views"]; ok {
						for _, item := range v.([]interface{}) {
							viewsMap := item.(map[string]interface{})
							viewItem := dts.ViewItem{}
							if v, ok := viewsMap["view_name"]; ok && v.(string) != "" {
								viewItem.ViewName = helper.String(v.(string))
							}
							if v, ok := viewsMap["new_view_name"]; ok && v.(string) != "" {
								viewItem.NewViewName = helper.String(v.(string))
							}
							dBItem.Views = append(dBItem.Views, &viewItem)
						}
					}
					if v, ok := databasesMap["role_mode"]; ok && v.(string) != "" {
						dBItem.RoleMode = helper.String(v.(string))
					}
					if v, ok := databasesMap["roles"]; ok {
						for _, item := range v.([]interface{}) {
							rolesMap := item.(map[string]interface{})
							roleItem := dts.RoleItem{}
							if v, ok := rolesMap["role_name"]; ok && v.(string) != "" {
								roleItem.RoleName = helper.String(v.(string))
							}
							if v, ok := rolesMap["new_role_name"]; ok && v.(string) != "" {
								roleItem.NewRoleName = helper.String(v.(string))
							}
							dBItem.Roles = append(dBItem.Roles, &roleItem)
						}
					}
					if v, ok := databasesMap["function_mode"]; ok && v.(string) != "" {
						dBItem.FunctionMode = helper.String(v.(string))
					}
					if v, ok := databasesMap["trigger_mode"]; ok && v.(string) != "" {
						dBItem.TriggerMode = helper.String(v.(string))
					}
					if v, ok := databasesMap["event_mode"]; ok && v.(string) != "" {
						dBItem.EventMode = helper.String(v.(string))
					}
					if v, ok := databasesMap["procedure_mode"]; ok && v.(string) != "" {
						dBItem.ProcedureMode = helper.String(v.(string))
					}
					log.Printf("[DEBUG]%s modify databases.Functions: databasesMap[\"functions\"]:[%v]", logId, databasesMap["functions"])
					if v, ok := databasesMap["functions"]; ok {
						functionsSet := v.(*schema.Set).List()
						log.Printf("[DEBUG]%s modify databases.Functions: i'm in. functionsSet:[%v]", logId, functionsSet)
						for _, funcc := range functionsSet {
							functions := funcc.(*string)
							dBItem.Functions = append(dBItem.Functions, functions)
							log.Printf("[DEBUG]%s modify databases.Functions: iterate functions:[%s]", logId, *functions)
						}
					}
					if v, ok := databasesMap["procedures"]; ok {
						proceduresSet := v.(*schema.Set).List()
						for _, proc := range proceduresSet {
							procedures := proc.(string)
							dBItem.Procedures = append(dBItem.Procedures, &procedures)
						}
					}
					if v, ok := databasesMap["events"]; ok {
						eventsSet := v.(*schema.Set).List()
						for i := range eventsSet {
							events := eventsSet[i].(string)
							dBItem.Events = append(dBItem.Events, &events)
						}
					}
					if v, ok := databasesMap["triggers"]; ok {
						triggersSet := v.(*schema.Set).List()
						for i := range triggersSet {
							triggers := triggersSet[i].(string)
							dBItem.Triggers = append(dBItem.Triggers, &triggers)
						}
					}
					databaseTableObject.Databases = append(databaseTableObject.Databases, &dBItem)
				}
			}
			if v, ok := databaseTableMap["advanced_objects"]; ok {
				advancedObjectsSet := v.(*schema.Set).List()
				for i := range advancedObjectsSet {
					advancedObjects := advancedObjectsSet[i].(string)
					databaseTableObject.AdvancedObjects = append(databaseTableObject.AdvancedObjects, &advancedObjects)
				}
			}
			migrateOption.DatabaseTable = &databaseTableObject
		}
		if v, ok := dMap["migrate_type"]; ok && v.(string) != "" {
			migrateOption.MigrateType = helper.String(v.(string))
		}
		log.Printf("[DEBUG]%s update  migrateJob.MigrateOption.Consistency dMap(consistency):[%v]", logId, dMap["consistency"])
		if consistencyMap, ok := helper.InterfaceToMap(dMap, "consistency"); ok {
			log.Printf("[DEBUG]%s update  migrateJob.MigrateOption.Consistency:[%v]", logId, consistencyMap)
			consistencyOption := dts.ConsistencyOption{}
			if v, ok := consistencyMap["mode"]; ok && v.(string) != "" {
				consistencyOption.Mode = helper.String(v.(string))
			}
			migrateOption.Consistency = &consistencyOption
		}
		if v, ok := dMap["is_migrate_account"]; ok {
			migrateOption.IsMigrateAccount = helper.Bool(v.(bool))
		}
		if v, ok := dMap["is_override_root"]; ok {
			migrateOption.IsOverrideRoot = helper.Bool(v.(bool))
		}
		if v, ok := dMap["is_dst_read_only"]; ok {
			migrateOption.IsDstReadOnly = helper.Bool(v.(bool))
		}
		if v, ok := dMap["extra_attr"]; ok {
			for _, item := range v.([]interface{}) {
				extraAttrMap := item.(map[string]interface{})
				keyValuePairOption := dts.KeyValuePairOption{}
				if v, ok := extraAttrMap["key"]; ok && v.(string) != "" {
					keyValuePairOption.Key = helper.String(v.(string))
				}
				if v, ok := extraAttrMap["value"]; ok && v.(string) != "" {
					keyValuePairOption.Value = helper.String(v.(string))
				}
				migrateOption.ExtraAttr = append(migrateOption.ExtraAttr, &keyValuePairOption)
			}
		}
		configMigrationJobRequest.MigrateOption = &migrateOption
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "src_info"); ok {
		dBEndpointInfo := dts.DBEndpointInfo{}
		if v, ok := dMap["region"]; ok && v.(string) != "" {
			dBEndpointInfo.Region = helper.String(v.(string))
		}
		if v, ok := dMap["access_type"]; ok && v.(string) != "" {
			dBEndpointInfo.AccessType = helper.String(v.(string))
		}
		if v, ok := dMap["database_type"]; ok && v.(string) != "" {
			dBEndpointInfo.DatabaseType = helper.String(v.(string))
		}
		if v, ok := dMap["node_type"]; ok && v.(string) != "" {
			dBEndpointInfo.NodeType = helper.String(v.(string))
		}
		if v, ok := dMap["info"]; ok {
			for _, item := range v.([]interface{}) {
				srcInfoMap := item.(map[string]interface{})
				dBInfo := dts.DBInfo{}
				if v, ok := srcInfoMap["role"]; ok && v.(string) != "" {
					dBInfo.Role = helper.String(v.(string))
				}
				if v, ok := srcInfoMap["db_kernel"]; ok && v.(string) != "" {
					dBInfo.DbKernel = helper.String(v.(string))
				}
				if v, ok := srcInfoMap["host"]; ok && v.(string) != "" {
					dBInfo.Host = helper.String(v.(string))
				}
				if v, ok := srcInfoMap["port"]; ok {
					dBInfo.Port = helper.IntUint64(v.(int))
				}
				if v, ok := srcInfoMap["user"]; ok && v.(string) != "" {
					dBInfo.User = helper.String(v.(string))
				}
				if v, ok := srcInfoMap["password"]; ok && v.(string) != "" {
					dBInfo.Password = helper.String(v.(string))
				}
				if v, ok := srcInfoMap["cvm_instance_id"]; ok && v.(string) != "" {
					dBInfo.CvmInstanceId = helper.String(v.(string))
				}
				if v, ok := srcInfoMap["uniq_vpn_gw_id"]; ok && v.(string) != "" {
					dBInfo.UniqVpnGwId = helper.String(v.(string))
				}
				if v, ok := srcInfoMap["uniq_dcg_id"]; ok && v.(string) != "" {
					dBInfo.UniqDcgId = helper.String(v.(string))
				}
				if v, ok := srcInfoMap["instance_id"]; ok && v.(string) != "" {
					dBInfo.InstanceId = helper.String(v.(string))
				}
				if v, ok := srcInfoMap["ccn_gw_id"]; ok && v.(string) != "" {
					dBInfo.CcnGwId = helper.String(v.(string))
				}
				if v, ok := srcInfoMap["vpc_id"]; ok && v.(string) != "" {
					dBInfo.VpcId = helper.String(v.(string))
				}
				if v, ok := srcInfoMap["subnet_id"]; ok && v.(string) != "" {
					dBInfo.SubnetId = helper.String(v.(string))
				}
				if v, ok := srcInfoMap["engine_version"]; ok && v.(string) != "" {
					dBInfo.EngineVersion = helper.String(v.(string))
				}
				if v, ok := srcInfoMap["account"]; ok && v.(string) != "" {
					dBInfo.Account = helper.String(v.(string))
				}
				if v, ok := srcInfoMap["account_role"]; ok && v.(string) != "" {
					dBInfo.AccountRole = helper.String(v.(string))
				}
				if v, ok := srcInfoMap["account_mode"]; ok && v.(string) != "" {
					dBInfo.AccountMode = helper.String(v.(string))
				}
				if v, ok := srcInfoMap["tmp_secret_id"]; ok && v.(string) != "" {
					dBInfo.TmpSecretId = helper.String(v.(string))
				}
				if v, ok := srcInfoMap["tmp_secret_key"]; ok && v.(string) != "" {
					dBInfo.TmpSecretKey = helper.String(v.(string))
				}
				if v, ok := srcInfoMap["tmp_token"]; ok && v.(string) != "" {
					dBInfo.TmpToken = helper.String(v.(string))
				}
				dBEndpointInfo.Info = append(dBEndpointInfo.Info, &dBInfo)
			}
		}
		if v, ok := dMap["supplier"]; ok && v.(string) != "" {
			dBEndpointInfo.Supplier = helper.String(v.(string))
		}
		if v, ok := dMap["extra_attr"]; ok {
			for _, item := range v.([]interface{}) {
				extraAttrMap := item.(map[string]interface{})
				keyValuePairOption := dts.KeyValuePairOption{}
				if v, ok := extraAttrMap["key"]; ok && v.(string) != "" {
					keyValuePairOption.Key = helper.String(v.(string))
				}
				if v, ok := extraAttrMap["value"]; ok && v.(string) != "" {
					keyValuePairOption.Value = helper.String(v.(string))
				}
				dBEndpointInfo.ExtraAttr = append(dBEndpointInfo.ExtraAttr, &keyValuePairOption)
			}
		}
		configMigrationJobRequest.SrcInfo = &dBEndpointInfo
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "dst_info"); ok {
		dBEndpointInfo := dts.DBEndpointInfo{}
		if v, ok := dMap["region"]; ok && v.(string) != "" {
			dBEndpointInfo.Region = helper.String(v.(string))
		}
		if v, ok := dMap["access_type"]; ok && v.(string) != "" {
			dBEndpointInfo.AccessType = helper.String(v.(string))
		}
		if v, ok := dMap["database_type"]; ok && v.(string) != "" {
			dBEndpointInfo.DatabaseType = helper.String(v.(string))
		}
		if v, ok := dMap["node_type"]; ok && v.(string) != "" {
			dBEndpointInfo.NodeType = helper.String(v.(string))
		}
		if v, ok := dMap["info"]; ok {
			for _, item := range v.([]interface{}) {
				dstInfoMap := item.(map[string]interface{})
				dBInfo := dts.DBInfo{}
				if v, ok := dstInfoMap["role"]; ok && v.(string) != "" {
					dBInfo.Role = helper.String(v.(string))
				}
				if v, ok := dstInfoMap["db_kernel"]; ok && v.(string) != "" {
					dBInfo.DbKernel = helper.String(v.(string))
				}
				if v, ok := dstInfoMap["host"]; ok && v.(string) != "" {
					dBInfo.Host = helper.String(v.(string))
				}
				if v, ok := dstInfoMap["port"]; ok {
					dBInfo.Port = helper.IntUint64(v.(int))
				}
				if v, ok := dstInfoMap["user"]; ok && v.(string) != "" {
					dBInfo.User = helper.String(v.(string))
				}
				if v, ok := dstInfoMap["password"]; ok && v.(string) != "" {
					dBInfo.Password = helper.String(v.(string))
				}
				if v, ok := dstInfoMap["cvm_instance_id"]; ok && v.(string) != "" {
					dBInfo.CvmInstanceId = helper.String(v.(string))
				}
				if v, ok := dstInfoMap["uniq_vpn_gw_id"]; ok && v.(string) != "" {
					dBInfo.UniqVpnGwId = helper.String(v.(string))
				}
				if v, ok := dstInfoMap["uniq_dcg_id"]; ok && v.(string) != "" {
					dBInfo.UniqDcgId = helper.String(v.(string))
				}
				if v, ok := dstInfoMap["instance_id"]; ok && v.(string) != "" {
					dBInfo.InstanceId = helper.String(v.(string))
				}
				if v, ok := dstInfoMap["ccn_gw_id"]; ok && v.(string) != "" {
					dBInfo.CcnGwId = helper.String(v.(string))
				}
				if v, ok := dstInfoMap["vpc_id"]; ok && v.(string) != "" {
					dBInfo.VpcId = helper.String(v.(string))
				}
				if v, ok := dstInfoMap["subnet_id"]; ok && v.(string) != "" {
					dBInfo.SubnetId = helper.String(v.(string))
				}
				if v, ok := dstInfoMap["engine_version"]; ok && v.(string) != "" {
					dBInfo.EngineVersion = helper.String(v.(string))
				}
				if v, ok := dstInfoMap["account"]; ok && v.(string) != "" {
					dBInfo.Account = helper.String(v.(string))
				}
				if v, ok := dstInfoMap["account_role"]; ok && v.(string) != "" {
					dBInfo.AccountRole = helper.String(v.(string))
				}
				if v, ok := dstInfoMap["account_mode"]; ok && v.(string) != "" {
					dBInfo.AccountMode = helper.String(v.(string))
				}
				if v, ok := dstInfoMap["tmp_secret_id"]; ok && v.(string) != "" {
					dBInfo.TmpSecretId = helper.String(v.(string))
				}
				if v, ok := dstInfoMap["tmp_secret_key"]; ok && v.(string) != "" {
					dBInfo.TmpSecretKey = helper.String(v.(string))
				}
				if v, ok := dstInfoMap["tmp_token"]; ok && v.(string) != "" {
					dBInfo.TmpToken = helper.String(v.(string))
				}
				dBEndpointInfo.Info = append(dBEndpointInfo.Info, &dBInfo)
			}
		}
		if v, ok := dMap["supplier"]; ok && v.(string) != "" {
			dBEndpointInfo.Supplier = helper.String(v.(string))
		}
		if v, ok := dMap["extra_attr"]; ok {
			for _, item := range v.([]interface{}) {
				extraAttrMap := item.(map[string]interface{})
				keyValuePairOption := dts.KeyValuePairOption{}
				if v, ok := extraAttrMap["key"]; ok && v.(string) != "" {
					keyValuePairOption.Key = helper.String(v.(string))
				}
				if v, ok := extraAttrMap["value"]; ok && v.(string) != "" {
					keyValuePairOption.Value = helper.String(v.(string))
				}
				dBEndpointInfo.ExtraAttr = append(dBEndpointInfo.ExtraAttr, &keyValuePairOption)
			}
		}
		configMigrationJobRequest.DstInfo = &dBEndpointInfo
	}

	if v, ok := d.GetOk("expect_run_time"); ok && v.(string) != "" {
		configMigrationJobRequest.ExpectRunTime = helper.String(v.(string))
	}

	if v, _ := d.GetOk("auto_retry_time_range_minutes"); v != nil {
		configMigrationJobRequest.AutoRetryTimeRangeMinutes = helper.IntInt64(v.(int))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := tcClient.UseDtsClient().ModifyMigrationJob(configMigrationJobRequest)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, configMigrationJobRequest.GetAction(), configMigrationJobRequest.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create dts migrateJob failed, reason:%+v", logId, err)
		return err
	}
	return nil
}

func handleCheckMigrate(d *schema.ResourceData, tcClient *connectivity.TencentCloudClient, logId, jobId string) error {
	checkMigrateJobRequest := dts.NewCreateMigrateCheckJobRequest()
	checkMigrateJobRequest.JobId = helper.String(jobId)

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := tcClient.UseDtsClient().CreateMigrateCheckJob(checkMigrateJobRequest)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, checkMigrateJobRequest.GetAction(), checkMigrateJobRequest.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s check dts migrateJob failed, reason:%+v", logId, err)
		return err
	}

	return nil
}

// reserve for future implementation
// func handleResumeMigrate(d *schema.ResourceData, tcClient *connectivity.TencentCloudClient, logId, jobId string) error {
// 	resumeMigrateJobRequest := dts.NewResumeMigrateJobRequest()
// 	resumeMigrateJobRequest.JobId = helper.String(jobId)
// 	service := DtsService{client: tcClient}

// 	if d.HasChange("resume_option") {
// 		if v, ok := d.GetOk("resume_option"); ok {
// 			resumeMigrateJobRequest.ResumeOption = helper.String(v.(string))
// 		}
// 	}

// 	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
// 		result, e := tcClient.UseDtsClient().ResumeMigrateJob(resumeMigrateJobRequest)
// 		if e != nil {
// 			return tccommon.RetryError(e)
// 		} else {
// 			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, resumeMigrateJobRequest.GetAction(), resumeMigrateJobRequest.ToJsonString(), result.ToJsonString())
// 		}
// 		return nil
// 	})
// 	if err != nil {
// 		log.Printf("[CRITAL]%s resume dts migrateJob failed, reason:%+v", logId, err)
// 		return err
// 	}

// 	conf := tccommon.BuildStateChangeConf([]string{}, []string{"readyComplete", "success", "failed"}, 3*tccommon.ReadRetryTimeout, time.Second, service.DtsMigrateJobStateRefreshFunc(jobId, []string{}))
// 	if _, e := conf.WaitForState(); e != nil {
// 		return e
// 	}

// 	return nil
// }

// func handleCompleteMigrate(d *schema.ResourceData, tcClient *connectivity.TencentCloudClient, logId, jobId string) error {
// 	completeMigrateJobRequest := dts.NewCompleteMigrateJobRequest()
// 	completeMigrateJobRequest.JobId = helper.String(jobId)
// 	service := DtsService{client: tcClient}

// 	if d.HasChange("complete_mode") {
// 		if v, ok := d.GetOk("complete_mode"); ok {
// 			completeMigrateJobRequest.CompleteMode = helper.String(v.(string))
// 		}
// 	}

// 	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
// 		result, e := tcClient.UseDtsClient().CompleteMigrateJob(completeMigrateJobRequest)
// 		if e != nil {
// 			return tccommon.RetryError(e)
// 		} else {
// 			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, completeMigrateJobRequest.GetAction(), completeMigrateJobRequest.ToJsonString(), result.ToJsonString())
// 		}
// 		return nil
// 	})
// 	if err != nil {
// 		log.Printf("[CRITAL]%s complete dts migrateJob failed, reason:%+v", logId, err)
// 		return err
// 	}

// 	conf := tccommon.BuildStateChangeConf([]string{}, []string{"success", "error", "failed"}, 3*tccommon.ReadRetryTimeout, time.Second, service.DtsMigrateJobStateRefreshFunc(jobId, []string{}))
// 	if _, e := conf.WaitForState(); e != nil {
// 		return e
// 	}

// 	return nil
// }

// func handleCompareMigrate(d *schema.ResourceData, tcClient *connectivity.TencentCloudClient, logId, jobId string) error {
// 	startCompareRequest := dts.NewStartCompareRequest()
// 	startCompareRequest.JobId = helper.String(jobId)

// 	if d.HasChange("compare_task_id") {
// 		if v, ok := d.GetOk("compare_task_id"); ok {
// 			startCompareRequest.CompareTaskId = helper.String(v.(string))
// 		}
// 	}

// 	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
// 		result, e := tcClient.UseDtsClient().StartCompare(startCompareRequest)
// 		if e != nil {
// 			return tccommon.RetryError(e)
// 		} else {
// 			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, startCompareRequest.GetAction(), startCompareRequest.ToJsonString(), result.ToJsonString())
// 		}
// 		return nil
// 	})
// 	if err != nil {
// 		log.Printf("[CRITAL]%s compare dts migrate job failed, reason:%+v", logId, err)
// 		return err
// 	}

// 	return nil
// }

func resourceTencentCloudDtsMigrateJobDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dts_migrate_job.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
