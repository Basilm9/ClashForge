package clashgame

import (
	"fmt"
	"image/color"
	"math"
	"sort"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// TroopSelectionSystem manages troop selection and deployment
type TroopSelectionSystem struct {
	// Currently selected troop name
	SelectedTroop string
	
	// Available troop names
	TroopNames []string
	
	// Category filters (Common, Rare, Epic, Legendary)
	CurrentFilter string
	
	// For scrolling through troops
	ScrollIndex int
	
	// Visual properties
	UIHeight        int
	UIWidth         int
	UIX             int
	UIY             int
	CardWidth       int
	CardHeight      int
	CardSpacing     int
	MaxVisibleCards int
}

// NewTroopSelectionSystem creates a new selection system
func NewTroopSelectionSystem() *TroopSelectionSystem {
	system := &TroopSelectionSystem{
		SelectedTroop:   "",
		TroopNames:      make([]string, 0),
		CurrentFilter:   "All",
		ScrollIndex:     0,
		UIHeight:        80,
		UIWidth:         screenWidth,
		UIX:             0,
		UIY:             screenHeight - 80,
		CardWidth:       70,
		CardHeight:      70,
		CardSpacing:     5,
		MaxVisibleCards: 8,
	}
	
	// Populate troop names from template map
	system.ReloadTroopNames()
	
	// Select first troop by default if available
	if len(system.TroopNames) > 0 {
		system.SelectedTroop = system.TroopNames[0]
	}
	
	return system
}

// ReloadTroopNames updates the troop names list from the template map
func (ts *TroopSelectionSystem) ReloadTroopNames() {
	ts.TroopNames = make([]string, 0, len(TroopTemplateMap))
	
	// Add all troop names from template map
	for name := range TroopTemplateMap {
		// Skip NOTINUSE troops
		if !strings.Contains(name, "NOTINUSE") {
			ts.TroopNames = append(ts.TroopNames, name)
		}
	}
	
	// Sort alphabetically
	sort.Strings(ts.TroopNames)
}

// Update handles input and selection changes
func (ts *TroopSelectionSystem) Update() {
	// Handle scrolling through troops
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowRight) {
		ts.ScrollIndex = min(ts.ScrollIndex+1, max(0, len(ts.TroopNames)-ts.MaxVisibleCards))
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowLeft) {
		ts.ScrollIndex = max(0, ts.ScrollIndex-1)
	}
	
	// Handle category filtering with number keys
	if inpututil.IsKeyJustPressed(ebiten.Key0) {
		ts.CurrentFilter = "All"
		ts.ReloadTroopNames()
	}
	if inpututil.IsKeyJustPressed(ebiten.Key1) {
		ts.FilterByRarity("Common")
	}
	if inpututil.IsKeyJustPressed(ebiten.Key2) {
		ts.FilterByRarity("Rare")
	}
	if inpututil.IsKeyJustPressed(ebiten.Key3) {
		ts.FilterByRarity("Epic")
	}
	if inpututil.IsKeyJustPressed(ebiten.Key4) {
		ts.FilterByRarity("Legendary")
	}
	
	// Handle numeric troop selection (5-9 keys for first 5 troops)
	for i := 0; i < 5; i++ {
		if inpututil.IsKeyJustPressed(ebiten.Key(int(ebiten.Key5) + i)) {
			idx := ts.ScrollIndex + i
			if idx < len(ts.TroopNames) {
				ts.SelectedTroop = ts.TroopNames[idx]
			}
		}
	}
	
	// Handle mouse selection in the UI area
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		
		// Check if click is in the UI area
		if y >= ts.UIY && y <= ts.UIY+ts.UIHeight {
			// Calculate which card was clicked
			for i := 0; i < min(ts.MaxVisibleCards, len(ts.TroopNames)-ts.ScrollIndex); i++ {
				cardX := ts.UIX + i*(ts.CardWidth+ts.CardSpacing) + ts.CardSpacing
				if x >= cardX && x <= cardX+ts.CardWidth {
					troopIndex := ts.ScrollIndex + i
					if troopIndex < len(ts.TroopNames) {
						ts.SelectedTroop = ts.TroopNames[troopIndex]
					}
					break
				}
			}
		}
	}
}

// Draw renders the troop selection UI
func (ts *TroopSelectionSystem) Draw(screen *ebiten.Image) {
	// Draw UI background
	ebitenutil.DrawRect(
		screen,
		float64(ts.UIX),
		float64(ts.UIY),
		float64(ts.UIWidth),
		float64(ts.UIHeight),
		color.RGBA{40, 40, 40, 220},
	)
	
	// Draw visible troop cards
	for i := 0; i < min(ts.MaxVisibleCards, len(ts.TroopNames)-ts.ScrollIndex); i++ {
		troopIndex := ts.ScrollIndex + i
		troopName := ts.TroopNames[troopIndex]
		
		// Card position
		cardX := ts.UIX + i*(ts.CardWidth+ts.CardSpacing) + ts.CardSpacing
		cardY := ts.UIY + ts.CardSpacing
		
		// Draw card background (with highlight if selected)
		cardColor := color.RGBA{60, 60, 60, 255}
		if troopName == ts.SelectedTroop {
			cardColor = color.RGBA{100, 150, 200, 255}
		}
		
		// Get troop template for info
		template, exists := TroopTemplateMap[troopName]
		if !exists {
			continue
		}
		
		// Color based on rarity
		switch template.Rarity {
		case "Common":
			// Gray
			if troopName != ts.SelectedTroop {
				cardColor = color.RGBA{100, 100, 100, 255}
			}
		case "Rare":
			// Blue
			if troopName != ts.SelectedTroop {
				cardColor = color.RGBA{60, 100, 180, 255}
			}
		case "Epic":
			// Purple
			if troopName != ts.SelectedTroop {
				cardColor = color.RGBA{140, 60, 180, 255}
			}
		case "Legendary":
			// Gold
			if troopName != ts.SelectedTroop {
				cardColor = color.RGBA{200, 150, 50, 255}
			}
		}
		
		ebitenutil.DrawRect(
			screen,
			float64(cardX),
			float64(cardY),
			float64(ts.CardWidth),
			float64(ts.CardHeight),
			cardColor,
		)
		
		// Draw troop icon/circle
		iconSize := float64(ts.CardWidth) * 0.4
		iconX := float64(cardX) + float64(ts.CardWidth)*0.5
		iconY := float64(cardY) + float64(ts.CardHeight)*0.4
		
		// Determine base color by rarity
		var troopIconColor color.RGBA
		switch template.Rarity {
		case "Common":
			troopIconColor = color.RGBA{220, 220, 220, 255}
		case "Rare":
			troopIconColor = color.RGBA{100, 180, 255, 255}
		case "Epic":
			troopIconColor = color.RGBA{180, 100, 255, 255}
		case "Legendary":
			troopIconColor = color.RGBA{255, 200, 80, 255}
		default:
			troopIconColor = color.RGBA{220, 220, 220, 255}
		}
		
		// Special visual indicators
		if template.FlyingHeight > 0 {
			// Draw a "wing" shape for flying troops
			ebitenutil.DrawLine(
				screen,
				iconX-iconSize*0.6, iconY-iconSize*0.2,
				iconX, iconY-iconSize*0.4,
				color.RGBA{220, 220, 220, 255},
			)
			ebitenutil.DrawLine(
				screen,
				iconX, iconY-iconSize*0.4,
				iconX+iconSize*0.6, iconY-iconSize*0.2,
				color.RGBA{220, 220, 220, 255},
			)
		}
		
		// Draw the troop icon
		ebitenutil.DrawCircle(
			screen,
			iconX,
			iconY,
			iconSize*0.5,
			troopIconColor,
		)
		
		// Add indicator for area damage
		if template.AreaDamageRadius > 0 {
			ebitenutil.DrawCircle(
				screen,
				iconX,
				iconY,
				iconSize*0.7,
				color.RGBA{255, 255, 255, 80},
			)
		}
		
		// Draw troop name
		displayName := GetTroopDisplayName(troopName)
		if len(displayName) > 10 {
			displayName = displayName[:9] + "."
		}
		ebitenutil.DebugPrintAt(
			screen,
			displayName,
			cardX+2,
			cardY+ts.CardHeight-15,
		)
		
		// Add hotkey indicator
		if i < 5 {
			hotkey := fmt.Sprintf("%d", i+5)
			ebitenutil.DebugPrintAt(
				screen,
				hotkey,
				cardX+ts.CardWidth-12,
				cardY+12,
			)
		}
	}
	
	// Draw scroll indicators if needed
	if ts.ScrollIndex > 0 {
		ebitenutil.DebugPrintAt(
			screen,
			"◀",
			ts.UIX+5,
			ts.UIY+ts.UIHeight/2-5,
		)
	}
	if ts.ScrollIndex+ts.MaxVisibleCards < len(ts.TroopNames) {
		ebitenutil.DebugPrintAt(
			screen,
			"▶",
			ts.UIX+ts.UIWidth-15,
			ts.UIY+ts.UIHeight/2-5,
		)
	}
	
	// Draw filter info
	filterText := fmt.Sprintf("Filter: %s [0-4]", ts.CurrentFilter)
	ebitenutil.DebugPrintAt(
		screen,
		filterText,
		ts.UIX+ts.UIWidth-150,
		ts.UIY+10,
	)
}

// FilterByRarity filters troops by rarity
func (ts *TroopSelectionSystem) FilterByRarity(rarity string) {
	// Set current filter
	ts.CurrentFilter = rarity
	
	// Create filtered list
	filteredNames := make([]string, 0)
	for name, template := range TroopTemplateMap {
		// Skip NOTINUSE troops
		if strings.Contains(name, "NOTINUSE") {
			continue
		}
		
		// Add if matches rarity
		if template.Rarity == rarity {
			filteredNames = append(filteredNames, name)
		}
	}
	
	// Sort filtered names
	sort.Strings(filteredNames)
	ts.TroopNames = filteredNames
	
	// Reset scroll and update selection if needed
	ts.ScrollIndex = 0
	if len(ts.TroopNames) > 0 {
		ts.SelectedTroop = ts.TroopNames[0]
	} else {
		ts.SelectedTroop = ""
	}
}

// Add this function to deploy the selected troop at a specified location
func (ts *TroopSelectionSystem) DeploySelectedTroop(game *Game, col, row, team int) {
	if ts.SelectedTroop == "" {
		return
	}
	
	// Get grid position
	pos := game.Grid.CellToPosition(col, row)
	
	// Spawn the troop from the template
	err := SpawnExtendedTroop(ts.SelectedTroop, pos.X, pos.Y, team, game)
	if err != nil {
		fmt.Printf("Error spawning troop: %v\n", err)
	}
}

// Helper min/max functions
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// EnhancedTroopDrawer is responsible for drawing troops with more visual details
type EnhancedTroopDrawer struct {
	// Reference to game for accessing templates
	Game *Game
}

// NewEnhancedTroopDrawer creates a new troop drawer
func NewEnhancedTroopDrawer(game *Game) *EnhancedTroopDrawer {
	return &EnhancedTroopDrawer{
		Game: game,
	}
}

// Helper functions to draw different shapes

// Draw a triangle (for flying troops)
func drawTriangle(screen *ebiten.Image, x, y, size float64, clr color.RGBA) {
	// Calculate vertices
	x1 := x
	y1 := y - size
	x2 := x - size*0.866 // cos(60°) = 0.866
	y2 := y + size*0.5   // sin(60°) = 0.5
	x3 := x + size*0.866
	y3 := y + size*0.5
	
	// Draw the triangle (three lines)
	ebitenutil.DrawLine(screen, x1, y1, x2, y2, clr)
	ebitenutil.DrawLine(screen, x2, y2, x3, y3, clr)
	ebitenutil.DrawLine(screen, x3, y3, x1, y1, clr)
	
	// Fill with semi-transparent color
	fillTriangle(screen, x1, y1, x2, y2, x3, y3, color.RGBA{clr.R, clr.G, clr.B, 150})
}

// Draw a square (for melee fighter troops)
func drawSquare(screen *ebiten.Image, x, y, size float64, clr color.RGBA) {
	ebitenutil.DrawRect(
		screen,
		x-size,
		y-size,
		size*2,
		size*2,
		clr,
	)
}

// Draw a diamond (for ranged troops)
func drawDiamond(screen *ebiten.Image, x, y, size float64, clr color.RGBA) {
	// Calculate vertices
	x1 := x
	y1 := y - size
	x2 := x + size
	y2 := y
	x3 := x
	y3 := y + size
	x4 := x - size
	y4 := y
	
	// Draw the diamond (four lines)
	ebitenutil.DrawLine(screen, x1, y1, x2, y2, clr)
	ebitenutil.DrawLine(screen, x2, y2, x3, y3, clr)
	ebitenutil.DrawLine(screen, x3, y3, x4, y4, clr)
	ebitenutil.DrawLine(screen, x4, y4, x1, y1, clr)
}

// Draw a hexagon (for tanky troops)
func drawHexagon(screen *ebiten.Image, x, y, size float64, clr color.RGBA) {
	// Calculate vertices
	vertices := make([]struct{ x, y float64 }, 6)
	for i := 0; i < 6; i++ {
		angle := float64(i) * (math.Pi / 3)
		vertices[i].x = x + size*math.Cos(angle)
		vertices[i].y = y + size*math.Sin(angle)
	}
	
	// Draw the hexagon (six lines)
	for i := 0; i < 6; i++ {
		next := (i + 1) % 6
		ebitenutil.DrawLine(
			screen,
			vertices[i].x, vertices[i].y,
			vertices[next].x, vertices[next].y,
			clr,
		)
	}
}

// FillTriangle fills a triangle with a color (simplified implementation)
func fillTriangle(screen *ebiten.Image, x1, y1, x2, y2, x3, y3 float64, clr color.RGBA) {
	// Simple implementation - draw horizontal lines between edges
	// Find top, middle, and bottom points
	pts := [][2]float64{{x1, y1}, {x2, y2}, {x3, y3}}
	
	// Sort by y-coordinate
	if pts[0][1] > pts[1][1] {
		pts[0], pts[1] = pts[1], pts[0]
	}
	if pts[1][1] > pts[2][1] {
		pts[1], pts[2] = pts[2], pts[1]
	}
	if pts[0][1] > pts[1][1] {
		pts[0], pts[1] = pts[1], pts[0]
	}
	
	// Extract sorted coordinates
	x1, y1 = pts[0][0], pts[0][1]
	x2, y2 = pts[1][0], pts[1][1]
	x3, y3 = pts[2][0], pts[2][1]
	
	// Draw horizontal lines for the top half
	if y2 > y1 {
		slope1 := (x2 - x1) / (y2 - y1)
		slope2 := (x3 - x1) / (y3 - y1)
		
		for y := math.Ceil(y1); y <= y2; y++ {
			startX := x1 + slope1*(y-y1)
			endX := x1 + slope2*(y-y1)
			
			if startX > endX {
				startX, endX = endX, startX
			}
			
			ebitenutil.DrawLine(
				screen,
				startX, y,
				endX, y,
				clr,
			)
		}
	}
	
	// Draw horizontal lines for the bottom half
	if y3 > y2 {
		slope1 := (x3 - x2) / (y3 - y2)
		slope2 := (x3 - x1) / (y3 - y1)
		
		for y := math.Ceil(y2); y <= y3; y++ {
			startX := x2 + slope1*(y-y2)
			endX := x1 + slope2*(y-y1)
			
			if startX > endX {
				startX, endX = endX, startX
			}
			
			ebitenutil.DrawLine(
				screen,
				startX, y,
				endX, y,
				clr,
			)
		}
	}
}

// DrawTroop draws a troop with its visual properties
func (etd *EnhancedTroopDrawer) DrawTroop(screen *ebiten.Image, troop *Troop) {
	if troop == nil {
		return
	}

	// Get the template for additional visual properties
	template := GetTroopTemplate(troop)
	if template == nil {
		return
	}

	// Draw the troop based on its type
	x, y := troop.Position.X, troop.Position.Y
	size := troop.Size

	// Draw the main shape
	switch template.Tribe {
	case "Air":
		drawHexagon(screen, x, y, size, troop.Color)
	case "Ground":
		drawSquare(screen, x, y, size, troop.Color)
	default:
		drawTriangle(screen, x, y, size, troop.Color)
	}

	// Draw health bar
	healthBarWidth := size * 1.2
	healthBarHeight := size * 0.1
	healthBarX := x - healthBarWidth/2
	healthBarY := y - size/2 - healthBarHeight - 2

	// Background (red)
	ebitenutil.DrawRect(screen, healthBarX, healthBarY, healthBarWidth, healthBarHeight, color.RGBA{200, 0, 0, 255})
	
	// Health (green)
	healthPercent := float64(troop.Health) / float64(troop.MaxHealth)
	ebitenutil.DrawRect(screen, healthBarX, healthBarY, healthBarWidth*healthPercent, healthBarHeight, color.RGBA{0, 200, 0, 255})
}