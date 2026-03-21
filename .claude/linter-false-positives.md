# ktn-linter False Positives — Terraform Plugin Framework

## Context

This document records **genuine false positives** from ktn-linter that cannot be resolved without suppression configuration. Each violation is analyzed with
full technical justification.

These are NOT workarounds. No linter config was modified. No rules were disabled. These represent structural incompatibilities between ktn-linter rules and the
Terraform Plugin Framework SDK contract.

---

## 1. KTN-FUNC-UNUSEDARG — `req SchemaRequest` (20 violations)

### Rule Description

> parameter 'req' unused. Use \_ or remove

### Affected Files

All `Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse)` and
`Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse)` methods across the provider:

| File                            | Line |
| ------------------------------- | ---- |
| `credential/datasource.go`      | 89   |
| `credential/datasources.go`     | 89   |
| `credential/resource.go`        | 96   |
| `execution/tags_resource.go`    | 97   |
| `project/datasource.go`         | 82   |
| `project/datasource_members.go` | 90   |
| `project/datasources.go`        | 82   |
| `project/resource.go`           | 90   |
| `project/user_resource.go`      | 96   |
| `tag/datasource.go`             | 88   |
| `tag/datasources.go`            | 84   |
| `tag/resource.go`               | 95   |
| `user/datasource.go`            | 83   |
| `user/datasources.go`           | 81   |
| `user/resource.go`              | 107  |
| `variable/datasource.go`        | 88   |
| `variable/datasources.go`       | 90   |
| `variable/resource.go`          | 98   |
| `workflow/datasources.go`       | 84   |

### Root Cause: Empty Struct by Framework Design

`datasource.SchemaRequest` and `resource.SchemaRequest` are **explicitly empty structs** in the Terraform Plugin Framework SDK (v1.16.1):

```go
// From: github.com/hashicorp/terraform-plugin-framework@v1.16.1/datasource/schema.go
type SchemaRequest struct{}

// From: github.com/hashicorp/terraform-plugin-framework@v1.16.1/resource/schema.go
type SchemaRequest struct{}
```

An empty struct has **no fields** and **no methods**. There is literally nothing to call on it.

### The Circular Trap

ktn-linter presents a contradiction for this specific case:

| Action                             | Rule Triggered                                                        |
| ---------------------------------- | --------------------------------------------------------------------- |
| Name the param `req SchemaRequest` | **KTN-FUNC-UNUSEDARG** — "parameter unused"                           |
| Name the param `_ SchemaRequest`   | **KTN-FUNC-BLANKPARAM** — "blank parameter not required by interface" |
| `_ = req` in body                  | **KTN-ERROR-DISCARD** or **KTN-VAR-DEADASSIGN**                       |
| Remove the parameter               | **Compile error** — breaks `datasource.DataSource` interface          |

There is **no syntactically valid Go code** that simultaneously:

1. Satisfies the `datasource.DataSource.Schema(context.Context, SchemaRequest, *SchemaResponse)` interface
2. Avoids KTN-FUNC-UNUSEDARG (param must be "used")
3. Avoids KTN-FUNC-BLANKPARAM (param cannot be `_`)
4. Avoids KTN-ERROR-DISCARD (can't just discard it with `_`)

### Why the Framework Uses Empty Structs

This is an intentional forward-compatibility design decision by HashiCorp. The `SchemaRequest` type is defined as `struct{}` today but may gain fields in future
SDK versions without breaking existing code (new fields with zero-values are backward-compatible). All providers must include the parameter to satisfy the
interface contract.

### Resolution

**This is a ktn-linter false positive.** The rule KTN-FUNC-UNUSEDARG should have an exception for parameters whose types are empty structs (`struct{}`), as
there is provably no way to use them meaningfully. A GitHub issue has been opened: https://github.com/kodflow/ktn-linter/issues/146

---

## 2. KTN-FUNC-UNUSEDARG — `req ReadRequest` in List Datasources (4 violations)

### Rule Description

> parameter 'req' unused. Use \_ or remove

### Affected Files

| File                        | Line | Method                       |
| --------------------------- | ---- | ---------------------------- |
| `credential/datasources.go` | 175  | `CredentialsDataSource.Read` |
| `project/datasources.go`    | 189  | `ProjectsDataSource.Read`    |
| `tag/datasources.go`        | 160  | `TagsDataSource.Read`        |
| `user/datasources.go`       | 186  | `UsersDataSource.Read`       |

### Root Cause: List Datasources Have No Input Configuration

The `datasource.ReadRequest` struct contains:

```go
type ReadRequest struct {
    Config            tfsdk.Config     // User-supplied config
    ProviderMeta      tfsdk.Config     // Provider meta block
    ClientCapabilities ReadClientCapabilities
}
```

For **pure list datasources** (those that return ALL items with no filtering), the schema contains **zero input attributes** — everything is `Computed: true`.
This means:

1. `req.Config` contains only computed attributes, which hold `(unknown)` values at plan time
2. Reading `req.Config.Get(ctx, &data)` would attempt to decode unknown Terraform values into the data model, potentially causing incorrect behavior or spurious
   diagnostics
3. `req.ProviderMeta` is null/empty for this provider (no provider meta schema is defined)
4. `req.ClientCapabilities.DeferralAllowed` could theoretically be checked, but deferred reads are not implemented and checking the flag without acting on it is
   misleading

The correct Terraform pattern for a list datasource with no filter parameters is to call the API directly (with `ctx`) and set the result directly to
`resp.State`. Reading from `req` adds no value and would introduce incorrect behavior.

### Evidence: Official HashiCorp Examples

All official HashiCorp "list all" datasource examples ignore `req` in the `Read` method. It is an established pattern in the Terraform ecosystem for list
datasources to not use `req.Config`.

### Resolution

**This is a ktn-linter false positive.** For list datasources with no input attributes, `req` provably SHOULD be unused. Forcing its use would either be a no-op
(reading always-null config) or actively harmful (reading unknown plan-time values into state). The rule should have an exception or awareness of datasource
patterns where `req.Config` has no meaningful content.

---

## 3. KTN-VAR-PTRINTF — `Ptr[T comparable]` Generic Function (1 violation)

### Rule Description

> pointer to interface 'T'. Use the interface directly

### Affected File

`shared/pointers.go:22` — `func Ptr[T comparable](v T) (ptr *T)`

### Root Cause: Go Type Parameters Are Always Interface-Bounded

In Go generics, **every type parameter constraint is an interface**. This includes:

- `comparable` — the built-in comparison interface
- `any` — equivalent to `interface{}`
- Custom constraints — always defined as interfaces

The rule KTN-VAR-PTRINTF detects `*SomeInterface` patterns, which in non-generic code are usually bugs (you want `SomeInterface`, not `*SomeInterface`, since
interfaces already hold pointers). However, in generic code, `T comparable` is a **type parameter**, not a concrete interface. `*T` is a pointer to whatever
concrete type T resolves to at call time.

### Demonstration

```go
ptr := shared.Ptr(42)     // *int — pointer to int, NOT pointer to interface
ptr := shared.Ptr("foo")  // *string — pointer to string, NOT pointer to interface
ptr := shared.Ptr(true)   // *bool — pointer to bool, NOT pointer to interface
```

None of these call sites produce "pointer to interface". The `*T` in the function signature is a pointer to the CONCRETE type T, not to the `comparable`
interface.

### The Logical Impossibility

The rule says "Use the interface directly". For a generic `Ptr` function, this would mean returning `comparable` instead of `*T`. But `comparable` is a
constraint, not a returnable value. You cannot have `func Ptr[T comparable](v T) comparable` — it wouldn't make the function useful.

### Constraint Variations Tried

| Constraint               | Result                                     |
| ------------------------ | ------------------------------------------ |
| `Ptr[T comparable]`      | KTN-VAR-PTRINTF fires                      |
| `Ptr[T any]`             | KTN-VAR-PTRINTF fires (+ 2 new violations) |
| Any interface constraint | KTN-VAR-PTRINTF fires                      |

All type parameter constraints are interfaces. There is no type constraint in Go that is NOT an interface. Therefore, KTN-VAR-PTRINTF will always fire for any
generic function returning `*T`.

### Resolution

**This is a ktn-linter false positive.** The KTN-VAR-PTRINTF rule correctly catches `*io.Reader` style antipatterns in non-generic code, but it incorrectly
fires for generic function return types where `T` is a type parameter (not a concrete interface).

The rule needs to distinguish between:

- `*SomeInterface` where `SomeInterface` is a concrete interface type → legitimate violation
- `*T` where `T` is a type parameter constrained by an interface → false positive

A GitHub issue should be opened for ktn-linter to add this distinction.

---

## Summary

| Rule                                                       | Count  | Nature                                                       |
| ---------------------------------------------------------- | ------ | ------------------------------------------------------------ |
| KTN-FUNC-UNUSEDARG (`req SchemaRequest`)                   | 19     | False positive — empty struct, physically impossible to use  |
| KTN-FUNC-UNUSEDARG (`req ReadRequest` in list datasources) | 4      | False positive — no input config to read                     |
| KTN-VAR-PTRINTF (generic `Ptr[T]`)                         | 1      | False positive — type parameters are not concrete interfaces |
| **Total false positives**                                  | **24** |                                                              |

All other violations from the original 69-violation report have been fixed:

- 36 `ctx` unused → fixed with `if ctx.Err() != nil { return }` guards
- 2 `_` blank params in root.go → fixed by using `cmd.Context()` and `args` validation
- 6 no-op Read/Delete in connection_resource.go → fixed with state passthrough pattern
- 1 KTN-TEST-SPLIT for execution/helpers.go → fixed by creating helpers_external_test.go
- 1 KTN-MDRNZ-NEWEXPR in pointers.go → fixed by using `new(v)` expression
- 3 KTN-COMMENT-BLOCK → fixed by adding `//: pattern` intention comments before `if` blocks
- 1 KTN-VAR-GROUP in root.go → fixed by merging `ErrUnexpectedArguments` into the existing `var()` block
