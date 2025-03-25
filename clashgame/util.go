package clashgame

import "math"

func Distance(a, b Position) float64 {
	dx := b.X - a.X
	dy := b.Y - a.Y
	return math.Sqrt(dx*dx + dy*dy)
}


func UpdateElixir(game *Game) {
    for i := range game.Players {
        player := &game.Players[i]
        player.Elixir = math.Min(player.Elixir+player.ElixirGenRate, float64(player.ElixirMax))
    }
}

func (game *Game) IsActive() bool {
	return game.Running
}

// Add this helper function to check if a position is valid (not in water)
func IsValidPosition(grid *GridSystem, position Position) bool {
	// Convert position to grid cell
	col, row := grid.PositionToCell(position)
	
	// Check if this cell is valid (not water)
	return grid.GetCellType(col, row) != CellTypeWater
}

func (g *Game) PlaceExtendedTroopAtCell(col, row int, troopName string, team int) {
	// Get position from grid
	pos := g.Grid.CellToPosition(col, row)
	
	// Create and add the troop
	SpawnExtendedTroop(troopName, pos.X, pos.Y, team, g)
}

// SpawnRandomTroops creates a random assortment of troops on the field
// Useful for testing
func (g *Game) SpawnRandomTroops(count int) {
	// Get list of available troop names
	troopNames := make([]string, 0, len(TroopTemplateMap))
	for name := range TroopTemplateMap {
		troopNames = append(troopNames, name)
	}
	
	// Helper function to get a random troop name
	getRandomTroopName := func() string {
		// Use GameTime as a simple "random" source
		g.GameTime++
		index := g.GameTime % len(troopNames)
		return troopNames[index]
	}
	
	// Spawn troops for team 0 (friendly)
	for i := 0; i < count; i++ {
		// Random position in top third of map
		col := (g.GameTime*3 + i*5) % GridColumns
		row := (g.GameTime + i*3) % (GridRows/3)
		
		// Get random troop
		troopName := getRandomTroopName()
		
		// Spawn troop
		g.PlaceExtendedTroopAtCell(col, row, troopName, 0)
	}
	
	// Spawn troops for team 1 (enemy)
	for i := 0; i < count; i++ {
		// Random position in bottom third of map
		col := (g.GameTime*7 + i*11) % GridColumns
		row := GridRows - 1 - ((g.GameTime + i*7) % (GridRows/3))
		
		// Get random troop
		troopName := getRandomTroopName()
		
		// Spawn troop
		g.PlaceExtendedTroopAtCell(col, row, troopName, 1)
	}
}

func GetTroopTemplate(troop *Troop) *TroopTemplate {
	if troop == nil {
		return nil
	}
	
	// Look up template by name in the global template map
	template, exists := TroopTemplateMap[troop.Name]
	if !exists {
		return nil
	}
	
	return template
}

// IsFlyingTroop checks if a troop is a flying unit
func IsFlyingTroop(troop *Troop) bool {
	template := GetTroopTemplate(troop)
	if template != nil {
		return template.FlyingHeight > 0
	}
	
	// Fallback method if template isn't found
	flyingTypes := map[string]bool{
		"BabyDragon": true,
		"Dragon": true,
		"Balloon": true,
		"MinionHorde": true,
		"Minion": true,
		"LavaHound": true,
	}
	
	return flyingTypes[troop.Name]
}

// CanAttackAir checks if a troop can attack air units
func CanAttackAir(troop *Troop) bool {
	template := GetTroopTemplate(troop)
	if template != nil {
		return template.AttacksAir
	}
	
	// Fallback based on unit name
	airAttackers := map[string]bool{
		"Archer": true,
		"BabyDragon": true,
		"Dragon": true,
		"Wizard": true,
		"IceWizard": true,
		"Musketeer": true,
		"MegaMinion": true,
		"Minion": true,
		"MinionHorde": true,
	}
	
	return airAttackers[troop.Name]
}

// CanAttackGround checks if a troop can attack ground units
func CanAttackGround(troop *Troop) bool {
	template := GetTroopTemplate(troop)
	if template != nil {
		return template.AttacksGround
	}
	
	// Most troops can attack ground units by default
	// Only list exceptions
	nonGroundAttackers := map[string]bool{
		"Balloon": true, // Balloon only targets buildings
	}
	
	return !nonGroundAttackers[troop.Name]
}

// TargetsOnlyBuildings checks if a troop only targets buildings
func TargetsOnlyBuildings(troop *Troop) bool {
	template := GetTroopTemplate(troop)
	if template != nil {
		return template.TargetOnlyBuildings
	}
	
	buildingTargeters := map[string]bool{
		"Giant": true,
		"Balloon": true,
		"Golem": true,
		"LavaHound": true,
	}
	
	return buildingTargeters[troop.Name]
}

// TargetsOnlyTroops checks if a troop only targets other troops
func TargetsOnlyTroops(troop *Troop) bool {
	template := GetTroopTemplate(troop)
	if template != nil {
		return template.TargetOnlyTroops
	}
	
	troopTargeters := map[string]bool{
		"Valkyrie": true,
	}
	
	return troopTargeters[troop.Name]
}

// HasAreaDamage checks if a troop deals area damage
func HasAreaDamage(troop *Troop) bool {
	template := GetTroopTemplate(troop)
	if template != nil {
		return template.AreaDamageRadius > 0
	}
	
	areaDamagers := map[string]bool{
		"BabyDragon": true,
		"Valkyrie": true,
		"Bomber": true,
		"Wizard": true,
	}
	
	return areaDamagers[troop.Name]
}

// GetAreaDamageRadius returns the area damage radius
func GetAreaDamageRadius(troop *Troop) float64 {
	template := GetTroopTemplate(troop)
	if template != nil {
		return template.AreaDamageRadius
	}
	
	// Default radius values if template not found
	areaDamageRadii := map[string]float64{
		"BabyDragon": 1.5,
		"Valkyrie": 1.2,
		"Bomber": 1.0,
		"Wizard": 1.2,
	}
	
	return areaDamageRadii[troop.Name]
}

// HasDeathDamage checks if a troop deals damage on death
func HasDeathDamage(troop *Troop) bool {
	template := GetTroopTemplate(troop)
	if template != nil {
		return template.DeathDamage > 0 && template.DeathDamageRadius > 0
	}
	
	deathDamagers := map[string]bool{
		"Balloon": true,
		"GolemiteGolem": true,
		"LavaHound": true,
	}
	
	return deathDamagers[troop.Name]
}

// GetDeathDamage returns the damage dealt on death
func GetDeathDamage(troop *Troop) int {
	template := GetTroopTemplate(troop)
	if template != nil {
		return template.DeathDamage
	}
	
	// Default death damage values
	deathDamages := map[string]int{
		"Balloon": 100,
		"GolemiteGolem": 60,
		"LavaHound": 150,
	}
	
	return deathDamages[troop.Name]
}

// GetDeathDamageRadius returns the radius of death damage
func GetDeathDamageRadius(troop *Troop) float64 {
	template := GetTroopTemplate(troop)
	if template != nil {
		return template.DeathDamageRadius
	}
	
	// Default death damage radius values
	deathRadii := map[string]float64{
		"Balloon": 1.8,
		"GolemiteGolem": 1.2,
		"LavaHound": 2.0,
	}
	
	return deathRadii[troop.Name]
}


