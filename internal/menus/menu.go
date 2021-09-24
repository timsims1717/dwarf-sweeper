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
	arrow  *pixel.Sprite
)

func Initialize() {
	corner = img.Batchers[cfg.MenuSprites].Sprites["menu_corner"]
	sideV = img.Batchers[cfg.MenuSprites].Sprites["menu_side_v"]
	sideH = img.Batchers[cfg.MenuSprites].Sprites["menu_side_h"]
	inner = img.Batchers[cfg.MenuSprites].Sprites["menu_inner"]
	arrow = img.Batchers[cfg.MenuSprites].Sprites["menu_arrow"]
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
	Roll    bool
	backFn  func()

	Center *transform.Transform
	CTUL   *transform.Transform
	CTUR   *transform.Transform
	CTDR   *transform.Transform
	CTDL   *transform.Transform
	STU    *transform.Transform
	STR    *transform.Transform
	STD    *transform.Transform
	STL    *transform.Transform
	ArrowT *transform.Transform
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
	AT := transform.NewTransform()
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
		ArrowT:  AT,
	}
}

func (m *DwarfMenu) AddItem(key, raw string) *Item {
	if _, ok := m.ItemMap[key]; ok {
		panic(fmt.Errorf("menu '%s' already has item '%s'", m.Key, key))
	}
	item := NewItem(key, raw)
	m.ItemMap[key] = item
	m.Items = append(m.Items, item)
	return item
}

func (m *DwarfMenu) InsertItem(key, raw string, i int) *Item {
	if _, ok := m.ItemMap[key]; ok {
		panic(fmt.Errorf("menu '%s' already has item '%s'", m.Key, key))
	}
	item := NewItem(key, raw)
	m.ItemMap[key] = item
	if i < 0 {
		i = 0
	}
	if i >= len(m.Items) {
		m.Items = append(m.Items, item)
	} else {
		m.Items = append(m.Items[:i], append([]*Item{item}, m.Items[i:]...)...)
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
	delete(m.ItemMap, key)
}

func (m *DwarfMenu) Open() {
	m.Closed = false
	m.closing = false
	m.opened = false
	hover := false
	for i, item := range m.Items {
		if !hover && !item.disabled && !item.NoHover {
			m.setHover(i)
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
	sameLine := false
	for i, item := range m.Items {
		bW := item.Text.BoundsOf(item.Raw).W()
		sW := 0.
		if !item.Right && i+1 < len(m.Items) && m.Items[i+1].Right {
			next := m.Items[i+1]
			sW = next.Text.BoundsOf(next.Raw).W() + next.Text.BoundsOf("   ").W()
			sameLine = true
		}
		minWidth = math.Max((bW + sW) * 1.4, minWidth)
		if !sameLine {
			minHeight += item.Text.LineHeight
		}
		sameLine = false
	}
	minWidth = math.Floor(math.Max(minWidth + 30., m.Rect.W()))
	minHeight = math.Floor(math.Max(minHeight, m.Rect.H()))
	m.Rect = pixel.R(0., 0., minWidth, minHeight)
	line := 0
	for i, item := range m.Items {
		if item.Right {
			item.Transform.Pos.Y = minHeight*0.5 - float64(line+1)*item.Text.LineHeight
			item.Transform.Pos.X = minWidth * 0.5 - 10.
		} else {
			item.Transform.Pos.Y = minHeight*0.5 - float64(line+1)*item.Text.LineHeight
			item.Transform.Pos.X = minWidth*-0.5 + 20.
		}
		if item.Right || i >= len(m.Items)-1 || !m.Items[i+1].Right {
			line++
		}
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
		m.ArrowT.UIZoom = m.Cam.GetZoomScale()
		m.ArrowT.UIPos = m.Cam.APos
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
	hovered := m.Items[m.Hovered]
	m.ArrowT.Pos.Y = hovered.Transform.Pos.Y + hovered.Text.BoundsOf(hovered.Raw).H() * 0.5
	if hovered.Right {
		m.ArrowT.Pos.X = hovered.Transform.Pos.X - hovered.Text.BoundsOf(hovered.Raw).W() * 1.45 - 10.
	} else {
		m.ArrowT.Pos.X = hovered.Transform.Pos.X - 10.
	}
	m.ArrowT.Scalar = pixel.V(1.4, 1.4)
	m.ArrowT.Update()
	m.Trans.Update()
}

func (m *DwarfMenu) UpdateItems(in *input.Input) {
	dir := -1
	if in.Get("menuUp").JustPressed() {
		dir = 0
	} else if in.Get("menuDown").JustPressed() {
		dir = 1
	} else if in.Get("menuRight").JustPressed() {
		dir = 2
	} else if in.Get("menuLeft").JustPressed() {
		dir = 3
	}
	if dir != -1 {
		m.GetNextHover(dir, m.Hovered, in)
	} else if in.MouseMoved {
		//point := m.Trans.Mat.Unproject(in.World)
		//point.X += m.Rect.W() * 0.5
		//point.X += m.Rect.H() * 0.5
		for i, item := range m.Items {
			if !item.Hovered && !item.Disabled && !item.NoHover {
				b := item.Text.BoundsOf(item.Raw)
				point := in.World
				if item.Right {
					point.X += b.W() * 0.5
				} else {
					point.X -= b.W() * 0.5
				}
				point.Y -= b.H() * 2.
				if util.PointInside(point, b, item.Transform.Mat) {
					m.setHover(i)
				}
			}
		}
	}
	if in.Get("menuBack").JustPressed() {
		m.Back()
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

func (m *DwarfMenu) setHover(nextI int) {
	next := m.Items[nextI]
	m.Items[m.Hovered].Hovered = false
	next.Hovered = true
	m.Hovered = nextI
	sfx.SoundPlayer.PlaySound("click", 2.0)
}

func (m *DwarfMenu) GetNextHover(dir, curr int, in *input.Input) {
	if dir == 0 || dir == 1 {
		m.GetNextHoverVert(dir, curr, m.Items[curr].Right, in)
	} else {
		m.GetNextHoverHor(dir, curr, in)
	}
}

func (m *DwarfMenu) GetNextHoverHor(dir, curr int, in *input.Input) {
	this := m.Items[curr]
	nextI := -1
	if dir == 2 && !this.Right && curr < len(m.Items)-1 {
		nextI = curr+1
	} else if dir == 3 && this.Right && curr > 0 {
		nextI = curr-1
	}
	if nextI != -1 {
		next := m.Items[nextI]
		if next.Right != this.Right && !next.Disabled && !next.NoHover {
			m.setHover(nextI)
			if dir == 2 {
				in.Get("menuRight").Consume()
			} else {
				in.Get("menuLeft").Consume()
			}
		}
	}
}

func (m *DwarfMenu) GetNextHoverVert(dir, curr int, right bool, in *input.Input) {
	nextI := curr
	if dir == 0 {
		nextI--
	} else {
		nextI++
	}
	if !m.Roll && (nextI >= len(m.Items) || nextI < 0) {
		return
	}
	if nextI < 0 {
		nextI += len(m.Items)
	}
	nextI %= len(m.Items)
	next := m.Items[nextI]
	if next.Disabled || next.NoHover || next.Right != right {
		m.GetNextHoverVert(dir, nextI, right, in)
	} else {
		m.setHover(nextI)
		if dir == 0 {
			in.Get("menuUp").Consume()
		} else {
			in.Get("menuDown").Consume()
		}
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
			if item.Hovered {
				arrow.Draw(target, m.ArrowT.Mat)
			}
		}
	}
}

func (m *DwarfMenu) Back() {
	if m.backFn != nil {
		m.backFn()
	} else {
		m.Close()
	}
}

func (m *DwarfMenu) SetBackFn(fn func()) {
	m.backFn = fn
}