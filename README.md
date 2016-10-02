# Gucci

A simple cli templating tool written in golang

[![GitHub version](https://badge.fury.io/gh/noqcks%2Fgucci.svg)](https://badge.fury.io/gh/noqcks%2Fgucci)
[![License](https://img.shields.io/github/license/noqcks/gucci.svg)](https://github.com/noqcks/gucci/blob/master/LICENSE)
[![Build Status](https://travis-ci.org/noqcks/gucci.svg?branch=master)](https://travis-ci.org/noqcks/gucci)

# Get

```
go get github.com/noqcks/gucci
```

To install the binary, download the latest release and move it into your path.

```
wget -q https://github.com/noqcks/gucci/releases/download/v0.0.1/gucci-v0.0.1-darwin-amd64
mv gucci-v0.0.1-darwin-amd64 /usr/local/bin/gucci
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
