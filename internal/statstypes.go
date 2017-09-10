package internal

type FuncMeasurement string

const (
	Lines        FuncMeasurement = "lines"
	TotalLines                   = "total-lines"
	Len                          = "len"
	TotalLen                     = "total-len"
	CommentLines                 = "comment-lines"
	Comments                     = "comments"
	Complexity                   = "complexity"
	MaxNesting                   = "max-nesting"
	TotalNesting                 = "total-nesting"
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
}

func isValidBasicType(ty FuncMeasurement) bool {
	for _, t := range AllTypes {
		if t == ty {
			return true
		}
	}
	return false
}
