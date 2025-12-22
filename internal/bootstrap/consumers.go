package bootstrap

import (
	"fmt"

	"github.com/DmitriySama/teammate_search/config"
	consumer "github.com/DmitriySama/teammate_search/internal/consumer"
	studentsinfoprocessor "github.com/DmitriySama/teammate_search/internal/services/processors/students_info_processor"
)

func InitStudentInfoUpsertConsumer(cfg *config.Config, studentsInfoProcessor *studentsinfoprocessor.StudentsInfoProcessor) *studentsinfoupsertconsumer.StudentInfoUpsertConsumer {
	kafkaBrockers := []string{fmt.Sprintf("%v:%v", cfg.Kafka.Host, cfg.Kafka.Port)}
	return consumer.NewStudentInfoUpsertConsumer(studentsInfoProcessor, kafkaBrockers, cfg.Kafka.StudentInfoUpsertTopicName)
}