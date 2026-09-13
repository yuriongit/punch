import { IconChevronRight } from "@tabler/icons-react"
import { createFileRoute, Link } from "@tanstack/react-router"
import { AnimatePresence, motion } from "motion/react"
import { JsonCodeBlock } from "../../../components/JsonCodeBlock"
import { PAGE_SLIDE_FWD_VARIANTS } from "../../../constants/animations"

export const Route = createFileRoute("/start/walkthrough/example-config")({
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

  return (
    <AnimatePresence mode="wait">
      <motion.main
        variants={PAGE_SLIDE_FWD_VARIANTS}
        initial="initial"
        animate="animate"
        exit="exit"
        className="tracking-tight flex-col px-10 flex min-w-full items-center min-h-screen text-left justify-center py-10"
      >
        <div className="w-full max-w-4xl gap-8 flex justify-center items-center">
          <div className="flex flex-col w-full min-w-3xs -space-y-1 justify-between h-full">
            <header className="text-3xl flex items-center justify-between font-bold">
              Configuration
            </header>
            <p className="py-5">
              This is a{" "}
              <span className="font-semibold text-purple-600">punch.json</span>
              {"; "}
              it's the configuration file Punch uses to load test your
              application.
            </p>
            <Link
              to="/start/test/create-config"
              className="hover:cursor-pointer p-1.5 px-3.5 rounded-md bg-purple-600 hover:bg-purple-500 text-white transition-all hover:rotate-z-3 hover:scale-105 gap-1.5 justify-center active:scale-95 w-fit duration-300 flex items-center mt-2.5"
            >
              <span className="ml-2.5">Next</span>
              <IconChevronRight stroke={1.5} />
            </Link>
          </div>
          <div className="border-2 bg-transparent border-purple-600/50 pl-10 pr-12.5 p-7.5 rounded-md">
            <JsonCodeBlock
              comments={{
                "children.1.expected_status": "// defaults to 200",
                "children.1.knockout_attempt": "// defaults to false",
              }}
              data={exampleConfig}
            />
          </div>
        </div>
        <footer className="bottom-10 fixed italic">
          <p className="max-w-lg text-sm  text-stone-700 text-center">
            Currently, configuration options are limited to the example's
            settings. More options will be added, alongside a dedicated page for
            configuration options.
          </p>
        </footer>
      </motion.main>
    </AnimatePresence>
  )
}
