// types.go
package clashgame

import (
	"image/color"
	"time"

	"github.com/google/uuid"
)

// troopState represents the current state/behavior of a troop
type troopState int

const (
	Seeking troopState = iota // Looking for targets
	Attacking               // Attacking a building
	Fighting                // Fighting another troop
)

// building grid positions and dimensions
const (
    // Screen dimensions
    screenWidth  = 1177 / 2
    screenHeight = 1687 / 2
    
    // Board ratios
    MAP_WIDTH_RATIO  = 0.7
    MAP_HEIGHT_RATIO = 1

    GridColumns = 36
    GridRows    = 64

    DURATION = 10

    // Building dimensions
    kingBuildingWidth = 6.0 
    kingBuildingHeight = 6.0 

    princessWidth = 4.0
    princessHeight = 5.0
)

type Vector struct {
    X, Y float64
}

type GridSystem struct {
	CellWidth  float64
	CellHeight float64
	ShowGrid   bool
	CellTypes  [][]int // Store the type of each cell
	TileMap    *TileMap // Add tilemap field
}

const (
    
    // New constants for movement
    PositionHistoryLength = 5      // Increased history length
    VelocityDamping = 0.95          // Damping factor to prevent oscillation
    DirectionSmoothingFactor = 0.2 // How quickly troop changes direction
    MinVelocityThreshold = 0.01    // Minimum velocity to prevent micro-jitters
)

// Add this struct to your Troop type (or modify your existing one)
type TroopPositionHistory struct {
	Positions []Position
	Index     int
}

type Troop struct {
    // Existing fields
    Name          string
    Position      Position
    PrevPosition  Position
    Health        int
    MaxHealth     int
    Damage        int
    Speed         float64
    Range         float64
    AggroDistance float64
    Color         color.RGBA
    SizeInCells   float64
    Size          float64
    Active        bool
    LastAttack    int
    Team          int
    AttackDelay   int
    TargetIndex   int
    IsNearTarget  bool
    LastTargetChange int
    ID            int
    // Movement smoothing fields
    PositionHistory TroopPositionHistory
    IsAttacking   bool
    // New fields for smooth movement
    Velocity      Position    // Current velocity vector
    TargetVelocity Position   // Desired velocity vector
    MaxAcceleration float64   // How quickly the troop can change direction
    TargetBuilding *Building // Current target building
}

type Game struct {
    Players            [2]Player
    Troops             []Troop
    Projectiles        []Projectile
    GameTime           int
    Running            bool
    Ticker             *time.Ticker
    StopChannel        chan bool
    ShowDebugGrid      bool
    Grid               *GridSystem
    LeftMousePressed   bool
    RightMousePressed  bool
    TroopSelection     *TroopSelectionSystem
    TroopDrawer        *EnhancedTroopDrawer
    ShowTroopInfo      bool
    SelectedTroopID    int
    ShowCSVPath        bool    // New field to control CSV path display
    CSVPath            string  // New field to store the CSV path
    
    // Building map for quick access (key: building ID, value: reference to building)
    BuildingMap        map[int]*Building
    NextBuildingID     int // To assign unique IDs to buildings
}

// Add decks to the Player struct
type Player struct {
    Id            uuid.UUID
    Elixir        float64
    KingBuilding     KingBuilding
    Buildings        []Building
    Color         color.RGBA
    NextCard      int       // Index of next card to draw
    ElixirMax     int       // Maximum elixir capacity
    ElixirGenRate float64   // Elixir generated per second
}

// Position represents a 2D position
type Position struct {
	X, Y float64
}

// building represents a basic defensive building
type Building struct {
    Position      Position
    Health        int
    MaxHealth     int
    Damage        int
    Range         float64       // Range in grid cells
    Color         color.RGBA
    WidthCells    float64       // Width in grid cells
    HeightCells   float64       // Height in grid cells
    Active        bool
    ProjectileType string
    LastAttack    int
    ID            int           // Unique identifier for the building
    Team          int           // Team ID (0 or 1)
}

// Kingbuilding represents the main building for each player
type KingBuilding struct {
	Building             // Embedding building struct
	ActivatesEndgame  bool
}