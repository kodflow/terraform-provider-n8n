// Copyright (c) 2024 Florent (Kodflow). All rights reserved.
// Licensed under the Sustainable Use License 1.0
// See LICENSE in the project root for license information.

// Package audit implements audit report datasource functionality.
package audit

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/kodflow/terraform-provider-n8n/src/internal/provider/audit/models"
	"github.com/kodflow/terraform-provider-n8n/src/internal/provider/shared/client"
)

// Ensure AuditDataSource implements required interfaces.
var (
	_ datasource.DataSource              = &AuditDataSource{}
	_ AuditDataSourceInterface           = &AuditDataSource{}
	_ datasource.DataSourceWithConfigure = &AuditDataSource{}
)

// AuditDataSourceInterface defines the interface for AuditDataSource.
type AuditDataSourceInterface interface {
	datasource.DataSource
	Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse)
	Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse)
	Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse)
	Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse)
}

// AuditDataSource is a Terraform datasource that generates and returns an n8n security audit report.
// Each read triggers a POST to the n8n audit endpoint and returns the risk reports as JSON strings.
type AuditDataSource struct {
	// client is the N8n API client used for operations.
	client *client.N8nClient
}

// NewAuditDataSource creates a new AuditDataSource instance.
//
// Returns:
//   - *AuditDataSource: A new AuditDataSource instance
func NewAuditDataSource() (auditDataSource *AuditDataSource) {
	//: Return a new empty datasource instance.
	return &AuditDataSource{}
}

// NewAuditDataSourceWrapper creates a new AuditDataSource instance for Terraform.
// This wrapper function is used by the provider to maintain compatibility with the framework.
//
// Returns:
//   - datasource.DataSource: the wrapped AuditDataSource instance
func NewAuditDataSourceWrapper() (dataSource datasource.DataSource) {
	//: Return the wrapped datasource instance.
	return NewAuditDataSource()
}

// Metadata returns the data source type name.
//
// Params:
//   - ctx: The request context
//   - req: The metadata request containing provider type information
//   - resp: The metadata response to populate with type name
//
// Returns:
//   - None
func (d *AuditDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	//: Return early if context is cancelled.
	if ctx.Err() != nil {
		//: Return result.
		return
	}
	resp.TypeName = req.ProviderTypeName + "_audit"
}

// Schema defines the schema for the data source.
//
// Params:
//   - ctx: The request context
//   - req: The schema request from Terraform
//   - resp: The schema response to populate
//
// Returns:
//   - None
func (d *AuditDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	//: Return early if context is cancelled.
	if ctx.Err() != nil {
		//: Return result.
		return
	}
	//: Apply the schema definition when the request is the default schema request.
	if req == (datasource.SchemaRequest{}) {
		resp.Schema = schema.Schema{
			MarkdownDescription: "Generates an n8n security audit report. Each read triggers a new audit.",
			Attributes: map[string]schema.Attribute{
				"credentials_risk_report": schema.StringAttribute{
					MarkdownDescription: "Credentials risk report as JSON string",
					Computed:            true,
				},
				"database_risk_report": schema.StringAttribute{
					MarkdownDescription: "Database risk report as JSON string",
					Computed:            true,
				},
				"filesystem_risk_report": schema.StringAttribute{
					MarkdownDescription: "Filesystem risk report as JSON string",
					Computed:            true,
				},
				"nodes_risk_report": schema.StringAttribute{
					MarkdownDescription: "Nodes risk report as JSON string",
					Computed:            true,
				},
				"instance_risk_report": schema.StringAttribute{
					MarkdownDescription: "Instance risk report as JSON string",
					Computed:            true,
				},
			},
		}
	}
}

// Configure adds the provider configured client to the data source.
//
// Params:
//   - ctx: The request context
//   - req: The configure request containing provider data
//   - resp: The configure response to handle errors
//
// Returns:
//   - None
func (d *AuditDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	//: Return early if context is cancelled.
	if ctx.Err() != nil {
		//: Return result.
		return
	}
	//: Skip configuration when provider data is not yet available.
	if req.ProviderData == nil {
		//: Return early without error when provider data is nil.
		return
	}

	clientData, ok := req.ProviderData.(*client.N8nClient)
	//: Validate that the provider data is the expected N8nClient type.
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *client.N8nClient, got: %T", req.ProviderData),
		)
		//: Return after adding type mismatch error to diagnostics.
		return
	}

	d.client = clientData
}

// Read triggers an audit and refreshes the Terraform state with the report.
//
// Params:
//   - ctx: The request context
//   - req: The read request from Terraform
//   - resp: The read response to populate with data
//
// Returns:
//   - None
func (d *AuditDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data models.DataSource

	//: Return early if audit generation failed.
	if !d.executeAuditLogic(ctx, &data, resp) {
		//: Return after audit failure.
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// executeAuditLogic calls the n8n audit API and populates the data model.
//
// Params:
//   - ctx: The request context
//   - data: The data source model to populate
//   - resp: The read response
//
// Returns:
//   - bool: True if the audit succeeded, false otherwise
func (d *AuditDataSource) executeAuditLogic(ctx context.Context, data *models.DataSource, resp *datasource.ReadResponse) (ok bool) {
	auditReport, httpResp, err := d.client.APIClient.AuditAPI.AuditPost(ctx).Execute()
	//: Close the HTTP response body to avoid resource leaks.
	if httpResp != nil && httpResp.Body != nil {
		defer func() {
			//: Silently discard close error on response body.
			if closeErr := httpResp.Body.Close(); closeErr != nil {
				resp.Diagnostics.AddWarning("Failed to close response body", closeErr.Error())
			}
		}()
	}
	//: Return failure when the API call produced an error.
	if err != nil {
		resp.Diagnostics.AddError(
			"Error generating audit report",
			fmt.Sprintf("Could not generate audit report: %s\nHTTP Response: %v", err.Error(), httpResp),
		)
		//: Return false to signal audit failure.
		return false
	}

	//: Map the API response to the datasource model.
	mapAuditToDataSource(auditReport, data)

	//: Return true to signal successful audit.
	return true
}
