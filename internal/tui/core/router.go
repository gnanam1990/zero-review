package core

// Screen is the current TUI screen.
type Screen int

const (
	ScreenWelcome Screen = iota
	ScreenPRInput
	ScreenLoadingReview
	ScreenDashboard
	ScreenFindings
	ScreenFindingDetail
	ScreenDiff
	ScreenChat
	ScreenApproval
	ScreenReport
	ScreenSettings
)

// ScreenName returns the human-readable name for a screen.
func ScreenName(s Screen) string {
	switch s {
	case ScreenWelcome:
		return "Welcome"
	case ScreenPRInput:
		return "New Review"
	case ScreenLoadingReview:
		return "Reviewing"
	case ScreenDashboard:
		return "Overview"
	case ScreenFindings:
		return "Findings"
	case ScreenFindingDetail:
		return "Finding Detail"
	case ScreenDiff:
		return "Diff"
	case ScreenChat:
		return "Chat"
	case ScreenApproval:
		return "Approval"
	case ScreenReport:
		return "Report"
	case ScreenSettings:
		return "Settings"
	default:
		return "Unknown"
	}
}

// NavigableScreens returns the ordered list shown in the sidebar.
func NavigableScreens() []Screen {
	return []Screen{
		ScreenDashboard,
		ScreenFindings,
		ScreenDiff,
		ScreenChat,
		ScreenApproval,
		ScreenReport,
		ScreenSettings,
	}
}

// ScreenLabel returns the short label used in the sidebar.
func ScreenLabel(s Screen) string {
	switch s {
	case ScreenDashboard:
		return "Overview"
	case ScreenFindingDetail:
		return "Findings"
	default:
		return ScreenName(s)
	}
}
