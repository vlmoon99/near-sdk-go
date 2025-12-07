// codegen.go - Multi-file recursive code generator with state management
package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
)

// MethodInfo holds parsed method information
type MethodInfo struct {
	Name         string
	ReceiverType string
	Params       []Param
	Returns      []string
	IsPublic     bool
	IsPrivate    bool
	IsView       bool
	IsMutating   bool
	IsPayable    bool
	MinDeposit   string
	FilePath     string
	RelativePath string
}

type Param struct {
	Name string
	Type string
}

// StateInfo holds contract state struct information
type StateInfo struct {
	Name         string
	Fields       []FieldInfo
	FilePath     string
	RelativePath string
}

type FieldInfo struct {
	Name string
	Type string
}

func main() {
	rootDir := "./new_format"
	outputDir := "./new_format"

	fmt.Println("🚀 NEAR Contract Code Generator")
	fmt.Println("================================")
	fmt.Printf("📂 Scanning from: %s\n\n", rootDir)

	// Parse all files recursively
	allMethods, stateStructs, err := parseAllFilesRecursive(rootDir)
	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
		os.Exit(1)
	}

	// Display state structs
	if len(stateStructs) > 0 {
		fmt.Printf("📦 State Structs:\n")
		for _, s := range stateStructs {
			fmt.Printf("  🗄️  %s (%s)\n", s.Name, s.RelativePath)
			for _, f := range s.Fields {
				fmt.Printf("      - %s: %s\n", f.Name, f.Type)
			}
		}
		fmt.Println()
	}

	if len(allMethods) == 0 {
		fmt.Println("⚠️  No methods with @contract annotations found")
		os.Exit(0)
	}

	fmt.Printf("✓ Found %d methods total\n\n", len(allMethods))

	// Display methods
	displayMethods(allMethods)

	// Generate the unified code
	generatedCode := generateCode(allMethods, stateStructs)

	// Write to output file
	outputFile := filepath.Join(outputDir, "generated_exports.go")

	if err := os.MkdirAll(outputDir, 0755); err != nil {
		fmt.Printf("❌ Failed to create output directory: %v\n", err)
		os.Exit(1)
	}

	if err := os.WriteFile(outputFile, []byte(generatedCode), 0644); err != nil {
		fmt.Printf("❌ Failed to write output: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("\n✅ Generated: %s\n", outputFile)
	fmt.Println("\n💡 Next step: Build with TinyGo")
	fmt.Println("   cd new_format && tinygo build -size short -no-debug -o main.wasm -target wasm-unknown generated_exports.go main.go")
}

// parseAllFilesRecursive recursively scans all directories
func parseAllFilesRecursive(rootDir string) ([]*MethodInfo, []*StateInfo, error) {
	var allMethods []*MethodInfo
	var stateStructs []*StateInfo

	err := filepath.WalkDir(rootDir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			name := d.Name()
			if strings.HasPrefix(name, ".") ||
				name == "vendor" ||
				name == "node_modules" ||
				name == "testdata" {
				return filepath.SkipDir
			}
			return nil
		}

		if !strings.HasSuffix(path, ".go") {
			return nil
		}

		if strings.HasSuffix(path, "_test.go") {
			return nil
		}

		fileName := filepath.Base(path)
		if strings.HasPrefix(fileName, "generated_") {
			fmt.Printf("  ⏭️  Skipping: %s (generated file)\n", path)
			return nil
		}

		relPath, err := filepath.Rel(rootDir, path)
		if err != nil {
			relPath = path
		}

		fmt.Printf("  📄 Parsing: %s\n", relPath)

		methods, states, err := parseContract(path, relPath)
		if err != nil {
			fmt.Printf("  ⚠️  Warning: failed to parse %s: %v\n", relPath, err)
			return nil
		}

		if len(methods) > 0 {
			allMethods = append(allMethods, methods...)
		}

		if len(states) > 0 {
			stateStructs = append(stateStructs, states...)
		}

		return nil
	})

	if err != nil {
		return nil, nil, fmt.Errorf("failed to walk directory tree: %w", err)
	}

	return allMethods, stateStructs, nil
}

// parseContract parses a single Go file
func parseContract(filePath string, relativePath string) ([]*MethodInfo, []*StateInfo, error) {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
	if err != nil {
		return nil, nil, err
	}

	var methods []*MethodInfo
	var stateStructs []*StateInfo

	ast.Inspect(file, func(n ast.Node) bool {
		switch node := n.(type) {
		case *ast.GenDecl:
			// Look for type declarations (structs)
			if node.Tok == token.TYPE {
				for _, spec := range node.Specs {
					typeSpec, ok := spec.(*ast.TypeSpec)
					if !ok {
						continue
					}

					structType, ok := typeSpec.Type.(*ast.StructType)
					if !ok {
						continue
					}

					// Check for @contract:state annotation
					if node.Doc != nil && hasStateAnnotation(node.Doc) {
						state := extractStateInfo(typeSpec, structType)
						state.FilePath = filePath
						state.RelativePath = relativePath
						stateStructs = append(stateStructs, state)
					}
				}
			}

		case *ast.FuncDecl:
			// Only process methods
			if node.Recv == nil || len(node.Recv.List) == 0 {
				return true
			}

			method := extractMethod(node)
			method.FilePath = filePath
			method.RelativePath = relativePath

			if method.IsPublic || method.IsPrivate {
				methods = append(methods, method)
			}
		}

		return true
	})

	return methods, stateStructs, nil
}

// hasStateAnnotation checks if comment group has @contract:state
func hasStateAnnotation(doc *ast.CommentGroup) bool {
	for _, comment := range doc.List {
		text := strings.TrimSpace(comment.Text)
		text = strings.TrimPrefix(text, "//")
		text = strings.TrimSpace(text)
		if text == "@contract:state" {
			return true
		}
	}
	return false
}

// extractStateInfo extracts struct field information
func extractStateInfo(typeSpec *ast.TypeSpec, structType *ast.StructType) *StateInfo {
	state := &StateInfo{
		Name: typeSpec.Name.Name,
	}

	if structType.Fields != nil {
		for _, field := range structType.Fields.List {
			fieldType := typeToString(field.Type)
			for _, name := range field.Names {
				state.Fields = append(state.Fields, FieldInfo{
					Name: name.Name,
					Type: fieldType,
				})
			}
		}
	}

	return state
}

// extractMethod extracts method information from AST
func extractMethod(fn *ast.FuncDecl) *MethodInfo {
	method := &MethodInfo{
		Name: fn.Name.Name,
	}

	if fn.Recv != nil && len(fn.Recv.List) > 0 {
		method.ReceiverType = extractReceiverType(fn.Recv.List[0].Type)
	}

	if fn.Doc != nil {
		for _, comment := range fn.Doc.List {
			parseAnnotation(comment.Text, method)
		}
	}

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

	if fn.Type.Results != nil {
		for _, field := range fn.Type.Results.List {
			method.Returns = append(method.Returns, typeToString(field.Type))
		}
	}

	return method
}

func extractReceiverType(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.Ident:
		return t.Name
	case *ast.StarExpr:
		return extractReceiverType(t.X)
	default:
		return "Unknown"
	}
}

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
	case *ast.MapType:
		return "map[" + typeToString(t.Key) + "]" + typeToString(t.Value)
	default:
		return "unknown"
	}
}

func displayMethods(methods []*MethodInfo) {
	publicCount := 0
	privateCount := 0
	viewCount := 0
	payableCount := 0

	fileGroups := make(map[string][]*MethodInfo)
	for _, m := range methods {
		fileGroups[m.RelativePath] = append(fileGroups[m.RelativePath], m)
	}

	for filePath, fileMethods := range fileGroups {
		fmt.Printf("📄 %s\n", filePath)
		for _, m := range fileMethods {
			if m.IsPrivate {
				privateCount++
				fmt.Printf("  ⊗ %s.%s() (private - no export)\n", m.ReceiverType, m.Name)
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

				fmt.Printf("%s %s.%s(", icon, m.ReceiverType, m.Name)
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
		fmt.Println()
	}

	fmt.Printf("📊 Summary:\n")
	fmt.Printf("   Public: %d | Private: %d | View: %d | Payable: %d\n",
		publicCount, privateCount, viewCount, payableCount)
}

// generateCode generates the unified exports file with state management
func generateCode(methods []*MethodInfo, stateStructs []*StateInfo) string {
	var sb strings.Builder

	// Header
	sb.WriteString("// Code generated by NEAR contract generator. DO NOT EDIT.\n")
	sb.WriteString("// This file includes automatic state management.\n\n")
	sb.WriteString("package main\n\n")
	sb.WriteString("import (\n")
	sb.WriteString("\tcontractBuilder \"github.com/vlmoon99/near-sdk-go/contract\"\n")
	sb.WriteString("\t\"github.com/vlmoon99/near-sdk-go/env\"\n")
	sb.WriteString("\tnearjson \"github.com/vlmoon99/near-sdk-go/json\"\n")
	sb.WriteString(")\n\n")

	// Generate state read/write functions for each state struct
	for _, state := range stateStructs {
		sb.WriteString(generateStateReadFunction(state))
		sb.WriteString("\n")
		sb.WriteString(generateStateWriteFunction(state))
		sb.WriteString("\n")
	}

	// Generate exports for public methods
	for _, m := range methods {
		if !m.IsPublic {
			continue
		}

		exportName := toSnakeCase(m.Name)
		sb.WriteString(fmt.Sprintf("// From: %s\n", m.RelativePath))
		sb.WriteString(fmt.Sprintf("//go:export %s\n", exportName))
		sb.WriteString(fmt.Sprintf("func %s() {\n", exportName))
		sb.WriteString("\tcontractBuilder.HandleClientJSONInput(func(input *contractBuilder.ContractInput) error {\n")

		// Read state before method execution
		sb.WriteString("\t\t// Read contract state\n")
		sb.WriteString(fmt.Sprintf("\t\tcontract := GetContractState%s()\n\n", m.ReceiverType))

		// Add payment validation
		if m.IsPayable {
			sb.WriteString(fmt.Sprintf("\t\t// Validate payment: %s NEAR\n", m.MinDeposit))
			sb.WriteString(fmt.Sprintf("\t\tif !validatePayment(\"%s\") {\n", m.MinDeposit))
			sb.WriteString("\t\t\tenv.PanicStr(\"Insufficient payment attached\")\n")
			sb.WriteString("\t\t}\n\n")
		}

		// Parse parameters
		for _, p := range m.Params {
			sb.WriteString(generateParamParser(p))
		}

		// Call the actual method
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

		// Write state after mutating methods
		if m.IsMutating {
			sb.WriteString("\t\t// Save state\n")
			sb.WriteString(fmt.Sprintf("\t\tSetContractState%s(contract)\n\n", m.ReceiverType))
		}

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
	sb.WriteString("// validatePayment checks if sufficient NEAR is attached\n")
	sb.WriteString("func validatePayment(minDeposit string) bool {\n")
	sb.WriteString("\t// TODO: Implement payment validation\n")
	sb.WriteString("\t// attached := env.AttachedDeposit()\n")
	sb.WriteString("\t// required := parseNEAR(minDeposit)\n")
	sb.WriteString("\t// return attached >= required\n")
	sb.WriteString("\treturn true\n")
	sb.WriteString("}\n")

	return sb.String()
}

// generateStateReadFunction generates GetContractState for a struct
func generateStateReadFunction(state *StateInfo) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("// GetContractState%s reads state from blockchain\n", state.Name))
	sb.WriteString(fmt.Sprintf("func GetContractState%s() *%s {\n", state.Name, state.Name))
	sb.WriteString("\treadBytes, err := env.StateRead()\n")
	sb.WriteString("\tif err != nil {\n")
	sb.WriteString("\t\tenv.PanicStr(\"Failed to read state\")\n")
	sb.WriteString("\t}\n\n")

	sb.WriteString("\tif len(readBytes) == 0 {\n")
	sb.WriteString(fmt.Sprintf("\t\treturn &%s{}\n", state.Name))
	sb.WriteString("\t}\n\n")

	sb.WriteString("\tp := nearjson.NewParser(readBytes)\n\n")

	// Parse each field
	for _, field := range state.Fields {
		sb.WriteString(generateFieldParser(field))
	}

	sb.WriteString(fmt.Sprintf("\n\treturn &%s{\n", state.Name))
	for _, field := range state.Fields {
		sb.WriteString(fmt.Sprintf("\t\t%s: %s,\n", field.Name, field.Name))
	}
	sb.WriteString("\t}\n")
	sb.WriteString("}\n")

	return sb.String()
}

// generateStateWriteFunction generates SetContractState for a struct
func generateStateWriteFunction(state *StateInfo) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("// SetContractState%s writes state to blockchain\n", state.Name))
	sb.WriteString(fmt.Sprintf("func SetContractState%s(contract *%s) {\n", state.Name, state.Name))
	sb.WriteString("\tb := nearjson.NewBuilder()\n\n")

	// Add each field
	for _, field := range state.Fields {
		sb.WriteString(generateFieldSerializer(field))
	}

	sb.WriteString("\n\tstateBytes := b.Build()\n\n")
	sb.WriteString("\terr := env.StateWrite(stateBytes)\n")
	sb.WriteString("\tif err != nil {\n")
	sb.WriteString("\t\tenv.PanicStr(\"Failed to write state\")\n")
	sb.WriteString("\t}\n")
	sb.WriteString("}\n")

	return sb.String()
}

// generateFieldParser generates parsing code for a field
func generateFieldParser(field FieldInfo) string {
	var sb strings.Builder

	switch field.Type {
	case "int":
		sb.WriteString(fmt.Sprintf("\t%sValue, err := p.GetInt(\"%s\")\n", field.Name, field.Name))
		sb.WriteString("\tif err != nil {\n")
		sb.WriteString("\t\tenv.PanicStr(\"Failed to parse state\")\n")
		sb.WriteString("\t}\n")
		sb.WriteString(fmt.Sprintf("\t%s := int(%sValue)\n", field.Name, field.Name))

	case "string":
		sb.WriteString(fmt.Sprintf("\t%s, err := p.GetString(\"%s\")\n", field.Name, field.Name))
		sb.WriteString("\tif err != nil {\n")
		sb.WriteString("\t\tenv.PanicStr(\"Failed to parse state\")\n")
		sb.WriteString("\t}\n")

	case "bool":
		sb.WriteString(fmt.Sprintf("\t%s, err := p.GetBool(\"%s\")\n", field.Name, field.Name))
		sb.WriteString("\tif err != nil {\n")
		sb.WriteString("\t\tenv.PanicStr(\"Failed to parse state\")\n")
		sb.WriteString("\t}\n")

	case "uint64":
		sb.WriteString(fmt.Sprintf("\t%sValue, err := p.GetInt(\"%s\")\n", field.Name, field.Name))
		sb.WriteString("\tif err != nil {\n")
		sb.WriteString("\t\tenv.PanicStr(\"Failed to parse state\")\n")
		sb.WriteString("\t}\n")
		sb.WriteString(fmt.Sprintf("\t%s := uint64(%sValue)\n", field.Name, field.Name))

	default:
		sb.WriteString(fmt.Sprintf("\t// TODO: Complex type %s (%s)\n", field.Name, field.Type))
		sb.WriteString(fmt.Sprintf("\tvar %s %s\n", field.Name, field.Type))
	}

	return sb.String()
}

// generateFieldSerializer generates serialization code for a field
func generateFieldSerializer(field FieldInfo) string {
	switch field.Type {
	case "int", "int64":
		return fmt.Sprintf("\tb.AddInt(\"%s\", int(contract.%s))\n", field.Name, field.Name)
	case "string":
		return fmt.Sprintf("\tb.AddString(\"%s\", contract.%s)\n", field.Name, field.Name)
	case "bool":
		return fmt.Sprintf("\tb.AddBool(\"%s\", contract.%s)\n", field.Name, field.Name)
	case "uint64":
		return fmt.Sprintf("\tb.AddInt(\"%s\", int(contract.%s))\n", field.Name, field.Name)
	default:
		return fmt.Sprintf("\t// TODO: Complex type %s (%s)\n", field.Name, field.Type)
	}
}

func generateParamParser(p Param) string {
	var sb strings.Builder

	switch p.Type {
	case "string":
		sb.WriteString(fmt.Sprintf("\t\t%s, err := input.JSON.GetString(\"%s\")\n", p.Name, p.Name))
		sb.WriteString("\t\tif err != nil {\n")
		sb.WriteString("\t\t\tenv.PanicStr(\"Failed to parse parameter\")\n")
		sb.WriteString("\t\t}\n\n")

	case "int":
		sb.WriteString(fmt.Sprintf("\t\t%sInt64, err := input.JSON.GetInt(\"%s\")\n", p.Name, p.Name))
		sb.WriteString("\t\tif err != nil {\n")
		sb.WriteString("\t\t\tenv.PanicStr(\"Failed to parse parameter\")\n")
		sb.WriteString("\t\t}\n")
		sb.WriteString(fmt.Sprintf("\t\t%s := int(%sInt64)\n\n", p.Name, p.Name))

	case "bool":
		sb.WriteString(fmt.Sprintf("\t\t%s, err := input.JSON.GetBool(\"%s\")\n", p.Name, p.Name))
		sb.WriteString("\t\tif err != nil {\n")
		sb.WriteString("\t\t\tenv.PanicStr(\"Failed to parse parameter\")\n")
		sb.WriteString("\t\t}\n\n")

	case "uint64":
		sb.WriteString(fmt.Sprintf("\t\t%sInt64, err := input.JSON.GetInt(\"%s\")\n", p.Name, p.Name))
		sb.WriteString("\t\tif err != nil {\n")
		sb.WriteString("\t\t\tenv.PanicStr(\"Failed to parse parameter\")\n")
		sb.WriteString("\t\t}\n")
		sb.WriteString(fmt.Sprintf("\t\t%s := uint64(%sInt64)\n\n", p.Name, p.Name))

	default:
		sb.WriteString(fmt.Sprintf("\t\t// TODO: Parse complex type %s\n", p.Type))
		sb.WriteString(fmt.Sprintf("\t\tvar %s %s\n\n", p.Name, p.Type))
	}

	return sb.String()
}

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