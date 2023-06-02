/*
Provides a resource to create a vpc local_gateway

Example Usage

```hcl
resource "tencentcloud_vpc_local_gateway" "local_gateway" {
  local_gateway_name = "local-gw-test"
  vpc_id             = "vpc-lh4nqig9"
  cdc_id             = "cluster-j9gyu1iy"
}
```

Import

vpc local_gateway can be imported using the id, e.g.

```
terraform import tencentcloud_vpc_local_gateway.local_gateway local_gateway_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudVpcLocalGateway() *schema.Resource {
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
	defer logElapsed("resource.tencentcloud_vpc_local_gateway.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

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

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseVpcClient().CreateLocalGateway(request)
		if e != nil {
			return retryError(e)
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
	d.SetId(cdcId + FILED_SP + localGatewayId)

	return resourceTencentCloudVpcLocalGatewayRead(d, meta)
}

func resourceTencentCloudVpcLocalGatewayRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_vpc_local_gateway.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := VpcService{client: meta.(*TencentCloudClient).apiV3Conn}

	idSplit := strings.Split(d.Id(), FILED_SP)
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
	defer logElapsed("resource.tencentcloud_vpc_local_gateway.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := vpc.NewModifyLocalGatewayRequest()

	idSplit := strings.Split(d.Id(), FILED_SP)
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

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseVpcClient().ModifyLocalGateway(request)
		if e != nil {
			return retryError(e)
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
	defer logElapsed("resource.tencentcloud_vpc_local_gateway.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := VpcService{client: meta.(*TencentCloudClient).apiV3Conn}
	idSplit := strings.Split(d.Id(), FILED_SP)
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
