package main

/*
	Zooming, panning, dragging, etc; anything that handles the representation of the game on the screen as opposed to the window itself.
*/

var (
// when dragging we set the xpos,ypos at that moment, and then calculate the drag from the difference
// DRAG_XPOS_ORIGIN float64
// DRAG_YPOS_ORIGIN float64
)

const (
	ZOOM_IN  = -1
	ZOOM_OUT = 1

	ZOOM_STEP_SIZE = 0.02
)

// zoom in or out
func zoom(direction int, step_size float64) {

	if direction == ZOOM_IN {
		ZOOM = max(ZOOM-step_size, 0.0)
		return
	}

	if direction == ZOOM_OUT {
		ZOOM = min(ZOOM+step_size, 2.0)
		return
	}
}

// move the screen
func pan(xpan, ypan float64) {
	if abs(xpan) > 0.0 {
		t := ZOOM * xpan
		if t < 0.0 && t > -MIN_PAN_DISTANCE {
			t = -MIN_PAN_DISTANCE
		}
		if t > 0.0 && t < MIN_PAN_DISTANCE {
			t = MIN_PAN_DISTANCE
		}
		X_TRANSLATE += t
	}
	if abs(ypan) > 0.0 {
		t := ZOOM * ypan
		if t < 0.0 && t > -MIN_PAN_DISTANCE {
			t = -MIN_PAN_DISTANCE
		}
		if t > 0.0 && t < MIN_PAN_DISTANCE {
			t = MIN_PAN_DISTANCE
		}
		Y_TRANSLATE += t
	}
}

// keep track of how much the cursor moved from the starting position
func update_drag_delta() {
	if MOUSE_LB_PRESSED == false {
		return
	}

	// drag since last update
	xd := CURSOR_XPOS_SCREEN - DRAG_XPOS_ORIGIN
	yd := CURSOR_YPOS_SCREEN - DRAG_YPOS_ORIGIN

	// reset drag origin
	DRAG_XPOS_ORIGIN = CURSOR_XPOS_SCREEN
	DRAG_YPOS_ORIGIN = CURSOR_YPOS_SCREEN

	// add drag delta to drag globals. (drag delta's are zeroed when they are applied)
	DRAG_X += xd
	DRAG_Y += yd
}

// pan the screen based on the difference between the cursor position and the drag origin
func drag() {
	// get difference with last known location (this is set in update_drag_delta())
	xd := DRAG_X
	yd := DRAG_Y

	// normalize drag distance from pixels to uv (0.0 - 1.0)
	px := xd / float64(Width)
	py := yd / float64(Height)

	// don't pan if pan distance is too small
	// also: reset drag distance if we do decide to pan
	tx := ZOOM * abs(px)
	ty := ZOOM * abs(py)
	if tx > MIN_PAN_DISTANCE || tx < -MIN_PAN_DISTANCE {
		DRAG_X = 0.0
	} else {
		px = 0.0
	}
	if ty > MIN_PAN_DISTANCE || ty < -MIN_PAN_DISTANCE {
		DRAG_Y = 0.0
	} else {
		py = 0.0
	}

	// pan
	pan(-px, py)

}

// set the drag position to the current cursor location
func reset_drag_origin() {
	DRAG_XPOS_ORIGIN = CURSOR_XPOS_SCREEN
	DRAG_YPOS_ORIGIN = CURSOR_YPOS_SCREEN
}
