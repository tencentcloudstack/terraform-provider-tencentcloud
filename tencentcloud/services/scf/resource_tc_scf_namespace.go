package scf

import (
	"context"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceTencentCloudScfNamespace() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudScfNamespaceCreate,
		Read:   resourceTencentCloudScfNamespaceRead,
		Update: resourceTencentCloudScfNamespaceUpdate,
		Delete: resourceTencentCloudScfNamespaceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"namespace": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name of the SCF namespace.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Description of the SCF namespace.",
			},

			// computed
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "SCF namespace creation time.",
			},
			"modify_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "SCF namespace last modified time.",
			},
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "SCF namespace type.",
			},
		},
	}
}

func resourceTencentCloudScfNamespaceCreate(d *schema.ResourceData, m interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_scf_namespace.create")()
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := ScfService{client: m.(tccommon.ProviderMeta).GetAPIV3Conn()}

	namespace := d.Get("namespace").(string)
	desc := d.Get("description").(string)

	if err := service.CreateNamespace(ctx, namespace, desc); err != nil {
		log.Printf("[CRITAL]%s create namespace failed: %+v", logId, err)
		return err
	}

	d.SetId(namespace)

	return resourceTencentCloudScfNamespaceRead(d, m)
}

func resourceTencentCloudScfNamespaceRead(d *schema.ResourceData, m interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_scf_namespace.read")()
	defer tccommon.InconsistentCheck(d, m)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := ScfService{client: m.(tccommon.ProviderMeta).GetAPIV3Conn()}

	id := d.Id()

	namespace, err := service.DescribeNamespace(ctx, id)
	if err != nil {
		log.Printf("[CRITAL]%s read namespace failed: %+v", logId, err)
		return err
	}

	if namespace == nil {
		d.SetId("")
		return nil
	}

	_ = d.Set("namespace", namespace.Name)
	_ = d.Set("description", namespace.Description)
	_ = d.Set("create_time", namespace.AddTime)
	_ = d.Set("modify_time", namespace.ModTime)
	_ = d.Set("type", namespace.Type)

	return nil
}

func resourceTencentCloudScfNamespaceUpdate(d *schema.ResourceData, m interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_scf_namespace.update")()
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := ScfService{client: m.(tccommon.ProviderMeta).GetAPIV3Conn()}

	if err := service.ModifyNamespace(ctx, d.Id(), d.Get("description").(string)); err != nil {
		log.Printf("[CRITAL]%s update namespace description failed: %+v", logId, err)
		return err
	}

	return resourceTencentCloudScfNamespaceRead(d, m)
}

func resourceTencentCloudScfNamespaceDelete(d *schema.ResourceData, m interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_scf_namespace.delete")()
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := ScfService{client: m.(tccommon.ProviderMeta).GetAPIV3Conn()}

	return service.DeleteNamespace(ctx, d.Id())
}
