// Code generated by go generate; DO NOT EDIT.
package protocol

// The definition of a symbol represented as one or many {@link Location locations}.
// For most programming languages there is only one location at which a symbol is
// defined.
// 
// Servers should prefer returning `DefinitionLink` over `Definition` if supported
// by the client.
type Definition any

// Information about where a symbol is defined.
// 
// Provides additional metadata over normal {@link Location location} definitions, including the range of
// the defining symbol
type DefinitionLink LocationLink

// LSP arrays.
// @since 3.17.0
type LSPArray any

// The LSP any type.
// Please note that strictly speaking a property with the value `undefined`
// can't be converted into JSON preserving the property name. However for
// convenience it is allowed and assumed that all these properties are
// optional as well.
// @since 3.17.0
type LSPAny any

// The declaration of a symbol representation as one or many {@link Location locations}.
type Declaration any

// Information about where a symbol is declared.
// 
// Provides additional metadata over normal {@link Location location} declarations, including the range of
// the declaring symbol.
// 
// Servers should prefer returning `DeclarationLink` over `Declaration` if supported
// by the client.
type DeclarationLink LocationLink

// Inline value information can be provided by different means:
// - directly as a text value (class InlineValueText).
// - as a name to use for a variable lookup (class InlineValueVariableLookup)
// - as an evaluatable expression (class InlineValueEvaluatableExpression)
// The InlineValue types combines all inline value types into one type.
// 
// @since 3.17.0
type InlineValue any

// The result of a document diagnostic pull request. A report can
// either be a full report containing all diagnostics for the
// requested document or an unchanged report indicating that nothing
// has changed in terms of diagnostics in comparison to the last
// pull request.
// 
// @since 3.17.0
type DocumentDiagnosticReport any

type PrepareRenameResult any

// A document selector is the combination of one or many document filters.
// 
// @sample `let sel:DocumentSelector = [{ language: 'typescript' }, { language: 'json', pattern: '**∕tsconfig.json' }]`;
// 
// The use of a string as a document filter is deprecated @since 3.16.0.
type DocumentSelector any

type ProgressToken any

// An identifier to refer to a change annotation stored with a workspace edit.
type ChangeAnnotationIdentifier string

// A workspace diagnostic document report.
// 
// @since 3.17.0
type WorkspaceDocumentDiagnosticReport any

// An event describing a change to a text document. If only a text is provided
// it is considered to be the full content of the document.
type TextDocumentContentChangeEvent any

// MarkedString can be used to render human readable text. It is either a markdown string
// or a code-block that provides a language and a code snippet. The language identifier
// is semantically equal to the optional language identifier in fenced code blocks in GitHub
// issues. See https://help.github.com/articles/creating-and-highlighting-code-blocks/#syntax-highlighting
// 
// The pair of a language and a value is an equivalent to markdown:
// ```${language}
// ${value}
// ```
// 
// Note that markdown strings will be sanitized - that means html will be escaped.
// @deprecated use MarkupContent instead.
type MarkedString any

// A document filter describes a top level text document or
// a notebook cell document.
// 
// @since 3.17.0 - proposed support for NotebookCellTextDocumentFilter.
type DocumentFilter any

// LSP object definition.
// @since 3.17.0
type LSPObject any

// The glob pattern. Either a string pattern or a relative pattern.
// 
// @since 3.17.0
type GlobPattern any

// A document filter denotes a document by different properties like
// the {@link TextDocument.languageId language}, the {@link Uri.scheme scheme} of
// its resource, or a glob-pattern that is applied to the {@link TextDocument.fileName path}.
// 
// Glob patterns can have the following syntax:
// - `*` to match one or more characters in a path segment
// - `?` to match on one character in a path segment
// - `**` to match any number of path segments, including none
// - `{}` to group sub patterns into an OR expression. (e.g. `**​/*.{ts,js}` matches all TypeScript and JavaScript files)
// - `[]` to declare a range of characters to match in a path segment (e.g., `example.[0-9]` to match on `example.0`, `example.1`, …)
// - `[!...]` to negate a range of characters to match in a path segment (e.g., `example.[!0-9]` to match on `example.a`, `example.b`, but not `example.0`)
// 
// @sample A language filter that applies to typescript files on disk: `{ language: 'typescript', scheme: 'file' }`
// @sample A language filter that applies to all package.json paths: `{ language: 'json', pattern: '**package.json' }`
// 
// @since 3.17.0
type TextDocumentFilter any

// A notebook document filter denotes a notebook document by
// different properties. The properties will be match
// against the notebook's URI (same as with documents)
// 
// @since 3.17.0
type NotebookDocumentFilter any

// The glob pattern to watch relative to the base path. Glob patterns can have the following syntax:
// - `*` to match one or more characters in a path segment
// - `?` to match on one character in a path segment
// - `**` to match any number of path segments, including none
// - `{}` to group conditions (e.g. `**​/*.{ts,js}` matches all TypeScript and JavaScript files)
// - `[]` to declare a range of characters to match in a path segment (e.g., `example.[0-9]` to match on `example.0`, `example.1`, …)
// - `[!...]` to negate a range of characters to match in a path segment (e.g., `example.[!0-9]` to match on `example.a`, `example.b`, but not `example.0`)
// 
// @since 3.17.0
type Pattern string
