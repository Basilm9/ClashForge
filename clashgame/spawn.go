// spawn.go
package clashgame

import (
	"image/color"

	"github.com/google/uuid"
)

// NewBuilding creates a new Building instance
func NewBuilding(x, y float64, health, damage int, attackRangeInCells float64, clr color.RGBA, widthCells, heightCells float64, grid *GridSystem) Building {
    if (int(widthCells) % 2) == 0{
        x += grid.CellWidth / 2
        y += grid.CellWidth / 2
    }
    return Building{
        Position:      Position{X: x, Y: y},
        Health:        health,
        MaxHealth:     health,
        Damage:        damage,
        Range:         attackRangeInCells, // Range in grid cells
        Color:         clr,
        WidthCells:    widthCells,         // Width in grid cells
        HeightCells:   heightCells,        // Height in grid cells
        Active:        true,
        ProjectileType: "normal",
        LastAttack:    0,
    }
}


// NewKingBuilding creates a new king Building instance
func NewKingBuilding(x, y float64, clr color.RGBA, grid *GridSystem) KingBuilding {
    
    return KingBuilding{
        Building: NewBuilding(x, y, 2000, 50, 10, clr, kingBuildingWidth, kingBuildingHeight, grid),
        ActivatesEndgame: true,
    }
}

func (b *Building) GetPixelDimensions(grid *GridSystem) (width, height float64) {
    width = b.WidthCells * grid.CellWidth
    height = b.HeightCells * grid.CellHeight
    return
}

func NewPlayer(color color.RGBA, isTopPlayer bool, playerIndex int, grid *GridSystem) Player {
    id, _ := uuid.NewRandom()
    
    // Initialize player with new attributes
    player := Player{
        Id:            id,
        Elixir:        4.0,
        Color:         color,
        NextCard:      0,
        ElixirMax:     10,
        ElixirGenRate: 0.1,
    }
    
    // Position Buildings based on player position using grid cells
    if isTopPlayer {
        // Top player
        kingPos := grid.CellToPosition(17, 5)  // Was 18,6
        
        leftBuildingPos := grid.CellToPosition(6, 12)   // Was 7,13
        rightBuildingPos := grid.CellToPosition(28, 12) // Was 29,13
        
        player.KingBuilding = NewKingBuilding(kingPos.X, kingPos.Y, color, grid)
        player.Buildings = []Building{
            NewBuilding(leftBuildingPos.X, leftBuildingPos.Y, 1200, 30, 3.5, color, princessWidth, princessHeight, grid),
            NewBuilding(rightBuildingPos.X, rightBuildingPos.Y, 1200, 30, 3.5, color, princessWidth, princessHeight, grid),
        }
    } else {
        // Bottom player
        kingPos := grid.CellToPosition(17, 57) // Was 18,58
        
        leftBuildingPos := grid.CellToPosition(6, 50)   // Was 7,51
        rightBuildingPos := grid.CellToPosition(28, 50) // Was 29,51
        
        player.KingBuilding = NewKingBuilding(kingPos.X, kingPos.Y, color, grid)
        player.Buildings = []Building{
            NewBuilding(leftBuildingPos.X, leftBuildingPos.Y, 1200, 30, 3.5, color, princessWidth, princessHeight, grid),
            NewBuilding(rightBuildingPos.X, rightBuildingPos.Y, 1200, 30, 3.5, color, princessWidth, princessHeight, grid),
        }
    }
    
    return player
}

func (g *Game) PlaceTroopAtCell(col, row int, health, damage int, speed, attackRange, aggroDistance float64, clr color.RGBA, team int) {
    pos := g.Grid.CellToPosition(col, row)
    
    // Mobs are typically smaller than one cell
    mobSizeInCells := 0.8 // In grid cells
    
    mob := NewTroop(
        pos.X, 
        pos.Y, 
        health,
        damage,
        speed,
        attackRange,
        aggroDistance,
        clr,
        g.Grid,         // Pass the grid system as a parameter
        mobSizeInCells,  // Pass the size in cells as a parameter
    )

    // Initialize movement fields
    mob.Position = pos
    mob.PrevPosition = pos
    mob.Velocity = Position{X: 0, Y: 0}
    mob.TargetVelocity = Position{X: 0, Y: 0}
    mob.TargetBuilding = nil
    mob.IsAttacking = false
    mob.Active = true
    mob.Team = team

    // Set initial target based on position
    if col < GridColumns/2 {
        // Troops on left side target left Building (index 1)
        mob.TargetIndex = 1
    } else {
        // Troops on right side target right Building (index 2)
        mob.TargetIndex = 2
    }

    g.Troops = append(g.Troops, mob)
}

var nextTroopID = 0
// Update the NewTroop function in spawn.go to initialize PrevPosition
func NewTroop(x, y float64, health, damage int, speedInCells, attackRangeInCells, aggroDistanceInCells float64, clr color.RGBA, grid *GridSystem, sizeInCells float64) Troop {
    // Calculate actual pixel size from grid cells
    gridCellSize := grid.CellWidth
    if grid.CellHeight < grid.CellWidth {
        gridCellSize = grid.CellHeight
    }
    actualSize := gridCellSize * sizeInCells
    
    // Default attack delay is 20 ticks
    attackDelay := 20

    pos := Position{X: x, Y: y}
    nextTroopID++
    return Troop{
        Position:      pos,
        PrevPosition:  pos,     // Initialize previous position to current position
        Health:        health,
        MaxHealth:     health,
        Damage:        damage,
        Speed:         speedInCells,       // Speed in grid cells per tick
        Range:         attackRangeInCells, // Range in grid cells
        AggroDistance: aggroDistanceInCells, // New aggro distance in grid cells
        Color:         clr,
        Size:          actualSize,         // Pixel size for rendering
        SizeInCells:   sizeInCells,        // Store the size in cells as the primary measurement
        Active:        true,
        AttackDelay:   attackDelay,
        IsNearTarget:  false,              // Initialize as not near target
        LastTargetChange: 0,               // Initialize last target change time
		ID: nextTroopID,
    }
}

// SpawnMob adds a new mob to the game
func SpawnTroop(mob Troop, team int, g *Game) {
    mob.Team = team  // Set the team explicitly
    g.Troops = append(g.Troops, mob)
}
