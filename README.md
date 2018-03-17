# jsonTail
tail command for json log file

## How to install
Retrieve last version
```bash
go get github.com/yoanm/jsonTail
```

Then, install it
```bash
go install github.com/yoanm/jsonTail
```

## Options
Following command will output available options
```bash
jsonTail --help
```

 * `-date` : will print line handling date in front of each processed lines
 
### File

 * `-f FILE_PATH` : will open `FILE_PATH` and wait for update. *Behaves like `tail -f` command*
 * `-F FILE_PATH` : Same behavior than `-f` but file will be reopened if recreated. *Behaves like `tail -F` command*
 
### Fields matching/exclusion

**Fields matching or exclusion only work with objects**

 * `-only FIELD` : will output only specified field. Multiple fields could be specified by using multiple `-only`.
   Field matching is made thanks to `github.com/tidwall/gjson` package.
   See [`github.com/tidwall/gjson` path syntac](https://github.com/tidwall/gjson#path-syntax) for more information about paths

 * `-exclude FIELD` : will excluded specified field from output. Multiple fields could be specified by using multiple `-exclude`.
   Field matching is made thanks to `github.com/tidwall/sjson` package.
   See [`github.com/tidwall/sjson` path syntac](https://github.com/tidwall/sjson#path-syntax) for more information about paths
   

## Dependencies

Following external packages are used under the hood : 

 * [github.com/hpcloud/tail](https://github.com/hpcloud/tail)
 * [github.com/tidwall/sjson](https://github.com/tidwall/sjson)
 * [github.com/tidwall/gjson](https://github.com/tidwall/gjson)
 * [github.com/hokaccha/go-prettyjson](https://github.com/hokaccha/go-prettyjson)
 * [github.com/fatih/color](https://github.com/fatih/color)
