package main

import (
	"fmt"
	"os"

	"github.com/zackslash/goviz/dotwriter"
	"github.com/zackslash/goviz/goimport"
	"github.com/zackslash/goviz/metrics"

	cli "gopkg.in/alecthomas/kingpin.v2"
)

var (
	inputDir            = cli.Flag("input", "project directory").Required().Short('i').String()
	outputFile          = cli.Flag("output", "output file").Default("STDOUT").Short('o').String()
	depth               = cli.Flag("depth", "max plot depth of the dependency tree").Default("2").Short('d').Int()
	reversed            = cli.Flag("focus", "focus on the specific module").Default("").Short('f').String()
	seekPath            = cli.Flag("search", "top directory of searching").Default("").Short('s').String()
	plotLeaf            = cli.Flag("leaf", "if leaf nodes are plotted").Default("false").Short('l').Bool()
	useMetrics          = cli.Flag("metrics", "display module metrics").Default("false").Short('m').Bool()
	projectPackagesOnly = cli.Flag("ppackage", "only include packages from immediate project").Default("true").Short('p').Bool()
)

func main() {
	cli.Version("1.0")
	cli.Parse()
	res := process()
	os.Exit(res)
}

func process() int {
	factory := goimport.ParseRelation(
		*inputDir,
		*seekPath,
		*plotLeaf,
	)

	if factory == nil {
		fmt.Errorf("inputdir does not exist.\n go get %s", *inputDir)
		return 1
	}

	root := factory.GetRoot()
	if !root.HasFiles() {
		fmt.Errorf("%s has no go files", root.ImportPath)
		return 1
	}

	if 0 > *depth {
		fmt.Errorf("-d or --depth should have positive int")
		return 1
	}

	output := getOutputWriter(*outputFile)
	if *useMetrics {
		metricsWriter := metrics.New(output)
		metricsWriter.Plot(pathToNode(factory.GetAll()))
		return 0
	}

	writer := dotwriter.New(output)
	writer.MaxDepth = *depth
	if *reversed == "" {
		writer.PlotGraph(root, *projectPackagesOnly)
		return 0
	}

	writer.Reversed = true

	rroot := factory.Get(*reversed)
	if rroot == nil {
		fmt.Errorf("-r %s does not exist", *reversed)
		return 1
	}

	if !rroot.HasFiles() {
		fmt.Errorf("-r %s has no go files", *reversed)
		return 1
	}

	writer.PlotGraph(rroot, *projectPackagesOnly)
	return 0
}

func pathToNode(pathes []*goimport.ImportPath) []dotwriter.IDotNode {
	r := make([]dotwriter.IDotNode, len(pathes))

	for i := range pathes {
		r[i] = pathes[i]
	}

	return r
}
func getOutputWriter(name string) *os.File {
	if name == "STDOUT" {
		return os.Stdout
	}
	if name == "STDERR" {
		return os.Stderr
	}
	f, _ := os.Create(name)

	return f
}
