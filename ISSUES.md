# Codacy Issues - terraform-provider-n8n

**Total Issues**: 329

## Summary

### By Severity
- ❌ **Error**: 28 issues
- 🔴 **High**: 96 issues
- ⚠️ **Warning**: 92 issues
- ℹ️ **Info**: 113 issues

### By Category
- **BestPractice**: 104 issues
- **CodeStyle**: 97 issues
- **ErrorProne**: 64 issues
- **Security**: 50 issues
- **Complexity**: 8 issues
- **Documentation**: 3 issues
- **UnusedCode**: 2 issues
- **Performance**: 1 issues

### By Tool
- **Pylint**: 108 issues
- **markdownlint**: 68 issues
- **Semgrep**: 55 issues
- **ShellCheck**: 33 issues
- **Prospector**: 28 issues
- **Bandit**: 20 issues
- **Ruff**: 15 issues
- **Hadolint**: 1 issues
- **Lizard**: 1 issues

---

## Issues by Pattern


### ❌ Error Issues

#### Semgrep_python.lang.security.audit.dangerous-subprocess-use-audit.dangerous-subprocess-use-audit (7 occurrences)

**Category**: Security  
**Description**: Detected subprocess function without a static string

**Affected files**:

- [ ] `codegen/patch-openapi.py:14`
- [ ] `codegen/download-only.py:16`
- [ ] `codegen/generate-sdk.py:14`
- [ ] `codegen/generate-sdk.py:72`
- [ ] `codegen/patch-only.py:13`
- [ ] `codegen/patch-only.py:36`
- [ ] `codegen/update-n8n-version.py:14`

#### Semgrep_python.lang.security.audit.subprocess-shell-true.subprocess-shell-true (6 occurrences)

**Category**: Security  
**Description**: Found subprocess function with shell=True

**Affected files**:

- [ ] `codegen/patch-openapi.py:14`
- [ ] `codegen/download-only.py:16`
- [ ] `codegen/generate-sdk.py:14`
- [ ] `codegen/generate-sdk.py:72`
- [ ] `codegen/patch-only.py:13`
- [ ] `codegen/update-n8n-version.py:14`

#### shellcheck_SC1073 (1 occurrences)

**Category**: ErrorProne  
**Description**: Couldn't parse this explicit subshell

**Affected files**:

- [ ] `.devcontainer/p10k.sh:24`

#### Semgrep_python.lang.security.audit.dangerous-subprocess-use-tainted-env-args.dangerous-subprocess-use-tainted-env-args (1 occurrences)

**Category**: Security  
**Description**: Detected subprocess function with user controlled data

**Affected files**:

- [ ] `codegen/patch-openapi.py:14`

#### Lizard_ccn-critical (1 occurrences)

**Category**: Complexity  
**Description**: Method has a cyclomatic complexity of 16 (limit is 12)

**Affected files**:

- [ ] `codegen/download-only.py:44`


### 🔴 High Issues

#### Semgrep_yaml.semgrep.duplicate-id.duplicate-id (18 occurrences)

**Category**: ErrorProne  
**Description**: The 'id' field was used multiple times

**Affected files**:

- [ ] `.golangci.yml:68`
- [ ] `.golangci.yml:69`
- [ ] `.golangci.yml:70`
- [ ] `.golangci.yml:71`
- [ ] `.golangci.yml:72`
- [ ] `.golangci.yml:73`
- [ ] `.golangci.yml:74`
- [ ] `.golangci.yml:75`
- [ ] `.golangci.yml:76`
- [ ] `.golangci.yml:78`
- [ ] `.golangci.yml:79`
- [ ] `.golangci.yml:82`
- [ ] `.golangci.yml:84`
- [ ] `.golangci.yml:85`
- [ ] `.golangci.yml:87`
- [ ] `.golangci.yml:88`
- [ ] `.golangci.yml:91`
- [ ] `.golangci.yml:94`

#### Prospector_bandit (12 occurrences)

**Category**: Security  
**Description**: Various Prospector bandit security issues

**Affected files**:

- [ ] `Multiple Python files`

#### Bandit_B602 (11 occurrences)

**Category**: Security  
**Description**: subprocess call with shell=True identified, security issue

**Affected files**:

- [ ] `codegen/patch-openapi.py:14`
- [ ] `codegen/download-only.py:16`
- [ ] `codegen/generate-sdk.py:14`
- [ ] `codegen/generate-sdk.py:74`
- [ ] `codegen/generate-sdk.py:139`
- [ ] `codegen/patch-only.py:13`
- [ ] `codegen/patch-only.py:38`
- [ ] `codegen/update-n8n-version.py:14`

#### shellcheck_SC2046 (8 occurrences)

**Category**: ErrorProne  
**Description**: Quote this to prevent word splitting

**Affected files**:

- [ ] `scripts/generate-coverage.sh:55`
- [ ] `scripts/generate-coverage.sh:117`
- [ ] `scripts/generate-coverage.sh:119`
- [ ] `scripts/generate-coverage.sh:244`
- [ ] `scripts/generate-coverage.sh:246`
- [ ] `scripts/generate-coverage.sh:248`
- [ ] `scripts/generate-coverage.sh:408`

#### PyLintPython3_W1510 (7 occurrences)

**Category**: ErrorProne  
**Description**: subprocess.run used without explicitly defining the value for 'check'

**Affected files**:

- [ ] `codegen/generate-sdk.py:14`
- [ ] `codegen/generate-sdk.py:72`
- [ ] `codegen/generate-sdk.py:139`
- [ ] `codegen/download-only.py:16`
- [ ] `codegen/patch-openapi.py:14`
- [ ] `codegen/patch-only.py:13`
- [ ] `codegen/patch-only.py:36`
- [ ] `codegen/update-n8n-version.py:14`

#### Semgrep_python_exec_rule-subprocess-popen-shell-true (7 occurrences)

**Category**: Security  
**Description**: Found subprocess function with shell=True

**Affected files**:

- [ ] `codegen/patch-openapi.py:14`
- [ ] `codegen/download-only.py:16`
- [ ] `codegen/generate-sdk.py:14`
- [ ] `codegen/generate-sdk.py:72`
- [ ] `codegen/patch-only.py:13`
- [ ] `codegen/patch-only.py:36`
- [ ] `codegen/update-n8n-version.py:14`

#### shellcheck_SC2155 (6 occurrences)

**Category**: ErrorProne  
**Description**: Declare and assign separately to avoid masking return values

**Affected files**:

- [ ] `scripts/generate-coverage.sh:203`
- [ ] `scripts/generate-coverage.sh:225`
- [ ] `scripts/generate-coverage.sh:226`
- [ ] `scripts/generate-coverage.sh:231`
- [ ] `scripts/generate-coverage.sh:232`
- [ ] `scripts/generate-coverage.sh:243`

#### Semgrep_python_tmpdir_rule-hardcodedtmp (5 occurrences)

**Category**: Security  
**Description**: Application creating files in /tmp without using tempfile.TemporaryFile

**Affected files**:

- [ ] `codegen/download-only.py:56`
- [ ] `codegen/download-only.py:57`
- [ ] `codegen/download-only.py:67`
- [ ] `codegen/download-only.py:85`
- [ ] `codegen/generate-sdk.py:55`

#### Semgrep_bash.lang.correctness.unquoted-expansion.unquoted-variable-expansion-in-command (4 occurrences)

**Category**: ErrorProne  
**Description**: Variable expansions must be double-quoted

**Affected files**:

- [ ] `scripts/generate-coverage.sh:293`
- [ ] `scripts/generate-coverage.sh:342`
- [ ] `scripts/generate-coverage.sh:365`
- [ ] `scripts/generate-coverage.sh:379`

#### PyLintPython3_W0702 (2 occurrences)

**Category**: ErrorProne  
**Description**: No exception type(s) specified

**Affected files**:

- [ ] `codegen/download-only.py:31`
- [ ] `codegen/update-n8n-version.py:27`

#### PyLintPython3_W0718 (2 occurrences)

**Category**: ErrorProne  
**Description**: Catching too general exception Exception

**Affected files**:

- [ ] `codegen/patch-openapi.py:145`
- [ ] `codegen/update-n8n-version.py:52`

#### Ruff_E722_bare-except (2 occurrences)

**Category**: ErrorProne  
**Description**: Do not use bare `except`

**Affected files**:

- [ ] `codegen/download-only.py:31`
- [ ] `codegen/update-n8n-version.py:27`

#### Prospector_mypy (2 occurrences)

**Category**: ErrorProne  
**Description**: Library stubs not installed for yaml

**Affected files**:

- [ ] `codegen/download-only.py:95`
- [ ] `codegen/patch-openapi.py:141`

#### shellcheck_SC2124 (1 occurrences)

**Category**: ErrorProne  
**Description**: Assigning an array to a string

**Affected files**:

- [ ] `scripts/generate-coverage.sh:194`

#### Semgrep_go.lang.correctness.looppointer.exported_loop_pointer (1 occurrences)

**Category**: ErrorProne  
**Description**: Loop pointer that may be exported from the loop

**Affected files**:

- [ ] `src/internal/provider/credential/helpers.go:97`

#### Semgrep_go_memory_rule-memoryaliasing (1 occurrences)

**Category**: ErrorProne  
**Description**: Go's for...range statements create an iteration variable for each iteration

**Affected files**:

- [ ] `src/internal/provider/credential/helpers.go:97`

#### Hadolint_DL4006 (1 occurrences)

**Category**: ErrorProne  
**Description**: Set the SHELL option -o pipefail before RUN with a pipe in it

**Affected files**:

- [ ] `.devcontainer/Dockerfile:5`


### ⚠️ Warning Issues

#### Ruff_F541_f-string-missing-placeholders (11 occurrences)

**Category**: BestPractice  
**Description**: f-string without any placeholders

**Affected files**:

- [ ] `codegen/patch-openapi.py:83`
- [ ] `codegen/patch-openapi.py:114`
- [ ] `codegen/patch-openapi.py:127`
- [ ] `codegen/download-only.py:75`
- [ ] `codegen/download-only.py:76`
- [ ] `codegen/download-only.py:78`
- [ ] `codegen/download-only.py:158`
- [ ] `codegen/generate-sdk.py:79`
- [ ] `codegen/update-n8n-version.py:87`

#### PyLintPython3_W1514 (9 occurrences)

**Category**: BestPractice  
**Description**: Using open without explicitly specifying an encoding

**Affected files**:

- [ ] `codegen/patch-openapi.py:62`
- [ ] `codegen/patch-openapi.py:142`
- [ ] `codegen/download-only.py:28`
- [ ] `codegen/download-only.py:102`
- [ ] `codegen/download-only.py:130`
- [ ] `codegen/patch-only.py:56`
- [ ] `codegen/patch-only.py:83`
- [ ] `codegen/update-n8n-version.py:58`
- [ ] `codegen/update-n8n-version.py:84`

#### Semgrep_python.lang.best-practice.unspecified-open-encoding.unspecified-open-encoding (9 occurrences)

**Category**: BestPractice  
**Description**: Missing 'encoding' parameter in open()

**Affected files**:

- [ ] `codegen/patch-openapi.py:62`
- [ ] `codegen/patch-openapi.py:142`
- [ ] `codegen/download-only.py:28`
- [ ] `codegen/download-only.py:102`
- [ ] `codegen/download-only.py:130`
- [ ] `codegen/patch-only.py:56`
- [ ] `codegen/patch-only.py:83`
- [ ] `codegen/update-n8n-version.py:58`
- [ ] `codegen/update-n8n-version.py:84`

#### shellcheck_SC2001 (5 occurrences)

**Category**: Performance  
**Description**: See if you can use ${variable//search/replace} instead

**Affected files**:

- [ ] `scripts/generate-coverage.sh:29`
- [ ] `scripts/generate-coverage.sh:114`
- [ ] `scripts/generate-coverage.sh:243`
- [ ] `scripts/generate-coverage.sh:309`
- [ ] `scripts/generate-coverage.sh:336`

#### PyLintPython3_R0915 (4 occurrences)

**Category**: Complexity  
**Description**: Too many statements

**Affected files**:

- [ ] `codegen/download-only.py:44`
- [ ] `codegen/generate-sdk.py:21`
- [ ] `codegen/patch-openapi.py:70`
- [ ] `codegen/patch-only.py:20`

#### shellcheck_SC2086 (4 occurrences)

**Category**: BestPractice  
**Description**: Double quote to prevent globbing and word splitting

**Affected files**:

- [ ] `scripts/generate-coverage.sh:293`
- [ ] `scripts/generate-coverage.sh:342`
- [ ] `scripts/generate-coverage.sh:365`
- [ ] `scripts/generate-coverage.sh:379`

#### PyLintPython3_R0912 (2 occurrences)

**Category**: Complexity  
**Description**: Too many branches

**Affected files**:

- [ ] `codegen/generate-sdk.py:21`
- [ ] `codegen/patch-only.py:20`

#### Bandit_B108 (2 occurrences)

**Category**: Security  
**Description**: Probable insecure usage of temp file/directory

**Affected files**:

- [ ] `codegen/download-only.py:49`
- [ ] `codegen/generate-sdk.py:27`

#### shellcheck_SC2034 (2 occurrences)

**Category**: ErrorProne  
**Description**: Variable appears unused

**Affected files**:

- [ ] `scripts/generate-coverage.sh:8`
- [ ] `scripts/add-copyright-headers.sh:11`

#### PyLintPython3_R0914 (1 occurrences)

**Category**: Complexity  
**Description**: Too many local variables

**Affected files**:

- [ ] `codegen/download-only.py:44`

#### PyLintPython3_R0801 (1 occurrences)

**Category**: Complexity  
**Description**: Similar lines in 2 files

**Affected files**:

- [ ] `codegen/generate-sdk.py:1`

#### Bandit_B607 (1 occurrences)

**Category**: Security  
**Description**: Starting a process with a partial executable path

**Affected files**:

- [ ] `codegen/generate-sdk.py:139`

#### Bandit_B110 (1 occurrences)

**Category**: Security  
**Description**: Try, Except, Pass detected

**Affected files**:

- [ ] `codegen/download-only.py:31`

#### Prospector_vulture (1 occurrences)

**Category**: UnusedCode  
**Description**: Unused variable

**Affected files**:

- [ ] `codegen/download-only.py:23`

#### shellcheck_SC1091 (1 occurrences)

**Category**: ErrorProne  
**Description**: Not following file

**Affected files**:

- [ ] `.devcontainer/post-create.sh:24`

#### shellcheck_SC1072 (1 occurrences)

**Category**: ErrorProne  
**Description**: Unexpected keyword/token

**Affected files**:

- [ ] `.devcontainer/p10k.sh:24`

#### shellcheck_SC2010 (1 occurrences)

**Category**: BestPractice  
**Description**: Don't use ls | grep. Use a glob or a for loop with a condition

**Affected files**:

- [ ] `scripts/generate-coverage.sh:305`

#### shellcheck_SC2064 (1 occurrences)

**Category**: BestPractice  
**Description**: Use single quotes, otherwise this expands now rather than when signalled

**Affected files**:

- [ ] `scripts/generate-coverage.sh:147`

#### shellcheck_SC2016 (1 occurrences)

**Category**: BestPractice  
**Description**: Expressions don't expand in single quotes, use double quotes for that

**Affected files**:

- [ ] `.devcontainer/post-create.sh:108`


### ℹ️ Info Issues

#### markdownlint_MD043 (28 occurrences)

**Category**: BestPractice  
**Description**: Required heading structure

**Affected files**:

- [ ] `docs/resources/workflow.md:9`
- [ ] `docs/resources/credential.md:12`
- [ ] `docs/resources/execution_retry.md:9`
- [ ] `docs/data-sources/user.md:9`
- [ ] `README.md:1`
- [ ] `docs/data-sources/workflow.md:9`
- [ ] `docs/data-sources/variables.md:9`
- [ ] `docs/resources/user.md:9`
- [ ] `docs/resources/credential_transfer.md:9`
- [ ] `docs/data-sources/tags.md:9`
- [ ] `docs/resources/project.md:9`
- [ ] `docs/index.md:8`
- [ ] `LICENSE.md:1`
- [ ] `docs/resources/source_control_pull.md:9`
- [ ] `docs/resources/variable.md:9`
- [ ] `docs/data-sources/workflows.md:9`
- [ ] `docs/resources/workflow_transfer.md:9`
- [ ] `.github/hooks/README.md:1`
- [ ] `docs/data-sources/variable.md:9`
- [ ] `docs/data-sources/execution.md:9`
- [ ] `docs/resources/tag.md:9`
- [ ] `CLAUDE.md:1`
- [ ] `docs/resources/project_user.md:9`
- [ ] `docs/data-sources/project.md:9`
- [ ] `docs/data-sources/projects.md:9`
- [ ] `docs/data-sources/tag.md:9`
- [ ] `.github/CONTRIBUTING.md:5`
- [ ] `docs/data-sources/users.md:9`
- [ ] `.github/SECURITY.md:5`
- [ ] `COVERAGE.MD:1`
- [ ] `docs/data-sources/executions.md:9`
- [ ] `.github/PULL_REQUEST_TEMPLATE.md:1`

#### PyLintPython3_C0301 (15 occurrences)

**Category**: CodeStyle  
**Description**: Line too long

**Affected files**:

- [ ] `codegen/patch-openapi.py:66`
- [ ] `codegen/patch-openapi.py:67`
- [ ] `codegen/patch-openapi.py:115`
- [ ] `codegen/patch-openapi.py:116`
- [ ] `codegen/download-only.py:48`
- [ ] `codegen/download-only.py:59`
- [ ] `codegen/download-only.py:74`
- [ ] `codegen/download-only.py:91`
- [ ] `codegen/download-only.py:127`
- [ ] `codegen/download-only.py:128`
- [ ] `codegen/download-only.py:151`
- [ ] `codegen/download-only.py:154`
- [ ] `codegen/download-only.py:156`
- [ ] `codegen/generate-sdk.py:57`
- [ ] `codegen/generate-sdk.py:73`
- [ ] `codegen/generate-sdk.py:106`
- [ ] `codegen/generate-sdk.py:110`
- [ ] `codegen/update-n8n-version.py:42`
- [ ] `codegen/update-n8n-version.py:62`
- [ ] `codegen/update-n8n-version.py:79`

#### PyLintPython3_C0103 (10 occurrences)

**Category**: CodeStyle  
**Description**: Variable name doesn't conform to snake_case naming style

**Affected files**:

- [ ] `codegen/generate-sdk.py:1`
- [ ] `codegen/generate-sdk.py:25`
- [ ] `codegen/generate-sdk.py:26`
- [ ] `codegen/generate-sdk.py:27`
- [ ] `codegen/generate-sdk.py:28`
- [ ] `codegen/download-only.py:1`
- [ ] `codegen/download-only.py:48`
- [ ] `codegen/download-only.py:49`
- [ ] `codegen/download-only.py:50`
- [ ] `codegen/patch-openapi.py:1`
- [ ] `codegen/patch-only.py:1`
- [ ] `codegen/patch-only.py:23`
- [ ] `codegen/update-n8n-version.py:1`

#### PyLintPython3_W1309 (9 occurrences)

**Category**: CodeStyle  
**Description**: Using an f-string that does not have any interpolated variables

**Affected files**:

- [ ] `codegen/patch-openapi.py:83`
- [ ] `codegen/patch-openapi.py:114`
- [ ] `codegen/patch-openapi.py:127`
- [ ] `codegen/download-only.py:75`
- [ ] `codegen/download-only.py:76`
- [ ] `codegen/download-only.py:78`
- [ ] `codegen/download-only.py:158`
- [ ] `codegen/generate-sdk.py:79`
- [ ] `codegen/update-n8n-version.py:87`

#### markdownlint_MD033 (6 occurrences)

**Category**: BestPractice  
**Description**: Inline HTML

**Affected files**:

- [ ] `docs/data-sources/tags.md:21`
- [ ] `docs/data-sources/projects.md:21`
- [ ] `docs/data-sources/variables.md:26`
- [ ] `docs/data-sources/users.md:21`
- [ ] `docs/data-sources/executions.md:28`
- [ ] `docs/data-sources/workflows.md:25`

#### PyLintPython3_C1805 (5 occurrences)

**Category**: CodeStyle  
**Description**: Simplify comparison to falsey value

**Affected files**:

- [ ] `codegen/generate-sdk.py:15`
- [ ] `codegen/download-only.py:17`
- [ ] `codegen/patch-openapi.py:31`
- [ ] `codegen/patch-only.py:14`
- [ ] `codegen/patch-only.py:42`
- [ ] `codegen/update-n8n-version.py:15`

#### Bandit_B404 (5 occurrences)

**Category**: Security  
**Description**: Consider possible security implications associated with the subprocess module

**Affected files**:

- [ ] `codegen/patch-openapi.py:9`
- [ ] `codegen/download-only.py:8`
- [ ] `codegen/generate-sdk.py:7`
- [ ] `codegen/patch-only.py:7`
- [ ] `codegen/update-n8n-version.py:6`

#### markdownlint_MD040 (4 occurrences)

**Category**: CodeStyle  
**Description**: Fenced code blocks should have a language specified

**Affected files**:

- [ ] `.github/CONTRIBUTING.md:126`
- [ ] `CLAUDE.md:60`
- [ ] `README.md:162`

#### PyLintPython3_C0415 (4 occurrences)

**Category**: CodeStyle  
**Description**: Import outside toplevel

**Affected files**:

- [ ] `codegen/patch-openapi.py:141`
- [ ] `codegen/download-only.py:95`
- [ ] `codegen/download-only.py:96`

#### Ruff_E741_ambiguous-variable-name (4 occurrences)

**Category**: CodeStyle  
**Description**: Ambiguous variable name: `l`

**Affected files**:

- [ ] `codegen/patch-openapi.py:66`
- [ ] `codegen/patch-openapi.py:67`
- [ ] `codegen/patch-openapi.py:115`
- [ ] `codegen/patch-openapi.py:116`

#### Prospector_pycodestyle (4 occurrences)

**Category**: CodeStyle  
**Description**: Expected 2 blank lines after class or function definition

**Affected files**:

- [ ] `codegen/patch-only.py:93`
- [ ] `codegen/patch-openapi.py:158`
- [ ] `codegen/download-only.py:167`
- [ ] `codegen/update-n8n-version.py:98`
- [ ] `codegen/generate-sdk.py:155`

#### PyLintPython3_C0116 (3 occurrences)

**Category**: Documentation  
**Description**: Missing function or method docstring

**Affected files**:

- [ ] `codegen/download-only.py:44`
- [ ] `codegen/generate-sdk.py:21`
- [ ] `codegen/patch-only.py:20`
- [ ] `codegen/update-n8n-version.py:30`

#### shellcheck_SC2129 (3 occurrences)

**Category**: Performance  
**Description**: Consider using { cmd1; cmd2; } >> file instead of individual redirects

**Affected files**:

- [ ] `scripts/add-copyright-headers.sh:94`
- [ ] `scripts/generate-coverage.sh:279`
- [ ] `scripts/generate-coverage.sh:297`

#### markdownlint_MD032 (2 occurrences)

**Category**: CodeStyle  
**Description**: Lists should be surrounded by blank lines

**Affected files**:

- [ ] `COVERAGE.MD:34`
- [ ] `COVERAGE.MD:6`

#### markdownlint_MD034 (2 occurrences)

**Category**: BestPractice  
**Description**: Bare URL used

**Affected files**:

- [ ] `docs/index.md:19`
- [ ] `LICENSE.md:63`

#### markdownlint_MD041 (1 occurrences)

**Category**: BestPractice  
**Description**: First line in a file should be a top-level heading

**Affected files**:

- [ ] `.github/PULL_REQUEST_TEMPLATE.md:1`

#### markdownlint_MD037 (1 occurrences)

**Category**: CodeStyle  
**Description**: Spaces inside emphasis markers

**Affected files**:

- [ ] `CLAUDE.md:23`

#### markdownlint_MD012 (1 occurrences)

**Category**: CodeStyle  
**Description**: Multiple consecutive blank lines

**Affected files**:

- [ ] `COVERAGE.MD:159`

#### markdownlint_MD036 (1 occurrences)

**Category**: BestPractice  
**Description**: Emphasis used instead of a heading

**Affected files**:

- [ ] `CLAUDE.md:48`

#### PyLintPython3_W0613 (1 occurrences)

**Category**: UnusedCode  
**Description**: Unused argument

**Affected files**:

- [ ] `codegen/download-only.py:23`

#### PyLintPython3_R1714 (1 occurrences)

**Category**: Performance  
**Description**: Consider merging these comparisons with 'in'

**Affected files**:

- [ ] `codegen/download-only.py:163`

#### PyLintPython3_W0404 (1 occurrences)

**Category**: CodeStyle  
**Description**: Reimport

**Affected files**:

- [ ] `codegen/download-only.py:140`

#### shellcheck_SC2126 (1 occurrences)

**Category**: CodeStyle  
**Description**: Consider using 'grep -c' instead of 'grep|wc -l'

**Affected files**:

- [ ] `scripts/generate-coverage.sh:232`

#### Semgrep_bash.lang.best-practice.useless-cat.useless-cat (1 occurrences)

**Category**: BestPractice  
**Description**: Useless call to 'cat' in a pipeline

**Affected files**:

- [ ] `scripts/generate-coverage.sh:203`


---

*Generated from Codacy analysis*
