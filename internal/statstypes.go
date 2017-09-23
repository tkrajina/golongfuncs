package internal

type FuncMeasurement string

const (
	Lines                FuncMeasurement = "lines"
	TotalLines                           = "total_lines"
	Len                                  = "len"
	TotalLen                             = "total_len"
	CommentLines                         = "comment_lines"
	Comments                             = "comments"
	Complexity                           = "complexity"
	MaxNesting                           = "max_nesting"
	TotalNesting                         = "total_nesting"
	InputParams                          = "in_params"
	OutputParams                         = "out_params"
	Variables                            = "variables"
	Assignments                          = "assignments"
	Control                              = "control"
	Todos                                = "todos"
	TodosCaseinsensitive                 = "todos_case_insensitive"
)

var AllTypes = []FuncMeasurement{
	Lines,
	TotalLines,
	Len,
	TotalLen,
	CommentLines,
	Comments,
	Complexity,
	MaxNesting,
	TotalNesting,
	InputParams,
	OutputParams,
	Variables,
	Assignments,
	Control,
	Todos,
	TodosCaseinsensitive,
}

func isValidBasicType(ty FuncMeasurement) bool {
	for _, t := range AllTypes {
		if t == ty {
			return true
		}
	}
	return false
}
