package postgresql

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svctag "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/tag"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	postgresql "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/postgres/v20170312"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudPostgresqlReadonlyGroupNetworkAccessAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudPostgresqlReadonlyGroupNetworkAccessAttachmentCreate,
		Read:   resourceTencentCloudPostgresqlReadonlyGroupNetworkAccessAttachmentRead,
		Delete: resourceTencentCloudPostgresqlReadonlyGroupNetworkAccessAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"db_instance_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Master database instance ID.",
			},

			"readonly_group_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "RO group identifier.",
			},

			"vpc_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Unified VPC ID.",
			},

			"subnet_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Subnet ID.",
			},

			"is_assign_vip": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeBool,
				Description: "Whether to manually assign the VIP. Valid values:true (manually assign), false (automatically assign).",
			},

			"vip": {
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Target VIP.",
			},

			"tags": {
				Type:        schema.TypeMap,
				ForceNew:    true,
				Optional:    true,
				Description: "Tag description list.",
			},
		},
	}
}

func resourceTencentCloudPostgresqlReadonlyGroupNetworkAccessAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_postgresql_readonly_group_network_access_attachment.create")()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		request      = postgresql.NewCreateReadOnlyGroupNetworkAccessRequest()
		dbInstanceId string
		roGroupId    string
		vpcId        string
		vip          string
		port         string
		isUserAssign bool
	)
	if v, ok := d.GetOk("db_instance_id"); ok {
		dbInstanceId = v.(string)
	}

	if v, ok := d.GetOk("readonly_group_id"); ok {
		request.ReadOnlyGroupId = helper.String(v.(string))
		roGroupId = v.(string)
	}

	if v, ok := d.GetOk("vpc_id"); ok {
		request.VpcId = helper.String(v.(string))
		vpcId = v.(string)
	}

	if v, ok := d.GetOk("subnet_id"); ok {
		request.SubnetId = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("is_assign_vip"); ok {
		request.IsAssignVip = helper.Bool(v.(bool))
		isUserAssign = v.(bool)
	}

	if v, ok := d.GetOk("vip"); ok {
		request.Vip = helper.String(v.(string))
		vip = v.(string)
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UsePostgresqlClient().CreateReadOnlyGroupNetworkAccess(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create postgresql ReadonlyGroupNetworkAccessAttachment failed, reason:%+v", logId, err)
		return err
	}

	id := strings.Join([]string{dbInstanceId, roGroupId, vpcId, vip, port}, tccommon.FILED_SP)

	service := PostgresqlService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	conf := tccommon.BuildStateChangeConf([]string{}, []string{"opened"}, 180*tccommon.ReadRetryTimeout, time.Second, service.PostgresqlReadonlyGroupNetworkAccessAttachmentStateRefreshFunc(id, []string{}))

	var ret interface{}
	var e error
	if ret, e = conf.WaitForState(); e != nil {
		return e
	} else {
		object := ret.(*postgresql.DBInstanceNetInfo)
		// fill out the port and vip
		if object != nil {
			if isUserAssign {
				// find the port
				if *object.VpcId == vpcId && *object.Ip == vip {
					port = helper.UInt64ToStr(*object.Port)
				}
			} else {
				// find the port and vip when is_assign_vip is false
				if *object.VpcId == vpcId {
					port = helper.UInt64ToStr(*object.Port)
					vip = *object.Ip
				}
			}
		}
	}

	id = strings.Join([]string{dbInstanceId, roGroupId, vpcId, vip, port}, tccommon.FILED_SP)
	d.SetId(id)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	if tags := helper.GetTags(d, "tags"); len(tags) > 0 {
		tagService := svctag.NewTagService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())
		region := meta.(tccommon.ProviderMeta).GetAPIV3Conn().Region
		resourceName := fmt.Sprintf("qcs::postgres:%s:uin/:dbInstanceId/%s", region, d.Id())
		if err := tagService.ModifyTags(ctx, resourceName, tags, nil); err != nil {
			return err
		}
	}

	return resourceTencentCloudPostgresqlReadonlyGroupNetworkAccessAttachmentRead(d, meta)
}

func resourceTencentCloudPostgresqlReadonlyGroupNetworkAccessAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_postgresql_readonly_group_network_access_attachment.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := PostgresqlService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 5 {
		return fmt.Errorf("id is broken,%s, location:%s", d.Id(), "resource.tencentcloud_postgresql_readonly_group_network_access_attachment.read")
	}

	dbInstanceId := idSplit[0]
	roGroupId := idSplit[1]
	vpcId := idSplit[2]
	vip := idSplit[3]
	port := idSplit[4]

	ReadonlyGroupNetworkAccessAttachment, err := service.DescribePostgresqlReadonlyGroupNetworkAccessAttachmentById(ctx, dbInstanceId, roGroupId)
	if err != nil {
		return err
	}

	if ReadonlyGroupNetworkAccessAttachment == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `PostgresqlReadonlyGroupNetworkAccessAttachment` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if ReadonlyGroupNetworkAccessAttachment.MasterDBInstanceId != nil {
		_ = d.Set("db_instance_id", ReadonlyGroupNetworkAccessAttachment.MasterDBInstanceId)
	}

	if ReadonlyGroupNetworkAccessAttachment.ReadOnlyGroupId != nil {
		_ = d.Set("readonly_group_id", ReadonlyGroupNetworkAccessAttachment.ReadOnlyGroupId)
	}

	if ReadonlyGroupNetworkAccessAttachment.VpcId != nil {
		_ = d.Set("vpc_id", ReadonlyGroupNetworkAccessAttachment.VpcId)
	}

	if ReadonlyGroupNetworkAccessAttachment.SubnetId != nil {
		_ = d.Set("subnet_id", ReadonlyGroupNetworkAccessAttachment.SubnetId)
	}

	if vip == "" {
		// That's mean isUserAssign is false and need to set vip assigned by system
		if ReadonlyGroupNetworkAccessAttachment.DBInstanceNetInfo != nil {
			for _, info := range ReadonlyGroupNetworkAccessAttachment.DBInstanceNetInfo {
				if *info.VpcId == vpcId && helper.UInt64ToStr(*info.Port) == port {
					if info.Ip != nil {
						vip = *info.Ip
						log.Printf("[DEBUG]%s the id:[%s]'s filed vip[%s] updated successfully!\n", logId, d.Id(), vip)
						break
					}
				}
			}
		}
		// update the vip into unique id
		id := strings.Join([]string{dbInstanceId, roGroupId, vpcId, vip, port}, tccommon.FILED_SP)
		d.SetId(id)
	}
	_ = d.Set("vip", vip)

	tcClient := meta.(tccommon.ProviderMeta).GetAPIV3Conn()
	tagService := svctag.NewTagService(tcClient)
	tags, err := tagService.DescribeResourceTags(ctx, "postgres", "dbInstanceId", tcClient.Region, d.Id())
	if err != nil {
		return err
	}
	_ = d.Set("tags", tags)

	return nil
}

func resourceTencentCloudPostgresqlReadonlyGroupNetworkAccessAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_postgresql_readonly_group_network_access_attachment.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := PostgresqlService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 5 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	var subnetId string
	roGroupId := idSplit[1]
	vpcId := idSplit[2]
	vip := idSplit[3]
	if v, ok := d.GetOk("subnet_id"); ok {
		subnetId = v.(string)
	}

	if err := service.DeletePostgresqlReadonlyGroupNetworkAccessAttachmentById(ctx, roGroupId, vpcId, subnetId, vip); err != nil {
		return err
	}

	conf := tccommon.BuildStateChangeConf([]string{}, []string{"closed"}, 180*tccommon.ReadRetryTimeout, time.Second, service.PostgresqlReadonlyGroupNetworkAccessAttachmentStateRefreshFunc(d.Id(), []string{}))

	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	return nil
}
