import Cookies from "js-cookie";

async function fetchWithAuth(url, options) {
  if (!options) options = {};
  const token = Cookies.get("token");
  if (token) {
    options.headers = {
      ...options.headers,
      Authorization: `Bearer ${token}`,
    };
  }

  return fetch(url, options);
}

async function isLoggedIn() {
  if (!Cookies.get("token")) return false;
  return (await fetchWithAuth("/api/ping")).ok;
}

async function logIn(username, password, persist) {
  const res = await fetch("/api/auth/login", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({
      username,
      password,
    }),
  });

  if (!res.ok) {
    return false;
  }

  const token = (await res.json()).token;
  if (!persist) {
    Cookies.set("token", token);
  } else {
    Cookies.set("token", token, { expires: 30 });
  }

  return true;
}

async function logOut() {
  Cookies.remove("token");
}

async function register(username, password) {
  const res = await fetch("/api/auth/register", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({
      username,
      password,
    }),
  });

  return res.ok;
}

export { fetchWithAuth, isLoggedIn, logIn, logOut, register };
