gen-base64-images
==================

#### Install

```bash
go get github.com/ririjin/gen-base64-images
go install github.com/ririjin/gen-base64-images
```


#### Usage

```bash
gen-base64-images -h
NAME:
   main - A new cli application

USAGE:
   main [global options] command [command options] [arguments...]

VERSION:
   0.0.0

COMMANDS:
     help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --src value    images source dir
   --dst value    base64-images js generate to
   --tpl value    base64-image file template (default: "module.exports = { uri:'{{.data}}' }")
   --help, -h     show help
   --version, -v  print the version

```


```go
gen-base64-images --src /home/user/pictures --dst /home/user/js-pictures
```

```bash
> ls -l /home/user/pictures
-rw-r--r--  1 ririjin  staff    39K Jun  8 14:46 xxxx.jpg

> ls -l /home/user/js-pictures
-rw-r--r--  1 ririjin  staff    39K Jun  8 14:46 xxxx_jpg.js

> cat /home/user/js-pictures/xxxx_jpg.js
uri:'data:image/jpeg;base64,/9j/4AAQSkZJRgABA....
```