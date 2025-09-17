package evmlisa

type output struct {
	Bytecode         string `json:"bytecode,omitempty"`
	Address          string `json:"address,omitempty"`
	EventsSignature  []any  `json:"events_signature,omitempty"`
	MnemonicBytecode string `json:"mnemonic_bytecode,omitempty"`
	Abi              []any  `json:"abi,omitempty"`
	BasicBlocks      []struct {
		Instructions []struct {
			Pc          int    `json:"pc,omitempty"`
			Instruction string `json:"instruction,omitempty"`
		} `json:"instructions,omitempty"`
		BackgroundColor string `json:"background_color,omitempty"`
		OutgoingEdges   []struct {
			Color  string `json:"color,omitempty"`
			Target int    `json:"target,omitempty"`
		} `json:"outgoing_edges,omitempty"`
		LastInstruction string `json:"last_instruction,omitempty"`
		ID              int    `json:"id,omitempty"`
		Label           string `json:"label,omitempty"`
	} `json:"basic_blocks,omitempty"`
	WorkingDirectory         string `json:"working_directory,omitempty"`
	FunctionsSignature       []any  `json:"functions_signature,omitempty"`
	LastPc                   int    `json:"last_pc,omitempty"`
	AbiFilePath              []any  `json:"abi_file_path,omitempty"`
	ExecutionTime            int    `json:"execution_time,omitempty"`
	BytecodeFilePath         string `json:"bytecode_file_path,omitempty"`
	MnemonicBytecodeFilePath string `json:"mnemonic_bytecode_file_path,omitempty"`
	Vulnerabilities          struct {
		TxOrigin                     int `json:"tx_origin,omitempty"`
		Reentrancy                   int `json:"reentrancy,omitempty"`
		RandomnessDependency         int `json:"randomness_dependency,omitempty"`
		RandomnessDependencyPossible int `json:"randomness_dependency_possible,omitempty"`
		TxOriginPossible             int `json:"tx_origin_possible,omitempty"`
	} `json:"vulnerabilities,omitempty"`
	BasicBlocksPc string `json:"basic_blocks_pc,omitempty"`
	Statistics    struct {
		DefinitelyUnreachableJumps int `json:"definitely_unreachable_jumps,omitempty"`
		TotalEdges                 int `json:"total_edges,omitempty"`
		UnsoundJumps               int `json:"unsound_jumps,omitempty"`
		TotalOpcodes               int `json:"total_opcodes,omitempty"`
		MaybeUnsoundJumps          int `json:"maybe_unsound_jumps,omitempty"`
		MaybeUnreachableJumps      int `json:"maybe_unreachable_jumps,omitempty"`
		TotalJumps                 int `json:"total_jumps,omitempty"`
		ResolvedJumps              int `json:"resolved_jumps,omitempty"`
	} `json:"statistics,omitempty"`
}
