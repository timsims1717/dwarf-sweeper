package menus

import (
	"dwarf-sweeper/internal/cfg"
	"dwarf-sweeper/pkg/camera"
	"dwarf-sweeper/pkg/img"
	"dwarf-sweeper/pkg/input"
	"dwarf-sweeper/pkg/sfx"
	"dwarf-sweeper/pkg/timing"
	"dwarf-sweeper/pkg/transform"
	"dwarf-sweeper/pkg/util"
	"fmt"
	"github.com/faiface/pixel"
	"golang.org/x/image/colornames"
	"image/color"
	"math"
)

var (
	DefaultColor  color.RGBA
	HoverColor    color.RGBA
	DisabledColor color.RGBA
	DefaultSize   pixel.Vec
	HoverSize     pixel.Vec

	corner *pixel.Sprite
	sideV  *pixel.Sprite
	sideH  *pixel.Sprite
	inner  *pixel.Sprite
)

func Initialize() {
	corner = img.Batchers[cfg.MenuSprites].Sprites["menu_corner"]
	sideV = img.Batchers[cfg.MenuSprites].Sprites["menu_side_v"]
	sideH = img.Batchers[cfg.MenuSprites].Sprites["menu_side_h"]
	inner = img.Batchers[cfg.MenuSprites].Sprites["menu_inner"]
	DefaultColor = color.RGBA{
		R: 74,
		G: 84,
		B: 98,
		A: 255,
	}
	HoverColor = colornames.Mediumblue
	DisabledColor = colornames.Darkgray
	DefaultSize = pixel.V(1.4, 1.4)
	HoverSize = pixel.V(1.45, 1.45)
}

type DwarfMenu struct {
	Key     string
	ItemMap map[string]*Item
	Items   []*Item
	Hovered int
	Trans   *transform.Transform
	Rect    pixel.Rect
	Closed  bool
	closing bool
	opened  bool
	StepV   float64
	StepH   float64
	Cam     *camera.Camera

	Center *transform.Transform
	CTUL   *transform.Transform
	CTUR   *transform.Transform
	CTDR   *transform.Transform
	CTDL   *transform.Transform
	STU    *transform.Transform
	STR    *transform.Transform
	STD    *transform.Transform
	STL    *transform.Transform
}

func New(key string, rect pixel.Rect, cam *camera.Camera) *DwarfMenu {
	tran := transform.NewTransform()
	tran.Anchor = transform.Anchor{
		H: transform.Center,
		V: transform.Center,
	}
	tran.SetRect(rect)
	Center := transform.NewTransform()
	CTUL := transform.NewTransform()
	CTUR := transform.NewTransform()
	CTDR := transform.NewTransform()
	CTDL := transform.NewTransform()
	STU := transform.NewTransform()
	STR := transform.NewTransform()
	STD := transform.NewTransform()
	STL := transform.NewTransform()
	CTUR.Flip = true
	CTDR.Flip = true
	CTDR.Flop = true
	CTDL.Flop = true
	STR.Flip = true
	STD.Flop = true
	return &DwarfMenu{
		ItemMap: map[string]*Item{},
		Items:   []*Item{},
		Trans:   tran,
		Cam:     cam,
		Rect:    rect,
		StepV:   16.,
		StepH:   16.,
		Center:  Center,
		CTUL:    CTUL,
		CTUR:    CTUR,
		CTDR:    CTDR,
		CTDL:    CTDL,
		STU:     STU,
		STR:     STR,
		STD:     STD,
		STL:     STL,
	}
}

func (m *DwarfMenu) AddItem(key, raw, sRaw string) *Item {
	if _, ok := m.ItemMap[key]; ok {
		panic(fmt.Errorf("menu '%s' already has item '%s'", m.Key, key))
	}
	item := NewItem(key, raw, sRaw)
	m.ItemMap[key] = item
	m.Items = append(m.Items, item)
	return item
}

func (m *DwarfMenu) InsertItem(key, raw, sRaw string, i int) *Item {
	if _, ok := m.ItemMap[key]; ok {
		panic(fmt.Errorf("menu '%s' already has item '%s'", m.Key, key))
	}
	item := NewItem(key, raw, sRaw)
	m.ItemMap[key] = item
	if i < 0 {
		i = 0
	}
	if i >= len(m.Items) {
		m.Items = append(m.Items, item)
	} else {
		second := m.Items[i:]
		m.Items = append(m.Items[:i], item)
		m.Items = append(m.Items, second...)
	}
	return item
}

func (m *DwarfMenu) RemoveItem(key string) {
	index := -1
	for i, item := range m.Items {
		if item.Key == key {
			index = i
			break
		}
	}
	if index != -1 {
		if len(m.Items) > 1 {
			m.Items = append(m.Items[:index], m.Items[index+1:]...)
		} else {
			m.Items = []*Item{}
		}
	}
	m.ItemMap[key] = nil
}

func (m *DwarfMenu) Open() {
	m.Closed = false
	m.closing = false
	m.opened = false
	hover := false
	for i, item := range m.Items {
		if !hover && !item.disabled && !item.NoHover {
			item.Hovered = true
			m.Hovered = i
			hover = true
		} else {
			item.Hovered = false
		}
	}
}

func (m *DwarfMenu) Close() {
	m.closing = true
	m.opened = false
}

func (m *DwarfMenu) CloseInstant() {
	m.closing = true
	m.Closed = true
	m.opened = false
	m.StepV = 16.
	m.StepH = 16.
}

func (m *DwarfMenu) Update(in *input.Input) {
	m.UpdateSize()
	m.UpdateView()
	m.UpdateTransforms()
	if !m.closing {
		m.UpdateItems(in)
	}
}

func (m *DwarfMenu) UpdateSize() {
	minWidth := 8.
	minHeight := 8.
	for _, item := range m.Items {
		bW := item.Text.BoundsOf(item.Raw).W()
		sW := 0.
		if item.SRaw != "" {
			sW = item.SText.BoundsOf(item.SRaw).W() + item.Text.BoundsOf("   ").W()
		}
		minWidth = math.Max((bW + sW) * 1.4, minWidth)
		minHeight += item.Text.LineHeight
	}
	minWidth = math.Floor(math.Max(minWidth + 30., m.Rect.W()))
	minHeight = math.Floor(math.Max(minHeight, m.Rect.H()))
	m.Rect = pixel.R(0., 0., minWidth, minHeight)
	for i, item := range m.Items {
		item.Transform.Pos.Y = minHeight * 0.5 - float64(i + 1) * item.Text.LineHeight
		item.Transform.Pos.X = minWidth * -0.5 + 20.
		item.STransform.Pos.Y = item.Transform.Pos.Y
		item.STransform.Pos.X = minWidth * 0.5 - 10.
	}
}

func (m *DwarfMenu) UpdateView() {
	if !m.closing {
		if m.StepV < m.Rect.H() * 0.5 {
			m.StepV += timing.DT * 300.
		}
		if m.StepV > m.Rect.H() * 0.5 {
			m.StepV = m.Rect.H() * 0.5
		}
		if m.StepH < m.Rect.W() * 0.5 {
			m.StepH += timing.DT * 400.
		}
		if m.StepH > m.Rect.W() * 0.5 {
			m.StepH = m.Rect.W() * 0.5
		}
		if m.StepH >= m.Rect.W() * 0.5 && m.StepV >= m.Rect.H() * 0.5 {
			m.opened = true
		}
	} else {
		if m.StepV > 16. {
			m.StepV -= timing.DT * 300.
		}
		if m.StepV < 16. {
			m.StepV = 16.
		}
		if m.StepH > 16. {
			m.StepH -= timing.DT * 400.
		}
		if m.StepH < 16. {
			m.StepH = 16.
		}
		if m.StepH < 20. && m.StepV < 20. {
			m.Closed = true
		}
	}
}

func (m *DwarfMenu) UpdateTransforms() {
	if m.Cam != nil {
		m.Trans.UIZoom = m.Cam.GetZoomScale()
		m.Trans.UIPos = m.Cam.APos
		m.CTUL.UIZoom = m.Cam.GetZoomScale()
		m.CTUL.UIPos = m.Cam.APos
		m.CTUR.UIZoom = m.Cam.GetZoomScale()
		m.CTUR.UIPos = m.Cam.APos
		m.CTDR.UIZoom = m.Cam.GetZoomScale()
		m.CTDR.UIPos = m.Cam.APos
		m.CTDL.UIZoom = m.Cam.GetZoomScale()
		m.CTDL.UIPos = m.Cam.APos
		m.STU.UIZoom = m.Cam.GetZoomScale()
		m.STU.UIPos = m.Cam.APos
		m.STR.UIZoom = m.Cam.GetZoomScale()
		m.STR.UIPos = m.Cam.APos
		m.STD.UIZoom = m.Cam.GetZoomScale()
		m.STD.UIPos = m.Cam.APos
		m.STL.UIZoom = m.Cam.GetZoomScale()
		m.STL.UIPos = m.Cam.APos
		m.Center.UIZoom = m.Cam.GetZoomScale()
		m.Center.UIPos = m.Cam.APos
	}
	m.CTUL.Pos = pixel.V(-m.StepH, m.StepV)
	m.CTUL.Scalar = pixel.V(1.4, 1.4)
	m.CTUL.Update()
	m.CTUR.Pos = pixel.V(m.StepH, m.StepV)
	m.CTUR.Scalar = pixel.V(1.4, 1.4)
	m.CTUR.Update()
	m.CTDR.Pos = pixel.V(m.StepH, -m.StepV)
	m.CTDR.Scalar = pixel.V(1.4, 1.4)
	m.CTDR.Update()
	m.CTDL.Pos = pixel.V(-m.StepH, -m.StepV)
	m.CTDL.Scalar = pixel.V(1.4, 1.4)
	m.CTDL.Update()
	m.STU.Pos = pixel.V(0., m.StepV)
	m.STU.Scalar = pixel.V(1.4 * m.StepH * 0.1725, 1.4)
	m.STU.Update()
	m.STR.Pos = pixel.V(m.StepH, 0.)
	m.STR.Scalar = pixel.V(1.4, 1.4 * m.StepV * 0.1725)
	m.STR.Update()
	m.STD.Pos = pixel.V(0., -m.StepV)
	m.STD.Scalar = pixel.V(1.4 * m.StepH * 0.1725, 1.4)
	m.STD.Update()
	m.STL.Pos = pixel.V(-m.StepH, 0.)
	m.STL.Scalar = pixel.V(1.4, 1.4 * m.StepV * 0.1725)
	m.STL.Update()
	m.Center.Scalar = pixel.V(1.4 * m.StepH * 0.1725, 1.4 * m.StepV * 0.1725)
	m.Center.Update()
	m.Trans.Update()
}

func (m *DwarfMenu) UpdateItems(in *input.Input) {
	dir := -1
	if in.Get("menuUp").JustPressed() {
		in.Get("menuUp").Consume()
		dir = 0
	} else if in.Get("menuDown").JustPressed() {
		in.Get("menuDown").Consume()
		dir = 1
	}
	if dir != -1 {
		m.GetNextHover(dir, m.Hovered)
	} else if in.MouseMoved {
		//point := m.Trans.Mat.Unproject(in.World)
		//point.X += m.Rect.W() * 0.5
		//point.X += m.Rect.H() * 0.5
		for i, item := range m.Items {
			if !item.Hovered && !item.Disabled && !item.NoHover {
				b := item.Text.BoundsOf(item.Raw)
				point := in.World
				point.X -= b.W() * 0.5
				point.Y -= b.H() * 2.
				if util.PointInside(point, b, item.Transform.Mat) {
					m.Items[m.Hovered].Hovered = false
					item.Hovered = true
					m.Hovered = i
					sfx.SoundPlayer.PlaySound("click", 2.0)
				}
			}
		}
	}
	if in.Get("menuBack").JustPressed() {
		m.Close()
		in.Get("menuBack").Consume()
	} else if in.Get("menuSelect").JustPressed() {
		if m.Items[m.Hovered].clickFn != nil {
			m.Items[m.Hovered].clickFn()
		}
		in.Get("menuSelect").Consume()
	} else if in.Get("click").JustPressed() {
		if m.Items[m.Hovered].clickFn != nil {
			m.Items[m.Hovered].clickFn()
		}
		in.Get("click").Consume()
	} else if in.Get("menuRight").JustPressed() {
		if m.Items[m.Hovered].rightFn != nil {
			m.Items[m.Hovered].rightFn()
		}
		in.Get("menuRight").Consume()
	} else if in.Get("menuLeft").JustPressed() {
		if m.Items[m.Hovered].leftFn != nil {
			m.Items[m.Hovered].leftFn()
		}
		in.Get("menuLeft").Consume()
	}
	for _, item := range m.Items {
		//point := m.Trans.Mat.Unproject(world)
		item.Transform.UIZoom = m.Cam.GetZoomScale()
		item.Transform.UIPos = m.Cam.APos
		item.Update()
	}
}

func (m *DwarfMenu) GetNextHover(dir, curr int) {
	nextI := curr
	if dir == 0 {
		nextI += len(m.Items) - 1
	} else {
		nextI++
	}
	nextI %= len(m.Items)
	next := m.Items[nextI]
	if next.Disabled || next.NoHover {
		m.GetNextHover(dir, nextI)
	} else {
		m.Items[m.Hovered].Hovered = false
		next.Hovered = true
		m.Hovered = nextI
		sfx.SoundPlayer.PlaySound("click", 2.0)
	}
}

func (m *DwarfMenu) Draw(target pixel.Target) {
	inner.Draw(target, m.Center.Mat)
	sideH.Draw(target, m.STU.Mat)
	sideV.Draw(target, m.STR.Mat)
	sideH.Draw(target, m.STD.Mat)
	sideV.Draw(target, m.STL.Mat)
	corner.Draw(target, m.CTUL.Mat)
	corner.Draw(target, m.CTUR.Mat)
	corner.Draw(target, m.CTDR.Mat)
	corner.Draw(target, m.CTDL.Mat)
	if !m.closing && m.opened {
		for _, item := range m.Items {
			item.Draw(target)
		}
	}
}