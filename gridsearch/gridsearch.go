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
	currentSetValueID partialParamSetValueID
	tailIterator      *partialParamSetsIterator
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
		currentSetValueID: -1,
		tailIterator:      _newPartialParamSetsIterator(param, paramSetID+1),
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
		if iterator.currentSetValueID < 0 {
			iterator.currentSetValueID++
		}

		set := iterator.param.getPartialParamSet(iterator.setID, iterator.currentSetValueID)
		tailSets := iterator.tailIterator.next()
		return append([]partialParamSet{set}, tailSets...)
	}

	iterator.tailIterator = _newPartialParamSetsIterator(iterator.param, iterator.setID+1)

	iterator.currentSetValueID++
	set := iterator.param.getPartialParamSet(iterator.setID, iterator.currentSetValueID)

	if iterator.tailIterator != nil && iterator.tailIterator.hasNext() {
		tailSets := iterator.tailIterator.next()
		return append([]partialParamSet{set}, tailSets...)
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
