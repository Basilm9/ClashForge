// collision.go
package clashgame

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
