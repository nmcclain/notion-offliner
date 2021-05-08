package tohtml

import (
	"github.com/kjk/notionapi"
)

func (c *Converter) RenderCVPage(page *notionapi.Page, crumbs []Crumb, title string) {
	// log.Printf("page = %+v\n", page)
	clsFont := "sans"
	/*
		fp := page.FormatPage()
		if fp != nil {
			if fp.PageFont != "" {
				clsFont = fp.PageFont
			}
		}
	*/
	c.RenderCrumbs(crumbs, title)
	c.Printf(`<article id="%s" class="page %s">`, page.ID, clsFont)
	c.renderCVPageHeader(page, title)
	{
		c.Printf(`<div class="page-body">`)
		c.RenderCollectionView(page.Root())
		c.Printf(`</div>`)
	}
	c.Printf(`</article>`)
}

func (c *Converter) renderCVPageHeader(page *notionapi.Page, title string) {
	coll := page.Root().TableViews[0].Collection
	c.Printf(`<header>`)
	if coll.Cover != "" {
		position := (1 - coll.Format.CoverPosition) * 100
		coverURL := FilePathFromPageCoverURL(coll.Cover, page.Root())
		coverURL = EscapeHTML(coverURL)
		c.Printf(`<img class="page-cover-image" src="%s" style="object-position:center %v%%"/>`, coverURL, position)
	}
	if coll.Icon != "" {
		// TODO: "undefined" is a bug in Notion export
		clsCover := "undefined"
		if coll.Cover != "" {
			clsCover = "page-header-icon-with-cover"
		}
		c.Printf(`<div class="page-header-icon %s">`, clsCover)
		if isURL(coll.Icon) {
			fileName := GetDownloadedFileName(coll.Icon, page.Root())
			c.Printf(`<img class="icon" src="%s"/>`, fileName)
		} else {
			c.Printf(`<span class="icon">%s</span>`, coll.Icon)
		}
		c.Printf(`</div>`)
	}

	c.Printf(`<h1 class="page-title">`)
	c.Printf(title)
	c.Printf(`</h1>`)

	if len(coll.Description) > 0 {
		c.Printf(`<h3 class="">`)
		c.Printf(coll.Description[0].([]interface{})[0].(string))
		c.Printf(`</h3>`)
	}

	c.Printf(`</header>`)
}
