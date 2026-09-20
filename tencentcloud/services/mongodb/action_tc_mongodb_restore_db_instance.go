package mongodb

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/action/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/framework/fw"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"

	mongodb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/mongodb/v20190725"
)

// _ action.Action is intentionally omitted: the factory function
// NewMongodbRestoreDbInstance returns action.Action, which already
// performs the implicit compile-time check at the return statement. Only the
// extended-interface assertion is kept here to verify the Configure method
// signature (promoted from fw.ActionWithConfigure) at compile time.
var _ action.ActionWithConfigure = &MongodbRestoreDbInstance{}

// NewMongodbRestoreDbInstance is the factory referenced by
// tencentcloud/framework/registry.go to register this action.
func NewMongodbRestoreDbInstance() action.Action {
	return &MongodbRestoreDbInstance{}
}

// MongodbRestoreDbInstance implements action.Action for
// tencentcloud_mongodb_restore_db_instance.
type MongodbRestoreDbInstance struct {
	fw.ActionWithConfigure
}

// MongodbRestoreDbInstanceModel maps the schema attributes.
type MongodbRestoreDbInstanceModel struct {
	InstanceId  types.String `tfsdk:"instance_id"`
	RestoreTime types.String `tfsdk:"restore_time"`
	Databases   types.List   `tfsdk:"databases"`
}

// MongodbRestoreDbDatabaseModel maps the databases nested block.
type MongodbRestoreDbDatabaseModel struct {
	Db          types.String `tfsdk:"db"`
	Collections types.List   `tfsdk:"collections"`
}

// MongodbRestoreDbCollectionModel maps the collections nested block.
type MongodbRestoreDbCollectionModel struct {
	OldCollection types.String `tfsdk:"old_collection"`
	NewCollection types.String `tfsdk:"new_collection"`
}

func (a *MongodbRestoreDbInstance) Metadata(_ context.Context, _ action.MetadataRequest, resp *action.MetadataResponse) {
	resp.TypeName = "tencentcloud_mongodb_restore_db_instance"
}

func (a *MongodbRestoreDbInstance) Schema(_ context.Context, _ action.SchemaRequest, resp *action.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provides an action to restore a MongoDB instance to a specified point in time via the RestoreDBInstance API. " +
			"This is a one-time operation; no cloud-side state is persisted after the action completes.",
		Attributes: map[string]schema.Attribute{
			"instance_id": schema.StringAttribute{
				Required:    true,
				Description: "Instance ID, e.g. `cmgo-xxxxxxxx`. Please log in to the MongoDB console and copy the instance ID from the instance list.",
			},
			"restore_time": schema.StringAttribute{
				Required:    true,
				Description: "Target point in time to restore. The time must be within the backup retention period of the instance. Format: `YYYY-MM-DD hh:mm:ss`.",
			},
		},
		Blocks: map[string]schema.Block{
			"databases": schema.ListNestedBlock{
				Description: "Database and collection information to restore.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"db": schema.StringAttribute{
							Required:    true,
							Description: "Database name.",
						},
					},
					Blocks: map[string]schema.Block{
						"collections": schema.ListNestedBlock{
							Description: "Collection information to restore.",
							NestedObject: schema.NestedBlockObject{
								Attributes: map[string]schema.Attribute{
									"old_collection": schema.StringAttribute{
										Required:    true,
										Description: "Original collection name to restore.",
									},
									"new_collection": schema.StringAttribute{
										Required:    true,
										Description: "Collection name after restore.",
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

// Invoke is called to run the logic of the action.
func (a *MongodbRestoreDbInstance) Invoke(ctx context.Context, req action.InvokeRequest, resp *action.InvokeResponse) {
	var data MongodbRestoreDbInstanceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if data.InstanceId.IsNull() || data.InstanceId.IsUnknown() || data.InstanceId.ValueString() == "" {
		resp.Diagnostics.AddError(
			"Missing instance_id",
			"instance_id is required and cannot be empty.",
		)
		return
	}

	if data.RestoreTime.IsNull() || data.RestoreTime.IsUnknown() || data.RestoreTime.ValueString() == "" {
		resp.Diagnostics.AddError(
			"Missing restore_time",
			"restore_time is required and cannot be empty.",
		)
		return
	}

	var databases []MongodbRestoreDbDatabaseModel
	if !data.Databases.IsNull() && !data.Databases.IsUnknown() {
		resp.Diagnostics.Append(data.Databases.ElementsAs(ctx, &databases, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}
	if len(databases) == 0 {
		resp.Diagnostics.AddError(
			"Missing databases",
			"databases is required and cannot be empty.",
		)
		return
	}

	var restoreDatabases []*mongodb.RestoreDatabases
	for _, db := range databases {
		if db.Db.IsNull() || db.Db.IsUnknown() || db.Db.ValueString() == "" {
			resp.Diagnostics.AddError(
				"Missing db",
				"db is required and cannot be empty in each databases block.",
			)
			return
		}

		var collections []MongodbRestoreDbCollectionModel
		if !db.Collections.IsNull() && !db.Collections.IsUnknown() {
			resp.Diagnostics.Append(db.Collections.ElementsAs(ctx, &collections, false)...)
			if resp.Diagnostics.HasError() {
				return
			}
		}
		if len(collections) == 0 {
			resp.Diagnostics.AddError(
				"Missing collections",
				"collections is required and cannot be empty in each databases block.",
			)
			return
		}

		var restoreCollections []*mongodb.RestoreCollection
		for _, c := range collections {
			if c.OldCollection.IsNull() || c.OldCollection.IsUnknown() || c.OldCollection.ValueString() == "" {
				resp.Diagnostics.AddError(
					"Missing old_collection",
					"old_collection is required and cannot be empty in each collections block.",
				)
				return
			}
			if c.NewCollection.IsNull() || c.NewCollection.IsUnknown() || c.NewCollection.ValueString() == "" {
				resp.Diagnostics.AddError(
					"Missing new_collection",
					"new_collection is required and cannot be empty in each collections block.",
				)
				return
			}
			oldCollection := c.OldCollection.ValueString()
			newCollection := c.NewCollection.ValueString()
			restoreCollections = append(restoreCollections, &mongodb.RestoreCollection{
				OldCollection: &oldCollection,
				NewCollection: &newCollection,
			})
		}

		dbName := db.Db.ValueString()
		restoreDatabases = append(restoreDatabases, &mongodb.RestoreDatabases{
			Db:          &dbName,
			Collections: restoreCollections,
		})
	}

	if a.Client() == nil {
		resp.Diagnostics.AddError(
			"Provider not configured",
			"The provider client is not available. This is unexpected; please report it.",
		)
		return
	}

	instanceId := data.InstanceId.ValueString()
	restoreTime := data.RestoreTime.ValueString()
	service := NewMongodbService(a.Client())

	flowId, err := service.RestoreDBInstance(ctx, instanceId, restoreTime, restoreDatabases)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("restoring mongodb instance (%s) to time (%s)", instanceId, restoreTime),
			err.Error(),
		)
		return
	}

	flowIdStr := helper.Int64ToStr(flowId)
	if err := service.DescribeAsyncRequestInfo(ctx, flowIdStr, 3*tccommon.ReadRetryTimeout); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("waiting for mongodb restore_db_instance task (%s) to complete", flowIdStr),
			err.Error(),
		)
		return
	}
}
