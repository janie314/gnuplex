import { nextui } from "@nextui-org/react";

/** @type {import('tailwindcss').Config} */

const slate500 = "#64748B";

export default {
  content: [
    "./index.html",
    "./src/**/*.{js,ts,jsx,tsx}",
    "./node_modules/@nextui-org/theme/dist/**/*.{js,ts,jsx,tsx}",
  ],
  darkMode: "class",
  plugins: [
    nextui({
      themes: {
        dark: {
          colors: {
            primary: {
              DEFAULT: slate500,
              // foreground: "#000000",
            },
            focus: slate500,
          },
        },
        light: {
          colors: {
            primary: {
              DEFAULT: slate500,
              // foreground: "#000000",
            },
            focus: slate500,
          },
        },
      },
    }),
  ],
};
