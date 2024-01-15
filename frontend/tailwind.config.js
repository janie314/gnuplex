import { nextui } from "@nextui-org/react";

/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./index.html",
    "./src/**/*.{js,ts,jsx,tsx}",
    "./node_modules/@nextui-org/theme/dist/**/*.{js,ts,jsx,tsx}",
  ],
  theme: {
    colors: {
      primary: "#93c5fd",
      secondary: "#e6f1fe",
    },
    extend: {
      colors: {
        "med-blue": "#cbd5e1", // hsl(213, 27%, 84%)
      },
    },
  },
  darkMode: "class",
  plugins: [nextui()],
};
