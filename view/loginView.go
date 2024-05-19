package view

import (
	"log"

	"github.com/gotk3/gotk3/gtk"
	"github.com/mjdiliscia/LemmeRead/data"
	"github.com/mjdiliscia/LemmeRead/utils"
)

type LoginView struct {
	LoginClicked func(string, string, string)
	Window       *gtk.Dialog

	server   *gtk.Entry
	username *gtk.Entry
	password *gtk.Entry
	login    *gtk.Button
}

func (lv *LoginView) SetupLoginView() (err error) {
	_, err = lv.buildAndSetReferences()
	if err != nil {
		return
	}

	lv.login.Connect("clicked", func() {
		if lv.LoginClicked != nil {
			server, err := lv.server.GetText()
			if err != nil {
				log.Panic(err)
			}

			username, err := lv.username.GetText()
			if err != nil {
				log.Panic(err)
			}

			password, err := lv.password.GetText()
			if err != nil {
				log.Panic(err)
			}

			lv.LoginClicked(server, username, password)
		}
	})

	lv.Window.Show()

	return nil
}

func (lv *LoginView) buildAndSetReferences() (builder *gtk.Builder, err error) {
	builder, err = gtk.BuilderNewFromString(string(data.LoginUI))
	if err != nil {
		return
	}

	lv.Window, err = utils.GetUIObject[gtk.Dialog](builder, "loginDialog")
	if err != nil {
		return
	}

	lv.login, err = utils.GetUIObject[gtk.Button](builder, "login")
	if err != nil {
		return
	}

	lv.server, err = utils.GetUIObject[gtk.Entry](builder, "serverUrl")
	if err != nil {
		return
	}

	lv.username, err = utils.GetUIObject[gtk.Entry](builder, "username")
	if err != nil {
		return
	}

	lv.password, err = utils.GetUIObject[gtk.Entry](builder, "password")
	if err != nil {
		return
	}

	return
}

func (lv *LoginView) DestroyWindow() {
	lv.Window.Destroy()
	lv.LoginClicked = nil
}
