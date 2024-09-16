import React from "react";
import ReactDOM from "react-dom/client";
import { App } from "./App.tsx";
import { NextUIProvider } from "@nextui-org/react";
import "./index.css";

ReactDOM.createRoot(document.getElementById("root") as HTMLElement).render(
  <React.StrictMode>
    <NextUIProvider>
      <main className="light text-foreground bg-background">
        <App />
      </main>
    </NextUIProvider>
  </React.StrictMode>,
);
