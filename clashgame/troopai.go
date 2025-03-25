package clashgame

import "math"

// FindNearestEnemyTroop finds the closest enemy troop within the aggro radius
// that isn't separated by water (on the same side of the river)
func FindNearestEnemyTroop(game *Game, troop *Troop) *Troop {
    // Convert aggro radius from tiles to pixels
    aggroRadius := troop.AggroDistance * game.Grid.CellWidth
    
    // Check if troop is flying
    troopIsFlying := IsFlyingTroop(troop)
    
    // Determine which side of the river the troop is on
    // We use the water rows (15 and 16) as dividing line
    waterStartRow := 15
    waterEndRow := 16
    
    _, troopRow := game.Grid.PositionToCell(troop.Position)
    var troopSide int
    if troopRow < waterStartRow {
        troopSide = 0 // Top side
    } else if troopRow > waterEndRow {
        troopSide = 1 // Bottom side
    } else {
        // Troop is on a bridge - can see both sides
        troopSide = 2
    }
    
    var closestTroop *Troop
    closestDistance := aggroRadius + 1 // Start outside aggro radius
    
    // Check all other troops
    for i := range game.Troops {
        otherTroop := &game.Troops[i]
        
        // Skip inactive, same team, or same troop
        if !otherTroop.Active || 
           otherTroop.Team == troop.Team || 
           otherTroop.ID == troop.ID {
            continue
        }
        
        // Calculate distance
        dist := Distance(troop.Position, otherTroop.Position)
        
        // Skip if outside aggro radius
        if dist > aggroRadius {
            continue
        }
        
        // Check if troop can attack this type of enemy
        otherIsFlying := IsFlyingTroop(otherTroop)
        if otherIsFlying && !CanAttackAir(troop) {
            continue // Skip flying troops if we can't attack air
        }
        if !otherIsFlying && !CanAttackGround(troop) {
            continue // Skip ground troops if we can't attack ground
        }
        
        // Check river crossing for ground troops
        if !troopIsFlying && !otherIsFlying {
            _, otherRow := game.Grid.PositionToCell(otherTroop.Position)
            
            var otherSide int
            if otherRow < waterStartRow {
                otherSide = 0 // Top side
            } else if otherRow > waterEndRow {
                otherSide = 1 // Bottom side
            } else {
                // On a bridge
                otherSide = 2
            }
            
            // Skip if on opposite sides of river (unless one is on a bridge)
            if troopSide != 2 && otherSide != 2 && troopSide != otherSide {
                continue
            }
        }
        
        // This troop is a valid target - check if it's the closest
        if dist < closestDistance {
            closestTroop = otherTroop
            closestDistance = dist
        }
    }
    
    return closestTroop
}

// FindNearestEnemyBuilding finds the closest enemy building within the aggro radius
// that isn't separated by water (on the same side of the river)
func FindNearestEnemyBuilding(game *Game, troop *Troop) (*Building, bool) {
    // Convert aggro radius from tiles to pixels
    aggroRadius := troop.AggroDistance * game.Grid.CellWidth * 3
    
    // Check if troop is flying (flying troops can target across river)
    troopIsFlying := IsFlyingTroop(troop)
    
    // Determine which side of the river the troop is on
    waterStartRow := 15
    waterEndRow := 16
    
    _, troopRow := game.Grid.PositionToCell(troop.Position)
    var troopSide int
    if troopRow < waterStartRow {
        troopSide = 0 // Top side
    } else if troopRow > waterEndRow {
        troopSide = 1 // Bottom side
    } else {
        // Troop is on a bridge - can see both sides
        troopSide = 2
    }
    
    var closestBuilding *Building
    var isKingBuilding bool
    closestDistance := aggroRadius + 1 // Start outside aggro radius
    
    // Get enemy team
    enemyTeam := 1 - troop.Team
    
    // Skip if troop only targets troops
    if TargetsOnlyTroops(troop) {
        return nil, false
    }
    
    // Check regular buildings first
    for i := range game.Players[enemyTeam].Buildings {
        building := &game.Players[enemyTeam].Buildings[i]
        
        // Skip inactive buildings
        if !building.Active {
            continue
        }
        
        // Calculate distance, taking into account building size
        buildingWidth, buildingHeight := building.GetPixelDimensions(game.Grid)
        buildingRadius := math.Max(buildingWidth, buildingHeight) / 2
        dist := Distance(troop.Position, building.Position) - buildingRadius
        if dist < 0 {
            dist = 0 // Already touching the building
        }
        
        // Skip if outside aggro radius
        if dist > aggroRadius {
            continue
        }
        
        // If ground troop, check river crossing
        if !troopIsFlying {
            _, buildingRow := game.Grid.PositionToCell(building.Position)
            
            var buildingSide int
            if buildingRow < waterStartRow {
                buildingSide = 0 // Top side
            } else if buildingRow > waterEndRow {
                buildingSide = 1 // Bottom side
            } else {
                // On a bridge (unlikely for a building, but just in case)
                buildingSide = 2
            }
            
            // Skip if on opposite sides of river (unless one is on a bridge)
            if troopSide != 2 && buildingSide != 2 && troopSide != buildingSide {
                continue
            }
        }
        
        // This building is a valid target - check if it's the closest
        if dist < closestDistance {
            closestBuilding = building
            closestDistance = dist
            isKingBuilding = false
        }
    }
    
    // Check king building
    kingBuilding := &game.Players[enemyTeam].KingBuilding.Building
    if kingBuilding.Active {
        // Calculate distance, taking into account building size
        buildingWidth, buildingHeight := kingBuilding.GetPixelDimensions(game.Grid)
        buildingRadius := math.Max(buildingWidth, buildingHeight) / 2
        dist := Distance(troop.Position, kingBuilding.Position) - buildingRadius
        if dist < 0 {
            dist = 0 // Already touching the building
        }
        
        // Check if this is a closer valid target
        if dist <= aggroRadius {
            // Check river crossing for ground troops
            if troopIsFlying {
                // Flying troops can always see the king
                if dist < closestDistance {
                    closestBuilding = kingBuilding
                    closestDistance = dist
                    isKingBuilding = true
                }
            } else {
                // Ground troops need to check river
                _, kingRow := game.Grid.PositionToCell(kingBuilding.Position)
                
                var kingSide int
                if kingRow < waterStartRow {
                    kingSide = 0 // Top side
                } else if kingRow > waterEndRow {
                    kingSide = 1 // Bottom side
                } else {
                    kingSide = 2 // On bridge (unlikely for king)
                }
                
                // Check if on same side of river
                if troopSide == 2 || kingSide == 2 || troopSide == kingSide {
                    if dist < closestDistance {
                        closestBuilding = kingBuilding
                        closestDistance = dist
                        isKingBuilding = true
                    }
                }
            }
        }
    }
    
    return closestBuilding, isKingBuilding
}