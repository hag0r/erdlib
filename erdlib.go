package main

import (
	"fmt"

	"github.com/alecthomas/participle"
	"github.com/alecthomas/participle/lexer"
)

// ERDExpression is the entry point for parsing textual ERDs
type ERDExpression struct {
	// Name is the relationship's name
	Name string `@Ident`
	// Attributes is the list of attributes defined for the relation
	EntitiesAndAttributes []Ref `"("@@ ("," @@)* ")"`
}

// IsRelationship tests if the parsed object represents a relationship
// If yes, returns true. Otherwise, if the object is a EntityType, returns false
func (e ERDExpression) IsRelationship() bool {
	for _, r := range e.EntitiesAndAttributes {
		if r.IsEntityRef() {
			return true
		}
	}
	return false
}

// IsEntityType returns true, if the object is a Entity-type definition
// Otherwise false
func (e ERDExpression) IsEntityType() bool {
	return !e.IsRelationship()
}

// Ref represents either an reference to an (existing)
// Entity or is an ordinary attribute
type Ref struct {
	RefName     Attribute `@@`
	Cardinality *MinMax   `("[" @@ "]")?`
}

// IsEntityRef returns true if the element is a reference to an attribute
// if the reference is an attribute, this is false
func (r Ref) IsEntityRef() bool {
	return r.Cardinality != nil
}

// Name returns the name of the reference/attribute
func (r Ref) Name() string {
	return r.RefName.Name()
}

// MinMax represents the cardinality informatin given in [min,max] notation
type MinMax struct {
	Min string `(@CardinalityNum | @"*") ","`
	Max string `(@CardinalityNum | @"*")`
}

// Attribute defines an attribute in a relation
type Attribute struct {
	// PK holds the field name if it is a primary key, otherwise nil
	PK *string `("_" @Ident "_"`
	// AttrName holds the field name if it is not a primary key, otherwise nil
	AttrName *string `| @Ident )`
}

// Name returns the attribute's name. Use this instead of direct access to AttrName or PK, which might be nil
func (attr Attribute) Name() string {
	if attr.IsPK() {
		return *(attr.PK)
	}

	return *(attr.AttrName)
}

// IsPK tests if the attribute is part of the primary key in its relation
func (attr Attribute) IsPK() bool {
	return attr.PK != nil
}

// ERDParser encapsules the parser generator and provides
// convenient methods to parse a string into a relation definition
type ERDParser struct {
	parser *participle.Parser
}

// CreateParser creates a new parser instance to parse relation definitions
func CreateParser() (*ERDParser, error) {

	rmLexer := lexer.Must(lexer.Regexp(`(\s+)` +
		`|(?P<Ident>[a-zA-Z][a-zA-Z0-9]*)` +
		`|(?P<Operators>->|[,()_\[\]\*])` +
		`|(?P<CardinalityNum>[0-9]+)`))

	parser, err := participle.Build(&ERDExpression{}, participle.Lexer(rmLexer))

	if err != nil {
		return nil, err
	}

	relParser := &ERDParser{parser}
	return relParser, nil
}

// Parse takes a string and tries to parse it into a RelationDefinition
// The given string should be a single line defining one relation and its attributes, like
// R(a, b, _c_,_d_, e -> R2(k))
func (parser *ERDParser) Parse(s string) (*ERDExpression, error) {
	erdElem := &ERDExpression{}

	err := parser.parser.ParseString(s, erdElem)

	if err != nil {
		return nil, err
	}

	return erdElem, nil
}

func main() {

	// s := "EntityType1(_a_,b,c)"
	s := "relationShip1(Entity1 [0,*], A[1,3], normalesattr, _pkattr_)"

	parser, err := CreateParser()

	if err != nil {
		fmt.Println("Failed to create parser", err)
		return
	}

	erdObject, err := parser.Parse(s)

	if err != nil {
		fmt.Println("got an error during parse", err)
		return
	}

	fmt.Printf("Successfully parsed relation: %v\n", erdObject)

	fmt.Println("Found elem: " + erdObject.Name)
	for _, attr := range erdObject.EntitiesAndAttributes {
		if attr.IsEntityRef() {
			fmt.Printf("--> %v is entity -> min: %v, max: %v\n", attr.Name(), attr.Cardinality.Min, attr.Cardinality.Max)
		} else {
			fmt.Printf("--> %v is attr. PK: %v \n", attr.Name(), attr.RefName.IsPK())
		}

	}

}
