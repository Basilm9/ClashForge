// movement.go
package clashgame

import (
	"container/heap"
	"fmt"
	"math"
)

// Node represents a position in the pathfinding grid
type Node struct {
	Col, Row int
	F, G, H  float64 // F = G + H
	Parent   *Node
	index    int // For priority queue
}

// PriorityQueue implementation for A* algorithm
type PriorityQueue []*Node

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].F < pq[j].F
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *PriorityQueue) Push(x interface{}) {
	n := len(*pq)
	node := x.(*Node)
	node.index = n
	*pq = append(*pq, node)
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	node := old[n-1]
	old[n-1] = nil
	node.index = -1
	*pq = old[0 : n-1]
	return node
}

// FindPath uses A* to find a path between two points
func FindPath(game *Game, start, target Position) []Position {
	grid := game.Grid
	
	// Convert positions to grid coordinates
	startCol, startRow := grid.PositionToCell(start)
	targetCol, targetRow := grid.PositionToCell(target)
	
	// Check if start or target is invalid
	if !grid.IsWalkableTile(startCol, startRow) {
		return nil
	}
	if !grid.IsWalkableTile(targetCol, targetRow) {
		return nil
	}
	
	// Initialize the open and closed sets
	openSet := &PriorityQueue{}
	heap.Init(openSet)
	closedSet := make(map[string]bool)
	
	// Create the start node
	startNode := &Node{
		Col: startCol,
		Row: startRow,
		G:   0,
		H:   heuristic(startCol, startRow, targetCol, targetRow),
	}
	startNode.F = startNode.G + startNode.H
	
	// Add start node to open set
	heap.Push(openSet, startNode)
	
	// Keep track of all nodes for path reconstruction
	allNodes := make(map[string]*Node)
	allNodes[nodeKey(startCol, startRow)] = startNode
	
	// Define possible movements (8 directions)
	directions := [][2]int{
		{-1, -1}, {0, -1}, {1, -1},
		{-1, 0}, {1, 0},
		{-1, 1}, {0, 1}, {1, 1},
	}
	
	// Main loop
	for openSet.Len() > 0 {
		current := heap.Pop(openSet).(*Node)
		
		// Check if we reached the target
		if current.Col == targetCol && current.Row == targetRow {
			return reconstructPath(current, grid)
		}
		
		// Add to closed set
		closedSet[nodeKey(current.Col, current.Row)] = true
		
		// Check each neighbor
		for _, dir := range directions {
			newCol := current.Col + dir[0]
			newRow := current.Row + dir[1]
			
			// Skip if neighbor is not walkable or in closed set
			if !grid.IsWalkableTile(newCol, newRow) || 
			   closedSet[nodeKey(newCol, newRow)] {
				continue
			}
			
			// Calculate cost to this neighbor
			moveCost := 1.0
			if dir[0] != 0 && dir[1] != 0 {
				moveCost = 1.414 // Diagonal movement cost
			}
			
			// Add additional cost for crossing water/bridges
			tileType := grid.GetTileType(newCol, newRow)
			if tileType == TileBridge1 || tileType == TileBridge2 {
				moveCost *= 1.5 // Make bridges slightly more "expensive"
			}
			
			tentativeG := current.G + moveCost
			
			neighbor := allNodes[nodeKey(newCol, newRow)]
			if neighbor == nil {
				// Create new neighbor node
				neighbor = &Node{
					Col:    newCol,
					Row:    newRow,
					G:      tentativeG,
					H:      heuristic(newCol, newRow, targetCol, targetRow),
					Parent: current,
				}
				neighbor.F = neighbor.G + neighbor.H
				allNodes[nodeKey(newCol, newRow)] = neighbor
				heap.Push(openSet, neighbor)
			} else if tentativeG < neighbor.G {
				// Update existing neighbor with better path
				neighbor.G = tentativeG
				neighbor.F = tentativeG + neighbor.H
				neighbor.Parent = current
				heap.Fix(openSet, neighbor.index)
			}
		}
	}
	
	// No path found
	return nil
}

// Helper function to calculate heuristic (Manhattan distance)
func heuristic(x1, y1, x2, y2 int) float64 {
	dx := math.Abs(float64(x2 - x1))
	dy := math.Abs(float64(y2 - y1))
	return dx + dy
}

// Helper function to generate unique key for a node
func nodeKey(col, row int) string {
	return fmt.Sprintf("%d,%d", col, row)
}

// Helper function to reconstruct path from A* result
func reconstructPath(endNode *Node, grid *GridSystem) []Position {
	path := make([]Position, 0)
	current := endNode
	
	for current != nil {
		pos := grid.CellToPosition(current.Col, current.Row)
		path = append([]Position{pos}, path...)
		current = current.Parent
	}
	
	return path
}

// Steering behaviors for natural movement and collision avoidance
type SteeringForce struct {
	Separation Position
	Alignment  Position
	Cohesion   Position
}

// Calculate separation force to avoid collisions with nearby troops
func calculateSeparation(game *Game, troop *Troop) Position {
	var steeringForce Position
	count := 0
	
	// Check radius for separation (smaller than aggro radius)
	separationRadius := troop.AggroDistance * 0.5 * game.Grid.CellWidth
	
	for i := range game.Troops {
		other := &game.Troops[i]
		
		// Skip if not active or same troop
		if !other.Active || other.ID == troop.ID {
			continue
		}
		
		// Calculate distance to other troop
		dist := Distance(troop.Position, other.Position)
		
		// If within separation radius, add repulsion force
		if dist < separationRadius {
			// Calculate vector pointing away from other troop
			dx := troop.Position.X - other.Position.X
			dy := troop.Position.Y - other.Position.Y
			
			// Normalize and scale by distance (closer = stronger force)
			if dist > 0 {
				strength := (separationRadius - dist) / separationRadius
				steeringForce.X += (dx / dist) * strength
				steeringForce.Y += (dy / dist) * strength
				count++
			}
		}
	}
	
	// Average the steering force
	if count > 0 {
		steeringForce.X /= float64(count)
		steeringForce.Y /= float64(count)
	}
	
	return steeringForce
}

// Calculate alignment force to match velocity with nearby troops
func calculateAlignment(game *Game, troop *Troop) Position {
	var steeringForce Position
	count := 0
	
	// Check radius for alignment (similar to separation)
	alignmentRadius := troop.AggroDistance * 0.7 * game.Grid.CellWidth
	
	for i := range game.Troops {
		other := &game.Troops[i]
		
		// Skip if not active, same troop, or different team
		if !other.Active || other.ID == troop.ID || other.Team != troop.Team {
			continue
		}
		
		// Calculate distance to other troop
		dist := Distance(troop.Position, other.Position)
		
		// If within alignment radius, add alignment force
		if dist < alignmentRadius {
			steeringForce.X += other.Velocity.X
			steeringForce.Y += other.Velocity.Y
			count++
		}
	}
	
	// Average the velocities
	if count > 0 {
		steeringForce.X /= float64(count)
		steeringForce.Y /= float64(count)
		
		// Calculate steering force to match average velocity
		steeringForce.X = (steeringForce.X - troop.Velocity.X) * 0.1
		steeringForce.Y = (steeringForce.Y - troop.Velocity.Y) * 0.1
	}
	
	return steeringForce
}

// Calculate cohesion force to move towards center of nearby troops
func calculateCohesion(game *Game, troop *Troop) Position {
	var centerOfMass Position
	count := 0
	
	// Check radius for cohesion (larger than alignment)
	cohesionRadius := troop.AggroDistance * 0.9 * game.Grid.CellWidth
	
	for i := range game.Troops {
		other := &game.Troops[i]
		
		// Skip if not active, same troop, or different team
		if !other.Active || other.ID == troop.ID || other.Team != troop.Team {
			continue
		}
		
		// Calculate distance to other troop
		dist := Distance(troop.Position, other.Position)
		
		// If within cohesion radius, add to center of mass
		if dist < cohesionRadius {
			centerOfMass.X += other.Position.X
			centerOfMass.Y += other.Position.Y
			count++
		}
	}
	
	// Calculate steering force towards center of mass
	if count > 0 {
		centerOfMass.X /= float64(count)
		centerOfMass.Y /= float64(count)
		
		// Calculate vector towards center
		dx := centerOfMass.X - troop.Position.X
		dy := centerOfMass.Y - troop.Position.Y
		dist := math.Sqrt(dx*dx + dy*dy)
		
		if dist > 0 {
			// Normalize and scale the steering force
			steeringForce := Position{
				X: (dx / dist) * 0.05,
				Y: (dy / dist) * 0.05,
			}
			return steeringForce
		}
	}
	
	return Position{X: 0, Y: 0}
}

// Calculate all steering forces for a troop
func calculateSteeringForces(game *Game, troop *Troop) SteeringForce {
	return SteeringForce{
		Separation: calculateSeparation(game, troop),
		Alignment:  calculateAlignment(game, troop),
		Cohesion:   calculateCohesion(game, troop),
	}
}

func UpdateTroopMovement(game *Game) {
	for i := range game.Troops {
		troop := &game.Troops[i]
		if !troop.Active {
			continue
		}
		
		// Skip if troop is attacking
		if troop.IsAttacking {
			continue
		}
		
		var targetPos Position
		var shouldMove bool = true
		var inAttackRange bool = false
		
		// First try to find nearest enemy or building within aggro range
		if enemyTroop := FindNearestEnemyTroop(game, troop); enemyTroop != nil {
			targetPos = enemyTroop.Position
			// Check if in attack range
			dist := Distance(troop.Position, enemyTroop.Position)
			
			// Calculate combined ranges (troop's collision radius + enemy's attack range)
			combinedRange := troop.Size/2 + enemyTroop.Range * game.Grid.CellWidth
			
			if dist <= combinedRange {
				// In attack range, stop moving but still apply push forces
				troop.Velocity = Position{X: 0, Y: 0}
				inAttackRange = true
			}
		} else if building, _ := FindNearestEnemyBuilding(game, troop); building != nil {
			targetPos = building.Position
			// Check if in attack range
			dist := Distance(troop.Position, building.Position)
			
			// Calculate combined ranges (troop's collision radius + building's attack range)
			combinedRange := troop.Size/2 + building.Range * game.Grid.CellWidth
			
			if dist <= combinedRange {
				// In attack range, stop moving but still apply push forces
				troop.Velocity = Position{X: 0, Y: 0}
				inAttackRange = true
			}
		} else {
			// No immediate targets in range, move towards enemy base
			shouldMove = true
			
			// Determine enemy team's buildings
			enemyTeam := 1 - troop.Team
			
			// First try to target nearest princess tower
			var nearestPrincessDist float64 = math.MaxFloat64
			var nearestPrincessPos Position
			
			for _, building := range game.Players[enemyTeam].Buildings {
				if !building.Active {
					continue
				}
				
				dist := Distance(troop.Position, building.Position)
				if dist < nearestPrincessDist {
					nearestPrincessDist = dist
					nearestPrincessPos = building.Position
				}
			}
			
			// If we found an active princess tower, target it
			if nearestPrincessDist != math.MaxFloat64 {
				// If we need to cross river to reach princess tower
				if NeedsToCrossBridge(game, troop.Position, nearestPrincessPos) {
					bridgePos := findNearestBridge(game, troop.Position)
					targetPos = bridgePos
				} else {
					targetPos = nearestPrincessPos
				}
			} else {
				// No princess towers left, go for king tower
				kingBuilding := &game.Players[enemyTeam].KingBuilding.Building
				
				// If we need to cross river to reach king tower
				if NeedsToCrossBridge(game, troop.Position, kingBuilding.Position) {
					bridgePos := findNearestBridge(game, troop.Position)
					targetPos = bridgePos
				} else {
					targetPos = kingBuilding.Position
				}
			}
		}
		
		// Always apply push forces, even when in attack range
		applyPushForces(game, troop)
		
		if shouldMove && !inAttackRange {
			// Find path to target
			path := FindPath(game, troop.Position, targetPos)
			
			var nextPos Position
			if path == nil || len(path) < 2 {
				// Pathfinding failed, use direct movement towards target
				dx := targetPos.X - troop.Position.X
				dy := targetPos.Y - troop.Position.Y
				dist := math.Sqrt(dx*dx + dy*dy)
				
				if dist > 0 {
					// Move directly towards target
					nextPos = Position{
						X: troop.Position.X + (dx/dist) * game.Grid.CellWidth,
						Y: troop.Position.Y + (dy/dist) * game.Grid.CellHeight,
					}
				} else {
					// Already at target, don't move
					continue
				}
			} else {
				nextPos = path[1]
			}
			
			// Calculate movement direction
			dx := nextPos.X - troop.Position.X
			dy := nextPos.Y - troop.Position.Y
			
			// Normalize direction
			dist := math.Sqrt(dx*dx + dy*dy)
			if dist > 0 {
				speed := troop.Speed * game.Grid.CellWidth
				
				// Calculate target velocity
				targetVelX := (dx / dist) * speed
				targetVelY := (dy / dist) * speed
				
				// Calculate steering forces
				steering := calculateSteeringForces(game, troop)
				
				// Combine steering forces with weights
				steeringForce := Position{
					X: steering.Separation.X*0.6 + // Reduced separation weight
					   steering.Alignment.X*0.7 +  // Increased alignment weight
					   steering.Cohesion.X*0.5,    // Slightly reduced cohesion
					Y: steering.Separation.Y*0.6 +
					   steering.Alignment.Y*0.7 +
					   steering.Cohesion.Y*0.5,
				}
				
				// Apply steering forces to target velocity
				targetVelX += steeringForce.X * speed
				targetVelY += steeringForce.Y * speed
				
				// Normalize final velocity
				finalDist := math.Sqrt(targetVelX*targetVelX + targetVelY*targetVelY)
				if finalDist > 0 {
					targetVelX = (targetVelX / finalDist) * speed
					targetVelY = (targetVelY / finalDist) * speed
				}
				
				// Add minimum velocity threshold to prevent stopping
				minVelocity := speed * 0.1 // 10% of max speed
				currentSpeed := math.Sqrt(troop.Velocity.X*troop.Velocity.X + troop.Velocity.Y*troop.Velocity.Y)
				
				// If current speed is too low, ensure we maintain minimum velocity
				if currentSpeed < minVelocity {
					// Use the last known good direction or target direction
					if currentSpeed > 0 {
						// Normalize current velocity
						troop.Velocity.X = (troop.Velocity.X / currentSpeed) * minVelocity
						troop.Velocity.Y = (troop.Velocity.Y / currentSpeed) * minVelocity
					} else {
						// Use target direction
						troop.Velocity.X = targetVelX
						troop.Velocity.Y = targetVelY
					}
				}
				
				// Smoothly adjust velocity with increased responsiveness
				accelerationFactor := troop.MaxAcceleration * 1.5 // Increase responsiveness
				troop.Velocity.X += (targetVelX - troop.Velocity.X) * accelerationFactor
				troop.Velocity.Y += (targetVelY - troop.Velocity.Y) * accelerationFactor
				
				// Apply velocity to position
				troop.Position.X += troop.Velocity.X
				troop.Position.Y += troop.Velocity.Y
			}
		}
		
		// Update position history for smooth rendering
		historyIndex := troop.PositionHistory.Index
		troop.PositionHistory.Positions[historyIndex] = troop.Position
		troop.PositionHistory.Index = (historyIndex + 1) % PositionHistoryLength
	}
}

// Helper function to find the nearest bridge position
func findNearestBridge(game *Game, pos Position) Position {
	// Bridge positions (from grid.go)
	leftBridgeCol := 5
	rightBridgeCol := 27
	bridgeWidth := 4
	waterStartRow := 30
	waterEndRow := 33
	
	// Calculate center row of water
	bridgeRow := (waterStartRow + waterEndRow) / 2
	
	// Convert current position to grid coordinates
	currentCol, _ := game.Grid.PositionToCell(pos)
	
	// Determine which bridge is closer
	var targetBridgeCol int
	if currentCol < GridColumns/2 {
		// Use left bridge
		targetBridgeCol = leftBridgeCol + bridgeWidth/2
	} else {
		// Use right bridge
		targetBridgeCol = rightBridgeCol + bridgeWidth/2
	}
	
	// Convert bridge position to world coordinates
	return game.Grid.CellToPosition(targetBridgeCol, bridgeRow)
}

// NeedsToCrossBridge checks if a path from troopPos to targetPos
// would need to cross the river
func NeedsToCrossBridge(game *Game, troopPos, targetPos Position) bool {
	// Flying troops don't need bridges
	if IsFlyingTroop(&Troop{Position: troopPos}) {
		return false
	}
	
	// Water row range (from grid.go)
	waterStartRow := 30
	waterEndRow := 33
	
	// Get troop and target rows
	_, troopRow := game.Grid.PositionToCell(troopPos)
	_, targetRow := game.Grid.PositionToCell(targetPos)
	
	// If troop is on a bridge, it doesn't need to find one
	troopCol, _ := game.Grid.PositionToCell(troopPos)
	if troopRow >= waterStartRow && troopRow <= waterEndRow {
		cellType := game.Grid.GetCellType(troopCol, troopRow)
		if cellType == CellTypeBridge {
			return false
		}
	}
	
	// Check if the path crosses the river by seeing if troop and target
	// are on opposite sides of the river section
	if troopRow < waterStartRow && targetRow > waterEndRow {
		return true // Need bridge to go from top to bottom
	}
	if troopRow > waterEndRow && targetRow < waterStartRow {
		return true // Need bridge to go from bottom to top
	}
	
	return false // No need for bridge if on same side
}

// InitTroopMovement - Call this when creating a new troop
func InitTroopMovement(troop *Troop) {
	// Initialize velocity
	troop.Velocity = Position{X: 0, Y: 0}
	troop.TargetVelocity = Position{X: 0, Y: 0}
	
	// Set acceleration based on troop type
	// Faster troops should have higher acceleration
	troop.MaxAcceleration = 0.2
	
	if troop.Speed > 0.15 {
		// Fast troops have more responsive movement
		troop.MaxAcceleration = 0.25
	} else if troop.Speed < 0.1 {
		// Slow, heavy troops have less responsive movement
		troop.MaxAcceleration = 0.1
	}
	
	// Initialize position history for smooth rendering
	troop.PositionHistory = TroopPositionHistory{
		Positions: make([]Position, PositionHistoryLength),
		Index:     0,
	}
	
	// Fill history with current position
	for i := 0; i < PositionHistoryLength; i++ {
		troop.PositionHistory.Positions[i] = troop.Position
	}
}

// Add these constants at the top of the file
const (
	// Existing constants
	SameTeamPushFactor    = 0.5
	DifferentTeamPushFactor = 0.5  
	
	// New constants for rigid body physics
	PushForceStrength = 0.8
	PushRadius = 0.8  // Slightly larger than separation radius
	PushDamping = 0.8 // How quickly push forces dissipate
)

// Add this new function for rigid body physics
func applyPushForces(game *Game, troop *Troop) {
	pushRadius := PushRadius * game.Grid.CellWidth
	
	for i := range game.Troops {
		other := &game.Troops[i]
		
		// Skip if not active or same troop
		if !other.Active || other.ID == troop.ID {
			continue
		}
		
		// Calculate distance and direction
		dx := other.Position.X - troop.Position.X
		dy := other.Position.Y - troop.Position.Y
		dist := math.Sqrt(dx*dx + dy*dy)
		
		// If within push radius, apply force
		if dist < pushRadius {
			// Calculate overlap
			overlap := pushRadius - dist
			
			// Determine push factor based on team
			pushFactor := DifferentTeamPushFactor
			if other.Team == troop.Team {
				pushFactor = SameTeamPushFactor
			}
			
			// Calculate push force (stronger when closer)
			force := (overlap / pushRadius) * PushForceStrength * pushFactor
			
			// Normalize direction
			if dist > 0 {
				dx = dx / dist
				dy = dy / dist
			}
			
			// Apply push force
			pushX := dx * force * game.Grid.CellWidth
			pushY := dy * force * game.Grid.CellHeight
			
			// Apply to both troops (equal and opposite forces)
			troop.Position.X -= pushX * PushDamping
			troop.Position.Y -= pushY * PushDamping
			other.Position.X += pushX * PushDamping
			other.Position.Y += pushY * PushDamping
			
			// Also affect velocities slightly
			if dist > 0 {
				// Transfer some momentum
				momentumTransfer := 0.3
				troop.Velocity.X -= dx * momentumTransfer * force
				troop.Velocity.Y -= dy * momentumTransfer * force
				other.Velocity.X += dx * momentumTransfer * force
				other.Velocity.Y += dy * momentumTransfer * force
			}
		}
	}
}