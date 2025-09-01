package constants

type View string

const (
	MenuView    View = "Menu"
	AddView     View = "Add"
	ListView    View = "View"
	JournalView View = "Journal"
	EditView    View = "Edit"
	ConfirmView View = "Confirm"

	TimeFormat = "2 Jan, 2006"
	Gap        = "\n\n"
	UnsavedId  = -1
)
