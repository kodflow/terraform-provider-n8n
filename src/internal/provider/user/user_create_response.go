// Copyright (c) 2024 Florent (Kodflow). All rights reserved.
// Licensed under the Sustainable Use License 1.0
// See LICENSE in the project root for license information.

// Package user implements user management resources and data sources.
package user

// userCreateResponse is array item from n8n user creation API.
type userCreateResponse struct {
	User  userCreateResponseUser `json:"user"`
	Error string                 `json:"error"`
}
