package constants

type ContextKey string

const (
	SepOfTags                         = ","
	SoftDeleteFieldTagName            = "soft_delete_field"
	CreateTimeFieldTagName            = "create_time_field"
	UpdateTimeFieldTagName            = "update_time_field"
	FieldTagKey                       = "crud_generator"
	ModelKey               ContextKey = "CURD_model"
)
