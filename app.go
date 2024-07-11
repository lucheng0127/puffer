package main

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/lucheng0127/virtuallan/pkg/client"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/websocket"
)

// App struct
type App struct {
	ctx    context.Context
	client *client.Client
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// Launch virtuallan client
func (a *App) Launch(target, user, password, key, logLevel string) {
	vl := client.NewClient(
		client.ClientSetKey(key),
		client.ClientSetTarget(target),
		client.ClientSetUser(user),
		client.ClientSetPasswd(password),
		client.ClientSetLogLevel(logLevel),
	)
	a.client = vl

	if err := vl.Launch(); err != nil {
		log.Errorf("launch virtuallan client %s", err.Error())
	}
}

// Check client login succeed or failed, true login succeed
func (a *App) CheckAuthed() bool {
	// Before check authed, maybe should sleep for 1 second wait for login succeed
	return a.client.CheckAuthed()
}

// Close client
func (a *App) Logout() {
	if err := a.client.Close(); err != nil {
		log.Errorf("virtuallan client logout %s", err.Error())
	}
}

// Redirect stdout to ws
func (a *App) MirrorLog(ws *websocket.Conn) {
	r, w, err := os.Pipe()
	if err != nil {
		fmt.Println(err)
	}

	// Redirect stdout to pipe
	os.Stdout = w

	for {
		io.Copy(ws, r)
	}
}
