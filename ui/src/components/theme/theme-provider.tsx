import { createContext, useContext, useEffect, useState } from "react"

export type Theme = "dark" | "light" | "system"

type ThemeProviderState = {
  theme: Theme
  setTheme: (theme: Theme) => void
}

const STORAGE_KEY = "waggle-theme"

const ThemeProviderContext = createContext<ThemeProviderState>({
  theme: "system",
  setTheme: () => null,
})

function applyTheme(theme: Theme) {
  const root = window.document.documentElement
  root.classList.remove("light", "dark")
  const resolved =
    theme === "system"
      ? window.matchMedia("(prefers-color-scheme: dark)").matches
        ? "dark"
        : "light"
      : theme
  root.classList.add(resolved)
}

export function ThemeProvider({ children }: { children: React.ReactNode }) {
  const [theme, setThemeState] = useState<Theme>(
    () => (localStorage.getItem(STORAGE_KEY) as Theme) || "system"
  )

  useEffect(() => {
    applyTheme(theme)
    if (theme !== "system") return
    // Keep "system" in sync with OS changes while selected.
    const mq = window.matchMedia("(prefers-color-scheme: dark)")
    const onChange = () => applyTheme("system")
    mq.addEventListener("change", onChange)
    return () => mq.removeEventListener("change", onChange)
  }, [theme])

  const setTheme = (next: Theme) => {
    localStorage.setItem(STORAGE_KEY, next)
    setThemeState(next)
  }

  return (
    <ThemeProviderContext.Provider value={{ theme, setTheme }}>
      {children}
    </ThemeProviderContext.Provider>
  )
}

export function useTheme() {
  return useContext(ThemeProviderContext)
}
