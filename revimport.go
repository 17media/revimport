/*Package main takes a list of go package and prints all reverse-imported packages.*/
package main

import (
	"flag"
	"fmt"
	"go/build"
	"sort"
	"strings"

	"golang.org/x/tools/refactor/importgraph"
)

var (
	pkgsStr = flag.String("pkgs", "", "a list of packages to be propagated, separated by comma.")
	layers  = flag.String("layers", "one", "prints \"one\" layer of reverse imported packages, "+
		"or propagate to \"all\" effected files.")
	omitStr = flag.String("omit", "htc.com/csi/deprecated", "a list of skip folder's term. We should skip the checking "+
		"path when these terms appear in the test.")
)

func main() {
	flag.Parse()

	if *pkgsStr == "" {
		return
	}
	// Build the import graph.
	// TODO: Cache this.
	_, reverse, err := importgraph.Build(&build.Default)
	if err != nil {
		panic(fmt.Sprintf("can't establish import dependencies. err %v", err))
	}

	// Generate package slice from a space separated string
	pkgs := strings.Split(*pkgsStr, ",")

	// Generate package slice from a space separated string
	omits := strings.Split(*omitStr, ",")

	// propagated stores propagated packages.
	propagated := make(map[string]bool)
	switch *layers {
	case "one":
		// Only propagate to the immediate imported packages.
		for _, p := range pkgs {
			// Put myself into propagation list
			propagated[p] = true
			for k, v := range reverse[p] {
				propagated[k] = v
			}
		}
	case "all":
		// Propagate all the way to all effected files.
		propagated = reverse.Search(pkgs...)
	default:
		panic(fmt.Sprintf("undefined layers %s", *layers))
	}

	var pSlice []string
	for ppath := range propagated {
		if !containsAny(omits, ppath) {
			pSlice = append(pSlice, ppath)
		}
	}
	sort.Strings(pSlice)
	// print propagated packages line by line
	for _, p := range pSlice {
		fmt.Println(p)
	}
}

func containsAny(srcs []string, target string) bool {
	for _, src := range srcs {
		if strings.Contains(src, target) {
			return true
		}
	}
	return false
}
