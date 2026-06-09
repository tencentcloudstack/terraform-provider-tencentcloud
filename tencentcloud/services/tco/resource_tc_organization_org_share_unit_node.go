package tco

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
)

func ResourceTencentCloudOrganizationOrgShareUnitNode() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudOrganizationOrgShareUnitNodeCreate,
		Read:   resourceTencentCloudOrganizationOrgShareUnitNodeRead,
		Delete: resourceTencentCloudOrganizationOrgShareUnitNodeDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"unit_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Shared unit ID.",
			},

			"node_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeInt,
				Description: "Organization department ID.",
			},
		},
	}
}

func resourceTencentCloudOrganizationOrgShareUnitNodeCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_organization_org_share_unit_node.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service = OrganizationService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		unitId  string
		nodeId  int64
	)

	if v, ok := d.GetOk("unit_id"); ok {
		unitId = v.(string)
	}

	if v, ok := d.GetOk("node_id"); ok {
		nodeId = int64(v.(int))
	}

	err := service.AddOrganizationOrgShareUnitNodeById(ctx, unitId, nodeId)
	if err != nil {
		log.Printf("[CRITAL]%s create organization org share unit node failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(strings.Join([]string{unitId, strconv.FormatInt(nodeId, 10)}, tccommon.FILED_SP))

	return resourceTencentCloudOrganizationOrgShareUnitNodeRead(d, meta)
}

func resourceTencentCloudOrganizationOrgShareUnitNodeRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_organization_org_share_unit_node.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service = OrganizationService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	unitId := idSplit[0]
	nodeIdStr := idSplit[1]

	nodeId, err := strconv.ParseInt(nodeIdStr, 10, 64)
	if err != nil {
		return fmt.Errorf("parse node_id failed, %s", err.Error())
	}

	orgShareUnitNode, err := service.DescribeOrganizationOrgShareUnitNodeById(ctx, unitId, nodeId)
	if err != nil {
		return err
	}

	if orgShareUnitNode == nil {
		log.Printf("[WARN]%s resource `tencentcloud_organization_org_share_unit_node` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		d.SetId("")
		return nil
	}

	_ = d.Set("unit_id", unitId)
	_ = d.Set("node_id", int(nodeId))

	return nil
}

func resourceTencentCloudOrganizationOrgShareUnitNodeDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_organization_org_share_unit_node.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service = OrganizationService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	unitId := idSplit[0]
	nodeIdStr := idSplit[1]

	nodeId, err := strconv.ParseInt(nodeIdStr, 10, 64)
	if err != nil {
		return fmt.Errorf("parse node_id failed, %s", err.Error())
	}

	if err := service.DeleteOrganizationOrgShareUnitNodeById(ctx, unitId, nodeId); err != nil {
		return err
	}

	return nil
}
