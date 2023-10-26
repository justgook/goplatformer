package ui

import (
	"image"
	"math"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

type children struct {
	absolute                 bool
	item                     *View
	bounds                   image.Rectangle
	isButtonPressed          bool
	isMouseLeftButtonHandler bool
	isMouseEntered           bool
	handledTouchID           ebiten.TouchID
	swipe
}

type swipe struct {
	downX, downY int
	upX, upY     int
	downTime     time.Time
	upTime       time.Time
	swipeDir     SwipeDirection
	swipeTouchID ebiten.TouchID
}

func (c *children) HandleJustPressedTouchID(
	frame *image.Rectangle, touchID ebiten.TouchID, x, y int) bool {
	var result = false
	if c.checkButtonHandlerStart(frame, touchID, x, y) {
		result = true
	}
	if !result && c.checkTouchHandlerStart(frame, touchID, x, y) {
		result = true
	}
	c.checkSwipeHandlerStart(frame, touchID, x, y)
	return result
}

func (c *children) HandleJustReleasedTouchID(
	frame *image.Rectangle, touchID ebiten.TouchID, x, y int) {
	c.checkTouchHandlerEnd(frame, touchID, x, y)
	c.checkButtonHandlerEnd(frame, touchID, x, y)
	c.checkSwipeHandlerEnd(frame, touchID, x, y)
}

func (c *children) checkTouchHandlerStart(frame *image.Rectangle, touchID ebiten.TouchID, x, y int) bool {
	touchHandler, ok := c.item.Handler.(TouchHandler)
	if !ok {
		return false
	}
	if !isInside(frame, x, y) {
		return false
	}
	if touchHandler.HandleJustPressedTouchID(touchID, x, y) {
		c.handledTouchID = touchID
		return true
	}
	return false
}

func (c *children) checkTouchHandlerEnd(_ *image.Rectangle, touchID ebiten.TouchID, x, y int) {
	touchHandler, ok := c.item.Handler.(TouchHandler)
	if !ok {
		return
	}
	if c.handledTouchID == touchID {
		touchHandler.HandleJustReleasedTouchID(touchID, x, y)
		c.handledTouchID = -1
	}

}

func (c *children) checkSwipeHandlerStart(frame *image.Rectangle, touchID ebiten.TouchID, x, y int) bool {
	if _, ok := c.item.Handler.(SwipeHandler); !ok {
		return false
	}
	if isInside(frame, x, y) {
		c.swipeTouchID = touchID
		c.swipe.downTime = time.Now()
		c.swipe.downX, c.swipe.downY = x, y
		return true
	}

	return false
}

func (c *children) checkSwipeHandlerEnd(_ *image.Rectangle, touchID ebiten.TouchID, x, y int) bool {
	swipeHandler, ok := c.item.Handler.(SwipeHandler)
	if !ok {
		return false
	}
	if c.swipeTouchID != touchID {
		return false
	}
	c.swipeTouchID = -1
	c.upTime = time.Now()
	c.upX, c.upY = x, y
	if c.checkSwipe() {
		swipeHandler.HandleSwipe(c.swipeDir)
		return true
	}

	return false
}

const swipeThresholdDist = 50.
const swipeThresholdTime = time.Millisecond * 300

func (c *children) checkSwipe() bool {
	dur := c.upTime.Sub(c.downTime)
	if dur > swipeThresholdTime {
		return false
	}

	deltaX := float64(c.downX - c.upX)
	if math.Abs(deltaX) >= swipeThresholdDist {
		if deltaX > 0 {
			c.swipeDir = SwipeDirectionLeft
		} else {
			c.swipeDir = SwipeDirectionRight
		}
		return true
	}

	deltaY := float64(c.downY - c.upY)
	if math.Abs(deltaY) >= swipeThresholdDist {
		if deltaY > 0 {
			c.swipeDir = SwipeDirectionUp
		} else {
			c.swipeDir = SwipeDirectionDown
		}
		return true
	}

	return false
}

func (c *children) checkButtonHandlerStart(frame *image.Rectangle, touchID ebiten.TouchID, x, y int) bool {
	button, ok := c.item.Handler.(ButtonHandler)
	if !ok {
		return false
	}

	if button2, ok2 := c.item.Handler.(NotButton); ok2 && !button2.IsButton() {
		return false
	}
	if isInside(frame, x, y) {
		if !c.isButtonPressed {
			c.isButtonPressed = true
			c.handledTouchID = touchID
			button.HandlePress(x, y, touchID)
		}
		return true
	}

	if c.handledTouchID == touchID {
		c.handledTouchID = -1
	}

	return false
}

func (c *children) checkButtonHandlerEnd(frame *image.Rectangle, touchID ebiten.TouchID, x, y int) {
	button, ok := c.item.Handler.(ButtonHandler)
	if !ok || c.handledTouchID != touchID || !c.isButtonPressed {
		return
	}
	c.isButtonPressed = false
	c.handledTouchID = -1
	if x == 0 && y == 0 {
		button.HandleRelease(x, y, false)
	} else {
		button.HandleRelease(x, y, !isInside(frame, x, y))
	}

}
