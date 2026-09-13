import type { Variants } from "motion/react"

export const PAGE_SLIDE_FWD_VARIANTS: Variants = {
  initial: {
    x: -150,
    opacity: 0,
  },
  animate: {
    x: 0,
    opacity: 1,
    transitionDuration: 0.55,
    animationDuration: 0.55,
    transition: {
      x: { type: "spring", stiffness: 185, damping: 24, mass: 0.8 },
    },
  },
  exit: {
    x: 150,
    opacity: 0,
  },
}

export const PAGE_SLIDE_UP_VARIANTS: Variants = {
  initial: {
    y: 125,
    opacity: 0,
  },
  animate: {
    y: 0,
    opacity: 1,
    transitionDuration: 0.35,
    animationDuration: 0.35,
    transition: {
      y: { type: "spring", stiffness: 185, damping: 24, mass: 0.8 },
    },
  },
  exit: {
    y: -125,
    opacity: 0,
  },
}

export const PAGE_SLIDE_DOWN_VARIANTS: Variants = {
  initial: {
    y: -125,
    opacity: 0,
  },
  animate: {
    y: 0,
    opacity: 1,
    transitionDuration: 0.35,
    animationDuration: 0.35,
    transition: {
      y: { type: "spring", stiffness: 185, damping: 24, mass: 0.8 },
    },
  },
  exit: {
    y: 125,
    opacity: 0,
  },
}
