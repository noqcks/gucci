# Gucci

A simple CLI templating tool written in golang. I created this because I wanted something that was more powerful than `envsubt` to use when templating files on the command line. I also used golang because I wanted a single binary.

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
VERSION=0.0.4
wget -q https://github.com/noqcks/gucci/releases/download/v${VERSION}/gucci-v${VERSION}-darwin-amd64
chmod +x gucci-v${VERSION}-darwin-amd64
mv gucci-v${VERSION}-darwin-amd64 /usr/local/bin/gucci
```


# Use

gucci can take input in multiple ways

### file

```
$ gucci template.tpl > template.conf
```

### stdin

```
$ gucci
Start typing stuff {{ print "here" }}
^d
Start typing stuff here
```

via piping

```
$ echo '{{ html "<escape-me/>" }}' | gucci
```

# Templating

### GoLang Functions

All of the existing [golang templating functions](https://golang.org/pkg/text/template/#hdr-Functions) are available for use.

### Built In Functions

This is a list of custom functions this tool adds that you can use:

- `split`: Used to split strings

  ```
  {{ range split .BACKENDS "," }}
    server {{ . }}
  {{ end }}
  ```

- `shell`: For arbitrary shell commands

   ```
   {{ shell "cat VERSION.txt" }}
   ```

## Example

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

### TODO

- add input sources (json & yaml files)
