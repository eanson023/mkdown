package main

import (
	"github.com/eanson023/mkdown"
)

const codeString = `
package main

import (
	"github.com/eanson023/mkdown"
)

const codeString = "fmt.Println(\"Hello,World\")"

func main() {
	md := mkdown.New("README.md")
	title1 := mkdown.NewTitleWithText(mkdown.Heading1, "🌟mkdown✔️")
	block := mkdown.NewBlock("author:" + mkdown.NewLink("eanson", "https://github.com/eanson023").String())
	title2 := mkdown.NewTitleWithText(mkdown.Heading2, "介绍")
	text := mkdown.NewText("📄 mkdown是一个用Go语言编写的库，提供了一组允许您简单写入md文件的功能 🦾")
	title3 := mkdown.NewTitleWithText(mkdown.Heading2, "基本用法")
	title4 := mkdown.NewTitleWithText(mkdown.Heading3, "安装")
	code1 := mkdown.NewCodeWithCodeBlock("go", "go get github.com/eanson023/mkdown")
	title5 := mkdown.NewTitleWithText(mkdown.Heading3, "使用")
	code2 := mkdown.NewCode("go")
	code2.SetCode(codeString)
	// 表格 两行两列
	table := mkdown.NewTable(2, 2).Add("性能").Add("简易度").Add("⭐").Update(2, 2, "⭐⭐⭐")
	// 嵌套列表举例
	const link = "https://github.com/eanson023/mkdown#"
	title6 := mkdown.NewTitleWithText(mkdown.Heading2, "List of Contents")
	// 创建有序列表
	ol := mkdown.NewOrderedList()
	// 创建无需列表
	ul := mkdown.NewUnOrderedList()
	ol.AppendNewLi(mkdown.NewLink("介绍", link+"介绍").String())
	li := ol.AppendNewLi(mkdown.NewLink("基本用法", link+"基本用法").String())
	ul.AppendNewLi(mkdown.NewLink("安装", link+"安装").String())
	ul.AppendNewLi(mkdown.NewLink("使用", link+"使用").String())
	// 讲ul加到ol中的li上
	li.AppendList(ul)
	md.Join(title1, block, title6, ol, title2, text, table, title3, title4, code1, title5, code2).Store()
}

`

func main() {
	md := mkdown.New("README.md")
	title1 := mkdown.NewTitleWithText(mkdown.Heading1, "🌟mkdown✔️")
	block := mkdown.NewBlock("author:" + mkdown.NewLink("eanson", "https://github.com/eanson023").String())
	title2 := mkdown.NewTitleWithText(mkdown.Heading2, "介绍")
	text := mkdown.NewText("📄 mkdown是一个用Go语言编写的库，提供了一组允许您简单写入markdown文件的功能 🦾")
	title3 := mkdown.NewTitleWithText(mkdown.Heading2, "基本用法")
	title4 := mkdown.NewTitleWithText(mkdown.Heading3, "安装")
	code1 := mkdown.NewCodeWithCodeBlock("go", "go get github.com/eanson023/mkdown")
	title5 := mkdown.NewTitleWithText(mkdown.Heading3, "使用")
	code2 := mkdown.NewCode("go")
	code2.SetCode(codeString)
	// 表格 两行两列
	table := mkdown.NewTable(2, 2).Add("性能").Add("简易度").Add("⭐").Update(2, 2, "⭐⭐⭐")
	// 嵌套列表举例
	const link = "https://github.com/eanson023/mkdown#"
	title6 := mkdown.NewTitleWithText(mkdown.Heading2, "List of Contents")
	// 创建有序列表
	ol := mkdown.NewOrderedList()
	// 创建无需列表
	ul := mkdown.NewUnOrderedList()
	ol.AppendNewLi(mkdown.NewLink("介绍", link+"介绍").String())
	li := ol.AppendNewLi(mkdown.NewLink("基本用法", link+"基本用法").String())
	ul.AppendNewLi(mkdown.NewLink("安装", link+"安装").String())
	ul.AppendNewLi(mkdown.NewLink("使用", link+"使用").String())
	// 讲ul加到ol中的li上
	li.AppendList(ul)
	md.Join(title1, block, title6, ol, title2, text, table, title3, title4, code1, title5, code2).Store()
}
