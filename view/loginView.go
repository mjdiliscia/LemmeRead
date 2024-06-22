package view

import (
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
	"github.com/mjdiliscia/LemmeRead/data"
	"github.com/mjdiliscia/LemmeRead/utils"
)

type LoginView struct {
	LoginClicked func(string, string, string, string)
	Window       *gtk.Dialog

	server   *gtk.Entry
	username *gtk.Entry
	password *gtk.Entry
	totp     *gtk.Entry
	login    *gtk.Button
}

func (lv *LoginView) SetupLoginView() (err error) {
	_, err = lv.buildAndSetReferences()
	if err != nil {
		return
	}

	lv.login.Connect("clicked", func() {
		if lv.LoginClicked != nil {
			server := lv.server.Text()
			username := lv.username.Text()
			password := lv.password.Text()
			totp := lv.totp.Text()

			lv.LoginClicked(server, username, password, totp)
		}
	})

	lv.Window.Show()

	return nil
}

func (lv *LoginView) buildAndSetReferences() (builder *gtk.Builder, err error) {
	builder = gtk.NewBuilderFromString(string(data.LoginUI), -1)

	lv.Window, err = utils.GetUIObject[*gtk.Dialog](builder, "loginDialog")
	if err != nil {
		return
	}

	lv.login, err = utils.GetUIObject[*gtk.Button](builder, "login")
	if err != nil {
		return
	}

	lv.server, err = utils.GetUIObject[*gtk.Entry](builder, "serverUrl")
	if err != nil {
		return
	}

	lv.username, err = utils.GetUIObject[*gtk.Entry](builder, "username")
	if err != nil {
		return
	}

	lv.password, err = utils.GetUIObject[*gtk.Entry](builder, "password")
	if err != nil {
		return
	}

	lv.totp, err = utils.GetUIObject[*gtk.Entry](builder, "totp")
	if err != nil {
		return
	}

	return
}

func (lv *LoginView) DestroyWindow() {
	lv.Window.Destroy()
	lv.LoginClicked = nil
}
