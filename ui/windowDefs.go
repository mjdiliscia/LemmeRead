package ui

import _ "embed"

//go:embed mainWindow.glade
var MainWindowUI []byte

//go:embed post.glade
var PostUI []byte

//go:embed comment.glade
var CommentUI []byte

//go:embed style.css
var StyleCSS []byte
