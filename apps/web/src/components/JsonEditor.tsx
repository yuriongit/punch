import { Editor, type Monaco } from "@monaco-editor/react"
import type { Dispatch, SetStateAction } from "react"

type Props = {
  jsonInput: string
  setJsonInput: Dispatch<SetStateAction<string>>
  size: number
}

export const JsonEditor = ({ jsonInput, setJsonInput, size }: Props) => {
  const handleEditorWillMount = (monaco: Monaco) => {
    monaco.editor.defineTheme("punch-theme", {
      base: "vs", // Light base since background is white
      inherit: true,
      rules: [
        // --- SYNTAX TOKENS ---
        { token: "string.key.json", foreground: "262626", fontStyle: "bold" }, // Deep Purple (Tailwind purple-800)
        { token: "string.value.json", foreground: "9810fa" }, // Base Purple (Tailwind purple-500)
        { token: "number", foreground: "0EA5E9" }, // Cotton Candy Blue (Tailwind sky-500)
        { token: "keyword.json", foreground: "0D59FF", fontStyle: "bold" }, // Dark Blue (Tailwind blue-800)
        { token: "delimiter.colon", foreground: "3F3F46" }, // Light Black (Tailwind zinc-700)
        { token: "delimiter.comma", foreground: "3F3F46" }, // Light Black (Tailwind zinc-700)
        { token: "delimiter.bracket", foreground: "581C87" }, // Dark Purple Brackets (purple-900)
      ],
      colors: {
        // --- INDENT GUIDE LINES ---
        "editorIndentGuide.background1": "#E9D5FF", // Soft Purple-200 line
        "editorIndentGuide.activeBackground1": "#b758fc", // Active purple-500 line

        // --- BRACKET PAIR COLORIZATION (Purple & Black spectrum) ---
        "editorBracketHighlight.foreground1": "#A855F7", // Purple-500
        "editorBracketHighlight.foreground2": "#2e054b", // Zinc-900 (Black)
        "editorBracketHighlight.foreground3": "#7E22CE", // Purple-700
        "editorBracketHighlight.foreground4": "#3F3F46", // Zinc-700 (Dark Charcoal)
        "editorBracketHighlight.foreground5": "#C084FC", // Purple-400
        "editorBracketHighlight.foreground6": "#2e054b", // Purple-900
        "editorBracketHighlight.unexpectedBracket.foreground": "#DC2626", // Red error

        // --- BRACKET MATCHING ---
        "editor.bracketMatch.background": "#F3E8FF", // Very light purple box fill
        "editor.bracketMatch.border": "#A855F7", // Solid purple-500 border

        // --- SELECTION & EDITOR HIGHLIGHTS ---
        "editor.background": "#00000000", // Transparent background
        "editor.lineHighlightBackground": "#FAF5FF", // Subtle purple-50 line tint
        "editor.selectionBackground": "#E9D5FF", // Light purple-200 selection
        "editor.inactiveSelectionBackground": "#F3E8FF",
        "editor.selectionHighlightBackground": "#F3E8FF80", // Subtle tint for matching symbols elsewhere
        "editor.wordHighlightBackground": "#F3E8FFB0", // Light tint when clicking/hovering a symbol
        "editor.wordHighlightStrongBackground": "#E9D5FF", // Light tint for write-access symbols
      },
    })
  }

  return (
    <Editor
      height={`${size}px`}
      defaultLanguage="json"
      theme="punch-theme"
      beforeMount={handleEditorWillMount}
      value={jsonInput}
      onChange={(value) => setJsonInput(value || "")}
      options={{
        extraEditorClassName: "font-black",
        minimap: { enabled: false },
        fontSize: 17,
        stickyScroll: {
          enabled: false,
        },
        fontWeight: "700",
        fontFamily: '"JetBrains Mono", monospace',
        formatOnType: true,
        fontLigatures: true,
        scrollBeyondLastLine: false,
        smoothScrolling: true,
        allowOverflow: true,
        cursorSmoothCaretAnimation: "on",
        cursorStyle: "block-outline",
        scrollbar: {
          vertical: "hidden",
          verticalScrollbarSize: 0,
        },
        insertSpaces: true,
        tabCompletion: "on",
        tabSize: 2,
        lineNumbers: "off",
        bracketPairColorization: {
          enabled: true,
        },
        roundedSelection: true,
      }}
    />
  )
}
