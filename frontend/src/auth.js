async function isLoggedIn() {
  return (await fetch("/api/ping")).ok;
}

async function logIn(username, password, remember) {
  const res = await fetch(`/api/auth/login?remember=${remember}`, {
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

async function logOut() {
  await fetch("/api/auth/revoke", {
    method: "POST",
  });
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

export { isLoggedIn, logIn, logOut, register };
