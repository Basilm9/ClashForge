// troop_wrapper.go
package clashgame

import (
	"encoding/csv"
	"fmt"
	"image/color"
	"io"
	"os"
	"strconv"
	"strings"
)

// TroopTemplate contains the extended properties from the CSV file
// Only including the most relevant fields from the 336 columns
type TroopTemplate struct {
	// Basic identifiers
	Name            string
	Rarity          string
	Tribe           string
	
	// Core stats from CSV
	SightRange      float64
	DeployTime      float64
	ChargeRange     float64
	Speed           float64
	Hitpoints       int
	HitSpeed        float64
	LoadTime        float64
	Damage          int
	DamageSpecial   int
	Range           float64
	MinimumRange    float64
	
	// Attack properties
	AttacksGround   bool
	AttacksAir      bool
	AreaDamageRadius float64
	TargetOnlyBuildings bool
	TargetOnlyTroops bool
	
	// Special abilities
	DeathDamage     int
	DeathDamageRadius float64
	LifeTime        float64
	SpawnInterval   float64
	SpawnNumber     int
	
	// Visual properties
	Scale           float64
	CollisionRadius float64
	FlyingHeight    float64
	
	// Additional properties can be added as needed
	Projectile		ProjectileTemplate
}

// ExtendedTroop wraps the base Troop struct with extended properties
type ExtendedTroop struct {
	Troop          // Embed the base Troop struct
	Template    *TroopTemplate
}

// TroopTemplateMap is a map of troop names to their templates
var TroopTemplateMap map[string]*TroopTemplate

// LoadTroopTemplates loads troop templates from a CSV file
func LoadTroopTemplates(filepath string) error {
	// Initialize the map
	TroopTemplateMap = make(map[string]*TroopTemplate)
	
	// Try to load from CSV
	err := loadTroopTemplatesFromCSV(filepath)
	if err != nil {
		fmt.Printf("Warning: Failed to load troop templates from CSV: %v\n", err)
		fmt.Println("Falling back to default troop templates")
		InitializeWithDefaultTroops()
		return nil
	}
	
	// If the CSV was loaded but no valid troops were found, use defaults
	if len(TroopTemplateMap) == 0 {
		fmt.Println("Warning: No valid troop templates found in CSV")
		fmt.Println("Falling back to default troop templates")
		InitializeWithDefaultTroops()
	}
	
	return nil
}

// NewExtendedTroop creates a new ExtendedTroop with a template
func NewExtendedTroop(x, y float64, troopName string, team int, grid *GridSystem) (*ExtendedTroop, error) {
	// Get the template
	template, exists := TroopTemplateMap[troopName]
	if !exists {
		return nil, fmt.Errorf("troop template not found: %s", troopName)
	}
	
	// Calculate appropriate size based on collision radius or scale
	sizeInCells := 0.8 // Default size
	if template.CollisionRadius > 0 {
		// Use collision radius if available
		sizeInCells = template.CollisionRadius / 500
		if sizeInCells < 0.5 {
			sizeInCells = 1 // Minimum size
		} else if sizeInCells > 2 {
			sizeInCells = 3 // Maximum size
		}
	} else if template.Scale > 0 {
		// Otherwise use scale
		sizeInCells = 0.8 * template.Scale
	}
	
	// Get color based on team and rarity
	troopColor := GetTroopColorByTeam(team, template.Rarity)
	
	// Create the base troop using the template values
	troop := NewTroop(
		x,
		y,
		template.Hitpoints,
		template.Damage,
		template.Speed,
		template.Range,
		template.SightRange, // Using SightRange as AggroDistance
		troopColor,
		grid,
		sizeInCells,
	)

	InitTroopMovement(&troop)
	
	// Set name from template
	troop.Name = template.Name
	
	// Set team
	troop.Team = team
	
	// Set attack delay based on hit speed
	if template.HitSpeed > 0 {
		troop.AttackDelay = int(60 / template.HitSpeed) // Assuming 60 ticks per second
	}
	
	// Create the extended troop
	extendedTroop := &ExtendedTroop{
		Troop:    troop,
		Template: template,
	}
	
	return extendedTroop, nil
}

// Methods to expose and use template properties

// GetAttackDelay returns the attack delay based on HitSpeed
func (et *ExtendedTroop) GetAttackDelay() int {
	if et.Template.HitSpeed <= 0 {
		return 20 // Default attack delay
	}
	// Convert hit speed to game ticks (assuming 60 ticks per second)
	return int(60 / et.Template.HitSpeed)
}

// SpawnExtendedTroop adds an extended troop to the game
func SpawnExtendedTroop(troopName string, x, y float64, team int, g *Game) error {
	extendedTroop, err := NewExtendedTroop(x, y, troopName, team, g.Grid)
	if err != nil {
		return err
	}
	
	// Add to game's troop list
	g.Troops = append(g.Troops, extendedTroop.Troop)
	
	return nil
}

// InitializeTroopSystem sets up the troop system by loading templates
func InitializeTroopSystem(csvPath string) error {
	return LoadTroopTemplates(csvPath)
}

// Default troops when CSV loading fails
var defaultTroops = map[string]*TroopTemplate{
	"Knight": {
		Name:           "Knight",
		Rarity:         "Common",
		Tribe:          "Ground",
		SightRange:     5.0,
		DeployTime:     1.0,
		Speed:          0.15,
		Hitpoints:      150,
		HitSpeed:       1.2,
		Damage:         75,
		Range:          0.8,
		AttacksGround:  true,
		AttacksAir:     false,
		Scale:          1.0,
		CollisionRadius: 0.7,
	},
	"Archer": {
		Name:           "Archer",
		Rarity:         "Common",
		Tribe:          "Ground",
		SightRange:     5.5,
		DeployTime:     0.5,
		Speed:          0.2,
		Hitpoints:      80,
		HitSpeed:       0.7,
		Damage:         40,
		Range:          4.0,
		AttacksGround:  true,
		AttacksAir:     true,
		Scale:          0.9,
		CollisionRadius: 0.6,
	},
	"Skeleton": {
		Name:           "Skeleton",
		Rarity:         "Common",
		Tribe:          "Ground",
		SightRange:     4.0,
		DeployTime:     0.0,
		Speed:          0.25,
		Hitpoints:      40,
		HitSpeed:       0.5,
		Damage:         25,
		Range:          0.5,
		AttacksGround:  true,
		AttacksAir:     false,
		Scale:          0.7,
		CollisionRadius: 0.4,
	},
	"Giant": {
		Name:           "Giant",
		Rarity:         "Rare",
		Tribe:          "Ground",
		SightRange:     5.0,
		DeployTime:     1.0,
		Speed:          0.1,
		Hitpoints:      800,
		HitSpeed:       1.5,
		Damage:         100,
		Range:          0.8,
		AttacksGround:  true,
		AttacksAir:     false,
		TargetOnlyBuildings: true,
		Scale:          1.5,
		CollisionRadius: 1.0,
	},
	"BabyDragon": {
		Name:           "BabyDragon",
		Rarity:         "Epic",
		Tribe:          "Air",
		SightRange:     6.0,
		DeployTime:     1.0,
		Speed:          0.2,
		Hitpoints:      200,
		HitSpeed:       1.3,
		Damage:         60,
		Range:          3.0,
		AttacksGround:  true,
		AttacksAir:     true,
		AreaDamageRadius: 1.5,
		Scale:          1.2,
		CollisionRadius: 0.8,
		FlyingHeight:   1.0,
	},
}

// InitializeWithDefaultTroops sets up default troop templates when CSV loading fails
func InitializeWithDefaultTroops() {
	// Initialize the map if it doesn't exist
	if TroopTemplateMap == nil {
		TroopTemplateMap = make(map[string]*TroopTemplate)
	}
	
	// Add default troops
	for name, template := range defaultTroops {
		TroopTemplateMap[name] = template
	}
	
	fmt.Println("Initialized with default troop templates")
}

// Convenience method to create team-specific colors for troops
func GetTroopColorByTeam(team int, rarity string) color.RGBA {
	// Base colors by team
	var baseRed, baseGreen, baseBlue uint8
	
	if team == 0 {
		// Team 0 (red team)
		baseRed = 220
		baseGreen = 40
		baseBlue = 40
	} else {
		// Team 1 (blue team)
		baseRed = 40
		baseGreen = 40
		baseBlue = 220
	}
	
	// Modify color based on rarity
	switch rarity {
	case "Common":
		// No change for common
	case "Rare":
		// Add green tint for rare
		if team == 0 {
			baseGreen = 120
		} else {
			baseGreen = 180
		}
	case "Epic":
		// Add purple tint for epic
		if team == 0 {
			baseBlue = 150
		} else {
			baseRed = 150
		}
	case "Legendary":
		// Add gold tint for legendary
		baseRed = 230
		baseGreen = 180
		if team == 0 {
			baseBlue = 40
		} else {
			baseBlue = 160
		}
	}
	
	return color.RGBA{baseRed, baseGreen, baseBlue, 255}
}

// Helper function to get troop display name
func GetTroopDisplayName(templateName string) string {
	template, exists := TroopTemplateMap[templateName]
	if !exists {
		return templateName
	}
	
	return template.Name
}

func loadTroopTemplatesFromCSV(filepath string) error {
	// Open the CSV file
	file, err := os.Open(filepath)
	if err != nil {
		return fmt.Errorf("failed to open CSV file: %v", err)
	}
	defer file.Close()

	// Parse the CSV file
	reader := csv.NewReader(file)
	
	// Read the header
	header, err := reader.Read()
	if err != nil {
		return fmt.Errorf("failed to read CSV header: %v", err)
	}
	
	// Read the type information (second row, defines data types)
	_, err = reader.Read()
	if err != nil {
		return fmt.Errorf("failed to read CSV types: %v", err)
	}
	
	// Read the type descriptor row (third row)
	_, err = reader.Read()
	if err != nil {
		return fmt.Errorf("failed to read CSV type descriptors: %v", err)
	}
	
	// Create a map to store column indices for easier access
	columnMap := make(map[string]int)
	for i, col := range header {
		columnMap[col] = i
	}
	
	// Initialize the template map if it doesn't exist yet
	if TroopTemplateMap == nil {
		TroopTemplateMap = make(map[string]*TroopTemplate)
	}
	
	// Read the remaining rows (actual troop data)
	rowCount := 0
	loadedCount := 0
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Printf("Warning: failed to read CSV row %d: %v\n", rowCount+4, err)
			continue
		}
		rowCount++
		
		// Skip empty rows or rows with "NOTINUSE" in the name
		if len(record) == 0 || record[0] == "" || strings.Contains(record[0], "NOTINUSE") {
			continue
		}
		
		// Extract troop data from the record
		troopName := record[columnMap["Name"]]
		if troopName == "" {
			continue // Skip rows without a name
		}

		projectileName := getStringValue(record, columnMap, "Projectile")
		
		// Create a new template
		template := &TroopTemplate{
			Name:               troopName,
			Rarity:             getStringValue(record, columnMap, "Rarity"),
			Tribe:              getStringValue(record, columnMap, "Tribe"),
			SightRange:         getFloatValue(record, columnMap, "SightRange") / 1000, 
			DeployTime:         getFloatValue(record, columnMap, "DeployTime") / 1000, 
			ChargeRange:        getFloatValue(record, columnMap, "ChargeRange") / 1000, 
			Speed:              getFloatValue(record, columnMap, "Speed") / 400, 
			Hitpoints:          getIntValue(record, columnMap, "Hitpoints"),
			HitSpeed:           getFloatValue(record, columnMap, "HitSpeed") / 1000, 
			LoadTime:           getFloatValue(record, columnMap, "LoadTime") / 1000, 
			Damage:             getIntValue(record, columnMap, "Damage"),
			DamageSpecial:      getIntValue(record, columnMap, "DamageSpecial"),
			Range:              getFloatValue(record, columnMap, "Range") / 1500, 
			MinimumRange:       getFloatValue(record, columnMap, "MinimumRange") / 1500, 
			AttacksGround:      getBoolValue(record, columnMap, "AttacksGround"),
			AttacksAir:         getBoolValue(record, columnMap, "AttacksAir"),
			AreaDamageRadius:   getFloatValue(record, columnMap, "AreaDamageRadius") / 1000, 
			TargetOnlyBuildings: getBoolValue(record, columnMap, "TargetOnlyBuildings"),
			TargetOnlyTroops:   getBoolValue(record, columnMap, "TargetOnlyTroops"),
			DeathDamage:        getIntValue(record, columnMap, "DeathDamage"),
			DeathDamageRadius:  getFloatValue(record, columnMap, "DeathDamageRadius") / 1000, 
			LifeTime:           getFloatValue(record, columnMap, "LifeTime") / 1000, 
			SpawnInterval:      getFloatValue(record, columnMap, "SpawnInterval") / 1000, 
			SpawnNumber:        getIntValue(record, columnMap, "SpawnNumber"),
			Scale:              getFloatValue(record, columnMap, "Scale") / 100, 
			CollisionRadius:    getFloatValue(record, columnMap, "CollisionRadius") / 100, 
			FlyingHeight:       getFloatValue(record, columnMap, "FlyingHeight") / 100,
			// Projectile field will be set below
		}
		
		linkTroopToProjectile(template, projectileName)
		
		// Add the template to the map
		TroopTemplateMap[troopName] = template
		loadedCount++
		
	}
	
	// Check if we loaded any templates
	if loadedCount == 0 {
		return fmt.Errorf("no valid troop templates found in CSV")
	}
	
	return nil
}

// Helper functions to extract typed values from CSV records
func getStringValue(record []string, columnMap map[string]int, column string) string {
	if index, exists := columnMap[column]; exists && index < len(record) {
		return record[index]
	}
	return ""
}

func getIntValue(record []string, columnMap map[string]int, column string) int {
	if index, exists := columnMap[column]; exists && index < len(record) {
		if record[index] != "" {
			if val, err := strconv.Atoi(record[index]); err == nil {
				return val
			}
		}
	}
	return 0
}

func getFloatValue(record []string, columnMap map[string]int, column string) float64 {
	if index, exists := columnMap[column]; exists && index < len(record) {
		if record[index] != "" {
			if val, err := strconv.ParseFloat(record[index], 64); err == nil {
				return val
			}
		}
	}
	return 0
}

func getBoolValue(record []string, columnMap map[string]int, column string) bool {
	if index, exists := columnMap[column]; exists && index < len(record) {
		return strings.ToLower(record[index]) == "true"
	}
	return false
}