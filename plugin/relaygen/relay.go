// Plugin to ease relay.go integration: https://relay.dev/docs/en/graphql-server-specification
package relaygen

import (
	"fmt"
	"github.com/99designs/gqlgen/codegen"
	"github.com/99designs/gqlgen/codegen/config"
	"github.com/99designs/gqlgen/plugin"
	"github.com/vektah/gqlparser/ast"
)

func New() plugin.Plugin {
	return &Plugin{}
}

type Plugin struct {
	Nodes []*Node
}

// Object that implements the Node type
type Node struct {
	Name string // The same name as the type declaration
	Def  *ast.Definition
}

type Connection struct {
	ForType string
	Type    string
}

var _ plugin.CodeGenerator = &Plugin{}

func (r *Plugin) Name() string {
	return "relay"
}

type RelayBuild struct {
	*codegen.Data

	PackageName string
	Nodes       []*Node
}

func (r *Plugin) GenerateCode(data *codegen.Data) error {
	//relayBuild := &RelayBuild{
	//	Data:         data,
	//	PackageName:  data.Config.Resolver.Package,
	//	Nodes: r.Nodes,
	//}

	//return templates.Render(templates.Options{
	//	PackageName: data.Config.Resolver.Package,
	//	Filename:    "resolve_relay.go",
	//	Data: &relayBuild,
	//	GeneratedHeader: true,
	//})

	return nil
}

func (r *Plugin) MutateConfig(cfg *config.Config) error {
	if err := cfg.Check(); err != nil {
		return err
	}

	_, err := cfg.LoadSchema()
	if err != nil {
		return err
	}

	cfg.Directives["relayID"] = config.DirectiveConfig{
		SkipRuntime: false,
	}

	cfg.Directives["connection"] = config.DirectiveConfig{
		SkipRuntime: true,
	}

	return nil
}

func (r *Plugin) MutateSchema(s *ast.Schema) error {
	//for _, ent := range r.Nodes {
	//	definition := s.Types[ent.Name]
	//
	//	definition.Interfaces = append(definition.Interfaces, "Node")
	//}

	return nil
}

func (r *Plugin) InjectSources(cfg *config.Config) {
	cfg.AdditionalSources = append(cfg.AdditionalSources, r.getSource(false))

	schema, err := cfg.LoadSchema()
	if err != nil {
		panic(err)
	}

	var connections []Connection

	for _, schemaType := range schema.Types {
		if schemaType.Kind == ast.Object {
			for _, field := range schemaType.Fields {
				if field.Name == "id" {
					// Match objects that satisfy node interface
					if field.Type.NamedType == "ID" && field.Type.NonNull {
						r.Nodes = append(r.Nodes, &Node{
							Name: schemaType.Name,
							Def:  schemaType,
						})
					}
				}

				directive := field.Directives.ForName("connection")

				if directive != nil {
					connections = append(connections, Connection{
						Type: field.Type.Name(),
						// TODO: validate
						ForType: directive.Arguments.ForName("for").Value.Raw,
					})
				}
			}
		}
	}

	s := ""
	for _, connection := range connections {
		edgeType := fmt.Sprintf("%sEdge", connection.ForType)

		s += fmt.Sprintf("type %s {\n\tcursor: String!\n\tnode: %s\n}\n", edgeType, connection.ForType)
		s += fmt.Sprintf("extend type %s {\n\tedges: [%s]\n\tpageInfo: PageInfo!\n}\n", connection.Type, edgeType)
	}

	cfg.AdditionalSources = append(cfg.AdditionalSources, &ast.Source{Name: "connections.graphql", Input: s, BuiltIn: true})
}

func (r *Plugin) getSource(builtin bool) *ast.Source {
	return &ast.Source{
		Name: "relay.graphql",
		Input: `# Declarations as required by the relay spec 
# See: https://relay.dev/docs/en/graphql-server-specification
interface Node {
  id: ID!
}

type PageInfo {
  hasNextPage: Boolean!
  hasPreviousPage: Boolean!
  startCursor: String
  endCursor: String
}

extend type Query {
  node(id: ID!): Node
}

directive @relayID(types: [String!]) on FIELD_DEFINITION
directive @connection(for: String!) on FIELD_DEFINITION
`,
		BuiltIn: builtin,
	}
}
