// Package metaresources provides framework-side resource implementations
// for the "meta" product family — i.e. resources that are cross-product or
// not bound to any specific cloud product.
//
// Currently includes:
//   - tencentcloud_local_note: a pure in-memory local resource that
//     demonstrates the full framework resource CRUD + Configure pipeline
//     and does not call any cloud API.
//
// Wiring: append metaresources.NewLocalNoteResource to
// frameworkResources() in tencentcloud/framework/registry.go to expose the
// resource through the framework provider.
package metaresources

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"sync"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/sharedmeta"
)

// Compile-time assertions: localNoteResource implements
// resource.Resource and resource.ResourceWithConfigure, so it can be
// registered with the framework.
var (
	_ resource.Resource              = &localNoteResource{}
	_ resource.ResourceWithConfigure = &localNoteResource{}
)

// noteRecord is the in-memory representation of a single note.
type noteRecord struct {
	Title       string
	Content     string
	LastUpdated string
}

// notesStore is the process-wide in-memory store. key=resource id,
// value=noteRecord.
//
// We use sync.Map rather than a bare map to keep behaviour thread-safe
// under concurrent acceptance tests / function reference invocations.
var notesStore sync.Map

// NewLocalNoteResource is the framework resource factory expected by the
// registration convention. It is referenced from the frameworkResources()
// aggregator slice in tencentcloud/framework/registry.go.
func NewLocalNoteResource() resource.Resource {
	return &localNoteResource{}
}

// localNoteResource is a **pure local in-memory** framework resource. It
// does not depend on any remote cloud API; every CRUD operation only
// reads from / writes to the in-process notesStore.
//
// Purpose: serve as the in-repo reference implementation of the framework
// resource type, covering the full Metadata / Schema / Configure /
// Create / Read / Update / Delete lifecycle so that future business
// resources have a concrete template to follow.
type localNoteResource struct {
	// client is fetched from *sharedmeta.ProviderMeta during Configure.
	// CRUD on this resource does not actually use it; the field exists
	// solely to demonstrate "how to obtain the shared client".
	client interface{}
}

// localNoteModel is the Go counterpart of this resource's schema.
type localNoteModel struct {
	ID          types.String `tfsdk:"id"`
	Title       types.String `tfsdk:"title"`
	Content     types.String `tfsdk:"content"`
	LastUpdated types.String `tfsdk:"last_updated"`
}

// Metadata returns this resource's Terraform type name
// tencentcloud_local_note.
func (r *localNoteResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_local_note"
}

// Schema declares the resource's attribute set.
func (r *localNoteResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "An in-memory note managed entirely inside the provider process. " +
			"This resource does NOT call any TencentCloud API; it serves as a reference " +
			"implementation of a framework resource (Metadata / Schema / Configure / CRUD).",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "Auto-generated, immutable identifier of the note.",
			},
			"title": schema.StringAttribute{
				Required:    true,
				Description: "Human-readable title of the note.",
			},
			"content": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Free-form content of the note. Defaults to empty string when not specified.",
			},
			"last_updated": schema.StringAttribute{
				Computed:    true,
				Description: "RFC3339 timestamp of the last successful Create/Update.",
			},
		},
	}
}

// Configure retrieves the shared client injected by the framework provider
// during its own Configure phase. As with the data source, a nil
// ProviderData is silently ignored, while a type mismatch appends a
// diagnostic.
func (r *localNoteResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	meta, ok := req.ProviderData.(*sharedmeta.ProviderMeta)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Provider Data Type",
			"Expected *sharedmeta.ProviderMeta, please report this issue to the provider maintainers.",
		)
		return
	}
	r.client = meta.Client
}

// Create inserts a new note into the in-memory store and returns the
// generated id and last_updated timestamp.
func (r *localNoteResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan localNoteModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := newNoteID()
	now := time.Now().UTC().Format(time.RFC3339)
	content := plan.Content.ValueString()

	notesStore.Store(id, noteRecord{
		Title:       plan.Title.ValueString(),
		Content:     content,
		LastUpdated: now,
	})

	plan.ID = types.StringValue(id)
	plan.Content = types.StringValue(content)
	plan.LastUpdated = types.StringValue(now)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Read fetches a note from the in-memory store. If the entry has been
// removed by another path, the caller (Terraform) is expected to mark the
// resource as drifted and recreate it.
func (r *localNoteResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state localNoteModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	raw, ok := notesStore.Load(state.ID.ValueString())
	if !ok {
		// Removed externally: drop the entry from state so Terraform
		// recreates the resource.
		resp.State.RemoveResource(ctx)
		return
	}
	rec := raw.(noteRecord)

	state.Title = types.StringValue(rec.Title)
	state.Content = types.StringValue(rec.Content)
	state.LastUpdated = types.StringValue(rec.LastUpdated)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Update overwrites the same-id record in the store with the new plan and
// refreshes last_updated.
func (r *localNoteResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan localNoteModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := plan.ID.ValueString()
	if _, ok := notesStore.Load(id); !ok {
		resp.Diagnostics.AddAttributeError(
			path.Root("id"),
			"Note Not Found",
			"The note referenced by id no longer exists in the in-memory store.",
		)
		return
	}

	now := time.Now().UTC().Format(time.RFC3339)
	content := plan.Content.ValueString()

	notesStore.Store(id, noteRecord{
		Title:       plan.Title.ValueString(),
		Content:     content,
		LastUpdated: now,
	})

	plan.Content = types.StringValue(content)
	plan.LastUpdated = types.StringValue(now)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Delete removes the target record from the store. After Delete returns,
// the framework calls resp.State.RemoveResource() automatically; this
// method does not need to invoke it directly.
func (r *localNoteResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state localNoteModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	notesStore.Delete(state.ID.ValueString())
}

// newNoteID generates a 16-byte hex note id. On failure it falls back to a
// timestamp, which is still unique within the process (the in-memory
// store does not require cross-process consistency).
func newNoteID() string {
	var b [16]byte
	if _, err := rand.Read(b[:]); err != nil {
		return "note-" + time.Now().UTC().Format("20060102150405.000000000")
	}
	return "note-" + hex.EncodeToString(b[:])
}
