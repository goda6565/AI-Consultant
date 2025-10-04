package service

import (
	"testing"

	actionValue "github.com/goda6565/ai-consultant/backend/internal/domain/action/value"
)

func TestOrchestrator_selectLeastFrequentAction(t *testing.T) {
	orchestrator := &Orchestrator{}

	tests := []struct {
		name               string
		history            []actionValue.ActionType
		expectedCandidates []actionValue.ActionType
	}{
		{
			name:               "empty_history_returns_plan",
			history:            []actionValue.ActionType{},
			expectedCandidates: []actionValue.ActionType{actionValue.ActionTypePlan},
		},
		{
			name:    "single_action_returns_least_frequent",
			history: []actionValue.ActionType{actionValue.ActionTypeSearch},
			expectedCandidates: []actionValue.ActionType{
				actionValue.ActionTypePlan,
				actionValue.ActionTypeAnalyze,
				actionValue.ActionTypeWrite,
				actionValue.ActionTypeReview,
			},
		},
		{
			name: "all_actions_equal_frequency_returns_any_candidate",
			history: []actionValue.ActionType{
				actionValue.ActionTypePlan,
				actionValue.ActionTypeSearch,
				actionValue.ActionTypeAnalyze,
				actionValue.ActionTypeWrite,
				actionValue.ActionTypeReview,
			},
			expectedCandidates: []actionValue.ActionType{
				actionValue.ActionTypePlan,
				actionValue.ActionTypeSearch,
				actionValue.ActionTypeAnalyze,
				actionValue.ActionTypeWrite,
				actionValue.ActionTypeReview,
			},
		},
		{
			name: "plan_most_frequent_returns_least_frequent",
			history: []actionValue.ActionType{
				actionValue.ActionTypePlan,
				actionValue.ActionTypePlan,
				actionValue.ActionTypeSearch,
			},
			expectedCandidates: []actionValue.ActionType{
				actionValue.ActionTypeAnalyze,
				actionValue.ActionTypeWrite,
				actionValue.ActionTypeReview,
			},
		},
		{
			name: "search_most_frequent_returns_least_frequent",
			history: []actionValue.ActionType{
				actionValue.ActionTypeSearch,
				actionValue.ActionTypeSearch,
				actionValue.ActionTypePlan,
			},
			expectedCandidates: []actionValue.ActionType{
				actionValue.ActionTypeAnalyze,
				actionValue.ActionTypeWrite,
				actionValue.ActionTypeReview,
			},
		},
		{
			name: "analyze_most_frequent_returns_least_frequent",
			history: []actionValue.ActionType{
				actionValue.ActionTypeAnalyze,
				actionValue.ActionTypeAnalyze,
				actionValue.ActionTypePlan,
			},
			expectedCandidates: []actionValue.ActionType{
				actionValue.ActionTypeSearch,
				actionValue.ActionTypeWrite,
				actionValue.ActionTypeReview,
			},
		},
		{
			name: "write_most_frequent_returns_least_frequent",
			history: []actionValue.ActionType{
				actionValue.ActionTypeWrite,
				actionValue.ActionTypeWrite,
				actionValue.ActionTypePlan,
			},
			expectedCandidates: []actionValue.ActionType{
				actionValue.ActionTypeSearch,
				actionValue.ActionTypeAnalyze,
				actionValue.ActionTypeReview,
			},
		},
		{
			name: "review_most_frequent_returns_least_frequent",
			history: []actionValue.ActionType{
				actionValue.ActionTypeReview,
				actionValue.ActionTypeReview,
				actionValue.ActionTypePlan,
			},
			expectedCandidates: []actionValue.ActionType{
				actionValue.ActionTypeSearch,
				actionValue.ActionTypeAnalyze,
				actionValue.ActionTypeWrite,
			},
		},
		{
			name: "complex_scenario_with_multiple_frequencies",
			history: []actionValue.ActionType{
				actionValue.ActionTypePlan,
				actionValue.ActionTypePlan,
				actionValue.ActionTypePlan,
				actionValue.ActionTypeSearch,
				actionValue.ActionTypeSearch,
				actionValue.ActionTypeAnalyze,
				actionValue.ActionTypeWrite,
				actionValue.ActionTypeReview,
			},
			expectedCandidates: []actionValue.ActionType{
				actionValue.ActionTypeAnalyze,
				actionValue.ActionTypeWrite,
				actionValue.ActionTypeReview,
			},
		},
		{
			name: "recent_10_actions_only_considered",
			history: []actionValue.ActionType{
				// First 5 actions (older than recent 10)
				actionValue.ActionTypePlan,
				actionValue.ActionTypePlan,
				actionValue.ActionTypePlan,
				actionValue.ActionTypePlan,
				actionValue.ActionTypePlan,
				// Recent 10 actions: Plan=5, Search=2, Analyze=1, Write=1, Review=1
				actionValue.ActionTypeSearch,
				actionValue.ActionTypeSearch,
				actionValue.ActionTypeAnalyze,
				actionValue.ActionTypeWrite,
				actionValue.ActionTypeReview,
				actionValue.ActionTypePlan,
				actionValue.ActionTypePlan,
				actionValue.ActionTypePlan,
				actionValue.ActionTypePlan,
				actionValue.ActionTypePlan,
			},
			expectedCandidates: []actionValue.ActionType{
				actionValue.ActionTypeAnalyze,
				actionValue.ActionTypeWrite,
				actionValue.ActionTypeReview,
			},
		},
		{
			name: "tie_breaking_returns_any_least_frequent",
			history: []actionValue.ActionType{
				actionValue.ActionTypeSearch,
				actionValue.ActionTypeAnalyze,
				actionValue.ActionTypeWrite,
				actionValue.ActionTypeReview,
			},
			expectedCandidates: []actionValue.ActionType{
				actionValue.ActionTypePlan,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := orchestrator.selectLeastFrequentAction(tt.history)
			found := false
			for _, candidate := range tt.expectedCandidates {
				if result == candidate {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("selectLeastFrequentAction() = %v, want one of %v", result, tt.expectedCandidates)
			}
		})
	}
}

func TestOrchestrator_selectLeastFrequentAction_EdgeCases(t *testing.T) {
	orchestrator := &Orchestrator{}

	t.Run("history_longer_than_10_considers_only_recent_10", func(t *testing.T) {
		// Create a history with 15 actions where the first 5 are all 'plan'
		// and the recent 10 have different frequencies
		history := make([]actionValue.ActionType, 15)
		for i := 0; i < 5; i++ {
			history[i] = actionValue.ActionTypePlan
		}
		// Recent 10: search appears 3 times, others appear 1-2 times
		history[5] = actionValue.ActionTypeSearch
		history[6] = actionValue.ActionTypeSearch
		history[7] = actionValue.ActionTypeSearch
		history[8] = actionValue.ActionTypeAnalyze
		history[9] = actionValue.ActionTypeWrite
		history[10] = actionValue.ActionTypeReview
		history[11] = actionValue.ActionTypePlan
		history[12] = actionValue.ActionTypeAnalyze
		history[13] = actionValue.ActionTypeWrite
		history[14] = actionValue.ActionTypeReview

		result := orchestrator.selectLeastFrequentAction(history)
		// Should return one of the least frequent in recent 10 (plan, write, review)
		expectedCandidates := []actionValue.ActionType{
			actionValue.ActionTypePlan,
			actionValue.ActionTypeWrite,
			actionValue.ActionTypeReview,
		}

		found := false
		for _, candidate := range expectedCandidates {
			if result == candidate {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected result to be one of %v, got %v", expectedCandidates, result)
		}
	})

	t.Run("all_actions_same_type_returns_different_action", func(t *testing.T) {
		history := []actionValue.ActionType{
			actionValue.ActionTypePlan,
			actionValue.ActionTypePlan,
			actionValue.ActionTypePlan,
		}
		result := orchestrator.selectLeastFrequentAction(history)
		if result == actionValue.ActionTypePlan {
			t.Errorf("Expected different action than plan, got %v", result)
		}
	})
}
