import {
  Box,
  Button,
  Callout,
  Card,
  Checkbox,
  Flex,
  Heading,
  TextField,
} from "@radix-ui/themes";
import {
  InfoCircledIcon,
  LockClosedIcon,
  PersonIcon,
} from "@radix-ui/react-icons";
import { useState } from "react";
import { logIn, register } from "../auth.js";
import { useNavigate } from "react-router-dom";

function Auth() {
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");
  const [remember, setRemember] = useState(false);

  const [infoMsg, setInfoMsg] = useState("");
  const [infoColor, setInfoColor] = useState("");

  const navigate = useNavigate();

  function showError(msg) {
    setInfoColor("red");
    setInfoMsg(msg);
  }

  function showInfo(msg) {
    setInfoColor("");
    setInfoMsg(msg);
  }

  async function tryLogIn() {
    const res = await logIn(username, password, remember);
    if (!res) {
      showError("Failed to log in");
    }
    navigate("/");
  }

  async function tryRegister() {
    const res = await register(username, password);
    if (!res) {
      showError("Failed to register");
      return;
    }

    showInfo("Registered successfully; you may log in");
    setUsername("");
    setPassword("");
  }

  return (
    <Flex align="center" direction="column" gap="3" justify="center" mt="9">
      <Card size="4" style={{ width: 400 }}>
        <Box mb="6">
          <Heading>Authorize</Heading>
        </Box>
        <Flex direction="column" gap="3" mb="6">
          <Box>
            <Heading mb="1" size="2">
              Username
            </Heading>
            <TextField.Root
              tabIndex="1"
              value={username}
              autoComplete="username"
              onChange={(v) => setUsername(v.target.value)}
            >
              <TextField.Slot>
                <PersonIcon />
              </TextField.Slot>
            </TextField.Root>
          </Box>
          <Box>
            <Heading mb="1" size="2">
              Password
            </Heading>
            <TextField.Root
              tabIndex="2"
              type="password"
              value={password}
              onChange={(v) => setPassword(v.target.value)}
            >
              <TextField.Slot>
                <LockClosedIcon />
              </TextField.Slot>
            </TextField.Root>
          </Box>
        </Flex>
        <Flex justify="between" align="end">
          <Flex gap="2" align="center">
            <Checkbox onCheckedChange={setRemember} />
            Remember
          </Flex>
          <Flex gap="2" align="center">
            <Button
              onClick={tryLogIn}
              disabled={username === "" || password === ""}
            >
              Log In
            </Button>
            <Button
              variant="surface"
              onClick={tryRegister}
              disabled={username === "" || password === ""}
            >
              Register
            </Button>
          </Flex>
        </Flex>
      </Card>
      {infoMsg !== "" && (
        <Callout.Root color={infoColor} variant="surface">
          <Callout.Icon>
            <InfoCircledIcon />
          </Callout.Icon>
          <Callout.Text>{infoMsg}</Callout.Text>
        </Callout.Root>
      )}
    </Flex>
  );
}

export default Auth;
