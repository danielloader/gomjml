package mjml

import (
	"strings"
	"testing"
)

// Callers commonly escape data (html/template, html.EscapeString) before
// putting it into a template. MJML keeps that escaping, so the data reaches the
// email as text, never as markup.
func TestEscapedMarkupStaysEscaped(t *testing.T) {
	const escaped = `&lt;b&gt;bold&lt;/b&gt; &amp; &#60;i&#62;`
	column := func(content string) string {
		return `<mjml><mj-body><mj-section><mj-column>` + content + `</mj-column></mj-section></mj-body></mjml>`
	}
	tests := []struct {
		name string
		mjml string
	}{
		{"mj-text", column(`<mj-text>` + escaped + `</mj-text>`)},
		{"mj-button", column(`<mj-button href="https://example.com">` + escaped + `</mj-button>`)},
		{"mj-table", column(`<mj-table><tr><td>` + escaped + `</td></tr></mj-table>`)},
		{"mj-raw", column(`<mj-raw><p>` + escaped + `</p></mj-raw>`)},
		{"mj-social-element", column(`<mj-social><mj-social-element name="facebook" href="https://example.com">` + escaped + `</mj-social-element></mj-social>`)},
		{"mj-navbar-link", column(`<mj-navbar><mj-navbar-link href="/a">` + escaped + `</mj-navbar-link></mj-navbar>`)},
		{"mj-accordion-title", column(`<mj-accordion><mj-accordion-element><mj-accordion-title>` + escaped + `</mj-accordion-title><mj-accordion-text>x</mj-accordion-text></mj-accordion-element></mj-accordion>`)},
		{"mj-accordion-text", column(`<mj-accordion><mj-accordion-element><mj-accordion-title>x</mj-accordion-title><mj-accordion-text>` + escaped + `</mj-accordion-text></mj-accordion-element></mj-accordion>`)},
		{"mj-title", `<mjml><mj-head><mj-title>` + escaped + `</mj-title></mj-head><mj-body></mj-body></mjml>`},
		{"mj-preview", `<mjml><mj-head><mj-preview>` + escaped + `</mj-preview></mj-head><mj-body></mj-body></mjml>`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			html, err := Render(tt.mjml)
			if err != nil {
				t.Fatalf("Render: %v", err)
			}
			if strings.Contains(html, "<b>bold</b>") || strings.Contains(html, "<i>") {
				t.Errorf("escaped markup became live markup:\n%s", html)
			}
			if !strings.Contains(html, escaped) {
				t.Errorf("output does not contain %s\n%s", escaped, html)
			}
		})
	}
}

func TestTextBetweenComponentsStaysEscaped(t *testing.T) {
	html, err := Render(`<mjml><mj-body><mj-section>&lt;b&gt;bold&lt;/b&gt;</mj-section></mj-body></mjml>`)
	if err != nil {
		t.Fatalf("Render: %v", err)
	}
	if strings.Contains(html, "<b>bold</b>") {
		t.Errorf("escaped markup became live markup:\n%s", html)
	}
}

// MJML writes attribute values as they are written in the template.
func TestAttributeValuesAreRenderedAsWritten(t *testing.T) {
	html, err := Render(`<mjml><mj-body><mj-section><mj-column>
		<mj-image src="https://example.com/i.png?a=1&amp;b=2" alt="say &quot;hi&quot;" title="x&#34; onerror=&#34;y" href="https://example.com/?a=1&amp;b=2" />
		<mj-image src="https://example.com/j.png?a=1&b=2" href="https://example.com/?c=1&d=2" />
		<mj-button href="https://example.com/?e=1&amp;f=2" title="&lt;b&gt;">Go</mj-button>
	</mj-column></mj-section></mj-body></mjml>`)
	if err != nil {
		t.Fatalf("Render: %v", err)
	}
	for _, want := range []string{
		`src="https://example.com/i.png?a=1&amp;b=2"`,
		`alt="say &quot;hi&quot;"`,
		`title="x&#34; onerror=&#34;y"`,
		`href="https://example.com/?a=1&amp;b=2"`,
		`src="https://example.com/j.png?a=1&b=2"`,
		`href="https://example.com/?c=1&d=2"`,
		`href="https://example.com/?e=1&amp;f=2"`,
		`title="&lt;b&gt;"`,
	} {
		if !strings.Contains(html, want) {
			t.Errorf("output does not contain %s", want)
		}
	}
	if strings.Contains(html, `onerror="y"`) {
		t.Errorf("an escaped quote ended the attribute:\n%s", html)
	}
}
