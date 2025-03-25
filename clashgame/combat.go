package clashgame

import (
	"fmt"
	"math"
)

func ProcessCombat(game *Game, troop1, troop2 *Troop) {
    // Recompute attack capabilities inside the function
    canAttack1to2 := CanAttackTroop(troop1, troop2, game.Grid)
    canAttack2to1 := CanAttackTroop(troop2, troop1, game.Grid)

    // Set attacking state based on individual ability to attack
    troop1.IsAttacking = canAttack1to2
    troop2.IsAttacking = canAttack2to1
    
    // Get current game time
    currentTime := game.GameTime
    
    // Process attack from troop1 to troop2 ONLY if troop1 can attack troop2
    if canAttack1to2 && currentTime - troop1.LastAttack >= troop1.AttackDelay {
        // Reset attack timer
        troop1.LastAttack = currentTime
        
        // Check if troop1 is melee - if so, never use projectiles
        if IsMeleeTroop(troop1) {
            // Melee troops apply damage directly
            troop2.Health -= troop1.Damage
            fmt.Printf("Melee attack! Troop ID=%d deals %d damage to Troop ID=%d (health now: %d)\n", 
                       troop1.ID, troop1.Damage, troop2.ID, troop2.Health)
            
            // Check if troop2 is defeated
            if troop2.Health <= 0 {
                troop2.Active = false
                fmt.Printf("Troop ID=%d defeated by Troop ID=%d\n", troop2.ID, troop1.ID)
                // Clear attacking state of troop1 when target is defeated
                troop1.IsAttacking = false
            }
        } else {
            // Get template for projectile information
            template1 := GetTroopTemplate(troop1)
            
            // Check if troop1 has a projectile defined
            hasProjectile := template1 != nil && template1.Projectile.Name != ""
            
            if hasProjectile {
                // Create the projectile
                projectile := CreateProjectile(
                    template1.Projectile.Name,
                    troop1.Position,
                    troop2.Position,
                    troop1.Damage,  // USE TROOP'S DAMAGE
                    troop1.Team,
                    troop1.ID,
                )
                
                fmt.Printf("Troop ID=%d fires projectile at Troop ID=%d\n", troop1.ID, troop2.ID)
                
                // Set target entity for homing projectiles
                if projectile != nil && projectile.IsHoming {
                    projectile.TargetEntity = troop2
                }
                
                // Add projectile to game if it was created successfully
                if projectile != nil {
                    game.Projectiles = append(game.Projectiles, *projectile)
                } else {
                    // Fallback to direct damage if projectile creation failed
                    troop2.Health -= troop1.Damage
                    fmt.Printf("Direct fallback! Troop ID=%d deals %d damage to Troop ID=%d (health now: %d)\n", 
                               troop1.ID, troop1.Damage, troop2.ID, troop2.Health)
                }
            } else {
                // Ranged troops without projectiles defined fall back to direct damage
                troop2.Health -= troop1.Damage
                fmt.Printf("Ranged attack! Troop ID=%d deals %d damage to Troop ID=%d (health now: %d)\n", 
                           troop1.ID, troop1.Damage, troop2.ID, troop2.Health)
                
                // Check if troop2 is defeated
                if troop2.Health <= 0 {
                    troop2.Active = false
                    fmt.Printf("Troop ID=%d defeated by Troop ID=%d\n", troop2.ID, troop1.ID)
                    // Clear attacking state of troop1 when target is defeated
                    troop1.IsAttacking = false
                }
            }
        }
    }
    
    // Process attack from troop2 to troop1 ONLY if troop2 can attack troop1
    if canAttack2to1 && currentTime - troop2.LastAttack >= troop2.AttackDelay {
        // Reset attack timer
        troop2.LastAttack = currentTime
        
        // Check if troop2 is melee - if so, never use projectiles
        if IsMeleeTroop(troop2) {
            // Melee troops apply damage directly
            troop1.Health -= troop2.Damage
            fmt.Printf("Melee attack! Troop ID=%d deals %d damage to Troop ID=%d (health now: %d)\n", 
                       troop2.ID, troop2.Damage, troop1.ID, troop1.Health)
            
            // Check if troop1 is defeated
            if troop1.Health <= 0 {
                troop1.Active = false
                fmt.Printf("Troop ID=%d defeated by Troop ID=%d\n", troop1.ID, troop2.ID)
                // Clear attacking state of troop2 when target is defeated
                troop2.IsAttacking = false
            }
        } else {
            // Get template for projectile information
            template2 := GetTroopTemplate(troop2)
            
            // Check if troop2 has a projectile defined
            hasProjectile := template2 != nil && template2.Projectile.Name != ""
            
            if hasProjectile {
                // Create the projectile
                projectile := CreateProjectile(
                    template2.Projectile.Name,
                    troop2.Position,
                    troop1.Position,
                    troop2.Damage,  // USE TROOP'S DAMAGE
                    troop2.Team,
                    troop2.ID,
                )
                
                fmt.Printf("Troop ID=%d fires projectile at Troop ID=%d\n", troop2.ID, troop1.ID)
                
                // Set target entity for homing projectiles
                if projectile != nil && projectile.IsHoming {
                    projectile.TargetEntity = troop1
                }
                
                // Add projectile to game if created successfully
                if projectile != nil {
                    game.Projectiles = append(game.Projectiles, *projectile)
                } else {
                    // Fallback to direct damage if projectile creation failed
                    troop1.Health -= troop2.Damage
                    fmt.Printf("Direct fallback! Troop ID=%d deals %d damage to Troop ID=%d (health now: %d)\n", 
                               troop2.ID, troop2.Damage, troop1.ID, troop1.Health)
                }
            } else {
                // Ranged troops without projectiles defined fall back to direct damage
                troop1.Health -= troop2.Damage
                fmt.Printf("Ranged attack! Troop ID=%d deals %d damage to Troop ID=%d (health now: %d)\n", 
                           troop2.ID, troop2.Damage, troop1.ID, troop1.Health)
                
                // Check if troop1 is defeated
                if troop1.Health <= 0 {
                    troop1.Active = false
                    fmt.Printf("Troop ID=%d defeated by Troop ID=%d\n", troop1.ID, troop2.ID)
                    // Clear attacking state of troop2 when target is defeated
                    troop2.IsAttacking = false
                }
            }
        }
    }
}

// Also update ProcessTroopBuildingCombat to ensure consistent behavior
func ProcessTroopBuildingCombat(game *Game, troop *Troop, building *Building) {
    // Set troop as attacking
    troop.IsAttacking = true
    
    // Get current game time
    currentTime := game.GameTime
    
    // Process attack from troop to building
    if currentTime - troop.LastAttack >= troop.AttackDelay {
        // Reset attack timer
        troop.LastAttack = currentTime
        
        // Check if troop is melee - if so, never use projectiles
        if IsMeleeTroop(troop) {
            // Melee troops apply damage directly
            building.Health -= troop.Damage
            fmt.Printf("Melee attack! Troop ID=%d deals %d damage to Building ID=%d (health now: %d)\n", 
                      troop.ID, troop.Damage, building.ID, building.Health)
            
            // Check if building is destroyed
            if building.Health <= 0 {
                building.Active = false
                fmt.Printf("Building ID=%d destroyed by Troop ID=%d\n", building.ID, troop.ID)
                // Clear attacking state when target is destroyed
                troop.IsAttacking = false
            }
        } else {
            // Get template for projectile information
            template := GetTroopTemplate(troop)
            
            // Check if troop has a projectile defined
            hasProjectile := template != nil && template.Projectile.Name != ""
            
            if hasProjectile {
                // Create the projectile
                projectile := CreateProjectile(
                    template.Projectile.Name,
                    troop.Position,
                    building.Position,
                    troop.Damage,  // USE TROOP'S DAMAGE
                    troop.Team,
                    troop.ID,
                )
                
                fmt.Printf("Troop ID=%d fires projectile at Building ID=%d\n", troop.ID, building.ID)
                
                // Set target entity for homing projectiles
                if projectile != nil && projectile.IsHoming {
                    projectile.TargetEntity = building
                }
                
                // Add projectile to game if created successfully
                if projectile != nil {
                    game.Projectiles = append(game.Projectiles, *projectile)
                } else {
                    // Fallback to direct damage if projectile creation failed
                    building.Health -= troop.Damage
                    fmt.Printf("Direct fallback! Troop ID=%d deals %d damage to Building ID=%d (health now: %d)\n", 
                               troop.ID, troop.Damage, building.ID, building.Health)
                }
            } else {
                // Ranged troops without projectiles defined fall back to direct damage
                building.Health -= troop.Damage
                fmt.Printf("Ranged attack! Troop ID=%d deals %d damage to Building ID=%d (health now: %d)\n", 
                           troop.ID, troop.Damage, building.ID, building.Health)
                
                // Check if building is destroyed
                if building.Health <= 0 {
                    building.Active = false
                    fmt.Printf("Building ID=%d destroyed by Troop ID=%d\n", building.ID, troop.ID)
                    // Clear attacking state when target is destroyed
                    troop.IsAttacking = false
                }
            }
        }
    }
}

// Add this to util.go
// IsMeleeTroop checks if a troop is a melee unit
func IsMeleeTroop(troop *Troop) bool {
    template := GetTroopTemplate(troop)
    if template != nil {
        // Use template range if available
        return template.Range <= 1.0
    }
    
    // Fallback to troop's range attribute
    return troop.Range <= 1.0
}

func ProcessBuildingTroopCombat(game *Game, building *Building, troop *Troop, buildingTeam int) {
    // Get current game time
    currentTime := game.GameTime
    
    // Buildings attack more slowly (every 60 ticks = 1 second)
    if currentTime - building.LastAttack >= 60 {
        // Reset attack timer
        building.LastAttack = currentTime
        
        // Buildings use projectiles
        projectileType := building.ProjectileType
        if projectileType == "" {
            projectileType = "normal" // Default
        }

        // Create projectile
        projectile := CreateProjectile(
            projectileType,
            building.Position,
            troop.Position,
            building.Damage,
            buildingTeam,
            0, // Buildings don't have IDs
        )
        
        // Only add valid projectiles to game
        if projectile != nil {
            game.Projectiles = append(game.Projectiles, *projectile)
            fmt.Printf("Building ID=%d fires projectile at Troop ID=%d\n", building.ID, troop.ID)
        } else {
            // Apply damage directly if projectile creation failed
            troop.Health -= building.Damage
            fmt.Printf("Direct attack! Building ID=%d deals %d damage to Troop ID=%d (health now: %d)\n", 
                      building.ID, building.Damage, troop.ID, troop.Health)
            
            // Check if troop is defeated
            if troop.Health <= 0 {
                troop.Active = false
                fmt.Printf("Troop ID=%d defeated by Building ID=%d\n", troop.ID, building.ID)
            }
        }
    }
}


func ClearInvalidAttackStates(game *Game) {
    // For each troop, check if it's still in combat range of any enemy
    for i := range game.Troops {
        troop := &game.Troops[i]
        
        // Skip inactive troops or those not in attacking state
        if !troop.Active || !troop.IsAttacking {
            continue
        }
        
        // Assume not in combat until we find a valid target
        inCombat := false
        
        // Check against other troops
        for j := range game.Troops {
            // Skip self, inactive troops, or same team troops
            if i == j || !game.Troops[j].Active || game.Troops[i].Team == game.Troops[j].Team {
                continue
            }
            
            // Check if in attack range
            if CanAttackTroop(troop, &game.Troops[j], game.Grid) {
                inCombat = true
                break
            }
        }
        
        // Also check for valid building targets
        if !inCombat && !TargetsOnlyTroops(troop) {
            // Get enemy team
            enemyTeam := 1 - troop.Team
            
            // Check king building
            kingBuilding := &game.Players[enemyTeam].KingBuilding.Building
            if kingBuilding.Active && CanTroopAttackBuilding(troop, kingBuilding, game.Grid) {
                inCombat = true
            }
            
            // Check regular buildings
            if !inCombat {
                for j := range game.Players[enemyTeam].Buildings {
                    building := &game.Players[enemyTeam].Buildings[j]
                    if building.Active && CanTroopAttackBuilding(troop, building, game.Grid) {
                        inCombat = true
                        break
                    }
                }
            }
        }
        
        // If not in combat with any valid target, clear attacking state
        if !inCombat {
            troop.IsAttacking = false
            fmt.Printf("Cleared attacking state for Troop ID=%d (no valid targets in range)\n", troop.ID)
        }
    }
}

// Add this new function to detect if a troop can attack another troop
func CanAttackTroop(troop1 *Troop, troop2 *Troop, grid *GridSystem) bool {
	// Check if troop types are compatible for combat
	// i.e., flying troops can only be attacked by troops that can attack air
	// ground troops can only be attacked by troops that can attack ground
	troop2Flying := IsFlyingTroop(troop2)
	if troop2Flying && !CanAttackAir(troop1) {
		return false
	}
	
	if !troop2Flying && !CanAttackGround(troop1) {
		return false
	}
	
	// Check if troop only targets buildings
	if TargetsOnlyBuildings(troop1) {
		return false
	}
	
	// Calculate distance between troops
	dist := Distance(troop1.Position, troop2.Position)
	
	// Convert attack range from grid cells to pixels
	attackRangePixels := troop1.Range * grid.CellWidth
	
	// Consider the target troop's size in the calculation
	targetTroopRadius := troop2.Size / 2
	
	// Check if troop1 can attack troop2
	// We can attack if distance - targetRadius <= attackRange
	return dist - targetTroopRadius <= attackRangePixels
}

// CanTroopAttackBuilding checks if a troop can attack a specific building
func CanTroopAttackBuilding(troop *Troop, building *Building, grid *GridSystem) bool {
    // Skip if troop only targets other troops
    if TargetsOnlyTroops(troop) {
        return false
    }
    
    // Calculate distance between troop and building
    dist := Distance(troop.Position, building.Position)
    
    // Get attack range in pixels
    attackRange := troop.Range * grid.CellWidth
    
    // Get building dimensions
    width, height := building.GetPixelDimensions(grid)
    buildingRadius := math.Max(width, height) / 2
    
    // A troop can attack if it's close enough (considering building size)
    return dist - buildingRadius <= attackRange
}

// FindTroopInBuildingRange finds the closest enemy troop in a building's attack range
func FindTroopInBuildingRange(game *Game, building *Building, buildingTeam int) *Troop {
    // Calculate attack range in pixels
    attackRange := building.Range * game.Grid.CellWidth
    
    var closestTroop *Troop
    closestDistance := attackRange + 1 // Start just outside range
    
    for i := range game.Troops {
        troop := &game.Troops[i]
        
        // Skip inactive or friendly troops
        if !troop.Active || troop.Team == buildingTeam {
            continue
        }
        
        // Skip flying troops if building can't attack air (placeholder logic)
        // For now assume all buildings can attack all troop types
        
        // Calculate distance
        dist := Distance(building.Position, troop.Position)
        
        // Check if in range and closer than current closest
        if dist <= attackRange && dist < closestDistance {
            closestTroop = troop
            closestDistance = dist
        }
    }
    
    return closestTroop
}

// CheckBuildingCombat handles buildings attacking troops
func CheckBuildingCombat(game *Game) {
    // Check each team's buildings
    for team := 0; team < 2; team++ {
        // Check king building
        kingBuilding := &game.Players[team].KingBuilding.Building
        if kingBuilding.Active {
            // Find closest enemy troop in range
            target := FindTroopInBuildingRange(game, kingBuilding, team)
            if target != nil {
                ProcessBuildingTroopCombat(game, kingBuilding, target, team)
            }
        }
        
        // Check regular buildings
        for i := range game.Players[team].Buildings {
            building := &game.Players[team].Buildings[i]
            if building.Active {
                // Find closest enemy troop in range
                target := FindTroopInBuildingRange(game, building, team)
                if target != nil {
                    ProcessBuildingTroopCombat(game, building, target, team)
                }
            }
        }
    }
}

// CheckProjectileBuildingCollisions checks if projectiles hit buildings
func CheckProjectileBuildingCollisions(game *Game) {
    for i := range game.Projectiles {
        projectile := &game.Projectiles[i]
        if !projectile.Active {
            continue
        }
        
        // Check against buildings for the team opposite to the projectile
        enemyTeam := 1 - projectile.Team
        
        // Check king building
        kingBuilding := &game.Players[enemyTeam].KingBuilding.Building
        if kingBuilding.Active {
            width, height := kingBuilding.GetPixelDimensions(game.Grid)
            
            // Calculate distance to building center
            dist := Distance(projectile.Position, kingBuilding.Position)
            
            // Building radius (approximate as half the largest dimension)
            buildingRadius := math.Max(width, height) / 2
            
            // Check if projectile hits building
            if dist <= buildingRadius + projectile.Size/2 {
                // Handle impact
                projectile.HandleImpact(game, projectile.Position)
                break // Projectile can only hit one target
            }
        }
        
        // Check regular buildings
        for j := range game.Players[enemyTeam].Buildings {
            building := &game.Players[enemyTeam].Buildings[j]
            if building.Active {
                width, height := building.GetPixelDimensions(game.Grid)
                
                // Calculate distance to building center
                dist := Distance(projectile.Position, building.Position)
                
                // Building radius (approximate as half the largest dimension)
                buildingRadius := math.Max(width, height) / 2
                
                // Check if projectile hits building
                if dist <= buildingRadius + projectile.Size/2 {
                    // Handle impact
                    projectile.HandleImpact(game, projectile.Position)
                    break // Projectile can only hit one target
                }
            }
        }
    }
}