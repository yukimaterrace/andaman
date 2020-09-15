package gridsearch

type (
	param interface {
		partialParamSetSize() int
		partialParamSetValueSize(setID partialParamSetID) int
		getPartialParamSet(setID partialParamSetID, setValueID partialParamSetValueID) partialParamSet
		createParam(paramSets []partialParamSet) param
	}

	partialParamSet interface{}

	partialParamSetID int

	partialParamSetValueID int
)

type partialParamSetsIterator struct {
	param             param
	setID             partialParamSetID
	tailIterator      *partialParamSetsIterator
	currentSetValueID partialParamSetValueID
	tailSets          []partialParamSet
}

func newPartialParamSetsIterator(param param) *partialParamSetsIterator {
	if param.partialParamSetSize() == 0 {
		panic("invalid param")
	}

	return _newPartialParamSetsIterator(param, 0)
}

func _newPartialParamSetsIterator(param param, paramSetID partialParamSetID) *partialParamSetsIterator {
	if int(paramSetID) >= param.partialParamSetSize() {
		return nil
	}

	return &partialParamSetsIterator{
		param:             param,
		setID:             paramSetID,
		tailIterator:      _newPartialParamSetsIterator(param, paramSetID+1),
		currentSetValueID: 0,
		tailSets:          nil,
	}
}

func (iterator *partialParamSetsIterator) hasNext() bool {
	cond1 := int(iterator.currentSetValueID) < iterator.param.partialParamSetValueSize(iterator.setID)-1
	cond2 := false
	if iterator.tailIterator != nil {
		cond2 = iterator.tailIterator.hasNext()
	}
	return cond1 || cond2
}

func (iterator *partialParamSetsIterator) next() []partialParamSet {
	if iterator.tailIterator != nil && iterator.tailIterator.hasNext() {
		tailSets := iterator.tailIterator.next()
		iterator.tailSets = tailSets
		set := iterator.param.getPartialParamSet(iterator.setID, iterator.currentSetValueID)
		return append([]partialParamSet{set}, tailSets...)
	}

	iterator.currentSetValueID++
	set := iterator.param.getPartialParamSet(iterator.setID, iterator.currentSetValueID)
	if iterator.tailSets != nil {
		return append([]partialParamSet{set}, iterator.tailSets...)
	}

	return []partialParamSet{set}
}

func paramsForGridSearch(p param) []param {
	result := []param{}

	iterator := newPartialParamSetsIterator(p)
	for iterator.hasNext() {
		paramSets := iterator.next()
		_param := p.createParam(paramSets)
		result = append(result, _param)
	}

	return result
}
