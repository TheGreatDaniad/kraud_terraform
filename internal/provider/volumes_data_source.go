package provider

import (
	"context"
	"fmt"
	"time"

	API "github.com/kraudcloud/cli/api"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func NewVolumesDataSource() datasource.DataSource {
	return &volumesDataSource{}
}

type KraudVolume struct {
	AID          types.String `json:"AID"`
	Class        types.String `json:"Class"`
	DeletionLock types.String `json:"DeletionLock,omitempty"`
	ExpiresAt    *time.Time   `json:"ExpiresAt,omitempty"`
	ID           types.String `json:"ID,omitempty"`
	IOPS         types.Int64  `json:"IOPS,omitempty"`
	Name         types.String `json:"Name"`
	Size         types.Int64  `json:"Size"`
	Version      types.String `json:"Version,omitempty"`
	Zone         types.String `json:"Zone,omitempty"`
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

	client, ok := req.ProviderData.(API.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *hashicups.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}
	d.client = client
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
			volumeState.ExpiresAt = volume.ExpiresAt
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

}
