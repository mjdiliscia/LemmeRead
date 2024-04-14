package ui

import _ "embed"

var MainWindowUI string = `
<interface>
  <object class="GtkApplicationWindow" id="window">
    <property name="visible">TRUE</property>
  </object>
</interface>
`

//go:embed post.glade
var PostUI []byte

