import { IconBoom, IconBrandGithub } from "@tabler/icons-react"
import {
  createRootRoute,
  Link,
  Outlet,
  useLocation,
} from "@tanstack/react-router"
import { TanStackRouterDevtools } from "@tanstack/react-router-devtools"
import { motion } from "framer-motion"
import ReactLenis, { useLenis } from "lenis/react"
import { AnimatePresence } from "motion/react"
import { useEffect } from "react"

// Default root transitions
const defaultVariants = {
  initial: { opacity: 0 },
  animate: { opacity: 1 },
}

function RootLayout() {
  const location = useLocation()

  const lenis = useLenis()
  
    // Reset scroll position on route change
    useEffect(() => {
      lenis?.scrollTo(0, { immediate: true })
    }, [lenis])
  
    return (
      <>
        <ReactLenis
          root
          options={{
            duration: 1.35,
            easing: (t) => Math.min(1, 1.001 - 2 ** (-10 * t)),
            smoothWheel: true,
            touchMultiplier: 2,
          }}
        >
          <div className="flex tracking-tight w-full fixed top-10 px-10 min-w-full justify-between z-50">
            <Link to={"/"} className="hover:text-purple-600">
              <IconBoom stroke={1.25} />
            </Link>
            <a
              href="https://github.com/yuriongit/punch"
              rel="noopener noreferrer"
              target="_blank"
              className="flex items-center gap-2.5 hover:text-purple-600 transition-all underline underline-offset-4 text-stone-500"
            >
              github.com/yuriongit/punch
              <span>
                <IconBrandGithub stroke={1.25} />
              </span>
            </a>
          </div>
          <AnimatePresence mode="wait">
            <motion.div
              key={location.pathname}
              variants={defaultVariants}
              initial="initial"
              animate="animate"
              transition={{ duration: 0.65 }}
              className="tracking-tight"
            >
              <Outlet />
            </motion.div>
          </AnimatePresence>
          
        </ReactLenis>
      <hr />
      <TanStackRouterDevtools />
    </>
  )
}

export const Route = createRootRoute({ component: RootLayout })
