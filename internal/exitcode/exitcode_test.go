package exitcode

import "testing"

func TestSuccessConstant(t *testing.T) {
	if Success != 0 {
		t.Errorf("Success = %d, want 0", Success)
	}
}

func TestFailureConstant(t *testing.T) {
	if Failure != 1 {
		t.Errorf("Failure = %d, want 1", Failure)
	}
}

func TestExitVariableOverride(t *testing.T) {
	original := Exit
	defer func() { Exit = original }()

	var captured int
	Exit = func(code int) { captured = code }

	Exit(Failure)

	if captured != 1 {
		t.Errorf("Exit(Failure) captured = %d, want 1", captured)
	}
}

func TestExitWithSuccess(t *testing.T) {
	original := Exit
	defer func() { Exit = original }()

	var captured int
	Exit = func(code int) { captured = code }

	Exit(Success)

	if captured != 0 {
		t.Errorf("Exit(Success) captured = %d, want 0", captured)
	}
}
