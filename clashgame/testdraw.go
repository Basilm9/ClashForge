// testdraw.go
package clashgame

import (
	"fmt"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// Draw draws the Building on the screen with enhanced visuals
func (b *Building) Draw(screen *ebiten.Image, grid *GridSystem) {
    if !b.Active {
        return
    }
    
    // Get pixel dimensions
    width, height := b.GetPixelDimensions(grid)
    
    // Draw building base (main body)
    baseColor := b.Color
    ebitenutil.DrawRect(
        screen, 
        b.Position.X - width/2,  
        b.Position.Y - height/2,
        width, 
        height, 
        baseColor,
    )
    
    // Draw building details based on size
    isKingTower := width >= 3.5*grid.CellWidth
    
    // Create darker color for details
    darkerColor := color.RGBA{
        R: uint8(float64(baseColor.R) * 0.7),
        G: uint8(float64(baseColor.G) * 0.7),
        B: uint8(float64(baseColor.B) * 0.7),
        A: baseColor.A,
    }
    
    // Create lighter color for windows/highlights
    lighterColor := color.RGBA{
        R: uint8(math.Min(float64(baseColor.R)+40, 255)),
        G: uint8(math.Min(float64(baseColor.G)+40, 255)),
        B: uint8(math.Min(float64(baseColor.B)+40, 255)),
        A: baseColor.A,
    }
    
    // Draw building top (towers, etc.)
    if isKingTower {
        // King tower - draw a crown-like structure
        towerWidth := width * 0.2
        centerX := b.Position.X
        
        // Draw central taller tower
        ebitenutil.DrawRect(
            screen,
            centerX - towerWidth/2,
            b.Position.Y - height/2 - towerWidth*1.2,
            towerWidth,
            towerWidth*1.2,
            darkerColor,
        )
        
        // Draw side towers
        sideOffset := width * 0.3
        for _, x := range []float64{centerX - sideOffset, centerX + sideOffset} {
            ebitenutil.DrawRect(
                screen,
                x - towerWidth*0.4,
                b.Position.Y - height/2 - towerWidth*0.8,
                towerWidth*0.8,
                towerWidth*0.8,
                darkerColor,
            )
        }
    } else {
        // Regular tower - just a simpler top
        towerWidth := width * 0.3
        ebitenutil.DrawRect(
            screen,
            b.Position.X - towerWidth/2,
            b.Position.Y - height/2 - towerWidth*0.8,
            towerWidth,
            towerWidth*0.8,
            darkerColor,
        )
    }
    
    // Draw windows
    windowSize := width * 0.15
    if isKingTower {
        // More windows for king tower
        for i := -1; i <= 1; i++ {
            for j := -1; j <= 1; j++ {
                if i == 0 && j == 0 {
                    continue // Skip center
                }
                
                windowX := b.Position.X + float64(i)*width*0.25
                windowY := b.Position.Y + float64(j)*height*0.25
                
                ebitenutil.DrawRect(
                    screen,
                    windowX - windowSize/2,
                    windowY - windowSize/2,
                    windowSize,
                    windowSize,
                    lighterColor,
                )
            }
        }
    } else {
        // Just two windows for regular towers
        ebitenutil.DrawRect(
            screen,
            b.Position.X - width*0.2,
            b.Position.Y,
            windowSize,
            windowSize,
            lighterColor,
        )
        
        ebitenutil.DrawRect(
            screen,
            b.Position.X + width*0.2 - windowSize,
            b.Position.Y,
            windowSize,
            windowSize,
            lighterColor,
        )
    }
    
    // Draw building outline
    borderWidth := 1.0
    borderColor := darkerColor
    
    // Top border
    ebitenutil.DrawRect(
        screen,
        b.Position.X - width/2,
        b.Position.Y - height/2,
        width,
        borderWidth,
        borderColor,
    )
    
    // Bottom border
    ebitenutil.DrawRect(
        screen,
        b.Position.X - width/2,
        b.Position.Y + height/2 - borderWidth,
        width,
        borderWidth,
        borderColor,
    )
    
    // Left border
    ebitenutil.DrawRect(
        screen,
        b.Position.X - width/2,
        b.Position.Y - height/2,
        borderWidth,
        height,
        borderColor,
    )
    
    // Right border
    ebitenutil.DrawRect(
        screen,
        b.Position.X + width/2 - borderWidth,
        b.Position.Y - height/2,
        borderWidth,
        height,
        borderColor,
    )
    
    // Draw health bar
    healthBarWidth := width
    healthBarHeight := 5.0
    healthPercent := float64(b.Health) / float64(b.MaxHealth)
    
    // Health bar background (gray)
    ebitenutil.DrawRect(
        screen,
        b.Position.X - width/2,
        b.Position.Y - height/2 - 10,
        healthBarWidth,
        healthBarHeight,
        color.RGBA{40, 40, 40, 220},
    )
    
    // Health bar fill (green to red based on health)
    healthColor := color.RGBA{
        uint8(255 * (1 - healthPercent)),
        uint8(255 * healthPercent),
        0,
        255,
    }
    
    ebitenutil.DrawRect(
        screen,
        b.Position.X - width/2,
        b.Position.Y - height/2 - 10,
        healthBarWidth * healthPercent,
        healthBarHeight,
        healthColor,
    )
}

// Helper function to draw building details
func drawBuildingDetails(screen *ebiten.Image, b *Building, width, height float64, grid *GridSystem) {
    // Draw a castle-like structure on top for visual interest
    
    // Is this a king tower?
    isKing := width >= 4*grid.CellWidth
    
    // Color variations
    baseColor := b.Color
    darkerColor := color.RGBA{
        R: uint8(float64(baseColor.R) * 0.7),
        G: uint8(float64(baseColor.G) * 0.7),
        B: uint8(float64(baseColor.B) * 0.7),
        A: baseColor.A,
    }
    lighterColor := color.RGBA{
        R: uint8(math.Min(float64(baseColor.R) * 1.3, 255)),
        G: uint8(math.Min(float64(baseColor.G) * 1.3, 255)),
        B: uint8(math.Min(float64(baseColor.B) * 1.3, 255)),
        A: baseColor.A,
    }
    
    if isKing {
        // Draw a crown-like structure for king tower
        centerX := b.Position.X
        baseY := b.Position.Y - height/2
        towerWidth := width * 0.15
        towerSpacing := width * 0.25
        
        // Draw central tower (tallest)
        ebitenutil.DrawRect(
            screen,
            centerX - towerWidth/2,
            baseY - towerWidth*1.5,
            towerWidth,
            towerWidth*1.5,
            darkerColor,
        )
        
        // Draw left tower
        ebitenutil.DrawRect(
            screen,
            centerX - towerSpacing - towerWidth/2,
            baseY - towerWidth,
            towerWidth,
            towerWidth,
            darkerColor,
        )
        
        // Draw right tower
        ebitenutil.DrawRect(
            screen,
            centerX + towerSpacing - towerWidth/2,
            baseY - towerWidth,
            towerWidth,
            towerWidth,
            darkerColor,
        )
        
        // Draw windows on the main body
        windowSize := width * 0.1
        windowSpacing := width * 0.2
        windowY := b.Position.Y
        
        // Draw a row of windows
        for i := -1; i <= 1; i++ {
            ebitenutil.DrawRect(
                screen,
                centerX + float64(i)*windowSpacing - windowSize/2,
                windowY - windowSize/2,
                windowSize,
                windowSize,
                lighterColor,
            )
        }
    } else {
        // Regular tower - simpler design
        centerX := b.Position.X
        baseY := b.Position.Y - height/2
        towerWidth := width * 0.3
        
        // Draw a single tower
        ebitenutil.DrawRect(
            screen,
            centerX - towerWidth/2,
            baseY - towerWidth,
            towerWidth,
            towerWidth,
            darkerColor,
        )
        
        // Draw a window
        windowSize := width * 0.15
        ebitenutil.DrawRect(
            screen,
            centerX - windowSize/2,
            b.Position.Y - windowSize/2,
            windowSize,
            windowSize,
            lighterColor,
        )
    }
    
    // Draw a border around the building
    borderWidth := 1.0
    ebitenutil.DrawRect(
        screen,
        b.Position.X - width/2,
        b.Position.Y - height/2,
        width,
        borderWidth,
        darkerColor,
    )
    ebitenutil.DrawRect(
        screen,
        b.Position.X - width/2,
        b.Position.Y + height/2 - borderWidth,
        width,
        borderWidth,
        darkerColor,
    )
    ebitenutil.DrawRect(
        screen,
        b.Position.X - width/2,
        b.Position.Y - height/2,
        borderWidth,
        height,
        darkerColor,
    )
    ebitenutil.DrawRect(
        screen,
        b.Position.X + width/2 - borderWidth,
        b.Position.Y - height/2,
        borderWidth,
        height,
        darkerColor,
    )
}

// Draw draws the troop on the screen
func (m *Troop) Draw(screen *ebiten.Image) {
	if m.Active {
		// Base troop color
		troopColor := m.Color
		
		// Determine if the troop has recently attacked (is in combat)
		inCombat := false
		
		// If troop has attacked within the last 5 game ticks, highlight it
		if m.LastAttack > 0 && (getGame().GameTime - m.LastAttack) < 5 {
			inCombat = true
		}
		
		// Draw a combat indicator (red glow) if in combat
		if inCombat {
			// Draw a slightly larger red circle behind the troop
			ebitenutil.DrawCircle(
				screen,
				m.Position.X,
				m.Position.Y,
				m.Size/2 + 2,
				color.RGBA{255, 0, 0, 150},
			)
		}
		
		// Draw the troop
		ebitenutil.DrawCircle(
			screen,
			m.Position.X,
			m.Position.Y,
			m.Size/2,
			troopColor,
		)
		
		// Draw health bar
		healthBarWidth := m.Size
		healthBarHeight := 4.0
		healthPercent := float64(m.Health) / float64(m.MaxHealth)
		
		// Background (gray)
		ebitenutil.DrawRect(
			screen,
			m.Position.X-m.Size/2,
			m.Position.Y-m.Size/2-8,
			healthBarWidth,
			healthBarHeight,
			color.RGBA{100, 100, 100, 255},
		)
		
		// Health (green to red based on health)
		healthColor := color.RGBA{
			uint8(255 * (1 - healthPercent)),
			uint8(255 * healthPercent),
			0,
			255,
		}
		ebitenutil.DrawRect(
			screen,
			m.Position.X-m.Size/2,
			m.Position.Y-m.Size/2-8,
			healthBarWidth*healthPercent,
			healthBarHeight,
			healthColor,
		)
	}
}

// Helper function to get the game instance
// We need to add this as a global variable or use a singleton pattern
var gameInstance *Game

func SetGameInstance(game *Game) {
	gameInstance = game
}

func getGame() *Game {
	return gameInstance
}

// Add this to the end of projectiles.go or in a helper file
func DrawCircle(screen *ebiten.Image, x, y, radius float64, clr color.RGBA) {
    ebitenutil.DrawCircle(screen, x, y, radius, clr)
}

// Layout is required for ebiten.Game interface
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func (g *Game) SpawnSkeletonArmy(centerCol, centerRow, count int, team int) {
    // Skeleton stats - smaller, faster, weaker troops
    health := 50           // Less health than regular troops
    damage := 20           // Lower damage
    speed := .2          // Faster movement
    attackRange := .1     // Shorter attack range
    aggroDistance := 4.0   // Standard aggro distance
    
    // Get color based on team
    var troopColor color.RGBA
    if team == 0 {
        // Team 0 skeletons (red with white tint)
        troopColor = color.RGBA{255, 120, 120, 255}
    } else {
        // Team 1 skeletons (blue with white tint)
        troopColor = color.RGBA{120, 120, 255, 255}
    }
    
    // Spawn pattern radius - how far from center to place skeletons
    spawnRadius := count / 3
    if spawnRadius < 1 {
        spawnRadius = 1
    }
    
    // Spawn skeletons in a loose group around the center point
    for i := 0; i < count; i++ {
        // Calculate spread pattern
        // Use a spiral or circular pattern to place troops
        angle := float64(i) * (2.0 * math.Pi / float64(count))
        distance := float64(i % spawnRadius) + 0.5
        
        // Calculate offset from center
        offsetX := int(math.Round(math.Cos(angle) * distance))
        offsetY := int(math.Round(math.Sin(angle) * distance))
        
        // Calculate final position, ensuring it's within grid bounds
        spawnCol := centerCol + offsetX
        spawnRow := centerRow + offsetY
        
        // Ensure spawn position is within grid bounds
        if spawnCol < 0 {
            spawnCol = 0
        }
        if spawnCol >= GridColumns {
            spawnCol = GridColumns - 1
        }
        if spawnRow < 0 {
            spawnRow = 0
        }
        if spawnRow >= GridRows {
            spawnRow = GridRows - 1
        }
        
        // Add a tiny random offset to health to help with collision resolution
        randomizedHealth := health + (i % 10)
        
        // Spawn the troop
        g.PlaceTroopAtCell(
            spawnCol, 
            spawnRow,
            randomizedHealth,  // Health with small random variation
            damage,            // Damage
            speed,             // Speed (in grid cells per tick)
            attackRange,       // Range (in grid cells)
            aggroDistance,     // Aggro distance (in grid cells)
            troopColor,        // Team-specific skeleton color
            team,              // Team
        )
    }
}

func (g *Game) Update() error {
    if !g.IsActive() {
        return nil
    }

    // Toggle grid visibility with G key
    if inpututil.IsKeyJustPressed(ebiten.KeyG) {
        g.Grid.ToggleGrid()
    }
    
    // Toggle troop info display with T key
    if inpututil.IsKeyJustPressed(ebiten.KeyT) {
        g.ShowTroopInfo = !g.ShowTroopInfo
    }
    
    // Update troop selection system if initialized
    if g.TroopSelection != nil {
        g.TroopSelection.Update()
    }
    
    // Handle left mouse clicks to spawn friendly troops
    if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
        // Only spawn on initial press (not continuous holding)
        if !g.LeftMousePressed {
            g.LeftMousePressed = true
            
            // Get mouse position and convert to grid cell
            x, y := ebiten.CursorPosition()
            mousePos := Position{X: float64(x), Y: float64(y)}
            col, row := g.Grid.PositionToCell(mousePos)
            
            // Only deploy if click is outside the UI area
            if g.TroopSelection != nil && float64(y) >= float64(g.TroopSelection.UIY) {
                // Click is in the UI area - don't spawn troops
                
                // Check if we clicked on a troop card
                if g.TroopSelection != nil {
                    for i := 0; i < min(g.TroopSelection.MaxVisibleCards, len(g.TroopSelection.TroopNames)-g.TroopSelection.ScrollIndex); i++ {
                        cardX := g.TroopSelection.UIX + i*(g.TroopSelection.CardWidth+g.TroopSelection.CardSpacing) + g.TroopSelection.CardSpacing
                        if x >= cardX && x <= cardX+g.TroopSelection.CardWidth {
                            troopIndex := g.TroopSelection.ScrollIndex + i
                            if troopIndex < len(g.TroopSelection.TroopNames) {
                                g.TroopSelection.SelectedTroop = g.TroopSelection.TroopNames[troopIndex]
                            }
                            break
                        }
                    }
                }
            } else {
                // Click is in the game area - spawn troop
                
                // Check for Skeleton Army modifier key (Shift)
                if ebiten.IsKeyPressed(ebiten.KeyShift) {
                    // Spawn skeleton army (12 troops) for team 0
                    g.SpawnSkeletonArmy(col, row, 12, 0)
                } else if g.TroopSelection != nil && g.TroopSelection.SelectedTroop != "" {
                    // Use selected troop from the troop selection system
                    g.TroopSelection.DeploySelectedTroop(g, col, row, 0)
                } else {
                    // Default behavior - spawn a single regular troop
                    g.PlaceTroopAtCell(
                        col, row,
                        100,        // Health
                        40,         // Damage
                        0.2,        // Speed (in grid cells per tick)
                        1.5,        // Range (in grid cells)
                        5.0,        // Aggro distance (in grid cells)
                        g.Players[0].Color, // Same color as player
                        0,          // Team 0 (friendly)
                    )
                }
                
                // Also check for troop selection (for detailed info)
                nearestTroopID := -1
                nearestDistance := 20.0 // Maximum selection distance
                
                for i, troop := range g.Troops {
                    if troop.Active {
                        dist := Distance(mousePos, troop.Position)
                        if dist < troop.Size/2 && dist < nearestDistance {
                            nearestDistance = dist
                            nearestTroopID = i
                        }
                    }
                }
                
                g.SelectedTroopID = nearestTroopID
            }
        }
    } else {
        g.LeftMousePressed = false
    }
    
    // Handle right mouse clicks to spawn enemy troops
    if ebiten.IsMouseButtonPressed(ebiten.MouseButtonRight) {
        // Only spawn on initial press
        if !g.RightMousePressed {
            g.RightMousePressed = true
            
            // Get mouse position and convert to grid cell
            x, y := ebiten.CursorPosition()
            mousePos := Position{X: float64(x), Y: float64(y)}
            col, row := g.Grid.PositionToCell(mousePos)
            
            // Only deploy if click is outside the UI area
            if g.TroopSelection != nil && float64(y) >= float64(g.TroopSelection.UIY) {
                // Click is in the UI area - don't spawn troops
            } else {
                // Click is in the game area - spawn troop
                
                // Check for Skeleton Army modifier key (Shift)
                if ebiten.IsKeyPressed(ebiten.KeyShift) {
                    // Spawn skeleton army (12 troops) for team 1
                    g.SpawnSkeletonArmy(col, row, 12, 1)
                } else if g.TroopSelection != nil && g.TroopSelection.SelectedTroop != "" {
                    // Use selected troop from the troop selection system
                    g.TroopSelection.DeploySelectedTroop(g, col, row, 1)
                } else {
                    // Default behavior - spawn a single regular enemy troop
                    g.PlaceTroopAtCell(
                        col, row,
                        100,        // Health
                        40,         // Damage
                        0.2,        // Speed (in grid cells per tick)
                        1.5,        // Range (in grid cells)
                        5.0,        // Aggro distance (in grid cells)
                        g.Players[1].Color, // Same color as enemy player
                        1,          // Team 1 (enemy)
                    )
                }
            }
        }
    } else {
        g.RightMousePressed = false
    }
    
    return nil
}

// Draw draws the projectile on the screen
func (p *Projectile) Draw(screen *ebiten.Image, game Game) {
	if p.Active {
		// Draw the projectile as a circle
		ebitenutil.DrawCircle(
			screen,
			p.Position.X,
			p.Position.Y,
			p.Size/2,
			p.Color,
		)
		
		// Draw trail effect if specified
		if p.Template.TrailEffect != "" {
			// Calculate trail points based on direction
			trailLength := 3
			
			for i := 1; i <= trailLength; i++ {
				// Calculate position
				trailX := p.Position.X - p.Direction.X * p.Speed * game.Grid.CellWidth * float64(i) * 0.5
				trailY := p.Position.Y - p.Direction.Y * p.Speed * game.Grid.CellHeight * float64(i) * 0.5
				
				// Calculate opacity and size
				opacity := 200 - i*40
				if opacity < 0 {
					opacity = 0
				}
				
				size := (p.Size/2) * (1.0 - float64(i)*0.2)
				
				// Draw trail segment
				trailColor := p.Color
				trailColor.A = uint8(opacity)
				
				ebitenutil.DrawCircle(
					screen,
					trailX,
					trailY,
					size,
					trailColor,
				)
			}
		}
	}
}

func (g *Game) Draw(screen *ebiten.Image) {
    // Draw background
    screen.Fill(color.RGBA{200, 200, 200, 255})
    
    // Draw grid
    g.Grid.Draw(screen)
    
    // IMPORTANT: Draw buildings
    for _, player := range g.Players {
        // Draw king building
        player.KingBuilding.Draw(screen, g.Grid)
        
        // Draw regular buildings
        for i := range player.Buildings {
            player.Buildings[i].Draw(screen, g.Grid)
        }
    }
    
    // Draw projectiles
    for _, projectile := range g.Projectiles {
        projectile.Draw(screen, *g)
    }
    
    // Draw troops using enhanced visuals if available
    if g.TroopDrawer != nil {
        for i, troop := range g.Troops {
            if troop.Active {
                g.TroopDrawer.DrawTroop(screen, &troop)
                
                // Draw selection highlight if this troop is selected
                if i == g.SelectedTroopID {
                    ebitenutil.DrawCircle(
                        screen,
                        troop.Position.X,
                        troop.Position.Y,
                        troop.Size/2+5,
                        color.RGBA{255, 255, 0, 150},
                    )
                }
            }
        }
    } else {
        // Fallback to original troop drawing
        for _, troop := range g.Troops {
            troop.Draw(screen)
        }
    }
    
    // Draw troop selection UI if available
    if g.TroopSelection != nil {
        g.TroopSelection.Draw(screen)
    }
    
    // Draw rest of UI (selected troop info, etc.)
    if g.ShowTroopInfo && g.SelectedTroopID >= 0 && g.SelectedTroopID < len(g.Troops) {
        // (existing code for troop info display)
    }
}

// DebugCombatSystem prints information about troops and projectiles
func DebugCombatSystem(game *Game) {
    // Only run debug every 60 ticks to avoid spamming console
    if game.GameTime % 60 != 0 {
        return
    }
    
    fmt.Println("\n----- COMBAT SYSTEM DEBUG -----")
    fmt.Printf("Game Time: %d\n", game.GameTime)
    
    // Count active troops by type and team
    meleeTroops := [2]int{0, 0} // [team0, team1]
    rangedTroops := [2]int{0, 0}
    flyingTroops := [2]int{0, 0}
    
    // Print info for up to 5 troops
    fmt.Println("Active Troops (sample):")
    sampleCount := 0
    
    for _, troop := range game.Troops {
        if !troop.Active {
            continue
        }
        
        template := GetTroopTemplate(&troop)
        isMelee := troop.Range <= 1.0
        isFlying := IsFlyingTroop(&troop)
        
        // Update counters
        if isFlying {
            flyingTroops[troop.Team]++
        } else if isMelee {
            meleeTroops[troop.Team]++
        } else {
            rangedTroops[troop.Team]++
        }
        
        // Print sample troop info
        if sampleCount < 5 {
            hasProjectile := "No"
            if template != nil && template.Projectile.Name != "" {
                hasProjectile = fmt.Sprintf("Yes (%s)", template.Projectile.Name)
            }
            
            troopType := "Ranged"
            if isMelee {
                troopType = "Melee"
            }
            if isFlying {
                troopType = "Flying"
            }
            
            fmt.Printf("  ID=%d, Name=%s, Team=%d, Type=%s, Range=%.1f, Health=%d/%d, HasProjectile=%s\n",
                      troop.ID, troop.Name, troop.Team, troopType, troop.Range, 
                      troop.Health, troop.MaxHealth, hasProjectile)
            
            sampleCount++
        }
    }
    
    // Print troop counts
    fmt.Printf("Team 0 Troops: %d Melee, %d Ranged, %d Flying\n", 
               meleeTroops[0], rangedTroops[0], flyingTroops[0])
    fmt.Printf("Team 1 Troops: %d Melee, %d Ranged, %d Flying\n", 
               meleeTroops[1], rangedTroops[1], flyingTroops[1])
    
    // Print projectile info
    activeProjectiles := 0
    for _, p := range game.Projectiles {
        if p.Active {
            activeProjectiles++
        }
    }
    
    fmt.Printf("Active Projectiles: %d\n", activeProjectiles)
    
    // Print sample projectiles
    if len(game.Projectiles) > 0 {
        fmt.Println("Projectile Samples:")
        
        sampleCount = 0
        for i, p := range game.Projectiles {
            if !p.Active || sampleCount >= 3 {
                continue
            }
            
            fmt.Printf("  ID=%d, Name=%s, Team=%d, Damage=%d, Speed=%.2f, Size=%.1f, Pos=(%.1f,%.1f)\n",
                      i, p.Name, p.Team, p.Damage, p.Speed, p.Size, p.Position.X, p.Position.Y)
            
            sampleCount++
        }
    }
    
    fmt.Println("-------------------------------")
}

// Add this helper function to check template integrity
func CheckTroopTemplateIntegrity() {
    fmt.Println("\n----- CHECKING TROOP TEMPLATES -----")
    
    // Count how many templates we have
    fmt.Printf("Total troop templates: %d\n", len(TroopTemplateMap))
    
    meleeTroops := 0
    rangedTroops := 0
    flyingTroops := 0
    
    // Check each template
    for name, template := range TroopTemplateMap {
        isMelee := template.Range <= 1.0
        isFlying := template.FlyingHeight > 0
        
        // Categorize the troop
        if isFlying {
            flyingTroops++
        } else if isMelee {
            meleeTroops++
        } else {
            rangedTroops++
        }
        
        // Check if projectile assignment is correct
        hasProjectile := template.Projectile.Name != "" && template.Projectile.Name != "none"
        
        if isMelee && hasProjectile {
            fmt.Printf("WARNING: Melee troop '%s' has projectile '%s' assigned\n", 
                       name, template.Projectile.Name)
        } else if !isMelee && !hasProjectile {
            fmt.Printf("WARNING: Ranged troop '%s' has no projectile assigned\n", name)
        }
    }
    
    fmt.Printf("Template counts: %d Melee, %d Ranged, %d Flying\n", 
               meleeTroops, rangedTroops, flyingTroops)
               
    fmt.Println("----------------------------------")
}

// Call this from your game.Update() method to enable debugging
func (g *Game) EnableCombatDebugging() {
    // Check if D key is pressed
    if inpututil.IsKeyJustPressed(ebiten.KeyD) {
        CheckTroopTemplateIntegrity()
    }
    
    // Run combat debug every 60 frames
    if g.GameTime % 60 == 0 {
        DebugCombatSystem(g)
    }
}