package datatype

import (
	"crypto/sha256"
	"encoding/hex"
	"evalevm/internal/export"
	"fmt"
	"log"
	"strings"

	jsoniter "github.com/json-iterator/go"
	"github.com/zerjioang/rooftop/v2/common/io/json"
)

// DockerTask represents a Docker command to be executed
type DockerTask struct {
	id          TaskId
	sampleId    string
	opts        BytecodeScanOpts
	hexBytecode string
	cmd         []string
	finishCh    chan struct{}
	result      *Result
	parser      ResultParser
	debug       bool
}

func (s *DockerTask) TrackerId() string {
	return s.sampleId
}

func NewDockerTask(id TaskId, opts BytecodeScanOpts, hexBytecode string, filename string, cmd []string) *DockerTask {
	sampleId := filename
	if filename == "" {
		sampleId = sha256Hex(hexBytecode)
	}
	return &DockerTask{
		id:          id,
		cmd:         cmd,
		opts:        opts,
		hexBytecode: hexBytecode,
		sampleId:    sampleId,
		finishCh:    make(chan struct{}, 1),
		debug:       false,
	}
}

// Sha256Hex returns the SHA-256 hash of the input string as a hex string.
func sha256Hex(input string) string {
	hash := sha256.Sum256([]byte(input)) // returns [32]byte
	return hex.EncodeToString(hash[:])   // convert array to slice and encode
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
	return s.result.Error != nil || s.result.Timeout || s.result.Measurements.ExitStatus != 0
}

func (s *DockerTask) Parse() {
	// if present, extract the measure.sh output
	output := string(s.result.Output)
	idxStart := strings.LastIndex(output, "evalevm_perf_metrics_start")
	idxEnd := strings.LastIndex(output, "evalevm_perf_metrics_end")
	if idxStart != -1 && idxEnd != -1 && idxEnd > idxStart {
		measureJsonStr := strings.TrimSpace(string(s.result.Output[idxStart+len("evalevm_perf_metrics_start") : idxEnd]))
		if err := json.Unmarshal([]byte(measureJsonStr), &s.result.Measurements); err != nil {
			log.Println("failed to parse measure.sh output: " + err.Error() + ". JSON: " + measureJsonStr)
		}
		s.result.Output = removeBytes(s.result.Output, idxStart, idxEnd+len("evalevm_perf_metrics_end"))
	} else {
		log.Println("failed to parse measure.sh output: missing output")
		log.Println(s.result.Output)
		s.result.Measurements = Measurements{
			ExitStatus:    1,
			FailedToParse: true,
		}
	}
	if s.parser != nil {
		// parse the output of the app being executed
		if s.debug {
			filename := fmt.Sprintf("debug_%s_out_%s.txt", s.result.Task.ID().App(), s.result.Task.ID().identifier)
			_ = export.DebugFile(filename, s.result.Output)

			filenameErr := fmt.Sprintf("debug_%s_err_%s.txt", s.result.Task.ID().App(), s.result.Task.ID().identifier)
			_ = export.DebugFile(filenameErr, s.result.OutputErr)
		}
		if err := s.parser.ParseOutput(s.result); err != nil {
			log.Println("failed to parse output: " + err.Error())
		}
		_ = s.writeResult(s.result)
	}
}

func (s *DockerTask) Result() *Result {
	return s.result
}

func (s *DockerTask) writeResult(result *Result) error {
	return export.DebugFile(s.result.Filename(), []byte(result.String()))
}

// MarshalJSON implements custom JSON marshalling for DockerTask
func (s *DockerTask) MarshalJSON() ([]byte, error) {
	// Inline struct with only the fields we want to expose
	type Alias struct {
		ID          TaskId           `json:"id"`
		SampleID    string           `json:"sample_id"`
		Opts        BytecodeScanOpts `json:"opts"`
		HexBytecode string           `json:"hex_bytecode"`
		Cmd         []string         `json:"cmd"`
		Result      *Result          `json:"result,omitempty"`
		Debug       bool             `json:"debug"`
	}

	// Note: we intentionally pass nil for Result to avoid circular reference.
	// Do NOT mutate s.result here — that would cause a data race.
	return jsoniter.Marshal(&Alias{
		ID:          s.id,
		SampleID:    s.sampleId,
		Opts:        s.opts,
		HexBytecode: s.hexBytecode,
		Cmd:         s.cmd,
		Result:      nil,
		Debug:       s.debug,
	})
}

func (s *DockerTask) ContainerName() string {
	return fmt.Sprintf("%s_%s", s.id.app, s.sampleId)
}

// removeBytes removes the subset of b between start (inclusive) and end (exclusive)
func removeBytes(b []byte, start, end int) []byte {
	// Validate indices
	if start < 0 || end > len(b) || start > end {
		return b // return original if invalid
	}

	// Concatenate before start and after end
	result := append([]byte{}, b[:start]...) // copy before start
	result = append(result, b[end:]...)      // add after end
	return result
}
