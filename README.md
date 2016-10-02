# Gucci

A simple golang templating cli tool

[![Build Status](https://travis-ci.org/noqcks/gucci.svg?branch=master)](https://travis-ci.org/noqcks/gucci)

# Get

```
go get github.com/noqcks/gucci
```

# Use

```
gucci template.tpl > template.conf
```

# Templates

### Single

This follows the same syntax as golang templating.

For example an ENV var $LOCALHOST = 127.0.0.1

```
{{ .LOCALHOST }}
```

-->

```
127.0.0.1
```

simple enough!

### Iteration

For iteration of ENV vars, you can set $BACKENDS=server1.com,server2.com

```
{{- range split .BACKENDS "," }}
server {{ . }}
{{- end }}
```

-->

```
server server1.com
server server2.com
```

# TODO

- a more complete function list
- return template output to template, but errs to STDOUT
