package constant

const (
	// Date formats
	DateFormat              = "2006-01-02"      // YYYY-MM-DD
	DateFormatSlash         = "2006/01/02"      // YYYY/MM/DD
	DateFormatHumanReadable = "January 2, 2006" // Month day, year

	// Time formats
	TimeFormat24Hour = "15:04:05"    // HH:MM:SS (24-hour format)
	TimeFormat12Hour = "03:04:05 PM" // HH:MM:SS AM/PM (12-hour format)

	// Combined date and time formats
	DateTimeFormat              = "2006-01-02 15:04:05"         // YYYY-MM-DD HH:MM:SS
	DateTimeFormat12Hour        = "2006-01-02 03:04:05 PM"      // YYYY-MM-DD HH:MM:SS AM/PM
	DateTimeFormatHumanReadable = "January 2, 2006 03:04:05 PM" // Month day, year HH:MM:SS AM/PM

	// Other formats
	ShortDateFormat  = "2006-01-02"              // YYYY-MM-DD
	LongDateFormat   = "Monday, January 2, 2006" // Day, Month day, year
	TimeFormatWithTZ = "15:04:05 MST"            // HH:MM:SS Timezone
)
