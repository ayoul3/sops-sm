# Sops-SM
![Build](https://github.com/ayoul3/sops-sm/workflows/Go/badge.svg)

Sops-SM is a leightweight fork of [Sops by Mozilla](https://github.com/mozilla/sops) that handles AWS SecretsManager and AWS Parameter Store.

Given a file with key values containing SecertsManager IDs, its spits out a plain text:
```yaml
values:
  - key1: arn:aws:ssm:eu-west-1:123456789123:parameter/key1
  - nested:
    - key2: arn:aws:secretsmanager:eu-west-3:123456789123:secret/key2
extra: top
```
```bash
$ sops-sm decrypt test.yaml --overwrite
```
Output:
```yaml
values:
  - key1: secretvalue
  - nested:
    - key2: secretvalue
extra: top
```
It currently handles JSON and YAML formats coupled with AWS SecretsManager and ParameterStore. It auto-discovers the secret's provider, region and file format.

## Which problem does it solve?

Sops by Mozilla is a fantastic tool that has many applications, but the secret providers it supports lack one central securtiy requirement: **Secret centralization**.

Let's say you use classic Sops to encrypt your Helm charts. You have 50 micro-services, each with their own sops-encrypted `secrets.yaml` file containg dozens of encrypted secrets. Try rotating that database password they all share....

Sops-SM relies on SecretsManager and ParameterStore, so as long as you use that same secret ID in all 50 charts, you simply need to rotate that password in two locations: SecretsManager/ParameterStore and the DB itself.

# Download

## Stable release
Binaries and packages of the latest stable release are available at https://github.com/ayoul3/sops-sm/releases.

## Build from source
1. [Prepare a Go environment](https://golang.org/dl/).
2. clone and build
```zsh
$ git clone https://github.com/ayoul3/sops-sm
$ make build
$ ./sops-sm version
sops-sm version 0.1
```

# Usage
## Decrypt a file
```zsh
$ sops-sm decrypt test.yaml
```
This will produce a `test.plain.yaml` with the decoded secrets. You can overwrite the input file with the `-o` or `--overwrite` option.

sops-sm builds a local cache of decrypted secrets to avoid fecthing them multiple times. That being said, for files with more than 20 secret you can activate the multi-threaded feature:
```zsh
$ sops-sm decrypt test.yaml --threads 10
```
This will start 10 go routines to fetch secrets concurrently. Be wary of AWS's throttling.

## How should secrets be formatted?
Secrets stored in SecretsManager or ParameterStore can be of two formats:
* Simple strings
* Flat JSON

You can specify which JSON key to fetch by adding it after the character `@`.

*Example*: let's create a secret with a JSON structure:
```
aws ssm put-parameter --name complex-secret --value '{"user":"name","pass":"secret"}' --type SecureString
```
Our encoded file would look like:
```yaml
values:
  - name: arn:aws:ssm:eu-west-1:123456789123:parameter/complex-secret@user
  - password: arn:aws:ssm:eu-west-1:123456789123:parameter/complex-secret@pass
```

## Encode the file again
Each time a file is decoded by sops-sm it generates a cache file containig the path of each secret key along with its original value:
```zsh
$ sops-sm decrypt test.yaml --overwrite
```
This command will generate test.yaml.cache. You edit the test.yaml file, or test.yaml.cache then produce the encoded version again:
```zsh
$ sops-sm encrypt test.yaml --overwrite
```

## Prepare the input file
sops-sm cannot currently create secrets in SecretsManager or Parameter. So you cannot start with a plain YAML file and end up with file containing secret IDs.

I may work on it in the future, but there are some challenges to be solved first. How do we avoid secret duplication across multiple files? what naming convention to settle on? how much flexibility should we offer? and so on.

## Detailed Usage
```
Sops-SM decrypts a yaml or json file that contain references to secrets
stored in AWS SecretsManager or Parameter Store

Usage:
  sops-sm [command]

Available Commands:
  decrypt     Decrypt input file
  encrypt     Encrypt input file - requires .cache file generated from the decryption phase
  help        Help about any command

Flags:
  -h, --help        help for sops-sm
  -o, --overwrite   Overwrite input file
  -v, --verbose     Show info message

Use "sops-sm [command] --help" for more information about a command.
```

## License
Mozilla Public License Version 2.0

## Author
@ayoul3

## Credit
Obviously the original sops project for solving much if not all the yaml and json parsing part.