package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/basilm9/clash/clashgame"
	"github.com/hajimehoshi/ebiten/v2"
)

const (
	// Screen dimensions
	screenWidth = 1177 / 2
	screenHeight = 1687 / 2
	DURATION = 10
)

func main() {
    // Set window size and title
    ebiten.SetWindowSize(screenWidth, screenHeight)
    ebiten.SetWindowTitle("Tower Defense Game")

    // Get executable directory
    exePath, err := os.Executable()
    if err != nil {
        log.Fatalf("Failed to get executable path: %v", err)
    }
    exeDir := filepath.Dir(exePath)

    // Path to troops.csv, projectiles.csv, and tilemap.csv files
    troopsCsvPath := filepath.Join(exeDir, "clashgame/csv/troops.csv")
    projectilesCsvPath := filepath.Join(exeDir, "clashgame/csv/projectiles.csv")
    tilemapCsvPath := filepath.Join(exeDir, "clashgame/csv/tilemap.csv")

    // Initialize projectile system
    err = clashgame.LoadProjectileTemplates(projectilesCsvPath)
    if err != nil {
        clashgame.InitializeProjectileSystem()
    }

    // Initialize the troop system by loading templates from CSV
    err = clashgame.InitializeTroopSystem(troopsCsvPath)
    if err != nil {
        clashgame.InitializeWithDefaultTroops()
    }

    // Create the game
    game := clashgame.NewGame()

    // Set the global game instance
    clashgame.SetGameInstance(game)

    // Create troop selection UI
    game.TroopSelection = clashgame.NewTroopSelectionSystem()

    // Create enhanced troop drawer
    game.TroopDrawer = clashgame.NewEnhancedTroopDrawer(game)

    // Initialize game state
    game.GameTime = 0
    game.Running = true
    game.StopChannel = make(chan bool)

    // Start the game loop in a goroutine with the tilemap CSV path
    clashgame.StartGameLoop(game, tilemapCsvPath)

    // Run the game
    if err := ebiten.RunGame(game); err != nil {
        log.Fatal(err)
    }
}