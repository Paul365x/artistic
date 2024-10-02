// package notify contains elements that tell the user what is happening
// this is a progress wheel whose concept is based on old style radar screen
package notify

import (
	"image/color"
	"math"
	"time"

	"github.com/crazy3lf/colorconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
)

type Progress_type struct {
	animate      *fyne.Animation  // ptr to animation object
	current      int              // leftmost index into lines matrix
	drawn        int              // number of lines drawn
	face         *canvas.Circle   // background circle
	lines        [][]*canvas.Line // matrix of lines: phases by steps
	steps        int              // number of steps/lines in a phase
	window       fyne.Window      // Window to overlay
	Border_width float32          // width of border circle for face
	Degrees      float64          // which hsl color we are up to
	Phases       float32          // number of segments/colors on circle
	PopUp        *fyne.Container  // container that is our popup
	Radius       float32          // radius of face
	Wayback      int              // number of phases back to de-emphasize
	Timeout      time.Duration    // time to run the animation before timing out
}

var Progress Progress_type

// calc_pos calculates the coordinates of the line end on the circle.
func (l *Progress_type) calc_pos(phases, nibbles float32) fyne.Position {
	arc_len := (2 * math.Pi * l.Radius) / l.Phases
	phase_angle := arc_len / l.Radius
	nibble_len := float32(1.0)
	nibble_angle := nibble_len / l.Radius

	angle_offset := float64((phase_angle * phases) + (nibble_angle * nibbles))
	length := l.Radius - l.Border_width
	x := l.Radius + (length * float32(math.Cos(angle_offset)))
	y := l.Radius + (length * float32(math.Sin(angle_offset)))
	return fyne.NewPos(x, y)
}

// SetupProgress initialises the progress popup widget
// should setup before each run
func (l *Progress_type) SetupProgress(w fyne.Window) *fyne.Container {

	Progress = Progress_type{
		current:      0,
		drawn:        0,
		window:       w,
		Border_width: 2.0,
		Degrees:      0.0,
		Phases:       160,
		Radius:       25.0,
		Wayback:      90,
		Timeout:      5 * time.Minute,
	}
	// setup the background or face                                                          // overcome rounding errors
	l.face = canvas.NewCircle(color.White)
	l.face.StrokeColor = color.Black
	l.face.StrokeWidth = l.Border_width
	l.PopUp = container.NewWithoutLayout(
		l.face,
	)

	// setup the lines that do the animation and color changes
	l.steps = int(math.Round(float64((2.0 * math.Pi * l.Radius) / l.Phases))) // this gives us a step size of 180th of circle
	r := int(l.Phases)
	l.lines = make([][]*canvas.Line, r)
	for i := range r {
		l.lines[i] = make([]*canvas.Line, l.steps)
		for j := range l.steps {
			l.lines[i][j] = canvas.NewLine(color.White)
			l.lines[i][j].Position1 = fyne.NewPos(l.Radius, l.Radius)
			l.lines[i][j].Position2 = l.calc_pos(float32(i), float32(j))
			l.PopUp.Add(l.lines[i][j])
		}
	}
	fyne.OverlayStack.Add(w.Canvas().Overlays(), l.PopUp)
	return l.PopUp
}

// RunProgress lines the elements up and starts the animation of the progress gizmo
func (l *Progress_type) RunProgress(p *fyne.Container) {

	// set up the face and progress
	w := l.window
	s := w.Canvas().Size()
	tly := (s.Height / 2) - l.Radius
	tlx := (s.Width / 2) - l.Radius
	loc := p.Position().AddXY(tlx, tly)
	l.face.Position1 = fyne.NewPos(0, 0)
	l.face.Position2 = fyne.NewPos(l.Radius*2, l.Radius*2)

	// setup the popup
	p.Resize(fyne.NewSize(l.Radius*2, l.Radius*2))
	p.Move(loc)
	p.Show()

	// create and start the animation
	l.animate = fyne.NewAnimation(l.Timeout, func(err float32) {
		// the animation tick supplies a float between 0 and 1. 1 means we are finished so cleanup
		if err == 1.0 {
			l.animate.Stop()
			p.Hide()
			fyne.OverlayStack.Remove(w.Canvas().Overlays(), l.PopUp)
			return
		}

		// current controls which phase or segment of the circle we paint
		// degrees controls the color in HSL color space
		history := l.current - l.Wayback - 1
		if history < 0 {
			history = int(l.Phases) + history
		}
		l.current++
		l.Degrees++
		l.drawn++
		if l.current >= int(l.Phases) {
			l.current = 0
		}
		if l.Degrees >= 360.0 {
			l.Degrees = 0.0
		}
		// only want to de-emphasize after the lines are coloured
		if l.drawn > l.Wayback {
			// de-emphasize previous phases gradually
			for range l.Wayback {
				history++
				if history >= int(l.Phases) {
					history = 0
				}
				for s := range l.steps {
					r, g, b, a := l.lines[history][s].StrokeColor.RGBA()
					h, sat, lum := colorconv.RGBToHSL(uint8(r), uint8(g), uint8(b))
					sat -= 0.005
					lum += 0.002
					rx, gx, bx, _ := colorconv.HSLToRGB(h, sat, lum)
					l.lines[history][s].StrokeColor = color.RGBA{
						R: uint8(rx),
						G: uint8(gx),
						B: uint8(bx),
						A: uint8(a),
					}
					l.lines[l.current][s].Refresh()
				}
			}
		}

		// Paint the new segment/phase
		for s := range l.steps {
			r, g, b, _ := colorconv.HSLToRGB(l.Degrees, 0.9, 0.5)
			l.lines[l.current][s].StrokeColor = color.RGBA{R: r, G: g, B: b, A: 255}
			l.lines[l.current][s].Refresh()
		}
	})
	l.animate.Start()
}

func (l *Progress_type) Start(w fyne.Window) {
	l.SetupProgress(w)
	l.RunProgress(l.PopUp)
}

func (l *Progress_type) Stop() {
	w := l.window
	l.animate.Stop()
	l.PopUp.Hide()
	fyne.OverlayStack.Remove(w.Canvas().Overlays(), l.PopUp)
}
