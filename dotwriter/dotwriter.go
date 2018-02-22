package dotwriter

import (
	"fmt"
	"io"
	"strconv"
	"strings"
)

type IDotNode interface {
	Name() string
	Label() string
	Shape() string
	Style() string
	Children() []IDotNode
	Parents() []IDotNode
}

type DotWriter struct {
	output   io.Writer
	MaxDepth int
	Reversed bool
}

type plotCtx struct {
	nodeFlags map[string]bool
	edgeFlags map[string]bool
	level     int
}

func (ctx *plotCtx) isPlottedNode(node IDotNode) bool {
	_, ok := ctx.nodeFlags[node.Name()]
	return ok
}
func (ctx *plotCtx) setPlotted(node IDotNode) {
	_, ok := ctx.nodeFlags[node.Name()]
	if !ok {
		ctx.nodeFlags[node.Name()] = true
	}

}

func (ctx *plotCtx) isDepthOver() bool {
	return (ctx.level <= 0)
}
func (ctx *plotCtx) Deeper() *plotCtx {
	return &plotCtx{
		nodeFlags: ctx.nodeFlags,
		edgeFlags: ctx.edgeFlags,
		level:     ctx.level - 1,
	}
}
func newPlotContext(level int) *plotCtx {
	return &plotCtx{
		level:     level,
		nodeFlags: make(map[string]bool),
		edgeFlags: make(map[string]bool),
	}
}
func (ctx *plotCtx) isPlottedEdge(nodeA, nodeB IDotNode) bool {
	edgeName := fmt.Sprintf("%s->%s", nodeA.Name(), nodeB.Name())
	_, ok := ctx.edgeFlags[edgeName]
	if !ok {
		ctx.edgeFlags[edgeName] = true
	}
	return ok
}

func New(output io.Writer) *DotWriter {
	return &DotWriter{output: output}
}

func (dw *DotWriter) PlotGraph(root IDotNode, projectPackageOnly bool) {
	dw.printLine("digraph main{")
	dw.printLine("\tedge[arrowhead=vee]")
	dw.printLine("\tgraph [rankdir=LR,compound=true,ranksep=1.0];")
	dw.plotNode(newPlotContext(dw.MaxDepth), root, projectPackageOnly)
	dw.printLine("}")
}

var rootPackage string

func setRootProjectPackage(p string) {
	rootPackage = p
}

func getRootProjectPackage() string {
	return rootPackage
}

// checks if given package is within scope
func isWithinScope(p string) bool {
	rp := getRootProjectPackage()
	r := strings.Split(rp, "/")

	// get 1st three levels i.e github.com/two/three
	// todo: setter for scope level
	sc := rp
	scopeLevel := 3
	if len(r) >= scopeLevel-1 {
		sc = strings.Join(r[0:scopeLevel], "/")
	}
	return strings.HasPrefix(p, sc)
}

func (dw *DotWriter) plotNode(ctx *plotCtx, node IDotNode, projectPackageOnly bool) {
	if projectPackageOnly {
		if getRootProjectPackage() == "" {
			setRootProjectPackage(node.Name())
		}

		if !isWithinScope(node.Name()) {
			return
		}
	}

	if ctx.isPlottedNode(node) {
		return
	}
	if ctx.isDepthOver() {
		return
	}
	ctx.setPlotted(node)
	dw.plotNodeStyle(node)
	for _, s := range dw.getDependency(node) {
		if projectPackageOnly {
			if !isWithinScope(s.Name()) {
				continue
			}
		}
		dw.plotEdge(ctx, node, s)
		dw.plotNode(ctx.Deeper(), s, projectPackageOnly)
	}
}

func (dw *DotWriter) getDependency(node IDotNode) []IDotNode {
	if dw.Reversed {
		return node.Parents()
	}
	return node.Children()
}
func (dw *DotWriter) plotNodeStyle(node IDotNode) {
	dw.printFormat("\t/* plot %s */\n", node.Name())
	dw.printFormat("\t%s[shape=%s,label=\"%s\",style=%s]\n",
		escape(node.Name()),
		escape(node.Shape()),
		node.Label(),
		escape(node.Style()),
	)
}

func (dw *DotWriter) plotEdge(ctx *plotCtx, nodeA, nodeB IDotNode) {
	if ctx.isPlottedEdge(nodeA, nodeB) {
		return
	}
	dir := "forward"
	if dw.Reversed {
		dir = "back"
	}
	dw.printFormat("\t%s -> %s[dir=%s]\n", escape(nodeA.Name()), escape(nodeB.Name()), dir)
}

func (dw *DotWriter) printLine(str string) {
	fmt.Fprintln(dw.output, str)
}

func (dw *DotWriter) printFormat(pattern string, args ...interface{}) {
	fmt.Fprintf(dw.output, pattern, args...)
}

func escape(target string) string {
	return strconv.Quote(target)
}
