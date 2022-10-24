# Gucci

A simple CLI templating tool written in golang.

[![GitHub version](https://badge.fury.io/gh/noqcks%2Fgucci.svg)](https://badge.fury.io/gh/noqcks%2Fgucci)
[![Build Status](https://travis-ci.org/noqcks/gucci.svg?branch=master)](https://travis-ci.org/noqcks/gucci)

## Installation

If you have `go` installed:

```
$ go get github.com/noqcks/gucci
```

Or you can just download the binary and move it into your `PATH`:

```
VERSION=1.6.5
wget -q "https://github.com/noqcks/gucci/releases/download/${VERSION}/gucci-v${VERSION}-darwin-amd64"
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

- A JSON or YAML file
- Environment variables
- Variable command options

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

### Options

Existing [golang templating options](https://golang.org/pkg/text/template/#Template.Option) can be used for templating.

If no option is specified, the `missingkey=error` option will be used (execution stops immediately with an error if a
key used in the template is not present in the supplied values).

One might want a different value for `missingkey` when using conditionals and having keys that won't be
used at all.

For instance, given the following template, containing two docker-compose services `service1` and `service2`:

```tpl
# template.tpl
version: "3.8"

services:
{{- if .service1 }}
  service1:
    image: {{ .service1.image }}
    restart: "always"
    ports: {{ toYaml .service1.ports | nindent 6}}
{{- end }}
{{- if .service2 }}
  service2:
    image: {{ .service2.image }}
    restart: "unless-stopped"
    ports: {{ toYaml .service2.ports | nindent 6}}
{{- end }}
```

And imagine a scenario where whe only need `service2`. By using the following values file:

```yaml
# values.yaml
service2:
  image: "myservice:latest"
  ports:
    - "80"
    - "443"
```

And using a different `missingkey=error`, we can actually get the desired result without having to define the values
for `service1`:

```shell
$ gucci -o missingkey=zero -f values.yaml  template.tpl
version: "3.8"

services:
  service2:
    image: myservice:latest
    restart: "unless-stopped"
    ports:
      - "80"
      - "443"
```

### GoLang Functions

All of the existing [golang templating functions](https://golang.org/pkg/text/template/#hdr-Functions) are available for use.

### Sprig Functions

gucci ships with the [sprig templating functions library](http://masterminds.github.io/sprig/) offering a wide variety of template helpers.

### Built In Functions

Furthermore, this tool also includes custom functions:

- `shell`: For arbitrary shell commands

  ```
  {{ shell "echo hello world" }}
  ```

  and

  ```
  # guest: world
  {{ shell "echo hello " .guest }}
  ```

  Both produce:

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
