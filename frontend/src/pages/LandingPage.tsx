import { useAuth0 } from "@auth0/auth0-react";
import { Button, Flex } from "@chakra-ui/react";
import { useEffect } from "react";
import { useNavigate } from "react-router";

const LandingPage = () => {
  const { loginWithRedirect, isAuthenticated } = useAuth0();
  const goto = useNavigate();

  useEffect(() => {
    if (isAuthenticated) {
      goto("/jobs");
    }
  }, [isAuthenticated]);

  return (
    <Flex w="full" h="full" justify="center" align="center">
      <Button w="1/5" onClick={() => loginWithRedirect()}>
        Log In
      </Button>
    </Flex>
  );
};

export default LandingPage;
