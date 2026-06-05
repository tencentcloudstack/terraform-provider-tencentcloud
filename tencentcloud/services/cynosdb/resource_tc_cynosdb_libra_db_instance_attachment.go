package cynosdb

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cynosdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cynosdb/v20190107"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudCynosdbLibraDbInstanceAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCynosdbLibraDbInstanceAttachmentCreate,
		Read:   resourceTencentCloudCynosdbLibraDbInstanceAttachmentRead,
		Update: resourceTencentCloudCynosdbLibraDbInstanceAttachmentUpdate,
		Delete: resourceTencentCloudCynosdbLibraDbInstanceAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Cluster ID.",
			},
			"zone": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Availability zone.",
			},
			"cpu": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Number of CPU cores.",
			},
			"mem": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Memory size in GB.",
			},
			"storage_size": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Disk size.",
			},
			"pay_mode": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Payment mode.",
			},
			"objects": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "Sync object list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"database_tables": {
							Type:        schema.TypeList,
							Optional:    true,
							MaxItems:    1,
							Description: "Database table information.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"migrate_db_mode": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Database migration mode.",
									},
									"databases": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "Database information list.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"db_name": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Database name.",
												},
												"migrate_table_mode": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Table migration mode.",
												},
												"tables": {
													Type:        schema.TypeList,
													Optional:    true,
													Description: "Table information list.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"table_name": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Table name.",
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
					},
				},
			},
			"port": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Port for the new RO group, value range [0, 65535).",
			},
			"goods_num": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Number of new read-only instances, value range (0, 15].",
			},
			"instance_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Instance name.",
			},
			"replicas_num": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Number of replicas.",
			},
			"instance_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Instance type.",
			},
			"storage_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Disk type.",
			},
			"auto_voucher": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Whether to automatically select vouchers: 1 yes, 0 no, default 0.",
			},
			"order_source": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Order source.",
			},
			"deal_mode": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Transaction mode: 0 - place order and pay, 1 - place order only.",
			},
			"vpc_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "VPC network ID.",
			},
			"subnet_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Subnet ID. Required if VpcId is set.",
			},
			"security_group_ids": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Security group IDs for the new read-only instance.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"libra_db_version": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Analytics engine version.",
			},
			"time_span": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Purchase duration, effective with TimeUnit.",
			},
			"time_unit": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Purchase duration unit. Options: d (day), m (month).",
			},
			"src_instance_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Source instance ID.",
			},
			"isolate_reason_types": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Isolation reason types for delete operation.",
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},
			"isolate_reason": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Isolation reason for delete operation.",
			},
			// computed attributes
			"big_deal_ids": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Big deal IDs.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"tran_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Frozen transaction ID.",
			},
			"deal_names": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Post-paid order names.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"resource_ids": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Resource ID list.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"instance_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Instance ID.",
			},
		},
	}
}

func resourceTencentCloudCynosdbLibraDbInstanceAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cynosdb_libra_db_instance_attachment.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId     = tccommon.GetLogId(tccommon.ContextNil)
		request   = cynosdb.NewAddLibraDBInstancesRequest()
		response  *cynosdb.AddLibraDBInstancesResponse
		clusterId string
	)

	if v, ok := d.GetOk("cluster_id"); ok {
		clusterId = v.(string)
		request.ClusterId = helper.String(clusterId)
	}

	if v, ok := d.GetOk("zone"); ok {
		request.Zone = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("cpu"); ok {
		request.Cpu = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("mem"); ok {
		request.Mem = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("storage_size"); ok {
		request.StorageSize = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("pay_mode"); ok {
		request.PayMode = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("objects"); ok {
		objectsList := v.([]interface{})
		if len(objectsList) > 0 {
			objectsMap := objectsList[0].(map[string]interface{})
			objects := &cynosdb.Objects{}
			if v, ok := objectsMap["database_tables"]; ok {
				dbTablesList := v.([]interface{})
				if len(dbTablesList) > 0 {
					dbTablesMap := dbTablesList[0].(map[string]interface{})
					migrateObject := &cynosdb.MigrateObject{}
					if v, ok := dbTablesMap["migrate_db_mode"]; ok && v.(string) != "" {
						migrateObject.MigrateDBMode = helper.String(v.(string))
					}
					if v, ok := dbTablesMap["databases"]; ok {
						for _, item := range v.([]interface{}) {
							dbItemMap := item.(map[string]interface{})
							migrateDBItem := &cynosdb.MigrateDBItem{}
							if v, ok := dbItemMap["db_name"]; ok && v.(string) != "" {
								migrateDBItem.DbName = helper.String(v.(string))
							}
							if v, ok := dbItemMap["migrate_table_mode"]; ok && v.(string) != "" {
								migrateDBItem.MigrateTableMode = helper.String(v.(string))
							}
							if v, ok := dbItemMap["tables"]; ok {
								for _, tableItem := range v.([]interface{}) {
									tableMap := tableItem.(map[string]interface{})
									migrateTableItem := &cynosdb.MigrateTableItem{}
									if v, ok := tableMap["table_name"]; ok && v.(string) != "" {
										migrateTableItem.TableName = helper.String(v.(string))
									}
									migrateDBItem.Tables = append(migrateDBItem.Tables, migrateTableItem)
								}
							}
							migrateObject.Databases = append(migrateObject.Databases, migrateDBItem)
						}
					}
					objects.DatabaseTables = migrateObject
				}
			}
			request.Objects = objects
		}
	}

	if v, ok := d.GetOkExists("port"); ok {
		request.Port = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("goods_num"); ok {
		request.GoodsNum = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("instance_name"); ok {
		request.InstanceName = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("replicas_num"); ok {
		request.ReplicasNum = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("instance_type"); ok {
		request.InstanceType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("storage_type"); ok {
		request.StorageType = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("auto_voucher"); ok {
		request.AutoVoucher = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("order_source"); ok {
		request.OrderSource = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("deal_mode"); ok {
		request.DealMode = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("vpc_id"); ok {
		request.VpcId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("subnet_id"); ok {
		request.SubnetId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("security_group_ids"); ok {
		for _, item := range v.([]interface{}) {
			request.SecurityGroupIds = append(request.SecurityGroupIds, helper.String(item.(string)))
		}
	}

	if v, ok := d.GetOk("libra_db_version"); ok {
		request.LibraDBVersion = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("time_span"); ok {
		request.TimeSpan = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("time_unit"); ok {
		request.TimeUnit = helper.String(v.(string))
	}

	if v, ok := d.GetOk("src_instance_id"); ok {
		request.SrcInstanceId = helper.String(v.(string))
	}

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCynosdbClient().AddLibraDBInstances(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Create cynosdb libra db instance failed, Response is nil."))
		}

		response = result
		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s create cynosdb libra db instance failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	if response.Response.ResourceIds == nil || len(response.Response.ResourceIds) == 0 {
		return fmt.Errorf("Create cynosdb libra db instance failed, ResourceIds is empty.")
	}

	instanceId := *response.Response.ResourceIds[0]
	d.SetId(strings.Join([]string{clusterId, instanceId}, tccommon.FILED_SP))

	if response.Response.BigDealIds != nil {
		bigDealIds := make([]string, 0, len(response.Response.BigDealIds))
		for _, v := range response.Response.BigDealIds {
			bigDealIds = append(bigDealIds, *v)
		}
		_ = d.Set("big_deal_ids", bigDealIds)
	}

	if response.Response.TranId != nil {
		_ = d.Set("tran_id", response.Response.TranId)
	}

	if response.Response.DealNames != nil {
		dealNames := make([]string, 0, len(response.Response.DealNames))
		for _, v := range response.Response.DealNames {
			dealNames = append(dealNames, *v)
		}
		_ = d.Set("deal_names", dealNames)
	}

	if response.Response.ResourceIds != nil {
		resourceIds := make([]string, 0, len(response.Response.ResourceIds))
		for _, v := range response.Response.ResourceIds {
			resourceIds = append(resourceIds, *v)
		}
		_ = d.Set("resource_ids", resourceIds)
	}

	// Poll DescribeLibraDBInstanceDetail until instance is ready
	descRequest := cynosdb.NewDescribeLibraDBInstanceDetailRequest()
	descRequest.ClusterId = helper.String(clusterId)
	descRequest.InstanceId = helper.String(instanceId)

	pollErr := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCynosdbClient().DescribeLibraDBInstanceDetail(descRequest)
		if e != nil {
			return tccommon.RetryError(e)
		}

		if result == nil || result.Response == nil {
			return resource.RetryableError(fmt.Errorf("Waiting for cynosdb libra db instance to be ready."))
		}

		if result.Response.Status != nil && *result.Response.Status == "running" {
			return nil
		}

		return resource.RetryableError(fmt.Errorf("Cynosdb libra db instance is still in status: %s", helper.PString(result.Response.Status)))
	})

	if pollErr != nil {
		log.Printf("[WARN]%s poll cynosdb libra db instance status failed, reason:%+v", logId, pollErr)
	}

	return resourceTencentCloudCynosdbLibraDbInstanceAttachmentRead(d, meta)
}

func resourceTencentCloudCynosdbLibraDbInstanceAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cynosdb_libra_db_instance_attachment.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		request = cynosdb.NewDescribeLibraDBInstanceDetailRequest()
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	clusterId := idSplit[0]
	instanceId := idSplit[1]

	request.ClusterId = helper.String(clusterId)
	request.InstanceId = helper.String(instanceId)

	var response *cynosdb.DescribeLibraDBInstanceDetailResponse
	reqErr := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCynosdbClient().DescribeLibraDBInstanceDetail(request)
		if e != nil {
			return tccommon.RetryError(e)
		}

		response = result
		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s read cynosdb libra db instance failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	if response == nil || response.Response == nil || response.Response.InstanceId == nil {
		log.Printf("[WARN]%s resource `tencentcloud_cynosdb_libra_db_instance_attachment` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		d.SetId("")
		return nil
	}

	respData := response.Response

	_ = d.Set("cluster_id", clusterId)
	_ = d.Set("instance_id", instanceId)

	if respData.Zone != nil {
		_ = d.Set("zone", respData.Zone)
	}

	if respData.Cpu != nil {
		_ = d.Set("cpu", respData.Cpu)
	}

	if respData.Memory != nil {
		_ = d.Set("mem", respData.Memory)
	}

	if respData.Storage != nil {
		_ = d.Set("storage_size", respData.Storage)
	}

	if respData.PayMode != nil {
		_ = d.Set("pay_mode", respData.PayMode)
	}

	if respData.InstanceName != nil {
		_ = d.Set("instance_name", respData.InstanceName)
	}

	if respData.InstanceType != nil {
		_ = d.Set("instance_type", respData.InstanceType)
	}

	if respData.StorageType != nil {
		_ = d.Set("storage_type", respData.StorageType)
	}

	if respData.VpcId != nil {
		_ = d.Set("vpc_id", respData.VpcId)
	}

	if respData.SubnetId != nil {
		_ = d.Set("subnet_id", respData.SubnetId)
	}

	if respData.LibraDBVersion != nil {
		_ = d.Set("libra_db_version", respData.LibraDBVersion)
	}

	return nil
}

func resourceTencentCloudCynosdbLibraDbInstanceAttachmentUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cynosdb_libra_db_instance_attachment.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	immutableArgs := []string{
		"zone", "cpu", "mem", "storage_size", "pay_mode", "objects",
		"port", "goods_num", "instance_name", "replicas_num", "instance_type",
		"storage_type", "auto_voucher", "order_source", "deal_mode",
		"vpc_id", "subnet_id", "security_group_ids", "libra_db_version",
		"time_span", "time_unit", "src_instance_id",
		"isolate_reason_types", "isolate_reason",
	}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	return resourceTencentCloudCynosdbLibraDbInstanceAttachmentRead(d, meta)
}

func resourceTencentCloudCynosdbLibraDbInstanceAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cynosdb_libra_db_instance_attachment.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		request = cynosdb.NewIsolateLibraDBClusterRequest()
	)

	_ = ctx

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	clusterId := idSplit[0]
	request.ClusterId = helper.String(clusterId)

	if v, ok := d.GetOk("isolate_reason_types"); ok {
		for _, item := range v.([]interface{}) {
			request.IsolateReasonTypes = append(request.IsolateReasonTypes, helper.IntInt64(item.(int)))
		}
	}

	if v, ok := d.GetOk("isolate_reason"); ok {
		request.IsolateReason = helper.String(v.(string))
	}

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCynosdbClient().IsolateLibraDBCluster(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s delete cynosdb libra db instance failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	return nil
}
