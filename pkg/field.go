package pkg

const (
	FieldType_String FieldType = iota
	FieldType_Int
	FieldType_Bool
	FieldType_Map
	FieldType_List
)

type FieldInterface interface {
	isFieldInterface()
}

type FieldInterface_Int struct {
	Value int64
}

type FieldInterface_String struct {
	Value string
}

type FieldInterface_Bool struct {
	Value bool
}

type FieldInterface_Map struct {
	Value map[string]interface{}
}

type FieldInterface_List struct {
	Value []interface{}
}

func (value FieldInterface_Int) isFieldInterface()    {}
func (value FieldInterface_String) isFieldInterface() {}
func (value FieldInterface_Bool) isFieldInterface()   {}
func (value FieldInterface_Map) isFieldInterface()    {}
func (value FieldInterface_List) isFieldInterface()   {}
