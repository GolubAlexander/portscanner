# Port scanner

Port scanner CLI tool, another one.

### Compile

To compile a scanner you must have golang compiler.
> go build -o portscanner

### Usage
```
-i string
    IP or domain names to scan
-o string
    Result list (default "output.txt")
-t int
    Timeout of a port scan (default 1 second)
```

### Input file format

Every line is separated host. For example,
```
192.168.0.1
192.168.0.3
ya.ru
```