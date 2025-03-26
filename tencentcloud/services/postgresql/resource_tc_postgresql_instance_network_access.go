package postgresql

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	postgresqlv20170312 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/postgres/v20170312"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudPostgresqlInstanceNetworkAccess() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudPostgresqlInstanceNetworkAccessCreate,
		Read:   resourceTencentCloudPostgresqlInstanceNetworkAccessRead,
		Delete: resourceTencentCloudPostgresqlInstanceNetworkAccessDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"db_instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Instance ID in the format of postgres-6bwgamo3.",
			},

			"vpc_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Unified VPC ID.",
			},

			"subnet_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Subnet ID.",
			},

			"vip": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "Target VIP.",
			},
		},
	}
}

func resourceTencentCloudPostgresqlInstanceNetworkAccessCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_postgresql_instance_network_access.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId        = tccommon.GetLogId(tccommon.ContextNil)
		ctx          = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request      = postgresqlv20170312.NewCreateDBInstanceNetworkAccessRequest()
		response     = postgresqlv20170312.NewCreateDBInstanceNetworkAccessResponse()
		dbInsntaceId string
		vpcId        string
		subnetId     string
		vip          string
	)

	if v, ok := d.GetOk("db_instance_id"); ok {
		request.DBInstanceId = helper.String(v.(string))
		dbInsntaceId = v.(string)
	}

	if v, ok := d.GetOk("vpc_id"); ok {
		request.VpcId = helper.String(v.(string))
		vpcId = v.(string)
	}

	if v, ok := d.GetOk("subnet_id"); ok {
		request.SubnetId = helper.String(v.(string))
		subnetId = v.(string)
	}

	request.IsAssignVip = helper.Bool(false)

	if v, ok := d.GetOk("vip"); ok {
		request.Vip = helper.String(v.(string))
		request.IsAssignVip = helper.Bool(true)
		vip = v.(string)
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UsePostgresqlV20170312Client().CreateDBInstanceNetworkAccessWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Create postgresql instance network access failed, Response is nil."))
		}

		response = result
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create postgresql instance network access failed, reason:%+v", logId, err)
		return err
	}

	if response.Response.FlowId == nil {
		return fmt.Errorf("FlowId is nil.")
	}

	// wait & get vip
	flowId := *response.Response.FlowId
	flowRequest := postgresqlv20170312.NewDescribeTasksRequest()
	flowRequest.TaskId = helper.Int64Uint64(flowId)
	err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UsePostgresqlV20170312Client().DescribeTasksWithContext(ctx, flowRequest)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, flowRequest.GetAction(), flowRequest.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil || result.Response.TaskSet == nil {
			return resource.NonRetryableError(fmt.Errorf("Describe tasks failed, Response is nil."))
		}

		if len(result.Response.TaskSet) == 0 {
			return resource.RetryableError(fmt.Errorf("wait TaskSet init."))
		}

		if result.Response.TaskSet[0].Status != nil && *result.Response.TaskSet[0].Status == "Success" {
			if result.Response.TaskSet[0].TaskDetail != nil && result.Response.TaskSet[0].TaskDetail.Output != nil {
				outPutObj := make(map[string]interface{})
				outputStr := *result.Response.TaskSet[0].TaskDetail.Output
				e := json.Unmarshal([]byte(outputStr), &outPutObj)
				if e != nil {
					return resource.NonRetryableError(fmt.Errorf("Json unmarshall output error: %s.", e.Error()))
				}

				dBInstanceNetInfo := outPutObj["DBInstanceNetInfo"].(map[string]interface{})
				vip = dBInstanceNetInfo["Ip"].(string)
				return nil
			}
		}

		return resource.RetryableError(fmt.Errorf("postgresql instance network access is running, status is %s.", *result.Response.TaskSet[0].Status))
	})

	if err != nil {
		log.Printf("[CRITAL]%s create postgresql instance network access failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(strings.Join([]string{dbInsntaceId, vpcId, subnetId, vip}, tccommon.FILED_SP))

	return resourceTencentCloudPostgresqlInstanceNetworkAccessRead(d, meta)
}

func resourceTencentCloudPostgresqlInstanceNetworkAccessRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_postgresql_instance_network_access.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = PostgresqlService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 4 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	dbInsntaceId := idSplit[0]
	vpcId := idSplit[1]
	subnetId := idSplit[2]
	vip := idSplit[3]

	_ = d.Set("db_instance_id", dbInsntaceId)

	_ = d.Set("vpc_id", vpcId)

	_ = d.Set("subnet_id", subnetId)

	respData, err := service.DescribePostgresqlInstanceNetworkAccessById(ctx, dbInsntaceId)
	if err != nil {
		return err
	}

	if respData == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `postgresql_instance_network_access` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	var checkFlag bool
	if respData.DBInstanceNetInfo != nil && len(respData.DBInstanceNetInfo) > 0 {
		for _, item := range respData.DBInstanceNetInfo {
			if *item.Ip == vip {
				_ = d.Set("vip", item.Ip)
				checkFlag = true
				break
			}
		}
	}

	if checkFlag == false {
		return fmt.Errorf("Not found vip %s, please check if it has been deleted.", vip)
	}

	return nil
}

func resourceTencentCloudPostgresqlInstanceNetworkAccessDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_postgresql_instance_network_access.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId    = tccommon.GetLogId(tccommon.ContextNil)
		ctx      = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request  = postgresqlv20170312.NewDeleteDBInstanceNetworkAccessRequest()
		response = postgresqlv20170312.NewDeleteDBInstanceNetworkAccessResponse()
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 4 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	dbInsntaceId := idSplit[0]
	vpcId := idSplit[1]
	subnetId := idSplit[2]
	vip := idSplit[3]

	request.DBInstanceId = helper.String(dbInsntaceId)
	request.VpcId = helper.String(vpcId)
	request.SubnetId = helper.String(subnetId)
	request.Vip = helper.String(vip)
	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UsePostgresqlV20170312Client().DeleteDBInstanceNetworkAccessWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Delete postgresql instance network access failed, Response is nil."))
		}

		response = result
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s delete postgresql instance network access failed, reason:%+v", logId, err)
		return err
	}

	// wait
	flowId := *response.Response.FlowId
	flowRequest := postgresqlv20170312.NewDescribeTasksRequest()
	flowRequest.TaskId = helper.Int64Uint64(flowId)
	err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UsePostgresqlV20170312Client().DescribeTasksWithContext(ctx, flowRequest)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil || result.Response.TaskSet == nil {
			return resource.NonRetryableError(fmt.Errorf("Describe tasks failed, Response is nil."))
		}

		if len(result.Response.TaskSet) == 0 {
			return resource.RetryableError(fmt.Errorf("wait TaskSet init."))
		}

		if result.Response.TaskSet[0].Status != nil && *result.Response.TaskSet[0].Status == "Success" {
			return nil
		}

		return resource.RetryableError(fmt.Errorf("postgresql instance network access is running, status is %s.", *result.Response.TaskSet[0].Status))
	})

	if err != nil {
		log.Printf("[CRITAL]%s delete postgresql instance network access failed, reason:%+v", logId, err)
		return err
	}

	return nil
}
