# Kopilot 🧑‍✈️: Your AI Kubernetes Expert

![GitHub](https://img.shields.io/github/license/knight42/kopilot)
![GitHub last commit](https://img.shields.io/github/last-commit/knight42/kopilot)

## Highlight

* Diagnose any unhealthy workload in your cluster and tell you what might be the cause
  * ![](https://user-images.githubusercontent.com/4237254/226960414-a343b624-b95f-479c-840f-10fb9dc5de05.gif)
* Audit Kubernetes resource and find the security misconfigurations
  * ![](https://user-images.githubusercontent.com/4237254/226959542-57193653-0afe-4a8b-bee6-ab96bacdfa83.gif)

## Installation

|         Distribution         |                         Command / Link                          |
|:----------------------------:|:---------------------------------------------------------------:|
|            macOS             |               `brew install knight42/tap/kopilot`               |
| Pre-built binaries for Linux | [GitHub releases](https://github.com/knight42/kopilot/releases) |

## Usage

Currently, you need to set two ENVs to run Kopilot:
* Set `KOPILOT_TOKEN` to specify your token.
* Set `KOPILOT_LANG` to specify the language, defaults to `English`. Valid options are `Chinese`, `French`, `Spain`, etc.
* `KOPILOT_TOKEN_TYPE` will be available soon to let you specify AI services other than ChatGPT. Please stay tuned.

### Diagnose

```bash
# Diagnose a CrashLoopBackOff pod
kopilot diagnose pod my-pod
```

### Audit

```bash
# Audit a deployment named nginx
kopilot audit deploy nginx
```
