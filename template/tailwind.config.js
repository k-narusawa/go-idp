module.exports = {
  purge: ["./src/**/*.{js,jsx,ts,tsx}", "./public/index.html"],
  darkMode: false, // or 'media' or 'class'
  theme: {
    extend: {
      backgroundImage: {
        "gradient-radial": "radial-gradient(var(--tw-gradient-stops))",
        "gradient-conic":
          "conic-gradient(from 180deg at 50% 50%, var(--tw-gradient-stops))",
      },
      colors: {
        transparent: "transparent",
        current: "currentColor",
        red: {
          light: "#fca5a5",
          DEFAULT: "#f87171",
          dark: "#ef4444",
        },
        blue: {
          light: "#93c5fd",
          DEFAULT: "#4285F4",
          dark: "#2563eb",
        },
        green: {
          light: "#86efac",
          DEFAULT: "#4ade80",
          dark: "#22c55e",
        },
        orange: {
          light: "#fed7aa",
          DEFAULT: "#fdba74",
          dark: "#f97316",
        },
        gray: {
          light: "#f3f4f6",
          DEFAULT: "#d1d5db",
          dark: "#6b7280",
        },
        white: {
          DEFAULT: "#ffffff",
          dark: "#f9fafb",
        },
      },
    },
    fontFamily: {
      sans: ["ui-sans-serif", "system-ui", "Noto Sans JP", "sans-serif"],
      serif: ["ui-serif", "Georgia", "Noto Serif JP", "serif"],
      mono: ["ui-monospace", "SFMono-Regular", "Noto Sans JP", "sans-serif"],
    },
  },
  variants: {
    extend: {},
  },
  plugins: [],
};
