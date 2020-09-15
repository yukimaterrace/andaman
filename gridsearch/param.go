package gridsearch

type (
	// Param is an interface for param
	Param interface {
		PartialParamSetSize() int
		PartialParamSetValueSize(partialParamSetID PartialParamSetID) int
		GetPartialParamSet(partialParamSetID PartialParamSetID, partialParamSetValueID PartialParamSetValueID) PartialParamSet
		CreateParam(partialParamSet ...PartialParamSet) Param
	}

	// PartialParamSet is an inteface for partial param set
	PartialParamSet interface{}

	// PartialParamSetID is a definition for partial param set ID
	PartialParamSetID int

	// PartialParamSetValueID is a definition for partioal param set value ID
	PartialParamSetValueID int
)

type partialParamSetsIterator struct {
	param             Param
	setID             PartialParamSetID
	tailIterator      *partialParamSetsIterator
	currentSetValueID PartialParamSetValueID
	tailSets          []PartialParamSet
}

func newPartialParamSetsIterator(param Param) *partialParamSetsIterator {
	if param.PartialParamSetSize() == 0 {
		panic("invalid param")
	}

	return _newPartialParamSetsIterator(param, 0)
}

func _newPartialParamSetsIterator(param Param, partialParamSetID PartialParamSetID) *partialParamSetsIterator {
	if int(partialParamSetID) >= param.PartialParamSetSize() {
		return nil
	}

	return &partialParamSetsIterator{
		param:             param,
		setID:             partialParamSetID,
		tailIterator:      _newPartialParamSetsIterator(param, partialParamSetID+1),
		currentSetValueID: 0,
		tailSets:          nil,
	}
}

func (iterator *partialParamSetsIterator) hasNext() bool {
	cond1 := int(iterator.currentSetValueID) < iterator.param.PartialParamSetValueSize(iterator.setID)-1
	cond2 := false
	if iterator.tailIterator != nil {
		cond2 = iterator.tailIterator.hasNext()
	}
	return cond1 || cond2
}

func (iterator *partialParamSetsIterator) next() []PartialParamSet {
	if iterator.tailIterator != nil && iterator.tailIterator.hasNext() {
		tailSets := iterator.tailIterator.next()
		iterator.tailSets = tailSets
		set := iterator.param.GetPartialParamSet(iterator.setID, iterator.currentSetValueID)
		return append([]PartialParamSet{set}, tailSets...)
	}

	iterator.currentSetValueID++
	set := iterator.param.GetPartialParamSet(iterator.setID, iterator.currentSetValueID)
	if iterator.tailSets != nil {
		return append([]PartialParamSet{set}, iterator.tailSets...)
	}

	return []PartialParamSet{set}
}
