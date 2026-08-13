package teo

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	teov20220901 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudTeoFunctionReplica() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTeoFunctionReplicaCreate,
		Read:   resourceTencentCloudTeoFunctionReplicaRead,
		Update: resourceTencentCloudTeoFunctionReplicaUpdate,
		Delete: resourceTencentCloudTeoFunctionReplicaDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"zone_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Zone ID.",
			},
			"function_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Function ID.",
			},
			"replica_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Edge function replica name. Limited to 1-50 characters, allowed characters are a-z, 0-9, -, and - cannot be used alone or consecutively, nor at the beginning or end. Replica names must be unique under the same FunctionId.",
			},
			"content": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Edge function replica content. Currently only supports JavaScript code, maximum 5MB.",
			},
			"remark": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Edge function replica description. Maximum 50 characters.",
			},
		},
	}
}

func resourceTencentCloudTeoFunctionReplicaCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_function_replica.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId       = tccommon.GetLogId(tccommon.ContextNil)
		ctx         = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request     = teov20220901.NewCreateFunctionReplicaRequest()
		zoneId      string
		functionId  string
		replicaName string
	)

	if v, ok := d.GetOk("zone_id"); ok {
		request.ZoneId = helper.String(v.(string))
		zoneId = v.(string)
	}

	if v, ok := d.GetOk("function_id"); ok {
		request.FunctionId = helper.String(v.(string))
		functionId = v.(string)
	}

	if v, ok := d.GetOk("replica_name"); ok {
		request.ReplicaName = helper.String(v.(string))
		replicaName = v.(string)
	}

	if v, ok := d.GetOk("content"); ok {
		request.Content = helper.String(v.(string))
	}

	if v, ok := d.GetOk("remark"); ok {
		request.Remark = helper.String(v.(string))
	}

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoV20220901Client().CreateFunctionReplicaWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Create teo function replica failed, Response is nil."))
		}

		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s create teo function replica failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	d.SetId(strings.Join([]string{zoneId, functionId, replicaName}, tccommon.FILED_SP))
	return resourceTencentCloudTeoFunctionReplicaRead(d, meta)
}

func resourceTencentCloudTeoFunctionReplicaRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_function_replica.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request = teov20220901.NewDescribeFunctionReplicasRequest()
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken, id is %s", d.Id())
	}

	zoneId := idSplit[0]
	functionId := idSplit[1]
	replicaName := idSplit[2]

	request.ZoneId = helper.String(zoneId)
	request.FunctionId = helper.String(functionId)
	request.Filters = []*teov20220901.AdvancedFilter{
		{
			Name:   helper.String("replica-name"),
			Values: []*string{helper.String(replicaName)},
		},
	}
	request.Limit = helper.Int64(200)

	var response *teov20220901.DescribeFunctionReplicasResponse
	reqErr := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoV20220901Client().DescribeFunctionReplicasWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		response = result
		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s read teo function replica failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	if response == nil || response.Response == nil {
		log.Printf("[WARN]%s resource `tencentcloud_teo_function_replica` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		d.SetId("")
		return nil
	}

	var targetReplica *teov20220901.FunctionReplica
	for _, replica := range response.Response.FunctionReplicas {
		if replica.ReplicaName != nil && *replica.ReplicaName == replicaName {
			targetReplica = replica
			break
		}
	}

	if targetReplica == nil {
		log.Printf("[WARN]%s resource `tencentcloud_teo_function_replica` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		d.SetId("")
		return nil
	}

	_ = d.Set("zone_id", zoneId)
	_ = d.Set("function_id", functionId)
	_ = d.Set("replica_name", replicaName)

	if targetReplica.Content != nil {
		_ = d.Set("content", targetReplica.Content)
	}

	if targetReplica.Remark != nil {
		_ = d.Set("remark", targetReplica.Remark)
	}

	return nil
}

func resourceTencentCloudTeoFunctionReplicaUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_function_replica.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request = teov20220901.NewModifyFunctionReplicaRequest()
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken, id is %s", d.Id())
	}

	zoneId := idSplit[0]
	functionId := idSplit[1]
	replicaName := idSplit[2]

	request.ZoneId = helper.String(zoneId)
	request.FunctionId = helper.String(functionId)
	request.ReplicaName = helper.String(replicaName)

	needChange := false
	mutableArgs := []string{"content", "remark"}
	for _, v := range mutableArgs {
		if d.HasChange(v) {
			needChange = true
			break
		}
	}

	if needChange {
		if v, ok := d.GetOk("content"); ok {
			request.Content = helper.String(v.(string))
		}

		if v, ok := d.GetOk("remark"); ok {
			request.Remark = helper.String(v.(string))
		}

		reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoV20220901Client().ModifyFunctionReplicaWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}

			return nil
		})

		if reqErr != nil {
			log.Printf("[CRITAL]%s update teo function replica failed, reason:%+v", logId, reqErr)
			return reqErr
		}
	}

	return resourceTencentCloudTeoFunctionReplicaRead(d, meta)
}

func resourceTencentCloudTeoFunctionReplicaDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_function_replica.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request = teov20220901.NewDeleteFunctionReplicaRequest()
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken, id is %s", d.Id())
	}

	zoneId := idSplit[0]
	functionId := idSplit[1]
	replicaName := idSplit[2]

	request.ZoneId = helper.String(zoneId)
	request.FunctionId = helper.String(functionId)
	request.ReplicaNames = []*string{helper.String(replicaName)}

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoV20220901Client().DeleteFunctionReplicaWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s delete teo function replica failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	return nil
}
