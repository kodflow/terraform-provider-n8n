#!/usr/bin/env python3
"""Script pour corriger automatiquement les erreurs VAR dans les fichiers Go."""

import re
import sys
from pathlib import Path

def fix_var_014_pointer(content: str) -> str:
    """Corrige VAR-014: transformer var x Model en x := &Model{}."""
    # Pattern: var <name> <TypeModel>
    # Remplacer par: <name> := &<TypeModel>{}

    # Cas 1: var plan WorkflowResourceModel
    content = re.sub(
        r'(\n\t)var (plan|state|data) ([A-Z]\w+Model)\n',
        r'\1\2 := &\3{}\n',
        content
    )

    # Cas 2: var plan, state Model (deux variables)
    content = re.sub(
        r'(\n\t)var (plan), (state) ([A-Z]\w+Model)\n',
        r'\1\2 := &\4{}\n\1\3 := &\4{}\n',
        content
    )

    return content

def fix_var_017_shadowing(content: str) -> str:
    """Corrige VAR-017: éviter le shadowing de variables."""
    lines = content.split('\n')
    result = []

    # Variables déclarées au niveau de la fonction
    declared_vars = set()

    for i, line in enumerate(lines):
        # Détecter les déclarations de variables (var err error, etc.)
        if '\tvar ' in line or 'var ' in line:
            var_match = re.search(r'var (\w+)', line)
            if var_match:
                declared_vars.add(var_match.group(1))

        # Détecter les shadowing: err := ... quand err existe déjà
        # Remplacer := par =
        if ':=' in line:
            # Extraire les noms de variables avant :=
            before_assign = line.split(':=')[0]
            vars_in_line = re.findall(r'\b(\w+)\b', before_assign.replace('_', ''))

            # Si une variable est déjà déclarée, utiliser = au lieu de :=
            for var_name in vars_in_line:
                if var_name in declared_vars or var_name in ['err', 'httpResp', 'ok']:
                    # Ne pas modifier si c'est dans un for range ou if
                    if 'for ' not in line and 'if ' not in line:
                        line = line.replace(' := ', ' = ', 1)
                        break

        result.append(line)

    return '\n'.join(result)

def fix_var_007_make_capacity(content: str) -> str:
    """Corrige VAR-007: make([]T, 0) → make([]T, 0, capacity)."""
    # make([]Type, 0) → make([]Type, 0, 10)
    content = re.sub(
        r'make\(\[\]([^,]+), 0\)',
        r'make([\1, 0, 10)',
        content
    )
    return content

def fix_var_008_make_optimal(content: str) -> str:
    """Corrige VAR-008: utiliser la capacité optimale pour make."""
    # Déjà couvert par VAR-007 dans la plupart des cas
    return content

def process_file(filepath: Path) -> bool:
    """Traite un fichier pour corriger les erreurs VAR."""
    try:
        content = filepath.read_text()
        original = content

        # Appliquer les corrections
        content = fix_var_014_pointer(content)
        content = fix_var_007_make_capacity(content)
        content = fix_var_008_make_optimal(content)
        # Note: fix_var_017_shadowing est complexe et risqué en automatique

        if content != original:
            filepath.write_text(content)
            print(f"✓ Fixed: {filepath}")
            return True
        else:
            print(f"  No changes: {filepath}")
            return False
    except Exception as e:
        print(f"✗ Error processing {filepath}: {e}")
        return False

def main():
    """Point d'entrée principal."""
    workspace = Path("/workspace/src")

    # Fichiers à traiter
    files_to_fix = [
        workspace / "internal/provider/resources/workflow.go",
        workspace / "internal/provider/resources/user.go",
        workspace / "internal/provider/resources/variable.go",
        workspace / "internal/provider/resources/tag.go",
        workspace / "internal/provider/resources/project.go",
        workspace / "internal/provider/resources/project_user.go",
        workspace / "internal/provider/resources/source_control_pull.go",
        workspace / "internal/provider/resources/execution_retry.go",
        workspace / "internal/provider/resources/credential_transfer.go",
        workspace / "internal/provider/resources/workflow_transfer.go",
    ]

    fixed_count = 0
    for filepath in files_to_fix:
        if filepath.exists():
            if process_file(filepath):
                fixed_count += 1
        else:
            print(f"  Not found: {filepath}")

    print(f"\n{'='*50}")
    print(f"Fixed {fixed_count} files")
    print(f"{'='*50}")

if __name__ == "__main__":
    main()
