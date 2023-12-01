/*
Provides a resource to create a cvm program_fpga_image

Example Usage

```hcl
resource "tencentcloud_cvm_program_fpga_image" "program_fpga_image" {
  instance_id = "ins-xxxxxx"
  fpga_url = ""
  dbd_fs = ""
}
```

Import

cvm program_fpga_image can be imported using the id, e.g.

```
terraform import tencentcloud_cvm_program_fpga_image.program_fpga_image program_fpga_image_id
```
*/
package tencentcloud

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudCvmProgramFpgaImage() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCvmProgramFpgaImageCreate,
		Read:   resourceTencentCloudCvmProgramFpgaImageRead,
		Delete: resourceTencentCloudCvmProgramFpgaImageDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "The ID information of the instance.",
			},

			"fpga_url": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "COS URL address of the FPGA image file.",
			},

			"dbd_fs": {
				Optional: true,
				ForceNew: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "The DBDF number of the FPGA card on the instance, if left blank, the FPGA image will be burned to all FPGA cards owned by the instance by default.",
			},

			"dry_run": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeBool,
				Description: "Trial run, will not perform the actual burning action, the default is False.",
			},
		},
	}
}

func resourceTencentCloudCvmProgramFpgaImageCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cvm_program_fpga_image.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request = cvm.NewProgramFpgaImageRequest()
	)
	instanceId := d.Get("instance_id").(string)
	request.InstanceId = helper.String(instanceId)

	if v, ok := d.GetOk("fpga_url"); ok {
		request.FPGAUrl = helper.String(v.(string))
	}

	if v, ok := d.GetOk("dbd_fs"); ok {
		dBDFsSet := v.(*schema.Set).List()
		for i := range dBDFsSet {
			dBDFs := dBDFsSet[i].(string)
			request.DBDFs = append(request.DBDFs, &dBDFs)
		}
	}

	if v, _ := d.GetOk("dry_run"); v != nil {
		request.DryRun = helper.Bool(v.(bool))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseCvmClient().ProgramFpgaImage(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate cvm programFpgaImage failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(instanceId)

	return resourceTencentCloudCvmProgramFpgaImageRead(d, meta)
}

func resourceTencentCloudCvmProgramFpgaImageRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cvm_program_fpga_image.read")()
	defer inconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudCvmProgramFpgaImageDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cvm_program_fpga_image.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
