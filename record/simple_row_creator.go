package record

type simpleRowCreator struct{}

func (creator *simpleRowCreator) header() []string {
	return nil
}

func (creator *simpleRowCreator) row(material *Material) []string {
	return nil
}
