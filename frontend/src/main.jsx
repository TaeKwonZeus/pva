import { StrictMode } from "react";
import { createRoot } from "react-dom/client";
import App from "./App.jsx";
import {
  createBrowserRouter,
  redirect,
  RouterProvider,
} from "react-router-dom";
import Index from "./pages/Index.jsx";
import "@radix-ui/themes/styles.css";
import { Theme } from "@radix-ui/themes";
import { isLoggedIn } from "./auth.js";
import Auth from "./pages/Auth.jsx";
import Passwords from "./pages/Passwords.jsx";

async function authLoader() {
  console.log(await isLoggedIn());
  return (await isLoggedIn()) ? null : redirect("/auth");
}

const router = createBrowserRouter([
  {
    path: "/",
    children: [
      {
        element: <App />,
        loader: authLoader,
        children: [
          { index: true, element: <Index /> },
          {
            path: "/passwords",
            element: <Passwords />,
          },
        ],
      },
      {
        path: "auth",
        element: <Auth />,
      },
    ],
  },
  // {
  //   path: "/login",
  //   element: <Auth />,
  // },
]);

document.body.style.margin = "0";

createRoot(document.getElementById("root")).render(
  <StrictMode>
    <Theme appearance="dark">
      <RouterProvider router={router} />
    </Theme>
  </StrictMode>,
);
