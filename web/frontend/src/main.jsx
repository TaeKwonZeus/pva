import { StrictMode } from "react";
import { createRoot } from "react-dom/client";
import App from "./App.jsx";
import { createBrowserRouter, RouterProvider } from "react-router-dom";
import Index from "./pages/Index.jsx";
import "@radix-ui/themes/styles.css";

const router = createBrowserRouter([
  {
    path: "/",
    element: <Index />,
  },
]);

document.body.style.margin = "0";

createRoot(document.getElementById("root")).render(
  <StrictMode>
    <App>
      <RouterProvider router={router} />
    </App>
  </StrictMode>,
);
