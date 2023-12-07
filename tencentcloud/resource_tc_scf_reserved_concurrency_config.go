package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	scf "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/scf/v20180416"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudScfReservedConcurrencyConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudScfReservedConcurrencyConfigCreate,
		Read:   resourceTencentCloudScfReservedConcurrencyConfigRead,
		Delete: resourceTencentCloudScfReservedConcurrencyConfigDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"function_name": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Specifies the function of which you want to configure the reserved quota.",
			},

			"reserved_concurrency_mem": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeInt,
				Description: "Reserved memory quota of the function. Note: the upper limit for the total reserved quota of the function is the user's total concurrency memory minus 12800.",
			},

			"namespace": {
				Optional:    true,
				ForceNew:    true,
				Default:     "default",
				Type:        schema.TypeString,
				Description: "Function namespace. Default value: default.",
			},
		},
	}
}

func resourceTencentCloudScfReservedConcurrencyConfigCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_scf_reserved_concurrency_config.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request      = scf.NewPutReservedConcurrencyConfigRequest()
		namespace    string
		functionName string
	)
	if v, ok := d.GetOk("function_name"); ok {
		functionName = v.(string)
		request.FunctionName = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("reserved_concurrency_mem"); ok {
		request.ReservedConcurrencyMem = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("namespace"); ok {
		namespace = v.(string)
		request.Namespace = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseScfClient().PutReservedConcurrencyConfig(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create scf ReservedConcurrencyConfig failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(namespace + FILED_SP + functionName)

	return resourceTencentCloudScfReservedConcurrencyConfigRead(d, meta)
}

func resourceTencentCloudScfReservedConcurrencyConfigRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_scf_reserved_concurrency_config.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := ScfService{client: meta.(*TencentCloudClient).apiV3Conn}

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	namespace := idSplit[0]
	functionName := idSplit[1]

	reservedConcurrencyConfig, err := service.DescribeScfReservedConcurrencyConfigById(ctx, namespace, functionName)
	if err != nil {
		return err
	}

	if reservedConcurrencyConfig == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `ScfReservedConcurrencyConfig` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	_ = d.Set("namespace", namespace)
	_ = d.Set("function_name", functionName)

	if reservedConcurrencyConfig.Response.ReservedMem != nil {
		_ = d.Set("reserved_concurrency_mem", reservedConcurrencyConfig.Response.ReservedMem)
	}

	return nil
}

func resourceTencentCloudScfReservedConcurrencyConfigDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_scf_reserved_concurrency_config.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := ScfService{client: meta.(*TencentCloudClient).apiV3Conn}
	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	namespace := idSplit[0]
	functionName := idSplit[1]

	if err := service.DeleteScfReservedConcurrencyConfigById(ctx, namespace, functionName); err != nil {
		return err
	}

	return nil
}
