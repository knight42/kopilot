# Kopilot 🧑‍✈️: Your AI Kubernetes Expert

![GitHub](https://img.shields.io/github/license/knight42/kopilot)
![GitHub last commit](https://img.shields.io/github/last-commit/knight42/kopilot)

## Highlight

* Diagnose any unhealthy workload in your cluster and tell you what might be the cause
  * ![](https://user-images.githubusercontent.com/15977536/227199412-0a4b4967-a456-4b81-a1e8-47891bf08af2.gif)
* Audit Kubernetes resource and find the security misconfigurations
  * ![](https://user-images.githubusercontent.com/15977536/227199946-04ea075a-787b-4871-aa25-4a8180c84f84.gif)

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
