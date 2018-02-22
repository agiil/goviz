goviz
=====

Goviz is a golang project dependency visualization tool (see example output below)

'zackslash/goviz' is an updated version of the original hirokidaichi project.

![](https://raw.githubusercontent.com/zackslash/goviz/master/images/own.png)


## Install

```
brew install graphviz
```

```
go get github.com/zackslash/goviz
go install github.com/zackslash/goviz
```

## Usage

```
goviz -i <your_project_path> | dot -Tpng -o <output_file_name>.png

Example: goviz -i github.com/hashicorp/serf | dot -Tpng -o diagram.png
```

```
usage: goviz --input=INPUT [<flags>]

Flags:
      --help             Show context-sensitive help (also try --help-long and --help-man).
  -i, --input=INPUT      project directory
  -o, --output="STDOUT"  output file
  -d, --depth=2          max plot depth of the dependency tree
  -f, --focus=""         focus on the specific module
  -s, --search=""        top directory of searching
  -l, --leaf             if leaf nodes are plotted
  -m, --metrics          display module metrics
  -p, --ppackage         only include packages from immediate project
      --version          Show application version.

```

## License

MIT

## Original goviz Author

hirokidaichi [at] gmail.com
