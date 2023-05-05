/*
Provide a resource to create a SCF namespace.

Example Usage

```hcl
resource "tencentcloud_scf_namespace" "foo" {
  namespace = "ci-test-scf"
}
```

Import

SCF namespace can be imported, e.g.

```
$ terraform import tencentcloud_scf_function.test default
```
*/
package tencentcloud

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceTencentCloudScfNamespace() *schema.Resource {
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
	defer logElapsed("resource.tencentcloud_scf_namespace.create")()
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := ScfService{client: m.(*TencentCloudClient).apiV3Conn}

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
	defer logElapsed("resource.tencentcloud_scf_namespace.read")()
	defer inconsistentCheck(d, m)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := ScfService{client: m.(*TencentCloudClient).apiV3Conn}

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
	defer logElapsed("resource.tencentcloud_scf_namespace.update")()
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := ScfService{client: m.(*TencentCloudClient).apiV3Conn}

	if err := service.ModifyNamespace(ctx, d.Id(), d.Get("description").(string)); err != nil {
		log.Printf("[CRITAL]%s update namespace description failed: %+v", logId, err)
		return err
	}

	return resourceTencentCloudScfNamespaceRead(d, m)
}

func resourceTencentCloudScfNamespaceDelete(d *schema.ResourceData, m interface{}) error {
	defer logElapsed("resource.tencentcloud_scf_namespace.delete")()
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := ScfService{client: m.(*TencentCloudClient).apiV3Conn}

	return service.DeleteNamespace(ctx, d.Id())
}
