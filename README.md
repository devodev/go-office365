# go-office365
`go-office365` provides both a client library as well as a CLI tool to interact with the `Microsoft Office365 Management Activity API`.

# Table of Contents

- [Overview](#overview)
- [Configuration file](#configuration-file)
- [Roadmap](#roadmap)
- [License](#license)

# Overview
`go-office365` is a client library for the `Microsoft Office365 Management Activity API` written in Go. It follows the Microsoft API Reference available [here](https://docs.microsoft.com/en-us/office/office-365-management-api/office-365-management-activity-api-reference). It is inspired by the [go-github](https://github.com/google/go-github) project.

`go-office365` also provides a CLI application.</br>

# Configuration file
The default configuration file name is `.go-office365.yaml`. It is expected to be found in the current working directory.

An alternative location can be provided using the `--config` flag.

Here is the format of the configuration file. These attributes can be found in `Azure Active Directory`, under: `Installed apps`.

```
// .go-office365.yaml
---
Credentials:
  ClientID: 00000000-0000-0000-0000-000000000000
  ClientSecret: 00000000000000000000000000000000
  TenantID: 00000000-0000-0000-0000-000000000000
  TenantDomain: some-company.onmicrosoft.com
```

# Roadmap
Overall
- Better README
- More documentation
- More tests

# License
go-office365 is released under the MIT license. See [LICENSE.txt](LICENSE.txt)
