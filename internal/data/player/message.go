package player

import (
	"dwarf-sweeper/internal/constants"
	"dwarf-sweeper/internal/menubox"
	"dwarf-sweeper/pkg/typeface"
	"fmt"
	"github.com/faiface/pixel"
)

type RawMessage struct {
	Raw     string
	OnClose func()
}

type Message struct {
	Player  *Player
	Box     *menubox.MenuBox
	Text    *typeface.Text
	Done    bool
	OnClose func()
	Prompt  *typeface.Text
}

func NewMessage(player *Player) *Message {
	box := menubox.NewBox(&player.CamPos, 1.0)
	box.SetSize(pixel.R(0., 0., 16., 16.))
	tex := typeface.New(&player.CamPos, "main", typeface.NewAlign(typeface.Left, typeface.Center), 1.2, constants.ActualHintSize, 0., 0.)
	tex.SetColor(constants.DefaultColor)
	pTex := typeface.New(&player.CamPos, "main", typeface.NewAlign(typeface.Right, typeface.Bottom), 1.2, constants.ActualHintSize, 0., 0.)
	pTex.SetColor(constants.DefaultColor)
	pTex.SetText(fmt.Sprintf("{symbol:%s-jump}", player.Code))
	return &Message{
		Player: player,
		Box:    box,
		Text:   tex,
		Done:   true,
		Prompt: pTex,
	}
}

func (m *Message) Update() {
	if m.Done && m.Box.IsOpen() {
		m.Box.Close()
		m.Text.SetText("")
		if m.OnClose != nil {
			m.OnClose()
		}
		m.OnClose = nil
	}
	m.Box.Pos.Y = m.Player.Canvas.Bounds().H() * -0.25
	w := m.Player.Canvas.Bounds().W()
	h := m.Text.Text.LineHeight * m.Text.RelativeSize * 2.
	m.Box.SetSize(pixel.R(0., 0., w - 20., h+30.))
	m.Box.Update()
	m.Text.SetWidth(w - 40.)
	m.Text.SetHeight(h + 10.)
	m.Text.SetPos(pixel.V(w*-0.5+10., m.Box.Pos.Y))
	m.Text.Update()
	m.Prompt.SetPos(pixel.V(w*0.5-10., m.Box.Pos.Y - m.Box.Rect.H() * 0.5 + 5.))
	m.Prompt.Update()
	if !m.Done && m.Box.IsOpen() && m.Player.Input.Get("jump").JustPressed() {
		m.Player.Input.Get("jump").Consume()
		m.Done = true
	}
}

func (m *Message) Draw() {
	if !m.Box.IsClosed() {
		m.Box.Draw(m.Player.Canvas)
		if m.Box.IsOpen() {
			m.Text.Draw(m.Player.Canvas)
			m.Prompt.Draw(m.Player.Canvas)
		}
	}
}