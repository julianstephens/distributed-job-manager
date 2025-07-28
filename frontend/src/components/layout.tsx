import { useAuth0 } from "@auth0/auth0-react";
import { Box, Flex, Heading, Link, Text } from "@chakra-ui/react";
import { useEffect, type JSX, type PropsWithChildren } from "react";
import { useNavigate } from "react-router";

export const Layout = ({
  actionButton,
  children,
  title,
}: PropsWithChildren & { title: string; actionButton?: JSX.Element }) => {
  const { isAuthenticated, isLoading, logout } = useAuth0();
  const goto = useNavigate();

  useEffect(() => {
    if (!isLoading && !isAuthenticated) {
      logout();
      goto("/");
    }
  }, [isAuthenticated, isLoading]);

  return (
    <Box m={10}>
      <Flex w="full" justify="center" align="center" mb="10">
        <Link
          href="/jobs"
          display="flex"
          flexDirection="column"
          justifyContent="center"
          alignItems="center"
        >
          <Heading size="5xl" mx="auto" color="purple.solid">
            DJM
          </Heading>
          <Text>Distributed Job Manager</Text>
        </Link>
      </Flex>
      <Flex w="full" justify="space-between">
        <Heading size="2xl" mb="5">
          {title}
        </Heading>
        {actionButton && actionButton}
      </Flex>
      {children}
    </Box>
  );
};
