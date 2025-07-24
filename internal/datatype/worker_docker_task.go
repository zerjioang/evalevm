package datatype

import (
	"github.com/zerjioang/rooftop/v2/common/io/json"
	"log"
	"strings"
)

// DockerTask represents a Docker command to be executed
type DockerTask struct {
	id          TaskId
	opts        BytecodeScanOpts
	hexBytecode string
	cmd         []string
	finishCh    chan struct{}
	result      *Result
	parser      ResultParser
}

func NewDockerTask(id TaskId, opts BytecodeScanOpts, hexBytecode string, cmd []string) *DockerTask {
	return &DockerTask{
		id:          id,
		cmd:         cmd,
		opts:        opts,
		hexBytecode: hexBytecode,
		finishCh:    make(chan struct{}, 1),
	}
}

func (s *DockerTask) ID() TaskId {
	return s.id
}

func (s *DockerTask) Command() []string {
	return s.cmd
}

func (s *DockerTask) FinishChan() chan struct{} {
	return s.finishCh
}

func (s *DockerTask) WithResult(result *Result) {
	s.result = result
}

func (s *DockerTask) WithResultParser(parser ResultParser) {
	s.parser = parser
}

func (s *DockerTask) Failed() bool {
	return s.result.Error != nil || s.result.Measurements.ExitStatus != 0
}

func (s *DockerTask) Parse() {
	//fmt.Println(s.result.Error)
	//fmt.Println(s.result.Output)
	// if present, extract the measure.sh output
	idxStart := strings.LastIndex(s.result.Output, "evalevm_perf_metrics_start")
	idxEnd := strings.LastIndex(s.result.Output, "evalevm_perf_metrics_end")
	if idxStart != -1 && idxEnd != -1 && idxEnd > idxStart {
		measureJsonStr := s.result.Output[idxStart+len("evalevm_perf_metrics_start") : idxEnd]
		if err := json.Unmarshal([]byte(measureJsonStr), &s.result.Measurements); err != nil {
			log.Println("failed to parse measure.sh output: " + err.Error() + ". Probably caused due to missing output")
		}
	} else {
		log.Println("failed to parse measure.sh output: missing output")
		log.Println(s.result.Output)
		s.result.Measurements = Measurements{
			ExitStatus: 1,
		}
	}
	if s.parser != nil {
		s.parser.ParseOutput(s.result.Output)
	}
}

func (s *DockerTask) Result() *Result {
	return s.result
}
