# ERD Lib 
ERDlib - short for Entity-Relationship-Diagram library - allows to parse textual representations of database ERDs.

The library does not make any sanity checks of the parsed ERD text (PKs on relationships, only one entity type reference in a relationship, etc.).

## Syntax
The supported syntax is something like
```
Entity1 ( Attribute1, Attribute2, _PrimaryKey_, ... )
Entity2 ( _PrimaryKey_, Attribute1, Attribute2)

relationship(Entity1[0,*], Entity2 [1,4], anotherattribute)
```

## Dependencies

We use the great [`https://github.com/alecthomas/participle`](https://github.com/alecthomas/participle) parser library


## Limitations

- Attribute names **must not** contain underscores `_`.


## Example 

```golang
func main() {

	s := "relationShip1(Entity1 [0,*], Entity2[1,3], attribute1, _pkattr_)"

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
```
