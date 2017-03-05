# revimport
`$ go get -u github.com/csigo/revimport`

**revimport** is to revert packages which import given packages in goimports path.
For example, if run `revimport` with `fmt` package, then it will print out all packages that import fmt.

**revimport** allows you to check if the modified packages affect those packages importing them.

## Usage

`$ revimport -pkgs='pkgname1[,pkgname2...]'`
