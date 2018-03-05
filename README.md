# jsonTail
tail command for json log file

## How to use
`jsonTail` require [`github.com/hpcloud/tail`](https://github.com/hpcloud/tail/blob/master/tail.go).
You need to execute following command to install it : 
```bash
go get github.com/hpcloud/tail
```

Then, run following command to create `bin/jsonTail` binary
```bash
make build
```

## Options
Following command will output available options
```bash
jsonTail --help
```
