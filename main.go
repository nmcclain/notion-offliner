package main

// TODO: linked tables not working at all
// TODO: render all properties for List view
// TODO: handle inline links to "in-branch" notion pages
// TODO: render inline bookmarks

import (
	_ "embed"
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/docopt/docopt.go"
	"github.com/kjk/notionapi"
	"github.com/nmcclain/notion-offliner/tohtml"
)

const Version = "1.0"

var usage = `notion-offliner
Usage:
  notion-offliner [options] <pageid>
  notion-offliner -h | --help
  notion-offliner --version

Options:
  -m <message>  Add message to top of pages.
  -h --help     Show this screen.
  --version     Show version.
`

//go:embed "notion.tmpl"
var rawTemplate string

func main() {
	client := &notionapi.Client{
		AuthToken: os.Getenv("NOTION_TOKEN"),
		DebugLog:  true,
	}

	t, err := template.New("page").Parse(rawTemplate)
	if err != nil {
		log.Fatalf("template parse:  %v", err)
	}
	args, _ := docopt.Parse(usage, nil, true, Version, false)
	pageID := cleanupPageID(args["<pageid>"].(string))
	msg := ""
	if args["-m"] != nil {
		msg = args["-m"].(string)
	}

	dopage(client, pageID, ".", []tohtml.Crumb{}, t, msg)
}

func cleanupPageID(pageID string) string {
	// see https://github.com/jamalex/notion-py/blob/b7041ade477c1f59edab1b6fc025326d406dd92a/notion/utils.py#L20
	if strings.HasPrefix(pageID, "https://www.notion.so/") {
		parts := strings.Split(pageID, "?")
		pageID = parts[0]
		parts = strings.Split(pageID, "/")
		pageID = parts[len(parts)-1]
		parts = strings.Split(pageID, "-")
		pageID = parts[len(parts)-1]
	}
	return pageID
}

func dopage(client *notionapi.Client, pageID string, dir string, crumbs []tohtml.Crumb, t *template.Template, msg string) {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		os.Mkdir(dir, 0775)
	}
	page, err := client.DownloadPage(pageID)
	if err != nil {
		log.Printf("DownloadPage() failed BAILING: %s\n", err)
		return
	}
	log.Printf("dopage = %s [%s] %s", pageID, dir, page.Root().Type)

	c := tohtml.NewConverter(page)
	var h []byte
	var title string
	var filename string

	if page.Root().Type == notionapi.BlockCollectionViewPage || page.Root().Type == notionapi.BlockCollectionView {
		tv := page.Root().TableViews[0]
		title = page.TableViews[0].Collection.GetName()
		log.Printf("\tdopage CV title = %s", title)

		path := dir + "/" + safeName(title)
		filename = path + ".html"

		if page.Root().Type == notionapi.BlockCollectionViewPage {
			c.PushNewBuffer()
			c.RenderCVPage(page, crumbs, title)
			buf := c.PopBuffer()
			h = buf.Bytes()
			crumbs = append(crumbs, tohtml.Crumb{Name: title, Link: filename})
		} else {
			crumbs = append(crumbs, tohtml.Crumb{Name: title, Link: filename, Skip: true})
		}

		// sub-pages
		for _, row := range tv.Rows {
			if len(row.Page.ContentIDs) > 0 {
				dopage(client, row.Page.ID, path, crumbs, t, msg)
			}
		}
		if page.Root().Type == notionapi.BlockCollectionView {
			// don't save full-page view of inline tables.
			return
		}

	} else {
		title = page.Root().Title
		path := dir + "/" + safeName(title)
		filename = path + ".html"

		log.Printf("\tpage title = %s", title)
		for _, block := range page.Root().Content {
			doblock(client, block, title, dir, crumbs, t, msg)
		}

		h, err = c.ToHTML(crumbs)
		if err != nil {
			log.Fatalf("DownloadPage() failed with %s\n", err)
		}
	}

	f, err := os.Create(filename)
	if err != nil {
		log.Fatalf("create file:  %v", err)
	}
	err = t.Execute(f, struct {
		Title string
		Msg   string
		Body  template.HTML
	}{
		Title: title,
		Msg:   msg,
		Body:  template.HTML(h),
	})
	if err != nil {
		log.Fatalf("template execute:  %v", err)
	}
	f.Close()
	log.Printf("Wrote %s", filename)
}

func doblock(client *notionapi.Client, block *notionapi.Block, title, dir string, crumbs []tohtml.Crumb, t *template.Template, msg string) {
	path := dir + "/" + safeName(title)
	filename := path + ".html"
	if block.Type == notionapi.BlockPage || block.Type == notionapi.BlockCollectionViewPage || block.Type == notionapi.BlockCollectionView {
		dopage(client, block.ID, path, append(crumbs, tohtml.Crumb{Name: title, Link: filename}), t, msg)
	} else if block.Type == notionapi.BlockImage || block.Type == notionapi.BlockFile {
		if _, err := os.Stat(path); os.IsNotExist(err) {
			os.Mkdir(path, 0775)
		}
		uri := block.Source
		if strings.HasPrefix(uri, "https://s3-us-west-2.amazonaws.com/secure.notion-static.com/") {
			log.Printf("Downloading: %s", uri)
			img, err := client.DownloadFile(uri, block.ID)
			if err != nil {
				log.Fatalf("Download failed: %v", err)
			}
			err = ioutil.WriteFile(dir+"/"+tohtml.GetDownloadedFileName(uri, block), img.Data, 0644)
			if err != nil {
				log.Fatalf("Write failed: %v", err)
			}
		}
	} else if block.Type == notionapi.BlockColumn {
		log.Fatalf("Unexpected BlockColumn")
	} else if block.Type == notionapi.BlockColumnList {
		for _, column := range block.Content {
			if column.Type != notionapi.BlockColumn {
				log.Fatalf("Unexpected type in BlockColumnList: %s", column.Type)
			}
			for _, content := range column.Content {
				doblock(client, content, title, dir, crumbs, t, msg)
			}
		}
	}
}

func isSafeChar(r rune) bool {
	if r >= '0' && r <= '9' {
		return true
	}
	if r >= 'a' && r <= 'z' {
		return true
	}
	if r >= 'A' && r <= 'Z' {
		return true
	}
	return false
}

// safeName returns a file-system safe name
func safeName(s string) string {
	var res string
	for _, r := range s {
		if !isSafeChar(r) {
			res += " "
		} else {
			res += string(r)
		}
	}
	// replace multi-dash with single dash
	for strings.Contains(res, "  ") {
		res = strings.Replace(res, "  ", " ", -1)
	}
	res = strings.TrimLeft(res, " ")
	res = strings.TrimRight(res, " ")
	if len(res) > 64 {
		res = res[:64]
	}
	return res
}
