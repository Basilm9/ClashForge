// projectile.go
package clashgame

import (
	"encoding/csv"
	"fmt"
	"image/color"
	"io"
	"math"
	"os"
	"strings"
)

func clearAttackingStateOfSource(game *Game, sourceID int) {
    if sourceID <= 0 {
        return // Invalid source ID
    }
    
    // Find the source troop and clear its attacking state
    for i := range game.Troops {
        if game.Troops[i].ID == sourceID && game.Troops[i].Active {
            game.Troops[i].IsAttacking = false
            fmt.Printf("Cleared attacking state for Troop ID=%d after its projectile defeated a target\n", sourceID)
            break
        }
    }
}

// ProjectileTemplate holds properties from the projectile CSV
type ProjectileTemplate struct {
	Name                 string
	Rarity               string
	Speed                float64
	Scale                float64
	Homing               bool
	HomingTime           float64
	HomingMinDistance    float64
	Damage               int
	CrownTowerDamagePercent float64
	Pushback             float64
	PushbackAll          bool
	Radius               float64
	AoeToAir             bool
	AoeToGround          bool
	OnlyEnemies          bool
	MaximumTargets       int
	ProjectileRadius     float64
	TrailEffect          string
	ConstantHeight       bool
}

// Projectile represents a projectile in the game
type Projectile struct {
	Name		   string
	Position       Position    // Current position
	StartPosition  Position    // Starting position
	TargetPosition Position    // Target position (for non-homing projectiles)
	TargetEntity   interface{} // Target troop or building (for homing projectiles)
	Direction      Position    // Normalized direction vector
	Speed          float64     // Movement speed (grid cells per tick)
	Damage         int         // Damage to deal on hit
	Radius         float64     // Explosion radius (for area damage)
	Color          color.RGBA  // Visual color
	Size           float64     // Visual size
	Active         bool        // Whether the projectile is active
	Team           int         // Team that fired this projectile
	IsHoming       bool        // Whether the projectile homes in on target
	HomingTime     float64     // Time the projectile has been homing (in ticks)
	MaxHomingTime  float64     // Maximum homing time (in ticks)
	AoeToAir       bool        // Whether area damage affects air units
	AoeToGround    bool        // Whether area damage affects ground units
	LifeTime       int         // How long the projectile has existed
	MaxLifeTime    int         // Maximum lifetime of projectile (prevents infinite projectiles)
	SourceID       int         // ID of the entity that fired this projectile (to prevent self-hits)
	Template       *ProjectileTemplate // Reference to the template
}

// ProjectileTemplateMap is a map of projectile names to their templates
var ProjectileTemplateMap map[string]*ProjectileTemplate
// Modify the InitializeProjectileSystem function in projectile.go to ensure "normal" is properly added
func InitializeProjectileSystem() {
    // Initialize the map
    ProjectileTemplateMap = make(map[string]*ProjectileTemplate)
    
    // Add "normal" projectile for buildings
    ProjectileTemplateMap["normal"] = &ProjectileTemplate{
        Name:              "normal",
        Rarity:            "Common",
        Speed:             0.7,
        Scale:             1.0,
        Homing:            false,
        Damage:            45,
        Radius:            0,
        AoeToGround:       true,
        AoeToAir:          true,
        OnlyEnemies:       true,
        ProjectileRadius:  0.2,
    }
    
    // Add "ArcherArrow" specific projectile
    ProjectileTemplateMap["ArcherArrow"] = &ProjectileTemplate{
        Name:              "ArcherArrow",
        Rarity:            "Common",
        Speed:             1.0,  // Faster than regular arrow
        Scale:             0.9,
        Homing:            false,
        Damage:            40,    // Base damage
        Radius:            0,
        AoeToGround:       true,
        AoeToAir:          true,
        OnlyEnemies:       true,
        ProjectileRadius:  0.15,
    }

}

// Updated CreateProjectile function to consider template damage
func CreateProjectile(templateName string, source Position, target Position, damage int, team int, sourceID int) *Projectile {
    // Skip if template name is empty or "none"
    if templateName == "" || templateName == "none" {
        return nil
    }
    
    // Ensure ProjectileTemplateMap is initialized
    if ProjectileTemplateMap == nil {
        InitializeProjectileSystem()
    }
    
    // Get the template
    template, exists := ProjectileTemplateMap[templateName]
    if !exists {
        // Log warning
        fmt.Printf("Warning: Projectile template '%s' not found, using default\n", templateName)
        
        // Create a default template for missing projectiles
        template = &ProjectileTemplate{
            Name:              templateName,
            Rarity:            "Common",
            Speed:             0.7,
            Scale:             1.0,
            Homing:            false,
            Damage:            45,
            Radius:            0,
            AoeToGround:       true,
            AoeToAir:          true,
            OnlyEnemies:       true,
            ProjectileRadius:  0.2,
        }
        
        // Add it to the map for future use
        ProjectileTemplateMap[templateName] = template
    }
    
    // Calculate direction vector
    dx := target.X - source.X
    dy := target.Y - source.Y
    distance := math.Sqrt(dx*dx + dy*dy)
    
    // Normalize direction
    var direction Position
    if distance > 0 {
        direction = Position{
            X: dx / distance,
            Y: dy / distance,
        }
    } else {
        // Default direction if source and target are the same
        direction = Position{X: 0, Y: 1}
    }
    
    // Ensure the speed is reasonable (not too fast or too slow)
    speed := template.Speed
    if speed < 0.1 {
        speed = 0.1  // Minimum speed
    } else if speed > 2.0 {
        speed = 2.0  // Maximum speed
    }
    
    // Ensure size is reasonable
    size := template.ProjectileRadius * 20 // Convert to pixels
    if size < 4.0 {
        size = 4.0  // Minimum size for visibility
    }
    
    // Determine color based on team and rarity
    var projectileColor color.RGBA
    switch template.Rarity {
    case "Common":
        if team == 0 {
            projectileColor = color.RGBA{255, 100, 100, 255} // Light red
        } else {
            projectileColor = color.RGBA{100, 100, 255, 255} // Light blue
        }
    case "Rare":
        if team == 0 {
            projectileColor = color.RGBA{255, 160, 80, 255} // Orange
        } else {
            projectileColor = color.RGBA{80, 160, 255, 255} // Light blue
        }
    case "Epic":
        if team == 0 {
            projectileColor = color.RGBA{255, 80, 255, 255} // Pink
        } else {
            projectileColor = color.RGBA{80, 80, 255, 255} // Blue
        }
    default:
        if team == 0 {
            projectileColor = color.RGBA{255, 100, 100, 255} // Light red
        } else {
            projectileColor = color.RGBA{100, 100, 255, 255} // Light blue
        }
    }
    
    // IMPORTANT FIX: Use damage from the template if available
    // Otherwise fall back to the provided damage parameter
    projectileDamage := damage
    if template.Damage > 0 {
        projectileDamage = template.Damage
        fmt.Printf("Using template damage %d for projectile %s (instead of troop damage %d)\n", 
                 template.Damage, templateName, damage)
    }
    
    // Create the projectile
    projectile := &Projectile{
        Name:           templateName,
        Position:       source,
        StartPosition:  source,
        TargetPosition: target,
        Direction:      direction,
        Speed:          speed,
        Damage:         projectileDamage,  // Now using template damage
        Radius:         template.Radius,
        Color:          projectileColor,
        Size:           size,
        Active:         true,
        Team:           team,
        IsHoming:       template.Homing,
        MaxHomingTime:  template.HomingTime,
        AoeToAir:       template.AoeToAir,
        AoeToGround:    template.AoeToGround,
        LifeTime:       0,
        MaxLifeTime:    300, // 5 seconds at 60 ticks per second
        SourceID:       sourceID,
        Template:       template,
    }
    
    // Log projectile creation
    fmt.Printf("Created projectile: Name=%s, Damage=%d, Speed=%.2f, Size=%.2f\n", 
             templateName, projectileDamage, speed, size)
    
    return projectile
}

func (p *Projectile) Update(game *Game) {}

// HandleImpact deals damage when a projectile hits something
func (p *Projectile) HandleImpact(game *Game, impactPos Position) {}

// Add this to UpdateGame or gameloop.go to update projectiles
func UpdateProjectiles(game *Game) {
    activeProjectiles := []Projectile{}
    
    // Update each projectile
    for i := range game.Projectiles {
        projectile := &game.Projectiles[i]
        if projectile.Active {
            projectile.Update(game)
            
            if projectile.Active {
                activeProjectiles = append(activeProjectiles, *projectile)
            }
        }
    }
    
    // Replace projectiles list with active ones
    game.Projectiles = activeProjectiles
}


// LoadProjectileTemplates loads projectile templates from a CSV file
func LoadProjectileTemplates(filepath string) error {
    // Initialize the map
    ProjectileTemplateMap = make(map[string]*ProjectileTemplate)
    
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
    
    // Create a map to store column indices for easier access
    columnMap := make(map[string]int)
    for i, col := range header {
        columnMap[col] = i
    }
    
    // Skip type rows (if needed, similar to troops.csv)
    _, err = reader.Read() // Skip type information row
    if err != nil {
        return fmt.Errorf("failed to read CSV types: %v", err)
    }
    
    _, err = reader.Read() // Skip type descriptor row
    if err != nil {
        return fmt.Errorf("failed to read CSV type descriptors: %v", err)
    }
    
    // Read all projectile data rows
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
        
        // Extract projectile data from the record
        projectileName := record[columnMap["Name"]]
        if projectileName == "" {
            continue // Skip rows without a name
        }

        // Create a new template
        template := &ProjectileTemplate{
            Name:                 projectileName,
            Rarity:               getStringValue(record, columnMap, "Rarity"),
            Speed:                getFloatValue(record, columnMap, "Speed") / 200, // Adjust divisor as needed
            Scale:                getFloatValue(record, columnMap, "Scale") / 100,
            Homing:               getBoolValue(record, columnMap, "Homing"),
            HomingTime:           getFloatValue(record, columnMap, "HomingTime") / 10, // Adjust divisor as needed
            HomingMinDistance:    getFloatValue(record, columnMap, "HomingMinDistance"),
            Damage:               getIntValue(record, columnMap, "Damage"),
            CrownTowerDamagePercent: getFloatValue(record, columnMap, "CrownTowerDamagePercent"),
            Pushback:             getFloatValue(record, columnMap, "Pushback"),
            PushbackAll:          getBoolValue(record, columnMap, "PushbackAll"),
            Radius:               getFloatValue(record, columnMap, "Radius") / 1000,
            AoeToAir:             getBoolValue(record, columnMap, "AoeToAir"),
            AoeToGround:          getBoolValue(record, columnMap, "AoeToGround"),
            OnlyEnemies:          getBoolValue(record, columnMap, "OnlyEnemies"),
            MaximumTargets:       getIntValue(record, columnMap, "MaximumTargets"),
            ProjectileRadius:     getFloatValue(record, columnMap, "ProjectileRadius") / 1000,
            TrailEffect:          getStringValue(record, columnMap, "TrailEffect"),
            ConstantHeight:       getBoolValue(record, columnMap, "ConstantHeight"),
        }
        
        // Add the template to the map
        ProjectileTemplateMap[projectileName] = template
        loadedCount++
        
    }
    
    // Check if we loaded any templates
    if loadedCount == 0 {
        return fmt.Errorf("no valid projectile templates found in CSV")
    }
    
    return nil
}

// Replace the linkTroopToProjectile function in projectile.go with this improved version
func linkTroopToProjectile(troopTemplate *TroopTemplate, projectileName string) {
    // Determine if this troop should be melee based on range
    isMeleeTroop := troopTemplate.Range <= 1.0
    
    // If no projectile specified (melee troop) OR this is a melee troop, leave it without a projectile
    if projectileName == "" || isMeleeTroop {
        // For melee troops, ensure the Projectile.Name is empty to signal no projectile
        troopTemplate.Projectile.Name = ""
        return
    }
    
    // Check if projectile template map is initialized
    if ProjectileTemplateMap == nil {
        // Projectile system not initialized yet, call the initialization
        InitializeProjectileSystem()
    }
    
    // Look up the projectile template
    if projectileTemplate, exists := ProjectileTemplateMap[projectileName]; exists {
        // Deep copy the projectile template to avoid sharing references
        troopTemplate.Projectile = *projectileTemplate
    } else {
        // Only assign default arrow if this is definitely a ranged troop
        if !isMeleeTroop {
            if defTemplate, exists := ProjectileTemplateMap["arrow"]; exists {
                troopTemplate.Projectile = *defTemplate
                // Update name to match the requested projectile
                troopTemplate.Projectile.Name = projectileName
                fmt.Printf("Using default arrow for ranged troop %s (requested %s)\n", troopTemplate.Name, projectileName)
            }else {
            // This is a melee troop, don't assign a projectile
            troopTemplate.Projectile.Name = ""
            fmt.Printf("Melee troop %s will not use projectiles\n", troopTemplate.Name)
        }}
    }
}

// This function should be added to troop_wrapper.go to handle projectile firing
func (et *ExtendedTroop) FireProjectile(game *Game, target Position) {
	// Don't fire a projectile if troop doesn't have one
	if et.Template.Projectile.Name == "" {
		return
	}
	
	// Create a new projectile based on the troop's template
	projectile := CreateProjectile(
		et.Template.Projectile.Name,
		et.Position,
		target,
		et.Damage,
		et.Team,
		et.ID,
	)
	
	// Add the projectile to the game
	if projectile != nil {
		game.Projectiles = append(game.Projectiles, *projectile)
	}
}