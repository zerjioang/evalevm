package paper

import (
	_ "embed"
	"evalevm/internal/datatype"
	"evalevm/internal/parser"
	"fmt"
	"log"
	"strings"

	"github.com/zerjioang/rooftop/v2/common/io/json"
)

type Paper struct {
	datatype.BytecodeAnalyzer
}

var (
	//go:embed Dockerfile
	paperDockerfile string
)

var _ datatype.Analyzer = (*Paper)(nil)

func NewPaper() Paper {
	app := Paper{}
	app.BytecodeAnalyzer = app.SetupDockerPlatform()
	app.AppName = "paper"
	app.WebsiteUrl = "https://github.com/zerjioang"
	app.Desc = "research tool"
	app.Options = datatype.BytecodeScanOpts{
		ForceRemoveHexPrefix: false,
		ForceSplitRuntime:    false,
	}
	app.Deprecated = false
	app.LastCommit = "today"
	app.Language = "go"
	app.Dockerfile = paperDockerfile
	return app
}

func (scan Paper) CreateTask(uid string, bytecode string) []datatype.Task {
	// docker run --rm paper-algorithm:latest 6080604052600436106100745763ffffffff60e060020a600035041663025313a2811461022e5780633ad06d161461025f57806354fd4d50146102855780635c60da1b146102ac5780636fde8202146102c1578063a9c45fcb146102d6578063d784d42614610332578063f1739cae14610353575b600080600080610082610374565b9350600160a060020a038416151561009957600080fd5b30905083600160a060020a0316635c60da1b6040518163ffffffff1660e060020a0281526004016000604051808303816000875af192505050156101e15783600160a060020a0316635c60da1b6040518163ffffffff1660e060020a028152600401602060405180830381600087803b15801561011557600080fd5b505af1158015610129573d6000803e3d6000fd5b505050506040513d602081101561013f57600080fd5b5051604080517fd784d426000000000000000000000000000000000000000000000000000000008152600160a060020a03831660048201529051919450309163d784d4269160248082019260009290919082900301818387803b1580156101a557600080fd5b505af11580156101b9573d6000803e3d6000fd5b507fd784d4260000000000000000000000000000000000000000000000000000000094505050505b60405136600082376000803683885af43d82016040523d6000833e3d84801561021e5760405186815288600482015260008160248360008a5af150505b5081801561022a578184f35b8184fd5b34801561023a57600080fd5b50610243610383565b60408051600160a060020a039092168252519081900360200190f35b34801561026b57600080fd5b50610283600435600160a060020a0360243516610392565b005b34801561029157600080fd5b5061029a6103bc565b60408051918252519081900360200190f35b3480156102b857600080fd5b50610243610374565b3480156102cd57600080fd5b506102436103c2565b604080516020600460443581810135601f81018490048402850184019095528484526102839482359460248035600160a060020a0316953695946064949201919081908401838280828437509497506103d19650505050505050565b34801561033e57600080fd5b50610283600160a060020a036004351661047f565b34801561035f57600080fd5b50610283600160a060020a03600435166104ba565b600254600160a060020a031690565b600061038d6103c2565b905090565b61039a610383565b600160a060020a031633146103ae57600080fd5b6103b88282610542565b5050565b60015490565b600054600160a060020a031690565b6103d9610383565b600160a060020a031633146103ed57600080fd5b6103f78383610392565b30600160a060020a0316348260405180828051906020019080838360005b8381101561042d578181015183820152602001610415565b50505050905090810190601f16801561045a5780820380516001836020036101000a031916815260200191505b5091505060006040518083038185875af192505050151561047a57600080fd5b505050565b33301461048b57600080fd5b6002805473ffffffffffffffffffffffffffffffffffffffff1916600160a060020a0392909216919091179055565b6104c2610383565b600160a060020a031633146104d657600080fd5b600160a060020a03811615156104eb57600080fd5b7f5a3e66efaa1e445ebd894728a69d6959842ea1e97bd79b892797106e270efcd9610514610383565b60408051600160a060020a03928316815291841660208301528051918290030190a161053f816105d3565b50565b600254600160a060020a038281169116141561055d57600080fd5b600154821161056b57600080fd5b600182905560028054600160a060020a03831673ffffffffffffffffffffffffffffffffffffffff1990911681179091556040805184815290517f4289d6195cf3c2d2174adf98d0e19d4d2d08887995b99cb7b100e7ffe795820e9181900360200190a25050565b6000805473ffffffffffffffffffffffffffffffffffffffff1916600160a060020a03929092169190911790555600a165627a7a723058209b97a6159e34e4a58bc0963f5c3bc078781a49f9a924713fa88b8a830bb1e6c90029
	return []datatype.Task{
		datatype.NewDockerTask(
			scan.CreateTaskId(uid),
			scan.Options,
			bytecode,
			[]string{
				// docker run command already defined. customize the flags here
				"local/paper:latest",
				"-c",
				fmt.Sprintf(`./measure.sh /app/detector-linux.elf %s`, bytecode),
			},
		),
	}
}

type paperOutput struct {
	Name                string      `json:"name"`
	ExecTimeMs          int         `json:"exec_time_ms"`
	ExecTimeMicros      int         `json:"exec_time_micros"`
	Errored             bool        `json:"errored"`
	Error               string      `json:"error"`
	Timeout             bool        `json:"timeout"`
	MetadataDetected    bool        `json:"metadata_detected"`
	MetadataSection     string      `json:"metadata_section"`
	MetadataHash        string      `json:"metadata_hash"`
	MetadataSectionSize int         `json:"metadata_section_size"`
	SolidityVersion     string      `json:"solidity_version"`
	CfgNodeCount        int         `json:"cfg_node_count"`
	Bytecode            string      `json:"bytecode"`
	IsOnlyRuntime       bool        `json:"is_only_runtime"`
	Vulnerabilities     []any       `json:"vulnerabilities"`
	Graphs              Graphs      `json:"graphs"`
	SectionData         SectionData `json:"section_data"`
	Coverage            float64     `json:"coverage"`
}
type Graphs struct {
	Runtime     string `json:".runtime"`
	Constructor string `json:".constructor"`
}
type SectionStats struct {
	RenderedBlocks           int     `json:"rendered_blocks"`
	ExecutedBlocks           int     `json:"executed_blocks"`
	TotalBlocks              int     `json:"total_blocks"`
	Coverage                 float64 `json:"coverage"`
	HiddenBlocks             int     `json:"hidden_blocks"`
	OpcodeRenderInstruction  int     `json:"opcode_render_instruction"`
	YulRenderInstruction     int     `json:"yul_render_instruction"`
	YulImprovementPercentage float64 `json:"yul_improvement_percentage"`
}
type SectionData struct {
	Runtime     SectionStats `json:".runtime"`
	Constructor SectionStats `json:".constructor"`
}

func (scan Paper) ParseOutput(output *datatype.Result) error {

	outJson, err := parser.ExtractBetween(string(output.Output), "paper_output_begin", "paper_output_end")
	if err != nil {
		return fmt.Errorf("failed to parse output: %w", err)
	}

	var dst paperOutput
	if err := json.Unmarshal([]byte(outJson), &dst); err != nil {
		log.Println("failed to parse output: ", err)
	}

	// count how many nodes are defined in the CFG
	// example: 119 [label="119: EXIT BLOCK\l" fillcolor=crimson ];
	nodesDetected := strings.Count(dst.Graphs.Constructor, ` [label="0x`) + strings.Count(dst.Graphs.Runtime, ` [label="0x`)

	// count how many edges are defined in the CFG
	// example: 119 -> 118 [label="119 -> 118\l" ];
	edgesDetected := strings.Count(dst.Graphs.Constructor, " -> ") + strings.Count(dst.Graphs.Runtime, " -> ")

	var asPtrBool = func(b bool) *bool { return &b }
	output.ParsedOutput = &datatype.ScanResult{
		Vulnerable:    asPtrBool(len(dst.Vulnerabilities) > 0),
		Error:         nil,
		EdgesDetected: edgesDetected,
		NodesDetected: nodesDetected,
		Coverage:      &dst.Coverage,
	}
	if dst.Graphs.Constructor != "" {
		filename := fmt.Sprintf("cfg_%s_%s_constructor.svg", output.Task.ID().App(), output.Task.TrackerId())
		output.ParsedOutput.WithGraph(dst.Graphs.Constructor)
		output.ParsedOutput.SaveGraph(dst.Graphs.Constructor, filename)
		output.AddFileReference(output.Task.ID().App(), output.Task.TrackerId(), filename)
	}

	filename := fmt.Sprintf("cfg_%s_%s_runtime.svg", output.Task.ID().App(), output.Task.TrackerId())
	output.ParsedOutput.WithGraph(dst.Graphs.Runtime)
	output.ParsedOutput.SaveGraph(dst.Graphs.Runtime, filename)
	output.AddFileReference(output.Task.ID().App(), output.Task.TrackerId(), filename)

	return nil
}
