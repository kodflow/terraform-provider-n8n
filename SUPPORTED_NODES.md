# N8N Terraform Provider - Supported Nodes

**Generated**: 2025-11-17T15:23:04.863Z
**Provider Version**: Latest
**N8N Version**: unknown
**Last Sync**: 2025-11-17T14:22:04.436Z

## Overview

This document lists all **296 n8n nodes** currently supported by the Terraform provider.

### Statistics

- **Total Nodes**: 296
- **Core Nodes**: 5
- **Trigger Nodes**: 25
- **Integration Nodes**: 266
- **Require Credentials**: 252
- **No Credentials**: 44

### Testing Status

All 296 nodes have been tested with `terraform init` and `terraform validate`:
- ✅ **296/296 workflows passed** (100% success rate)
- Each node has a complete example workflow in `examples/nodes/{category}/{node-slug}/`
- Full test results available in `WORKFLOWS_TEST_RESULTS.md`

---

## Quick Navigation

- [Core Nodes](#core-nodes) (5)
- [Trigger Nodes](#trigger-nodes) (25)
- [Integration Nodes](#integration-nodes) (266)
- [Credential Requirements](#credential-requirements)
- [Usage Examples](#usage-examples)

---

## Core Nodes

Essential workflow building blocks for data manipulation, flow control, and logic.

**Total**: 5 nodes

| Node | Type | Description | Credentials | Example |
|------|------|-------------|-------------|----------|
| **Code** | `code` | Run custom JavaScript or Python code | ✅ None | [`examples/nodes/core/code/`](examples/nodes/core/code/) |
| **If** | `if` | Route items to different branches (true/false) | ✅ None | [`examples/nodes/core/if/`](examples/nodes/core/if/) |
| **Merge** | `merge` | Merges data of multiple streams once data from both is available | ✅ None | [`examples/nodes/core/merge/`](examples/nodes/core/merge/) |
| **Set** | `set` | Add or edit fields on an input item and optionally remove other fields | ✅ None | [`examples/nodes/core/set/`](examples/nodes/core/set/) |
| **Switch** | `switch` | Route items depending on defined expression or rules | ✅ None | [`examples/nodes/core/switch/`](examples/nodes/core/switch/) |

## Trigger Nodes

Event-based nodes that initiate workflow execution.

**Total**: 25 nodes

| Node | Type | Description | Credentials | Example |
|------|------|-------------|-------------|----------|
| **Acuity Scheduling Trigger** | `acuitySchedulingTrigger` | Handle Acuity Scheduling events via webhooks | ✅ None | [`examples/nodes/trigger/acuity-scheduling-trigger/`](examples/nodes/trigger/acuity-scheduling-trigger/) |
| **Bitbucket Trigger** | `bitbucketTrigger` | Handle Bitbucket events via webhooks | ✅ None | [`examples/nodes/trigger/bitbucket-trigger/`](examples/nodes/trigger/bitbucket-trigger/) |
| **Cal.com Trigger** | `calTrigger` | Handle Cal.com events via webhooks | ✅ None | [`examples/nodes/trigger/cal-com-trigger/`](examples/nodes/trigger/cal-com-trigger/) |
| **Calendly Trigger** | `calendlyTrigger` | Starts the workflow when Calendly events occur | ✅ None | [`examples/nodes/trigger/calendly-trigger/`](examples/nodes/trigger/calendly-trigger/) |
| **Email Trigger (IMAP)** | `emailReadImap` | Triggers the workflow when a new email is received | ✅ None | [`examples/nodes/trigger/email-trigger-imap/`](examples/nodes/trigger/email-trigger-imap/) |
| **Error Trigger** | `errorTrigger` | Triggers the workflow when another workflow has an error | ✅ None | [`examples/nodes/trigger/error-trigger/`](examples/nodes/trigger/error-trigger/) |
| **Eventbrite Trigger** | `eventbriteTrigger` | Handle Eventbrite events via webhooks | ✅ None | [`examples/nodes/trigger/eventbrite-trigger/`](examples/nodes/trigger/eventbrite-trigger/) |
| **Facebook Lead Ads Trigger** | `facebookLeadAdsTrigger` | Handle Facebook Lead Ads events via webhooks | ⚠️ Authentication Required | [`examples/nodes/trigger/facebook-lead-ads-trigger/`](examples/nodes/trigger/facebook-lead-ads-trigger/) |
| **Figma Trigger (Beta)** | `figmaTrigger` | Starts the workflow when Figma events occur | ✅ None | [`examples/nodes/trigger/figma-trigger-beta/`](examples/nodes/trigger/figma-trigger-beta/) |
| **Form.io Trigger** | `formIoTrigger` | Handle form.io events via webhooks | ✅ None | [`examples/nodes/trigger/form-io-trigger/`](examples/nodes/trigger/form-io-trigger/) |
| **Formstack Trigger** | `formstackTrigger` | Starts the workflow on a Formstack form submission. | ✅ None | [`examples/nodes/trigger/formstack-trigger/`](examples/nodes/trigger/formstack-trigger/) |
| **Gumroad Trigger** | `gumroadTrigger` | Handle Gumroad events via webhooks | ✅ None | [`examples/nodes/trigger/gumroad-trigger/`](examples/nodes/trigger/gumroad-trigger/) |
| **Jotform Trigger** | `jotFormTrigger` | Handle Jotform events via webhooks | ✅ None | [`examples/nodes/trigger/jotform-trigger/`](examples/nodes/trigger/jotform-trigger/) |
| **Local File Trigger** | `localFileTrigger` | Triggers a workflow on file system changes | ✅ None | [`examples/nodes/trigger/local-file-trigger/`](examples/nodes/trigger/local-file-trigger/) |
| **Manual Trigger** | `manualTrigger` | Runs the flow on clicking a button in n8n | ✅ None | [`examples/nodes/trigger/manual-trigger/`](examples/nodes/trigger/manual-trigger/) |
| **n8n Trigger** | `n8nTrigger` | Handle events and perform actions on your n8n instance | ✅ None | [`examples/nodes/trigger/n8n-trigger/`](examples/nodes/trigger/n8n-trigger/) |
| **Postmark Trigger** | `postmarkTrigger` | Starts the workflow when Postmark events occur | ✅ None | [`examples/nodes/trigger/postmark-trigger/`](examples/nodes/trigger/postmark-trigger/) |
| **SSE Trigger** | `sseTrigger` | Triggers the workflow when Server-Sent Events occur | ✅ None | [`examples/nodes/trigger/sse-trigger/`](examples/nodes/trigger/sse-trigger/) |
| **SurveyMonkey Trigger** | `surveyMonkeyTrigger` | Starts the workflow when Survey Monkey events occur | ⚠️ Authentication Required | [`examples/nodes/trigger/surveymonkey-trigger/`](examples/nodes/trigger/surveymonkey-trigger/) |
| **Toggl Trigger** | `togglTrigger` | Starts the workflow when Toggl events occur | ✅ None | [`examples/nodes/trigger/toggl-trigger/`](examples/nodes/trigger/toggl-trigger/) |
| **Typeform Trigger** | `typeformTrigger` | Starts the workflow on a Typeform form submission | ✅ None | [`examples/nodes/trigger/typeform-trigger/`](examples/nodes/trigger/typeform-trigger/) |
| **Webhook** | `webhook` | Starts the workflow when a webhook is called | ✅ None | [`examples/nodes/trigger/webhook/`](examples/nodes/trigger/webhook/) |
| **Workable Trigger** | `workableTrigger` | Starts the workflow when Workable events occur | ✅ None | [`examples/nodes/trigger/workable-trigger/`](examples/nodes/trigger/workable-trigger/) |
| **Workflow Trigger** | `workflowTrigger` | Triggers based on various lifecycle events, like when a workflow is activated | ✅ None | [`examples/nodes/trigger/workflow-trigger/`](examples/nodes/trigger/workflow-trigger/) |
| **Wufoo Trigger** | `wufooTrigger` | Handle Wufoo events via webhooks | ✅ None | [`examples/nodes/trigger/wufoo-trigger/`](examples/nodes/trigger/wufoo-trigger/) |

## Integration Nodes

Third-party service integrations for connecting to external platforms.

**Total**: 266 nodes

| Node | Type | Description | Credentials | Example |
|------|------|-------------|-------------|----------|
| **Action Network** | `actionNetwork` | Consume the Action Network API | ⚠️ Authentication Required | [`examples/nodes/integration/action-network/`](examples/nodes/integration/action-network/) |
| **ActiveCampaign** | `activeCampaign` | Create and edit data in ActiveCampaign | ⚠️ Authentication Required | [`examples/nodes/integration/activecampaign/`](examples/nodes/integration/activecampaign/) |
| **Adalo** | `adalo` | Consume Adalo API | ⚠️ Authentication Required | [`examples/nodes/integration/adalo/`](examples/nodes/integration/adalo/) |
| **Affinity** | `affinity` | Consume Affinity API | ⚠️ Authentication Required | [`examples/nodes/integration/affinity/`](examples/nodes/integration/affinity/) |
| **Agile CRM** | `agileCrm` | Consume Agile CRM API | ⚠️ Authentication Required | [`examples/nodes/integration/agile-crm/`](examples/nodes/integration/agile-crm/) |
| **AI Transform** | `aiTransform` | Modify data based on instructions written in plain english | ⚠️ Authentication Required | [`examples/nodes/integration/ai-transform/`](examples/nodes/integration/ai-transform/) |
| **Airtable** | `airtable` | Read, update, write and delete data from Airtable | ⚠️ Authentication Required | [`examples/nodes/integration/airtable/`](examples/nodes/integration/airtable/) |
| **Airtop** | `airtop` | Scrape and control any site with Airtop | ⚠️ Authentication Required | [`examples/nodes/integration/airtop/`](examples/nodes/integration/airtop/) |
| **AMQP Sender** | `amqp` | Sends a raw-message via AMQP 1.0, executed once per item | ⚠️ Authentication Required | [`examples/nodes/integration/amqp-sender/`](examples/nodes/integration/amqp-sender/) |
| **APITemplate.io** | `apiTemplateIo` | Consume the APITemplate.io API | ⚠️ API Key | [`examples/nodes/integration/apitemplate-io/`](examples/nodes/integration/apitemplate-io/) |
| **Asana** | `asana` | Consume Asana REST API | ⚠️ Authentication Required | [`examples/nodes/integration/asana/`](examples/nodes/integration/asana/) |
| **Automizy** | `automizy` | Consume Automizy API | ⚠️ Authentication Required | [`examples/nodes/integration/automizy/`](examples/nodes/integration/automizy/) |
| **Autopilot** | `autopilot` | Consume Autopilot API | ⚠️ Authentication Required | [`examples/nodes/integration/autopilot/`](examples/nodes/integration/autopilot/) |
| **AWS Lambda** | `awsLambda` | Invoke functions on AWS Lambda | ⚠️ Authentication Required | [`examples/nodes/integration/aws-lambda/`](examples/nodes/integration/aws-lambda/) |
| **Background Color** | `Blur` | Adds a blur to the image and so makes it less sharp | ⚠️ Authentication Required | [`examples/nodes/integration/background-color/`](examples/nodes/integration/background-color/) |
| **BambooHr** | `n8n-nodes-base.bamboohr` | N/A | ⚠️ Authentication Required | [`examples/nodes/integration/bamboohr/`](examples/nodes/integration/bamboohr/) |
| **Bannerbear** | `bannerbear` | Consume Bannerbear API | ⚠️ Authentication Required | [`examples/nodes/integration/bannerbear/`](examples/nodes/integration/bannerbear/) |
| **Baserow** | `baserow` | Consume the Baserow API | ⚠️ Authentication Required | [`examples/nodes/integration/baserow/`](examples/nodes/integration/baserow/) |
| **Beeminder** | `beeminder` | Consume Beeminder API | ⚠️ Authentication Required | [`examples/nodes/integration/beeminder/`](examples/nodes/integration/beeminder/) |
| **Bitly** | `bitly` | Consume Bitly API | ⚠️ Authentication Required | [`examples/nodes/integration/bitly/`](examples/nodes/integration/bitly/) |
| **Bitwarden** | `bitwarden` | Consume the Bitwarden API | ⚠️ Authentication Required | [`examples/nodes/integration/bitwarden/`](examples/nodes/integration/bitwarden/) |
| **Box** | `box` | Consume Box API | ⚠️ Authentication Required | [`examples/nodes/integration/box/`](examples/nodes/integration/box/) |
| **Brandfetch** | `Brandfetch` | Consume Brandfetch API | ⚠️ Authentication Required | [`examples/nodes/integration/brandfetch/`](examples/nodes/integration/brandfetch/) |
| **Brevo** | `sendInBlue` | Consume Brevo API | ⚠️ Authentication Required | [`examples/nodes/integration/brevo/`](examples/nodes/integration/brevo/) |
| **Bubble** | `bubble` | Consume the Bubble Data API | ⚠️ Authentication Required | [`examples/nodes/integration/bubble/`](examples/nodes/integration/bubble/) |
| **Chargebee** | `chargebee` | Retrieve data from Chargebee API | ⚠️ Authentication Required | [`examples/nodes/integration/chargebee/`](examples/nodes/integration/chargebee/) |
| **CircleCI** | `circleCi` | Consume CircleCI API | ⚠️ Authentication Required | [`examples/nodes/integration/circleci/`](examples/nodes/integration/circleci/) |
| **Clearbit** | `clearbit` | Consume Clearbit API | ⚠️ Authentication Required | [`examples/nodes/integration/clearbit/`](examples/nodes/integration/clearbit/) |
| **ClickUp** | `clickUp` | Consume ClickUp API (Beta) | ⚠️ Authentication Required | [`examples/nodes/integration/clickup/`](examples/nodes/integration/clickup/) |
| **Clockify** | `clockify` | Consume Clockify REST API | ✅ None | [`examples/nodes/integration/clockify/`](examples/nodes/integration/clockify/) |
| **Cloudflare** | `cloudflare` | Consume Cloudflare API | ⚠️ Authentication Required | [`examples/nodes/integration/cloudflare/`](examples/nodes/integration/cloudflare/) |
| **Cockpit** | `cockpit` | Consume Cockpit API | ⚠️ Authentication Required | [`examples/nodes/integration/cockpit/`](examples/nodes/integration/cockpit/) |
| **Coda** | `coda` | Consume Coda API | ⚠️ Authentication Required | [`examples/nodes/integration/coda/`](examples/nodes/integration/coda/) |
| **CoinGecko** | `coinGecko` | Consume CoinGecko API | ⚠️ Authentication Required | [`examples/nodes/integration/coingecko/`](examples/nodes/integration/coingecko/) |
| **Compare Datasets** | `compareDatasets` | Compare two inputs for changes | ✅ None | [`examples/nodes/integration/compare-datasets/`](examples/nodes/integration/compare-datasets/) |
| **Compression** | `compression` | Compress and decompress files | ⚠️ Authentication Required | [`examples/nodes/integration/compression/`](examples/nodes/integration/compression/) |
| **Contentful** | `contentful` | Consume Contentful API | ⚠️ Authentication Required | [`examples/nodes/integration/contentful/`](examples/nodes/integration/contentful/) |
| **Convert to/from binary data** | `moveBinaryData` | Move data between binary and JSON properties | ⚠️ Authentication Required | [`examples/nodes/integration/convert-to-from-binary-data/`](examples/nodes/integration/convert-to-from-binary-data/) |
| **ConvertKit** | `convertKit` | Consume ConvertKit API | ⚠️ Authentication Required | [`examples/nodes/integration/convertkit/`](examples/nodes/integration/convertkit/) |
| **Copper** | `copper` | Consume the Copper API | ⚠️ Authentication Required | [`examples/nodes/integration/copper/`](examples/nodes/integration/copper/) |
| **Cortex** | `cortex` | Apply the Cortex analyzer/responder on the given entity | ⚠️ Authentication Required | [`examples/nodes/integration/cortex/`](examples/nodes/integration/cortex/) |
| **CrateDB** | `crateDb` | Add and update data in CrateDB | ⚠️ Authentication Required | [`examples/nodes/integration/cratedb/`](examples/nodes/integration/cratedb/) |
| **Cron** | `cron` | Triggers the workflow at a specific time | ✅ None | [`examples/nodes/integration/cron/`](examples/nodes/integration/cron/) |
| **crowd.dev** | `crowdDev` | crowd.dev is an open-source suite of community and data tools built to unlock community-led growth f | ⚠️ Authentication Required | [`examples/nodes/integration/crowd-dev/`](examples/nodes/integration/crowd-dev/) |
| **Crypto** | `crypto` | Provide cryptographic utilities | ⚠️ Authentication Required | [`examples/nodes/integration/crypto/`](examples/nodes/integration/crypto/) |
| **Customer Datastore (n8n training)** | `Jay Gatsby` | Dummy node used for n8n training | ⚠️ Authentication Required | [`examples/nodes/integration/customer-datastore-n8n-training/`](examples/nodes/integration/customer-datastore-n8n-training/) |
| **Customer Messenger (n8n training)** | `n8nTrainingCustomerMessenger` | Dummy node used for n8n training | ⚠️ Authentication Required | [`examples/nodes/integration/customer-messenger-n8n-training/`](examples/nodes/integration/customer-messenger-n8n-training/) |
| **Customer.io** | `customerIo` | Consume Customer.io API | ⚠️ Authentication Required | [`examples/nodes/integration/customer-io/`](examples/nodes/integration/customer-io/) |
| **Data table** | `dataTable` | Permanently save data across workflow executions in a table | ⚠️ Authentication Required | [`examples/nodes/integration/data-table/`](examples/nodes/integration/data-table/) |
| **Date & Time** | `dateTime` | Allows you to manipulate date and time values | ⚠️ Authentication Required | [`examples/nodes/integration/date-time/`](examples/nodes/integration/date-time/) |
| **DebugHelper** | `debugHelper` | Causes problems intentionally and generates useful data for debugging | ⚠️ Authentication Required | [`examples/nodes/integration/debughelper/`](examples/nodes/integration/debughelper/) |
| **DeepL** | `deepL` | Translate data using DeepL | ⚠️ Authentication Required | [`examples/nodes/integration/deepl/`](examples/nodes/integration/deepl/) |
| **Demio** | `demio` | Consume the Demio API | ⚠️ Authentication Required | [`examples/nodes/integration/demio/`](examples/nodes/integration/demio/) |
| **DHL** | `dhl` | Consume DHL API | ⚠️ Authentication Required | [`examples/nodes/integration/dhl/`](examples/nodes/integration/dhl/) |
| **Discord** | `discord` | Sends data to Discord | ⚠️ Authentication Required | [`examples/nodes/integration/discord/`](examples/nodes/integration/discord/) |
| **Discourse** | `discourse` | Consume Discourse API | ⚠️ Authentication Required | [`examples/nodes/integration/discourse/`](examples/nodes/integration/discourse/) |
| **Disqus** | `disqus` | Access data on Disqus | ⚠️ Authentication Required | [`examples/nodes/integration/disqus/`](examples/nodes/integration/disqus/) |
| **Drift** | `drift` | Consume Drift API | ✅ None | [`examples/nodes/integration/drift/`](examples/nodes/integration/drift/) |
| **Dropbox** | `dropbox` | Access data on Dropbox | ⚠️ Authentication Required | [`examples/nodes/integration/dropbox/`](examples/nodes/integration/dropbox/) |
| **Dropcontact** | `dropcontact` | Find B2B emails and enrich contacts | ⚠️ Authentication Required | [`examples/nodes/integration/dropcontact/`](examples/nodes/integration/dropcontact/) |
| **E-goi** | `egoi` | Consume E-goi API | ⚠️ Authentication Required | [`examples/nodes/integration/e-goi/`](examples/nodes/integration/e-goi/) |
| **E2E Test** | `e2eTest` | Dummy node used for e2e testing | ⚠️ Authentication Required | [`examples/nodes/integration/e2e-test/`](examples/nodes/integration/e2e-test/) |
| **Emelia** | `emelia` | Consume the Emelia API | ⚠️ Authentication Required | [`examples/nodes/integration/emelia/`](examples/nodes/integration/emelia/) |
| **ERPNext** | `erpNext` | Consume ERPNext API | ⚠️ Authentication Required | [`examples/nodes/integration/erpnext/`](examples/nodes/integration/erpnext/) |
| **Execute Command** | `executeCommand` | Executes a command on the host | ⚠️ Authentication Required | [`examples/nodes/integration/execute-command/`](examples/nodes/integration/execute-command/) |
| **Execution Data** | `executionData` | Add execution data for search | ⚠️ Authentication Required | [`examples/nodes/integration/execution-data/`](examples/nodes/integration/execution-data/) |
| **Extraction Values** | `extractionValues` | The key under which the extracted value should be saved | ⚠️ Authentication Required | [`examples/nodes/integration/extraction-values/`](examples/nodes/integration/extraction-values/) |
| **Facebook Graph API** | `facebookGraphApi` | Interacts with Facebook using the Graph API | ⚠️ API Key | [`examples/nodes/integration/facebook-graph-api/`](examples/nodes/integration/facebook-graph-api/) |
| **FileMaker** | `filemaker` | Retrieve data from the FileMaker data API | ⚠️ Authentication Required | [`examples/nodes/integration/filemaker/`](examples/nodes/integration/filemaker/) |
| **Filter** | `filter` | Remove items matching a condition | ✅ None | [`examples/nodes/integration/filter/`](examples/nodes/integration/filter/) |
| **Flow** | `flow` | Consume Flow API | ⚠️ Authentication Required | [`examples/nodes/integration/flow/`](examples/nodes/integration/flow/) |
| **Freshdesk** | `freshdesk` | Consume Freshdesk API | ⚠️ Authentication Required | [`examples/nodes/integration/freshdesk/`](examples/nodes/integration/freshdesk/) |
| **Freshservice** | `freshservice` | Consume the Freshservice API | ⚠️ Authentication Required | [`examples/nodes/integration/freshservice/`](examples/nodes/integration/freshservice/) |
| **Freshworks CRM** | `freshworksCrm` | Consume the Freshworks CRM API | ⚠️ Authentication Required | [`examples/nodes/integration/freshworks-crm/`](examples/nodes/integration/freshworks-crm/) |
| **FTP** | `ftp` | Transfer files via FTP or SFTP | ⚠️ Authentication Required | [`examples/nodes/integration/ftp/`](examples/nodes/integration/ftp/) |
| **Function** | `function` | Run custom function code which gets executed once and allows you to add, remove, change and replace  | ✅ None | [`examples/nodes/integration/function/`](examples/nodes/integration/function/) |
| **Function Item** | `functionItem` | Run custom function code which gets executed once per item | ✅ None | [`examples/nodes/integration/function-item/`](examples/nodes/integration/function-item/) |
| **GetResponse** | `getResponse` | Consume GetResponse API | ⚠️ Authentication Required | [`examples/nodes/integration/getresponse/`](examples/nodes/integration/getresponse/) |
| **Ghost** | `ghost` | Consume Ghost API | ⚠️ Authentication Required | [`examples/nodes/integration/ghost/`](examples/nodes/integration/ghost/) |
| **Git** | `git` | Control git. | ⚠️ Authentication Required | [`examples/nodes/integration/git/`](examples/nodes/integration/git/) |
| **GitHub** | `github` | Consume GitHub API | ⚠️ Authentication Required | [`examples/nodes/integration/github/`](examples/nodes/integration/github/) |
| **GitLab** | `gitlab` | Retrieve data from GitLab API | ⚠️ Authentication Required | [`examples/nodes/integration/gitlab/`](examples/nodes/integration/gitlab/) |
| **Gong** | `gong` | Interact with Gong API | ⚠️ Authentication Required | [`examples/nodes/integration/gong/`](examples/nodes/integration/gong/) |
| **Gotify** | `gotify` | Consume Gotify API | ✅ None | [`examples/nodes/integration/gotify/`](examples/nodes/integration/gotify/) |
| **GoToWebinar** | `goToWebinar` | Consume the GoToWebinar API | ⚠️ Authentication Required | [`examples/nodes/integration/gotowebinar/`](examples/nodes/integration/gotowebinar/) |
| **Grafana** | `grafana` | Consume the Grafana API | ⚠️ Authentication Required | [`examples/nodes/integration/grafana/`](examples/nodes/integration/grafana/) |
| **GraphQL** | `graphql` | Makes a GraphQL request and returns the received data | ⚠️ Authentication Required | [`examples/nodes/integration/graphql/`](examples/nodes/integration/graphql/) |
| **Grist** | `grist` | Consume the Grist API | ⚠️ Authentication Required | [`examples/nodes/integration/grist/`](examples/nodes/integration/grist/) |
| **Hacker News** | `hackerNews` | Consume Hacker News API | ⚠️ Authentication Required | [`examples/nodes/integration/hacker-news/`](examples/nodes/integration/hacker-news/) |
| **HaloPSA** | `haloPSA` | Consume HaloPSA API | ⚠️ Authentication Required | [`examples/nodes/integration/halopsa/`](examples/nodes/integration/halopsa/) |
| **Harvest** | `harvest` | Access data on Harvest | ⚠️ Authentication Required | [`examples/nodes/integration/harvest/`](examples/nodes/integration/harvest/) |
| **Help Scout** | `helpScout` | Consume Help Scout API | ⚠️ Authentication Required | [`examples/nodes/integration/help-scout/`](examples/nodes/integration/help-scout/) |
| **HighLevel** | `highLevel` | Consume HighLevel API | ⚠️ Authentication Required | [`examples/nodes/integration/highlevel/`](examples/nodes/integration/highlevel/) |
| **Home Assistant** | `homeAssistant` | Consume Home Assistant API | ⚠️ Authentication Required | [`examples/nodes/integration/home-assistant/`](examples/nodes/integration/home-assistant/) |
| **HTML Extract** | `htmlExtract` | Extracts data from HTML | ⚠️ Authentication Required | [`examples/nodes/integration/html-extract/`](examples/nodes/integration/html-extract/) |
| **HTTP Request** | `httpRequest` | Makes an HTTP request and returns the response data | ⚠️ Authentication Required | [`examples/nodes/integration/http-request/`](examples/nodes/integration/http-request/) |
| **HubSpot** | `hubspot` | Consume HubSpot API | ⚠️ Authentication Required | [`examples/nodes/integration/hubspot/`](examples/nodes/integration/hubspot/) |
| **Humantic AI** | `humanticAi` | Consume Humantic AI API | ⚠️ Authentication Required | [`examples/nodes/integration/humantic-ai/`](examples/nodes/integration/humantic-ai/) |
| **Hunter** | `hunter` | Consume Hunter API | ⚠️ Authentication Required | [`examples/nodes/integration/hunter/`](examples/nodes/integration/hunter/) |
| **iCalendar** | `iCal` | Create iCalendar file | ⚠️ Authentication Required | [`examples/nodes/integration/icalendar/`](examples/nodes/integration/icalendar/) |
| **Interact with Telegram using our pre-built** | `preBuiltAgentsCalloutTelegram` | Sends data to Telegram | ⚠️ Authentication Required | [`examples/nodes/integration/interact-with-telegram-using-our-pre-built/`](examples/nodes/integration/interact-with-telegram-using-our-pre-built/) |
| **Intercom** | `intercom` | Consume Intercom API | ⚠️ Authentication Required | [`examples/nodes/integration/intercom/`](examples/nodes/integration/intercom/) |
| **Interval** | `interval` | Triggers the workflow in a given interval | ⚠️ Authentication Required | [`examples/nodes/integration/interval/`](examples/nodes/integration/interval/) |
| **Invoice Ninja** | `invoiceNinja` | Consume Invoice Ninja API | ⚠️ Authentication Required | [`examples/nodes/integration/invoice-ninja/`](examples/nodes/integration/invoice-ninja/) |
| **Item Lists** | `itemLists` | Helper for working with lists of items and transforming arrays | ⚠️ Authentication Required | [`examples/nodes/integration/item-lists/`](examples/nodes/integration/item-lists/) |
| **Iterable** | `iterable` | Consume Iterable API | ⚠️ Authentication Required | [`examples/nodes/integration/iterable/`](examples/nodes/integration/iterable/) |
| **Jenkins** | `jenkins` | Consume Jenkins API | ⚠️ Authentication Required | [`examples/nodes/integration/jenkins/`](examples/nodes/integration/jenkins/) |
| **Jina AI** | `jinaAi` | Interact with Jina AI API | ⚠️ Authentication Required | [`examples/nodes/integration/jina-ai/`](examples/nodes/integration/jina-ai/) |
| **Jira Software** | `jira` | Consume Jira Software API | ⚠️ Authentication Required | [`examples/nodes/integration/jira-software/`](examples/nodes/integration/jira-software/) |
| **JWT** | `jwt` | Be sure to add a valid JWT token to the  | ⚠️ Authentication Required | [`examples/nodes/integration/jwt/`](examples/nodes/integration/jwt/) |
| **Kafka** | `kafka` | Sends messages to a Kafka topic | ⚠️ Authentication Required | [`examples/nodes/integration/kafka/`](examples/nodes/integration/kafka/) |
| **Keap** | `keap` | Consume Keap API | ⚠️ Authentication Required | [`examples/nodes/integration/keap/`](examples/nodes/integration/keap/) |
| **Kitemaker** | `kitemaker` | Consume the Kitemaker GraphQL API | ⚠️ Authentication Required | [`examples/nodes/integration/kitemaker/`](examples/nodes/integration/kitemaker/) |
| **KoBoToolbox** | `koBoToolbox` | Work with KoBoToolbox forms and submissions | ⚠️ Authentication Required | [`examples/nodes/integration/kobotoolbox/`](examples/nodes/integration/kobotoolbox/) |
| **Ldap** | `ldap` | Interact with LDAP servers | ⚠️ Authentication Required | [`examples/nodes/integration/ldap/`](examples/nodes/integration/ldap/) |
| **Lemlist** | `lemlist` | Consume the Lemlist API | ⚠️ Authentication Required | [`examples/nodes/integration/lemlist/`](examples/nodes/integration/lemlist/) |
| **Limit Wait Time** | `limitWaitTime` | Whether to limit the time this node should wait for a user response before execution resumes | ✅ None | [`examples/nodes/integration/limit-wait-time/`](examples/nodes/integration/limit-wait-time/) |
| **Line** | `line` | Consume Line API | ⚠️ Authentication Required | [`examples/nodes/integration/line/`](examples/nodes/integration/line/) |
| **Linear** | `linear` | Consume Linear API | ⚠️ Authentication Required | [`examples/nodes/integration/linear/`](examples/nodes/integration/linear/) |
| **LingvaNex** | `lingvaNex` | Consume LingvaNex API | ⚠️ Authentication Required | [`examples/nodes/integration/lingvanex/`](examples/nodes/integration/lingvanex/) |
| **LinkedIn** | `linkedIn` | Consume LinkedIn API | ⚠️ Authentication Required | [`examples/nodes/integration/linkedin/`](examples/nodes/integration/linkedin/) |
| **LoneScale** | `loneScale` | Create List, add / delete items | ⚠️ Authentication Required | [`examples/nodes/integration/lonescale/`](examples/nodes/integration/lonescale/) |
| **Magento 2** | `magento2` | Consume Magento API | ⚠️ Authentication Required | [`examples/nodes/integration/magento-2/`](examples/nodes/integration/magento-2/) |
| **Mailcheck** | `mailcheck` | Consume Mailcheck API | ⚠️ Authentication Required | [`examples/nodes/integration/mailcheck/`](examples/nodes/integration/mailcheck/) |
| **Mailchimp** | `mailchimp` | Consume Mailchimp API | ⚠️ Authentication Required | [`examples/nodes/integration/mailchimp/`](examples/nodes/integration/mailchimp/) |
| **MailerLite** | `mailerLite` | Consume MailerLite API | ⚠️ Authentication Required | [`examples/nodes/integration/mailerlite/`](examples/nodes/integration/mailerlite/) |
| **Mailgun** | `mailgun` | Sends an email via Mailgun | ⚠️ Authentication Required | [`examples/nodes/integration/mailgun/`](examples/nodes/integration/mailgun/) |
| **Mailjet** | `mailjet` | Consume Mailjet API | ⚠️ Authentication Required | [`examples/nodes/integration/mailjet/`](examples/nodes/integration/mailjet/) |
| **Mandrill** | `mandrill` | Consume Mandrill API | ⚠️ Authentication Required | [`examples/nodes/integration/mandrill/`](examples/nodes/integration/mandrill/) |
| **Markdown** | `markdown` | Convert data between Markdown and HTML | ⚠️ Authentication Required | [`examples/nodes/integration/markdown/`](examples/nodes/integration/markdown/) |
| **Marketstack** | `marketstack` | Consume Marketstack API | ⚠️ Authentication Required | [`examples/nodes/integration/marketstack/`](examples/nodes/integration/marketstack/) |
| **Matrix** | `matrix` | Consume Matrix API | ⚠️ Authentication Required | [`examples/nodes/integration/matrix/`](examples/nodes/integration/matrix/) |
| **Mattermost** | `mattermost` | Sends data to Mattermost | ⚠️ Authentication Required | [`examples/nodes/integration/mattermost/`](examples/nodes/integration/mattermost/) |
| **Mautic** | `mautic` | Consume Mautic API | ⚠️ Authentication Required | [`examples/nodes/integration/mautic/`](examples/nodes/integration/mautic/) |
| **Medium** | `medium` | Consume Medium API | ⚠️ Authentication Required | [`examples/nodes/integration/medium/`](examples/nodes/integration/medium/) |
| **MessageBird** | `messageBird` | Sends SMS via MessageBird | ⚠️ Authentication Required | [`examples/nodes/integration/messagebird/`](examples/nodes/integration/messagebird/) |
| **Metabase** | `metabase` | Use the Metabase API | ⚠️ Authentication Required | [`examples/nodes/integration/metabase/`](examples/nodes/integration/metabase/) |
| **Mindee** | `mindee` | Consume Mindee API | ⚠️ Authentication Required | [`examples/nodes/integration/mindee/`](examples/nodes/integration/mindee/) |
| **MISP** | `misp` | Consume the MISP API | ⚠️ Authentication Required | [`examples/nodes/integration/misp/`](examples/nodes/integration/misp/) |
| **Mistral AI** | `mistralAi` | Consume Mistral AI API | ⚠️ Authentication Required | [`examples/nodes/integration/mistral-ai/`](examples/nodes/integration/mistral-ai/) |
| **Mocean** | `mocean` | Send SMS and voice messages via Mocean | ⚠️ Authentication Required | [`examples/nodes/integration/mocean/`](examples/nodes/integration/mocean/) |
| **Monday.com** | `mondayCom` | Consume Monday.com API | ⚠️ Authentication Required | [`examples/nodes/integration/monday-com/`](examples/nodes/integration/monday-com/) |
| **MongoDB** | `mongoDb` | Find, insert and update documents in MongoDB | ⚠️ Authentication Required | [`examples/nodes/integration/mongodb/`](examples/nodes/integration/mongodb/) |
| **Monica CRM** | `monicaCrm` | Consume the Monica CRM API | ⚠️ Authentication Required | [`examples/nodes/integration/monica-crm/`](examples/nodes/integration/monica-crm/) |
| **MQTT** | `mqtt` | Push messages to MQTT | ⚠️ Authentication Required | [`examples/nodes/integration/mqtt/`](examples/nodes/integration/mqtt/) |
| **MSG91** | `msg91` | Sends transactional SMS via MSG91 | ⚠️ Authentication Required | [`examples/nodes/integration/msg91/`](examples/nodes/integration/msg91/) |
| **MySQL** | `mySql` | Get, add and update data in MySQL | ⚠️ Authentication Required | [`examples/nodes/integration/mysql/`](examples/nodes/integration/mysql/) |
| **n8n** | `n8n` | Handle events and perform actions on your n8n instance | ⚠️ Authentication Required | [`examples/nodes/integration/n8n/`](examples/nodes/integration/n8n/) |
| **NASA** | `nasa` | Retrieve data from the NASA API | ⚠️ Authentication Required | [`examples/nodes/integration/nasa/`](examples/nodes/integration/nasa/) |
| **Netlify** | `netlify` | Consume Netlify API | ✅ None | [`examples/nodes/integration/netlify/`](examples/nodes/integration/netlify/) |
| **Nextcloud** | `nextCloud` | Access data on Nextcloud | ⚠️ Authentication Required | [`examples/nodes/integration/nextcloud/`](examples/nodes/integration/nextcloud/) |
| **No Operation, do nothing** | `noOp` | No Operation | ⚠️ Authentication Required | [`examples/nodes/integration/no-operation-do-nothing/`](examples/nodes/integration/no-operation-do-nothing/) |
| **NocoDB** | `nocoDb` | Read, update, write and delete data from NocoDB | ⚠️ Authentication Required | [`examples/nodes/integration/nocodb/`](examples/nodes/integration/nocodb/) |
| **Notion** | `notion` | Consume Notion API | ⚠️ Authentication Required | [`examples/nodes/integration/notion/`](examples/nodes/integration/notion/) |
| **Npm** | `npm` | Consume NPM registry API | ⚠️ Authentication Required | [`examples/nodes/integration/npm/`](examples/nodes/integration/npm/) |
| **Odoo** | `odoo` | Consume Odoo API | ⚠️ Authentication Required | [`examples/nodes/integration/odoo/`](examples/nodes/integration/odoo/) |
| **Okta** | `okta` | Use the Okta API | ⚠️ Authentication Required | [`examples/nodes/integration/okta/`](examples/nodes/integration/okta/) |
| **One Simple API** | `oneSimpleApi` | A toolbox of no-code utilities | ⚠️ API Key | [`examples/nodes/integration/one-simple-api/`](examples/nodes/integration/one-simple-api/) |
| **Onfleet** | `onfleet` | Consume Onfleet API | ⚠️ Authentication Required | [`examples/nodes/integration/onfleet/`](examples/nodes/integration/onfleet/) |
| **OpenAI** | `openAi` | Consume Open AI | ⚠️ Authentication Required | [`examples/nodes/integration/openai/`](examples/nodes/integration/openai/) |
| **OpenThesaurus** | `openThesaurus` | Get synonmns for German words using the OpenThesaurus API | ⚠️ Authentication Required | [`examples/nodes/integration/openthesaurus/`](examples/nodes/integration/openthesaurus/) |
| **OpenWeatherMap** | `openWeatherMap` | Gets current and future weather information | ⚠️ Authentication Required | [`examples/nodes/integration/openweathermap/`](examples/nodes/integration/openweathermap/) |
| **Orbit** | `orbit` | Consume Orbit API | ⚠️ Authentication Required | [`examples/nodes/integration/orbit/`](examples/nodes/integration/orbit/) |
| **Oura** | `oura` | Consume Oura API | ⚠️ Authentication Required | [`examples/nodes/integration/oura/`](examples/nodes/integration/oura/) |
| **Paddle** | `paddle` | Consume Paddle API | ⚠️ Authentication Required | [`examples/nodes/integration/paddle/`](examples/nodes/integration/paddle/) |
| **PagerDuty** | `pagerDuty` | Consume PagerDuty API | ⚠️ Authentication Required | [`examples/nodes/integration/pagerduty/`](examples/nodes/integration/pagerduty/) |
| **PayPal** | `payPal` | Consume PayPal API | ⚠️ Authentication Required | [`examples/nodes/integration/paypal/`](examples/nodes/integration/paypal/) |
| **Peekalink** | `peekalink` | Consume the Peekalink API | ⚠️ Authentication Required | [`examples/nodes/integration/peekalink/`](examples/nodes/integration/peekalink/) |
| **Perplexity** | `perplexity` | Interact with the Perplexity API to generate AI responses with citations | ⚠️ Authentication Required | [`examples/nodes/integration/perplexity/`](examples/nodes/integration/perplexity/) |
| **Phantombuster** | `phantombuster` | Consume Phantombuster API | ⚠️ Authentication Required | [`examples/nodes/integration/phantombuster/`](examples/nodes/integration/phantombuster/) |
| **Philips Hue** | `philipsHue` | Consume Philips Hue API | ⚠️ Authentication Required | [`examples/nodes/integration/philips-hue/`](examples/nodes/integration/philips-hue/) |
| **Pipedrive** | `pipedrive` | Create and edit data in Pipedrive | ⚠️ Authentication Required | [`examples/nodes/integration/pipedrive/`](examples/nodes/integration/pipedrive/) |
| **Plivo** | `plivo` | Send SMS/MMS messages or make phone calls | ⚠️ Authentication Required | [`examples/nodes/integration/plivo/`](examples/nodes/integration/plivo/) |
| **PostBin** | `postBin` | Consume PostBin API | ⚠️ Authentication Required | [`examples/nodes/integration/postbin/`](examples/nodes/integration/postbin/) |
| **Postgres** | `postgres` | Get, add and update data in Postgres | ⚠️ Authentication Required | [`examples/nodes/integration/postgres/`](examples/nodes/integration/postgres/) |
| **PostHog** | `postHog` | Consume PostHog API | ⚠️ Authentication Required | [`examples/nodes/integration/posthog/`](examples/nodes/integration/posthog/) |
| **ProfitWell** | `profitWell` | Consume ProfitWell API | ⚠️ Authentication Required | [`examples/nodes/integration/profitwell/`](examples/nodes/integration/profitwell/) |
| **Pushbullet** | `pushbullet` | Consume Pushbullet API | ⚠️ Authentication Required | [`examples/nodes/integration/pushbullet/`](examples/nodes/integration/pushbullet/) |
| **Pushcut** | `pushcut` | Consume Pushcut API | ⚠️ Authentication Required | [`examples/nodes/integration/pushcut/`](examples/nodes/integration/pushcut/) |
| **Pushover** | `pushover` | Consume Pushover API | ⚠️ Authentication Required | [`examples/nodes/integration/pushover/`](examples/nodes/integration/pushover/) |
| **QuestDB** | `questDb` | Get, add and update data in QuestDB | ⚠️ Authentication Required | [`examples/nodes/integration/questdb/`](examples/nodes/integration/questdb/) |
| **Quick Base** | `quickbase` | Integrate with the Quick Base RESTful API | ⚠️ Authentication Required | [`examples/nodes/integration/quick-base/`](examples/nodes/integration/quick-base/) |
| **QuickBooks Online** | `quickbooks` | Consume the QuickBooks Online API | ⚠️ Authentication Required | [`examples/nodes/integration/quickbooks-online/`](examples/nodes/integration/quickbooks-online/) |
| **QuickChart** | `quickChart` | Create a chart via QuickChart | ⚠️ Authentication Required | [`examples/nodes/integration/quickchart/`](examples/nodes/integration/quickchart/) |
| **RabbitMQ** | `rabbitmq` | Sends messages to a RabbitMQ topic | ⚠️ Authentication Required | [`examples/nodes/integration/rabbitmq/`](examples/nodes/integration/rabbitmq/) |
| **Raindrop** | `raindrop` | Consume the Raindrop API | ⚠️ Authentication Required | [`examples/nodes/integration/raindrop/`](examples/nodes/integration/raindrop/) |
| **Read Binary File** | `readBinaryFile` | Reads a binary file from disk | ⚠️ Authentication Required | [`examples/nodes/integration/read-binary-file/`](examples/nodes/integration/read-binary-file/) |
| **Read Binary Files** | `readBinaryFiles` | Reads binary files from disk | ⚠️ Authentication Required | [`examples/nodes/integration/read-binary-files/`](examples/nodes/integration/read-binary-files/) |
| **Read PDF** | `readPDF` | Reads a PDF and extracts its content | ⚠️ Authentication Required | [`examples/nodes/integration/read-pdf/`](examples/nodes/integration/read-pdf/) |
| **Reddit** | `reddit` | Consume the Reddit API | ⚠️ Authentication Required | [`examples/nodes/integration/reddit/`](examples/nodes/integration/reddit/) |
| **Redis** | `redis` | Get, send and update data in Redis | ⚠️ Authentication Required | [`examples/nodes/integration/redis/`](examples/nodes/integration/redis/) |
| **Rename Keys** | `renameKeys` | Update item field names | ⚠️ Authentication Required | [`examples/nodes/integration/rename-keys/`](examples/nodes/integration/rename-keys/) |
| **Respond With** | `respondWith` | Respond with all input JSON items | ⚠️ Authentication Required | [`examples/nodes/integration/respond-with/`](examples/nodes/integration/respond-with/) |
| **RocketChat** | `rocketchat` | Consume RocketChat API | ⚠️ Authentication Required | [`examples/nodes/integration/rocketchat/`](examples/nodes/integration/rocketchat/) |
| **RSS Read** | `rssFeedRead` | Reads data from an RSS Feed | ⚠️ Authentication Required | [`examples/nodes/integration/rss-read/`](examples/nodes/integration/rss-read/) |
| **Rundeck** | `rundeck` | Manage Rundeck API | ⚠️ Authentication Required | [`examples/nodes/integration/rundeck/`](examples/nodes/integration/rundeck/) |
| **S3** | `s3` | Sends data to any S3-compatible service | ⚠️ Authentication Required | [`examples/nodes/integration/s3/`](examples/nodes/integration/s3/) |
| **Salesforce** | `salesforce` | Consume Salesforce API | ⚠️ Authentication Required | [`examples/nodes/integration/salesforce/`](examples/nodes/integration/salesforce/) |
| **Salesmate** | `salesmate` | Consume Salesmate API | ⚠️ Authentication Required | [`examples/nodes/integration/salesmate/`](examples/nodes/integration/salesmate/) |
| **Schedule Trigger** | `scheduleTrigger` | Triggers the workflow on a given schedule | ✅ None | [`examples/nodes/integration/schedule-trigger/`](examples/nodes/integration/schedule-trigger/) |
| **SeaTable** | `seaTable` | Read, update, write and delete data from SeaTable | ⚠️ Authentication Required | [`examples/nodes/integration/seatable/`](examples/nodes/integration/seatable/) |
| **SecurityScorecard** | `securityScorecard` | Consume SecurityScorecard API | ⚠️ Authentication Required | [`examples/nodes/integration/securityscorecard/`](examples/nodes/integration/securityscorecard/) |
| **Segment** | `segment` | Consume Segment API | ⚠️ Authentication Required | [`examples/nodes/integration/segment/`](examples/nodes/integration/segment/) |
| **Send Email** | `emailSend` | Sends an email using SMTP protocol | ⚠️ Authentication Required | [`examples/nodes/integration/send-email/`](examples/nodes/integration/send-email/) |
| **SendGrid** | `sendGrid` | Consume SendGrid API | ⚠️ Authentication Required | [`examples/nodes/integration/sendgrid/`](examples/nodes/integration/sendgrid/) |
| **Sendy** | `sendy` | Consume Sendy API | ⚠️ Authentication Required | [`examples/nodes/integration/sendy/`](examples/nodes/integration/sendy/) |
| **Sentry.io** | `sentryIo` | Consume Sentry.io API | ⚠️ Authentication Required | [`examples/nodes/integration/sentry-io/`](examples/nodes/integration/sentry-io/) |
| **ServiceNow** | `serviceNow` | Consume ServiceNow API | ⚠️ Authentication Required | [`examples/nodes/integration/servicenow/`](examples/nodes/integration/servicenow/) |
| **seven** | `sms77` | Send SMS and make text-to-speech calls | ⚠️ Authentication Required | [`examples/nodes/integration/seven/`](examples/nodes/integration/seven/) |
| **Shopify** | `shopify` | Consume Shopify API | ✅ None | [`examples/nodes/integration/shopify/`](examples/nodes/integration/shopify/) |
| **SIGNL4** | `signl4` | Consume SIGNL4 API | ⚠️ Authentication Required | [`examples/nodes/integration/signl4/`](examples/nodes/integration/signl4/) |
| **Simulate** | `simulate` | Simulate a node | ⚠️ Authentication Required | [`examples/nodes/integration/simulate/`](examples/nodes/integration/simulate/) |
| **Slack** | `slack` | Consume Slack API | ⚠️ Authentication Required | [`examples/nodes/integration/slack/`](examples/nodes/integration/slack/) |
| **Snowflake** | `snowflake` | Get, add and update data in Snowflake | ⚠️ Authentication Required | [`examples/nodes/integration/snowflake/`](examples/nodes/integration/snowflake/) |
| **Split In Batches** | `splitInBatches` | Split data into batches and iterate over each batch | ✅ None | [`examples/nodes/integration/split-in-batches/`](examples/nodes/integration/split-in-batches/) |
| **Splunk** | `splunk` | Consume the Splunk Enterprise API | ⚠️ Authentication Required | [`examples/nodes/integration/splunk/`](examples/nodes/integration/splunk/) |
| **Spontit** | `spontit` | Consume Spontit API | ⚠️ Authentication Required | [`examples/nodes/integration/spontit/`](examples/nodes/integration/spontit/) |
| **Spotify** | `spotify` | Access public song data via the Spotify API | ✅ None | [`examples/nodes/integration/spotify/`](examples/nodes/integration/spotify/) |
| **Spreadsheet File** | `spreadsheetFile` | Reads and writes data from a spreadsheet file like CSV, XLS, ODS, etc | ⚠️ Authentication Required | [`examples/nodes/integration/spreadsheet-file/`](examples/nodes/integration/spreadsheet-file/) |
| **SSH** | `ssh` | Execute commands via SSH | ⚠️ Authentication Required | [`examples/nodes/integration/ssh/`](examples/nodes/integration/ssh/) |
| **Stackby** | `stackby` | Read, write, and delete data in Stackby | ⚠️ Authentication Required | [`examples/nodes/integration/stackby/`](examples/nodes/integration/stackby/) |
| **Start** | `start` | Starts the workflow execution from this node | ⚠️ Authentication Required | [`examples/nodes/integration/start/`](examples/nodes/integration/start/) |
| **Sticky Note** | `stickyNote` | Make your workflow easier to understand | ⚠️ Authentication Required | [`examples/nodes/integration/sticky-note/`](examples/nodes/integration/sticky-note/) |
| **Stop and Error** | `stopAndError` | Throw an error in the workflow | ✅ None | [`examples/nodes/integration/stop-and-error/`](examples/nodes/integration/stop-and-error/) |
| **Storyblok** | `storyblok` | Consume Storyblok API | ⚠️ Authentication Required | [`examples/nodes/integration/storyblok/`](examples/nodes/integration/storyblok/) |
| **Strapi** | `strapi` | Consume Strapi API | ⚠️ API Key | [`examples/nodes/integration/strapi/`](examples/nodes/integration/strapi/) |
| **Strava** | `strava` | Consume Strava API | ⚠️ Authentication Required | [`examples/nodes/integration/strava/`](examples/nodes/integration/strava/) |
| **Stripe** | `stripe` | Consume the Stripe API | ⚠️ Authentication Required | [`examples/nodes/integration/stripe/`](examples/nodes/integration/stripe/) |
| **Supabase** | `supabase` | Add, get, delete and update data in a table | ⚠️ Authentication Required | [`examples/nodes/integration/supabase/`](examples/nodes/integration/supabase/) |
| **SyncroMSP** | `syncroMsp` | Manage contacts, tickets and more from Syncro MSP | ⚠️ Authentication Required | [`examples/nodes/integration/syncromsp/`](examples/nodes/integration/syncromsp/) |
| **Taiga** | `taiga` | Consume Taiga API | ⚠️ Authentication Required | [`examples/nodes/integration/taiga/`](examples/nodes/integration/taiga/) |
| **Tapfiliate** | `tapfiliate` | Consume Tapfiliate API | ⚠️ Authentication Required | [`examples/nodes/integration/tapfiliate/`](examples/nodes/integration/tapfiliate/) |
| **TheHive** | `theHive` | Consume TheHive API | ⚠️ Authentication Required | [`examples/nodes/integration/thehive/`](examples/nodes/integration/thehive/) |
| **TheHiveProject** | `n8n-nodes-base.thehiveproject` | N/A | ⚠️ Authentication Required | [`examples/nodes/integration/thehiveproject/`](examples/nodes/integration/thehiveproject/) |
| **TimescaleDB** | `timescaleDb` | Add and update data in TimescaleDB | ⚠️ Authentication Required | [`examples/nodes/integration/timescaledb/`](examples/nodes/integration/timescaledb/) |
| **Todoist** | `todoist` | Consume Todoist API | ⚠️ Authentication Required | [`examples/nodes/integration/todoist/`](examples/nodes/integration/todoist/) |
| **TOTP** | `totp` | Generate a time-based one-time password | ⚠️ Authentication Required | [`examples/nodes/integration/totp/`](examples/nodes/integration/totp/) |
| **TravisCI** | `travisCi` | Consume TravisCI API | ⚠️ Authentication Required | [`examples/nodes/integration/travisci/`](examples/nodes/integration/travisci/) |
| **Trello** | `trello` | Create, change and delete boards and cards | ⚠️ Authentication Required | [`examples/nodes/integration/trello/`](examples/nodes/integration/trello/) |
| **Twake** | `twake` | Consume Twake API | ⚠️ Authentication Required | [`examples/nodes/integration/twake/`](examples/nodes/integration/twake/) |
| **Twilio** | `twilio` | Send SMS and WhatsApp messages or make phone calls | ⚠️ Authentication Required | [`examples/nodes/integration/twilio/`](examples/nodes/integration/twilio/) |
| **Twist** | `twist` | Consume Twist API | ⚠️ Authentication Required | [`examples/nodes/integration/twist/`](examples/nodes/integration/twist/) |
| **Unleashed Software** | `unleashedSoftware` | Consume Unleashed Software API | ⚠️ Authentication Required | [`examples/nodes/integration/unleashed-software/`](examples/nodes/integration/unleashed-software/) |
| **Uplead** | `uplead` | Consume Uplead API | ⚠️ Authentication Required | [`examples/nodes/integration/uplead/`](examples/nodes/integration/uplead/) |
| **uProc** | `uproc` | Consume uProc API | ⚠️ Authentication Required | [`examples/nodes/integration/uproc/`](examples/nodes/integration/uproc/) |
| **UptimeRobot** | `uptimeRobot` | Consume UptimeRobot API | ⚠️ Authentication Required | [`examples/nodes/integration/uptimerobot/`](examples/nodes/integration/uptimerobot/) |
| **urlscan.io** | `urlScanIo` | Provides various utilities for monitoring websites like health checks or screenshots | ⚠️ Authentication Required | [`examples/nodes/integration/urlscan-io/`](examples/nodes/integration/urlscan-io/) |
| **Vero** | `vero` | Consume Vero API | ⚠️ Authentication Required | [`examples/nodes/integration/vero/`](examples/nodes/integration/vero/) |
| **Vonage** | `vonage` | Consume Vonage API | ⚠️ Authentication Required | [`examples/nodes/integration/vonage/`](examples/nodes/integration/vonage/) |
| **Wait Amount** | `amount` | The time to wait | ✅ None | [`examples/nodes/integration/wait-amount/`](examples/nodes/integration/wait-amount/) |
| **Webflow** | `webflow` | Consume the Webflow API | ⚠️ Authentication Required | [`examples/nodes/integration/webflow/`](examples/nodes/integration/webflow/) |
| **Wekan** | `wekan` | Consume Wekan API | ⚠️ Authentication Required | [`examples/nodes/integration/wekan/`](examples/nodes/integration/wekan/) |
| **WhatsApp Business Cloud** | `whatsApp` | Access WhatsApp API | ⚠️ Authentication Required | [`examples/nodes/integration/whatsapp-business-cloud/`](examples/nodes/integration/whatsapp-business-cloud/) |
| **Wise** | `wise` | Consume the Wise API | ⚠️ Authentication Required | [`examples/nodes/integration/wise/`](examples/nodes/integration/wise/) |
| **WooCommerce** | `wooCommerce` | Consume WooCommerce API | ⚠️ Authentication Required | [`examples/nodes/integration/woocommerce/`](examples/nodes/integration/woocommerce/) |
| **Wordpress** | `wordpress` | Consume Wordpress API | ⚠️ Authentication Required | [`examples/nodes/integration/wordpress/`](examples/nodes/integration/wordpress/) |
| **Write Binary File** | `writeBinaryFile` | Writes a binary file to disk | ⚠️ Authentication Required | [`examples/nodes/integration/write-binary-file/`](examples/nodes/integration/write-binary-file/) |
| **X (Formerly Twitter)** | `twitter` | Consume the X API | ⚠️ Authentication Required | [`examples/nodes/integration/x-formerly-twitter/`](examples/nodes/integration/x-formerly-twitter/) |
| **Xero** | `xero` | Consume Xero API | ⚠️ Authentication Required | [`examples/nodes/integration/xero/`](examples/nodes/integration/xero/) |
| **XML** | `xml` | Convert data from and to XML | ⚠️ Authentication Required | [`examples/nodes/integration/xml/`](examples/nodes/integration/xml/) |
| **Yourls** | `yourls` | Consume Yourls API | ⚠️ Authentication Required | [`examples/nodes/integration/yourls/`](examples/nodes/integration/yourls/) |
| **Zammad** | `zammad` | Consume the Zammad API | ⚠️ Authentication Required | [`examples/nodes/integration/zammad/`](examples/nodes/integration/zammad/) |
| **Zendesk** | `zendesk` | Consume Zendesk API | ⚠️ Authentication Required | [`examples/nodes/integration/zendesk/`](examples/nodes/integration/zendesk/) |
| **Zoho CRM** | `zohoCrm` | Consume Zoho CRM API | ⚠️ Authentication Required | [`examples/nodes/integration/zoho-crm/`](examples/nodes/integration/zoho-crm/) |
| **Zoom** | `zoom` | Consume Zoom API | ⚠️ Authentication Required | [`examples/nodes/integration/zoom/`](examples/nodes/integration/zoom/) |
| **Zulip** | `zulip` | Consume Zulip API | ⚠️ Authentication Required | [`examples/nodes/integration/zulip/`](examples/nodes/integration/zulip/) |

---

## Credential Requirements

Some nodes require external service credentials (API keys, OAuth tokens, etc.).

### Nodes Requiring Credentials (252)

#### Authentication Required (248)

- **Action Network** - [`actionNetwork`](examples/nodes/integration/action-network/)
- **ActiveCampaign** - [`activeCampaign`](examples/nodes/integration/activecampaign/)
- **Adalo** - [`adalo`](examples/nodes/integration/adalo/)
- **Affinity** - [`affinity`](examples/nodes/integration/affinity/)
- **Agile CRM** - [`agileCrm`](examples/nodes/integration/agile-crm/)
- **AI Transform** - [`aiTransform`](examples/nodes/integration/ai-transform/)
- **Airtable** - [`airtable`](examples/nodes/integration/airtable/)
- **Airtop** - [`airtop`](examples/nodes/integration/airtop/)
- **AMQP Sender** - [`amqp`](examples/nodes/integration/amqp-sender/)
- **Asana** - [`asana`](examples/nodes/integration/asana/)
- **Automizy** - [`automizy`](examples/nodes/integration/automizy/)
- **Autopilot** - [`autopilot`](examples/nodes/integration/autopilot/)
- **AWS Lambda** - [`awsLambda`](examples/nodes/integration/aws-lambda/)
- **Background Color** - [`Blur`](examples/nodes/integration/background-color/)
- **BambooHr** - [`n8n-nodes-base.bamboohr`](examples/nodes/integration/bamboohr/)
- **Bannerbear** - [`bannerbear`](examples/nodes/integration/bannerbear/)
- **Baserow** - [`baserow`](examples/nodes/integration/baserow/)
- **Beeminder** - [`beeminder`](examples/nodes/integration/beeminder/)
- **Bitly** - [`bitly`](examples/nodes/integration/bitly/)
- **Bitwarden** - [`bitwarden`](examples/nodes/integration/bitwarden/)
- **Box** - [`box`](examples/nodes/integration/box/)
- **Brandfetch** - [`Brandfetch`](examples/nodes/integration/brandfetch/)
- **Brevo** - [`sendInBlue`](examples/nodes/integration/brevo/)
- **Bubble** - [`bubble`](examples/nodes/integration/bubble/)
- **Chargebee** - [`chargebee`](examples/nodes/integration/chargebee/)
- **CircleCI** - [`circleCi`](examples/nodes/integration/circleci/)
- **Clearbit** - [`clearbit`](examples/nodes/integration/clearbit/)
- **ClickUp** - [`clickUp`](examples/nodes/integration/clickup/)
- **Cloudflare** - [`cloudflare`](examples/nodes/integration/cloudflare/)
- **Cockpit** - [`cockpit`](examples/nodes/integration/cockpit/)
- **Coda** - [`coda`](examples/nodes/integration/coda/)
- **CoinGecko** - [`coinGecko`](examples/nodes/integration/coingecko/)
- **Compression** - [`compression`](examples/nodes/integration/compression/)
- **Contentful** - [`contentful`](examples/nodes/integration/contentful/)
- **Convert to/from binary data** - [`moveBinaryData`](examples/nodes/integration/convert-to-from-binary-data/)
- **ConvertKit** - [`convertKit`](examples/nodes/integration/convertkit/)
- **Copper** - [`copper`](examples/nodes/integration/copper/)
- **Cortex** - [`cortex`](examples/nodes/integration/cortex/)
- **CrateDB** - [`crateDb`](examples/nodes/integration/cratedb/)
- **crowd.dev** - [`crowdDev`](examples/nodes/integration/crowd-dev/)
- **Crypto** - [`crypto`](examples/nodes/integration/crypto/)
- **Customer Datastore (n8n training)** - [`Jay Gatsby`](examples/nodes/integration/customer-datastore-n8n-training/)
- **Customer Messenger (n8n training)** - [`n8nTrainingCustomerMessenger`](examples/nodes/integration/customer-messenger-n8n-training/)
- **Customer.io** - [`customerIo`](examples/nodes/integration/customer-io/)
- **Data table** - [`dataTable`](examples/nodes/integration/data-table/)
- **Date & Time** - [`dateTime`](examples/nodes/integration/date-time/)
- **DebugHelper** - [`debugHelper`](examples/nodes/integration/debughelper/)
- **DeepL** - [`deepL`](examples/nodes/integration/deepl/)
- **Demio** - [`demio`](examples/nodes/integration/demio/)
- **DHL** - [`dhl`](examples/nodes/integration/dhl/)
- **Discord** - [`discord`](examples/nodes/integration/discord/)
- **Discourse** - [`discourse`](examples/nodes/integration/discourse/)
- **Disqus** - [`disqus`](examples/nodes/integration/disqus/)
- **Dropbox** - [`dropbox`](examples/nodes/integration/dropbox/)
- **Dropcontact** - [`dropcontact`](examples/nodes/integration/dropcontact/)
- **E-goi** - [`egoi`](examples/nodes/integration/e-goi/)
- **E2E Test** - [`e2eTest`](examples/nodes/integration/e2e-test/)
- **Emelia** - [`emelia`](examples/nodes/integration/emelia/)
- **ERPNext** - [`erpNext`](examples/nodes/integration/erpnext/)
- **Execute Command** - [`executeCommand`](examples/nodes/integration/execute-command/)
- **Execution Data** - [`executionData`](examples/nodes/integration/execution-data/)
- **Extraction Values** - [`extractionValues`](examples/nodes/integration/extraction-values/)
- **Facebook Lead Ads Trigger** - [`facebookLeadAdsTrigger`](examples/nodes/trigger/facebook-lead-ads-trigger/)
- **FileMaker** - [`filemaker`](examples/nodes/integration/filemaker/)
- **Flow** - [`flow`](examples/nodes/integration/flow/)
- **Freshdesk** - [`freshdesk`](examples/nodes/integration/freshdesk/)
- **Freshservice** - [`freshservice`](examples/nodes/integration/freshservice/)
- **Freshworks CRM** - [`freshworksCrm`](examples/nodes/integration/freshworks-crm/)
- **FTP** - [`ftp`](examples/nodes/integration/ftp/)
- **GetResponse** - [`getResponse`](examples/nodes/integration/getresponse/)
- **Ghost** - [`ghost`](examples/nodes/integration/ghost/)
- **Git** - [`git`](examples/nodes/integration/git/)
- **GitHub** - [`github`](examples/nodes/integration/github/)
- **GitLab** - [`gitlab`](examples/nodes/integration/gitlab/)
- **Gong** - [`gong`](examples/nodes/integration/gong/)
- **GoToWebinar** - [`goToWebinar`](examples/nodes/integration/gotowebinar/)
- **Grafana** - [`grafana`](examples/nodes/integration/grafana/)
- **GraphQL** - [`graphql`](examples/nodes/integration/graphql/)
- **Grist** - [`grist`](examples/nodes/integration/grist/)
- **Hacker News** - [`hackerNews`](examples/nodes/integration/hacker-news/)
- **HaloPSA** - [`haloPSA`](examples/nodes/integration/halopsa/)
- **Harvest** - [`harvest`](examples/nodes/integration/harvest/)
- **Help Scout** - [`helpScout`](examples/nodes/integration/help-scout/)
- **HighLevel** - [`highLevel`](examples/nodes/integration/highlevel/)
- **Home Assistant** - [`homeAssistant`](examples/nodes/integration/home-assistant/)
- **HTML Extract** - [`htmlExtract`](examples/nodes/integration/html-extract/)
- **HTTP Request** - [`httpRequest`](examples/nodes/integration/http-request/)
- **HubSpot** - [`hubspot`](examples/nodes/integration/hubspot/)
- **Humantic AI** - [`humanticAi`](examples/nodes/integration/humantic-ai/)
- **Hunter** - [`hunter`](examples/nodes/integration/hunter/)
- **iCalendar** - [`iCal`](examples/nodes/integration/icalendar/)
- **Interact with Telegram using our pre-built** - [`preBuiltAgentsCalloutTelegram`](examples/nodes/integration/interact-with-telegram-using-our-pre-built/)
- **Intercom** - [`intercom`](examples/nodes/integration/intercom/)
- **Interval** - [`interval`](examples/nodes/integration/interval/)
- **Invoice Ninja** - [`invoiceNinja`](examples/nodes/integration/invoice-ninja/)
- **Item Lists** - [`itemLists`](examples/nodes/integration/item-lists/)
- **Iterable** - [`iterable`](examples/nodes/integration/iterable/)
- **Jenkins** - [`jenkins`](examples/nodes/integration/jenkins/)
- **Jina AI** - [`jinaAi`](examples/nodes/integration/jina-ai/)
- **Jira Software** - [`jira`](examples/nodes/integration/jira-software/)
- **JWT** - [`jwt`](examples/nodes/integration/jwt/)
- **Kafka** - [`kafka`](examples/nodes/integration/kafka/)
- **Keap** - [`keap`](examples/nodes/integration/keap/)
- **Kitemaker** - [`kitemaker`](examples/nodes/integration/kitemaker/)
- **KoBoToolbox** - [`koBoToolbox`](examples/nodes/integration/kobotoolbox/)
- **Ldap** - [`ldap`](examples/nodes/integration/ldap/)
- **Lemlist** - [`lemlist`](examples/nodes/integration/lemlist/)
- **Line** - [`line`](examples/nodes/integration/line/)
- **Linear** - [`linear`](examples/nodes/integration/linear/)
- **LingvaNex** - [`lingvaNex`](examples/nodes/integration/lingvanex/)
- **LinkedIn** - [`linkedIn`](examples/nodes/integration/linkedin/)
- **LoneScale** - [`loneScale`](examples/nodes/integration/lonescale/)
- **Magento 2** - [`magento2`](examples/nodes/integration/magento-2/)
- **Mailcheck** - [`mailcheck`](examples/nodes/integration/mailcheck/)
- **Mailchimp** - [`mailchimp`](examples/nodes/integration/mailchimp/)
- **MailerLite** - [`mailerLite`](examples/nodes/integration/mailerlite/)
- **Mailgun** - [`mailgun`](examples/nodes/integration/mailgun/)
- **Mailjet** - [`mailjet`](examples/nodes/integration/mailjet/)
- **Mandrill** - [`mandrill`](examples/nodes/integration/mandrill/)
- **Markdown** - [`markdown`](examples/nodes/integration/markdown/)
- **Marketstack** - [`marketstack`](examples/nodes/integration/marketstack/)
- **Matrix** - [`matrix`](examples/nodes/integration/matrix/)
- **Mattermost** - [`mattermost`](examples/nodes/integration/mattermost/)
- **Mautic** - [`mautic`](examples/nodes/integration/mautic/)
- **Medium** - [`medium`](examples/nodes/integration/medium/)
- **MessageBird** - [`messageBird`](examples/nodes/integration/messagebird/)
- **Metabase** - [`metabase`](examples/nodes/integration/metabase/)
- **Mindee** - [`mindee`](examples/nodes/integration/mindee/)
- **MISP** - [`misp`](examples/nodes/integration/misp/)
- **Mistral AI** - [`mistralAi`](examples/nodes/integration/mistral-ai/)
- **Mocean** - [`mocean`](examples/nodes/integration/mocean/)
- **Monday.com** - [`mondayCom`](examples/nodes/integration/monday-com/)
- **MongoDB** - [`mongoDb`](examples/nodes/integration/mongodb/)
- **Monica CRM** - [`monicaCrm`](examples/nodes/integration/monica-crm/)
- **MQTT** - [`mqtt`](examples/nodes/integration/mqtt/)
- **MSG91** - [`msg91`](examples/nodes/integration/msg91/)
- **MySQL** - [`mySql`](examples/nodes/integration/mysql/)
- **n8n** - [`n8n`](examples/nodes/integration/n8n/)
- **NASA** - [`nasa`](examples/nodes/integration/nasa/)
- **Nextcloud** - [`nextCloud`](examples/nodes/integration/nextcloud/)
- **No Operation, do nothing** - [`noOp`](examples/nodes/integration/no-operation-do-nothing/)
- **NocoDB** - [`nocoDb`](examples/nodes/integration/nocodb/)
- **Notion** - [`notion`](examples/nodes/integration/notion/)
- **Npm** - [`npm`](examples/nodes/integration/npm/)
- **Odoo** - [`odoo`](examples/nodes/integration/odoo/)
- **Okta** - [`okta`](examples/nodes/integration/okta/)
- **Onfleet** - [`onfleet`](examples/nodes/integration/onfleet/)
- **OpenAI** - [`openAi`](examples/nodes/integration/openai/)
- **OpenThesaurus** - [`openThesaurus`](examples/nodes/integration/openthesaurus/)
- **OpenWeatherMap** - [`openWeatherMap`](examples/nodes/integration/openweathermap/)
- **Orbit** - [`orbit`](examples/nodes/integration/orbit/)
- **Oura** - [`oura`](examples/nodes/integration/oura/)
- **Paddle** - [`paddle`](examples/nodes/integration/paddle/)
- **PagerDuty** - [`pagerDuty`](examples/nodes/integration/pagerduty/)
- **PayPal** - [`payPal`](examples/nodes/integration/paypal/)
- **Peekalink** - [`peekalink`](examples/nodes/integration/peekalink/)
- **Perplexity** - [`perplexity`](examples/nodes/integration/perplexity/)
- **Phantombuster** - [`phantombuster`](examples/nodes/integration/phantombuster/)
- **Philips Hue** - [`philipsHue`](examples/nodes/integration/philips-hue/)
- **Pipedrive** - [`pipedrive`](examples/nodes/integration/pipedrive/)
- **Plivo** - [`plivo`](examples/nodes/integration/plivo/)
- **PostBin** - [`postBin`](examples/nodes/integration/postbin/)
- **Postgres** - [`postgres`](examples/nodes/integration/postgres/)
- **PostHog** - [`postHog`](examples/nodes/integration/posthog/)
- **ProfitWell** - [`profitWell`](examples/nodes/integration/profitwell/)
- **Pushbullet** - [`pushbullet`](examples/nodes/integration/pushbullet/)
- **Pushcut** - [`pushcut`](examples/nodes/integration/pushcut/)
- **Pushover** - [`pushover`](examples/nodes/integration/pushover/)
- **QuestDB** - [`questDb`](examples/nodes/integration/questdb/)
- **Quick Base** - [`quickbase`](examples/nodes/integration/quick-base/)
- **QuickBooks Online** - [`quickbooks`](examples/nodes/integration/quickbooks-online/)
- **QuickChart** - [`quickChart`](examples/nodes/integration/quickchart/)
- **RabbitMQ** - [`rabbitmq`](examples/nodes/integration/rabbitmq/)
- **Raindrop** - [`raindrop`](examples/nodes/integration/raindrop/)
- **Read Binary File** - [`readBinaryFile`](examples/nodes/integration/read-binary-file/)
- **Read Binary Files** - [`readBinaryFiles`](examples/nodes/integration/read-binary-files/)
- **Read PDF** - [`readPDF`](examples/nodes/integration/read-pdf/)
- **Reddit** - [`reddit`](examples/nodes/integration/reddit/)
- **Redis** - [`redis`](examples/nodes/integration/redis/)
- **Rename Keys** - [`renameKeys`](examples/nodes/integration/rename-keys/)
- **Respond With** - [`respondWith`](examples/nodes/integration/respond-with/)
- **RocketChat** - [`rocketchat`](examples/nodes/integration/rocketchat/)
- **RSS Read** - [`rssFeedRead`](examples/nodes/integration/rss-read/)
- **Rundeck** - [`rundeck`](examples/nodes/integration/rundeck/)
- **S3** - [`s3`](examples/nodes/integration/s3/)
- **Salesforce** - [`salesforce`](examples/nodes/integration/salesforce/)
- **Salesmate** - [`salesmate`](examples/nodes/integration/salesmate/)
- **SeaTable** - [`seaTable`](examples/nodes/integration/seatable/)
- **SecurityScorecard** - [`securityScorecard`](examples/nodes/integration/securityscorecard/)
- **Segment** - [`segment`](examples/nodes/integration/segment/)
- **Send Email** - [`emailSend`](examples/nodes/integration/send-email/)
- **SendGrid** - [`sendGrid`](examples/nodes/integration/sendgrid/)
- **Sendy** - [`sendy`](examples/nodes/integration/sendy/)
- **Sentry.io** - [`sentryIo`](examples/nodes/integration/sentry-io/)
- **ServiceNow** - [`serviceNow`](examples/nodes/integration/servicenow/)
- **seven** - [`sms77`](examples/nodes/integration/seven/)
- **SIGNL4** - [`signl4`](examples/nodes/integration/signl4/)
- **Simulate** - [`simulate`](examples/nodes/integration/simulate/)
- **Slack** - [`slack`](examples/nodes/integration/slack/)
- **Snowflake** - [`snowflake`](examples/nodes/integration/snowflake/)
- **Splunk** - [`splunk`](examples/nodes/integration/splunk/)
- **Spontit** - [`spontit`](examples/nodes/integration/spontit/)
- **Spreadsheet File** - [`spreadsheetFile`](examples/nodes/integration/spreadsheet-file/)
- **SSH** - [`ssh`](examples/nodes/integration/ssh/)
- **Stackby** - [`stackby`](examples/nodes/integration/stackby/)
- **Start** - [`start`](examples/nodes/integration/start/)
- **Sticky Note** - [`stickyNote`](examples/nodes/integration/sticky-note/)
- **Storyblok** - [`storyblok`](examples/nodes/integration/storyblok/)
- **Strava** - [`strava`](examples/nodes/integration/strava/)
- **Stripe** - [`stripe`](examples/nodes/integration/stripe/)
- **Supabase** - [`supabase`](examples/nodes/integration/supabase/)
- **SurveyMonkey Trigger** - [`surveyMonkeyTrigger`](examples/nodes/trigger/surveymonkey-trigger/)
- **SyncroMSP** - [`syncroMsp`](examples/nodes/integration/syncromsp/)
- **Taiga** - [`taiga`](examples/nodes/integration/taiga/)
- **Tapfiliate** - [`tapfiliate`](examples/nodes/integration/tapfiliate/)
- **TheHive** - [`theHive`](examples/nodes/integration/thehive/)
- **TheHiveProject** - [`n8n-nodes-base.thehiveproject`](examples/nodes/integration/thehiveproject/)
- **TimescaleDB** - [`timescaleDb`](examples/nodes/integration/timescaledb/)
- **Todoist** - [`todoist`](examples/nodes/integration/todoist/)
- **TOTP** - [`totp`](examples/nodes/integration/totp/)
- **TravisCI** - [`travisCi`](examples/nodes/integration/travisci/)
- **Trello** - [`trello`](examples/nodes/integration/trello/)
- **Twake** - [`twake`](examples/nodes/integration/twake/)
- **Twilio** - [`twilio`](examples/nodes/integration/twilio/)
- **Twist** - [`twist`](examples/nodes/integration/twist/)
- **Unleashed Software** - [`unleashedSoftware`](examples/nodes/integration/unleashed-software/)
- **Uplead** - [`uplead`](examples/nodes/integration/uplead/)
- **uProc** - [`uproc`](examples/nodes/integration/uproc/)
- **UptimeRobot** - [`uptimeRobot`](examples/nodes/integration/uptimerobot/)
- **urlscan.io** - [`urlScanIo`](examples/nodes/integration/urlscan-io/)
- **Vero** - [`vero`](examples/nodes/integration/vero/)
- **Vonage** - [`vonage`](examples/nodes/integration/vonage/)
- **Webflow** - [`webflow`](examples/nodes/integration/webflow/)
- **Wekan** - [`wekan`](examples/nodes/integration/wekan/)
- **WhatsApp Business Cloud** - [`whatsApp`](examples/nodes/integration/whatsapp-business-cloud/)
- **Wise** - [`wise`](examples/nodes/integration/wise/)
- **WooCommerce** - [`wooCommerce`](examples/nodes/integration/woocommerce/)
- **Wordpress** - [`wordpress`](examples/nodes/integration/wordpress/)
- **Write Binary File** - [`writeBinaryFile`](examples/nodes/integration/write-binary-file/)
- **X (Formerly Twitter)** - [`twitter`](examples/nodes/integration/x-formerly-twitter/)
- **Xero** - [`xero`](examples/nodes/integration/xero/)
- **XML** - [`xml`](examples/nodes/integration/xml/)
- **Yourls** - [`yourls`](examples/nodes/integration/yourls/)
- **Zammad** - [`zammad`](examples/nodes/integration/zammad/)
- **Zendesk** - [`zendesk`](examples/nodes/integration/zendesk/)
- **Zoho CRM** - [`zohoCrm`](examples/nodes/integration/zoho-crm/)
- **Zoom** - [`zoom`](examples/nodes/integration/zoom/)
- **Zulip** - [`zulip`](examples/nodes/integration/zulip/)

#### API Key (4)

- **APITemplate.io** - [`apiTemplateIo`](examples/nodes/integration/apitemplate-io/)
- **Facebook Graph API** - [`facebookGraphApi`](examples/nodes/integration/facebook-graph-api/)
- **One Simple API** - [`oneSimpleApi`](examples/nodes/integration/one-simple-api/)
- **Strapi** - [`strapi`](examples/nodes/integration/strapi/)

### Nodes Without Credentials (44)

These nodes work out-of-the-box without external credentials:

- **Acuity Scheduling Trigger** (`acuitySchedulingTrigger`) - [Example](examples/nodes/trigger/acuity-scheduling-trigger/)
- **Bitbucket Trigger** (`bitbucketTrigger`) - [Example](examples/nodes/trigger/bitbucket-trigger/)
- **Cal.com Trigger** (`calTrigger`) - [Example](examples/nodes/trigger/cal-com-trigger/)
- **Calendly Trigger** (`calendlyTrigger`) - [Example](examples/nodes/trigger/calendly-trigger/)
- **Clockify** (`clockify`) - [Example](examples/nodes/integration/clockify/)
- **Code** (`code`) - [Example](examples/nodes/core/code/)
- **Compare Datasets** (`compareDatasets`) - [Example](examples/nodes/integration/compare-datasets/)
- **Cron** (`cron`) - [Example](examples/nodes/integration/cron/)
- **Drift** (`drift`) - [Example](examples/nodes/integration/drift/)
- **Email Trigger (IMAP)** (`emailReadImap`) - [Example](examples/nodes/trigger/email-trigger-imap/)
- **Error Trigger** (`errorTrigger`) - [Example](examples/nodes/trigger/error-trigger/)
- **Eventbrite Trigger** (`eventbriteTrigger`) - [Example](examples/nodes/trigger/eventbrite-trigger/)
- **Figma Trigger (Beta)** (`figmaTrigger`) - [Example](examples/nodes/trigger/figma-trigger-beta/)
- **Filter** (`filter`) - [Example](examples/nodes/integration/filter/)
- **Form.io Trigger** (`formIoTrigger`) - [Example](examples/nodes/trigger/form-io-trigger/)
- **Formstack Trigger** (`formstackTrigger`) - [Example](examples/nodes/trigger/formstack-trigger/)
- **Function** (`function`) - [Example](examples/nodes/integration/function/)
- **Function Item** (`functionItem`) - [Example](examples/nodes/integration/function-item/)
- **Gotify** (`gotify`) - [Example](examples/nodes/integration/gotify/)
- **Gumroad Trigger** (`gumroadTrigger`) - [Example](examples/nodes/trigger/gumroad-trigger/)
- **If** (`if`) - [Example](examples/nodes/core/if/)
- **Jotform Trigger** (`jotFormTrigger`) - [Example](examples/nodes/trigger/jotform-trigger/)
- **Limit Wait Time** (`limitWaitTime`) - [Example](examples/nodes/integration/limit-wait-time/)
- **Local File Trigger** (`localFileTrigger`) - [Example](examples/nodes/trigger/local-file-trigger/)
- **Manual Trigger** (`manualTrigger`) - [Example](examples/nodes/trigger/manual-trigger/)
- **Merge** (`merge`) - [Example](examples/nodes/core/merge/)
- **n8n Trigger** (`n8nTrigger`) - [Example](examples/nodes/trigger/n8n-trigger/)
- **Netlify** (`netlify`) - [Example](examples/nodes/integration/netlify/)
- **Postmark Trigger** (`postmarkTrigger`) - [Example](examples/nodes/trigger/postmark-trigger/)
- **Schedule Trigger** (`scheduleTrigger`) - [Example](examples/nodes/integration/schedule-trigger/)
- **Set** (`set`) - [Example](examples/nodes/core/set/)
- **Shopify** (`shopify`) - [Example](examples/nodes/integration/shopify/)
- **Split In Batches** (`splitInBatches`) - [Example](examples/nodes/integration/split-in-batches/)
- **Spotify** (`spotify`) - [Example](examples/nodes/integration/spotify/)
- **SSE Trigger** (`sseTrigger`) - [Example](examples/nodes/trigger/sse-trigger/)
- **Stop and Error** (`stopAndError`) - [Example](examples/nodes/integration/stop-and-error/)
- **Switch** (`switch`) - [Example](examples/nodes/core/switch/)
- **Toggl Trigger** (`togglTrigger`) - [Example](examples/nodes/trigger/toggl-trigger/)
- **Typeform Trigger** (`typeformTrigger`) - [Example](examples/nodes/trigger/typeform-trigger/)
- **Wait Amount** (`amount`) - [Example](examples/nodes/integration/wait-amount/)
- **Webhook** (`webhook`) - [Example](examples/nodes/trigger/webhook/)
- **Workable Trigger** (`workableTrigger`) - [Example](examples/nodes/trigger/workable-trigger/)
- **Workflow Trigger** (`workflowTrigger`) - [Example](examples/nodes/trigger/workflow-trigger/)
- **Wufoo Trigger** (`wufooTrigger`) - [Example](examples/nodes/trigger/wufoo-trigger/)

---

## Usage Examples

### Basic Node Usage

Each node can be used in a Terraform workflow:

```hcl
# Example: Using the Code node
resource "n8n_workflow_node" "my_code" {
  name     = "Process Data"
  type     = "n8n-nodes-base.code"
  position = [250, 300]

  parameters = jsonencode({
    mode = "runOnceForAllItems"
    jsCode = "return items;"
  })
}
```

### Complete Workflow Examples

Every node has a complete, tested workflow example in:

```
examples/nodes/{category}/{node-slug}/
  ├── main.tf         # Complete workflow with the node
  ├── variables.tf    # Provider configuration
  └── README.md       # Node-specific documentation
```

### Testing Your Workflow

```bash
cd examples/nodes/core/code
terraform init
terraform validate
terraform plan
```

---

## Contributing

### Reporting Issues

If you find issues with any node:
1. Check the node's example in `examples/nodes/`
2. Review the [n8n node documentation](https://docs.n8n.io/integrations/)
3. Report issues on [GitHub](https://github.com/yourusername/terraform-provider-n8n/issues)

### Adding New Nodes

New nodes are automatically synchronized from the official n8n repository:
1. Run `make nodes` to sync latest nodes
2. Review `NODES_SYNC.md` for changes
3. Test new nodes: `make nodes/test-workflows`
4. Update documentation: `make nodes/docs`

---

## Node Registry

The complete node registry with all properties is available at:
- **Current Registry**: `data/n8n-nodes-registry.json`
- **Previous Registry**: `data/n8n-nodes-registry.previous.json`
- **Sync Report**: `NODES_SYNC.md` (generated, not committed)

### Registry Structure

```json
{
  "version": "n8n@1.x.x",
  "last_sync": "2024-01-01T00:00:00Z",
  "nodes": [
    {
      "name": "Code",
      "type": "n8n-nodes-base.code",
      "category": "core",
      "latest_version": 2,
      "description": "Execute custom JavaScript code",
      "inputs": ["main"],
      "outputs": ["main"],
      "file": "dist/nodes/Code/Code.node.js"
    }
  ]
}
```

---

## Maintenance

This documentation is auto-generated from the node registry. To regenerate:

```bash
make nodes/docs
```

**Last Generated**: 2025-11-17T15:23:04.863Z
