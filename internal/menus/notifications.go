package menus

import (
	"dwarf-sweeper/internal/constants"
	"dwarf-sweeper/internal/menubox"
	"dwarf-sweeper/pkg/camera"
	gween "dwarf-sweeper/pkg/gween64"
	"dwarf-sweeper/pkg/gween64/ease"
	"dwarf-sweeper/pkg/timing"
	"dwarf-sweeper/pkg/typeface"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

var NotificationHandler *notificationHandler

var (
	notificationX = constants.ActualW * -0.5 + 10.
	notificationY = constants.BaseH * 0.5 - 50.
	showTime      = 8.
)

type notificationHandler struct {
	Notifications []*notification
	MenuBox       *menubox.MenuBox
}

type notification struct {
	Raw    string
	Shown  *timing.Timer
	InterX *gween.Tween
	InterY *gween.Tween
	InterA *gween.Tween
	InterR *gween.Tween
	InterG *gween.Tween
	InterB *gween.Tween
	Text   *typeface.Text
}

func (nh *notificationHandler) Update() {
	notificationX = constants.ActualW * -0.5 + 10.
	notificationY = constants.BaseH * 0.5 - 50.
	if len(nh.Notifications) > 0 {
		w := 0.
		i := 0
		for i < 4 && i < len(nh.Notifications) {
			n := nh.Notifications[i]
			if n.Text.Width > w {
				w = n.Text.Width
			}
			i++
		}
		nh.MenuBox.SetSize(pixel.R(0., 0., (w + 15.) * 2., float64(i) * 10.))
		nh.MenuBox.Pos = pixel.V(constants.ActualW * -0.5, notificationY + 5. - float64(i) * 5.)
		nh.MenuBox.Open()
	} else {
		nh.MenuBox.Close()
	}
	nh.MenuBox.Update()
	removeTop := false
	for i := 0; i < 4 && i < len(nh.Notifications); i++ {
		n := nh.Notifications[i]
		if n.Shown == nil {
			n.Shown = timing.New(showTime)
			n.Text.Transform.Pos.Y = notificationY - 10. * float64(i)
		} else {
			n.Shown.Update()
		}
		if n.Text.Transform.Pos.X < notificationX && n.InterX == nil {
			n.InterX = gween.New(n.Text.Transform.Pos.X, notificationX, 0.5, ease.Linear)
		}
		if n.Text.Transform.Pos.Y != notificationY- 10. * float64(i) && n.InterY == nil {
			n.InterY = gween.New(n.Text.Transform.Pos.Y, notificationY- 10. * float64(i), 0.2, ease.Linear)
		}
		if n.Shown.Done() {
			//if n.Text.Color.A == 0 {
				removeTop = true
			//} else {
			//	if n.InterA == nil {
			//		n.InterA = gween.New(float64(n.Text.Color.A), 0., 0.15, ease.Linear)
			//	}
			//	if n.InterR == nil {
			//		n.InterR = gween.New(float64(n.Text.Color.R), 0., 0.15, ease.Linear)
			//	}
			//	if n.InterG == nil {
			//		n.InterG = gween.New(float64(n.Text.Color.G), 0., 0.15, ease.Linear)
			//	}
			//	if n.InterB == nil {
			//		n.InterB = gween.New(float64(n.Text.Color.B), 0., 0.15, ease.Linear)
			//	}
			//}
		}
		if n.InterX != nil {
			x, fin := n.InterX.Update(timing.DT)
			n.Text.SetPos(pixel.V(x, n.Text.Transform.Pos.Y))
			if fin {
				n.InterX = nil
			}
		}
		if n.InterY != nil {
			y, fin := n.InterY.Update(timing.DT)
			n.Text.SetPos(pixel.V(n.Text.Transform.Pos.X, y))
			if fin {
				n.InterY = nil
			}
		}
		//if n.InterA != nil {
		//	a, fin := n.InterA.Update(timing.DT)
		//	col := n.Text.Color
		//	col.A = uint8(a)
		//	n.Text.SetColor(col)
		//	if fin {
		//		n.InterA = nil
		//	}
		//}
		//if n.InterR != nil {
		//	r, fin := n.InterR.Update(timing.DT)
		//	col := n.Text.Color
		//	col.R = uint8(r)
		//	n.Text.SetColor(col)
		//	if fin {
		//		n.InterR = nil
		//	}
		//}
		//if n.InterG != nil {
		//	g, fin := n.InterG.Update(timing.DT)
		//	col := n.Text.Color
		//	col.G = uint8(g)
		//	n.Text.SetColor(col)
		//	if fin {
		//		n.InterG = nil
		//	}
		//}
		//if n.InterB != nil {
		//	b, fin := n.InterB.Update(timing.DT)
		//	col := n.Text.Color
		//	col.A = uint8(b)
		//	n.Text.SetColor(col)
		//	if fin {
		//		n.InterB = nil
		//	}
		//}
	}
	if removeTop {
		if len(nh.Notifications) > 1 {
			nh.Notifications = nh.Notifications[1:]
		} else {
			nh.Notifications = []*notification{}
		}
	}
}

func (nh *notificationHandler) Draw(win *pixelgl.Window) {
	for i := 0; i < 4 && i < len(nh.Notifications); i++ {
		if i == 0 {
			nh.MenuBox.Draw(win)
		}
		n := nh.Notifications[i]
		n.Text.Update()
		n.Text.Draw(win)
	}
}

func (nh *notificationHandler) AddMessage(raw string) {
	txt := typeface.New(&camera.Cam.APos, "main", typeface.NewAlign(typeface.Left, typeface.Center), 1., constants.ActualHintSize, 0., 0.)
	txt.SetText(raw)
	txt.SetColor(constants.DefaultColor)
	txt.SetPos(pixel.V(-(constants.ActualW * 0.5 + txt.Text.BoundsOf(raw).W() * txt.RelativeSize), notificationY))
	nh.Notifications = append(nh.Notifications, &notification{
		Raw:    raw,
		Text:   txt,
	})
}