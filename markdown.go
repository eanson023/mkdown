package mkdown

import "strconv"
import "os"
import "errors"
import "strings"
import "fmt"
import "bytes"

// Heading 标题
type Heading byte

const (
	// Heading1 一级标题
	Heading1 Heading = 1 + iota
	Heading2
	Heading3
	Heading4
	Heading5
	Heading6
)

// Handler 接口用于将内容写入bytes.Buffer
type Handler interface {
	// 创建
	Build(buffer *bytes.Buffer) error
}

// Sort 可排序接口 golang无法实现类的多态，我只有使用组合的方式来实现两种不同的排序方式
type Sort interface {
	// buffer:需要写入的缓存 list:列表数据源 recursiveNum：递归的次数，决定嵌套列表的前缀缩进
	createSort(buffer *bytes.Buffer, list *List, recursiveNum int) error
}

// Markdown 整个程序的主体
type Markdown struct {
	filename string
	buf      *bytes.Buffer
	handlers []Handler
}

// Table 表格
type Table struct {
	row   int
	col   int
	texts []*Text
	// 实际存储大小
	size int
	// 表格中内容最大长度大小
	maxLength int
}

// Title 标题
type Title struct {
	heading Heading
	text    *Text
}

// Text 基本文本
type Text struct {
	line string
}

// List 列表
type List struct {
	// 头节点
	head *Li
	//尾节点
	tail *Li
	//排序类型
	sorter Sort
}

// Li 列表里面的元素
type Li struct {
	//孩子列表
	child *List
	// 下一个
	next *Li
	//内容
	text *Text
}

// UnOrderedList 无序列表
type UnOrderedList struct{}

// OrderedList 有序列表
type OrderedList struct{}

// Block 区块
type Block struct {
	text *Text
}

// Code 代码块
type Code struct {
	language string
	code     *Text
}

// Link 链接
type Link struct {
	description string
	link        string
}

// New 创建新的markdown文档
func New(filename string) *Markdown {
	return &Markdown{
		filename: filename,
		buf:      &bytes.Buffer{},
	}
}

// NewText 创建新的文本
func NewText(data string) *Text {
	return &Text{
		line: data,
	}
}

// NewTitle 创建新的标题
func NewTitle(heading Heading) *Title {
	return &Title{
		heading: heading,
	}
}

// NewTitleWithText 创建携带文本的标题
func NewTitleWithText(heading Heading, title string) *Title {
	return &Title{
		heading: heading,
		text:    NewText(title),
	}
}

// NewTable 创建新的表格
func NewTable(row int, col int) *Table {
	return &Table{
		row:       row,
		col:       col,
		texts:     make([]*Text, row*col),
		size:      0,
		maxLength: 0,
	}
}

// NewOrderedList 创建一个空的有序列表 ol
func NewOrderedList() *List {
	return &List{
		sorter: new(OrderedList),
	}
}

// NewUnOrderedList 创建一个空的无序列表 ul
func NewUnOrderedList() *List {
	return &List{
		sorter: new(UnOrderedList),
	}
}

// NewLi 创建一个li
func NewLi(data string) *Li {
	return &Li{
		text: NewText(data),
	}
}

// NewBlock 创建一个区块
func NewBlock(data string) *Block {
	return &Block{
		text: NewText(data),
	}
}

// NewCode 创建一个空的代码块
func NewCode(language string) *Code {
	return &Code{
		language: language,
	}
}

// NewCodeWithCodeBlock 创建一个携带代码的代码块
func NewCodeWithCodeBlock(language string, code string) *Code {
	return &Code{
		language: language,
		code:     NewText(code),
	}
}

// NewLink 创建一个新的链接
func NewLink(description string, link string) *Link {
	return &Link{
		description: description,
		link:        link,
	}
}

// Join 添加handler到存储链
func (md *Markdown) Join(handlers ...Handler) *Markdown {
	for _, handle := range handlers {
		md.handlers = append(md.handlers, handle)
	}
	return md
}

// Store 进行存储
func (md *Markdown) Store() {
	for _, handler := range md.handlers {
		err := handler.Build(md.buf)
		if err != nil {
			panic(err)
		}
	}
	file, err := os.Create(md.filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	if _, err = md.buf.WriteTo(file); err != nil {
		panic(err)
	}

}

// Append 追加文字
func (tt *Text) Append(data string) {
	tt.line += data
}

// SetTitle 设置文本
func (title *Title) SetTitle(data string) {
	if title.text == nil {
		title.text = NewText(data)
	}
	title.text.line = data
}

// AppendNewLi 添加一个新的Li到list的尾部
func (list *List) AppendNewLi(data string) *Li {
	li := NewLi(data)
	// 如果头节点为空 同时代表整个list为空 所以需要同时把li挂到head和tail上
	if list.head == nil {
		list.head = li
		list.tail = li
	} else {
		list.tail.next = li
		list.tail = li
	}
	return li
}

// AppendList 添加List到某一个Li上
func (li *Li) AppendList(list *List) {
	li.child = list
}

// AppendList 讲新的list添加到目标list的尾节点li上
func (list *List) AppendList(l *List) {
	list.tail.AppendList(l)
}

func writeLine(data string, buffer *bytes.Buffer) error {
	_, err := buffer.WriteString(data + "\r\n")
	return err
}

func newLine(buffer *bytes.Buffer) {
	buffer.WriteString("\r\n")
}

// Build 写入文本
func (tt *Text) Build(buf *bytes.Buffer) error {
	if err := writeLine(tt.line, buf); err != nil {
		return err
	}
	newLine(buf)
	return nil
}

// Build 写入标题
func (title *Title) Build(buf *bytes.Buffer) error {
	orig := title.text.line
	title.text.line = fmt.Sprintf("%s %s", strings.Repeat("#", int(title.heading)), orig)
	if err := title.text.Build(buf); err != nil {
		return err
	}
	newLine(buf)
	return nil
}

// Build 创建表格
func (table *Table) Build(buf *bytes.Buffer) (err error) {
	//header
	//格式化对齐 eg:%-99s
	var format string = "| %-" + strconv.Itoa(table.maxLength) + "s "
	for i := 0; i < table.col; i++ {
		_, err = buf.WriteString(fmt.Sprintf(format, table.texts[i].line))
	}
	if err != nil {
		return err
	}
	err = writeLine("|", buf)
	// 分割线
	buf.WriteString(strings.Repeat(fmt.Sprintf("| %s ", strings.Repeat("-", table.maxLength)), table.col) + "|\r\n")
	// 内容
	for r := 1; r < table.row; r++ {
		for c := 0; c < table.col; c++ {
			if _, err = buf.WriteString(fmt.Sprintf(format, table.texts[r*table.col+c].line)); err != nil {
				return err
			}
		}
		err = writeLine("|", buf)
	}
	newLine(buf)
	return err
}

// Build 创建list
func (list *List) Build(buf *bytes.Buffer) (err error) {
	return list.sorter.createSort(buf, list, 0)
}

// Build 创建block
func (bk *Block) Build(buf *bytes.Buffer) (err error) {
	err = writeLine(fmt.Sprintf("> %s", bk.text.line), buf)
	newLine(buf)
	return err
}

// Build 创建代码块
func (c *Code) Build(buf *bytes.Buffer) (err error) {
	_, err = buf.WriteString(fmt.Sprintf("```%s\r\n%s\r\n```\r\n", c.language, c.code.line))
	if err != nil {
		return err
	}
	newLine(buf)
	return nil
}

// Build 创建链接
func (lk *Link) Build(buf *bytes.Buffer) (err error) {
	res := lk.generateLink()
	_, err = buf.WriteString(res + "\r\n")
	if err != nil {
		return err
	}
	newLine(buf)
	return nil
}

// createSort 无须列表排序
func (ul *UnOrderedList) createSort(buffer *bytes.Buffer, list *List, recursiveNum int) (err error) {
	node := list.head
	for node != nil {
		_, err = buffer.WriteString(fmt.Sprintf("%s* %s\r\n", strings.Repeat(" ", recursiveNum*3), node.text.line))
		if err != nil {
			return err
		}
		child := node.child
		if node.child != nil {
			// 递归排序孩子列表
			err = child.sorter.createSort(buffer, child, recursiveNum+1)
			if err != nil {
				return err
			}
		}
		node = node.next
	}
	newLine(buffer)
	return nil
}

// createSort 有序列表排序
func (ol *OrderedList) createSort(buffer *bytes.Buffer, list *List, recursiveNum int) (err error) {
	node := list.head
	index := 1
	for node != nil {
		_, err = buffer.WriteString(fmt.Sprintf("%s%d. %s\r\n", strings.Repeat(" ", recursiveNum*3), index, node.text.line))
		if err != nil {
			return err
		}
		child := node.child
		if node.child != nil {
			// 递归排序孩子列表
			err = child.sorter.createSort(buffer, child, recursiveNum+1)
			if err != nil {
				return err
			}
		}
		// 下标加1
		index++
		node = node.next
	}
	newLine(buffer)
	return nil
}

// Add 添加到表格得格子中 多了将panic异常
func (table *Table) Add(data string) *Table {
	length := table.row * table.col
	if table.size >= length {
		panic(errors.New("the table size will overflow"))
	}
	if len(data) > table.maxLength {
		table.maxLength = len(data)
	}
	index := table.size
	table.texts[index] = NewText(data)
	table.size++
	return table
}

// AddIgnoreError 直接链式编程但忽略了错误 多了会返回nil
func (table *Table) AddIgnoreError(data string) *Table {
	length := table.row * table.col
	if table.size >= length {
		return nil
	}
	if len(data) > table.maxLength {
		table.maxLength = len(data)
	}
	index := table.size
	table.texts[index] = NewText(data)
	table.size++
	return table
}

// Update 更新某个格子内容 注意这里是填人类的常规意识的数组下标
func (table *Table) Update(rowIdx int, colIdx int, data string) *Table {
	rowIdx--
	colIdx--
	if rowIdx < 0 || rowIdx >= table.row || colIdx < 0 || colIdx >= table.col {
		panic(errors.New("the width or length index of the table must be within the specified range of the table"))
	}
	if len(data) > table.maxLength {
		table.maxLength = len(data)
	}
	table.texts[rowIdx*table.row+colIdx] = NewText(data)
	return table
}

// SetCode 设置代码
func (c *Code) SetCode(code string) {
	c.code = NewText(code)
}

// AppendCode 添加代码块 会自动换行
func (c *Code) AppendCode(code string) {
	if c.code == nil {
		c.code = NewText(code)
	} else {
		c.code.line += "\r\n" + code
	}
}

func (lk *Link) generateLink() string {
	return fmt.Sprintf("[%s](%s)", lk.description, lk.link)
}

// String 返回链接的markdown格式
func (lk *Link) String() string {
	return lk.generateLink()
}
