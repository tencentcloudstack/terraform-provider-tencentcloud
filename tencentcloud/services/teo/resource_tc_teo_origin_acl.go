package teo

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	teov20220901 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudTeoOriginAcl() *schema.Resource {
	return &schema.Resource{
		Create: ResourceTencentCloudTeoOriginAclCreate,
		Read:   ResourceTencentCloudTeoOriginAclRead,
		Update: ResourceTencentCloudTeoOriginAclUpdate,
		Delete: ResourceTencentCloudTeoOriginAclDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"zone_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Specifies the site ID.",
			},

			"l7_hosts": {
				Type:        schema.TypeSet,
				Optional:    true,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "The list of L7 acceleration domains that require enabling the origin ACLs. This list must be empty when the request parameter L7EnableMode is set to 'all'.",
			},

			"l4_proxy_ids": {
				Type:        schema.TypeSet,
				Optional:    true,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "he list of L4 proxy Instances that require enabling origin ACLs. This list must be empty when the request parameter L4EnableMode is set to 'all'.",
			},
		},
	}
}

func ResourceTencentCloudTeoOriginAclCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_origin_acl.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = TeoService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		request = teov20220901.NewEnableOriginACLRequest()
		zoneId  string
	)

	if v, ok := d.GetOk("zone_id"); ok {
		request.ZoneId = helper.String(v.(string))
		zoneId = v.(string)
	}

	tmpL7Hosts := make([]interface{}, 0)
	tmpL4ProxyIds := make([]interface{}, 0)
	if v, ok := d.GetOk("l7_hosts"); ok {
		l7Hosts := v.(*schema.Set).List()
		if len(l7Hosts) > 1 {
			l7Hosts = v.(*schema.Set).List()[:1]
			tmpL7Hosts = v.(*schema.Set).List()[1:]
		}

		for i := range l7Hosts {
			if v, ok := l7Hosts[i].(string); ok && v != "" {
				l7Host := l7Hosts[i].(string)
				request.L7Hosts = append(request.L7Hosts, &l7Host)
			}
		}
	}

	if v, ok := d.GetOk("l4_proxy_ids"); ok {
		l4ProxyIds := v.(*schema.Set).List()
		if len(l4ProxyIds) > 1 {
			l4ProxyIds = v.(*schema.Set).List()[:1]
			tmpL4ProxyIds = v.(*schema.Set).List()[1:]
		}

		for i := range l4ProxyIds {
			if v, ok := l4ProxyIds[i].(string); ok && v != "" {
				l4ProxyId := l4ProxyIds[i].(string)
				request.L4ProxyIds = append(request.L4ProxyIds, &l4ProxyId)
			}
		}
	}

	request.L7EnableMode = helper.String("specific")
	request.L4EnableMode = helper.String("specific")
	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoV20220901Client().EnableOriginACLWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Enable teo origin acl failed, Response is nil."))
		}

		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s enable teo origin acl failed, reason:%+v", logId, err)
		return err
	}

	// wait
	err = service.WaitTeoOriginACLById(ctx, d.Timeout(schema.TimeoutCreate), zoneId, "online")
	if err != nil {
		return err
	}

	if len(tmpL7Hosts) > 0 || len(tmpL4ProxyIds) > 0 {
		batchSize := 200
		for i := 0; i < len(tmpL7Hosts); i += batchSize {
			end := i + batchSize
			if end > len(tmpL7Hosts) {
				end = len(tmpL7Hosts)
			}

			batchHosts := tmpL7Hosts[i:end]
			request := teov20220901.NewModifyOriginACLRequest()
			request.ZoneId = &zoneId
			request.OriginACLEntities = append(request.OriginACLEntities, &teov20220901.OriginACLEntity{
				Type:          helper.String("l7"),
				Instances:     helper.InterfacesStringsPoint(batchHosts),
				OperationMode: helper.String("enable"),
			})

			err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
				result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoV20220901Client().ModifyOriginACLWithContext(ctx, request)
				if e != nil {
					return tccommon.RetryError(e)
				} else {
					log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
				}

				if result == nil || result.Response == nil {
					return resource.NonRetryableError(fmt.Errorf("Modify teo origin acl failed, Response is nil."))
				}

				return nil
			})

			if err != nil {
				log.Printf("[CRITAL]%s modify teo origin acl (L7 batch %d) failed, reason:%+v", logId, i/batchSize, err)
				return err
			}

			// wait
			err = service.WaitTeoOriginACLById(ctx, d.Timeout(schema.TimeoutCreate), zoneId, "online")
			if err != nil {
				return err
			}
		}

		for i := 0; i < len(tmpL4ProxyIds); i += batchSize {
			end := i + batchSize
			if end > len(tmpL4ProxyIds) {
				end = len(tmpL4ProxyIds)
			}

			batchProxyIds := tmpL4ProxyIds[i:end]
			request := teov20220901.NewModifyOriginACLRequest()
			request.ZoneId = &zoneId
			request.OriginACLEntities = append(request.OriginACLEntities, &teov20220901.OriginACLEntity{
				Type:          helper.String("l4"),
				Instances:     helper.InterfacesStringsPoint(batchProxyIds),
				OperationMode: helper.String("enable"),
			})

			err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
				result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoV20220901Client().ModifyOriginACLWithContext(ctx, request)
				if e != nil {
					return tccommon.RetryError(e)
				} else {
					log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
				}

				if result == nil || result.Response == nil {
					return resource.NonRetryableError(fmt.Errorf("Modify teo origin acl failed, Response is nil."))
				}

				return nil
			})

			if err != nil {
				log.Printf("[CRITAL]%s modify teo origin acl (L4 batch %d) failed, reason:%+v", logId, i/batchSize, err)
				return err
			}

			// wait
			err = service.WaitTeoOriginACLById(ctx, d.Timeout(schema.TimeoutCreate), zoneId, "online")
			if err != nil {
				return err
			}
		}
	}

	d.SetId(zoneId)
	return ResourceTencentCloudTeoOriginAclRead(d, meta)
}

func ResourceTencentCloudTeoOriginAclRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_origin_acl.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = TeoService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		zoneId  = d.Id()
	)

	respData, err := service.DescribeTeoOriginACLById(ctx, zoneId)
	if err != nil {
		return err
	}

	if respData == nil {
		log.Printf("[WARN]%s resource `tencentcloud_teo_origin_acl` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		d.SetId("")
		return nil
	}

	_ = d.Set("zone_id", zoneId)

	if respData.L7Hosts != nil {
		tmpList := make([]string, 0, len(respData.L7Hosts))
		for _, item := range respData.L7Hosts {
			tmpList = append(tmpList, *item)
		}

		_ = d.Set("l7_hosts", tmpList)
	}

	if respData.L4ProxyIds != nil {
		tmpList := make([]string, 0, len(respData.L4ProxyIds))
		for _, item := range respData.L4ProxyIds {
			tmpList = append(tmpList, *item)
		}

		_ = d.Set("l4_proxy_ids", tmpList)
	}

	return nil
}

func ResourceTencentCloudTeoOriginAclUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_origin_acl.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId          = tccommon.GetLogId(tccommon.ContextNil)
		ctx            = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service        = TeoService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		zoneId         = d.Id()
		l7List, l4List []*teov20220901.OriginACLEntity
	)

	if d.HasChange("l7_hosts") {
		o, n := d.GetChange("l7_hosts")
		os := o.(*schema.Set)
		ns := n.(*schema.Set)
		add := ns.Difference(os).List()
		remove := os.Difference(ns).List()
		if len(add) > 0 {
			l7List = append(l7List, &teov20220901.OriginACLEntity{
				Type:          helper.String("l7"),
				Instances:     helper.InterfacesStringsPoint(add),
				OperationMode: helper.String("enable"),
			})
		}

		if len(remove) > 0 {
			l7List = append(l7List, &teov20220901.OriginACLEntity{
				Type:          helper.String("l7"),
				Instances:     helper.InterfacesStringsPoint(remove),
				OperationMode: helper.String("disable"),
			})
		}
	}

	if d.HasChange("l4_proxy_ids") {
		o, n := d.GetChange("l4_proxy_ids")
		os := o.(*schema.Set)
		ns := n.(*schema.Set)
		add := ns.Difference(os).List()
		remove := os.Difference(ns).List()
		if len(add) > 0 {
			l4List = append(l4List, &teov20220901.OriginACLEntity{
				Type:          helper.String("l4"),
				Instances:     helper.InterfacesStringsPoint(add),
				OperationMode: helper.String("enable"),
			})
		}

		if len(remove) > 0 {
			l4List = append(l4List, &teov20220901.OriginACLEntity{
				Type:          helper.String("l4"),
				Instances:     helper.InterfacesStringsPoint(remove),
				OperationMode: helper.String("disable"),
			})
		}
	}

	if len(l7List) > 0 || len(l4List) > 0 {
		batchSize := 200
		for i := 0; i < len(l7List); i += batchSize {
			end := i + batchSize
			if end > len(l7List) {
				end = len(l7List)
			}

			batchHosts := l7List[i:end]
			request := teov20220901.NewModifyOriginACLRequest()
			request.ZoneId = &zoneId
			request.OriginACLEntities = append(request.OriginACLEntities, batchHosts...)
			err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
				result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoV20220901Client().ModifyOriginACLWithContext(ctx, request)
				if e != nil {
					return tccommon.RetryError(e)
				} else {
					log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
				}

				if result == nil || result.Response == nil {
					return resource.NonRetryableError(fmt.Errorf("Modify teo origin acl failed, Response is nil."))
				}

				return nil
			})

			if err != nil {
				log.Printf("[CRITAL]%s modify teo origin acl (L7 batch %d) failed, reason:%+v", logId, i/batchSize, err)
				return err
			}

			// wait
			err = service.WaitTeoOriginACLById(ctx, d.Timeout(schema.TimeoutCreate), zoneId, "online")
			if err != nil {
				return err
			}
		}

		for i := 0; i < len(l4List); i += batchSize {
			end := i + batchSize
			if end > len(l4List) {
				end = len(l4List)
			}

			batchProxyIds := l4List[i:end]
			request := teov20220901.NewModifyOriginACLRequest()
			request.ZoneId = &zoneId
			request.OriginACLEntities = append(request.OriginACLEntities, batchProxyIds...)
			err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
				result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoV20220901Client().ModifyOriginACLWithContext(ctx, request)
				if e != nil {
					return tccommon.RetryError(e)
				} else {
					log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
				}

				if result == nil || result.Response == nil {
					return resource.NonRetryableError(fmt.Errorf("Modify teo origin acl failed, Response is nil."))
				}

				return nil
			})

			if err != nil {
				log.Printf("[CRITAL]%s modify teo origin acl (L4 batch %d) failed, reason:%+v", logId, i/batchSize, err)
				return err
			}

			// wait
			err = service.WaitTeoOriginACLById(ctx, d.Timeout(schema.TimeoutCreate), zoneId, "online")
			if err != nil {
				return err
			}
		}
	}

	return ResourceTencentCloudTeoOriginAclRead(d, meta)
}

func ResourceTencentCloudTeoOriginAclDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_origin_acl.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = TeoService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		request = teov20220901.NewDisableOriginACLRequest()
		zoneId  = d.Id()
	)

	request.ZoneId = &zoneId
	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoV20220901Client().DisableOriginACLWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s delete teo origin acl failed, reason:%+v", logId, err)
		return err
	}

	// wait
	err = service.WaitTeoOriginACLById(ctx, d.Timeout(schema.TimeoutDelete), zoneId, "offline")
	if err != nil {
		return err
	}

	return nil
}
