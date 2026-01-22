package ckafka

import (
	"context"
	"fmt"
	"strings"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceTencentCloudCkafkaUser() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCkafkaUserCreate,
		Read:   resourceTencentCloudCkafkaUserRead,
		Update: resourceTencentCloudCkafkaUserUpdate,
		Delete: resourceTencentCloudCkafkaUserDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ID of the ckafka instance.",
			},
			"account_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Account name used to access to ckafka instance.",
			},
			"password": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				Description: "Password of the account.",
			},
			// computed
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Creation time of the account.",
			},
			"update_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The last update time of the account.",
			},
		},
	}
}

func resourceTencentCloudCkafkaUserCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ckafka_user.create")()

	var (
		logId         = tccommon.GetLogId(tccommon.ContextNil)
		ctx           = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		ckafkaService = CkafkaService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	instanceId := d.Get("instance_id").(string)
	accountName := d.Get("account_name").(string)
	password := d.Get("password").(string)

	if err := ckafkaService.CreateUser(ctx, instanceId, accountName, password); err != nil {
		return fmt.Errorf("[CRITAL]%s create ckafka user failed, reason:%+v", logId, err)
	}

	d.SetId(strings.Join([]string{instanceId, accountName}, tccommon.FILED_SP))
	return resourceTencentCloudCkafkaUserRead(d, meta)
}

func resourceTencentCloudCkafkaUserRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ckafka_user.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId         = tccommon.GetLogId(tccommon.ContextNil)
		ctx           = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		ckafkaService = CkafkaService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		id            = d.Id()
	)

	info, has, err := ckafkaService.DescribeUserByUserId(ctx, id)
	if err != nil {
		return err
	}

	if !has {
		d.SetId("")
		return nil
	}

	items := strings.Split(id, tccommon.FILED_SP)
	_ = d.Set("instance_id", items[0])
	_ = d.Set("account_name", info.Name)
	_ = d.Set("create_time", info.CreateTime)
	_ = d.Set("update_time", info.UpdateTime)

	return nil
}

func resourceTencentCloudCkafkaUserUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ckafka_user.update")()

	var (
		logId         = tccommon.GetLogId(tccommon.ContextNil)
		ctx           = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		ckafkaService = CkafkaService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	if d.HasChange("password") {
		idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
		if len(idSplit) != 2 {
			return fmt.Errorf("id is broken, id is %s", d.Id())
		}

		instanceId, user := idSplit[0], idSplit[1]
		old, new := d.GetChange("password")
		if err := ckafkaService.ModifyPassword(ctx, instanceId, user, old.(string), new.(string)); err != nil {
			return err
		}
	}

	return resourceTencentCloudCkafkaUserRead(d, meta)
}

func resourceTencentCloudCkafkaUserDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ckafka_user.delete")()

	var (
		logId         = tccommon.GetLogId(tccommon.ContextNil)
		ctx           = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		ckafkaService = CkafkaService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	if err := ckafkaService.DeleteUser(ctx, d.Id()); err != nil {
		return err
	}

	return nil
}
