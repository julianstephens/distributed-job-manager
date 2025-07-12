import type { ChildrenProps } from "@/types";
import { Box, Heading } from "@chakra-ui/react";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";

export const Layout = ({
  children,
  title,
}: ChildrenProps & { title: string }) => {
  const queryClient = new QueryClient();

  return (
    <QueryClientProvider client={queryClient}>
      <Box m={10}>
        <Heading size="2xl" mb="5">
          {title}
        </Heading>
        {children}
      </Box>
    </QueryClientProvider>
  );
};
