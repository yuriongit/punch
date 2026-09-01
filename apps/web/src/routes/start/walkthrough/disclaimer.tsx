import { createFileRoute, Link } from "@tanstack/react-router"

export const Route = createFileRoute("/start/walkthrough/disclaimer")({
  component: RouteComponent,
})

function RouteComponent() {
  return (
    <main className="tracking-tight flex-col px-10 flex min-w-full items-center min-h-screen text-center justify-between py-10">
      <div />
      <div className="space-y-3 flex flex-col items-center justify-center">
        <h1 className="text-3xl font-bold">
          Disclaimer:
          <span className="text-purple-600"> Currently Under Development!</span>
        </h1>
        <p className="max-w-lg">
          Punch is not yet an entirely feature-rich project. Although, that
          doesn't mean that the core functionality and more isn't served :)
        </p>
      </div>
      <Link
        to="/start/walkthrough/example-config"
        className="hover:cursor-pointer p-1.5 px-5 rounded-md bg-purple-600 hover:bg-purple-500 text-white transition-all hover:rotate-z-3 hover:scale-105 gap-2.5 justify-center active:scale-95 w-fit duration-300"
      >
        Okay
      </Link>
    </main>
  )
}
