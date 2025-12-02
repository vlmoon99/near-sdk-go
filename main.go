// codegen.go - Simple all-in-one code generator
package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"strings"
)

// MethodInfo holds parsed method information
type MethodInfo struct {
	Name       string
	Params     []Param
	Returns    []string
	IsPublic   bool
	IsPrivate  bool
	IsView     bool
	IsMutating bool
	IsPayable  bool
	MinDeposit string
}

type Param struct {
	Name string
	Type string
}

func main() {

	sourceFile := "./new_format/main.go"
	
	fmt.Println("🚀 NEAR Contract Code Generator")
	fmt.Println("================================")
	fmt.Printf("📖 Parsing: %s\n", sourceFile)

	// Parse the source file
	methods, err := parseContract(sourceFile)
	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("✓ Found %d methods\n\n", len(methods))

	// Display what we found
	displayMethods(methods)

	// Generate the code
	generatedCode := generateCode(methods)

	// Write to output file
	outputFile := "./new_format/generated_exports.go"
	if err := os.WriteFile(outputFile, []byte(generatedCode), 0644); err != nil {
		fmt.Printf("❌ Failed to write output: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("\n✅ Generated: %s\n", outputFile)
	fmt.Println("\n💡 Next step: Build with TinyGo")
	fmt.Println("   tinygo build -o contract.wasm -target=wasi .")
}

// parseContract parses the Go source file and extracts methods
func parseContract(filename string) ([]*MethodInfo, error) {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, filename, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	var methods []*MethodInfo

	// Walk through the AST
	ast.Inspect(file, func(n ast.Node) bool {
		fn, ok := n.(*ast.FuncDecl)
		if !ok {
			return true
		}

		// Only process methods (functions with receivers)
		if fn.Recv == nil || len(fn.Recv.List) == 0 {
			return true
		}

		method := extractMethod(fn)
		methods = append(methods, method)

		return true
	})

	return methods, nil
}

// extractMethod extracts method information from AST
func extractMethod(fn *ast.FuncDecl) *MethodInfo {
	method := &MethodInfo{
		Name: fn.Name.Name,
	}

	// Parse annotations from comments
	if fn.Doc != nil {
		for _, comment := range fn.Doc.List {
			parseAnnotation(comment.Text, method)
		}
	}

	// Extract parameters
	if fn.Type.Params != nil {
		for _, field := range fn.Type.Params.List {
			typeName := typeToString(field.Type)
			for _, name := range field.Names {
				method.Params = append(method.Params, Param{
					Name: name.Name,
					Type: typeName,
				})
			}
		}
	}

	// Extract return types
	if fn.Type.Results != nil {
		for _, field := range fn.Type.Results.List {
			method.Returns = append(method.Returns, typeToString(field.Type))
		}
	}

	return method
}

// parseAnnotation parses a comment line for annotations
func parseAnnotation(text string, method *MethodInfo) {
	text = strings.TrimSpace(text)
	text = strings.TrimPrefix(text, "//")
	text = strings.TrimSpace(text)

	if !strings.HasPrefix(text, "@contract:") {
		return
	}

	annotation := strings.TrimPrefix(text, "@contract:")
	parts := strings.Fields(annotation)
	if len(parts) == 0 {
		return
	}

	switch parts[0] {
	case "public":
		method.IsPublic = true
	case "private":
		method.IsPrivate = true
	case "view":
		method.IsView = true
	case "mutating":
		method.IsMutating = true
	case "payable":
		method.IsPayable = true
		for _, part := range parts[1:] {
			if strings.HasPrefix(part, "min_deposit=") {
				method.MinDeposit = strings.TrimPrefix(part, "min_deposit=")
			}
		}
	}
}

// typeToString converts AST type to string
func typeToString(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.Ident:
		return t.Name
	case *ast.StarExpr:
		return "*" + typeToString(t.X)
	case *ast.ArrayType:
		return "[]" + typeToString(t.Elt)
	case *ast.SelectorExpr:
		return typeToString(t.X) + "." + t.Sel.Name
	default:
		return "unknown"
	}
}

// displayMethods shows what was found
func displayMethods(methods []*MethodInfo) {
	publicCount := 0
	privateCount := 0
	viewCount := 0
	payableCount := 0

	for _, m := range methods {
		if m.IsPrivate {
			privateCount++
			fmt.Printf("  ⊗ %s (private - no export)\n", m.Name)
			continue
		}

		if m.IsPublic {
			publicCount++
			icon := "  ✓"
			if m.IsPayable {
				icon = "  💰"
				payableCount++
			}
			if m.IsView {
				icon = "  👁"
				viewCount++
			}

			fmt.Printf("%s %s(", icon, m.Name)
			for i, p := range m.Params {
				if i > 0 {
					fmt.Print(", ")
				}
				fmt.Printf("%s: %s", p.Name, p.Type)
			}
			fmt.Print(")")

			if len(m.Returns) > 0 {
				fmt.Print(" → ")
				fmt.Print(strings.Join(m.Returns, ", "))
			}

			tags := []string{}
			if m.IsView {
				tags = append(tags, "view")
			}
			if m.IsMutating {
				tags = append(tags, "mutating")
			}
			if m.IsPayable {
				tags = append(tags, fmt.Sprintf("payable[%s NEAR]", m.MinDeposit))
			}

			if len(tags) > 0 {
				fmt.Printf(" [%s]", strings.Join(tags, ", "))
			}
			fmt.Println()
		}
	}

	fmt.Printf("\n📊 Summary:\n")
	fmt.Printf("   Public: %d | Private: %d | View: %d | Payable: %d\n",
		publicCount, privateCount, viewCount, payableCount)
}

// generateCode generates the exports file
func generateCode(methods []*MethodInfo) string {
	var sb strings.Builder

	// Header
	sb.WriteString("// Code generated by NEAR contract generator. DO NOT EDIT.\n\n")
	sb.WriteString("package main\n\n")
	sb.WriteString("import (\n")
	sb.WriteString("\t\"encoding/json\"\n")
	sb.WriteString("\t\"fmt\"\n")
	sb.WriteString("\t\n")
	sb.WriteString("\tcontractBuilder \"github.com/vlmoon99/near-sdk-go/contract\"\n")
	sb.WriteString(")\n\n")

	// Generate exports for public methods only
	for _, m := range methods {
		if !m.IsPublic {
			continue
		}

		exportName := toSnakeCase(m.Name)
		sb.WriteString(fmt.Sprintf("//go:export %s\n", exportName))
		sb.WriteString(fmt.Sprintf("func %s() {\n", exportName))
		sb.WriteString("\tcontractBuilder.HandleClientJSONInput(func(input *contractBuilder.ContractInput) error {\n")

		// Add view comment
		if m.IsView {
			sb.WriteString("\t\t// View function - read-only\n")
		}

		// Add payment validation
		if m.IsPayable {
			sb.WriteString(fmt.Sprintf("\t\t// Validate payment requirement: %s NEAR\n", m.MinDeposit))
			sb.WriteString(fmt.Sprintf("\t\tif err := validatePayment(\"%s\"); err != nil {\n", m.MinDeposit))
			sb.WriteString("\t\t\treturn fmt.Errorf(\"insufficient payment: %%w\", err)\n")
			sb.WriteString("\t\t}\n\n")
		}

		// Parse parameters
		for _, p := range m.Params {
			sb.WriteString(fmt.Sprintf("\t\t// Parse parameter: %s\n", p.Name))
			sb.WriteString(generateParamParser(p))
		}

		// Call the actual method
		sb.WriteString("\t\t// Get contract instance\n")
		sb.WriteString("\t\tcontract := GetContract().(*Counter)\n\n")
		sb.WriteString("\t\t// Call method\n")

		if len(m.Returns) > 0 {
			sb.WriteString("\t\tresult := contract.")
		} else {
			sb.WriteString("\t\tcontract.")
		}

		sb.WriteString(m.Name)
		sb.WriteString("(")
		for i, p := range m.Params {
			if i > 0 {
				sb.WriteString(", ")
			}
			sb.WriteString(p.Name)
		}
		sb.WriteString(")\n\n")

		// Return result
		if len(m.Returns) > 0 {
			sb.WriteString("\t\tcontractBuilder.ReturnValue(result)\n")
		} else {
			sb.WriteString("\t\tcontractBuilder.ReturnValue(\"Success\")\n")
		}

		sb.WriteString("\t\treturn nil\n")
		sb.WriteString("\t})\n")
		sb.WriteString("}\n\n")
	}

	// Add helper functions
	sb.WriteString("// Helper: validate payment for payable methods\n")
	sb.WriteString("func validatePayment(minDeposit string) error {\n")
	sb.WriteString("\t// TODO: Implement actual payment validation\n")
	sb.WriteString("\t// Check that context.AttachedDeposit() >= parseNEAR(minDeposit)\n")
	sb.WriteString("\treturn nil\n")
	sb.WriteString("}\n")

	return sb.String()
}

// generateParamParser generates code to parse a parameter
func generateParamParser(p Param) string {
	var sb strings.Builder

	switch p.Type {
	case "string":
		sb.WriteString(fmt.Sprintf("\t\t%s, err := input.JSON.GetString(\"%s\")\n", p.Name, p.Name))
		sb.WriteString("\t\tif err != nil {\n")
		sb.WriteString(fmt.Sprintf("\t\t\treturn fmt.Errorf(\"failed to parse '%s': %%w\", err)\n", p.Name))
		sb.WriteString("\t\t}\n\n")

	case "int":
		sb.WriteString(fmt.Sprintf("\t\t%sInt64, err := input.JSON.GetInt(\"%s\")\n", p.Name, p.Name))
		sb.WriteString("\t\tif err != nil {\n")
		sb.WriteString(fmt.Sprintf("\t\t\treturn fmt.Errorf(\"failed to parse '%s': %%w\", err)\n", p.Name))
		sb.WriteString("\t\t}\n")
		sb.WriteString(fmt.Sprintf("\t\t%s := int(%sInt64)\n\n", p.Name, p.Name))

	case "bool":
		sb.WriteString(fmt.Sprintf("\t\t%s, err := input.JSON.GetBool(\"%s\")\n", p.Name, p.Name))
		sb.WriteString("\t\tif err != nil {\n")
		sb.WriteString(fmt.Sprintf("\t\t\treturn fmt.Errorf(\"failed to parse '%s': %%w\", err)\n", p.Name))
		sb.WriteString("\t\t}\n\n")

	default:
		// Complex types - use JSON unmarshaling
		sb.WriteString(fmt.Sprintf("\t\tvar %s %s\n", p.Name, p.Type))
		sb.WriteString("\t\tvar params struct {\n")
		sb.WriteString(fmt.Sprintf("\t\t\tValue %s `json:\"%s\"`\n", p.Type, p.Name))
		sb.WriteString("\t\t}\n")
		sb.WriteString("\t\tif err := json.Unmarshal(input.RawJSON, &params); err != nil {\n")
		sb.WriteString(fmt.Sprintf("\t\t\treturn fmt.Errorf(\"failed to parse '%s': %%w\", err)\n", p.Name))
		sb.WriteString("\t\t}\n")
		sb.WriteString(fmt.Sprintf("\t\t%s = params.Value\n\n", p.Name))
	}

	return sb.String()
}

// toSnakeCase converts PascalCase to snake_case
func toSnakeCase(s string) string {
	var result strings.Builder
	for i, r := range s {
		if i > 0 && r >= 'A' && r <= 'Z' {
			result.WriteRune('_')
		}
		result.WriteRune(r)
	}
	return strings.ToLower(result.String())
}