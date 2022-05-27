package surface

type Viewer struct {
	Width  int
	Height int
	Top    int
	Left   int
	Title  string
}

func NewViewer(width, height int, title string) *Viewer {
	return &Viewer{
		Width:  width,
		Height: height,
		Top:    0,
		Left:   0,
		Title:  title,
	}
}
