import type { ReactNode } from "react"

interface JsonCodeBlockProps {
  data: unknown
  /**
   * Maps a dotted key-path to a trailing comment, rendered like JSONC:
   *   "children.1.expected_status": "defaults to 200"
   * Array indices are path segments too, e.g. "children.1.name".
   */
  comments?: Record<string, string>
}

const INDENT = "  "

function Key({ children }: { children: string }) {
  return <span className="text-[#262626]">{JSON.stringify(children)}</span>
}

function Colon() {
  return <span className="text-[#3F3F46]">: </span>
}

function Comma() {
  return <span className="text-[#3F3F46]">,</span>
}

function Bracket({ children }: { children: string }) {
  return <span className="text-[#581C87]">{children}</span>
}

function StringValue({ children }: { children: string }) {
  return <span className="text-[#9810fa]">{JSON.stringify(children)}</span>
}

function NumberValue({ children }: { children: number }) {
  return <span className="text-[#0EA5E9]">{String(children)}</span>
}

function KeywordValue({ children }: { children: string }) {
  return <span className="text-[#0D59FF]">{children}</span>
}

function CommentText({ children }: { children: string }) {
  return <span className="text-[#8C8C8C] italic"> {children}</span>
}

function renderValue(
  value: unknown,
  depth: number,
  path: string,
  comments: Record<string, string>,
): ReactNode {
  if (value === null) return <KeywordValue key={path}>null</KeywordValue>
  if (typeof value === "boolean") {
    return <KeywordValue key={path}>{value ? "true" : "false"}</KeywordValue>
  }
  if (typeof value === "number") {
    return <NumberValue key={path}>{value}</NumberValue>
  }
  if (typeof value === "string") {
    return <StringValue key={path}>{value}</StringValue>
  }

  const pad = INDENT.repeat(depth)
  const innerPad = INDENT.repeat(depth + 1)

  if (Array.isArray(value)) {
    if (value.length === 0) {
      return (
        <span key={path}>
          <Bracket>[</Bracket>
          <Bracket>]</Bracket>
        </span>
      )
    }
    return (
      <span key={path}>
        <Bracket>[</Bracket>
        {"\n"}
        {value.map((item, i) => {
          const itemPath = `${path}.${i}`
          const comment = comments[itemPath]
          return (
            <span key={itemPath}>
              {innerPad}
              {renderValue(item, depth + 1, itemPath, comments)}
              {i < value.length - 1 ? <Comma /> : null}
              {comment ? <CommentText>{comment}</CommentText> : null}
              {"\n"}
            </span>
          )
        })}
        {pad}
        <Bracket>]</Bracket>
      </span>
    )
  }

  if (typeof value === "object") {
    const entries = Object.entries(value as Record<string, unknown>)
    if (entries.length === 0) {
      return (
        <span key={path}>
          <Bracket>{"{"}</Bracket>
          <Bracket>{"}"}</Bracket>
        </span>
      )
    }
    return (
      <span key={path}>
        <Bracket>{"{"}</Bracket>
        {"\n"}
        {entries.map(([k, v], i) => {
          const entryPath = path ? `${path}.${k}` : k
          const comment = comments[entryPath]
          return (
            <span key={k}>
              {innerPad}
              <Key>{k}</Key>
              <Colon />
              {renderValue(v, depth + 1, entryPath, comments)}
              {i < entries.length - 1 ? <Comma /> : null}
              {comment ? <CommentText>{comment}</CommentText> : null}
              {"\n"}
            </span>
          )
        })}
        {pad}
        <Bracket>{"}"}</Bracket>
      </span>
    )
  }

  // Fallback for anything JSON can't represent (undefined, function, symbol)
  return <span key={path}>{String(value)}</span>
}

export function JsonCodeBlock({ data, comments = {} }: JsonCodeBlockProps) {
  return (
    <code className="font-code font-[750] text-sm tracking-wide selection:bg-purple-500/20 whitespace-pre">
      {renderValue(data, 0, "", comments)}
    </code>
  )
}
