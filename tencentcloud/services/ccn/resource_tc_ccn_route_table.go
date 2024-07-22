package ccn

import (
	"context"
	"fmt"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"
)

func ResourceTencentCloudCcnRouteTable() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCcnRouteTableCreate,
		Read:   resourceTencentCloudCcnRouteTableRead,
		Update: resourceTencentCloudCcnRouteTableUpdate,
		Delete: resourceTencentCloudCcnRouteTableDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"ccn_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "CCN Instance ID.",
			},
			"name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "CCN Route table name.",
			},
			"description": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Description of CCN Route table.",
			},
			// computed
			"is_default_table": {
				Computed:    true,
				Type:        schema.TypeBool,
				Description: "True: default routing table False: non default routing table.",
			},
			"create_time": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "create time.",
			},
		},
	}
}

func resourceTencentCloudCcnRouteTableCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ccn_route_table.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId       = tccommon.GetLogId(tccommon.ContextNil)
		request     = vpc.NewCreateCcnRouteTablesRequest()
		response    = vpc.NewCreateCcnRouteTablesResponse()
		ccnId       string
		name        string
		description string
	)

	if v, ok := d.GetOk("ccn_id"); ok {
		ccnId = v.(string)
	}

	if v, ok := d.GetOk("name"); ok {
		name = v.(string)
	}

	if v, ok := d.GetOk("description"); ok {
		description = v.(string)
	}

	request.RouteTable = []*vpc.CcnBatchRouteTable{
		{
			CcnId:       helper.String(ccnId),
			Name:        helper.String(name),
			Description: helper.String(description),
		},
	}
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseVpcClient().CreateCcnRouteTables(request)
		if e != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), e.Error())
			return tccommon.RetryError(e)
		}

		if result == nil || result.Response == nil || len(result.Response.CcnRouteTableSet) != 1 {
			e = fmt.Errorf("create ccn route table failed")
			return resource.NonRetryableError(e)
		}

		response = result
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create ccn route table failed, reason:%s\n", logId, err.Error())
		return err
	}

	routeTableId := *response.Response.CcnRouteTableSet[0].CcnRouteTableId
	d.SetId(routeTableId)

	return resourceTencentCloudCcnRouteTableRead(d, meta)
}

func resourceTencentCloudCcnRouteTableRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ccn_route_table.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service = VpcService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	routeTableId := d.Id()
	ccnRouteTable, err := service.DescribeVpcCcnRouteTablesById(ctx, routeTableId)
	if err != nil {
		return err
	}

	if ccnRouteTable == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `VpcCcnRouteTable` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if ccnRouteTable.CcnId != nil {
		_ = d.Set("ccn_id", ccnRouteTable.CcnId)
	}

	if ccnRouteTable.RouteTableName != nil {
		_ = d.Set("name", ccnRouteTable.RouteTableName)
	}

	if ccnRouteTable.RouteTableDescription != nil {
		_ = d.Set("description", ccnRouteTable.RouteTableDescription)
	}

	if ccnRouteTable.IsDefaultTable != nil {
		_ = d.Set("is_default_table", ccnRouteTable.IsDefaultTable)
	}

	if ccnRouteTable.CreateTime != nil {
		_ = d.Set("create_time", ccnRouteTable.CreateTime)
	}

	return nil
}

func resourceTencentCloudCcnRouteTableUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ccn_route_table.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId        = tccommon.GetLogId(tccommon.ContextNil)
		routeTableId = d.Id()
	)

	if d.HasChange("name") || d.HasChange("description") {
		var (
			request     = vpc.NewModifyCcnRouteTablesRequest()
			name        string
			description string
		)

		if v, ok := d.GetOk("name"); ok {
			name = v.(string)
		}

		if v, ok := d.GetOk("description"); ok {
			description = v.(string)
		}

		request.RouteTableInfo = []*vpc.ModifyRouteTableInfo{
			{
				RouteTableId: helper.String(routeTableId),
				Name:         helper.String(name),
				Description:  helper.String(description),
			},
		}
		err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			_, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseVpcClient().ModifyCcnRouteTables(request)
			if e != nil {
				log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
					logId, request.GetAction(), request.ToJsonString(), e.Error())
				return tccommon.RetryError(e)
			}

			return nil
		})

		if err != nil {
			log.Printf("[CRITAL]%s modify ccn route table failed, reason:%s\n", logId, err.Error())
			return err
		}
	}

	return resourceTencentCloudCcnRouteTableRead(d, meta)
}

func resourceTencentCloudCcnRouteTableDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ccn_route_table.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId        = tccommon.GetLogId(tccommon.ContextNil)
		request      = vpc.NewDeleteCcnRouteTablesRequest()
		routeTableId = d.Id()
	)

	request.RouteTableId = helper.Strings([]string{routeTableId})
	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		_, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseVpcClient().DeleteCcnRouteTables(request)
		if e != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), e.Error())
			return tccommon.RetryError(e)
		}

		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s delete ccn route table failed, reason:%s\n", logId, err.Error())
		return err
	}

	return nil
}
