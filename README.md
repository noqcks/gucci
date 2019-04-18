# Gucci

A simple CLI templating tool written in golang.

[![GitHub version](https://badge.fury.io/gh/noqcks%2Fgucci.svg)](https://badge.fury.io/gh/noqcks%2Fgucci)
[![License](https://img.shields.io/github/license/noqcks/gucci.svg)](https://github.com/noqcks/gucci/blob/master/LICENSE)
[![Build Status](https://travis-ci.org/noqcks/gucci.svg?branch=master)](https://travis-ci.org/noqcks/gucci)

## Installation

If you have `go` installed:

```
$ go get github.com/noqcks/gucci
```

Or you can just download the binary and move it into your `PATH`:

```
VERSION=0.1.0
wget -q https://github.com/noqcks/gucci/releases/download/v${VERSION}/gucci-v${VERSION}-darwin-amd64
chmod +x gucci-v${VERSION}-darwin-amd64
mv gucci-v${VERSION}-darwin-amd64 /usr/local/bin/gucci
```

## Use

### Locating Templates

`gucci` can locate a template in multiple ways.

#### File

Pass the template file path as the first argument:

```
$ gucci template.tpl > template.out
```

#### Stdin

Supply the template through standard input:

```bash
$ gucci
Start typing stuff {{ print "here" }}
^d
Start typing stuff here
```

Via piping:

```bash
$ echo '{{ html "<escape-me/>" }}' | gucci
```

### Supplying Variable Inputs

`gucci` can receive variables for use in templates in the following ways (in order of lowest to highest precedence):
* A JSON or YAML file
* Environment variables
* Variable command options

#### Variables File

Given an example variables file:
```yaml
# vars.yaml
hosts:
  - name: bastion
  - name: app
```

Pass it into `gucci` with `-f` or `--vars-file`:
```bash
$ gucci -f vars.yaml template.tpl
```

#### Environment Variables

Here, `MY_HOST` is available to the template:
```bash
$ export MY_HOST=localhost
$ gucci template.tpl
```

#### Variable Options

Pass variable options into `gucci` with `-s` or `--set-var`, which can be repeated:
```bash
$ gucci -s foo.bar=baz template.tpl
```

Variable option keys are split on the `.` character, and nested such that
the above example would equate to the following yaml variable input:

```yaml
foo:
  bar: baz
```

## Templating

#### GoLang Functions

All of the existing [golang templating functions](https://golang.org/pkg/text/template/#hdr-Functions) are available for use.

#### Sprig Functions

gucci ships with the [sprig templating functions library](http://masterminds.github.io/sprig/) offering a wide variety of template helpers. Sprig's `env` and `expandenv` functions are disabled in favor of gucci's own environment variable parsing (see below).

#### Built In Functions

Furthermore, this tool also includes custom functions:

- `shell`: For arbitrary shell commands

   ```
   {{ shell "echo hello world" }}
   ```

   Produces:

   ```
   hello world
   ```


- `toYaml`: Print items in YAML format

   ```
   {{ $myList := list "a" "b" "c" }}
   {{ toYaml $myList }}
   ```

   Produces:

   ```
   - a
   - b
   - c
   ```

### Example

**NOTE**: gucci reads and makes available all environment variables.

For example a var $LOCALHOST = 127.0.0.1

gucci template.tpl > template.conf


```
# template.tpl
{{ .LOCALHOST }}
```

`gucci template.tpl > template.conf` -->

```
# template.conf
127.0.0.1
```

simple enough!

For an iteration example, you have $BACKENDS=server1.com,server2.com

```
# template.tpl
{{ range split .BACKENDS "," }}
server {{ . }}
{{ end }}
```

`gucci template.tpl > template.conf` -->


```
# template.conf
server server1.com
server server2.com
```

## Testing

Setup:

```bash
go get github.com/noqcks/gucci
go get github.com/onsi/ginkgo/ginkgo
go get github.com/onsi/gomega
```

Run tests:

```bash
ginkgo ./...
```
