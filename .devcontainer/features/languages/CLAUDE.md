<!-- updated: 2026-03-17T12:15:00Z -->

# Language Features

## Purpose

Language-specific installation scripts. Conventions handled by specialist agents.

## Available Languages

| Language      | Version          | Agent                           |
| ------------- | ---------------- | ------------------------------- |
| Python        | >= 3.14.0        | `developer-specialist-python`   |
| C             | C23              | `developer-specialist-c`        |
| C++           | C++23            | `developer-specialist-cpp`      |
| Java          | >= 25            | `developer-specialist-java`     |
| C#            | .NET 9+ / C# 13+ | `developer-specialist-csharp`   |
| Node.js       | >= 25.0.0        | `developer-specialist-nodejs`   |
| VB.NET        | .NET 9+          | `developer-specialist-vbnet`    |
| R             | >= 4.5.0         | `developer-specialist-r`        |
| Pascal        | FPC 3.2+         | `developer-specialist-pascal`   |
| Perl          | >= 5.40.0        | `developer-specialist-perl`     |
| Fortran       | Fortran 2023     | `developer-specialist-fortran`  |
| PHP           | >= 8.5.0         | `developer-specialist-php`      |
| Rust          | >= 1.92.0        | `developer-specialist-rust`     |
| Go            | >= 1.26.0        | `developer-specialist-go`       |
| Ada           | Ada 2022         | `developer-specialist-ada`      |
| MATLAB/Octave | Octave 9+        | `developer-specialist-matlab`   |
| Assembly      | NASM 2.16+       | `developer-specialist-assembly` |
| Kotlin        | >= 2.2.0         | `developer-specialist-kotlin`   |
| Swift         | >= 6.2.0         | `developer-specialist-swift`    |
| COBOL         | GnuCOBOL 3.2+    | `developer-specialist-cobol`    |
| Ruby          | >= 4.0.0         | `developer-specialist-ruby`     |
| Dart/Flutter  | >= 3.10/3.38     | `developer-specialist-dart`     |
| Lua           | >= 5.4.0         | `developer-specialist-lua`      |
| Scala         | >= 3.7.0         | `developer-specialist-scala`    |
| Elixir        | >= 1.19.0        | `developer-specialist-elixir`   |

## VS Code Extensions (LSP)

Each `devcontainer-feature.json` declares `customizations.vscode.extensions` with the LSP extension(s) for that language. When a language feature is enabled in
`devcontainer.json`, its extensions install automatically. Commented-out features contribute nothing — no bloat.

| Language      | Extension(s)                                              |
| ------------- | --------------------------------------------------------- |
| Go            | `golang.go`                                               |
| Python        | `ms-python.python`, `ms-python.vscode-pylance`            |
| Node.js       | `dbaeumer.vscode-eslint`                                  |
| Rust          | `rust-lang.rust-analyzer`                                 |
| Java          | `redhat.java`                                             |
| C#            | `ms-dotnettools.csharp`                                   |
| C++           | `llvm-vs-code-extensions.vscode-clangd`                   |
| C             | `llvm-vs-code-extensions.vscode-clangd`                   |
| Ruby          | `Shopify.ruby-lsp`                                        |
| PHP           | `bmewburn.vscode-intelephense-client`                     |
| Elixir        | `JakeBecker.elixir-ls`                                    |
| Scala         | `scalameta.metals`                                        |
| Dart/Flutter  | `Dart-Code.dart-code`, `Dart-Code.flutter`                |
| Kotlin        | `fwcd.kotlin`                                             |
| Swift         | `swiftlang.swift-vscode`                                  |
| R             | `REditorSupport.r`                                        |
| Perl          | `bscan.perlnavigator`                                     |
| Lua           | `sumneko.lua`                                             |
| Fortran       | `fortran-lang.linter-gfortran`                            |
| Ada           | `AdaCore.ada`                                             |
| COBOL         | `broadcomMFD.cobol-language-support`                      |
| Pascal        | `alefragnani.pascal`                                      |
| VB.NET        | `ms-dotnettools.csharp`                                   |
| MATLAB/Octave | `MathWorks.language-matlab`, `toasty-technologies.octave` |
| Assembly      | `13xforever.language-x86-64-assembly`                     |

## Per-Language Structure

```text
languages/
├── shared/             # Shared utility library
│   └── feature-utils.sh  # Colors, logging, arch detection, GitHub API
└── <language>/
    └── install.sh      # Installation script (sources shared/feature-utils.sh)
```

## Version Discovery

Agents use WebFetch on official sources to get latest versions dynamically. No static version files needed.

## Conventions

- All code in /src regardless of language
- Tests in /tests (except Go: alongside code)
- Specialist agents enforce academic standards

## Learned Conventions

- **set -e in all scripts**: Every install.sh uses `set -e` (25 occurrences)
- **FEATURE_DIR pattern**: All scripts start with `FEATURE_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"` (25 occurrences)
- **shellcheck directive**: All scripts include `# shellcheck source=../shared/feature-utils.sh` (25 occurrences)
- **devcontainer-feature.json**: Every language directory has both `install.sh` and `devcontainer-feature.json` (25 occurrences)
