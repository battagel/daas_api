# DaaS - Documentation-as-a-Service API

Have you every struggled to remember your company's jargon? DaaS might be the solution!

The DaaS suite is a collection of software extensions based around a central
API. This API is the single source of truth for company specific terms. Each
extension can be installed provide explanations no matter where you work.

## Usage

### Overview

Ensure that the DaaS API is running within a secure network. For every
applicable application install a DaaS extension. Hover over any terms you are
unfamiliar with and a helpful explanation should pop up. 

See extensions
- [Chromium Browsers](https://github.com/battagel/daas_chrome)
- [Visual Studio Code](https://github.com/battagel/daas_vscode)

## Development

### Backlog

- Seperate explanations into separate table
- Add JSONB for faster reads in SQL lite
- Unittest and mocking
- Dockerise
- Hosting?
  - HTTPS
- Author field
- Investigate speed of getAllPhrases
- Add ping for verification of if API is reachable

