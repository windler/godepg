# godepg
`godepg` generates a dependency graph for a go package using `graphviz`.

## Prerequisites
In order to generate graphs you have to install [graphviz](https://graphviz.gitlab.io/)
## Installation
`go get github.com/windler/godepg`

## Usage
Simplest way to use the tool is to type `godepg -p <package> -o <file.png>`. There are some options available which you can list via `godepg -h`:

```(bash)
GLOBAL OPTIONS:
   -o file, --output file         destination file to write png to
   -p package, --package package  the package to analyze
   -n, --no-go-packages           hide gos buildin packages
   -d value, --depth value        limit the depth of the graph (default: -1)
   -f value, --filter value       filter package name
   -m, --my-packages-only         show only subpackages of scanned package
   --help, -h                     show help
   --version, -v                  print the version
```

## Output
All graphs are written to `~/godepg/<pkg>_timestamp.png` if option `-o` is not present. You can change the home directory by setting the env `GODEPG_HOME`.
There is also a small build in webserver you can launch to browse the generated graphs in browser. Launch it with command `ws`:
```(bash)
godepg ws
Started webserver on port 8000...
http://localhost:8000
```

## Examples
Following you can find sample outputs of the [asd package](https://github.com/windler/asd).

### Without go internal packages and specific output file
```(bash)
godepg -p github.com/windler/asd -o ~/asd_package.png --no-go-packages
```
![asd no go packages](asd_no_go_pkgs.png)

### Only sub packages
```(bash)
godepg -p github.com/windler/asd --my-packages-only
```
![asd only sub](asd_my_only.png)

### Without go internal packages and custom filter
```(bash)
godepg -p github.com/windler/asd --no-go-packages -f internal/test
```
![asd custom filter](asd_custom_filter.png)

### Other package
```(bash)
godepg -p github.com/windler/godepg
```
![godepg graph](godepg_graph.png)

# Future work
- recognize and render isolated subgraphs
- render dependencies of go-file 

