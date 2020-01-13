package ansi

import (
	"bytes"
	"io"
	"text/template"

	"github.com/muesli/termenv"
)

// BaseElement renders a styled primitive element.
type BaseElement struct {
	Token  string
	Prefix string
	Suffix string
	Style  StylePrimitive
}

func formatToken(format string, token string) (string, error) {
	var b bytes.Buffer

	v := make(map[string]interface{})
	v["text"] = token

	tmpl, err := template.New(format).Funcs(TemplateFuncMap).Parse(format)
	if err != nil {
		return "", err
	}

	err = tmpl.Execute(&b, v)
	return b.String(), err
}

func renderText(w io.Writer, rules StylePrimitive, s string) {
	if len(s) == 0 {
		return
	}

	p := termenv.SupportedColorProfile()
	out := termenv.String(s)

	if rules.Color != nil {
		out = out.Foreground(p.Color(*rules.Color))
	}
	if rules.BackgroundColor != nil {
		out = out.Background(p.Color(*rules.BackgroundColor))
	}

	if rules.Underline != nil && *rules.Underline {
		out = out.Underline()
	}
	if rules.Bold != nil && *rules.Bold {
		out = out.Bold()
	}
	if rules.Italic != nil && *rules.Italic {
		out = out.Italic()
	}
	if rules.CrossedOut != nil && *rules.CrossedOut {
		out = out.CrossOut()
	}
	if rules.Overlined != nil && *rules.Overlined {
		out = out.Overline()
	}
	if rules.Inverse != nil && *rules.Inverse {
		out = out.Reverse()
	}
	if rules.Blink != nil && *rules.Blink {
		out = out.Blink()
	}

	_, _ = w.Write([]byte(out.String()))
}

func (e *BaseElement) Render(w io.Writer, ctx RenderContext) error {
	bs := ctx.blockStack

	renderText(w, bs.Current().Style.StylePrimitive, e.Prefix)
	defer func() {
		renderText(w, bs.Current().Style.StylePrimitive, e.Suffix)
	}()

	rules := bs.With(e.Style)
	// render unstyled prefix/suffix
	renderText(w, bs.Current().Style.StylePrimitive, rules.BlockPrefix)
	defer func() {
		renderText(w, bs.Current().Style.StylePrimitive, rules.BlockSuffix)
	}()

	// render styled prefix/suffix
	renderText(w, rules, rules.Prefix)
	defer func() {
		renderText(w, rules, rules.Suffix)
	}()

	s := e.Token
	if len(rules.Format) > 0 {
		var err error
		s, err = formatToken(rules.Format, s)
		if err != nil {
			return err
		}
	}
	renderText(w, rules, s)
	return nil
}
