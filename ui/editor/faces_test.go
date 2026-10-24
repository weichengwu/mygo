package editor

import (
	"bytes"
	"testing"
)

// A token style's Italic and Bold draw its glyphs in that face, where
// the regular ones would be, and leave the other tokens as they were.
func TestTokenFaces(t *testing.T) {
	const src = "x := 1\n// a comment\ny := 2|"
	rows := func(theme *Theme) [2][]byte {
		ct := newEditorTest(t, "go", src)
		ct.st.Theme, ct.st.NoActiveLine = theme, true
		ct.tt.Frame()
		img := ct.tt.Image()
		s := float32(img.Bounds().Dx()) / 600
		lh := ct.st.ed.view.lineHeight
		var out [2][]byte // the caret's row, which blinks, left out
		for i := range out {
			top, bottom := int(float32(i)*lh*s), int(float32(i+1)*lh*s)
			out[i] = bytes.Clone(img.Pix[top*img.Stride : bottom*img.Stride])
		}
		return out
	}
	plain := rows(&Theme{})
	italic := &Theme{}
	italic.Tokens[TokenComment].Italic = true
	slanted := rows(italic)
	if !bytes.Equal(plain[0], slanted[0]) {
		t.Error("the code above the comment changed")
	}
	if bytes.Equal(plain[1], slanted[1]) {
		t.Error("the comment is not drawn in italics")
	}
	bold := &Theme{}
	bold.Tokens[TokenComment].Bold = true
	if heavy := rows(bold); bytes.Equal(plain[1], heavy[1]) || bytes.Equal(slanted[1], heavy[1]) {
		t.Error("the comment is not drawn in bold")
	}
}
