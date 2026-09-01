import { IconCaretRightFilled } from "@tabler/icons-react"
import { createFileRoute, Link } from "@tanstack/react-router"

export const Route = createFileRoute("/")({
  component: Index,
})

function Index() {
  return (
    <div className="flex-col p-10 flex min-w-full items-center min-h-screen text-center justify-between">
      <div />
      <div className="flex flex-col gap-7.5 max-w-4xl">
        <div className="flex flex-col gap-2.5">
          <h1 className="text-7xl tracking-tight font-black font-title">
            Punch
          </h1>
          <p className="text-xl">Distributed Load Tester</p>
        </div>
        <div className="min-w-sm w-full flex items-center gap-5">
          <a
            href="https://github.com/yuriongit/punch"
            rel="noopener noreferrer"
            target="_blank"
            className="border hover:cursor-pointer border-black/50 p-1.5 px-5 rounded-md w-full text-black hover:bg-black/5 transition-all duration-300 hover:-rotate-z-2 hover:scale-105 active:scale-95"
          >
            Visit Repo
          </a>
          <Link
            to="/start/walkthrough/disclaimer"
            className="hover:cursor-pointer p-1.5 px-5 rounded-md w-full bg-purple-600 hover:bg-purple-500 text-white transition-all hover:rotate-z-2 hover:scale-105 flex items-center gap-2.5 justify-center uppercase active:scale-95 duration-300 font-bold"
          >
            <span className="pl-3">Start!</span>
            <IconCaretRightFilled />
          </Link>
        </div>
      </div>
      <p className="max-w-md text-sm italic text-stone-700">
        A project focused on learning Kubernetes, deepening my understanding of
        concurrency and concurrent systems, and improving my skills with Docker,
        GitHub Actions, and Go.
      </p>
    </div>
  )
}
