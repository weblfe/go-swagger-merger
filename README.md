# go Swagger merger
To merge a few swagger YAML/JSON files into one.

Install Go if you don't have one.

`https://golang.org/doc/install`

Install the command line tool first.
```shell
go install github.com/weblfe/go-swagger-merger
```

The command below will merge /data/swagger1.yaml /data/swagger2.yaml and save result file in the /data/swagger.yaml. The library supports more than two files to merge. You can add more paths to the list /data/swagger3.yaml, /data/swaggerN.yaml.
```shell
go-swagger-merger -o /data/swagger.yaml  /data/swagger1.yaml /data/swagger2.json
```
Attention. The order of the files is essential, and the following file overwrites the same fields from the previous file.