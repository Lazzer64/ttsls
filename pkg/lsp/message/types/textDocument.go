package types

import "github.com/lazzer64/ttsls/pkg/uri"

type TextDocumentIdentifier struct {
	Uri uri.URI `json:"uri"`
}

type TextDocumentItem struct {
	Uri        uri.URI `json:"uri"`
	LanguageId string  `json:"languageId"`
	Version    int     `json:"version"`
	Text       string  `json:"text"`
}

type Position struct {
	Character int `json:"character"`
	Line      int `json:"line"`
}

type Range struct {
	Start Position `json:"start"`
	End   Position `json:"end"`
}

type TextDocumentContentChangeEvent struct {
	Range Range  `json:"range"`
	Text  string `json:"text"`
}

type TextEdit struct {
	Range   Range  `json:"range"`
	NewText string `json:"newText"`
}

type Command struct {
	Title     string        `json:"title"`
	Command   string        `json:"command"`
	Arguments []interface{} `json:"arguments,omitempty"`
}

type CompletionTriggerKind int

const (
	InvokedTrigger                         CompletionTriggerKind = 1
	TriggerCharacterTrigger                CompletionTriggerKind = 2
	TriggerForIncompleteCompletionsTrigger CompletionTriggerKind = 3
)

type CompletionItemKind int

const (
	TextCompletion          CompletionItemKind = 1
	MethodCompletion        CompletionItemKind = 2
	FunctionCompletion      CompletionItemKind = 3
	ConstructorCompletion   CompletionItemKind = 4
	FieldCompletion         CompletionItemKind = 5
	VariableCompletion      CompletionItemKind = 6
	ClassCompletion         CompletionItemKind = 7
	InterfaceCompletion     CompletionItemKind = 8
	ModuleCompletion        CompletionItemKind = 9
	PropertyCompletion      CompletionItemKind = 10
	UnitCompletion          CompletionItemKind = 11
	ValueCompletion         CompletionItemKind = 12
	EnumCompletion          CompletionItemKind = 13
	KeywordCompletion       CompletionItemKind = 14
	SnippetCompletion       CompletionItemKind = 15
	ColorCompletion         CompletionItemKind = 16
	FileCompletion          CompletionItemKind = 17
	ReferenceCompletion     CompletionItemKind = 18
	FolderCompletion        CompletionItemKind = 19
	EnumMemberCompletion    CompletionItemKind = 20
	ConstantCompletion      CompletionItemKind = 21
	StructCompletion        CompletionItemKind = 22
	EventCompletion         CompletionItemKind = 23
	OperatorCompletion      CompletionItemKind = 24
	TypeParameterCompletion CompletionItemKind = 25
)

type SignatureHelpTriggerKind int

const (
	SignatureHelpTriggerKindInvoked          SignatureHelpTriggerKind = 1
	SignatureHelpTriggerKindTriggerCharacter SignatureHelpTriggerKind = 2
	SignatureHelpTriggerKindContentChange    SignatureHelpTriggerKind = 3
)

type SignatureHelp struct {
	Signatures []struct{
		Label string `json:"label"`
		Documentation string `json:"documentation,omitempty"`
		Parameters []struct{
			Label string `json:"label"`
			Documentation string `json:"documentation,omitempty"`
		} `json:"parameters,omitempty"`
		ActiveParameter uint `json:"activeParameter,omitempty"`
	} `json:"signatures"`
	ActiveSignature uint `json:"activeSignature,omitempty"`
	ActiveParameter uint `json:"activeParameter,omitempty"`
}
