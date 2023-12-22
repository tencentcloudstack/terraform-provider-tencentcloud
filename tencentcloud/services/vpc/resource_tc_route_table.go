package vpc

import (
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svctag "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/tag"

	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudVpcRouteTable() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudVpcRouteTableCreate,
		Read:   resourceTencentCloudVpcRouteTableRead,
		Update: resourceTencentCloudVpcRouteTableUpdate,
		Delete: resourceTencentCloudVpcRouteTableDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"vpc_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ID of VPC to which the route table should be associated.",
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: tccommon.ValidateStringLengthInRange(1, 60),
				Description:  "The name of routing table.",
			},
			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "The tags of routing table.",
			},

			// Computed values
			"subnet_ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "ID list of the subnets associated with this route table.",
			},
			"route_entry_ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "ID list of the routing entries.",
			},
			"is_default": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Indicates whether it is the default routing table.",
			},
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Creation time of the routing table.",
			},
		},
	}
}

func resourceTencentCloudVpcRouteTableCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_route_table.create")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	vpcService := VpcService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var (
		vpcId string
		name  string
		tags  map[string]string
	)
	if temp, ok := d.GetOk("vpc_id"); ok {
		vpcId = temp.(string)
		if len(vpcId) < 1 {
			return fmt.Errorf("vpc_id should be not empty string")
		}
	}
	if temp, ok := d.GetOk("name"); ok {
		name = temp.(string)
	}

	if temp := helper.GetTags(d, "tags"); len(temp) > 0 {
		tags = temp
	}

	routeTableId, err := vpcService.CreateRouteTable(ctx, name, vpcId, tags)
	if err != nil {
		return err
	}
	d.SetId(routeTableId)

	if tags := helper.GetTags(d, "tags"); len(tags) > 0 {
		tagService := svctag.NewTagService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())

		region := meta.(tccommon.ProviderMeta).GetAPIV3Conn().Region
		resourceName := fmt.Sprintf("qcs::vpc:%s:uin/:rtb/%s", region, routeTableId)

		if err := tagService.ModifyTags(ctx, resourceName, tags, nil); err != nil {
			return err
		}
	}

	return resourceTencentCloudVpcRouteTableRead(d, meta)
}

func resourceTencentCloudVpcRouteTableRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_route_table.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	id := d.Id()

	var (
		info VpcRouteTableBasicInfo
		has  int
		e    error
	)
	service := VpcService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		info, has, e = service.DescribeRouteTable(ctx, id)
		if e != nil {
			return tccommon.RetryError(e)
		}
		// deleted
		if has == 0 {
			d.SetId("")
			return nil
		}
		if has != 1 {
			errRet := fmt.Errorf("one route_table_id read get %d route_table info", has)
			log.Printf("[CRITAL]%s %s", logId, errRet.Error())
			return resource.NonRetryableError(errRet)
		}
		return nil
	})
	if err != nil {
		return err
	}
	if has == 0 {
		return nil
	}
	routeEntryIds := make([]string, 0, len(info.entryInfos))
	for _, v := range info.entryInfos {
		tfRouteEntryId := fmt.Sprintf("%d.%s", v.routeEntryId, id)
		routeEntryIds = append(routeEntryIds, tfRouteEntryId)
	}

	tagService := svctag.NewTagService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())

	region := meta.(tccommon.ProviderMeta).GetAPIV3Conn().Region
	tags, err := tagService.DescribeResourceTags(ctx, "vpc", "rtb", region, id)
	if err != nil {
		return err
	}

	_ = d.Set("vpc_id", info.vpcId)
	_ = d.Set("name", info.name)
	_ = d.Set("subnet_ids", info.subnetIds)
	_ = d.Set("route_entry_ids", routeEntryIds)
	_ = d.Set("is_default", info.isDefault)
	_ = d.Set("create_time", info.createTime)
	_ = d.Set("tags", tags)

	return nil
}

func resourceTencentCloudVpcRouteTableUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_route_table.update")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	id := d.Id()

	service := VpcService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	d.Partial(true)

	if d.HasChange("name") {
		name := d.Get("name").(string)
		err := service.ModifyRouteTableAttribute(ctx, id, name)
		if err != nil {
			return err
		}

	}

	if d.HasChange("tags") {
		oldTags, newTags := d.GetChange("tags")
		replaceTags, deleteTags := svctag.DiffTags(oldTags.(map[string]interface{}), newTags.(map[string]interface{}))

		tagService := svctag.NewTagService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())

		region := meta.(tccommon.ProviderMeta).GetAPIV3Conn().Region
		resourceName := fmt.Sprintf("qcs::vpc:%s:uin/:rtb/%s", region, id)

		if err := tagService.ModifyTags(ctx, resourceName, replaceTags, deleteTags); err != nil {
			return err
		}

	}

	d.Partial(false)

	return resourceTencentCloudVpcRouteTableRead(d, meta)
}

func resourceTencentCloudVpcRouteTableDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_route_table.delete")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := VpcService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		if err := service.DeleteRouteTable(ctx, d.Id()); err != nil {
			if sdkErr, ok := err.(*errors.TencentCloudSDKError); ok {
				if sdkErr.Code == VPCNotFound {
					return nil
				}
			}
			return resource.RetryableError(err)
		}
		return nil
	})

	return err
}
