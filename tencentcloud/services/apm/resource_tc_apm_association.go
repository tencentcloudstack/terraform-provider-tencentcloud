package apm

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	apmv20210622 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/apm/v20210622"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudApmAssociation() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudApmAssociationCreate,
		Read:   resourceTencentCloudApmAssociationRead,
		Update: resourceTencentCloudApmAssociationUpdate,
		Delete: resourceTencentCloudApmAssociationDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"product_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Associated product name. currently only supports Prometheus.",
			},

			"status": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: tccommon.ValidateAllowedIntValue([]int{1, 2}),
				Description:  "Status of the association relationship: // association status: 1 (enabled), 2 (disabled).",
			},

			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Business system ID.",
			},

			"peer_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Associated product instance ID.",
			},

			"topic": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the CKafka message topic.",
			},
		},
	}
}

func resourceTencentCloudApmAssociationCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_apm_association.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId       = tccommon.GetLogId(tccommon.ContextNil)
		ctx         = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request     = apmv20210622.NewModifyApmAssociationRequest()
		instanceId  string
		productName string
	)

	if v, ok := d.GetOk("product_name"); ok {
		request.ProductName = helper.String(v.(string))
		productName = v.(string)
	}

	if v, ok := d.GetOkExists("status"); ok {
		request.Status = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("instance_id"); ok {
		request.InstanceId = helper.String(v.(string))
		instanceId = v.(string)
	}

	if v, ok := d.GetOk("peer_id"); ok {
		request.PeerId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("topic"); ok {
		request.Topic = helper.String(v.(string))
	}

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseApmClient().ModifyApmAssociationWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s create apm association failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	d.SetId(strings.Join([]string{instanceId, productName}, tccommon.FILED_SP))
	return resourceTencentCloudApmAssociationRead(d, meta)
}

func resourceTencentCloudApmAssociationRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_apm_association.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = ApmService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	instanceId := idSplit[0]
	productName := idSplit[1]

	respData, err := service.DescribeApmAssociationById(ctx, instanceId, productName)
	if err != nil {
		return err
	}

	if respData == nil {
		log.Printf("[WARN]%s resource `tencentcloud_apm_association` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		d.SetId("")
		return nil
	}

	_ = d.Set("instance_id", instanceId)
	_ = d.Set("product_name", productName)

	if respData.PeerId != nil {
		_ = d.Set("peer_id", respData.PeerId)
	}

	if respData.Status != nil {
		_ = d.Set("status", respData.Status)
	}

	if respData.Topic != nil {
		_ = d.Set("topic", respData.Topic)
	}

	return nil
}

func resourceTencentCloudApmAssociationUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_apm_association.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId = tccommon.GetLogId(tccommon.ContextNil)
		ctx   = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	instanceId := idSplit[0]
	productName := idSplit[1]

	needChange := false
	mutableArgs := []string{"status", "peer_id", "topic"}
	for _, v := range mutableArgs {
		if d.HasChange(v) {
			needChange = true
			break
		}
	}

	if needChange {
		request := apmv20210622.NewModifyApmAssociationRequest()
		if v, ok := d.GetOkExists("status"); ok {
			request.Status = helper.IntUint64(v.(int))
		}

		if v, ok := d.GetOk("peer_id"); ok {
			request.PeerId = helper.String(v.(string))
		}

		if v, ok := d.GetOk("topic"); ok {
			request.Topic = helper.String(v.(string))
		}

		request.InstanceId = &instanceId
		request.ProductName = &productName
		reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseApmClient().ModifyApmAssociationWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}

			return nil
		})

		if reqErr != nil {
			log.Printf("[CRITAL]%s update apm association failed, reason:%+v", logId, reqErr)
			return reqErr
		}
	}

	return resourceTencentCloudApmAssociationRead(d, meta)
}

func resourceTencentCloudApmAssociationDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_apm_association.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request = apmv20210622.NewModifyApmAssociationRequest()
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	instanceId := idSplit[0]
	productName := idSplit[1]

	request.InstanceId = &instanceId
	request.ProductName = &productName
	request.Status = helper.IntUint64(4)
	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseApmClient().ModifyApmAssociationWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s delete apm association failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	return nil
}
