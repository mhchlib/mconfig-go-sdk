package mconfig

const (
	// FieldType_String ...
	FieldType_String FieldType = iota
	// FieldType_Int ...
	FieldType_Int
	// FieldType_Bool ...
	FieldType_Bool
	// FieldType_Map ...
	FieldType_Map
	// FieldType_List ...
	FieldType_List
	// FieldType_Interface ...
	FieldType_Interface
)

// FieldInterface ...
type FieldInterface interface {
	isFieldInterface()
}

// FieldInterface_Int ...
type FieldInterface_Int struct {
	Value int64
}

// FieldInterface_String ...
type FieldInterface_String struct {
	Value string
}

// FieldInterface_Bool ...
type FieldInterface_Bool struct {
	Value bool
}

// FieldInterface_Map ...
type FieldInterface_Map struct {
	Value map[string]interface{}
}

// FieldInterface_List ...
type FieldInterface_List struct {
	Value []interface{}
}

// FieldInterface_Interface ...
type FieldInterface_Interface struct {
	Value interface{}
}

func (value FieldInterface_Int) isFieldInterface()       {}
func (value FieldInterface_String) isFieldInterface()    {}
func (value FieldInterface_Bool) isFieldInterface()      {}
func (value FieldInterface_Map) isFieldInterface()       {}
func (value FieldInterface_List) isFieldInterface()      {}
func (value FieldInterface_Interface) isFieldInterface() {}
