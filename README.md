## comments
comments 主要用于把readme.md里面的中文转成英文

## usage
```console
Usage of ./comments:
  -in string
    	(must)input file
  -out string
    	(must)output file
  -overwrite
    	(must)Can overwrite files
```

## example
```bash
go build comments
./comments -in  README.md -out README_EN.md -overwrite
```
