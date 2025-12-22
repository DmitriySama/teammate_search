package teammateSearchService
	
import (
	"time"

	"github.com/DmitriySama/teammate_search/internal/models"
)

func validateStatus(status string) bool {
	_, ok := models.ValidStatuses[status]
	return ok
}

func validateDeadline(deadline string) bool {
	if deadline == "" {
		return true
	}

	_, err := time.Parse("02.01.2006", deadline)
	return err == nil
}
