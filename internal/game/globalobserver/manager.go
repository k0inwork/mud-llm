package globalobserver

import (
	"github.com/sirupsen/logrus"
	"mud/internal/dal"
	"mud/internal/game"
	"mud/internal/game/events"
	"mud/internal/models"
)

// GlobalObserverManager handles asynchronous propagation of events to global observers.
type GlobalObserverManager struct {
	eventBus         *events.EventBus
	perceptionFilter game.PerceptionFilterInterface
	ownerDAL         dal.OwnerDALInterface
	raceDAL          dal.RaceDALInterface
	professionDAL    dal.ProfessionDALInterface
}

// NewGlobalObserverManager creates a new GlobalObserverManager.
func NewGlobalObserverManager(
	eventBus *events.EventBus,
	perceptionFilter game.PerceptionFilterInterface,
	ownerDAL dal.OwnerDALInterface,
	raceDAL dal.RaceDALInterface,
	professionDAL dal.ProfessionDALInterface,
) *GlobalObserverManager {
	m := &GlobalObserverManager{
		eventBus:         eventBus,
		perceptionFilter: perceptionFilter,
		ownerDAL:         ownerDAL,
		raceDAL:          raceDAL,
		professionDAL:    professionDAL,
	}

	// Subscribe to ActionEvents for asynchronous processing
	actionEventChannel := make(chan interface{})
	eventBus.Subscribe(events.ActionEventType, actionEventChannel)
	go func() {
		for event := range actionEventChannel {
			if actionEvent, ok := event.(*events.ActionEvent); ok {
				m.HandleActionEvent(actionEvent)
			} else {
				logrus.Errorf("GlobalObserverManager: received unexpected event type on ActionEventType channel: %T", event)
			}
		}
	}()
	return m
}

func (gom *GlobalObserverManager) HandleActionEvent(event interface{}) {
	actionEvent, ok := event.(*events.ActionEvent)
	if !ok {
		logrus.Errorf("GlobalObserverManager: received unexpected event type: %T", event)
		return
	}

	// Find all Owners that are global observers (race-based, profession-based)
	owners, err := gom.ownerDAL.GetAllOwners()
	if err != nil {
		logrus.Errorf("GlobalObserverManager: failed to get all owners: %v", err)
		return
	}

	for _, owner := range owners {
		isGlobalObserver := false
		switch owner.MonitoredAspect {
		case "race":
			if actionEvent.Player != nil && actionEvent.Player.RaceID == owner.AssociatedID {
				isGlobalObserver = true
			}
		case "profession":
			if actionEvent.Player != nil && actionEvent.Player.ProfessionID == owner.AssociatedID {
				isGlobalObserver = true
			}
		}

		if isGlobalObserver {
			// Process asynchronously to avoid blocking the main event loop
			go gom.processGlobalObservation(actionEvent, owner)
		}
	}
}

func (gom *GlobalObserverManager) processGlobalObservation(event *events.ActionEvent, owner *models.Owner) {
	perceivedAction, err := gom.perceptionFilter.Filter(event, owner)
	if err != nil {
		logrus.Errorf("GlobalObserverManager: failed to filter perception for owner %s: %v", owner.ID, err)
		return
	}

	// Calculate significance score
	// For now, no additive bonuses or multipliers are implemented, so it's just BaseSignificance * Clarity
	significance := perceivedAction.BaseSignificance * perceivedAction.Clarity

	// Update owner's influence budget based on significance
	// This is a simplified model. More complex logic might involve decay, caps, etc.
	owner.CurrentInfluenceBudget = owner.CurrentInfluenceBudget + significance
	// Ensure budget doesn't exceed max
	if owner.CurrentInfluenceBudget > owner.MaxInfluenceBudget {
		owner.CurrentInfluenceBudget = owner.MaxInfluenceBudget
	}

	// Persist the updated owner (this should ideally be batched or handled by a separate persistence layer)
	err = gom.ownerDAL.UpdateOwner(owner)
	if err != nil {
		logrus.Errorf("GlobalObserverManager: failed to update owner %s budget: %v", owner.ID, err)
	}

	logrus.Infof("GlobalObserverManager: Owner %s (Monitors: %s %s) perceived action '%s' with significance %.2f. New budget: %.2f",
		owner.Name, owner.MonitoredAspect, owner.AssociatedID, perceivedAction.PerceivedActionType, significance, owner.CurrentInfluenceBudget)

	// TODO: Potentially trigger other global reactions here, e.g., new quests, global messages.
}