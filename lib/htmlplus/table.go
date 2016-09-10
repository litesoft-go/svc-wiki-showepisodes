package htmlplus

import (
	"golang.org/x/net/html"

	"lib-builtin/lib/slices"
	"lib-builtin/lib/lines"

	"errors"
	"fmt"
	"lib-builtin/lib/ints"
	"strings"
)

type HeaderRow []string

func (this HeaderRow) String() string {
	return slices.AsOptions([]string(this))
}

func (this HeaderRow) Equals(them HeaderRow) bool {
	return !slices.Equals([]string(this), []string(them)...)
}

func addHeaders(pCollector *lines.Collector, pWhat string, pHeaderRows []HeaderRow) {
	pCollector.Indent()
	pCollector.Line(pWhat)
	for _, zRow := range pHeaderRows {
		pCollector.Indent()
		pCollector.Line(zRow.String())
		pCollector.Outdent()
	}
	pCollector.Outdent()
}

type Table struct {
	mIdForTable string
	mCaption    string
	mHeaderRows []HeaderRow
	mBodyRows   []*Row
	mFooterRows []*Row // should probably be 'HeaderRow'
}

func NewTable() *Table {
	return &Table{}
}

func (this *Table) GetId() string {
	return this.mIdForTable
}

func (this *Table) SetId(pIdForTable string) {
	this.mIdForTable = pIdForTable
}

func (this *Table) HeaderMatches(pHeaderRows []HeaderRow) bool {
	if len(this.mHeaderRows) != len(pHeaderRows) {
		return false
	}
	for i, zRow := range pHeaderRows {
		if this.mHeaderRows[i].Equals(zRow) {
			return false
		}
	}
	return true
}

func (this *Table) ErrorHeaderNotMatched() error {
	zCollector := lines.NewCollector()
	zCollector.Line("headers don't match any option for Table: " + this.mIdForTable)
	addHeaders(zCollector, "Actual:", this.mHeaderRows)
	return errors.New(zCollector.String())
}

func (this *Table) Populate(pNode *html.Node) (err error) {
	if pNode != nil {
		zPTable := &populationTable{}
		err = zPTable.populate(pNode, InferredRowType)
		if err == nil {
			this.mCaption = zPTable.mCaption
			this.mHeaderRows, this.mBodyRows, this.mFooterRows = zPTable.splitRows()
		}
	}
	return
}

type populationTable struct {
	mCaption             string
	mRowsExplicitlyTyped bool
	mRows                []*Row
}

func (this *populationTable) splitRows() (rHeaderRows []HeaderRow, rBodyRows []*Row, rFooterRows []*Row) {
	if len(this.mRows) == 0 {
		return
	}
	if this.mRowsExplicitlyTyped {
		this.typeInferredRows()
	} else {
		this.typeAllRows()
	}
	rFooterRows = this.extractTypeRows(FootRowType)
	rBodyRows = this.extractTypeRows(BodyRowType)
	rHeaderRows = convertHeaderRows(this.extractTypeRows(HeaderRowType))
	return
}

func (this *populationTable) typeInferredRows() {
	for _, zRow := range this.mRows {
		if zRow.mRowType == InferredRowType {
			zRow.mRowType = BodyRowType
		}
	}
}

func (this *populationTable) typeAllRows() {
	this.convertTrailingRowsToFooterRowTypes() // do this first so that if the table is ALL THs then they will all end up as Header Rows
	this.convertLeadingRowsToHeaderRowTypes()
	this.typeInferredRows()
}

func (this *populationTable) convertLeadingRowsToHeaderRowTypes() {
	for _, zRow := range this.mRows {
		if !zRow.isAllTH() {
			return
		}
		zRow.mRowType = HeaderRowType
	}
}

func (this *populationTable) convertTrailingRowsToFooterRowTypes() {
	for i := len(this.mRows) - 1; 0 <= i; i-- {
		zRow := this.mRows[i]
		if !zRow.isAllTH() {
			return
		}
		zRow.mRowType = FootRowType
	}
}

func (this *populationTable) extractTypeRows(pRowType RowType) (rRows []*Row) {
	for i := 0; i < len(this.mRows); {
		zRow := this.mRows[i]
		if zRow.mRowType != pRowType {
			i++
		} else {
			rRows = append(rRows, zRow)
			this.mRows = Remove(this.mRows, i)
		}
	}
	// fmt.Println("type:", pRowType, "| count:", len(rRows))
	return
}

type RowType int

const (
	InferredRowType RowType = iota
	HeaderRowType
	BodyRowType
	FootRowType
)

func (this *populationTable) populate(pNode *html.Node, pRowType RowType) error {
	for zRowNode := pNode.FirstChild; zRowNode != nil; zRowNode = zRowNode.NextSibling {
		zRow, err := this.parseRow(zRowNode, pRowType)
		if err != nil {
			return err
		}
		if zRow != nil {
			this.mRowsExplicitlyTyped = (this.mRowsExplicitlyTyped || (zRow.mRowType != InferredRowType))
			this.mRows = append(this.mRows, zRow)
		}
	}
	return nil
}

func (this *populationTable) parseRow(pNode *html.Node, pRowType RowType) (rRow *Row, err error) {
	if pNode.Type != html.ElementNode {
		return
	}
	if pNode.Data == "colgroup" {
		return
	}
	if pNode.Data == "caption" {
		this.mCaption = extractAllText(pNode)
		return
	}
	if pNode.Data == "thead" {
		err = this.populate(pNode, HeaderRowType)
		return
	}
	if pNode.Data == "tfoot" {
		err = this.populate(pNode, FootRowType)
		return
	}
	if pNode.Data == "tbody" {
		err = this.populate(pNode, pRowType) // because the tbody may have been virtually added, we leave the type as Inferred
		return
	}
	if pNode.Data != "tr" {
		err = errors.New("unexpected element, expected 'tr', but got: " + pNode.Data)
		return
	}
	rRow = &Row{mRowType:pRowType}
	err = rRow.parseRow(pNode)
	return
}

func Remove(pRows []*Row, pIndexToRemove int) (rRows []*Row) {
	zLastIndex := len(pRows) - 1
	if (pIndexToRemove < 0) || (zLastIndex < pIndexToRemove) {
		panic(fmt.Sprintf("index (%d) out of bounds for Rows w/ length: %d", pIndexToRemove, len(pRows)))
	}
	if pIndexToRemove == 0 {
		if zLastIndex != 0 {
			rRows = pRows[1:]
		}
		return
	}
	if pIndexToRemove == zLastIndex {
		return pRows[:zLastIndex]
	}
	return append(pRows[:pIndexToRemove], pRows[(pIndexToRemove + 1):]...)
}

type Row struct {
	mRowType RowType
	mCells   []*Cell
}

func (this *Row) isAllTH() bool {
	for _, zCell := range this.mCells {
		if !zCell.mHeader {
			return false
		}
	}
	return true
}

func (this *Row) parseRow(pNode *html.Node) error {
	for zCellNode := pNode.FirstChild; zCellNode != nil; zCellNode = zCellNode.NextSibling {
		zCell, err := this.parseCell(zCellNode)
		if err != nil {
			return err
		}
		if zCell != nil {
			this.mCells = append(this.mCells, zCell)
		}
	}
	return nil
}

func (this *Row) parseCell(pNode *html.Node) (rCell *Cell, err error) {
	if pNode.Type != html.ElementNode {
		return
	}
	if pNode.Data == "th" {
		rCell = &Cell{mHeader:true}
		err = rCell.parseCell(pNode)
		return
	}
	if pNode.Data != "tr" {
		rCell = &Cell{mHeader:false}
		err = rCell.parseCell(pNode)
		return
	}
	err = errors.New("unexpected element, expected 'th' / 'td', but got: " + pNode.Data)
	return
}

type Cell struct {
	mHeader  bool
	mColspan int
	mRowspan int
	mText    string
}

// Extract the spans, and then get only the text from all the children using "|||" as separator.
func (this *Cell) parseCell(pNode *html.Node) error {
	this.mColspan = getIntAttributeValue(pNode, "colspan", 1)
	this.mRowspan = getIntAttributeValue(pNode, "rowspan", 1)
	this.mText = extractAllText(pNode)
	return nil
}

// Get only the text from all the children using "|||" as separator.
func extractAllText(pNode *html.Node) string {
	zCollector := &textCollector{}
	zCollector.from(pNode)
	return zCollector.mText
}

type textCollector struct {
	mText string
}

// Get only the text from all the children using "|||" as separator.
func (this *textCollector) from(pNode *html.Node) {
	this.addText(pNode)
	for zSubNode := pNode.FirstChild; zSubNode != nil; zSubNode = zSubNode.NextSibling {
		this.from(zSubNode)
	}
}

// Add the text (if the node is a TextNode) using "|||" as separator.
func (this *textCollector) addText(pNode *html.Node) {
	if pNode.Type == html.TextNode {
		zText := strings.Trim(pNode.Data, whitespace)
		if len(zText) != 0 {
			if len(this.mText) != 0 {
				this.mText += "|||"
			}
			this.mText += zText
		}
	}
}

func convertHeaderRows(pRows []*Row) (rHeaderRows []HeaderRow) {
	if len(pRows) != 0 {
		zHeaders := &headers{}
		for zRowIndex, zRow := range pRows {
			for _, zCell := range zRow.mCells {
				zHeaders.addCell(zRowIndex, zCell)
			}
		}
		rHeaderRows = zHeaders.makeGridRectangular().fillHoles().convert()
	}
	return
}

type headers struct {
	mRows []*headerRow
}

func (this *headers) makeGridRectangular() *headers {
	return this.padAllRowsTo(this.calculateMaxColumns())
}

func (this *headers) calculateMaxColumns() (rMaxCells int) {
	for _, zRow := range this.mRows {
		rMaxCells = ints.Max(rMaxCells, len(zRow.mCells))
	}
	return
}

func (this *headers) padAllRowsTo(pLength int) *headers {
	for _, zRow := range this.mRows {
		zRow.padTo(pLength)
	}
	return this
}

func (this *headers) fillHoles() *headers {
	for _, zRow := range this.mRows {
		zRow.fillHoles()
	}
	return this
}

func (this *headers) convert() []HeaderRow {
	zHeaderRows := make([]HeaderRow, len(this.mRows))
	for i, zRow := range this.mRows {
		zHeaderRows[i] = zRow.convert()
	}
	return zHeaderRows
}

func (this *headers) addCell(pRowIndex int, pCell *Cell) {
	for zRowspan := 0; zRowspan < pCell.mRowspan; zRowspan++ {
		this.getRow(pRowIndex + zRowspan).addCell(pCell)
	}
}

func (this *headers) getRow(pRowIndex int) *headerRow {
	for len(this.mRows) <= pRowIndex {
		this.mRows = append(this.mRows, &headerRow{})
	}
	return this.mRows[pRowIndex]
}

type headerRow struct {
	mCells []*headerCell
}

func (this *headerRow) padTo(pLength int) {
	for pLength < len(this.mCells) {
		this.mCells = append(this.mCells, &headerCell{})
	}
}

func (this *headerRow) fillHoles() {
	for _, zCell := range this.mCells {
		zCell.fillHole()
	}
}

func (this *headerRow) addCell(pCell *Cell) {
	for zColspan := 0; zColspan < pCell.mColspan; zColspan++ {
		this.findUnsetCell().setCell(pCell)
	}
}

func (this *headerRow) findUnsetCell() *headerCell {
	for _, zCell := range this.mCells {
		if !zCell.mSet {
			return zCell
		}
	}
	zCell := &headerCell{}
	this.mCells = append(this.mCells, zCell)
	return zCell
}

func (this *headerRow) convert() HeaderRow {
	zHeaderRow := make([]string, len(this.mCells))
	for i, zCell := range this.mCells {
		zHeaderRow[i] = zCell.mText
	}
	return HeaderRow(zHeaderRow)
}

type headerCell struct {
	mSet  bool
	mText string
}

func (this *headerCell) fillHole() {
	if !this.mSet {
		this.mText = "{NotSet}";
	}
}

func (this *headerCell) setCell(pCell *Cell) {
	this.mSet = true
	this.mText = pCell.mText
}

