package provider

import (
	"context"
	"fmt"

	API "github.com/kraudcloud/cli/api"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func NewVolumesDataSource() datasource.DataSource {
	return &volumesDataSource{}
}

type KraudVolume struct {
	AID          types.String `json:"aid" tfsdk:"aid"`
	Class        types.String `json:"class" tfsdk:"class"`
	DeletionLock types.String `json:"deletion_lock,omitempty" tfsdk:"deletion_lock"`
	ExpiresAt    types.String `json:"expires_at,omitempty" tfsdk:"expires_at"`
	ID           types.String `json:"id,omitempty" tfsdk:"id"`
	IOPS         types.Int64  `json:"iops,omitempty" tfsdk:"iops"`
	Name         types.String `json:"name" tfsdk:"name"`
	Size         types.Int64  `json:"size" tfsdk:"size"`
	Version      types.String `json:"version,omitempty" tfsdk:"version"`
	Zone         types.String `json:"zone,omitempty" tfsdk:"zone"`
}

type volumesDataSourceModel struct {
	Volumes []KraudVolume `tfsdk:"volumes"`
}

type volumesDataSource struct {
	client API.Client
}

func (d *volumesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*API.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *API.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}
	d.client = *client
}

// Metadata returns the data source type name.
func (d *volumesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_volumes"
}

// Schema defines the schema for the data source.
func (d *volumesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"volumes": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"aid": schema.StringAttribute{
							Computed: true,
						},
						"class": schema.StringAttribute{
							Computed: true,
						},
						"deletion_lock": schema.StringAttribute{
							Computed: true,
							Optional: true,
						},
						"expires_at": schema.StringAttribute{
							Computed: true,
							Optional: true,
						},
						"id": schema.StringAttribute{
							Computed: true,
							Optional: true,
						},
						"iops": schema.NumberAttribute{
							Computed: true,
							Optional: true,
						},
						"name": schema.StringAttribute{
							Computed: true,
						},
						"size": schema.NumberAttribute{
							Computed: true,
						},
						"version": schema.StringAttribute{
							Computed: true,
							Optional: true,
						},
						"zone": schema.StringAttribute{
							Computed: true,
							Optional: true,
						},
					},
				},
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *volumesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state volumesDataSourceModel

	volumes, err := d.client.ListVolumes(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to read volumes",
			err.Error(),
		)
		return
	}
	fmt.Println(volumes)
	// Map response body to KraudVolume model
	for _, volume := range volumes.Items {
		volumeState := KraudVolume{
			AID:   types.StringValue(volume.AID),
			Class: types.StringValue(volume.Class),
			Name:  types.StringValue(volume.Name),
			Size:  types.Int64Value(int64(volume.Size)),
		}

		if volume.DeletionLock != nil {
			volumeState.DeletionLock = types.StringValue(*volume.DeletionLock)
		}

		if volume.ExpiresAt != nil {
			volumeState.ExpiresAt = types.StringValue(*volume.ID)
		}

		if volume.ID != nil {
			volumeState.ID = types.StringValue(*volume.ID)
		}

		if volume.IOPS != nil {
			volumeState.IOPS = types.Int64Value(int64(*volume.IOPS))
		}

		if volume.Version != nil {
			volumeState.Version = types.StringValue(*volume.Version)
		}

		if volume.Zone != nil {
			volumeState.Zone = types.StringValue(*volume.Zone)
		}

		state.Volumes = append(state.Volumes, volumeState)
	}
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}
