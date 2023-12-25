package vpc

import (
	"context"
	"fmt"
	"log"
	"strings"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudVpcLocalGateway() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudVpcLocalGatewayCreate,
		Read:   resourceTencentCloudVpcLocalGatewayRead,
		Update: resourceTencentCloudVpcLocalGatewayUpdate,
		Delete: resourceTencentCloudVpcLocalGatewayDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"local_gateway_name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Local gateway name.",
			},

			"vpc_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "VPC instance ID.",
			},

			"cdc_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "CDC instance ID.",
			},
		},
	}
}

func resourceTencentCloudVpcLocalGatewayCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_vpc_local_gateway.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		request        = vpc.NewCreateLocalGatewayRequest()
		response       = vpc.NewCreateLocalGatewayResponse()
		cdcId          string
		localGatewayId string
	)
	if v, ok := d.GetOk("local_gateway_name"); ok {
		request.LocalGatewayName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("vpc_id"); ok {
		request.VpcId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("cdc_id"); ok {
		cdcId = v.(string)
		request.CdcId = helper.String(v.(string))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseVpcClient().CreateLocalGateway(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create vpc localGateway failed, reason:%+v", logId, err)
		return err
	}

	localGatewayId = *response.Response.LocalGateway.UniqLocalGwId
	d.SetId(cdcId + tccommon.FILED_SP + localGatewayId)

	return resourceTencentCloudVpcLocalGatewayRead(d, meta)
}

func resourceTencentCloudVpcLocalGatewayRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_vpc_local_gateway.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := VpcService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	localGatewayId := idSplit[1]

	localGateway, err := service.DescribeVpcLocalGatewayById(ctx, localGatewayId)
	if err != nil {
		return err
	}

	if localGateway == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `VpcLocalGateway` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if localGateway.LocalGatewayName != nil {
		_ = d.Set("local_gateway_name", localGateway.LocalGatewayName)
	}

	if localGateway.VpcId != nil {
		_ = d.Set("vpc_id", localGateway.VpcId)
	}

	if localGateway.CdcId != nil {
		_ = d.Set("cdc_id", localGateway.CdcId)
	}

	return nil
}

func resourceTencentCloudVpcLocalGatewayUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_vpc_local_gateway.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	request := vpc.NewModifyLocalGatewayRequest()

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	cdcId := idSplit[0]
	localGatewayId := idSplit[1]

	request.CdcId = &cdcId
	request.LocalGatewayId = &localGatewayId

	if v, ok := d.GetOk("local_gateway_name"); ok {
		request.LocalGatewayName = helper.String(v.(string))
	}

	if d.HasChange("vpc_id") {
		if v, ok := d.GetOk("vpc_id"); ok {
			request.VpcId = helper.String(v.(string))
		}
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseVpcClient().ModifyLocalGateway(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update vpc localGateway failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudVpcLocalGatewayRead(d, meta)
}

func resourceTencentCloudVpcLocalGatewayDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_vpc_local_gateway.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := VpcService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	cdcId := idSplit[0]
	localGatewayId := idSplit[1]

	if err := service.DeleteVpcLocalGatewayById(ctx, cdcId, localGatewayId); err != nil {
		return err
	}

	return nil
}
