package ghostscript

// Document represents the entire pdf document
type Document struct {
	// TODO add more metadata like title, author, ...
	NumPages int    `json:"numPages"`
	Pages    []Page `json:"pages"`
}

// NewDocument creates a new document with the given amount of changes
func NewDocument(numPages int) *Document {
	var doc Document
	doc.NumPages = numPages
	doc.Pages = make([]Page, numPages)
	return &doc
}

// Page one page of a pdf
type Page struct {
	Texts        []Text `json:"texts"`
	ThumbnailKey string `json:"thumbnailKey"` // minio/S3 key of the thumbnail
}

// AddText adds a text to the pdf page
func (p *Page) AddText(content string, x, y, w, h float64) {
	text := Text{
		Position:  Coordinate{X: x, Y: y},
		Dimension: Dimension{Width: w, Height: h},
		Content:   content,
	}

	p.Texts = append(p.Texts, text)
}

// Text fragment of test on the pdf
type Text struct {
	Position  Coordinate `json:"position"`
	Dimension Dimension  `json:"dimension"`
	Content   string     `json:"content"`
}

// Coordinate X and Y in percent
type Coordinate struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

// Dimension Width and Height in percent
type Dimension struct {
	Width  float64 `json:"width"`
	Height float64 `json:"height"`
}
