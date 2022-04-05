package menus

import (
	"dwarf-sweeper/internal/constants"
	"dwarf-sweeper/internal/menubox"
	"dwarf-sweeper/pkg/camera"
	"dwarf-sweeper/pkg/img"
	"dwarf-sweeper/pkg/input"
	"dwarf-sweeper/pkg/sfx"
	"dwarf-sweeper/pkg/transform"
	"dwarf-sweeper/pkg/util"
	"dwarf-sweeper/pkg/world"
	"fmt"
	"github.com/faiface/pixel"
	"math"
)

const (
	MaxLines = 10
)

var (
	arrow *pixel.Sprite
)


func Initialize() {
	DefaultDist = world.TileSize * 4.
	arrow = img.Batchers[constants.MenuSprites].Sprites["menu_arrow"]
}

type DwarfMenu struct {
	Key       string
	ItemMap   map[string]*Item
	Items     []*Item
	Hovered   int
	Top       int
	Title     bool
	Roll      bool
	TLines    int
	HideArrow bool

	Box  *menubox.MenuBox
	Hint *HintBox
	Tran *transform.Transform
	Cam  *camera.Camera

	backFn   func()
	openFn   func()
	closeFn  func()
	updateFn func(*input.Input)

	ArrowT *transform.Transform
}

func New(key string, cam *camera.Camera) *DwarfMenu {
	tran := transform.New()
	hint := NewHintBox("", cam)
	hint.Box.SetEntry(menubox.Left)
	AT := transform.New()
	return &DwarfMenu{
		Key:     key,
		ItemMap: map[string]*Item{},
		Items:   []*Item{},
		Box:     menubox.NewBox(&cam.APos, 1.0),
		Hint:    hint,
		Tran:    tran,
		Cam:     cam,
		ArrowT:  AT,
	}
}

func (m *DwarfMenu) AddItem(key, raw string, right bool) *Item {
	if _, ok := m.ItemMap[key]; ok {
		panic(fmt.Errorf("menu '%s' already has item '%s'", m.Key, key))
	}
	item := NewItem(key, raw, right)
	m.ItemMap[key] = item
	m.Items = append(m.Items, item)
	return item
}

func (m *DwarfMenu) InsertItem(key, raw, after string, right bool) *Item {
	if _, ok := m.ItemMap[key]; ok {
		panic(fmt.Errorf("menu '%s' already has item '%s'", m.Key, key))
	}
	item := NewItem(key, raw, right)
	m.ItemMap[key] = item
	i := 0
	for j, itemAfter := range m.Items {
		if itemAfter.Key == after {
			i = j+1
			break
		}
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
	m.Box.Open()
	m.setHover(-1)
	if m.openFn != nil {
		m.openFn()
	}
}

func (m *DwarfMenu) IsOpen() bool {
	return m.Box.IsOpen()
}

func (m *DwarfMenu) IsClosed() bool {
	return m.Box.IsClosed()
}

func (m *DwarfMenu) Close() {
	m.Box.Close()
	if m.closeFn != nil {
		m.closeFn()
	}
}

func (m *DwarfMenu) CloseInstant() {
	m.Box.CloseInstant()
	if m.closeFn != nil {
		m.closeFn()
	}
}

func (m *DwarfMenu) Update(in *input.Input) {
	if m.Box.IsOpen() && in != nil {
		m.UpdateView(in)
	}
	m.UpdateSize()
	m.Box.Update()
	m.UpdateTransforms()
	if m.Box.IsOpen() && in != nil {
		m.UpdateItems(in)
	}
}

func (m *DwarfMenu) UpdateView(in *input.Input) {
	if in.Get("scrollUp").JustPressed() {
		m.menuUp()
	}
	if in.Get("scrollDown").JustPressed() {
		m.menuDown()
	}
	dir := -1
	if in.Get("menuUp").JustPressed() || in.Get("menuUp").Repeated() {
		dir = 0
	} else if in.Get("menuDown").JustPressed() || in.Get("menuDown").Repeated() {
		dir = 1
	} else if in.Get("menuRight").JustPressed() || in.Get("menuRight").Repeated() {
		dir = 2
	} else if in.Get("menuLeft").JustPressed() || in.Get("menuLeft").Repeated() {
		dir = 3
	}
	if dir != -1 {
		m.GetNextHover(dir, m.Hovered, in)
	} else if in.MouseMoved {
		for i, item := range m.Items {
			if !item.Hovered && !item.Disabled && !item.NoHover && !item.noShowT {
				b := item.Text.Text.BoundsOf(item.Raw)
				if !m.HideArrow {
					b.Min.X -= 30. / constants.ActualMenuSize
				}
				point := in.World
				if item.Right {
					point.X += 15.
				} else {
					point.X -= b.W() * 0.5 * constants.ActualMenuSize
				}
				point.Y -= b.H() * 1.45 * constants.ActualMenuSize
				if util.PointInside(point, b, item.Text.Transform.Mat) {
					m.setHover(i)
				}
			}
		}
	}
}

func (m *DwarfMenu) UpdateSize() {
	minWidth := 8.
	minHeight := 8.
	sameLine := false
	lines := 0
	tLines := 0
	for i, item := range m.Items {
		if item.Ignore {
			item.noShowT = true
			continue
		}
		visible := (m.Title && i == 0) || (tLines >= m.Top && lines < MaxLines)
		//if (m.Title && i == 0) || (tLines >= m.Top && lines < MaxLines) {
		item.CurrLine = tLines
		bW := item.Text.Text.Bounds().W() * constants.ActualMenuSize
		sW := 0.
		if !item.Right && i+1 < len(m.Items) && m.Items[i+1].Right {
			next := m.Items[i+1]
			sW = (next.Text.Text.Bounds().W() + next.Text.Text.BoundsOf("   ").W()) * constants.ActualMenuSize
			sameLine = true
		}
		minWidth = math.Max(bW+sW, minWidth)
		if !sameLine {
			if visible {
				minHeight += item.Text.Text.LineHeight * constants.ActualMenuSize
				lines++
			}
			tLines++
		}
		sameLine = false
		item.noShowT = !visible
		//} else {
		//	item.CurrLine = tLines
		//	if !item.Right && i+1 < len(m.Items) && m.Items[i+1].Right {
		//		sameLine = true
		//	}
		//	if !sameLine {
		//		tLines++
		//	}
		//	sameLine = false
		//	item.noShowT = true
		//}
	}
	m.TLines = tLines
	minWidth += 15.
	if !m.HideArrow {
		minWidth += 30.
	}
	m.Box.SetSize(pixel.R(0., 0., minWidth, minHeight))
	line := 0
	for i, item := range m.Items {
		if !item.noShowT {
			if item.Right {
				item.Text.SetPos(pixel.V(minWidth*0.5 - 10., minHeight*0.5 - float64(line+1)*item.Text.Text.LineHeight * constants.ActualMenuSize))
			} else {
				nextY := minHeight*0.5 - float64(line+1)*item.Text.Text.LineHeight * constants.ActualMenuSize
				var nextX float64
				if !m.HideArrow {
					nextX = minWidth*-0.5 + 20.
				} else {
					nextX = minWidth*-0.5 + 5.
				}
				item.Text.SetPos(pixel.V(nextX, nextY))
			}
			if item.Right || i >= len(m.Items)-1 || !m.Items[i+1].Right {
				line++
			}
		}
	}
}

func (m *DwarfMenu) UpdateTransforms() {
	if m.Cam != nil {
		m.Tran.UIZoom = m.Cam.GetZoomScale()
		m.Tran.UIPos = m.Cam.APos
		m.ArrowT.UIZoom = m.Cam.GetZoomScale()
		m.ArrowT.UIPos = m.Cam.APos
	}
	m.Tran.Update()
	if m.Hovered != -1 {
		hovered := m.Items[m.Hovered]
		m.ArrowT.Pos.Y = hovered.Text.Transform.Pos.Y + hovered.Text.Height*0.25
		if hovered.Right {
			m.ArrowT.Pos.X = hovered.Text.Transform.Pos.X - hovered.Text.Width - 10.
		} else {
			m.ArrowT.Pos.X = hovered.Text.Transform.Pos.X - 10.
		}
		m.ArrowT.Scalar = pixel.V(1.4, 1.4)
		m.ArrowT.Update()
		if hovered.Hint != "" && m.Box.IsOpen() {
			m.Hint.SetText(hovered.Hint)
			m.Hint.Tran.Pos.X = m.Box.STR.Pos.X + m.Hint.Box.Rect.W() * 0.5
			m.Hint.Tran.Pos.Y = m.ArrowT.Pos.Y
			m.Hint.Update()
			m.Hint.Display = true
		} else {
			m.Hint.SetText("")
			m.Hint.Update()
			m.Hint.Display = false
		}
	} else {
		m.Hint.SetText("")
		m.Hint.Update()
		m.Hint.Display = false
	}
}

func (m *DwarfMenu) UpdateItems(in *input.Input) {
	if in.Get("menuBack").JustPressed() {
		m.Back()
		in.Get("menuBack").Consume()
	} else if in.Get("menuSelect").JustPressed() && m.Hovered != -1 {
		if m.Items[m.Hovered].clickFn != nil {
			m.Items[m.Hovered].clickFn()
		}
		in.Get("menuSelect").Consume()
	} else if in.Get("click").JustPressed() && m.Hovered != -1 {
		if m.Items[m.Hovered].clickFn != nil {
			m.Items[m.Hovered].clickFn()
		}
		in.Get("click").Consume()
	} else if in.Get("menuRight").JustPressed() && m.Hovered != -1 {
		if m.Items[m.Hovered].rightFn != nil {
			m.Items[m.Hovered].rightFn()
		}
		in.Get("menuRight").Consume()
	} else if in.Get("menuLeft").JustPressed() && m.Hovered != -1 {
		if m.Items[m.Hovered].leftFn != nil {
			m.Items[m.Hovered].leftFn()
		}
		in.Get("menuLeft").Consume()
	}
	for _, item := range m.Items {
		item.Update()
	}
}

func (m *DwarfMenu) setHover(nextI int) {
	if nextI == -1 {
		hover := false
		for i, item := range m.Items {
			if !hover && !item.disabled && !item.NoHover {
				m.setHover(i)
				hover = true
			} else {
				item.Hovered = false
			}
		}
	} else {
		if m.Hovered != -1 {
			prev := m.Items[m.Hovered]
			prev.Hovered = false
			if prev.unHoverFn != nil {
				prev.unHoverFn()
			}
		}
		next := m.Items[nextI]
		next.Hovered = true
		m.Hovered = nextI
		sfx.SoundPlayer.PlaySound("click", 2.0)
		if next.hoverFn != nil {
			next.hoverFn()
		}
		m.setTop(next.CurrLine)
	}
}

func (m *DwarfMenu) setTop(line int) {
	if line < m.Top {
		m.Top = line
	} else if m.Title && line >= m.Top+MaxLines-1 {
		m.Top = line - MaxLines + 2
	} else if line >= m.Top+MaxLines {
		m.Top = line - MaxLines + 1
	}
}

func (m *DwarfMenu) menuUp() {
	m.Top--
	if m.Top < 0 {
		m.Top = 0
	}
}

func (m *DwarfMenu) menuDown() {
	m.Top++
	if m.Top > m.TLines-MaxLines+1 {
		m.Top = m.TLines - MaxLines + 1
	}
}

func (m *DwarfMenu) UnhoverAll() {
	for _, item := range m.Items {
		if item.Hovered && item.unHoverFn != nil {
			item.unHoverFn()
		}
		item.Hovered = false
	}
	m.Hovered = -1
}

func (m *DwarfMenu) GetNextHover(dir, curr int, in *input.Input) {
	if curr == -1 {
		m.setHover(-1)
	}
	if dir == 0 || dir == 1 {
		r := false
		if curr != -1 {
			r = m.Items[curr].Right
		}
		m.GetNextHoverVert(dir, curr, r, in)
	} else {
		m.GetNextHoverHor(dir, curr, in)
	}
}

func (m *DwarfMenu) GetNextHoverHor(dir, curr int, in *input.Input) {
	this := m.Items[curr]
	nextI := -1
	if dir == 2 && !this.Right && curr < len(m.Items)-1 {
		nextI = curr + 1
	} else if dir == 3 && this.Right && curr > 0 {
		nextI = curr - 1
	}
	if nextI != -1 {
		next := m.Items[nextI]
		if next.Right != this.Right && !next.Disabled && !next.NoHover && !next.noShowT {
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
	if next.Disabled || next.NoHover || next.Ignore || next.Right != right {
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
	m.Box.Draw(target)
	if m.Box.IsOpen() {
		for _, item := range m.Items {
			item.Draw(target)
			if item.Hovered && !item.Ignore && !item.noShowT && !item.NoDraw && !m.HideArrow {
				arrow.Draw(target, m.ArrowT.Mat)
			}
		}
	}
	m.Hint.Draw(target)
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

func (m *DwarfMenu) SetOpenFn(fn func()) {
	m.openFn = fn
}

func (m *DwarfMenu) SetCloseFn(fn func()) {
	m.closeFn = fn
}
