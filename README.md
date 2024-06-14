

<!-- markdownlint-disable -->
<a href="https://cpco.io/homepage"><img src="https://github.com/cloudposse/atmos/blob/main/.github/banner.png?raw=true" alt="Project Banner"/></a><br/>
    <p align="right">
<a href="https://github.com/cloudposse/atmos/releases/latest"><img src="https://img.shields.io/github/release/cloudposse/atmos.svg?style=for-the-badge" alt="Latest Release"/></a><a href="https://github.com/cloudposse/atmos/commits/main/"><img src="https://img.shields.io/github/last-commit/cloudposse/atmos/main?style=for-the-badge" alt="Last Updated"/></a><a href="https://github.com/cloudposse/atmos/actions/workflows/test.yml"><img src="https://img.shields.io/github/actions/workflow/status/cloudposse/atmos/test.yml?style=for-the-badge" alt="Tests"/></a><a href="https://slack.cloudposse.com"><img src="https://slack.cloudposse.com/for-the-badge.svg" alt="Slack Community"/></a></p>
<!-- markdownlint-restore -->

<!--




  ** DO NOT EDIT THIS FILE
  **
  ** This file was automatically generated by the `cloudposse/build-harness`.
  ** 1) Make all changes to `README.yaml`
  ** 2) Run `make init` (you only need to do this once)
  ** 3) Run`make readme` to rebuild this file.
  **
  ** (We maintain HUNDREDS of open source projects. This is how we maintain our sanity.)
  **





-->


Atmos is a framework that simplifies complex cloud architectures and DevOps workflows into [intuitive CLI commands](https://atmos.tools/category/cli).
Its strength in managing [DRY configurations at scale](https://atmos.tools/core-concepts/) for Terraform and is supported by robust
[design patterns](https://atmos.tools/design-patterns/), comprehensive [documentation](https://atmos.tools/), and a
[passionate community](https://slack.cloudposse.com/), making it a versatile [tool for both startups and enterprises](https://cloudposse.com/).
Atmos is extensible to accommodate any tooling, including enterprise-scale Terraform, and includes custom
[policy controls](https://atmos.tools/core-concepts/validate), [vendoring](https://atmos.tools/core-concepts/vendor/),
and [GitOps capabilities](https://atmos.tools/integrations/github-actions) out of the box. Everything is open source and free.


> [!TIP]
> ### You can try out `atmos` directly in your browser using GitHub Codespaces
>
> [![Open in GitHub Codespaces](https://github.com/codespaces/badge.svg)](https://github.com/codespaces/new?hide_repo_select=true&ref=reorg&repo=cloudposse/atmos&skip_quickstart=true) 
> <i>Already start one? Find it [here](https://github.com/codespaces).</i>
>

## Screenshots

<img src="docs/demo.gif" alt="Demo" />*<br/>Example of running atmos to describe infrastructure.*




## Introduction


[Atmos](https://atmos.tools) centralizes the DevOps chain and cloud automation/orchestration into a robust command-line tool,
streamlining environments and workflows into straightforward CLI commands. Leveraging advanced hierarchical configurations,
it efficiently orchestrates both local and CI/CD pipeline tasks, optimizing infrastructure management for engineers and cloud 
architects alike. You can then run the CLI anywhere, such as locally or in CI/CD.

The Atmos project consists of a command-line tool, a `Go` library, and even a terraform provider.  It provides numerous
[conventions](https://atmos.tools/design-patterns/) to help you provision, manage, and orchestrate workflows across various toolchains.
You can even access the configurations natively from within terraform using our [`terraform-provider-utils`](https://github.com/cloudposse/terraform-provider-utils/).

[Cloud Posse](https://cloudposse.com/) uses this tool extensively for automating cloud infrastructure with
[Terraform](https://www.hashicorp.com/products/terraform) and [Kubernetes](https://kubernetes.io/), but it can be used to automate any complex workflow.

> [!TIP]
> ### Did you know?
>
> By leveraging Atmos in conjunction with Cloud Posse's [*expertise in AWS*](https://cloudposse.com),
> [*terraform blueprints*](https://cloudposse.com/services/), and our [*knowledgeable community*](https://slack.cloudposse.com), teams can achieve
> operational mastery and innovation faster, transforming their infrastructure management practices into a competitive advantage.

## Core Features

Atmos streamlines Terraform orchestration, environment, and configuration management, offering developers and DevOps a set of
powerful tools to tackle deployment challenges. Designed to be cloud agnostic, it enables you to operate consistently across
various cloud platforms. These features boost efficiency, clarity, and control across various environments, making it an
indispensable asset for managing complex infrastructures with confidence.

- [**Terminal UI**](https://atmos.tools/cli) Polished interface for easier interaction with Terraform, workflows, and commands.
- [**Native Terraform Support:**](https://atmos.tools/cli/commands/terraform/usage) Orchestration, backend generation, varfile generation, ensuring compatibility with vanilla Terraform.
- [**Stacks:**](https://atmos.tools/core-concepts/stacks) Powerful abstraction layer defined in YAML for orchestrating and deploying components.
- [**Components:**](https://atmos.tools/core-concepts/components) A generic abstraction for deployable units, such as Terraform "root" modules.
- [**Vendoring:**](https://atmos.tools/core-concepts/vendor) Pulls dependencies from remote sources, supporting immutable infrastructure practices.
- [**Custom Commands:**](https://atmos.tools/core-concepts/custom-commands) Extends Atmos's functionality, allowing integration of any command with stack configurations.
- [**Workflow Orchestration:**](https://atmos.tools/core-concepts/workflows) Comprehensive support for managing the lifecycle of cloud infrastructure from initiation to maintenance.

See [all features of Atmos](https://atmos.tools/features).

## Use Cases

Atmos has consistently demonstrated its effectiveness in addressing these key use-cases, showcasing its adaptability and
strength in the cloud infrastructure and DevOps domains:

- **Managing Large Multi-Account Cloud Environments:** Suitable for organizations using multiple cloud accounts to separate different
  projects or stages of development.
- **Cross-Platform Cloud Architectures:** Ideal for businesses that need to manage configuration of services across AWS, GCP, Azure, etc., to
  build a cohesive system.
- **Multi-Tenant Systems for SaaS:** Perfect for SaaS companies looking to host multiple customers within a unified infrastructure.
  Simply define a baseline tenant configuration once, and then seamlessly onboard new tenants by reusing this baseline through pure
  configuration, bypassing the need for further code development.
- **Efficient Multi-Region Deployments:** Atmos facilitates streamlined multi-region deployments by enabling businesses to define baseline
  configurations with [stacks](https://atmos.tools/core-concepts/stacks/) and extend them across regions with DRY principles through
  [imports](https://atmos.tools/core-concepts/stacks/imports) and [inheritance](https://atmos.tools/core-concepts/stacks/inheritance).
- **Compliant Infrastructure for Regulated Industries:** Atmos empowers DevOps and SecOps teams to create vetted configurations that comply
  with SOC2, HIPAA, HITRUST, PCI, and other regulatory standards. These configurations can then be efficiently shared and reused across the
  organization via [service catalogs](https://atmos.tools/core-concepts/stacks/catalogs), [component libraries](https://atmos.tools/core-concepts/components/library),
  [vendoring](https://atmos.tools/core-concepts/vendor), and [OPA policies](https://atmos.tools/core-concepts/validate/opa),
  simplifying the process of achieving and maintaining rigorous compliance.
- **Empowering Teams with Self-Service Infrastructure:** Allows teams to manage their infrastructure needs independently, using
  predefined templates and policies.
- **Streamlining Deployment with Service Catalogs, Landing Zones, and Blueprints:** Provides ready-to-use templates and guidelines for
  setting up cloud environments quickly and consistently.

> [!TIP]
> Don't see your use-case listed? Ask us in the [`#atmos`](https://slack.cloudposse.com) Slack channel,
> or [join us for "Office Hours"](https://cloudposse.com/office-hours/) every week.


Moreover, `atmos` is not only a command-line interface for managing clouds and clusters. It provides many useful patterns
and best practices, such as:
- Enforces a project structure convention, so everybody knows where to find things.
- Provides clear separation of configuration from code, so the same code is easily deployed to different regions, environments and stages
- It can be extended to include new features, commands, and workflows
- The commands have a clean, consistent and easy to understand syntax
- The CLI code is modular and self-documenting

## Documentation

Find all documentation at: [atmos.tools](https://atmos.tools)













## ✨ Contributing

This project is under active development, and we encourage contributions from our community.



Many thanks to our outstanding contributors:

<a href="https://github.com/cloudposse/atmos/graphs/contributors">
  <img src="https://contrib.rocks/image?repo=cloudposse/atmos&max=24" />
</a>

For 🐛 bug reports & feature requests, please use the [issue tracker](https://github.com/cloudposse/atmos/issues).

In general, PRs are welcome. We follow the typical "fork-and-pull" Git workflow.
 1. Review our [Code of Conduct](https://github.com/cloudposse/atmos/?tab=coc-ov-file#code-of-conduct) and [Contributor Guidelines](https://github.com/cloudposse/.github/blob/main/CONTRIBUTING.md).
 2. **Fork** the repo on GitHub
 3. **Clone** the project to your own machine
 4. **Commit** changes to your own branch
 5. **Push** your work back up to your fork
 6. Submit a **Pull Request** so that we can review your changes

**NOTE:** Be sure to merge the latest changes from "upstream" before making a pull request!

### 🌎 Slack Community

Join our [Open Source Community](https://cpco.io/slack?utm_source=github&utm_medium=readme&utm_campaign=cloudposse/atmos&utm_content=slack) on Slack. It's **FREE** for everyone! Our "SweetOps" community is where you get to talk with others who share a similar vision for how to rollout and manage infrastructure. This is the best place to talk shop, ask questions, solicit feedback, and work together as a community to build totally *sweet* infrastructure.

### 📰 Newsletter

Sign up for [our newsletter](https://cpco.io/newsletter?utm_source=github&utm_medium=readme&utm_campaign=cloudposse/atmos&utm_content=newsletter) and join 3,000+ DevOps engineers, CTOs, and founders who get insider access to the latest DevOps trends, so you can always stay in the know.
Dropped straight into your Inbox every week — and usually a 5-minute read.

### 📆 Office Hours <a href="https://cloudposse.com/office-hours?utm_source=github&utm_medium=readme&utm_campaign=cloudposse/atmos&utm_content=office_hours"><img src="https://img.cloudposse.com/fit-in/200x200/https://cloudposse.com/wp-content/uploads/2019/08/Powered-by-Zoom.png" align="right" /></a>

[Join us every Wednesday via Zoom](https://cloudposse.com/office-hours?utm_source=github&utm_medium=readme&utm_campaign=cloudposse/atmos&utm_content=office_hours) for your weekly dose of insider DevOps trends, AWS news and Terraform insights, all sourced from our SweetOps community, plus a _live Q&A_ that you can’t find anywhere else.
It's **FREE** for everyone!
## License

<a href="https://opensource.org/licenses/Apache-2.0"><img src="https://img.shields.io/badge/License-Apache%202.0-blue.svg?style=for-the-badge" alt="License"></a>

<details>
<summary>Preamble to the Apache License, Version 2.0</summary>
<br/>
<br/>

Complete license is available in the [`LICENSE`](LICENSE) file.

```text
Licensed to the Apache Software Foundation (ASF) under one
or more contributor license agreements.  See the NOTICE file
distributed with this work for additional information
regarding copyright ownership.  The ASF licenses this file
to you under the Apache License, Version 2.0 (the
"License"); you may not use this file except in compliance
with the License.  You may obtain a copy of the License at

  https://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing,
software distributed under the License is distributed on an
"AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
KIND, either express or implied.  See the License for the
specific language governing permissions and limitations
under the License.
```
</details>

## Trademarks

All other trademarks referenced herein are the property of their respective owners.


---
Copyright © 2017-2024 [Cloud Posse, LLC](https://cpco.io/copyright)


<a href="https://cloudposse.com/readme/footer/link?utm_source=github&utm_medium=readme&utm_campaign=cloudposse/atmos&utm_content=readme_footer_link"><img alt="README footer" src="https://cloudposse.com/readme/footer/img"/></a>

<img alt="Beacon" width="0" src="https://ga-beacon.cloudposse.com/UA-76589703-4/cloudposse/atmos?pixel&cs=github&cm=readme&an=atmos"/>
