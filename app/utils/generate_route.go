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
		Param  interface{}
	}
)

func parseGroup(input string) string {
	path := input[len("@RouterGroup "):]
	return strings.TrimSpace(path)
}

func parseEndpoint(name string, input string, param interface{}) APIEndpoint {
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
			path = strings.TrimSpace(v)
		}
	}

	return APIEndpoint{
		Name:   name,
		Path:   path,
		Auth:   auth,
		Method: method,
		Param:  param,
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
		src, _ := ioutil.ReadFile(path)
		node, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
		if err != nil {
			log.Fatal(err)
		}

		var routeGroup APIGroup

		reEnd := regexp.MustCompile(`@Router.+`)
		reGroup := regexp.MustCompile(`@RouterGroup.+`)

		offset := node.Pos()
		stringSrc := string(src)

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
					var paramName interface{}
					for _, param := range fn.Type.Params.List {
						typeName := param.Type
						typeStr := stringSrc[typeName.Pos()-offset : typeName.End()-offset]
						if !strings.Contains(typeStr, "gin.Context") {
							paramName = typeStr
						}
					}
					meta := reEnd.FindString(fn.Doc.Text())
					child := append(routeGroup.Child, parseEndpoint(fn.Name.String(), meta, paramName))
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
func GenerateRoutes() bool {
	generated := false
	var contents string
	var routeConst string = fmt.Sprintf("package endpoint\n\nconst (\n")
	routeFile := "./app/router/router.go"
	routeConstFile := "./tests/endpoint/route_const.go"
	routes := readEndpoints()

	for _, route := range routes {
		var body string
		var constLine string
		serviceName := strings.ReplaceAll(route.GroupName, "New", "")
		varName := strcase.ToLowerCamel(serviceName)
		groupName := varName + "Group"
		body = fmt.Sprintf("\n\t\t// Generate route for %s\n", serviceName)
		body += fmt.Sprintf("\t\t%s := service.%s()\n", varName, route.GroupName)
		body += fmt.Sprintf("\t\t%s := apiGroup.Group(\"%s\")\n", groupName, route.Base)
		body += fmt.Sprintf("\t\t{\n")

		for _, e := range route.Child {
			var validator string
			if e.Param != nil {
				param := fmt.Sprintf("service.%s", e.Param)
				if strings.HasPrefix(e.Param.(string), "*repo") {
					param = e.Param.(string)
				}

				param = strings.ReplaceAll(param, "*", "")

				validator = "func(c *gin.Context) {\n"
				// validator += fmt.Sprintf("\tquery := &%s{}\n", param)
				validator += fmt.Sprintf("\tmid.RequestValidator(c, &%s{})\n", param)
				validator += "\t}, "
				validator += "func(c *gin.Context) {\n"
				validator += fmt.Sprintf("\tquery, ok := c.MustGet(\"validated\").(*%s)\n", param)
				validator += "if !ok {\n log.Println(\"validated not set\")\n}\n"
				validator += fmt.Sprintf("\t%s.%s(c, query)\n}", varName, e.Name)
			} else {
				validator = fmt.Sprintf("%s.%s", varName, e.Name)
			}
			if e.Auth {
				body += fmt.Sprintf("\t\t\t%s.%s(\"%s\", mid.RequiresUserAuth, %s)\n", groupName, e.Method, e.Path, validator)
			} else {
				body += fmt.Sprintf("\t\t\t%s.%s(\"%s\", %s)\n", groupName, strings.ToTitle(e.Method), e.Path, validator)
			}

			constLine += fmt.Sprintf("\t// %sEndpoint for testing only\n", e.Name)
			constLine += fmt.Sprintf("\t%s = \"%s%s\"\n", e.Name, route.Base, e.Path)
		}

		body += fmt.Sprintf("\t\t}\n")
		contents += body
		routeConst += constLine
	}

	routeConst += fmt.Sprintf(")")

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
	ioutil.WriteFile(routeConstFile, []byte(routeConst), 0777)
	if errWrite != nil {
		generated = true
	}

	return generated
}
