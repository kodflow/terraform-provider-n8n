#!/bin/bash
# Script to add Params and Returns to all functions missing them

cd /workspace/src

# Function to add doc to NewXxxDataSource functions
add_datasource_factory_doc() {
    local file="$1"

    # Add doc for NewXxxDataSource functions
    sed -i '/^func New.*DataSource() datasource\.DataSource {$/ {
        i\
//\
// Params:\
//   - None\
//\
// Returns:\
//   - datasource.DataSource: New data source instance
    }' "$file"
}

# Function to add doc to Metadata methods
add_metadata_doc() {
    local file="$1"

    sed -i '/^func (d \*.*) Metadata(ctx context\.Context, req datasource\.MetadataRequest, resp \*datasource\.MetadataResponse) {$/ {
        i\
//\
// Params:\
//   - ctx: The context for the request\
//   - req: Metadata request with provider information\
//   - resp: Metadata response to populate\
//\
// Returns:\
//   - None
    }' "$file"
}

# Function to add doc to Schema methods
add_schema_doc() {
    local file="$1"

    sed -i '/^func (d \*.*) Schema(ctx context\.Context, req datasource\.SchemaRequest, resp \*datasource\.SchemaResponse) {$/ {
        i\
//\
// Params:\
//   - ctx: The context for the request\
//   - req: Schema request\
//   - resp: Schema response to populate\
//\
// Returns:\
//   - None
    }' "$file"
}

# Function to add doc to Configure methods
add_configure_doc() {
    local file="$1"

    sed -i '/^func (d \*.*) Configure(ctx context\.Context, req datasource\.ConfigureRequest, resp \*datasource\.ConfigureResponse) {$/ {
        i\
//\
// Params:\
//   - ctx: The context for the request\
//   - req: Configure request with provider data\
//   - resp: Configure response for diagnostics\
//\
// Returns:\
//   - None
    }' "$file"
}

# Function to add doc to Read methods
add_read_doc() {
    local file="$1"

    sed -i '/^func (d \*.*) Read(ctx context\.Context, req datasource\.ReadRequest, resp \*datasource\.ReadResponse) {$/ {
        i\
//\
// Params:\
//   - ctx: The context for the request\
//   - req: Read request with configuration\
//   - resp: Read response to populate\
//\
// Returns:\
//   - None
    }' "$file"
}

# Fix all datasource files
for file in ./internal/provider/datasources/*.go; do
    echo "Processing $file..."
    add_datasource_factory_doc "$file"
    add_metadata_doc "$file"
    add_schema_doc "$file"
    add_configure_doc "$file"
    add_read_doc "$file"
done

# Similar for resources
add_resource_factory_doc() {
    local file="$1"

    sed -i '/^func New.*Resource() resource\.Resource {$/ {
        i\
//\
// Params:\
//   - None\
//\
// Returns:\
//   - resource.Resource: New resource instance
    }' "$file"
}

add_resource_metadata_doc() {
    local file="$1"

    sed -i '/^func (r \*.*) Metadata(ctx context\.Context, req resource\.MetadataRequest, resp \*resource\.MetadataResponse) {$/ {
        i\
//\
// Params:\
//   - ctx: The context for the request\
//   - req: Metadata request with provider information\
//   - resp: Metadata response to populate\
//\
// Returns:\
//   - None
    }' "$file"
}

add_resource_schema_doc() {
    local file="$1"

    sed -i '/^func (r \*.*) Schema(ctx context\.Context, req resource\.SchemaRequest, resp \*resource\.SchemaResponse) {$/ {
        i\
//\
// Params:\
//   - ctx: The context for the request\
//   - req: Schema request\
//   - resp: Schema response to populate\
//\
// Returns:\
//   - None
    }' "$file"
}

add_resource_configure_doc() {
    local file="$1"

    sed -i '/^func (r \*.*) Configure(ctx context\.Context, req resource\.ConfigureRequest, resp \*resource\.ConfigureResponse) {$/ {
        i\
//\
// Params:\
//   - ctx: The context for the request\
//   - req: Configure request with provider data\
//   - resp: Configure response for diagnostics\
//\
// Returns:\
//   - None
    }' "$file"
}

add_resource_create_doc() {
    local file="$1"

    sed -i '/^func (r \*.*) Create(ctx context\.Context, req resource\.CreateRequest, resp \*resource\.CreateResponse) {$/ {
        i\
//\
// Params:\
//   - ctx: The context for the request\
//   - req: Create request with plan data\
//   - resp: Create response to populate\
//\
// Returns:\
//   - None
    }' "$file"
}

add_resource_read_doc() {
    local file="$1"

    sed -i '/^func (r \*.*) Read(ctx context\.Context, req resource\.ReadRequest, resp \*resource\.ReadResponse) {$/ {
        i\
//\
// Params:\
//   - ctx: The context for the request\
//   - req: Read request with state data\
//   - resp: Read response to populate\
//\
// Returns:\
//   - None
    }' "$file"
}

add_resource_update_doc() {
    local file="$1"

    sed -i '/^func (r \*.*) Update(ctx context\.Context, req resource\.UpdateRequest, resp \*resource\.UpdateResponse) {$/ {
        i\
//\
// Params:\
//   - ctx: The context for the request\
//   - req: Update request with plan data\
//   - resp: Update response to populate\
//\
// Returns:\
//   - None
    }' "$file"
}

add_resource_delete_doc() {
    local file="$1"

    sed -i '/^func (r \*.*) Delete(ctx context\.Context, req resource\.DeleteRequest, resp \*resource\.DeleteResponse) {$/ {
        i\
//\
// Params:\
//   - ctx: The context for the request\
//   - req: Delete request with state data\
//   - resp: Delete response to populate\
//\
// Returns:\
//   - None
    }' "$file"
}

# Fix all resource files
for file in ./internal/provider/resources/*.go; do
    echo "Processing $file..."
    add_resource_factory_doc "$file"
    add_resource_metadata_doc "$file"
    add_resource_schema_doc "$file"
    add_resource_configure_doc "$file"
    add_resource_create_doc "$file"
    add_resource_read_doc "$file"
    add_resource_update_doc "$file"
    add_resource_delete_doc "$file"
done

echo "Done!"
