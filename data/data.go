package data

import _ "embed"

//go:embed login.ui
var LoginUI []byte

//go:embed mainWindow.ui
var MainWindowUI []byte

//go:embed post.ui
var PostUI []byte

//go:embed comment.ui
var CommentUI []byte

//go:embed style.css
var StyleCSS []byte
