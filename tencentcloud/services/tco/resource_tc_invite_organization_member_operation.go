package tco

import (
	"context"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	organization "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/organization/v20210331"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudInviteOrganizationMemberOperation() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudInviteOrganizationMemberOperationCreate,
		Read:   resourceTencentCloudInviteOrganizationMemberOperationRead,
		Delete: resourceTencentCloudInviteOrganizationMemberOperationDelete,
		Schema: map[string]*schema.Schema{
			"member_uin": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: "Invited account Uin.",
			},

			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Member name. The maximum length is 25 characters and supports English letters, numbers, Chinese characters, symbols `+`, `@`, `&`, `.`, `[`, `]`, `-`, `:`, `,` and enumeration comma.",
			},

			"policy_type": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Relationship strategies. Value taken: Financial.",
			},

			"permission_ids": {
				Type:        schema.TypeSet,
				Required:    true,
				ForceNew:    true,
				Description: "List of member financial authority IDs. Values: 1-View bill, 2-View balance, 3-Fund transfer, 4-Consolidated disbursement, 5-Invoice, 6-Benefit inheritance, 7-Proxy payment, 1 and 2 must be default.",
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},

			"node_id": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: "Node ID of the member's department.",
			},

			"remark": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Remark.",
			},

			"is_allow_quit": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Whether to allow members to withdraw. Allow: Allow, Disallow: Denied.",
			},

			"pay_uin": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Payer Uin. Member needs to pay on behalf of.",
			},

			"relation_auth_name": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Name of the real-name subject of mutual trust.",
			},

			"auth_file": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				Description: "List of supporting documents of mutual trust entities.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "File name.",
						},
						"url": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "File path.",
						},
					},
				},
			},

			"tags": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				Description: "List of member tags. Maximum 10.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"tag_key": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Tag key.",
						},
						"tag_value": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Tag value.",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudInviteOrganizationMemberOperationCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_invite_organization_member_operation.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	var (
		memberUin string
		request   = organization.NewInviteOrganizationMemberRequest()
		response  = organization.NewInviteOrganizationMemberResponse()
	)

	if v, ok := d.GetOkExists("member_uin"); ok {
		memberUin = strconv.Itoa(v.(int))
		request.MemberUin = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("name"); ok {
		request.Name = helper.String(v.(string))
	}

	if v, ok := d.GetOk("policy_type"); ok {
		request.PolicyType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("permission_ids"); ok {
		permissionIdsSet := v.(*schema.Set).List()
		for i := range permissionIdsSet {
			permissionIds := permissionIdsSet[i].(int)
			request.PermissionIds = append(request.PermissionIds, helper.IntUint64(permissionIds))
		}
	}

	if v, ok := d.GetOkExists("node_id"); ok {
		request.NodeId = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("remark"); ok {
		request.Remark = helper.String(v.(string))
	}

	if v, ok := d.GetOk("is_allow_quit"); ok {
		request.IsAllowQuit = helper.String(v.(string))
	}

	if v, ok := d.GetOk("pay_uin"); ok {
		request.PayUin = helper.String(v.(string))
	}

	if v, ok := d.GetOk("relation_auth_name"); ok {
		request.RelationAuthName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("auth_file"); ok {
		for _, item := range v.([]interface{}) {
			authFileMap := item.(map[string]interface{})
			authRelationFile := organization.AuthRelationFile{}
			if v, ok := authFileMap["name"]; ok {
				authRelationFile.Name = helper.String(v.(string))
			}
			if v, ok := authFileMap["url"]; ok {
				authRelationFile.Url = helper.String(v.(string))
			}
			request.AuthFile = append(request.AuthFile, &authRelationFile)
		}
	}

	if v, ok := d.GetOk("tags"); ok {
		for _, item := range v.([]interface{}) {
			tagsMap := item.(map[string]interface{})
			tag := organization.Tag{}
			if v, ok := tagsMap["tag_key"]; ok {
				tag.TagKey = helper.String(v.(string))
			}
			if v, ok := tagsMap["tag_value"]; ok {
				tag.TagValue = helper.String(v.(string))
			}
			request.Tags = append(request.Tags, &tag)
		}
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseOrganizationClient().InviteOrganizationMemberWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create invite organization member operation failed, reason:%+v", logId, err)
		return err
	}

	_ = response

	d.SetId(memberUin)

	return resourceTencentCloudInviteOrganizationMemberOperationRead(d, meta)
}

func resourceTencentCloudInviteOrganizationMemberOperationRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_invite_organization_member_operation.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudInviteOrganizationMemberOperationDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_invite_organization_member_operation.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
