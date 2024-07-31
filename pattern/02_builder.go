package builder

import (
	"fmt"
)

// Документ определяет интерфейс для документа
type Document interface {
	GetContent() string
}

// JSONDocument реализует интерфейс Document для JSON документов
type JSONDocument struct {
	content string
}

func (d *JSONDocument) GetContent() string {
	return d.content
}

// XMLDocument реализует интерфейс Document для XML документов
type XMLDocument struct {
	content string
}

func (d *XMLDocument) GetContent() string {
	return d.content
}

// DocumentBuilder определяет интерфейс для построителя документов
type DocumentBuilder interface {
	SetTitle(title string) DocumentBuilder
	SetBody(body string) DocumentBuilder
	SetFooter(footer string) DocumentBuilder
	Build() Document
}

// JSONDocumentBuilder реализует интерфейс DocumentBuilder для создания JSON документов
type JSONDocumentBuilder struct {
	document *JSONDocument
}

func NewJSONDocumentBuilder() *JSONDocumentBuilder {
	return &JSONDocumentBuilder{document: &JSONDocument{}}
}

func (b *JSONDocumentBuilder) SetTitle(title string) DocumentBuilder {
	b.document.content += fmt.Sprintf(`{"title": "%s",`, title)
	return b
}

func (b *JSONDocumentBuilder) SetBody(body string) DocumentBuilder {
	b.document.content += fmt.Sprintf(`"body": "%s",`, body)
	return b
}

func (b *JSONDocumentBuilder) SetFooter(footer string) DocumentBuilder {
	b.document.content += fmt.Sprintf(`"footer": "%s"}`, footer)
	return b
}

func (b *JSONDocumentBuilder) Build() Document {
	return b.document
}

// XMLDocumentBuilder реализует интерфейс DocumentBuilder для создания XML документов
type XMLDocumentBuilder struct {
	document *XMLDocument
}

func NewXMLDocumentBuilder() *XMLDocumentBuilder {
	return &XMLDocumentBuilder{document: &XMLDocument{}}
}

func (b *XMLDocumentBuilder) SetTitle(title string) DocumentBuilder {
	b.document.content += fmt.Sprintf(`<title>%s</title>`, title)
	return b
}

func (b *XMLDocumentBuilder) SetBody(body string) DocumentBuilder {
	b.document.content += fmt.Sprintf(`<body>%s</body>`, body)
	return b
}

func (b *XMLDocumentBuilder) SetFooter(footer string) DocumentBuilder {
	b.document.content += fmt.Sprintf(`<footer>%s</footer>`, footer)
	return b
}

func (b *XMLDocumentBuilder) Build() Document {
	return b.document
}

// Director управляет построением документа с использованием построителя
type Director struct {
	builder DocumentBuilder
}

func NewDirector(builder DocumentBuilder) *Director {
	return &Director{builder: builder}
}

func (d *Director) Construct(title, body, footer string) Document {
	d.builder.SetTitle(title).SetBody(body).SetFooter(footer)
	return d.builder.Build()
}

// Основная функция демонстрирует использование паттерна Строитель
func main() {
	// Создание JSON документа
	jsonBuilder := NewJSONDocumentBuilder()
	director := NewDirector(jsonBuilder)
	jsonDocument := director.Construct("JSON Title", "This is the body of the JSON document", "JSON Footer")
	fmt.Println("JSON Document:")
	fmt.Println(jsonDocument.GetContent())

	// Создание XML документа
	xmlBuilder := NewXMLDocumentBuilder()
	director = NewDirector(xmlBuilder)
	xmlDocument := director.Construct("XML Title", "This is the body of the XML document", "XML Footer")
	fmt.Println("\nXML Document:")
	fmt.Println(xmlDocument.GetContent())
}
