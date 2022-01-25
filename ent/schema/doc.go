package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

// Doc schema
type Doc struct {
	ent.Schema
}

func (Doc) Fields() []ent.Field {
	return []ent.Field{
		field.Uint32("id").Positive().Unique(),
		field.Bytes("wordInts"),
		field.String("inputDocId"),
		field.String("userId"),
		field.String("businessId"),
		field.Float32("stars"),
		field.Int16("useful"),
		field.Int16("funny"),
		field.Int16("cool"),
		field.String("text"),
		field.String("date"),
	}
}
