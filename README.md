# notion-offliner
A CLI tool that you can use create regular backups of your Notion.so Pages.

Perfect for disaster scenarios and offline usage. MacOS and Linux.

## Usage
* Download it for [Linux AMD64](https://github.com/nmcclain/notion-offliner/releases/download/v0.1/notion-offliner-linux-amd64), [MacOS](https://github.com/nmcclain/notion-offliner/releases/download/v0.1/notion-offliner-macos), or [MacOS M1](https://github.com/nmcclain/notion-offliner/releases/download/v0.1/notion-offliner-macos-m1).
* `./notion-offliner https://www.notion.so/Learn-the-shortcuts-66e28cec810548c3a4061513126766b0`
* If you want to access private Pages, you must set the `NOTION_TOKEN` environment variable.
  * For example: `export NOTION_TOKEN=XXXXXXX`
  * You can find the token in your logged-in browser's cookies. It's named `token_v2`.
  * Or here's an [automated way to get your token](https://gist.github.com/nmcclain/30e4f5ee231e606bc4a1dcedef1c551f).

## Features
* Recursively captures all Pages beneath Page you provide.
* Includes image and file attachments.
* Automatically adds simple "breadcrumb"-style navigation to all Pages.
* Optional `-m "message"` to add a static footer (ex: `-m "Offline copy 5/8/2021"`)
* Based on the awesome [notionapi project](https://github.com/kjk/notionapi/).

## Known Limitations (help appreciated!)
* Linked tables not working at all.
* Only Table views work for Collections. List view properties are not all rendered, and board/timeline/calendar/gallery views are not supported at all.
* Inline Bookmark blocks aren't rendered properly.
* Inline Page links to "in-branch" Pages are downloaded twice - should be linked to first download.
* Any issue with the [notionapi project](https://github.com/kjk/notionapi/issues).

## Contributing

**Something bugging you?** Please open an Issue or Pull Request - we're here to help!

**New Feature Ideas?** Please open Pull Request!
 
**All Humans Are Equal In This Project And Will Be Treated With Respect.**
