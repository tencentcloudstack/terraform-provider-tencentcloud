package tencentcloud

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	mariadb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/mariadb/v20170312"
)

func resourceTencentCloudMariadbInstanceConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudMariadbInstanceConfigCreate,
		Read:   resourceTencentCloudMariadbInstanceConfigRead,
		Update: resourceTencentCloudMariadbInstanceConfigUpdate,
		Delete: resourceTencentCloudMariadbInstanceConfigDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "instance id.",
			},
			"rs_access_strategy": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validateAllowedIntValue([]int{0, 1}),
				Description:  "RS proximity mode, 0- no strategy, 1- access to the nearest available zone.",
			},
			"extranet_access": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validateAllowedIntValue([]int{0, 1}),
				Description:  "External network status, 0-closed; 1- Opening; Default not enabled.",
			},
		},
	}
}

func resourceTencentCloudMariadbInstanceConfigCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mariadb_instance_config.create")()
	defer inconsistentCheck(d, meta)()

	var instanceId string

	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
	}

	d.SetId(instanceId)

	return resourceTencentCloudMariadbInstanceConfigUpdate(d, meta)
}

func resourceTencentCloudMariadbInstanceConfigRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mariadb_instance_config.read")()
	defer inconsistentCheck(d, meta)()

	var (
		logId      = getLogId(contextNil)
		ctx        = context.WithValue(context.TODO(), logIdKey, logId)
		service    = MariadbService{client: meta.(*TencentCloudClient).apiV3Conn}
		instanceId = d.Id()
	)

	dbDetail, err := service.DescribeDBInstanceDetailById(ctx, instanceId)
	if err != nil {
		return err
	}

	if dbDetail == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `MariadbInstanceConfig` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if dbDetail.InstanceId != nil {
		_ = d.Set("instance_id", dbDetail.InstanceId)
	}

	if dbDetail.RsAccessStrategy != nil {
		_ = d.Set("rs_access_strategy", dbDetail.RsAccessStrategy)
	}

	if dbDetail.WanStatus != nil {
		_ = d.Set("extranet_access", dbDetail.WanStatus)
	}

	return nil
}

func resourceTencentCloudMariadbInstanceConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mariadb_instance_config.update")()
	defer inconsistentCheck(d, meta)()

	var (
		logId      = getLogId(contextNil)
		ctx        = context.WithValue(context.TODO(), logIdKey, logId)
		service    = MariadbService{client: meta.(*TencentCloudClient).apiV3Conn}
		request    = mariadb.NewModifyAccountPrivilegesRequest()
		instanceId = d.Id()
	)

	needChange := false

	mutableArgs := []string{"rs_access_strategy", "extranet_access"}

	for _, v := range mutableArgs {
		if d.HasChange(v) {
			needChange = true
			break
		}
	}

	if needChange {
		// set rs_access_strategy
		if v, ok := d.GetOkExists("rs_access_strategy"); ok {
			rsAccessStrategy := int64(v.(int))
			if rsAccessStrategy == RSACCESSSTRATEGY_ENABLE {
				rsAccessStrategyRequest := mariadb.NewModifyRealServerAccessStrategyRequest()
				rsAccessStrategyRequest.InstanceId = &instanceId
				rsAccessStrategyRequest.RsAccessStrategy = &rsAccessStrategy

				err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
					result, e := meta.(*TencentCloudClient).apiV3Conn.UseMariadbClient().ModifyRealServerAccessStrategy(rsAccessStrategyRequest)
					if e != nil {
						return retryError(e)
					} else {
						log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
					}

					return nil
				})

				if err != nil {
					log.Printf("[CRITAL]%s operate mariadb accessStrategy failed, reason:%+v", logId, err)
					return err
				}
			}
		}

		// set ExtranetAccess
		if v, ok := d.GetOkExists("extranet_access"); ok {
			extranetAccess := v.(int)
			var extranetAccessFlowId int64
			if extranetAccess == ExtranetAccess_ENABLE {
				extranetAccessRequest := mariadb.NewOpenDBExtranetAccessRequest()
				extranetAccessRequest.InstanceId = &instanceId
				err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
					result, e := meta.(*TencentCloudClient).apiV3Conn.UseMariadbClient().OpenDBExtranetAccess(extranetAccessRequest)
					if e != nil {
						return retryError(e)
					} else {
						log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
					}

					extranetAccessFlowId = *result.Response.FlowId
					return nil
				})

				if err != nil {
					log.Printf("[CRITAL]%s operate mariadb openDBExtranetAccess failed, reason:%+v", logId, err)
					return err
				}

			} else {
				extranetAccessRequest := mariadb.NewCloseDBExtranetAccessRequest()
				extranetAccessRequest.InstanceId = &instanceId
				err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
					result, e := meta.(*TencentCloudClient).apiV3Conn.UseMariadbClient().CloseDBExtranetAccess(extranetAccessRequest)
					if e != nil {
						return retryError(e)
					} else {
						log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
					}

					extranetAccessFlowId = *result.Response.FlowId
					return nil
				})
				if err != nil {
					log.Printf("[CRITAL]%s operate mariadb closeDBExtranetAccess failed, reason:%+v", logId, err)
					return err
				}
			}

			// wait
			if extranetAccessFlowId != NONE_FLOW_TASK {
				err := resource.Retry(10*writeRetryTimeout, func() *resource.RetryError {
					result, e := service.DescribeFlowById(ctx, extranetAccessFlowId)
					if e != nil {
						return retryError(e)
					}

					if *result.Status == MARIADB_TASK_SUCCESS {
						return nil
					} else if *result.Status == MARIADB_TASK_RUNNING {
						return resource.RetryableError(fmt.Errorf("operate mariadb DBExtranetAccess status is running"))
					} else if *result.Status == MARIADB_TASK_FAIL {
						return resource.NonRetryableError(fmt.Errorf("operate mariadb DBExtranetAccess status is fail"))
					} else {
						e = fmt.Errorf("operate mariadb DBExtranetAccess status illegal")
						return resource.NonRetryableError(e)
					}
				})

				if err != nil {
					log.Printf("[CRITAL]%s operate mariadb DBExtranetAccess task failed, reason:%+v", logId, err)
					return err
				}
			}
		}
	}

	return resourceTencentCloudMariadbInstanceConfigRead(d, meta)
}

func resourceTencentCloudMariadbInstanceConfigDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mariadb_instance_config.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
