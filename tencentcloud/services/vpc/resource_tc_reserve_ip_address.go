package vpc

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	svctag "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/tag"
)

func ResourceTencentCloudReserveIpAddress() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudReserveIpAddressCreate,
		Read:   resourceTencentCloudReserveIpAddressRead,
		Update: resourceTencentCloudReserveIpAddressUpdate,
		Delete: resourceTencentCloudReserveIpAddressDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"vpc_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "VPC unique ID.",
			},

			"ip_address": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Specify the reserved IP address of the intranet for which the IP application is requested.",
			},

			"subnet_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Subnet ID.",
			},

			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The IP name is reserved for the intranet.",
			},

			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The IP description is retained on the intranet.",
			},

			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Tags.",
			},
			"reserve_ip_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Reserve ip ID.",
			},
			"resource_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The intranet retains the resource instance ID bound to the IPs.",
			},
			"ip_type": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Ip type for product application.",
			},
			"state": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Binding status.",
			},
			"created_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Created time.",
			},
		},
	}
}

func resourceTencentCloudReserveIpAddressCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_reserve_ip_address.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	var (
		vpcId       string
		reserveIpId string
	)
	var (
		request  = vpc.NewCreateReserveIpAddressesRequest()
		response = vpc.NewCreateReserveIpAddressesResponse()
	)

	if v, ok := d.GetOk("vpc_id"); ok {
		vpcId = v.(string)
		request.VpcId = helper.String(vpcId)
	}

	if v, ok := d.GetOk("ip_address"); ok {
		ipAddress := v.(string)
		request.IpAddresses = append(request.IpAddresses, helper.String(ipAddress))
	} else {
		request.IpAddressCount = helper.IntUint64(1)
	}

	if v, ok := d.GetOk("subnet_id"); ok {
		request.SubnetId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("name"); ok {
		request.Name = helper.String(v.(string))
	}

	if v, ok := d.GetOk("description"); ok {
		request.Description = helper.String(v.(string))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseVpcClient().CreateReserveIpAddressesWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		if len(result.Response.ReserveIpAddressSet) > 0 {
			reserveIpId = *result.Response.ReserveIpAddressSet[0].ReserveIpId
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create reserve ip addresses failed, reason:%+v", logId, err)
		return err
	}

	_ = response

	if tags := helper.GetTags(d, "tags"); len(tags) > 0 {
		tcClient := meta.(tccommon.ProviderMeta).GetAPIV3Conn()
		tagService := svctag.NewTagService(tcClient)
		resourceName := tccommon.BuildTagResourceName("vpc", "rsvip", tcClient.Region, reserveIpId)
		if err := tagService.ModifyTags(ctx, resourceName, tags, nil); err != nil {
			return err
		}
	}

	d.SetId(strings.Join([]string{vpcId, reserveIpId}, tccommon.FILED_SP))

	return resourceTencentCloudReserveIpAddressRead(d, meta)
}

func resourceTencentCloudReserveIpAddressRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_reserve_ip_address.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	service := VpcService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	reserveIpId := idSplit[1]

	respData, err := service.DescribeReserveIpAddressesById(ctx, reserveIpId)
	if err != nil {
		return err
	}

	if respData == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `reserve_ip_addresses` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}
	if len(respData.ReserveIpAddressSet) < 1 {
		d.SetId("")
		log.Printf("[WARN]%s resource `reserve_ip_addresses` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}
	reserveIpAddress := respData.ReserveIpAddressSet[0]

	_ = d.Set("vpc_id", reserveIpAddress.VpcId)
	_ = d.Set("ip_address", reserveIpAddress.ReserveIpAddress)
	_ = d.Set("name", reserveIpAddress.Name)
	_ = d.Set("description", reserveIpAddress.Description)
	_ = d.Set("reserve_ip_id", reserveIpAddress.ReserveIpId)
	_ = d.Set("resource_id", reserveIpAddress.ResourceId)
	_ = d.Set("ip_type", reserveIpAddress.IpType)
	_ = d.Set("state", reserveIpAddress.State)
	_ = d.Set("created_time", reserveIpAddress.CreatedTime)

	tcClient := meta.(tccommon.ProviderMeta).GetAPIV3Conn()
	tagService := svctag.NewTagService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())
	tags, err := tagService.DescribeResourceTags(ctx, "vpc", "rsvip", tcClient.Region, reserveIpId)
	if err != nil {
		return err
	}
	_ = d.Set("tags", tags)
	_ = reserveIpId
	return nil
}

func resourceTencentCloudReserveIpAddressUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_reserve_ip_address.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	immutableArgs := []string{"vpc_id", "ip_address", "subnet_id"}
	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}
	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	vpcId := idSplit[0]
	reserveIpId := idSplit[1]

	needChange := false
	mutableArgs := []string{"name", "description"}
	for _, v := range mutableArgs {
		if d.HasChange(v) {
			needChange = true
			break
		}
	}

	if needChange {
		request := vpc.NewModifyReserveIpAddressRequest()

		request.VpcId = helper.String(vpcId)

		request.ReserveIpId = helper.String(reserveIpId)

		if v, ok := d.GetOk("name"); ok {
			request.Name = helper.String(v.(string))
		}

		if v, ok := d.GetOk("description"); ok {
			request.Description = helper.String(v.(string))
		}

		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseVpcClient().ModifyReserveIpAddressWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s update reserve ip addresses failed, reason:%+v", logId, err)
			return err
		}
	}

	if d.HasChange("tags") {
		tcClient := meta.(tccommon.ProviderMeta).GetAPIV3Conn()
		tagService := svctag.NewTagService(tcClient)
		oldTags, newTags := d.GetChange("tags")
		replaceTags, deleteTags := svctag.DiffTags(oldTags.(map[string]interface{}), newTags.(map[string]interface{}))
		resourceName := tccommon.BuildTagResourceName("vpc", "rsvip", tcClient.Region, reserveIpId)
		if err := tagService.ModifyTags(ctx, resourceName, replaceTags, deleteTags); err != nil {
			return err
		}
	}

	return resourceTencentCloudReserveIpAddressRead(d, meta)
}

func resourceTencentCloudReserveIpAddressDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_reserve_ip_address.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	vpcId := idSplit[0]
	reserveIpId := idSplit[1]

	var (
		request  = vpc.NewDeleteReserveIpAddressesRequest()
		response = vpc.NewDeleteReserveIpAddressesResponse()
	)

	request.VpcId = helper.String(vpcId)
	request.ReserveIpIds = []*string{&reserveIpId}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseVpcClient().DeleteReserveIpAddressesWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s delete reserve ip addresses failed, reason:%+v", logId, err)
		return err
	}

	_ = response
	_ = reserveIpId
	return nil
}
