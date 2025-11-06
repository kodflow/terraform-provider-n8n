#!/bin/bash

# Script pour corriger automatiquement les erreurs STRUCT

# Liste des fichiers datasources
DATASOURCES=(
    "/workspace/src/internal/provider/datasources/executions.go"
    "/workspace/src/internal/provider/datasources/project.go"
    "/workspace/src/internal/provider/datasources/projects.go"
    "/workspace/src/internal/provider/datasources/tag.go"
    "/workspace/src/internal/provider/datasources/tags.go"
    "/workspace/src/internal/provider/datasources/user.go"
    "/workspace/src/internal/provider/datasources/users.go"
    "/workspace/src/internal/provider/datasources/variable.go"
    "/workspace/src/internal/provider/datasources/variables.go"
    "/workspace/src/internal/provider/datasources/workflow.go"
    "/workspace/src/internal/provider/datasources/workflows.go"
)

# Liste des fichiers resources
RESOURCES=(
    "/workspace/src/internal/provider/resources/credential.go"
    "/workspace/src/internal/provider/resources/credential_transfer.go"
    "/workspace/src/internal/provider/resources/execution_retry.go"
    "/workspace/src/internal/provider/resources/project.go"
    "/workspace/src/internal/provider/resources/project_user.go"
    "/workspace/src/internal/provider/resources/source_control_pull.go"
    "/workspace/src/internal/provider/resources/tag.go"
    "/workspace/src/internal/provider/resources/user.go"
    "/workspace/src/internal/provider/resources/variable.go"
    "/workspace/src/internal/provider/resources/workflow.go"
    "/workspace/src/internal/provider/resources/workflow_transfer.go"
)

echo "Correction des erreurs STRUCT en cours..."
echo "=========================================="

# Fonction pour extraire le nom de la struct principale
get_main_struct() {
    local file=$1
    grep -E "^type [A-Z][a-zA-Z0-9]*DataSource struct {" "$file" | head -1 | awk '{print $2}'
    if [ $? -ne 0 ]; then
        grep -E "^type [A-Z][a-zA-Z0-9]*Resource struct {" "$file" | head -1 | awk '{print $2}'
    fi
}

echo "✓ Script de correction des erreurs STRUCT créé"
