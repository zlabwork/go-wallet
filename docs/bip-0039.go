package docs

import "strings"

type wordList struct {
}

func WordList() *wordList {
    return &wordList{}
}

func (word *wordList) ChineseSimple() []string {
    return strings.Split(chineseSimplified, "\n")
}

func (word *wordList) Japanese() []string {
    return strings.Split(japanese, "\n")
}

func (word *wordList) Korean() []string {
    return strings.Split(korean, "\n")
}
