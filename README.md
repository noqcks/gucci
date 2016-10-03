# Gucci

A simple cli templating tool written in golang. I created this because I wanted something that was more powerful than `envsubt` to use when templating files at Docker container start.

[![GitHub version](https://badge.fury.io/gh/noqcks%2Fgucci.svg)](https://badge.fury.io/gh/noqcks%2Fgucci)
[![License](https://img.shields.io/github/license/noqcks/gucci.svg)](https://github.com/noqcks/gucci/blob/master/LICENSE)
[![Build Status](https://travis-ci.org/noqcks/gucci.svg?branch=master)](https://travis-ci.org/noqcks/gucci)

# Get

If you have go installed

```
go get github.com/noqcks/gucci
```

Or you can just download the binary and move it into your path

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

This follows the same syntax as [golang templating](https://golang.org/pkg/text/template/).

For example an ENV var $LOCALHOST = 127.0.0.1

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

### Iteration

For iteration of ENV vars, you can set $BACKENDS=server1.com,server2.com

```
# template.tpl
{{- range split .BACKENDS "," }}
server {{ . }}
{{- end }}
```

`gucci template.tpl > template.conf` -->


```
# template.conf
server server1.com
server server2.com
```

# TODO

- a more complete function list
- return template output to template, but errs to terminal (not template)
