// Code generated by go generate; DO NOT EDIT.
package protocol

// A set of predefined token types. This set is not fixed
// an clients can specify additional token types via the
// corresponding client capabilities.
// 
// @since 3.16.0
type SemanticTokenTypes string

const (
	SemanticTokenTypesnamespace SemanticTokenTypes = "namespace"
	// Represents a generic type. Acts as a fallback for types which can't be mapped to
	// a specific type like class or enum.
	SemanticTokenTypestype SemanticTokenTypes = "type"
	SemanticTokenTypesclass SemanticTokenTypes = "class"
	SemanticTokenTypesenum SemanticTokenTypes = "enum"
	SemanticTokenTypesinterface SemanticTokenTypes = "interface"
	SemanticTokenTypesstruct SemanticTokenTypes = "struct"
	SemanticTokenTypestypeParameter SemanticTokenTypes = "typeParameter"
	SemanticTokenTypesparameter SemanticTokenTypes = "parameter"
	SemanticTokenTypesvariable SemanticTokenTypes = "variable"
	SemanticTokenTypesproperty SemanticTokenTypes = "property"
	SemanticTokenTypesenumMember SemanticTokenTypes = "enumMember"
	SemanticTokenTypesevent SemanticTokenTypes = "event"
	SemanticTokenTypesfunction SemanticTokenTypes = "function"
	SemanticTokenTypesmethod SemanticTokenTypes = "method"
	SemanticTokenTypesmacro SemanticTokenTypes = "macro"
	SemanticTokenTypeskeyword SemanticTokenTypes = "keyword"
	SemanticTokenTypesmodifier SemanticTokenTypes = "modifier"
	SemanticTokenTypescomment SemanticTokenTypes = "comment"
	SemanticTokenTypesstring SemanticTokenTypes = "string"
	SemanticTokenTypesnumber SemanticTokenTypes = "number"
	SemanticTokenTypesregexp SemanticTokenTypes = "regexp"
	SemanticTokenTypesoperator SemanticTokenTypes = "operator"
	// @since 3.17.0
	SemanticTokenTypesdecorator SemanticTokenTypes = "decorator"
)

// A set of predefined token modifiers. This set is not fixed
// an clients can specify additional token types via the
// corresponding client capabilities.
// 
// @since 3.16.0
type SemanticTokenModifiers string

const (
	SemanticTokenModifiersdeclaration SemanticTokenModifiers = "declaration"
	SemanticTokenModifiersdefinition SemanticTokenModifiers = "definition"
	SemanticTokenModifiersreadonly SemanticTokenModifiers = "readonly"
	SemanticTokenModifiersstatic SemanticTokenModifiers = "static"
	SemanticTokenModifiersdeprecated SemanticTokenModifiers = "deprecated"
	SemanticTokenModifiersabstract SemanticTokenModifiers = "abstract"
	SemanticTokenModifiersasync SemanticTokenModifiers = "async"
	SemanticTokenModifiersmodification SemanticTokenModifiers = "modification"
	SemanticTokenModifiersdocumentation SemanticTokenModifiers = "documentation"
	SemanticTokenModifiersdefaultLibrary SemanticTokenModifiers = "defaultLibrary"
)

// The document diagnostic report kinds.
// 
// @since 3.17.0
type DocumentDiagnosticReportKind string

const (
	// A diagnostic report with a full
	// set of problems.
	DocumentDiagnosticReportKindFull DocumentDiagnosticReportKind = "full"
	// A report indicating that the last
	// returned report is still accurate.
	DocumentDiagnosticReportKindUnchanged DocumentDiagnosticReportKind = "unchanged"
)

// Predefined error codes.
type ErrorCodes int

const (
	ErrorCodesParseError ErrorCodes = -32700
	ErrorCodesInvalidRequest ErrorCodes = -32600
	ErrorCodesMethodNotFound ErrorCodes = -32601
	ErrorCodesInvalidParams ErrorCodes = -32602
	ErrorCodesInternalError ErrorCodes = -32603
	// Error code indicating that a server received a notification or
	// request before the server has received the `initialize` request.
	ErrorCodesServerNotInitialized ErrorCodes = -32002
	ErrorCodesUnknownErrorCode ErrorCodes = -32001
)

type LSPErrorCodes int

const (
	// A request failed but it was syntactically correct, e.g the
	// method name was known and the parameters were valid. The error
	// message should contain human readable information about why
	// the request failed.
	// 
	// @since 3.17.0
	LSPErrorCodesRequestFailed LSPErrorCodes = -32803
	// The server cancelled the request. This error code should
	// only be used for requests that explicitly support being
	// server cancellable.
	// 
	// @since 3.17.0
	LSPErrorCodesServerCancelled LSPErrorCodes = -32802
	// The server detected that the content of a document got
	// modified outside normal conditions. A server should
	// NOT send this error code if it detects a content change
	// in it unprocessed messages. The result even computed
	// on an older state might still be useful for the client.
	// 
	// If a client decides that a result is not of any use anymore
	// the client should cancel the request.
	LSPErrorCodesContentModified LSPErrorCodes = -32801
	// The client has canceled a request and a server as detected
	// the cancel.
	LSPErrorCodesRequestCancelled LSPErrorCodes = -32800
)

// A set of predefined range kinds.
type FoldingRangeKind string

const (
	// Folding range for a comment
	FoldingRangeKindComment FoldingRangeKind = "comment"
	// Folding range for an import or include
	FoldingRangeKindImports FoldingRangeKind = "imports"
	// Folding range for a region (e.g. `#region`)
	FoldingRangeKindRegion FoldingRangeKind = "region"
)

// A symbol kind.
type SymbolKind uint32

const (
	SymbolKindFile SymbolKind = 1
	SymbolKindModule SymbolKind = 2
	SymbolKindNamespace SymbolKind = 3
	SymbolKindPackage SymbolKind = 4
	SymbolKindClass SymbolKind = 5
	SymbolKindMethod SymbolKind = 6
	SymbolKindProperty SymbolKind = 7
	SymbolKindField SymbolKind = 8
	SymbolKindConstructor SymbolKind = 9
	SymbolKindEnum SymbolKind = 10
	SymbolKindInterface SymbolKind = 11
	SymbolKindFunction SymbolKind = 12
	SymbolKindVariable SymbolKind = 13
	SymbolKindConstant SymbolKind = 14
	SymbolKindString SymbolKind = 15
	SymbolKindNumber SymbolKind = 16
	SymbolKindBoolean SymbolKind = 17
	SymbolKindArray SymbolKind = 18
	SymbolKindObject SymbolKind = 19
	SymbolKindKey SymbolKind = 20
	SymbolKindNull SymbolKind = 21
	SymbolKindEnumMember SymbolKind = 22
	SymbolKindStruct SymbolKind = 23
	SymbolKindEvent SymbolKind = 24
	SymbolKindOperator SymbolKind = 25
	SymbolKindTypeParameter SymbolKind = 26
)

// Symbol tags are extra annotations that tweak the rendering of a symbol.
// 
// @since 3.16
type SymbolTag uint32

const (
	// Render a symbol as obsolete, usually using a strike-out.
	SymbolTagDeprecated SymbolTag = 1
)

// Moniker uniqueness level to define scope of the moniker.
// 
// @since 3.16.0
type UniquenessLevel string

const (
	// The moniker is only unique inside a document
	UniquenessLeveldocument UniquenessLevel = "document"
	// The moniker is unique inside a project for which a dump got created
	UniquenessLevelproject UniquenessLevel = "project"
	// The moniker is unique inside the group to which a project belongs
	UniquenessLevelgroup UniquenessLevel = "group"
	// The moniker is unique inside the moniker scheme.
	UniquenessLevelscheme UniquenessLevel = "scheme"
	// The moniker is globally unique
	UniquenessLevelglobal UniquenessLevel = "global"
)

// The moniker kind.
// 
// @since 3.16.0
type MonikerKind string

const (
	// The moniker represent a symbol that is imported into a project
	MonikerKindimport MonikerKind = "import"
	// The moniker represents a symbol that is exported from a project
	MonikerKindexport MonikerKind = "export"
	// The moniker represents a symbol that is local to a project (e.g. a local
	// variable of a function, a class not visible outside the project, ...)
	MonikerKindlocal MonikerKind = "local"
)

// Inlay hint kinds.
// 
// @since 3.17.0
type InlayHintKind uint32

const (
	// An inlay hint that for a type annotation.
	InlayHintKindType InlayHintKind = 1
	// An inlay hint that is for a parameter.
	InlayHintKindParameter InlayHintKind = 2
)

// The message type
type MessageType uint32

const (
	// An error message.
	MessageTypeError MessageType = 1
	// A warning message.
	MessageTypeWarning MessageType = 2
	// An information message.
	MessageTypeInfo MessageType = 3
	// A log message.
	MessageTypeLog MessageType = 4
)

// Defines how the host (editor) should sync
// document changes to the language server.
type TextDocumentSyncKind uint32

const (
	// Documents should not be synced at all.
	TextDocumentSyncKindNone TextDocumentSyncKind = 0
	// Documents are synced by always sending the full content
	// of the document.
	TextDocumentSyncKindFull TextDocumentSyncKind = 1
	// Documents are synced by sending the full content on open.
	// After that only incremental updates to the document are
	// send.
	TextDocumentSyncKindIncremental TextDocumentSyncKind = 2
)

// Represents reasons why a text document is saved.
type TextDocumentSaveReason uint32

const (
	// Manually triggered, e.g. by the user pressing save, by starting debugging,
	// or by an API call.
	TextDocumentSaveReasonManual TextDocumentSaveReason = 1
	// Automatic after a delay.
	TextDocumentSaveReasonAfterDelay TextDocumentSaveReason = 2
	// When the editor lost focus.
	TextDocumentSaveReasonFocusOut TextDocumentSaveReason = 3
)

// The kind of a completion entry.
type CompletionItemKind uint32

const (
	CompletionItemKindText CompletionItemKind = 1
	CompletionItemKindMethod CompletionItemKind = 2
	CompletionItemKindFunction CompletionItemKind = 3
	CompletionItemKindConstructor CompletionItemKind = 4
	CompletionItemKindField CompletionItemKind = 5
	CompletionItemKindVariable CompletionItemKind = 6
	CompletionItemKindClass CompletionItemKind = 7
	CompletionItemKindInterface CompletionItemKind = 8
	CompletionItemKindModule CompletionItemKind = 9
	CompletionItemKindProperty CompletionItemKind = 10
	CompletionItemKindUnit CompletionItemKind = 11
	CompletionItemKindValue CompletionItemKind = 12
	CompletionItemKindEnum CompletionItemKind = 13
	CompletionItemKindKeyword CompletionItemKind = 14
	CompletionItemKindSnippet CompletionItemKind = 15
	CompletionItemKindColor CompletionItemKind = 16
	CompletionItemKindFile CompletionItemKind = 17
	CompletionItemKindReference CompletionItemKind = 18
	CompletionItemKindFolder CompletionItemKind = 19
	CompletionItemKindEnumMember CompletionItemKind = 20
	CompletionItemKindConstant CompletionItemKind = 21
	CompletionItemKindStruct CompletionItemKind = 22
	CompletionItemKindEvent CompletionItemKind = 23
	CompletionItemKindOperator CompletionItemKind = 24
	CompletionItemKindTypeParameter CompletionItemKind = 25
)

// Completion item tags are extra annotations that tweak the rendering of a completion
// item.
// 
// @since 3.15.0
type CompletionItemTag uint32

const (
	// Render a completion as obsolete, usually using a strike-out.
	CompletionItemTagDeprecated CompletionItemTag = 1
)

// Defines whether the insert text in a completion item should be interpreted as
// plain text or a snippet.
type InsertTextFormat uint32

const (
	// The primary text to be inserted is treated as a plain string.
	InsertTextFormatPlainText InsertTextFormat = 1
	// The primary text to be inserted is treated as a snippet.
	// 
	// A snippet can define tab stops and placeholders with `$1`, `$2`
	// and `${3:foo}`. `$0` defines the final tab stop, it defaults to
	// the end of the snippet. Placeholders with equal identifiers are linked,
	// that is typing in one will update others too.
	// 
	// See also: https://microsoft.github.io/language-server-protocol/specifications/specification-current/#snippet_syntax
	InsertTextFormatSnippet InsertTextFormat = 2
)

// How whitespace and indentation is handled during completion
// item insertion.
// 
// @since 3.16.0
type InsertTextMode uint32

const (
	// The insertion or replace strings is taken as it is. If the
	// value is multi line the lines below the cursor will be
	// inserted using the indentation defined in the string value.
	// The client will not apply any kind of adjustments to the
	// string.
	InsertTextModeasIs InsertTextMode = 1
	// The editor adjusts leading whitespace of new lines so that
	// they match the indentation up to the cursor of the line for
	// which the item is accepted.
	// 
	// Consider a line like this: <2tabs><cursor><3tabs>foo. Accepting a
	// multi line completion item is indented using 2 tabs and all
	// following lines inserted will be indented using 2 tabs as well.
	InsertTextModeadjustIndentation InsertTextMode = 2
)

// A document highlight kind.
type DocumentHighlightKind uint32

const (
	// A textual occurrence.
	DocumentHighlightKindText DocumentHighlightKind = 1
	// Read-access of a symbol, like reading a variable.
	DocumentHighlightKindRead DocumentHighlightKind = 2
	// Write-access of a symbol, like writing to a variable.
	DocumentHighlightKindWrite DocumentHighlightKind = 3
)

// A set of predefined code action kinds
type CodeActionKind string

const (
	// Empty kind.
	CodeActionKindEmpty CodeActionKind = ""
	// Base kind for quickfix actions: 'quickfix'
	CodeActionKindQuickFix CodeActionKind = "quickfix"
	// Base kind for refactoring actions: 'refactor'
	CodeActionKindRefactor CodeActionKind = "refactor"
	// Base kind for refactoring extraction actions: 'refactor.extract'
	// 
	// Example extract actions:
	// 
	// - Extract method
	// - Extract function
	// - Extract variable
	// - Extract interface from class
	// - ...
	CodeActionKindRefactorExtract CodeActionKind = "refactor.extract"
	// Base kind for refactoring inline actions: 'refactor.inline'
	// 
	// Example inline actions:
	// 
	// - Inline function
	// - Inline variable
	// - Inline constant
	// - ...
	CodeActionKindRefactorInline CodeActionKind = "refactor.inline"
	// Base kind for refactoring rewrite actions: 'refactor.rewrite'
	// 
	// Example rewrite actions:
	// 
	// - Convert JavaScript function to class
	// - Add or remove parameter
	// - Encapsulate field
	// - Make method static
	// - Move method to base class
	// - ...
	CodeActionKindRefactorRewrite CodeActionKind = "refactor.rewrite"
	// Base kind for source actions: `source`
	// 
	// Source code actions apply to the entire file.
	CodeActionKindSource CodeActionKind = "source"
	// Base kind for an organize imports source action: `source.organizeImports`
	CodeActionKindSourceOrganizeImports CodeActionKind = "source.organizeImports"
	// Base kind for auto-fix source actions: `source.fixAll`.
	// 
	// Fix all actions automatically fix errors that have a clear fix that do not require user input.
	// They should not suppress errors or perform unsafe fixes such as generating new types or classes.
	// 
	// @since 3.15.0
	CodeActionKindSourceFixAll CodeActionKind = "source.fixAll"
)

type TraceValues string

const (
	// Turn tracing off.
	TraceValuesOff TraceValues = "off"
	// Trace messages only.
	TraceValuesMessages TraceValues = "messages"
	// Verbose message tracing.
	TraceValuesVerbose TraceValues = "verbose"
)

// Describes the content type that a client supports in various
// result literals like `Hover`, `ParameterInfo` or `CompletionItem`.
// 
// Please note that `MarkupKinds` must not start with a `$`. This kinds
// are reserved for internal usage.
type MarkupKind string

const (
	// Plain text is supported as a content format
	MarkupKindPlainText MarkupKind = "plaintext"
	// Markdown is supported as a content format
	MarkupKindMarkdown MarkupKind = "markdown"
)

// A set of predefined position encoding kinds.
// 
// @since 3.17.0
type PositionEncodingKind string

const (
	// Character offsets count UTF-8 code units.
	PositionEncodingKindUTF8 PositionEncodingKind = "utf-8"
	// Character offsets count UTF-16 code units.
	// 
	// This is the default and must always be supported
	// by servers
	PositionEncodingKindUTF16 PositionEncodingKind = "utf-16"
	// Character offsets count UTF-32 code units.
	// 
	// Implementation note: these are the same as Unicode code points,
	// so this `PositionEncodingKind` may also be used for an
	// encoding-agnostic representation of character offsets.
	PositionEncodingKindUTF32 PositionEncodingKind = "utf-32"
)

// The file event type
type FileChangeType uint32

const (
	// The file got created.
	FileChangeTypeCreated FileChangeType = 1
	// The file got changed.
	FileChangeTypeChanged FileChangeType = 2
	// The file got deleted.
	FileChangeTypeDeleted FileChangeType = 3
)

type WatchKind uint32

const (
	// Interested in create events.
	WatchKindCreate WatchKind = 1
	// Interested in change events
	WatchKindChange WatchKind = 2
	// Interested in delete events
	WatchKindDelete WatchKind = 4
)

// The diagnostic's severity.
type DiagnosticSeverity uint32

const (
	// Reports an error.
	DiagnosticSeverityError DiagnosticSeverity = 1
	// Reports a warning.
	DiagnosticSeverityWarning DiagnosticSeverity = 2
	// Reports an information.
	DiagnosticSeverityInformation DiagnosticSeverity = 3
	// Reports a hint.
	DiagnosticSeverityHint DiagnosticSeverity = 4
)

// The diagnostic tags.
// 
// @since 3.15.0
type DiagnosticTag uint32

const (
	// Unused or unnecessary code.
	// 
	// Clients are allowed to render diagnostics with this tag faded out instead of having
	// an error squiggle.
	DiagnosticTagUnnecessary DiagnosticTag = 1
	// Deprecated or obsolete code.
	// 
	// Clients are allowed to rendered diagnostics with this tag strike through.
	DiagnosticTagDeprecated DiagnosticTag = 2
)

// How a completion was triggered
type CompletionTriggerKind uint32

const (
	// Completion was triggered by typing an identifier (24x7 code
	// complete), manual invocation (e.g Ctrl+Space) or via API.
	CompletionTriggerKindInvoked CompletionTriggerKind = 1
	// Completion was triggered by a trigger character specified by
	// the `triggerCharacters` properties of the `CompletionRegistrationOptions`.
	CompletionTriggerKindTriggerCharacter CompletionTriggerKind = 2
	// Completion was re-triggered as current completion list is incomplete
	CompletionTriggerKindTriggerForIncompleteCompletions CompletionTriggerKind = 3
)

// How a signature help was triggered.
// 
// @since 3.15.0
type SignatureHelpTriggerKind uint32

const (
	// Signature help was invoked manually by the user or by a command.
	SignatureHelpTriggerKindInvoked SignatureHelpTriggerKind = 1
	// Signature help was triggered by a trigger character.
	SignatureHelpTriggerKindTriggerCharacter SignatureHelpTriggerKind = 2
	// Signature help was triggered by the cursor moving or by the document content changing.
	SignatureHelpTriggerKindContentChange SignatureHelpTriggerKind = 3
)

// The reason why code actions were requested.
// 
// @since 3.17.0
type CodeActionTriggerKind uint32

const (
	// Code actions were explicitly requested by the user or by an extension.
	CodeActionTriggerKindInvoked CodeActionTriggerKind = 1
	// Code actions were requested automatically.
	// 
	// This typically happens when current selection in a file changes, but can
	// also be triggered when file content changes.
	CodeActionTriggerKindAutomatic CodeActionTriggerKind = 2
)

// A pattern kind describing if a glob pattern matches a file a folder or
// both.
// 
// @since 3.16.0
type FileOperationPatternKind string

const (
	// The pattern matches a file only.
	FileOperationPatternKindfile FileOperationPatternKind = "file"
	// The pattern matches a folder only.
	FileOperationPatternKindfolder FileOperationPatternKind = "folder"
)

// A notebook cell kind.
// 
// @since 3.17.0
type NotebookCellKind uint32

const (
	// A markup-cell is formatted source that is used for display.
	NotebookCellKindMarkup NotebookCellKind = 1
	// A code-cell is source code.
	NotebookCellKindCode NotebookCellKind = 2
)

type ResourceOperationKind string

const (
	// Supports creating new files and folders.
	ResourceOperationKindCreate ResourceOperationKind = "create"
	// Supports renaming existing files and folders.
	ResourceOperationKindRename ResourceOperationKind = "rename"
	// Supports deleting existing files and folders.
	ResourceOperationKindDelete ResourceOperationKind = "delete"
)

type FailureHandlingKind string

const (
	// Applying the workspace change is simply aborted if one of the changes provided
	// fails. All operations executed before the failing operation stay executed.
	FailureHandlingKindAbort FailureHandlingKind = "abort"
	// All operations are executed transactional. That means they either all
	// succeed or no changes at all are applied to the workspace.
	FailureHandlingKindTransactional FailureHandlingKind = "transactional"
	// If the workspace edit contains only textual file changes they are executed transactional.
	// If resource changes (create, rename or delete file) are part of the change the failure
	// handling strategy is abort.
	FailureHandlingKindTextOnlyTransactional FailureHandlingKind = "textOnlyTransactional"
	// The client tries to undo the operations already executed. But there is no
	// guarantee that this is succeeding.
	FailureHandlingKindUndo FailureHandlingKind = "undo"
)

type PrepareSupportDefaultBehavior uint32

const (
	// The client's default behavior is to select the identifier
	// according the to language's syntax rule.
	PrepareSupportDefaultBehaviorIdentifier PrepareSupportDefaultBehavior = 1
)

type TokenFormat string

const (
	TokenFormatRelative TokenFormat = "relative"
)
