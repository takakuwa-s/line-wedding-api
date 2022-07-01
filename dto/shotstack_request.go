package dto

type TemplateRender struct {
	Id    string       `json:"id"`
	Merge []MergeField `json:"merge"`
}

type MergeField struct {
	Find    string `json:"find"`
	Replace string `json:"replace"`
}

func NewTemplateRender(id string) *TemplateRender {
	return &TemplateRender{
		Id:    id,
		Merge: []MergeField{},
	}
}

func NewMergeField(find, replace string) MergeField {
	return MergeField{
		Find:    find,
		Replace: replace,
	}
}

func (tr *TemplateRender) ApendMerge(m MergeField) {
	tr.Merge = append(tr.Merge, m)
}
