package htmlplus

import (
	"golang.org/x/net/html"

	"lib-builtin/lib/slices"

	"strings"
	"errors"
	"bytes"
	"strconv"
)

type NodeNotFoundError struct {
	mWhy string
}

func newNodeNotFoundError(pWhyPreMessageWithOptionalFormatting ...interface{}) *NodeNotFoundError {
	return &NodeNotFoundError{mWhy:slices.Format(pWhyPreMessageWithOptionalFormatting...)}
}

func (this *NodeNotFoundError) Error() string {
	return this.mWhy
}

type matcher func(pNode *html.Node) (ok bool, err error)

type Document struct {
	mDoc *html.Node
}

func NewDocument(pBody []byte) (rDocument *Document, err error) {
	zDoc, err := html.Parse(bytes.NewBuffer(pBody))
	if err == nil {
		rDocument = &Document{mDoc:zDoc}
	}
	return
}

func (this *Document) GetTableWithId(pId string) (*Table, error) {
	zTable := NewTable()
	return this.getTableAtOrAfterMatchedNode(zTable,
		func(pNode *html.Node) (ok bool, err error) {
			if pNode == nil {
				err = newNodeNotFoundError("no node found w/ id '%s'", pId)
				return
			}
			zId, ok := getAttributeValue(pNode, "id")
			ok = ok && (zId == pId)
			if ok {
				zTable.SetId(zId)
			}
			return
		})
}

func (this *Document) GetTableWithIdStartingWith(pIdPrefix string) (*Table, error) {
	zTable := NewTable()
	return this.getTableAtOrAfterMatchedNode(zTable,
		func(pNode *html.Node) (ok bool, err error) {
			if pNode == nil {
				err = newNodeNotFoundError("no node found w/ id that starts with '%s'", pIdPrefix)
				return
			}
			zId, ok := getAttributeValue(pNode, "id")
			ok = ok && strings.HasPrefix(zId, pIdPrefix)
			if ok {
				zTable.SetId(zId)
			}
			return
		})
}

func (this *Document) getTableAtOrAfterMatchedNode(pTable *Table, pMatcher matcher) (rTable *Table, err error) {
	zLocator := newNodeLocator(pMatcher, createFindNodeNamedAtOrAfterNodeWithId(pTable, "table"))
	zNode, err := zLocator.find(this.mDoc)
	if err == nil {
		err = pTable.Populate(zNode)
		if err == nil {
			rTable = pTable
		}
	}
	return
}

func getIntAttributeValue(pNode *html.Node, pKey string, pDefault int) int {
	zValue, ok := getAttributeValue(pNode, pKey)
	if ok {
		zValue = strings.TrimSpace(zValue)
		if len(zValue) != 0 {
			i, err := strconv.Atoi(zValue)
			if (i > 0) && (err == nil) {
				return i
			}
		}
	}
	return pDefault
}

func getAttributeValue(pNode *html.Node, pKey string) (rValue string, ok bool) {
	if pNode.Type == html.ElementNode {
		for _, zAttr := range pNode.Attr {
			if zAttr.Key == pKey {
				rValue = zAttr.Val
				ok = true
				return
			}
		}
	}
	return
}

type IdAccessor interface {
	GetId() string
}

func createFindNodeNamedAtOrAfterNodeWithId(pIdAccessor IdAccessor, pName string) matcher {
	zMatcherFuncStruct := &nodeTypeAtOrAfterNodeWithIdMatcher{mIdAccessor:pIdAccessor, mName:pName}
	return zMatcherFuncStruct.match
}

type nodeTypeAtOrAfterNodeWithIdMatcher struct {
	mIdAccessor              IdAccessor
	mName                    string
	mElementNodesEncountered int
}

func (this *nodeTypeAtOrAfterNodeWithIdMatcher) match(pNode *html.Node) (ok bool, err error) {
	if pNode == nil {
		err = newNodeNotFoundError("no '%s' node found after node w/ id '%s'", this.mName, this.mIdAccessor.GetId())
		return
	}
	zName, err := this.checkForIdOnElements(pNode)
	ok = (this.mName == zName) && (err == nil)
	return
}

func (this *nodeTypeAtOrAfterNodeWithIdMatcher) checkForIdOnElements(pNode *html.Node) (rName string, err error) {
	if pNode.Type != html.ElementNode {
		return
	}
	z1stNode := (this.mElementNodesEncountered == 0)
	this.mElementNodesEncountered++

	rName = pNode.Data // For Elements, the Data is the Name

	// We "assume" that the 1st Node has an "id" and is OK!
	if !z1stNode {
		_, zHasId := getAttributeValue(pNode, "id")
		if zHasId {
			err = newNodeNotFoundError("'%s' node w/ an id encountered while searching for a '%s' node (after node w/ id '%s')",
				rName, this.mName, this.mIdAccessor.GetId())
		}
	}
	return
}

type nodeLocator struct {
	mMatchers []matcher
	//mRecursiveFunc func(pNode *html.Node) (rFoundNode *html.Node, err error)
}

func newNodeLocator(pMatchers ...matcher) (rLocator *nodeLocator) {
	rLocator = &nodeLocator{mMatchers:pMatchers}
	//rLocator.mRecursiveFunc = rLocator.find
	return
}

func (this *nodeLocator) find(pNode *html.Node) (rNode *html.Node, err error) {
	rNode, err = this.recursiveFind(pNode)
	if (rNode == nil) && (err == nil) {
		_, err = this.match(rNode) // Force Appropriate Error
	}
	return
}

func (this *nodeLocator) recursiveFind(pNode *html.Node) (rNode *html.Node, err error) {
	ok, err := this.match(pNode)
	if ok || (err != nil) {
		rNode = pNode // Note: returning !nil even when there is an Error!
		return
	}
	for zNode := pNode.FirstChild; zNode != nil; zNode = zNode.NextSibling {
		rNode, err = this.recursiveFind(zNode)
		if (rNode != nil) || (err != nil) {
			return
		}
	}
	return nil, nil // 4 Clarity
}

func (this *nodeLocator) match(pNode *html.Node) (ok bool, err error) {
	if len(this.mMatchers) == 0 {
		err = errors.New("locator w/ NO matchers")
		return
	}
	for len(this.mMatchers) != 0 {
		ok, err = this.mMatchers[0](pNode)
		if !ok || (err != nil) {
			return
		}
		this.mMatchers = this.mMatchers[1:] // 1st OK, so drop it and move on to the next
	}
	return // ok true!
}
