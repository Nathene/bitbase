package states

import (
	"github.com/hajimehoshi/ebiten/v2"
)

// GameState represents a discrete state in the game (menu, gameplay, etc.)
type GameState interface {
	// Initialize is called when a state is first pushed onto the stack
	Initialize() error

	// Enter is called whenever this state becomes the active state
	Enter() error

	// Exit is called when this state is no longer the active state
	Exit() error

	// Update handles state logic, animations, etc.
	Update() error

	// Draw renders the state to the screen
	Draw(screen *ebiten.Image)

	// HandleInput processes all input for this state
	HandleInput() error

	// GetStateID returns a unique identifier for this state
	GetStateID() string
}

// StateManager controls the flow between different game states
type StateManager struct {
	states       []GameState   // Stack of game states
	stateChanges []StateChange // Queue of pending state changes
	isProcessing bool          // Flag to prevent recursive state changes
}

// StateChangeType represents different ways to change states
type StateChangeType int

const (
	Push    StateChangeType = iota // Add a new state on top
	Pop                            // Remove the top state
	Replace                        // Replace the top state
	Clear                          // Clear all states
)

// StateChange represents a pending change to the state stack
type StateChange struct {
	changeType StateChangeType
	state      GameState
}

// NewStateManager creates a new state manager
func NewStateManager() *StateManager {
	return &StateManager{
		states:       make([]GameState, 0),
		stateChanges: make([]StateChange, 0),
		isProcessing: false,
	}
}

// Update processes any pending state changes and updates the active state
func (sm *StateManager) Update() error {
	// Process any pending state changes
	if len(sm.stateChanges) > 0 && !sm.isProcessing {
		sm.isProcessing = true
		change := sm.stateChanges[0]
		sm.stateChanges = sm.stateChanges[1:]

		switch change.changeType {
		case Push:
			// Exit current state if it exists
			if len(sm.states) > 0 {
				sm.states[len(sm.states)-1].Exit()
			}

			// Initialize and enter new state
			change.state.Initialize()
			change.state.Enter()

			// Add to stack
			sm.states = append(sm.states, change.state)

		case Pop:
			// Exit and remove current state
			if len(sm.states) > 0 {
				sm.states[len(sm.states)-1].Exit()
				sm.states = sm.states[:len(sm.states)-1]

				// Enter the now-active state
				if len(sm.states) > 0 {
					sm.states[len(sm.states)-1].Enter()
				}
			}

		case Replace:
			// Exit and remove current state
			if len(sm.states) > 0 {
				sm.states[len(sm.states)-1].Exit()
				sm.states = sm.states[:len(sm.states)-1]
			}

			// Initialize and enter new state
			change.state.Initialize()
			change.state.Enter()

			// Add to stack
			sm.states = append(sm.states, change.state)

		case Clear:
			// Exit all states
			for i := len(sm.states) - 1; i >= 0; i-- {
				sm.states[i].Exit()
			}
			sm.states = make([]GameState, 0)

			// Initialize and enter new state if provided
			if change.state != nil {
				change.state.Initialize()
				change.state.Enter()
				sm.states = append(sm.states, change.state)
			}
		}

		sm.isProcessing = false
	}

	// Update active state
	if len(sm.states) > 0 {
		return sm.states[len(sm.states)-1].Update()
	}

	return nil
}

// RequestStateChange queues a state change
func (sm *StateManager) RequestStateChange(change StateChange) {
	sm.stateChanges = append(sm.stateChanges, change)
}

// PushState adds a new state to the top of the stack
func (sm *StateManager) PushState(state GameState) {
	sm.RequestStateChange(StateChange{
		changeType: Push,
		state:      state,
	})
}

// PopState removes the top state from the stack
func (sm *StateManager) PopState() {
	sm.RequestStateChange(StateChange{
		changeType: Pop,
	})
}

// ReplaceState replaces the top state with a new state
func (sm *StateManager) ReplaceState(state GameState) {
	sm.RequestStateChange(StateChange{
		changeType: Replace,
		state:      state,
	})
}

// ClearStates removes all states from the stack
func (sm *StateManager) ClearStates() {
	sm.RequestStateChange(StateChange{
		changeType: Clear,
	})
}

// Draw draws the active state
func (sm *StateManager) Draw(screen *ebiten.Image) {
	if len(sm.states) > 0 {
		sm.states[len(sm.states)-1].Draw(screen)
	}
}

// GetActiveState returns the currently active state
func (sm *StateManager) GetActiveState() GameState {
	if len(sm.states) > 0 {
		return sm.states[len(sm.states)-1]
	}
	return nil
}
