package utils

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/iancoleman/strcase"
)

type (
	// APIGroup struct for api grouping
	APIGroup struct {
		GroupName string
		Base      string
		Child     []APIEndpoint
	}

	// APIEndpoint struct for api grouping
	APIEndpoint struct {
		Name   string
		Path   string
		Auth   bool
		Method string
	}
)

func parseGroup(input string) string {
	return input[len("@RouterBase "):]
}

func parseEndpoint(name string, input string) APIEndpoint {
	reMethods := regexp.MustCompile(`get|post|delete|put|patch`)
	args := strings.Split(input[len("@Router "):], " ")

	var path string
	var method string
	var auth = false

	for _, v := range args {
		// find methods
		if reMethods.FindString(v) != "" {
			method = strings.ToUpper(reMethods.FindString(v))
		}

		// find auth is enabled
		if v == "[auth]" {
			auth = true
		}

		// find path
		if strings.HasPrefix(v, "/") {
			path = v
		}
	}

	return APIEndpoint{
		Name:   name,
		Path:   path,
		Auth:   auth,
		Method: method,
	}
}

func readEndpoints() []APIGroup {
	var apiFiles []string
	var routes []APIGroup
	apiRoot := "./app/service"

	fset := token.NewFileSet()

	filepath.Walk(apiRoot, func(path string, info os.FileInfo, err error) error {
		if strings.HasSuffix(path, "_service.go") {
			apiFiles = append(apiFiles, path)
		}
		return nil
	})

	for _, path := range apiFiles {
		node, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
		if err != nil {
			log.Fatal(err)
		}

		var routeGroup APIGroup

		reEnd := regexp.MustCompile(`@Router.+`)
		reGroup := regexp.MustCompile(`@RouterGroup.+`)

		comments := []*ast.CommentGroup{}
		ast.Inspect(node, func(n ast.Node) bool {
			// collect comments
			c, ok := n.(*ast.CommentGroup)
			if ok {
				comments = append(comments, c)
			}

			switch fn := n.(type) {
			case *ast.FuncDecl:
				if reGroup.MatchString(fn.Doc.Text()) {
					routeGroup = APIGroup{
						GroupName: fn.Name.String(),
						Base:      parseGroup(reGroup.FindString(fn.Doc.Text())),
					}
				} else if reEnd.MatchString(fn.Doc.Text()) {
					meta := reEnd.FindString(fn.Doc.Text())
					child := append(routeGroup.Child, parseEndpoint(fn.Name.String(), meta))
					routeGroup.Child = child
				}
			}

			return true
		})

		routes = append(routes, routeGroup)
	}

	return routes
}

// GenerateRoutes automatically
func GenerateRoutes() {
	var contents string
	routeFile := "./app/router/router.go"
	routes := readEndpoints()

	for _, route := range routes {
		var body string
		serviceName := strings.ReplaceAll(route.GroupName, "New", "")
		varName := strcase.ToLowerCamel(serviceName)
		groupName := varName + "Group"
		body = fmt.Sprintf("\n\t\t// Generate route for %s\n", serviceName)
		body += fmt.Sprintf("\t\t%s := service.%s(models.DB)\n", varName, route.GroupName)
		body += fmt.Sprintf("\t\t%s := apiGroup.Group(\"%s\")\n", groupName, route.Base)
		body += fmt.Sprintf("\t\t{\n")

		for _, e := range route.Child {
			if e.Auth {
				body += fmt.Sprintf("\t\t\t%s.%s(\"%s\", mid.RequiresUserAuth, %s.%s)\n", groupName, e.Method, e.Path, varName, e.Name)
			} else {
				body += fmt.Sprintf("\t\t\t%s.%s(\"%s\", %s.%s)\n", groupName, strings.ToTitle(e.Method), e.Path, varName, e.Name)
			}
		}

		body += fmt.Sprintf("\t\t}\n")
		contents += body
	}

	codes, err := ioutil.ReadFile(routeFile)
	if err != nil {
		fmt.Println(err)
	}

	splittedString := strings.Split(string(codes), "\n")
	var startIndex int
	var endIndex int
	for index, line := range splittedString {
		if strings.Contains(line, "@StartCodeBlocks") {
			startIndex = index + 1
		} else if strings.Contains(line, "@EndCodeBlocks") {
			endIndex = index
		}
	}

	splittedString[endIndex] = contents
	splittedString[endIndex] += "\n\t\t// @EndCodeBlocks"
	splittedString = append(splittedString[:startIndex], splittedString[endIndex:]...)

	errWrite := ioutil.WriteFile(routeFile, []byte(strings.Join(splittedString, "\n")), 0777)
	if errWrite != nil {
		fmt.Println(errWrite)
	}
}
