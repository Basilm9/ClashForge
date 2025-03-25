// grid_system.go
package clashgame

import (
	"encoding/csv"
	"fmt"
	"image/color"
	"io"
	"math"
	"os"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

// Cell types for different terrain
const (
	CellTypeGround = iota
	CellTypeWater
	CellTypeBridge
)

// TileType represents different types of tiles
const (
	TileEmpty = 0
	TileTeam1Territory = 1
	TileTeam2Territory = 2
	TileBoundary = 16
	TileSpecialTerrain = 32
	TileTransition1 = 17
	TileTransition2 = 18
	TileBridge1 = 257
	TileBridge2 = 258
)

// TileMap stores the tile data from CSV
type TileMap struct {
	Data [][]int
}

// NewGridSystem creates a grid with default cell types
func NewGridSystem() *GridSystem {
	// Create the grid
	grid := &GridSystem{
		CellWidth:  float64(screenWidth) / GridColumns,
		CellHeight: float64(screenHeight) / GridRows,
		ShowGrid:   false,
		CellTypes:  make([][]int, GridRows),
	}
	
	// Initialize all cells as ground
	for i := range grid.CellTypes {
		grid.CellTypes[i] = make([]int, GridColumns)
		// Default to ground
		for j := range grid.CellTypes[i] {
			grid.CellTypes[i][j] = CellTypeGround
		}
	}
	
	// Set up the Clash Royale style map with water and bridges
	// Water in the middle section
	waterStartRow := 30
	waterEndRow := 33
	
	// Bridge positions
	leftBridgeCol := 5
	rightBridgeCol := 27
	bridgeWidth := 4
	
	// Create water and bridges
	for row := waterStartRow; row <= waterEndRow; row++ {
		for col := 0; col < GridColumns; col++ {
			// Check if this cell is part of a bridge
			isLeftBridge := col >= leftBridgeCol && col < leftBridgeCol+bridgeWidth
			isRightBridge := col >= rightBridgeCol && col < rightBridgeCol+bridgeWidth
			
			if isLeftBridge || isRightBridge {
				grid.CellTypes[row][col] = CellTypeBridge
			} else {
				grid.CellTypes[row][col] = CellTypeWater
			}
		}
	}
	
	return grid
}

// GetCellType returns the type of cell at the given coordinates
func (g *GridSystem) GetCellType(col, row int) int {
	if row < 0 || row >= GridRows || col < 0 || col >= GridColumns {
		return CellTypeGround // Default for out-of-bounds
	}
	return g.CellTypes[row][col]
}

// SetCellType changes the type of a cell
func (g *GridSystem) SetCellType(col, row, cellType int) {
	if row >= 0 && row < GridRows && col >= 0 && col < GridColumns {
		g.CellTypes[row][col] = cellType
	}
}

// CellToPosition converts grid cell to position coordinates
func (g *GridSystem) CellToPosition(col, row int) Position {
	x := float64(col)*g.CellWidth + g.CellWidth/2
	y := float64(row)*g.CellHeight + g.CellHeight/2
	return Position{X: x, Y: y}
}

// PositionToCell converts screen position to grid cell coordinates
func (g *GridSystem) PositionToCell(pos Position) (int, int) {
	col := int(math.Floor(pos.X / g.CellWidth))
	row := int(math.Floor(pos.Y / g.CellHeight))

	return col, row
}

// Draw renders the grid and different cell types
func (g *GridSystem) Draw(screen *ebiten.Image) {
	// Draw tiles from CSV first if available
	if g.TileMap != nil {
		for row := 0; row < GridRows; row++ {
			for col := 0; col < GridColumns; col++ {
				tileType := g.TileMap.Data[row][col]
				
				// Calculate tile position
				x := float64(col) * g.CellWidth
				y := float64(row) * g.CellHeight
				
				// Choose color based on tile type
				var tileColor color.RGBA
				switch tileType {
				case TileEmpty:
					continue // Don't draw empty tiles
				case TileTeam1Territory:
					tileColor = color.RGBA{255, 200, 200, 100} // Light red
				case TileTeam2Territory:
					tileColor = color.RGBA{200, 200, 255, 100} // Light blue
				case TileBoundary:
					tileColor = color.RGBA{100, 100, 100, 180} // Dark gray
				case TileSpecialTerrain:
					tileColor = color.RGBA{150, 150, 150, 150} // Gray
				case TileTransition1:
					tileColor = color.RGBA{255, 220, 220, 100} // Lighter red
				case TileTransition2:
					tileColor = color.RGBA{220, 220, 255, 100} // Lighter blue
				case TileBridge1, TileBridge2:
					tileColor = color.RGBA{139, 69, 19, 255} // Brown for bridges
				default:
					continue // Skip unknown tile types
				}
				
				// Draw the tile
				ebitenutil.DrawRect(
					screen,
					x, y,
					g.CellWidth, g.CellHeight,
					tileColor,
				)
			}
		}
	}
	
	// Draw cell backgrounds based on type
	for row := 0; row < GridRows; row++ {
		for col := 0; col < GridColumns; col++ {
			cellType := g.GetCellType(col, row)
			
			// Only draw water and bridge cells with special colors
			var cellColor color.RGBA
			switch cellType {
			case CellTypeWater:
				cellColor = color.RGBA{0, 100, 200, 150} // Blue for water
				x := float64(col) * g.CellWidth
				y := float64(row) * g.CellHeight
				ebitenutil.DrawRect(screen, x, y, g.CellWidth, g.CellHeight, cellColor)
			
			case CellTypeBridge:
				cellColor = color.RGBA{139, 69, 19, 255} // Brown for bridge
				x := float64(col) * g.CellWidth
				y := float64(row) * g.CellHeight
				ebitenutil.DrawRect(screen, x, y, g.CellWidth, g.CellHeight, cellColor)
			}
		}
	}
	
	// Draw grid lines if enabled
	if !g.ShowGrid {
		return
	}
	
	// Grid line color
	gridColor := color.RGBA{100, 100, 255, 255}
	
	// Draw vertical lines (columns)
	for i := 0; i <= GridColumns; i++ {
		x := float64(i) * g.CellWidth
		ebitenutil.DrawLine(
			screen,
			x, 0,
			x, float64(screenHeight),
			gridColor,
		)
		
		// Draw column numbers at the top
		if i < GridColumns {
			ebitenutil.DebugPrintAt(
				screen,
				fmt.Sprintf("%d", i),
				int(x+2), 10,
			)
		}
	}
	
	// Draw horizontal lines (rows)
	for i := 0; i <= GridRows; i++ {
		y := float64(i) * g.CellHeight
		ebitenutil.DrawLine(
			screen,
			0, y,
			float64(screenWidth), y,
			gridColor,
		)
		
		// Draw row numbers on the left
		if i < GridRows {
			ebitenutil.DebugPrintAt(
				screen,
				fmt.Sprintf("%d", i),
				2, int(y+12),
			)
		}
	}
}

func (g *GridSystem) ToggleGrid() {
	g.ShowGrid = !g.ShowGrid
}

// LoadTileMap loads the tilemap from a CSV file
func (g *GridSystem) LoadTileMap(filepath string) error {
	// Initialize the map
	g.TileMap = &TileMap{
		Data: make([][]int, GridRows),
	}
	
	// Open and read the CSV file
	file, err := os.Open(filepath)
	if err != nil {
		return fmt.Errorf("failed to open tilemap: %v", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	
	// Skip header rows
	_, err = reader.Read() // Skip "Map" row
	if err != nil {
		return err
	}
	_, err = reader.Read() // Skip x coordinates
	if err != nil {
		return err
	}
	_, err = reader.Read() // Skip type row
	if err != nil {
		return err
	}

	// Initialize all rows
	for i := range g.TileMap.Data {
		g.TileMap.Data[i] = make([]int, GridColumns)
	}

	// Read each row
	for row := 0; row < GridRows; row++ {
		record, err := reader.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			return fmt.Errorf("error reading row %d: %v", row, err)
		}

		// Parse each column, starting from index 1 (skip first column)
		for col := 0; col < GridColumns; col++ {
			if col+1 >= len(record) || record[col+1] == "" {
				g.TileMap.Data[row][col] = TileEmpty
				continue
			}

			value, err := strconv.Atoi(record[col+1])
			if err != nil {
				fmt.Printf("Warning: Invalid tile value at row %d, col %d: %s\n", row, col, record[col+1])
				g.TileMap.Data[row][col] = TileEmpty
				continue
			}
			g.TileMap.Data[row][col] = value
		}
	}

	fmt.Printf("Successfully loaded tilemap with %d rows\n", len(g.TileMap.Data))
	return nil
}

// IsWalkableTile checks if a tile can be walked on
func (g *GridSystem) IsWalkableTile(col, row int) bool {
	if row < 0 || row >= GridRows || col < 0 || col >= GridColumns {
		return false
	}

	tileType := g.TileMap.Data[row][col]
	
	// Define walkable tiles
	switch tileType {
	case TileEmpty, TileTeam1Territory, TileTeam2Territory,
		 TileTransition1, TileTransition2, TileBridge1, TileBridge2:
		return true
	default:
		return false
	}
}

// GetTileType returns the type of tile at the given position
func (g *GridSystem) GetTileType(col, row int) int {
	if row < 0 || row >= GridRows || col < 0 || col >= GridColumns {
		return TileBoundary
	}
	return g.TileMap.Data[row][col]
}