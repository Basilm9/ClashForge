// Modified gameloop.go
package clashgame

import (
	"fmt"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

func StartGameLoop(game *Game, tilemapPath string) {
    fmt.Println("starting game loop...")
    
    game.Running = true
    game.Ticker = time.NewTicker(time.Millisecond * 40)
    gameTimer := time.NewTimer(time.Duration(DURATION) * time.Minute)
    broadcastStateTicker := time.NewTicker(time.Millisecond * 33)
    elixirTicker := time.NewTicker(time.Second)
    
    go func() {
        defer game.Ticker.Stop()
        defer gameTimer.Stop()
        defer broadcastStateTicker.Stop()
        defer elixirTicker.Stop()
        
        for {
            if !game.IsActive() {
                break
            }
            
            select {
            case <-game.Ticker.C:
                game.GameTime++
                
                // Process these updates in an improved order:
                // 1. Update projectiles first to ensure they hit targets before they move
                UpdateProjectiles(game)
                
                // 2. Now update troops with the old projectiles cleared
                UpdateTroopMovement(game)
                
                // 3. Clear any invalid attack states
                ClearInvalidAttackStates(game)
                
            case <-broadcastStateTicker.C:
                // broadcastStateToClients(game)
                // broadcastStateToSpectators(game)
            case <-elixirTicker.C:
                UpdateElixir(game)
            case <-gameTimer.C:
                fmt.Println("Game over: Time's up!")
                game.Running = false
            case <-game.StopChannel:
                return
            }
        }
    }()
}

// DrawGame draws the game state
func DrawGame(screen *ebiten.Image, game *Game) {
    // Draw the grid first
    game.Grid.Draw(screen)
    
    // Draw buildings
    for _, building := range game.BuildingMap {
        if building.Active {
            // Draw building logic here
            width, height := building.GetPixelDimensions(game.Grid)
            ebitenutil.DrawRect(screen, 
                building.Position.X - width/2,
                building.Position.Y - height/2,
                width, height,
                building.Color)
        }
    }
    
    // Draw troops
    for _, troop := range game.Troops {
        if troop.Active {
            game.TroopDrawer.DrawTroop(screen, &troop)
        }
    }
    
    // Draw projectiles
    for _, projectile := range game.Projectiles {
        if projectile.Active {
            // Draw projectile logic here
            ebitenutil.DrawCircle(screen, projectile.Position.X, projectile.Position.Y, projectile.Size/2, projectile.Color)
        }
    }
    
    // Draw CSV path if enabled
    if game.ShowCSVPath {
        ebitenutil.DebugPrintAt(screen, fmt.Sprintf("CSV Path: %s", game.CSVPath), 10, 30)
    }
}