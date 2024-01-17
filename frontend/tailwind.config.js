import { nextui } from "@nextui-org/react";

/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./index.html",
    "./src/**/*.{js,ts,jsx,tsx}",
    "./node_modules/@nextui-org/theme/dist/**/*.{js,ts,jsx,tsx}",
  ],
  theme: {
    extend: {
      colors: {
        steel: {
          100: "#f1f5f9", // slate 100
          500: "#64748b", // slate 500 https://tailwindcss.com/docs/customizing-colors
        },
      },
    },
  },
  darkMode: "class",
  plugins: [nextui()],
};
