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
	if !grid.IsWalkableTile(startCol, startRow) || !grid.IsWalkableTile(targetCol, targetRow) {
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
		
		// First try to find nearest enemy or building within aggro range
		if enemyTroop := FindNearestEnemyTroop(game, troop); enemyTroop != nil {
			targetPos = enemyTroop.Position
		} else if building, _ := FindNearestEnemyBuilding(game, troop); building != nil {
			targetPos = building.Position
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
		
		if shouldMove {
			// Find path to target
			path := FindPath(game, troop.Position, targetPos)
			if path == nil || len(path) < 2 {
				continue
			}
			
			// Move towards next point in path
			nextPos := path[1]
			dx := nextPos.X - troop.Position.X
			dy := nextPos.Y - troop.Position.Y
			
			// Normalize direction
			dist := math.Sqrt(dx*dx + dy*dy)
			if dist > 0 {
				speed := troop.Speed * game.Grid.CellWidth
				
				// Update velocity with smoothing
				targetVelX := (dx / dist) * speed
				targetVelY := (dy / dist) * speed
				
				// Smoothly adjust velocity
				troop.Velocity.X += (targetVelX - troop.Velocity.X) * troop.MaxAcceleration
				troop.Velocity.Y += (targetVelY - troop.Velocity.Y) * troop.MaxAcceleration
				
				// Apply velocity to position
				troop.Position.X += troop.Velocity.X
				troop.Position.Y += troop.Velocity.Y
				
				// Update position history for smooth rendering
				historyIndex := troop.PositionHistory.Index
				troop.PositionHistory.Positions[historyIndex] = troop.Position
				troop.PositionHistory.Index = (historyIndex + 1) % PositionHistoryLength
			}
		}
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
