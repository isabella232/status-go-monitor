package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/jroimartin/gocui"
)

func GenRenderFunc(g *gocui.Gui, sc *StateController) func() {
	return func() {
		ps := sc.State.GetData()
		renderPeerList(g, ps.Peers)
		renderPeerInfo(g, sc.State.GetCurrent())
		updatePeerCursor(g, ps.Current)
	}
}

func renderPeerList(g *gocui.Gui, peers []Peer) {
	if len(peers) == 0 {
		return
	}
	g.Update(func(g *gocui.Gui) error {
		v, err := g.View("main")
		if err != nil {
			return err
		}
		v.Clear()
		maxWidth, _ := g.Size()
		for _, peer := range peers {
			fmt.Fprintf(v, "%s\n", peer.AsTable(maxWidth))
		}
		return nil
	})
}

func renderPeerInfo(g *gocui.Gui, peer *Peer) {
	if peer == nil {
		return
	}
	g.Update(func(g *gocui.Gui) error {
		v, err := g.View("info")
		if err != nil {
			return err
		}
		v.Clear()
		fmt.Fprintf(v, strings.Repeat("%-8s: %v\n", 8),
			"Name", peer.Name,
			"ID", string(peer.ID),
			"Enode", peer.Enode,
			"Static", peer.Network.Static,
			"Trusted", peer.Network.Trusted,
			"Local", peer.Network.LocalAddress,
			"Remote", peer.Network.RemoteAddress,
			"Caps", strings.Join(peer.Caps, ", "))
		return nil
	})
}

func updatePeerCursor(g *gocui.Gui, current int) {
	v, err := g.View("main")
	if err != nil {
		log.Panicln("unable to find main view")
	}
	cx, _ := v.Cursor()

	if err := v.SetCursor(cx, current); err != nil {
		ox, _ := v.Origin()
		if err := v.SetOrigin(ox, current); err != nil {
			log.Panicln("unable to scroll")
		}
	}
}

func (p Peer) AsTable(maxWidth int) string {
	var id string
	if maxWidth > 160 {
		id = string(p.ID)
	} else {
		id = p.ID.String()
	}
	return fmt.Sprintf("%s ｜  %-15s ｜  %-21s ｜  %-7s ｜  %-8s",
		id, p.Name,
		p.Network.RemoteAddress,
		boolToString(p.Network.Trusted, "trusted", "normal"),
		boolToString(p.Network.Static, "static", "dynamic"))
}

func boolToString(v bool, yes string, no string) string {
	if v {
		return yes
	} else {
		return no
	}
}
