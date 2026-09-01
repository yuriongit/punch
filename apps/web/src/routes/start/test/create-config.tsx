import { IconArrowBackUp, IconLockCheck } from "@tabler/icons-react"
import { createFileRoute, Link } from "@tanstack/react-router"
import { useState } from "react"
import { JsonEditor } from "../../../components/JsonEditor"

export const Route = createFileRoute("/start/test/create-config")({
  component: RouteComponent,
})

function RouteComponent() {
  const exampleConfig = {
    protocol: "http",
    victim: "localhost:3000/api/v1",
    children: [
      {
        name: "/urls",
        method: "GET",
        duration_secs: 30,
        concurrency_rate: 75,
      },
      {
        name: "/urls/shorten",
        method: "POST",
        expected_status: 301,
        duration_secs: 60,
        concurrency_rate: 100,
        knockout_attempt: true,
      },
    ],
  }

  const jsonExampleConfig = JSON.stringify(exampleConfig, null, 2)
  const [jsonInput, setJsonInput] = useState(jsonExampleConfig)
  const [error, setError] = useState<string | null>(null)

  const handleFormat = () => {
    try {
      const parsed = JSON.parse(jsonInput)
      setJsonInput(JSON.stringify(parsed, null, 2))
      setError(null)
    } catch (err) {
      if (err instanceof Error) {
        setError(err.message)
      }
    }
  }

  return (
    <main className="tracking-tight flex-col px-10 flex min-w-full items-center min-h-screen text-left justify-center py-10">
      <div className="max-w-4xl w-full flex flex-col gap-5">
        <div className="flex w-full justify-between items-end">
          <div className="space-y-2">
            <h1 className="text-3xl w-fit flex items-center justify-between font-bold">
              Configure Punch
            </h1>
            <p className="italic font-code text-stone-600">punch.json</p>
          </div>
          <button
            type="button"
            onClick={handleFormat}
            className="border text-xs hover:cursor-pointer border-black/50 p-1.5 px-5 rounded-md flex items-center gap-1.5 w-fit text-black hover:bg-black/5 transition-all duration-300"
          >
            Format JSON
          </button>
        </div>
        <div className="border-2 bg-transparent border-purple-600/50 px-1.5 pt-7.5 rounded-md justify-center">
          <div className="flex bg-white flex-col rounded overflow-hidden">
            {/* Validation Feedback */}
            <JsonEditor
              size={500}
              jsonInput={jsonInput}
              setJsonInput={setJsonInput}
            />
            {error && (
              <div className="p-3 mx-5 mb-5 text-xs bg-red-500/10 text-red-400 border rounded-md border-red-500/20 font-code">
                Invalid JSON: {error}
              </div>
            )}
          </div>
        </div>
        <div className="w-full flex justify-between items-center">
          <Link
            to="/start/walkthrough/example-config"
            className="border hover:cursor-pointer border-black/50 p-1.5 px-5 rounded-md flex items-center gap-1.5 w-fit text-black hover:bg-black/5 transition-all duration-300 hover:-rotate-z-2 hover:scale-105 active:scale-95"
          >
            <span>
              <IconArrowBackUp stroke={1.5} />
            </span>
          </Link>
          <Link
            to="/"
            // to="/start/test/punching"
            className="hover:cursor-pointer p-1.5 px-5 rounded-md bg-purple-600 hover:bg-purple-500 text-white transition-all hover:rotate-z-3 hover:scale-105 gap-2 justify-center active:scale-95 w-fit duration-300 flex items-center uppercase font-bold pl-6"
          >
            Confirm <IconLockCheck stroke={2.75} size={21} />
          </Link>
        </div>
      </div>
    </main>
  )
}
