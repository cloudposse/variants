
<!-- markdownlint-disable -->
# atmos<a href="https://cpco.io/homepage"><img align="right" src="https://cloudposse.com/logo-300x69.svg" width="150" /></a> [![Latest Release](https://img.shields.io/github/release/cloudposse/atmos.svg)](https://github.com/cloudposse/atmos/releases/latest) [![Slack Community](https://slack.cloudposse.com/badge.svg)](https://slack.cloudposse.com)
<!-- markdownlint-restore -->


<!--




  ** DO NOT EDIT THIS FILE
  **
  ** This file was automatically generated by the `build-harness`.
  ** 1) Make all changes to `README.yaml`
  ** 2) Run `make init` (you only need to do this once)
  ** 3) Run`make readme` to rebuild this file.
  **
  ** (We maintain HUNDREDS of open source projects. This is how we maintain our sanity.)
  **





-->

![atmos](docs/img/atmos-logo-128.svg)

`atmos` - Universal Tool for DevOps and Cloud Automation.

---
> [!NOTE]
> This project is part of Cloud Posse's comprehensive ["SweetOps"](https://cpco.io/sweetops) approach towards DevOps.
>
> It's 100% Open Source and licensed under the [APACHE2](LICENSE).
>

[![README Header][readme_header_img]][readme_header_link]



## Screenshots


![Demo](docs/demo.gif)
*<br/>Example of running atmos to describe infrastructure.*


## Introduction


[Atmos](https://atmos.tools) is a workflow automation tool. Atmos is lets you take a complex workflow and turn it into a parameterized CLI. 
You can then run the CLI anywhere, such as locally or in CI/CD.

The Atmos project consists of a command-line tool, a `Go` library, and even terraform provider.  It provides numerous conventions to help you provision, manage, and orchestrate workflows across various toolchains. You can 
even access the configurations natively from within terraform using our [`terraform-provider-utils`](https://github.com/cloudposse/terraform-provider-utils/).

Cloud Posse uses this tool extensively for automating cloud infrastructure with [Terraform](https://hashicorp.com/) and [Kubernetes](https://kubernetes.io/), but it can be used to automate any complex workflow.

## Common Use-cases

  `atmos` includes workflows for dealing with:

   - Provision large, multi-account[Terraform](https://www.terraform.io/) environments
   - Deploy [helm](https://helm.sh/) [charts](https://helm.sh/docs/topics/charts/) to Kubernetes clusters with [helmfile](https://github.com/roboll/helmfile)
   - Execute [helm](https://helm.sh/) commands on Kubernetes clusters
   - Executing [kubectl](https://kubernetes.io/docs/reference/kubectl/overview/) commands on Kubernetes clusters

   Here are some other ways you can extend atmos:
   - Running [AWS CDK](https://aws.amazon.com/cdk/) constructs to define cloud resources
   - Executing commands for the [serverless](https://www.serverless.com/) framework
   - Executing shell commands
   - Combining various commands into workflows to execute many commands sequentially in just one step
   - ... and many more

In essence, it's a tool that orchestrates the other CLI tools in a consistent and self-explaining manner.

It's a superset of all other tools and task runners (e.g. `make`, `terragrunt`, `terraform`, `aws` cli, `gcloud`, etc.)
and is intended to be used to tie everything together, so you can provide a simple CLI interface for your developers.

Moreover, `atmos` is not only a command-line interface for managing clouds and clusters. It provides many useful patterns and best practices, such as:

- Enforces a project structure convention, so everybody knows where to find things.
- Provides clear separation of configuration from code, so the same code is easily deployed to different regions, environments and stages)
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

### 🐛 Bug Reports & Feature Requests

Please use the [issue tracker](https://github.com/cloudposse/atmos/issues) to report any bugs or file feature requests.

### 💻 Developing

If you are interested in being a contributor and want to get involved in developing this project or [help out](https://cpco.io/help-out) with Cloud Posse's other projects, we would love to hear from you! Shoot us an [email][email].

In general, PRs are welcome. We follow the typical "fork-and-pull" Git workflow.

 1. **Fork** the repo on GitHub
 2. **Clone** the project to your own machine
 3. **Commit** changes to your own branch
 4. **Push** your work back up to your fork
 5. Submit a **Pull Request** so that we can review your changes

**NOTE:** Be sure to merge the latest changes from "upstream" before making a pull request!

### 🌎 Slack Community

Join our [Open Source Community][slack] on Slack. It's **FREE** for everyone! Our "SweetOps" community is where you get to talk with others who share a similar vision for how to rollout and manage infrastructure. This is the best place to talk shop, ask questions, solicit feedback, and work together as a community to build totally *sweet* infrastructure.

### 📰 Newsletter

Sign up for [our newsletter][newsletter] that covers everything on our technology radar.  Receive updates on what we're up to on GitHub as well as awesome new projects we discover.

### 📆 Office Hours <img src="https://img.cloudposse.com/fit-in/200x200/https://cloudposse.com/wp-content/uploads/2019/08/Powered-by-Zoom.png" align="right" />

[Join us every Wednesday via Zoom][office_hours] for our weekly "Lunch & Learn" sessions. It's **FREE** for everyone!

## About 

This project is maintained and funded by [Cloud Posse, LLC][website]. 
<a href="https://cpco.io/homepage"><img src="https://cloudposse.com/logo-300x69.svg" align="right" /></a>

We are a [**DevOps Accelerator**][commercial_support]. We'll help you build your cloud infrastructure from the ground up so you can own it. Then we'll show you how to operate it and stick around for as long as you need us.

[![Learn More](https://img.shields.io/badge/learn%20more-success.svg?style=for-the-badge)][commercial_support]

Work directly with our team of DevOps experts via email, slack, and video conferencing.

We deliver 10x the value for a fraction of the cost of a full-time engineer. Our track record is not even funny. If you want things done right and you need it done FAST, then we're your best bet.

- **Reference Architecture.** You'll get everything you need from the ground up built using 100% infrastructure as code.
- **Release Engineering.** You'll have end-to-end CI/CD with unlimited staging environments.
- **Site Reliability Engineering.** You'll have total visibility into your apps and microservices.
- **Security Baseline.** You'll have built-in governance with accountability and audit logs for all changes.
- **GitOps.** You'll be able to operate your infrastructure via Pull Requests.
- **Training.** You'll receive hands-on training so your team can operate what we build.
- **Questions.** You'll have a direct line of communication between our teams via a Shared Slack channel.
- **Troubleshooting.** You'll get help to triage when things aren't working.
- **Code Reviews.** You'll receive constructive feedback on Pull Requests.
- **Bug Fixes.** We'll rapidly work with you to fix any bugs in our projects.

[![README Commercial Support][readme_commercial_support_img]][readme_commercial_support_link]
## License

[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

See [LICENSE](LICENSE) for full details.

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

## Trademarks

All other trademarks referenced herein are the property of their respective owners.
---
Copyright © 2017-2024 [Cloud Posse, LLC](https://cpco.io/copyright)
[![README Footer][readme_footer_img]][readme_footer_link]
[![Beacon][beacon]][website]
<!-- markdownlint-disable -->
  [logo]: https://cloudposse.com/logo-300x69.svg
  [docs]: https://cpco.io/docs?utm_source=github&utm_medium=readme&utm_campaign=cloudposse/atmos&utm_content=docs
  [website]: https://cpco.io/homepage?utm_source=github&utm_medium=readme&utm_campaign=cloudposse/atmos&utm_content=website
  [github]: https://cpco.io/github?utm_source=github&utm_medium=readme&utm_campaign=cloudposse/atmos&utm_content=github
  [jobs]: https://cpco.io/jobs?utm_source=github&utm_medium=readme&utm_campaign=cloudposse/atmos&utm_content=jobs
  [hire]: https://cpco.io/hire?utm_source=github&utm_medium=readme&utm_campaign=cloudposse/atmos&utm_content=hire
  [slack]: https://cpco.io/slack?utm_source=github&utm_medium=readme&utm_campaign=cloudposse/atmos&utm_content=slack
  [twitter]: https://cpco.io/twitter?utm_source=github&utm_medium=readme&utm_campaign=cloudposse/atmos&utm_content=twitter
  [office_hours]: https://cloudposse.com/office-hours?utm_source=github&utm_medium=readme&utm_campaign=cloudposse/atmos&utm_content=office_hours
  [newsletter]: https://cpco.io/newsletter?utm_source=github&utm_medium=readme&utm_campaign=cloudposse/atmos&utm_content=newsletter
  [email]: https://cpco.io/email?utm_source=github&utm_medium=readme&utm_campaign=cloudposse/atmos&utm_content=email
  [commercial_support]: https://cpco.io/commercial-support?utm_source=github&utm_medium=readme&utm_campaign=cloudposse/atmos&utm_content=commercial_support
  [we_love_open_source]: https://cpco.io/we-love-open-source?utm_source=github&utm_medium=readme&utm_campaign=cloudposse/atmos&utm_content=we_love_open_source
  [terraform_modules]: https://cpco.io/terraform-modules?utm_source=github&utm_medium=readme&utm_campaign=cloudposse/atmos&utm_content=terraform_modules
  [readme_header_img]: https://cloudposse.com/readme/header/img
  [readme_header_link]: https://cloudposse.com/readme/header/link?utm_source=github&utm_medium=readme&utm_campaign=cloudposse/atmos&utm_content=readme_header_link
  [readme_footer_img]: https://cloudposse.com/readme/footer/img
  [readme_footer_link]: https://cloudposse.com/readme/footer/link?utm_source=github&utm_medium=readme&utm_campaign=cloudposse/atmos&utm_content=readme_footer_link
  [readme_commercial_support_img]: https://cloudposse.com/readme/commercial-support/img
  [readme_commercial_support_link]: https://cloudposse.com/readme/commercial-support/link?utm_source=github&utm_medium=readme&utm_campaign=cloudposse/atmos&utm_content=readme_commercial_support_link
  [beacon]: https://ga-beacon.cloudposse.com/UA-76589703-4/cloudposse/atmos?pixel&cs=github&cm=readme&an=atmos
<!-- markdownlint-restore -->
