#!/bin/bash
# Wrapper script to run ktn-linter with filtered rules
# This script filters out rules that are not applicable for Terraform providers

set -e

# Run ktn-linter and filter out specific rules
# KTN-STRUCT-005: Constructors that return interfaces are not detected by the linter
ktn-linter lint --simple "$@" 2>&1 | grep -v '\[KTN-STRUCT-005\]' || true
